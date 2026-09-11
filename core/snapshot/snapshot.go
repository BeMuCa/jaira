// Package snapshot keeps a branch that holds the board as files, and reaps the
// refs whose tickets have arrived in a branch that matters.
//
// It exists because a ticket nobody is working lives only on its ref, which
// means an untouched backlog lives only on the remote. Every participant's
// clone holds the refs it has fetched, so the board survives any one machine —
// but somebody cloning for the first time has no refs at all, and neither has
// somebody who never fetched. This branch is for them. It is a backup, not the
// storage: the working state is always on the refs.
//
// The branch is built with plumbing and never checked out. That is not a trick
// to save time: the run happens in the background of somebody else's command,
// and a run that touched the working tree could not be allowed to happen while
// they were in the middle of anything.
//
// Two properties are load-bearing:
//
//   - The files live under board/, not under .jaira/tickets/. If they sat at
//     the same path, the first person to merge this branch by accident would
//     get an add/add conflict across every ticket at once — the worse version
//     of the problem this design exists to avoid. A different path makes the
//     collision impossible by construction rather than by care.
//   - The reaping happens here, immediately after the write. At that moment
//     the ticket is in two places: in this snapshot and in the branch it
//     landed in. That is the difference between tidying up and throwing away,
//     and it is why removing a ref belongs in this one place rather than in
//     logbook, in archive and in a timer.
package snapshot

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/BeMuCa/jaira/core/gitref"
)

// DefaultBranch is where the snapshot goes unless told otherwise.
const DefaultBranch = "jaira/board"

// Dir is the directory the ticket files live in inside the snapshot branch.
const Dir = "board"

// DefaultEvery is how old a snapshot may get before another is taken.
//
// Three days rather than minutes, and the reason is not thrift: every
// participant's clone holds the refs locally, so the board is reconstructible
// from any of them. This is the backup for the first-time cloner, which is a
// slow case. Each ticket ref also carries its own history, so the snapshot is
// not the record of what happened either — a commit per write would buy
// nothing and cost hundreds of commits a day.
const DefaultEvery = 72 * time.Hour

// Result reports one run.
type Result struct {
	// Branch is the snapshot branch.
	Branch string `json:"branch"`
	// Commit is the snapshot that was written, or "" when nothing changed.
	Commit string `json:"commit,omitempty"`
	// Tickets is how many tickets the snapshot holds.
	Tickets int `json:"tickets"`
	// Added and Removed are the ids this run put in and took out, against the
	// previous snapshot.
	Added   []string `json:"added,omitempty"`
	Removed []string `json:"removed,omitempty"`
	// Reaped are the refs deleted because their ticket has landed.
	Reaped []string `json:"reaped,omitempty"`
	// Unchanged says the set of refs was identical to the last snapshot, so no
	// commit was written and nothing was pushed.
	Unchanged bool `json:"unchanged"`
}

// Runner takes snapshots for one board.
type Runner struct {
	Repo   *gitref.Repo
	Branch string

	// LandingBranches are the branches whose contents mean "this ticket has
	// arrived": a ticket filed away in one of them is finished for everybody,
	// so its ref has done its job. Empty means nothing is reaped — an
	// unresolved answer must not delete anything, because a ref left standing
	// costs nothing while one removed by mistake takes away exactly the
	// visibility this is for.
	LandingBranches []string
}

func (r *Runner) branch() string {
	if b := strings.TrimSpace(r.Branch); b != "" {
		return b
	}
	return DefaultBranch
}

