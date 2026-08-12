package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func makeBoard(t *testing.T, root, name string) string {
	t.Helper()
	dir := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Join(dir, ".jaira", "tickets"), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func newHome(t *testing.T, roots []string, w, h int) *Home {
	t.Helper()
	t.Setenv("JAIRA_HOME", t.TempDir())
	m, err := NewHome(roots)
	if err != nil {
		t.Fatal(err)
	}
	m.width, m.height = w, h
	return m
}

// The launcher exists to answer one question across every board at once: which
// of these needs me. So the list has to render, and it has to name the projects.
func TestHomeListsProjects(t *testing.T) {
	root := t.TempDir()
	a := makeBoard(t, root, "alpha")
	b := makeBoard(t, root, "beta")
	h := newHome(t, []string{a, b}, 120, 30)

	out := stripANSI(h.render())
	for _, want := range []string{"alpha", "beta", "enter open"} {
		if !strings.Contains(out, want) {
			t.Errorf("home screen is missing %q:\n%s", want, out)
		}
	}
}

// The icon is decoration and the list is the point, so a narrow terminal drops
// the icon rather than the thing the screen is for.
func TestHomeDropsTheIconWhenNarrow(t *testing.T) {
	root := t.TempDir()
	a := makeBoard(t, root, "alpha")

	wide := stripANSI(newHome(t, []string{a}, 120, 30).render())
	narrow := stripANSI(newHome(t, []string{a}, 40, 30).render())

	if !strings.ContainsAny(wide, "▀▄") {
		t.Error("the icon is missing from a wide terminal")
	}
	if strings.ContainsAny(narrow, "▀▄") {
		t.Errorf("the icon was drawn in a narrow terminal:\n%s", narrow)
	}
	if !strings.Contains(narrow, "alpha") {
		t.Error("the project list was dropped instead of the icon")
	}
}

func TestHomeSelectsAProject(t *testing.T) {
	root := t.TempDir()
	a := makeBoard(t, root, "alpha")
	b := makeBoard(t, root, "beta")
	h := newHome(t, []string{a, b}, 120, 30)

	h.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	h.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if h.Chosen == "" {
		t.Fatal("enter did not choose a project")
	}
	if filepath.Base(h.Chosen) != "beta" {
		t.Errorf("chose %s, want beta", h.Chosen)
	}
}

func TestHomeQuits(t *testing.T) {
	h := newHome(t, nil, 120, 30)
	h.Update(tea.KeyPressMsg{Code: 'q', Text: "q"})
	if !h.Quit {
		t.Error("q did not quit")
	}
}

// An empty launcher must say what to do rather than showing a blank screen.
func TestHomeWithNoProjectsExplainsItself(t *testing.T) {
	h := newHome(t, nil, 120, 30)
	out := stripANSI(h.render())
	if !strings.Contains(out, "No boards yet") || !strings.Contains(out, "Press a to find") {
		t.Errorf("empty home screen gives no guidance:\n%s", out)
	}
}

// Nothing on this screen may run past the edge of the terminal.
func TestHomeDoesNotOverflow(t *testing.T) {
	root := t.TempDir()
	a := makeBoard(t, root, "a-project-with-a-fairly-long-name")
	for _, w := range []int{20, 40, 80, 120} {
		h := newHome(t, []string{a}, w, 30)
		for _, line := range strings.Split(stripANSI(h.render()), "\n") {
			if len(the(line)) > w {
				t.Errorf("width %d: line is %d cols: %q", w, len(the(line)), line)
			}
		}
	}
}

// the returns the runes of a line, so width is counted in characters.
func the(s string) []rune { return []rune(s) }

// The launcher must be able to add a board without dropping to a shell. This was
// the one of the three chosen ways that did not exist, and it is the one a person
// actually reaches for.
func TestHomeOpensTheDirectoryBrowser(t *testing.T) {
	h := newHome(t, nil, 120, 30)
	h.Update(tea.KeyPressMsg{Code: 'a', Text: "a"})
	if h.browse == nil {
		t.Fatal("a did not open the directory browser")
	}
	out := stripANSI(h.render())
	if !strings.Contains(out, "Add a board") || !strings.Contains(out, "esc cancel") {
		t.Errorf("browser did not render:\n%s", out)
	}
}

func TestBrowserAddsABoardAndItAppears(t *testing.T) {
	root := t.TempDir()
	makeBoard(t, root, "found-me")
	h := newHome(t, nil, 120, 30)
	if len(h.entries) != 0 {
		t.Fatalf("expected an empty launcher, got %d entries", len(h.entries))
	}

	h.Update(tea.KeyPressMsg{Code: 'a', Text: "a"})
	h.browse.goTo(root)
	// The board sorts first, so it is already highlighted.
	h.Update(tea.KeyPressMsg{Code: 'a', Text: "a"})

	if h.browse != nil {
		t.Error("the browser stayed open after adding")
	}
	if len(h.entries) != 1 || h.entries[0].Name != "found-me" {
		t.Fatalf("the added board did not appear on the home screen: %+v", h.entries)
	}
}

// Adding a directory that is not a board must say so rather than registering
// something that will never load.
func TestBrowserRefusesANonBoard(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "just-a-folder"), 0o755); err != nil {
		t.Fatal(err)
	}
	h := newHome(t, nil, 120, 30)
	h.Update(tea.KeyPressMsg{Code: 'a', Text: "a"})
	h.browse.goTo(root)
	h.Update(tea.KeyPressMsg{Code: 'a', Text: "a"})

	if h.browse == nil {
		t.Fatal("the browser closed on a directory that is not a board")
	}
	if !strings.Contains(h.browse.msg, "no board yet") {
		t.Errorf("no explanation given: %q", h.browse.msg)
	}
	if len(h.entries) != 0 {
		t.Error("a non-board was registered")
	}
}

