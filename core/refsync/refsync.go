// Package refsync connects the ticket store to the ref a ticket travels on.
//
// It is the only place that knows all three: core/ticket writes the file,
// core/gitref carries it on refs/jaira/tickets/<id>, core/outbox holds what has
// not been sent. Keeping that knowledge here is what lets core/ticket stay
// unaware of git and core/gitref stay unaware of tickets.
//
// Nothing here pushes on the write path. A ticket write is recorded into the
// outbox and returns immediately, because the write path runs inside the TUI
// and inside the git merge driver as well as in the CLI, and none of those may
// wait for a network round trip. Sending is a separate act (Flush), run where
// waiting is acceptable.
package refsync

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/merge"
	"github.com/BeMuCa/jaira/core/outbox"
	"github.com/BeMuCa/jaira/core/ticket"
)

// Syncer records ticket writes for the remote and sends them.
type Syncer struct {
	Repo *gitref.Repo
	Box  *outbox.Box

	// Actor is who the queued write is from, so a flush can say whose write it
	// is sending and a rejection can name both sides.
	Actor string

	// Store is the working tree's tickets, used only to tell a ticket that
	// arrived on a ref alone from one this clone already has a file for.
	Store *ticket.Store

	// IsMine decides whether an assignee is this user. It is a function rather
	// than a name because "me" includes the aliases a person has recorded, and
	// core/identity already knows how to answer it — this package should not
	// hold a second opinion.
	IsMine func(assignee string) bool

	// SeenPath records which ref SHA was last reported for each ticket, so
	// being handed a ticket is announced once instead of on every fetch.
	SeenPath string

	// usable caches whether this checkout can carry refs at all. It is asked
	// once per process rather than per write: the answer is a property of the
	// checkout, and every ticket write would otherwise pay two git invocations
	// to learn the same thing.
	once   sync.Once
	usable error

	// dirty records that this process queued something. It is what lets a read
	// command stay entirely off the network: only a command that actually
	// wrote a ticket has a reason to send, and jaira list must never wait for
	// a remote.
	dirty bool
}

// Dirty reports whether this process queued a write. A caller uses it to decide
// whether flushing is worth a round trip at all.
func (y *Syncer) Dirty() bool { return y != nil && y.dirty }

// Usable reports whether this board carries tickets on refs. A board with no
// repository or no remote is not an error — it is jaira used the way it also
// works, and the answer is simply no.
func (y *Syncer) Usable() error {
	if y == nil {
		return gitref.ErrNoRepo
	}
	y.once.Do(func() { y.usable = y.Repo.Usable() })
	return y.usable
}

// New builds a syncer for a store, or returns nil when this board has no
// business talking to a remote: no git, no repository, or no remote configured.
// A nil *Syncer is a valid recorder that records nothing, so callers do not
// need to branch.
func New(s *ticket.Store, remote, actor string) *Syncer {
	repo := &gitref.Repo{Dir: s.Root, Remote: remote, AuthorName: actor}
	return &Syncer{
		Repo:     repo,
		Box:      outbox.At(s),
		Actor:    actor,
		Store:    s,
		SeenPath: filepath.Join(s.StateDir(), "refs-seen.json"),
	}
}

// Record queues the ticket's current bytes for its ref. It satisfies
// ticket.WriteRecorder.
//
// The lease is the local ref's SHA, which is only ever moved after the remote
// has accepted a push — so it is exactly "the state the remote last confirmed
// to us", which is what a compare-and-swap has to be checked against. A ticket
// with no ref yet leases the empty string, meaning "I expect this ref not to
// exist".
func (y *Syncer) Record(id string, content []byte) error {
	if y == nil || y.Usable() != nil {
		// A board without a repository or without a remote does not carry
		// tickets on refs. Saying so on every single write would be noise, and
		// queueing a write that can never be sent would be worse: the card
		// would show "unsent" forever.
		return nil
	}
	lease, err := y.Repo.SHA(id)
	if err != nil && !errors.Is(err, gitref.ErrNoRef) {
		return err
	}
	if err := y.Box.Queue(id, outbox.OpWrite, content, lease, y.Actor); err != nil {
		return err
	}
	y.dirty = true
	return nil
}

