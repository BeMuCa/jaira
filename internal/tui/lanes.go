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
// in this project, publish it to a teammate, adopt one of theirs. Following
// browse.go's shape: its own state, a key method, a render method, no
// bubbletea imports of its own.
type laneScreen struct {
	store *ticket.Store
	lanes []*lane.Lane
	idx   int

	// drift implements D-02: lanes whose project copy no longer matches their
	// catalogue copy of the same id, keyed by lane id. Computed once when the
	// screen opens, per the decided design — not on every command.
	drift map[string]lane.DriftEntry

	// shared is every lane found under .jaira/shared/, visible and adoptable
	// but never loaded onto the board (see lane.Shared). sharedWarnings names
	// the files that failed to parse and were skipped.
	shared         []lane.SharedLane
	sharedIdx      int
	sharedWarnings []string

	// focus picks which list j/k and the action keys apply to: 0 is this
	// installation's own lanes, 1 is the shared section.
	focus int

	// confirmAdoptID is set when an adopt collided with an existing catalogue
	// id, holding the id until the next 'a' confirms the overwrite or esc
	// cancels it.
	confirmAdoptID string

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
	if sh, warn, err := lane.Shared(store.Root); err == nil {
		ls.shared = sh
		ls.sharedWarnings = warn
	}
	return ls
}

func (ls *laneScreen) selected() *lane.Lane {
	if ls.idx < 0 || ls.idx >= len(ls.lanes) {
		return nil
	}
	return ls.lanes[ls.idx]
}

func (ls *laneScreen) selectedShared() *lane.SharedLane {
	if ls.sharedIdx < 0 || ls.sharedIdx >= len(ls.shared) {
		return nil
	}
	return &ls.shared[ls.sharedIdx]
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
		if ls.confirmAdoptID != "" {
			ls.confirmAdoptID = ""
			return false
		}
		return true
	case "tab":
		if len(ls.shared) > 0 {
			ls.focus = 1 - ls.focus
		}
	case "j", "down":
		ls.moveDown()
	case "k", "up":
		ls.moveUp()
	case "u":
		if ls.focus == 0 {
			ls.use()
		}
	case "p":
		if ls.focus == 0 {
			ls.publish()
		}
	case "R":
		if ls.focus == 0 {
			ls.refreshDrift()
		}
	case "a":
		if ls.focus == 1 {
			ls.adopt()
		}
	}
	return false
}

func (ls *laneScreen) moveDown() {
	ls.confirmAdoptID = ""
	if ls.focus == 1 {
		if ls.sharedIdx < len(ls.shared)-1 {
			ls.sharedIdx++
		}
		return
	}
	if ls.idx < len(ls.lanes)-1 {
		ls.idx++
	}
}

func (ls *laneScreen) moveUp() {
	ls.confirmAdoptID = ""
	if ls.focus == 1 {
		if ls.sharedIdx > 0 {
			ls.sharedIdx--
		}
		return
	}
	if ls.idx > 0 {
		ls.idx--
	}
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

// adopt copies a teammate's shared lane into this catalogue. The prompt was
// already shown before this key could be pressed — adopting means agreeing
// to run someone else's instructions at whatever tier the file declares, and
// that agreement has to come after reading it, not before.
//
// A collision with an existing catalogue id is not overwritten on the first
// press: it is held in confirmAdoptID until the next 'a' confirms it, or esc
// cancels.
func (ls *laneScreen) adopt() {
	sl := ls.selectedShared()
	if sl == nil {
		return
	}
	overwrite := ls.confirmAdoptID == sl.Lane.ID
	_, dst, err := lane.Adopt(sl.Path, lane.UserLanesDir(), overwrite)
	if err != nil {
		if !overwrite {
			ls.confirmAdoptID = sl.Lane.ID
			ls.msg, ls.isErr = sl.Lane.ID+" already exists in your catalogue; press a again to overwrite, esc to cancel", true
			return
		}
		ls.msg, ls.isErr = err.Error(), true
		return
	}
	ls.confirmAdoptID = ""
	ls.msg, ls.isErr = "adopted into " + dst, false
}

// promptOf returns the name, id, creator and prompt to show in the bottom
// pane for whichever list currently has focus, so the lane settings screen
// and lane settings adoption preview read from exactly the same fields.
func (ls *laneScreen) promptOf() (name, id, creator, prompt string, ok bool) {
	if ls.focus == 1 {
		sl := ls.selectedShared()
		if sl == nil {
			return "", "", "", "", false
		}
		return sl.Lane.Name, sl.Lane.ID, sl.Lane.Creator, sl.Lane.Prompt, true
	}
	l := ls.selected()
	if l == nil {
		return "", "", "", "", false
	}
	return l.Name, l.ID, l.Creator, l.Prompt, true
}

func (ls *laneScreen) render(width, height int) string {
	w := max(20, min(width, 78))
	var sb strings.Builder
	sb.WriteString(styLaneTitle.Render("Lanes") + "\n")
	sb.WriteString(styBar.Render(strings.Repeat("─", w)) + "\n")

	listHeight := max(3, (height-14)/3)
	first := 0
	if ls.idx >= listHeight {
		first = ls.idx - listHeight + 1
	}
	for i := first; i < len(ls.lanes) && i < first+listHeight; i++ {
		l := ls.lanes[i]
		lead := "  "
		name := l.Name
		if ls.focus == 0 && i == ls.idx {
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

	if len(ls.shared) > 0 {
		sb.WriteString(styBar.Render(strings.Repeat("─", w)) + "\n")
		sb.WriteString(styLaneTitle.Render("Shared by teammates") + "\n")
		for i, sl := range ls.shared {
			lead := "  "
			label := fmt.Sprintf("%s/%s", sl.Folder, sl.Lane.ID)
			if ls.focus == 1 && i == ls.sharedIdx {
				lead = stySelected.Render("▌ ")
				label = stySelected.Render(label)
			}
			tail := ""
			if sl.Lane.Creator != "" {
				tail = styMeta.Render("  creator: " + sl.Lane.Creator)
			}
			sb.WriteString(truncate(lead+label+tail, w) + "\n")
		}
		for _, warn := range ls.sharedWarnings {
			sb.WriteString(styWarn.Render(truncate("  ⚠ "+warn, w)) + "\n")
		}
	}

	sb.WriteString(styBar.Render(strings.Repeat("─", w)) + "\n")
	if name, id, creator, prompt, ok := ls.promptOf(); ok {
		sb.WriteString(styLaneTitle.Render(fmt.Sprintf("%s (%s)", name, id)) + "\n")
		if creator != "" {
			sb.WriteString(styMeta.Render("creator: "+creator) + "\n")
		}
		promptHeight := max(3, height-listHeight-16)
		lines := strings.Split(prompt, "\n")
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
	help := "jk select · u use here · p publish · R refresh drift · esc back"
	if len(ls.shared) > 0 {
		help = "jk select · tab switch list · u use · p publish · a adopt · R refresh drift · esc back"
	}
	sb.WriteString("\n" + styMeta.Render(truncate(help, w)))
	return sb.String()
}
