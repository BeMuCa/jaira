package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/BeMuCa/jaira/core/board"
	"github.com/BeMuCa/jaira/core/identity"
	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// laneColWidth is the width of one column in the small board this screen
// draws — narrower than the main board's columns (view.go's minColWidth):
// a column here shows only a lane's name and a short label, never a ticket
// card, so it needs far less room.
const laneColWidth = 16

// laneScreen is the lane settings screen: the project's lanes drawn as the
// board's own shape at a smaller size — narrow columns, navigable left and
// right, with a '+' column at the far right that opens the catalogue. It
// also carries the pre-existing use/publish/adopt/new-lane/drift-refresh
// actions from before this screen became a board, all reusing the same
// core/lane calls the CLI does — one implementation, not two.
type laneScreen struct {
	store *ticket.Store
	// set is kept alongside lanes (set.Lanes, unpacked for convenience
	// everywhere below) because lane.Add/Remove/MoveLane need the whole Set,
	// not just its slice, to resolve an id — see core/lane.Set.Get.
	set   *lane.Set
	lanes []*lane.Lane
	// idx selects a column: 0..len(lanes)-1 is a lane, len(lanes) is the '+'
	// column.
	idx int

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

	// focus picks which list h/l/j/k and the action keys apply to: 0 is this
	// installation's own lanes (the board), 1 is the shared section.
	focus int

	// available is every catalogue lane not part of this board — shown as
	// dimmed columns after the installed ones, because a lane that exists but
	// is invisible until someone finds the '+' column is a lane nobody knows
	// about. Selecting one offers enter to add it and E to edit it.
	available []*lane.Lane

	// catalogue, catalogueIdx and catalogueOpen back the '+' column: pressing
	// enter on it lists every lane not already part of this project
	// (lane.Installable) for one to be chosen and added.
	catalogue     []*lane.Lane
	catalogueIdx  int
	catalogueOpen bool

	// confirmAdoptID is set when an adopt collided with an existing catalogue
	// id, holding the id until the next 'a' confirms the overwrite or esc
	// cancels it.
	confirmAdoptID string

	// confirm holds the yes/no question every 'x' now opens before removing
	// anything, so hammering enter never deletes anything by accident: it
	// starts unset, is set by x, and is cleared by enter (acting on it) or
	// esc (cancelling).
	confirm *confirmRemove

	// pendingCmd carries a tea.Cmd out of key() for the one action that needs
	// one — 'n' launching $EDITOR — without changing key()'s existing
	// bool-only signature that every caller and test already relies on.
	pendingCmd tea.Cmd

	msg   string
	isErr bool
}

// confirmRemove is the pending yes/no question a still-open 'x' has raised.
// path empty means a board removal (lane.Remove, keyed by id); path
// non-empty means a catalogue file delete (os.Remove on path). yes starts
// false so a bare enter answers no.
type confirmRemove struct {
	id   string
	path string
	yes  bool
}

