package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/ticket"
)

// secondClone clones the board's bare repository again — a teammate's
// checkout, or the main checkout beside a worktree — and gives it a board of its
// own. It has fetched nothing the first clone created after this moment.
func secondClone(t *testing.T, first string) string {
	t.Helper()
	bare := strings.TrimSpace(gitOut(t, first, "remote", "get-url", "origin"))
	other := filepath.Join(filepath.Dir(first), "other-"+filepath.Base(t.TempDir()))
	gitRun(t, filepath.Dir(first), "clone", "--quiet", bare, other)
	gitRun(t, other, "config", "user.name", "grace")
	gitRun(t, other, "config", "user.email", "grace@example.test")
	if out, err := runCLI(t, other, "init"); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	return other
}

// onlyRefID is the full id of the one ticket ref in the clone.
func onlyRefID(t *testing.T, clone string) string {
	t.Helper()
	ids, err := (&gitref.Repo{Dir: clone}).List()
	if err != nil || len(ids) != 1 {
		t.Fatalf("ticket refs: %v (%v)", ids, err)
	}
	return ids[0]
}

// remoteTicket reads the ticket the remote's ref carries, through a fresh fetch.
func remoteTicket(t *testing.T, clone, id string) *ticket.Ticket {
	t.Helper()
	r := &gitref.Repo{Dir: clone}
	if err := r.Fetch(); err != nil {
		t.Fatal(err)
	}
	content, _, err := r.Read(id)
	if err != nil {
		t.Fatalf("read ref %s: %v", id, err)
	}
	d, err := ticket.ParseDoc(content)
	if err != nil {
		t.Fatal(err)
	}
	tk, err := ticket.Decode(d, "")
	if err != nil {
		t.Fatal(err)
	}
	return tk
}

// filesUnder lists every .md under .jaira — tickets, archive, logbook — so a
// test can say no file was written anywhere in the checkout.
func filesUnder(t *testing.T, dir string) []string {
	t.Helper()
	var out []string
	_ = filepath.Walk(filepath.Join(dir, ".jaira"), func(p string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(p, ".md") && strings.HasPrefix(filepath.Base(p), "01") {
			out = append(out, p)
		}
		return nil
	})
	return out
}

// The HS3K1E accident: the ticket was created in another checkout, this one's
// refs had not seen it, and it never was a file here. Rejecting it by its short
// handle now fetches, writes the rejection onto the ref and leaves no file in
// this checkout; the creator's board drops it after a fetch and show still
// says why.
func TestRejectATicketThatLivesOnlyOnAStaleRef(t *testing.T) {
	for _, form := range []string{"handle", "ulid"} {
		t.Run(form, func(t *testing.T) {
			t.Setenv("JAIRA_USER", "ada")
			creator := cloneWithRemotes(t, "origin")
			h := ticketIn(t, creator, "nobody needs this")
			other := secondClone(t, creator)
			id := onlyRefID(t, creator)
			arg := h
			if form == "ulid" {
				arg = id
			}

			out, err := runCLI(t, other, "reject", arg, "--why", "not  worth\nit", "--force")
			if err != nil {
				t.Fatalf("reject: %v\n%s", err, out)
			}
			if !strings.Contains(out, "next 'jaira fetch'") {
				t.Errorf("output does not say where the rejection went:\n%s", out)
			}
			if got := filesUnder(t, other); len(got) != 0 {
				t.Errorf("reject left a file in the checkout: %v", got)
			}
			r := remoteTicket(t, creator, id).Rejected
			if r.Kind != ticket.RejectWontfix || r.Why != "not worth it" || r.By == "" || r.At.IsZero() {
				t.Errorf("the ref does not carry the rejection in full: %+v", r)
			}

			// The creator fetches: the ticket is off the board, and show still
			// tells what it was and why it will not be done.
			if out, err := runCLI(t, creator, "fetch"); err != nil {
				t.Fatalf("fetch: %v\n%s", err, out)
			}
			s, err := ticket.At(creator)
			if err != nil {
				t.Fatal(err)
			}
			all, _ := s.List()
			for _, tk := range all {
				if tk.ID == id {
					t.Errorf("the rejected ticket is still on the creator's board")
				}
			}
			out, err = runCLI(t, creator, "show", h)
			if err != nil || !strings.Contains(out, "not worth it") {
				t.Errorf("show does not say why it was rejected: %v\n%s", err, out)
			}
		})
	}
}

