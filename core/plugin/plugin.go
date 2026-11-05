// Package plugin starts the programs docs/PLUGINS.md describes: `jaira <id>`
// runs ~/.jaira/plugins/<id>/jaira-<id>, and a write a plugin subscribed to
// runs `jaira-<id> event <name>`.
//
// What runs is only what sits in the plugins directory. jaira never searches
// PATH for a jaira-<id>: a program somebody dropped on PATH is not a plugin a
// person chose to install, and the directory is the whitelist.
package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"time"
)

// InEventVar is set to "1" for a plugin run by an event. A write made from
// inside one fires no events, or a plugin that answers a set with a set — or
// two plugins answering each other — would never stop.
const InEventVar = "JAIRA_IN_EVENT"

// eventTimeout is the cap docs/PLUGINS.md promises, the same as a hook's: an
// event is a courtesy on the end of a command somebody is waiting for.
const eventTimeout = 5 * time.Second

// idRe is the id pattern of docs/schema/jaira-plugin.schema.json.
var idRe = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)

// manifestFile is the manifest an install leaves beside the program.
const manifestFile = "jaira-plugin.json"

// Dir is where installed plugins live: plugins/ beside settings.json, which is
// JAIRA_HOME when that is set, so tests and a sandboxed run never reach the
// real one.
func Dir() string {
	if v := os.Getenv("JAIRA_HOME"); v != "" {
		return filepath.Join(v, "plugins")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".jaira", "plugins")
}

// Program is the path of id's program, or "" when id is not a well-formed id
// or no such plugin is installed.
func Program(id string) string {
	if !idRe.MatchString(id) {
		return ""
	}
	dir := Dir()
	if dir == "" {
		return ""
	}
	name := "jaira-" + id
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	path := filepath.Join(dir, id, name)
	if fi, err := os.Stat(path); err != nil || !fi.Mode().IsRegular() {
		return ""
	}
	return path
}

// Env is what a plugin is told about the call, in the three variables the
// standard names.
type Env struct {
	// Root is the directory holding .jaira/, "" outside a repository.
	Root string
	// Board is the board the call works on, "default" for today's.
	Board string
	// Version is the running jaira's version.
	Version string
}

func (e Env) vars() []string {
	return append(os.Environ(),
		"JAIRA_ROOT="+e.Root,
		"JAIRA_BOARD="+e.Board,
		"JAIRA_VERSION="+e.Version,
	)
}

// Run starts prog with args on jaira's own stdin, stdout and stderr, in the
// directory jaira was started in, and returns the plugin's exit code. The error
// is only for a program that could not be started at all.
func Run(prog string, args []string, env Env) (int, error) {
	cmd := exec.Command(prog, args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.Env = env.vars()
	err := cmd.Run()
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return ee.ExitCode(), nil
	}
	return 0, err
}

// Subscribers are the ids of the installed plugins whose manifest names
// event, in order of their ids. A plugin without a readable manifest, or whose
// manifest names another id, hears nothing: the manifest is the subscription.
//
// Inside an event — JAIRA_IN_EVENT=1 — there are none, so a write a plugin
// makes from its own event calls nobody.
func Subscribers(event string) []string {
	if os.Getenv(InEventVar) == "1" {
		return nil
	}
	dir := Dir()
	if dir == "" {
		return nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var ids []string
	for _, e := range entries {
		id := e.Name()
		if !e.IsDir() || Program(id) == "" {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, id, manifestFile))
		if err != nil {
			continue
		}
		var m struct {
			ID     string   `json:"id"`
			Events []string `json:"events"`
		}
		if json.Unmarshal(raw, &m) != nil || m.ID != id {
			continue
		}
		if slices.Contains(m.Events, event) {
			ids = append(ids, id)
		}
	}
	// ReadDir returns the entries sorted by name, which is the promised order.
	return ids
}

// Fire runs `jaira-<id> event <name>` for each of ids — what Subscribers
// returned — one after the other, with payload on stdin. Nothing about a plugin's run reaches the
// caller: the write has already happened, so a plugin that fails, times out
// or is missing changes nothing about the command's result. Its output is
// discarded, as a hook's is, so it never mixes into JSON an agent is parsing.
func Fire(ids []string, event string, payload []byte, env Env) {
	for _, id := range ids {
		fireOne(id, event, payload, env)
	}
}

func fireOne(id, event string, payload []byte, env Env) {
	prog := Program(id)
	if prog == "" {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), eventTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, prog, "event", event)
	cmd.Dir = env.Root
	cmd.Env = append(env.vars(), InEventVar+"=1")
	cmd.Stdin = bytes.NewReader(payload)
	cmd.Stdout, cmd.Stderr = nil, nil
	// A plugin that left a child holding its stdin open would keep Wait
	// blocked past the kill; WaitDelay caps that too.
	cmd.WaitDelay = time.Second
	_ = cmd.Run()
}
