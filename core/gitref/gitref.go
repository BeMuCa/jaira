// Package gitref carries a ticket on a git ref of its own, outside
// refs/heads.
//
// The problem it solves: a ticket file lives in the branch of whoever wrote it.
// If that branch is unpushed, stale or force-pushed, the person the ticket was
// assigned to sees nothing — not the title, not even the id. Scanning every
// branch for ticket files is expensive and still wrong.
//
// A ref under refs/jaira/tickets/<id> belongs to no branch. The default
// refspec does not fetch it, so a clone by someone who does not use jaira is
// unaffected, and the push onto it is a compare-and-swap that uses the ref lock
// the hosting provider already implements. No server, no daemon: the remote
// that is already there is the channel.
//
// The ref's tree carries the ticket file itself rather than only its id, so a
// colleague reads the whole ticket without a branch and without a checkout:
//
//	git show refs/jaira/tickets/<id>:<id>.md
//
// The ticket file also stays in the code commit, and the two answer different
// questions: the file in the commit says what this change was and why, the ref
// says where the ticket stands right now and with whom. Same bytes, one of them
// on a channel that is visible without the branch.
package gitref

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

// Prefix is the ref namespace. It is outside refs/heads deliberately: the
// default refspec does not fetch it, so these refs are invisible to anyone who
// does not ask for them, and they never show up as branches or in a hosting
// provider's web UI.
const Prefix = "refs/jaira/tickets/"

// DefaultRemote is where refs go when nothing is configured.
const DefaultRemote = "origin"

var (
	// ErrRaceLost means someone else wrote the ref after we read it. The
	// caller has stale bytes and must re-read; it is not a transport failure
	// and retrying the same push would only lose again.
	ErrRaceLost = errors.New("gitref: the ticket moved since it was read")

	// ErrOffline means the remote could not be reached at all. The write is
	// not lost and not refused — it is unsent, which is what the outbox is
	// for.
	ErrOffline = errors.New("gitref: the remote is unreachable")

	// ErrNoRef means the ticket has no ref yet.
	ErrNoRef = errors.New("gitref: no ref for this ticket")

	// ErrNoRepo means this board is not in a git repository, or the remote it
	// was told to use does not exist. It is not a failure to report at every
	// write: a board without a repository is a supported way to use jaira, it
	// simply does not carry tickets on refs.
	ErrNoRepo = errors.New("gitref: no repository or no such remote")

	// ErrNoGit means the git binary is unavailable.
	ErrNoGit = errors.New("gitref: git is not available on PATH")
)

// Repo is a working tree together with the remote its ticket refs live on.
type Repo struct {
	Dir    string
	Remote string

	// Author names the identity written onto the ref commit. It is separate
	// from git config on purpose: these commits are made in the background of
	// somebody else's command, and a machine with no user.name configured (CI,
	// a fresh container) must still be able to write a ticket rather than fail
	// with git's "please tell me who you are".
	AuthorName  string
	AuthorEmail string
}

// Usable reports whether this working tree can carry tickets on refs at all:
// git present, inside a repository, and the remote actually configured.
//
// Asked once by the caller rather than inferred from a failed write, because
// the two are different answers to different questions — "this board does not
// use refs" is a permanent property of the checkout, while a refused push is
// about one ticket at one moment.
func (r *Repo) Usable() error {
	out, _, err := r.run("", "rev-parse", "--is-inside-work-tree")
	if err != nil {
		if errors.Is(err, ErrNoGit) {
			return err
		}
		return ErrNoRepo
	}
	if strings.TrimSpace(out) != "true" {
		return ErrNoRepo
	}
	if _, _, err := r.run("", "remote", "get-url", r.remote()); err != nil {
		if errors.Is(err, ErrNoGit) {
			return err
		}
		return fmt.Errorf("%w: no remote %q", ErrNoRepo, r.remote())
	}
	return nil
}

// RefName is the ref a ticket travels on.
func RefName(id string) string { return Prefix + id }

func (r *Repo) remote() string {
	if strings.TrimSpace(r.Remote) == "" {
		return DefaultRemote
	}
	return r.Remote
}

func (r *Repo) author() (string, string) {
	name, email := r.AuthorName, r.AuthorEmail
	if strings.TrimSpace(name) == "" {
		name = "jaira"
	}
	if strings.TrimSpace(email) == "" {
		email = "jaira@localhost"
	}
	return name, email
}

