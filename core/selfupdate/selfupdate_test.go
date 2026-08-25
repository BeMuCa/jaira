package selfupdate

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// buildTarGz builds a real gzip'd tar archive containing one entry "jaira"
// with the given content — real bytes, not a checked-in blob, so extraction
// is exercised against the same shape goreleaser actually produces.
func buildTarGz(t *testing.T, content []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	if err := tw.WriteHeader(&tar.Header{Name: "jaira", Mode: 0o755, Size: int64(len(content))}); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// buildZip builds a real zip archive containing one entry "jaira.exe",
// mirroring the Windows archive shape.
func buildZip(t *testing.T, content []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	fw, err := zw.Create("jaira.exe")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fw.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// sha256Hex computes the checksum the way checksums.txt records it, so the
// test never hardcodes a digest that could pass for the wrong reason if the
// fixture generation changes.
func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// releaseFixture stands up an httptest.Server serving one release: a
// releases/latest JSON body, one archive, and a checksums.txt computed from
// that archive's real bytes.
type releaseFixture struct {
	server  *httptest.Server
	tag     string
	asset   string
	archive []byte
}

func newReleaseFixture(t *testing.T, tag, asset string, archive []byte) *releaseFixture {
	t.Helper()
	rf := &releaseFixture{tag: tag, asset: asset, archive: archive}
	mux := http.NewServeMux()
	mux.HandleFunc("/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"tag_name":%q}`, tag)
	})
	mux.HandleFunc("/"+tag+"/"+asset, func(w http.ResponseWriter, r *http.Request) {
		w.Write(archive)
	})
	mux.HandleFunc("/"+tag+"/checksums.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s  %s\n", sha256Hex(archive), asset)
	})
	rf.server = httptest.NewServer(mux)
	t.Cleanup(rf.server.Close)
	t.Setenv("JAIRA_RELEASE_API", rf.server.URL)
	t.Setenv("JAIRA_RELEASE_DOWNLOADS", rf.server.URL)
	return rf
}

// noChecksumFixture is like newReleaseFixture but omits the asset's entry
// from checksums.txt entirely, for the "missing digest" failure mode.
func newNoChecksumFixture(t *testing.T, tag, asset string, archive []byte) *releaseFixture {
	t.Helper()
	rf := &releaseFixture{tag: tag, asset: asset, archive: archive}
	mux := http.NewServeMux()
	mux.HandleFunc("/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"tag_name":%q}`, tag)
	})
	mux.HandleFunc("/"+tag+"/"+asset, func(w http.ResponseWriter, r *http.Request) {
		w.Write(archive)
	})
	mux.HandleFunc("/"+tag+"/checksums.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "")
	})
	rf.server = httptest.NewServer(mux)
	t.Cleanup(rf.server.Close)
	t.Setenv("JAIRA_RELEASE_API", rf.server.URL)
	t.Setenv("JAIRA_RELEASE_DOWNLOADS", rf.server.URL)
	return rf
}

// newBadChecksumFixture is like newReleaseFixture but records a digest that
// does not match the served archive, for the "checksum mismatch" failure
// mode.
func newBadChecksumFixture(t *testing.T, tag, asset string, archive []byte) *releaseFixture {
	t.Helper()
	rf := &releaseFixture{tag: tag, asset: asset, archive: archive}
	mux := http.NewServeMux()
	mux.HandleFunc("/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"tag_name":%q}`, tag)
	})
	mux.HandleFunc("/"+tag+"/"+asset, func(w http.ResponseWriter, r *http.Request) {
		w.Write(archive)
	})
	mux.HandleFunc("/"+tag+"/checksums.txt", func(w http.ResponseWriter, r *http.Request) {
		wrongDigest := sha256Hex([]byte("not the archive"))
		fmt.Fprintf(w, "%s  %s\n", wrongDigest, asset)
	})
	rf.server = httptest.NewServer(mux)
	t.Cleanup(rf.server.Close)
	t.Setenv("JAIRA_RELEASE_API", rf.server.URL)
	t.Setenv("JAIRA_RELEASE_DOWNLOADS", rf.server.URL)
	return rf
}

func TestLatestReturnsTagAndVersion(t *testing.T) {
	content := []byte("binary-contents-for-latest-test")
	archive := buildTarGz(t, content)
	newReleaseFixture(t, "v1.3.0", AssetName("1.3.0", "linux", "amd64"), archive)

	c := New()
	rel, err := c.Latest(context.Background())
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if rel.Tag != "v1.3.0" {
		t.Errorf("Tag = %q, want %q", rel.Tag, "v1.3.0")
	}
	if rel.Version != "1.3.0" {
		t.Errorf("Version = %q, want %q", rel.Version, "1.3.0")
	}
}

func TestAssetNameFormatsPerPlatform(t *testing.T) {
	cases := []struct {
		version, goos, goarch, want string
	}{
		{"1.3.0", "linux", "amd64", "jaira_1.3.0_linux_amd64.tar.gz"},
		{"1.3.0", "windows", "arm64", "jaira_1.3.0_windows_arm64.zip"},
	}
	for _, tc := range cases {
		got := AssetName(tc.version, tc.goos, tc.goarch)
		if got != tc.want {
			t.Errorf("AssetName(%q, %q, %q) = %q, want %q", tc.version, tc.goos, tc.goarch, got, tc.want)
		}
	}
}

