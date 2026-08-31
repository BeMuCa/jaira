package tui

import (
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/lane"
)

// boardEditorDoneMsg reports that $EDITOR, opened on a lane's own file from
// the default board screen, has exited.
type boardEditorDoneMsg struct{ err error }

// defaultBoardScreen is the per-user default board screen, reached from the
// home screen with 'd' — the only per-user surface there is. Following
// browse.go's shape: its own state, a key method, a render method.
//
// It selects lanes and pre-ticks options; it does not reorder them and holds
// no form for a lane's prompt, tier or contract. A reorder control here would
// be a second way to position a lane on top of precedence, which is exactly
// the duplicate mechanism the ordering model exists to avoid, and a second
// form for a lane's contract would be a second way for it to disagree with
// the file. 'e' opens the file instead — full editing power without a
// second write path.
type defaultBoardScreen struct {
	set   *lane.Set
	board *lane.DefaultBoard

	lanes   map[string]bool // lane id -> selected
	options map[string]bool // option name -> pre-ticked

	idx   int // cursor within whichever list has focus
	focus int // 0 = lanes, 1 = options

	msg   string
	isErr bool
}

func newDefaultBoardScreen(set *lane.Set, board *lane.DefaultBoard) *defaultBoardScreen {
	d := &defaultBoardScreen{set: set, board: board, lanes: map[string]bool{}, options: map[string]bool{}}

	// An absent or empty selection means the built-ins (see DefaultBoard's own
	// doc comment), so a first visit to this screen must show every built-in
	// already ticked — otherwise saving straight away would narrow the board
	// down to nothing instead of leaving it unchanged.
	if len(board.Lanes) == 0 {
		for _, l := range set.Lanes {
			d.lanes[l.ID] = l.Builtin
		}
	} else {
		selected := make(map[string]bool, len(board.Lanes))
		for _, id := range board.Lanes {
			selected[id] = true
		}
		for _, l := range set.Lanes {
			d.lanes[l.ID] = selected[l.ID]
		}
	}

	optedIn := make(map[string]bool, len(board.Options))
	for _, o := range board.Options {
		optedIn[o] = true
	}
	for _, o := range set.Options() {
		d.options[o] = optedIn[o]
	}
	return d
}

func (d *defaultBoardScreen) laneCount() int   { return len(d.set.Lanes) }
func (d *defaultBoardScreen) optionCount() int { return len(d.set.Options()) }

// key drives the screen. It reports whether the screen is finished (esc/q)
// and a command to run, non-nil only when 'e' launches $EDITOR.
func (d *defaultBoardScreen) key(s string) (done bool, cmd tea.Cmd) {
	d.msg = ""
	switch s {
	case "esc", "q":
		return true, nil
	case "tab":
		d.focus = 1 - d.focus
		d.idx = 0
	case "j", "down":
		d.move(1)
	case "k", "up":
		d.move(-1)
	case " ", "space":
		d.toggle()
	case "e":
		return false, d.openSelectedInEditor()
	case "s":
		d.save()
	}
	return false, nil
}

func (d *defaultBoardScreen) move(delta int) {
	n := d.laneCount()
	if d.focus == 1 {
		n = d.optionCount()
	}
	if n == 0 {
		d.idx = 0
		return
	}
	d.idx += delta
	if d.idx < 0 {
		d.idx = 0
	}
	if d.idx >= n {
		d.idx = n - 1
	}
}

func (d *defaultBoardScreen) toggle() {
	if d.focus == 0 {
		if d.idx < len(d.set.Lanes) {
			id := d.set.Lanes[d.idx].ID
			d.lanes[id] = !d.lanes[id]
		}
		return
	}
	opts := d.set.Options()
	if d.idx < len(opts) {
		d.options[opts[d.idx]] = !d.options[opts[d.idx]]
	}
}

// openSelectedInEditor opens the selected lane's own file, not a temp copy —
// unlike a ticket body, a lane file is the thing itself, and this screen
// deliberately has no second form to disagree with it.
func (d *defaultBoardScreen) openSelectedInEditor() tea.Cmd {
	if d.focus != 0 || d.idx >= len(d.set.Lanes) {
		return nil
	}
	l := d.set.Lanes[d.idx]
	if l.Builtin {
		d.msg, d.isErr = "built-in lanes have no file to edit; use it here from the lane settings screen first", true
		return nil
	}
	argv := append(editorCommand(), l.Source)
	cmd := exec.Command(argv[0], argv[1:]...)
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return boardEditorDoneMsg{err: err}
	})
}

// save writes the current selection to disk via lane.SaveDefaultBoard.
func (d *defaultBoardScreen) save() {
	var lanes []string
	for _, l := range d.set.Lanes {
		if d.lanes[l.ID] {
			lanes = append(lanes, l.ID)
		}
	}
	var opts []string
	for _, o := range d.set.Options() {
		if d.options[o] {
			opts = append(opts, o)
		}
	}
	d.board.Lanes, d.board.Options = lanes, opts
	if err := lane.SaveDefaultBoard(d.board); err != nil {
		d.msg, d.isErr = err.Error(), true
		return
	}
	d.msg, d.isErr = "saved "+d.board.Path, false
}

func (d *defaultBoardScreen) render(width, height int) string {
	w := max(20, width)
	var sb strings.Builder
	sb.WriteString(styLaneTitle.Render("Default board") + "\n")
	sb.WriteString(styMeta.Render(wrap(d.board.Path, w, 0)) + "\n")
	sb.WriteString(styBar.Render(strings.Repeat("─", w)) + "\n")

	sb.WriteString(styLaneTitle.Render("Lanes") + "\n")
	for i, l := range d.set.Lanes {
		lead := "  "
		name := l.Name
		if d.focus == 0 && i == d.idx {
			lead = stySelected.Render("▌ ")
			name = stySelected.Render(name)
		}
		mark := "[ ]"
		if d.lanes[l.ID] {
			mark = "[x]"
		}
		sb.WriteString(truncate(lead+mark+" "+name, w) + "\n")
	}

	if opts := d.set.Options(); len(opts) > 0 {
		sb.WriteString(styBar.Render(strings.Repeat("─", w)) + "\n")
		sb.WriteString(styLaneTitle.Render("Options") + "\n")
		for i, o := range opts {
			lead := "  "
			name := o
			if d.focus == 1 && i == d.idx {
				lead = stySelected.Render("▌ ")
				name = stySelected.Render(name)
			}
			mark := "[ ]"
			if d.options[o] {
				mark = "[x]"
			}
			sb.WriteString(truncate(lead+mark+" "+name, w) + "\n")
		}
	}

	if d.msg != "" {
		style := styOK
		if d.isErr {
			style = styErr
		}
		sb.WriteString("\n" + style.Render(wrap(d.msg, w, 0)) + "\n")
	}
	for _, l := range wrapHints([]string{"tab switch list", "space toggle", "e edit file", "s save", "esc back (selects lanes, does not reorder them)"}, max(1, w)) {
		sb.WriteString("\n" + styMeta.Render(l))
	}
	return sb.String()
}