// RecordDelete queues the removal of a ticket's ref, for a ticket that has been
// logged or archived. It goes through the same queue as a write so that taking
// a ticket off the board works offline too — otherwise the last act of a
// finished ticket would be the one act that needs a network.
func (y *Syncer) RecordDelete(id string) error {
	if y == nil || y.Usable() != nil {
		return nil
	}
	lease, err := y.Repo.SHA(id)
	if err != nil {
		if errors.Is(err, gitref.ErrNoRef) {
			// Nothing to delete: the ticket never reached a ref.
			return nil
		}
		return err
	}
	if err := y.Box.Queue(id, outbox.OpDelete, nil, lease, y.Actor); err != nil {
		return err
	}
	y.dirty = true
	return nil
}

// RecordFiled queues a ticket's final state for its ref, for a ticket that has
// been logged or archived here.
//
// Filing does not remove the ref, and that is the point. At this moment the
// ticket file exists only in this clone's branch; until that branch is merged,
// nobody else can see the ticket. Taking the ref down as well would open a
// window as long as a review, in which the ticket is invisible to everybody —
// and somebody who notices the same problem writes it down a second time,
// with the same solution.
//
// The ref goes later, when the ticket has arrived where everybody can see it.
// That is the snapshot run's job (core/snapshot), which removes a ref only
// after writing the ticket into the snapshot branch, so at the moment of
// deletion the ticket exists in two places.
func (y *Syncer) RecordFiled(id string, content []byte) error {
	return y.Record(id, content)
}

// Winner describes the state that beat a rejected write, read from the ref
// itself. A rejection that only says "you lost" leaves the user to go and find
// out who and what — which they cannot do without knowing the ref exists.
type Winner struct {
	ID        string `json:"id"`
	Assignee  string `json:"assignee,omitempty"`
	Status    string `json:"status,omitempty"`
	UpdatedBy string `json:"updated-by,omitempty"`
	UpdatedAt string `json:"updated-at,omitempty"`
}

// Describe renders the winner as the one line a command can print.
func (w Winner) Describe() string {
	who := w.UpdatedBy
	if who == "" {
		who = w.Assignee
	}
	if who == "" {
		who = "someone else"
	}
	parts := []string{who + " wrote it first"}
	if w.Status != "" {
		parts = append(parts, "it is now in "+w.Status)
	}
	if w.Assignee != "" && !strings.EqualFold(w.Assignee, who) {
		parts = append(parts, "assigned to "+w.Assignee)
	}
	if w.UpdatedAt != "" {
		parts = append(parts, "at "+w.UpdatedAt)
	}
	return strings.Join(parts, ", ")
}

// Holds renders the same state as a statement of ownership rather than of a
// race. The two refusals are different facts and must not borrow each other's
// wording: "grace wrote it first" is about who won a push, "grace has it" is
// about who the ticket belongs to.
func (w Winner) Holds() string {
	who := w.Assignee
	if who == "" {
		who = w.UpdatedBy
	}
	if who == "" {
		who = "somebody else"
	}
	out := who + " has it"
	if w.Status != "" {
		out += ", it is in " + w.Status
	}
	return out
}

// Report is one queued write's fate, with the other side named when we lost.
type Report struct {
	outbox.Result
	Winner *Winner `json:"winner,omitempty"`
}

