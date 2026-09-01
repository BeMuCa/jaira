// Package tui renders the board.
//
// The TUI is a peer of the CLI, not a wrapper around it: both link the same core
// packages directly. Shelling out to the binary for every keypress would add
// process overhead to an interaction that must feel instant, and would mean two
// different paths into the store — exactly the drift the single-write-path rule
// exists to prevent. Every mutation here goes through the same core calls the
// CLI uses, so the gates apply identically.
package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/fsnotify/fsnotify"

	"github.com/BeMuCa/jaira/core/gate"
	"github.com/BeMuCa/jaira/core/gitrepo"
	"github.com/BeMuCa/jaira/core/identity"
	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/project"
	"github.com/BeMuCa/jaira/core/session"
	"github.com/BeMuCa/jaira/core/ticket"
)

type mode int

const (
	modeBoard mode = iota
	modeDetail
	modeFilter
	modeHelp
	modeMove
	modeCreate
	modeMessage
	modeProjects
	modeEdit
	modePipeline
	modeLanes
	modeLaneFocus
	modeSettings
	modeDefaultBoard
	modeDelete
	modeDropBoard
)

// Model is the board's state.
type Model struct {
	store *ticket.Store
	lanes *lane.Set

	tickets []*ticket.Ticket
	cols    []column

	laneIdx int
	cardIdx int

	mode     mode
	filter   string
	input    string // shared line-editor buffer for filter/create
	message  string
	isErr    bool
	warnings []string

	// scroll tracks the first visible card per lane so long lanes can be paged
	// without losing the cursor.
	scroll map[string]int

	// detailScroll is the first visible line of the open ticket. The detail pane
	// is the one view whose content has no upper bound — several checklists, a
	// body, notes — so it is the one view that must be clipped to the window and
	// scrolled rather than rendered whole.
	detailScroll int
	// derivedFor/derivedShas memoise one git derivation per open ticket: the
	// sign-off screen falls back to deriving commits when none are recorded
	// (the normal case — only done demands recording), and doing that per
	// render would shell out on every keypress.
	derivedFor  string
	derivedShas []string

	// detail holds the fully loaded ticket, since the board only reads
	// frontmatter for speed.
	detail *ticket.Ticket

	// detailFrom is the mode the detail pane was opened from, so backing out of
	// it returns there instead of always landing on the board. It defaults to
	// modeBoard, which is correct for every caller that predates modeLaneFocus.
	detailFrom mode

	// returnTo is the same idea for the overlays — the move picker and a
	// message. A dialog is not a reason to lose the page you were on, and a
	// refusal you can act on is worth nothing if acting on it moves you
	// somewhere else. Zero value is modeBoard, which is what these dismissed to
	// before.
	returnTo mode

	// copied marks that the id was just put on the clipboard. OSC52 gives no
	// feedback of its own, so the pane has to say the copy happened — it is
	// transient and clears on the next keypress rather than being real state.
	copied bool

	moveTarget int // lane index highlighted in the move picker

	// follow is the split view: a follow-up being written beside the ticket it
	// follows. Non-nil means both halves of the screen are in use.
	follow *follow

	// pending is a move the gate refused, kept so the refusal itself can offer
	// to override it. Sending the user to the CLI for --force means leaving the
	// board and retyping the move; the refusal is already on screen, so the
	// override belongs here.
	pending *pendingMove

	// editIdx is the field being edited in the detail pane and editBuf its
	// working value. The buffer is separate from the ticket so an abandoned edit
	// leaves the file untouched.
	editIdx int
	editBuf string

	// projects is the switcher's list, loaded when the model is built, not
	// only when 'p' is pressed, so every screen that wants it has it from the
	// first frame. Picking one swaps the store in place (switchBoard).
	projects []project.Project
	projIdx  int

	// me is this machine's identity (identity.Current), cached because it
	// shells out to git config and the card renderer asks for it per card.
	me string

	// myAliases is every name this person goes by, cached for the same reason
	// and refreshed with me. identity.Aliases shells out to git twice and reads
	// a file; asking it per card per frame would be paid on every keypress.
	myAliases []string

	// liveBoards marks, by root, which recorded boards have a non-stale session
	// working in them right now. It is recomputed whenever the project list
	// loads or the board reloads — not on a timer, since this is a board that
	// gets glanced at, not a monitor.
	liveBoards map[string]bool

	// laneScreen is the lane settings screen, non-nil while modeLanes is active.
	laneScreen *laneScreen

	// settingsScreen is the menu behind S, non-nil while modeSettings is
	// active. Both laneScreen and board are opened from it and, on esc,
	// return to it rather than straight to the board.
	settingsScreen *settingsScreen

	// drop is the board-removal screen, non-nil while modeDropBoard is active.
	drop *dropBoard

	// board is the default board screen, non-nil while modeDefaultBoard is
	// active. Reached only through settings — the launcher's own 'd' builds
	// its own copy for the case where no board is open yet.
	board *defaultBoardScreen

	// sessions is what any agent working this tree last checkpointed — the
	// board's view of agent memory.
	sessions []session.Session

	// watch carries filesystem events. A watcher is more responsive than the
	// timer, but the timer stays as a backstop because change notifications are
	// unreliable on some filesystems, notably Windows drives mounted into WSL2.
	watch  chan struct{}
	closer func()

	// thinEmpty is whether lanes holding no tickets are drawn thin on the
	// multi-column board. Session-only, reset on every launch: nothing about
	// the board's *view* is persisted today — the filter, the cursor and the
	// per-lane scroll all start fresh — and the only saved per-user state
	// (lane.SaveDefaultBoard) is a workflow definition rather than a view
	// preference. A new file, or a new key in someone else's file, for one
	// boolean would be a second write path bought for nothing.
	thinEmpty bool

	// laneStart is the first rendered column: the window, stored rather than
	// recomputed from the cursor every frame, which is what lets the columns
	// hold still while the cursor moves across them. boardFit reads and
	// clamps it on every render.
	laneStart int

	// versionLine is the persistent "which version, is there an update"
	// indicator, computed once at construction — see versionLine() in
	// updatecheck.go for why not on every render. It survives switchBoard,
	// since which jaira binary is running is a machine-level fact, not a
	// per-board one.
	versionLine string

	width, height int
}

