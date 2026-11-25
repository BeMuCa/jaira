package lane

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"testing/fstest"
)

// The code type names the built-in lanes a fresh board gets, in their
// shipped order — so a lane that joins or leaves the default selection
// cannot leave 'jaira boards create' building yesterday's board.
func TestCodeTypeIsTheBuiltinDefaultLanes(t *testing.T) {
	builtins, err := Builtins()
	if err != nil {
		t.Fatal(err)
	}
	var want []string
	for _, l := range builtins {
		if l.Default {
			want = append(want, l.ID)
		}
	}
	types, err := BoardTypes()
	if err != nil {
		t.Fatal(err)
	}
	i := slices.IndexFunc(types, func(bt BoardType) bool { return bt.Name == DefaultBoardType })
	if i < 0 {
		t.Fatalf("no %q type among %v", DefaultBoardType, types)
	}
	if !slices.Equal(types[i].Lanes, want) {
		t.Errorf("code type lanes = %v, want the built-in defaults %v", types[i].Lanes, want)
	}
	if types[i].Description == "" {
		t.Error("code type has no description")
	}
}

// A type with lanes of its own — the shape an ideas board takes — writes
// those files, takes the ids it has no file for from the built-ins, and the
// board it builds loads with exactly those lanes in that order.
func TestCreateBoardFromATypeWithLanesOfItsOwn(t *testing.T) {
	old := boardTypesFS
	t.Cleanup(func() { boardTypesFS = old })
	boardTypesFS = fstest.MapFS{
		"boardtypes/ideas/description": {Data: []byte("A topic, several ideas with a roadmap each, a person picks one.\n")},
		"boardtypes/ideas/order":       {Data: []byte("backlog\nideate\npicked\n")},
		"boardtypes/ideas/ideate.md":   {Data: []byte("---\nid: ideate\nname: Ideate\nagentic: true\nmodel-tier: strong\n---\n\n# Prompt\n\nWrite three ideas.\n")},
		"boardtypes/ideas/picked.md":   {Data: []byte("---\nid: picked\nname: Picked\nterminal: true\n---\n")},
	}
	root := newBoardRoot(t)

	dir, ids, err := CreateBoard(root, "ideas", "ideas", "Hackathon ideas")
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"backlog", "ideate", "picked"}; !slices.Equal(ids, want) {
		t.Errorf("ids = %v, want %v", ids, want)
	}
	set, err := LoadBoard(root, "ideas")
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, l := range set.Lanes {
		got = append(got, l.ID)
	}
	if !slices.Equal(got, ids) {
		t.Errorf("loaded board lanes = %v, want %v", got, ids)
	}
	if d := Description(root, "ideas"); d != "Hackathon ideas" {
		t.Errorf("description = %q, want the one given, not the type's", d)
	}
	if _, err := os.Stat(filepath.Join(dir, "backlog.md")); err != nil {
		t.Errorf("built-in lane not written: %v", err)
	}

	// A second create never overwrites the board it would replace.
	before, _ := os.ReadFile(filepath.Join(dir, "ideate.md"))
	if _, _, err := CreateBoard(root, "ideas", "ideas", "other"); !errors.Is(err, ErrBoardExists) {
		t.Errorf("create over an existing board: %v, want ErrBoardExists", err)
	}
	after, _ := os.ReadFile(filepath.Join(dir, "ideate.md"))
	if string(before) != string(after) || Description(root, "ideas") != "Hackathon ideas" {
		t.Error("refused create still changed the board")
	}
}

func TestCreateBoardRefusesTodaysBoardAndUnknownTypes(t *testing.T) {
	root := newBoardRoot(t)
	if _, _, err := CreateBoard(root, "default", DefaultBoardType, "x"); !errors.Is(err, ErrBoardExists) {
		t.Errorf("create default: %v, want ErrBoardExists", err)
	}
	if _, _, err := CreateBoard(root, "video", "nope", "x"); !errors.Is(err, ErrNoSuchBoardType) {
		t.Errorf("unknown type: %v, want ErrNoSuchBoardType", err)
	}
	if Boards(root) != nil {
		t.Errorf("a refused create left a board behind: %v", Boards(root))
	}
}