// run executes git and keeps stdout and stderr apart, because the classification
// of a failed push lives in stderr and must not be mixed into the value a
// successful call returns.
func (r *Repo) run(stdin string, args ...string) (string, string, error) {
	if _, err := exec.LookPath("git"); err != nil {
		return "", "", ErrNoGit
	}
	cmd := exec.Command("git", append([]string{"-C", r.Dir}, args...)...)
	name, email := r.author()
	cmd.Env = append(cmd.Environ(),
		"GIT_AUTHOR_NAME="+name, "GIT_AUTHOR_EMAIL="+email,
		"GIT_COMMITTER_NAME="+name, "GIT_COMMITTER_EMAIL="+email,
		// A push that waits for a credential prompt would hang a command the
		// user ran for another reason entirely. Failing is the correct answer:
		// an unsent write goes to the outbox.
		"GIT_TERMINAL_PROMPT=0",
	)
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	err := cmd.Run()
	return out.String(), errb.String(), err
}

// value runs git for its stdout and reports a failure as an error naming the
// command, matching gitrepo's shape for everything that is not a push.
func (r *Repo) value(args ...string) (string, error) {
	out, errb, err := r.run("", args...)
	if err != nil {
		if errors.Is(err, ErrNoGit) {
			return "", err
		}
		msg := strings.TrimSpace(errb)
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), msg)
	}
	return strings.TrimSpace(out), nil
}

// SHA returns the commit the ticket's local ref points at, or ErrNoRef. This
// SHA is the lease: every write is a compare-and-swap against the value the
// writer last read.
func (r *Repo) SHA(id string) (string, error) {
	out, _, err := r.run("", "rev-parse", "--verify", "--quiet", RefName(id))
	if err != nil {
		if errors.Is(err, ErrNoGit) {
			return "", err
		}
		return "", ErrNoRef
	}
	if sha := strings.TrimSpace(out); sha != "" {
		return sha, nil
	}
	return "", ErrNoRef
}

// Read returns the ticket file carried by the local ref.
func (r *Repo) Read(id string) ([]byte, string, error) {
	sha, err := r.SHA(id)
	if err != nil {
		return nil, "", err
	}
	out, errb, runErr := r.run("", "show", sha+":"+id+".md")
	if runErr != nil {
		if errors.Is(runErr, ErrNoGit) {
			return nil, "", runErr
		}
		return nil, "", fmt.Errorf("gitref: ref %s carries no %s.md: %s", RefName(id), id, strings.TrimSpace(errb))
	}
	return []byte(out), sha, nil
}

// ReadParent returns the ticket file as it was one commit earlier on the ref —
// the state this clone had already seen when it last wrote or fetched.
//
// This is what makes reconciling the local file with the ref a real three-way
// merge rather than a guess about which side is newer: because every ref commit
// takes the leased SHA as its parent, the parent's blob is exactly "the state
// both sides started from". A ref with no parent (the ticket's first write) has
// no such state, and reports ErrNoRef so the caller can fall back to two-way.
func (r *Repo) ReadParent(id string) ([]byte, error) {
	sha, err := r.SHA(id)
	if err != nil {
		return nil, err
	}
	parent, _, err := r.run("", "rev-parse", "--verify", "--quiet", sha+"^")
	if err != nil || strings.TrimSpace(parent) == "" {
		return nil, ErrNoRef
	}
	out, _, runErr := r.run("", "show", strings.TrimSpace(parent)+":"+id+".md")
	if runErr != nil {
		return nil, ErrNoRef
	}
	return []byte(out), nil
}

// Write puts the ticket file on the ref and pushes it, refusing if the remote
// moved since lease was read.
//
// The commit takes the leased SHA as its parent when there is one, so the ref
// keeps the ticket's history (git log on the ref is the ticket's movement) and
// a stale write is a non-fast-forward that git rejects on its own.
// --force-with-lease is still passed: it is the check that holds when a ticket
// is written for the first time by two people at once, where neither commit has
// a parent to compare.
//
// An empty lease means "no ref yet" and requires the remote to have none
// either.
func (r *Repo) Write(id string, content []byte, lease string) (string, error) {
	blob, err := r.hashObject(content)
	if err != nil {
		return "", err
	}
	tree, err := r.mktree(id+".md", blob)
	if err != nil {
		return "", err
	}
	args := []string{"commit-tree", tree, "-m", "jaira: " + id}
	if lease != "" {
		args = append(args, "-p", lease)
	}
	commit, err := r.value(args...)
	if err != nil {
		return "", err
	}
	if err := r.push(id, commit, lease); err != nil {
		return "", err
	}
	// Only now is the local ref moved: it records what the remote accepted, so
	// the next lease is a SHA the remote has really seen.
	if _, err := r.value("update-ref", RefName(id), commit); err != nil {
		return commit, err
	}
	return commit, nil
}

