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
	"sort"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/fsnotify/fsnotify"

	"github.com/berk/jaira/core/gate"
	"github.com/berk/jaira/core/gitrepo"
	"github.com/berk/jaira/core/lane"
	"github.com/berk/jaira/core/project"
	"github.com/berk/jaira/core/session"
	"github.com/berk/jaira/core/ticket"
)

type mode int

const (
	modeBoard mode = iota
	modeDetail
	modeFilter
	modeHelp
	modeDiff
	modeMove
	modeCreate
	modeMessage
	modeProjects
	modeEdit
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

	// detail holds the fully loaded ticket, since the board only reads
	// frontmatter for speed.
	detail     *ticket.Ticket
	diffText   string
	diffScroll int

	moveTarget int // lane index highlighted in the move picker

	// editIdx is the field being edited in the detail pane and editBuf its
	// working value. The buffer is separate from the ticket so an abandoned edit
	// leaves the file untouched.
	editIdx int
	editBuf string

	// projects is the switcher's list; SwitchTo is set when the user picks one,
	// and the caller reopens the board there. Reopening rather than swapping the
	// store keeps every piece of per-board state (watcher, cursor, filter) from
	// having to be individually reset.
	projects []project.Project
	projIdx  int
	SwitchTo string

	// sessions is what any agent working this tree last checkpointed — the
	// board's view of agent memory.
	sessions []session.Session

	// watch carries filesystem events. A watcher is more responsive than the
	// timer, but the timer stays as a backstop because change notifications are
	// unreliable on some filesystems, notably Windows drives mounted into WSL2.
	watch  chan struct{}
	closer func()

	width, height int
}

type column struct {
	lane    *lane.Lane
	tickets []*ticket.Ticket
}

// New builds a board model.
func New(s *ticket.Store) (*Model, error) {
	m := &Model{store: s, scroll: map[string]int{}}
	if err := m.reload(); err != nil {
		return nil, err
	}
	return m, nil
}

// reload rebuilds the whole view from disk.
//
// A full rescan is used rather than applying incremental changes: at any
// plausible ticket count it is cheap, and an incremental path would be a second
// source of truth that has to stay consistent with edits made entirely outside
// this process (a git pull, another session, a hand edit).
func (m *Model) reload() error {
	lanes, err := lane.Load()
	if err != nil {
		return err
	}
	m.lanes = lanes
	m.warnings = lanes.Warnings

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
	return nil
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
	for _, f := range []string{t.ID, t.Title, t.Goal, t.Context, t.DoD, t.Assignee, t.Status, t.ModelTier} {
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

	case modeMove:
		switch s {
		case "enter":
			m.applyMove()
		case "esc", "q":
			m.mode = modeBoard
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

	case modeHelp, modeMessage:
		if s == "esc" || s == "q" || s == "enter" || s == "?" {
			m.mode = modeBoard
		}
		return m, nil

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
				m.SwitchTo = m.projects[m.projIdx].Root
				m.Close()
				return m, tea.Quit
			}
		}
		return m, nil

	case modeDiff:
		switch s {
		case "esc", "q", "d":
			m.mode = modeDetail
		case "j", "down":
			m.diffScroll++
		case "k", "up":
			if m.diffScroll > 0 {
				m.diffScroll--
			}
		case "ctrl+d", "pgdown":
			m.diffScroll += 10
		case "ctrl+u", "pgup":
			m.diffScroll -= 10
			if m.diffScroll < 0 {
				m.diffScroll = 0
			}
		}
		return m, nil

	case modeDetail:
		switch s {
		case "esc", "q", "enter":
			m.mode = modeBoard
			m.detail = nil
		case "e":
			m.startEdit()
		case "d":
			m.openDiff()
		case "m":
			m.openMove()
		case "j", "down":
			m.mode = modeBoard
			m.detail = nil
			m.moveCard(1)
		case "k", "up":
			m.mode = modeBoard
			m.detail = nil
			m.moveCard(-1)
		}
		return m, nil
	}

	// Board mode.
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
	case "d":
		m.openDetail()
		if m.detail != nil {
			m.openDiff()
		}
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
	case "p":
		m.projects = project.Load()
		m.projIdx = 0
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