func newLaneScreen(store *ticket.Store, set *lane.Set) *laneScreen {
	ls := &laneScreen{store: store, set: set, lanes: set.Lanes}
	if av, err := lane.Installable(set); err == nil {
		ls.available = av
	}
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

// reload re-Loads the lane set after add/remove/move so the board reflects
// the write immediately, without waiting for the model's own reload (which
// only runs when this screen closes).
func (ls *laneScreen) reload() error {
	set, err := lane.Load(ls.store.Root)
	if err != nil {
		return err
	}
	ls.set = set
	ls.lanes = set.Lanes
	ls.available = nil
	if av, err := lane.Installable(set); err == nil {
		ls.available = av
	}
	if ls.idx > len(ls.lanes)+len(ls.available) {
		ls.idx = len(ls.lanes) + len(ls.available)
	}
	// The agent note names this board's lanes, and this function runs after
	// every change this screen makes to them — adding, removing, reordering,
	// adopting. Regenerating here is what keeps the note true at the moment the
	// pipeline changes rather than at the next 'jaira update'. It writes only
	// when the content actually differs, so the reload that merely opens this
	// screen costs nothing.
	if _, err := board.AnnounceInAgentFiles(ls.store.Root, laneFacts(set)); err != nil {
		ls.msg, ls.isErr = "the lanes changed but the agent note could not be updated: "+err.Error(), true
	}
	return nil
}

func (ls *laneScreen) selected() *lane.Lane {
	if ls.idx < 0 || ls.idx >= len(ls.lanes) {
		return nil
	}
	return ls.lanes[ls.idx]
}

// isPlusColumn reports whether the far-right '+' column is selected.
func (ls *laneScreen) isPlusColumn() bool {
	return ls.idx == len(ls.lanes)+len(ls.available)
}

// selectedAvailable returns the not-installed lane under the cursor, when the
// cursor sits past the installed ones.
func (ls *laneScreen) selectedAvailable() *lane.Lane {
	i := ls.idx - len(ls.lanes)
	if i < 0 || i >= len(ls.available) {
		return nil
	}
	return ls.available[i]
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

func indexOfLane(lanes []*lane.Lane, id string) int {
	for i, l := range lanes {
		if l.ID == id {
			return i
		}
	}
	return -1
}

// key drives the screen. It reports whether the screen is finished (esc/q).
func (ls *laneScreen) key(s string) (done bool) {
	ls.msg = ""
	ls.pendingCmd = nil

	if ls.catalogueOpen {
		switch s {
		case "esc":
			ls.catalogueOpen = false
		case "j", "down":
			if ls.catalogueIdx < len(ls.catalogue)-1 {
				ls.catalogueIdx++
			}
		case "k", "up":
			if ls.catalogueIdx > 0 {
				ls.catalogueIdx--
			}
		case "enter":
			ls.addFromCatalogue()
		}
		return false
	}

	// Every 'x' opens this confirmation before anything is removed, so
	// hammering enter never deletes anything: while it is open, h/l (and
	// left/right) only switch the selected answer, enter acts on it, and
	// esc/q cancel without finishing the screen. Every other key is ignored
	// rather than falling through to the normal board navigation below.
	if ls.confirm != nil {
		switch s {
		case "h", "left":
			ls.confirm.yes = false
		case "l", "right":
			ls.confirm.yes = true
		case "enter":
			c := ls.confirm
			ls.confirm = nil
			if c.yes {
				if c.path != "" {
					ls.deleteCatalogueLane(c.path)
				} else {
					ls.removeBoardLane(c.id)
				}
			}
		case "esc", "q":
			ls.confirm = nil
		}
		return false
	}

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
		if ls.focus == 1 {
			ls.moveSharedDown()
		}
	case "k", "up":
		if ls.focus == 1 {
			ls.moveSharedUp()
		}
	case "h", "left":
		if ls.focus == 0 {
			ls.moveColumn(-1)
		}
	case "l", "right":
		if ls.focus == 0 {
			ls.moveColumn(1)
		}
	case "H":
		if ls.focus == 0 {
			ls.moveLane(-1)
		}
	case "L":
		if ls.focus == 0 {
			ls.moveLane(1)
		}
	case "x":
		if ls.focus == 0 {
			ls.startRemove()
		}
	case "enter":
		if ls.focus == 0 && ls.isPlusColumn() {
			ls.openCatalogue()
		} else if ls.focus == 0 && ls.selectedAvailable() != nil {
			ls.addAvailable()
		}
	case "p":
		if ls.focus == 0 {
			ls.publish()
		}
	case "n":
		if ls.focus == 0 {
			ls.pendingCmd = ls.newLane()
		}
	case "R":
		if ls.focus == 0 {
			ls.refreshDrift()
		}
	case "E":
		if ls.focus == 0 {
			ls.pendingCmd = ls.editLane()
		}
	case "a":
		if ls.focus == 1 {
			ls.adopt()
		}
	}
	return false
}

// moveColumn shifts which column (a lane, or the '+' column past the last
// lane) is selected — h/l, the board's own navigation vocabulary.
func (ls *laneScreen) moveColumn(delta int) {
	ls.confirmAdoptID = ""
	next := ls.idx + delta
	if next < 0 || next > len(ls.lanes)+len(ls.available) {
		return
	}
	ls.idx = next
}

func (ls *laneScreen) moveSharedDown() {
	if ls.sharedIdx < len(ls.shared)-1 {
		ls.sharedIdx++
	}
}

func (ls *laneScreen) moveSharedUp() {
	if ls.sharedIdx > 0 {
		ls.sharedIdx--
	}
}

// startRemove opens the yes/no confirmation every 'x' now requires, for
// either an installed lane under the cursor (board removal) or an available
// (not-installed) catalogue lane (file delete). A built-in has no file — it
// lives inside the binary — so it errors immediately rather than opening a
// confirmation there is nothing to act on. The '+' column has nothing to
// remove either.
func (ls *laneScreen) startRemove() {
	if ls.isPlusColumn() {
		return
	}
	if l := ls.selected(); l != nil {
		ls.confirm = &confirmRemove{id: l.ID}
		return
	}
	l := ls.selectedAvailable()
	if l == nil {
		return
	}
	if l.Builtin {
		ls.msg, ls.isErr = fmt.Sprintf("built-in lane %q lives inside the binary — nothing to delete", l.ID), true
		return
	}
	ls.confirm = &confirmRemove{id: l.ID, path: l.Source}
}

// removeBoardLane takes id out of this project through lane.Remove — the
// same call 'jaira lanes remove' makes — refusing with the same message when
// a ticket currently sits in it. Runs once the confirmation startRemove
// opened for a board lane is answered yes.
func (ls *laneScreen) removeBoardLane(id string) {
	if _, err := lane.Remove(ls.store.Root, ls.set, ls.store, id); err != nil {
		ls.msg, ls.isErr = err.Error(), true
		return
	}
	if err := ls.reload(); err != nil {
		ls.msg, ls.isErr = err.Error(), true
		return
	}
	ls.msg, ls.isErr = "removed "+id, false
}

// deleteCatalogueLane deletes path — an available (not-installed) catalogue
// lane's file — from disk. Available lanes only come from builtins (which
// have no file, and are refused before a confirmation ever opens, see
// startRemove) and the UserLanesDir() glob (core/lane.Installable's other
// source), so path is always a real catalogue file here — no extra path
// guard needed. Runs once the confirmation startRemove opened for an
// available lane is answered yes.
func (ls *laneScreen) deleteCatalogueLane(path string) {
	if err := os.Remove(path); err != nil {
		ls.msg, ls.isErr = err.Error(), true
		return
	}
	if err := ls.reload(); err != nil {
		ls.msg, ls.isErr = err.Error(), true
		return
	}
	ls.msg, ls.isErr = "deleted "+path, false
}

// moveLane shifts the selected lane one step through lane.MoveLane — the
// same call 'jaira lanes move' makes — and keeps it selected at its new
// position.
func (ls *laneScreen) moveLane(delta int) {
	l := ls.selected()
	if l == nil {
		return
	}
	id := l.ID
	if err := lane.MoveLane(ls.store.Root, ls.set, id, delta); err != nil {
		ls.msg, ls.isErr = err.Error(), true
		return
	}
	if err := ls.reload(); err != nil {
		ls.msg, ls.isErr = err.Error(), true
		return
	}
	if i := indexOfLane(ls.lanes, id); i >= 0 {
		ls.idx = i
	}
}

// openCatalogue lists every lane not already part of this project
// (lane.Installable) for the '+' column's enter key — the only way in for a
// user who does not know 'jaira lanes template' exists.
func (ls *laneScreen) openCatalogue() {
	installable, err := lane.Installable(ls.set)
	if err != nil {
		ls.msg, ls.isErr = err.Error(), true
		return
	}
	ls.catalogue = installable
	ls.catalogueIdx = 0
	ls.catalogueOpen = true
}

// addFromCatalogue adds the highlighted catalogue lane through lane.Add —
// the same call 'jaira lanes add' makes — appending it at the end of the
// order, and selects it.
func (ls *laneScreen) addFromCatalogue() {
	if ls.catalogueIdx < 0 || ls.catalogueIdx >= len(ls.catalogue) {
		return
	}
	l := ls.catalogue[ls.catalogueIdx]
	ls.catalogueOpen = false
	if _, err := lane.Add(ls.store.Root, ls.set, l.ID); err != nil {
		ls.msg, ls.isErr = err.Error(), true
		return
	}
	id := l.ID
	if err := ls.reload(); err != nil {
		ls.msg, ls.isErr = err.Error(), true
		return
	}
	if i := indexOfLane(ls.lanes, id); i >= 0 {
		ls.idx = i
	}
	ls.msg, ls.isErr = "added "+id, false
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
	ls.msg, ls.isErr = "published to "+dst, false
}

// newLaneDoneMsg reports that $EDITOR, opened on a freshly written lane
// skeleton from the lane settings screen's 'n' key, has exited.
type newLaneDoneMsg struct{ err error }

// newLane writes a fresh lane skeleton into the catalogue and opens it in
// $EDITOR — until now, doing this meant knowing 'jaira lanes template'
// exists and combining it with tools outside jaira.
// addAvailable installs the not-installed lane under the cursor, the same
// call the '+' catalogue and 'jaira lanes add' make.
func (ls *laneScreen) addAvailable() {
	l := ls.selectedAvailable()
	if l == nil {
		return
	}
	if _, err := lane.Add(ls.store.Root, ls.set, l.ID); err != nil {
		ls.msg, ls.isErr = err.Error(), true
		return
	}
	id := l.ID
	if err := ls.reload(); err != nil {
		ls.msg, ls.isErr = err.Error(), true
		return
	}
	ls.idx = indexOfLane(ls.lanes, id)
	ls.msg, ls.isErr = "added "+id, false
}

// editLane opens the selected lane's file in $EDITOR. A built-in has no file —
// it lives inside the binary — so editing one first writes a catalogue
// override copy under ~/.jaira/lanes (or opens the copy a previous edit left)
// and edits that, which is the same override mechanism a hand-written file of
// the same id uses.
func (ls *laneScreen) editLane() tea.Cmd {
	l := ls.selected()
	if l == nil {
		l = ls.selectedAvailable()
	}
	if l == nil {
		return nil
	}
	path := l.Source
	if l.Builtin {
		dst := filepath.Join(lane.UserLanesDir(), l.ID+".md")
		if _, err := os.Stat(dst); os.IsNotExist(err) {
			exported, err := lane.Export(l, lane.UserLanesDir(), false)
			if err != nil {
				ls.msg, ls.isErr = err.Error(), true
				return nil
			}
			dst = exported
		}
		path = dst
	}
	argv := append(editorCommand(), path)
	cmd := exec.Command(argv[0], argv[1:]...)
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return newLaneDoneMsg{err: err}
	})
}