// Flush sends what the outbox holds and reports each entry, naming who was
// quicker for every write the remote refused.
//
// The fetch comes first, and deliberately: a queued write leases a SHA this
// clone has seen, so without fetching, the very first push would be refused
// for a race that has in fact already been resolved on the remote. Fetching
// also makes the winner readable, which is the difference between "you lost"
// and "Grace moved it to review".
func (y *Syncer) Flush() ([]Report, error) {
	if y == nil || y.Usable() != nil {
		return nil, nil
	}
	entries, err := y.Box.List()
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, nil
	}
	if err := y.Repo.Fetch(); err != nil {
		// An unreachable remote is not an error here: it is the state the
		// outbox exists for. The queue stays as it is and the next command
		// tries again.
		if errors.Is(err, gitref.ErrOffline) || errors.Is(err, gitref.ErrNoGit) {
			return nil, nil
		}
		// A remote that exists but is not a repository is the opposite case:
		// no amount of waiting fixes it, and staying quiet would leave every
		// write queued with nobody told why.
		return nil, err
	}

	results, err := y.Box.Flush(y.Repo)
	if err != nil {
		return nil, err
	}
	reports := make([]Report, 0, len(results))
	for _, r := range results {
		rep := Report{Result: r}
		if r.Outcome == outbox.Rejected {
			if w, err := y.winner(r.ID); err == nil {
				rep.Winner = w
			}
		}
		reports = append(reports, rep)
	}
	return reports, nil
}

// Pending reports an unsent write for one ticket, for the marker on a board
// card.
func (y *Syncer) Pending(id string) (outbox.Entry, bool) {
	if y == nil {
		return outbox.Entry{}, false
	}
	return y.Box.Pending(id)
}

// PendingAge is how long this ticket's write has been waiting, or zero.
func (y *Syncer) PendingAge(id string) time.Duration {
	e, ok := y.Pending(id)
	if !ok {
		return 0
	}
	return e.Age()
}

// winner reads the state now on the ref.
func (y *Syncer) winner(id string) (*Winner, error) {
	content, _, err := y.Repo.Read(id)
	if err != nil {
		return nil, err
	}
	d, err := ticket.ParseDoc(content)
	if err != nil {
		return nil, fmt.Errorf("refsync: the ref for %s is not a ticket: %w", id, err)
	}
	w := &Winner{ID: id}
	for _, f := range []struct {
		key string
		dst *string
	}{
		{ticket.FieldAssignee, &w.Assignee},
		{ticket.FieldStatus, &w.Status},
		{ticket.FieldUpdatedBy, &w.UpdatedBy},
		{ticket.FieldUpdatedAt, &w.UpdatedAt},
	} {
		if v, ok, err := d.Scalar(f.key); err == nil && ok {
			*f.dst = v
		}
	}
	return w, nil
}

// Arrival is a ticket as the ref now has it, seen from this clone.
type Arrival struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Assignee string `json:"assignee"`
	SHA      string `json:"sha"`

	// Mine says the ticket is assigned to this user. It is the whole reason
	// this exists: being handed a ticket is the event nobody currently learns
	// about.
	Mine bool `json:"mine"`

	// Changed says the ref moved since this clone last looked. A ticket that
	// has been sitting there assigned to me for a week is not news, and
	// announcing it on every fetch would train the user to ignore the thing.
	Changed bool `json:"changed"`

	// Local says a ticket file for this id exists in the working tree. False
	// means the ticket reached this clone only on its ref — the case a
	// teammate's unmerged branch cannot deliver.
	Local bool `json:"local"`
}

// seen is the per-clone record of which ref SHA was last reported, so an
// arrival is announced once rather than on every fetch.
type seen map[string]string

