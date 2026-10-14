package cli

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/web"
)

func webGet(t *testing.T, h http.Handler, path string, v any) {
	t.Helper()
	r := httptest.NewRequest("GET", path, nil)
	r.Host = "127.0.0.1:1"
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("GET %s: %d %s", path, w.Code, w.Body)
	}
	if err := json.Unmarshal(w.Body.Bytes(), v); err != nil {
		t.Fatal(err)
	}
}

// The page reads the epics the CLI wrote, counts them the way 'jaira epic
// show' does, and a drop on it reaches the other clone on its next fetch —
// written to the epic's ref and pushed, never to the ticket's.
func TestWebShowsEpicsAndADropTravels(t *testing.T) {
	ada, grace := twoBoards(t)
	a := createTicket(t, ada, "Login form")
	createTicket(t, ada, "Fix the form", "--follows", a)
	c := createTicket(t, ada, "Unrelated")
	jsonOf(t, ada, "epic", "create", "magic-login", "--title", "Anmeldung per Link", "--description", "ohne Passwort")
	jsonOf(t, ada, "epic", "add", "magic-login", a)

	g.dir = ada
	s, err := openStore()
	if err != nil {
		t.Fatal(err)
	}
	srv := &web.Server{Source: newWebSource(s), Push: func() error { _, err := refs.Flush(); return err }}
	h := srv.Handler()

	var epics []web.Epic
	webGet(t, h, "/api/epics", &epics)
	if len(epics) != 1 || epics[0].Title != "Anmeldung per Link" || epics[0].Description != "ohne Passwort" || epics[0].Total != 2 {
		t.Fatalf("epics = %+v; want one epic of 2 (member and its follow-up)", epics)
	}
	var loose []web.Ticket
	webGet(t, h, "/api/unassigned", &loose)
	if len(loose) != 1 || loose[0].ID != c {
		t.Fatalf("unassigned = %+v; want only %s", loose, c)
	}

	r := httptest.NewRequest("POST", "/api/epics/magic-login/tickets", strings.NewReader(`{"id":"`+c+`"}`))
	r.Host = "127.0.0.1:1"
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 || strings.Contains(w.Body.String(), "push_error") {
		t.Fatalf("drop: %d %s", w.Code, w.Body)
	}
	webGet(t, h, "/api/unassigned", &loose)
	if len(loose) != 0 {
		t.Fatalf("after the drop unassigned = %+v", loose)
	}

	if out, err := runCLI(t, grace, "fetch"); err != nil {
		t.Fatalf("fetch: %v\n%s", err, out)
	}
	if got := jsonOf(t, grace, "epic", "show", "magic-login"); got["total"].(float64) != 3 {
		t.Fatalf("grace sees %v tickets in the epic; want 3", got["total"])
	}

	r = httptest.NewRequest("POST", "/api/epics/no-such/tickets", strings.NewReader(`{"id":"`+c+`"}`))
	r.Host = "127.0.0.1:1"
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusNotFound {
		t.Fatalf("unknown epic: %d %s", w.Code, w.Body)
	}
}
