package cli

import (
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/move"
	"github.com/BeMuCa/jaira/core/plugin"
	"github.com/BeMuCa/jaira/core/release"
	"github.com/BeMuCa/jaira/core/ticket"
)

// runPlugin starts the installed plugin args[0] names, when it names one and
// jaira has no command of that name. ok is false when nothing was started and
// the arguments belong to cobra — which then says "unknown command" exactly as
// before, so a typo with no plugin behind it reads the same.
//
// jaira's own commands always win: a plugin can never cover one, whatever it
// is called. Only the first argument is looked at; every argument after it is
// the plugin's, --json and --board included.
func runPlugin(root *cobra.Command, args []string) (code int, ok bool) {
	if len(args) == 0 || isBuiltin(root, args[0]) {
		return 0, false
	}
	prog := plugin.Program(args[0])
	if prog == "" {
		return 0, false
	}
	code, err := plugin.Run(prog, args[1:], pluginEnv(boardArg(args[1:])))
	if err != nil {
		fmt.Fprintf(os.Stderr, "jaira: plugin %s: %v\n", args[0], err)
		return ExitError, true
	}
	return code, true
}

// isBuiltin reports whether name is one of jaira's commands or an alias of
// one. help and completion are added by cobra only when it executes, so they
// are added here first, or a plugin called "help" would take the word over.
func isBuiltin(root *cobra.Command, name string) bool {
	root.InitDefaultHelpCmd()
	root.InitDefaultCompletionCmd()
	for _, c := range root.Commands() {
		if c.Name() == name || c.HasAlias(name) {
			return true
		}
	}
	return false
}

// boardArg is the value of a --board flag among a plugin's arguments, which
// jaira reads without taking it away: the plugin gets every argument as typed.
func boardArg(args []string) string {
	for i, a := range args {
		if a == "--" {
			break
		}
		if v, found := strings.CutPrefix(a, "--board="); found {
			return v
		}
		if a == "--board" && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

// pluginEnv is the environment a plugin is started in. JAIRA_ROOT stays empty
// outside a repository rather than refusing the call: a plugin's login or
// setup command has no board to work on.
func pluginEnv(board string) plugin.Env {
	root := ""
	if s, err := ticket.Discover(g.dir); err == nil {
		root = s.Root
	}
	return plugin.Env{Root: root, Board: cmp.Or(board, lane.DefaultBoardName), Version: release.Current}
}

// A move is fired from core/move itself, so the TUI's moves reach a plugin
// the same way the CLI's do.
func init() {
	move.Moved = func(s *ticket.Store, t *ticket.Ticket, from string, lanes *lane.Set) {
		firePlugins("move", s, t, from, lanes)
	}
}

// events, while it is set, takes plugin events off the caller's path: the
// terminal board sets it, so a move there never waits for a plugin, which may
// take up to five seconds. The CLI leaves it nil and runs them before it
// exits, as docs/PLUGINS.md describes.
var events chan func()

// queueEvents starts the one goroutine that runs queued events in the order
// they were queued. stop drains the queue and waits for it, so quitting the
// board loses no event.
func queueEvents() (stop func()) {
	ch := make(chan func(), 64)
	done := make(chan struct{})
	go func() {
		for fire := range ch {
			fire()
		}
		close(done)
	}()
	events = ch
	return func() {
		events = nil
		close(ch)
		<-done
	}
}

// runBoard runs the terminal board with its plugin events queued.
func runBoard(m tea.Model) error {
	stop := queueEvents()
	defer stop()
	_, err := tea.NewProgram(m).Run()
	return err
}

// firePlugins tells every plugin subscribed to event about t, after the write
// has succeeded. stdin is what 'jaira show <id> --json' prints, plus who made
// the change and the lane the ticket left (empty for every event but move).
// lanes may be nil: they are then loaded only when a plugin is listening, so a
// write with no subscriber pays for nothing. Nothing about it reaches the
// command: see plugin.Fire.
func firePlugins(event string, s *ticket.Store, t *ticket.Ticket, from string, lanes *lane.Set) {
	if t == nil {
		return
	}
	ids := plugin.Subscribers(event)
	if len(ids) == 0 {
		return
	}
	if lanes == nil {
		env, err := envFor(s, t)
		if err != nil {
			return
		}
		lanes = env.Lanes
	}
	j := showJSON(t, lanes, t.Path)
	j["actor"] = s.Actor
	j["from_lane"] = from
	payload, err := json.Marshal(j)
	if err != nil {
		return
	}
	env := plugin.Env{Root: s.Root, Board: cmp.Or(t.Board, lane.DefaultBoardName), Version: release.Current}
	fire := func() { plugin.Fire(ids, event, payload, env) }
	if events != nil {
		events <- fire
		return
	}
	fire()
}
