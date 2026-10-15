package tui

import (
	"os"
	"path/filepath"
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

func TestNotesBase(t *testing.T) {
	for _, c := range []struct {
		name, current, seen, stamped, want string
	}{
		{"dev build", "dev", "0.5.0", "0.5.0", ""},
		{"seen equals current", "0.5.1", "0.5.1", "0.4.0", ""},
		{"seen is older", "0.5.1", "0.5.0", "0.5.1", "0.5.0"},
		{"seen wins over an older stamp", "0.5.1", "0.5.0", "0.3.0", "0.5.0"},
		{"no seen, old stamp", "0.5.1", "", "0.5.0", "0.5.0"},
		{"no seen, stamp equals current", "0.5.1", "", "0.5.1", ""},
		{"neither", "0.5.1", "", "", ""},
		{"stamp is dev", "0.5.1", "", "dev", ""},
		{"seen is dev", "0.5.1", "dev", "0.5.0", ""},
		{"downgrade", "0.5.0", "0.5.1", "0.5.1", ""},
		{"numbers, not text", "0.10.0", "0.9.0", "", "0.9.0"},
	} {
		if got := notesBase(c.current, c.seen, c.stamped); got != c.want {
			t.Errorf("%s: notesBase(%q, %q, %q) = %q, want %q", c.name, c.current, c.seen, c.stamped, got, c.want)
		}
	}
}

// popupBoard builds a board whose state directory holds the given stamp and
// notes-seen record ("" leaves the file out) and opens it as the real
// constructor does, under a binary of the given version.
func popupBoard(t *testing.T, current, stamp, seen string) *Model {
	t.Helper()
	setReleaseCurrent(t, current)
	s := newTestStore(t)
	dir := s.StateDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, v := range map[string]string{"version": stamp, seenFile: seen} {
		if v != "" {
			if err := os.WriteFile(filepath.Join(dir, name), []byte(v+"\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	m, err := New(s)
	if err != nil {
		t.Fatal(err)
	}
	m.width, m.height = 120, 40
	return m
}

func TestPopupAfterUpgradeShowsNotesSinceTheStamp(t *testing.T) {
	m := popupBoard(t, "0.5.1", "0.5.0", "")
	if m.mode != modeReleaseNotes || m.notes == nil {
		t.Fatalf("mode = %v, want the release-notes popup", m.mode)
	}
	out := stripANSI(m.render())
	if !strings.Contains(out, "What is new since 0.5.0") || !strings.Contains(out, "0.5.1") {
		t.Errorf("popup does not name the versions:\n%s", out)
	}
}

func TestPopupAfterUpdateReadsTheSeenRecord(t *testing.T) {
	// 'jaira update' already stamped the board with the running version; the
	// seen record is what remembers where the user stopped reading.
	m := popupBoard(t, "0.5.1", "0.5.1", "0.5.0")
	if m.mode != modeReleaseNotes {
		t.Fatalf("mode = %v, want the release-notes popup", m.mode)
	}
	if !strings.Contains(stripANSI(m.render()), "What is new since 0.5.0") {
		t.Errorf("popup should start at the seen version:\n%s", stripANSI(m.render()))
	}
}

func TestPopupClosesAndStaysClosedUntilTheNextVersion(t *testing.T) {
	m := popupBoard(t, "0.5.1", "0.5.0", "")
	dir := m.store.StateDir()

	press(m, "esc")
	if m.mode != modeBoard || m.notes != nil {
		t.Fatalf("esc left mode %v, notes %v", m.mode, m.notes)
	}
	if got := readSeen(dir); got != "0.5.1" {
		t.Fatalf("notes-seen = %q, want 0.5.1", got)
	}
	again, err := New(m.store)
	if err != nil {
		t.Fatal(err)
	}
	if again.mode == modeReleaseNotes {
		t.Fatal("the popup came back for the same version")
	}

	setReleaseCurrent(t, "0.5.2")
	next, err := New(m.store)
	if err != nil {
		t.Fatal(err)
	}
	if next.mode != modeReleaseNotes {
		t.Fatal("the next version did not show the popup")
	}
}

func TestNoPopupOnAFreshBoardOrDevBuild(t *testing.T) {
	// A board nothing ever stamped has no "before" to show changes since.
	fresh := popupBoard(t, "0.5.1", "", "")
	if fresh.mode == modeReleaseNotes {
		t.Error("a fresh board opened the popup")
	}
	if got := readSeen(fresh.store.StateDir()); got != "0.5.1" {
		t.Errorf("a fresh board should be marked seen, notes-seen = %q", got)
	}

	dev := popupBoard(t, "dev", "0.5.0", "0.4.0")
	if dev.mode == modeReleaseNotes {
		t.Error("a dev build opened the popup")
	}
	if got := readSeen(dev.store.StateDir()); got != "0.4.0" {
		t.Errorf("a dev build must not write notes-seen, got %q", got)
	}
}
