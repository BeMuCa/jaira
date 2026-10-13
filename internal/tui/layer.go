package tui

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/link"
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
	id   string // the bubble's milestone ID, or the bubble a [+] belongs to; "" for onFar
}

// layerStage is what the layer is doing on top of the picture.
type layerStage int

const (
	stageLayer  layerStage = iota
	stageChoose            // enter on a [+]: a new milestone or an existing one
	stageForm              // name and goal, or name, goal and predecessors
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

	stage layerStage
	form  layerForm

	pick    []string // the names the connect picker offers
	pickIdx int

	// The milestone detail: whose it is, its ticket rows, the cursor and the
	// first row on screen.
	detailID string
	rows     []memberRow
	rowIdx   int
	rowTop   int

	err string // why the last choice was refused, shown where it was made
}

// layerForm is the form behind a new milestone (name, goal) and behind e
// (name, goal, predecessors).
type layerForm struct {
	edit   string // the name of the milestone e is editing; "" for a new one
	editID string // its ID
	labels []string
	vals   []string
	field  int
}

// memberRow is one ticket in the milestone detail.
type memberRow struct {
	id, title string
	place     link.Place
	where     string // the lane's name, "logbook", or why it is nowhere
}

// openLayer shows the layer, the cursor on the first bubble — or on the one
// [+] an empty layer has.
func (m *Model) openLayer() {
	m.layer = &layerView{}
	m.mode = modeLayer
	m.layerScene() // settles the cursor on something that exists
}

// closeLayer goes back to the board.
func (m *Model) closeLayer() {
	m.layer = nil
	m.mode = modeBoard
}

// refreshLayer is reload's part for an open layer: an open detail lists what
// its milestone now holds.
func (m *Model) refreshLayer() {
	l := m.layer
	if l == nil {
		return
	}
	if l.stage == stageDetail {
		ms := m.milestoneByID(l.detailID)
		if ms == nil {
			// Taken off the board elsewhere: there is nothing left to list.
			l.stage = stageLayer
			return
		}
		l.rows = m.memberRows(ms)
		l.rowIdx = min(l.rowIdx, max(0, len(l.rows)-1))
	}
}

// spellAsID is an after: token as the ID it names, or as written when it names
// no milestone.
func (m *Model) spellAsID(token string) string {
	if id := milestone.Resolve(m.milestones, token); id != "" {
		return id
	}
	return token
}

// milestoneLabel is what a milestone is called on screen.
func (m *Model) milestoneLabel(id string) string {
	if ms := m.milestoneByID(id); ms != nil {
		return ms.Name
	}
	return id
}

func (m *Model) milestoneByID(id string) *milestone.Milestone {
	for _, ms := range m.milestones {
		if ms.ID == id {
			return ms
		}
	}
	return nil
}

// milestoneDone counts a milestone's tickets and how many of them are done,
// asking the relation index the link window and the gate ask: a ticket in a
// terminal lane, one filed into the logbook on any day, or one archived from
// a terminal lane. The board loads only the last days of the logbook; the
// index sees all of it.
func (m *Model) milestoneDone(ms *milestone.Milestone) (done, total int) {
	ix := m.linkIndex()
	members := ms.Members()
	for _, id := range members {
		if ix.Satisfied(id) {
			done++
		}
	}
	return done, len(members)
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
			m.openMilestoneDetail(l.sel.id)
		} else {
			l.stage = stageChoose
		}
	case "e":
		if l.sel.kind == onBubble {
			m.startLayerEdit(l.sel.id)
		}
	case "c":
		if l.sel.kind == onBubble {
			m.nextMilestoneColour(l.sel.id)
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
		if ms := m.milestoneByID(t.id); ms != nil && len(ms.After) > 0 {
			return []string{m.spellAsID(ms.After[0])}
		}
	case onBelow:
		return []string{t.id}
	}
	return []string{}
}

// placementText says in words what a [+] would create, for the header and
// the dialogs behind it.
func (m *Model) placementText(t layerTarget) string {
	if after := m.placement(t); len(after) > 0 {
		return "after " + m.milestoneLabel(after[0])
	}
	return "a swimlane of its own"
}

