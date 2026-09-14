package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/hook"
	coreidentity "github.com/BeMuCa/jaira/core/identity"
	"github.com/BeMuCa/jaira/core/outbox"
	"github.com/BeMuCa/jaira/core/refsync"
	"github.com/BeMuCa/jaira/core/settings"
	"github.com/BeMuCa/jaira/core/ticket"
)

// refs is this process's link between the ticket store and the refs tickets
// travel on. It is a package variable for the same reason the store's Actor is
// set in one place: every command opens its store through openStore, and a
// command that had to remember to wire this up itself would be a command that
// silently stops carrying tickets to the team.
var refs *refsync.Syncer

// attachRefs makes every write through this store also queue the ticket for its
// ref.
func attachRefs(s *ticket.Store) {
	set := settings.Load()
	// The machine-wide setting in settings.json is only a default, and this is
	// the one place per process that resolves it against the repository actually
	// in front of us. Everything downstream reads the answer
	// off refs.Repo.Remote rather than resolving it again, so a command costs the
	// git calls once.
	refs = refsync.New(s, set.RemoteFor(s.Root), s.Actor)
	// "Me" includes the aliases a person recorded, so core/identity answers it
	// rather than this package holding a second opinion about who someone is.
	refs.IsMine = func(assignee string) bool {
		return assignee != "" && coreidentity.IsMe(s.Root, assignee)
	}
	s.Recorder = refs
	// And the same syncer supplies what the board can see without having a
	// file for it, so list, next and show are never half a board.
	s.Source = refs
	openedStore = s
}

// fireHook calls the user's script for a move or a claim, for delivery that
// does not wait for the other side to fetch. It is best effort by design: see
// core/hook.
func fireHook(name string, s *ticket.Store, t *ticket.Ticket) {
	script := settings.Load().Hook
	if script == "" || t == nil {
		return
	}
	hook.Run(script, hook.Event{
		Name: name, ID: t.ID, Title: t.Title, Status: t.Status,
		Assignee: t.Assignee, Actor: s.Actor, Root: s.Root,
	})
}

