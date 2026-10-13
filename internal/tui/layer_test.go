package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/milestone"
	"github.com/BeMuCa/jaira/core/ticket"
)

// mkPlaced writes a milestone placed in the layer after the given names (none:
// a swimlane top), created at the given minute so lanes and children come out
// in a known order.
func mkPlaced(t *testing.T, m *Model, name string, minute int, after []string, ids ...string) {
	t.Helper()
	ms := milestone.New(name, milestone.Palette[minute%len(milestone.Palette)], nowish().Add(time.Duration(minute)*time.Minute))
	ms.SetAfter(after)
	for _, id := range ids {
		ms.Add(id)
	}
	if err := ms.Save(m.store.Root); err != nil {
		t.Fatal(err)
	}
}

func loadMilestone(t *testing.T, m *Model, name string) *milestone.Milestone {
	t.Helper()
	ms, err := milestone.Load(m.store.Root, name)
	if err != nil {
		t.Fatal(err)
	}
	return ms
}

func reloadOrFail(t *testing.T, m *Model) {
	t.Helper()
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
}

var (
	keyUp    = tea.KeyPressMsg{Code: tea.KeyUp}
	keyDown  = tea.KeyPressMsg{Code: tea.KeyDown}
	keyLeft  = tea.KeyPressMsg{Code: tea.KeyLeft}
	keyRight = tea.KeyPressMsg{Code: tea.KeyRight}
)

// L opens the layer from the board; L, esc and q all put the board back.
func TestLOpensTheLayerAndClosesIt(t *testing.T) {
	m := newTestModel(t, 120, 40)
	for _, closeKey := range []string{"L", "esc", "q"} {
		m.key(key("L"))
		if m.mode != modeLayer || m.layer == nil {
			t.Fatalf("L must open the layer, mode = %v", m.mode)
		}
		if out := stripANSI(m.render()); !strings.Contains(out, "Milestone Roadmap") {
			t.Errorf("the layer must name itself the Milestone Roadmap:\n%s", out)
		}
		m.key(key(closeKey))
		if m.mode != modeBoard || m.layer != nil {
			t.Fatalf("%s must close the layer, mode = %v", closeKey, m.mode)
		}
	}
}

// A layer with nothing placed in it shows one thing: the [+] that places the
// first milestone, and says so.
func TestEmptyLayerShowsOnlyThePlus(t *testing.T) {
	m := newTestModel(t, 120, 40)
	mkMilestone(t, m, "unplaced", 45) // in the board, not in the layer
	reloadOrFail(t, m)
	m.key(key("L"))
	out := stripANSI(m.render())
	if strings.Count(out, "[+]") != 2 { // the picture's one, and the header's
		t.Errorf("an empty layer shows exactly one [+]:\n%s", out)
	}
	if !strings.Contains(out, "enter: place the first milestone") {
		t.Errorf("the empty layer must say what enter does:\n%s", out)
	}
	if strings.Contains(out, "unplaced") {
		t.Errorf("a milestone without an after line is not in the layer:\n%s", out)
	}
	if m.layer.sel.kind != onFar {
		t.Errorf("the cursor must be on the [+], got %+v", m.layer.sel)
	}
}

