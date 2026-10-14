package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/ticket"
)

// jsonOf runs a command with --json and decodes its result. Warnings go to the
// same buffer in tests, so the object is read from its first brace.
func jsonOf(t *testing.T, dir string, args ...string) map[string]any {
	t.Helper()
	out, err := runAndSend(t, dir, append(args, "--json")...)
	if err != nil {
		t.Fatalf("%v: %v\n%s", args, err, out)
	}
	var got map[string]any
	start := strings.Index(out, "{")
	if start < 0 {
		t.Fatalf("%v printed no JSON:\n%s", args, out)
	}
	if err := json.NewDecoder(strings.NewReader(out[start:])).Decode(&got); err != nil {
		t.Fatalf("%v: %v\n%s", args, err, out)
	}
	return got
}

func createTicket(t *testing.T, dir, title string, extra ...string) string {
	t.Helper()
	args := append([]string{"create", title, "--goal", "g", "--context", "c", "--dod", "d"}, extra...)
	return jsonOf(t, dir, args...)["id"].(string)
}

// An epic is the description plus the list, on its own ref and nowhere else.
// It reaches the other clone on a fetch, counts the follow-up of a member
// without a line for it, and sorting a ticket in never writes that ticket's
// ref — the member list is in the epic, so a teamlead sorting work cannot
// collide with an agent working it.
func TestAnEpicTravelsOnItsRefAndCountsFollowUps(t *testing.T) {
	ada, grace := twoBoards(t)
	a := createTicket(t, ada, "Login form")
	b := createTicket(t, ada, "Fix the form", "--follows", a)
	c := createTicket(t, ada, "Unrelated")

	repo := &gitref.Repo{Dir: ada}
	before, err := repo.SHA(a)
	if err != nil {
		t.Fatal(err)
	}
	jsonOf(t, ada, "epic", "create", "magic-login", "--title", "Anmeldung per Link", "--description", "ohne Passwort")
	if got := jsonOf(t, ada, "epic", "add", "magic-login", ticket.Handle(a)); len(got["changed"].([]any)) != 1 {
		t.Fatalf("add changed nothing: %v", got)
	}
	if after, _ := repo.SHA(a); after != before {
		t.Errorf("sorting %s into an epic wrote its ticket ref: %s -> %s", ticket.Handle(a), before, after)
	}
	if _, err := os.Stat(filepath.Join(ada, ".jaira", "epics")); !os.IsNotExist(err) {
		t.Errorf("an epic left a working file behind; it lives on its ref only (stat: %v)", err)
	}
	file := gitOut(t, ada, "show", gitref.EpicRefName("magic-login")+":magic-login.md")
	for _, want := range []string{"title: Anmeldung per Link", "description: ohne Passwort", "- " + a} {
		if !strings.Contains(file, want) {
			t.Errorf("the ref's file lacks %q:\n%s", want, file)
		}
	}

	if out, err := runCLI(t, grace, "fetch"); err != nil {
		t.Fatalf("fetch: %v\n%s", err, out)
	}
	got := jsonOf(t, grace, "epic", "show", "magic-login")
	if got["total"].(float64) != 2 || got["done"].(float64) != 0 {
		t.Fatalf("grace counts %v of %v; want 0 of 2 (member and its follow-up)", got["done"], got["total"])
	}
	var direct, follow int
	for _, raw := range got["tickets"].([]any) {
		tk := raw.(map[string]any)
		switch {
		case tk["id"] == a && tk["direct"] == true:
			direct++
		case tk["id"] == b && tk["direct"] == false:
			follow++
		}
	}
	if direct != 1 || follow != 1 {
		t.Errorf("tickets = %v; want %s direct and %s as follow-up", got["tickets"], ticket.Handle(a), ticket.Handle(b))
	}

	un := jsonOf(t, grace, "epic", "ls", "--unassigned")
	if un["count"].(float64) != 1 || un["unassigned"].([]any)[0].(map[string]any)["id"] != c {
		t.Errorf("unassigned = %v; want only %s", un, ticket.Handle(c))
	}
	ls := jsonOf(t, grace, "epic", "ls")
	if ls["count"].(float64) != 1 {
		t.Errorf("ls = %v; want one epic", ls)
	}
}

// Two people sorting tickets into one epic at once must both get theirs in.
// The epic has no file on disk, so a lost race that simply dropped the queued
// write would lose a ticket somebody sorted in, with nothing left to merge
// from.
func TestTwoClonesSortingIntoOneEpicBothLand(t *testing.T) {
	ada, grace := twoBoards(t)
	a := createTicket(t, ada, "One")
	c := createTicket(t, ada, "Two")
	jsonOf(t, ada, "epic", "create", "feature")
	if out, err := runCLI(t, grace, "fetch"); err != nil {
		t.Fatalf("fetch: %v\n%s", err, out)
	}

	jsonOf(t, ada, "epic", "add", "feature", a)
	// grace has not fetched ada's add: her write is leased against the
	// state before it and loses the race.
	jsonOf(t, grace, "epic", "add", "feature", c)

	if out, err := runCLI(t, ada, "fetch"); err != nil {
		t.Fatalf("fetch: %v\n%s", err, out)
	}
	file := gitOut(t, ada, "show", gitref.EpicRefName("feature")+":feature.md")
	if !strings.Contains(file, "- "+a) || !strings.Contains(file, "- "+c) {
		t.Fatalf("one of the two additions was lost:\n%s", file)
	}
}

// Epics live on refs only, so a board without a remote has none, and says so
// rather than pretending to keep one.
func TestAnEpicNeedsARemote(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))
	gitRun(t, dir, "init", "--quiet")
	if out, err := runCLI(t, dir, "init"); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	out, err := runCLI(t, dir, "epic", "create", "feature")
	if err == nil || !strings.Contains(out+err.Error(), "remote") {
		t.Fatalf("create without a remote = %v\n%s; want a refusal naming the remote", err, out)
	}
}
