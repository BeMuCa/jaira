package gitrepo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// fixtureRepo builds a real git repository at test runtime rather than
// committing a nested .git as a fixture, which the repo's own test
// conventions treat as unsafe: git and most tooling handle a nested .git as
// ambiguous, and a runtime-built repo is the only way to get byte-identical
// behaviour to what git itself produces.
func fixtureRepo(t *testing.T) (dir string, repo *Repo) {
	t.Helper()
	if !Available() {
		t.Skip("git is not on PATH")
	}
	dir = t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
		}
	}
	run("init", "-q")
	run("config", "user.email", "fixture@example.com")
	run("config", "user.name", "Fixture")
	return dir, &Repo{Dir: dir}
}

// fixtureCommit writes files, stages, and commits them with an explicit,
// strictly increasing author/committer date, so the commit-date sort this
// package relies on has no ties and tests are deterministic rather than
// machine-speed-dependent. It returns the new commit's sha.
func fixtureCommit(t *testing.T, dir string, files map[string]string, message string, when time.Time) string {
	t.Helper()
	for name, content := range files {
		p := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatalf("mkdir for %s: %v", name, err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	addArgs := append([]string{"-C", dir, "add"}, names...)
	if out, err := exec.Command("git", addArgs...).CombinedOutput(); err != nil {
		t.Fatalf("git add: %v\n%s", err, out)
	}
	stamp := when.UTC().Format(time.RFC3339)
	cmd := exec.Command("git", "-C", dir, "commit", "-q", "-m", message)
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_DATE="+stamp, "GIT_COMMITTER_DATE="+stamp)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git commit: %v\n%s", err, out)
	}
	out, err := exec.Command("git", "-C", dir, "rev-parse", "HEAD").Output()
	if err != nil {
		t.Fatalf("git rev-parse HEAD: %v", err)
	}
	return strings.TrimSpace(string(out))
}

func TestCommitsForTicketUnionOldestFirstNoDuplicates(t *testing.T) {
	dir, repo := fixtureRepo(t)
	id := ticket.NewID(time.Now())
	ticketPath := "tickets/" + id + "-t.md"
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	// A: message only, does not touch the ticket file — the oldest commit, so
	// a naive "follow-list then grep-list" concatenation (which would put the
	// file-touching commits first) fails this test.
	a := fixtureCommit(t, dir, map[string]string{"unrelated.txt": "1"},
		"fix: mentions "+id+" in passing", base)
	// B: touches the ticket file only, unrelated message.
	b := fixtureCommit(t, dir, map[string]string{ticketPath: "v1"},
		"chore: bump ticket body", base.Add(time.Minute))
	// C: touches the ticket file and names it in the message.
	c := fixtureCommit(t, dir, map[string]string{ticketPath: "v2"},
		"feat: implements "+id, base.Add(2*time.Minute))
	// D: unrelated file, unrelated message — must be absent from the result.
	fixtureCommit(t, dir, map[string]string{"other.txt": "2"},
		"chore: something else entirely", base.Add(3*time.Minute))

	got, err := repo.CommitsForTicket(ticketPath, id)
	if err != nil {
		t.Fatalf("CommitsForTicket: %v", err)
	}
	want := []string{a, b, c}
	if !equalSlices(got, want) {
		t.Errorf("CommitsForTicket() = %v, want %v (oldest-first, deduped, no unrelated commit)", got, want)
	}
}

func TestCommitsForTicketHandleMatchedAsBoundedReference(t *testing.T) {
	dir, repo := fixtureRepo(t)
	id := ticket.NewID(time.Now())
	ticketPath := "tickets/" + id + "-t.md"
	handle := ticket.Handle(id)
	base := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	// Seed the ticket file so the repo is not empty and the file exists.
	seed := fixtureCommit(t, dir, map[string]string{ticketPath: "v1"}, "chore: create ticket", base)

	// A commit naming the handle in a scope, touching no ticket file, IS
	// found — the handle is the only form of the id a person ever writes.
	scoped := fixtureCommit(t, dir, map[string]string{"scoped.txt": "1"},
		fmt.Sprintf("fix(%s): narrow the lane", handle), base.Add(time.Minute))

	// A commit whose message embeds the same six characters inside a longer
	// alphanumeric run must NOT match — the handle is a reference, not a
	// substring.
	fixtureCommit(t, dir, map[string]string{"substring.txt": "1"},
		fmt.Sprintf("refs X%sZ1 upstream", handle), base.Add(2*time.Minute))

	got, err := repo.CommitsForTicket(ticketPath, id)
	if err != nil {
		t.Fatalf("CommitsForTicket: %v", err)
	}
	want := []string{seed, scoped}
	if !equalSlices(got, want) {
		t.Errorf("CommitsForTicket() = %v, want %v (handle matched as a bounded reference, not a substring)", got, want)
	}
}

