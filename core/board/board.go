// Package board holds the logic for making a freshly created board private and
// telling the agents that work in this repository about it. It depends only on
// the stdlib so both internal/cli and internal/tui can call it without a cycle.
package board

// Prepared reports what preparing a new board did.
type Prepared struct {
	Ignored   bool // .gitignore gained the entry on this run
	IgnoreErr error
	Notes     []string // agent files touched, as "NAME (action)"
	NoteErr   error
}

// Prepare makes a freshly created board private and tells the agents about it.
//
// lanes is this board's own pipeline, which the note renders so an agent reads
// the route it is actually on rather than a generic description of lanes. A nil
// slice writes exactly the note that was written before boards could describe
// themselves, which is what a caller with no board loaded yet should pass.
func Prepare(root string, lanes []LaneFact) Prepared {
	var p Prepared
	p.Ignored, p.IgnoreErr = AddIgnore(root)
	p.Notes, p.NoteErr = AnnounceInAgentFiles(root, lanes)
	return p
}
