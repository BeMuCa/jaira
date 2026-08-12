package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/berk/jaira/core/project"
	"github.com/berk/jaira/core/ticket"
)

// browser is the directory picker for adding a board without leaving jaira.
//
// It lists directories only. The thing being chosen is a project root, and
// showing every file in a source tree would bury the two or three directories
// that are actually candidates.
type browser struct {
	dir string
	// pending is a directory the user tried to add that has no board yet, held
	// so that pressing i creates one there rather than somewhere else.
	pending string
	entries []browseEntry
	idx     int
	msg     string
}

type browseEntry struct {
	name    string
	isBoard bool
}

func newBrowser(start string) *browser {
	b := &browser{}
	if start == "" {
		start, _ = os.Getwd()
	}
	b.goTo(start)
	return b
}

func (b *browser) goTo(dir string) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		b.msg = err.Error()
		return
	}
	entries, err := os.ReadDir(abs)
	if err != nil {
		b.msg = err.Error()
		return
	}
	b.dir, b.idx, b.msg, b.pending = abs, 0, "", ""
	b.entries = nil
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		p := filepath.Join(abs, e.Name())
		b.entries = append(b.entries, browseEntry{name: e.Name(), isBoard: isBoardDir(p)})
	}
	sort.Slice(b.entries, func(i, j int) bool {
		// Boards first: they are what the screen is for.
		if b.entries[i].isBoard != b.entries[j].isBoard {
			return b.entries[i].isBoard
		}
		return b.entries[i].name < b.entries[j].name
	})
}

func isBoardDir(dir string) bool {
	fi, err := os.Stat(filepath.Join(dir, ".jaira", "tickets"))
	return err == nil && fi.IsDir()
}

func (b *browser) selected() string {
	if b.idx < 0 || b.idx >= len(b.entries) {
		return ""
	}
	return filepath.Join(b.dir, b.entries[b.idx].name)
}

// key drives the browser. It reports the roots to register, if any, and whether
// the browser is finished.
func (b *browser) key(s string) (added []string, done bool) {
	switch s {
	case "esc", "q":
		return nil, true
	case "j", "down":
		if b.idx < len(b.entries)-1 {
			b.idx++
		}
	case "k", "up":
		if b.idx > 0 {
			b.idx--
		}
	case "h", "left", "backspace":
		b.goTo(filepath.Dir(b.dir))
	case "l", "right", "enter":
		if sel := b.selected(); sel != "" {
			b.goTo(sel)
		}
	case "a":
		// Add whatever is highlighted, or this directory if nothing is.
		target := b.selected()
		if target == "" {
			target = b.dir
		}
		if !isBoardDir(target) {
			// Offering to create it is the obvious next step, and being told to go
			// and run a command is exactly the round trip this screen removes.
			b.pending = target
			b.msg = filepath.Base(target) + " has no board yet — press i to create one here"
			return nil, false
		}
		project.Remember(target)
		return []string{target}, true
	case "i":
		target := b.pending
		if target == "" {
			target = b.selected()
		}
		if target == "" {
			target = b.dir
		}
		if isBoardDir(target) {
			b.msg = filepath.Base(target) + " already has a board"
			return nil, false
		}
		st, err := ticket.At(target)
		if err != nil {
			b.msg = err.Error()
			return nil, false
		}
		if _, err := st.Init(); err != nil {
			b.msg = err.Error()
			return nil, false
		}
		project.Remember(target)
		b.pending = ""
		return []string{target}, true

	case "s":
		found := project.Discover(b.dir, project.MaxScanDepth)
		if len(found) == 0 {
			b.pending = b.dir
			b.msg = fmt.Sprintf("no boards within %d levels of here — press i to create one in %s",
				project.MaxScanDepth, filepath.Base(b.dir))
			return nil, false
		}
		for _, f := range found {
			project.Remember(f)
		}
		return found, true
	}
	return nil, false
}

func (b *browser) render(width, height int) string {
	w := max(20, min(width, 78))
	var sb strings.Builder
	sb.WriteString(styLaneTitle.Render("Add a board") + "\n")
	sb.WriteString(styMeta.Render(truncate(b.dir, w)) + "\n")
	sb.WriteString(styBar.Render(strings.Repeat("─", w)) + "\n")

	if len(b.entries) == 0 {
		sb.WriteString(styMeta.Render("  (no subdirectories)") + "\n")
	}
	// Keep the cursor on screen in a directory with many children.
	visible := max(3, height-8)
	first := 0
	if b.idx >= visible {
		first = b.idx - visible + 1
	}
	for i := first; i < len(b.entries) && i < first+visible; i++ {
		e := b.entries[i]
		lead := "  "
		name := e.name
		if i == b.idx {
			lead = stySelected.Render("▌ ")
			name = stySelected.Render(name)
		}
		mark := ""
		if e.isBoard {
			mark = styOK.Render("  ● board")
		}
		sb.WriteString(truncate(lead+name+mark, w) + "\n")
	}
	if rest := len(b.entries) - (first + visible); rest > 0 {
		sb.WriteString(styMeta.Render(fmt.Sprintf("  +%d more", rest)) + "\n")
	}
	if b.msg != "" {
		sb.WriteString("\n" + styWarn.Render(truncate(b.msg, w)) + "\n")
	}
	sb.WriteString("\n" + styMeta.Render(truncate(
		"enter open dir · h up · a add · i init a board · s scan here · esc cancel", w)))
	return sb.String()
}
