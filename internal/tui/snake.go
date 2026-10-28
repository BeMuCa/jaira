package tui

import (
	"hash/fnv"
	"strconv"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/session"
	"github.com/BeMuCa/jaira/core/ticket"
)

// A card someone is working on right now has a short band of colour running
// round its edge — the snake. The colour says who the ticket is waiting on: an
// agent at work, an agent reviewing with fresh eyes, or a person. It is the
// board's answer to "is anybody on this?", and so also the sign to everyone
// else to leave the ticket alone.

type snakeKind int

const (
	snakeNone snakeKind = iota
	snakeAgent
	snakeReview
	snakeHuman
)

// snakeColour is each kind's 256-colour index: mid tones, so the card's text
// stays readable where the snake passes under it.
var snakeColour = map[snakeKind]int{
	snakeAgent:  30,  // teal
	snakeReview: 97,  // violet
	snakeHuman:  130, // amber
}

const (
	// snakeStep is how long the snake takes to move one cell. The frame rate
	// does not set the speed: a slow frame moves it several cells at once.
	snakeStep = 250 * time.Millisecond

	// The next frame waits this many times as long as the last render took,
	// so drawing stays under 1% of one core however large the board is: the
	// terminal output costs about as much again as the render, which a factor
	// of 100 left out (measured: 2% with 100 cards). The bounds keep a tiny
	// board from spinning and a huge one from freezing.
	frameCostFactor = 250
	minFrameGap     = 120 * time.Millisecond
	maxFrameGap     = time.Second
)

type animFrameMsg struct{}

// snakeOf is the snake a ticket's card carries. A ticket in a lane that waits
// on a person carries one whoever holds it: the work is stopped until somebody
// answers. Otherwise only a live claim does, in review's colour when the lane
// runs in a fresh context. A finished ticket never does.
func (m *Model) snakeOf(t *ticket.Ticket) snakeKind {
	l, ok := m.lanes.Get(t.Status)
	if !ok || l.Terminal {
		return snakeNone
	}
	if l.RequiresQuestion || l.RequiresHumanExit {
		return snakeHuman
	}
	if !m.claimLive(t) {
		return snakeNone
	}
	if l.Context == lane.ContextFresh {
		return snakeReview
	}
	return snakeAgent
}

// claimLive is a claim inside its lease, or one past it whose agent process is
// still running on this machine (liveClaims, read on every reload).
func (m *Model) claimLive(t *ticket.Ticket) bool {
	return claimInLease(t) || (t.ClaimedBy != "" && m.liveClaims[t.ID])
}

func claimInLease(t *ticket.Ticket) bool {
	return t.ClaimedBy != "" && !t.ClaimedAt.IsZero() && time.Since(t.ClaimedAt) <= ticket.ClaimTTL
}

// refreshLiveClaims asks, for every claim past its lease, whether the process
// that took it is still alive here. Done once per reload, not per frame.
func (m *Model) refreshLiveClaims() {
	live := map[string]bool{}
	for _, t := range m.tickets {
		if t.ClaimedBy == "" || claimInLease(t) {
			continue
		}
		if session.ClaimProcessAlive(t.ID, t.ClaimedBy) {
			live[t.ID] = true
		}
	}
	m.liveClaims = live
}

// animating says whether a frame is worth drawing: the animation is on, the
// board is what is on screen, and some card shown there carries a snake.
func (m *Model) animating() bool {
	if !m.anim || m.mode != modeBoard {
		return false
	}
	for _, c := range m.cols {
		for _, t := range c.tickets {
			if m.snakeOf(t) != snakeNone {
				return true
			}
		}
	}
	return false
}

// animKick starts the frame loop when it should run and is not running yet.
func (m *Model) animKick() tea.Cmd {
	if m.animPending || !m.animating() {
		return nil
	}
	return m.nextFrame()
}

func (m *Model) nextFrame() tea.Cmd {
	m.animPending = true
	gap := m.lastRender * frameCostFactor
	if gap < minFrameGap {
		gap = minFrameGap
	}
	if gap > maxFrameGap {
		gap = maxFrameGap
	}
	return tea.Tick(gap, func(time.Time) tea.Msg { return animFrameMsg{} })
}

// frame advances the snake by the clock and asks for the next frame, or lets
// the loop end once nothing is left to animate.
func (m *Model) frame() tea.Cmd {
	m.animPending = false
	if !m.animating() {
		return nil
	}
	m.animPhase = int(time.Since(m.animStart) / snakeStep)
	return m.nextFrame()
}

// toggleAnim flips the snake between running and standing still. A still
// snake is still drawn — who is working on what is worth knowing either way —
// it just costs nothing.
func (m *Model) toggleAnim() {
	m.anim = !m.anim
}

// snakeCells is which cells of a card w columns wide and h rows tall the
// snake covers, keyed by row then column. It runs clockwise round the edge;
// each ticket starts at its own point so two snakes do not move in lockstep.
func (m *Model) snakeCells(id string, w, h int) map[int]map[int]bool {
	if w < 2 || h < 2 {
		return nil
	}
	var ring [][2]int
	for c := 0; c < w; c++ {
		ring = append(ring, [2]int{0, c})
	}
	for r := 1; r < h; r++ {
		ring = append(ring, [2]int{r, w - 1})
	}
	for c := w - 2; c >= 0; c-- {
		ring = append(ring, [2]int{h - 1, c})
	}
	for r := h - 2; r >= 1; r-- {
		ring = append(ring, [2]int{r, 0})
	}
	n := len(ring)
	length := max(3, n/5)
	f := fnv.New32a()
	f.Write([]byte(id))
	phase := 0
	if m.anim {
		phase = m.animPhase
	}
	head := (int(f.Sum32()%uint32(n)) + phase) % n
	cells := map[int]map[int]bool{}
	for k := 0; k < length; k++ {
		p := ring[((head-k)%n+n)%n]
		if cells[p[0]] == nil {
			cells[p[0]] = map[int]bool{}
		}
		cells[p[0]][p[1]] = true
	}
	return cells
}

// paintCells fills a row of width cells, each in the background paramAt
// returns for it: runs of one background are cut out of the styled row whole
// (ansi.Cut keeps the row's own colours) and painted like any other row.
func paintCells(line string, width int, paramAt func(int) string) string {
	var b strings.Builder
	for a := 0; a < width; {
		p := paramAt(a)
		e := a + 1
		for e < width && paramAt(e) == p {
			e++
		}
		b.WriteString(paintRow(p, ansi.Cut(line, a, e)))
		a = e
	}
	return b.String()
}

// snakeLegend is the help screen's sample of a kind's colour.
func snakeLegend(k snakeKind, label string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(strconv.Itoa(snakeColour[k]))).Render(label)
}

func snakeParams(k snakeKind) string {
	return "5;" + strconv.Itoa(snakeColour[k])
}
