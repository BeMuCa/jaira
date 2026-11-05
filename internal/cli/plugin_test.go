//go:build unix

package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/release"
	"github.com/BeMuCa/jaira/core/ticket"
)

// installPlugin writes a shell script where 'jaira plugin install' will leave
// jaira-<id>, with a manifest subscribing it to events, under a JAIRA_HOME of
// the test's own.
func installPlugin(t *testing.T, id, body string, events ...string) {
	t.Helper()
	home := os.Getenv("JAIRA_HOME")
	dir := filepath.Join(home, "plugins", id)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "jaira-"+id), []byte("#!/bin/sh\n"+body+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(map[string]any{"id": id, "events": events})
	if err := os.WriteFile(filepath.Join(dir, "jaira-plugin.json"), raw, 0o644); err != nil {
		t.Fatal(err)
	}
}

func pluginBoard(t *testing.T) string {
	t.Helper()
	t.Setenv("JAIRA_HOME", t.TempDir())
	t.Setenv("JAIRA_IN_EVENT", "")
	t.Setenv("JAIRA_USER", "berk")
	dir := t.TempDir()
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	return dir
}

// 'jaira <id>' starts the installed plugin with every argument as typed, the
// three variables, and returns its exit code as jaira's.
func TestUnknownCommandStartsTheInstalledPlugin(t *testing.T) {
	dir := pluginBoard(t)
	t.Chdir(dir)
	seen := filepath.Join(t.TempDir(), "seen")
	installPlugin(t, "jira", `printf '%s|' "$@" > `+seen+`; echo "$JAIRA_ROOT $JAIRA_BOARD $JAIRA_VERSION" >> `+seen+`; exit 4`)
	release.Current = "9.9.9"

	code, ok := runPlugin(newRoot("test"), []string{"jira", "pull", "--board", "ops", "--json"})
	if !ok || code != 4 {
		t.Fatalf("runPlugin = %d, %v; want 4, true", code, ok)
	}
	raw, err := os.ReadFile(seen)
	if err != nil {
		t.Fatal(err)
	}
	root, _ := filepath.EvalSymlinks(dir)
	if got, want := string(raw), "pull|--board|ops|--json|"+root+" ops 9.9.9\n"; got != want {
		t.Errorf("plugin saw %q, want %q", got, want)
	}

	// Without --board the board is today's.
	if _, ok := runPlugin(newRoot("test"), []string{"jira"}); !ok {
		t.Fatal("the plugin did not start without arguments")
	}
	raw, _ = os.ReadFile(seen)
	if got, want := string(raw), "|"+root+" default 9.9.9\n"; got != want {
		t.Errorf("plugin saw %q, want %q", got, want)
	}
}

// jaira's own commands win: a plugin named like one — or like cobra's help —
// is never started, and a name nobody installed stays an unknown command.
func TestBuiltinCommandWinsOverAPluginOfTheSameName(t *testing.T) {
	pluginBoard(t)
	ran := filepath.Join(t.TempDir(), "ran")
	for _, id := range []string{"list", "help", "completion"} {
		installPlugin(t, id, "touch "+ran)
		if _, ok := runPlugin(newRoot("test"), []string{id}); ok {
			t.Errorf("a plugin called %q was started instead of jaira's command", id)
		}
	}
	if _, err := os.Stat(ran); err == nil {
		t.Fatal("a plugin ran in place of a built-in command")
	}
	if _, ok := runPlugin(newRoot("test"), []string{"youtrack"}); ok {
		t.Error("a plugin that is not installed was started")
	}
	// What it said before plugins existed.
	if _, err := runCLI(t, t.TempDir(), "youtrack"); err == nil || !strings.Contains(err.Error(), "takes no arguments") {
		t.Errorf("an unknown command without a plugin: %v, want cobra's refusal", err)
	}
}

