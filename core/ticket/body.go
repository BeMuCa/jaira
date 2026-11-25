package ticket

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// BodyOption is one entry in a new ticket's Options checklist.
type BodyOption struct {
	Name   string
	Ticked bool
}

// NewBody is the starting shape of a ticket's markdown, shared by the CLI and
// the TUI so a ticket created either way looks the same.
//
// options is resolved by the caller (core/ticket must not import core/lane,
// which itself imports core/ticket) from the installed lanes'
// requires-option fields, with a per-user default board deciding which ones
// start ticked. Passing nil options reproduces today's two-option shape with
// nothing ticked — the regression baseline this function must not silently
// drift from.
func NewBody(title string, dod []string, options []BodyOption) string {
	var b strings.Builder
	b.WriteString("# " + title + "\n\n")
	b.WriteString("## Definition of Done\n\n")
	// Several criteria are several boxes. One box per criterion is the whole
	// point: the terminal lane's gate counts boxes, so criteria folded into one
	// line are criteria the gate cannot check separately, and a ticket created
	// with six of them used to arrive carrying one.
	wrote := false
	for _, item := range dod {
		if strings.TrimSpace(item) == "" {
			continue
		}
		b.WriteString("- [ ] " + strings.TrimSpace(item) + "\n")
		wrote = true
	}
	if !wrote {
		b.WriteString("- [ ] <A checkable statement, readable by someone who was not here>\n")
	}
	// Options turn steps on and off for this one ticket. NewBody writes the
	// ticks it is given; the caller decides them. Without a default board that
	// names options:, a new ticket gets lane.DefaultOptions (planning ticked),
	// so every ticket goes through the pre-process planner; a per-user default
	// board can untick it. They are listed anyway, because an option nobody
	// knows exists is one nobody ticks.
	b.WriteString("\n## Options\n\n")
	for _, o := range options {
		mark := " "
		if o.Ticked {
			mark = "x"
		}
		b.WriteString("- [" + mark + "] " + o.Name + "\n")
	}

	// The Plan is how the work will be done, as opposed to the criteria for
	// accepting it. It is seeded empty rather than omitted: a heading that is
	// already there gets filled in, and one that has to be remembered does not.
	// The heading is seeded but deliberately holds no checkbox. A placeholder
	// item would count as a plan, which would make the pre-process lane's
	// promise to produce one satisfied by every ticket the moment it is created.
	b.WriteString("\n## Plan\n\n")
	b.WriteString("<Steps, in order — filled in by the pre-process step, or by you.>\n")
	// "Progress", not "Notes": this is where 'jaira note' appends its
	// timestamped entries, and seeding the same heading the writer targets is
	// what keeps them landing here instead of in a second section at the end
	// of the file. A human wanting free-form notes can add their own heading.
	b.WriteString("\n## Progress\n\n")
	return b.String()
}

// NoteHeading is the body section progress notes live under.
//
// Notes go in the ticket body rather than a sidecar file on purpose: a separate
// memory file would be one more thing to lose, would not travel with the ticket
// through git, and could drift out of sync with it.
const NoteHeading = "## Progress"

// AppendNote adds one timestamped note under the Progress heading, creating the
// heading if the body has none.
//
// It lives here rather than in the command that types notes because it now has
// two callers — a person writing one, and a take-over recording that it
// happened — and two copies of "where does a note go" would eventually disagree
// about the answer.
func AppendNote(d *Doc, text string, now time.Time, who string) error {
	text = strings.TrimSpace(text)
	if text == "" {
		return errors.New("ticket: a note with no text records nothing")
	}
	entry := fmt.Sprintf("- **%s · %s** — %s", now.UTC().Format("2006-01-02 15:04"), who, text)
	body := d.Body()
	if idx := strings.Index(body, NoteHeading); idx >= 0 {
		// Append under the existing heading, after anything already there, and
		// before whatever section follows it.
		rest := body[idx:]
		next := strings.Index(rest[len(NoteHeading):], "\n## ")
		if next < 0 {
			body = strings.TrimRight(body, "\n") + "\n" + entry + "\n"
		} else {
			at := idx + len(NoteHeading) + next
			body = strings.TrimRight(body[:at], "\n") + "\n" + entry + "\n" + body[at:]
		}
	} else {
		body = strings.TrimRight(body, "\n") + "\n\n" + NoteHeading + "\n\n" + entry + "\n"
	}
	d.SetBody(body)
	return nil
}

// Note is one progress note as AppendNote wrote it.
type Note struct {
	At     time.Time
	Author string
	// Text is everything after the heading, continuation lines included,
	// verbatim: it is markdown a person wrote.
	Text string
	// Raw is the whole entry, heading included, as it stands in the body.
	Raw string
}

// noteStart is where one note starts: its dated, bold heading. Every line up
// to the next one belongs to it, whatever it starts with — a continuation
// line reading "#79 ..." or "- a list" is still the same note.
var noteStart = regexp.MustCompile(`(?m)^- \*\*\d{4}-\d\d-\d\d \d\d:\d\d `)

// noteHead is the heading AppendNote writes: "- **<date> <time> · <who>** — ".
var noteHead = regexp.MustCompile(`^- \*\*(\d{4}-\d\d-\d\d \d\d:\d\d) · (.*?)\*\* — `)

// Notes reads the progress notes back out of a body, oldest first, the order
// AppendNote adds them in. Text before the first note is not a note, and the
// section ends at the next heading. A dated entry written by hand in another
// shape still counts as a note; it just has no author, and its Text is Raw.
func Notes(body string) []Note {
	idx := strings.Index(body, NoteHeading)
	if idx < 0 {
		return nil
	}
	section := body[idx+len(NoteHeading):]
	if next := strings.Index(section, "\n## "); next >= 0 {
		section = section[:next]
	}
	starts := noteStart.FindAllStringIndex(section, -1)
	out := make([]Note, 0, len(starts))
	for i, at := range starts {
		end := len(section)
		if i+1 < len(starts) {
			end = starts[i+1][0]
		}
		n := Note{Raw: strings.TrimSpace(section[at[0]:end])}
		n.Text = n.Raw
		if m := noteHead.FindStringSubmatchIndex(n.Raw); m != nil {
			n.At, _ = time.Parse("2006-01-02 15:04", n.Raw[m[2]:m[3]])
			n.Author = n.Raw[m[4]:m[5]]
			n.Text = strings.TrimSpace(n.Raw[m[1]:])
		}
		out = append(out, n)
	}
	return out
}
