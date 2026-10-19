package web

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// panicSource panics in Epics, the way a bug in a Source would.
type panicSource struct{ fakeSource }

func (panicSource) Epics() ([]Epic, error) { panic("boom in Epics") }

func TestPanicAnswers500AndIsLogged(t *testing.T) {
	var buf bytes.Buffer
	h := (&Server{Source: &panicSource{}, Log: NewLogger(&buf)}).Handler()

	w := do(t, h, "GET", "/api/epics", "")
	if w.Code != 500 || !strings.Contains(w.Body.String(), "boom in Epics") {
		t.Fatalf("panic answered %d %s, want 500 naming it", w.Code, w.Body)
	}
	got := buf.String()
	if !strings.Contains(got, "GET /api/epics: panic: boom in Epics") || !strings.Contains(got, "goroutine") {
		t.Fatalf("log has no panic with its stack:\n%s", got)
	}
	// The same handler serves on after the panic.
	if w := do(t, h, "GET", "/api/unassigned", ""); w.Code != 200 {
		t.Fatalf("after the panic: %d %s", w.Code, w.Body)
	}
}

func TestErrorAnswersAreLogged(t *testing.T) {
	var buf bytes.Buffer
	h := (&Server{Source: &fakeSource{}, Log: NewLogger(&buf)}).Handler()

	do(t, h, "GET", "/api/tickets/Z", "")
	do(t, h, "GET", "/api/epics", "")
	got := buf.String()
	if !strings.Contains(got, `GET /api/tickets/Z: 404 {"error":"ticket Z: not found"}`) {
		t.Fatalf("404 not logged with its cause:\n%s", got)
	}
	if strings.Contains(got, "/api/epics") {
		t.Fatalf("a 200 was logged:\n%s", got)
	}
}

func TestFailedFetchAndPushAreLogged(t *testing.T) {
	var buf bytes.Buffer
	s := &Server{
		Source: &fakeSource{},
		Fetch:  func() error { return errors.New("remote gone") },
		Push:   func() error { panic("boom in Push") },
		Log:    NewLogger(&buf),
	}
	s.FetchNow()
	if err := s.pushNow(); err == nil || !strings.Contains(err.Error(), "boom in Push") {
		t.Fatalf("a panicking push returned %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "fetch failed: remote gone") || !strings.Contains(got, "push failed: panic: boom in Push") {
		t.Fatalf("log misses the fetch or the push:\n%s", got)
	}
	if !strings.Contains(s.SyncState().PushError, "boom in Push") {
		t.Fatalf("sync record misses the push: %+v", s.SyncState())
	}
}

func TestOpenLogMovesALargeLogAside(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state", "web.log")
	f, err := OpenLog(path)
	if err != nil {
		t.Fatal(err)
	}
	f.Write(bytes.Repeat([]byte("x"), logKeep+1))
	f.Close()

	f, err = OpenLog(path)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	if fi, err := os.Stat(path); err != nil || fi.Size() != 0 {
		t.Fatalf("log not started fresh: %v %v", fi, err)
	}
	if fi, err := os.Stat(path + ".1"); err != nil || fi.Size() != logKeep+1 {
		t.Fatalf("old log not kept: %v %v", fi, err)
	}
}
