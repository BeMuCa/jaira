package web

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
)

// logKeep is how large the log may grow before the next start moves it aside
// to <path>.1. One old file is kept: enough to read what happened before a
// restart, and a page left open for weeks never fills the disk.
const logKeep = 1 << 20

// OpenLog opens the log 'jaira web' writes at path, for appending.
//
// The page runs where nobody watches a terminal, so stderr is not a record: a
// tab closed, a session detached, and what the process said on the way down
// is gone. The file is what is left to read afterwards.
func OpenLog(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	if fi, err := os.Stat(path); err == nil && fi.Size() > logKeep {
		_ = os.Rename(path, path+".1")
	}
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
}

// NewLogger writes to w with the date and time on every line.
func NewLogger(w io.Writer) *log.Logger {
	return log.New(w, "", log.LstdFlags|log.Lmicroseconds)
}

func (s *Server) logf(format string, args ...any) {
	if s.Log != nil {
		s.Log.Printf(format, args...)
	}
}

// recovered turns a panic in f into an error carrying its stack, so a fetch or
// a push that panics is recorded and reported like one that failed, and the
// loop that called it keeps running.
func recovered(f func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v\n%s", r, debug.Stack())
		}
	}()
	return f()
}

// logged records every answer with an error status, and every panic, then
// answers the panic with 500 instead of dropping the connection.
//
// net/http would survive the panic on its own, but it writes the stack to
// stderr and closes the connection without a status, so the page sees a
// network error and the cause is nowhere.
func (s *Server) logged(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &recorder{ResponseWriter: w}
		defer func() {
			p := recover()
			if p == nil {
				if rec.status >= 400 {
					s.logf("%s %s: %d %s", r.Method, r.URL.Path, rec.status, bytes.TrimSpace(rec.body.Bytes()))
				}
				return
			}
			if p == http.ErrAbortHandler {
				panic(p)
			}
			s.logf("%s %s: panic: %v\n%s", r.Method, r.URL.Path, p, debug.Stack())
			if rec.status == 0 {
				writeError(w, http.StatusInternalServerError, fmt.Errorf("jaira web: internal error: %v", p))
			}
		}()
		next.ServeHTTP(rec, r)
	})
}

// recorder keeps the status, and for an error the start of the body, which
// is where the handlers put the cause.
type recorder struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (r *recorder) WriteHeader(status int) {
	if r.status == 0 {
		r.status = status
	}
	r.ResponseWriter.WriteHeader(status)
}

func (r *recorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	if r.status >= 400 && r.body.Len() < 512 {
		r.body.Write(b[:min(len(b), 512-r.body.Len())])
	}
	return r.ResponseWriter.Write(b)
}
