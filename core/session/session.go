// Package session records what each agent session is currently doing.
//
// This is scratch state, stored per working tree and never committed: it changes
// every few minutes and would swamp the history of a repository whose whole point
// is readable diffs. The consequence, stated plainly, is that a teammate on their
// own clone cannot see your session — there is no sync channel for it, by design.
package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/berk/jaira/core/ticket"
)

// Session is one agent session's current focus.
type Session struct {
	ID        string `json:"id"`
	TicketID  string `json:"ticket_id,omitempty"`
	Focus     string `json:"focus"`
	Reasoning string `json:"reasoning,omitempty"`
	Model     string `json:"model,omitempty"`
	Host      string `json:"host,omitempty"`
	PID       int    `json:"pid,omitempty"`
	UpdatedAt string `json:"updated_at"`
}

// StaleAfter is how long a session may go without a checkpoint before the board
// dims it. Sessions are never deleted automatically: a crashed session's last
// known focus is useful information, not garbage.
const StaleAfter = 20 * time.Minute

// Path is where a session's state lives.
func Path(s *ticket.Store, id string) string {
	return filepath.Join(s.SessionsDir(), id+".json")
}

// Load reads every session file for this working tree.
func Load(s *ticket.Store) ([]Session, error) {
	matches, err := filepath.Glob(filepath.Join(s.SessionsDir(), "*.json"))
	if err != nil {
		return nil, err
	}
	var out []Session
	for _, m := range matches {
		if filepath.Base(m) == "task-map.json" {
			continue
		}
		b, err := os.ReadFile(m)
		if err != nil {
			continue
		}
		var sess Session
		if err := json.Unmarshal(b, &sess); err != nil {
			continue
		}
		if sess.ID == "" {
			sess.ID = filepath.Base(m)
		}
		out = append(out, sess)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].UpdatedAt > out[j].UpdatedAt })
	return out, nil
}

// Stale reports whether a session has gone quiet.
func (s Session) Stale() bool {
	t, err := time.Parse(time.RFC3339, s.UpdatedAt)
	if err != nil {
		return true
	}
	return time.Since(t) > StaleAfter
}

// Age is how long since the last checkpoint.
func (s Session) Age() time.Duration {
	t, err := time.Parse(time.RFC3339, s.UpdatedAt)
	if err != nil {
		return 0
	}
	return time.Since(t)
}

// Save writes a session record.
func Save(s *ticket.Store, sess Session) error {
	path := Path(s, sess.ID)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(sess, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

// Read loads one session, returning a zero value when absent.
func Read(s *ticket.Store, id string) Session {
	sess := Session{ID: id}
	if b, err := os.ReadFile(Path(s, id)); err == nil {
		_ = json.Unmarshal(b, &sess)
		sess.ID = id
	}
	return sess
}

// Remove deletes a session record.
func Remove(s *ticket.Store, id string) error {
	err := os.Remove(Path(s, id))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// DefaultID prefers an identifier supplied by the harness, falling back to
// something stable for the shell so repeated calls update one record rather than
// littering the directory.
func DefaultID() string {
	for _, k := range []string{"CLAUDE_SESSION_ID", "JAIRA_SESSION_ID", "TERM_SESSION_ID"} {
		if v := os.Getenv(k); v != "" {
			return sanitize(v)
		}
	}
	host, _ := os.Hostname()
	return sanitize(fmt.Sprintf("%s-%d", host, os.Getppid()))
}

func sanitize(s string) string {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_':
			out = append(out, r)
		default:
			out = append(out, '-')
		}
	}
	if len(out) > 48 {
		out = out[:48]
	}
	return string(out)
}
