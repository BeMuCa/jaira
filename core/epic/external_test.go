package epic

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// The keys are a block list in the frontmatter, one per line; adding,
// relinking and taking one out touch that block and nothing else, and the last
// key going takes the block with it, so the file reads as it did before.
func TestExternalKeysAreABlockListInTheFrontmatter(t *testing.T) {
	s := &ticket.Store{Root: t.TempDir()}
	refs := mapRefs{}
	if _, err := Create(s, refs, "f", "F", "d", []string{idA}); err != nil {
		t.Fatal(err)
	}
	plain := string(refs["f"])
	if strings.Contains(plain, "external") {
		t.Fatalf("an epic without keys wrote the field:\n%s", plain)
	}

	jira := External{Key: "JIRA-1", URL: "https://jira.example/browse/JIRA-1"}
	e, changed, err := Set(s, refs, "f", Patch{AddExternal: []External{jira, {Key: "YT-7"}}})
	if err != nil || !changed {
		t.Fatalf("add: changed %v, %v", changed, err)
	}
	want := "external:\n  - JIRA-1 https://jira.example/browse/JIRA-1\n  - YT-7\n---\n"
	if got := string(refs["f"]); !strings.Contains(got, want) || !strings.Contains(got, "- "+idA) {
		t.Fatalf("file:\n%s\nwant the block %q", got, want)
	}
	if got := FromBytes("f", refs["f"]).Externals(); len(got) != 2 || got[0] != jira || got[1].Key != "YT-7" || got[1].URL != "" {
		t.Fatalf("read back = %+v", got)
	}
	if !e.HasExternal("jira-1") {
		t.Error("a key is matched in any case")
	}

	// The same key again, in another case, relinks it in its place.
	relinked := External{Key: "jira-1", URL: "https://jira.example/JIRA-1"}
	if _, changed, _ := Set(s, refs, "f", Patch{AddExternal: []External{relinked}}); !changed {
		t.Fatal("relink wrote nothing")
	}
	if got := FromBytes("f", refs["f"]).Externals(); len(got) != 2 || got[0] != relinked {
		t.Fatalf("after relink = %+v", got)
	}
	if _, changed, _ := Set(s, refs, "f", Patch{AddExternal: []External{relinked}}); changed {
		t.Error("the same key and URL again wrote")
	}

	if _, changed, _ := Set(s, refs, "f", Patch{DropExternal: []string{"JIRA-1", "yt-7"}}); !changed {
		t.Fatal("drop wrote nothing")
	}
	if got := string(refs["f"]); got != plain {
		t.Errorf("after the last key went:\n%s\nwant it as it was:\n%s", got, plain)
	}
	if _, changed, _ := Set(s, refs, "f", Patch{DropExternal: []string{"JIRA-1"}}); changed {
		t.Error("dropping a key that is not there wrote")
	}
}

// A file a person edited stays theirs: a quoted line, a key typed on the field's
// own line and a line that is no URL are all read, and editing one key keeps
// every other line where it was.
func TestExternalKeysReadAHandEditedFile(t *testing.T) {
	src := "---\ntitle: T\nexternal: A-1 https://a.example/1\n  - \"B-2 https://b.example/2\"\n  - C-3 not a link\ndescription: d\n---\n\n- " + idA + "\n"
	e := FromBytes("f", []byte(src))
	got := e.Externals()
	if len(got) != 3 || got[0].Key != "A-1" || got[1].URL != "https://b.example/2" || got[2].URL != "not a link" {
		t.Fatalf("read = %+v", got)
	}
	if e.Description != "d" || !e.Has(idA) {
		t.Fatalf("the lines around the list: %q %v", e.Description, e.Members())
	}
	e.dropExternal("B-2")
	e.writeExternal()
	want := "---\ntitle: T\nexternal:\n  - A-1 https://a.example/1\n  - C-3 not a link\ndescription: d\n---\n\n- " + idA + "\n"
	if got := string(e.Bytes()); got != want {
		t.Errorf("file:\n%s\nwant:\n%s", got, want)
	}
}

// What a person types is checked: one word for the key, and a link only when
// it is one a page may open.
func TestParseExternal(t *testing.T) {
	for raw, want := range map[string]External{
		"JIRA-1 https://j.example/JIRA-1": {Key: "JIRA-1", URL: "https://j.example/JIRA-1"},
		"  YT-7  ":                        {Key: "YT-7"},
	} {
		if got, err := ParseExternal(raw); err != nil || got != want {
			t.Errorf("%q = %+v, %v", raw, got, err)
		}
	}
	for _, raw := range []string{"", "A B C", "A javascript:alert(1)", "A ftp://x", "A https://", "-A https://x.example"} {
		if _, err := ParseExternal(raw); !errors.Is(err, ErrBadExternal) {
			t.Errorf("%q accepted: %v", raw, err)
		}
	}
}

// Two clones change the keys at once: each keeps what it did.
func TestMergeKeepsBothSidesKeyChanges(t *testing.T) {
	base := New("f", "F", "", time.Unix(0, 0))
	base.putExternal(External{Key: "A-1"})
	base.putExternal(External{Key: "B-2"})
	base.writeExternal()
	theirs := FromBytes("f", base.Bytes())
	theirs.putExternal(External{Key: "C-3"})
	theirs.writeExternal()
	ours := FromBytes("f", base.Bytes())
	ours.dropExternal("A-1")
	ours.putExternal(External{Key: "B-2", URL: "https://b.example"})
	ours.putExternal(External{Key: "D-4"})
	ours.writeExternal()

	got := FromBytes("f", Merge("f", base.Bytes(), ours.Bytes(), theirs.Bytes())).Externals()
	want := []External{{Key: "B-2", URL: "https://b.example"}, {Key: "C-3"}, {Key: "D-4"}}
	if len(got) != len(want) {
		t.Fatalf("merged = %+v; want %+v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("merged = %+v; want %+v", got, want)
		}
	}
}
