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

// Source is where the page's data comes from and where a change goes.
//
// Assign only records the change locally (and queues its ref); sending it is
// Push's job, so a remote that is down never makes the drop itself fail.
type Source interface {
	Epics() ([]Epic, error)
	Unassigned() ([]Ticket, error)
	Assign(epic, ticketID string) error
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
	// Now is the clock, replaceable in tests.
	Now func() time.Time

	mu   sync.Mutex
	sync Sync
}

func (s *Server) now() time.Time {
	if s.Now != nil {
		return s.Now()
	}
	return time.Now()
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
	err := s.Fetch()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sync.LastFetch = s.now()
	s.sync.FetchError = errText(err)
}

func (s *Server) pushNow() error {
	if s.Push == nil {
		return nil
	}
	err := s.Push()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sync.LastPush = s.now()
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
	mux.HandleFunc("POST /api/epics/{name}/tickets", s.handleAssign)
	mux.HandleFunc("GET /api/sync", s.handleSync)
	mux.HandleFunc("POST /api/sync", s.handleSyncNow)
	dist, err := fs.Sub(distFS, "ui/dist")
	if err != nil {
		panic(err)
	}
	mux.Handle("GET /", http.FileServerFS(dist))
	return guard(mux)
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

func (s *Server) handleUnassigned(w http.ResponseWriter, _ *http.Request) {
	tickets, err := s.Source.Unassigned()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if tickets == nil {
		tickets = []Ticket{}
	}
	writeJSON(w, http.StatusOK, tickets)
}

// ErrNotFound is what a Source wraps for an epic or ticket that does not
// exist; the handler answers it with 404 rather than 500.
var ErrNotFound = errors.New("not found")

func (s *Server) handleAssign(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ID == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf(`want {"id": "<ticket id>"}`))
		return
	}
	if err := s.Source.Assign(r.PathValue("name"), body.ID); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, ErrNotFound) {
			status = http.StatusNotFound
		}
		writeError(w, status, err)
		return
	}
	// The drop is on disk and queued whatever the push does; a failed push
	// is reported in the sync record and retried by the next one.
	s.pushNow()
	writeJSON(w, http.StatusOK, s.SyncState())
}

func (s *Server) handleSync(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.SyncState())
}

func (s *Server) handleSyncNow(w http.ResponseWriter, _ *http.Request) {
	s.FetchNow()
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
