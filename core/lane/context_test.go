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
		if l.Fresh() {
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

// TestContextRejectsAnUnknownValue: a typo must not silently mean "not fresh"
// — that is the reviewer quietly turning into the implementer.
func TestContextRejectsAnUnknownValue(t *testing.T) {
	_, err := Parse([]byte("---\nid: x\ncontext: frsh\n---\n"), "test")
	if err == nil || !strings.Contains(err.Error(), "context") {
		t.Fatalf("parse error = %v, want one naming context", err)
	}
	l, err := Parse([]byte("---\nid: x\ncontext: \"\"\n---\n"), "test")
	if err != nil || l.Context != "" {
		t.Fatalf("an empty context must parse as undeclared: lane=%+v err=%v", l, err)
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