// Incoming fetches the ticket refs and reports what they now say.
//
// It is the read half of the feature: a fetch of refs/jaira/tickets/* costs one
// round trip and needs no branch, no checkout and no merge, so a clone can be
// told "this is yours now" without anyone having pushed a branch.
func (y *Syncer) Incoming() ([]Arrival, error) {
	if y == nil || y.Usable() != nil {
		return nil, nil
	}
	if err := y.Repo.Fetch(); err != nil {
		if errors.Is(err, gitref.ErrOffline) || errors.Is(err, gitref.ErrNoGit) {
			return nil, nil
		}
		return nil, err
	}
	ids, err := y.Repo.List()
	if err != nil {
		return nil, err
	}
	known := y.loadSeen()
	next := make(seen, len(ids))
	var out []Arrival
	for _, id := range ids {
		content, sha, err := y.Repo.Read(id)
		if err != nil {
			// A ref whose tree is not a ticket is not worth failing a fetch
			// over; it is also not something this tool wrote.
			continue
		}
		next[id] = sha
		a := Arrival{ID: id, SHA: sha, Changed: known[id] != sha}
		d, err := ticket.ParseDoc(content)
		if err != nil {
			continue
		}
		for _, f := range []struct {
			key string
			dst *string
		}{
			{ticket.FieldTitle, &a.Title},
			{ticket.FieldStatus, &a.Status},
			{ticket.FieldAssignee, &a.Assignee},
		} {
			if v, ok, err := d.Scalar(f.key); err == nil && ok {
				*f.dst = v
			}
		}
		if y.IsMine != nil {
			a.Mine = y.IsMine(a.Assignee)
		}
		a.Local = y.hasLocal(id)
		out = append(out, a)
	}
	// A ref that has vanished while a file for it is still here is unfinished
	// business, not history: it is kept in the seen record so it can still be
	// reported afterwards. Forgetting it here is what made the report silent —
	// the fetch that should have raised it had already overwritten the only
	// evidence.
	local := y.localIDs()
	for id, sha := range known {
		if _, stillThere := next[id]; !stillThere && local[id] {
			next[id] = sha
		}
	}
	y.saveSeen(next)
	return out, nil
}

// isMe answers for the caller's own name, falling back to a plain comparison
// when nobody told this syncer how identity works.
func (y *Syncer) isMe(who string) bool {
	if y.IsMine != nil {
		return y.IsMine(who)
	}
	return strings.EqualFold(strings.TrimSpace(who), strings.TrimSpace(y.Actor))
}

// localPath is the file this working tree holds for a ticket, if any.
//
// Asked from the filenames rather than through Store.Load, and that is not a
// detail: Load answers for a ticket on a ref too, with no path, so using it to
// mean "is it here" would report every visible ticket as present.
func (y *Syncer) localPath(id string) (string, bool) {
	if y.Store == nil {
		return "", false
	}
	paths, err := y.Store.Paths()
	if err != nil {
		return "", false
	}
	for _, p := range paths {
		if ticket.IDFromFilename(filepath.Base(p)) == id {
			return p, true
		}
	}
	return "", false
}

// localIDs is the set of ticket ids this working tree holds a file for, read
// from the filenames alone.
func (y *Syncer) localIDs() map[string]bool {
	out := map[string]bool{}
	if y.Store == nil {
		return out
	}
	paths, err := y.Store.Paths()
	if err != nil {
		return out
	}
	for _, p := range paths {
		if id := ticket.IDFromFilename(filepath.Base(p)); id != "" {
			out[id] = true
		}
	}
	return out
}

// hasLocal reports whether the working tree holds a file for this ticket.
func (y *Syncer) hasLocal(id string) bool { return y.localIDs()[id] }

func (y *Syncer) loadSeen() seen {
	s := seen{}
	if y.SeenPath == "" {
		return s
	}
	b, err := os.ReadFile(y.SeenPath)
	if err != nil {
		return s
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return seen{}
	}
	return s
}

// saveSeen records what was just reported. A failure here is silent on purpose:
// the worst it costs is one repeated notification, and refusing a fetch over it
// would cost the feature.
func (y *Syncer) saveSeen(s seen) {
	if y.SeenPath == "" {
		return
	}
	if err := os.MkdirAll(filepath.Dir(y.SeenPath), 0o755); err != nil {
		return
	}
	if b, err := json.MarshalIndent(s, "", "  "); err == nil {
		_ = os.WriteFile(y.SeenPath, append(b, '\n'), 0o644)
	}
}

