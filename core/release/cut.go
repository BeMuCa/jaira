package release

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var versionRe = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)

// Cut releases every note file directly in <pkgDir>/notes by moving it into
// <pkgDir>/notes/<version>/, and returns the version it used. A leading "v"
// is dropped, so "v0.3.7" and "0.3.7" cut the same release. It refuses a
// version that is not X.Y.Z, a release that already has a folder, and a
// release with no notes: that one would be a pull request without a commit.
//
// The release workflow (.github/workflows/release-notes.yaml) runs it on a
// branch of its own; nothing here touches git.
func Cut(pkgDir, version string) (string, error) {
	version = strings.TrimPrefix(version, "v")
	if !versionRe.MatchString(version) {
		return "", fmt.Errorf("version %q is not X.Y.Z", version)
	}
	notes := filepath.Join(pkgDir, notesDir)
	dest := filepath.Join(notes, version)
	if _, err := os.Stat(dest); err == nil {
		return "", fmt.Errorf("%s already exists: %s was released before", dest, version)
	}
	files, err := filepath.Glob(filepath.Join(notes, "*.md"))
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", fmt.Errorf("no note files directly in %s: nothing to release", notes)
	}
	if err := os.Mkdir(dest, 0o755); err != nil {
		return "", err
	}
	for _, f := range files {
		if err := os.Rename(f, filepath.Join(dest, filepath.Base(f))); err != nil {
			return "", err
		}
	}
	return version, nil
}
