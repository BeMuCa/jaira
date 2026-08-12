package ticket

import (
	"fmt"
	"strings"
)

// Section names one of the two checklists a ticket can carry.
type Section int

const (
	// SectionDoD is the definition of done: the criteria that decide whether
	// the work is acceptable. Only this one gates the terminal lane.
	SectionDoD Section = iota
	// SectionPlan is the method being followed to get there.
	SectionPlan
)

func (s Section) String() string {
	if s == SectionPlan {
		return "plan"
	}
	return "definition of done"
}

func (s Section) headings() []string {
	if s == SectionPlan {
		return planHeadings
	}
	return dodHeadings
}

// SetItemState rewrites the marker of the i-th checkbox (0-based) of one section
// and returns the new body.
//
// Only the single marker character is replaced. The item's text, the surrounding
// prose, and every line ending are left exactly as they were, so a tick produces
// a one-character diff rather than a reflowed body — which is what keeps the file
// readable in a diff and cheap to merge.
//
// Marking an item as in progress clears any other in-progress item in the same
// section: the marker answers "which one is being worked on", and that question
// has one answer.
func SetItemState(body string, sec Section, i int, st State) (string, error) {
	lines := strings.Split(body, "\n")
	idx := checklistLineIndexes(lines, sec)
	if len(idx) == 0 {
		return "", fmt.Errorf("this ticket has no %s checklist", sec)
	}
	if i < 0 || i >= len(idx) {
		return "", fmt.Errorf("%s has %d item(s); there is no item %d", sec, len(idx), i+1)
	}

	if st == StateDoing {
		for n, ln := range idx {
			if n == i {
				continue
			}
			if item, ok := parseCheckbox(strings.TrimSpace(lines[ln])); ok && item.State == StateDoing {
				lines[ln] = replaceMarker(lines[ln], StateTodo)
			}
		}
	}
	lines[idx[i]] = replaceMarker(lines[idx[i]], st)
	return strings.Join(lines, "\n"), nil
}

// checklistLineIndexes returns the line numbers of the checkboxes belonging to a
// section, in order. Indexes are per-section, so a definition-of-done index can
// never address a checkbox that lives under some other heading.
func checklistLineIndexes(lines []string, sec Section) []int {
	headings := sec.headings()
	inSection := false
	var out []int
	for n, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			heading := strings.ToLower(strings.TrimLeft(trimmed, "# "))
			inSection = false
			for _, h := range headings {
				if strings.Contains(heading, h) {
					inSection = true
					break
				}
			}
			continue
		}
		if !inSection {
			continue
		}
		if _, ok := parseCheckbox(trimmed); ok {
			out = append(out, n)
		}
	}
	return out
}

// replaceMarker swaps the character between the brackets, leaving the rest of the
// line — indentation, bullet style, text, and any trailing carriage return —
// exactly as it was.
func replaceMarker(line string, st State) string {
	open := strings.Index(line, "[")
	if open < 0 || open+2 >= len(line) || line[open+2] != ']' {
		return line
	}
	return line[:open+1] + st.Marker() + line[open+2:]
}
