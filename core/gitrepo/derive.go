package gitrepo

import (
	"fmt"
	"sort"
	"strings"

	"github.com/BeMuCa/jaira/core/ticket"
)

// CommitsForTicket works out which commits carry a ticket, so a person never
// has to type a sha by hand for the tool to know it. The result is the union
// of two independent sources, because neither alone finds every commit that
// belongs to a ticket: a commit can touch the ticket file without naming the
// ticket in its message (a rename, a bulk edit), and a commit can name the
// ticket without touching its file (the source change the ticket describes).
// The returned shas are oldest-first, since the recorded list is read as the
// history of the ticket and every caller records it the same way.
func (r *Repo) CommitsForTicket(ticketPath, ticketID string) ([]string, error) {
	if !Available() {
		return nil, ErrNoGit
	}
	// An id this tool would not have written is not searched for. Both forms go
	// into an extended regular expression below, and nothing escapes them —
	// a hand-edited "id: .*" would match every commit in the repository and
	// stamp the lot onto the ticket, silently, since archive derives
	// best-effort and discards its error. Load does not reject a malformed id
	// (only 'jaira validate' reports one, as an error), so the file format
	// being hand-editable makes this reachable rather than theoretical.
	//
	// Deriving nothing is the right answer here, not deriving something: an id
	// that is not a ULID is damage, and this list is what somebody reads months
	// later to recover why a change was made.
	if !ticket.ValidID(ticketID) {
		return nil, nil
	}

	found := map[string]bool{}

	// First source: every commit that touched the ticket file, followed
	// through renames so a `git mv`d ticket does not lose its history.
	if out, err := r.run("log", "--format=%H", "--follow", "--", ticketPath); err == nil {
		for _, line := range strings.Split(out, "\n") {
			if sha := strings.TrimSpace(line); sha != "" {
				found[sha] = true
			}
		}
	}

	// Second source: commits that name the ticket in their message without
	// necessarily touching its file. Both forms of the id are searched as
	// alternatives in one pattern, because the handle — not the 26-character
	// ULID — is the only form of the id this tool ever shows a person: every
	// gate refusal, every card, every line the CLI prints uses it. A human or
	// an agent writing a commit message writes the handle, not the ULID, so
	// searching the ULID alone would leave this half of the union finding
	// almost nothing and quietly reduce two sources to one.
	//
	// Six random characters are a real false-positive risk, so the handle
	// (and the ULID, for symmetry) is matched as a reference bounded by a
	// non-alphanumeric character or the edge of the line, never as a bare
	// substring — the list this produces is what somebody reads months later
	// to recover why a change was made, and a stranger's commit sitting in it
	// is worse than a missing one. Matching is case-insensitive because a
	// ULID is uppercase Crockford base32 and nobody types it that way, the
	// same reason ticket.NormalizeIDPrefix exists.
	if out, err := r.run("log", "--format=%H", "--extended-regexp", "--regexp-ignore-case", "--grep="+namesTicket(ticketID)); err == nil {
		for _, line := range strings.Split(out, "\n") {
			if sha := strings.TrimSpace(line); sha != "" {
				found[sha] = true
			}
		}
	}

	if len(found) == 0 {
		return nil, nil
	}
	shas := make([]string, 0, len(found))
	for sha := range found {
		shas = append(shas, sha)
	}
	if len(shas) == 1 {
		return shas, nil
	}

	// rev-list --no-walk=sorted returns the given shas ordered by commit date,
	// newest first; reversed here rather than asked for in the other
	// direction, so the flag stays the boring, well-known one and the
	// ordering choice stays visible in one place rather than hidden in a git
	// flag.
	out, err := r.run(append([]string{"rev-list", "--no-walk=sorted"}, shas...)...)
	if err != nil {
		// The union was already found by two log queries that succeeded; a
		// failure here is not a reason to report nothing, only to report it
		// unsorted.
		return shas, nil
	}
	sorted := make([]string, 0, len(shas))
	for _, line := range strings.Split(out, "\n") {
		if sha := strings.TrimSpace(line); sha != "" {
			sorted = append(sorted, sha)
		}
	}
	if len(sorted) != len(shas) {
		// Should not happen, but a length mismatch means the sort is not
		// trustworthy; fall back to a stable order rather than a partial one.
		sort.Strings(shas)
		return shas, nil
	}
	// sorted is newest-first; reverse in place for oldest-first.
	for i, j := 0, len(sorted)-1; i < j; i, j = i+1, j-1 {
		sorted[i], sorted[j] = sorted[j], sorted[i]
	}
	return sorted, nil
}

// namesTicket is the --extended-regexp pattern for a commit message that
// names the ticket: its ULID or its handle, bounded as a reference. The id
// must be a ValidID — nothing in it is escaped.
func namesTicket(ticketID string) string {
	return fmt.Sprintf(`(^|[^0-9A-Za-z])(%s|%s)([^0-9A-Za-z]|$)`, ticketID, ticket.Handle(ticketID))
}

// BranchesNamingTicket names the local branches that hold a commit naming the
// ticket which HEAD does not reach, sorted. It is what a refusal for "git has
// none for <id>" points at: the work was committed, only on another branch —
// a follow-up worked on the branch of the ticket it followed, a worktree the
// move was not run from. Commits HEAD reaches are already in
// CommitsForTicket's answer and name no branch here. Nothing found, or no
// git, is an empty answer: the refusal then stands as it always read.
func (r *Repo) BranchesNamingTicket(ticketID string) ([]string, error) {
	if !Available() {
		return nil, ErrNoGit
	}
	if !ticket.ValidID(ticketID) {
		return nil, nil
	}
	out, err := r.run("log", "--branches", "--not", "HEAD", "--format=%H", "--extended-regexp", "--regexp-ignore-case", "--grep="+namesTicket(ticketID))
	if err != nil {
		return nil, err
	}
	// Newest first; the first 64 are plenty to name the branches and keep the
	// one git branch call below its argument limit on any platform.
	args := []string{"branch", "--format=%(refname:short)"}
	for _, line := range strings.Split(out, "\n") {
		if sha := strings.TrimSpace(line); sha != "" && len(args) < 2+2*64 {
			args = append(args, "--contains", sha)
		}
	}
	if len(args) == 2 {
		return nil, nil
	}
	// Several --contains list every branch holding any of them, in one call.
	listed, err := r.run(args...)
	if err != nil {
		return nil, err
	}
	var branches []string
	for _, line := range strings.Split(listed, "\n") {
		if b := strings.TrimSpace(line); b != "" && !strings.HasPrefix(b, "(") {
			branches = append(branches, b)
		}
	}
	sort.Strings(branches)
	return branches, nil
}
