package tui

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/milestone"
	"github.com/BeMuCa/jaira/core/ticket"
)

// The milestone layer is the board's second view: every placed milestone as a
// bubble, arrows for "first this, then that". It answers what the board does
// not — what comes after this round of work, and how far along each round
// is. The arrows are a picture only: nothing on this screen blocks or orders
// a ticket.
//
// Every write goes through milestone.Create and milestone.Edit, the path the
// CLI writes through, so what a person does here is exactly what an agent
// does with 'jaira milestone create' and 'jaira milestone set'.

// layerKind says what the cursor rests on: a bubble or one of the [+] targets.
type layerKind int

const (
	onBubble layerKind = iota
	onRight            // [+] right of a bubble: a sibling, same first predecessor
	onBelow            // [+] under a bubble: a successor
	onFar              // [+] after the last swimlane: a swimlane of its own
)

// layerTarget is one thing the cursor can rest on. A [+] is named after the
// bubble it belongs to, so it survives a reload the way the bubble does.
type layerTarget struct {
	kind layerKind
	name string // the bubble, or the bubble a [+] belongs to; "" for onFar
}

// layerStage is what the layer is doing on top of the picture.
type layerStage int

const (
	stageLayer  layerStage = iota
	stageChoose            // enter on a [+]: a new milestone or an existing one
	stageForm              // name and goal, or goal and predecessors
	stagePick              // the milestones that can be connected here
	stageRemove            // x: asked before the bubble leaves the layer
	stageDetail            // enter on a bubble: its goal and its tickets
)

// layerView is the open layer.
type layerView struct {
	sel layerTarget

	// panX and panY are the first canvas column and row on screen, stored
	// rather than derived from the cursor so the picture holds still while the
	// cursor moves inside it; the render moves them only when the selection
	// would otherwise leave the window.
	panX, panY int

	// filed maps every ticket id filed away — logbook or archive, any day —
	// to its file. Read off the file names when the layer opens and on every
	// reload while it is open: the board itself loads only the last days of
	// the logbook, and a milestone whose tickets were filed a month ago is
	// still done.
	filed map[string]string

	stage layerStage
	form  layerForm

	pick    []string // the names the connect picker offers
	pickIdx int

	// The milestone detail: whose it is, its ticket rows, the cursor and the
	// first row on screen.
	detailName string
	rows       []memberRow
	rowIdx     int
	rowTop     int

	err string // why the last choice was refused, shown where it was made
}

// layerForm is the two-field form behind a new milestone (name, goal) and
// behind e (goal, predecessors).
type layerForm struct {
	edit   string // the milestone e is editing; "" for a new one
	labels [2]string
	vals   [2]string
	field  int
}

// memberRow is one ticket in the milestone detail.
type memberRow struct {
	id, title string
	where     string // the lane's name, "logbook", or why it is nowhere
}

// openLayer shows the layer, the cursor on the first bubble — or on the one
// [+] an empty layer has.
func (m *Model) openLayer() {
	m.layer = &layerView{filed: m.store.FiledAwayIDs()}
	m.mode = modeLayer
	m.layerScene() // settles the cursor on something that exists
}

// closeLayer goes back to the board.
func (m *Model) closeLayer() {
	m.layer = nil
	m.mode = modeBoard
}

// refreshLayer is reload's part for an open layer: the logbook may have grown,
// and an open detail lists what its milestone now holds.
func (m *Model) refreshLayer() {
	l := m.layer
	if l == nil {
		return
	}
	l.filed = m.store.FiledAwayIDs()
	if l.stage == stageDetail {
		ms := m.milestoneNamed(l.detailName)
		if ms == nil {
			// Taken off the board elsewhere: there is nothing left to list.
			l.stage = stageLayer
			return
		}
		l.rows = m.memberRows(ms)
		l.rowIdx = min(l.rowIdx, max(0, len(l.rows)-1))
	}
}

func (m *Model) milestoneNamed(name string) *milestone.Milestone {
	for _, ms := range m.milestones {
		if ms.Name == name {
			return ms
		}
	}
	return nil
}

// loggedFile reports whether a filed-away file lies in the logbook rather
// than in the archive. An archived ticket was dropped, not finished, so it
// counts for nothing.
func (m *Model) loggedFile(path string) bool {
	return filepath.Dir(path) != m.store.ArchiveDir()
}