func TestCommitsForTicketHandleMatchedCaseInsensitively(t *testing.T) {
	dir, repo := fixtureRepo(t)
	id := ticket.NewID(time.Now())
	ticketPath := "tickets/" + id + "-t.md"
	handle := ticket.Handle(id)
	base := time.Date(2026, 2, 5, 0, 0, 0, 0, time.UTC)

	seed := fixtureCommit(t, dir, map[string]string{ticketPath: "v1"}, "chore: create ticket", base)

	// A commit naming the handle in lowercase IS found — matching is
	// case-insensitive because nobody types a ULID in its stored uppercase
	// Crockford base32 form, the same reason ticket.NormalizeIDPrefix exists.
	lower := fixtureCommit(t, dir, map[string]string{"lower.txt": "1"},
		fmt.Sprintf("fix(%s): tidy up", strings.ToLower(handle)), base.Add(time.Minute))

	got, err := repo.CommitsForTicket(ticketPath, id)
	if err != nil {
		t.Fatalf("CommitsForTicket: %v", err)
	}
	want := []string{seed, lower}
	if !equalSlices(got, want) {
		t.Errorf("CommitsForTicket() = %v, want %v (handle matched case-insensitively)", got, want)
	}
}

func TestCommitsForTicketFollowsRename(t *testing.T) {
	dir, repo := fixtureRepo(t)
	id := ticket.NewID(time.Now())
	oldPath := "tickets/" + id + "-old.md"
	newPath := "tickets/" + id + "-new.md"
	base := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)

	created := fixtureCommit(t, dir, map[string]string{oldPath: "v1"}, "chore: create", base)

	run := func(args ...string) {
		t.Helper()
		if out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
		}
	}
	run("mv", oldPath, newPath)
	stamp := base.Add(time.Minute).UTC().Format(time.RFC3339)
	cmd := exec.Command("git", "-C", dir, "commit", "-q", "-m", "chore: rename ticket file")
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_DATE="+stamp, "GIT_COMMITTER_DATE="+stamp)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git commit: %v\n%s", err, out)
	}
	renamed, err := exec.Command("git", "-C", dir, "rev-parse", "HEAD").Output()
	if err != nil {
		t.Fatalf("git rev-parse HEAD: %v", err)
	}

	got, err := repo.CommitsForTicket(newPath, id)
	if err != nil {
		t.Fatalf("CommitsForTicket: %v", err)
	}
	want := []string{created, strings.TrimSpace(string(renamed))}
	if !equalSlices(got, want) {
		t.Errorf("CommitsForTicket() after rename = %v, want %v (history followed through the rename)", got, want)
	}
}

func TestCommitsForTicketEmptyRepoOrUncommittedFile(t *testing.T) {
	dir, repo := fixtureRepo(t)
	id := ticket.NewID(time.Now())
	ticketPath := "tickets/" + id + "-t.md"

	// A repo with no commits at all.
	got, err := repo.CommitsForTicket(ticketPath, id)
	if err != nil {
		t.Fatalf("CommitsForTicket on an empty repo: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("CommitsForTicket on an empty repo = %v, want empty", got)
	}

	// A commit exists, but not for this ticket's file or id.
	fixtureCommit(t, dir, map[string]string{"other.md": "x"}, "chore: unrelated", time.Now())
	got, err = repo.CommitsForTicket(ticketPath, id)
	if err != nil {
		t.Fatalf("CommitsForTicket for an uncommitted ticket file: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("CommitsForTicket for an uncommitted ticket file = %v, want empty", got)
	}
}

func TestCommitsForTicketNoGit(t *testing.T) {
	t.Setenv("PATH", "")

	repo := &Repo{Dir: t.TempDir()}
	_, err := repo.CommitsForTicket("tickets/x.md", "01KZTT3XZ2YQBX93TTSR7BVRCT")
	if err != ErrNoGit {
		t.Errorf("CommitsForTicket with no git on PATH: err = %v, want ErrNoGit", err)
	}
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Both forms of the id are interpolated into an extended regular expression, so
// an id carrying metacharacters would match commits that have nothing to do
// with the ticket — and archive stamps what this returns without asking. A
// hand-edited ticket is the reachable path: Load accepts any id, and only
// 'jaira validate' reports a non-ULID one, as an error.
func TestCommitsForTicketRefusesAnIDItWouldNotHaveWritten(t *testing.T) {
	dir, repo := fixtureRepo(t)
	base := time.Date(2026, 8, 23, 10, 0, 0, 0, time.UTC)
	for i, msg := range []string{"unrelated one", "unrelated two", "unrelated three"} {
		fixtureCommit(t, dir, map[string]string{fmt.Sprintf("f%d.txt", i): "x"}, msg, base.Add(time.Duration(i)*time.Minute))
	}

	for _, id := range []string{".*", "^", "[A-Z]", "x|y", ""} {
		shas, err := repo.CommitsForTicket(filepath.Join(dir, ".jaira", "tickets", "x.md"), id)
		if err != nil {
			t.Fatalf("id %q: %v", id, err)
		}
		if len(shas) != 0 {
			t.Errorf("id %q derived %d commit(s) it has nothing to do with: %v", id, len(shas), shas)
		}
	}
}
