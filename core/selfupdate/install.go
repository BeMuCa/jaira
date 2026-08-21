package selfupdate

import (
	"os"
	"path/filepath"
	"strings"
)

// Kind classifies how the running binary got onto this machine, which
// decides whether 'jaira self upgrade' is allowed to touch it.
type Kind string

const (
	// SelfManaged means nothing else claims to own this file; upgrading it
	// in place is safe.
	SelfManaged Kind = "self-managed"
	// Managed means a package manager (Homebrew) owns this file.
	Managed Kind = "homebrew"
	// GoInstall means 'go install' put this file where it is.
	GoInstall Kind = "go-install"
	// Unwritable means neither of the above applies, but the directory
	// cannot be written to anyway.
	Unwritable Kind = "unwritable"
)

// Detect classifies target and, for anything other than SelfManaged,
// returns the instruction to print instead of upgrading.
//
// This exists because replacing a file a package manager owns leaves that
// manager's database lying about what is installed — its next operation
// (a `brew upgrade`, a `go install` of something else) will silently step on
// or contradict what this command just did. Refusing is the correct
// behavior here, not a limitation: `uv self update` draws exactly this same
// line for the same reason.
//
// The checks run in this order — Homebrew, then go-install, then
// writability — because ownership is a stronger fact than writability. A
// Homebrew prefix is usually writable by its owning user, so checking
// writability first would misclassify a Homebrew-owned file as merely
// self-managed-and-writable.
func Detect(target string) (Kind, string) {
	if kind, instr, ok := detectHomebrew(target); ok {
		return kind, instr
	}
	if kind, instr, ok := detectGoInstall(target); ok {
		return kind, instr
	}
	dir := filepath.Dir(target)
	if !writable(dir) {
		return Unwritable, "the directory " + dir + " is not writable — install to a directory you own, or re-run with permission to write there"
	}
	return SelfManaged, ""
}

// detectHomebrew reports whether target lives under a Homebrew-managed
// Cellar. A literal "Cellar" path element covers the common macOS/Linuxbrew
// layouts without needing either env var set; $HOMEBREW_CELLAR and
// $HOMEBREW_PREFIX are also honored when set, using filepath.Rel rather than
// strings.HasPrefix so a directory that merely starts with the same
// characters (e.g. "/opt/homebrew-mine" next to "/opt/homebrew") is not
// mistaken for a child of it.
func detectHomebrew(target string) (Kind, string, bool) {
	for _, elem := range splitPath(target) {
		if elem == "Cellar" {
			return Managed, "this binary is managed by Homebrew — run `brew upgrade jaira` instead", true
		}
	}
	for _, envVar := range []string{"HOMEBREW_CELLAR", "HOMEBREW_PREFIX"} {
		if prefix := os.Getenv(envVar); prefix != "" && isWithin(target, prefix) {
			return Managed, "this binary is managed by Homebrew — run `brew upgrade jaira` instead", true
		}
	}
	return "", "", false
}

// detectGoInstall reports whether target lives under the directory 'go
// install' places binaries in: $GOBIN if set, else $GOPATH/bin, else
// $HOME/go/bin.
func detectGoInstall(target string) (Kind, string, bool) {
	var candidates []string
	if v := os.Getenv("GOBIN"); v != "" {
		candidates = append(candidates, v)
	}
	if v := os.Getenv("GOPATH"); v != "" {
		candidates = append(candidates, filepath.Join(v, "bin"))
	}
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, "go", "bin"))
	}
	for _, dir := range candidates {
		if isWithin(target, dir) {
			return GoInstall, "this binary was installed with `go install` — run `go install github.com/BeMuCa/jaira/cmd/jaira@latest` instead", true
		}
	}
	return "", "", false
}

// isWithin reports whether target is prefix or a descendant of it, using
// filepath.Rel rather than strings.HasPrefix so that a directory sharing a
// string prefix with prefix but not actually nested inside it (e.g.
// "/opt/homebrew-mine" vs "/opt/homebrew") is not misclassified.
func isWithin(target, prefix string) bool {
	rel, err := filepath.Rel(prefix, target)
	if err != nil {
		return false
	}
	if filepath.IsAbs(rel) {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}

// splitPath splits an absolute path into its elements, for the Homebrew
// Cellar-element check. Converting to forward slashes first makes this
// correct on Windows too, where filepath.Separator is '\\'.
func splitPath(p string) []string {
	return strings.Split(filepath.ToSlash(filepath.Clean(p)), "/")
}

// writable reports whether dir can actually be written to, by creating and
// immediately removing a file in it. Permission bits are deliberately not
// consulted: they lie under root, and they lie on Windows (which maps them
// onto a read-only attribute with no real relationship to write access) — an
// actual write attempt is the only answer that is true on every platform
// this ships to.
func writable(dir string) bool {
	f, err := os.CreateTemp(dir, ".jaira-writable-check-*")
	if err != nil {
		return false
	}
	name := f.Name()
	f.Close()
	os.Remove(name)
	return true
}
