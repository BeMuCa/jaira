package cli

import (
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

func claimFixture(t *testing.T) (string, *ticket.Store, *ticket.Ticket) {
	t.Helper()
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
		ticket.FieldID:    ticket.NewID(time.Date(2026, 9, 5, 8, 0, 0, 0, time.UTC)),
		ticket.FieldTitle: "t", ticket.FieldStatus: "backlog",
	}, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	return dir, s, tk
}

func TestReleaseDropsYourOwnClaim(t *testing.T) {
	dir, s, tk := claimFixture(t)
	h := ticket.Handle(tk.ID)
	if out, err := runCLI(t, dir, "claim", h, "--session", "me"); err != nil {
		t.Fatalf("claim: %v\n%s", err, out)
	}
	if out, err := runCLI(t, dir, "claim", h, "--release", "--session", "me"); err != nil {
		t.Fatalf("release: %v\n%s", err, out)
	}
	if got, _ := s.Load(tk.ID); got.ClaimedBy != "" {
		t.Errorf("claim still on the ticket: %q", got.ClaimedBy)
	}
}

func TestReleaseClearsAnExpiredForeignClaim(t *testing.T) {
	dir, s, tk := claimFixture(t)
	h := ticket.Handle(tk.ID)
	if _, err := runCLI(t, dir, "claim", h, "--session", "dead-session"); err != nil {
		t.Fatal(err)
	}
	// Age the claim past the TTL — set validates nothing, by design.
	if _, err := runCLI(t, dir, "set", h, "claimed-at=2026-09-01T08:00:00Z"); err != nil {
		t.Fatal(err)
	}
	if out, err := runCLI(t, dir, "claim", h, "--release", "--session", "me"); err != nil {
		t.Fatalf("release of an expired foreign claim: %v\n%s", err, out)
	}
	if got, _ := s.Load(tk.ID); got.ClaimedBy != "" {
		t.Errorf("expired claim still on the ticket: %q", got.ClaimedBy)
	}
}

func TestReleaseRefusesALiveForeignClaim(t *testing.T) {
	dir, s, tk := claimFixture(t)
	h := ticket.Handle(tk.ID)
	if _, err := runCLI(t, dir, "claim", h, "--session", "other-session"); err != nil {
		t.Fatal(err)
	}
	out, err := runCLI(t, dir, "claim", h, "--release", "--session", "me")
	if err == nil {
		t.Fatalf("a live foreign claim was released:\n%s", out)
	}
	if !strings.Contains(out+err.Error(), "still live") {
		t.Errorf("the refusal does not say the claim is live: %v\n%s", err, out)
	}
	if got, _ := s.Load(tk.ID); got.ClaimedBy != "other-session" {
		t.Errorf("the live claim was cleared: %q", got.ClaimedBy)
	}
}