type column struct {
	lane    *lane.Lane
	tickets []*ticket.Ticket
}

// pendingMove is everything a refused move needs to be retried as a forced one.
// confirm records that f has already been pressed: overriding a gate is not a
// single keystroke, so the notice asks a second time before it writes.
type pendingMove struct {
	ticketID string
	target   *lane.Lane
	actor    string // who the gate was asked on behalf of, and who claims
	claiming bool
	refusals gate.Violations
	confirm  bool
}

// New builds a board model.
func New(s *ticket.Store) (*Model, error) {
	// Loaded here, not only inside the 'p' handler, so the switcher's tabs and
	// its 1-9 binding work from the very first frame — a board that has to be
	// opened once with 'p' before it knows its own neighbours is the bug this
	// exists to fix.
	m := &Model{store: s, scroll: map[string]int{}, projects: project.Load(), versionLine: versionLine()}
	m.me = identity.Current(s.Root)
	// Mutations from the board record who made them, the same as the CLI's.
	s.Actor = m.me
	m.myAliases = identity.Aliases(s.Root)
	if err := m.reload(); err != nil {
		return nil, err
	}
	return m, nil
}

// switchBoard swaps the store underneath the running program. This used to
// quit the program and let the CLI loop start a new one, which dropped the
// alternate screen between the two — the terminal flashing through on every
// board switch. Swapping in place keeps the screen up; the price is resetting
// the per-board view state by hand, which is exactly the state a restart
// used to throw away.
//
// The old watcher's pending waitForChange command stays blocked on the old
// channel forever — one parked goroutine per switch, the same leak profile the
// process restart had, and bounded by how often a person can press a key.
func (m *Model) switchBoard(root string) tea.Cmd {
	s, err := ticket.Discover(root)
	if err != nil {
		m.notify(err.Error(), true)
		return nil
	}
	project.Remember(s.Root)
	m.Close()
	m.store = s
	m.me = identity.Current(s.Root)
	// Mutations from the board record who made them, the same as the CLI's.
	s.Actor = m.me
	m.myAliases = identity.Aliases(s.Root)
	m.scroll = map[string]int{}
	m.laneIdx, m.cardIdx = 0, 0
	m.detail = nil
	m.detailScroll = 0
	// A refused move and an unwritten follow-up both name tickets in the store
	// being swapped out.
	m.pending = nil
	m.follow = nil
	m.filter, m.input = "", ""
	m.projects = project.Load()
	if err := m.reload(); err != nil {
		m.notify(err.Error(), true)
		return nil
	}
	m.startWatching()
	return waitForChange(m.watch)
}

// reload rebuilds the whole view from disk.
//
// A full rescan is used rather than applying incremental changes: at any
// plausible ticket count it is cheap, and an incremental path would be a second
// source of truth that has to stay consistent with edits made entirely outside
// this process (a git pull, another session, a hand edit).
func (m *Model) reload() error {
	lanes, err := lane.Load(m.store.Root)
	if err != nil {
		return err
	}
	m.lanes = lanes
	m.warnings = lanes.Warnings
	// A reload is exactly the moment the git state behind a derived commit
	// list may have changed — a teammate's commit naming the handle arrives,
	// the board refreshes, and a stale memo would keep the sign-off screen
	// on yesterday's answer (or on "no commits" forever).
	m.derivedFor = ""

	tickets, err := m.store.List()
	if err != nil {
		var pe *ticket.PartialError
		if ok := asPartial(err, &pe); !ok {
			return err
		}
		m.warnings = append(m.warnings, pe.Problems...)
	}
	m.tickets = tickets
	if sess, err := session.Load(m.store); err == nil {
		m.sessions = sess
	}
	m.rebuild()
	m.refreshLiveBoards()
	return nil
}

// refreshLiveBoards recomputes which recorded boards, this one included, have
// a session actively working in them. It reads this board's own state from
// m.sessions rather than the disk a second time — reload already loaded it.
func (m *Model) refreshLiveBoards() {
	live := make(map[string]bool, len(m.projects))
	for _, p := range m.projects {
		if p.Root == m.store.Root {
			for _, s := range m.sessions {
				if !s.Stale() {
					live[p.Root] = true
					break
				}
			}
			continue
		}
		if boardHasLiveSession(p.Root) {
			live[p.Root] = true
		}
	}
	m.liveBoards = live
}

// boardHasLiveSession reads another board's session state to answer "is an
// agent working there right now". It must never take the render down: a board
// removed since it was recorded, or sitting on a disconnected drive, simply
// reports not-live. ticket.At and session.Load already tolerate a missing or
// unreadable directory (Glob on a directory that is not there matches
// nothing, no error), so recover() here is a second line of defence against
// anything neither of them anticipated, not the primary guard.
func boardHasLiveSession(root string) (live bool) {
	defer func() {
		if recover() != nil {
			live = false
		}
	}()
	st, err := ticket.At(root)
	if err != nil {
		return false
	}
	sessions, err := session.Load(st)
	if err != nil {
		return false
	}
	for _, s := range sessions {
		if !s.Stale() {
			return true
		}
	}
	return false
}

