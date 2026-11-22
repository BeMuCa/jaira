package ticket

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/oklog/ulid/v2"
)

// NewID mints a ticket identifier.
//
// ULIDs are used rather than sequential numbers (JAIRA-1, JAIRA-2, ...) because
// a sequential scheme needs a coordinator to hand out the next value. In a store
// with no server, that coordinator would be a counter file — the single line in
// the repository that every ticket creation touches, and therefore the worst
// possible merge-conflict hotspot. A ULID needs no coordinator, and it sorts
// lexicographically by creation time, so the tickets directory lists in order.
func NewID(now time.Time) string {
	return ulid.MustNew(ulid.Timestamp(now), rand.Reader).String()
}

// ValidID reports whether id is a well-formed ULID.
//
// Every reference between tickets — dependencies, follow-ups, the handle shown
// on a card — resolves through the id, so an id that does not parse makes the
// ticket unreachable rather than merely untidy.
func ValidID(id string) bool {
	_, err := ulid.ParseStrict(id)
	return err == nil
}

// slugMax is the most bytes a slug holds.
const slugMax = 48

// Slug renders a title into a filename-safe breadcrumb. It is a static hint for
// humans running ls or grep, deliberately not kept in sync with later title
// edits: renaming files on every title change would churn git history for no
// benefit. The frontmatter title is the only field callers should trust.
func Slug(title string) string {
	var b strings.Builder
	lastDash := true // suppress a leading dash
	for _, r := range strings.ToLower(title) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			lastDash = false
		default:
			if !lastDash {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	s := strings.Trim(b.String(), "-")
	if len(s) > slugMax {
		// Cut at a rune boundary: a cut through a multi-byte letter (ü, ß,
		// Cyrillic) leaves a name that is not valid UTF-8, and APFS refuses to
		// check such a file out at all.
		cut := slugMax
		for cut > 0 && !utf8.RuneStart(s[cut]) {
			cut--
		}
		s = strings.Trim(s[:cut], "-")
	}
	if s == "" {
		s = "untitled"
	}
	return s
}

// Filename is the on-disk name for a ticket.
func Filename(id, title string) string {
	return fmt.Sprintf("%s-%s.md", id, Slug(title))
}

// IDFromFilename recovers the identifier from a ticket filename.
func IDFromFilename(name string) string {
	name = strings.TrimSuffix(name, ".md")
	if i := strings.IndexByte(name, '-'); i > 0 {
		return name[:i]
	}
	return name
}

// NormalizeIDPrefix upper-cases a user-supplied prefix, since ULIDs are encoded
// in uppercase Crockford base32 but nobody types them that way.
func NormalizeIDPrefix(p string) string { return strings.ToUpper(strings.TrimSpace(p)) }

// Handle is the short reference shown to people and printed by the CLI.
//
// It is taken from the end of the ULID rather than the beginning, because the
// leading ten characters encode only the millisecond timestamp: tickets created
// in the same burst share them entirely. The trailing characters are the random
// component, so a tail is distinctive where a head is not.
//
// A pull request in blocked-by (pr#12) is its own short reference and comes
// back whole: its tail would be a bare number.
func Handle(id string) string {
	if _, isPR, _ := ParsePRRef(id); isPR || len(id) <= handleLen {
		return id
	}
	return id[len(id)-handleLen:]
}

const handleLen = 6
