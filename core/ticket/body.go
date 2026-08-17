package ticket

import "strings"

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
func NewBody(title, dod string, options []BodyOption) string {
	var b strings.Builder
	b.WriteString("# " + title + "\n\n")
	b.WriteString("## Definition of Done\n\n")
	if strings.TrimSpace(dod) != "" {
		b.WriteString("- [ ] " + strings.TrimSpace(dod) + "\n")
	} else {
		b.WriteString("- [ ] <A checkable statement, readable by someone who was not here>\n")
	}
	// Options turn steps on and off for this one ticket. Unticked by default:
	// most tickets need neither a brainstorm nor a separate planning pass, and
	// a step every ticket must traverse stops being a decision and becomes
	// ceremony. They are listed anyway, because an option nobody knows exists
	// is one nobody ticks. A per-user default board can pre-tick one, which is
	// what turns "always brainstorm" into a setting rather than a habit.
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