// Scanning registers every board within the depth limit in one keypress.
func TestBrowserScanAddsEverythingFound(t *testing.T) {
	root := t.TempDir()
	makeBoard(t, root, "one")
	makeBoard(t, filepath.Join(root, "nested"), "two")
	h := newHome(t, nil, 120, 30)
	h.Update(tea.KeyPressMsg{Code: 'a', Text: "a"})
	h.browse.goTo(root)
	h.Update(tea.KeyPressMsg{Code: 's', Text: "s"})

	if len(h.entries) != 2 {
		t.Fatalf("scan registered %d boards, want 2: %+v", len(h.entries), h.entries)
	}
}

// Being told "there is no board here, go and run a command" is the round trip
// this screen exists to remove, so it offers to create one instead.
func TestBrowserOffersToInitAndDoesIt(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "fresh-repo"), 0o755); err != nil {
		t.Fatal(err)
	}
	h := newHome(t, nil, 120, 30)
	h.Update(tea.KeyPressMsg{Code: 'a', Text: "a"})
	h.browse.goTo(root)
	h.Update(tea.KeyPressMsg{Code: 'a', Text: "a"}) // no board here

	if h.browse == nil {
		t.Fatal("the browser closed instead of offering to create a board")
	}
	if !strings.Contains(h.browse.msg, "press i") {
		t.Fatalf("no offer to create one: %q", h.browse.msg)
	}

	h.Update(tea.KeyPressMsg{Code: 'i', Text: "i"})
	if h.browse != nil {
		t.Error("the browser stayed open after creating the board")
	}
	if _, err := os.Stat(filepath.Join(root, "fresh-repo", ".jaira", "tickets")); err != nil {
		t.Errorf("no board was created: %v", err)
	}
	if len(h.entries) != 1 || h.entries[0].Name != "fresh-repo" {
		t.Errorf("the new board did not appear: %+v", h.entries)
	}
}

// A fruitless scan offers the same way out.
func TestBrowserScanOffersToInit(t *testing.T) {
	root := t.TempDir()
	h := newHome(t, nil, 120, 30)
	h.Update(tea.KeyPressMsg{Code: 'a', Text: "a"})
	h.browse.goTo(root)
	h.Update(tea.KeyPressMsg{Code: 's', Text: "s"})

	if h.browse == nil || !strings.Contains(h.browse.msg, "press i") {
		t.Fatalf("a fruitless scan gave no way forward: %q", h.browse.msg)
	}
	h.Update(tea.KeyPressMsg{Code: 'i', Text: "i"})
	if _, err := os.Stat(filepath.Join(root, ".jaira", "tickets")); err != nil {
		t.Errorf("scan-then-init did not create a board: %v", err)
	}
}

// A scan that registers boards already on the list changes nothing visible,
// which reads as the key not working. It has to say what it did either way.
func TestScanReportsWhatItDid(t *testing.T) {
	root := t.TempDir()
	makeBoard(t, root, "one")
	makeBoard(t, root, "two")

	h := newHome(t, nil, 120, 30)
	h.Update(tea.KeyPressMsg{Code: 'a', Text: "a"})
	h.browse.goTo(root)
	h.Update(tea.KeyPressMsg{Code: 's', Text: "s"})

	if len(h.entries) != 2 {
		t.Fatalf("scan registered %d boards, want 2", len(h.entries))
	}
	if !strings.Contains(h.msg, "added 2") {
		t.Errorf("scan said %q, expected it to report adding 2", h.msg)
	}

	// Scanning the same place again must say so rather than looking inert.
	h.Update(tea.KeyPressMsg{Code: 'a', Text: "a"})
	h.browse.goTo(root)
	h.Update(tea.KeyPressMsg{Code: 's', Text: "s"})
	if !strings.Contains(h.msg, "already") {
		t.Errorf("a repeat scan said %q, expected it to say they were already listed", h.msg)
	}
	if !strings.Contains(stripANSI(h.render()), "already") {
		t.Error("the message is not shown on the home screen")
	}
}