// A ticket somebody pulled into this branch is rejected the way it always was:
// archived here, where its file is.
func TestRejectATicketWithAFileHereStillArchivesIt(t *testing.T) {
	t.Setenv("JAIRA_USER", "ada")
	clone := cloneWithRemotes(t, "origin")
	h := ticketIn(t, clone, "pulled and dropped")
	if out, err := runCLI(t, clone, "pull", h); err != nil {
		t.Fatalf("pull: %v\n%s", err, out)
	}
	out, err := runCLI(t, clone, "reject", h, "--why", "changed our mind", "--force", "--json")
	if err != nil {
		t.Fatalf("reject: %v\n%s", err, out)
	}
	var got map[string]any
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("not json: %v\n%s", err, out)
	}
	if got["path"] == "" || got["on-ref-only"] == true {
		t.Errorf("a ticket with a file was not archived here: %v", got)
	}
	id := onlyRefID(t, clone)
	if r := archivedTicket(t, clone, id).Rejected; r.Why != "changed our mind" {
		t.Errorf("archived without the rejection: %+v", r)
	}
}

// installRaceHook makes the next push of the ticket's ref lose a race: a
// pre-push hook pushes another commit onto the same ref first, after the
// command has fetched and read its lease — the window a teammate's write lands
// in.
func installRaceHook(t *testing.T, clone, id string) {
	t.Helper()
	ref := gitref.RefName(id)
	guard := filepath.Join(t.TempDir(), "raced")
	script := "#!/bin/sh\n" +
		"[ -f '" + guard + "' ] && exit 0\n" +
		"touch '" + guard + "'\n" +
		"tree=$(git rev-parse " + ref + "^{tree}) || exit 1\n" +
		"c=$(git commit-tree \"$tree\" -p " + ref + " -m race) || exit 1\n" +
		"git push --quiet --no-verify origin \"$c:" + ref + "\" || exit 1\n"
	hooks := filepath.Join(clone, ".git", "hooks")
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(hooks, "pre-push"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
}

// Every way a reject of a ref-only ticket can fail says so in one line and
// with an exit code a script can branch on — and writes nothing.
func TestRejectOnTheRefRefusals(t *testing.T) {
	cases := []struct {
		name   string
		setup  func(t *testing.T, clone, h, id string) string // returns the argument
		code   int
		reason string
		says   string
	}{
		{
			name: "unknown ticket",
			setup: func(t *testing.T, clone, h, id string) string {
				return "ZZZZZZ"
			},
			code: ExitNotFound, says: "ZZZZZZ",
		},
		{
			name: "already rejected on its ref",
			setup: func(t *testing.T, clone, h, id string) string {
				if out, err := runCLI(t, clone, "reject", h, "--why", "first no", "--force"); err != nil {
					t.Fatalf("first reject: %v\n%s", err, out)
				}
				return h
			},
			code: ExitValidation, reason: "already_rejected", says: "first no",
		},
		{
			name: "already rejected and archived here",
			setup: func(t *testing.T, clone, h, id string) string {
				if out, err := runCLI(t, clone, "pull", h); err != nil {
					t.Fatalf("pull: %v\n%s", err, out)
				}
				if out, err := runCLI(t, clone, "reject", h, "--why", "first no", "--force"); err != nil {
					t.Fatalf("first reject: %v\n%s", err, out)
				}
				return id
			},
			code: ExitValidation, reason: "already_rejected", says: "first no",
		},
		{
			name: "archived here, not rejected",
			setup: func(t *testing.T, clone, h, id string) string {
				if out, err := runCLI(t, clone, "pull", h); err != nil {
					t.Fatalf("pull: %v\n%s", err, out)
				}
				if out, err := runCLI(t, clone, "archive", h); err != nil {
					t.Fatalf("archive: %v\n%s", err, out)
				}
				return h
			},
			code: ExitValidation, reason: "filed_away", says: "jaira restore",
		},
		{
			name: "remote unreachable",
			setup: func(t *testing.T, clone, h, id string) string {
				gitRun(t, clone, "remote", "set-url", "origin", filepath.Join(t.TempDir(), "gone.git"))
				return h
			},
			code: ExitError, reason: "no_remote", says: "nothing was written",
		},
		{
			name: "remote unreachable, ticket not known here",
			setup: func(t *testing.T, clone, h, id string) string {
				gitRun(t, clone, "remote", "set-url", "origin", filepath.Join(t.TempDir(), "gone.git"))
				return "QQQQQQ"
			},
			code: ExitError, reason: "no_remote", says: "could not be asked",
		},
		{
			name: "ref moved while rejecting",
			setup: func(t *testing.T, clone, h, id string) string {
				installRaceHook(t, clone, id)
				return h
			},
			code: ExitValidation, reason: "race", says: "run the same command again",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("JAIRA_USER", "ada")
			clone := cloneWithRemotes(t, "origin")
			h := ticketIn(t, clone, "refused")
			id := onlyRefID(t, clone)
			arg := c.setup(t, clone, h, id)
			before := filesUnder(t, clone)

			out, err := runCLI(t, clone, "reject", arg, "--why", "second no", "--force")
			if exitCode(err) != c.code {
				t.Fatalf("exit %d, want %d: %v\n%s", exitCode(err), c.code, err, out)
			}
			var ce *codedError
			if c.reason != "" && (!errors.As(err, &ce) || ce.reason != c.reason) {
				t.Errorf("reason %v, want %q", err, c.reason)
			}
			if !strings.Contains(err.Error(), c.says) {
				t.Errorf("message %q does not say %q", err.Error(), c.says)
			}
			if after := filesUnder(t, clone); strings.Join(after, ",") != strings.Join(before, ",") {
				t.Errorf("a refused reject changed the checkout: %v -> %v", before, after)
			}
			if c.name == "ref moved while rejecting" {
				if r := remoteTicket(t, clone, id).Rejected; r.Kind != "" {
					t.Errorf("the lost race still wrote the rejection: %+v", r)
				}
			}
		})
	}
}