// Bubbles carry their name and done/total; the selected one has the thick
// border; a parent with two children splits its arrow; another lane is set
// off by the dim separator; and the [+] are drawn for the selection only.
func TestLayerDrawsBubblesArrowsAndPlus(t *testing.T) {
	m := newTestModel(t, 160, 40)
	done := titled(t, m, "Add session test harness")
	open := titled(t, m, "Fix session cookie dropped on 302")
	mkPlaced(t, m, "api", 0, nil, done, open)
	mkPlaced(t, m, "docs", 1, []string{"api"})
	mkPlaced(t, m, "tests", 2, []string{"api"})
	mkPlaced(t, m, "ui", 3, nil)
	reloadOrFail(t, m)
	m.key(key("L"))
	out := stripANSI(m.render())

	for _, want := range []string{"api", "docs", "tests", "ui", "1/2", "0/0", "┏", "┗━━━━━━┯", "├", "┐", "▼", "┊"} {
		if !strings.Contains(out, want) {
			t.Errorf("the layer is missing %q:\n%s", want, out)
		}
	}
	// api is a swimlane top with children: no [+] beside it, none under it,
	// only the far one.
	if n := strings.Count(out, "[+]"); n != 1 {
		t.Errorf("a selected top with children offers only the far [+], got %d:\n%s", n, out)
	}

	m.key(keyDown) // onto docs: a child without children of its own
	if m.layer.sel != (layerTarget{onBubble, "docs"}) {
		t.Fatalf("↓ from api must land on docs, got %+v", m.layer.sel)
	}
	if n := strings.Count(stripANSI(m.render()), "[+]"); n != 3 {
		t.Errorf("docs offers [+] beside it, under it and far right, got %d", n)
	}
}

// travelsLeftTo reports whether one line of the picture carries an extra
// arrow that comes from the right (┘), runs left over nothing but line and
// crossings, and turns down (┌) into a ▼ of its own — with at least `crossings`
// drops crossed on the way.
func travelsLeftTo(out string, crossings int) bool {
	lines := strings.Split(out, "\n")
	for k, line := range lines {
		r := []rune(line)
		for i, c := range r {
			if c != '┌' {
				continue
			}
			for j := i + 1; j < len(r); j++ {
				if r[j] == '┘' {
					run := string(r[i+1 : j])
					if strings.Trim(run, "─┼┊") == "" && strings.Count(run, "┼") >= crossings &&
						k+1 < len(lines) && i < len([]rune(lines[k+1])) && []rune(lines[k+1])[i] == '▼' {
						return true
					}
					break
				}
				if !strings.ContainsRune("─┼┊", r[j]) {
					break
				}
			}
		}
	}
	return false
}

// A join into the next row (docs→final beside tests→final) arrives with a ▼
// of its own beside the parent arrow's, crossing the parent's drop on the
// way: ending in the drop made an arrival look like one more crossing
// (Berk, 2026-10-13). It never loops round the side of the lane.
func TestJoinArrowArrivesBesideTheParentArrow(t *testing.T) {
	m := newTestModel(t, 160, 40)
	mkPlaced(t, m, "api", 0, nil)
	mkPlaced(t, m, "docs", 1, []string{"api"})
	mkPlaced(t, m, "tests", 2, []string{"api"})
	mkPlaced(t, m, "final", 3, []string{"docs", "tests"})
	reloadOrFail(t, m)
	m.key(key("L"))
	out := stripANSI(m.render())

	if !travelsLeftTo(out, 1) {
		t.Errorf("the join must run ┌──┼──┘ from tests across docs' drop into a ▼ of its own:\n%s", out)
	}
	if !strings.Contains(out, "▼  ▼") {
		t.Errorf("final must show the join's ▼ three left of the parent arrow's:\n%s", out)
	}
	if strings.Contains(out, "└") {
		t.Errorf("no arrow may loop round the side of the lane:\n%s", out)
	}
	if n := strings.Count(out, "▼"); n != 4 {
		t.Errorf("docs, tests and final get one ▼ each, and the join one of its own: %d\n%s", n, out)
	}
}

// An extra arrow into the next row must not run along the split of the bubble
// above its target: with a sibling between, it would end at the sibling's
// junction and name the wrong successor (review of 2026-10-13). It travels a
// row of its own, crosses the sibling's drop and the target's, and arrives
// in a ▼ of its own.
func TestExtraArrowNeverEndsAtASiblingsJunction(t *testing.T) {
	m := newTestModel(t, 200, 40)
	mkPlaced(t, m, "runde-1", 0, nil)
	mkPlaced(t, m, "runde-2", 1, []string{"runde-1", "parallel"})
	mkPlaced(t, m, "runde-2b", 2, []string{"runde-1"})
	mkPlaced(t, m, "parallel", 3, nil)
	reloadOrFail(t, m)
	m.key(key("L"))
	out := stripANSI(m.render())

	if !travelsLeftTo(out, 2) {
		t.Errorf("the arrow from parallel must cross runde-2b's and runde-2's drops (┼ ┼) into a ▼ of its own:\n%s", out)
	}
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "├") && strings.Contains(line, "┬") && strings.Contains(line, "┘") {
			t.Errorf("a split line must not run on into another lane's arrow:\n%s", out)
		}
	}
}

