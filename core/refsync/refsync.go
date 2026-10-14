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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/BeMuCa/jaira/core/gitref"
	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/merge"
	"github.com/BeMuCa/jaira/core/milestone"
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

	// extraMu guards extraCache: the decoded ticket of each ref-only ticket,
	// with the commit hash it was decoded from. A hash is a content address, so
	// an equal hash needs no read and no decode, and no safety net behind it.
	extraMu    sync.Mutex
	extraCache map[string]extraEntry
}

type extraEntry struct {
	hash string
	t    *ticket.Ticket // nil: the ref held something that is not a ticket
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

// RecordMilestone queues a milestone file's current bytes for its ref, under
// the same rules as a ticket: the whole file, leased against the SHA the
// remote last confirmed, sent later.
//
// A milestone travels on a ref for the reason it exists at all. It is the one
// file everybody plans from, and the whole point of planning in a file is that
// twenty tickets are grouped in one edit — if that edit then waited for a
// branch to be merged, the file everybody is supposed to read would be the
// file nobody has.
func (y *Syncer) RecordMilestone(key, file string, content []byte) error {
	if y == nil || y.Usable() != nil {
		return nil
	}
	lease, err := y.Repo.MilestoneSHA(key)
	if err != nil && !errors.Is(err, gitref.ErrNoRef) {
		return err
	}
	if err := y.Box.QueueMilestone(key, file, outbox.OpWrite, content, lease, y.Actor); err != nil {
		return err
	}
	y.dirty = true
	return nil
}

// RecordMigration queues a milestone's move to its ID as the commit of its own
// the ref's base is derived from: the migrated content, and from, the old name
// ref's tip it came from.
func (y *Syncer) RecordMigration(key, file string, migrated []byte, from string) error {
	if y == nil || y.Usable() != nil {
		return nil
	}
	if from == "" {
		from = "none"
	}
	if err := y.Box.QueueMigration(key, file, migrated, from, y.Actor); err != nil {
		return err
	}
	y.dirty = true
	return nil
}

// MigrationBase is the base the ID ref's migration commit gives: the migrated
// content and the old ref's tip it came from. A ref without one was created
// with its ID, or moved there by a build of this change from before the
// migration commit existed; milestone.AdoptOld merges an old side of such a
// milestone against an empty base.
func (y *Syncer) MigrationBase(id string) (milestone.Base, bool) {
	if y == nil || y.Usable() != nil {
		return milestone.Base{}, false
	}
	content, from, err := y.Repo.ReadMigration(id)
	if err != nil {
		return milestone.Base{}, false
	}
	return milestone.Base{Content: string(content), OldSHA: from}, true
}

// IncomingMilestones writes to disk every milestone file the refs carry that
// this working tree does not already have in the same state, and reports the
// names it wrote.
//
// It writes rather than merely reporting because a milestone has no local
// editor in jaira the way a ticket has a lane: the file IS the interface, and
// a group somebody planned is of no use as a ref nobody has checked out. The
// ref is the newer state by construction — it is what the remote accepted —
// so a local file that differs is overwritten, and the previous content is in
// git if it was ever committed. The one exception is a milestone with a write
// still queued in the outbox: that file is newer than the ref and is left
// alone until the write has been sent or rejected.
//
// A milestone's ref is named after its ID and carries the file under its
// current name, so a rename arrives as the same ref with another file name.
// That is the one thing a fetch removes: every other file on this disk with
// the ref's ID is the milestone under its old name — or a copy git merged in
// beside it — and the ref decides which is right, before the file is compared.
// A file with another ID is another milestone, whatever it is called, and is
// never touched. The names removed this way come back as renamed, old name to
// new.
//
// A ref named after an old name — written before IDs existed — is read as the
// ID LegacyID gives that name. It is skipped once a milestone has moved on
// from it (the ID ref is here, or is on this disk); the mark a rename of that
// era left removes the old file this disk may still have and is reported as a
// rename; otherwise it is taken like an old file and migrated by the next
// write.
//
// A ref whose file says it has been filed is written over a file this tree
// already has: the marked line is what keeps a milestone off the board, so
// writing it is how a clone that already has the file learns it was filed at
// all — which is the one thing that ref exists to carry. Skipping it would
// leave that clone showing a round of work everybody else has closed. Those
// names come back separately, because "this group left the board" is
// different news from "this group changed".
//
// A filed milestone this tree does NOT have is not written: there is nothing
// here to correct, and putting the file on disk would leave a second copy
// beside the one in the logbook of the tree that filed it — which is the tree
// 'jaira restore' would then refuse, because the milestone is back on the
// board.
func (y *Syncer) IncomingMilestones(root string) (wrote, filed []string, renamed map[string]string, conflicts []milestone.Conflict, err error) {
	if y == nil || y.Usable() != nil {
		return nil, nil, nil, nil, nil
	}
	keys, err := y.Repo.ListMilestones()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	renamed = map[string]string{}
	// ID keys first: an old name ref's edit is adopted into the ID file in the
	// same fetch that writes it, whatever the name sorts like.
	sort.SliceStable(keys, func(i, j int) bool { return milestone.IsID(keys[i]) && !milestone.IsID(keys[j]) })
	for _, key := range keys {
		file, content, sha, err := y.Repo.ReadMilestone(key)
		if err != nil {
			continue
		}
		ref := milestone.FromBytes(file, content)
		if milestone.IsID(key) {
			if ref.ID != key {
				continue
			}
		} else if to := milestone.RenamedTo(content); to != "" {
			// The mark a rename of the name-keyed era left: the milestone is
			// under its new name now, so the old file goes if this disk still
			// has it, and nothing is written.
			if err := os.Remove(filepath.Join(milestone.Dir(root), file)); err == nil {
				renamed[ref.Name] = to
			}
			continue
		} else if y.HasMilestoneRef(ref.ID) || diskHasIDFile(root, ref.ID) {
			// The milestone has moved on to its ID, but an older build may
			// still write this ref: what it changed since the move is merged
			// into the ID milestone rather than ignored.
			adopted, cs, err := y.adoptOldRef(root, key, ref.ID, content, sha)
			if err != nil {
				return wrote, filed, renamed, conflicts, err
			}
			if adopted && !slices.Contains(wrote, ref.Name) {
				wrote = append(wrote, ref.Name)
			}
			conflicts = append(conflicts, cs...)
			continue
		}
		// A write of this milestone that has not been sent is newer than the
		// ref. It goes out on the next flush under the lease it was queued
		// with: if nobody wrote the ref since, it lands; if somebody did, the
		// flush is refused and mergeRefused merges the edit into theirs.
		if y.pendingMilestone(key, ref.Name, ref.ID) {
			continue
		}
		path := filepath.Join(milestone.Dir(root), file)
		others, hadFile := otherFilesWithID(root, ref.ID, file)
		// A file of the old format that lies beside the ID: the older build's
		// edits in it are merged into the ID file once the ref is on disk, and
		// only then is it removed — a merge that fails or a fetch that is
		// killed leaves it where the next fetch finds it. Another file with
		// this ID in its name carries nothing but a name, and goes as soon as
		// the ref is on disk.
		type oldSide struct {
			name    string
			content []byte
		}
		var oldSides []oldSide
		for _, other := range others {
			if !other.IDNamed {
				if b, err := os.ReadFile(filepath.Join(milestone.Dir(root), other.File)); err == nil {
					oldSides = append(oldSides, oldSide{strings.TrimSuffix(other.File, ".md"), b})
				}
			}
		}
		removeOthers := func(idNamed bool) error {
			for _, other := range others {
				if other.IDNamed != idNamed {
					continue
				}
				if err := os.Remove(filepath.Join(milestone.Dir(root), other.File)); err != nil && !os.IsNotExist(err) {
					return err
				}
				if other.Name != ref.Name {
					renamed[other.Name] = ref.Name
				}
			}
			return nil
		}
		adoptOld := func() error {
			for _, o := range oldSides {
				if milestone.StaleOld(y, o.name, o.content) {
					continue
				}
				adopted, cs, err := milestone.AdoptOld(root, y, ref.ID, ref.Name, o.content, "")
				if err != nil {
					return err
				}
				if adopted && !slices.Contains(wrote, ref.Name) {
					wrote = append(wrote, ref.Name)
				}
				conflicts = append(conflicts, cs...)
			}
			return removeOthers(false)
		}
		have, haveErr := os.ReadFile(path)
		if haveErr == nil && bytes.Equal(have, content) {
			if err := removeOthers(true); err != nil {
				return wrote, filed, renamed, conflicts, err
			}
			if err := adoptOld(); err != nil {
				return wrote, filed, renamed, conflicts, err
			}
			continue
		}
		isFiled := ref.Filed()
		if isFiled && !hadFile && len(others) == 0 {
			continue
		}
		if err := os.MkdirAll(milestone.Dir(root), 0o755); err != nil {
			return wrote, filed, renamed, conflicts, err
		}
		if err := ticket.WriteAtomic(path, content); err != nil {
			return wrote, filed, renamed, conflicts, err
		}
		if err := removeOthers(true); err != nil {
			return wrote, filed, renamed, conflicts, err
		}
		if err := adoptOld(); err != nil {
			return wrote, filed, renamed, conflicts, err
		}
		if isFiled {
			filed = append(filed, ref.Name)
			continue
		}
		if !slices.Contains(wrote, ref.Name) {
			wrote = append(wrote, ref.Name)
		}
	}
	return wrote, filed, renamed, conflicts, nil
}

// pendingMilestone reports whether a write of this milestone is queued, under
// its ref key, its ID or its old name.
func (y *Syncer) pendingMilestone(keys ...string) bool {
	for _, k := range keys {
		if _, ok := y.Box.PendingMilestone(k); ok {
			return true
		}
	}
	return false
}

// MilestonesFetched is what one fetch of milestones did on disk.
type MilestonesFetched struct {
	Wrote, Filed, Moved []string
	Renamed             map[string]string
	Conflicts           []milestone.Conflict
}

// FetchMilestones writes the milestone files the refs carry, then moves old
// files to their ID names (milestone.Migrate), under the milestone lock every
// other milestone writer takes. It is what 'jaira fetch' and the board's
// background fetch both run: the lock is only held here, never across the
// network round trip the caller made before.
//
// The error of Migrate comes back as it is, a *milestone.RecordError included:
// the caller decides how loud a ref that could not be queued is.
func (y *Syncer) FetchMilestones(s *ticket.Store) (MilestonesFetched, error) {
	unlock, err := s.Lock(milestone.LockName)
	if err != nil {
		return MilestonesFetched{}, err
	}
	defer unlock()
	var out MilestonesFetched
	out.Wrote, out.Filed, out.Renamed, out.Conflicts, err = y.IncomingMilestones(s.Root)
	if out.Renamed == nil {
		out.Renamed = map[string]string{}
	}
	if err != nil {
		return out, err
	}
	// Taken in after the refs, so a name-keyed ref that only now brought an old
	// file is moved to its ID in the same fetch.
	res, err := milestone.Migrate(s.Root, y)
	out.Moved = res.Moved
	out.Conflicts = append(out.Conflicts, res.Conflicts...)
	return out, err
}

// adoptOldRef merges what an older build wrote to an old name ref since the
// base recorded when its milestone moved to its ID into the ID file
// (milestone.AdoptOld). A ref that has not moved since the base is skipped
// without reading it again.
func (y *Syncer) adoptOldRef(root, name, id string, content []byte, sha string) (bool, []milestone.Conflict, error) {
	if base, ok := milestone.BaseFor(root, y, id); ok && sha != "" && sha == base.OldSHA {
		return false, nil, nil
	}
	return milestone.AdoptOld(root, y, id, name, content, sha)
}

// OldRefSHA is the tip of the ref an old milestone name had before IDs, "" when
// this tree has none. Migrate records it with the base.
func (y *Syncer) OldRefSHA(name string) string {
	if y == nil || y.Usable() != nil {
		return ""
	}
	sha, err := y.Repo.MilestoneSHA(name)
	if err != nil {
		return ""
	}
	return sha
}

// OldRefHad reports whether content is a state the ref of an old milestone
// name has carried. Such an old file is a stale copy, not an edit.
func (y *Syncer) OldRefHad(name string, content []byte) bool {
	if y == nil || y.Usable() != nil {
		return false
	}
	return y.Repo.MilestoneHad(name, content)
}

// FetchRefs brings the refs here and reports whether it did. An unreachable
// remote is not an error: whoever asks carries on without it.
func (y *Syncer) FetchRefs() bool {
	if y == nil || y.Usable() != nil {
		return false
	}
	return y.Repo.Fetch() == nil
}

// HasMilestoneRef reports whether this tree has the ref of the milestone with
// this ID. A syncer that is nil or without a remote has no refs.
func (y *Syncer) HasMilestoneRef(id string) bool {
	if y == nil || y.Usable() != nil {
		return false
	}
	_, err := y.Repo.MilestoneSHA(id)
	return err == nil
}

// diskHasIDFile reports whether a file named <name>-<ID>.md for this ID lies
// on the board: the milestone has moved to its ID on this disk already.
func diskHasIDFile(root, id string) bool {
	for _, m := range milestoneFiles(root) {
		if m.ID == id && m.IDNamed {
			return true
		}
	}
	return false
}

// otherFilesWithID lists the files on the board that carry this ID under a
// file name other than file, and whether file itself is there.
func otherFilesWithID(root, id, file string) (others []diskMilestone, hadFile bool) {
	for _, m := range milestoneFiles(root) {
		switch {
		case m.ID != id:
		case m.File == file:
			hadFile = true
		default:
			others = append(others, m)
		}
	}
	return others, hadFile
}

// diskMilestone is one file under .jaira/milestones read for what its name
// says: the milestone name and ID, and whether the ID is in the file name or
// was computed for an old one.
type diskMilestone struct {
	File, Name, ID string
	IDNamed        bool
}

func milestoneFiles(root string) []diskMilestone {
	entries, err := os.ReadDir(milestone.Dir(root))
	if err != nil {
		return nil
	}
	var out []diskMilestone
	for _, e := range entries {
		base, ok := strings.CutSuffix(e.Name(), ".md")
		if !ok || e.IsDir() {
			continue
		}
		m := diskMilestone{File: e.Name()}
		m.Name, m.ID, m.IDNamed = milestone.FileIdentity(base)
		out = append(out, m)
	}
	return out
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
	// Merged is set on a refused milestone write whose edit was merged into
	// the version that won, and queued again; Conflicts are the fields both
	// sides changed, where the winner's value was kept.
	Merged    bool                 `json:"merged,omitempty"`
	Conflicts []milestone.Conflict `json:"conflicts,omitempty"`
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

	if err := y.tidyMilestoneEntries(); err != nil {
		return nil, err
	}
	queued := map[string]outbox.Entry{}
	entries, _ = y.Box.List() // again: tidyMilestoneEntries may have rekeyed some
	for _, e := range entries {
		if e.Kind == outbox.KindMilestone {
			queued[e.ID] = e
		}
	}
	results, err := y.Box.Flush(y.Repo)
	if err != nil {
		return nil, err
	}
	reports := make([]Report, 0, len(results))
	for _, r := range results {
		rep := Report{Result: r}
		if r.Outcome == outbox.Rejected && r.Kind == outbox.KindMilestone {
			if e, ok := queued[r.ID]; ok {
				rep.Merged, rep.Conflicts = y.mergeRefused(e)
			}
		}
		// Only a ticket has a winner to name: the lookup parses the ref's
		// content as a ticket, and a milestone file is not one.
		if r.Outcome == outbox.Rejected && r.Kind != outbox.KindMilestone {
			if w, err := y.winner(r.ID); err == nil {
				rep.Winner = w
			}
		}
		reports = append(reports, rep)
	}
	return reports, nil
}

// mergeRefused keeps the edit of a milestone write the remote refused, because
// a teammate wrote the ref first. Without it the next fetch would replace the
// file with the teammate's version: the edit would be gone. It is merged into
// the version that won — the ref as this clone has just fetched it — from the
// state the write was leased on (for a moving write that has no lease, the base
// recorded when the file moved), and the result is written to the file and
// queued again, leased on the winner. A field both sides changed keeps the
// winner's value and the other is written into the file as a comment; the
// file name is merged like a field (see mergedFile). It reports whether it
// merged.
//
// It runs under the milestone lock every other milestone writer takes, and the
// refused entry stays queued until the merged one replaces it (Flush keeps it);
// a write another session queued since is the newer state and is what gets
// merged. On any failure the file is left as it is and the entry stays queued.
func (y *Syncer) mergeRefused(e outbox.Entry) (bool, []milestone.Conflict) {
	if y.Store == nil {
		_ = y.Box.DropMilestone(e.ID)
		return false, nil
	}
	unlock, err := y.Store.Lock(milestone.LockName)
	if err != nil {
		return false, nil
	}
	defer unlock()
	if now, ok := y.Box.PendingMilestone(e.ID); ok {
		e = now
	}
	file, cur, _, err := y.Repo.ReadMilestone(e.ID)
	if err != nil {
		return false, nil
	}
	var base []byte
	baseFile := ""
	if e.Lease != "" {
		if baseFile, base, err = y.Repo.ReadMilestoneFileAt(e.Lease); err != nil {
			base, baseFile = nil, ""
		}
	} else if b, ok := milestone.BaseFor(y.Store.Root, y, e.ID); ok {
		base = []byte(b.Content)
	}
	if e.Lease == "" && e.Migration != nil {
		// A moving write has no lease, but it carries the name the file had
		// when it moved: the name a later rename is measured against.
		baseFile = e.Migration.File
	}
	all, err := milestone.LoadAll(y.Store.Root)
	if err != nil {
		return false, nil
	}
	name := milestone.FromBytes(file, cur).Name
	merged, conflicts := milestone.Adopt(name, base, []byte(e.Content), cur, all)
	chosen, clash := mergedFile(baseFile, e.File, file)
	if clash != "" {
		c := milestone.Conflict{Milestone: name, Field: "name", Old: strconv.Quote(clash)}
		merged = milestone.NoteConflict(merged, c)
		conflicts = append(conflicts, c)
	}
	if bytes.Equal(merged, cur) && chosen == file {
		// Nothing of the edit is left to send: the winner already says it.
		_ = y.Box.DropMilestone(e.ID)
		return false, nil
	}
	dir := milestone.Dir(y.Store.Root)
	if err := ticket.WriteAtomic(filepath.Join(dir, chosen), merged); err != nil {
		return false, nil
	}
	for _, f := range []string{file, e.File} {
		if f != "" && f != chosen {
			_ = os.Remove(filepath.Join(dir, f))
		}
	}
	// The merged write replaces the refused one. The box keeps the lease of an
	// entry that is there, so the old one goes first, and comes back if the new
	// one cannot be queued.
	if err := y.Box.DropMilestone(e.ID); err != nil {
		return false, nil
	}
	if err := y.RecordMilestone(e.ID, chosen, merged); err != nil {
		_ = y.Box.QueueMilestone(e.ID, e.File, outbox.OpWrite, []byte(e.Content), e.Lease, e.By)
		return false, nil
	}
	return true, conflicts
}

// mergedFile is the file name a refused write is merged under: the file name
// is a field like the others. base is the name the write was leased on, local
// the name this clone wrote, winner the name on the ref now. A rename made on
// one side only is taken; one made on both sides differently keeps the winner's
// name, and the local name comes back as clash to be reported; when the base
// name is not known the winner's is kept. A moving write's base name is the
// file its migration commit carries.
func mergedFile(base, local, winner string) (file, clash string) {
	name := func(f string) string {
		n, _, _ := milestone.FileIdentity(strings.TrimSuffix(f, ".md"))
		return n
	}
	if local == "" || local == winner || base == "" {
		return winner, ""
	}
	switch {
	case winner == base:
		return local, ""
	case local == base:
		return winner, ""
	}
	return winner, name(local)
}

// tidyMilestoneEntries settles the milestone writes the queue may still hold
// from before IDs, on the refs this flush has just fetched. An entry keyed by a
// name is rekeyed to the ID of the file it names (FileIdentity: the ID in a
// <name>-<ID> key, else the one LegacyID computes) — as the file would be
// moved — unless it is the mark of a rename, which is dropped, or the ID
// already has an entry. And a first write of a milestone (no lease) that
// the fetched ref already carries, file and bytes, is dropped: two clones moved
// the same old file to the same ID, and the second one to push must not be told
// that the remote wins over a write that said the same thing.
func (y *Syncer) tidyMilestoneEntries() error {
	entries, err := y.Box.List()
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.Kind != outbox.KindMilestone {
			continue
		}
		if !milestone.IsID(e.ID) {
			if milestone.RenamedTo([]byte(e.Content)) != "" {
				if err := y.Box.DropMilestone(e.ID); err != nil {
					return err
				}
				continue
			}
			name, id, _ := milestone.FileIdentity(e.ID)
			if err := y.Box.RekeyMilestone(e.ID, id, name+"-"+id+".md", milestone.MigrateContent([]byte(e.Content)), ""); err != nil {
				return err
			}
			continue
		}
		if e.Lease != "" {
			continue
		}
		file, content, _, err := y.Repo.ReadMilestone(e.ID)
		if err == nil && file == e.File && string(content) == e.Content {
			if err := y.Box.DropMilestone(e.ID); err != nil {
				return err
			}
		}
	}
	return nil
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

// StaleError is the refusal for a ticket file that is older than its ref. It
// names both states, because "your file is stale" alone does not tell anyone
// what writing it would have thrown away.
type StaleError struct {
	ID        string
	File, Ref *ticket.Ticket
}

func (e *StaleError) Error() string {
	h := ticket.Handle(e.ID)
	return fmt.Sprintf("%s: the ticket file here is older than its ref; writing it would put that old state back on the ref\n"+
		"  file: %s\n  ref:  %s\n  run 'jaira pull %s' to bring the file up to the ref, then run this again",
		h, stateOf(e.File), stateOf(e.Ref), h)
}

func (e *StaleError) Unwrap() error { return ticket.ErrStale }

// stateOf is one side of a StaleError: the lane and when it was last written.
func stateOf(t *ticket.Ticket) string {
	return fmt.Sprintf("%s, updated %s", t.Status, t.UpdatedAt.UTC().Format("2006-01-02 15:04 UTC"))
}

// refTicket reads the ticket the local ref carries, with its bytes, or nil when
// there is no ref or it does not hold a ticket.
func (y *Syncer) refTicket(id string) (*ticket.Ticket, []byte) {
	content, _, err := y.Repo.Read(id)
	if err != nil {
		return nil, nil
	}
	d, err := ticket.ParseDoc(content)
	if err != nil {
		return nil, nil
	}
	t, err := ticket.Decode(d, "")
	if err != nil {
		return nil, nil
	}
	return t, content
}

// isStale is the one definition of "older than its ref": the ref was written
// strictly after the file. A file or ref without a timestamp gives no verdict.
func isStale(file, ref *ticket.Ticket) bool {
	return ref != nil && !file.UpdatedAt.IsZero() && !ref.UpdatedAt.IsZero() &&
		ref.UpdatedAt.After(file.UpdatedAt)
}

// CheckFresh refuses a write whose file is older than the ticket's ref. It
// satisfies ticket.FreshnessChecker, so every write through the store asks it.
//
// The local ref is the yardstick, not the remote: it is exactly the state a
// queued write leases against, so a file older than it is a file whose write
// the compare-and-swap would accept — and that is the whole problem. A file
// read from an older branch carries a perfectly current lease. Reading the
// local ref also keeps a write off the network.
//
// The comparison trusts clocks, and skew between machines can refuse a write
// that is in fact current: a teammate whose clock runs ahead writes the ref,
// this clone then writes offline with an earlier stamp, and since the local ref
// only moves after a push, the next write here reads as stale — and pull
// refuses too, over the unsent write — until the remote is reachable again.
func (y *Syncer) CheckFresh(t *ticket.Ticket) error {
	if y == nil || t == nil || y.Usable() != nil {
		return nil
	}
	if ref, _ := y.refTicket(t.ID); isStale(t, ref) {
		return &StaleError{ID: t.ID, File: t, Ref: ref}
	}
	return nil
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

	// Was is set when the file was already here but older than its ref, and
	// the pull brought it up: the state the file was in, beside Now, the state
	// it carries since. Nothing on the ref changes in that case, so nobody's
	// claim on the ticket does either.
	Was string `json:"was,omitempty"`
	Now string `json:"now,omitempty"`

	// TakenFrom names who the ticket was taken from, when it was taken with
	// --steal. A silent steal would leave the other person's board saying the
	// ticket is theirs, with nothing anywhere saying why it stopped being.
	TakenFrom string `json:"taken-from,omitempty"`
}

// ErrUnsent means the file here holds a write that has not reached the remote,
// so bringing it up to the ref would throw that write away.
var ErrUnsent = errors.New("refsync: the ticket has a write here that has not been sent")

// refresh is a pull of a ticket whose file is already here: it brings a file
// that is older than the ref up to the ref, and leaves any other file alone.
//
// It overwrites rather than merges. The file is behind by definition — the ref
// was written after it — so the ref already holds everything the file does,
// and a merge could only turn the gap into conflicts. The one thing the file
// can hold that the ref does not is an unsent write, and that is refused.
//
// The fetch is best effort: offline, the local ref is still the newest state
// this clone knows of, and a repeated pull has always worked without a network.
func (y *Syncer) refresh(id, path string) (*Pulled, error) {
	out := &Pulled{ID: id, Path: path, AlreadyHere: true}
	if err := y.Repo.Fetch(); err != nil && !errors.Is(err, gitref.ErrOffline) && !errors.Is(err, gitref.ErrNoGit) {
		return nil, err
	}
	file, err := y.Store.Load(id)
	if err != nil {
		return nil, err
	}
	out.Title = file.Title
	ref, content := y.refTicket(id)
	if !isStale(file, ref) {
		return out, nil
	}
	if _, pending := y.Pending(id); pending {
		return out, fmt.Errorf("%w: nothing was overwritten, and the ref has moved on since (%s here, %s on the ref)",
			ErrUnsent, stateOf(file), stateOf(ref))
	}
	dst := filepath.Join(filepath.Dir(path), ticket.Filename(id, ref.Title))
	if err := ticket.WriteAtomic(dst, content); err != nil {
		return nil, err
	}
	// A title changed on the ref changes the filename; the old one would be a
	// second file claiming the same id.
	if dst != path {
		if err := os.Remove(path); err != nil {
			return nil, err
		}
	}
	out.AlreadyHere = false
	out.Path, out.Title = dst, ref.Title
	out.Was, out.Now = stateOf(file), stateOf(ref)
	return out, nil
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
	// ticket to the lane it is already in — unless the file is behind its ref.
	// The question is whether there is a FILE, not whether the board can see
	// the ticket — it can see every ref.
	if path, ok := y.localPath(id); ok {
		return y.refresh(id, path)
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
	hashes, err := y.Repo.ListHashes()
	if err != nil {
		return nil, err
	}
	// Read from the filenames, never through Store.Load: Load falls back to
	// this very function for a ticket with no file, and asking it here would
	// recurse.
	local := y.localIDs()

	// A ticket filed away here is off this board on purpose. Its ref outlives
	// the filing — it goes when the ticket lands in a landing branch — so
	// without this the logbook would hand every finished ticket straight back
	// as a card, which is the opposite of what filing it meant.
	filed := map[string]bool{}
	if y.Store != nil {
		for id := range y.Store.FiledAwayIDs() {
			filed[id] = true
		}
	}
	y.extraMu.Lock()
	defer y.extraMu.Unlock()
	if y.extraCache == nil {
		y.extraCache = map[string]extraEntry{}
	}
	for id := range y.extraCache {
		if _, listed := hashes[id]; !listed {
			delete(y.extraCache, id)
		}
	}
	var kept, wanted []string
	// By id: the order List gave, which callers have always seen.
	for _, id := range slices.Sorted(maps.Keys(hashes)) {
		if have[id] {
			continue
		}
		if local[id] || filed[id] {
			// Either a file for it exists here — the file is the board's copy,
			// and the ref is not a second ticket — or it has been filed away,
			// which means it left this board deliberately.
			continue
		}
		kept = append(kept, id)
		if e, ok := y.extraCache[id]; ok && e.hash == hashes[id] {
			continue
		}
		wanted = append(wanted, id)
	}
	contents, err := y.Repo.ReadMany(wanted)
	if err != nil {
		return nil, err
	}
	for _, id := range wanted {
		content, ok := contents[id]
		if !ok {
			delete(y.extraCache, id)
			continue
		}
		var t *ticket.Ticket
		if d, err := ticket.ParseDoc(content); err == nil { // a ref whose tree is not a ticket is not this tool's
			if tk, err := ticket.Decode(d, ""); err == nil {
				t = tk
			}
		}
		y.extraCache[id] = extraEntry{hash: hashes[id], t: t}
	}
	var out []*ticket.Ticket
	for _, id := range kept {
		if e, ok := y.extraCache[id]; ok && e.t != nil {
			out = append(out, e.t.Clone())
		}
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
