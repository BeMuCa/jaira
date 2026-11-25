package tui

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// pressShift sends shift+down or shift+up through the real dispatch.
func pressShift(m *Model, dir string) {
	code := tea.KeyDown
	if dir == "up" {
		code = tea.KeyUp
	}
	m.key(tea.KeyPressMsg{Code: code, Mod: tea.ModShift})
}

// addBoard gives the test model's repository a board called name with the
// lanes named.
func addBoard(t *testing.T, m *Model, name string, lanes ...string) {
	t.Helper()
	dir := filepath.Join(m.store.Root, ticket.DirName, "boards", name, "lanes")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for i, l := range lanes {
		id := strings.ToLower(l)
		body := "---\nid: " + id + "\nname: " + l + "\nprecedence: " + string(rune('1'+i)) + "0\n---\n"
		if err := os.WriteFile(filepath.Join(dir, id+".md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// shift+up/down steps through the active boards only, both ways, wrapping
// round; an inactive board is passed by.
func TestShiftArrowsStepThroughActiveBoardsOnly(t *testing.T) {
	m := newTestModel(t, 150, 32)
	addBoard(t, m, "docs", "Draft")
	addBoard(t, m, "video", "Script")
	if err := lane.SetBoardActive(m.store.Root, "docs", false); err != nil {
		t.Fatal(err)
	}

	for _, step := range []struct{ dir, want string }{
		{"down", "video"}, {"down", ""}, {"up", "video"}, {"up", ""},
	} {
		pressShift(m, step.dir)
		if m.boardName != step.want {
			t.Fatalf("shift+%s: board %q, want %q", step.dir, m.boardName, step.want)
		}
	}

	// On an inactive board the next active one in that direction comes next.
	m.boardName = "docs"
	pressShift(m, "down")
	if m.boardName != "video" {
		t.Errorf("shift+down from inactive docs: board %q, want video", m.boardName)
	}
}

// b lists every board: its name as the heading, who created it, its lanes,
// a tick on the active ones; active first, then the most recently shown,
// then the rest. enter shows the board under the cursor.
func TestBOpensTheBoardWindow(t *testing.T) {
	m := newTestModel(t, 150, 50)
	root := m.store.Root
	addBoard(t, m, "alpha", "One")
	addBoard(t, m, "beta", "Two")
	addBoard(t, m, "video", "Script", "Voice")
	commitAs(t, root, "Vera Video", filepath.Join(ticket.DirName, "boards", "video"))
	for _, b := range []string{"alpha", "beta"} {
		if err := lane.SetBoardActive(root, b, false); err != nil {
			t.Fatal(err)
		}
	}
	// beta was shown more recently than alpha.
	m.markBoardsSeen("alpha")
	time.Sleep(10 * time.Millisecond)
	m.markBoardsSeen("beta")

	press(m, "b")
	if m.mode != modeBoards {
		t.Fatalf("b: mode %v, want the board window", m.mode)
	}
	var order []string
	for _, r := range m.boards.rows {
		order = append(order, r.name)
	}
	if got, want := strings.Join(order, ","), ",video,beta,alpha"; got != want {
		t.Errorf("order = %q, want %q (active, most recently shown, rest)", got, want)
	}

	out := stripANSI(m.render())
	for _, want := range []string{
		"Boards", "✓ default", "✓ video", "○ beta", "○ alpha",
		"created by Vera Video", "Script → Voice", "creator unknown",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the window lacks %q:\n%s", want, out)
		}
	}

	press(m, "j") // video
	press(m, "enter")
	if m.mode != modeBoard || m.boardName != "video" {
		t.Errorf("enter on video: mode %v, board %q", m.mode, m.boardName)
	}
}

// space switches the board under the cursor on and off in .jaira/active-boards
// — the one setting the teamlead reads — and refuses to switch off the last.
func TestSpaceTogglesABoardActive(t *testing.T) {
	m := newTestModel(t, 150, 50)
	root := m.store.Root
	addBoard(t, m, "video", "Script")

	press(m, "b")
	press(m, "j") // video
	press(m, " ")
	b, err := os.ReadFile(lane.ActiveBoardsPath(root))
	if err != nil || string(b) != "default\n" {
		t.Fatalf("after space on video: active-boards = %q, %v; want only default", b, err)
	}
	press(m, "k") // default, now the only active board
	press(m, " ")
	if !strings.Contains(stripANSI(m.render()), "last active board stays active") {
		t.Errorf("switching off the last active board was not refused on screen")
	}
	press(m, "j")
	press(m, " ")
	if got := lane.ActiveBoards(root); len(got) != 2 {
		t.Errorf("after switching video back on: active %q, want both", got)
	}
}

// shift+arrows in an open ticket keep scrolling it rather than switching the
// board behind it.
func TestShiftArrowsInAnOpenTicketDoNotSwitchBoards(t *testing.T) {
	m := newTestModel(t, 150, 32)
	addBoard(t, m, "video", "Script")
	m.openDetail()
	if m.mode != modeDetail {
		t.Fatal("detail did not open")
	}
	pressShift(m, "down")
	if m.boardName != "" || m.mode != modeDetail {
		t.Errorf("shift+down in an open ticket: board %q, mode %v", m.boardName, m.mode)
	}
}

// commitAs makes root a git repository and commits path as author.
func commitAs(t *testing.T, root, author, path string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on PATH")
	}
	for _, args := range [][]string{
		{"init", "-q"},
		{"add", "--", path},
		{"-c", "user.name=" + author, "-c", "user.email=v@example.com", "commit", "-q", "-m", "board"},
	} {
		if out, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
}
