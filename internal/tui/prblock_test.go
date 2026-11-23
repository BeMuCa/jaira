package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/gate"
	"github.com/BeMuCa/jaira/core/pr"
	"github.com/BeMuCa/jaira/core/ticket"
)

// The board answers pr#<n> from the cache list and next fill, read once per
// reload, and never asks the forge; a ticket waiting only on a pull request has no blocker to open.
func TestBoardReadsPRStateFromTheCacheOnly(t *testing.T) {
	m := newTestModel(t, 150, 32)
	now := time.Now()
	tk, err := m.store.Create(map[string]string{
		ticket.FieldID:       ticket.NewID(now),
		ticket.FieldTitle:    "waits on a pull request",
		ticket.FieldStatus:   "todo",
		ticket.FieldAssignee: "berk",
		ticket.FieldGoal:     "g",
		ticket.FieldDoD:      "d",
		ticket.FieldContext:  "c",
	}, map[string][]string{ticket.FieldBlockedBy: {"pr#12"}}, "")
	if err != nil {
		t.Fatal(err)
	}
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	if gate.Actionable(m.gateEnv(), tk) {
		t.Fatal("with nothing cached the ticket must wait")
	}
	if w := gate.WaitingOnPR(m.gateEnv(), tk); len(w) != 1 || !strings.Contains(w[0], "not asked here") {
		t.Fatalf("WaitingOnPR = %v", w)
	}

	file := filepath.Join(m.store.RepoStateDir(), pr.CacheFile)
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte(`{"12":{"state":"merged","at":"2026-11-22T12:00:00Z"}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if gate.Actionable(m.gateEnv(), tk) {
		t.Fatal("the cache is read once per reload, not per render")
	}
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	if !gate.Actionable(m.gateEnv(), tk) {
		t.Fatal("after a reload a merged pull request in the cache must clear the ticket on the board")
	}

	m.selectByID(tk.ID)
	m.mode = modeDetail
	m.detail = tk
	press(m, "b")
	if m.detail.ID != tk.ID || m.isErr || !strings.Contains(m.message, "pr#12") {
		t.Fatalf("b on a ticket waiting only on a pull request: detail %s, err %v, message %q", ticket.Handle(m.detail.ID), m.isErr, m.message)
	}
}

// The cache follows the board's repository: after a switch to another one,
// a pull request is answered from that repository's cache, never from the
// one the board showed before.
func TestBoardSwitchReadsTheNewRepositorysPRCache(t *testing.T) {
	m := newTestModel(t, 150, 32)
	writeCache := func(root string, body string) {
		t.Helper()
		s, err := ticket.Discover(root)
		if err != nil {
			t.Fatal(err)
		}
		file := filepath.Join(s.RepoStateDir(), pr.CacheFile)
		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	writeCache(m.store.Root, `{"12":{"state":"merged","at":"2026-11-22T12:00:00Z"}}`)
	if w := m.prWait(12); !w.Merged {
		t.Fatalf("first repository: %+v, want merged", w)
	}

	other := t.TempDir()
	s2, err := ticket.At(other)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s2.Init(); err != nil {
		t.Fatal(err)
	}
	m.switchBoard(other)
	if m.store.Root != s2.Root {
		t.Fatalf("switch did not take: %s", m.store.Root)
	}
	if w := m.prWait(12); w.Merged {
		t.Fatal("after the switch PR 12 was answered from the previous repository's cache")
	}
	writeCache(other, `{"12":{"state":"merged","at":"2026-11-22T12:00:00Z"}}`)
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	if w := m.prWait(12); !w.Merged {
		t.Fatalf("second repository's own cache not read: %+v", w)
	}
}
