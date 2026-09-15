package milestone

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// write puts a milestone file on disk verbatim, the way a person editing it
// by hand would leave it.
func write(t *testing.T, root, name, body string) {
	t.Helper()
	if err := os.MkdirAll(Dir(root), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(Path(root, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

const idA = "01M2GGKMFB1XXK7V8AFW0YGWXQ"
const idB = "01M2GGKMFB1XXK7V8AFW0S1VM4"
const idC = "01M2GGKMFB1XXK7V8AFW0ZZSFT"

// A milestone file is edited by hand, so a write must give the file back the
// way it was found: comments, blank lines and a chosen order included. This is
// what makes the format the API rather than an implementation detail.
func TestSaveKeepsHandEditsVerbatim(t *testing.T) {
	root := t.TempDir()
	body := "---\n" +
		"name: round-one\n" +
		"colour: 45\n" +
		"created-at: 2026-09-15T10:00:00Z\n" +
		"---\n" +
		"\n" +
		"# round one\n" +
		"\n" +
		"<!-- the two that matter first -->\n" +
		"- " + idC + "  # the tail of the queue\n" +
		"\n" +
		"- " + idA + "\n" +
		"\n" +
		"not a member line at all\n"
	write(t, root, "round-one", body)

	m, err := Load(root, "round-one")
	if err != nil {
		t.Fatal(err)
	}
	if m.Name != "round-one" || m.Colour != 45 {
		t.Errorf("Load() = name %q colour %d, want round-one/45", m.Name, m.Colour)
	}
	if want := "2026-09-15T10:00:00Z"; m.CreatedAt.UTC().Format(time.RFC3339) != want {
		t.Errorf("created-at = %v, want %s", m.CreatedAt, want)
	}
	got := m.Members()
	if len(got) != 2 || got[0] != idC || got[1] != idA {
		t.Errorf("Members() = %v, want file order [%s %s]", got, idC, idA)
	}

	if err := m.Save(root); err != nil {
		t.Fatal(err)
	}
	back, err := os.ReadFile(Path(root, "round-one"))
	if err != nil {
		t.Fatal(err)
	}
	if string(back) != body {
		t.Errorf("Save() rewrote a hand-edited file.\n got:\n%s\nwant:\n%s", back, body)
	}
}

// Add appends one line and Remove drops one line. Everything around them —
// the comment, the blank lines, the order — has to be where it was, because
// twenty tickets grouped in one edit is the whole point of the file.
func TestAddAndRemoveTouchOneLine(t *testing.T) {
	root := t.TempDir()
	write(t, root, "round-two", "---\nname: round-two\ncolour: 33\n---\n\n<!-- keep me -->\n- "+idA+"\n")

	m, err := Load(root, "round-two")
	if err != nil {
		t.Fatal(err)
	}
	if !m.Add(idB) {
		t.Fatal("Add() of a new member reported no change")
	}
	if m.Add(idB) {
		t.Error("Add() of a member already there reported a change")
	}
	if !m.Remove(idA) {
		t.Fatal("Remove() of a member reported no change")
	}
	if m.Remove(idA) {
		t.Error("Remove() of a non-member reported a change")
	}
	if err := m.Save(root); err != nil {
		t.Fatal(err)
	}
	want := "---\nname: round-two\ncolour: 33\n---\n\n<!-- keep me -->\n- " + idB + "\n"
	back, err := os.ReadFile(Path(root, "round-two"))
	if err != nil {
		t.Fatal(err)
	}
	if string(back) != want {
		t.Errorf("after Add/Remove:\n got:\n%s\nwant:\n%s", back, want)
	}
}

// A board with no milestones is the state every board starts in, and one
// unreadable file must not be able to stop the board from opening.
func TestLoadAllToleratesAbsenceAndJunk(t *testing.T) {
	root := t.TempDir()
	all, err := LoadAll(root)
	if err != nil || len(all) != 0 {
		t.Fatalf("LoadAll() on a fresh board = %v, %v; want no milestones and no error", all, err)
	}
	write(t, root, "beta", "---\nname: beta\ncolour: 40\n---\n- "+idA+"\n")
	write(t, root, "alpha", "nothing here is frontmatter\n")
	if err := os.WriteFile(filepath.Join(Dir(root), "notes.txt"), []byte("ignored"), 0o644); err != nil {
		t.Fatal(err)
	}
	all, err = LoadAll(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 || all[0].Name != "alpha" || all[1].Name != "beta" {
		t.Fatalf("LoadAll() = %d milestones, want alpha then beta", len(all))
	}
}

// The index runs the other way from the file, and a ticket may sit in more
// than one milestone — which is exactly what carrying work over looks like
// until somebody removes the old line.
func TestIndexAnswersPerTicketAndAllowsMultipleMembership(t *testing.T) {
	root := t.TempDir()
	write(t, root, "alpha", "---\nname: alpha\ncolour: 33\n---\n- "+idA+"\n- "+idB+"\n")
	write(t, root, "beta", "---\nname: beta\ncolour: 40\n---\n- "+idA+"\n")

	all, err := LoadAll(root)
	if err != nil {
		t.Fatal(err)
	}
	idx := Build(all)
	if got := idx.Names(idA); len(got) != 2 || got[0] != "alpha" || got[1] != "beta" {
		t.Errorf("Names(%s) = %v, want [alpha beta]", idA, got)
	}
	if got := idx.Names(idC); len(got) != 0 {
		t.Errorf("Names(%s) = %v, want nothing", idC, got)
	}
	if !idx.Matches(idB, "Alpha") {
		t.Error("Matches() should normalize the wanted name, like a tag filter does")
	}
	if idx.Matches(idB, "alph") {
		t.Error("Matches() must compare whole names, not substrings")
	}
	if idx.Matches(idB, "") {
		t.Error("Matches() on an impossible name should answer no, not panic")
	}
}

// Nobody picks a milestone colour, and two milestones on one board must not
// share one while the palette still has room.
func TestAssignColourAvoidsTheOnesInUse(t *testing.T) {
	var existing []*Milestone
	for _, c := range Palette[:len(Palette)-1] {
		existing = append(existing, &Milestone{Colour: c})
	}
	got := AssignColour(existing, "last-one")
	if got != Palette[len(Palette)-1] {
		t.Errorf("AssignColour() = %d, want the one free colour %d", got, Palette[len(Palette)-1])
	}
	full := append(existing, &Milestone{Colour: Palette[len(Palette)-1]})
	a := AssignColour(full, "spent")
	b := AssignColour(full, "spent")
	if a != b {
		t.Errorf("an exhausted palette gave %d then %d; a repeat has to be stable per name", a, b)
	}
}

// A milestone jaira just created has to be one it can read back, or the first
// thing anybody does with the feature breaks.
func TestNewRoundTrips(t *testing.T) {
	root := t.TempDir()
	m := New("round-three", 71, time.Date(2026, 9, 15, 12, 0, 0, 0, time.UTC))
	if !m.Add(idA) {
		t.Fatal("Add() on a fresh milestone reported no change")
	}
	if err := m.Save(root); err != nil {
		t.Fatal(err)
	}
	back, err := Load(root, "round-three")
	if err != nil {
		t.Fatal(err)
	}
	if back.Name != "round-three" || back.Colour != 71 || !back.Has(idA) {
		t.Errorf("round trip lost something: %+v members %v", back, back.Members())
	}
}
