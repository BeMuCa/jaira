package gitrepo

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// HuntCheck is what the gate reads before signoff (G8HNHT): a ticket whose
// commits carry only board data changed no code; one with a code commit and
// no hunt has none; a hunt at HEAD covers it, and a commit of the ticket on a
// hunted file after it leaves it failing.
func TestHuntCheck(t *testing.T) {
	const id = "01M7ZYEXTPHSC0PNW5STG8HNHT"
	at := time.Date(2026, 11, 21, 8, 0, 0, 0, time.UTC)
	dir, repo := fixtureRepo(t)
	tk := &ticket.Ticket{ID: id, Path: ".jaira/tickets/" + id + ".md"}

	fixtureCommit(t, dir, map[string]string{"a.go": "a"}, "start", at)
	fixtureCommit(t, dir, map[string]string{tk.Path: "x", ".jaira/tags": "gates"}, "chore(G8HNHT): board only", at.Add(time.Minute))
	if got := repo.HuntCheck(tk); got.CodeChanged {
		t.Errorf("board data only: %+v, want no code changed", got)
	}

	head := fixtureCommit(t, dir, map[string]string{"a.go": "a2"}, "fix(G8HNHT): a", at.Add(2*time.Minute))
	got := repo.HuntCheck(tk)
	if want := (HuntReport{CodeChanged: true, Head: head}); got != want {
		t.Errorf("code, no hunt: %+v, want %+v", got, want)
	}

	tk.Body = ticket.NoteHeading + "\n- **2026-11-21 08:03 · t** — bughunt: done at " + head[:9] + " - 0 findings; files: a.go\n"
	if got := repo.HuntCheck(tk); got.State != "ok" || got.At != head[:9] {
		t.Errorf("hunt at HEAD: %+v, want ok at %s", got, head[:9])
	}

	fixtureCommit(t, dir, map[string]string{"a.go": "a3"}, "fix(G8HNHT): a again", at.Add(3*time.Minute))
	if got := repo.HuntCheck(tk); got.State != "fail" {
		t.Errorf("commit after the hunt: %+v, want fail", got)
	}
}

// A board may live in a subdirectory of its repository, and git names paths
// from the root: sub/.jaira/tickets/X.md is board data all the same. A commit
// carrying only the ticket file changed no code, and a hunt with no file
// list is not made stale by one; code outside the board's directory is code.
func TestHuntCheckBoardInSubdirectory(t *testing.T) {
	const id = "01M7ZYEXTPHSC0PNW5STG8HNHT"
	at := time.Date(2026, 11, 21, 8, 0, 0, 0, time.UTC)
	dir, _ := fixtureRepo(t)
	repo := &Repo{Dir: filepath.Join(dir, "sub")}
	tk := &ticket.Ticket{ID: id, Path: filepath.Join(dir, "sub", ".jaira", "tickets", id+".md")}

	start := fixtureCommit(t, dir, map[string]string{"sub/a.go": "a"}, "start", at)
	fixtureCommit(t, dir, map[string]string{"sub/.jaira/tickets/" + id + ".md": "x"}, "chore(G8HNHT): board only", at.Add(time.Minute))
	if got := repo.HuntCheck(tk); got.CodeChanged {
		t.Errorf("ticket file only, board in sub/: %+v, want no code changed", got)
	}
	if got, err := repo.HuntCoverage(start, id, nil); err != nil || got != HuntCovers {
		t.Errorf("HuntCoverage over all code after a ticket-file commit = %v, %v; want covers", got, err)
	}

	fixtureCommit(t, dir, map[string]string{"lib.go": "l"}, "fix(G8HNHT): outside the board's directory", at.Add(2*time.Minute))
	if got := repo.HuntCheck(tk); !got.CodeChanged {
		t.Errorf("code outside sub/: %+v, want code changed", got)
	}
}
