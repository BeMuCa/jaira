package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BeMuCa/jaira/core/milestone"
	"github.com/BeMuCa/jaira/core/ticket"
	"github.com/spf13/cobra"
)

func newArchiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archive [id]",
		Short: "Take a ticket off the board, or list what has been taken off",
		Long: `Moves a ticket into .jaira/archive/. With no argument, lists the archive.

The file is moved, not deleted. The board exists so you can still answer what a
task was for months later, and on a private board a deleted file is not in git
history to recover from — so this command removes nothing, ever. 'jaira restore'
moves a ticket back.

Before the move, every commit git can find for the ticket is stamped onto it,
best-effort: a ticket with no commits is archived all the same. This works
from any lane; for finished work, 'jaira logbook' files the ticket under the
day it was finished and refuses one that has not reached the terminal lane.

'jaira delete' is the one command that does remove a file, for a ticket that
should never have existed. It asks for the handle typed back first.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return fail(ExitUsage, "usage", "archive takes at most one ticket id, received %d", len(args))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			w := cmd.OutOrStdout()

			if len(args) == 0 {
				names, err := archivedNames(s)
				if err != nil {
					return err
				}
				rejected := archivedRejections(s, names)
				if g.jsonOut {
					// rejected is keyed by file and holds only the rejected
					// ones, so archived stays the list of names it always was.
					byFile := map[string]map[string]string{}
					for n, r := range rejected {
						byFile[n] = rejectionJSON(r)
					}
					return emit(w, map[string]any{"archived": names, "count": len(names), "rejected": byFile})
				}
				if len(names) == 0 {
					fmt.Fprintf(w, "The archive is empty.\n")
					return nil
				}
				for _, n := range names {
					if r, ok := rejected[n]; ok {
						fmt.Fprintf(w, "%s\n    rejected %s\n", n, rejectionLine(r))
						continue
					}
					fmt.Fprintf(w, "%s\n", n)
				}
				fmt.Fprintf(w, "\n%d archived. Bring one back with 'jaira restore <file>'.\n", len(names))
				return nil
			}

			t, err := s.Load(args[0])
			if err != nil {
				return err
			}
			// Stamp what git can find before the ticket leaves the board —
			// the same moment 'jaira logbook' stamps at. Best-effort here:
			// archive has never required commits and must not start
			// refusing because a derivation found nothing.
			if env, _, envErr := loadEnv(s); envErr == nil {
				_, _ = s.StampCommits(t, env.DeriveCommits)
			}
			dst, err := s.Archive(args[0])
			if err != nil {
				return err
			}
			if g.jsonOut {
				return emit(w, map[string]any{
					"archived": true, "id": t.ID, "handle": ticket.Handle(t.ID),
					"path": dst, "file": filepath.Base(dst),
				})
			}
			fmt.Fprintf(w, "Archived %s  %s\n", ticket.Handle(t.ID), t.Title)
			fmt.Fprintf(w, "Moved to %s — restore it with 'jaira restore %s'.\n",
				filepath.Join(ticket.DirName, ticket.ArchiveSubdir), filepath.Base(dst))
			return nil
		},
	}
	return cmd
}

func newRestoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restore <file|milestone>",
		Short: "Put an archived ticket, or a filed milestone, back on the board",
		Long: `Moves a file out of .jaira/archive/ or .jaira/logbook/ and back onto the board,
named by its filename.

