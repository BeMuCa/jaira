package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/milestone"
)

// milestoneID is the ID of the one milestone a board holds under name.
func milestoneID(t *testing.T, dir, name string) string {
	t.Helper()
	all, err := milestone.LoadAll(dir)
	if err != nil {
		t.Fatalf("load milestones: %v", err)
	}
	for _, m := range all {
		if m.Name == name {
			return m.ID
		}
	}
	t.Fatalf("no milestone %q in %s", name, dir)
	return ""
}

// A road is an object on its own ref: it reaches another clone on a fetch with
// no branch merged, lists its milestones by name there, and leaves no file in
// the working tree.
func TestARoadTravelsOnItsRefAndListsItsMilestones(t *testing.T) {
	ada, grace := twoBoards(t)
	for _, name := range []string{"round-one", "round-two"} {
		if out, err := runAndSend(t, ada, "milestone", "create", name); err != nil {
			t.Fatalf("milestone create: %v\n%s", err, out)
		}
	}
	first := milestoneID(t, ada, "round-one")

	if out, err := runAndSend(t, ada, "road", "create", "Auth Rollout", "round-one", "--title", "Auth rollout", "--description", "login first"); err != nil {
		t.Fatalf("road create: %v\n%s", err, out)
	}
	if _, err := os.Stat(filepath.Join(ada, ".jaira", "roads")); !os.IsNotExist(err) {
		t.Errorf("a road left a working file behind; it lives on its ref only (stat: %v)", err)
	}
	file := gitOut(t, ada, "show", "refs/jaira/roads/auth-rollout:auth-rollout.md")
	for _, want := range []string{"title: Auth rollout", "description: login first", "- " + first} {
		if !strings.Contains(file, want) {
			t.Errorf("the ref's file lacks %q:\n%s", want, file)
		}
	}

	if out, err := runCLI(t, grace, "fetch"); err != nil {
		t.Fatalf("fetch: %v\n%s", err, out)
	}
	out, err := runCLI(t, grace, "road", "ls")
	if err != nil {
		t.Fatalf("road ls: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Auth rollout") || !strings.Contains(out, "round-one-"+first) {
		t.Errorf("grace's 'road ls' does not show the road and its milestone:\n%s", out)
	}
	got := jsonOf(t, grace, "road", "ls")
	roads := got["roads"].([]any)
	if len(roads) != 1 {
		t.Fatalf("road ls --json = %v; want one road", got)
	}
	ms := roads[0].(map[string]any)["milestones"].([]any)
	if len(ms) != 1 || ms[0].(map[string]any)["name"] != "round-one" {
		t.Errorf("milestones = %v; want round-one", ms)
	}
}

// A milestone is in one road at a time. A second road refuses it with exit 3,
// naming the road it is in, and writes nothing.
func TestAMilestoneInTwoRoadsIsRefused(t *testing.T) {
	ada, _ := twoBoards(t)
	if out, err := runAndSend(t, ada, "milestone", "create", "round-one"); err != nil {
		t.Fatalf("milestone create: %v\n%s", err, out)
	}
	if out, err := runAndSend(t, ada, "road", "create", "auth", "round-one"); err != nil {
		t.Fatalf("road create: %v\n%s", err, out)
	}
	out, err := runAndSend(t, ada, "road", "create", "billing", "round-one")
	if err == nil || exitCode(err) != ExitValidation || !strings.Contains(out+err.Error(), `"auth"`) {
		t.Fatalf("second road = %v\n%s; want exit 3 naming road auth", err, out)
	}
	if got := jsonOf(t, ada, "road", "ls"); got["count"].(float64) != 1 {
		t.Errorf("a refused road was written anyway: %v", got)
	}
}

// Renaming a road changes its title and keeps the milestones it lists. An
// unknown road is not found, and a milestone that does not exist is not found.
func TestRoadSetKeepsItsMembersAndUnknownsAreNotFound(t *testing.T) {
	ada, _ := twoBoards(t)
	if out, err := runAndSend(t, ada, "milestone", "create", "round-one"); err != nil {
		t.Fatalf("milestone create: %v\n%s", err, out)
	}
	if out, err := runAndSend(t, ada, "road", "create", "auth", "round-one"); err != nil {
		t.Fatalf("road create: %v\n%s", err, out)
	}
	if out, err := runAndSend(t, ada, "road", "set", "auth", "--title", "Login"); err != nil {
		t.Fatalf("road set: %v\n%s", err, out)
	}
	got := jsonOf(t, ada, "road", "ls")
	r := got["roads"].([]any)[0].(map[string]any)
	if r["title"] != "Login" || len(r["milestones"].([]any)) != 1 {
		t.Errorf("after set: %v; want title Login and its one milestone", r)
	}

	out, err := runAndSend(t, ada, "road", "add", "nope", "round-one")
	if err == nil || exitCode(err) != ExitNotFound {
		t.Errorf("add to an unknown road = %v\n%s; want exit 5", err, out)
	}
	out, err = runAndSend(t, ada, "road", "add", "auth", "no-such-milestone")
	if err == nil || exitCode(err) != ExitNotFound {
		t.Errorf("add an unknown milestone = %v\n%s; want exit 5", err, out)
	}
}

// Roads live on refs only, so a board without a remote has none and says so.
func TestARoadNeedsARemote(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))
	gitRun(t, dir, "init", "--quiet")
	if out, err := runCLI(t, dir, "init"); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	out, err := runCLI(t, dir, "road", "create", "auth")
	if err == nil || !strings.Contains(out+err.Error(), "remote") {
		t.Fatalf("create without a remote = %v\n%s; want a refusal naming the remote", err, out)
	}
}