// ↑ on a bubble of the top row has nowhere to go; it must not jump to the far
// [+], which sits on the same row.
func TestUpOnTheTopRowStays(t *testing.T) {
	m := newTestModel(t, 200, 40)
	mkPlaced(t, m, "a", 0, nil)
	mkPlaced(t, m, "b", 1, nil)
	reloadOrFail(t, m)
	m.key(key("L"))
	m.key(keyUp)
	if m.layer.sel != (layerTarget{onBubble, "a"}) {
		t.Errorf("↑ on the top row must stay put, got %+v", m.layer.sel)
	}
}

// The fill: proportional over 14 cells, never empty or full for work that is
// partly done.
func TestFilledCells(t *testing.T) {
	for _, c := range []struct{ done, total, want int }{
		{0, 0, 0}, {0, 5, 0}, {5, 5, 14}, {1, 2, 7}, {1, 40, 1}, {39, 40, 13}, {6, 7, 12},
	} {
		if got := filledCells(c.done, c.total); got != c.want {
			t.Errorf("filledCells(%d, %d) = %d, want %d", c.done, c.total, got, c.want)
		}
	}
}

// The fill colour is the milestone colour 55 % toward white: as RGB on a
// 24-bit terminal, as the nearest of the 256 colours elsewhere, with black
// text on a light fill.
func TestBubbleFillColourCodes(t *testing.T) {
	for _, c := range []struct {
		trueColor bool
		want      string
	}{
		{true, "48;2;255;201;183"},
		{false, "48;5;223"},
	} {
		m := newTestModel(t, 120, 40)
		m.trueColor = c.trueColor
		ms := milestone.New("salmon", 209, nowish())
		ms.SetAfter(nil)
		ms.Add(titled(t, m, "Add session test harness"))
		ms.Add(titled(t, m, "Fix session cookie dropped on 302"))
		if err := ms.Save(m.store.Root); err != nil {
			t.Fatal(err)
		}
		reloadOrFail(t, m)
		m.key(key("L"))
		out := m.render()
		if !strings.Contains(out, c.want) {
			t.Errorf("trueColor=%v: the fill must be %s:\n%q", c.trueColor, c.want, out)
		}
		if !strings.Contains(out, "38;5;209") {
			t.Errorf("trueColor=%v: the border must be in the milestone colour", c.trueColor)
		}
		if fill, ink := m.bubbleFill(209); ink != inkBlack || fill == nil {
			t.Errorf("a light fill takes black text, got %v", ink)
		}
	}
}

// The arrows move to the nearest target in their direction: down to the
// child under the bubble, right to the [+] beside it and on to the sibling,
// up to the parent, and right across the lanes to the far [+].
func TestLayerNavigationIsGeometric(t *testing.T) {
	m := newTestModel(t, 200, 40)
	mkPlaced(t, m, "a", 0, nil)
	mkPlaced(t, m, "b", 1, []string{"a"})
	mkPlaced(t, m, "c", 2, []string{"a"})
	mkPlaced(t, m, "d", 3, nil)
	reloadOrFail(t, m)
	m.key(key("L"))

	steps := []struct {
		k    tea.KeyPressMsg
		want layerTarget
	}{
		{keyDown, layerTarget{onBubble, "b"}},
		{keyRight, layerTarget{onRight, "b"}},
		{keyRight, layerTarget{onBubble, "c"}},
		{keyDown, layerTarget{onBelow, "c"}},
		{keyUp, layerTarget{onBubble, "c"}},
		{keyUp, layerTarget{onBubble, "a"}},
		{keyRight, layerTarget{onBubble, "d"}},
		{keyRight, layerTarget{onFar, ""}},
		{keyLeft, layerTarget{onBubble, "d"}},
	}
	for i, s := range steps {
		m.key(s.k)
		if m.layer.sel != s.want {
			t.Fatalf("step %d: cursor on %+v, want %+v", i, m.layer.sel, s.want)
		}
	}
}

