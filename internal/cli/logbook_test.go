package cli

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/gitrepo"
	"github.com/BeMuCa/jaira/core/ticket"
)

// syncoutFixture builds a real git repository with a jaira store at its root
// and one ticket sitting in the terminal lane, its file committed alongside a
// commit naming the ticket by id — so 'jaira logbook' has something for git to
// find. Built at test runtime rather than as a committed fixture, per this
// repo's test conventions.
func syncoutFixture(t *testing.T) (dir, id string) {
	t.Helper()
	if !gitrepo.Available() {
		t.Skip("git is not on PATH")
	}
	dir = t.TempDir()
	t.Setenv("JAIRA_USER", "Alexander Sacharov")
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))

	run := func(args ...string) {
		t.Helper()
		if out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
		}
	}
	run("init", "-q")
	run("config", "user.email", "fixture@example.com")
	run("config", "user.name", "Fixture")

	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	// Bypasses the gate deliberately: this test exercises 'jaira logbook', not
	// the pipeline that gets a ticket into the terminal lane.
	id = ticket.NewID(time.Now())
	tk, err := s.Create(map[string]string{
		ticket.FieldID:     id,
		ticket.FieldTitle:  "t",
		ticket.FieldStatus: "done",
	}, nil, "")
	if err != nil {
		t.Fatal(err)
	}

	run("add", ".")
	run("commit", "-q", "-m", "feat: implements "+id)
	return dir, tk.ID
}

