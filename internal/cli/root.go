// Package cli implements the jaira command surface.
//
// The CLI is the only write path into the store. Everything — the TUI, an agent
// working the board, a shell script — goes through the same core calls, so no
// caller can invent its own idea of what a valid ticket is. It is also designed
// to be driven by an agent over bash rather than by a human alone: read commands
// speak JSON on request, errors are machine-readable, and exit codes are a
// documented contract rather than an accident.
package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/gate"
	"github.com/BeMuCa/jaira/core/gitrepo"
	coreidentity "github.com/BeMuCa/jaira/core/identity"
	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/release"
	"github.com/BeMuCa/jaira/core/ticket"
)

// Exit codes are part of the contract with any agent driving this tool.
// Changing one is a breaking change.
const (
	ExitOK         = 0
	ExitError      = 1 // unexpected failure
	ExitUsage      = 2 // bad flags or arguments (cobra's own convention)
	ExitValidation = 3 // a gate refused the operation
	ExitBlocked    = 4 // unresolved dependencies
	ExitNotFound   = 5 // no such ticket, or an ambiguous prefix
)

// codedError carries an exit code and a machine-readable reason.
type codedError struct {
	code       int
	reason     string
	message    string
	violations gate.Violations
}

func (e *codedError) Error() string { return e.message }

func fail(code int, reason, format string, args ...any) error {
	return &codedError{code: code, reason: reason, message: fmt.Sprintf(format, args...)}
}

// exactArgs and minArgs mirror cobra's validators but report a usage error, so a
// wrong argument count exits 2 like a bad flag does rather than 1. Agents branch
// on these codes, so "usage mistake" and "something went wrong" must not collide.
func exactArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != n {
			return fail(ExitUsage, "usage", "%s accepts %d argument(s), received %d\n\nUsage: %s",
				cmd.Name(), n, len(args), cmd.UseLine())
		}
		return nil
	}
}

func minArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < n {
			return fail(ExitUsage, "usage", "%s needs at least %d argument(s), received %d\n\nUsage: %s",
				cmd.Name(), n, len(args), cmd.UseLine())
		}
		return nil
	}
}

func noArgs() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fail(ExitUsage, "usage", "%s takes no arguments, received %d\n\nUsage: %s",
				cmd.Name(), len(args), cmd.UseLine())
		}
		return nil
	}
}

// globals are flags shared by every subcommand.
type globals struct {
	jsonOut bool
	dir     string
}

var g globals

// Execute runs the CLI and returns the process exit code.
func Execute(version string) int {
	// The TUI creates boards too, and stamps them the same way, so this is set
	// once here rather than threaded as a parameter through every call that
	// might need it.
	release.Current = version
	root := newRoot(version)
	err := root.Execute()
	if err == nil {
		return ExitOK
	}
	return report(err)
}

// report renders an error in the requested format and maps it to an exit code.
func report(err error) int {
	var ce *codedError
	code := ExitError
	reason := "error"
	if errors.As(err, &ce) {
		code, reason = ce.code, ce.reason
	} else {
		switch {
		case errors.Is(err, ticket.ErrNotFound), errors.Is(err, ticket.ErrAmbiguous):
			code, reason = ExitNotFound, "not_found"
		case errors.Is(err, ticket.ErrNoStore):
			code, reason = ExitError, "no_store"
		case errors.Is(err, ticket.ErrUnsafeYAML):
			code, reason = ExitValidation, "unsafe_yaml"
		}
	}

	if g.jsonOut {
		// Structured errors on stderr so an agent never has to regex a sentence
		// to find out why something failed.
		payload := map[string]any{
			"error": map[string]any{
				"reason":  reason,
				"code":    code,
				"message": err.Error(),
			},
		}
		if ce != nil && len(ce.violations) > 0 {
			payload["error"].(map[string]any)["violations"] = ce.violations
		}
		enc := json.NewEncoder(os.Stderr)
		enc.SetIndent("", "  ")
		_ = enc.Encode(payload)
	} else {
		fmt.Fprintf(os.Stderr, "jaira: %v\n", err)
	}
	return code
}

