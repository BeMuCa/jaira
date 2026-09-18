package ticket

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// trimStore builds a store whose done lane holds n tickets, updated a minute
// apart, and returns them oldest first.
func trimStore(t *testing.T, n int) (*Store, []*Ticket) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))
	s, err := At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	base := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	var out []*Ticket
	for i := 0; i < n; i++ {
		stamp := base.Add(time.Duration(i) * time.Minute)
		tk, err := s.Create(map[string]string{
			FieldID:        NewID(stamp),
			FieldTitle:     fmt.Sprintf("t%02d", i),
			FieldStatus:    "done",
			FieldCreatedAt: FormatTime(stamp),
			FieldUpdatedAt: FormatTime(stamp),
		}, nil, "")
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, tk)
	}
	return s, out
}

// tiedStore builds a store with n done tickets created a second apart but all
// carrying the SAME updated-at — the routine case: several moves inside the
// stamp's one-second resolution.
func tiedStore(t *testing.T, n int) (*Store, []*Ticket) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))
	s, err := At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	base := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	stamp := FormatTime(base.Add(time.Hour))
	var out []*Ticket
	for i := 0; i < n; i++ {
		tk, err := s.Create(map[string]string{
			FieldID:        NewID(base.Add(time.Duration(i) * time.Second)),
			FieldTitle:     fmt.Sprintf("tie%d", i),
			FieldStatus:    "done",
			FieldCreatedAt: stamp,
			FieldUpdatedAt: stamp,
		}, nil, "")
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, tk)
	}
	return s, out
}

func TestOverflowReturnsExactlyTheOldestBeyondTheCap(t *testing.T) {
	s, ts := trimStore(t, 12)
	over, err := s.Overflow("done", 10, "")
	if err != nil {
		t.Fatalf("Overflow: %v", err)
	}
	if len(over) != 2 {
		t.Fatalf("overflow of 12 over 10 = %d tickets, want 2", len(over))
	}
	if over[0].ID != ts[0].ID || over[1].ID != ts[1].ID {
		t.Errorf("overflow = %s, %s — want the two oldest %s, %s",
			over[0].Title, over[1].Title, ts[0].Title, ts[1].Title)
	}
}

func TestOverflowAtOrUnderTheCapIsEmpty(t *testing.T) {
	s, _ := trimStore(t, 10)
	over, err := s.Overflow("done", 10, "")
	if err != nil {
		t.Fatalf("Overflow: %v", err)
	}
	if len(over) != 0 {
		t.Errorf("exactly 10 over a cap of 10 overflows %d tickets, want none", len(over))
	}
	// keep <= 0 means no cap at all.
	over, err = s.Overflow("done", 0, "")
	if err != nil {
		t.Fatalf("Overflow: %v", err)
	}
	if len(over) != 0 {
		t.Errorf("a cap of 0 (unlimited) overflows %d tickets, want none", len(over))
	}
}

// updated-at has second resolution, so scripted moves tie routinely. Among
// ties the newer ULID is the newer ticket — the oldest-created leaves, never
// the latest arrival.
func TestOverflowBreaksTiesTowardTheNewerULID(t *testing.T) {
	s, ts := tiedStore(t, 3)
	over, err := s.Overflow("done", 2, "")
	if err != nil {
		t.Fatalf("Overflow: %v", err)
	}
	if len(over) != 1 {
		t.Fatalf("overflow of 3 tied over 2 = %d tickets, want 1", len(over))
	}
	if over[0].ID != ts[0].ID {
		t.Errorf("overflow picked %s, want the oldest-created %s", over[0].Title, ts[0].Title)
	}
}

// The ticket whose arrival triggered the trim is pinned as newest: even with
// the oldest ULID and a tied stamp, it never leaves.
func TestOverflowNeverPicksTheJustMovedTicket(t *testing.T) {
	s, ts := tiedStore(t, 3)
	over, err := s.Overflow("done", 2, ts[0].ID)
	if err != nil {
		t.Fatalf("Overflow: %v", err)
	}
	if len(over) != 1 {
		t.Fatalf("overflow = %d tickets, want 1", len(over))
	}
	if over[0].ID == ts[0].ID {
		t.Fatalf("the just-moved ticket was trimmed")
	}
	if over[0].ID != ts[1].ID {
		t.Errorf("overflow picked %s, want the oldest unpinned %s", over[0].Title, ts[1].Title)
	}
}

func TestTrimLaneFilesTheOverflowAndRestoreBringsItBack(t *testing.T) {
	s, ts := trimStore(t, 11)
	trimmed, err := s.TrimLane("done", 10, "bc-20260903", "")
	if err != nil {
		t.Fatalf("TrimLane: %v", err)
	}
	if len(trimmed) != 1 || trimmed[0].ID != ts[0].ID {
		t.Fatalf("trimmed %v, want exactly the oldest %s", trimmed, ts[0].Title)
	}
	if _, err := os.Stat(trimmed[0].Path); err != nil {
		t.Errorf("trimmed file not in the logbook at %q: %v", trimmed[0].Path, err)
	}
	if _, err := os.Stat(ts[0].Path); !os.IsNotExist(err) {
		t.Errorf("oldest ticket still on the board at %q", ts[0].Path)
	}
	left, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(left) != 10 {
		t.Errorf("%d tickets remain on the board, want 10", len(left))
	}
	// The way back exists: restore returns the ticket to the board.
	if _, err := s.Restore(filepath.Base(trimmed[0].Path)); err != nil {
		t.Fatalf("Restore: %v", err)
	}
	left, err = s.List()
	if err != nil {
		t.Fatalf("List after restore: %v", err)
	}
	if len(left) != 11 {
		t.Errorf("%d tickets after restore, want 11", len(left))
	}
}