// milestoneDone counts a milestone's tickets and how many of them are done:
// in the board's terminal lane, or filed into the logbook on any day.
func (m *Model) milestoneDone(ms *milestone.Milestone, status map[string]string) (done, total int) {
	terminal := ""
	if t := m.lanes.Terminal(); t != nil {
		terminal = t.ID
	}
	for _, id := range ms.Members() {
		total++
		if path, ok := m.layer.filed[id]; ok && m.loggedFile(path) {
			done++
		} else if st, ok := status[id]; ok && terminal != "" && st == terminal {
			done++
		}
	}
	return done, total
}

// keyLayer drives the layer and whatever is open on top of it.
func (m *Model) keyLayer(k tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	s := cmdKey(k)
	l := m.layer
	if l == nil {
		m.mode = modeBoard
		return m, nil
	}
	if s == "ctrl+c" {
		m.Close()
		return m, tea.Quit
	}
	switch l.stage {
	case stageChoose:
		m.keyLayerChoose(s)
		return m, nil
	case stageForm:
		m.keyLayerForm(s, k.Text)
		return m, nil
	case stagePick:
		m.keyLayerPick(s)
		return m, nil
	case stageRemove:
		m.keyLayerRemove(s)
		return m, nil
	case stageDetail:
		m.keyMilestoneDetail(s)
		return m, nil
	}

	switch s {
	case "L", "esc", "q":
		m.closeLayer()
	case "left", "right", "up", "down":
		l.sel = m.layerScene().step(l.sel, s)
	case "enter":
		l.err = ""
		if l.sel.kind == onBubble {
			m.openMilestoneDetail(l.sel.name)
		} else {
			l.stage = stageChoose
		}
	case "e":
		if l.sel.kind == onBubble {
			m.startLayerEdit(l.sel.name)
		}
	case "c":
		if l.sel.kind == onBubble {
			m.nextMilestoneColour(l.sel.name)
		}
	case "x":
		if l.sel.kind == onBubble {
			l.stage = stageRemove
		}
	}
	return m, nil
}

// placement is where a [+] puts what is chosen on it: the predecessor list a
// new milestone is written with. A sibling shares the bubble's first
// predecessor, a successor comes after the bubble itself, and the far [+]
// starts a swimlane of its own — no predecessor at all.
func (m *Model) placement(t layerTarget) []string {
	switch t.kind {
	case onRight:
		if ms := m.milestoneNamed(t.name); ms != nil && len(ms.After) > 0 {
			return []string{ms.After[0]}
		}
	case onBelow:
		return []string{t.name}
	}
	return []string{}
}

// placementText says in words what a [+] would create, for the header and
// the dialogs behind it.
func (m *Model) placementText(t layerTarget) string {
	if after := m.placement(t); len(after) > 0 {
		return "after " + after[0]
	}
	return "a swimlane of its own"
}

func (m *Model) keyLayerChoose(s string) {
	l := m.layer
	switch s {
	case "esc", "q":
		l.stage = stageLayer
	case "n":
		l.form = layerForm{labels: [2]string{"name", "goal"}}
		l.err = ""
		l.stage = stageForm
	case "c":
		l.pick = m.connectable(l.sel)
		l.pickIdx = 0
		l.err = ""
		l.stage = stagePick
	}
}

// startLayerEdit opens the form behind e on a milestone: its goal and its
// predecessors as they stand.
func (m *Model) startLayerEdit(name string) {
	ms := m.milestoneNamed(name)
	if ms == nil {
		return
	}
	l := m.layer
	l.form = layerForm{
		edit:   name,
		labels: [2]string{"goal", "after"},
		vals:   [2]string{ms.Goal, strings.Join(ms.After, ", ")},
	}
	l.err = ""
	l.stage = stageForm
}

func (m *Model) keyLayerForm(s, text string) {
	l := m.layer
	f := &l.form
	switch s {
	case "esc":
		l.stage = m.formReturn()
	case "tab", "shift+tab", "up", "down":
		f.field = 1 - f.field
	case "backspace":
		if r := []rune(f.vals[f.field]); len(r) > 0 {
			f.vals[f.field] = string(r[:len(r)-1])
		}
	case "enter":
		m.submitLayerForm()
	default:
		m.insertText(text)
	}
}

// formReturn is where a closed form goes back to: e can be pressed in the
// detail as well as on the bubble.
func (m *Model) formReturn() layerStage {
	if m.layer.detailName != "" && m.layer.form.edit == m.layer.detailName {
		return stageDetail
	}
	return stageLayer
}