// A wide layer pans to keep the cursor on screen, and holds still while the
// cursor moves inside the window.
func TestLayerPansToTheSelection(t *testing.T) {
	m := newTestModel(t, 60, 20)
	names := []string{"one", "two", "three", "four", "five"}
	for i, n := range names {
		mkPlaced(t, m, n, i, nil)
	}
	reloadOrFail(t, m)
	m.key(key("L"))
	m.render()
	for _, n := range names[1:] {
		m.key(keyRight)
		if out := stripANSI(m.render()); !strings.Contains(out, n) {
			t.Fatalf("moving onto %s must bring it on screen:\n%s", n, out)
		}
	}
	pan := m.layer.panX
	m.key(keyLeft) // four is still on screen beside five
	m.render()
	if m.layer.panX != pan {
		t.Errorf("a step inside the window must not move it: panX %d -> %d", pan, m.layer.panX)
	}
}

// enter on the empty layer's [+], then n, opens the form; what is typed
// becomes a milestone at the top of a swimlane of its own, with its goal.
func TestNewMilestoneFromTheFarPlus(t *testing.T) {
	m := newTestModel(t, 120, 40)
	m.key(key("L"))
	m.key(key("enter"))
	if m.layer.stage != stageChoose {
		t.Fatalf("enter on a [+] must offer new or connect, stage %v", m.layer.stage)
	}
	m.key(key("n"))
	typeInto(m, "Release Zero")
	m.key(key("tab"))
	typeInto(m, "zero is out")
	m.key(key("enter"))

	ms := loadMilestone(t, m, "release-zero")
	if !ms.Placed || len(ms.After) != 0 || ms.Goal != "zero is out" {
		t.Errorf("want a placed swimlane top with its goal, got placed=%v after=%v goal=%q", ms.Placed, ms.After, ms.Goal)
	}
	if !ms.HasColour() {
		t.Error("a new milestone is given a colour, as the CLI gives one")
	}
	if m.mode != modeLayer || m.layer.stage != stageLayer || m.layer.sel != (layerTarget{onBubble, "release-zero"}) {
		t.Errorf("the cursor must land on the new bubble, mode %v stage %v sel %+v", m.mode, m.layer.stage, m.layer.sel)
	}
}

// The [+] under a bubble makes a successor, the one beside it a sibling.
func TestNewMilestoneBelowAndBeside(t *testing.T) {
	m := newTestModel(t, 160, 40)
	mkPlaced(t, m, "a", 0, nil)
	reloadOrFail(t, m)
	m.key(key("L"))

	m.key(keyDown) // the [+] under a
	if m.layer.sel != (layerTarget{onBelow, "a"}) {
		t.Fatalf("want the [+] under a, got %+v", m.layer.sel)
	}
	m.key(key("enter"))
	m.key(key("n"))
	typeInto(m, "b")
	m.key(key("enter"))
	if got := loadMilestone(t, m, "b").After; len(got) != 1 || got[0] != "a" {
		t.Fatalf("the [+] under a makes a successor of a, after = %v", got)
	}

	m.key(keyRight) // b's sibling [+]
	if m.layer.sel != (layerTarget{onRight, "b"}) {
		t.Fatalf("want the [+] beside b, got %+v", m.layer.sel)
	}
	m.key(key("enter"))
	m.key(key("n"))
	typeInto(m, "c")
	m.key(key("enter"))
	if got := loadMilestone(t, m, "c").After; len(got) != 1 || got[0] != "a" {
		t.Errorf("the [+] beside b makes a sibling: after a, got %v", got)
	}
}

