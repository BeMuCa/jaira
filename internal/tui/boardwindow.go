package tui

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/BeMuCa/jaira/core/gitrepo"
	"github.com/BeMuCa/jaira/core/lane"
)

// boardRow is one board in the window b opens: its name ("" is today's), who
// created it, its lanes in order, and whether it is active.
type boardRow struct {
	name    string
	creator string
	lanes   []string
	active  bool
}

// boardsWindow is the open board window: its rows in the order they are
// shown, the cursor, and what went wrong with the last space.
type boardsWindow struct {
	rows []boardRow
	idx  int
	err  string
}

// boardsSeenFile records, per machine and per repository, when each board was
// last shown. It orders the inactive boards in the window: which board a
// person last looked at is their own navigation, not something the team
// shares, so it lives beside the other per-repository state rather than in
// .jaira/.
const boardsSeenFile = "boards-seen.json"

func (m *Model) boardsSeenPath() string {
	return filepath.Join(m.store.RepoStateDir(), boardsSeenFile)
}

// boardsSeen reads when each board was last shown; a missing or unreadable
// file is no history at all.
func (m *Model) boardsSeen() map[string]time.Time {
	seen := map[string]time.Time{}
	b, err := os.ReadFile(m.boardsSeenPath())
	if err == nil {
		_ = json.Unmarshal(b, &seen)
	}
	return seen
}

// markBoardsSeen stamps the boards named with now. Best-effort: a board that
// cannot write its history still switches.
func (m *Model) markBoardsSeen(names ...string) {
	seen := m.boardsSeen()
	now := time.Now().UTC()
	for _, n := range names {
		seen[n] = now
	}
	b, err := json.Marshal(seen)
	if err != nil {
		return
	}
	path := m.boardsSeenPath()
	if os.MkdirAll(filepath.Dir(path), 0o755) == nil {
		_ = os.WriteFile(path, b, 0o644)
	}
}

// showBoard switches the board on screen to name, stamping the board left and
// the board shown as just used.
func (m *Model) showBoard(name string) {
	if name == m.boardName {
		return
	}
	m.markBoardsSeen(m.boardName, name)
	m.boardName = name
	m.laneIdx, m.cardIdx = 0, 0
	if err := m.reload(); err != nil {
		m.notify(err.Error(), true)
	}
}

// stepBoard shows the next active board, dir 1 down and -1 up, in the order
// the board window lists them, wrapping around. On a board that is not active
// it goes to the nearest active one in that direction.
func (m *Model) stepBoard(dir int) {
	root := m.store.Root
	all := lane.AllBoards(root)
	active := map[string]bool{}
	for _, b := range lane.ActiveBoards(root) {
		active[b] = true
	}
	at := 0
	for i, b := range all {
		if b == m.boardName {
			at = i
		}
	}
	n := len(all)
	for step := 1; step < n; step++ {
		if b := all[((at+dir*step)%n+n)%n]; active[b] {
			m.showBoard(b)
			return
		}
	}
	if len(all) == 1 {
		m.notify("This repository has one board.\n\nAnother is a directory of lane files: .jaira/boards/<name>/lanes/.", false)
		return
	}
	m.notify("No other board is active.\n\nb lists every board; space there makes one active.", false)
}

// openBoards opens the board window: active boards first in the fixed order
// shift+up/down steps through, then the others most recently shown first,
// then the rest by name. The cursor starts on the board on screen.
func (m *Model) openBoards() {
	root := m.store.Root
	active := map[string]bool{}
	for _, b := range lane.ActiveBoards(root) {
		active[b] = true
	}
	seen := m.boardsSeen()
	repo := &gitrepo.Repo{Dir: root}
	var rows []boardRow
	for _, name := range lane.AllBoards(root) {
		r := boardRow{name: name, active: active[name]}
		dir := filepath.Dir(lane.DescriptionPath(root, name))
		if rel, err := filepath.Rel(root, dir); err == nil {
			r.creator = repo.FirstAuthor(rel)
		}
		if set, err := lane.LoadBoard(root, name); err == nil {
			for _, l := range set.Lanes {
				r.lanes = append(r.lanes, l.Name)
			}
		}
		rows = append(rows, r)
	}
	// Stable, so the active ones and the never-shown ones keep AllBoards'
	// order: today's board first, then by name.
	sort.SliceStable(rows, func(i, j int) bool {
		a, b := rows[i], rows[j]
		if a.active != b.active {
			return a.active
		}
		if a.active {
			return false
		}
		return seen[a.name].After(seen[b.name])
	})
	m.boards = boardsWindow{rows: rows}
	for i, r := range rows {
		if r.name == m.boardName {
			m.boards.idx = i
		}
	}
	m.returnTo = m.mode
	m.mode = modeBoards
}

