package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// dropPart is one thing the removal screen can take off a board.
//
// Each is a directory rather than a vague category, because what a person is
// agreeing to delete should be nameable as a path — and because deleting some
// of a board and keeping the rest has to leave something coherent behind.
type dropPart struct {
	label string
	dir   string
	count int // files inside, so the screen says how much is at stake
	on    bool
}

// dropBoard is the "remove this board" screen.
//
// Deleting a board is the second irreversible thing this tool does, after
// deleting a ticket, so it is built like that one: a question you have to
// answer deliberately rather than dismiss. The cursor starts on No and has to
// be moved to Yes; a stray return key cancels.
type dropBoard struct {
	root    string
	name    string
	current bool
	parts   []dropPart
	idx     int  // 0..len(parts)-1 are the checkboxes, len(parts) is the answer row
	yes     bool // which half of the answer row the cursor is on
	msg     string
}

// newDropBoard describes what removing root would take, listing only what is
// actually there — an empty archive is not a decision anyone needs to make.
func newDropBoard(root, name string, current bool) *dropBoard {
	d := &dropBoard{root: root, name: name, current: current}
	jaira := filepath.Join(root, ticket.DirName)
	for _, p := range []struct{ label, dir string }{
		{"tickets", filepath.Join(jaira, ticket.TicketsSubdir)},
		{"archive", filepath.Join(jaira, ticket.ArchiveSubdir)},
		{"logbook", filepath.Join(jaira, ticket.LogbookSubdir)},
		{"lanes published to teammates", filepath.Join(jaira, ticket.SharedSubdir)},
		{"this project's lane files", lane.ProjectLanesDir(root)},
	} {
		if n, ok := countFiles(p.dir); ok {
			d.parts = append(d.parts, dropPart{label: p.label, dir: p.dir, count: n, on: true})
		}
	}
	// Whatever else .jaira holds — its .gitignore, its .gitattributes, the
	// directory itself. Listed last and last to go, since removing it is what
	// actually takes the board off the home screen.
	d.parts = append(d.parts, dropPart{label: "the rest of .jaira (the board itself)", dir: jaira, on: true})
	d.idx = len(d.parts) // start on the answer row, on No
	return d
}

func countFiles(dir string) (int, bool) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, false
	}
	n := 0
	for _, e := range entries {
		if e.IsDir() {
			sub, _ := countFiles(filepath.Join(dir, e.Name()))
			n += sub
			continue
		}
		n++
	}
	return n, true
}

// key handles one press. It reports whether the screen is finished, and if so
// whether anything was removed.
func (d *dropBoard) key(s string) (done, removed bool) {
	switch s {
	case "esc", "q":
		return true, false
	case "j", "down":
		if d.idx < len(d.parts) {
			d.idx++
		}
	case "k", "up":
		if d.idx > 0 {
			d.idx--
		}
	case "h", "left":
		d.yes = false
	case "l", "right":
		if !d.current {
			d.yes = true
		}
	case " ", "space", "x":
		if d.idx < len(d.parts) {
			d.parts[d.idx].on = !d.parts[d.idx].on
		}
	case "enter":
		if d.idx < len(d.parts) {
			d.parts[d.idx].on = !d.parts[d.idx].on
			return false, false
		}
		if !d.yes {
			return true, false
		}
		if err := d.remove(); err != nil {
			d.msg = err.Error()
			return false, false
		}
		return true, true
	}
	return false, false
}

// remove deletes what is ticked. The board directory goes last: if it went
// first there would be nothing left to delete the rest from.
func (d *dropBoard) remove() error {
	jaira := filepath.Join(d.root, ticket.DirName)
	var last string
	for _, p := range d.parts {
		if !p.on {
			continue
		}
		if p.dir == jaira {
			last = p.dir
			continue
		}
		if err := os.RemoveAll(p.dir); err != nil {
			return err
		}
	}
	if last != "" {
		return os.RemoveAll(last)
	}
	return nil
}

func (d *dropBoard) render(width, height int) string {
	var b strings.Builder
	b.WriteString(styLaneTitle.Render("Remove board") + "\n")
	b.WriteString(styBar.Render(strings.Repeat("─", min(width, 78))) + "\n\n")
	fmt.Fprintf(&b, "  %s\n  %s\n\n", stySelected.Render(d.name), styMeta.Render(truncate(d.root, max(10, width-4))))

	if d.current {
		b.WriteString(styWarn.Render("  This is the board you are looking at. Switch to another one first.") + "\n\n")
	}

	for i, p := range d.parts {
		mark := "[ ]"
		sty := styMeta
		if p.on {
			mark, sty = "[x]", styErr
		}
		cursor := "  "
		if i == d.idx {
			cursor = stySelected.Render("▌ ")
		}
		files := styMeta.Render(fmt.Sprintf(" (%d file(s))", p.count))
		if p.dir == filepath.Join(d.root, ticket.DirName) {
			files = ""
		}
		fmt.Fprintf(&b, "%s%s %s%s\n", cursor, sty.Render(mark), sty.Render(p.label), files)
	}

	b.WriteString("\n")
	no, yes := " No ", " Yes "
	if d.idx == len(d.parts) {
		if d.yes {
			yes = stySelected.Render(" Yes ")
			no = styMeta.Render(" No ")
		} else {
			no = stySelected.Render(" No ")
			yes = styMeta.Render(" Yes ")
		}
	} else {
		no, yes = styMeta.Render(no), styMeta.Render(yes)
	}
	fmt.Fprintf(&b, "  %s  %s   %s\n", no, yes, styWarn.Render("this cannot be undone"))

	b.WriteString("\n" + styMeta.Render(
		"Your catalogue lanes in ~/.jaira/lanes are never touched by this.") + "\n")
	if d.msg != "" {
		b.WriteString("\n" + styErr.Render(d.msg) + "\n")
	}
	b.WriteString("\n" + styMeta.Render("j k move · space toggle · h l no/yes · enter choose · esc cancel"))
	return clampBlock(b.String(), width, height)
}