// A taken name is refused in the form and nothing is written.
func TestNewMilestoneRefusesATakenName(t *testing.T) {
	m := newTestModel(t, 120, 40)
	mkMilestone(t, m, "taken", 45)
	reloadOrFail(t, m)
	m.key(key("L"))
	m.key(key("enter"))
	m.key(key("n"))
	typeInto(m, "taken")
	m.key(key("enter"))
	if m.layer.stage != stageForm || !strings.Contains(m.layer.err, "already exists") {
		t.Errorf("want the form kept open with the refusal, stage %v err %q", m.layer.stage, m.layer.err)
	}
	if loadMilestone(t, m, "taken").Placed {
		t.Error("a refused name must not touch the milestone that has it")
	}

	// A name filed into this tree's logbook stays taken: the layer asks the
	// same question 'jaira milestone create' asks.
	dir := filepath.Join(m.store.LogbookDir(), "bc-20200101", ticket.MilestonesSubdir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "gone.md"), []byte("---\ncolor: 45\nstatus: filed\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	m.layer.form.vals[0] = "gone"
	m.key(key("enter"))
	if !strings.Contains(m.layer.err, "filed into the logbook") {
		t.Errorf("want the filed refusal, err %q", m.layer.err)
	}
}

// e edits name, goal and predecessors; an unknown name and a loop are refused
// with the CLI's own words, and an empty list makes the bubble a swimlane top.
func TestEditGoalAndPredecessors(t *testing.T) {
	m := newTestModel(t, 160, 40)
	mkPlaced(t, m, "a", 0, nil)
	mkPlaced(t, m, "b", 1, []string{"a"})
	reloadOrFail(t, m)
	m.key(key("L"))

	m.key(key("e")) // on a
	if m.layer.stage != stageForm || strings.Join(m.layer.form.vals, "|") != "a||" {
		t.Fatalf("e must open the form with a's name, goal and predecessors, got %+v", m.layer.form)
	}
	m.key(key("tab"))
	m.key(key("tab"))
	typeInto(m, "zzz")
	m.key(key("enter"))
	if want := `no milestone "zzz" on this board to come after`; m.layer.err != want {
		t.Errorf("unknown predecessor: err %q, want %q", m.layer.err, want)
	}
	for range 3 {
		m.key(key("backspace"))
	}
	typeInto(m, "b")
	m.key(key("enter"))
	if want := `milestone "a" cannot come after "b": that closes a cycle, a -> b -> a`; m.layer.err != want {
		t.Errorf("loop: err %q, want %q", m.layer.err, want)
	}
	if loadMilestone(t, m, "a").Placed && len(loadMilestone(t, m, "a").After) != 0 {
		t.Fatal("a refused edit must write nothing")
	}
	m.key(key("esc"))

	m.key(keyDown) // b
	m.key(key("e"))
	if m.layer.form.vals[2] != "a" {
		t.Fatalf("the form starts from b's predecessors, got %q", m.layer.form.vals[2])
	}
	m.key(key("tab"))
	typeInto(m, "write the docs")
	m.key(key("tab"))
	m.key(key("backspace"))
	m.key(key("enter"))
	b := loadMilestone(t, m, "b")
	if b.Goal != "write the docs" || !b.Placed || len(b.After) != 0 {
		t.Errorf("want goal set and b a swimlane top, got goal %q placed %v after %v", b.Goal, b.Placed, b.After)
	}
}

// A changed name in e renames the milestone the way 'jaira milestone rename'
// does: the file moves, what came after it names the new name, and the
// cursor and a board narrowed to it stay with it. A taken name is refused
// with nothing written — the goal typed beside it included.
func TestEditRenamesAMilestone(t *testing.T) {
	m := newTestModel(t, 160, 40)
	mkPlaced(t, m, "a", 0, nil)
	mkPlaced(t, m, "b", 1, []string{"a"})
	reloadOrFail(t, m)
	m.setFilterTerm("milestone", "a")
	m.key(key("L"))

	m.key(key("e")) // on a
	m.key(key("backspace"))
	typeInto(m, "b")
	m.key(key("tab"))
	typeInto(m, "new goal")
	m.key(key("enter"))
	if m.layer.stage != stageForm || !strings.Contains(m.layer.err, `"b" already exists`) {
		t.Fatalf("want the form kept open with the refusal, stage %v err %q", m.layer.stage, m.layer.err)
	}
	if loadMilestone(t, m, "a").Goal != "" {
		t.Error("a refused rename must write nothing, the goal included")
	}

	m.key(tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift})
	m.key(key("backspace"))
	typeInto(m, "First Round")
	m.key(key("enter"))
	if m.layer.stage != stageLayer || m.layer.err != "" {
		t.Fatalf("the rename must close the form, stage %v err %q", m.layer.stage, m.layer.err)
	}
	if _, err := os.Stat(milestone.Path(m.store.Root, "a")); !os.IsNotExist(err) {
		t.Errorf("the old file is still there: %v", err)
	}
	if got := loadMilestone(t, m, "first-round"); got.Goal != "new goal" || !got.Placed || len(got.After) != 0 {
		t.Errorf("first-round must carry the goal and a's place, got goal %q placed %v after %v", got.Goal, got.Placed, got.After)
	}
	if b := loadMilestone(t, m, "b"); strings.Join(b.After, ",") != "first-round" {
		t.Errorf("b must come after first-round, got %v", b.After)
	}
	if m.layer.sel != (layerTarget{kind: onBubble, name: "first-round"}) {
		t.Errorf("the cursor must stay on the renamed bubble, got %+v", m.layer.sel)
	}
	if m.filter != "milestone:first-round" {
		t.Errorf("the board filter must follow the rename, got %q", m.filter)
	}
}

