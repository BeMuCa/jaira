package cli

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/BeMuCa/jaira/core/selfupdate"
)

// runSelf drives the real 'self' command tree, the same path a user or agent
// invokes it through, and returns its combined stdout/stderr.
func runSelf(t *testing.T, dir string, args ...string) (string, error) {
	t.Helper()
	root := newRoot("test")
	var out strings.Builder
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs(append([]string{"-C", dir, "self"}, args...))
	err := root.Execute()
	return out.String(), err
}

// selfArchiveContent builds the archive bytes 'self upgrade' would extract
// on the current runtime platform, carrying a distinguishable payload so a
// test can assert on exactly which version's bytes landed.
func selfArchiveContent(t *testing.T, payload string) []byte {
	t.Helper()
	if runtime.GOOS == "windows" {
		var buf bytes.Buffer
		zw := zip.NewWriter(&buf)
		fw, err := zw.Create("jaira.exe")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := fw.Write([]byte(payload)); err != nil {
			t.Fatal(err)
		}
		if err := zw.Close(); err != nil {
			t.Fatal(err)
		}
		return buf.Bytes()
	}
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	if err := tw.WriteHeader(&tar.Header{Name: "jaira", Mode: 0o755, Size: int64(len(payload))}); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write([]byte(payload)); err != nil {
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

func selfSHA256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// selfTestRelease stands up an httptest.Server serving 'latest' at
// latestTag with latestContent, plus a second, older release at pinTag with
// pinContent for '--version' to pin to, and points
// JAIRA_RELEASE_API/JAIRA_RELEASE_DOWNLOADS at it. latestHits counts calls
// to releases/latest, so a test can assert '--version' never makes one.
func selfTestRelease(t *testing.T, latestTag, latestContent, pinTag, pinContent string) (latestHits *int32) {
	t.Helper()
	// self upgrade / self upgrade --check write the release-check cache
	// (core/selfupdate/cache.go) on success. Without this, every test here
	// would write straight into the real user's ~/.jaira — exactly the kind
	// of test-writes-to-the-real-home bug the rest of this project isolates
	// against everywhere else.
	t.Setenv("JAIRA_HOME", t.TempDir())
	latestAsset := selfupdate.AssetName(strings.TrimPrefix(latestTag, "v"), runtime.GOOS, runtime.GOARCH)
	pinAsset := selfupdate.AssetName(strings.TrimPrefix(pinTag, "v"), runtime.GOOS, runtime.GOARCH)
	latestArchive := selfArchiveContent(t, latestContent)
	pinArchive := selfArchiveContent(t, pinContent)

	var hits int32
	mux := http.NewServeMux()
	mux.HandleFunc("/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		fmt.Fprintf(w, `{"tag_name":%q}`, latestTag)
	})
	mux.HandleFunc("/"+latestTag+"/"+latestAsset, func(w http.ResponseWriter, r *http.Request) {
		w.Write(latestArchive)
	})
	mux.HandleFunc("/"+latestTag+"/checksums.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s  %s\n", selfSHA256Hex(latestArchive), latestAsset)
	})
	if pinTag != latestTag {
		mux.HandleFunc("/"+pinTag+"/"+pinAsset, func(w http.ResponseWriter, r *http.Request) {
			w.Write(pinArchive)
		})
		mux.HandleFunc("/"+pinTag+"/checksums.txt", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s  %s\n", selfSHA256Hex(pinArchive), pinAsset)
		})
	}

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	t.Setenv("JAIRA_RELEASE_API", server.URL)
	t.Setenv("JAIRA_RELEASE_DOWNLOADS", server.URL)
	return &hits
}

// selfUpgradeTarget resolves the same path 'self upgrade' would operate on,
// so a test can read its bytes before/after to assert nothing was written.
func selfUpgradeTarget(t *testing.T) string {
	t.Helper()
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	target, err := filepath.EvalSymlinks(exe)
	if err != nil {
		t.Fatal(err)
	}
	return target
}

func TestSelfUpgradeUpToDateWritesNoFile(t *testing.T) {
	setCurrent(t, "1.3.0")
	selfTestRelease(t, "v1.3.0", "latest-bytes", "v1.1.0", "pin-bytes")
	target := selfUpgradeTarget(t)
	before, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}

	out, err := runSelf(t, t.TempDir(), "upgrade", "--json")
	if err != nil {
		t.Fatalf("self upgrade --json: %v\n%s", err, out)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("stdout was not a single JSON object: %v\n%s", err, out)
	}
	if payload["up_to_date"] != true {
		t.Errorf("up_to_date = %v, want true", payload["up_to_date"])
	}
	if payload["upgraded"] != false {
		t.Errorf("upgraded = %v, want false", payload["upgraded"])
	}

	after, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Error("an already-current binary must not be rewritten")
	}
}

func TestSelfUpgradeDevBuildRefusesButCheckReports(t *testing.T) {
	setCurrent(t, "dev")
	selfTestRelease(t, "v1.3.0", "latest-bytes", "v1.1.0", "pin-bytes")
	dir := t.TempDir()

	out, err := runSelf(t, dir, "upgrade")
	if err == nil {
		t.Fatalf("expected a dev build to refuse; got %s", out)
	}
	var ce *codedError
	if !errors.As(err, &ce) {
		t.Fatalf("error is not a codedError: %v", err)
	}
	if ce.code != ExitValidation {
		t.Errorf("code = %d, want %d", ce.code, ExitValidation)
	}
	if ce.reason != "dev_build" {
		t.Errorf("reason = %q, want %q", ce.reason, "dev_build")
	}
	if out != "" {
		t.Errorf("refusal wrote %q, want nothing written before the error is returned", out)
	}

	checkOut, err := runSelf(t, dir, "upgrade", "--check")
	if err != nil {
		t.Fatalf("self upgrade --check on a dev build must not refuse: %v\n%s", err, checkOut)
	}
	if !strings.Contains(checkOut, "dev build") {
		t.Errorf("--check output = %q, want it to say this is a dev build", checkOut)
	}
}