// Reconciled is one ticket as the board should show it: the local file and the
// ref brought together, with the two state markers the user needs.
type Reconciled struct {
	ID string `json:"id"`

	// Content is what to show. On a clean merge it is the two sides merged; a
	// ticket that exists on only one side is that side unchanged.
	Content []byte `json:"-"`

	// Conflicts are the prose fields the merge could not settle. They are the
	// existing conflict path, not a new one: the same fields 'jaira resolve'
	// handles after a git merge.
	Conflicts []merge.Conflict `json:"conflicts,omitempty"`

	// RefOnly says the ticket reached this clone on its ref alone — it is in
	// no branch here. This is the case the whole feature exists for.
	RefOnly bool `json:"ref-only"`

	// Unsent says this clone holds a write for the ticket that has not reached
	// the remote yet.
	Unsent bool `json:"unsent"`
}

// Reconcile brings the local ticket file and the ref together for one ticket.
//
// It is deliberately not "show whichever is newer". The two sides are merged
// field by field by core/merge, which is the same resolver the git merge driver
// uses: status by progress along the lane chain (never backwards), lists by
// union, other scalars by updated-at, prose as a real conflict. A newer local
// edit must not drag a ticket back out of review because someone touched it
// after the reviewer did.
//
// The base for that merge is the ref's parent blob — the state both sides
// started from — which is available precisely because a ref commit keeps its
// leased SHA as parent.
func (y *Syncer) Reconcile(id string, lanes *lane.Set) (*Reconciled, error) {
	out := &Reconciled{ID: id}
	if y == nil {
		return out, nil
	}
	_, out.Unsent = y.Pending(id)

	var local []byte
	if path, ok := y.localPath(id); ok {
		if b, err := os.ReadFile(path); err == nil {
			local = b
		}
	}
	if y.Usable() != nil {
		out.Content = local
		return out, nil
	}
	theirs, _, err := y.Repo.Read(id)
	if err != nil {
		// No ref: the local file is the whole story.
		out.Content = local
		return out, nil
	}
	if local == nil {
		out.RefOnly = true
		out.Content = theirs
		return out, nil
	}
	base, err := y.Repo.ReadParent(id)
	if err != nil {
		// The ticket's first write has no earlier state; the two sides are all
		// there is, and core/merge treats an empty base as "both added it".
		base = local
	}
	res, err := merge.Merge(base, local, theirs, lanes)
	if err != nil {
		// A merge that cannot be performed must not blank the card: showing
		// the local file is always defensible, since it is what this clone
		// wrote.
		out.Content = local
		return out, nil
	}
	out.Content = res.Merged
	out.Conflicts = res.Conflicts
	return out, nil
}

// ErrTaken means the ticket on the ref already belongs to somebody else.
//
// It is a separate refusal from a lost race, and both are needed. The
// compare-and-swap only rules out two writes landing in the same instant; it
// cannot say "this ticket is not yours", because a pull re-reads the ref
// immediately before writing and would therefore always hold a current lease.
// Ownership is a fact recorded in the ticket, so it is checked as one.
var ErrTaken = errors.New("refsync: the ticket already belongs to someone else")

// Pulled reports what a pull did.
type Pulled struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Path  string `json:"path,omitempty"`

	// AlreadyHere means the file was already on this machine and nothing
	// needed doing. A pull is safe to repeat, so this is a success.
	AlreadyHere bool `json:"already-here"`

	// Winner is set when the ref refused the take-over: somebody else got
	// there first, and no file was written.
	Winner *Winner `json:"winner,omitempty"`

	// TakenFrom names who the ticket was taken from, when it was taken with
	// --steal. A silent steal would leave the other person's board saying the
	// ticket is theirs, with nothing anywhere saying why it stopped being.
	TakenFrom string `json:"taken-from,omitempty"`
}

