package cli

// This file implements 'jaira lanes market' — the lanes published in the jaira
// repository's lanes/ directory on GitHub, listed and adopted from the command
// line. The directory is the marketplace: anyone adds a lane with a pull
// request, CI parses it before it lands, and it reaches a machine only because
// someone adopted it on purpose, exactly like a teammate's shared lane.

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/market"
)

func newLanesMarketCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "market",
		Short: "List the lanes published in jaira's own repository, or adopt one",
		Long: `Lists the lanes in the lanes/ directory of github.com/BeMuCa/jaira — the
catalogue anyone can add to with a pull request. Each is fetched and parsed, so
the listing shows what a lane is, not a filename.

'jaira lanes market adopt <id>' downloads one into your catalogue
(~/.jaira/lanes); 'jaira lanes add <id>' then puts it on a board. Adopting means
agreeing to run whatever the file's prompt says at whatever model tier it
declares — read it first: the listing shows the description, and the file is
one click away in the repository.

Needs the network. Without it the command says so and nothing changes.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			entries, warnings, err := market.New().List(context.Background())
			if err != nil {
				return fail(ExitError, "market_unreachable", "%v", err)
			}
			w := cmd.OutOrStdout()
			if g.jsonOut {
				arr := make([]map[string]any, 0, len(entries))
				for _, e := range entries {
					arr = append(arr, map[string]any{
						"id": e.Lane.ID, "name": e.Lane.Name, "description": e.Lane.Description,
						"agentic": e.Lane.Agentic, "model_tier": e.Lane.ModelTier,
						"after": e.Lane.After, "rejects_to": e.Lane.RejectsTo,
						"path": e.Path, "url": e.URL,
					})
				}
				return emit(w, map[string]any{"lanes": arr, "warnings": warnings, "override": market.Overridden()})
			}
			if o := market.Overridden(); o != "" {
				fmt.Fprintf(w, "note: listing from %s\n", o)
			}
			if len(entries) == 0 {
				fmt.Fprintln(w, "the marketplace holds no lanes")
			} else {
				fmt.Fprintf(w, "%-14s %-18s %s\n", "ID", "NAME", "DESCRIPTION")
				for _, e := range entries {
					fmt.Fprintf(w, "%-14s %-18s %s\n", e.Lane.ID, e.Lane.Name, e.Lane.Description)
				}
				fmt.Fprintf(w, "\n%d lane(s). 'jaira lanes market adopt <id>' puts one in your catalogue.\n", len(entries))
			}
			for _, warn := range warnings {
				fmt.Fprintf(os.Stderr, "jaira: warning: %s\n", warn)
			}
			return nil
		},
	}
	cmd.AddCommand(newLanesMarketAdoptCmd())
	return cmd
}

func newLanesMarketAdoptCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "adopt <id>",
		Short: "Download a marketplace lane into your catalogue",
		Long: `Downloads the named lane from the marketplace and copies it into your
catalogue under its own id, the same move 'jaira lanes adopt <path>' makes for
a teammate's file. Refuses to overwrite an existing catalogue entry unless
--force is given. Then 'jaira lanes add <id>' puts it on a board.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := market.New()
			entries, _, err := client.List(context.Background())
			if err != nil {
				return fail(ExitError, "market_unreachable", "%v", err)
			}
			var pick *market.Entry
			var ids []string
			for i := range entries {
				ids = append(ids, entries[i].Lane.ID)
				if entries[i].Lane.ID == args[0] {
					pick = &entries[i]
				}
			}
			if pick == nil {
				return fail(ExitUsage, "no_such_lane", "no lane %q in the marketplace; available: %v", args[0], ids)
			}
			raw, err := client.Fetch(context.Background(), *pick)
			if err != nil {
				return fail(ExitError, "market_unreachable", "%v", err)
			}
			// Adopt parses a file and copies it under the parsed id; the download
			// goes through a temporary file so that one code path owns "a lane
			// from outside enters the catalogue".
			tmp, err := os.CreateTemp("", "jaira-market-*.md")
			if err != nil {
				return err
			}
			defer os.Remove(tmp.Name())
			if _, err := tmp.Write(raw); err != nil {
				tmp.Close()
				return err
			}
			tmp.Close()
			l, dst, err := lane.Adopt(tmp.Name(), lane.UserLanesDir(), force)
			if err != nil {
				return writeConflictError(err)
			}
			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, map[string]any{"id": l.ID, "path": dst, "url": pick.URL, "override": market.Overridden()})
			}
			if o := market.Overridden(); o != "" {
				fmt.Fprintf(w, "note: fetched from %s\n", o)
			}
			fmt.Fprintf(w, "adopted %s\nRead it before you run it: 'jaira lanes show %s'. Then 'jaira lanes add %s' puts it on this board.\n",
				filepath.Clean(dst), l.ID, l.ID)
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "overwrite an existing catalogue entry")
	return cmd
}