func (m *Model) keyLayerChoose(s string) {
	l := m.layer
	switch s {
	case "esc", "q":
		l.stage = stageLayer
	case "n":
		l.form = layerForm{labels: []string{"name", "goal"}, vals: make([]string, 2)}
		l.err = ""
		l.stage = stageForm
	case "c":
		l.pick = m.connectable(l.sel)
		l.pickIdx = 0
		l.err = ""
		l.stage = stagePick
	}
}

// startLayerEdit opens the form behind e on a milestone: its name, its goal
// and its predecessors as they stand.
func (m *Model) startLayerEdit(id string) {
	ms := m.milestoneByID(id)
	if ms == nil {
		return
	}
	l := m.layer
	l.form = layerForm{
		edit:   ms.Name,
		editID: ms.ID,
		labels: []string{"name", "goal", "after"},
		vals:   []string{ms.Name, ms.Goal, strings.Join(m.afterLabels(ms.After), ", ")},
	}
	l.err = ""
	l.stage = stageForm
}

// afterLabels is a milestone's predecessors as the form shows them: by name,
// whatever after: holds.
func (m *Model) afterLabels(tokens []string) []string {
	out := make([]string, len(tokens))
	for i, tok := range tokens {
		out[i] = tok
		if id := milestone.Resolve(m.milestones, tok); id != "" {
			out[i] = m.milestoneLabel(id)
		}
	}
	return out
}