// create, set, claim and move each tell a subscribed plugin, with what show
// --json prints plus actor and from_lane on stdin.
func TestWritesFireThePluginsEvents(t *testing.T) {
	dir := pluginBoard(t)
	log := filepath.Join(t.TempDir(), "events")
	installPlugin(t, "rec", `{ echo "== $1 $2"; cat; echo; } >> `+log, "create", "move", "claim", "set")

	out, err := runCLI(t, dir, "create", "plugged", "--goal", "g", "--context", "c", "--dod", "d", "--json")
	if err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	var created map[string]any
	if err := json.Unmarshal([]byte(out), &created); err != nil {
		t.Fatalf("create --json: %v\n%s", err, out)
	}
	h := created["handle"].(string)
	for _, args := range [][]string{
		{"set", h, "external=jira:PROJ-1"},
		{"claim", h, "--session", "me"},
		{"move", h, "--to", "todo", "--force"},
	} {
		if out, err := runCLI(t, dir, args...); err != nil {
			t.Fatalf("%v: %v\n%s", args, err, out)
		}
	}

	raw, err := os.ReadFile(log)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]map[string]any{}
	var order []string
	for _, part := range strings.Split(string(raw), "== ")[1:] {
		head, body, _ := strings.Cut(part, "\n")
		var p map[string]any
		if err := json.Unmarshal([]byte(body), &p); err != nil {
			t.Fatalf("%s: stdin is not JSON: %v\n%s", head, err, body)
		}
		order = append(order, head)
		got[head] = p
	}
	if want := "event create,event set,event claim,event move"; strings.Join(order, ",") != want {
		t.Fatalf("events %v, want %s", order, want)
	}
	for head, p := range got {
		if p["handle"] != h || p["actor"] != "berk" {
			t.Errorf("%s: handle %v actor %v", head, p["handle"], p["actor"])
		}
	}
	if got["event set"]["external"] != "jira:PROJ-1" {
		t.Errorf("set payload external = %v", got["event set"]["external"])
	}
	if got["event move"]["from_lane"] != "backlog" || got["event move"]["status"] != "todo" {
		t.Errorf("move payload from_lane %v status %v", got["event move"]["from_lane"], got["event move"]["status"])
	}
	if got["event create"]["from_lane"] != "" {
		t.Errorf("create payload from_lane = %v, want empty", got["event create"]["from_lane"])
	}
}

// A plugin that fails changes nothing about the command: the write stands and
// the exit code is the command's.
func TestAFailingPluginDoesNotFailTheCommand(t *testing.T) {
	dir := pluginBoard(t)
	installPlugin(t, "broken", "echo noise; echo more >&2; exit 1", "create")
	out, err := runCLI(t, dir, "create", "still filed", "--goal", "g", "--context", "c", "--dod", "d", "--json")
	if err != nil {
		t.Fatalf("create failed because a plugin did: %v\n%s", err, out)
	}
	if strings.Contains(out, "noise") || strings.Contains(out, "more") {
		t.Errorf("the plugin's output reached the command's:\n%s", out)
	}
}

// The external field is in show --json and in every list --json row, so a
// plugin reads its key without parsing the ticket file.
func TestExternalIsInShowAndListJSON(t *testing.T) {
	dir := pluginBoard(t)
	out, err := runCLI(t, dir, "create", "mirrored", "--goal", "g", "--context", "c", "--dod", "d", "--json")
	if err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	var created map[string]any
	_ = json.Unmarshal([]byte(out), &created)
	h := created["handle"].(string)
	if created["external"] != "" {
		t.Errorf("a new ticket's external = %v, want empty", created["external"])
	}
	if out, err := runCLI(t, dir, "set", h, "external=jira:PROJ-7"); err != nil {
		t.Fatalf("set: %v\n%s", err, out)
	}

	out, err = runCLI(t, dir, "show", h, "--json")
	if err != nil {
		t.Fatalf("show: %v\n%s", err, out)
	}
	var shown map[string]any
	_ = json.Unmarshal([]byte(out), &shown)
	if shown["external"] != "jira:PROJ-7" {
		t.Errorf("show --json external = %v", shown["external"])
	}

	out, err = runCLI(t, dir, "list", "--json")
	if err != nil {
		t.Fatalf("list: %v\n%s", err, out)
	}
	var listed struct {
		Tickets []map[string]any `json:"tickets"`
	}
	if err := json.Unmarshal([]byte(out), &listed); err != nil {
		t.Fatalf("list --json: %v\n%s", err, out)
	}
	if len(listed.Tickets) != 1 || listed.Tickets[0]["external"] != "jira:PROJ-7" {
		t.Errorf("list --json rows %v", listed.Tickets)
	}
}
