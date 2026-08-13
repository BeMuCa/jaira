// Package identity answers "who is acting" the same way for every caller —
// the CLI, the TUI, and the core gate — so a ticket is attributed and owned
// consistently no matter which one wrote to it.
package identity

import (
	"os"
	"os/exec"
	"strings"
)

// Current determines who is acting, preferring git's configured name so
// tickets are attributed the same way commits are. dir scopes the git lookup
// to a particular repository.
func Current(dir string) string {
	if v := strings.TrimSpace(os.Getenv("JAIRA_USER")); v != "" {
		return v
	}
	if out, err := exec.Command("git", "-C", dir, "config", "user.name").Output(); err == nil {
		if name := strings.TrimSpace(string(out)); name != "" {
			return name
		}
	}
	for _, k := range []string{"USER", "USERNAME", "LOGNAME"} {
		if v := strings.TrimSpace(os.Getenv(k)); v != "" {
			return v
		}
	}
	return "unknown"
}

// Slug reduces name to a string safe as a path component: lowercase,
// [a-z0-9-] only, with runs of anything else collapsed to a single dash and
// no leading or trailing dash. A name that reduces to nothing gets a
// non-empty fallback, since an empty path component would silently write to
// the parent directory.
func Slug(name string) string {
	var b strings.Builder
	dashed := false
	for _, r := range strings.ToLower(name) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
			dashed = false
			continue
		}
		if b.Len() > 0 && !dashed {
			b.WriteByte('-')
			dashed = true
		}
	}
	s := strings.TrimRight(b.String(), "-")
	if s == "" {
		return "unnamed"
	}
	return s
}
