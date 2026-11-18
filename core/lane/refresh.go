package lane

import (
	"fmt"
	"os"
	"strings"

	"github.com/BeMuCa/jaira/core/shipped"
	"github.com/BeMuCa/jaira/core/ticket"
)

// Results of RefreshShipped, the same words 'jaira roles install' uses.
const (
	RefreshUnchanged = "unchanged"
	RefreshUpdated   = "updated"
	RefreshSkipped   = "skipped"
)

// Refreshed is what RefreshShipped did, or would do, to one lane file on a board.
// Board is "" for today's board and the name of a named one otherwise.
type Refreshed struct {
	ID, Board, Path, Action string
}

// RefreshShipped brings every lane file at root — today's board and every named board under
// .jaira/boards — that is an older version of a lane this binary ships up to this binary's
// version, verbatim. A file equal to the shipped one is
// unchanged; one that is no version jaira ever shipped was edited by a person and is skipped,
// untouched. Lanes jaira does not ship, and shipped lanes the board does not hold, are not
// visited: adding a lane is 'jaira lanes add'. dryRun reports and writes nothing.
// Load has its usual side effects (see Load), also in a dry run.
// The core that ticket 0S1A7M ('lanes use --all') can expose later with a force parameter.
func RefreshShipped(root string, dryRun bool) ([]Refreshed, error) {
	builtins, err := Builtins()
	if err != nil {
		return nil, err
	}
	byID := make(map[string]*Lane, len(builtins))
	for _, b := range builtins {
		byID[b.ID] = b
	}
	set, err := Load(root)
	if err != nil {
		return nil, err
	}
	type boardLane struct {
		board string
		lane  *Lane
	}
	var lanes []boardLane
	for _, l := range set.Lanes {
		lanes = append(lanes, boardLane{"", l})
	}
	for _, name := range Boards(root) {
		bs, err := LoadBoard(root, name)
		if err != nil {
			return nil, err
		}
		for _, l := range bs.Lanes {
			lanes = append(lanes, boardLane{name, l})
		}
	}
	var out []Refreshed
	for _, bl := range lanes {
		l := bl.lane
		base, ok := byID[l.ID]
		if !ok || l.Builtin {
			continue
		}
		have, err := os.ReadFile(l.Source)
		if err != nil {
			return out, err
		}
		want, err := Bytes(base)
		if err != nil {
			return out, err
		}
		r := Refreshed{ID: l.ID, Board: bl.board, Path: l.Source}
		switch {
		case string(have) == string(want):
			r.Action = RefreshUnchanged
		case shipped.Known(shippedKey(base), have):
			r.Action = RefreshUpdated
			if !dryRun {
				// Atomic: two jaira processes refreshing one board at once must
				// never leave a half-written lane file for a third to load.
				if err := ticket.WriteAtomic(l.Source, want); err != nil {
					return out, fmt.Errorf("lane %s: %w", l.ID, err)
				}
			}
		default:
			r.Action = RefreshSkipped
		}
		out = append(out, r)
	}
	return out, nil
}

// shippedKey is the path of a built-in lane in the jaira repository, the key its
// versions are listed under in core/shipped/shipped.sha256.
func shippedKey(l *Lane) string {
	return "core/lane/builtin/" + strings.TrimPrefix(l.Source, "builtin:")
}

// PreviewRefresh is the set Load(root) returns once the refresh rs reports has
// been written: every lane reported updated reads as this binary's version of
// it, at its own path. A dry run renders the agent note from it, so that what
// it says it would write is what the real run, which reloads after writing,
// does write — previewed from the lanes still on disk, a lane whose outputs
// the refresh changes made the note look changed when it would not be, or
// unchanged when it would.
func PreviewRefresh(root string, rs []Refreshed) (*Set, error) {
	set, err := Load(root)
	if err != nil {
		return nil, err
	}
	// By path, not by lane id: a named board's lane of the same id is
	// another file, and its refresh says nothing about today's board.
	updated := map[string]bool{}
	for _, r := range rs {
		if r.Action == RefreshUpdated && r.Board == "" {
			updated[r.Path] = true
		}
	}
	if len(updated) == 0 {
		return set, nil
	}
	builtins, err := Builtins()
	if err != nil {
		return nil, err
	}
	lanes := make([]*Lane, 0, len(set.Lanes))
	for _, l := range set.Lanes {
		if path := l.Source; updated[path] {
			for _, b := range builtins {
				if b.ID != l.ID {
					continue
				}
				want, err := Bytes(b)
				if err != nil {
					return nil, err
				}
				if l, err = parse(want, path, false); err != nil {
					return nil, fmt.Errorf("lane %s: %w", b.ID, err)
				}
			}
		}
		lanes = append(lanes, l)
	}
	orderIDs, _ := LoadOrder(root)
	return assemble(lanes, orderIDs, nil), nil
}