func asPartial(err error, target **ticket.PartialError) bool {
	pe, ok := err.(*ticket.PartialError)
	if ok {
		*target = pe
	}
	return ok
}

// rebuild groups tickets into columns, applying the current filter.
func (m *Model) rebuild() {
	// Remember what was selected so a refresh does not move the cursor out from
	// under the user. Matching is by ticket ID, never by index.
	var selectedID string
	if t := m.selected(); t != nil {
		selectedID = t.ID
	}
	// Remember the lane too, by ID, wherever the lane is the page rather than a
	// cursor position. Lanes come and go between reloads, so an index would not
	// survive; the ID does.
	var heldLane string
	if m.holdsLane() && m.laneIdx >= 0 && m.laneIdx < len(m.cols) {
		heldLane = m.cols[m.laneIdx].lane.ID
	}

	statuses := make([]string, 0, len(m.tickets))
	for _, t := range m.tickets {
		statuses = append(statuses, t.Status)
	}
	lanes := m.lanes.Columns(statuses)

	byLane := map[string][]*ticket.Ticket{}
	for _, t := range m.tickets {
		if m.filter != "" && !matches(t, m.filter) {
			continue
		}
		byLane[t.Status] = append(byLane[t.Status], t)
	}
	for _, ts := range byLane {
		sort.Slice(ts, func(i, j int) bool {
			if ts[i].UpdatedAt.Equal(ts[j].UpdatedAt) {
				return ts[i].ID < ts[j].ID
			}
			return ts[i].UpdatedAt.After(ts[j].UpdatedAt)
		})
	}

	m.cols = make([]column, 0, len(lanes))
	for _, l := range lanes {
		m.cols = append(m.cols, column{lane: l, tickets: byLane[l.ID]})
	}

	m.clampCursor()
	switch {
	case heldLane != "":
		m.holdLane(heldLane, selectedID)
	case selectedID != "":
		m.selectByID(selectedID)
	}
}

// viewMode resolves whatever is on screen down to the page underneath it: a
// dialog or a message belongs to the page behind it, and an open ticket belongs
// to the view it was opened from, which is where esc puts the user back.
func (m *Model) viewMode() mode {
	md := m.mode
	if md == modeMove || md == modeMessage {
		md = m.returnTo
	}
	if md == modeDetail {
		md = m.detailFrom
	}
	return md
}

// holdsLane reports whether the view showing right now is one lane rather than
// all of them. In those views m.laneIdx is the page the user is on, so a ticket
// moved elsewhere must not drag the screen along with it — on the multi-column
// board the new lane is already visible and following the ticket is the point.
func (m *Model) holdsLane() bool {
	md := m.viewMode()
	return md == modePipeline || md == modeLaneFocus
}

// holdLane keeps a single-lane view pointed at the lane it was showing. Inside
// that lane the cursor still follows the selected ticket — cards are sorted by
// updated_at, so any unrelated edit reorders them — and if the ticket left the
// lane the cursor stays at its index, landing on whatever slid into the gap.
func (m *Model) holdLane(laneID, selectedID string) {
	for li, c := range m.cols {
		if c.lane.ID != laneID {
			continue
		}
		m.laneIdx = li
		for ci, t := range c.tickets {
			if t.ID == selectedID {
				m.cardIdx = ci
				break
			}
		}
		m.clampCursor()
		return
	}
	// The lane itself is gone from the board. Pointing at a lane that no longer
	// exists is worse than following the ticket.
	if selectedID != "" {
		m.selectByID(selectedID)
	}
}

func (m *Model) selectByID(id string) {
	for li, c := range m.cols {
		for ci, t := range c.tickets {
			if t.ID == id {
				m.laneIdx, m.cardIdx = li, ci
				return
			}
		}
	}
	m.clampCursor()
}

func (m *Model) clampCursor() {
	if len(m.cols) == 0 {
		m.laneIdx, m.cardIdx = 0, 0
		return
	}
	if m.laneIdx >= len(m.cols) {
		m.laneIdx = len(m.cols) - 1
	}
	if m.laneIdx < 0 {
		m.laneIdx = 0
	}
	n := len(m.cols[m.laneIdx].tickets)
	if m.cardIdx >= n {
		m.cardIdx = n - 1
	}
	if m.cardIdx < 0 {
		m.cardIdx = 0
	}
}

func (m *Model) selected() *ticket.Ticket {
	if m.laneIdx < 0 || m.laneIdx >= len(m.cols) {
		return nil
	}
	col := m.cols[m.laneIdx]
	if m.cardIdx < 0 || m.cardIdx >= len(col.tickets) {
		return nil
	}
	return col.tickets[m.cardIdx]
}

func (m *Model) currentLane() *lane.Lane {
	if m.laneIdx < 0 || m.laneIdx >= len(m.cols) {
		return nil
	}
	return m.cols[m.laneIdx].lane
}

