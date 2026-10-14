package milestone

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const adoptBase = "---\ncolor: 45\nafter: []\n---\n\n# alpha\n\n"

func adopt(base, old, cur string, all ...*Milestone) (string, []Conflict) {
	out, c := Adopt("alpha", []byte(base), []byte(old), []byte(cur), all)
	return string(out), c
}

func TestAdoptTakesWhatOnlyTheOldSideChanged(t *testing.T) {
	base := adoptBase + "- " + idA + "\n"
	old := strings.Replace(base, "after: []\n", "after: []\ngoal: theirs\n", 1) + "- " + idB + "\n"
	got, c := adopt(base, old, base)
	if len(c) != 0 {
		t.Fatalf("conflicts %v", c)
	}
	ms := FromBytes("alpha-HV9ZDB.md", []byte(got))
	if ms.Goal != "theirs" || !ms.Has(idA) || !ms.Has(idB) {
		t.Errorf("goal %q members %v\n%s", ms.Goal, ms.Members(), got)
	}
}

func TestAdoptAppliesTicketAddsAndRemovesToThisSide(t *testing.T) {
	base := adoptBase + "- " + idA + "\n"
	old := adoptBase + "- " + idB + "\n"                     // old removed A, added B
	cur := base + "- " + "01M2GGKMFB1XXK7V8AFW0YGWX3" + "\n" // this side added a third
	got, _ := adopt(base, old, cur)
	ms := FromBytes("alpha-HV9ZDB.md", []byte(got))
	if ms.Has(idA) || !ms.Has(idB) || !ms.Has("01M2GGKMFB1XXK7V8AFW0YGWX3") {
		t.Errorf("members %v", ms.Members())
	}
}

func TestAdoptKeepsThisSideOnAClashAndSaysWhatWasLeftOut(t *testing.T) {
	base := adoptBase
	old := strings.Replace(base, "after: []\n", "after: []\ngoal: theirs\n", 1)
	cur := strings.Replace(base, "after: []\n", "after: []\ngoal: ours\n", 1)
	got, c := adopt(base, old, cur)
	if len(c) != 1 || c[0].Field != "goal" || c[0].Old != `"theirs"` {
		t.Fatalf("conflicts %+v", c)
	}
	ms := FromBytes("alpha-HV9ZDB.md", []byte(got))
	if ms.Goal != "ours" || !strings.Contains(got, "theirs") {
		t.Errorf("goal %q, the old value must be written into the file:\n%s", ms.Goal, got)
	}
	// the same change on both sides is no clash
	if _, c := adopt(base, old, old); len(c) != 0 {
		t.Errorf("identical change reported: %v", c)
	}
}

func TestAdoptTakesAnOldBodyOnlyIfThisSideLeftItAlone(t *testing.T) {
	base := adoptBase + "- " + idA + "\n"
	old := strings.Replace(base, "# alpha\n", "# alpha\n\nold note\n", 1)
	got, c := adopt(base, old, base+"- "+idB+"\n")
	if len(c) != 0 || !strings.Contains(got, "old note") || !strings.Contains(got, idB) {
		t.Errorf("conflicts %v\n%s", c, got)
	}
	cur := strings.Replace(base, "# alpha\n", "# alpha\n\nour note\n", 1)
	got, c = adopt(base, old, cur)
	if len(c) != 1 || c[0].Field != "body" || strings.Contains(got, "\nold note\n") && !strings.Contains(got, "<!--") {
		t.Errorf("a body both sides changed must keep ours and say so: %v\n%s", c, got)
	}
	if !strings.Contains(got, "our note") {
		t.Errorf("our body was lost:\n%s", got)
	}
}

func TestAdoptSpellsAnOldAfterAsIDs(t *testing.T) {
	beta := FromBytes("beta-"+LegacyID("beta")+".md", []byte("---\nafter: []\n---\n"))
	old := strings.Replace(adoptBase, "after: []", "after: [beta]", 1)
	got, c := adopt(adoptBase, old, adoptBase, beta)
	if ms := FromBytes("alpha-HV9ZDB.md", []byte(got)); len(c) != 0 || strings.Join(ms.After, ",") != beta.ID {
		t.Errorf("after %v conflicts %v", ms.After, c)
	}
	// an old file that still holds names equals a base that was migrated
	migrated := string(MigrateContent([]byte(old)))
	if got, c := adopt(migrated, old, migrated); len(c) != 0 || got != migrated {
		t.Errorf("names against their IDs differ: %v\n%s", c, got)
	}
}

func TestAdoptChangesNothingWhenThereIsNothingToAdopt(t *testing.T) {
	cur := adoptBase + "- " + idA + "\n"
	if got, c := adopt(adoptBase, adoptBase, cur); got != cur || c != nil {
		t.Errorf("changed:\n%s", got)
	}
}