// Delete removes the ticket's ref, locally and on the remote. A logged or
// archived ticket is off the board, and leaving its ref behind would keep it on
// everybody else's.
func (r *Repo) Delete(id, lease string) error {
	if err := r.pushRefspec(id, ":"+RefName(id), lease); err != nil {
		return err
	}
	_, _, _ = r.run("", "update-ref", "-d", RefName(id))
	return nil
}

// ReadMany returns the ticket files carried by many refs, in two git
// invocations rather than two per ticket.
//
// Measured on this project: reading a hundred tickets one at a time costs
// about a quarter of a second, almost all of it process startup, while one
// cat-file --batch pass is not measurable. The board reads every ticket on
// every refresh, so the per-ticket shape would grow into the one thing this
// tool must never be — slow to open.
func (r *Repo) ReadMany(ids []string) (map[string][]byte, error) {
	out := make(map[string][]byte, len(ids))
	if len(ids) == 0 {
		return out, nil
	}
	var in strings.Builder
	for _, id := range ids {
		fmt.Fprintf(&in, "%s:%s.md\n", RefName(id), id)
	}
	stdout, _, err := r.run(in.String(), "cat-file", "--batch")
	if err != nil {
		if errors.Is(err, ErrNoGit) {
			return nil, err
		}
		return nil, fmt.Errorf("gitref: cat-file --batch: %w", err)
	}
	// Each answer is "<sha> blob <size>\n<size bytes>\n", or a line ending in
	// "missing" for a ref that has gone. Parsing by the declared size rather
	// than by scanning for the next header is what makes a ticket containing
	// the word "blob" harmless.
	rest := stdout
	for _, id := range ids {
		nl := strings.IndexByte(rest, '\n')
		if nl < 0 {
			break
		}
		header := rest[:nl]
		rest = rest[nl+1:]
		fields := strings.Fields(header)
		if len(fields) < 3 {
			continue // missing, ambiguous: not an error, just nothing to show
		}
		var size int
		if _, err := fmt.Sscanf(fields[2], "%d", &size); err != nil || size > len(rest) {
			break
		}
		out[id] = []byte(rest[:size])
		rest = rest[size:]
		rest = strings.TrimPrefix(rest, "\n")
	}
	return out, nil
}

// BlobsOf resolves the blob each ticket ref carries, in one git invocation.
//
// A snapshot reuses these blobs rather than hashing the content again: the
// objects are already in the repository, put there by the ref writes, so a
// snapshot adds a tree and a commit and nothing else.
func (r *Repo) BlobsOf(ids []string) (map[string]string, error) {
	out := make(map[string]string, len(ids))
	if len(ids) == 0 {
		return out, nil
	}
	var in strings.Builder
	for _, id := range ids {
		fmt.Fprintf(&in, "%s:%s.md\n", RefName(id), id)
	}
	stdout, _, err := r.run(in.String(), "cat-file", "--batch-check")
	if err != nil {
		if errors.Is(err, ErrNoGit) {
			return nil, err
		}
		return nil, fmt.Errorf("gitref: cat-file --batch-check: %w", err)
	}
	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	for i, id := range ids {
		if i >= len(lines) {
			break
		}
		fields := strings.Fields(lines[i])
		if len(fields) < 2 || fields[1] != "blob" {
			continue // missing or not a file: nothing to put in the snapshot
		}
		out[id] = fields[0]
	}
	return out, nil
}

// Tree writes a tree object holding the given entries, and returns its hash.
// The entries are name to blob, all regular files.
func (r *Repo) Tree(entries map[string]string) (string, error) {
	names := make([]string, 0, len(entries))
	for name := range entries {
		names = append(names, name)
	}
	sort.Strings(names)
	var in strings.Builder
	for _, name := range names {
		fmt.Fprintf(&in, "100644 blob %s\t%s\n", entries[name], name)
	}
	out, errb, err := r.run(in.String(), "mktree")
	if err != nil {
		if errors.Is(err, ErrNoGit) {
			return "", err
		}
		return "", fmt.Errorf("gitref: mktree: %s", strings.TrimSpace(errb))
	}
	return strings.TrimSpace(out), nil
}