func matches(t *ticket.Ticket, q string) bool {
	q = strings.ToLower(q)

	// A "key:value" query narrows the search to one field — "assignee:berk"
	// finds berk's tickets without also matching every ticket whose prose
	// mentions them. An unrecognized key is not an error: "http:" in a pasted
	// URL is a search term, not a field, so it falls through to full text.
	if key, val, ok := strings.Cut(q, ":"); ok {
		val = strings.TrimSpace(val)
		known := true
		var field string
		switch strings.TrimSpace(key) {
		case "id", "ticket":
			field = t.ID
		case "title":
			field = t.Title
		case "goal":
			field = t.Goal
		case "context":
			field = t.Context
		case "assignee":
			field = t.Assignee
		case "tag", "tags":
			// Joined rather than matched per tag, so "tag:ui" also finds
			// "ui-polish" — the same substring rule every other key here uses.
			field = strings.Join(t.Tags, " ")
		case "lane", "status":
			field = t.Status
		case "body":
			field = t.Body
		default:
			known = false
		}
		if known {
			return strings.Contains(strings.ToLower(field), val)
		}
	}

	// The body is included because half of what a ticket says lives there — the
	// description and both checklists — and searching only the frontmatter meant
	// the thing you remembered reading was the thing you could not find.
	fields := []string{t.ID, t.Title, t.Goal, t.Context, t.DoD, t.Assignee, t.Status, t.ModelTier, t.Body}
	fields = append(fields, t.Tags...)
	for _, it := range append(append([]ticket.DoDItem{}, t.DoDItems...), t.PlanItems...) {
		fields = append(fields, it.Text)
	}
	for _, f := range fields {
		if strings.Contains(strings.ToLower(f), q) {
			return true
		}
	}
	return false
}

// Init satisfies tea.Model.
func (m *Model) Init() tea.Cmd {
	m.startWatching()
	return tea.Batch(tick(), waitForChange(m.watch))
}

// startWatching subscribes to changes in the ticket and session directories.
// Failure is not fatal: the periodic rescan already keeps the board correct, so a
// platform without working notifications is merely less immediate.
func (m *Model) startWatching() {
	m.watch = make(chan struct{}, 1)
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	for _, d := range []string{m.store.TicketsDir(), m.store.SessionsDir()} {
		_ = w.Add(d)
	}
	m.closer = func() { _ = w.Close() }

	go func() {
		// Coalesce bursts: a git pull touches many files at once, and re-rendering
		// per event would thrash the screen for one logical change.
		var pending bool
		debounce := time.NewTimer(time.Hour)
		defer debounce.Stop()
		for {
			select {
			case _, ok := <-w.Events:
				if !ok {
					return
				}
				if !pending {
					pending = true
					debounce.Reset(200 * time.Millisecond)
				}
			case <-debounce.C:
				if pending {
					pending = false
					select {
					case m.watch <- struct{}{}:
					default:
					}
				}
			case _, ok := <-w.Errors:
				if !ok {
					return
				}
			}
		}
	}()
}

// Close releases the watcher.
func (m *Model) Close() {
	if m.closer != nil {
		m.closer()
	}
}

type changeMsg struct{}

func waitForChange(ch chan struct{}) tea.Cmd {
	if ch == nil {
		return nil
	}
	return func() tea.Msg {
		<-ch
		return changeMsg{}
	}
}

type tickMsg time.Time

// tick drives a periodic rescan. A file watcher is more responsive, but a timer
// is the backstop that keeps the board correct on filesystems where change
// notifications are unreliable — notably Windows drives mounted into WSL2.
func tick() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}

// Update satisfies tea.Model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case tickMsg:
		// Preserve the cursor and any open pane across a background refresh.
		_ = m.reload()
		return m, tick()

	case changeMsg:
		_ = m.reload()
		return m, waitForChange(m.watch)

	case editorDoneMsg:
		defer os.RemoveAll(msg.dir)
		if msg.err != nil {
			m.notify(msg.err.Error(), true)
			return m, nil
		}
		if err := m.applyExternalEdit(msg.id, msg.path); err != nil {
			m.notify(err.Error(), true)
		}
		return m, nil

	case boardEditorDoneMsg:
		if m.board != nil {
			if msg.err != nil {
				m.board.msg, m.board.isErr = msg.err.Error(), true
			} else {
				m.board.msg, m.board.isErr = "", false
			}
		}
		return m, nil

	case newLaneDoneMsg:
		if msg.err != nil {
			if m.laneScreen != nil {
				m.laneScreen.msg, m.laneScreen.isErr = msg.err.Error(), true
			}
			return m, nil
		}
		// The new file may parse to a fresh lane, or to a warning if it was
		// left unedited (an id collision with the last skeleton, say) — either
		// way a full reload is what makes the settings screen agree with it.
		if err := m.reload(); err != nil {
			m.notify(err.Error(), true)
			return m, nil
		}
		if m.laneScreen != nil {
			m.laneScreen = newLaneScreen(m.store, m.lanes)
		}
		return m, nil

	case tea.KeyPressMsg:
		return m.key(msg)
	}
	return m, nil
}

