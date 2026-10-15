package tui

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"

	"github.com/BeMuCa/jaira/core/release"
)

func fixtureNotes() []release.Entry {
	return []release.Entry{
		{Version: "0.5.2", Changes: []string{"! run jaira update now", "a plain change"}},
		{Version: "0.5.1", Changes: []string{"an older change"}},
	}
}

func TestNotesViewMarksImportantLines(t *testing.T) {
	v := newNotesView("Release notes", fixtureNotes())
	out := stripANSI(v.render(60, 30))
	for _, want := range []string{"! run jaira update now", "- a plain change", "- an older change", "0.5.2", "0.5.1"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "- ! ") {
		t.Errorf("the marker was printed as a bullet and a marker:\n%s", out)
	}
	// An important line is its own styled run, not part of the plain ones.
	if lines := v.lines(60); lipgloss.Width(lines[1]) != lipgloss.Width("! run jaira update now") {
		t.Errorf("important line was padded: %q", lines[1])
	}
}

func TestNotesViewWrapsWithinTheWidth(t *testing.T) {
	v := newNotesView("t", []release.Entry{{Version: "1.0.0", Changes: []string{"! " + strings.Repeat("word ", 30)}}})
	for _, l := range strings.Split(stripANSI(v.render(30, 30)), "\n") {
		if lipgloss.Width(l) > 30 {
			t.Errorf("line is %d wide: %q", lipgloss.Width(l), l)
		}
	}
}

func TestNotesViewScrollsAndCloses(t *testing.T) {
	var changes []string
	for i := 0; i < 40; i++ {
		changes = append(changes, "change "+string(rune('A'+i%26))+string(rune('a'+i/26)))
	}
	v := newNotesView("Release notes", []release.Entry{{Version: "1.0.0", Changes: changes}})

	first := stripANSI(v.render(60, 12))
	if !strings.Contains(first, "1.0.0") || strings.Contains(first, "change Ab") {
		t.Fatalf("first page wrong:\n%s", first)
	}
	if v.key("j") || v.key("down") {
		t.Fatal("scrolling closed the view")
	}
	if v.scroll != 2 {
		t.Fatalf("scroll = %d after two downs, want 2", v.scroll)
	}
	v.key("k")
	v.key("up")
	v.key("up") // stops at the top
	if v.scroll != 0 {
		t.Fatalf("scroll = %d, want 0 at the top", v.scroll)
	}
	v.key("pgdown")
	if v.scroll <= 1 {
		t.Fatalf("pgdown did not page, scroll = %d", v.scroll)
	}
	v.key("pgup")
	if v.scroll != 0 {
		t.Fatalf("pgup scroll = %d, want 0", v.scroll)
	}
	for i := 0; i < 200; i++ {
		v.key("j")
	}
	last := stripANSI(v.render(60, 12))
	if !strings.Contains(last, changes[len(changes)-1]) || v.scroll >= 200 {
		t.Fatalf("scrolling past the end was not clamped (scroll %d):\n%s", v.scroll, last)
	}
	v.key("g")
	if v.scroll != 0 {
		t.Fatalf("g scroll = %d, want 0", v.scroll)
	}
	if got := len(strings.Split(v.render(60, 12), "\n")); got > 12 {
		t.Errorf("render is %d lines, want at most 12", got)
	}
	for _, k := range []string{"esc", "q", "enter"} {
		if !v.key(k) {
			t.Errorf("%s did not close the view", k)
		}
	}
}
