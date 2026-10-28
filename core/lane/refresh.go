package lane

import (
	"fmt"
	"os"
	"strings"

	"github.com/BeMuCa/jaira/core/shipped"
)

// Results of RefreshShipped, the same words 'jaira roles install' uses.
const (
	RefreshUnchanged = "unchanged"
	RefreshUpdated   = "updated"
	RefreshSkipped   = "skipped"
)

// Refreshed is what RefreshShipped did, or would do, to one lane file on a board.
type Refreshed struct {
	ID, Path, Action string
}

// RefreshShipped brings every lane file on the board at root that is an older version of a lane
// this binary ships up to this binary's version, verbatim. A file equal to the shipped one is
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
	var out []Refreshed
	for _, l := range set.Lanes {
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
		r := Refreshed{ID: l.ID, Path: l.Source}
		switch {
		case string(have) == string(want):
			r.Action = RefreshUnchanged
		case shipped.Known(shippedKey(base), have):
			r.Action = RefreshUpdated
			if !dryRun {
				if err := os.WriteFile(l.Source, want, 0o644); err != nil {
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
