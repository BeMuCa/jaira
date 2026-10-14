package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type fakeSource struct {
	epics      []Epic
	unassigned []Ticket
	assigned   []string
	err        error
}

func (f *fakeSource) Epics() ([]Epic, error)        { return f.epics, f.err }
func (f *fakeSource) Unassigned() ([]Ticket, error) { return f.unassigned, f.err }
func (f *fakeSource) Assign(epic, id string) error {
	if f.err != nil {
		return f.err
	}
	if epic != "login" {
		return fmt.Errorf("epic %q: %w", epic, ErrNotFound)
	}
	f.assigned = append(f.assigned, epic+"/"+id)
	return nil
}

func do(t *testing.T, h http.Handler, method, target, body string, hdr ...string) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(method, target, strings.NewReader(body))
	r.Host = "127.0.0.1:7420"
	for i := 0; i+1 < len(hdr); i += 2 {
		r.Header.Set(hdr[i], hdr[i+1])
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestEpicsAndUnassigned(t *testing.T) {
	src := &fakeSource{
		epics: []Epic{{Name: "login", Title: "Login", Description: "Sign in.", Done: 1, Total: 2,
			Tickets: []Ticket{{ID: "A", Handle: "AAAAAA", Done: true, Direct: true}}}},
		unassigned: []Ticket{{ID: "B", Handle: "BBBBBB", Title: "loose", Lane: "todo"}},
	}
	h := (&Server{Source: src}).Handler()

	w := do(t, h, "GET", "/api/epics", "")
	if w.Code != 200 {
		t.Fatalf("epics: %d %s", w.Code, w.Body)
	}
	var epics []Epic
	if err := json.Unmarshal(w.Body.Bytes(), &epics); err != nil {
		t.Fatal(err)
	}
	if len(epics) != 1 || epics[0].Done != 1 || epics[0].Total != 2 || epics[0].Description != "Sign in." {
		t.Fatalf("epics = %+v", epics)
	}

	w = do(t, h, "GET", "/api/unassigned", "")
	var loose []Ticket
	if err := json.Unmarshal(w.Body.Bytes(), &loose); err != nil {
		t.Fatal(err)
	}
	if len(loose) != 1 || loose[0].Handle != "BBBBBB" {
		t.Fatalf("unassigned = %+v", loose)
	}
}

// An empty board answers [] rather than null, so the page never has to tell
// "nothing" from "broken".
func TestEmptyListsAreArrays(t *testing.T) {
	h := (&Server{Source: &fakeSource{epics: []Epic{{Name: "x"}}}}).Handler()
	if got := strings.TrimSpace(do(t, h, "GET", "/api/unassigned", "").Body.String()); got != "[]" {
		t.Fatalf("unassigned = %s", got)
	}
	if got := do(t, h, "GET", "/api/epics", "").Body.String(); !strings.Contains(got, `"tickets":[]`) {
		t.Fatalf("epics = %s", got)
	}
}

func TestAssignWritesThenPushes(t *testing.T) {
	src := &fakeSource{}
	var order []string
	pushed := errors.New("remote down")
	s := &Server{Source: src, Push: func() error {
		order = append(order, fmt.Sprint(len(src.assigned)))
		return pushed
	}}
	h := s.Handler()

	w := do(t, h, "POST", "/api/epics/login/tickets", `{"id":"B"}`)
	if w.Code != 200 {
		t.Fatalf("assign: %d %s", w.Code, w.Body)
	}
	if len(src.assigned) != 1 || src.assigned[0] != "login/B" {
		t.Fatalf("assigned = %v", src.assigned)
	}
	// The push ran after the write, and its failure is reported, not fatal.
	if len(order) != 1 || order[0] != "1" {
		t.Fatalf("push order = %v", order)
	}
	var st Sync
	if err := json.Unmarshal(w.Body.Bytes(), &st); err != nil {
		t.Fatal(err)
	}
	if st.PushError != "remote down" || st.LastPush.IsZero() {
		t.Fatalf("sync = %+v", st)
	}
}

func TestAssignErrors(t *testing.T) {
	h := (&Server{Source: &fakeSource{}}).Handler()
	if w := do(t, h, "POST", "/api/epics/nope/tickets", `{"id":"B"}`); w.Code != 404 {
		t.Fatalf("unknown epic: %d", w.Code)
	}
	if w := do(t, h, "POST", "/api/epics/login/tickets", `{}`); w.Code != 400 {
		t.Fatalf("no id: %d", w.Code)
	}
}

// A foreign site can reach a loopback port by rebinding its own name to
// 127.0.0.1, or post a form to it; both carry a host that is not loopback.
func TestForeignHostAndOriginRefused(t *testing.T) {
	src := &fakeSource{}
	h := (&Server{Source: src}).Handler()

	r := httptest.NewRequest("GET", "/api/epics", nil)
	r.Host = "evil.example:7420"
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusForbidden {
		t.Fatalf("rebinding host: %d", w.Code)
	}

	w = do(t, h, "POST", "/api/epics/login/tickets", `{"id":"B"}`, "Origin", "https://evil.example")
	if w.Code != http.StatusForbidden || len(src.assigned) != 0 {
		t.Fatalf("cross-origin write: %d, assigned %v", w.Code, src.assigned)
	}
	w = do(t, h, "POST", "/api/epics/login/tickets", `{"id":"B"}`, "Origin", "http://127.0.0.1:7420")
	if w.Code != 200 {
		t.Fatalf("same-origin write: %d %s", w.Code, w.Body)
	}
}

func TestSyncRecordsFetch(t *testing.T) {
	calls := 0
	s := &Server{Source: &fakeSource{}, Fetch: func() error {
		calls++
		if calls == 2 {
			return errors.New("offline")
		}
		return nil
	}}
	h := s.Handler()
	s.FetchNow()
	if st := s.SyncState(); st.LastFetch.IsZero() || st.FetchError != "" {
		t.Fatalf("first fetch: %+v", st)
	}
	s.FetchNow()
	if st := s.SyncState(); st.FetchError != "offline" {
		t.Fatalf("second fetch: %+v", st)
	}
	if got := do(t, h, "GET", "/api/sync", "").Body.String(); !strings.Contains(got, "offline") {
		t.Fatalf("GET /api/sync = %s", got)
	}
}

func TestLoopFetchesAtOnceAndStops(t *testing.T) {
	fetched := make(chan struct{}, 10)
	s := &Server{Source: &fakeSource{}, Fetch: func() error { fetched <- struct{}{}; return nil }}
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() { s.Loop(time.Hour, stop); close(done) }()
	select {
	case <-fetched:
	case <-time.After(5 * time.Second):
		t.Fatal("no fetch at start")
	}
	close(stop)
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("loop did not stop")
	}
}

// The page itself is embedded: 'go install' has no Node and must still ship it.
func TestServesEmbeddedPage(t *testing.T) {
	h := (&Server{Source: &fakeSource{}}).Handler()
	w := do(t, h, "GET", "/", "")
	if w.Code != 200 || !strings.Contains(w.Body.String(), `<div id="app">`) {
		t.Fatalf("index: %d %.200s", w.Code, w.Body)
	}
}

func TestListenIsLoopback(t *testing.T) {
	l, err := Listen(0)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	if !strings.HasPrefix(l.Addr().String(), "127.0.0.1:") {
		t.Fatalf("listening on %s", l.Addr())
	}
}
