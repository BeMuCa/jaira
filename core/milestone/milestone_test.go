package milestone

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// write puts a milestone file on disk verbatim, the way a person editing it
// by hand would leave it, under the ID file name its name gets from LegacyID.
func write(t *testing.T, root, name, body string) {
	t.Helper()
	if err := os.MkdirAll(Dir(root), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(testPath(root, name), []byte(body), 0o644); err != nil {
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
		"color: 45\n" +
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
	back, err := os.ReadFile(testPath(root, "round-one"))
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
	write(t, root, "round-two", "---\nname: round-two\ncolor: 33\n---\n\n<!-- keep me -->\n- "+idA+"\n")

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
	want := "---\nname: round-two\ncolor: 33\n---\n\n<!-- keep me -->\n- " + idB + "\n"
	back, err := os.ReadFile(testPath(root, "round-two"))
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
	write(t, root, "beta", "---\nname: beta\ncolor: 40\n---\n- "+idA+"\n")
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
	write(t, root, "alpha", "---\nname: alpha\ncolor: 33\n---\n- "+idA+"\n- "+idB+"\n")
	write(t, root, "beta", "---\nname: beta\ncolor: 40\n---\n- "+idA+"\n")

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

// The status is one line in the frontmatter, and setting or clearing it must
// leave every other line where it was — the same promise Add and Remove make,
// because the file a person filed is the file they get back on a restore.
func TestSetStatusTouchesOneLineOnly(t *testing.T) {
	root := t.TempDir()
	body := "---\n" +
		"color: 45\n" +
		"<!-- picked by hand -->\n" +
		"created-at: 2026-09-15T10:00:00Z\n" +
		"---\n" +
		"\n" +
		"# round one\n" +
		"\n" +
		"- " + idA + "\n"
	write(t, root, "round-one", body)

	m, err := Load(root, "round-one")
	if err != nil {
		t.Fatal(err)
	}
	if m.Filed() {
		t.Error("a file with no status line reports as filed; no line means on the board")
	}
	m.SetStatus(StatusFiled)
	if err := m.Save(root); err != nil {
		t.Fatal(err)
	}
	back, err := Load(root, "round-one")
	if err != nil {
		t.Fatal(err)
	}
	if !back.Filed() {
		t.Errorf("status did not read back: %q", back.Status)
	}
	if got := back.Members(); len(got) != 1 || got[0] != idA {
		t.Errorf("Members() = %v after SetStatus, want [%s]", got, idA)
	}
	raw, err := os.ReadFile(testPath(root, "round-one"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(raw), "\n") != strings.Count(body, "\n")+1 {
		t.Errorf("SetStatus changed more than one line:\n%s", raw)
	}
	for _, line := range strings.Split(strings.TrimSuffix(body, "\n"), "\n") {
		if !strings.Contains(string(raw), line+"\n") {
			t.Errorf("line %q did not survive SetStatus:\n%s", line, raw)
		}
	}

	// Clearing it puts the file back exactly as it was: restore has to undo
	// filing, not leave a spent line behind.
	back.SetStatus("")
	if err := back.Save(root); err != nil {
		t.Fatal(err)
	}
	raw, err = os.ReadFile(testPath(root, "round-one"))
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != body {
		t.Errorf("clearing the status did not restore the file.\n got:\n%s\nwant:\n%s", raw, body)
	}
}

// A hand-written file is allowed to have no frontmatter at all. It still has
// to be filable, so it is given one rather than losing the status silently.
func TestSetStatusGivesAFileWithoutFrontmatterOne(t *testing.T) {
	root := t.TempDir()
	body := "# round two\n\n- " + idB + "\n"
	write(t, root, "round-two", body)

	m, err := Load(root, "round-two")
	if err != nil {
		t.Fatal(err)
	}
	m.SetStatus(StatusFiled)
	if err := m.Save(root); err != nil {
		t.Fatal(err)
	}
	back, err := Load(root, "round-two")
	if err != nil {
		t.Fatal(err)
	}
	if !back.Filed() {
		t.Errorf("status did not read back from a file that had no frontmatter: %q", back.Status)
	}
	if got := back.Members(); len(got) != 1 || got[0] != idB {
		t.Errorf("Members() = %v, want [%s]", got, idB)
	}
	raw, err := os.ReadFile(testPath(root, "round-two"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(string(raw), body) {
		t.Errorf("the original lines did not survive:\n%s", raw)
	}
}

// The status line is what takes a milestone off the board, not where its file
// happens to lie: a filing that marked the file and got no further, and a
// fetch that wrote a marked file back, both leave one here, and neither is a
// milestone anybody is planning.
func TestLoadAllLeavesOutAFiledMilestone(t *testing.T) {
	root := t.TempDir()
	write(t, root, "open", "---\ncolor: 40\n---\n- "+idA+"\n")
	write(t, root, "closed", "---\ncolor: 33\nstatus: filed\n---\n- "+idB+"\n")

	all, err := LoadAll(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 1 || all[0].Name != "open" {
		t.Fatalf("LoadAll() = %v, want only the open milestone", names(all))
	}
	// Read by name it is still there, which is what a restore needs.
	ms, err := Load(root, "closed")
	if err != nil || !ms.Filed() {
		t.Fatalf("Load(closed) = %v, %v; want the filed milestone", ms, err)
	}
}

func names(all []*Milestone) []string {
	out := make([]string, 0, len(all))
	for _, m := range all {
		out = append(out, m.Name)
	}
	return out
}

// goal and after are the two keys the layer adds. They are read by the same
// frontmatter scan as the rest, so a file written by hand with either must
// parse, and a file with neither must read as "not placed".
func TestParseReadsGoalAndAfter(t *testing.T) {
	m := FromBytes("c", []byte("---\ncolor: 45\ngoal: Version 0.4: out\nafter: [B, a-b , b, ]\n---\n- "+idA+"\n"))
	if m.Goal != "Version 0.4: out" {
		t.Errorf("goal = %q (a colon inside the value must survive)", m.Goal)
	}
	if !m.Placed || strings.Join(m.After, ",") != "b,a-b" {
		t.Errorf("placed=%v after=%v, want normalized, deduplicated [b a-b]", m.Placed, m.After)
	}
	if bare := FromBytes("c", []byte("---\nafter: x, y\n---\n")); strings.Join(bare.After, ",") != "x,y" {
		t.Errorf("no brackets: after = %v", bare.After)
	}
	empty := FromBytes("c", []byte("---\nafter: []\n---\n"))
	if !empty.Placed || len(empty.After) != 0 {
		t.Errorf("after: [] must be placed with no predecessor, got placed=%v after=%v", empty.Placed, empty.After)
	}
	none := FromBytes("c", []byte("---\ncolor: 45\n---\n"))
	if none.Placed || none.Goal != "" {
		t.Errorf("no keys: placed=%v goal=%q", none.Placed, none.Goal)
	}
}

// The setters rewrite one line and leave the rest of the file as it was,
// which is the promise every other edit keeps.
func TestSettersTouchOnlyTheirOwnLine(t *testing.T) {
	body := "---\ncolor: 45\n# keep me\nweird: thing\ncreated-at: 2026-09-15T10:00:00Z\n---\n\n# c\n- " + idA + "\n"
	m := FromBytes("c", []byte(body))

	m.SetGoal("ship it")
	m.SetAfter([]string{"a", "b"})
	want := "---\ncolor: 45\n# keep me\nweird: thing\ncreated-at: 2026-09-15T10:00:00Z\ngoal: ship it\nafter: [a, b]\n---\n\n# c\n- " + idA + "\n"
	if got := string(m.Bytes()); got != want {
		t.Errorf("after set:\n%s\nwant:\n%s", got, want)
	}

	m.SetGoal("changed")
	m.SetAfter(nil)
	if got := string(m.Bytes()); !strings.Contains(got, "goal: changed\nafter: []\n") || strings.Count(got, "goal:") != 1 {
		t.Errorf("replace in place failed:\n%s", got)
	}
	if !m.Placed || len(m.After) != 0 {
		t.Errorf("SetAfter(nil) must place as a lane top, placed=%v", m.Placed)
	}

	m.SetGoal("")
	m.Unplace()
	if got := string(m.Bytes()); got != body {
		t.Errorf("removing both must give the original back:\n%s", got)
	}
	if m.Placed || m.Goal != "" {
		t.Errorf("fields not cleared: placed=%v goal=%q", m.Placed, m.Goal)
	}
}

// The layer's c rewrites the colour where it stands and nothing else.
func TestSetColourRewritesOnlyTheColourLine(t *testing.T) {
	body := "---\ncolor: 45\n# keep me\ncreated-at: 2026-09-15T10:00:00Z\n---\n\n# c\n- " + idA + "\n"
	m := FromBytes("c", []byte(body))
	m.SetColour(209)
	if got, want := string(m.Bytes()), strings.Replace(body, "color: 45", "color: 209", 1); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
	if m.Colour != 209 {
		t.Errorf("Colour = %d, want 209", m.Colour)
	}
}

// A file that began without frontmatter is given one rather than refused.
func TestSetGoalOnAFileWithoutFrontmatter(t *testing.T) {
	m := FromBytes("c", []byte("- "+idA+"\n"))
	m.SetGoal("g")
	if got := string(m.Bytes()); got != "---\ngoal: g\n---\n\n- "+idA+"\n" {
		t.Errorf("got:\n%s", got)
	}
}

// What a setter wrote reads back the same, through the file.
func TestGoalAndAfterRoundTripThroughSave(t *testing.T) {
	root := t.TempDir()
	m := New("c", 45, time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC))
	m.SetGoal("a: b")
	m.SetAfter([]string{"x"})
	if err := m.Save(root); err != nil {
		t.Fatal(err)
	}
	got, err := Load(root, "c")
	if err != nil {
		t.Fatal(err)
	}
	if got.Goal != "a: b" || !got.Placed || strings.Join(got.After, ",") != "x" {
		t.Errorf("round trip: goal=%q placed=%v after=%v", got.Goal, got.Placed, got.After)
	}
}

// ParseNames is what --after and the layer's form read: names stored the way
// milestone names are, repeats dropped, order kept, and a name that cannot be
// one refused rather than quietly left out.
func TestParseNames(t *testing.T) {
	got, err := ParseNames(" API Umbau, ui,,api-umbau ")
	if err != nil || strings.Join(got, ",") != "api-umbau,ui" {
		t.Errorf("got %v, %v", got, err)
	}
	if got, err := ParseNames(""); err != nil || len(got) != 0 {
		t.Errorf("empty: %v, %v", got, err)
	}
	if _, err := ParseNames("ok, ???"); err == nil {
		t.Error("a name that does not normalize must be refused")
	}
}

// testPath is where write puts a milestone: <name>-<LegacyID(name)>.md.
func testPath(root, name string) string {
	return filepath.Join(Dir(root), name+"-"+LegacyID(name)+".md")
}

func writeFile(t *testing.T, root, file, body string) {
	t.Helper()
	if err := os.MkdirAll(Dir(root), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(Dir(root), file), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// A file with an ID in its name and one from before IDs sit side by side, and
// both are read: the old one under the ID computed from its name, untouched.
func TestLoadAllReadsIDFilesAndOldFilesAndNeverWritesWhileReading(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "backend-K3F9QZ.md", "---\ncolor: 45\n---\n")
	writeFile(t, root, "q4.md", "---\ncolor: 46\n---\n")
	all, err := LoadAll(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 || all[0].Name != "backend" || all[0].ID != "K3F9QZ" || all[1].Name != "q4" || all[1].ID != LegacyID("q4") {
		t.Fatalf("LoadAll = %+v", all)
	}
	if all[0].FileName() != "backend-K3F9QZ.md" || all[1].FileName() != "q4.md" {
		t.Errorf("file names %q, %q", all[0].FileName(), all[1].FileName())
	}
	entries, _ := os.ReadDir(Dir(root))
	if len(entries) != 2 {
		t.Errorf("reading changed the directory: %v", entries)
	}
}

// Two files with one ID are one milestone: the first by file name counts.
func TestTwoFilesWithOneIDAreOneMilestone(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "alpha-K3F9QZ.md", "---\ncolor: 45\n---\n")
	writeFile(t, root, "beta-K3F9QZ.md", "---\ncolor: 46\n---\n")
	all, _ := LoadAll(root)
	if len(all) != 1 || all[0].Name != "alpha" {
		t.Fatalf("LoadAll = %+v", all)
	}
}

func TestLoadResolvesIDThenNameIDThenName(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "backend-K3F9QZ.md", "---\ncolor: 45\n---\n")
	writeFile(t, root, "backend-M7P2RT.md", "---\ncolor: 46\n---\n")
	writeFile(t, root, "ui-B4C5D6.md", "---\ncolor: 47\n---\n")
	writeFile(t, root, "old-filed.md", "---\nstatus: filed\n---\n")
	for key, want := range map[string]string{
		"K3F9QZ": "backend-K3F9QZ.md", "k3f9qz": "backend-K3F9QZ.md",
		"backend-M7P2RT": "backend-M7P2RT.md", "backend-m7p2rt": "backend-M7P2RT.md",
		"ui": "ui-B4C5D6.md", "UI": "ui-B4C5D6.md",
		"old-filed": "old-filed.md", // a filed milestone is loadable: restore needs it
	} {
		m, err := Load(root, key)
		if err != nil || m.FileName() != want {
			t.Errorf("Load(%q) = %v, %v; want %s", key, m, err, want)
		}
	}
	// a name and the wrong ID is nothing, not the name's milestone
	if _, err := Load(root, "ui-K3F9QZ"); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Load(ui-K3F9QZ) err = %v, want not-exist", err)
	}
	if _, err := Load(root, "nothing"); !errors.Is(err, os.ErrNotExist) || !os.IsNotExist(err) {
		t.Errorf("Load(nothing) err = %v, want not-exist", err)
	}
	_, err := Load(root, "backend")
	var amb *AmbiguousError
	if !errors.As(err, &amb) || amb.Name != "backend" || strings.Join(amb.Choices, " ") != "backend-K3F9QZ backend-M7P2RT" {
		t.Errorf("Load(backend) err = %v, want AmbiguousError naming both", err)
	}
}

// A name on the board wins over the same name on a file that is filed.
func TestLoadPrefersTheMilestoneOnTheBoardOverAFiledOne(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "backend-AAAAAA.md", "---\nstatus: filed\n---\n")
	writeFile(t, root, "backend-BBBBBB.md", "---\ncolor: 45\n---\n")
	m, err := Load(root, "backend")
	if err != nil || m.ID != "BBBBBB" {
		t.Fatalf("Load(backend) = %v, %v", m, err)
	}
}

func TestFromBytesTakesNameAndIDFromTheFileName(t *testing.T) {
	m := FromBytes("release-stages-K3F9QZ.md", []byte("---\ncolor: 45\n---\n"))
	if m.Name != "release-stages" || m.ID != "K3F9QZ" {
		t.Errorf("got %q %q", m.Name, m.ID)
	}
	if m = FromBytes("sprint-202611.md", nil); m.Name != "sprint-202611" || m.ID != LegacyID("sprint-202611") {
		t.Errorf("got %q %q", m.Name, m.ID)
	}
}

// A hand-named file keeps its spelling, and is found under it.
func TestLoadFindsAHandNamedFileAsItStands(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "Q4.md", "---\ncolor: 45\n---\n")
	if m, err := Load(root, "Q4"); err != nil || m.Name != "Q4" {
		t.Errorf("Load(Q4) = %v, %v", m, err)
	}
}

// Twins are the names the board holds more than once, in file order, each with
// the labels that tell its milestones apart and the command that renames one.
func TestTwinsListsEveryNameHeldTwice(t *testing.T) {
	mk := func(name, id string) *Milestone {
		m := New(name, 0, time.Now())
		m.ID = id
		return m
	}
	all := []*Milestone{mk("backend", "23G8MG"), mk("solo", "AAAAAA"), mk("backend", "4K2ZZ9"), mk("ui", "BBBBBB"), mk("ui", "CCCCCC"), mk("backend", "DDDDD1")}
	got := Twins(all)
	if len(got) != 2 {
		t.Fatalf("Twins = %+v, want backend and ui", got)
	}
	if got[0].Name != "backend" || strings.Join(got[0].Labels, " ") != "backend-23G8MG backend-4K2ZZ9 backend-DDDDD1" {
		t.Errorf("first twin = %+v", got[0])
	}
	if want := "jaira milestone rename backend-23G8MG <new-name>"; got[0].RenameCommand() != want {
		t.Errorf("RenameCommand = %q, want %q", got[0].RenameCommand(), want)
	}
	if got[1].Name != "ui" {
		t.Errorf("second twin = %+v", got[1])
	}
	if Twins(all[:2]) != nil || Twins(nil) != nil {
		t.Errorf("no twins must give nil")
	}
}