func TestSelfUpgradeCheckAgainstNewerReleaseDoesNotInstall(t *testing.T) {
	setCurrent(t, "1.0.0")
	selfTestRelease(t, "v1.3.0", "latest-bytes", "v1.1.0", "pin-bytes")
	target := selfUpgradeTarget(t)
	before, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}

	out, err := runSelf(t, t.TempDir(), "upgrade", "--check", "--json")
	if err != nil {
		t.Fatalf("self upgrade --check --json: %v\n%s", err, out)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("stdout was not a single JSON object: %v\n%s", err, out)
	}
	if payload["up_to_date"] != false {
		t.Errorf("up_to_date = %v, want false", payload["up_to_date"])
	}
	if payload["upgraded"] != false {
		t.Errorf("upgraded = %v, want false", payload["upgraded"])
	}
	if payload["latest"] != "1.3.0" {
		t.Errorf("latest = %v, want %q", payload["latest"], "1.3.0")
	}

	after, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Error("--check must not write the target file")
	}
}

// TestSelfUpgradeVersionPinsInstallsExactRelease asserts '--version' installs
// the named release even though a newer one is published, and never
// requests releases/latest to do it.
//
// This drives the real command, which means it really does replace the
// currently-running test binary on disk — safe on unix because a running
// process keeps executing from its old, now-unlinked inode (the same
// property Replace's own doc comment relies on), and because 'go test'
// links a fresh temporary binary for every invocation, so nothing outside
// this one test run is affected.
func TestSelfUpgradeVersionPinsInstallsExactRelease(t *testing.T) {
	setCurrent(t, "1.0.0")
	hits := selfTestRelease(t, "v1.3.0", "latest-bytes", "v1.1.0", "pin-bytes")
	target := selfUpgradeTarget(t)

	out, err := runSelf(t, t.TempDir(), "upgrade", "--version", "v1.1.0", "--json")
	if err != nil {
		t.Fatalf("self upgrade --version v1.1.0: %v\n%s", err, out)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("stdout was not a single JSON object: %v\n%s", err, out)
	}
	if payload["latest"] != "1.1.0" {
		t.Errorf("latest = %v, want %q (the pinned version)", payload["latest"], "1.1.0")
	}
	if payload["upgraded"] != true {
		t.Errorf("upgraded = %v, want true", payload["upgraded"])
	}

	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "pin-bytes" {
		t.Errorf("target contents = %q, want %q (the pinned release, not latest)", got, "pin-bytes")
	}
	if atomic.LoadInt32(hits) != 0 {
		t.Errorf("releases/latest was requested %d time(s), want 0 when --version is given", *hits)
	}
}

func TestSelfUpgradeVersionPinUnknownReleaseExitsValidation(t *testing.T) {
	setCurrent(t, "1.0.0")
	selfTestRelease(t, "v1.3.0", "latest-bytes", "v1.1.0", "pin-bytes")

	_, err := runSelf(t, t.TempDir(), "upgrade", "--version", "v9.9.9")
	if err == nil {
		t.Fatal("expected a nonexistent pinned release to fail")
	}
	var ce *codedError
	if !errors.As(err, &ce) {
		t.Fatalf("error is not a codedError: %v", err)
	}
	if ce.code != ExitValidation {
		t.Errorf("code = %d, want %d", ce.code, ExitValidation)
	}
	if ce.reason != "no_release" {
		t.Errorf("reason = %q, want %q", ce.reason, "no_release")
	}
}

func TestSelfUpgradeCheckAndVersionTogetherIsUsageError(t *testing.T) {
	setCurrent(t, "1.0.0")
	selfTestRelease(t, "v1.3.0", "latest-bytes", "v1.1.0", "pin-bytes")

	_, err := runSelf(t, t.TempDir(), "upgrade", "--check", "--version", "v1.1.0")
	if err == nil {
		t.Fatal("expected --check and --version together to fail")
	}
	var ce *codedError
	if !errors.As(err, &ce) {
		t.Fatalf("error is not a codedError: %v", err)
	}
	if ce.code != ExitUsage {
		t.Errorf("code = %d, want %d", ce.code, ExitUsage)
	}
}

// TestSelfUpgradeRefusalWritesNothingBeforeReturningTheError asserts a
// refusal (here: an unknown pinned release) never writes partial output —
// the combined stdout+stderr buffer is empty, and the error alone carries
// the reason. The JSON-on-stderr rendering itself is report()'s job, shared
// by every command, and is not re-tested per command here.
func TestSelfUpgradeRefusalWritesNothingBeforeReturningTheError(t *testing.T) {
	setCurrent(t, "1.0.0")
	selfTestRelease(t, "v1.3.0", "latest-bytes", "v1.1.0", "pin-bytes")

	out, err := runSelf(t, t.TempDir(), "upgrade", "--version", "v9.9.9", "--json")
	if err == nil {
		t.Fatal("expected an error")
	}
	if out != "" {
		t.Errorf("output = %q, want nothing written on a refusal", out)
	}
}