// A name left as it stood is the file's own: e on a hand-named Q4.md saves
// the goal and does not rename the file to the normalized q4.
func TestEditKeepsAHandWrittenName(t *testing.T) {
	m := newTestModel(t, 160, 40)
	mkPlaced(t, m, "Q4", 0, nil)
	reloadOrFail(t, m)
	m.key(key("L"))
	m.key(key("e"))
	m.key(key("tab"))
	typeInto(m, "ship")
	m.key(key("enter"))
	if m.layer.stage != stageLayer || m.layer.err != "" {
		t.Fatalf("the save must close the form, stage %v err %q", m.layer.stage, m.layer.err)
	}
	if got := loadMilestone(t, m, "Q4"); got.Goal != "ship" {
		t.Errorf("Q4 must carry the goal, got %q", got.Goal)
	}
	// The directory's own spelling, not os.Stat: a case-insensitive disk
	// (macOS, Windows) finds Q4.md under the name q4.md.
	entries, err := os.ReadDir(filepath.Dir(milestone.Path(m.store.Root, "Q4")))
	if err != nil {
		t.Fatal(err)
	}
	var upper, lower bool
	for _, e := range entries {
		upper = upper || e.Name() == "Q4.md"
		lower = lower || e.Name() == "q4.md"
	}
	if !upper || lower {
		t.Errorf("saving the goal must keep Q4.md as it is: Q4.md there %v, q4.md there %v", upper, lower)
	}
}

// c on a [+] connects a milestone that exists: an unplaced one is placed
// there, the bubble itself is never offered, and one that would close a loop
// is left out.
func TestConnectExistingMilestone(t *testing.T) {
	m := newTestModel(t, 160, 40)
	mkPlaced(t, m, "a", 0, nil)
	mkPlaced(t, m, "b", 1, []string{"a"})
	mkMilestone(t, m, "loose", 45)
	reloadOrFail(t, m)
	m.key(key("L"))
	m.key(keyDown)      // b
	m.key(keyDown)      // the [+] under b
	m.key(key("enter")) // new or connect
	m.key(key("c"))
	if got := strings.Join(m.layer.pick, ","); got != "loose" {
		t.Fatalf("under b the picker offers loose alone (not b, and a would close a loop), got %q", got)
	}
	m.key(key("enter"))
	if got := loadMilestone(t, m, "loose").After; len(got) != 1 || got[0] != "b" {
		t.Errorf("loose must be placed after b, got %v", got)
	}
}