// Nest wraps a tree under a single directory name and returns the outer tree.
func (r *Repo) Nest(dir, tree string) (string, error) {
	out, errb, err := r.run(fmt.Sprintf("040000 tree %s\t%s\n", tree, dir), "mktree")
	if err != nil {
		if errors.Is(err, ErrNoGit) {
			return "", err
		}
		return "", fmt.Errorf("gitref: mktree: %s", strings.TrimSpace(errb))
	}
	return strings.TrimSpace(out), nil
}

// Commit writes a commit object for a tree, optionally with a parent.
func (r *Repo) Commit(tree, parent, message string) (string, error) {
	args := []string{"commit-tree", tree, "-m", message}
	if parent != "" {
		args = append(args, "-p", parent)
	}
	return r.value(args...)
}

// TreeOf returns the tree a commit points at, or "" when the commit is unknown.
func (r *Repo) TreeOf(rev string) string {
	out, _, err := r.run("", "rev-parse", "--verify", "--quiet", rev+"^{tree}")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out)
}

// Rev resolves a revision to a SHA, or "" when it does not exist.
func (r *Repo) Rev(rev string) string {
	out, _, err := r.run("", "rev-parse", "--verify", "--quiet", rev)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out)
}

// PushBranch moves a branch on the remote to commit, refusing if the remote
// moved since lease was read. Used for the snapshot branch, where two people
// snapshotting at once must not overwrite one another.
func (r *Repo) PushBranch(branch, commit, lease string) error {
	local := "refs/heads/" + branch
	args := []string{"push", "--force-with-lease=" + local + ":" + lease, r.remote(), commit + ":" + local}
	_, errb, err := r.run("", args...)
	if err != nil {
		if errors.Is(err, ErrNoGit) {
			return err
		}
		return classify(errb, err)
	}
	_, _, _ = r.run("", "update-ref", local, commit)
	return nil
}

// FetchBranch brings one branch from the remote into a local ref of the same
// name, so a snapshot can be appended to what the remote already has.
func (r *Repo) FetchBranch(branch string) error {
	_, errb, err := r.run("", "fetch", "--quiet", r.remote(),
		"+refs/heads/"+branch+":refs/heads/"+branch)
	if err != nil {
		if errors.Is(err, ErrNoGit) {
			return err
		}
		return classify(errb, err)
	}
	return nil
}

// RemoteHead returns the remote's default branch as a revision (origin/master
// and the like), or "" if the clone does not know it.
//
// A clone of a repository that has commits knows this without asking anybody;
// a clone of an empty one does not, and cannot be told until the remote has a
// branch at all. So this answers from what is already local, and SetRemoteHead
// is the separate, network-touching way to fill it in later.
func (r *Repo) RemoteHead() string {
	out, _, err := r.run("", "symbolic-ref", "--quiet", "refs/remotes/"+r.remote()+"/HEAD")
	if err != nil {
		return ""
	}
	return strings.TrimPrefix(strings.TrimSpace(out), "refs/remotes/")
}

// SetRemoteHead asks the remote which branch it points HEAD at and records it,
// so RemoteHead can answer locally from then on. One round trip, ever.
func (r *Repo) SetRemoteHead() error {
	_, errb, err := r.run("", "remote", "set-head", r.remote(), "-a")
	if err != nil {
		if errors.Is(err, ErrNoGit) {
			return err
		}
		return classify(errb, err)
	}
	return nil
}

// LsTree lists the file names directly under dir in a revision.
func (r *Repo) LsTree(rev, dir string) ([]string, error) {
	out, err := r.value("ls-tree", "--name-only", rev+":"+dir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, line := range strings.Split(out, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			names = append(names, line)
		}
	}
	return names, nil
}

// Landed reports the commit in which a ticket was filed away on one of the
// given branches, or "" if it has not arrived there.
//
// Only the filed-away paths count. A ticket still sitting under tickets/ on
// that branch has not landed in this sense: it is on the board there, which is
// the opposite of finished.
func (r *Repo) Landed(id string, branches []string) string {
	for _, b := range branches {
		rev := r.Rev(b)
		if rev == "" {
			continue
		}
		out, _, err := r.run("", "rev-list", "-1", rev, "--",
			".jaira/logbook/*/"+id+"*", ".jaira/archive/"+id+"*")
		if err != nil {
			continue
		}
		if sha := strings.TrimSpace(out); sha != "" {
			return sha
		}
	}
	return ""
}

