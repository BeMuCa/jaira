package tui

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/lane"
)

// B on the lane settings screen opens the new-board form; a name, a
// description and the ticked lanes create the board and save it into the
// catalogue. Keys go through the real dispatch, so q typed into the name is
// text, not a way out.
func TestLaneScreenCreatesABoardFromTheForm(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())
	m := newTestModel(t, 150, 40)
	m.mode, m.laneScreen = modeLanes, newLaneScreen(m.store, m.lanes)

	press(m, "B")
	ls := m.laneScreen
	if ls == nil || ls.form == nil {
		t.Fatal("B did not open the new-board form")
	}
	for _, r := range "quick" {
		press(m, string(r))
	}
	press(m, "enter")

	// Without a description the board is refused and the form stays open.
	press(m, "tab")
	press(m, "enter")
	if ls.form == nil || !strings.Contains(ls.form.err, "description") {
		t.Fatalf("a board without a description was not refused: form %v", ls.form)
	}
	press(m, "tab") // round to the name, then on to the description
	press(m, "tab")
	pasteInto(m, "Quick fixes\n")
	press(m, "tab")
	if out := stripANSI(m.render()); !strings.Contains(out, "New board") || !strings.Contains(out, "Quick fixes") {
		t.Errorf("the form is not on screen:\n%s", out)
	}

	// Today's lanes start ticked; untick the first.
	want := m.lanes.IDs()[1:]
	press(m, " ")
	press(m, "enter")
	if ls.form != nil {
		t.Fatalf("the form stayed open: %s", ls.form.err)
	}
	set, err := lane.LoadBoard(m.store.Root, "quick")
	if err != nil {
		t.Fatal(err)
	}
	if got := set.IDs(); !slices.Equal(got, want) {
		t.Errorf("board lanes = %v, want %v", got, want)
	}
	if d := lane.Description(m.store.Root, "quick"); d != "Quick fixes" {
		t.Errorf("description = %q", d)
	}
	if _, err := os.Stat(filepath.Join(lane.UserBoardsDir(), "quick", "order")); err != nil {
		t.Errorf("the board was not saved to the catalogue: %v", err)
	}
	if !strings.Contains(ls.msg, "catalogue") || ls.isErr {
		t.Errorf("message = %q (err %v)", ls.msg, ls.isErr)
	}
}
