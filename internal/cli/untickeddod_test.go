package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// Leaving the implementing lane with not one definition-of-done item ticked
// warns on stderr — under --json too, where stdout stays the data — and the
// move lands; with an item ticked it does not warn (77B4GA).
func TestMoveOutOfImplementingWarnsOnUntickedDoD(t *testing.T) {
	for _, tick := range []bool{false, true} {
		for _, asJSON := range []bool{false, true} {
			dir, id := movableTicket(t)
			runSplit(t, dir, "move", id, "--to", "todo")
			runSplit(t, dir, "move", id, "--to", "in-progress")
			if tick {
				runSplit(t, dir, "dod", id, "1", "--done", "--proof", "x_test.go")
			}
			args := []string{"move", id, "--to", "review", "--what", "w", "--why", "y", "--resolves", "r"}
			if asJSON {
				args = append(args, "--json")
			}
			stdout, stderr := runSplit(t, dir, args...)

			warned := strings.Contains(stderr, "Definition of Done: 0 of 1 ticked")
			if warned == tick {
				t.Errorf("ticked=%v json=%v: warned=%v\nstderr:\n%s", tick, asJSON, warned, stderr)
			}
			if strings.Contains(stdout, "Definition of Done") {
				t.Errorf("ticked=%v json=%v: the warning reached stdout:\n%s", tick, asJSON, stdout)
			}
			if asJSON {
				var got struct {
					Moved bool `json:"moved"`
				}
				if err := json.Unmarshal([]byte(stdout), &got); err != nil || !got.Moved {
					t.Errorf("ticked=%v: stdout is not the move's JSON (%v):\n%s", tick, err, stdout)
				}
			} else if !strings.Contains(stdout, "→ review") {
				t.Errorf("ticked=%v: the move did not land:\n%s", tick, stdout)
			}
		}
	}
}