func (m *Model) key(k tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	s := k.String()

	// Line-editing modes consume most keys, so they are handled first.
	switch m.mode {
	case modeEdit:
		return m.editKey(k)

	case modeFilter:
		switch s {
		case "enter":
			m.filter = m.input
			m.mode = modeBoard
			m.rebuild()
		case "esc":
			m.input = ""
			m.mode = modeBoard
		case "backspace":
			if r := []rune(m.input); len(r) > 0 {
				m.input = string(r[:len(r)-1])
				m.filter = m.input
				m.rebuild()
			}
		default:
			// k.Text is what the key produced, so multi-byte characters survive.
			// Gating on a one-byte string dropped every umlaut.
			if k.Text != "" {
				m.input += k.Text
				m.filter = m.input
				m.rebuild()
			}
		}
		return m, nil

	case modeCreate:
		switch s {
		case "enter":
			title := strings.TrimSpace(m.input)
			m.input = ""
			m.mode = modeBoard
			if title != "" {
				m.createTicket(title)
			}
		case "esc":
			m.input = ""
			m.mode = modeBoard
		case "backspace":
			if r := []rune(m.input); len(r) > 0 {
				m.input = string(r[:len(r)-1])
			}
		default:
			if k.Text != "" {
				m.input += k.Text
			}
		}
		return m, nil

	case modeDelete:
		switch s {
		case "enter":
			m.confirmDelete()
		case "esc":
			m.input = ""
			m.mode = modeDetail
		case "backspace":
			if r := []rune(m.input); len(r) > 0 {
				m.input = string(r[:len(r)-1])
			}
		default:
			if k.Text != "" {
				m.input += k.Text
			}
		}
		return m, nil

	case modeMove:
		switch s {
		case "enter":
			m.applyMove()
		case "esc", "q":
			m.mode = m.returnTo
		case "j", "down", "l", "right":
			if m.moveTarget < len(m.lanes.Lanes)-1 {
				m.moveTarget++
			}
		case "k", "up", "h", "left":
			if m.moveTarget > 0 {
				m.moveTarget--
			}
		}
		return m, nil

	case modePipeline:
		if s == "q" || s == "ctrl+c" {
			m.Close()
			return m, tea.Quit
		}
		quit, cmd := m.pipelineKey(s)
		if quit {
			m.Close()
			return m, tea.Quit
		}
		return m, cmd

	case modeHelp:
		switch s {
		case "esc", "q", "enter", "?":
			m.mode = modeBoard
			m.detailScroll = 0
		// The help is taller than most terminals, so it scrolls with the same
		// keys an open ticket does rather than inventing a second vocabulary.
		case "down", "j":
			m.detailScroll++
		case "up", "k":
			if m.detailScroll > 0 {
				m.detailScroll--
			}
		case "pgdown", "ctrl+d":
			m.detailScroll += max(1, m.height-4)
		case "pgup", "ctrl+u":
			m.detailScroll = max(0, m.detailScroll-max(1, m.height-4))
		case "g":
			m.detailScroll = 0
		}
		return m, nil

	case modeMessage:
		// A refused move offers its own way out. f arms the override and y takes
		// it; every other exit drops the offer with the message, so a stale f can
		// never fire a move the user has already walked away from. enter is not a
		// yes — the second question needs a key of its own.
		if p := m.pending; p != nil {
			switch {
			case s == "f" && !p.confirm:
				m.armForce()
				return m, nil
			case s == "y" && p.confirm:
				m.forceMove()
				return m, nil
			case s == "n":
				m.pending = nil
				m.mode = m.returnTo
				return m, nil
			}
		}
		if s == "esc" || s == "q" || s == "enter" || s == "?" {
			m.pending = nil
			m.mode = m.returnTo
		}
		return m, nil

	case modeLanes:
		done := m.laneScreen.key(s)
		cmd := m.laneScreen.pendingCmd
		if done {
			m.laneScreen = nil
			// Reached only through settings now, so closing it goes back one
			// level to the menu, not all the way to the board.
			m.mode = modeSettings
			if err := m.reload(); err != nil {
				m.notify(err.Error(), true)
			}
		}
		return m, cmd

	case modeSettings:
		switch action := m.settingsScreen.key(s); action {
		case settingsActionBack:
			m.settingsScreen = nil
			m.mode = modeBoard
		case settingsActionLanes:
			m.laneScreen = newLaneScreen(m.store, m.lanes)
			m.mode = modeLanes
		case settingsActionDefaultBoard:
			db, err := lane.LoadDefaultBoard()
			if err != nil {
				m.notify(err.Error(), true)
			} else {
				m.board = newDefaultBoardScreen(m.lanes, db)
				m.mode = modeDefaultBoard
			}
		}
		return m, nil

	case modeDefaultBoard:
		done, cmd := m.board.key(s)
		if done {
			m.board = nil
			m.mode = modeSettings
		}
		return m, cmd

	case modeProjects:
		switch s {
		case "esc", "q", "p":
			m.mode = modeBoard
		case "j", "down":
			if m.projIdx < len(m.projects)-1 {
				m.projIdx++
			}
		case "k", "up":
			if m.projIdx > 0 {
				m.projIdx--
			}
		case "enter":
			if m.projIdx < len(m.projects) {
				cmd := m.switchBoard(m.projects[m.projIdx].Root)
				m.mode = modeBoard
				return m, cmd
			}
		case "x":
			// Removing a board is offered where you already stand looking at
			// the list of them, and it opens a screen rather than acting: this
			// is the one key on this screen that cannot be taken back.
			if m.projIdx < len(m.projects) {
				p := m.projects[m.projIdx]
				m.drop = newDropBoard(p.Root, p.Name, p.Root == m.store.Root)
				m.mode = modeDropBoard
			}
		}
		return m, nil

	case modeDropBoard:
		if m.drop == nil {
			m.mode = modeProjects
			return m, nil
		}
		done, removed := m.drop.key(s)
		if !done {
			return m, nil
		}
		m.drop = nil
		m.mode = modeProjects
		if removed {
			// Load drops any board whose .jaira is gone, so a fully removed one
			// leaves the list by itself; a partly removed one stays, correctly.
			m.projects = project.Load()
			if m.projIdx >= len(m.projects) {
				m.projIdx = max(0, len(m.projects)-1)
			}
			m.refreshLiveBoards()
		}
		return m, nil

	case modeDetail:
		m.copied = false
		switch s {
		case "esc", "q", "enter":
			// One level back: out of the split first, onto the ticket the follow-up
			// was written for, and only then out of the ticket.
			if m.follow != nil {
				m.closeFollowUp()
				return m, nil
			}
			m.mode = m.detailFrom
			m.detail = nil
		case "n":
			m.startFollowUp()
		case "tab":
			if m.follow != nil {
				m.follow.focusLeft = !m.follow.focusLeft
			}
		case "shift+down":
			m.scrollSrc(1)
		case "shift+up":
			m.scrollSrc(-1)
		case "e":
			m.startEdit()
		case "E":
			return m.openInEditor()
		case "a":
			if m.atHumanCheckpoint() {
				m.accept()
			}
		case "f":
			if m.atHumanCheckpoint() {
				m.followUp()
			}
		case "y":
			if m.detail != nil {
				m.copied = true
				return m, tea.SetClipboard(m.detail.ID)
			}
		case "m":
			m.openMove()
		case "X":
			// Not x: x is archive on the board and must keep meaning that
			// everywhere. Shift, then the handle typed back — the only
			// irreversible thing the board can do costs two deliberate acts.
			if m.detail != nil {
				m.mode = modeDelete
				m.input = ""
			}
		case "b":
			// A blocked ticket names its blocker; the reader's next question is
			// always "and what is that one waiting on" — so the link is walkable.
			// Repeated presses follow the chain; esc returns to the board.
			if m.detail != nil && len(m.detail.BlockedBy) > 0 {
				if full, err := m.store.Load(m.detail.BlockedBy[0]); err != nil {
					m.notify(err.Error(), true)
				} else {
					m.detail = full
					m.detailScroll = 0
				}
			}
		// Arrows scroll, j/k switch tickets: the arrows are the base movement
		// vocabulary and must work on every screen a ticket can be read on,
		// while jumping to the neighbouring ticket stays one key away. In the
		// split they move whichever pane tab has focused.
		case "down":
			m.scrollFocused(1)
		case "up":
			m.scrollFocused(-1)
		case "pgdown", "ctrl+d":
			m.scrollFocused(max(1, m.height-4))
		case "pgup", "ctrl+u":
			m.scrollFocused(-max(1, m.height-4))
		case "j":
			m.leaveDetail()
			m.moveCard(1)
		case "k":
			m.leaveDetail()
			m.moveCard(-1)
		}
		return m, nil

	case modeLaneFocus:
		if s == "ctrl+c" {
			m.Close()
			return m, tea.Quit
		}
		m.laneFocusKey(s)
		return m, nil
	}

	// Board mode.
	// 1-9 switches board here exactly as it does in the compact view — both
	// call switchToProject so "which board is number three" is decided once.
	if len(s) == 1 && s[0] >= '1' && s[0] <= '9' {
		return m, m.switchToProject(int(s[0] - '0'))
	}
	switch s {
	case "q", "ctrl+c":
		m.Close()
		return m, tea.Quit
	case "h", "left":
		m.moveLane(-1)
	case "l", "right":
		m.moveLane(1)
	case "j", "down":
		m.moveCard(1)
	case "k", "up":
		m.moveCard(-1)
	case "g":
		m.cardIdx = 0
	case "G":
		m.cardIdx = len(m.cols[m.laneIdx].tickets) - 1
		m.clampCursor()
	case "enter":
		m.openDetail()
	case "x":
		m.archiveSelected()
	case "v":
		// The compact view: the whole flow on one screen, for watching several
		// agents rather than working one ticket.
		m.mode = modePipeline
	case "z":
		// Next to v: both are about what the screen shows, not about a ticket.
		m.toggleEmptyLanes()
	case "/":
		m.mode = modeFilter
		m.input = m.filter
	case "esc":
		if m.filter != "" {
			m.filter, m.input = "", ""
			m.rebuild()
		}
	case "?":
		m.mode = modeHelp
		m.detailScroll = 0
	case "p":
		// The list can change while the TUI is open (another board opened
		// elsewhere), so this reload stays — it is just no longer the only place
		// the list gets loaded.
		m.projects = project.Load()
		m.projIdx = 0
		m.refreshLiveBoards()
		if len(m.projects) <= 1 {
			m.notify("No other boards recorded yet.\n\nOpen jaira inside another repository and it will appear here.", false)
		} else {
			m.mode = modeProjects
		}
	case "n":
		m.mode = modeCreate
		m.input = ""
	case "m":
		m.openMove()
	case "r":
		if err := m.reload(); err != nil {
			m.notify(err.Error(), true)
		} else {
			m.notify("Reloaded", false)
		}
	case "S":
		m.settingsScreen = newSettingsScreen()
		m.mode = modeSettings
	}
	return m, nil
}