// A doorway must not jam: what cannot be filed is skipped and named, the rest
// leaves anyway — otherwise one bad ticket blocks every later landing with the
// same error, and the lane never heals (found live by the WXQ9PT review).
func TestFileLaneSkipsWhatItCannotFileAndFilesTheRest(t *testing.T) {
	s, ts := trimStore(t, 3)
	// An unreadable ticket on the board: List reports it, the sweep continues.
	if err := os.WriteFile(filepath.Join(s.TicketsDir(), "broken.md"), []byte("not frontmatter at all"), 0o644); err != nil {
		t.Fatal(err)
	}
	// A name collision for the middle ticket: filed once already today.
	dir := filepath.Join(s.LogbookDir(), "bc-20260904")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, filepath.Base(ts[1].Path)), []byte("occupied"), 0o644); err != nil {
		t.Fatal(err)
	}

	out, err := s.FileLane("done", "bc-20260904", nil)
	if err == nil {
		t.Fatal("a sweep with a broken ticket and a collision reported no problem")
	}
	if len(out) != 2 || out[0].ID != ts[0].ID || out[1].ID != ts[2].ID {
		t.Fatalf("filed %d tickets, want the two fileable ones (t00, t02): %v", len(out), out)
	}
	if !strings.Contains(err.Error(), "broken.md") || !strings.Contains(err.Error(), Handle(ts[1].ID)) {
		t.Errorf("the problems do not name both the unreadable file and the collision: %v", err)
	}
	// The collided ticket is still on the board — skipped, not lost.
	if _, statErr := os.Stat(ts[1].Path); statErr != nil {
		t.Errorf("the skipped ticket left the board: %v", statErr)
	}
}

// refSource hands the store a finished ticket that exists for this board but
// has no file here — what a board on refs sees of somebody else's work.
type refSource struct{ tickets []*Ticket }

func (r refSource) Extra(have map[string]bool) ([]*Ticket, error) {
	var out []*Ticket
	for _, t := range r.tickets {
		if !have[t.ID] {
			out = append(out, t)
		}
	}
	return out, nil
}

// A finished ticket that is only on a ref is not this clone's to file: there is
// no file to move, so trying it only turns every one of them into a problem in
// the cut's warning stream. The board's counter already knew this; the command
// that does the filing has to know it too, or the two disagree about the set.
func TestFileLaneLeavesRefOnlyTicketsWhereTheyAre(t *testing.T) {
	s, ts := trimStore(t, 1)
	stamp := time.Date(2026, 9, 2, 12, 0, 0, 0, time.UTC)
	s.Source = refSource{tickets: []*Ticket{{
		ID:        NewID(stamp),
		Title:     "somebody else's finished work",
		Status:    "done",
		UpdatedAt: stamp,
	}}}

	out, err := s.FileLane("done", "bc-20260904", nil)
	if err != nil {
		t.Fatalf("a ref-only ticket in the lane made the cut report a problem: %v", err)
	}
	if len(out) != 1 || out[0].ID != ts[0].ID {
		t.Fatalf("filed %d tickets, want only the one with a file here: %v", len(out), out)
	}
}

// The holds cap reaches the same tickets the cut does, so it has to make the
// same distinction: a ref-only ticket has no file here to move, and counting
// it toward the cap would push a ticket that does have a file off the board to
// make room for one that was never on it.
func TestOverflowLeavesRefOnlyTicketsOutOfTheCap(t *testing.T) {
	s, ts := trimStore(t, 2)
	stamp := time.Date(2026, 9, 2, 12, 0, 0, 0, time.UTC)
	s.Source = refSource{tickets: []*Ticket{{
		ID:        NewID(stamp),
		Title:     "somebody else's finished work",
		Status:    "done",
		UpdatedAt: stamp,
	}}}

	over, err := s.Overflow("done", 2, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(over) != 0 {
		t.Fatalf("a ref-only ticket counted toward the cap and trimmed %d of ours: %v", len(over), over)
	}

	over, err = s.Overflow("done", 1, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(over) != 1 || over[0].ID != ts[0].ID {
		t.Fatalf("with a cap of one, want only the oldest ticket that has a file here, got %v", over)
	}
}

// Two screens label the same merged list — the lane payload and the signoff
// screen a person accepts work on. The case that matters is the third one: a
// sha only the ticket records must not be labelled as git's account of the
// branch, because that is the label a person reads before signing.
func TestCommitsSourceNamesWhoContributed(t *testing.T) {
	for _, c := range []struct {
		name     string
		derived  []string
		recorded []string
		want     string
	}{
		{"nothing at all", nil, nil, ""},
		{"git accounts for all of it", []string{"a", "b"}, []string{"a"}, "git"},
		{"only the field carries it", nil, []string{"a"}, "ticket"},
		{"the field adds a sha git lost", []string{"a"}, []string{"b"}, "git+ticket"},
	} {
		t.Run(c.name, func(t *testing.T) {
			got := CommitsSource(c.derived, MergeCommits(c.derived, c.recorded))
			if got != c.want {
				t.Fatalf("CommitsSource = %q, want %q", got, c.want)
			}
		})
	}
}