func (m *Model) submitLayerForm() {
	l := m.layer
	f := l.form
	if f.edit == "" {
		name, _, err := milestone.NormalizeName(strings.TrimSpace(f.vals[0]))
		if strings.TrimSpace(f.vals[0]) == "" {
			err = errors.New("a milestone needs a name")
		}
		if err != nil {
			l.err = err.Error()
			return
		}
		if err := m.createMilestone(name, f.vals[1], m.placement(l.sel)); err != nil {
			l.err = err.Error()
			return
		}
		m.layerWritten(name, stageLayer)
		return
	}

	preds, err := parseLayerNames(f.vals[1])
	if err != nil {
		l.err = err.Error()
		return
	}
	_, err = milestone.Edit(m.store.Root, m.store, m.refSync, f.edit, func(ms *milestone.Milestone) error {
		if err := refuseFiledMilestone(ms); err != nil {
			return err
		}
		all, err := milestone.LoadAll(m.store.Root)
		if err != nil {
			return err
		}
		if err := checkLayerPredecessors(all, f.edit, preds); err != nil {
			return err
		}
		ms.SetGoal(f.vals[0])
		ms.SetAfter(preds)
		return nil
	})
	if err := m.recordOnly(err); err != nil {
		l.err = err.Error()
		return
	}
	m.layerWritten(f.edit, m.formReturn())
}

// createMilestone writes a new milestone placed after the given names, as
// 'jaira milestone create --goal --after' does: under the lock, with a colour
// no other milestone uses, refusing a name that is taken — on the board, or
// filed.
func (m *Model) createMilestone(name, goal string, after []string) error {
	root := m.store.Root
	_, err := milestone.Create(root, m.store, m.refSync, name, func(existing []*milestone.Milestone) (*milestone.Milestone, error) {
		if have, err := milestone.Load(root, name); err == nil {
			if have.Filed() {
				return nil, fmt.Errorf("milestone %q has been filed; 'jaira restore %s.md', run where it was filed, brings it back", name, name)
			}
			return nil, fmt.Errorf("milestone %q already exists — c on a [+] connects it", name)
		} else if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		if where, ok := m.store.FiledMilestone(name); ok {
			return nil, fmt.Errorf("milestone %q has been filed into the logbook (%s) — 'jaira restore %s.md' brings it back", name, where, name)
		}
		if y := m.refSync; y != nil && y.Usable() == nil {
			if content, _, err := y.Repo.ReadMilestone(name); err == nil && milestone.FromBytes(name, content).Filed() {
				return nil, fmt.Errorf("milestone %q has been filed in another clone; 'jaira restore %s.md', run there, brings it back", name, name)
			}
		}
		if err := checkLayerPredecessors(existing, name, after); err != nil {
			return nil, err
		}
		ms := milestone.New(name, milestone.AssignColour(existing, name), time.Now())
		ms.SetGoal(goal)
		ms.SetAfter(after)
		return ms, nil
	})
	return m.recordOnly(err)
}

// parseLayerNames reads the predecessor field: comma-separated names, stored
// the way milestone names are. Empty is the empty list — a swimlane top.
func parseLayerNames(raw string) ([]string, error) {
	out := []string{}
	for _, part := range strings.Split(raw, ",") {
		if strings.TrimSpace(part) == "" {
			continue
		}
		n, _, err := milestone.NormalizeName(part)
		if err != nil {
			return nil, err
		}
		if !slices.Contains(out, n) {
			out = append(out, n)
		}
	}
	return out, nil
}

// checkLayerPredecessors refuses a predecessor no milestone on the board
// carries and an arrow that would close a loop, in the words the CLI's
// checkPredecessors uses — the drawing ignores both, so this is where they
// are stopped.
func checkLayerPredecessors(all []*milestone.Milestone, name string, preds []string) error {
	for _, p := range preds {
		if !slices.ContainsFunc(all, func(ms *milestone.Milestone) bool { return ms.Name == p }) {
			return fmt.Errorf("no milestone %q on this board to come after", p)
		}
		if path := milestone.CyclePath(all, name, p); path != nil {
			return fmt.Errorf("milestone %q cannot come after %q: that closes a cycle, %s",
				name, p, strings.Join(path, " -> "))
		}
	}
	return nil
}

// refuseFiledMilestone stops an edit of a milestone whose file still lies
// here but is marked filed — what a fetch of somebody else's filing leaves.
// Writing it would push the change onto its ref for everyone.
func refuseFiledMilestone(ms *milestone.Milestone) error {
	if ms.Filed() {
		return fmt.Errorf("milestone %q has been filed; 'jaira restore %s.md', run where it was filed, brings it back", ms.Name, ms.Name)
	}
	return nil
}