func (m *Model) moveLane(d int) {
	if len(m.cols) == 0 {
		return
	}
	m.laneIdx = (m.laneIdx + d + len(m.cols)) % len(m.cols)
	m.cardIdx = 0
	m.clampCursor()
}

// toggleEmptyLanes flips whether lanes holding no tickets are drawn thin.
// Nothing else changes: a thin lane is still on the board, so the cursor may
// sit on it and h/l step onto it like any other.
func (m *Model) toggleEmptyLanes() {
	m.thinEmpty = !m.thinEmpty
}

func (m *Model) moveCard(d int) {
	m.cardIdx += d
	m.clampCursor()
}

// notify puts a message on screen. Dismissing it goes back to the page it was
// raised on. A message raised by the move picker keeps the picker's own answer:
// the picker is on its way out, and the page behind it is what the message
// belongs to. The screens that own their own state (settings, the lane editor,
// the default-board picker) still dismiss to the board, which is what they did
// before this existed.
func (m *Model) notify(msg string, isErr bool) {
	m.message, m.isErr = msg, isErr
	switch m.mode {
	case modeBoard, modePipeline, modeLaneFocus, modeDetail:
		m.returnTo = m.mode
	case modeMove, modeMessage:
		// Leave returnTo alone.
	default:
		m.returnTo = modeBoard
	}
	m.mode = modeMessage
}