func TestLogbookStampsCommitsBeforeMoving(t *testing.T) {
	dir, id := syncoutFixture(t)

	out, err := runCLI(t, dir, "logbook", id, "--json")
	if err != nil {
		t.Fatalf("logbook: %v\n%s", err, out)
	}
	var payload struct {
		Logged  bool     `json:"logged"`
		ID      string   `json:"id"`
		Handle  string   `json:"handle"`
		Path    string   `json:"path"`
		File    string   `json:"file"`
		Commits []string `json:"commits"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("unmarshal %q: %v", out, err)
	}
	if !payload.Logged {
		t.Errorf("logged = false, want true")
	}
	if payload.ID != id {
		t.Errorf("id = %q, want %q", payload.ID, id)
	}
	if payload.Handle != ticket.Handle(id) {
		t.Errorf("handle = %q, want %q", payload.Handle, ticket.Handle(id))
	}
	if len(payload.Commits) == 0 {
		t.Errorf("commits = %v, want at least the commit naming the ticket", payload.Commits)
	}
	if _, err := os.Stat(payload.Path); err != nil {
		t.Errorf("logged file not found at %q: %v", payload.Path, err)
	}
	if !strings.Contains(payload.Path, filepath.Join(".jaira", "logbook")) {
		t.Errorf("path = %q, want it under .jaira/logbook/", payload.Path)
	}

	// The commits are on the file BEFORE it moves: read the moved file's raw
	// frontmatter and confirm it carries them, not just the JSON payload.
	// Load can no longer find it by id — it has left the board — so the file
	// is read directly at the path the command reported.
	raw, err := os.ReadFile(payload.Path)
	if err != nil {
		t.Fatalf("reading the synced file: %v", err)
	}
	if !strings.Contains(string(raw), "commits:") {
		t.Errorf("synced ticket frontmatter = %q, want it to carry the stamped commits field", raw)
	}
}

func TestLogbookFolderIsInitialsAndDate(t *testing.T) {
	dir, id := syncoutFixture(t)

	if out, err := runCLI(t, dir, "logbook", id); err != nil {
		t.Fatalf("logbook: %v\n%s", err, out)
	}
	wantFolder := "as-" + time.Now().Format("20060102")
	if _, err := os.Stat(filepath.Join(dir, ".jaira", "logbook", wantFolder)); err != nil {
		t.Errorf(".jaira/logbook/%s does not exist: %v", wantFolder, err)
	}
}

func TestLogbookRefusesNonTerminalTicket(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	tk, err := s.Create(map[string]string{
		ticket.FieldID:     ticket.NewID(time.Now()),
		ticket.FieldTitle:  "t",
		ticket.FieldStatus: "todo",
	}, nil, "")
	if err != nil {
		t.Fatal(err)
	}

	out, err := runCLI(t, dir, "logbook", tk.ID, "--json")
	if err == nil {
		t.Fatalf("expected logbook of a non-terminal ticket to fail, got %s", out)
	}
	var ce *codedError
	if !errors.As(err, &ce) {
		t.Fatalf("error is not a codedError: %v", err)
	}
	if ce.code != ExitValidation {
		t.Errorf("code = %d, want %d", ce.code, ExitValidation)
	}
	if ce.reason != "not_terminal" {
		t.Errorf("reason = %q, want %q", ce.reason, "not_terminal")
	}
	if !strings.Contains(err.Error(), "todo") || !strings.Contains(err.Error(), "done") || !strings.Contains(err.Error(), "jaira move") {
		t.Errorf("error %q does not name the current lane, the terminal lane, and the move command", err)
	}
}

func TestLogbookWithNoArgumentListsLikeArchive(t *testing.T) {
	dir, id := syncoutFixture(t)

	if out, err := runCLI(t, dir, "logbook"); err != nil {
		t.Fatalf("bare logbook on an empty area: %v\n%s", err, out)
	} else if !strings.Contains(out, "logbook is empty") {
		t.Errorf("bare logbook on an empty area = %q, want it to say the area is empty", out)
	}

	if out, err := runCLI(t, dir, "logbook", id); err != nil {
		t.Fatalf("logbook: %v\n%s", err, out)
	}

	out, err := runCLI(t, dir, "logbook")
	if err != nil {
		t.Fatalf("bare logbook: %v\n%s", err, out)
	}
	if !strings.Contains(out, "jaira restore") {
		t.Errorf("bare logbook listing = %q, want it to end with the restore hint", out)
	}

	jout, err := runCLI(t, dir, "logbook", "--json")
	if err != nil {
		t.Fatalf("bare logbook --json: %v\n%s", err, jout)
	}
	var payload struct {
		Logbook []string `json:"logbook"`
		Count   int      `json:"count"`
	}
	if err := json.Unmarshal([]byte(jout), &payload); err != nil {
		t.Fatalf("unmarshal %q: %v", jout, err)
	}
	if payload.Count != 1 || len(payload.Logbook) != 1 {
		t.Errorf("bare logbook --json = %+v, want exactly one entry", payload)
	}
}

func TestLogbookTwoArgumentsIsUsageError(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}

	_, err = runCLI(t, dir, "logbook", "a", "b")
	if err == nil {
		t.Fatal("expected 'sync a b' to fail")
	}
	var ce *codedError
	if !errors.As(err, &ce) {
		t.Fatalf("error is not a codedError: %v", err)
	}
	if ce.code != ExitUsage {
		t.Errorf("code = %d, want %d", ce.code, ExitUsage)
	}
}

func TestArchiveStampsCommitsBeforeMoving(t *testing.T) {
	dir, id := syncoutFixture(t)

	out, err := runCLI(t, dir, "archive", id)
	if err != nil {
		t.Fatalf("archive: %v\n%s", err, out)
	}
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(s.ArchiveDir())
	if err != nil {
		t.Fatalf("reading archive dir: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("archive dir has %d entries, want 1", len(entries))
	}
	raw, err := os.ReadFile(filepath.Join(s.ArchiveDir(), entries[0].Name()))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "commits:") {
		t.Errorf("archived ticket frontmatter = %q, want it to carry a stamped commits field", raw)
	}
}

func TestLogbookAndSyncTasksBothResolve(t *testing.T) {
	dir := t.TempDir()
	if out, err := runCLI(t, dir, "logbook", "--help"); err != nil {
		t.Fatalf("sync --help: %v\n%s", err, out)
	}
	if out, err := runCLI(t, dir, "sync-tasks", "--help"); err != nil {
		t.Fatalf("sync-tasks --help: %v\n%s", err, out)
	}
}
