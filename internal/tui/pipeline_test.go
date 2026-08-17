package tui

import (
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/BeMuCa/jaira/core/project"
)

func pipelineModel(t *testing.T, w, h int) *Model {
	t.Helper()
	m := newTestModel(t, w, h)
	m.mode = modePipeline
	return m
}

// The compact view exists to be glanced at while several agents run, so the
// whole flow and its order must be on one screen.
func TestPipelineShowsEveryStepInOrder(t *testing.T) {
	out := stripANSI(pipelineModel(t, 130, 34).renderPipeline())
	// "Human Review" does not fit the compact view's cell whole (see
	// renderStep's truncation), so this checks the prefix that survives it.
	for _, want := range []string{"Backlog", "Todo", "Pre-process", "HITL", "Review", "Human Rev", "Done"} {
		if !strings.Contains(out, want) {
			t.Errorf("step %q missing:\n%s", want, out)
		}
	}
	if !strings.Contains(out, "▶") {
		t.Error("no arrows between steps")
	}
	if strings.Index(out, "Backlog") > strings.Index(out, "Done") {
		t.Error("steps are not in pipeline order")
	}
}

// A cell whose name wraps is taller than its neighbours and the row stops
// lining up, which is why names are truncated rather than wrapped.
func TestPipelineCellsStayOneRowTall(t *testing.T) {
	m := pipelineModel(t, 130, 34)
	for _, line := range strings.Split(stripANSI(m.renderPipeline()), "\n") {
		if strings.Count(line, "╭") > 0 && strings.Count(line, "╮") != strings.Count(line, "╭") {
			t.Errorf("a cell border is broken across the row: %q", line)
		}
	}
}

func TestPipelineNavigationAndOpening(t *testing.T) {
	m := pipelineModel(t, 130, 34)
	m.laneIdx = 0
	m.pipelineKey("d")
	if m.laneIdx != 1 {
		t.Errorf("d did not move right: %d", m.laneIdx)
	}
	m.pipelineKey("a")
	if m.laneIdx != 0 {
		t.Errorf("a did not move left: %d", m.laneIdx)
	}
	m.pipelineKey("enter")
	if m.mode != modeLaneFocus {
		t.Error("enter did not open the step as a single-lane view")
	}
}

func TestPipelineTogglesFromTheBoard(t *testing.T) {
	m := newTestModel(t, 130, 34)
	m.mode = modeBoard
	m.key(tea.KeyPressMsg{Code: 'v', Text: "v"})
	if m.mode != modePipeline {
		t.Fatalf("v did not open the compact view (mode=%v)", m.mode)
	}
	m.key(tea.KeyPressMsg{Code: 'v', Text: "v"})
	if m.mode != modeBoard {
		t.Errorf("v did not return to the board (mode=%v)", m.mode)
	}
}

// A step that just received a ticket lights its incoming arrow, so movement is
// visible without reading counts.
func TestPipelineLightsTheArrowOnAFreshMove(t *testing.T) {
	m := pipelineModel(t, 130, 34)
	steps := m.pipelineSteps()
	var anyFresh bool
	for _, s := range steps {
		if s.fresh {
			anyFresh = true
		}
	}
	if !anyFresh {
		t.Fatal("the fixture has just-created tickets, so some step should be fresh")
	}
	// Age every ticket past the window and it should go quiet.
	for _, tk := range m.tickets {
		tk.UpdatedAt = time.Now().Add(-2 * recentMove)
	}
	for _, s := range m.pipelineSteps() {
		if s.fresh {
			t.Errorf("step %q still reads as fresh after ageing", s.name)
		}
	}
}

func TestPipelineDoesNotOverflow(t *testing.T) {
	for _, w := range []int{40, 80, 130, 200} {
		m := pipelineModel(t, w, 34)
		for _, line := range strings.Split(stripANSI(m.renderPipeline()), "\n") {
			if lipgloss.Width(line) > w {
				t.Errorf("width %d: line is %d cols: %q", w, lipgloss.Width(line), line)
			}
		}
	}
}

// Number keys switch project — in place, inside the running program. Quitting
// and restarting per switch dropped the alternate screen for a frame, which
// flashed the terminal through on every switch.
func TestPipelineSwitchesProjectByNumber(t *testing.T) {
	m := pipelineModel(t, 130, 34)
	other := newTestStore(t)
	m.projects = []project.Project{
		{Root: m.store.Root, Name: "current"},
		{Root: other.Root, Name: "other"},
	}
	quit, _ := m.pipelineKey("2")
	if quit {
		t.Fatal("switching boards must not quit the program")
	}
	if m.store.Root != other.Root {
		t.Errorf("store root = %q, want the other board %q", m.store.Root, other.Root)
	}
	if m.mode != modePipeline {
		t.Error("switching boards changed the view mode")
	}

	// Choosing the project already open is a no-op, not a pointless reload.
	m2 := pipelineModel(t, 130, 34)
	m2.projects = []project.Project{{Root: m2.store.Root, Name: "current"}}
	before := m2.store
	if _, cmd := m2.pipelineKey("1"); cmd != nil {
		t.Error("selecting the current project produced a command")
	}
	if m2.store != before {
		t.Error("selecting the current project swapped the store")
	}

	// A number with no project behind it does nothing.
	if quit, cmd := m2.pipelineKey("9"); quit || cmd != nil {
		t.Error("an unused number asked to switch")
	}
}

// The hints belong at the bottom of the terminal. Sitting them directly under
// the diagram left the rest of the screen blank beneath them, which reads as the
// view having failed to fill the space.
func TestPipelineFillsTheHeight(t *testing.T) {
	for _, h := range []int{20, 30, 44} {
		m := pipelineModel(t, 120, h)
		lines := strings.Split(m.renderPipeline(), "\n")
		if len(lines) < h-2 {
			t.Errorf("height %d: rendered only %d lines, leaving a gap below the hints", h, len(lines))
		}
		if len(lines) > h {
			t.Errorf("height %d: rendered %d lines, more than the terminal has", h, len(lines))
		}
	}
}

// The flow runs as a serpentine, so a wrapped row continues the path instead of
// jumping back to the left edge.
func TestPipelineSnakesBackOnTheSecondRow(t *testing.T) {
	out := stripANSI(pipelineModel(t, 60, 30).renderPipeline())
	if !strings.Contains(out, "◀") {
		t.Errorf("the second row does not run backwards:\n%s", out)
	}
	if !strings.Contains(out, "▼") {
		t.Errorf("no connector between rows:\n%s", out)
	}
}
