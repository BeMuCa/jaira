package cli

import (
	"encoding/json"
	"fmt"
	"os"

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
	refs = refsync.New(s, set.RemoteName(), s.Actor)
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

// fileOnRefOnly sends a freshly created ticket to its ref and, if the remote
// accepted it, takes the local file away again.
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
func fileOnRefOnly(s *ticket.Store, t *ticket.Ticket) (onRefOnly bool) {
	if refs.Usable() != nil || t == nil {
		return false
	}
	reports, err := refs.Flush()
	if err != nil {
		return false
	}
	for _, r := range reports {
		if r.ID != t.ID {
			continue
		}
		if r.Outcome != outbox.Sent {
			return false
		}
		if err := os.Remove(t.Path); err != nil {
			return false
		}
		return true
	}
	return false
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
	refs.SpawnFetch(s.StateDir(), s.Root, settings.Load().FetchInterval())
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
