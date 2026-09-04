package cli

import (
	"reflect"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// --append keeps history: loop rounds add to a review field instead of
// overwriting the previous round's reasoning.
func TestSetAppendKeepsTheOldValue(t *testing.T) {
	t.Setenv("JAIRA_USER", "berk")
	dir := t.TempDir()
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	tk, err := s.Create(map[string]string{
		ticket.FieldID:    ticket.NewID(time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)),
		ticket.FieldTitle: "t", ticket.FieldStatus: "backlog",
	}, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	h := ticket.Handle(tk.ID)

	// On an empty field --append behaves like a plain set.
	if out, err := runCLI(t, dir, "set", h, "review-gaps=runde eins", "--append"); err != nil {
		t.Fatalf("set --append on empty: %v\n%s", err, out)
	}
	// The second round lands below the first, not over it.
	if out, err := runCLI(t, dir, "set", h, "review-gaps=runde zwei", "--append"); err != nil {
		t.Fatalf("set --append: %v\n%s", err, out)
	}
	got, err := s.Load(tk.ID)
	if err != nil {
		t.Fatal(err)
	}
	if want := "runde eins\nrunde zwei"; got.ReviewGaps != want {
		t.Errorf("review-gaps = %q, want %q", got.ReviewGaps, want)
	}
	// Without the flag the old behaviour holds: overwrite.
	if out, err := runCLI(t, dir, "set", h, "review-gaps=nur das"); err != nil {
		t.Fatalf("set: %v\n%s", err, out)
	}
	if got, _ := s.Load(tk.ID); got.ReviewGaps != "nur das" {
		t.Errorf("plain set no longer overwrites: %q", got.ReviewGaps)
	}
}

func TestSetAppendExtendsAList(t *testing.T) {
	t.Setenv("JAIRA_USER", "berk")
	dir := t.TempDir()
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	tk, err := s.Create(map[string]string{
		ticket.FieldID:    ticket.NewID(time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)),
		ticket.FieldTitle: "t", ticket.FieldStatus: "backlog",
	}, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	h := ticket.Handle(tk.ID)
	if _, err := runCLI(t, dir, "set", h, "commits=aaa1111"); err != nil {
		t.Fatal(err)
	}
	if _, err := runCLI(t, dir, "set", h, "commits=bbb2222", "--append"); err != nil {
		t.Fatal(err)
	}
	got, err := s.Load(tk.ID)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"aaa1111", "bbb2222"}; !reflect.DeepEqual(got.Commits, want) {
		t.Errorf("commits = %v, want %v", got.Commits, want)
	}
}
