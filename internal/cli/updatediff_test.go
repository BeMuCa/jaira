package cli

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// treeState is every file under dirs with its content, so a test can tell whether a
// command wrote anything at all.
func treeState(t *testing.T, dirs ...string) map[string]string {
	t.Helper()
	out := map[string]string{}
	for _, d := range dirs {
		err := filepath.WalkDir(d, func(p string, e fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if e.IsDir() {
				return nil
			}
			b, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			out[p] = string(b)
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	return out
}

// linesWith returns the lines of out that name file and carry state.
func linesWith(out, file, state string) []string {
	var got []string
	for _, l := range strings.Split(out, "\n") {
		if strings.Contains(l, file) && strings.Contains(l, state) {
			got = append(got, l)
		}
	}
	return got
}

func TestUpdateDiffWritesNothingAndNamesWhatIsBehind(t *testing.T) {
	st := newStaleSetup(t)
	before := treeState(t, st.dir, st.home)

	out, err := runUpdate(t, st.dir, "diff")
	if err != nil {
		t.Fatalf("update diff: %v\n%s", err, out)
	}
	after := treeState(t, st.dir, st.home)
	if len(before) != len(after) {
		t.Errorf("update diff changed the number of files: %d -> %d", len(before), len(after))
	}
	for p, b := range before {
		if after[p] != b {
			t.Errorf("update diff wrote %s", p)
		}
	}
	if len(linesWith(out, "jaira-role-acceptance/SKILL.md", "behind (")) == 0 {
		t.Errorf("output does not list the stale role as behind:\n%s", out)
	}
	if len(linesWith(out, "testing", "behind (")) == 0 {
		t.Errorf("output does not list the stale testing lane as behind:\n%s", out)
	}
}

func TestUpdateDiffNamesAnEditedRoleAsEdited(t *testing.T) {
	st := newStaleSetup(t)
	mustWrite(t, st.projRole, []byte("my own words\n"))

	out, err := runUpdate(t, st.dir, "diff")
	if err != nil {
		t.Fatalf("update diff: %v\n%s", err, out)
	}
	if len(linesWith(out, "jaira-role-acceptance/SKILL.md", "edited")) == 0 {
		t.Errorf("output does not list the edited role as edited:\n%s", out)
	}
}

func TestUpdateDiffJSONCarriesEveryFileAndWritesNothing(t *testing.T) {
	st := newStaleSetup(t)

	out, err := runUpdate(t, st.dir, "diff", "--json")
	if err != nil {
		t.Fatalf("update diff --json: %v\n%s", err, out)
	}
	var payload struct {
		Roles          []map[string]any `json:"roles"`
		Lanes          []map[string]any `json:"lanes"`
		NothingWritten bool             `json:"nothing_written"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("stdout is not JSON: %v\n%s", err, out)
	}
	if len(payload.Roles) == 0 || len(payload.Lanes) == 0 || !payload.NothingWritten {
		t.Errorf("payload lacks roles, lanes or nothing_written: %s", out)
	}
}

func TestUpdateDiffOutsideABoardListsPersonalRolesOnly(t *testing.T) {
	newStaleSetup(t)
	out, err := runUpdate(t, t.TempDir(), "diff")
	if err != nil {
		t.Fatalf("update diff outside a board: %v\n%s", err, out)
	}
	if !hasLine(out, "no board here") {
		t.Errorf("output lacks the no-board line:\n%s", out)
	}
}