// keyBoards drives the open board window.
func (m *Model) keyBoards(s string) {
	w := &m.boards
	n := len(w.rows)
	switch s {
	case "esc", "q", "b":
		m.mode = m.returnTo
	case "j", "down":
		if n > 0 {
			w.idx = (w.idx + 1) % n
		}
	case "k", "up":
		if n > 0 {
			w.idx = (w.idx - 1 + n) % n
		}
	case "space":
		if w.idx < 0 || w.idx >= n {
			return
		}
		r := &w.rows[w.idx]
		err := lane.SetBoardActive(m.store.Root, r.name, !r.active)
		switch {
		case errors.Is(err, lane.ErrLastActiveBoard):
			w.err = "The last active board stays active: the teamlead needs one to put tickets on."
		case err != nil:
			w.err = err.Error()
		default:
			w.err = ""
			r.active = !r.active
		}
	case "enter":
		m.mode = m.returnTo
		if w.idx >= 0 && w.idx < n {
			m.showBoard(w.rows[w.idx].name)
		}
	}
}

// renderBoards draws the board window: each board under its name, with who
// created it and its lanes below, a tick on the active ones.
func (m *Model) renderBoards() string {
	w := m.boards
	width := max(20, m.width-14)
	var b strings.Builder
	b.WriteString(styLaneTitle.Render("Boards") + "\n")
	b.WriteString(styBar.Render(strings.Repeat("─", min(width, 60))) + "\n")

	blocks := make([]string, len(w.rows))
	for i, r := range w.rows {
		name := r.name
		if name == "" {
			name = lane.DefaultBoardName
		}
		tick := styMeta.Render("○")
		if r.active {
			tick = styOK.Render("✓")
		}
		marker := "  "
		head := styLaneTitle.Render(name)
		if i == w.idx {
			marker, head = stySelected.Render("▌ "), stySelected.Render(name)
		}
		cur := ""
		if r.name == m.boardName {
			cur = styMeta.Render("  (shown)")
		}
		creator := "created by " + r.creator
		if r.creator == "" {
			creator = "creator unknown: not committed yet"
		}
		blocks[i] = "\n" + marker + tick + " " + head + cur + "\n" +
			"    " + styMeta.Render(truncate(creator, width-4)) + "\n" +
			"    " + styleLines(styMeta, wrap(strings.Join(r.lanes, " → "), width-4, 4)) + "\n"
	}
	// Only the boards around the cursor that fit, so the cursor never falls
	// below the window's edge; the title, the hints and the error line take
	// the rest.
	budget := max(4, m.modalContent()-6)
	height := func(s string) int { return strings.Count(s, "\n") }
	start, used := 0, 0
	for i := 0; i <= w.idx && i < len(blocks); i++ {
		used += height(blocks[i])
		for used > budget && start < i {
			used -= height(blocks[start])
			start++
		}
	}
	used = 0
	for i := start; i < len(blocks); i++ {
		if i > w.idx && used+height(blocks[i]) > budget {
			break
		}
		used += height(blocks[i])
		b.WriteString(blocks[i])
	}
	if w.err != "" {
		b.WriteString("\n" + styWarn.Render(w.err) + "\n")
	}
	for _, l := range wrapHints([]string{"space active / not active", "enter show this board", "shift+↓↑ on the board: next active board", "esc close"}, width) {
		b.WriteString("\n" + styMeta.Render(l))
	}
	return b.String()
}
