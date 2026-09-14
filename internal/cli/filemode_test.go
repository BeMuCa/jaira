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

// The silence this ticket exists for: a board whose remote is not there writes
// the ticket to disk and says nothing, so seventeen of them can be created
// before anybody notices the board stopped carrying tickets on refs.
func TestCreateSaysWhenTheTicketStaysAFile(t *testing.T) {
	clone := cloneWithRemotes(t, "")
	gitRun(t, clone, "config", "--local", "jaira.remote", "nowhere")
	t.Setenv("JAIRA_USER", "ada")

	if out, err := runCLI(t, clone, "init"); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	out, err := runCLI(t, clone, "create", "stays here")
	if err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	for _, want := range []string{"not on a ref", `"nowhere"`, "config jaira.remote"} {
		if !strings.Contains(out, want) {
			t.Errorf("create did not say %s:\n%s", want, out)
		}
	}
}

// The same fact for an agent: --json has to carry it as a field, because an
// agent that only reads "on-ref-only": false learns nothing it can act on.
func TestCreateJSONCarriesTheFileModeReason(t *testing.T) {
	clone := cloneWithRemotes(t, "")
	gitRun(t, clone, "config", "--local", "jaira.remote", "nowhere")
	t.Setenv("JAIRA_USER", "ada")

	if out, err := runCLI(t, clone, "init"); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	out, err := runCLI(t, clone, "--json", "create", "stays here")
	if err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("create --json did not emit one object: %v\n%s", err, out)
	}
	if payload["on-ref-only"] != false {
		t.Fatalf("on-ref-only is %v on a board with no usable remote", payload["on-ref-only"])
	}
	reason, _ := payload["file-only-reason"].(string)
	if !strings.Contains(reason, `"nowhere"`) {
		t.Errorf("file-only-reason does not name the remote looked for: %q", reason)
	}
}

// The way back. A ticket created while the board had no usable remote is a file
// and nothing else; once the remote is there, one command puts it on its ref and
// the file goes.
func TestReleasePutsAFileTicketOnItsRef(t *testing.T) {
	clone := cloneWithRemotes(t, "")
	gitRun(t, clone, "config", "--local", "jaira.remote", "later")
	t.Setenv("JAIRA_USER", "ada")

	if out, err := runCLI(t, clone, "init"); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if out, err := runCLI(t, clone, "create", "catch up"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	st, err := ticket.At(clone)
	if err != nil {
		t.Fatal(err)
	}
	all, err := st.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 1 {
		t.Fatalf("expected the ticket to be a file here, found %d", len(all))
	}
	id, path := all[0].ID, all[0].Path
	if ids, _ := (&gitref.Repo{Dir: clone}).List(); len(ids) != 0 {
		t.Fatalf("the ticket already has a ref: %v", ids)
	}

	// The remote arrives — the repository's own origin was there all along, it
	// was the name that was wrong.
	gitRun(t, clone, "config", "--local", "jaira.remote", "origin")
	if out, err := runCLI(t, clone, "release", ticket.Handle(id)); err != nil {
		t.Fatalf("release refused a ticket with no ref: %v\n%s", err, out)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("the local file is still there: %v", err)
	}
	ids, err := (&gitref.Repo{Dir: clone}).List()
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 1 || ids[0] != id {
		t.Fatalf("the ticket did not reach its ref: %v", ids)
	}
}

// The board's git state in one call, on a board that carries tickets on refs.
func TestWhoamiShowsTheBoardIsOnRefs(t *testing.T) {
	clone := cloneWithRemotes(t, "")
	t.Setenv("JAIRA_USER", "ada")

	h := ticketIn(t, clone, "on a ref")
	_ = h
	out, err := runCLI(t, clone, "whoami")
	if err != nil {
		t.Fatalf("whoami: %v\n%s", err, out)
	}
	for _, want := range []string{"Remote:", "origin", "Remotes:", "Ref mode:    yes", "File only:   0"} {
		if !strings.Contains(out, want) {
			t.Errorf("whoami does not say %q:\n%s", want, out)
		}
	}
}

// And on a board whose remote is not there: the mode, the reason, and how many
// tickets are lying on this disk alone.
func TestWhoamiShowsTheBoardIsInFileMode(t *testing.T) {
	clone := cloneWithRemotes(t, "")
	gitRun(t, clone, "config", "--local", "jaira.remote", "nowhere")
	t.Setenv("JAIRA_USER", "ada")

	if out, err := runCLI(t, clone, "init"); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if out, err := runCLI(t, clone, "create", "stays here"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	out, err := runCLI(t, clone, "--json", "whoami")
	if err != nil {
		t.Fatalf("whoami: %v\n%s", err, out)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("whoami --json did not emit one object: %v\n%s", err, out)
	}
	if payload["remote"] != "nowhere" {
		t.Errorf("whoami names %v as the remote", payload["remote"])
	}
	if payload["ref_mode"] != false {
		t.Errorf("whoami says the board is on refs without the remote it wants")
	}
	if n, _ := payload["file_only"].(float64); n != 1 {
		t.Errorf("whoami counted %v tickets on this disk only, want 1", payload["file_only"])
	}
	if reason, _ := payload["ref_mode_reason"].(string); !strings.Contains(reason, "config jaira.remote") {
		t.Errorf("the reason does not say what to do about it: %q", reason)
	}
}

// Acceptance criteria come in several, and --dod has to take several. Before
// this, a ticket created with three of them arrived carrying one, and the
// terminal lane's gate then ran against a third of the conditions.
func TestCreateTakesSeveralDoDItems(t *testing.T) {
	clone := cloneWithRemotes(t, "")
	t.Setenv("JAIRA_USER", "ada")

	if out, err := runCLI(t, clone, "init"); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if out, err := runCLI(t, clone, "create", "three boxes",
		"--dod", "the first thing is true",
		"--dod", "the second thing is true",
		"--dod", "the third thing is true"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	// Read the ticket off its ref: nobody pulled it, which is the point — a
	// board on refs must not make you assign a ticket to yourself to give it
	// its acceptance criteria.
	repo := &gitref.Repo{Dir: clone}
	ids, err := repo.List()
	if err != nil || len(ids) != 1 {
		t.Fatalf("expected one ticket on a ref, got %v (%v)", ids, err)
	}
	content, _, err := repo.Read(ids[0])
	if err != nil {
		t.Fatal(err)
	}
	tk, err := ticket.Decode(mustParse(t, content), filepath.Join(clone, "unused.md"))
	if err != nil {
		t.Fatal(err)
	}
	if len(tk.DoDItems) != 3 {
		t.Fatalf("the ticket carries %d criteria, want 3", len(tk.DoDItems))
	}
	if tk.DoD != "the first thing is true" {
		t.Errorf("the frontmatter summary is %q", tk.DoD)
	}
}

func mustParse(t *testing.T, content []byte) *ticket.Doc {
	t.Helper()
	d, err := ticket.ParseDoc(content)
	if err != nil {
		t.Fatal(err)
	}
	return d
}