// Run writes a snapshot if the refs have changed since the last one, then reaps
// the refs whose tickets have landed.
func (r *Runner) Run() (*Result, error) {
	if r == nil || r.Repo == nil {
		return nil, errors.New("snapshot: no repository")
	}
	if err := r.Repo.Usable(); err != nil {
		return nil, err
	}
	branch := r.branch()
	res := &Result{Branch: branch}

	// A landing branch nobody could resolve means nothing gets reaped, and the
	// commonest reason is a clone that simply never learned the remote's own
	// default (a clone of an empty repository cannot). This is the one place
	// allowed to ask, because it is the background run: it already touches the
	// network and nobody is waiting for it.
	if len(r.LandingBranches) == 0 && r.Repo.RemoteHead() == "" {
		if err := r.Repo.SetRemoteHead(); err == nil {
			if head := r.Repo.RemoteHead(); head != "" {
				r.LandingBranches = []string{head}
			}
		}
	}

	// Start from what the remote has, so two people snapshotting produce one
	// history rather than two that overwrite each other.
	if err := r.Repo.FetchBranch(branch); err != nil {
		if !errors.Is(err, gitref.ErrOffline) && !isNoSuchBranch(err) {
			return nil, err
		}
	}
	if err := r.Repo.Fetch(); err != nil && !errors.Is(err, gitref.ErrOffline) {
		return nil, err
	}

	ids, err := r.Repo.List()
	if err != nil {
		return nil, err
	}
	sort.Strings(ids)
	blobs, err := r.Repo.BlobsOf(ids)
	if err != nil {
		return nil, err
	}
	entries := make(map[string]string, len(blobs))
	for id, blob := range blobs {
		entries[id+".md"] = blob
	}
	res.Tickets = len(entries)

	local := "refs/heads/" + branch
	lease := r.Repo.Rev(local)
	previous := r.previousIDs(lease)

	var tree string
	if len(entries) > 0 {
		inner, err := r.Repo.Tree(entries)
		if err != nil {
			return nil, err
		}
		tree, err = r.Repo.Nest(Dir, inner)
		if err != nil {
			return nil, err
		}
	} else {
		// No refs at all: an empty tree, so a board that was emptied is
		// recorded as empty rather than frozen at its last populated state.
		tree, err = r.Repo.Tree(nil)
		if err != nil {
			return nil, err
		}
	}

	res.Added, res.Removed = diffIDs(previous, ids)

	if tree == r.Repo.TreeOf(local) {
		// Nothing changed. No commit, no push: a snapshot that says the same
		// thing as the last one is not a record, it is noise, and it would
		// cost a network round trip to add it.
		res.Unchanged = true
		res.Reaped = r.reap(ids)
		return res, nil
	}

	commit, err := r.Repo.Commit(tree, lease, r.message(res))
	if err != nil {
		return nil, err
	}
	if err := r.Repo.PushBranch(branch, commit, lease); err != nil {
		return nil, err
	}
	res.Commit = commit

	// Reaped only after the snapshot is on the remote. At this point every
	// ticket is in two places, which is what makes deleting a ref safe.
	res.Reaped = r.reap(ids)
	return res, nil
}

// message names what the run did, so 'git log' on the branch reads as the
// board's history rather than as a column of identical lines.
func (r *Runner) message(res *Result) string {
	var b strings.Builder
	fmt.Fprintf(&b, "board: %d ticket(s)", res.Tickets)
	if len(res.Added) > 0 {
		fmt.Fprintf(&b, ", +%d", len(res.Added))
	}
	if len(res.Removed) > 0 {
		fmt.Fprintf(&b, ", -%d", len(res.Removed))
	}
	b.WriteString("\n")
	for _, id := range res.Added {
		fmt.Fprintf(&b, "\nadded   %s", id)
	}
	for _, id := range res.Removed {
		fmt.Fprintf(&b, "\nremoved %s", id)
	}
	return b.String()
}

// previousIDs reads the ids the last snapshot held.
func (r *Runner) previousIDs(rev string) []string {
	if rev == "" {
		return nil
	}
	names, err := r.Repo.LsTree(rev, Dir)
	if err != nil {
		return nil
	}
	out := make([]string, 0, len(names))
	for _, n := range names {
		out = append(out, strings.TrimSuffix(n, ".md"))
	}
	sort.Strings(out)
	return out
}

// reap deletes the refs whose tickets have arrived in a landing branch.
//
// A ticket only counts as landed when it is filed away there — in the logbook
// or the archive. One still sitting under tickets/ on that branch is on the
// board, which is the opposite of finished.
func (r *Runner) reap(ids []string) []string {
	if len(r.LandingBranches) == 0 {
		return nil
	}
	var reaped []string
	for _, id := range ids {
		if r.Repo.Landed(id, r.LandingBranches) == "" {
			continue
		}
		lease, err := r.Repo.SHA(id)
		if err != nil {
			continue
		}
		if err := r.Repo.Delete(id, lease); err != nil {
			continue // a race or no route: the next run tries again
		}
		reaped = append(reaped, id)
	}
	return reaped
}

// Drop removes one ref without waiting for the ticket to land, for a branch
// that is never going to be merged.
//
// Deliberately a separate act with its own command: the automatic path only
// ever removes a ref whose ticket is provably somewhere else, and a tool that
// quietly relaxed that rule would recreate the invisible window this whole
// design closes.
func (r *Runner) Drop(id string) error {
	if r == nil || r.Repo == nil {
		return errors.New("snapshot: no repository")
	}
	lease, err := r.Repo.SHA(id)
	if err != nil {
		return err
	}
	return r.Repo.Delete(id, lease)
}

func diffIDs(before, after []string) (added, removed []string) {
	was := map[string]bool{}
	for _, id := range before {
		was[id] = true
	}
	is := map[string]bool{}
	for _, id := range after {
		is[id] = true
		if !was[id] {
			added = append(added, id)
		}
	}
	for _, id := range before {
		if !is[id] {
			removed = append(removed, id)
		}
	}
	sort.Strings(added)
	sort.Strings(removed)
	return added, removed
}

func isNoSuchBranch(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "couldn't find remote ref") ||
		strings.Contains(msg, "no such ref") ||
		strings.Contains(msg, "not our ref")
}