// newLane writes a lane skeleton into this board's own lane directory and
// opens it: a board is its lane directory, so a file there is on the board
// the moment the editor closes. To offer it to other boards, publish it or
// copy it into the catalogue.
func (ls *laneScreen) newLane() tea.Cmd {
	path, err := writeLaneSkeleton(lane.ProjectLanesDir(ls.store.Root))
	if err != nil {
		ls.msg, ls.isErr = err.Error(), true
		return nil
	}
	argv := append(editorCommand(), path)
	cmd := exec.Command(argv[0], argv[1:]...)
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return newLaneDoneMsg{err: err}
	})
}

// writeLaneSkeleton creates a new lane skeleton file in dir under the first
// free name — "new-lane.md", then "new-lane-2.md", "new-lane-3.md", … — so
// pressing n more than once writes another file rather than clobbering the
// last one. O_EXCL makes the creation itself the collision check, so a race
// with another writer cannot silently overwrite an existing lane.
func writeLaneSkeleton(dir string) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	for n := 1; ; n++ {
		name := "new-lane.md"
		if n > 1 {
			name = fmt.Sprintf("new-lane-%d.md", n)
		}
		path := filepath.Join(dir, name)
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if os.IsExist(err) {
			continue
		}
		if err != nil {
			return "", err
		}
		_, writeErr := f.WriteString(lane.Template)
		closeErr := f.Close()
		if writeErr != nil {
			return "", writeErr
		}
		if closeErr != nil {
			return "", closeErr
		}
		return path, nil
	}
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
	ls.msg, ls.isErr = "pulled the catalogue copy of "+l.ID, false
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
	ls.msg, ls.isErr = "adopted into "+dst, false
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