func (m *Model) moveCard(d int) {
	m.cardIdx += d
	m.clampCursor()
}

func (m *Model) notify(msg string, isErr bool) {
	m.message, m.isErr = msg, isErr
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
	m.mode = modeDetail
}

func (m *Model) openDiff() {
	if m.detail == nil {
		return
	}
	if len(m.detail.Commits) == 0 {
		m.notify("This ticket records no commits, so there is no diff to show.", false)
		return
	}
	repo := &gitrepo.Repo{Dir: m.store.Root}
	d, err := repo.Diff(m.detail.Commits)
	if err != nil {
		m.notify(err.Error(), true)
		return
	}
	m.diffText, m.diffScroll, m.mode = d, 0, modeDiff
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
	m.mode = modeMove
}

// applyMove runs the same gate checks and the same core mutation the CLI uses.
func (m *Model) applyMove() {
	t := m.selected()
	if t == nil || m.moveTarget >= len(m.lanes.Lanes) {
		m.mode = modeBoard
		return
	}
	target := m.lanes.Lanes[m.moveTarget]
	env := gate.Env{Lanes: m.lanes, All: m.tickets}

	full, err := m.store.Load(t.ID)
	if err != nil {
		m.notify(err.Error(), true)
		return
	}
	vs := gate.CheckAdvance(env, full, gate.Request{To: target.ID})
	if len(vs) > 0 {
		var b strings.Builder
		fmt.Fprintf(&b, "Cannot move to %s:\n\n", target.Name)
		for _, v := range vs {
			fmt.Fprintf(&b, "  • %s\n", v.Message)
		}
		b.WriteString("\nEdit the ticket to supply what is missing, or use the CLI with --force.")
		m.notify(b.String(), true)
		return
	}
	if _, err := m.store.Mutate(full.ID, func(t *ticket.Ticket) error {
		if err := t.Doc().SetScalar(ticket.FieldStatus, target.ID); err != nil {
			return err
		}
		return ticket.SetReady(t.Doc(), gate.Ready(t))
	}); err != nil {
		m.notify(err.Error(), true)
		return
	}
	m.mode = modeBoard
	if err := m.reload(); err != nil {
		m.notify(err.Error(), true)
		return
	}
	m.selectByID(full.ID)
}

// createTicket adds a ticket to the default lane straight from the board.
func (m *Model) createTicket(title string) {
	now := time.Now()
	me := identity(m.store.Root)
	def := m.lanes.Default()
	fields := map[string]string{
		ticket.FieldID:        ticket.NewID(now),
		ticket.FieldTitle:     title,
		ticket.FieldStatus:    def.ID,
		ticket.FieldReady:     "false",
		ticket.FieldCreator:   me,
		ticket.FieldAssignee:  me,
		ticket.FieldCreatedAt: ticket.FormatTime(now),
		ticket.FieldUpdatedAt: ticket.FormatTime(now),
	}
	lists := map[string][]string{ticket.FieldBlockedBy: nil, ticket.FieldCommits: nil}
	t, err := m.store.Create(fields, lists, "")
	if err != nil {
		m.notify(err.Error(), true)
		return
	}
	if err := m.reload(); err != nil {
		m.notify(err.Error(), true)
		return
	}
	m.selectByID(t.ID)
	m.notify(fmt.Sprintf(
		"Created %s in %s.\n\nIt needs a goal, a definition of done, context and an assignee\nbefore it can start. Use:\n\n  jaira set %s goal=… definition-of-done=… context=…",
		ticket.Handle(t.ID), def.Name, ticket.Handle(t.ID)), false)
}
