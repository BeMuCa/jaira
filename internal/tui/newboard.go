package tui

import (
	"fmt"
	"strings"

	"github.com/BeMuCa/jaira/core/lane"
)

// boardForm is the new-board form the lane settings screen opens with B:
// a name, a description and a checklist of lanes — today's board's and the
// catalogue's (lane.PickableLanes), today's ticked. Saving makes the same
// calls 'jaira boards create --lanes' does: lane.CreateBoardWithLanes, then
// lane.SaveToCatalogue.
type boardForm struct {
	field  int // 0 name, 1 description, 2 lanes
	name   string
	desc   string
	lanes  []*lane.Lane
	picked map[string]bool
	idx    int
	err    string
}

const (
	formName = iota
	formDesc
	formLanes
	formFields
)

func newBoardForm(set *lane.Set) (*boardForm, error) {
	lanes, err := lane.PickableLanes(set)
	if err != nil {
		return nil, err
	}
	f := &boardForm{lanes: lanes, picked: map[string]bool{}}
	for _, l := range set.Lanes {
		f.picked[l.ID] = true
	}
	return f, nil
}

// openBoardForm is B on the lane settings screen.
func (ls *laneScreen) openBoardForm() {
	f, err := newBoardForm(ls.set)
	if err != nil {
		ls.msg, ls.isErr = err.Error(), true
		return
	}
	ls.form = f
}

// formKey drives the form: s is the command name of the key, text what it
// typed. tab and shift+tab change field; enter goes on to the next field
// and, on the lanes, creates the board; esc leaves without writing.
func (ls *laneScreen) formKey(s, text string) {
	f := ls.form
	f.err = ""
	switch s {
	case "esc":
		ls.form = nil
		return
	case "tab":
		f.field = (f.field + 1) % formFields
		return
	case "shift+tab":
		f.field = (f.field + formFields - 1) % formFields
		return
	}
	if f.field == formLanes {
		switch s {
		case "j", "down":
			if f.idx < len(f.lanes)-1 {
				f.idx++
			}
		case "k", "up":
			if f.idx > 0 {
				f.idx--
			}
		case " ", "space", "x":
			if f.idx < len(f.lanes) {
				id := f.lanes[f.idx].ID
				f.picked[id] = !f.picked[id]
			}
		case "enter":
			ls.createBoard()
		}
		return
	}
	switch s {
	case "enter":
		f.field++
	case "backspace":
		buf := f.buffer()
		if r := []rune(*buf); len(r) > 0 {
			*buf = string(r[:len(r)-1])
		}
	default:
		f.insert(text)
	}
}

// buffer is the text field with focus.
func (f *boardForm) buffer() *string {
	if f.field == formDesc {
		return &f.desc
	}
	return &f.name
}

// insert types or pastes text into the text field with focus.
func (f *boardForm) insert(text string) {
	if f.field == formLanes {
		return
	}
	*f.buffer() += foldToOneLine(text)
}

// createBoard writes the board the form describes and saves it into the
// catalogue, or leaves the form open with what stopped it.
func (ls *laneScreen) createBoard() {
	f := ls.form
	name, err := lane.NormalizeBoard(f.name)
	switch {
	case err != nil:
		f.err = err.Error()
		return
	case name == "":
		f.err = "a board needs a name of its own; default is today's board"
		return
	case strings.TrimSpace(f.desc) == "":
		f.err = "a board needs a description: which tickets belong on it"
		return
	}
	var ids []string
	for _, l := range f.lanes {
		if f.picked[l.ID] {
			ids = append(ids, l.ID)
		}
	}
	if _, err := lane.CreateBoardWithLanes(ls.store.Root, name, ls.set, ids, f.desc); err != nil {
		f.err = err.Error()
		return
	}
	ls.form = nil
	ls.msg, ls.isErr = fmt.Sprintf("board %s created with %d lanes", name, len(ids)), false
	path, saved, err := lane.SaveToCatalogue(ls.store.Root, name)
	switch {
	case err != nil:
		ls.msg, ls.isErr = ls.msg+"; not saved to your catalogue: "+err.Error(), true
	case saved:
		ls.msg += "; saved to your catalogue, " + path
	}
}

// renderBoardForm draws the form in place of the screen's lower half.
func (ls *laneScreen) renderBoardForm(w int) string {
	f := ls.form
	var sb strings.Builder
	sb.WriteString(styLaneTitle.Render("New board") + "\n")
	field := func(i int, label, val string) {
		lead := "  "
		if f.field == i {
			lead = stySelected.Render("▌ ")
			val += "▏"
		}
		sb.WriteString(truncate(lead+label+": "+val, w) + "\n")
	}
	field(formName, "name", f.name)
	field(formDesc, "description", f.desc)
	sb.WriteString(styMeta.Render("  lanes, in this order:") + "\n")
	for i, l := range f.lanes {
		box := "[ ]"
		if f.picked[l.ID] {
			box = "[x]"
		}
		row := fmt.Sprintf("    %s %s (%s)", box, l.Name, ls.sourceLabel(l))
		if f.field == formLanes && i == f.idx {
			row = stySelected.Render(row)
		}
		sb.WriteString(truncate(row, w) + "\n")
	}
	if f.err != "" {
		sb.WriteString("\n" + styErr.Render(truncate(f.err, w)) + "\n")
	}
	for _, l := range wrapHints([]string{"tab next field", "space tick", "enter create", "esc cancel"}, max(1, w)) {
		sb.WriteString("\n" + styMeta.Render(l))
	}
	return sb.String()
}
