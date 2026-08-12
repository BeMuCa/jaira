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
	if !strings.Contains(out, "No boards yet") || !strings.Contains(out, "jaira init") {
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