func newRoot(version string) *cobra.Command {
	root := &cobra.Command{
		Use:   "jaira",
		Short: "A git-native kanban board for work done with coding agents",
		Long: `jaira tracks the tasks a coding agent generates, as markdown files
committed inside your repository.

Because the tickets are files in the repo, the board needs no server and no
accounts: a teammate clones and runs jaira, and sees the same board. Each lane
can bind a prompt and a model tier, which turns the board into a pipeline you
can watch rather than one opaque agent run.

Exit codes:
  0  success
  1  unexpected error
  2  usage error
  3  a gate refused the operation
  4  unresolved dependencies
  5  no such ticket, or an ambiguous id prefix`,
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			// Bare `jaira` opens the launcher, so one window oversees every board
			// rather than only the one you happen to be standing in. `jaira board`
			// is the direct route when you already know where you are.
			return runHome(cmd)
		},
	}
	root.PersistentFlags().BoolVar(&g.jsonOut, "json", false, "emit JSON on stdout and structured errors on stderr")
	root.PersistentFlags().StringVarP(&g.dir, "dir", "C", ".", "run as if jaira was started in this directory")

	root.AddCommand(
		newInitCmd(),
		newCreateCmd(),
		newListCmd(),
		newShowCmd(),
		newSetCmd(),
		newDoDCmd(),
		newValidateCmd(),
		newArchiveCmd(),
		newRestoreCmd(),
		newSyncOutCmd(),
		newNoteCmd(),
		newResumeCmd(),
		newMoveCmd(),
		newNextCmd(),
		newLanesCmd(),
		newBoardCmd(),
		newMergeDriverCmd(),
		newResolveCmd(),
		newSyncTasksCmd(),
		newTasksCmd(),
		newCheckpointCmd(),
		newSessionsCmd(),
		newClaimCmd(),
		newProjectsCmd(),
		newShareCmd(),
		newUpdateCmd(),
	)
	// Usage errors must exit 2 rather than 1.
	root.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return fail(ExitUsage, "usage", "%v", err)
	})
	return root
}

// openStore finds the store for the current directory.
//
// Opening a store is also where a shared board gets its merge driver bound, and
// where a stale board is nudged toward 'jaira update', since this is the one
// path every command goes through — including the first command a teammate
// runs after cloning.
func openStore() (*ticket.Store, error) {
	s, err := ticket.Discover(g.dir)
	if err != nil {
		return nil, err
	}
	bindDriverIfShared(s)
	nudgeIfStale(s)
	return s, nil
}

// bestEffortRoot returns the board root for the current directory, or "" if
// there is none. It is for callers that want to scope lane loading to a
// project when one is in hand (jaira lanes, the merge driver) but must still
// work — catalogue only — outside any board, unlike openStore which requires
// one.
func bestEffortRoot() string {
	s, err := ticket.Discover(g.dir)
	if err != nil {
		return ""
	}
	return s.Root
}

// loadEnv assembles the state gates need.
func loadEnv(s *ticket.Store) (gate.Env, []*ticket.Ticket, error) {
	lanes, err := lane.Load(s.Root)
	if err != nil {
		return gate.Env{}, nil, err
	}
	all, err := s.List()
	if err != nil {
		var pe *ticket.PartialError
		if !errors.As(err, &pe) {
			return gate.Env{}, nil, err
		}
		// Unreadable tickets are surfaced but must not stop the command.
		if !g.jsonOut {
			fmt.Fprintf(os.Stderr, "jaira: warning: %v\n", pe)
		}
	}
	for _, w := range lanes.Warnings {
		if !g.jsonOut {
			fmt.Fprintf(os.Stderr, "jaira: warning: %s\n", w)
		}
	}
	repo := &gitrepo.Repo{Dir: s.Root}
	return gate.Env{
		Lanes: lanes,
		All:   all,
		// Reaches git only when the gate actually needs it: CheckAdvance calls
		// this closure only for a ticket recording no commits of its own, at
		// the requires-commits lane, never merely to render the board.
		DeriveCommits: func(t *ticket.Ticket) []string {
			shas, err := repo.CommitsForTicket(t.Path, t.ID)
			if err != nil {
				return nil
			}
			return shas
		},
	}, all, nil
}

// emit writes a value as JSON.
func emit(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// identity determines who is acting, preferring git's configured name so tickets
// are attributed the same way commits are.
func identity() string {
	return coreidentity.Current(g.dir)
}
