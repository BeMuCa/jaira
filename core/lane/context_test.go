package lane

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// freshIDs lists the lanes of a set that run in a fresh context, each with
// "*" appended when the board's default put it there rather than its file.
func freshIDs(set *Set) []string {
	var out []string
	for _, l := range set.Lanes {
		if l.Context == ContextFresh {
			id := l.ID
			if l.ContextImplied {
				id += "*"
			}
			out = append(out, id)
		}
	}
	return out
}

// setUpBoard loads a fresh board from the built-ins, so its lane files exist
// on disk to be edited, and returns its root.
func setUpBoard(t *testing.T) string {
	t.Helper()
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	t.Setenv("JAIRA_DEFAULT_BOARD", filepath.Join(t.TempDir(), "none.md"))
	root := newBoardRoot(t)
	if _, err := Load(root); err != nil {
		t.Fatal(err)
	}
	return root
}

// addFrontmatter writes one line into the frontmatter of a board's lane file.
func addFrontmatter(t *testing.T, root, id, line string) {
	t.Helper()
	p := filepath.Join(ProjectLanesDir(root), id+".md")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	s := strings.Replace(string(b), "---\n", "---\n"+line+"\n", 1)
	if err := os.WriteFile(p, []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestContextDefaultsToReviewOnTheShippedBoard: no shipped lane says
// context:, so the default decides — and on every lane the binary ships,
// critique, optimize and testing included, it lands on review. HITL sits
// between testing and review and only takes a question; counting it as the
// first lane a person owns would make testing the fresh lane.
func TestContextDefaultsToReviewOnTheShippedBoard(t *testing.T) {
	t.Setenv("JAIRA_LANES_DIR", t.TempDir())
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(freshIDs(set), ","); got != "review*" {
		t.Fatalf("fresh lanes on the shipped offer = %q, want %q", got, "review*")
	}

	root := setUpBoard(t)
	board, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(freshIDs(board), ","); got != "review*" {
		t.Fatalf("fresh lanes on a board set up from the built-ins = %q, want %q", got, "review*")
	}
}

// TestContextDeclaredReplacesTheDefault: a lane that says context: fresh is
// fresh, and it switches the default off for the whole board — review is no
// longer fresh just because nobody wrote on it.
func TestContextDeclaredReplacesTheDefault(t *testing.T) {
	root := setUpBoard(t)
	addFrontmatter(t, root, "in-progress", "context: fresh")
	set, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(freshIDs(set), ","); got != "in-progress" {
		t.Fatalf("fresh lanes = %q, want only the declared %q", got, "in-progress")
	}
}

// TestContextSharedSwitchesTheDefaultOff: shared is how a board says it wants
// no fresh lane at all. Without it the default would always find one.
func TestContextSharedSwitchesTheDefaultOff(t *testing.T) {
	root := setUpBoard(t)
	addFrontmatter(t, root, "review", "context: shared")
	set, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if got := freshIDs(set); len(got) != 0 {
		t.Fatalf("fresh lanes = %v, want none: review says shared", got)
	}
	if l, _ := set.Get("review"); l.Context != ContextShared {
		t.Errorf("review context = %q, want %q", l.Context, ContextShared)
	}
}

// TestContextUnknownValueRunsFresh: a typo must neither silently mean "not
// fresh" — the reviewer quietly turning into the implementer — nor drop the
// lane, which would let the board's default move the fresh boundary onto
// another lane. The lane loads, runs fresh, and the value is kept for the
// warning.
func TestContextUnknownValueRunsFresh(t *testing.T) {
	l, err := Parse([]byte("---\nid: x\ncontext: frsh\n---\n"), "test")
	if err != nil {
		t.Fatalf("an unknown context must not fail the lane: %v", err)
	}
	if l.Context != ContextFresh || l.ContextInvalid != "frsh" {
		t.Fatalf("context=%q invalid=%q, want %q/%q", l.Context, l.ContextInvalid, ContextFresh, "frsh")
	}
	l, err = Parse([]byte("---\nid: x\ncontext: \"\"\n---\n"), "test")
	if err != nil || l.Context != "" || l.ContextInvalid != "" {
		t.Fatalf("an empty context must parse as undeclared: lane=%+v err=%v", l, err)
	}
}

// TestContextTypoKeepsTheFreshBoundary: review.md with "context: Fresh" stays
// on the board as the fresh lane, the default does not move to in-progress,
// and Load says what is wrong.
func TestContextTypoKeepsTheFreshBoundary(t *testing.T) {
	root := setUpBoard(t)
	addFrontmatter(t, root, "review", "context: Fresh")
	set, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(freshIDs(set), ","); got != "review" {
		t.Fatalf("fresh lanes = %q, want %q: the typo must neither drop review nor move the default", got, "review")
	}
	if !strings.Contains(strings.Join(set.Warnings, "\n"), `lane review: context "Fresh"`) {
		t.Errorf("no warning names the bad value:\n%s", strings.Join(set.Warnings, "\n"))
	}
}

// TestContextDefaultNeedsALanePersonOwns: a board where nothing waits for a
// person has no "last look before a person's", and the default marks nothing
// rather than picking a lane at random.
func TestContextDefaultNeedsALanePersonOwns(t *testing.T) {
	root := setUpBoard(t)
	if err := os.Remove(filepath.Join(ProjectLanesDir(root), "signoff.md")); err != nil {
		t.Fatal(err)
	}
	set, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if got := freshIDs(set); len(got) != 0 {
		t.Fatalf("fresh lanes = %v, want none on a board without signoff", got)
	}
}

// TestContextTypoAddsFreshnessNeverRemovesIt: an invalid value is no
// declaration. "context: Shared" on in-progress runs in-progress fresh and
// leaves the default's fresh review standing; Load names the value.
func TestContextTypoAddsFreshnessNeverRemovesIt(t *testing.T) {
	root := setUpBoard(t)
	addFrontmatter(t, root, "in-progress", "context: Shared")
	set, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(freshIDs(set), ","); got != "in-progress,review*" {
		t.Fatalf("fresh lanes = %q, want %q: a typo must not switch the default off", got, "in-progress,review*")
	}
	if !strings.Contains(strings.Join(set.Warnings, "\n"), `lane in-progress: context "Shared"`) {
		t.Errorf("no warning names the bad value:\n%s", strings.Join(set.Warnings, "\n"))
	}
}

// TestContextValidSharedWinsOverATypo: once a lane carries a valid value the
// board has declared itself. review saying shared stays shared, the typo on
// in-progress still runs that lane fresh, and nothing is implied.
func TestContextValidSharedWinsOverATypo(t *testing.T) {
	root := setUpBoard(t)
	addFrontmatter(t, root, "review", "context: shared")
	addFrontmatter(t, root, "in-progress", "context: Shared")
	set, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(freshIDs(set), ","); got != "in-progress" {
		t.Fatalf("fresh lanes = %q, want only the typo's %q", got, "in-progress")
	}
	if l, _ := set.Get("review"); l.Context != ContextShared {
		t.Errorf("review context = %q, want %q", l.Context, ContextShared)
	}
}
