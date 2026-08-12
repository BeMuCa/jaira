// Package board holds the logic for making a freshly created board private and
// telling the agents that work in this repository about it. It depends only on
// the stdlib so both internal/cli and internal/tui can call it without a cycle.
package board

// Prepared reports what preparing a new board did.
type Prepared struct {
	Ignored   bool     // .gitignore gained the entry on this run
	IgnoreErr error
	Notes     []string // agent files touched, as "NAME (action)"
	NoteErr   error
}

// Prepare makes a freshly created board private and tells the agents about it.
func Prepare(root string) Prepared {
	var p Prepared
	p.Ignored, p.IgnoreErr = AddIgnore(root)
	p.Notes, p.NoteErr = AnnounceInAgentFiles(root)
	return p
}
