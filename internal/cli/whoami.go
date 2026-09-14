package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/gitref"
	coreidentity "github.com/BeMuCa/jaira/core/identity"
	"github.com/BeMuCa/jaira/core/settings"
	"github.com/BeMuCa/jaira/core/ticket"
)

// newWhoamiCmd exists because the ownership rail compares strings and a person
// is several. When a move is refused with "belongs to <someone>" and that
// someone is you under another name, nothing else on the board would tell you
// which name jaira thinks you have.
//
// It answers the git half of the same question too — which remote this board
// uses, whether tickets travel on refs at all, and how many are lying on this
// disk only. Both halves are "what does jaira think about this environment
// before I start": the identity decides whose tickets these are, the remote
// decides where they go. Keeping them in one command is also what keeps this
// tool small; a second 'status' would be a second place to look for one answer.
func newWhoamiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Show the identity jaira acts as, and the state of this board",
		Long: `Prints the name jaira attributes work to, followed by every other string it
treats as the same person: git's user.email as well as user.name, plus anything
listed in the alias file.

A ticket assigned under a name that is not in this list is somebody else's as
far as the gates are concerned, and moving it will be refused. Add the name to
the alias file rather than passing --force to every move.

Inside a board it also prints where tickets go: the remote this board uses and
where that name came from, the remotes this repository actually has, whether
tickets travel on refs or stay as files, and how many tickets exist on this disk
only. A board that has quietly stopped carrying tickets on refs says so here,
before a ticket is written rather than after.`,
		Args: noArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			me := identity()
			aliases := coreidentity.Aliases(g.dir)
			path := coreidentity.AliasesPath()

			// A store is a bonus, never a requirement: whoami has always
			// answered outside a board and must keep doing so.
			st := boardState()

			if g.jsonOut {
				payload := map[string]any{
					"identity": me, "aliases": aliases, "alias_file": path,
				}
				if st != nil {
					payload["board"] = st.Root
					payload["remote"] = st.Remote
					payload["remote_source"] = st.Source
					payload["remotes"] = st.Remotes
					payload["ref_mode"] = st.RefMode
					if st.Reason != "" {
						payload["ref_mode_reason"] = st.Reason
					}
					payload["file_only"] = st.FileOnly
				}
				return emit(cmd.OutOrStdout(), payload)
			}
			w := cmd.OutOrStdout()
			fmt.Fprintf(w, "Acting as:   %s\n", me)
			for i, a := range aliases {
				if i == 0 {
					continue // the canonical name, already printed above
				}
				fmt.Fprintf(w, "Also me:     %s\n", a)
			}
			fmt.Fprintf(w, "Alias file:  %s\n", dash(path))
			if st == nil {
				return nil
			}
			fmt.Fprintf(w, "Board:       %s\n", st.Root)
			fmt.Fprintf(w, "Remote:      %s (%s)\n", st.Remote, st.Source)
			fmt.Fprintf(w, "Remotes:     %s\n", dash(strings.Join(st.Remotes, ", ")))
			if st.RefMode {
				fmt.Fprintf(w, "Ref mode:    yes — new tickets travel on their refs\n")
			} else {
				fmt.Fprintf(w, "Ref mode:    no — new tickets stay as files here\n")
				fmt.Fprintf(w, "             %s\n", st.Reason)
			}
			fmt.Fprintf(w, "File only:   %d %s on this disk and on no ref\n", st.FileOnly, plural(st.FileOnly, "ticket"))
			return nil
		},
	}
}

// boardGit is the git side of this board, gathered in one place so the text and
// the JSON cannot disagree about it.
type boardGit struct {
	Root     string
	Remote   string
	Source   string
	Remotes  []string
	RefMode  bool
	Reason   string
	FileOnly int
}

func plural(n int, word string) string {
	if n == 1 {
		return word
	}
	return word + "s"
}

// boardState collects the four facts, or nil when there is no board here.
func boardState() *boardGit {
	s, err := openStore()
	if err != nil {
		return nil
	}
	// The name comes off the repo this process actually talks to, never from a
	// second walk down the ladder: whoami exists to say which remote is used,
	// and a command that re-derived it could name a different one than the code
	// that fails. Only the step of the ladder that answered is asked for
	// separately, because the resolved Repo does not carry it.
	_, source := settings.Load().RemoteSourceFor(s.Root)
	name := gitref.DefaultRemote
	if refs != nil && refs.Repo != nil {
		name = refs.Repo.RemoteName()
	}
	st := &boardGit{
		Root:    s.Root,
		Remote:  name,
		Source:  source,
		Remotes: gitref.Remotes(s.Root),
	}
	if err := refs.Usable(); err != nil {
		st.Reason = noRefReason(err)
	} else {
		st.RefMode = true
	}
	st.FileOnly = fileOnlyCount(s)
	return st
}

// fileOnlyCount counts the tickets that exist as a file here and on no ref.
//
// Counted from the files rather than from the board, because that is the number
// that matters: on a board carrying tickets on refs a file is a ticket somebody
// pulled into work, and on a board that lost its remote every file is a ticket
// nobody else can see.
func fileOnlyCount(s *ticket.Store) int {
	paths, err := s.Paths()
	if err != nil {
		return 0
	}
	onRef := map[string]bool{}
	if refs.Usable() == nil {
		ids, err := refs.Repo.List()
		if err == nil {
			for _, id := range ids {
				onRef[id] = true
			}
		}
	}
	n := 0
	for _, p := range paths {
		if !onRef[ticket.IDFromFilename(filepath.Base(p))] {
			n++
		}
	}
	return n
}
