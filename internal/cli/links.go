package cli

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/link"
	"github.com/BeMuCa/jaira/core/ticket"
	"github.com/spf13/cobra"
)

// newLinksCmd prints everything connected to one ticket.
//
// It reads past the board on purpose. A link's other end is usually finished
// by the time anybody follows it, and a view that only knew the board
// answered "not found" for the exact tickets a reader most wanted: the
// blocker that cleared, the parent whose children shipped, the follow-up that
// closed.
func newLinksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "links <id>",
		Short: "Show every ticket connected to this one",
		Long: `Lists what this ticket waits on, blocks, is part of, contains, relates to
and follows — in both directions, and wherever the other end now lives: on the
board, on a ref nobody pulled yet, in the logbook, or in the archive.

Children are not a stored list. A ticket names its parent and the tree is read
back from that, so an epic can be nested as deep as it needs to be without any
ticket having to be edited twice.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			t, _, err := s.LoadAnywhere(args[0])
			if err != nil {
				return err
			}
			lanes, err := lane.Load(s.Root)
			if err != nil {
				return err
			}
			all, err := s.List()
			if err != nil {
				// An unreadable ticket must not hide every link of the ones
				// that are readable.
				var pe *ticket.PartialError
				if !errors.As(err, &pe) {
					return err
				}
			}
			ix := link.Build(s, lanes, all)
			entries := ix.Relations(t.ID)
			if g.jsonOut {
				return emit(cmd.OutOrStdout(), linksJSON(t, entries))
			}
			printLinks(cmd.OutOrStdout(), t, entries)
			return nil
		},
	}
	return cmd
}

func linksJSON(t *ticket.Ticket, entries []link.Entry) map[string]any {
	groups := map[string][]map[string]any{}
	for _, k := range link.Order {
		for _, e := range entries {
			if e.Kind != k {
				continue
			}
			groups[string(k)] = append(groups[string(k)], map[string]any{
				"id":     e.Ref.ID,
				"handle": ticket.Handle(e.Ref.ID),
				"title":  e.Ref.Title,
				"lane":   e.Ref.Lane,
				"place":  string(e.Ref.Place),
				"path":   e.Ref.Path,
				"done":   e.Ref.Done,
				"depth":  e.Ref.Depth,
			})
		}
	}
	return map[string]any{
		"id":     t.ID,
		"handle": ticket.Handle(t.ID),
		"title":  t.Title,
		"links":  groups,
	}
}

func printLinks(w io.Writer, t *ticket.Ticket, entries []link.Entry) {
	fmt.Fprintf(w, "%s  %s\n", ticket.Handle(t.ID), t.Title)
	fmt.Fprintf(w, "%s\n", strings.Repeat("─", 64))
	any := false
	for _, k := range link.Order {
		var group []link.Entry
		for _, e := range entries {
			if e.Kind == k {
				group = append(group, e)
			}
		}
		if len(group) == 0 {
			continue
		}
		any = true
		fmt.Fprintf(w, "\n%s\n", k.Title())
		for _, e := range group {
			fmt.Fprintf(w, "  %s%s\n", strings.Repeat("  ", e.Ref.Depth), e.Ref.Label()+"  — "+e.Ref.Whereabouts())
		}
	}
	if !any {
		fmt.Fprintf(w, "\nnothing is linked to this ticket yet\n")
	}
}
