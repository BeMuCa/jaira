package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// pickerModel is the seeded board — everything berk's — plus a ticket
// assigned to sam and one written by Alexander Sacharov, whose name has the
// space a filter has to keep together.
func pickerModel(t *testing.T) *Model {
	t.Helper()
	m := newTestModel(t, 160, 40)
	for i, f := range []map[string]string{
		{ticket.FieldTitle: "Sams ticket", ticket.FieldAssignee: "sam", ticket.FieldCreator: "berk"},
		{ticket.FieldTitle: "Alexanders ticket", ticket.FieldCreator: "Alexander Sacharov"},
	} {
		f[ticket.FieldID] = ticket.NewID(time.Now().Add(time.Duration(i) * time.Millisecond))
		f[ticket.FieldStatus] = "todo"
		if _, err := m.store.Create(f, nil, ""); err != nil {
			t.Fatal(err)
		}
	}
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	return m
}

// pickUser moves the picker's cursor to name and ticks it.
func pickUser(t *testing.T, m *Model, name string) {
	t.Helper()
	for i, u := range m.users {
		if u.name == name {
			m.userIdx = i
			m.key(key(" "))
			return
		}
	}
	t.Fatalf("%q is not in the picker: %v", name, m.users)
}

func boardTitles(m *Model) []string {
	var out []string
	for _, c := range m.cols {
		for _, tk := range c.tickets {
			out = append(out, tk.Title)
		}
	}
	return out
}

// u lists everyone a loaded ticket is assigned to or written by, the most
// frequent first, with how many tickets each has.
func TestUOpensTheUserPicker(t *testing.T) {
	m := pickerModel(t)
	m.key(key("u"))
	if m.mode != modeUsers {
		t.Fatalf("u must open the user picker, mode = %v", m.mode)
	}
	var got []string
	for _, u := range m.users {
		got = append(got, u.name)
	}
	if strings.Join(got, ",") != "berk,Alexander Sacharov,sam" {
		t.Errorf("picker lists %v, want berk first, then the others by name", got)
	}
	out := stripANSI(m.render())
	for _, want := range []string{"Users", "berk", "Alexander Sacharov", "sam"} {
		if !strings.Contains(out, want) {
			t.Errorf("the picker does not show %q", want)
		}
	}
}

// Ticking two people and pressing enter narrows the board to their tickets,
// through the ordinary filter.
func TestUserPickerNarrowsTheBoardToTheTickedPeople(t *testing.T) {
	m := pickerModel(t)
	m.key(key("u"))
	pickUser(t, m, "sam")
	pickUser(t, m, "Alexander Sacharov")
	m.key(key("enter"))

	if m.mode != modeBoard {
		t.Fatalf("enter must close the picker, mode = %v", m.mode)
	}
	if want := `user:"Alexander Sacharov",sam`; m.filter != want {
		t.Errorf("filter = %q, want %q", m.filter, want)
	}
	got := boardTitles(m)
	if strings.Join(got, ",") != "Alexanders ticket,Sams ticket" && strings.Join(got, ",") != "Sams ticket,Alexanders ticket" {
		t.Errorf("board shows %v, want only sam's and Alexander's tickets", got)
	}
}

// The picker owns only the user: part of the filter. What else was typed
// stays, the names already in it come back ticked, and x takes the people
// out again while keeping the rest.
func TestUserPickerKeepsTheRestOfTheFilter(t *testing.T) {
	m := pickerModel(t)
	m.filter, m.input = `"sams ticket"`, `"sams ticket"`
	m.rebuild()

	m.key(key("u"))
	pickUser(t, m, "sam")
	m.key(key("enter"))
	if want := `"sams ticket" user:sam`; m.filter != want {
		t.Fatalf("filter = %q, want %q", m.filter, want)
	}

	m.key(key("u"))
	for _, u := range m.users {
		if ticked := m.userPicked[strings.ToLower(u.name)]; ticked != (u.name == "sam") {
			t.Errorf("reopened, %q ticked = %v", u.name, ticked)
		}
	}
	m.key(key("x"))
	if want := `"sams ticket"`; m.filter != want || m.mode != modeBoard {
		t.Errorf("after x: filter = %q mode = %v, want %q on the board", m.filter, m.mode, want)
	}
}

// esc closes the picker and changes nothing, ticks included.
func TestUserPickerEscChangesNothing(t *testing.T) {
	m := pickerModel(t)
	m.key(key("u"))
	pickUser(t, m, "sam")
	m.key(key("esc"))
	if m.mode != modeBoard || m.filter != "" {
		t.Errorf("after esc: mode = %v filter = %q, want the board unfiltered", m.mode, m.filter)
	}
}

// The header shows the filter as it was written: a quoted name reads as one,
// not as escaped quotes.
func TestHeaderShowsTheFilterAsWritten(t *testing.T) {
	m := pickerModel(t)
	m.filter = `user:"Alexander Sacharov"`
	m.rebuild()
	out := stripANSI(m.render())
	if !strings.Contains(out, `filter: user:"Alexander Sacharov"`) || strings.Contains(out, `\"`) {
		t.Errorf("header does not show the filter as written:\n%s", strings.SplitN(out, "\n", 3)[1])
	}
}

// esc on the board clears a filter the picker wrote, like any other.
func TestEscOnTheBoardClearsTheUserFilter(t *testing.T) {
	m := pickerModel(t)
	all := len(boardTitles(m))
	m.key(key("u"))
	pickUser(t, m, "sam")
	m.key(key("enter"))
	if len(boardTitles(m)) >= all {
		t.Fatal("setup: the user filter did not narrow the board")
	}
	m.key(key("esc"))
	if m.filter != "" || len(boardTitles(m)) != all {
		t.Errorf("after esc: filter = %q, %d of %d tickets shown", m.filter, len(boardTitles(m)), all)
	}
}