// recordOnly drops the one error a write may end with although it happened:
// the file was saved and only its ref could not be queued. That is said, but
// not as a failure — the CLI warns in the same case and carries on.
func (m *Model) recordOnly(err error) error {
	var re *milestone.RecordError
	if errors.As(err, &re) {
		m.notify(fmt.Sprintf("Milestone %q was saved here but could not be queued for the remote: %v", re.Name, re.Err), false)
		return nil
	}
	return err
}

// layerWritten reloads after a write and puts the cursor on the milestone it
// was about, so the change is the thing in front of the reader.
func (m *Model) layerWritten(name string, back layerStage) {
	l := m.layer
	l.err = ""
	l.stage = back
	if err := m.reload(); err != nil {
		m.notify(err.Error(), true)
		return
	}
	if back == stageLayer {
		l.sel = layerTarget{kind: onBubble, name: name}
	}
}

// connectable lists what c on a [+] may connect there: every milestone but
// the bubble itself, without one that already comes after the same
// predecessor or that would close a loop. The far [+] starts a swimlane, so it
// offers only the milestones not in the layer yet — a placed one there would
// gain no arrow.
func (m *Model) connectable(t layerTarget) []string {
	after := m.placement(t)
	var out []string
	for _, ms := range m.milestones {
		if ms.Name == t.name {
			continue
		}
		if len(after) == 0 {
			if !ms.Placed {
				out = append(out, ms.Name)
			}
			continue
		}
		p := after[0]
		if ms.Placed && slices.Contains(ms.After, p) {
			continue
		}
		if milestone.WouldCycle(m.milestones, ms.Name, p) {
			continue
		}
		out = append(out, ms.Name)
	}
	return out
}

func (m *Model) keyLayerPick(s string) {
	l := m.layer
	n := len(l.pick)
	switch s {
	case "esc", "q":
		l.stage = stageLayer
	case "down", "j":
		if n > 0 {
			l.pickIdx = (l.pickIdx + 1) % n
		}
	case "up", "k":
		if n > 0 {
			l.pickIdx = (l.pickIdx - 1 + n) % n
		}
	case "enter":
		if l.pickIdx < n {
			m.connectMilestone(l.pick[l.pickIdx], m.placement(l.sel))
		}
	}
}

// connectMilestone puts an existing milestone where the [+] was. One that is
// not in the layer yet is placed there; one that is keeps its place and only
// gains the arrow, because moving a bubble somebody arranged is not what
// "connect" says.
func (m *Model) connectMilestone(name string, after []string) {
	l := m.layer
	_, err := milestone.Edit(m.store.Root, m.store, m.refSync, name, func(ms *milestone.Milestone) error {
		if err := refuseFiledMilestone(ms); err != nil {
			return err
		}
		all, err := milestone.LoadAll(m.store.Root)
		if err != nil {
			return err
		}
		next := after
		if ms.Placed {
			next = append(slices.Clone(ms.After), after...)
		}
		if err := checkLayerPredecessors(all, name, after); err != nil {
			return err
		}
		ms.SetAfter(next)
		return nil
	})
	if err := m.recordOnly(err); err != nil {
		l.err = err.Error()
		return
	}
	m.layerWritten(name, stageLayer)
}

// nextMilestoneColour gives the milestone the next colour of the palette, so
// two neighbours that came out alike can be told apart in a keypress or two.
func (m *Model) nextMilestoneColour(name string) {
	_, err := milestone.Edit(m.store.Root, m.store, m.refSync, name, func(ms *milestone.Milestone) error {
		if err := refuseFiledMilestone(ms); err != nil {
			return err
		}
		i := slices.Index(milestone.Palette, ms.Colour)
		ms.SetColour(milestone.Palette[(i+1)%len(milestone.Palette)])
		return nil
	})
	if err := m.recordOnly(err); err != nil {
		m.notify(err.Error(), true)
		return
	}
	m.layerWritten(name, stageLayer)
}

func (m *Model) keyLayerRemove(s string) {
	l := m.layer
	switch s {
	case "y":
		m.removeFromLayer(l.sel.name)
	case "n", "esc", "q":
		l.stage = stageLayer
	}
}

