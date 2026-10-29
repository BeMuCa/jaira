package cli

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

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

// A drop that does not reach the remote has to say so: the outbox keeps it,
// Flush reports nothing wrong, and without webPush the page would show the
// change as done while nobody else can see it.
func TestWebPushReportsAQueuedDrop(t *testing.T) {
	ada, _ := twoBoards(t)
	c := createTicket(t, ada, "Loose")
	jsonOf(t, ada, "epic", "create", "auth")
	// Port 1 on loopback refuses at once: unreachable, without DNS or a timeout.
	gitRun(t, ada, "remote", "set-url", "origin", "http://127.0.0.1:1/board.git")

	g.dir = ada
	s, err := openStore()
	if err != nil {
		t.Fatal(err)
	}
	h := (&web.Server{Source: newWebSource(s), Push: webPush}).Handler()
	r := httptest.NewRequest("POST", "/api/epics/auth/tickets", strings.NewReader(`{"id":"`+c+`"}`))
	r.Host = "127.0.0.1:1"
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	var st web.Sync
	if err := json.Unmarshal(w.Body.Bytes(), &st); err != nil {
		t.Fatal(err)
	}
	if w.Code != 200 || st.PushError == "" {
		t.Fatalf("drop with the remote down: %d %s; want 200 with a push_error", w.Code, w.Body)
	}
}

func webDo(t *testing.T, h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	r.Host = "127.0.0.1:1"
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

// The page creates an epic, moves a ticket from one epic to another, takes
// it out again, and opens a ticket with its notes — each through the calls
// the CLI makes, so the other clone sees the result on its next fetch.
func TestWebCreatesMovesRemovesAndOpens(t *testing.T) {
	ada, grace := twoBoards(t)
	a := createTicket(t, ada, "Login form")
	f := createTicket(t, ada, "Fix the form", "--follows", a)
	jsonOf(t, ada, "epic", "create", "magic-login")
	jsonOf(t, ada, "epic", "add", "magic-login", a)
	if out, err := runAndSend(t, ada, "pull", a); err != nil {
		t.Fatalf("pull: %v\n%s", err, out)
	}
	if out, err := runAndSend(t, ada, "note", a, "**bold** and <script>alert(1)</script>"); err != nil {
		t.Fatalf("note: %v\n%s", err, out)
	}

	g.dir = ada
	s, err := openStore()
	if err != nil {
		t.Fatal(err)
	}
	h := (&web.Server{Source: newWebSource(s), Push: webPush}).Handler()

	if w := webDo(t, h, "POST", "/api/epics", `{"name":"Pay Later","title":"Später zahlen","description":"Rechnung"}`); w.Code != 200 {
		t.Fatalf("create: %d %s", w.Code, w.Body)
	}
	if w := webDo(t, h, "POST", "/api/epics", `{"name":"pay-later"}`); w.Code != http.StatusConflict {
		t.Fatalf("create twice: %d %s", w.Code, w.Body)
	}
	if w := webDo(t, h, "POST", "/api/epics/pay-later/tickets", `{"id":"`+a+`","from":"magic-login"}`); w.Code != 200 || strings.Contains(w.Body.String(), "push_error") {
		t.Fatalf("move: %d %s", w.Code, w.Body)
	}
	if w := webDo(t, h, "DELETE", "/api/epics/pay-later/tickets/"+f, ""); w.Code != http.StatusConflict || !strings.Contains(w.Body.String(), "follows") {
		t.Fatalf("removing a follow-up: %d %s; want 409 naming follows", w.Code, w.Body)
	}
	if w := webDo(t, h, "POST", "/api/epics/no-such/tickets", `{"id":"`+a+`","from":"pay-later"}`); w.Code != http.StatusNotFound || !strings.Contains(w.Body.String(), "adding") {
		t.Fatalf("move to a missing epic: %d %s; want 404 naming the add", w.Code, w.Body)
	}

	if out, err := runCLI(t, grace, "fetch"); err != nil {
		t.Fatalf("fetch: %v\n%s", err, out)
	}
	if got := jsonOf(t, grace, "epic", "show", "pay-later"); got["total"].(float64) != 2 || got["title"] != "Später zahlen" {
		t.Fatalf("grace sees pay-later as %v", got)
	}
	if got := jsonOf(t, grace, "epic", "show", "magic-login"); got["total"].(float64) != 0 {
		t.Fatalf("magic-login still holds %v after the move", got["total"])
	}

	if w := webDo(t, h, "DELETE", "/api/epics/pay-later/tickets/"+a, ""); w.Code != 200 {
		t.Fatalf("remove: %d %s", w.Code, w.Body)
	}
	var loose []web.Ticket
	webGet(t, h, "/api/unassigned", &loose)
	if len(loose) != 2 {
		t.Fatalf("after the remove unassigned = %+v; want the ticket and its follow-up", loose)
	}

	var d web.TicketDetail
	webGet(t, h, "/api/tickets/"+a, &d)
	if d.Title != "Login form" || d.Goal != "g" || len(d.DoD) != 1 || len(d.Notes) == 0 {
		t.Fatalf("detail = %+v", d)
	}
	if last := d.Notes[len(d.Notes)-1]; last.Text != "**bold** and <script>alert(1)</script>" || last.Author == "" || last.At.IsZero() {
		t.Fatalf("last note = %+v; want the text verbatim with its author and time", last)
	}
}

// fileSays rewrites the ticket file's lane and updated-at by hand, the state a
// checkout is left in when later lanes wrote only the ref.
func fileSays(t *testing.T, path, lane, updated string) {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var lines []string
	for _, line := range strings.Split(string(b), "\n") {
		switch {
		case strings.HasPrefix(line, "status:"):
			line = "status: " + lane
		case strings.HasPrefix(line, "updated-at:"):
			line = "updated-at: " + updated
		}
		lines = append(lines, line)
	}
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
		t.Fatal(err)
	}
}

