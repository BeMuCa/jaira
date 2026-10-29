package cli

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
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

// fileAway moves a board ticket's file to dir under .jaira, as filing it
// does, with extra frontmatter lines put in before the closing delimiter.
func fileAway(t *testing.T, s *ticket.Store, id, dir, extra string) {
	t.Helper()
	var from string
	entries, err := os.ReadDir(s.TicketsDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if ticket.IDFromFilename(e.Name()) == id {
			from = filepath.Join(s.TicketsDir(), e.Name())
		}
	}
	raw, err := os.ReadFile(from)
	if err != nil {
		t.Fatal(err)
	}
	src := strings.Replace(string(raw), "\n---\n", "\n"+extra+"---\n", 1)
	to := filepath.Join(s.Root, ticket.DirName, dir, filepath.Base(from))
	if err := os.MkdirAll(filepath.Dir(to), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(to, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(from); err != nil {
		t.Fatal(err)
	}
}

// all=1 lists every ticket that ever was: the logbook and the archive beside
// the board, each saying where it lies, the logbook's with the day it was
// filed and a rejected one with its kind; each of them opens, and the opened
// archive ticket says why it was rejected.
func TestWebListsFiledTicketsAndOpensThem(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("JAIRA_USER", "Ada Lovelace")
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))
	if out, err := runCLI(t, dir, "init"); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	open := createTicket(t, dir, "Still open", "--tag", "web")
	logged := createTicket(t, dir, "Long done")
	rejected := createTicket(t, dir, "Never")

	g.dir = dir
	s, err := openStore()
	if err != nil {
		t.Fatal(err)
	}
	fileAway(t, s, logged, filepath.Join(ticket.LogbookSubdir, "al-20261001"), "")
	fileAway(t, s, rejected, ticket.ArchiveSubdir,
		"rejected-kind: wontfix\nrejected-why: nobody needs it\nrejected-by: Ada Lovelace\nrejected-at: 2026-10-02T09:00:00Z\n")
	h := (&web.Server{Source: newWebSource(s)}).Handler()

	var loose []web.Ticket
	webGet(t, h, "/api/unassigned", &loose)
	if len(loose) != 1 || loose[0].ID != open || loose[0].Place != "board" || loose[0].Creator != "Ada Lovelace" || loose[0].Goal != "g" || len(loose[0].Tags) != 1 {
		t.Fatalf("unassigned = %+v; want only %s, on the board, with its creator and tag", loose, open)
	}

	webGet(t, h, "/api/unassigned?all=1", &loose)
	got := map[string]web.Ticket{}
	for _, c := range loose {
		got[c.ID] = c
	}
	if len(loose) != 3 {
		t.Fatalf("all = %+v; want 3 tickets", loose)
	}
	if c := got[logged]; c.Place != "logbook" || c.FiledAt.Format("2006-01-02") != "2026-10-01" {
		t.Errorf("logged = %+v; want place logbook, filed 2026-10-01", c)
	}
	if c := got[rejected]; c.Place != "archive" || c.Rejected != "wontfix" || c.FiledAt.Format("2006-01-02") != "2026-10-02" {
		t.Errorf("rejected = %+v; want place archive, rejected wontfix, filed 2026-10-02", c)
	}

	var d web.TicketDetail
	webGet(t, h, "/api/tickets/"+logged, &d)
	if d.Place != "logbook" || d.Rejected != nil {
		t.Errorf("logged detail = %+v", d)
	}
	webGet(t, h, "/api/tickets/"+rejected, &d)
	if d.Place != "archive" || d.Rejected == nil || d.Rejected.Kind != "wontfix" || d.Rejected.Why != "nobody needs it" {
		t.Errorf("rejected detail = %+v; want its rejection", d)
	}
	webGet(t, h, "/api/tickets/"+open, &d)
	if d.Place != "board" {
		t.Errorf("open detail place = %q", d.Place)
	}
}

// The page colours lanes by the board's order and their kind, so the lanes
// come in that order with the kind of each.
func TestWebListsLanesInBoardOrderWithKind(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))
	if out, err := runCLI(t, dir, "init"); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	g.dir = dir
	s, err := openStore()
	if err != nil {
		t.Fatal(err)
	}
	var lanes []web.Lane
	webGet(t, (&web.Server{Source: newWebSource(s)}).Handler(), "/api/lanes", &lanes)
	kind := map[string]string{}
	var order []string
	for _, l := range lanes {
		kind[l.ID] = l.Kind
		order = append(order, l.ID)
	}
	if len(order) == 0 || order[0] != "backlog" {
		t.Fatalf("lanes = %v; want the board's order, backlog first", order)
	}
	for id, want := range map[string]string{"backlog": "work", "human": "human", "signoff": "human", "done": "terminal", "blocked": "parking"} {
		if kind[id] != want {
			t.Errorf("%s is %q, want %q", id, kind[id], want)
		}
	}
}

