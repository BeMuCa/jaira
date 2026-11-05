package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/outbox"
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

// An epic is addressed by its name in every message, never by the last six
// characters a ticket handle would cut it to.
func TestAnEpicIsReportedByItsName(t *testing.T) {
	if key, name := refSubject(outbox.KindEpic, "feature", ""); key != "epic" || name != "feature" {
		t.Errorf("refSubject = %q, %q; want epic, feature", key, name)
	}
	if a := rejectedAdvice(outbox.KindEpic); strings.Contains(a, "file") {
		t.Errorf("epic advice talks about a file it does not have: %q", a)
	}
}

// 'epic set' rewrites what an epic says and nothing else: the members stay,
// the text comes back out of 'epic show', and a name nobody created or a call
// that names nothing to set is refused with the code a script branches on.
func TestEpicSetRewritesTheTextNotTheMembers(t *testing.T) {
	ada, _ := twoBoards(t)
	a := createTicket(t, ada, "One")
	jsonOf(t, ada, "epic", "create", "feature", a, "--title", "Old")

	if out, err := runCLIStdin(t, ada, "## Goal\n\nlog in by link\n", "epic", "set", "feature", "--title", "New", "--description", "one line", "--body-file", "-"); err != nil {
		t.Fatalf("set: %v\n%s", err, out)
	}
	flushRefs()
	got := jsonOf(t, ada, "epic", "show", "feature")
	if got["title"] != "New" || got["description"] != "one line" || got["body"] != "## Goal\n\nlog in by link" {
		t.Errorf("after set: %v", got)
	}
	if m := got["members"].([]any); len(m) != 1 {
		t.Errorf("members = %v; want %s alone", m, ticket.Handle(a))
	}
	if out, _ := runCLI(t, ada, "epic", "show", "feature"); !strings.Contains(out, "log in by link") {
		t.Errorf("show does not print the text:\n%s", out)
	}
	if _, err := runCLI(t, ada, "epic", "set", "nothing-here", "--title", "x"); exitCode(err) != ExitNotFound {
		t.Errorf("unknown epic: exit %d, want %d", exitCode(err), ExitNotFound)
	}
	if _, err := runCLI(t, ada, "epic", "set", "feature", "--title", ""); exitCode(err) != ExitUsage {
		t.Errorf("empty title: exit %d, want %d", exitCode(err), ExitUsage)
	}
	if _, err := runCLI(t, ada, "epic", "set", "feature"); exitCode(err) != ExitUsage {
		t.Errorf("no flag: exit %d, want %d", exitCode(err), ExitUsage)
	}
}

// 'epic close' marks the epic closed on its ref and deletes nothing; a closed
// epic refuses 'epic add' with exit 3 naming --force, takes the ticket with
// --force and stays closed, is missing from 'epic ls' and listed marked by
// 'epic ls --all'. 'epic reopen' takes it all back. Closing with an
// unfinished ticket is refused without --force.
func TestEpicCloseAndReopen(t *testing.T) {
	ada, _ := twoBoards(t)
	a := createTicket(t, ada, "One")
	b := createTicket(t, ada, "Two")
	jsonOf(t, ada, "epic", "create", "feature", a)

	out, err := runCLI(t, ada, "epic", "close", "feature")
	if exitCode(err) != ExitValidation || !strings.Contains(out+err.Error(), "--force") {
		t.Fatalf("close with an open ticket: exit %d, %v\n%s; want 3 naming --force", exitCode(err), err, out)
	}
	if got := jsonOf(t, ada, "epic", "close", "feature", "--force"); got["closed"] != true || got["changed"] != true {
		t.Fatalf("forced close: %v", got)
	}
	flushRefs()
	if file := gitOut(t, ada, "show", gitref.EpicRefName("feature")+":feature.md"); !strings.Contains(file, "status: closed") || !strings.Contains(file, "- "+a) {
		t.Fatalf("the ref after close:\n%s", file)
	}

	out, err = runCLI(t, ada, "epic", "add", "feature", b)
	if exitCode(err) != ExitValidation || !strings.Contains(out+err.Error(), "--force") {
		t.Fatalf("add to closed: exit %d, %v\n%s; want 3 naming --force", exitCode(err), err, out)
	}
	if got := jsonOf(t, ada, "epic", "add", "feature", b, "--force"); got["closed"] != true || len(got["changed"].([]any)) != 1 {
		t.Fatalf("forced add: %v; want %s in and the epic still closed", got, ticket.Handle(b))
	}

	if ls := jsonOf(t, ada, "epic", "ls"); ls["count"].(float64) != 0 {
		t.Errorf("ls lists a closed epic: %v", ls)
	}
	if ls := jsonOf(t, ada, "epic", "ls", "--all"); ls["count"].(float64) != 1 || ls["epics"].([]any)[0].(map[string]any)["closed"] != true {
		t.Errorf("ls --all = %v; want the closed epic, marked", ls)
	}
	if out, _ := runCLI(t, ada, "epic", "ls", "--all"); !strings.Contains(out, "[closed]") {
		t.Errorf("ls --all does not mark it:\n%s", out)
	}

	if got := jsonOf(t, ada, "epic", "reopen", "feature"); got["closed"] != false || got["changed"] != true {
		t.Fatalf("reopen: %v", got)
	}
	if ls := jsonOf(t, ada, "epic", "ls"); ls["count"].(float64) != 1 {
		t.Errorf("ls after reopen = %v; want the epic back", ls)
	}
	if _, err := runCLI(t, ada, "epic", "close", "nothing-here"); exitCode(err) != ExitNotFound {
		t.Errorf("unknown epic: exit %d, want %d", exitCode(err), ExitNotFound)
	}
}