// laneColumnStyle borders one column of the small board this screen draws.
// It is a fresh, smaller-scale style rather than a reuse of view.go's
// renderColumn: renderColumn draws a ticket card list from the main Model's
// own state (tickets, scroll position, card styling), which has no
// equivalent here — this screen has only lane names, never tickets.
func laneColumnStyle(focused bool) lipgloss.Style {
	st := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(colFaint).
		Width(laneColWidth).Height(3)
	if focused {
		st = st.BorderForeground(colAccent)
	}
	return st
}

// renderBoard draws the project's lanes as narrow columns — the main
// board's own shape (bordered columns joined horizontally, per view.go's
// renderBoard) at a smaller size — plus one more column holding '+'. As many
// columns as fit are shown, scrolled so the selection stays visible.
func (ls *laneScreen) renderBoard(totalW int) string {
	cols := make([]string, 0, len(ls.lanes)+1)
	for i, l := range ls.lanes {
		focused := ls.focus == 0 && i == ls.idx
		name := truncate(l.Name, laneColWidth-2)
		if focused {
			name = stySelected.Render(name)
		}
		label := ls.sourceLabel(l)
		if l.Overrides != "" {
			label = "changed"
		}
		if _, drifted := ls.drift[l.ID]; drifted {
			label = "drifted"
		}
		body := name + "\n" + styMeta.Render(truncate(label, laneColWidth-2))
		cols = append(cols, laneColumnStyle(focused).Render(body))
	}
	// Lanes that exist but are not on this board are shown dimmed rather than
	// hidden behind the '+' column: a lane nobody can see is a lane nobody
	// knows exists. enter adds one, E edits it.
	for i, l := range ls.available {
		focused := ls.focus == 0 && ls.idx == len(ls.lanes)+i
		name := truncate(l.Name, laneColWidth-2)
		if focused {
			name = stySelected.Render(name)
		} else {
			name = styMeta.Render(name)
		}
		body := name + "\n" + styMeta.Render(truncate("not on board", laneColWidth-2))
		cols = append(cols, laneColumnStyle(focused).Render(body))
	}
	plusFocused := ls.focus == 0 && ls.isPlusColumn()
	plusLabel := "+"
	if plusFocused {
		plusLabel = stySelected.Render("+")
	}
	cols = append(cols, laneColumnStyle(plusFocused).Render(plusLabel))

	// Columns wrap onto further rows rather than being clipped at the right
	// edge: a narrow terminal used to silently hide the tail of the board —
	// review, human review, done — which read as those lanes not existing.
	perRow := max(1, totalW/(laneColWidth+2))
	var rows []string
	for start := 0; start < len(cols); start += perRow {
		end := min(len(cols), start+perRow)
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cols[start:end]...))
	}
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

