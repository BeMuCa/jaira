// Package selfupdate replaces the running jaira binary with a published
// release.
//
// core/release answers "what version am I, and what changed since an older
// one" — a question the binary can answer entirely from what is embedded in
// itself. This package answers a different question: "what version is
// published right now, and how do I become it" — which means touching the
// network and the filesystem's executable bit. Keeping the two apart keeps
// the blast radius of a bug in the self-replace path contained to one small
// package that release.go never has to import.
package selfupdate

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// repo is the GitHub repository releases are published under.
const repo = "BeMuCa/jaira"

// apiBase returns the GitHub API root to query for release metadata.
// JAIRA_RELEASE_API overrides it. This exists so the test suite can serve a
// whole release from an httptest.Server rather than requiring network access
// to run 'go test', and it incidentally lets someone point this at a fork.
// It is deliberately not a flag: which endpoint answers "what is published"
// is not something a user should be choosing per invocation, and a flag
// would let it be slipped into a copy-pasted command line.
func apiBase() string {
	if v := os.Getenv("JAIRA_RELEASE_API"); v != "" {
		return v
	}
	return "https://api.github.com/repos/" + repo
}

// downloadBase returns the root release archives and checksums.txt are
// downloaded from. Same env-only override and the same reasoning as
// apiBase.
func downloadBase() string {
	if v := os.Getenv("JAIRA_RELEASE_DOWNLOADS"); v != "" {
		return v
	}
	return "https://github.com/" + repo + "/releases/download"
}

// maxDownload bounds every remote response body this package reads. An
// unbounded io.ReadAll of a remote body is a memory bomb waiting on a
// misbehaving or tampered server; no jaira archive or JSON payload is
// remotely near this size.
const maxDownload = 128 << 20

// errNotFound is returned by fetch on an HTTP 404, so a caller (At, below)
// can tell "this release does not exist" apart from any other fetch failure.
var errNotFound = errors.New("not found")

// Client talks to the release host.
type Client struct {
	HTTP *http.Client
}

// New returns a Client with an explicit timeout. The zero-value http.Client
// has no timeout at all, so a hung TLS handshake would hang the command
// forever — and this command sometimes runs detached (see core/selfupdate's
// cache.go), where nobody would ever see it stuck.
func New() *Client {
	return &Client{HTTP: &http.Client{Timeout: 30 * time.Second}}
}

func (c *Client) client() *http.Client {
	if c.HTTP != nil {
		return c.HTTP
	}
	return New().HTTP
}

// Release identifies one published release.
//
// Tag carries the leading "v" — GitHub tags do, and DownloadURL's path needs
// it. Version does not — release.Current and the archive name never have
// one. Getting the two backwards is the obvious bug here, which is exactly
// why they are two named fields instead of one string reused two ways.
type Release struct {
	Tag     string
	Version string
}

// Latest fetches the most recently published release from the API.
func (c *Client) Latest(ctx context.Context) (Release, error) {
	body, err := c.fetch(ctx, apiBase()+"/releases/latest")
	if err != nil {
		return Release{}, fmt.Errorf("fetching latest release: %w", err)
	}
	// GitHub's release payload has dozens of fields; this reads exactly the
	// one field we depend on, so a field we never look at can change shape
	// without ever breaking this.
	var payload struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return Release{}, fmt.Errorf("decoding latest release: %w", err)
	}
	if payload.TagName == "" {
		return Release{}, fmt.Errorf(
			"no release found for %s — install from source instead: go install github.com/%s/cmd/jaira@latest",
			repo, repo)
	}
	return Release{Tag: payload.TagName, Version: strings.TrimPrefix(payload.TagName, "v")}, nil
}

// At resolves an explicit tag or version to a Release without calling the
// API at all. It normalizes the tag (adding a leading "v" if the caller gave
// a bare version), then confirms the release exists by fetching
// checksums.txt for it — mapping a 404 to a "no such release" error.
//
// This skips releases/tags/<tag> deliberately: that would be a second GitHub
// payload shape to model and to keep working, when the artifact this
// actually needs already proves the release exists.
func (c *Client) At(ctx context.Context, tag string) (Release, error) {
	if !strings.HasPrefix(tag, "v") {
		tag = "v" + tag
	}
	version := strings.TrimPrefix(tag, "v")
	if _, err := c.fetch(ctx, c.DownloadURL(tag, "checksums.txt")); err != nil {
		if errors.Is(err, errNotFound) {
			return Release{}, fmt.Errorf("no release %s found", tag)
		}
		return Release{}, fmt.Errorf("checking release %s: %w", tag, err)
	}
	return Release{Tag: tag, Version: version}, nil
}

// AssetName returns the archive name goreleaser publishes for a version and
// platform. This must stay in step with .goreleaser.yaml's default archive
// name_template — that file names the thing this function must keep
// predicting.
func AssetName(version, goos, goarch string) string {
	ext := "tar.gz"
	if goos == "windows" {
		ext = "zip"
	}
	return fmt.Sprintf("jaira_%s_%s_%s.%s", version, goos, goarch, ext)
}