// fileOnRefOnly sends a ticket to its ref and, if the remote accepted it, takes
// the local file away again.
//
// This is the one place the two storage modes are decided, and it is decided
// once: a board with a remote keeps an unworked ticket on its ref only, a board
// without one keeps the file, exactly as before. Every other command reads
// whichever is there and does not need to know which mode it is in.
//
// Why the file goes: while nobody is working a ticket, a file in somebody's
// checkout is a copy that a merge can duplicate and that hides who the ticket
// belongs to. 'jaira pull' is what brings it back, for exactly one clone.
//
// A ticket that could not be sent keeps its file. It is then correct here,
// carries the 'unsent' marker, and goes out with the next command — losing the
// file for a write that never left the machine would lose the ticket.
//
// The second return value is why the file stayed, and it is the whole point of
// the signature: a bare false is what let 'jaira create' choose the file mode in
// silence, so that seventeen tickets were written to disk and no command said
// so. The caller prints this.
func fileOnRefOnly(t *ticket.Ticket) (onRefOnly bool, why error) {
	if t == nil {
		return false, nil
	}
	if err := refs.Usable(); err != nil {
		return false, err
	}
	reports, err := refs.Flush()
	if err != nil {
		return false, err
	}
	for _, r := range reports {
		if r.ID != t.ID {
			continue
		}
		if r.Outcome != outbox.Sent {
			if r.Err != nil {
				return false, fmt.Errorf("the write is %s: %w", r.Outcome, r.Err)
			}
			return false, fmt.Errorf("the write is %s", r.Outcome)
		}
		if err := os.Remove(t.Path); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, errors.New("nothing was queued for this ticket")
}

// putOnRef is fileOnRefOnly for a ticket that already exists as a file: it
// queues the bytes that are on disk now and then goes through the same send and
// remove.
//
// It is the way back for a ticket created while the board had no usable remote.
// Nothing about that is specific to how the ticket got here — Record leases the
// empty string for a ticket with no ref, which is exactly "I expect this ref not
// to exist", the same compare-and-swap a freshly created ticket makes.
func putOnRef(t *ticket.Ticket) (onRefOnly bool, why error) {
	if t == nil {
		return false, nil
	}
	if err := refs.Usable(); err != nil {
		return false, err
	}
	content, err := os.ReadFile(t.Path)
	if err != nil {
		return false, err
	}
	if err := refs.Record(t.ID, content); err != nil {
		return false, err
	}
	return fileOnRefOnly(t)
}

// noRefReason renders why a ticket is not on a ref as the line a person can act
// on: which remote was looked for, and what git said about it. Every command
// that has to say it — 'jaira create' in its own words, 'jaira whoami', and the
// --json field beside them — says it from here, so the three cannot drift apart.
//
// The diagnostic itself is gitref's (Repo.NoRemoteHint names the remote, the
// remotes this repository does have, and the git config line that sets it).
// This only puts it where the state is created instead of at the end of the
// chain.
//
// It branches on the two halves of gitref.ErrNoRepo because they need opposite
// words. A repository whose remote is missing has a name that was looked for,
// remotes it does have, and a git config line that settles it — all of which
// gitref already writes. A directory that is no repository at all has none of
// that: naming a remote there is an invented problem, and the raw
// "gitref: no repository or no such remote" was the reader's only clue that the
// advice did not apply.
func noRefReason(why error) string {
	if why == nil {
		return "this board does not carry tickets on refs"
	}
	if errors.Is(why, gitref.ErrNoGitRepo) {
		return "this board is not in a git repository, so there is no remote to carry a ref"
	}
	if errors.Is(why, gitref.ErrNoGit) {
		return "git is not available on PATH, so nothing can be pushed to a ref"
	}
	// Everything a reader can act on is already in gitref's own sentence: the
	// remote name that was looked for, the remotes this repository really has,
	// and the git config line that settles it. So it is asked for by name —
	// Repo.NoRemoteHint — rather than cut back out of the error text, which tied
	// this line to how Usable concatenated it: a changed format would have missed
	// silently and left the reader with a line naming nothing. Nothing is wrapped
	// around it either: "no remote \"origin\" — this repository has no remotes"
	// already is the sentence, and a second framing in front of it only says the
	// same thing twice.
	if errors.Is(why, gitref.ErrNoRepo) && refs != nil && refs.Repo != nil {
		return refs.Repo.NoRemoteHint()
	}
	// Nothing to ask: no repo on this process (refsync.Syncer.Usable returns a
	// naked sentinel for a nil syncer), or an error from somewhere else. Say the
	// one thing still known instead of printing the sentinel at somebody.
	return "this board has no usable remote"
}

// canReachARef reports whether the board could carry tickets on refs once
// somebody configures a remote. It is false where there is no repository, and
// that is what keeps 'jaira create' from advising a command that cannot work
// there however the user answers.
func canReachARef(why error) bool {
	return !errors.Is(why, gitref.ErrNoGitRepo) && !errors.Is(why, gitref.ErrNoGit)
}

// openedStore is the store this command opened, remembered so the work that
// happens after the command — sending queued writes, and taking a backup when
// one is due — does not have to find the board a second time.
var openedStore *ticket.Store

// afterCommand is the background work that follows any command: send what this
// process queued, bring the refs up to date, and take the board's backup when
// one is due.
//
// None of it is on the command path. The flush only runs when this process
// actually wrote something, and the other two only decide — a file read — and
// leave the work to a detached child. That is what lets a read command stay
// instant while somebody who only ever reads still learns that a ticket was
// assigned to them.
func afterCommand() {
	flushRefs()
	maybeFetch(openedStore)
	maybeSnapshot(openedStore)
}

// maybeFetch spawns a detached 'jaira fetch' when the last one is older than
// the configured interval.
//
// Deliberately not inside 'jaira list': a read command must never wait for a
// remote. The fetch happens beside the command, in another process, and its
// result is there for the next one.
func maybeFetch(s *ticket.Store) {
	if s == nil || refs.Usable() != nil {
		return
	}
	refs.SpawnFetch(s.RepoStateDir(), s.Root, settings.Load().FetchInterval())
}

// resolveID turns whatever the user typed — a full id, a prefix, or the
// six-character handle the board prints everywhere — into the full id.
//
// The ref commands need this and the others do not: every other command reaches
// the ticket through the store, which resolves handles itself, while these
// address a ref by name and a ref named after a handle does not exist. Written
// as its own step rather than inside each command, because a handle failing for
// three commands and working for twenty is worse than it failing for all of
// them.
func resolveID(s *ticket.Store, arg string) string {
	if t, err := s.Load(arg); err == nil && t.ID != "" {
		return t.ID
	}
	return ticket.NormalizeIDPrefix(arg)
}

// flushRefs sends what this command queued, and is called once after the
// command has finished — including after it failed, because a ticket the
// command did manage to write is a ticket the team should see.
//
// It sends nothing, and touches no network, unless this process actually queued
// something. A read command therefore stays entirely off the network, which is
// what keeps 'jaira list' instant.
func flushRefs() {
	if !refs.Dirty() {
		return
	}
	reports, err := refs.Flush()
	if err != nil {
		warnRef(map[string]any{"error": err.Error()},
			"jaira: warning: the ticket could not be sent to the remote: %v", err)
		return
	}
	for _, r := range reports {
		switch r.Outcome {
		case outbox.Sent:
			// The normal case says nothing. A line on every write would be
			// noise on the one path every command takes.
		case outbox.Rejected:
			line := "someone else wrote it first"
			if r.Winner != nil {
				line = r.Winner.Describe()
			}
			warnRef(map[string]any{"ticket": r.ID, "outcome": "rejected", "winner": r.Winner},
				"jaira: %s was not sent: %s\n  your file is unchanged; the board shows both sides once the ref is fetched",
				ticket.Handle(r.ID), line)
		case outbox.Unsent:
			warnRef(map[string]any{"ticket": r.ID, "outcome": "unsent"},
				"jaira: %s is written locally and waiting to be sent; it goes out with your next command",
				ticket.Handle(r.ID))
		case outbox.Failed:
			warnRef(map[string]any{"ticket": r.ID, "outcome": "failed", "error": errText(r.Err)},
				"jaira: %s could not be sent: %v", ticket.Handle(r.ID), r.Err)
		}
	}
}

// warnRef writes to stderr in whichever shape the caller asked for. Never
// stdout: a command's stdout is its result, and an agent parsing --json output
// must not find a status line in the middle of it.
func warnRef(payload map[string]any, format string, args ...any) {
	if g.jsonOut {
		if b, err := json.Marshal(payload); err == nil {
			fmt.Fprintf(os.Stderr, "%s\n", b)
			return
		}
	}
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func errText(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