// A file older than its ref shows the ref's lane, on the card and when the
// ticket is opened; a file newer than its ref shows its own. Every card
// carries the ticket's created and updated stamps.
func TestWebShowsTheNewerOfFileAndRef(t *testing.T) {
	t.Setenv("JAIRA_USER", "ada")
	clone := cloneWithRemotes(t, "origin")
	h := ticketIn(t, clone, "stale card")
	if out, err := runCLI(t, clone, "pull", h); err != nil {
		t.Fatalf("pull: %v\n%s", err, out)
	}
	if out, err := runCLI(t, clone, "move", h, "--to", "signoff", "--force"); err != nil {
		t.Fatalf("move: %v\n%s", err, out)
	}
	flushRefs()

	g.dir = clone
	s, err := openStore()
	if err != nil {
		t.Fatal(err)
	}
	tk, err := s.Load(h)
	if err != nil {
		t.Fatal(err)
	}
	h2 := (&web.Server{Source: newWebSource(s)}).Handler()
	lanes := func() (string, string, web.Ticket) {
		var loose []web.Ticket
		webGet(t, h2, "/api/unassigned", &loose)
		if len(loose) != 1 {
			t.Fatalf("unassigned = %+v; want the one ticket", loose)
		}
		var d web.TicketDetail
		webGet(t, h2, "/api/tickets/"+tk.ID, &d)
		return loose[0].Lane, d.Lane, loose[0]
	}

	fileSays(t, tk.Path, "in-progress", "2026-09-10T08:00:00Z")
	if c, d, _ := lanes(); c != "signoff" || d != "signoff" {
		t.Errorf("file older than its ref: card %q, detail %q; want signoff from the ref", c, d)
	}

	fileSays(t, tk.Path, "in-progress", "2099-01-01T00:00:00Z")
	c, d, card := lanes()
	if c != "in-progress" || d != "in-progress" {
		t.Errorf("file newer than its ref: card %q, detail %q; want in-progress from the file", c, d)
	}
	if card.CreatedAt.IsZero() || !card.CreatedAt.Equal(tk.CreatedAt) {
		t.Errorf("card created_at = %v; want %v", card.CreatedAt, tk.CreatedAt)
	}
	if want := time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC); !card.UpdatedAt.Equal(want) {
		t.Errorf("card updated_at = %v; want %v", card.UpdatedAt, want)
	}
}