// A placed milestone connected under a [+] keeps its place and gains the
// arrow as a further predecessor.
func TestConnectPlacedMilestoneAddsAnArrow(t *testing.T) {
	m := newTestModel(t, 200, 40)
	mkPlaced(t, m, "a", 0, nil)
	mkPlaced(t, m, "b", 1, []string{"a"})
	mkPlaced(t, m, "x", 2, nil)
	mkPlaced(t, m, "y", 3, []string{"x"})
	reloadOrFail(t, m)
	m.key(key("L"))
	m.key(keyDown) // b
	m.key(keyDown) // the [+] under b
	m.key(key("enter"))
	m.key(key("c"))
	idx := -1
	for i, n := range m.layer.pick {
		if n == "y" {
			idx = i
		}
	}
	if idx < 0 {
		t.Fatalf("y must be offered, got %v", m.layer.pick)
	}
	m.layer.pickIdx = idx
	m.key(key("enter"))
	if got := loadMilestone(t, m, "y").After; len(got) != 2 || got[0] != "x" || got[1] != "b" {
		t.Errorf("y keeps x first and gains b, got %v", got)
	}
	if out := stripANSI(m.render()); !strings.Contains(out, "└") && !strings.Contains(out, "┘") {
		t.Errorf("the new arrow must be drawn:\n%s", out)
	}
}

// c steps the colour through the palette, wrapping at its end.
func TestColourKeyStepsThroughThePalette(t *testing.T) {
	m := newTestModel(t, 120, 40)
	mkPlaced(t, m, "a", len(milestone.Palette)-1, nil)
	reloadOrFail(t, m)
	m.key(key("L"))
	m.key(key("c"))
	if got := loadMilestone(t, m, "a").Colour; got != milestone.Palette[0] {
		t.Errorf("from the last palette colour c wraps to the first, got %d", got)
	}
	m.key(key("c"))
	if got := loadMilestone(t, m, "a").Colour; got != milestone.Palette[1] {
		t.Errorf("c steps to the next colour, got %d", got)
	}
}

// c skips a colour another milestone already wears, so the keypress meant to
// tell two bubbles apart never makes them alike.
func TestColourKeySkipsAColourInUse(t *testing.T) {
	m := newTestModel(t, 120, 40)
	mkPlaced(t, m, "a", 0, nil) // Palette[0]
	mkPlaced(t, m, "b", 1, nil) // Palette[1]
	reloadOrFail(t, m)
	m.key(key("L"))
	if m.layer.sel != (layerTarget{onBubble, "a"}) {
		t.Fatalf("the cursor must start on a, got %+v", m.layer.sel)
	}
	m.key(key("c"))
	if got := loadMilestone(t, m, "a").Colour; got != milestone.Palette[2] {
		t.Errorf("c must pass over b's colour %d to %d, got %d", milestone.Palette[1], milestone.Palette[2], got)
	}
}

// x asks first; y takes the bubble out of the layer, keeps its tickets, and
// hangs what came after it on what came before it.
func TestRemoveReattachesSuccessors(t *testing.T) {
	m := newTestModel(t, 200, 40)
	id := titled(t, m, "Refactor auth middleware")
	mkPlaced(t, m, "a", 0, nil)
	mkPlaced(t, m, "b", 1, []string{"a"}, id)
	mkPlaced(t, m, "c", 2, []string{"b"})
	mkPlaced(t, m, "d", 3, []string{"b"})
	reloadOrFail(t, m)
	m.key(key("L"))
	m.key(keyDown) // b

	m.key(key("x"))
	m.key(key("n"))
	if m.layer.stage != stageLayer || !loadMilestone(t, m, "b").Placed {
		t.Fatal("n must keep b in the layer")
	}
	m.key(key("x"))
	if out := stripANSI(m.render()); !strings.Contains(out, "off the roadmap?") {
		t.Errorf("x must ask first:\n%s", out)
	}
	m.key(key("y"))

	b := loadMilestone(t, m, "b")
	if b.Placed || !b.Has(id) {
		t.Errorf("b must leave the layer and keep its ticket, placed %v", b.Placed)
	}
	for _, n := range []string{"c", "d"} {
		if got := loadMilestone(t, m, n).After; len(got) != 1 || got[0] != "a" {
			t.Errorf("%s must now come after a, got %v", n, got)
		}
	}
}