func (m *Model) keyLayerForm(s, text string) {
	l := m.layer
	f := &l.form
	switch s {
	case "esc":
		l.stage = m.formReturn()
	case "tab", "down":
		f.field = (f.field + 1) % len(f.vals)
	case "shift+tab", "up":
		f.field = (f.field + len(f.vals) - 1) % len(f.vals)
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
	if m.layer.detailID != "" && m.layer.form.editID == m.layer.detailID {
		return stageDetail
	}
	return stageLayer
}

func (m *Model) submitLayerForm() {
	l := m.layer
	f := l.form
	// A name typed here is filed the way 'jaira milestone create' files it. A
	// name left as it stood is the file's own, kept as it is: normalizing it
	// would rename a hand-named Q4.md to q4 on every save of its goal.
	name := strings.TrimSpace(f.vals[0])
	if f.edit == "" || name != f.edit {
		n, _, err := milestone.NormalizeName(name)
		if name == "" {
			err = errors.New("a milestone needs a name")
		}
		if err != nil {
			l.err = err.Error()
			return
		}
		name = n
	}
	if f.edit == "" {
		created, err := m.createMilestone(name, f.vals[1], m.placement(l.sel))
		if err != nil {
			l.err = err.Error()
			return
		}
		m.layerWritten(created.ID, stageLayer)
		return
	}

	// Goal and predecessors are checked before the name changes, so a form
	// these checks refuse has written nothing; the rename goes the way 'jaira
	// milestone rename' goes, and the rest is written under the new name.
	preds, err := milestone.ParseNames(f.vals[2])
	if err != nil {
		l.err = err.Error()
		return
	}
	all, err := milestone.LoadAll(m.store.Root)
	if err == nil {
		var ids []string
		if ids, err = milestone.ResolveRefs(all, preds); err == nil {
			err = milestone.CheckAfter(all, f.editID, ids)
		}
	}
	if err != nil {
		l.err = err.Error()
		return
	}
	id := f.editID
	if name != f.edit {
		renamed, err := m.renameMilestone(f.editID, name)
		if err != nil {
			l.err = err.Error()
			return
		}
		id = renamed.ID
	}
	_, err = milestone.Edit(m.store.Root, m.store, m.refSync, id, func(ms *milestone.Milestone) error {
		all, err := milestone.LoadAll(m.store.Root)
		if err != nil {
			return err
		}
		ids, err := milestone.ResolveRefs(all, preds)
		if err == nil {
			err = milestone.CheckAfter(all, id, ids)
		}
		if err != nil {
			return err
		}
		ms.SetGoal(f.vals[1])
		ms.SetAfter(ids)
		return nil
	})
	if err := m.recordOnly(err); err != nil {
		l.err = err.Error()
		if name != f.edit {
			// The rename went through and only the rest was refused: the
			// picture behind the open form shows the new name.
			if err := m.reload(); err != nil {
				m.notify(err.Error(), true)
			}
		}
		return
	}
	m.layerWritten(id, m.formReturn())
}

// renameMilestone renames through milestone.Rename, the path 'jaira milestone
// rename' takes, refusing a new name that is taken the way createMilestone
// refuses it. From then on the open form and detail are about the new name,
// and a board narrowed to the old name stays narrowed to the milestone.
func (m *Model) renameMilestone(fromID, to string) (*milestone.Milestone, error) {
	from := m.milestoneLabel(fromID)
	renamed, err := milestone.Rename(m.store.Root, m.store, m.refSync, fromID, to)
	if err := m.recordOnly(err); err != nil {
		return nil, err
	}
	l := m.layer
	if l.detailID == fromID {
		l.detailID = renamed.ID
	}
	if l.sel.id == fromID {
		l.sel.id = renamed.ID
	}
	l.form.edit, l.form.editID = to, renamed.ID
	_, names := splitTerm(m.filter, "milestone")
	if i := slices.Index(names, from); i >= 0 {
		names[i] = to
		m.setFilterTerm("milestone", names...)
	}
	return renamed, nil
}

// createMilestone writes a new milestone placed after the given names, as
// 'jaira milestone create --goal --after' does: under the lock, with a colour
// no other milestone uses, refusing a name a milestone on the board has.
func (m *Model) createMilestone(name, goal string, after []string) (*milestone.Milestone, error) {
	root := m.store.Root
	ms, err := milestone.Create(root, m.store, m.refSync, func(existing []*milestone.Milestone) (*milestone.Milestone, error) {
		var ee *milestone.ExistsError
		if err := milestone.CheckFree(root, name); errors.As(err, &ee) {
			return nil, fmt.Errorf("milestone %q already exists — c on a [+] connects it", name)
		} else if err != nil {
			return nil, err
		}
		ms := milestone.New(name, milestone.AssignColour(existing, name), time.Now())
		ids, err := milestone.ResolveRefs(existing, after)
		if err == nil {
			err = milestone.CheckAfter(existing, ms.ID, ids)
		}
		if err != nil {
			return nil, err
		}
		ms.SetGoal(goal)
		ms.SetAfter(ids)
		return ms, nil
	})
	return ms, m.recordOnly(err)
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
func (m *Model) layerWritten(id string, back layerStage) {
	l := m.layer
	l.err = ""
	l.stage = back
	if err := m.reload(); err != nil {
		m.notify(err.Error(), true)
		return
	}
	if back == stageLayer {
		l.sel = layerTarget{kind: onBubble, id: id}
	}
}

// connectable lists what c on a [+] may connect there: every milestone but
// the bubble itself, without one that already comes after the same
// predecessor or that the write's own check (CheckAfter) would refuse — a
// loop. The far [+] starts a swimlane, so it offers only the milestones not in
// the layer yet — a placed one there would gain no arrow.
func (m *Model) connectable(t layerTarget) []string {
	after := m.placement(t)
	var out []string
	for _, ms := range m.milestones {
		if ms.ID == t.id {
			continue
		}
		if len(after) == 0 {
			if !ms.Placed {
				out = append(out, ms.ID)
			}
			continue
		}
		p := after[0]
		pid := milestone.Resolve(m.milestones, p)
		if ms.Placed && slices.ContainsFunc(ms.After, func(tok string) bool { return milestone.Resolve(m.milestones, tok) == pid }) {
			continue
		}
		if milestone.CheckAfter(m.milestones, ms.ID, []string{p}) != nil {
			continue
		}
		out = append(out, ms.ID)
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
func (m *Model) connectMilestone(id string, after []string) {
	l := m.layer
	_, err := milestone.Edit(m.store.Root, m.store, m.refSync, id, func(ms *milestone.Milestone) error {
		all, err := milestone.LoadAll(m.store.Root)
		if err != nil {
			return err
		}
		ids, err := milestone.ResolveRefs(all, after)
		if err == nil {
			err = milestone.CheckAfter(all, id, ids)
		}
		if err != nil {
			return err
		}
		next := ids
		if ms.Placed {
			next = append(slices.Clone(ms.After), ids...)
		}
		ms.SetAfter(next)
		return nil
	})
	if err := m.recordOnly(err); err != nil {
		l.err = err.Error()
		return
	}
	m.layerWritten(id, stageLayer)
}

// nextMilestoneColour gives the milestone the next colour of the palette that
// no other milestone wears, so two neighbours that came out alike can be told
// apart in a keypress or two. Once every colour is taken it is simply the next.
func (m *Model) nextMilestoneColour(id string) {
	used := map[int]bool{}
	for _, other := range m.milestones {
		if other.ID != id {
			used[other.Colour] = true
		}
	}
	_, err := milestone.Edit(m.store.Root, m.store, m.refSync, id, func(ms *milestone.Milestone) error {
		n := len(milestone.Palette)
		i := slices.Index(milestone.Palette, ms.Colour)
		next := milestone.Palette[(i+1)%n]
		for k := 1; k < n; k++ {
			if c := milestone.Palette[(i+k)%n]; !used[c] {
				next = c
				break
			}
		}
		ms.SetColour(next)
		return nil
	})
	if err := m.recordOnly(err); err != nil {
		m.notify(err.Error(), true)
		return
	}
	m.layerWritten(id, stageLayer)
}

func (m *Model) keyLayerRemove(s string) {
	l := m.layer
	switch s {
	case "y":
		m.removeFromLayer(l.sel.id)
	case "n", "esc", "q":
		l.stage = stageLayer
	}
}

// removeFromLayer takes a milestone out of the layer — its after line goes,
// its tickets stay — and milestone.Edit closes the chain behind it under the
// same lock, exactly as 'jaira milestone set --no-roadmap' does: what came after
// it now comes after what came before it. The cursor goes to the first of
// those, so the change is the thing in front of the reader.
func (m *Model) removeFromLayer(id string) {
	m.layer.stage = stageLayer
	moved := slices.SortedFunc(maps.Keys(milestone.ReattachSuccessors(m.milestones, id)), func(a, b string) int {
		if c := strings.Compare(m.milestoneLabel(a), m.milestoneLabel(b)); c != 0 {
			return c
		}
		return strings.Compare(a, b)
	})
	_, err := milestone.Edit(m.store.Root, m.store, m.refSync, id, func(ms *milestone.Milestone) error {
		ms.Unplace()
		return nil
	})
	if err := m.recordOnly(err); err != nil {
		m.notify(err.Error(), true)
		return
	}
	first := ""
	if len(moved) > 0 {
		first = moved[0]
	}
	m.layerWritten(first, stageLayer)
}

// openMilestoneDetail shows one milestone: its goal, where it sits, and its
// tickets.
func (m *Model) openMilestoneDetail(id string) {
	ms := m.milestoneByID(id)
	if ms == nil {
		return
	}
	l := m.layer
	l.detailID = id
	l.rows = m.memberRows(ms)
	l.rowIdx, l.rowTop = 0, 0
	l.stage = stageDetail
}

// memberRows lists a milestone's tickets where they are now — the lane of
// one on the board, the logbook, the archive — as the relation index finds
// them, which reads a filed ticket's header for its title and nothing more.
func (m *Model) memberRows(ms *milestone.Milestone) []memberRow {
	ix := m.linkIndex()
	var out []memberRow
	for _, id := range ms.Members() {
		ref, _ := ix.Lookup(id)
		r := memberRow{id: id, title: ref.Title, place: ref.Place}
		switch ref.Place {
		case link.PlaceBoard, link.PlaceRef:
			r.where = ref.Lane
			if ln, ok := m.lanes.Get(ref.Lane); ok {
				r.where = ln.Name
			}
		case link.PlaceLogbook:
			r.where = "logbook"
		case link.PlaceArchive:
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
		l.detailID = ""
	case "down", "j":
		if l.rowIdx < n-1 {
			l.rowIdx++
		}
	case "up", "k":
		if l.rowIdx > 0 {
			l.rowIdx--
		}
	case "e":
		m.startLayerEdit(l.detailID)
	case "b":
		// The board narrowed to this round of work, through the one filter
		// M and / write as well, so esc on the board takes it off again.
		name := m.milestoneLabel(l.detailID)
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
	switch r.place {
	case link.PlaceLogbook:
		m.notify(fmt.Sprintf("%s was filed into the logbook before the days the board shows.\n\n'jaira show %s' reads it.", h, h), false)
	case link.PlaceArchive, link.PlaceUnknown:
		m.notify(fmt.Sprintf("%s is %s, so there is no card to open.", h, r.where), false)
	default:
		m.notify(fmt.Sprintf("%s is on the board but the filter %s hides it.\n\nClear the filter with esc on the board, then try again.", h, m.filter), false)
	}
}