// Pull takes a ticket that lives on its ref and makes it this clone's to work
// on: it claims the ticket on the ref and only then writes the file here.
//
// Two things gate it, and neither is enough alone:
//
//   - The ref's assignee. A ticket somebody else has already pulled is
//     refused with ErrTaken, naming them. This is the check that makes "one
//     clone has the file" true, because a pull re-reads the ref right before
//     writing and so would otherwise always hold a current lease.
//   - The compare-and-swap on the push. That covers the remaining case the
//     first check cannot see: two clones pulling in the same instant, each
//     having read an unassigned ticket.
//
// The order is the other half of the mechanism. The push goes first. Writing
// the file first would leave the loser holding a ticket file that belongs to
// somebody else — the duplicate this design exists to make impossible rather
// than to resolve.
//
// This is also the one write in jaira that is synchronous and does not go
// through the outbox. That is not an exception to the offline rule but a
// consequence of it: with no route to the remote there is no ref to read, so
// there is nothing to take over in the first place.
func (y *Syncer) Pull(id string, steal bool) (*Pulled, error) {
	if y == nil || y.Store == nil {
		return nil, errors.New("refsync: no board here")
	}
	if err := y.Usable(); err != nil {
		return nil, err
	}
	// Already here: a repeated pull is a successful no-op, the same as moving a
	// ticket to the lane it is already in. The question is whether there is a
	// FILE, not whether the board can see the ticket — it can see every ref.
	if path, ok := y.localPath(id); ok {
		out := &Pulled{ID: id, Path: path, AlreadyHere: true}
		if t, err := y.Store.Load(id); err == nil {
			out.Title = t.Title
		}
		return out, nil
	}
	if err := y.Repo.Fetch(); err != nil {
		return nil, err
	}
	content, lease, err := y.Repo.Read(id)
	if err != nil {
		return nil, err
	}
	d, err := ticket.ParseDoc(content)
	if err != nil {
		return nil, fmt.Errorf("refsync: the ref for %s does not carry a ticket: %w", id, err)
	}
	holder, _, _ := d.Scalar(ticket.FieldAssignee)
	if strings.TrimSpace(holder) != "" && !y.isMe(holder) && !steal {
		out := &Pulled{ID: id}
		if w, wErr := y.winner(id); wErr == nil {
			out.Winner = w
		}
		return out, fmt.Errorf("%w: %s has it", ErrTaken, holder)
	}
	stolenFrom := ""
	if h := strings.TrimSpace(holder); h != "" && !y.isMe(h) {
		// Taken with --steal. Recorded on the ticket itself, because the ref is
		// the only thing the other person will fetch, and a note there is the
		// only place they can read why the ticket left them.
		stolenFrom = h
		if err := ticket.AppendNote(d, fmt.Sprintf("%s took this ticket over from %s", y.Actor, h), time.Now(), y.Actor); err != nil {
			return nil, err
		}
	}
	if err := d.SetScalar(ticket.FieldAssignee, y.Actor); err != nil {
		return nil, err
	}
	if err := ticket.Touch(d, time.Now()); err != nil {
		return nil, err
	}
	if err := ticket.TouchBy(d, y.Actor); err != nil {
		return nil, err
	}
	taken := d.Bytes()

	if _, err := y.Repo.Write(id, taken, lease); err != nil {
		if errors.Is(err, gitref.ErrRaceLost) {
			out := &Pulled{ID: id}
			if w, wErr := y.winner(id); wErr == nil {
				out.Winner = w
			}
			return out, err
		}
		return nil, err
	}

	title, _, _ := d.Scalar(ticket.FieldTitle)
	path := filepath.Join(y.Store.TicketsDir(), ticket.Filename(id, title))
	if err := os.MkdirAll(y.Store.TicketsDir(), 0o755); err != nil {
		return nil, err
	}
	// Written with the bytes the remote accepted, not with a freshly built
	// document: what is on the board here and what the team sees on the ref are
	// then the same file, byte for byte.
	if err := ticket.WriteAtomic(path, taken); err != nil {
		return nil, err
	}
	return &Pulled{ID: id, Title: title, Path: path, TakenFrom: stolenFrom}, nil
}

