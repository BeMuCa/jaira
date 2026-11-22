package pr

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// CacheFile is the name of the file in a repository's shared state directory
// that remembers pull request states across processes and worktrees.
const CacheFile = "pr-state.json"

// OpenFor is how long an open or closed answer is trusted: long enough that a
// teamlead running list and next in a row asks the forge once, short enough
// that a merge shows up on the next pass. A merged pull request stays merged,
// and that answer is kept for good.
const OpenFor = 2 * time.Minute

// Checker answers what a blocked-by entry pr#<n> waits on, asking the forge
// at most once per number and process, and not at all while the cache file
// holds a fresh answer. Every failure to read or write the cache only costs
// a forge call: it is never an error of its own.
type Checker struct {
	Dir  string
	File string // the cache file; "" keeps no cache across processes
	Run  Run
	Look func(string) (string, error)
	Now  func() time.Time
	// Offline answers from the cache only and never runs gh or glab — for a
	// caller that must stay instant, such as the TUI.
	Offline bool

	memo  map[int]Answer
	cache map[string]cached
}

// Answer is a pull request's state, or why there is none.
type Answer struct {
	State string // open, merged or closed; "" with Err set, or offline and unknown
	Err   error
}

type cached struct {
	State string    `json:"state"`
	At    time.Time `json:"at"`
}

// Ask returns the state of pull request n.
func (c *Checker) Ask(n int) Answer {
	if a, ok := c.memo[n]; ok {
		return a
	}
	if c.memo == nil {
		c.memo = map[int]Answer{}
	}
	now := c.now()
	if e, ok := c.load()[strconv.Itoa(n)]; ok && (e.State == "merged" || now.Sub(e.At) < OpenFor) {
		c.memo[n] = Answer{State: e.State}
		return c.memo[n]
	}
	if c.Offline || c.Run == nil {
		c.memo[n] = Answer{}
		return c.memo[n]
	}
	state, err := State(c.Run, c.Look, c.Dir, n)
	a := Answer{State: state, Err: err}
	c.memo[n] = a
	if err == nil {
		c.store(n, cached{State: state, At: now})
	}
	return a
}

func (c *Checker) now() time.Time {
	if c.Now != nil {
		return c.Now()
	}
	return time.Now()
}

func (c *Checker) load() map[string]cached {
	if c.cache == nil {
		c.cache = readCache(c.File)
	}
	return c.cache
}

func readCache(path string) map[string]cached {
	m := map[string]cached{}
	if path == "" {
		return m
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return m
	}
	if json.Unmarshal(b, &m) != nil {
		return map[string]cached{}
	}
	return m
}

// store re-reads the file before it writes, so an answer another process
// wrote meanwhile is kept, and replaces it through a temp file and a rename:
// a reader sees the old file or the new one, never half of one.
func (c *Checker) store(n int, e cached) {
	c.load()[strconv.Itoa(n)] = e
	if c.File == "" {
		return
	}
	m := readCache(c.File)
	m[strconv.Itoa(n)] = e
	b, err := json.Marshal(m)
	if err != nil {
		return
	}
	dir := filepath.Dir(c.File)
	if os.MkdirAll(dir, 0o755) != nil {
		return
	}
	f, err := os.CreateTemp(dir, ".pr-state-*")
	if err != nil {
		return
	}
	tmp := f.Name()
	_, werr := f.Write(b)
	if cerr := f.Close(); werr != nil || cerr != nil {
		_ = os.Remove(tmp)
		return
	}
	if os.Rename(tmp, c.File) != nil {
		_ = os.Remove(tmp)
	}
}

// Merged reports whether the pull request is merged — the one answer that
// clears a blocked-by entry pr#<n>.
func (a Answer) Merged() bool { return a.Err == nil && a.State == "merged" }

// Why is what a ticket waiting on this pull request is held back by, as
// "waits for PR #12: <Why>" prints it; empty when merged.
func (a Answer) Why() string {
	switch {
	case a.Err != nil:
		return "state unknown (" + a.Err.Error() + ")"
	case a.State == "merged":
		return ""
	case a.State == "open":
		return "open"
	case a.State == "closed":
		return "closed without merge — clear blocked-by or reject the ticket"
	}
	return "not asked here — jaira list or next asks the forge"
}