// Fetch brings every ticket ref up to date in one roundtrip. The refspec is
// explicit because the default one does not carry these refs — which is the
// property that keeps them invisible to people who do not use jaira.
//
// --prune is part of the contract, not a tidy-up: a ticket that was logged or
// archived deletes its ref, and without pruning that ticket would sit on
// everyone else's board forever, deleted for its owner and immortal for
// everyone else.
func (r *Repo) Fetch() error {
	_, errb, err := r.run("", "fetch", "--quiet", "--prune", r.remote(), "+"+Prefix+"*:"+Prefix+"*")
	if err != nil {
		if errors.Is(err, ErrNoGit) {
			return err
		}
		return classify(errb, err)
	}
	return nil
}

// List returns the ticket ids that have a local ref.
func (r *Repo) List() ([]string, error) {
	out, err := r.value("for-each-ref", "--format=%(refname)", Prefix)
	if err != nil {
		return nil, err
	}
	return idsFrom(out, func(line string) string { return line }), nil
}

// ListRemote asks the remote which tickets exist, without fetching them. One
// roundtrip is cheap enough to run before a command needs the contents.
func (r *Repo) ListRemote() ([]string, error) {
	out, errb, err := r.run("", "ls-remote", r.remote(), Prefix+"*")
	if err != nil {
		if errors.Is(err, ErrNoGit) {
			return nil, err
		}
		return nil, classify(errb, err)
	}
	return idsFrom(out, func(line string) string {
		if i := strings.IndexAny(line, " \t"); i >= 0 {
			return strings.TrimSpace(line[i:])
		}
		return line
	}), nil
}

func idsFrom(out string, refOf func(string) string) []string {
	var ids []string
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		ref := refOf(line)
		if id := strings.TrimPrefix(ref, Prefix); id != ref && id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}

func (r *Repo) hashObject(content []byte) (string, error) {
	out, errb, err := r.run(string(content), "hash-object", "-w", "--stdin")
	if err != nil {
		if errors.Is(err, ErrNoGit) {
			return "", err
		}
		return "", fmt.Errorf("gitref: hash-object: %s", strings.TrimSpace(errb))
	}
	return strings.TrimSpace(out), nil
}

func (r *Repo) mktree(name, blob string) (string, error) {
	entry := fmt.Sprintf("100644 blob %s\t%s\n", blob, name)
	out, errb, err := r.run(entry, "mktree")
	if err != nil {
		if errors.Is(err, ErrNoGit) {
			return "", err
		}
		return "", fmt.Errorf("gitref: mktree: %s", strings.TrimSpace(errb))
	}
	return strings.TrimSpace(out), nil
}

func (r *Repo) push(id, commit, lease string) error {
	return r.pushRefspec(id, commit+":"+RefName(id), lease)
}

func (r *Repo) pushRefspec(id, refspec, lease string) error {
	args := []string{"push", "--force-with-lease=" + RefName(id) + ":" + lease, r.remote(), refspec}
	_, errb, err := r.run("", args...)
	if err != nil {
		if errors.Is(err, ErrNoGit) {
			return err
		}
		return classify(errb, err)
	}
	return nil
}

// classify separates losing a race from having no network, because the two need
// opposite handling: a lost race means re-read and tell the user who was
// quicker, an unreachable remote means keep the write and send it later.
// Retrying a lost race would only lose again; discarding an unsent write would
// lose the ticket.
func classify(stderr string, err error) error {
	s := strings.ToLower(stderr)
	switch {
	case strings.Contains(s, "stale info"),
		strings.Contains(s, "non-fast-forward"),
		strings.Contains(s, "fetch first"),
		strings.Contains(s, "rejected"):
		return fmt.Errorf("%w: %s", ErrRaceLost, firstLine(stderr))
	case strings.Contains(s, "not a git repository"),
		strings.Contains(s, "does not appear to be a git repository"),
		strings.Contains(s, "no such remote"):
		return fmt.Errorf("%w: %s", ErrNoRepo, firstLine(stderr))
	case strings.Contains(s, "could not resolve host"),
		strings.Contains(s, "could not read from remote"),
		strings.Contains(s, "unable to access"),
		strings.Contains(s, "connection"),
		strings.Contains(s, "network is unreachable"),
		strings.Contains(s, "timed out"),
		strings.Contains(s, "terminal prompts disabled"),
		strings.Contains(s, "authentication failed"):
		return fmt.Errorf("%w: %s", ErrOffline, firstLine(stderr))
	}
	if msg := firstLine(stderr); msg != "" {
		return fmt.Errorf("gitref: %s", msg)
	}
	return fmt.Errorf("gitref: %w", err)
}

func firstLine(s string) string {
	for _, line := range strings.Split(s, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			return line
		}
	}
	return ""
}
