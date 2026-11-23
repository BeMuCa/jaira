// Package people holds the board's list of who is who: each person once, under
// the one name jaira writes for them, with every other name that means them —
// a git user.name on another machine, a user.email, a handle somebody else
// typed into an assignee.
//
// It exists because a person used to appear on a board under as many names as
// they had git configurations, and nothing could say that "BeMuCa" and
// "berkc" were one human. The list travels on its own ref, gitref.PeopleRef,
// so everybody who fetches reads the same answer; it is not a per-machine
// setting, because a name that differs per machine is the problem it solves.
//
// The file is the API, as everywhere in jaira: one person a line,
//
//	berkc: BeMuCa, berk@example.com
//
// the jaira name, a colon, then the other names, comma-separated. Lines
// starting with # and blank lines are ignored, so a person may edit the ref's
// file by hand; a write through jaira rewrites it sorted, without them.
//
// Tickets are never rewritten. An old ticket keeps the name it was written
// under, and every reader asks Canonical for the name to show — which is what
// makes a rename one write instead of a migration of every ticket file.
package people

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
)

// Person is one human: the name jaira writes for them and every other name
// that means them.
type Person struct {
	Name    string   `json:"name"`
	Aliases []string `json:"aliases"`
}

// List is the whole file. The zero value and a nil *List are an empty list,
// so a board nobody has named a person on behaves exactly as before.
type List struct {
	People []Person
}

var (
	// ErrBadName means a name cannot be written into the file: empty, or
	// carrying the colon or comma the line format splits on.
	ErrBadName = errors.New("not a usable name")
	// ErrTaken means a name already belongs to somebody else. Nothing was
	// written.
	ErrTaken = errors.New("name belongs to another person")
	// ErrUnknown means no person goes by that name.
	ErrUnknown = errors.New("no person of that name")
)

// Parse reads the file. A line it cannot read — no colon, or nothing before
// it — is skipped rather than refused: a list half-understood still names
// everyone it can, and one bad hand edit must not blind every command.
func Parse(content []byte) *List {
	l := &List{}
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		name, rest, ok := strings.Cut(line, ":")
		name = strings.TrimSpace(name)
		if !ok || name == "" {
			continue
		}
		p := Person{Name: name}
		for _, a := range strings.Split(rest, ",") {
			if a = strings.TrimSpace(a); a != "" && !same(a, name) && !containsFold(p.Aliases, a) {
				p.Aliases = append(p.Aliases, a)
			}
		}
		l.People = append(l.People, p)
	}
	return l
}

// header opens every file jaira writes, so a person who finds the ref knows
// how to edit it without reading this package.
const header = `# jaira people: one person a line - the jaira name, a colon, then every other
# name that means them (git user.name, git user.email), comma-separated.
# 'jaira people alias' and 'jaira people rename' write this file.
`

// Bytes renders the list sorted, so two writes of one list are one file and a
// diff of the ref shows only what changed.
func (l *List) Bytes() []byte {
	var b strings.Builder
	b.WriteString(header)
	for _, p := range l.Sorted() {
		b.WriteString(p.Name)
		b.WriteString(":")
		for i, a := range p.Aliases {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(" ")
			b.WriteString(a)
		}
		b.WriteString("\n")
	}
	return []byte(b.String())
}

// Sorted is every person by name, each one's aliases sorted too, ignoring
// case. The list itself is left in file order.
func (l *List) Sorted() []Person {
	if l == nil {
		return nil
	}
	out := make([]Person, len(l.People))
	for i, p := range l.People {
		p.Aliases = slices.Clone(p.Aliases)
		sort.SliceStable(p.Aliases, func(a, b int) bool { return lessFold(p.Aliases[a], p.Aliases[b]) })
		out[i] = p
	}
	sort.SliceStable(out, func(a, b int) bool { return lessFold(out[a].Name, out[b].Name) })
	return out
}

