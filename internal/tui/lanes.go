package tui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/BeMuCa/jaira/core/identity"
	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// laneScreen is the lane settings screen: read a lane, see its prompt, use it
// in this project, publish it to a teammate. Following browse.go's shape: its
// own state, a key method, a render method, no bubbletea imports of its own.
type laneScreen struct {
	store *ticket.Store
	lanes []*lane.Lane
	idx   int

	// drift implements D-02: lanes whose project copy no longer matches their
	// catalogue copy of the same id, keyed by lane id. Computed once when the
	// screen opens, per the decided design — not on every command.
	drift map[string]lane.DriftEntry

	msg   string
	isErr bool
}

func newLaneScreen(store *ticket.Store, set *lane.Set) *laneScreen {
	ls := &laneScreen{store: store, lanes: set.Lanes}
	if d, err := lane.Drift(store.Root, set); err == nil {
		ls.drift = make(map[string]lane.DriftEntry, len(d))
		for _, e := range d {
			ls.drift[e.ID] = e
		}
	}
	return ls
}

func (ls *laneScreen) selected() *lane.Lane {
	if ls.idx < 0 || ls.idx >= len(ls.lanes) {
		return nil
	}
	return ls.lanes[ls.idx]
}

// sourceLabel names where a lane came from, the same three sources core/lane
// resolves: built-in, this project's own directory, or the catalogue.
func (ls *laneScreen) sourceLabel(l *lane.Lane) string {
	switch {
	case l.Builtin:
		return "built-in"
	case strings.HasPrefix(l.Source, lane.ProjectLanesDir(ls.store.Root)):
		return "project"
	default:
		return "catalogue"
	}
}

// key drives the screen. It reports whether the screen is finished (esc/q).
func (ls *laneScreen) key(s string) (done bool) {
	ls.msg = ""
	switch s {
	case "esc", "q":
		return true
	case "j", "down":
		if ls.idx < len(ls.lanes)-1 {
			ls.idx++
		}
	case "k", "up":
		if ls.idx > 0 {
			ls.idx--
		}
	case "u":
		ls.use()
	case "p":
		ls.publish()
	case "R":
		ls.refreshDrift()
	}
	return false
}

// use exports the selected lane into this project's own lane directory —
// "use this lane here" is a copy with a confirmation, not its own command.
func (ls *laneScreen) use() {
	l := ls.selected()
	if l == nil {
		return
	}
	dst, err := lane.Export(l, lane.ProjectLanesDir(ls.store.Root), false)
	if err != nil {
		ls.msg, ls.isErr = err.Error(), true
		return
	}
	ls.msg, ls.isErr = "wrote "+dst, false
}

// publish copies the selected lane to .jaira/shared/<slug>/, the deliberate,
// opt-in hand-off to teammates the design note describes.
func (ls *laneScreen) publish() {
	l := ls.selected()
	if l == nil {
		return
	}
	who := identity.Slug(identity.Current(ls.store.Root))
	dstDir := filepath.Join(ls.store.SharedDir(), who)
	dst, err := lane.Publish(l, dstDir, who, false)
	if err != nil {
		ls.msg, ls.isErr = err.Error(), true
		return
	}
	ls.msg, ls.isErr = "published to " + dst, false
}

// refreshDrift pulls the catalogue's copy of the selected lane into the
// project, in the direction the user asked for. Nothing syncs on its own.
func (ls *laneScreen) refreshDrift() {
	l := ls.selected()
	if l == nil {
		return
	}
	d, ok := ls.drift[l.ID]
	if !ok {
		return
	}
	if err := lane.RefreshDrift(d); err != nil {
		ls.msg, ls.isErr = err.Error(), true
		return
	}
	delete(ls.drift, l.ID)
	ls.msg, ls.isErr = "pulled the catalogue copy of " + l.ID, false
}

func (ls *laneScreen) render(width, height int) string {
	w := max(20, min(width, 78))
	var sb strings.Builder
	sb.WriteString(styLaneTitle.Render("Lanes") + "\n")
	sb.WriteString(styBar.Render(strings.Repeat("─", w)) + "\n")

	listHeight := max(3, (height-10)/2)
	first := 0
	if ls.idx >= listHeight {
		first = ls.idx - listHeight + 1
	}
	for i := first; i < len(ls.lanes) && i < first+listHeight; i++ {
		l := ls.lanes[i]
		lead := "  "
		name := l.Name
		if i == ls.idx {
			lead = stySelected.Render("▌ ")
			name = stySelected.Render(name)
		}
		tail := styMeta.Render("  " + ls.sourceLabel(l))
		if l.Overrides != "" {
			tail += styWarn.Render("  overrides " + l.Overrides)
		}
		if _, drifted := ls.drift[l.ID]; drifted {
			tail += styWarn.Render("  drifted from catalogue")
		}
		sb.WriteString(truncate(lead+name+tail, w) + "\n")
	}
	if rest := len(ls.lanes) - (first + listHeight); rest > 0 {
		sb.WriteString(styMeta.Render(fmt.Sprintf("  +%d more", rest)) + "\n")
	}

	sb.WriteString(styBar.Render(strings.Repeat("─", w)) + "\n")
	if l := ls.selected(); l != nil {
		sb.WriteString(styLaneTitle.Render(fmt.Sprintf("%s (%s)", l.Name, l.ID)) + "\n")
		if l.Creator != "" {
			sb.WriteString(styMeta.Render("creator: "+l.Creator) + "\n")
		}
		promptHeight := max(3, height-listHeight-12)
		lines := strings.Split(l.Prompt, "\n")
		for i, ln := range lines {
			if i >= promptHeight {
				sb.WriteString(styMeta.Render(fmt.Sprintf("  … +%d more line(s)", len(lines)-promptHeight)) + "\n")
				break
			}
			sb.WriteString(truncate(ln, w) + "\n")
		}
	}

	if ls.msg != "" {
		style := styOK
		if ls.isErr {
			style = styErr
		}
		sb.WriteString("\n" + style.Render(truncate(ls.msg, w)) + "\n")
	}
	sb.WriteString("\n" + styMeta.Render(truncate(
		"jk select · u use here · p publish · R refresh drift · esc back", w)))
	return sb.String()
}