// renderCatalogue lists the lanes the '+' column's enter key found
// installable — every built-in and catalogue lane not already part of this
// project.
func (ls *laneScreen) renderCatalogue(w int) string {
	var sb strings.Builder
	sb.WriteString(styLaneTitle.Render("Add a lane") + "\n")
	if len(ls.catalogue) == 0 {
		sb.WriteString(styMeta.Render(truncate(
			"  every installed lane is already part of this project", w)) + "\n")
		return sb.String()
	}
	for i, l := range ls.catalogue {
		lead := "  "
		name := l.Name
		if i == ls.catalogueIdx {
			lead = stySelected.Render("▌ ")
			name = stySelected.Render(name)
		}
		sb.WriteString(truncate(lead+name, w) + "\n")
	}
	return sb.String()
}

func (ls *laneScreen) render(width, height int) string {
	w := max(20, width)
	var sb strings.Builder
	sb.WriteString(styLaneTitle.Render("Lanes") + "\n")
	sb.WriteString(styBar.Render(strings.Repeat("─", w)) + "\n")

	sb.WriteString(ls.renderBoard(w))
	sb.WriteString("\n")

	if ls.catalogueOpen {
		sb.WriteString(styBar.Render(strings.Repeat("─", w)) + "\n")
		sb.WriteString(ls.renderCatalogue(w))
		sb.WriteString("\n" + styMeta.Render(truncate("jk select · enter add · esc cancel", w)))
		return sb.String()
	}

	sb.WriteString(styBar.Render(strings.Repeat("─", w)) + "\n")
	sb.WriteString(styLaneTitle.Render("Shared by teammates") + "\n")
	if len(ls.shared) == 0 {
		sb.WriteString(styMeta.Render(truncate(
			"  none yet — a teammate publishes a lane with p, then commits .jaira/shared/", w)) + "\n")
	} else {
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
		promptHeight := max(3, height-16)
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

	// While a removal is pending, the footer becomes the yes/no question
	// itself rather than the normal key list — nothing else on this screen
	// should be actionable until it is answered. "no" is always drawn first
	// so hammering enter (default no) reads the same as it behaves.
	if ls.confirm != nil {
		question := "remove " + ls.confirm.id + " from this board?"
		if ls.confirm.path != "" {
			question = "delete " + ls.confirm.path + "?"
		}
		sb.WriteString("\n" + truncate(question, w) + "\n")
		no, yes := "no", "yes"
		if ls.confirm.yes {
			no, yes = styMeta.Render(no), stySelected.Render(yes)
		} else {
			no, yes = stySelected.Render(no), styMeta.Render(yes)
		}
		sb.WriteString(no + "  " + yes + "\n")
		for _, l := range wrapHints([]string{"←/→ choose", "enter confirm", "esc cancel"}, max(1, w)) {
			sb.WriteString("\n" + styMeta.Render(l))
		}
		return sb.String()
	}

	// Not truncated to w like everything else on this screen: with add/move/
	// remove now joining the pre-existing use/publish/new/refresh/adopt keys,
	// the full footer runs longer than a narrow terminal's column width, and
	// cutting it off would silently hide a key's name rather than the column
	// content truncate exists to fit.
	help := []string{"E edit", "x remove", "p publish", "n new", "R refresh", "esc back"}
	if len(ls.shared) > 0 {
		help = []string{"E edit", "x remove", "tab switch", "p publish", "a adopt", "esc back"}
	}
	for _, l := range wrapHints(help, max(1, w)) {
		sb.WriteString("\n" + styMeta.Render(l))
	}
	return sb.String()
}
