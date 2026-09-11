// Package hook calls the user's own script when a ticket moves or is claimed.
//
// It is the delivery escape hatch. Fetching refs tells a clone what changed the
// next time it looks, which is enough for a board but not for "tell me now" —
// and the ways people want to be told (a Slack webhook, ntfy.sh, Telegram, a
// terminal bell) are neither few nor stable. So jaira calls a script and brings
// no dependency of its own: the integration is the user's, the invocation is
// ours.
//
// The contract is environment variables rather than arguments, so a script can
// read only what it cares about and adding a field later breaks nothing.
package hook

import (
	"os"
	"os/exec"
	"strings"
	"time"
)

// timeout caps the script. A hook is a courtesy on the end of a command the
// user is waiting for, and a script that blocks — a curl to a host that is
// down — must not become the command's runtime.
const timeout = 5 * time.Second

// Event is what happened, and to which ticket.
type Event struct {
	// Name is "move" or "claim".
	Name string
	// ID is the ticket's full id.
	ID string
	// Title, Status and Assignee describe the ticket after the change.
	Title    string
	Status   string
	Assignee string
	// Actor is who did it.
	Actor string
	// Root is the board's directory, so a script can work out which project
	// this is without being told separately.
	Root string
}

// Run calls the script, waiting no longer than the timeout.
//
// Everything about a failure is silent: a hook is optional, its exit code is
// its own business, and a board that reported a broken script on every move
// would push the user to remove the feature rather than fix the script. The
// return value says whether it ran cleanly, for a caller that wants to say so
// on request.
func Run(script string, ev Event) bool {
	script = strings.TrimSpace(script)
	if script == "" {
		return false
	}
	cmd := exec.Command(script)
	cmd.Dir = ev.Root
	cmd.Env = append(os.Environ(),
		"JAIRA_EVENT="+ev.Name,
		"JAIRA_TICKET="+ev.ID,
		"JAIRA_TITLE="+ev.Title,
		"JAIRA_STATUS="+ev.Status,
		"JAIRA_ASSIGNEE="+ev.Assignee,
		"JAIRA_ACTOR="+ev.Actor,
		"JAIRA_ROOT="+ev.Root,
	)
	// The script's own output belongs to the script. Discarding it keeps a
	// chatty hook from interleaving with a command's stdout, which an agent
	// may be parsing as JSON.
	cmd.Stdout, cmd.Stderr = nil, nil
	if err := cmd.Start(); err != nil {
		return false
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		return err == nil
	case <-time.After(timeout):
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		return false
	}
}
