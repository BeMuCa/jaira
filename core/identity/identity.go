// Package identity answers "who is acting" the same way for every caller —
// the CLI, the TUI, and the core gate — so a ticket is attributed and owned
// consistently no matter which one wrote to it.
package identity

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/people"
)

// Current determines who is acting: the jaira name the board's people list
// gives this person, found by git's user.name or user.email. A person the list
// does not know acts under git's configured name, so tickets are attributed
// the same way commits are. dir scopes the git lookup to a particular
// repository.
func Current(dir string) string {
	return currentIn(People(dir), raw(dir), gitConfig(dir, "user.email"))
}

// currentIn is Current from what was already read, so Aliases asks git for
// each of them once.
func currentIn(list *people.List, raw, email string) string {
	if p := list.Find(raw); p != nil {
		return p.Name
	}
	if p := list.Find(email); p != nil {
		return p.Name
	}
	return raw
}

// raw is who git and the environment say is acting, before the people list
// has been asked.
func raw(dir string) string {
	if v := strings.TrimSpace(os.Getenv("JAIRA_USER")); v != "" {
		return v
	}
	if name := gitConfig(dir, "user.name"); name != "" {
		return name
	}
	for _, k := range []string{"USER", "USERNAME", "LOGNAME"} {
		if v := strings.TrimSpace(os.Getenv(k)); v != "" {
			return v
		}
	}
	return "unknown"
}

func gitConfig(dir, key string) string {
	out, err := exec.Command("git", "-C", dir, "config", key).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// People is the board's people list as the local ref in dir carries it. No
// repository, no ref or an unreadable one is an empty list: who is acting must
// always have an answer, and without a list it is the one it always was.
func People(dir string) *people.List {
	out, err := exec.Command("git", "-C", dir, "show", gitref.PeopleRef+":"+gitref.PeopleFile).Output()
	if err != nil {
		return &people.List{}
	}
	return people.Parse(out)
}

// Initials renders name as the first letter of each word, lowercased — the
// form used to stamp a sync folder, so a teammate can tell whose sweep it was
// at a glance rather than reading a full slug. A name with only one word has
// no initials to take, so its whole Slug is used instead: a folder named
// "b-20260823" tells a teammate nothing, but one named "berk-20260823" does.
// Reuses Slug for that single-word case and for the empty-name fallback, so
// there is one definition of "no usable characters" rather than two.
func Initials(name string) string {
	fields := strings.Fields(name)
	if len(fields) < 2 {
		return Slug(name)
	}
	var b strings.Builder
	for _, f := range fields {
		for _, r := range strings.ToLower(f) {
			if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
				b.WriteRune(r)
				break
			}
		}
	}
	if b.Len() == 0 {
		return Slug(name)
	}
	return b.String()
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

// AliasesPath is the file listing the other names that mean you: one per line,
// blank lines and # comments ignored. It exists because jaira knew an identity
// as exactly one string, while a person is several — a git user.name on one
// machine, a work email in a ticket's assignee, a personal email in another
// repository. The ownership rail compared those strings and concluded that your
// own tickets were not yours, which trained the habit of passing --force to
// every move. $JAIRA_HOME relocates it, as it does the rest of the per-user
// state.
func AliasesPath() string {
	if v := os.Getenv("JAIRA_HOME"); v != "" {
		return filepath.Join(v, "identity")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".jaira", "identity")
}

// Aliases lists every string that means "me", starting with Current(dir) so the
// canonical name is always first. git's user.email is included without being
// configured, because a ticket assigned by a teammate's tooling routinely
// carries an address rather than a name, and nobody should have to discover that
// by being refused.
//
// The result is deduplicated case-insensitively and never empty.
//
// Every name the board's people list gives this person is included too, and
// git's own user.name when the list renamed it: a ticket written under any of
// them is this person's.
func Aliases(dir string) []string {
	list, rawName, email := People(dir), raw(dir), gitConfig(dir, "user.email")
	me := currentIn(list, rawName, email)
	out := []string{me}
	add := func(v string) {
		v = strings.TrimSpace(v)
		if v == "" {
			return
		}
		for _, have := range out {
			if strings.EqualFold(have, v) {
				return
			}
		}
		out = append(out, v)
	}
	for _, n := range list.Names(me)[1:] {
		add(n)
	}
	add(rawName)
	add(email)
	if p := AliasesPath(); p != "" {
		if raw, err := os.ReadFile(p); err == nil {
			for _, line := range strings.Split(string(raw), "\n") {
				if line = strings.TrimSpace(line); line != "" && !strings.HasPrefix(line, "#") {
					add(line)
				}
			}
		}
	}
	return out
}

// Same is the one comparison every caller uses for "are these the same person":
// case-insensitive, ignoring surrounding space. Two empty strings are not the
// same person — an unassigned ticket belongs to nobody, not to everybody.
func Same(a, b string) bool {
	a, b = strings.TrimSpace(a), strings.TrimSpace(b)
	if a == "" || b == "" {
		return false
	}
	return strings.EqualFold(a, b)
}

// IsMe reports whether who names the person acting in dir, under any of their
// aliases.
func IsMe(dir, who string) bool {
	for _, a := range Aliases(dir) {
		if Same(a, who) {
			return true
		}
	}
	return false
}
