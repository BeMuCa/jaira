// Package web is the review portal 'jaira web' serves: a page in the browser
// for a person who wants to see and sort the work without an agent and without
// the terminal board.
//
// There is no server in the deployment sense. The binary listens on loopback
// for as long as the command runs, and everything it shows comes from the same
// files and refs the CLI reads; sync stays git, a fetch in the background and a
// push after each change. That is also why this package holds no data of its
// own: a Source hands it the epics and the tickets, and the command that starts
// it decides where they come from.
//
// The page is a Svelte build under ui/dist, committed and embedded, so a
// 'go install' without Node still ships it. CI rebuilds it from ui/src and
// fails when the committed files differ.
package web

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

//go:embed all:ui/dist
var distFS embed.FS

// Ticket is one card on the page.
type Ticket struct {
	ID     string `json:"id"`
	Handle string `json:"handle"`
	Title  string `json:"title"`
	Lane   string `json:"lane"`
	Done   bool   `json:"done"`
	// Direct is false for a ticket that belongs to an epic only because it
	// follows, or is a child of, one of its members.
	Direct bool `json:"direct"`
	// Place is where the ticket lies: board, ref, logbook or archive. A
	// filed ticket is shown, never dragged.
	Place string `json:"place"`
	// The fields the page filters by.
	Goal      string    `json:"goal,omitzero"`
	Tags      []string  `json:"tags,omitempty"`
	Creator   string    `json:"creator,omitzero"`
	Assignee  string    `json:"assignee,omitzero"`
	UpdatedBy string    `json:"updated_by,omitzero"`
	CreatedAt time.Time `json:"created_at,omitzero"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
	// FiledAt is the day it left the board: the logbook folder's date, or
	// for the archive the day it was rejected. Zero when nobody knows.
	FiledAt time.Time `json:"filed_at,omitzero"`
	// Rejected is the kind of no, for a ticket 'jaira reject' archived.
	Rejected string `json:"rejected,omitzero"`
}

// Rejection is why nobody will do a ticket, as 'jaira reject' wrote it.
type Rejection struct {
	Kind string    `json:"kind"`
	Why  string    `json:"why"`
	By   string    `json:"by,omitzero"`
	At   time.Time `json:"at,omitzero"`
}

// Epic is one function of the product with its progress.
type Epic struct {
	Name        string   `json:"name"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Done        int      `json:"done"`
	Total       int      `json:"total"`
	Tickets     []Ticket `json:"tickets"`
}

// Lane is one lane of the board, in the board's order. Kind is what the
// page colours by beside the position: work, human (a person moves it on),
// terminal or parking.
type Lane struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}

// The kinds a Lane has.
const (
	LaneWork     = "work"
	LaneHuman    = "human"
	LaneTerminal = "terminal"
	LaneParking  = "parking"
)

// Check is one item of a ticket's definition of done.
type Check struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
}

// Note is one progress note; Text is markdown as its author wrote it.
type Note struct {
	At     time.Time `json:"at"`
	Author string    `json:"author"`
	Text   string    `json:"text"`
}

// TicketDetail is what the page shows of one ticket when it is opened.
type TicketDetail struct {
	ID      string  `json:"id"`
	Handle  string  `json:"handle"`
	Title   string  `json:"title"`
	Lane    string  `json:"lane"`
	Place   string  `json:"place"`
	Goal    string  `json:"goal"`
	Context string  `json:"context"`
	DoD     []Check `json:"dod"`
	Notes   []Note  `json:"notes"`
	// Rejected is nil for a ticket that was never rejected.
	Rejected *Rejection `json:"rejected,omitempty"`
}

// Source is where the page's data comes from and where a change goes.
//
// The writes only record the change locally (and queue its ref); sending it
// is Push's job, so a remote that is down never makes the change itself fail.
type Source interface {
	Epics() ([]Epic, error)
	// Unassigned is the open tickets in no epic; with all set, every ticket
	// in no epic that ever was — finished, logged and archived ones too.
	Unassigned(all bool) ([]Ticket, error)
	Ticket(id string) (TicketDetail, error)
	Lanes() ([]Lane, error)
	CreateEpic(name, title, description string) error
	// Assign puts a ticket into epic; with from set it takes it out of that
	// epic as well, which is a move.
	Assign(epic, from, ticketID string) error
	Unassign(epic, ticketID string) error
}

// Sync is what the page says about the last round trip to the remote.
type Sync struct {
	LastFetch  time.Time `json:"last_fetch,omitzero"`
	FetchError string    `json:"fetch_error,omitzero"`
	LastPush   time.Time `json:"last_push,omitzero"`
	PushError  string    `json:"push_error,omitzero"`
}

// Server serves the API and the embedded page.
type Server struct {
	Source Source
	// Fetch brings the remote's refs into this clone; Push sends what Assign
	// queued. Either may be nil on a board with no remote.
	Fetch func() error
	Push  func() error
	// Log records every error answer, failed fetch or push, and panic, with
	// its time and cause. nil records nothing.
	Log *log.Logger

	mu   sync.Mutex
	sync Sync
	// pushing keeps two drops in quick succession from flushing one outbox
	// at the same time: the outbox has no lock of its own.
	pushing sync.Mutex
}

// SyncState returns a copy of the sync record.
func (s *Server) SyncState() Sync {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sync
}