// Migrate records the base, and an old file turning up beside the ID file is
// merged into it instead of being removed.
func TestMigrateRecordsABaseAndAdoptsALaterOldFile(t *testing.T) {
	root := t.TempDir()
	rec := &fakeRec{}
	writeOld(t, root, "alpha", adoptBase)
	if _, err := Migrate(root, rec); err != nil {
		t.Fatal(err)
	}
	id := LegacyID("alpha")
	if b, ok := LoadBase(root, id); !ok || b.Content != adoptBase {
		t.Fatalf("base %+v %v", b, ok)
	}
	rec.names = nil
	writeOld(t, root, "alpha", adoptBase+"- "+idA+"\n")
	res, err := Migrate(root, rec)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(filepath.Join(Dir(root), "alpha-"+id+".md"))
	if !strings.Contains(string(b), idA) || len(res.Adopted) != 1 || len(rec.names) != 1 {
		t.Errorf("adopted %v, queued %v\n%s", res.Adopted, rec.names, b)
	}
	if dirNames(t, root) != "alpha-"+id+".md" {
		t.Errorf("directory %s", dirNames(t, root))
	}
	// the same old file again: nothing to adopt, nothing queued
	rec.names = nil
	writeOld(t, root, "alpha", adoptBase+"- "+idA+"\n")
	if res, _ := Migrate(root, rec); len(res.Adopted) != 0 || len(rec.names) != 0 {
		t.Errorf("adopted twice: %v %v", res.Adopted, rec.names)
	}
}

// One file under two spellings (a case-insensitive file system) is never
// removed as if it were the old one.
func TestMigrateLeavesAnOldNameThatIsTheIDFileItself(t *testing.T) {
	root := t.TempDir()
	id := LegacyID("alpha")
	target := filepath.Join(Dir(root), "alpha-"+id+".md")
	writeFile(t, root, "alpha-"+id+".md", adoptBase)
	if err := os.Link(target, filepath.Join(Dir(root), "alpha.md")); err != nil {
		t.Skip("no hard links here")
	}
	if _, err := Migrate(root, &fakeRec{}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(Dir(root), "alpha.md")); err != nil {
		t.Errorf("the old name was removed although it is the ID file: %v", err)
	}
}

func TestAdoptWithoutABaseUnionsAndKeepsThisSideOnAClash(t *testing.T) {
	cur := strings.Replace(adoptBase, "after: []\n", "after: []\ngoal: ours\n", 1) + "- " + idA + "\n"
	old := strings.Replace(adoptBase, "after: []\n", "after: []\ngoal: theirs\nowner: o\n", 1) + "- " + idB + "\n"
	got, c := adopt("", old, cur)
	ms := FromBytes("alpha-HV9ZDB.md", []byte(got))
	if ms.Goal != "ours" || !ms.Has(idA) || !ms.Has(idB) || !strings.Contains(got, "owner: o") {
		t.Errorf("goal %q members %v\n%s", ms.Goal, ms.Members(), got)
	}
	if len(c) != 1 || c[0].Field != "goal" || c[0].Old != `"theirs"` {
		t.Errorf("conflicts %+v", c)
	}
}

func conflictLines(s string) []string {
	var out []string
	for _, l := range strings.Split(s, "\n") {
		if strings.HasPrefix(l, "<!--") {
			out = append(out, l)
		}
	}
	return out
}

// The comment a clash leaves in the file is not a field and not body: it never
// makes a later merge see a change, and a second clash on the same field
// replaces its line instead of stacking another.
func TestAdoptCommentIsIgnoredByLaterMergesAndNotStacked(t *testing.T) {
	base := adoptBase
	cur := strings.Replace(base, "after: []\n", "after: []\ngoal: ours\n", 1)
	old1 := strings.Replace(base, "after: []\n", "after: []\ngoal: first\n", 1)
	once, c := adopt(base, old1, cur)
	if len(c) != 1 || len(conflictLines(once)) != 1 {
		t.Fatalf("first clash: %v\n%s", c, once)
	}
	if ms := FromBytes("alpha-HV9ZDB.md", []byte(once)); ms.Goal != "ours" {
		t.Errorf("goal %q", ms.Goal)
	}
	// the old side did not move: nothing to adopt, the file is untouched
	if again, c := adopt(old1, old1, once); again != once || len(c) != 0 {
		t.Errorf("an unchanged old side changed the file or conflicted: %v\n%s", c, again)
	}
	// a body only the old side changed is still taken, with the comment kept
	body := strings.Replace(old1, "# alpha", "# alpha\n\nnew text", 1)
	got, c := adopt(old1, body, once)
	if len(c) != 0 || !strings.Contains(got, "new text") || len(conflictLines(got)) != 1 {
		t.Errorf("body not taken or comment lost: %v\n%s", c, got)
	}
	// a second clash on the same field replaces the line
	old2 := strings.Replace(base, "after: []\n", "after: []\ngoal: second\n", 1)
	twice, c := adopt(old1, old2, once)
	lines := conflictLines(twice)
	if len(c) != 1 || len(lines) != 1 || !strings.Contains(lines[0], "second") || strings.Contains(lines[0], "first") {
		t.Errorf("second clash stacked or kept the old value: %v\n%s", c, twice)
	}
}

// A merge that takes the body from the old side, which already carries the
// clash comments this side has too, writes each comment once.
func TestAdoptTakingABodyFromASideWithCommentsDoesNotDoubleThem(t *testing.T) {
	note := "<!-- another edit changed goal too and its value was not taken: \"x\" -->\n"
	base := adoptBase + "text one\n"
	old := adoptBase + "text two\n" + note
	cur := base + note
	got, _ := adopt(base, old, cur)
	if n := strings.Count(got, note); n != 1 {
		t.Errorf("the comment is in the file %d times:\n%s", n, got)
	}
	if !strings.Contains(got, "text two") {
		t.Errorf("the old body was not taken:\n%s", got)
	}
	// a different comment for the same field on the old side does not stack
	other := "<!-- another edit changed goal too and its value was not taken: \"y\" -->\n"
	got, _ = adopt(base, adoptBase+"text two\n"+other, cur)
	if n := strings.Count(got, "changed goal too"); n != 1 {
		t.Errorf("%d comments for one field:\n%s", n, got)
	}
}
