package gitrepo

import (
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
