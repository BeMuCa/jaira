// Package shipped knows every version of every file jaira installs: the role
// prompts, the cockpit mod, the jaira skill and the built-in lanes.
//
// A file on someone's disk that is byte-identical to one of those versions was
// written by jaira and never touched by a person, so a newer jaira may replace
// it. Anything else was edited. The list is a manifest, not a copy of the old
// contents: one sha256 line per version, generated from git history by
// scripts/shipped-history.sh.
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

// lines is the manifest as a set of its exact "<hex>  <path>" lines. Blank
// lines and lines starting with '#' are not entries.
var lines = sync.OnceValue(func() map[string]bool {
	set := map[string]bool{}
	for _, l := range strings.Split(manifest, "\n") {
		l = strings.TrimRight(l, "\r")
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		set[l] = true
	}
	return set
})

// Known reports whether b is, byte for byte, a version of the file at path
// that some jaira build shipped. path is the file's path in the jaira
// repository, slash-separated, as it is embedded from.
func Known(path string, b []byte) bool {
	return lines()[Line(path, b)]
}

// Line is b's line in the manifest: lowercase hex sha256, two spaces, path.
func Line(path string, b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]) + "  " + path
}