// Find is the person a name means — their jaira name or any alias, ignoring
// case and surrounding space — or nil. A hand-edited file that gives one
// alias to two people answers with the first of them in file order.
func (l *List) Find(name string) *Person {
	if l == nil || strings.TrimSpace(name) == "" {
		return nil
	}
	for i := range l.People {
		if same(l.People[i].Name, name) {
			return &l.People[i]
		}
	}
	for i := range l.People {
		if containsFold(l.People[i].Aliases, name) {
			return &l.People[i]
		}
	}
	return nil
}

// Canonical is the name jaira shows and writes for name: the jaira name of
// the person it means, or name itself, trimmed, when it means nobody listed.
func (l *List) Canonical(name string) string {
	if p := l.Find(name); p != nil {
		return p.Name
	}
	return strings.TrimSpace(name)
}

// Names is every name of the person name means, the jaira name first; for a
// name the list does not know, just that name.
func (l *List) Names(name string) []string {
	p := l.Find(name)
	if p == nil {
		return []string{strings.TrimSpace(name)}
	}
	return append([]string{p.Name}, p.Aliases...)
}

// Same reports whether a and b mean the same person: equal ignoring case, or
// two names of one listed person. Two empty names are nobody, not the same.
func (l *List) Same(a, b string) bool {
	if strings.TrimSpace(a) == "" || strings.TrimSpace(b) == "" {
		return false
	}
	return same(l.Canonical(a), l.Canonical(b))
}

// Alias gives person one more name. person may itself be any name of a listed
// person; a person nobody listed yet is added under that name. A name that
// already means somebody else is refused with ErrTaken: two people are made
// one by renaming, not by stealing one's name. Reports whether anything
// changed.
func (l *List) Alias(person, alias string) (bool, error) {
	person, alias = strings.TrimSpace(person), strings.TrimSpace(alias)
	if err := check(person); err != nil {
		return false, err
	}
	if err := check(alias); err != nil {
		return false, err
	}
	p := l.Find(person)
	if other := l.Find(alias); other != nil {
		if p != nil && other == p {
			return false, nil
		}
		return false, fmt.Errorf("%w: %q is a name of %s", ErrTaken, alias, other.Name)
	}
	if p == nil {
		if same(person, alias) {
			return false, nil
		}
		l.People = append(l.People, Person{Name: person})
		p = &l.People[len(l.People)-1]
	}
	p.Aliases = append(p.Aliases, alias)
	return true, nil
}

// Rename makes to the jaira name of the person old means; the name they had
// becomes one of their aliases, so every ticket written under it is still
// theirs. to may already be one of their own aliases, never a name of
// anybody else.
func (l *List) Rename(old, to string) error {
	old, to = strings.TrimSpace(old), strings.TrimSpace(to)
	if err := check(to); err != nil {
		return err
	}
	p := l.Find(old)
	if p == nil {
		return fmt.Errorf("%w: %q", ErrUnknown, old)
	}
	if other := l.Find(to); other != nil && other != p {
		return fmt.Errorf("%w: %q is a name of %s", ErrTaken, to, other.Name)
	}
	if p.Name == to {
		return nil
	}
	aliases := []string{}
	if !same(p.Name, to) {
		aliases = append(aliases, p.Name)
	}
	for _, a := range p.Aliases {
		if !same(a, to) && !containsFold(aliases, a) {
			aliases = append(aliases, a)
		}
	}
	p.Name, p.Aliases = to, aliases
	return nil
}

// check refuses a name the line format cannot carry.
func check(name string) error {
	switch {
	case name == "":
		return fmt.Errorf("%w: empty", ErrBadName)
	case strings.ContainsAny(name, ":,\n\r"):
		return fmt.Errorf("%w: %q carries a colon, a comma or a line break, which the people file splits on", ErrBadName, name)
	case strings.HasPrefix(name, "#"):
		return fmt.Errorf("%w: %q starts with #, which the people file reads as a comment", ErrBadName, name)
	}
	return nil
}

func same(a, b string) bool { return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b)) }

func containsFold(list []string, v string) bool {
	return slices.ContainsFunc(list, func(x string) bool { return same(x, v) })
}

func lessFold(a, b string) bool {
	la, lb := strings.ToLower(a), strings.ToLower(b)
	if la != lb {
		return la < lb
	}
	return a < b
}