// The page rewrites an epic through the same call 'jaira epic set' makes and
// reads the new state straight back. An editor that sat open while the CLI
// changed the same epic is refused with 409 and writes nothing, so neither
// change is lost; the members are untouched by either.
func TestWebRewritesAnEpicAndRefusesAStaleEditor(t *testing.T) {
	ada, grace := twoBoards(t)
	a := createTicket(t, ada, "Login form")
	jsonOf(t, ada, "epic", "create", "magic-login", a, "--title", "Login")

	g.dir = ada
	s, err := openStore()
	if err != nil {
		t.Fatal(err)
	}
	h := (&web.Server{Source: newWebSource(s), Push: webPush}).Handler()
	var epics []web.Epic
	webGet(t, h, "/api/epics", &epics)
	read, _ := json.Marshal(web.EpicFields{Title: epics[0].Title, Description: epics[0].Description, Body: epics[0].Body})

	edit := `{"title":"Anmeldung","description":"per Link","body":"# Login\n\n## Ziel\n\nkein Passwort","base":` + string(read) + `}`
	if w := webDo(t, h, "PATCH", "/api/epics/magic-login", edit); w.Code != 200 {
		t.Fatalf("edit: %d %s", w.Code, w.Body)
	}
	webGet(t, h, "/api/epics", &epics)
	if e := epics[0]; e.Title != "Anmeldung" || e.Description != "per Link" || e.Body != "# Anmeldung\n\n## Ziel\n\nkein Passwort" || e.Text != "## Ziel\n\nkein Passwort" || len(e.Tickets) != 1 {
		t.Fatalf("after the edit: %+v", e)
	}

	if out, err := runAndSend(t, ada, "epic", "set", "magic-login", "--body", "from the CLI"); err != nil {
		t.Fatalf("set: %v\n%s", err, out)
	}
	stale := `{"body":"from a stale editor","base":` + string(read) + `}`
	if w := webDo(t, h, "PATCH", "/api/epics/magic-login", stale); w.Code != http.StatusConflict || !strings.Contains(w.Body.String(), "open it again") {
		t.Fatalf("stale editor: %d %s; want 409 saying to open it again", w.Code, w.Body)
	}
	if w := webDo(t, h, "PATCH", "/api/epics/no-such", `{"title":"x"}`); w.Code != http.StatusNotFound {
		t.Fatalf("unknown epic: %d %s", w.Code, w.Body)
	}
	if w := webDo(t, h, "PATCH", "/api/epics/magic-login", `{"title":"  "}`); w.Code != http.StatusBadRequest {
		t.Fatalf("empty title: %d %s", w.Code, w.Body)
	}

	if out, err := runCLI(t, grace, "fetch"); err != nil {
		t.Fatalf("fetch: %v\n%s", err, out)
	}
	got := jsonOf(t, grace, "epic", "show", "magic-login")
	if got["body"] != "from the CLI" || got["title"] != "Anmeldung" || len(got["members"].([]any)) != 1 {
		t.Fatalf("grace sees %v", got)
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
		if loose[0].BehindRef != d.BehindRef {
			t.Errorf("behind_ref: card %v, detail %v; want the same", loose[0].BehindRef, d.BehindRef)
		}
		return loose[0].Lane, d.Lane, loose[0]
	}

	fileSays(t, tk.Path, "in-progress", "2026-09-10T08:00:00Z")
	if c, d, card := lanes(); c != "signoff" || d != "signoff" || !card.BehindRef {
		t.Errorf("file older than its ref: card %q, detail %q, behind_ref %v; want signoff from the ref, marked", c, d, card.BehindRef)
	}

	fileSays(t, tk.Path, "in-progress", "2099-01-01T00:00:00Z")
	c, d, card := lanes()
	if c != "in-progress" || d != "in-progress" || card.BehindRef {
		t.Errorf("file newer than its ref: card %q, detail %q, behind_ref %v; want in-progress from the file, unmarked", c, d, card.BehindRef)
	}
	if card.CreatedAt.IsZero() || !card.CreatedAt.Equal(tk.CreatedAt) {
		t.Errorf("card created_at = %v; want %v", card.CreatedAt, tk.CreatedAt)
	}
	if want := time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC); !card.UpdatedAt.Equal(want) {
		t.Errorf("card updated_at = %v; want %v", card.UpdatedAt, want)
	}
}

// An epic counts the tickets waiting in signoff apart from the ones done: a
// person still has to accept them.
func TestWebCountsSignoffApartFromDone(t *testing.T) {
	ada, _ := twoBoards(t)
	a := createTicket(t, ada, "Waiting for acceptance")
	b := createTicket(t, ada, "Not started")
	jsonOf(t, ada, "epic", "create", "checkout")
	jsonOf(t, ada, "epic", "add", "checkout", a, b)
	if out, err := runAndSend(t, ada, "pull", a); err != nil {
		t.Fatalf("pull: %v\n%s", err, out)
	}
	if out, err := runAndSend(t, ada, "move", a, "--to", "signoff", "--force"); err != nil {
		t.Fatalf("move: %v\n%s", err, out)
	}

	g.dir = ada
	s, err := openStore()
	if err != nil {
		t.Fatal(err)
	}
	var epics []web.Epic
	webGet(t, (&web.Server{Source: newWebSource(s)}).Handler(), "/api/epics", &epics)
	if len(epics) != 1 || epics[0].Signoff != 1 || epics[0].Done != 0 || epics[0].Total != 2 {
		t.Fatalf("epics = %+v; want 1 in signoff, 0 done, of 2", epics)
	}
}
