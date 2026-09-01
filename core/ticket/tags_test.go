package ticket

import (
	"strings"
	"testing"
)

// tags is a plain list field, which is the whole design decision: everything a
// list field already does — 'jaira set tags=a,b', the merge driver's union, the
// one-line rewrite — comes for free, and none of it works if the field is a
// comma-joined scalar instead.
func TestTagsRoundTripAsAListField(t *testing.T) {
	d := NewDoc(
		map[string]string{FieldID: "01KZTT3XZ2YQBX93TTSR7BVRCT", FieldStatus: "todo", FieldTitle: "t"},
		map[string][]string{FieldTags: {"ui", "backend-api"}},
		"# t\n")

	tk, err := Decode(d, "t.md")
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(tk.Tags, ","); got != "ui,backend-api" {
		t.Fatalf("Tags = %v, want [ui backend-api]", tk.Tags)
	}

	raw := string(d.Bytes())
	if !strings.Contains(raw, "tags:\n  - ui\n  - backend-api\n") {
		t.Errorf("tags is not written as a YAML list:\n%s", raw)
	}
	// Order in a new file: beside blocked-by, the other list the board reads a
	// ticket's relationships from — not appended after the timestamps.
	if strings.Index(raw, "\ntags:") > strings.Index(raw, "\nblocked-by:") && strings.Contains(raw, "\nblocked-by:") {
		t.Errorf("tags is written after blocked-by:\n%s", raw)
	}

	// Rewriting the list touches that key's block and nothing else.
	if err := d.SetList(FieldTags, []string{"docs"}); err != nil {
		t.Fatal(err)
	}
	again, err := ParseDoc(d.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	tk2, err := Decode(again, "t.md")
	if err != nil {
		t.Fatal(err)
	}
	if len(tk2.Tags) != 1 || tk2.Tags[0] != "docs" {
		t.Errorf("Tags after SetList = %v, want [docs]", tk2.Tags)
	}
	if tk2.ID != tk.ID || tk2.Title != tk.Title || tk2.Status != tk.Status {
		t.Error("rewriting tags disturbed another field")
	}
}

// A ticket written before tags existed decodes with none, rather than failing.
func TestTicketWithoutTagsDecodes(t *testing.T) {
	d, err := ParseDoc([]byte("---\nid: 01KZTT3XZ2YQBX93TTSR7BVRCT\nstatus: todo\n---\n\n# t\n"))
	if err != nil {
		t.Fatal(err)
	}
	tk, err := Decode(d, "t.md")
	if err != nil {
		t.Fatal(err)
	}
	if len(tk.Tags) != 0 {
		t.Errorf("Tags = %v, want none", tk.Tags)
	}
}
