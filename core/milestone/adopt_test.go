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
