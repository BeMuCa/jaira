package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/release"
	"github.com/BeMuCa/jaira/core/role"
)

// roleRun is what an install did, or would do, in one skills directory.
type roleRun struct {
	target  role.Target
	results []role.Result
}

// updateSetup is what 'jaira update --roles --lanes' did, or with --dry-run would do,
// beyond the jaira section: the lane files and the skills directories. The zero value is
// a bare 'jaira update', which does neither.
type updateSetup struct {
	roles, lanes, dry bool
	laneResults       []lane.Refreshed
	roleRuns          []roleRun
}

// skipped reports whether a file was left alone because a person edited it.
func (u *updateSetup) skipped() bool {
	for _, r := range u.roleRuns {
		if role.SkippedAny(r.results) {
			return true
		}
	}
	for _, l := range u.laneResults {
		if l.Action == lane.RefreshSkipped {
			return true
		}
	}
	return false
}

// json adds the keys of the switches that were given to the payload a bare update emits.
// A bare update gains none, so what an agent parses stays as it was.
func (u *updateSetup) json(m map[string]any) map[string]any {
	if u.roles {
		m["roles"] = roleRunsJSON(u.roleRuns)
	}
	if u.lanes {
		m["lanes"] = laneResultsJSON(u.laneResults)
	}
	if u.roles || u.lanes {
		m["skipped"] = u.skipped()
	}
	if u.dry {
		m["dry_run"] = true
	}
	return m
}

// editedError is the refusal of an update that left an edited file alone.
func editedError() error {
	return fail(ExitValidation, "edited",
		"some files were edited here and left alone; 'jaira roles install --project --force' (or '--global --force') and 'jaira lanes use <id> --force' replace them")
}

func refreshRoles(root string, dryRun bool) ([]roleRun, error) {
	targets, err := role.Installed(root)
	if err != nil {
		return nil, err
	}
	runs := make([]roleRun, 0, len(targets))
	for _, t := range targets {
		results, err := t.Agent.Install(t.Dir, false, dryRun)
		if err != nil {
			return runs, err
		}
		runs = append(runs, roleRun{target: t, results: results})
	}
	return runs, nil
}

func printRoleRuns(w io.Writer, runs []roleRun, root string, dry bool) {
	if len(runs) == 0 {
		fmt.Fprintln(w, "roles: none installed where jaira writes them — 'jaira roles install --global' or '--project' installs them")
		return
	}
	for _, r := range runs {
		printRoleResults(w, r.target.Dir, r.results, dry)
		if role.SkippedAny(r.results) {
			if cmd, err := role.UpdateCommand(r.target.Dir, root, r.target.Global); err == nil {
				fmt.Fprintf(w, "edited files in %s: '%s --force' replaces them\n", r.target.Dir, cmd)
			}
		}
	}
}

func roleRunsJSON(runs []roleRun) []map[string]any {
	out := make([]map[string]any, 0, len(runs))
	for _, r := range runs {
		out = append(out, map[string]any{
			"agent": r.target.Agent.Name, "directory": r.target.Dir, "global": r.target.Global,
			"installed": roleResultsJSON(r.results),
		})
	}
	return out
}

func printLaneResults(w io.Writer, results []lane.Refreshed, dry bool) {
	counts := map[string]int{}
	for _, r := range results {
		counts[r.Action]++
		switch r.Action {
		case lane.RefreshUpdated:
			if dry {
				fmt.Fprintf(w, "would update %s\n", r.Path)
			} else {
				fmt.Fprintf(w, "updated %s\n", r.Path)
			}
		case lane.RefreshSkipped:
			fmt.Fprintf(w, "skipped %s (edited here; 'jaira lanes use %s --force' replaces it)\n", r.Path, r.ID)
		}
	}
	fmt.Fprintf(w, "lanes: %d updated, %d unchanged, %d skipped\n",
		counts[lane.RefreshUpdated], counts[lane.RefreshUnchanged], counts[lane.RefreshSkipped])
}

func laneResultsJSON(results []lane.Refreshed) []map[string]any {
	out := make([]map[string]any, 0, len(results))
	for _, r := range results {
		out = append(out, map[string]any{"lane": r.ID, "path": r.Path, "action": r.Action})
	}
	return out
}

// hintLines says, for the parts a bare update did not touch, that they are behind and which
// switch brings them up. It writes nothing and ignores every error of its own probes: a
// hint must never fail an update.
func hintLines(root string, lanes bool) string {
	var out string
	if !lanes {
		if res, err := lane.RefreshShipped(root, true); err == nil {
			n := 0
			for _, r := range res {
				if r.Action == lane.RefreshUpdated {
					n++
				}
			}
			if n > 0 {
				out += fmt.Sprintf("lanes: %d lane file(s) are an older version jaira shipped — 'jaira update --lanes' brings them up (add --dry-run to list them)\n", n)
			}
		}
	}
	return out
}

// updatePersonalRoles is 'jaira update --roles' outside a board: the personal skills
// directories only.
func updatePersonalRoles(cmd *cobra.Command, dryRun bool) (err error) {
	runs, err := refreshRoles("", dryRun)
	if err != nil {
		return err
	}
	u := &updateSetup{roles: true, dry: dryRun, roleRuns: runs}
	defer func() {
		if err == nil && u.skipped() {
			err = editedError()
		}
	}()
	w := cmd.OutOrStdout()
	if g.jsonOut {
		return emit(w, u.json(map[string]any{"root": "", "version": release.Current}))
	}
	fmt.Fprintln(w, "no board here: only your personal skills directories were checked")
	printRoleRuns(w, runs, "", dryRun)
	if dryRun {
		fmt.Fprintln(w, "nothing was written (--dry-run)")
	}
	return nil
}