// A ticket rejected on its ref has no archive file for 'jaira restore'.
// Pulling it is the way back, and it takes the rejection off — saying so in
// the output and in a note — rather than bringing back a live card that still
// carries a no.
func TestPullingATicketRejectedOnItsRefTakesTheRejectionOff(t *testing.T) {
	t.Setenv("JAIRA_USER", "ada")
	clone := cloneWithRemotes(t, "origin")
	h := ticketIn(t, clone, "wanted after all")
	id := onlyRefID(t, clone)
	if out, err := runCLI(t, clone, "reject", h, "--why", "too early", "--duplicate-of", h+"x", "--force"); err == nil {
		t.Fatalf("a bad --duplicate-of was accepted:\n%s", out)
	}
	if out, err := runCLI(t, clone, "reject", h, "--why", "too early", "--force"); err != nil {
		t.Fatalf("reject: %v\n%s", err, out)
	}
	out, err := runCLI(t, clone, "pull", h)
	if err != nil {
		t.Fatalf("pull: %v\n%s", err, out)
	}
	if !strings.Contains(out, "rejection is taken off") || !strings.Contains(out, "too early") {
		t.Errorf("pull did not say it took the rejection off:\n%s", out)
	}
	if r := remoteTicket(t, clone, id); r.Rejected.Kind != "" || !strings.Contains(r.Body, "too early") {
		t.Errorf("the ref still carries the rejection, or no note says it went: %+v", r.Rejected)
	}
	s, err := ticket.At(clone)
	if err != nil {
		t.Fatal(err)
	}
	tk, err := s.Load(h)
	if err != nil || tk.ReadOnly || tk.Rejected.Kind != "" {
		t.Errorf("the pulled file: %v readOnly=%v rejected=%+v", err, tk != nil && tk.ReadOnly, tk)
	}
}

