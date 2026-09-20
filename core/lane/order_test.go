package lane

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestMoveSwapsWithNeighbour(t *testing.T) {
	got := Move([]string{"a", "b", "c"}, "c", -1)
	want := []string{"a", "c", "b"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Move(-1) = %v, want %v", got, want)
	}

	got = Move([]string{"a", "c", "b"}, "c", 1)
	want = []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Move(+1) = %v, want %v", got, want)
	}
}

func TestMoveAtEitherEndIsNoOp(t *testing.T) {
	ids := []string{"a", "b", "c"}
	if got := Move(ids, "a", -1); !reflect.DeepEqual(got, ids) {
		t.Errorf("moving the first lane left = %v, want unchanged %v", got, ids)
	}
	if got := Move(ids, "c", 1); !reflect.DeepEqual(got, ids) {
		t.Errorf("moving the last lane right = %v, want unchanged %v", got, ids)
	}
}

func TestMoveSingleElementIsNoOp(t *testing.T) {
	ids := []string{"only"}
	if got := Move(ids, "only", -1); !reflect.DeepEqual(got, ids) {
		t.Errorf("Move(-1) on single element = %v, want %v", got, ids)
	}
	if got := Move(ids, "only", 1); !reflect.DeepEqual(got, ids) {
		t.Errorf("Move(+1) on single element = %v, want %v", got, ids)
	}
}

func TestMoveUnknownIDIsNoOp(t *testing.T) {
	ids := []string{"a", "b"}
	if got := Move(ids, "ghost", -1); !reflect.DeepEqual(got, ids) {
		t.Errorf("Move on unknown id = %v, want unchanged %v", got, ids)
	}
}

func TestLoadOrderAbsentFileIsNotAnError(t *testing.T) {
	root := t.TempDir()
	ids, err := LoadOrder(root)
	if err != nil {
		t.Fatalf("LoadOrder on absent file: %v", err)
	}
	if ids != nil {
		t.Errorf("ids = %v, want nil", ids)
	}
}

func TestSaveThenLoadOrderRoundTrips(t *testing.T) {
	root := t.TempDir()
	want := []string{"backlog", "in-progress", "done"}
	if err := SaveOrder(root, want); err != nil {
		t.Fatalf("SaveOrder: %v", err)
	}
	got, err := LoadOrder(root)
	if err != nil {
		t.Fatalf("LoadOrder: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("LoadOrder = %v, want %v", got, want)
	}
}

// TestBuiltinShadowsAnAdoptedCopyOfTheSameID pins a precedence that used to be
// invisible and is now reachable: Installable seeds `seen` with the built-ins
// before it globs UserLanesDir, so a file in the user's catalogue carrying a
// built-in's id is skipped — the embedded copy wins.
//
// This always held for the ten lanes a board starts with, where nobody had a
// catalogue copy to begin with. It began to matter when critique, optimize and
// testing moved into the binary: until then those three existed ONLY as
// catalogue files, so anybody who ran 'jaira lanes market adopt critique' and
// edited the result now has an edited file that no longer reaches a board.
// Renaming it to an id of their own is the way back.
//
// The behaviour is deliberate — a lane shipped with the binary is the one the
// binary's prompts and gates were written against — but it is deliberate only
// as long as a test says so out loud. Change the loop order in Installable and
// this test fails rather than a user's edited prompt silently coming back.
func TestBuiltinShadowsAnAdoptedCopyOfTheSameID(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	const mine = "MY OWN ADOPTED CRITIQUE"
	writeLane(t, catalogue, "critique.md",
		"---\nid: critique\nname: Critique\ndescription: "+mine+"\nafter: in-progress\nprecedence: 45\n---\nmy own prompt\n")

	offer, err := Installable(&Set{})
	if err != nil {
		t.Fatal(err)
	}
	var got *Lane
	for _, l := range offer {
		if l.ID == "critique" {
			if got != nil {
				t.Fatal("critique is offered twice; the offer must hold one lane per id")
			}
			got = l
		}
	}
	if got == nil {
		t.Fatal("critique is not on offer at all")
	}
	if got.Description == mine {
		t.Fatal("the adopted file won: an edited ~/.jaira/lanes/critique.md must not shadow the built-in")
	}

	builtins, err := Builtins()
	if err != nil {
		t.Fatal(err)
	}
	var want string
	for _, l := range builtins {
		if l.ID == "critique" {
			want = l.Description
		}
	}
	if want == "" {
		t.Fatal("no built-in critique — this test has lost its subject")
	}
	if got.Description != want {
		t.Fatalf("offered critique = %q, want the built-in's %q", got.Description, want)
	}

	// And what 'jaira lanes add critique' writes onto a board is that same file.
	root := newBoardRoot(t)
	set, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	dst, _, _, err := Add(root, set, "critique")
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), mine) {
		t.Fatalf("'lanes add critique' installed the adopted file, not the built-in: %s", dst)
	}
}
