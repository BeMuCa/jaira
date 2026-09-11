// Package bgrun runs a jaira subcommand in the background, at most every so
// often, without any command ever waiting for it.
//
// Three things in this program work that way — the release check, the board's
// snapshot, and the ref fetch — and they are the same five careful details
// every time. Written out three times, one of those details would eventually be
// missing from one of them, and the symptom would be a command that hangs, or a
// process that spawns itself for ever.
//
// The details, and why each is there:
//
//   - The stamp is written BEFORE the run, not after. A run that dies or hangs
//     must not make the next command start another one.
//   - The child is given the disable switch as an environment variable, so it
//     can never decide the task is due and spawn a grandchild. That is the same
//     switch a user flips to turn the task off, so there is one mechanism here
//     rather than two.
//   - Nothing waits. The parent is about to exit and the child is reparented to
//     init; a background task nobody asked for must never be something somebody
//     waits for.
//   - Never from a test binary, which has no such subcommand and would re-run
//     its own suite, detached, reaching the network on the way.
//   - Stdio goes to os.DevNull: a child writing to the terminal after the
//     parent has exited is the exact failure this design avoids.
package bgrun

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Task is one background job: what to run, how often, and how to switch it off.
type Task struct {
	// Stamp is the file recording when this task last ran, normally under the
	// board's state directory — never inside the repository, because when a
	// background job last ran is this machine's business.
	Stamp string

	// DisableEnv is the environment variable that turns the task off, and is
	// also handed to the child as its recursion guard.
	DisableEnv string

	// Args is the jaira subcommand to run, e.g. []string{"fetch", "--quiet"}.
	Args []string

	// Dir is the board the child should act on. It must be set explicitly: a
	// child inherits the parent's working directory, not the parent's -C flag,
	// so a task spawned from a command run with -C would otherwise quietly
	// operate on whatever directory the shell happened to be in.
	Dir string

	// Every is how old the last run may be before another is due. Zero means
	// the task is never due on its own.
	Every time.Duration
}

// Mark is what gets recorded.
type Mark struct {
	RanAt time.Time `json:"ran_at"`
	Note  string    `json:"note,omitempty"`
}

// Disabled reports whether the task is switched off for this process.
func (t Task) Disabled() bool { return t.DisableEnv != "" && os.Getenv(t.DisableEnv) == "1" }

// Read returns the last mark, or a zero Mark. Every error is swallowed: a
// missing or damaged mark is indistinguishable from "never ran" for every
// purpose this serves, and the cost of getting it wrong is one extra run.
func (t Task) Read() Mark {
	b, err := os.ReadFile(t.Stamp)
	if err != nil {
		return Mark{}
	}
	var m Mark
	if err := json.Unmarshal(b, &m); err != nil {
		return Mark{}
	}
	return m
}

// Write records a run.
func (t Task) Write(m Mark) error {
	if t.Stamp == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(t.Stamp), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	tmp := t.Stamp + ".tmp"
	if err := os.WriteFile(tmp, append(b, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, t.Stamp)
}

// Due reports whether the task should run now.
func (t Task) Due() bool {
	if t.Disabled() || t.Every <= 0 || t.Stamp == "" {
		return false
	}
	m := t.Read()
	if m.RanAt.IsZero() {
		return true
	}
	return time.Since(m.RanAt) >= t.Every
}

// Spawn stamps the task and starts a detached child, without waiting. It is a
// no-op when the task is not due, so a caller is one line.
func (t Task) Spawn() {
	if !t.Due() {
		return
	}
	if err := t.Write(Mark{RanAt: time.Now().UTC()}); err != nil {
		return
	}
	exe, err := os.Executable()
	if err != nil || isTestBinary(exe) {
		return
	}
	devNull, err := os.OpenFile(os.DevNull, os.O_RDWR, 0)
	if err != nil {
		return
	}
	defer devNull.Close()

	cmd := exec.Command(exe, t.Args...)
	cmd.Dir = t.Dir
	cmd.Stdin, cmd.Stdout, cmd.Stderr = devNull, devNull, devNull
	if t.DisableEnv != "" {
		cmd.Env = append(os.Environ(), t.DisableEnv+"=1")
	}
	_ = cmd.Start()
}

func isTestBinary(exe string) bool {
	base := strings.TrimSuffix(filepath.Base(exe), ".exe")
	return strings.HasSuffix(base, ".test")
}
