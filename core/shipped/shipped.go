// Package shipped knows every version of every file jaira installs: the role
// prompts, the cockpit mod, the jaira skill and the built-in lanes.
//
// A file on someone's disk that is byte-identical to one of those versions was
// written by jaira and never touched by a person, so a newer jaira may replace
// it. Anything else was edited. The list is a manifest, not a copy of the old
// contents: one line per version, "<sha256>  <path>  <release>", generated from
// git history by scripts/shipped-history.sh.
package shipped

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"strings"
	"sync"
)

//go:embed shipped.sha256
var manifest string

// releases maps each manifest entry's "<hex>  <path>" to its release column.
// Blank lines and lines starting with '#' are not entries.
var releases = sync.OnceValue(func() map[string]string {
	set := map[string]string{}
	for _, l := range strings.Split(manifest, "\n") {
		l = strings.TrimRight(l, "\r")
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		f := strings.Fields(l)
		if len(f) < 2 {
			continue
		}
		rel := ""
		if len(f) > 2 {
			rel = f[2]
		}
		set[f[0]+"  "+f[1]] = rel
	}
	return set
})

// Known reports whether b is, byte for byte, a version of the file at path
// that some jaira build shipped. path is the file's path in the jaira
// repository, slash-separated, as it is embedded from.
func Known(path string, b []byte) bool {
	_, ok := releases()[Line(path, b)]
	return ok
}

// Release names the first release that shipped b as the file at path: "vX.Y.Z",
// or "-" for a version no release holds yet. ok is false when b is not a shipped
// version of path at all.
func Release(path string, b []byte) (release string, ok bool) {
	release, ok = releases()[Line(path, b)]
	return release, ok
}

// Line is b's line in the manifest: lowercase hex sha256, two spaces, path.
func Line(path string, b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]) + "  " + path
}