// removeFromLayer takes a milestone out of the layer — its after line goes,
// its tickets stay — and closes the chain behind it: what came after it now
// comes after what came before it. The successors are worked out under the
// lock the removal holds, then each is written on its own, the way 'jaira
// milestone set' would write them one after another.
func (m *Model) removeFromLayer(name string) {
	l := m.layer
	l.stage = stageLayer
	var reattach map[string][]string
	_, err := milestone.Edit(m.store.Root, m.store, m.refSync, name, func(ms *milestone.Milestone) error {
		if err := refuseFiledMilestone(ms); err != nil {
			return err
		}
		all, err := milestone.LoadAll(m.store.Root)
		if err != nil {
			return err
		}
		reattach = milestone.ReattachSuccessors(all, name)
		ms.Unplace()
		return nil
	})
	if err := m.recordOnly(err); err != nil {
		m.notify(err.Error(), true)
		return
	}
	succ := make([]string, 0, len(reattach))
	for s := range reattach {
		succ = append(succ, s)
	}
	sort.Strings(succ)
	for _, s := range succ {
		_, err := milestone.Edit(m.store.Root, m.store, m.refSync, s, func(ms *milestone.Milestone) error {
			ms.SetAfter(reattach[s])
			return nil
		})
		if err := m.recordOnly(err); err != nil {
			m.notify(fmt.Sprintf("%q left the layer, but %q could not be moved up: %v", name, s, err), true)
			_ = m.reload()
			return
		}
	}
	if err := m.reload(); err != nil {
		m.notify(err.Error(), true)
		return
	}
	l.sel = layerTarget{}
	if len(succ) > 0 {
		l.sel = layerTarget{kind: onBubble, name: succ[0]}
	}
}

// openMilestoneDetail shows one milestone: its goal, where it sits, and its
// tickets.
func (m *Model) openMilestoneDetail(name string) {
	ms := m.milestoneNamed(name)
	if ms == nil {
		return
	}
	l := m.layer
	l.detailName = name
	l.rows = m.memberRows(ms)
	l.rowIdx, l.rowTop = 0, 0
	l.stage = stageDetail
}

// memberRows lists a milestone's tickets where they are now: the lane of one
// on the board, the logbook for a filed one. A ticket filed before the days
// the board loads is read for its title alone, here and only here — the
// layer's counts never read a file.
func (m *Model) memberRows(ms *milestone.Milestone) []memberRow {
	board := make(map[string]*ticket.Ticket, len(m.tickets))
	for _, t := range m.tickets {
		board[t.ID] = t
	}
	for _, lg := range m.logged {
		if _, ok := board[lg.Ticket.ID]; !ok {
			board[lg.Ticket.ID] = lg.Ticket
		}
	}
	var out []memberRow
	for _, id := range ms.Members() {
		r := memberRow{id: id}
		path, filed := m.layer.filed[id]
		t := board[id]
		switch {
		case t != nil && !m.isLogged(t):
			r.title, r.where = t.Title, t.Status
			if ln, ok := m.lanes.Get(t.Status); ok {
				r.where = ln.Name
			}
		case filed && m.loggedFile(path):
			r.where = "logbook"
			if t != nil {
				r.title = t.Title
			} else if h, err := ticket.LoadHeader(path); err == nil {
				r.title = h.Title
			}
		case filed:
			r.where = "archived"
		default:
			r.where = "not on this board"
		}
		out = append(out, r)
	}
	return out
}

func (m *Model) keyMilestoneDetail(s string) {
	l := m.layer
	n := len(l.rows)
	switch s {
	case "esc", "q":
		l.stage = stageLayer
		l.detailName = ""
	case "down", "j":
		if l.rowIdx < n-1 {
			l.rowIdx++
		}
	case "up", "k":
		if l.rowIdx > 0 {
			l.rowIdx--
		}
	case "e":
		m.startLayerEdit(l.detailName)
	case "b":
		// The board narrowed to this round of work, through the one filter
		// M and / write as well, so esc on the board takes it off again.
		name := l.detailName
		m.closeLayer()
		m.setFilterTerm("milestone", name)
	case "enter":
		if l.rowIdx < n {
			m.openMember(l.rows[l.rowIdx])
		}
	}
}

// openMember opens one of the milestone's tickets the way the board opens a
// card — by putting the cursor on it — so every key of the open ticket acts
// on the ticket that is shown. esc there comes back to the milestone.
func (m *Model) openMember(r memberRow) {
	if m.selectByID(r.id) {
		m.openDetail()
		return
	}
	h := ticket.Handle(r.id)
	switch r.where {
	case "logbook":
		m.notify(fmt.Sprintf("%s was filed into the logbook before the days the board shows.\n\n'jaira show %s' reads it.", h, h), false)
	case "archived", "not on this board":
		m.notify(fmt.Sprintf("%s is %s, so there is no card to open.", h, r.where), false)
	default:
		m.notify(fmt.Sprintf("%s is on the board but the filter %s hides it.\n\nClear the filter with esc on the board, then try again.", h, m.filter), false)
	}
}