func (m *Model) openDetail() {
	t := m.selected()
	if t == nil {
		return
	}
	full, err := m.store.Load(t.ID)
	if err != nil {
		m.notify(err.Error(), true)
		return
	}
	m.detail = full
	m.detailFrom = m.mode
	m.detailScroll = 0
	m.mode = modeDetail
}

// isMe reports whether a recorded name is this person, under any of the names
// they go by. One person is not one string — a work address in one ticket, a
// git user.name in another — and a marker that called your own change somebody
// else's would be worse than no marker at all.
func (m *Model) isMe(who string) bool {
	who = strings.TrimSpace(who)
	if who == "" {
		return false
	}
	for _, a := range m.myAliases {
		if identity.Same(a, who) {
			return true
		}
	}
	return identity.Same(m.me, who)
}

func (m *Model) openMove() {
	t := m.selected()
	if t == nil {
		return
	}
	m.moveTarget = 0
	for i, l := range m.lanes.Lanes {
		if l.ID == t.Status {
			m.moveTarget = i
		}
	}
	m.returnTo = m.mode
	m.mode = modeMove
}

// gateEnv assembles the state gate.CheckAdvance needs, the same way the CLI's
// loadEnv does, so both interfaces enforce identically — the promise
// core/gate's own package doc makes. The two render sites below (renderCard,
// renderDetail) pay nothing extra for the DeriveCommits closure: gate.Ready
// and gate.Actionable never call it, only CheckAdvance does, and that only
// runs at the moment a move is actually attempted.
func (m *Model) gateEnv() gate.Env {
	repo := &gitrepo.Repo{Dir: m.store.Root}
	return gate.Env{
		Lanes: m.lanes,
		All:   m.tickets,
		DeriveCommits: func(t *ticket.Ticket) []string {
			shas, err := repo.CommitsForTicket(t.Path, t.ID)
			if err != nil {
				return nil
			}
			return shas
		},
	}
}

// applyMove runs the same gate checks and the same core mutation the CLI uses.
func (m *Model) applyMove() {
	t := m.selected()
	if t == nil || m.moveTarget >= len(m.lanes.Lanes) {
		m.mode = m.returnTo
		return
	}
	target := m.lanes.Lanes[m.moveTarget]
	env := m.gateEnv()

	full, err := m.store.Load(t.ID)
	if err != nil {
		m.notify(err.Error(), true)
		return
	}
	// Moving an unassigned ticket claims it — the same rule the CLI's move
	// applies. Set before the gate check so the promotion gate's assignee
	// requirement is satisfied by the pull itself.
	me := identity.Current(m.store.Root)
	claiming := strings.TrimSpace(full.Assignee) == ""
	if claiming {
		full.Assignee = me
	}
	vs := gate.CheckAdvance(env, full, gate.Request{
		To: target.ID, Actor: me, ActorAliases: identity.Aliases(m.store.Root),
	})
	if len(vs) > 0 {
		m.pending = &pendingMove{
			ticketID: full.ID,
			target:   target,
			actor:    me,
			claiming: claiming,
			refusals: vs,
		}
		var b strings.Builder
		fmt.Fprintf(&b, "Cannot move to %s:\n\n", target.Name)
		b.WriteString(refusalBullets(vs))
		b.WriteString("\nEdit the ticket to supply what is missing, or press f to override — the same override the CLI spells --force.")
		m.notify(b.String(), true)
		return
	}
	m.pending = nil
	if _, err := m.store.Mutate(full.ID, moveMutation(target.ID, me, claiming)); err != nil {
		m.notify(err.Error(), true)
		return
	}
	m.finishMove(full.ID)
}

// moveMutation is the write a move performs. The gated path and the forced one
// share it so a forced move leaves behind exactly what a clean one would.
func moveMutation(to, me string, claiming bool) func(*ticket.Ticket) error {
	return func(t *ticket.Ticket) error {
		if claiming {
			if err := t.Doc().SetScalar(ticket.FieldAssignee, me); err != nil {
				return err
			}
			t.Assignee = me
		}
		if err := t.Doc().SetScalar(ticket.FieldStatus, to); err != nil {
			return err
		}
		return ticket.SetReady(t.Doc(), gate.Ready(t))
	}
}

func refusalBullets(vs gate.Violations) string {
	var b strings.Builder
	for _, v := range vs {
		fmt.Fprintf(&b, "  • %s\n", v.Message)
	}
	return b.String()
}

// armForce is the first f: it asks again rather than writing. A gate refusal is
// the project's own opinion about the ticket, so overriding it takes two keys.
func (m *Model) armForce() {
	p := m.pending
	p.confirm = true
	var b strings.Builder
	fmt.Fprintf(&b, "Move %s to %s anyway, overriding %d refusal(s)?\n\n",
		ticket.Handle(p.ticketID), p.target.Name, len(p.refusals))
	b.WriteString(refusalBullets(p.refusals))
	m.notify(b.String(), true)
}