// An epic carries the keys of the tasks it stands for in another system, and
// they find the work back: 'epic ls --external' the epics, 'list --external'
// their tickets, follow-ups included. A key nobody typed right is refused.
func TestEpicExternalKeysFindTheWork(t *testing.T) {
	ada, _ := twoBoards(t)
	a := createTicket(t, ada, "Login form")
	b := createTicket(t, ada, "Fix the form", "--follows", a)
	c := createTicket(t, ada, "Unrelated")
	jsonOf(t, ada, "epic", "create", "login", a, "--external", "JIRA-1 https://jira.example/browse/JIRA-1")
	jsonOf(t, ada, "epic", "create", "other", c)
	got := jsonOf(t, ada, "epic", "set", "other", "--external", "YT-7", "--external", "jira-1 https://jira.example/JIRA-1")
	if x := got["external"].([]any); len(x) != 2 {
		t.Fatalf("set --external: %v", got["external"])
	}
	jsonOf(t, ada, "epic", "set", "other", "--external-rm", "JIRA-1")

	show := jsonOf(t, ada, "epic", "show", "login")
	if x := show["external"].([]any); len(x) != 1 || x[0].(map[string]any)["url"] != "https://jira.example/browse/JIRA-1" {
		t.Errorf("show --json external = %v", show["external"])
	}
	if out, _ := runCLI(t, ada, "epic", "show", "login"); !strings.Contains(out, "external: JIRA-1 https://jira.example/browse/JIRA-1") {
		t.Errorf("show does not print the key:\n%s", out)
	}
	if out, _ := runCLI(t, ada, "epic", "ls"); !strings.Contains(out, "[JIRA-1]") || !strings.Contains(out, "[YT-7]") {
		t.Errorf("ls does not print the keys:\n%s", out)
	}
	ls := jsonOf(t, ada, "epic", "ls", "--external", "jira-1")
	if e := ls["epics"].([]any); len(e) != 1 || e[0].(map[string]any)["name"] != "login" {
		t.Errorf("ls --external jira-1 = %v; want login alone", ls["epics"])
	}
	if e := jsonOf(t, ada, "epic", "ls")["epics"].([]any); len(e) != 2 {
		t.Errorf("ls without a key lists %d epics, want 2", len(e))
	}

	list := jsonOf(t, ada, "list", "--external", "JIRA-1")
	var ids []string
	for _, row := range list["tickets"].([]any) {
		ids = append(ids, row.(map[string]any)["id"].(string))
	}
	if len(ids) != 2 || !slices.Contains(ids, a) || !slices.Contains(ids, b) {
		t.Errorf("list --external JIRA-1 = %v; want the member and its follow-up", ids)
	}
	if n := jsonOf(t, ada, "list", "--external", "NOPE-1")["count"].(float64); n != 0 {
		t.Errorf("list for a key nobody carries = %v tickets", n)
	}

	for _, bad := range [][]string{
		{"epic", "set", "login", "--external", "JIRA-2 javascript:alert(1)"},
		{"epic", "set", "login", "--external", "JIRA 2 https://x.example"},
		{"epic", "create", "third", "--external", "A ftp://x.example"},
	} {
		if _, err := runCLI(t, ada, bad...); exitCode(err) != ExitUsage {
			t.Errorf("%v: exit %d, want %d", bad, exitCode(err), ExitUsage)
		}
	}
}

