package cli

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// linkStore builds a board with an epic, a child, and a blocker that has been
// finished and filed into the logbook — the arrangement every link question
// here is about.
func linkStore(t *testing.T) (dir, epic, child, filed string) {
	t.Helper()
	dir = t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	mk := func(title, status string, extra map[string]string) string {
		f := map[string]string{
			ticket.FieldID:     ticket.NewID(now),
			ticket.FieldTitle:  title,
			ticket.FieldStatus: status,
		}
		for k, v := range extra {
			f[k] = v
		}
		tk, err := s.Create(f, nil, "")
		if err != nil {
			t.Fatal(err)
		}
		now = now.Add(time.Millisecond)
		return tk.ID
	}
	epic = mk("The epic", "todo", nil)
	child = mk("The child", "todo", map[string]string{ticket.FieldParent: epic})
	filed = mk("The finished blocker", "done", nil)
	if _, err := s.Logbook(filed, "as-20260911"); err != nil {
		t.Fatal(err)
	}
	return dir, epic, child, filed
}

// A handle is what every command prints, so a handle is what gets typed. It
// has to be stored as the full id, or the link resolves to nothing for good.
func TestSetResolvesAHandleToTheFullID(t *testing.T) {
	dir, epic, child, _ := linkStore(t)

	if _, err := runCLI(t, dir, "set", ticket.Handle(epic), "related="+ticket.Handle(child)); err != nil {
		t.Fatal(err)
	}

	out, err := runCLI(t, dir, "--json", "show", epic)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatal(err)
	}
	rel, _ := got["related"].([]any)
	if len(rel) != 1 || rel[0] != child {
		t.Errorf("related = %#v, want the full id %s", got["related"], child)
	}
}

// A reference that names nothing is refused rather than written: a dead link
// is worse than no link, because it looks like a trail exists.
func TestSetRefusesAReferenceThatNamesNothing(t *testing.T) {
	dir, epic, _, _ := linkStore(t)

	out, err := runCLI(t, dir, "set", epic, "parent=ZZZZZZ")
	if err == nil {
		t.Fatalf("expected a refusal, got none: %s", out)
	}
	if !strings.Contains(out+err.Error(), "ZZZZZZ") {
		t.Errorf("the refusal must name the reference, got %q / %v", out, err)
	}
}

// A parent that has been finished and filed is still a parent. Refusing to
// file a child under it would make the link useless exactly when the work
// grows past one ticket.
func TestCreateAcceptsAFiledParent(t *testing.T) {
	dir, _, _, filed := linkStore(t)

	if _, err := runCLI(t, dir, "create", "A late child",
		"--parent", ticket.Handle(filed), "--goal", "g", "--context", "c", "--dod", "d"); err != nil {
		t.Fatalf("a filed parent must be accepted: %v", err)
	}
}

// links reads past the board: the blocker is in the logbook and the reader
// has to be told that, not told it does not exist.
func TestLinksShowsWhereEachEndLives(t *testing.T) {
	dir, epic, _, filed := linkStore(t)
	if _, err := runCLI(t, dir, "set", epic, "blocked-by="+filed); err != nil {
		t.Fatal(err)
	}

	out, err := runCLI(t, dir, "links", epic)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"waiting on", "The finished blocker", "logbook · done",
		"contains", "The child", "on the board · todo",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("links is missing %q\n%s", want, out)
		}
	}
}

func TestLinksJSONGroupsByKind(t *testing.T) {
	dir, epic, child, _ := linkStore(t)

	out, err := runCLI(t, dir, "--json", "links", epic)
	if err != nil {
		t.Fatal(err)
	}
	var got struct {
		Links map[string][]struct {
			ID    string `json:"id"`
			Place string `json:"place"`
		} `json:"links"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatal(err)
	}
	kids := got.Links["child"]
	if len(kids) != 1 || kids[0].ID != child || kids[0].Place != "board" {
		t.Errorf("child group = %#v, want %s on the board", kids, child)
	}
}

// show used to answer "not found" for a file sitting in plain sight under
// .jaira/logbook/. It now finds it and says where it is.
func TestShowFindsAFiledTicketAndSaysWhereItIs(t *testing.T) {
	dir, _, _, filed := linkStore(t)

	out, err := runCLI(t, dir, "show", filed)
	if err != nil {
		t.Fatalf("a filed ticket must still be readable: %v", err)
	}
	if !strings.Contains(out, "off the board:") || !strings.Contains(out, "logbook") {
		t.Errorf("show must say where a filed ticket lives\n%s", out)
	}
	if !strings.Contains(out, "The finished blocker") {
		t.Errorf("show must print the filed ticket itself\n%s", out)
	}
}

// The children of a ticket are derived, not stored, so show has to ask for
// them — including the ones that have already shipped.
func TestShowListsWhatATicketContains(t *testing.T) {
	dir, epic, _, _ := linkStore(t)

	out, err := runCLI(t, dir, "show", epic)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "contains") || !strings.Contains(out, "The child") {
		t.Errorf("show must list the children\n%s", out)
	}
}