// forceMove carries out the move the gate refused. This is the TUI's --force:
// the CLI overrides by simply not returning the refusal and then running the
// same mutation (internal/cli/flow.go), so the two agree on what a forced move
// leaves behind — including that nothing is written to record the override. What
// it overrode is said out loud instead, the way the CLI prints it.
func (m *Model) forceMove() {
	p := m.pending
	m.pending = nil
	if _, err := m.store.Mutate(p.ticketID, moveMutation(p.target.ID, p.actor, p.claiming)); err != nil {
		m.notify(err.Error(), true)
		return
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%s → %s\n\nOverrode %d gate refusal(s):\n\n",
		ticket.Handle(p.ticketID), p.target.Name, len(p.refusals))
	b.WriteString(refusalBullets(p.refusals))

	m.finishMove(p.ticketID)
	if m.mode == modeMessage {
		// finishMove has worse news of its own; do not paper over it.
		return
	}
	m.notify(b.String(), false)
}

// finishMove puts the user back on the page the move was started from and
// refreshes what that page shows. An open ticket is reloaded rather than closed:
// the move was made while looking at it, so it has to show the lane it landed
// in, not the one it left.
func (m *Model) finishMove(id string) {
	m.mode = m.returnTo
	if err := m.reload(); err != nil {
		m.notify(err.Error(), true)
		return
	}
	// A move the user made themselves still follows the ticket — unlike one made
	// from another shell, this cursor jump is the answer to "where did it go".
	m.selectByID(id)
	if m.mode != modeDetail {
		return
	}
	full, err := m.store.Load(id)
	if err != nil {
		// Closing the ticket is better than leaving a stale one on screen.
		m.detail = nil
		m.mode = m.detailFrom
		return
	}
	m.detail = full
}

// createTicket adds a ticket to the default lane straight from the board.
func (m *Model) createTicket(title string) {
	now := time.Now()
	me := identity.Current(m.store.Root)
	def := m.lanes.Default()
	fields := map[string]string{
		ticket.FieldID:      ticket.NewID(now),
		ticket.FieldTitle:   title,
		ticket.FieldStatus:  def.ID,
		ticket.FieldReady:   "false",
		ticket.FieldCreator: me,
		// Unassigned on purpose, same as the CLI's create: capturing and
		// claiming are two acts, and the claim happens when someone pulls the
		// ticket into work.
		ticket.FieldCreatedAt: ticket.FormatTime(now),
		ticket.FieldUpdatedAt: ticket.FormatTime(now),
	}
	lists := map[string][]string{ticket.FieldBlockedBy: nil, ticket.FieldCommits: nil}
	// NewBody is the one starting shape for a ticket, CLI or TUI: without this,
	// a board-created ticket had no Options section at all, so a default
	// board's "always brainstorm" setting silently did not apply to half the
	// tickets created against it.
	db, _ := lane.LoadDefaultBoard()
	body := ticket.NewBody(title, "", lane.ResolveOptions(m.lanes, db))
	t, err := m.store.Create(fields, lists, body)
	if err != nil {
		m.notify(err.Error(), true)
		return
	}
	if err := m.reload(); err != nil {
		m.notify(err.Error(), true)
		return
	}
	m.selectByID(t.ID)
	// Capture stays cheap — the prompt still asks only for a title — but a ticket
	// with nothing but a title cannot leave the backlog, and sending the user to
	// the CLI to fix that was the reason the board could capture work but not
	// specify it. Hand straight over to the field editor instead.
	if full, err := m.store.Load(t.ID); err == nil {
		m.detail = full
		m.startEdit()
	}
}

// confirmDelete removes the open ticket if what was typed is its handle.
//
// A yes/no dialog would be wrong here. Everything else the board does is
// reversible — an archived ticket is restored, a forced move is reported and can
// be moved back — and this is not, so the second step asks for something only
// someone looking at the ticket can supply.
func (m *Model) confirmDelete() {
	t := m.detail
	if t == nil {
		m.mode = modeBoard
		return
	}
	handle := ticket.Handle(t.ID)
	typed := strings.TrimSpace(m.input)
	m.input = ""
	if !strings.EqualFold(typed, handle) {
		m.mode = modeDetail
		m.notify(fmt.Sprintf("That is not %s. Nothing was deleted.", handle), true)
		return
	}
	path, err := m.store.Delete(t.ID)
	if err != nil {
		m.mode = modeDetail
		m.notify(err.Error(), true)
		return
	}
	m.mode = modeBoard
	m.detail = nil
	if err := m.reload(); err != nil {
		m.notify(err.Error(), true)
		return
	}
	m.notify(fmt.Sprintf("Deleted %s.\n\n%s is gone; there is no restore for this one.",
		handle, filepath.Base(path)), false)
}

// archiveSelected takes a ticket off the board.
//
// It moves the file rather than deleting it, and it is offered here because the
// board is where you notice that a done ticket has stopped being worth looking
// at. Restoring is 'jaira restore', which is why nothing asks for confirmation:
// the action is reversible.
func (m *Model) archiveSelected() {
	t := m.selected()
	if t == nil {
		return
	}
	dst, err := m.store.Archive(t.ID)
	if err != nil {
		m.notify(err.Error(), true)
		return
	}
	if err := m.reload(); err != nil {
		m.notify(err.Error(), true)
		return
	}
	m.mode = modeBoard
	m.detail = nil
	m.notify(fmt.Sprintf("Archived %s.\n\nNothing was deleted — bring it back with:\n\n  jaira restore %s",
		ticket.Handle(t.ID), filepath.Base(dst)), false)
}