// The key is on the epic's ref and nowhere else: setting it writes no
// ticket's ref, no file in the working tree, and no commit outside
// refs/jaira/epics carries it — so a code commit never takes it along.
func TestAnExternalKeyStaysOnTheEpicRef(t *testing.T) {
	ada, _ := twoBoards(t)
	a := createTicket(t, ada, "Login form")
	if out, err := runCLI(t, ada, "pull", ticket.Handle(a)); err != nil {
		t.Fatalf("pull: %v\n%s", err, out)
	}
	jsonOf(t, ada, "epic", "create", "login", a)
	repo := &gitref.Repo{Dir: ada}
	before, err := repo.SHA(a)
	if err != nil {
		t.Fatal(err)
	}
	const key = "BILL-4711"
	jsonOf(t, ada, "epic", "set", "login", "--external", key+" https://jira.example/browse/"+key)
	flushRefs()

	if after, _ := repo.SHA(a); after != before {
		t.Errorf("setting a key wrote the ticket's ref: %s -> %s", before, after)
	}
	if file := gitOut(t, ada, "show", gitref.EpicRefName("login")+":login.md"); !strings.Contains(file, key) {
		t.Fatalf("the epic's ref lacks the key:\n%s", file)
	}
	err = filepath.WalkDir(ada, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			if d != nil && d.Name() == ".git" {
				return filepath.SkipDir
			}
			return err
		}
		raw, err := os.ReadFile(path)
		if err == nil && strings.Contains(string(raw), key) {
			t.Errorf("%s carries the key", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if log := gitOut(t, ada, "log", "-p", "--exclude=refs/jaira/epics/*", "--all"); strings.Contains(log, key) {
		t.Errorf("a commit outside the epic refs carries the key")
	}
}

// Wherever the CLI names an epic, it names it by its title, the name to type
// beside it; an epic titled by its name is named once.
func TestTheCLINamesAnEpicByItsTitle(t *testing.T) {
	ada, _ := twoBoards(t)
	a := createTicket(t, ada, "Login form")
	b := createTicket(t, ada, "Second")
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"epic", "create", "web-0-2", "--title", "Web 0.2"}, `Epic "Web 0.2" (web-0-2) created`},
		{[]string{"epic", "create", "plain"}, `Epic "plain" created`},
		{[]string{"epic", "add", "web-0-2", ticket.Handle(a)}, `added to "Web 0.2" (web-0-2)`},
		{[]string{"epic", "add", "web-0-2", ticket.Handle(b)}, `"Web 0.2" (web-0-2) now lists 2 ticket(s).`},
		{[]string{"epic", "rm", "web-0-2", ticket.Handle(b)}, `"Web 0.2" (web-0-2) now lists 1 ticket(s).`},
		{[]string{"epic", "ls"}, "  Web 0.2  (web-0-2)    0 of 1   done"},
		{[]string{"epic", "show", "web-0-2"}, "Web 0.2  (web-0-2)\n"},
		{[]string{"epic", "set", "web-0-2", "--description", "d"}, `Epic "Web 0.2" (web-0-2) rewritten`},
		{[]string{"epic", "close", "web-0-2", "--force"}, `Epic "Web 0.2" (web-0-2) closed`},
	} {
		out, err := runAndSend(t, ada, c.args...)
		if err != nil {
			t.Fatalf("%v: %v\n%s", c.args, err, out)
		}
		if !strings.Contains(out, c.want) {
			t.Errorf("%v printed:\n%s\nwant it to contain %q", c.args, out, c.want)
		}
	}
	if out, err := runCLI(t, ada, "epic", "add", "web-0-2", ticket.Handle(b)); err == nil || !strings.Contains(out+err.Error(), `epic "Web 0.2" (web-0-2) is closed`) {
		t.Errorf("closed epic refusal: %v\n%s", err, out)
	}
}
