package merge

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/lane"
)

func lanes(t *testing.T) *lane.Set {
	t.Helper()
	t.Setenv("JAIRA_HOME", filepath.Join(t.TempDir(), "home"))
	t.Setenv("JAIRA_LANES_DIR", filepath.Join(t.TempDir(), "none"))
	s, err := lane.Load("")
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func doc(status, assignee, goal, updated string, blockedBy, commits []string) string {
	var b strings.Builder
	b.WriteString("---\nid: 01TEST\ntitle: A ticket\n")
	b.WriteString("status: " + status + "\n")
	b.WriteString("assignee: " + assignee + "\n")
	b.WriteString("goal: " + goal + "\n")
	b.WriteString("blocked-by: [" + strings.Join(blockedBy, ", ") + "]\n")
	b.WriteString("commits: [" + strings.Join(commits, ", ") + "]\n")
	b.WriteString("updated-at: " + updated + "\n---\n\nbody\n")
	return b.String()
}

func mergeStr(t *testing.T, base, ours, theirs string) *Result {
	t.Helper()
	r, err := Merge([]byte(base), []byte(ours), []byte(theirs), lanes(t))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestConcurrentLaneMovesDoNotConflict(t *testing.T) {
	base := doc("todo", "berk", "g", "2026-08-11T10:00:00Z", nil, nil)
	ours := doc("review", "berk", "g", "2026-08-11T10:05:00Z", nil, nil)
	theirs := doc("in-progress", "berk", "g", "2026-08-11T10:09:00Z", nil, nil)

	r := mergeStr(t, base, ours, theirs)
	if !r.Clean() {
		t.Fatalf("expected a clean merge, got %v", r.Conflicts)
	}
	// Review is further along than in-progress, so it must win even though the
	// other side is more recent: reverting progress is the worse failure.
	if !strings.Contains(string(r.Merged), "status: review") {
		t.Errorf("expected status to stay at review:\n%s", r.Merged)
	}
}

func TestListsUnion(t *testing.T) {
	base := doc("todo", "berk", "g", "2026-08-11T10:00:00Z", nil, nil)
	ours := doc("todo", "berk", "g", "2026-08-11T10:01:00Z", []string{"01AAA"}, []string{"c1"})
	theirs := doc("todo", "berk", "g", "2026-08-11T10:02:00Z", []string{"01BBB"}, []string{"c2"})

	r := mergeStr(t, base, ours, theirs)
	if !r.Clean() {
		t.Fatalf("unexpected conflicts: %v", r.Conflicts)
	}
	got := string(r.Merged)
	for _, want := range []string{"01AAA", "01BBB", "c1", "c2"} {
		if !strings.Contains(got, want) {
			t.Errorf("union lost %q:\n%s", want, got)
		}
	}
}

func TestContestedScalarResolvesByRecency(t *testing.T) {
	base := doc("todo", "berk", "g", "2026-08-11T10:00:00Z", nil, nil)
	ours := doc("todo", "alice", "g", "2026-08-11T10:01:00Z", nil, nil)
	theirs := doc("todo", "bob", "g", "2026-08-11T10:30:00Z", nil, nil)

	r := mergeStr(t, base, ours, theirs)
	if !r.Clean() {
		t.Fatalf("unexpected conflicts: %v", r.Conflicts)
	}
	if !strings.Contains(string(r.Merged), "assignee: bob") {
		t.Errorf("expected the newer side to win:\n%s", r.Merged)
	}
	if len(r.Notes) == 0 {
		t.Error("an automatic resolution should be explained in Notes")
	}
}

func TestOneSidedChangeIsTakenWithoutConflict(t *testing.T) {
	base := doc("todo", "berk", "old goal", "2026-08-11T10:00:00Z", nil, nil)
	ours := doc("todo", "berk", "old goal", "2026-08-11T10:00:00Z", nil, nil)
	theirs := doc("todo", "berk", "new goal", "2026-08-11T10:05:00Z", nil, nil)

	r := mergeStr(t, base, ours, theirs)
	if !r.Clean() {
		t.Fatalf("a one-sided prose edit must not conflict: %v", r.Conflicts)
	}
	if !strings.Contains(string(r.Merged), "new goal") {
		t.Errorf("their edit was dropped:\n%s", r.Merged)
	}
}

func TestCompetingProseRewritesDoConflict(t *testing.T) {
	// This is the case no rule can settle, and it must not be silently resolved.
	base := doc("todo", "berk", "original", "2026-08-11T10:00:00Z", nil, nil)
	ours := doc("todo", "berk", "our rewrite", "2026-08-11T10:01:00Z", nil, nil)
	theirs := doc("todo", "berk", "their rewrite", "2026-08-11T10:02:00Z", nil, nil)

	r := mergeStr(t, base, ours, theirs)
	if r.Clean() {
		t.Fatal("competing prose rewrites should conflict rather than pick a winner")
	}
	if len(r.Conflicts) != 1 || r.Conflicts[0].Field != "goal" {
		t.Errorf("conflict should be scoped to the goal field, got %+v", r.Conflicts)
	}
	// And the report must name both sides so a human can choose.
	out := RenderConflicts(r.Conflicts)
	for _, want := range []string{"our rewrite", "their rewrite", "original"} {
		if !strings.Contains(out, want) {
			t.Errorf("conflict report omitted %q:\n%s", want, out)
		}
	}
}

func TestMergePreservesUnknownBlockAndFormatting(t *testing.T) {
	mk := func(status, updated string) string {
		return "---\n" +
			"# hand-written note\n" +
			"id: 01TEST\n" +
			"title: A ticket\n" +
			"status: " + status + "          # aligned comment\n" +
			"assignee: berk\n" +
			"goal: g\n" +
			"blocked-by: []\n" +
			"commits: []\n" +
			"external:\n  provider: jira\n  key: PROJ-1\n" +
			"updated-at: " + updated + "\n---\n\nbody\n"
	}
	r := mergeStr(t, mk("todo", "2026-08-11T10:00:00Z"),
		mk("todo", "2026-08-11T10:01:00Z"), mk("review", "2026-08-11T10:02:00Z"))
	if !r.Clean() {
		t.Fatalf("unexpected conflicts: %v", r.Conflicts)
	}
	got := string(r.Merged)
	if !strings.Contains(got, "# hand-written note") {
		t.Error("comment lost in merge")
	}
	if !strings.Contains(got, "external:\n  provider: jira\n  key: PROJ-1\n") {
		t.Errorf("opaque block damaged:\n%s", got)
	}
	if !strings.Contains(got, "# aligned comment") {
		t.Error("trailing comment lost")
	}
	if !strings.Contains(got, "status: review") {
		t.Error("status not merged")
	}
}

func TestMergeNeverProducesInvalidYAML(t *testing.T) {
	base := doc("todo", "berk", "g", "2026-08-11T10:00:00Z", nil, nil)
	ours := doc("review", "alice", "g", "2026-08-11T10:01:00Z", []string{"01AAA"}, []string{"c1"})
	theirs := doc("done", "bob", "g", "2026-08-11T10:02:00Z", []string{"01BBB"}, []string{"c2"})

	r := mergeStr(t, base, ours, theirs)
	// Exactly one status line, whatever the outcome.
	if n := strings.Count(string(r.Merged), "status:"); n != 1 {
		t.Errorf("expected exactly one status key, found %d:\n%s", n, r.Merged)
	}
	if strings.Contains(string(r.Merged), "<<<<<<<") {
		t.Error("merge emitted conflict markers into the frontmatter")
	}
}

func TestBodyConflictIsReported(t *testing.T) {
	mk := func(body string) string {
		return "---\nid: 01TEST\nstatus: todo\nupdated-at: 2026-08-11T10:00:00Z\n---\n\n" + body + "\n"
	}
	r := mergeStr(t, mk("base text"), mk("our text"), mk("their text"))
	if r.Clean() {
		t.Fatal("competing body edits should conflict")
	}
	if r.Conflicts[0].Field != "body" {
		t.Errorf("expected a body conflict, got %+v", r.Conflicts)
	}
}

// git hands the merge driver an empty base for an add/add conflict: two
// branches that each created the same ticket file. Rejecting that used to fail
// the driver, and the damage was silent — git left "ours" behind with no
// conflict markers, so it read as a clean merge while the other side's lane was
// gone.
func TestAnEmptyBaseIsAnAbsentAncestorNotABrokenFile(t *testing.T) {
	ls := lanes(t)
	ours := []byte("---\nid: 01AAA\nstatus: backlog\nassignee: ada\nupdated-by: ada\nupdated-at: 2026-09-10T12:00:00Z\ntags:\n  - concurrency\n---\n\n# a ticket\n")
	theirs := []byte("---\nid: 01AAA\nstatus: todo\nassignee: berk\nupdated-by: berk\nupdated-at: 2026-09-09T12:00:00Z\ntags:\n  - cli\n---\n\n# a ticket\n")

	for _, base := range [][]byte{nil, {}, []byte("   \n"), []byte("\n\n")} {
		res, err := Merge(base, ours, theirs, ls)
		if err != nil {
			t.Fatalf("empty base (%q) was rejected: %v", base, err)
		}
		got := string(res.Merged)
		// Progress is never reverted, even though ours carries the newer stamp.
		if !strings.Contains(got, "status: todo") {
			t.Errorf("base %q: the further lane was lost:\n%s", base, got)
		}
		// Neither side's tag is dropped: with no ancestor, nothing was removed
		// by anyone, so both are additions.
		for _, tag := range []string{"concurrency", "cli"} {
			if !strings.Contains(got, tag) {
				t.Errorf("base %q: tag %q was lost:\n%s", base, tag, got)
			}
		}
	}

	// Prose is the one thing no rule can settle, and an absent ancestor does
	// not change that: two different sentences in the same field stay a real
	// conflict for a person, rather than one of them being picked.
	oursProse := []byte("---\nid: 01AAA\nstatus: todo\ngoal: the cookie survives the redirect\nupdated-at: 2026-09-10T12:00:00Z\n---\n\nbody\n")
	theirsProse := []byte("---\nid: 01AAA\nstatus: todo\ngoal: sessions must not be dropped on 302\nupdated-at: 2026-09-09T12:00:00Z\n---\n\nbody\n")
	res, err := Merge(nil, oursProse, theirsProse, ls)
	if err != nil {
		t.Fatalf("empty base with differing prose: %v", err)
	}
	if res.Clean() {
		t.Error("two different goals merged silently; one of them was picked")
	}
	var found bool
	for _, c := range res.Conflicts {
		if c.Field == "goal" {
			found = true
		}
	}
	if !found {
		t.Errorf("the conflict is not on goal: %+v", res.Conflicts)
	}
}