A milestone filed with 'jaira logbook <name>' comes back by its name, its ID or
<name>-<ID>: the "filed" line comes off its ref, and with it every clone shows
it again after its next fetch. A milestone file an older build filed into the
logbook is named by its file name, as before; it goes back onto its ref, and
the file is deleted.`,
		Args: exactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			if done, err := restoreMilestone(cmd, s, args[0]); done || err != nil {
				return err
			}
			// The milestone lock covers the whole restore: a milestone file
			// an older build filed comes back through .jaira/milestones/ and
			// onto its ref, and which of the two this is only becomes clear
			// once the file has moved. Ordering a ticket restore behind a
			// milestone write costs nothing.
			refs.FetchRefs()
			unlock, err := s.Lock(milestoneLockName)
			if err != nil {
				return err
			}
			defer unlock()

			dst, err := s.Restore(args[0])
			if err != nil {
				return &codedError{code: ExitNotFound, reason: "not_archived", message: err.Error()}
			}
			restored, err := unfileMilestone(s, dst)
			if err != nil {
				return err
			}
			if g.jsonOut {
				return emit(cmd.OutOrStdout(), map[string]any{"restored": true, "path": dst, "milestone": restored})
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Restored %s to the board.\n", filepath.Base(dst))
			return nil
		},
	}
}

// restoreMilestone takes the "filed" line off the milestone arg names, by its
// name, its ID or <name>-<ID>. done is false when arg names no milestone on
// the refs, and the caller goes on to restore a file.
func restoreMilestone(cmd *cobra.Command, s *ticket.Store, arg string) (done bool, err error) {
	if strings.HasSuffix(arg, ".md") {
		return false, nil
	}
	if _, _, err := milestone.NormalizeName(arg); err != nil {
		return false, nil
	}
	importMilestones(s)
	var amb *milestone.AmbiguousError
	var onBoard bool
	ms, err := milestone.EditFiled(s, refs, strings.TrimSpace(arg), func(ms *milestone.Milestone) error {
		if !ms.Filed() {
			onBoard = true
			return nil
		}
		ms.SetStatus("")
		return nil
	})
	switch {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case errors.As(err, &amb):
		return true, refuseAmbiguous(amb)
	case err != nil:
		return true, err
	case onBoard:
		return true, fail(ExitValidation, "milestone_on_board",
			"milestone %q is on the board, not filed — a filed one of the same name is restored by its <name>-<ID>, which 'jaira logbook' printed when it filed it", ms.Name)
	}
	if g.jsonOut {
		return true, emit(cmd.OutOrStdout(), map[string]any{"restored": true, "milestone": ms.Name, "id": ms.ID, "ref": ms.Ref()})
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Restored milestone %s to the board (%s).\n", milestoneLabel(ms), ms.Ref())
	return true, nil
}

// unfileMilestone puts a milestone file an older build filed into the logbook,
// just moved back to .jaira/milestones/, onto its ref and takes the "filed"
// line off it there, reporting the milestone's name if that is what was
// restored. The caller holds milestoneLockName.
//
// A failure fails the restore rather than being swallowed: reporting success
// would put a milestone nobody can see on a board and call it restored.
func unfileMilestone(s *ticket.Store, dst string) (string, error) {
	if filepath.Base(filepath.Dir(dst)) != ticket.MilestonesSubdir {
		return "", nil
	}
	name, id, _ := milestone.FileIdentity(strings.TrimSuffix(filepath.Base(dst), ".md"))
	res, err := refs.ImportMilestones(s.Root)
	reportImported(res.Moved, res.Waiting)
	reportConflicts(res.Conflicts)
	if err != nil {
		return name, fmt.Errorf("%s was moved back to %s but could not be put on its ref, so it is still filed: %w",
			filepath.Base(dst), milestone.Dir(s.Root), err)
	}
	if len(res.Waiting) > 0 {
		return name, fail(ExitError, "milestone_waiting",
			"%s was moved back to %s, but this clone could not reach the remote to tell it from a teammate's milestone — 'jaira fetch' takes it onto its ref, then 'jaira restore %s-%s' unfiles it",
			filepath.Base(dst), milestone.Dir(s.Root), name, id)
	}
	ms, err := milestone.Load(refs, id)
	if err != nil {
		return name, fmt.Errorf("%s was put on its ref but could not be read back, so it is still filed: %w", filepath.Base(dst), err)
	}
	if ms.Filed() {
		ms.SetStatus("")
		if err := refs.RecordMilestone(ms.ID, ms.FileName(), ms.Bytes()); err != nil {
			return name, fmt.Errorf("the %q line could not be taken off %s, so it stays off the board: %w",
				milestone.StatusFiled, ms.Ref(), err)
		}
	}
	return name, nil
}

func archivedNames(s *ticket.Store) ([]string, error) {
	entries, err := os.ReadDir(s.ArchiveDir())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out, nil
}

// archivedRejections reads the rejection off every archived file that carries
// one. A file that will not parse is listed without one rather than failing
// the listing: the archive is where damaged tickets end up, too.
func archivedRejections(s *ticket.Store, names []string) map[string]ticket.Rejection {
	out := map[string]ticket.Rejection{}
	for _, n := range names {
		t, err := ticket.LoadHeader(filepath.Join(s.ArchiveDir(), n))
		if err != nil || t.Rejected.Kind == "" {
			continue
		}
		out[n] = t.Rejected
	}
	return out
}