// FetchNow runs one fetch and records its outcome. It is what the background
// loop calls every interval.
func (s *Server) FetchNow() {
	if s.Fetch == nil {
		return
	}
	err := recovered(s.Fetch)
	if err != nil {
		s.logf("fetch failed: %v", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sync.LastFetch = time.Now()
	s.sync.FetchError = errText(err)
}

func (s *Server) pushNow() error {
	if s.Push == nil {
		return nil
	}
	s.pushing.Lock()
	defer s.pushing.Unlock()
	err := recovered(s.Push)
	if err != nil {
		s.logf("push failed: %v", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sync.LastPush = time.Now()
	s.sync.PushError = errText(err)
	return err
}

// Loop fetches every interval until stop is closed. The first fetch runs at
// once, so the page opened right after the command shows today's state.
func (s *Server) Loop(every time.Duration, stop <-chan struct{}) {
	s.FetchNow()
	t := time.NewTicker(every)
	defer t.Stop()
	for {
		select {
		case <-stop:
			return
		case <-t.C:
			s.FetchNow()
		}
	}
}

func errText(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// Handler returns the routes: the JSON API under /api/ and the page for
// everything else.
//
// Every request must name a loopback host. Binding to 127.0.0.1 keeps other
// machines out, but not a web page in the same browser: through DNS
// rebinding a foreign site can reach a port on loopback under its own host
// name. Checking Host refuses that, and checking Origin on writes refuses a
// cross-site form post.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/epics", s.handleEpics)
	mux.HandleFunc("GET /api/unassigned", s.handleUnassigned)
	mux.HandleFunc("GET /api/tickets/{id}", s.handleTicket)
	mux.HandleFunc("GET /api/lanes", s.handleLanes)
	mux.HandleFunc("POST /api/epics", s.handleCreateEpic)
	mux.HandleFunc("POST /api/epics/{name}/tickets", s.handleAssign)
	mux.HandleFunc("DELETE /api/epics/{name}/tickets/{id}", s.handleUnassign)
	mux.HandleFunc("GET /api/sync", s.handleSync)
	dist, err := fs.Sub(distFS, "ui/dist")
	if err != nil {
		panic(err)
	}
	mux.Handle("GET /", http.FileServerFS(dist))
	return s.logged(guard(mux))
}

func guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !loopbackHost(r.Host) {
			http.Error(w, "jaira web answers only on 127.0.0.1", http.StatusForbidden)
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			if o := r.Header.Get("Origin"); o != "" && !loopbackHost(strings.TrimPrefix(strings.TrimPrefix(o, "http://"), "https://")) {
				http.Error(w, "cross-origin write refused", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func loopbackHost(hostport string) bool {
	host, _, err := net.SplitHostPort(hostport)
	if err != nil {
		host = hostport
	}
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(strings.Trim(host, "[]"))
	return ip != nil && ip.IsLoopback()
}

func (s *Server) handleEpics(w http.ResponseWriter, _ *http.Request) {
	epics, err := s.Source.Epics()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if epics == nil {
		epics = []Epic{}
	}
	for i := range epics {
		if epics[i].Tickets == nil {
			epics[i].Tickets = []Ticket{}
		}
	}
	writeJSON(w, http.StatusOK, epics)
}

func (s *Server) handleUnassigned(w http.ResponseWriter, r *http.Request) {
	tickets, err := s.Source.Unassigned(r.URL.Query().Get("all") == "1")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if tickets == nil {
		tickets = []Ticket{}
	}
	writeJSON(w, http.StatusOK, tickets)
}

// The errors a Source wraps so the handler answers with what went wrong
// rather than 500: no such epic or ticket, a name already taken or not
// usable, and a change that does not apply (a follow-up has no member line
// to take out).
var (
	ErrNotFound = errors.New("not found")
	ErrExists   = errors.New("exists already")
	ErrInvalid  = errors.New("invalid")
	ErrConflict = errors.New("does not apply")
)

func statusOf(err error) int {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrExists), errors.Is(err, ErrConflict):
		return http.StatusConflict
	case errors.Is(err, ErrInvalid):
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

func (s *Server) handleTicket(w http.ResponseWriter, r *http.Request) {
	t, err := s.Source.Ticket(r.PathValue("id"))
	if err != nil {
		writeError(w, statusOf(err), err)
		return
	}
	if t.DoD == nil {
		t.DoD = []Check{}
	}
	if t.Notes == nil {
		t.Notes = []Note{}
	}
	writeJSON(w, http.StatusOK, t)
}

func (s *Server) handleLanes(w http.ResponseWriter, _ *http.Request) {
	lanes, err := s.Source.Lanes()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if lanes == nil {
		lanes = []Lane{}
	}
	writeJSON(w, http.StatusOK, lanes)
}

func (s *Server) handleCreateEpic(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name        string `json:"name"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.Name) == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf(`want {"name": "<epic name>", "title": ..., "description": ...}`))
		return
	}
	s.write(w, s.Source.CreateEpic(body.Name, body.Title, body.Description))
}

func (s *Server) handleAssign(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID   string `json:"id"`
		From string `json:"from"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ID == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf(`want {"id": "<ticket id>", "from": "<epic it leaves, if any>"}`))
		return
	}
	s.write(w, s.Source.Assign(r.PathValue("name"), body.From, body.ID))
}

func (s *Server) handleUnassign(w http.ResponseWriter, r *http.Request) {
	s.write(w, s.Source.Unassign(r.PathValue("name"), r.PathValue("id")))
}

// write answers a change. Whatever part of it was recorded is pushed even
// when another part failed — a move whose second half failed has still
// written its first — and a failed push is reported in the sync record and
// retried by the next one.
func (s *Server) write(w http.ResponseWriter, err error) {
	s.pushNow()
	if err != nil {
		writeError(w, statusOf(err), err)
		return
	}
	writeJSON(w, http.StatusOK, s.SyncState())
}

func (s *Server) handleSync(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.SyncState())
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

// Listen opens a listener on loopback. port 0 picks a free one.
func Listen(port int) (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
}