// An assignee on the ref does not say anybody pulled the ticket — tickets get
// one to leave the backlog — so it refuses nothing. The rejection takes it off,
// so pull, the way back, does not answer "taken"; the output warns that a
// checkout which did pull it now has a file behind the ref. There the card
// leaves the board once the rejection is pulled in, and anybody can bring the
// ticket back by pulling it into a checkout without its file.
func TestRejectOnTheRefTakesTheAssigneeOffAndLeavesEveryBoard(t *testing.T) {
	t.Setenv("JAIRA_USER", "ada")
	clone := cloneWithRemotes(t, "origin")
	h := ticketIn(t, clone, "assigned elsewhere")
	id := onlyRefID(t, clone)
	other := secondClone(t, clone)
	t.Setenv("JAIRA_USER", "grace")
	if out, err := runCLI(t, other, "pull", h); err != nil {
		t.Fatalf("grace pulls: %v\n%s", err, out)
	}

	t.Setenv("JAIRA_USER", "ada")
	out, err := runCLI(t, clone, "reject", h, "--why", "not needed", "--force")
	if err != nil {
		t.Fatalf("reject: %v\n%s", err, out)
	}
	if !strings.Contains(out, "assigned to grace") || !strings.Contains(out, "behind the ref") {
		t.Errorf("no warning about the assignee's file:\n%s", out)
	}
	r := remoteTicket(t, clone, id)
	if r.Rejected.Why != "not needed" || r.Assignee != "" {
		t.Errorf("the ref: rejection %+v, assignee %q — want the rejection and no assignee", r.Rejected, r.Assignee)
	}

	// Grace's checkout: refused as stale after a fetch, then pulled up to the
	// ref — and the card is gone from her board, while show still finds it.
	t.Setenv("JAIRA_USER", "grace")
	if out, err := runCLI(t, other, "fetch"); err != nil {
		t.Fatalf("grace fetches: %v\n%s", err, out)
	}
	if out, err := runCLI(t, other, "note", h, "still working"); !errors.Is(err, ticket.ErrStale) {
		t.Errorf("grace's write on the file behind the rejection was not refused as stale: %v\n%s", err, out)
	}
	out, err = runCLI(t, other, "pull", h)
	if err != nil {
		t.Fatalf("grace refreshes: %v\n%s", err, out)
	}
	if !strings.Contains(out, "off the board") {
		t.Errorf("pull did not say the refreshed ticket is rejected:\n%s", out)
	}
	gs, err := ticket.At(other)
	if err != nil {
		t.Fatal(err)
	}
	onBoard := func(s *ticket.Store) bool {
		all, _ := s.List()
		for _, tk := range all {
			if tk.ID == id {
				return true
			}
		}
		return false
	}
	if onBoard(gs) {
		t.Errorf("the rejected ticket is still a card on grace's board, file and all")
	}
	if tk, err := gs.Load(h); err != nil || tk.Rejected.Why != "not needed" {
		t.Errorf("show/Load on grace's side lost the rejection: %v %+v", err, tk)
	}

	// Anybody without a file brings it back: ada, who rejected it.
	t.Setenv("JAIRA_USER", "ada")
	out, err = runCLI(t, clone, "pull", h)
	if err != nil {
		t.Fatalf("ada pulls it back: %v\n%s", err, out)
	}
	if !strings.Contains(out, "rejection is taken off") {
		t.Errorf("pull did not take the rejection off:\n%s", out)
	}
	as, err := ticket.At(clone)
	if err != nil {
		t.Fatal(err)
	}
	if !onBoard(as) {
		t.Errorf("the pulled-back ticket is not on ada's board")
	}
}

