package lane

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// Settle is the one place that decides which file rule a lane runs — doorway
// first, cap second, neither for a plain lane or a missing one.
func TestSettleRoutesDoorwayThenCap(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	base := time.Date(2026, 9, 5, 8, 0, 0, 0, time.UTC)
	var ids []string
	for i := 0; i < 3; i++ {
		stamp := base.Add(time.Duration(i) * time.Minute)
		tk, err := s.Create(map[string]string{
			ticket.FieldID: ticket.NewID(stamp), ticket.FieldTitle: "t",
			ticket.FieldStatus: "closed", ticket.FieldCreatedAt: ticket.FormatTime(stamp),
			ticket.FieldUpdatedAt: ticket.FormatTime(stamp),
		}, nil, "")
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, tk.ID)
	}

	// A nil lane settles nothing.
	if out, filed, err := Settle(s, nil, "x-20260905", "", nil); out != nil || filed || err != nil {
		t.Errorf("Settle(nil) = %v, %v, %v; want silence", out, filed, err)
	}
	// A plain lane settles nothing.
	if out, filed, err := Settle(s, &Lane{ID: "closed"}, "x-20260905", "", nil); out != nil || filed || err != nil {
		t.Errorf("Settle(plain) = %v, %v, %v; want silence", out, filed, err)
	}
	// A cap trims the overflow, pinning the just-moved ticket.
	capped := &Lane{ID: "closed", Holds: 2}
	out, filed, err := Settle(s, capped, "x-20260905", ids[0], nil)
	if err != nil || filed {
		t.Fatalf("cap settle: %v, filed=%v", err, filed)
	}
	if len(out) != 1 || out[0].ID == ids[0] {
		t.Fatalf("cap trimmed %v, want one ticket and never the pinned %s", out, ids[0])
	}
	// A doorway files everything that is left, and wins over a cap.
	door := &Lane{ID: "closed", LogbookOnEntry: true, Holds: 2}
	out, filed, err = Settle(s, door, "x-20260905", "", nil)
	if err != nil || !filed {
		t.Fatalf("doorway settle: %v, filed=%v", err, filed)
	}
	if len(out) != 2 {
		t.Fatalf("doorway filed %d tickets, want the remaining 2", len(out))
	}
}