// Released reports what a release did.
type Released struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Removed string `json:"removed,omitempty"`

	// NotYours is set when the ticket is not this user's to release. Handing
	// somebody else's work back is not a thing one person decides for another.
	NotYours *Winner `json:"not-yours,omitempty"`
}

// Release hands a ticket back: it clears the assignee on the ref and removes
// the file from this clone, in that order.
//
// It is the other half of Pull, and without it an assignment is a reservation
// nobody can give back. The order matters for the same reason: the ref is
// cleared first, so a failure leaves the ticket still yours rather than
// unowned-but-still-on-your-disk.
//
// The local file is removed rather than kept, because keeping it would put a
// second copy of a ticket somebody else may now pull into a branch — the
// duplicate this whole design exists to make impossible. Nothing is lost: the
// ref carries the bytes, and the work already committed stays in the history.
func (y *Syncer) Release(id string, force bool) (*Released, error) {
	if y == nil || y.Store == nil {
		return nil, errors.New("refsync: no board here")
	}
	if err := y.Usable(); err != nil {
		return nil, err
	}
	if err := y.Repo.Fetch(); err != nil {
		return nil, err
	}
	content, lease, err := y.Repo.Read(id)
	if err != nil {
		return nil, err
	}
	d, err := ticket.ParseDoc(content)
	if err != nil {
		return nil, fmt.Errorf("refsync: the ref for %s does not carry a ticket: %w", id, err)
	}
	holder, _, _ := d.Scalar(ticket.FieldAssignee)
	if strings.TrimSpace(holder) != "" && !y.isMe(holder) && !force {
		out := &Released{ID: id}
		if w, wErr := y.winner(id); wErr == nil {
			out.NotYours = w
		}
		return out, fmt.Errorf("%w: %s has it", ErrTaken, holder)
	}
	if err := d.SetScalar(ticket.FieldAssignee, ""); err != nil {
		return nil, err
	}
	if err := ticket.Touch(d, time.Now()); err != nil {
		return nil, err
	}
	if err := ticket.TouchBy(d, y.Actor); err != nil {
		return nil, err
	}
	if _, err := y.Repo.Write(id, d.Bytes(), lease); err != nil {
		return nil, err
	}

	title, _, _ := d.Scalar(ticket.FieldTitle)
	out := &Released{ID: id, Title: title}
	if path, ok := y.localPath(id); ok {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return out, err
		}
		out.Removed = path
	}
	return out, nil
}

// Extra hands the store the tickets this board can see but has no file for —
// the ones still travelling on their refs, unpulled. It satisfies
// ticket.TicketSource.
//
// This is what keeps 'jaira list', 'jaira next', the board and the task mirror
// showing the whole board once an unworked ticket lives only on its ref. They
// all read through the store, so none of them has to know that refs exist.
func (y *Syncer) Extra(have map[string]bool) ([]*ticket.Ticket, error) {
	if y == nil || y.Usable() != nil {
		return nil, nil
	}
	ids, err := y.Repo.List()
	if err != nil {
		return nil, err
	}
	// Read from the filenames, never through Store.Load: Load falls back to
	// this very function for a ticket with no file, and asking it here would
	// recurse.
	local := y.localIDs()
	wanted := make([]string, 0, len(ids))
	for _, id := range ids {
		if have[id] {
			continue
		}
		if local[id] {
			// A file for it exists here after all: the file is the board's
			// copy, and the ref is not a second ticket.
			continue
		}
		wanted = append(wanted, id)
	}
	contents, err := y.Repo.ReadMany(wanted)
	if err != nil {
		return nil, err
	}
	out := make([]*ticket.Ticket, 0, len(contents))
	for _, id := range wanted {
		content, ok := contents[id]
		if !ok {
			continue
		}
		d, err := ticket.ParseDoc(content)
		if err != nil {
			continue // a ref whose tree is not a ticket is not this tool's
		}
		t, err := ticket.Decode(d, "")
		if err != nil {
			continue
		}
		out = append(out, t)
	}
	return out, nil
}

