package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// stopHookSnippet is the settings.json fragment 'jaira hook print' emits.
//
// Everything about the shell command is deliberate:
//
//   - It checks "stop_hook_active" first and gets out of the way when it is
//     set. Claude Code sets that field on the payload once a stop hook has
//     already blocked once, and a hook that keeps blocking on a condition it
//     cannot itself resolve wedges the session. Reading it with grep costs the
//     one read of stdin the hook needs anyway.
//   - It parses nothing. jq is not installed everywhere (it is not on the
//     machine this was written on) and a hook that depends on it fails on the
//     first teammate who lacks it. The two matches it needs are single fields
//     with fixed names, and '": *true' covers both the indented and the
//     compact encoding of the same JSON.
//   - Every failure exits 0. No board, no jaira on PATH, a broken repo — the
//     stderr redirect swallows it, the grep finds nothing, and the session ends
//     normally. This copies hooks/sync-tasks.sh: a board that cannot be read is
//     a nuisance, a hook that traps someone in a session is much worse.
//   - It exits 2, not 1. Exit 2 is the only code Claude Code treats as
//     blocking; 1 is a non-blocking error that lets the stop through. The
//     stderr line is what Claude is handed as the reason to keep going, so it
//     names the command that answers "which lane".
//
// It reports on any agentic lane with waiting work, not only the ticket the
// session was last touching, because the lane nobody drives is the one that
// fills up — which is what this exists to stop.
const stopHookSnippet = `{
  "hooks": {
    "Stop": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "grep -q '\"stop_hook_active\": *true' && exit 0; jaira next --per-lane --json 2>/dev/null | grep -q '\"agentic\": *true' || exit 0; echo 'jaira: an agentic lane still has work waiting. Run jaira next --per-lane, work the lane, and move the ticket on before stopping.' >&2; exit 2"
          }
        ]
      }
    ]
  }
}
`

// newHookCmd groups the Claude Code hook snippets. It is a separate parent from
// 'jaira update', which writes the agent instruction block: the block is
// jaira's to maintain, whereas a hook lives in the user's own settings.json and
// is theirs to install or not.
func newHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hook",
		Short: "Claude Code hooks for this board",
		Args:  noArgs(),
	}
	cmd.AddCommand(newHookPrintCmd())
	return cmd
}

// newHookPrintCmd prints the Stop-hook snippet to stdout and nothing else — no
// file is written and no settings are touched, because settings.json is the
// user's file. 'jaira hook print >> ...' is not safe on a JSON file, so the
// snippet is meant to be read and merged by hand.
//
// This is the enforcement half of what 'jaira move' can only suggest. The move
// naming the next command is advice a model may take or drop; the hook is the
// environment refusing, which is why it is worth having both.
func newHookPrintCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "print",
		Short: "Print a Claude Code Stop-hook snippet that refuses to end a session with agentic lanes left unworked",
		Long: `Prints a settings.json snippet for Claude Code's Stop hook. Nothing is
installed: merge it into ~/.claude/settings.json (or the project's
.claude/settings.json) yourself.

Once installed, a session cannot end while 'jaira next --per-lane' reports an
agentic lane with work waiting. Claude is handed the reason and carries on
working that lane instead of stopping.

The hook stays out of the way when it cannot do its job: no board, no jaira on
PATH, or an already-blocked stop and it lets the session end. It never blocks
twice for the same stop, so it cannot trap you in a session.

That last guard reads the hook payload on stdin and has no other source: run
the command with stdin closed or pointed at /dev/null and it blocks every time,
unconditionally. It is a guard, not defence in depth — keep the payload on
stdin if you adapt the snippet.

This is opt-in and Claude Code specific, which is why it is a snippet you paste
rather than something 'jaira init' arranges for you.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Fprint(cmd.OutOrStdout(), stopHookSnippet)
			return nil
		},
	}
}