// TestFullFetchAndReplaceOverwritesTargetAtomically drives the whole
// download-verify-extract-replace path against a fake installed binary, and
// asserts the target ends up holding the archive's jaira entry byte for
// byte, executable, with no temp file left behind.
func TestFullFetchAndReplaceOverwritesTargetAtomically(t *testing.T) {
	content := []byte("the-new-binary-bytes")
	archive := buildTarGz(t, content)
	asset := AssetName("1.3.0", "linux", "amd64")
	rf := newReleaseFixture(t, "v1.3.0", asset, archive)

	dir := t.TempDir()
	target := filepath.Join(dir, "jaira")
	if err := os.WriteFile(target, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}

	c := New()
	rel, err := c.Latest(context.Background())
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	bin, err := c.Binary(context.Background(), rel, "linux", "amd64")
	if err != nil {
		t.Fatalf("Binary: %v", err)
	}
	if !bytes.Equal(bin, content) {
		t.Fatalf("Binary() = %q, want %q", bin, content)
	}
	if err := Replace(target, bin); err != nil {
		t.Fatalf("Replace: %v", err)
	}

	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, content) {
		t.Errorf("target contents = %q, want %q", got, content)
	}
	fi, err := os.Stat(target)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm()&0o111 == 0 {
		t.Errorf("target mode = %v, want it executable", fi.Mode())
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".jaira-upgrade-") {
			t.Errorf("leftover temp file %q in %s", e.Name(), dir)
		}
	}
	_ = rf
}

// TestBinaryChecksumMismatchLeavesTargetUnchanged asserts a wrong digest
// aborts the whole operation before anything is written, naming the archive.
func TestBinaryChecksumMismatchLeavesTargetUnchanged(t *testing.T) {
	content := []byte("mismatched-content")
	archive := buildTarGz(t, content)
	asset := AssetName("1.3.0", "linux", "amd64")
	newBadChecksumFixture(t, "v1.3.0", asset, archive)

	dir := t.TempDir()
	target := filepath.Join(dir, "jaira")
	if err := os.WriteFile(target, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}

	c := New()
	rel := Release{Tag: "v1.3.0", Version: "1.3.0"}
	_, err := c.Binary(context.Background(), rel, "linux", "amd64")
	if err == nil {
		t.Fatal("expected Binary to fail on a checksum mismatch")
	}
	if !strings.Contains(err.Error(), asset) {
		t.Errorf("error = %v, want it to name %q", err, asset)
	}

	got, readErr := os.ReadFile(target)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(got) != "old" {
		t.Errorf("target contents = %q, want %q (unchanged)", got, "old")
	}
}

// TestBinaryMissingChecksumEntryLeavesTargetUnchanged asserts an asset absent
// from checksums.txt aborts too, with a distinct message from the mismatch
// case, and the target unchanged.
func TestBinaryMissingChecksumEntryLeavesTargetUnchanged(t *testing.T) {
	content := []byte("no-checksum-entry-content")
	archive := buildTarGz(t, content)
	asset := AssetName("1.3.0", "linux", "amd64")
	newNoChecksumFixture(t, "v1.3.0", asset, archive)

	dir := t.TempDir()
	target := filepath.Join(dir, "jaira")
	if err := os.WriteFile(target, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}

	c := New()
	rel := Release{Tag: "v1.3.0", Version: "1.3.0"}
	_, err := c.Binary(context.Background(), rel, "linux", "amd64")
	if err == nil {
		t.Fatal("expected Binary to fail when checksums.txt has no entry for the asset")
	}
	if !strings.Contains(err.Error(), asset) {
		t.Errorf("error = %v, want it to name %q", err, asset)
	}
	if strings.Contains(err.Error(), "mismatch") {
		t.Errorf("error = %v, a missing entry must not read as a mismatch", err)
	}

	got, readErr := os.ReadFile(target)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(got) != "old" {
		t.Errorf("target contents = %q, want %q (unchanged)", got, "old")
	}
}

// TestExtractZipDrivesWindowsArchiveFormatOnAnyRunner exercises the zip
// extraction path directly, so the Windows archive format is covered even
// on a Linux CI runner.
func TestExtractZipDrivesWindowsArchiveFormatOnAnyRunner(t *testing.T) {
	content := []byte("windows-binary-contents")
	archive := buildZip(t, content)

	got, err := extract(archive, "windows")
	if err != nil {
		t.Fatalf("extract(windows): %v", err)
	}
	if !bytes.Equal(got, content) {
		t.Errorf("extract(windows) = %q, want %q", got, content)
	}
}

// TestExtractTarGzIgnoresDirectoryPrefix asserts extraction matches on the
// entry's base name, so a future wrap_in_directory change in goreleaser
// would not break this.
func TestExtractTarGzIgnoresDirectoryPrefix(t *testing.T) {
	content := []byte("nested-binary-contents")
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	if err := tw.WriteHeader(&tar.Header{Name: "jaira_1.3.0_linux_amd64/jaira", Mode: 0o755, Size: int64(len(content))}); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write(content); err != nil {
		t.Fatal(err)
	}
	tw.Close()
	gz.Close()

	got, err := extract(buf.Bytes(), "linux")
	if err != nil {
		t.Fatalf("extract: %v", err)
	}
	if !bytes.Equal(got, content) {
		t.Errorf("extract() = %q, want %q", got, content)
	}
}