// DownloadURL returns the address of a named asset in a tagged release.
func (c *Client) DownloadURL(tag, asset string) string {
	return downloadBase() + "/" + tag + "/" + asset
}

// Binary downloads, verifies and extracts the jaira executable for rel on
// the given platform. Nothing is written to disk anywhere in this function:
// a verification failure must not be able to leave a partial file that some
// other code path could pick up. The order below — checksums first, then the
// archive, then the hash comparison, then extraction — is the security
// property itself: no byte of the archive is ever interpreted before it is
// known to be the byte-for-byte published one.
func (c *Client) Binary(ctx context.Context, rel Release, goos, goarch string) ([]byte, error) {
	asset := AssetName(rel.Version, goos, goarch)

	checksums, err := c.fetch(ctx, c.DownloadURL(rel.Tag, "checksums.txt"))
	if err != nil {
		return nil, fmt.Errorf("fetching checksums.txt: %w", err)
	}
	archive, err := c.fetch(ctx, c.DownloadURL(rel.Tag, asset))
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", asset, err)
	}

	want, ok := digestFor(checksums, asset)
	if !ok {
		// A missing entry means the release itself is broken (or the asset
		// name and the goreleaser template have drifted apart); a mismatched
		// digest below means the download was corrupted or tampered with.
		// They are different failures, so they get different messages.
		return nil, fmt.Errorf("%s has no entry in checksums.txt — this release looks broken", asset)
	}
	got := sha256.Sum256(archive)
	if hex.EncodeToString(got[:]) != want {
		return nil, fmt.Errorf("checksum mismatch for %s — the download may be corrupted or tampered with", asset)
	}

	return extract(archive, goos)
}

// fetch does one GET, bounding the response body and mapping a 404 to
// errNotFound so At can distinguish "no such release" from any other
// failure.
func (c *Client) fetch(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, errNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %s", resp.Status)
	}
	return io.ReadAll(io.LimitReader(resp.Body, maxDownload))
}

// digestFor finds the sha256 digest checksums.txt records for asset. A
// missing entry and a mismatched digest are different failures to whoever
// reads the error, so callers get to tell them apart instead of seeing one
// generic "verification failed".
func digestFor(checksums []byte, asset string) (string, bool) {
	for _, line := range strings.Split(string(checksums), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		if fields[1] == asset {
			return fields[0], true
		}
	}
	return "", false
}

// extract pulls the jaira executable out of a downloaded archive, returning
// its bytes without ever writing the archive or the entry to disk.
func extract(archive []byte, goos string) ([]byte, error) {
	name := "jaira"
	if goos == "windows" {
		name = "jaira.exe"
	}
	if goos == "windows" {
		return extractZip(archive, name)
	}
	return extractTarGz(archive, name)
}

// extractTarGz walks a .tar.gz looking for the entry whose path.Base matches
// name, ignoring any directory prefix so a future change to goreleaser's
// wrap_in_directory setting does not break this. Matching on Base rather
// than joining the entry name to a filesystem path also means there is
// nothing here for a "../../" entry to traverse — the return value is bytes
// in memory, never a path.
func extractTarGz(archive []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(archive))
	if err != nil {
		return nil, fmt.Errorf("opening archive: %w", err)
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading archive: %w", err)
		}
		if path.Base(hdr.Name) != name {
			continue
		}
		b, err := io.ReadAll(io.LimitReader(tr, maxDownload))
		if err != nil {
			return nil, err
		}
		if len(b) == 0 {
			return nil, fmt.Errorf("%s in archive is empty", name)
		}
		return b, nil
	}
	return nil, fmt.Errorf("no %s entry found in archive", name)
}

// extractZip is extractTarGz's counterpart for the Windows archive format —
// same Base-name match, same in-memory-only return.
func extractZip(archive []byte, name string) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
	if err != nil {
		return nil, fmt.Errorf("opening archive: %w", err)
	}
	for _, f := range zr.File {
		if path.Base(f.Name) != name {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		b, err := io.ReadAll(io.LimitReader(rc, maxDownload))
		rc.Close()
		if err != nil {
			return nil, err
		}
		if len(b) == 0 {
			return nil, fmt.Errorf("%s in archive is empty", name)
		}
		return b, nil
	}
	return nil, fmt.Errorf("no %s entry found in archive", name)
}

// stage writes bin to a temp file beside target, marks it executable, and
// returns its path. The temp file goes in target's own directory rather
// than os.TempDir(): a rename across filesystems is not atomic, and on most
// machines /tmp is a different filesystem from, say, ~/.local/bin — so the
// cross-device copy fallback os.Rename would need could leave a
// half-written binary on a crash. Staying on one filesystem is what makes
// Replace's final step a single atomic rename instead of a copy.
func stage(target string, bin []byte) (string, error) {
	dir := filepath.Dir(target)
	f, err := os.CreateTemp(dir, ".jaira-upgrade-*")
	if err != nil {
		return "", err
	}
	tmp := f.Name()
	if _, err := f.Write(bin); err != nil {
		f.Close()
		os.Remove(tmp)
		return "", err
	}
	if err := f.Chmod(0o755); err != nil {
		f.Close()
		os.Remove(tmp)
		return "", err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return "", err
	}
	return tmp, nil
}