// Off the board is not out of existence. A ticket blocked by one rejected on
// its ref still names a ticket that is there: validate does not call it
// missing, and a holder's file carrying the rejection cannot be claimed as
// work nobody sees.
func TestATicketRejectedOnItsRefStillExistsForLinksAndCannotBeClaimed(t *testing.T) {
	t.Setenv("JAIRA_USER", "ada")
	clone := cloneWithRemotes(t, "origin")
	b := ticketIn(t, clone, "blocker")
	if out, err := runCLI(t, clone, "create", "blocked one"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	ids, err := (&gitref.Repo{Dir: clone}).List()
	if err != nil || len(ids) != 2 {
		t.Fatalf("refs: %v %v", ids, err)
	}
	a := ticket.Handle(ids[0])
	if a == b {
		a = ticket.Handle(ids[1])
	}
	if out, err := runCLI(t, clone, "pull", a); err != nil {
		t.Fatalf("pull: %v\n%s", err, out)
	}
	bID := ids[0]
	if ticket.Handle(bID) != b {
		bID = ids[1]
	}
	if out, err := runCLI(t, clone, "set", a, "blocked-by="+bID); err != nil {
		t.Fatalf("set blocked-by: %v\n%s", err, out)
	}
	if out, err := runCLI(t, clone, "reject", b, "--why", "not needed", "--force"); err != nil {
		t.Fatalf("reject: %v\n%s", err, out)
	}
	// Claiming it, with no file here, names only the pull: there is no file
	// to delete.
	out, err := runCLI(t, clone, "claim", b)
	if exitCode(err) != ExitValidation || !strings.Contains(fmt.Sprint(err), "'jaira pull "+b+"'") ||
		strings.Contains(fmt.Sprint(err), "delete its file") {
		t.Errorf("claim of a ticket rejected only on its ref: exit %d %v\n%s", exitCode(err), err, out)
	}
	out, err = runCLI(t, clone, "validate")
	if strings.Contains(out+fmt.Sprint(err), "exists nowhere") {
		t.Errorf("validate calls a ref-rejected blocker missing: %v\n%s", err, out)
	}
	// And the links view names it rejected, not "on its ref · jaira pull":
	// that pull would take the rejection off.
	out, err = runCLI(t, clone, "links", a)
	if err != nil || !strings.Contains(out, "rejected · off the board") || strings.Contains(out, "jaira pull "+b) {
		t.Errorf("links on a ticket blocked by a rejected one: %v\n%s", err, out)
	}

	// A holder's file that pulled the rejection in.
	if out, err := runCLI(t, clone, "create", "held elsewhere"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	all, _ := (&gitref.Repo{Dir: clone}).List()
	c := ""
	for _, id := range all {
		if h := ticket.Handle(id); h != a && h != b {
			c = h
		}
	}
	other := secondClone(t, clone)
	t.Setenv("JAIRA_USER", "grace")
	if out, err := runCLI(t, other, "pull", c); err != nil {
		t.Fatalf("grace pull: %v\n%s", err, out)
	}
	t.Setenv("JAIRA_USER", "ada")
	if out, err := runCLI(t, clone, "reject", c, "--why", "also not", "--force"); err != nil {
		t.Fatalf("reject c: %v\n%s", err, out)
	}
	t.Setenv("JAIRA_USER", "grace")
	if out, err := runCLI(t, other, "fetch"); err != nil {
		t.Fatalf("fetch: %v\n%s", err, out)
	}
	if out, err := runCLI(t, other, "pull", c); err != nil {
		t.Fatalf("grace refresh: %v\n%s", err, out)
	}
	out, err = runCLI(t, other, "claim", c)
	var ce *codedError
	if exitCode(err) != ExitValidation || !errors.As(err, &ce) || ce.reason != "rejected" ||
		!strings.Contains(ce.message, "delete its file here") {
		t.Errorf("claim of a rejected file: exit %d %v\n%s", exitCode(err), err, out)
	}
}