// Departure is a ticket whose ref is gone while the file is still here.
type Departure struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Path  string `json:"path"`
}

// Departed reports tickets somebody else has taken off the board while a file
// for them is still sitting here.
//
// The signal is a ref that this clone has seen before and cannot see now.
// Deleting the ref is exactly what logging or archiving a ticket does, so its
// absence is the board saying "this is finished elsewhere". Without noticing
// it, the next merge keeps both paths and the ticket lives twice — closed for
// its owner, open here.
//
// Nothing is moved. Which of the two is right is a question for the person, and
// a tool that guessed would sometimes reopen finished work; the caller reports
// this and names the command that settles it.
func (y *Syncer) Departed() []Departure {
	if y == nil || y.Store == nil || y.Usable() != nil {
		return nil
	}
	seen := y.loadSeen()
	if len(seen) == 0 {
		return nil
	}
	present := map[string]bool{}
	if ids, err := y.Repo.List(); err == nil {
		for _, id := range ids {
			present[id] = true
		}
	}
	var out []Departure
	for id := range seen {
		if present[id] {
			continue
		}
		path, ok := y.localPath(id)
		if !ok {
			continue // gone from the refs and not here either: nothing to say
		}
		d := Departure{ID: id, Path: path}
		if t, err := y.Store.Load(id); err == nil {
			d.Title = t.Title
		}
		out = append(out, d)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// ForgetDeparted drops a ticket from the seen record, so it stops being
// reported once the person has dealt with it.
func (y *Syncer) ForgetDeparted(id string) {
	if y == nil {
		return
	}
	seen := y.loadSeen()
	if _, ok := seen[id]; !ok {
		return
	}
	delete(seen, id)
	y.saveSeen(seen)
}

// Stranded is a ticket that was finished here and has not arrived anywhere
// everybody can see it.
type Stranded struct {
	ID    string        `json:"id"`
	Title string        `json:"title"`
	Path  string        `json:"path"`
	Since time.Duration `json:"-"`
	Days  int           `json:"days"`
}

// Stranded reports tickets filed away here whose ref is still standing because
// they have not landed in any landing branch.
//
// This is the cost of not deleting a ref at filing time, and it has to be
// visible: a branch that is never merged would otherwise keep its ref for
// ever, and the board would carry a finished ticket nobody can account for.
// Naming it is the whole fix — 'jaira snapshot --drop' is the way out, and it
// is deliberately a person's decision.
func (y *Syncer) Stranded(landing []string, grace time.Duration) []Stranded {
	if y == nil || y.Store == nil || y.Usable() != nil || len(landing) == 0 {
		return nil
	}
	if grace <= 0 {
		grace = 7 * 24 * time.Hour
	}
	var out []Stranded
	for id, path := range y.Store.FiledAwayIDs() {
		if _, err := y.Repo.SHA(id); err != nil {
			continue // no ref: it has already been cleared
		}
		if y.Repo.Landed(id, landing) != "" {
			continue // arrived; the next snapshot run clears the ref
		}
		fi, err := os.Stat(path)
		if err != nil {
			continue
		}
		age := time.Since(fi.ModTime())
		if age < grace {
			continue
		}
		st := Stranded{ID: id, Path: path, Since: age, Days: int(age.Hours() / 24)}
		if t, err := y.Store.Load(id); err == nil {
			st.Title = t.Title
		}
		out = append(out, st)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
