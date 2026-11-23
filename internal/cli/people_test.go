package cli

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/people"
)

// The whole feature across two real clones: names written on one arrive on
// the other, two writes from stale clones both survive, the acting name
// follows the list, and a ticket written under an old name is found and shown
// under the new one without its file changing.
func TestPeopleTravelBetweenClonesAndRenameOldTickets(t *testing.T) {
	t.Setenv("JAIRA_USER", "")
	ada, grace := twoBoards(t)

	// A ticket ada files and claims under her git name, before anybody named her.
	if out, err := runAndSend(t, ada, "create", "old work"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	h := handleFromList(t, ada, "old work")
	if out, err := runAndSend(t, ada, "pull", h); err != nil {
		t.Fatalf("pull: %v\n%s", err, out)
	}
	if out, err := runAndSend(t, ada, "set", h, "assignee=ada"); err != nil {
		t.Fatalf("set: %v\n%s", err, out)
	}

	// Two clones that have not fetched each other's names both write.
	if out, err := runCLI(t, ada, "people", "alias", "ada", "ada@example.test"); err != nil {
		t.Fatalf("alias on ada: %v\n%s", err, out)
	}
	if out, err := runCLI(t, grace, "people", "alias", "grace", "gh"); err != nil {
		t.Fatalf("alias on grace: %v\n%s", err, out)
	}
	if out, err := runCLI(t, ada, "fetch"); err != nil {
		t.Fatalf("fetch: %v\n%s", err, out)
	}
	out, err := runCLI(t, ada, "people", "--json")
	if err != nil {
		t.Fatalf("people: %v\n%s", err, out)
	}
	var listed struct{ People []people.Person }
	if err := json.Unmarshal([]byte(out), &listed); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if len(listed.People) != 2 || listed.People[0].Name != "ada" || listed.People[1].Name != "grace" {
		t.Fatalf("one write overwrote the other: %+v", listed.People)
	}

	// Renamed on grace's clone; ada acts under the new name once she fetches.
	if out, err := runCLI(t, grace, "people", "rename", "ada", "Ada Lovelace"); err != nil {
		t.Fatalf("rename: %v\n%s", err, out)
	}
	if out, err := runCLI(t, ada, "fetch"); err != nil {
		t.Fatalf("fetch: %v\n%s", err, out)
	}
	if got := strings.TrimSpace(gitOut(t, ada, "show", gitref.PeopleRef+":"+gitref.PeopleFile)); !strings.Contains(got, "Ada Lovelace: ada, ada@example.test") {
		t.Fatalf("the rename did not arrive:\n%s", got)
	}
	out, err = runCLI(t, ada, "whoami", "--json")
	if err != nil {
		t.Fatalf("whoami: %v\n%s", err, out)
	}
	var who struct {
		Identity string
		Aliases  []string
	}
	if err := json.Unmarshal([]byte(out), &who); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if who.Identity != "Ada Lovelace" || !containsName(who.Aliases, "ada") {
		t.Fatalf("whoami = %+v", who)
	}

	// The old ticket: found by the new name, shown under it, its file untouched.
	for _, filter := range []string{"Ada Lovelace", "ada", "ADA@example.test"} {
		out, err = runCLI(t, ada, "list", "--assignee", filter, "--json")
		if err != nil {
			t.Fatalf("list --assignee %s: %v\n%s", filter, err, out)
		}
		if !strings.Contains(out, `"assignee": "Ada Lovelace"`) {
			t.Errorf("list --assignee %q does not show the old ticket under the new name:\n%s", filter, out)
		}
	}
	out, err = runCLI(t, ada, "show", h, "--json")
	if err != nil {
		t.Fatalf("show: %v\n%s", err, out)
	}
	if !strings.Contains(out, `"assignee": "Ada Lovelace"`) {
		t.Errorf("show does not use the new name:\n%s", out)
	}
	var shown struct{ ID string }
	if err := json.Unmarshal([]byte(out), &shown); err != nil {
		t.Fatal(err)
	}
	if raw := gitOut(t, ada, "show", gitref.RefName(shown.ID)+":"+shown.ID+".md"); !strings.Contains(raw, "assignee: ada\n") {
		t.Errorf("the ticket file was rewritten:\n%s", raw)
	}

	// Somebody else's pull is refused under the name she goes by now.
	out, err = runCLI(t, grace, "pull", h)
	if err == nil || !strings.Contains(err.Error()+out, "Ada Lovelace") {
		t.Errorf("the refusal does not name the jaira name: %v\n%s", err, out)
	}

	// A new write by ada carries her jaira name.
	if out, err := runAndSend(t, ada, "create", "new work"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	out, err = runCLI(t, ada, "show", handleFromList(t, ada, "new work"), "--json")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"creator": "Ada Lovelace"`) {
		t.Errorf("a new ticket is not written under the jaira name:\n%s", out)
	}
}

// Refusals carry the exit code a script branches on, and write nothing.
func TestPeopleRefusals(t *testing.T) {
	t.Setenv("JAIRA_USER", "")
	ada, _ := twoBoards(t)
	if _, err := runCLI(t, ada, "people", "alias", "ada", "AL"); err != nil {
		t.Fatal(err)
	}
	if _, err := runCLI(t, ada, "people", "alias", "grace", "gh"); err != nil {
		t.Fatal(err)
	}
	before := gitOut(t, ada, "rev-parse", gitref.PeopleRef)
	cases := []struct {
		args []string
		code int
	}{
		{[]string{"people", "alias", "grace", "al"}, ExitValidation},
		{[]string{"people", "rename", "grace", "ada"}, ExitValidation},
		{[]string{"people", "rename", "nobody", "x"}, ExitNotFound},
		{[]string{"people", "alias", "ada", "a,b"}, ExitUsage},
	}
	for _, c := range cases {
		_, err := runCLI(t, ada, c.args...)
		if code := exitCode(err); code != c.code {
			t.Errorf("%v: exit %d, want %d (%v)", c.args, code, c.code, err)
		}
	}
	// Naming nobody new is not an error and writes nothing either.
	for _, args := range [][]string{{"people", "alias", "zed", "ZED"}, {"people", "alias", "ada", "al"}} {
		if out, err := runCLI(t, ada, args...); err != nil || !strings.Contains(out, "Nothing to change") {
			t.Errorf("%v: %v\n%s", args, err, out)
		}
	}
	if out, err := runCLI(t, ada, "people", "--json", "alias", "zed", "zed"); err != nil || !strings.Contains(out, `"changed": false`) {
		t.Errorf("json: %v\n%s", err, out)
	}
	if gitOut(t, ada, "rev-parse", gitref.PeopleRef) != before {
		t.Error("a refused write changed the list")
	}
}

func containsName(list []string, v string) bool {
	for _, x := range list {
		if strings.EqualFold(x, v) {
			return true
		}
	}
	return false
}
