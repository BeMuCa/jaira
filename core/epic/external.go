package epic

import (
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strings"
)

// External is the key of the task in another system that an epic stands for —
// JIRA-1234, a YouTrack issue — and where to open it. It sits on the epic and
// never on a ticket: the ticket file is committed with the code, the epic's ref
// is not, so a key kept for billing never ends up in a code commit.
//
// The frontmatter holds them as a block list, one key per line:
//
//	external:
//	  - JIRA-1234 https://jira.example.com/browse/JIRA-1234
//	  - YT-77
//
// URL may be empty: a key without a link is still something to filter by.
type External struct {
	Key string `json:"key"`
	URL string `json:"url"`
}

// externalField is the frontmatter key of the list.
const externalField = "external"

// ErrBadExternal means a key given on the command line or from the page is
// not "<KEY> [<URL>]" with an http or https URL. Nothing was written.
var ErrBadExternal = errors.New("an external key is '<KEY> [<URL>]', the URL http or https")

// String is the key as its line in the file has it.
func (x External) String() string { return strings.TrimSpace(x.Key + " " + x.URL) }

// ParseExternal reads "<KEY> [<URL>]" as a person typed it. It is strict
// where the file reader is lenient: what jaira writes itself is a key of one
// word and a link a page can open, never a javascript: or file: one.
func ParseExternal(raw string) (External, error) {
	f := strings.Fields(raw)
	if len(f) == 0 || len(f) > 2 || strings.HasPrefix(f[0], "-") {
		return External{}, fmt.Errorf("%w: got %q", ErrBadExternal, raw)
	}
	x := External{Key: f[0]}
	if len(f) == 2 {
		if !LinkableURL(f[1]) {
			return External{}, fmt.Errorf("%w: got %q", ErrBadExternal, raw)
		}
		x.URL = f[1]
	}
	return x, nil
}

// LinkableURL reports whether a URL is one a page may make a link of: http or
// https with a host. A hand-edited file may carry anything; this is what
// decides whether it is shown as a link or as text.
func LinkableURL(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

// parseItem reads one line of the list as the file has it, a quoted value
// included. It never refuses: a line a person wrote stays theirs.
func parseItem(raw string) (External, bool) {
	raw = strings.Trim(strings.TrimSpace(raw), `"'`)
	key, rest, _ := strings.Cut(raw, " ")
	if key == "" {
		return External{}, false
	}
	return External{Key: key, URL: strings.TrimSpace(rest)}, true
}

// sameKey compares keys the way a person means them: jira-12 is JIRA-12.
func sameKey(a, b string) bool { return strings.EqualFold(a, b) }

// Externals lists the epic's keys in file order.
func (e *Epic) Externals() []External { return slices.Clone(e.external) }

// HasExternal reports whether the epic carries a key, in any case.
func (e *Epic) HasExternal(key string) bool {
	return slices.ContainsFunc(e.external, func(x External) bool { return sameKey(x.Key, key) })
}

// putExternal adds a key, or gives a key the epic carries already the new URL
// in its place. It reports whether anything changed.
func (e *Epic) putExternal(x External) bool {
	for i, have := range e.external {
		if sameKey(have.Key, x.Key) {
			if have == x {
				return false
			}
			e.external = slices.Clone(e.external)
			e.external[i] = x
			return true
		}
	}
	e.external = append(slices.Clone(e.external), x)
	return true
}

// dropExternal takes a key out, in any case. It reports whether it was there.
func (e *Epic) dropExternal(key string) bool {
	if !e.HasExternal(key) {
		return false
	}
	e.external = slices.DeleteFunc(slices.Clone(e.external), func(x External) bool { return sameKey(x.Key, key) })
	return true
}

// isItem is a line of a block list: "- something" however indented.
func isItem(line string) bool {
	t := strings.TrimSpace(line)
	return t == "-" || strings.HasPrefix(t, "- ")
}

// writeExternal puts the list into the frontmatter as it stands, in place of
// the block that was there, and takes the block out when the list is empty:
// an epic that never had a key reads the same as one whose last key went.
func (e *Epic) writeExternal() {
	var block []string
	if len(e.external) > 0 {
		block = append(block, externalField+":")
		for _, x := range e.external {
			block = append(block, "  - "+x.String())
		}
	}
	end, hasFront := e.frontEnd()
	if !hasFront {
		if len(block) > 0 {
			e.lines = slices.Concat([]string{"---"}, block, []string{"---", ""}, e.lines)
		}
		return
	}
	start, stop := end, end
	for i := 1; i < end; i++ {
		k, _, ok := strings.Cut(strings.TrimSpace(e.lines[i]), ":")
		if ok && !isItem(e.lines[i]) && strings.TrimSpace(k) == externalField {
			start, stop = i, i+1
			for stop < end && isItem(e.lines[stop]) {
				stop++
			}
			break
		}
	}
	e.lines = slices.Concat(e.lines[:start], block, e.lines[stop:])
}