// enter on a bubble shows its goal and its tickets where they are, the
// logbook included even when it was filed long before the board's window;
// enter opens a ticket and esc comes back to the milestone; b shows the
// board filtered to it.
func TestMilestoneDetail(t *testing.T) {
	m := newTestModel(t, 140, 40)
	todo := titled(t, m, "Refactor auth middleware")
	done := titled(t, m, "Add session test harness")
	old := titled(t, m, "Decide on cookie SameSite policy")
	if _, err := m.store.Logbook(old, "bc-20200101"); err != nil {
		t.Fatal(err)
	}
	ms := milestone.New("round", 45, nowish())
	ms.SetGoal("the round is done")
	ms.SetAfter(nil)
	for _, id := range []string{todo, done, old} {
		ms.Add(id)
	}
	if err := ms.Save(m.store.Root); err != nil {
		t.Fatal(err)
	}
	reloadOrFail(t, m)
	m.key(key("L"))
	if out := stripANSI(m.render()); !strings.Contains(out, "2/3") {
		t.Errorf("done counts the terminal lane and the whole logbook: want 2/3\n%s", out)
	}

	m.key(key("enter"))
	if m.layer.stage != stageDetail {
		t.Fatalf("enter on a bubble opens it, stage %v", m.layer.stage)
	}
	out := stripANSI(m.render())
	for _, want := range []string{"round", "the round is done", "Refactor auth middleware", "Add session test harness", "Decide on cookie SameSite policy", "logbook", "2/3 done"} {
		if !strings.Contains(out, want) {
			t.Errorf("the detail is missing %q:\n%s", want, out)
		}
	}

	m.key(key("enter")) // the first row: todo
	if m.mode != modeDetail || m.detail == nil || m.detail.ID != todo {
		t.Fatalf("enter must open the ticket, mode %v", m.mode)
	}
	m.key(key("esc"))
	if m.mode != modeLayer || m.layer.stage != stageDetail {
		t.Fatalf("esc on the ticket must come back to the milestone, mode %v stage %v", m.mode, m.layer.stage)
	}

	m.key(key("b"))
	if m.mode != modeBoard || m.layer != nil || m.filter != "milestone:round" {
		t.Errorf("b must show the board filtered to the milestone, mode %v filter %q", m.mode, m.filter)
	}
}

// A bubble asks core/link whether a ticket is done, the answer the link window
// and the gate give: an archived ticket counts when it left from the terminal
// lane ("archive · done") and not when it was dropped on the way.
func TestArchivedTicketCountsAsTheLinkIndexSays(t *testing.T) {
	m := newTestModel(t, 140, 40)
	finished := titled(t, m, "Add session test harness") // in done
	dropped := titled(t, m, "Refactor auth middleware")  // in todo
	for _, id := range []string{finished, dropped} {
		if _, err := m.store.Archive(id); err != nil {
			t.Fatal(err)
		}
	}
	mkPlaced(t, m, "round", 0, nil, finished, dropped)
	reloadOrFail(t, m)
	m.key(key("L"))
	if out := stripANSI(m.render()); !strings.Contains(out, "1/2") {
		t.Errorf("want 1/2: the finished archived ticket is done, the dropped one is not\n%s", out)
	}
	m.key(key("enter"))
	if out := stripANSI(m.render()); !strings.Contains(out, "Add session test harness") || strings.Count(out, "archived") != 2 {
		t.Errorf("the detail must list both, titled, as archived:\n%s", out)
	}
}
