package cli

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/release"
	"github.com/BeMuCa/jaira/core/ticket"
)

// diffFile is one role, cockpit or lane file as 'jaira update diff' classifies it. For a
// lane, dir holds the lane's id and name its path.
type diffFile struct {
	dir, name, action, release string
}

// newUpdateDiffCmd is 'jaira update diff': for every role, cockpit file and lane, whether a
// bare update would leave it current, bring it up to date, or leave it alone because it was
// edited. It writes nothing.
func newUpdateDiffCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "diff",
		Short: "List which roles, cockpit files and lanes are current, behind or edited; writes nothing",
		Long: `Lists every role, cockpit file and lane jaira manages here, and says for each one
whether it is current, behind (with the release that shipped the version on
disk), or edited. A bare 'jaira update' brings the behind roles and cockpit files
up to date and leaves edited ones alone; lanes are brought up to date only with
--lanes.

Outside a board only your personal skills directories are listed. It writes
nothing, not even the version stamp.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			root := ""
			s, err := openStore()
			switch {
			case err == nil:
				root = s.Root
			case errors.Is(err, ticket.ErrNoStore):
			default:
				return err
			}
			runs, err := refreshRoles(root, true)
			if err != nil {
				return err
			}
			// lane.RefreshShipped loads the board's lanes, and loading a board with no
			// lane files of its own writes them. Such a board runs the built-in lanes,
			// which this binary already has, so there is nothing to compare.
			var lanes []lane.Refreshed
			hasLanes := root != "" && lane.ProjectLanesActive(root)
			if hasLanes {
				if lanes, err = lane.RefreshShipped(root, true); err != nil {
					return err
				}
			}
			w := cmd.OutOrStdout()
			if g.jsonOut {
				return emit(w, diffJSON(root, runs, lanes, hasLanes))
			}
			writeDiff(w, root, runs, lanes, hasLanes)
			return nil
		},
	}
}

// diffState is the word a file gets in the diff. role and lane results use the same
// action words, so one mapping covers both.
func diffState(action, release string) string {
	switch action {
	case "unchanged":
		return "current"
	case "updated":
		if release == "" || release == "-" {
			return "behind (not in a release yet)"
		}
		return "behind (shipped in " + release + ")"
	case "written":
		return "behind (missing)"
	case "skipped":
		return "edited"
	}
	return action
}

// diffFiles flattens what a dry run found into one entry per file, roles first.
func diffFiles(runs []roleRun, lanes []lane.Refreshed) (roles, lanesOut []diffFile) {
	for _, r := range runs {
		for _, res := range r.results {
			rel, err := filepath.Rel(r.target.Dir, res.Path)
			if err != nil {
				rel = res.Path
			}
			roles = append(roles, diffFile{dir: r.target.Dir, name: rel, action: string(res.Action), release: res.Release})
		}
	}
	for _, l := range lanes {
		lanesOut = append(lanesOut, diffFile{dir: l.ID, name: l.Path, action: l.Action, release: l.Release})
	}
	return roles, lanesOut
}

func writeDiff(w io.Writer, root string, runs []roleRun, lanes []lane.Refreshed, hasLanes bool) {
	if root == "" {
		fmt.Fprintln(w, "no board here: only your personal skills directories were checked")
	}
	roles, lanesOut := diffFiles(runs, lanes)
	counts := map[string]int{}
	var lastDir string
	for _, f := range roles {
		if f.dir != lastDir {
			fmt.Fprintf(w, "%s\n", f.dir)
			lastDir = f.dir
		}
		state := diffState(f.action, f.release)
		counts[state]++
		fmt.Fprintf(w, "  %-44s %s\n", f.name, state)
	}
	if len(runs) == 0 {
		fmt.Fprintln(w, "roles: none installed where jaira writes them — 'jaira roles install --global' or '--project' installs them")
	}
	switch {
	case root == "":
	case !hasLanes:
		fmt.Fprintln(w, "lanes: this board has no lane files of its own; it runs the built-in lanes, so none is behind")
	default:
		fmt.Fprintln(w, "lanes:")
		for _, f := range lanesOut {
			state := diffState(f.action, f.release)
			counts[state]++
			fmt.Fprintf(w, "  %-44s %s\n", f.dir, state)
		}
	}
	behind, edited := 0, 0
	for state, n := range counts {
		switch {
		case state == "edited":
			edited += n
		case state != "current":
			behind += n
		}
	}
	fmt.Fprintf(w, "\n%d current, %d behind, %d edited (nothing was written)\n", counts["current"], behind, edited)
}

// diffJSON is the payload of 'jaira update diff --json'.
func diffJSON(root string, runs []roleRun, lanes []lane.Refreshed, hasLanes bool) map[string]any {
	roles, lanesOut := diffFiles(runs, lanes)
	roleEntries := make([]map[string]any, 0, len(roles))
	for _, f := range roles {
		roleEntries = append(roleEntries, map[string]any{
			"directory": f.dir, "file": f.name, "action": f.action,
			"state": diffState(f.action, f.release), "release": f.release,
		})
	}
	laneEntries := make([]map[string]any, 0, len(lanesOut))
	for _, f := range lanesOut {
		laneEntries = append(laneEntries, map[string]any{
			"lane": f.dir, "path": f.name, "action": f.action,
			"state": diffState(f.action, f.release), "release": f.release,
		})
	}
	return map[string]any{
		"root": root, "version": release.Current,
		"roles": roleEntries, "lanes": laneEntries, "lanes_on_board": hasLanes,
		"nothing_written": true,
	}
}
