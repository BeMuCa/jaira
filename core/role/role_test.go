package role

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The seven roles this ticket ships. Named rather than counted, because the
// point of the list is that every one of them arrives, not that seven of
// something did.
var wantRoles = []string{
	"jaira-dispatcher",
	"jaira-role-brainstorm",
	"jaira-role-lane",
	"jaira-role-pr",
	"jaira-role-research",
	"jaira-role-tester",
	"jaira-teamlead",
}

func TestBuiltinsAreTheSevenRoles(t *testing.T) {
	roles, err := Builtins()
	if err != nil {
		t.Fatal(err)
	}
	if len(roles) != len(wantRoles) {
		t.Fatalf("got %d roles, want %d", len(roles), len(wantRoles))
	}
	for i, r := range roles {
		if r.ID != wantRoles[i] {
			t.Errorf("role %d: got id %q, want %q", i, r.ID, wantRoles[i])
		}
		// The directory name is what the harness invokes, but Claude Code
		// refuses to load a skill whose frontmatter name disagrees with it.
		// Nothing in the package reads the field back, so this is the only
		// thing keeping the two in step in the embedded files.
		b, err := File(r.ID, skillFile)
		if err != nil {
			t.Fatalf("%s: %v", r.ID, err)
		}
		if want := "\nname: " + r.ID + "\n"; !strings.Contains(string(b), want) {
			t.Errorf("%s: %s frontmatter does not carry %q", r.ID, skillFile, strings.TrimSpace(want))
		}
		if r.Description == "" {
			t.Errorf("%s: no description", r.ID)
		}
		if len(r.Files) == 0 || r.Files[0] != skillFile {
			t.Errorf("%s: files %v, want %s first", r.ID, r.Files, skillFile)
		}
	}
}

// A role's supporting files travel with it: teamlead references a script, and
// shipping the prompt without the script ships a broken instruction.
func TestTeamleadShipsItsScript(t *testing.T) {
	roles, err := Builtins()
	if err != nil {
		t.Fatal(err)
	}
	var r Role
	for _, b := range roles {
		if b.ID == "jaira-teamlead" {
			r = b
		}
	}
	if r.ID == "" {
		t.Fatal("no jaira-teamlead role")
	}
	found := false
	for _, f := range r.Files {
		if f == "scripts/spawn.sh" {
			found = true
		}
	}
	if !found {
		t.Fatalf("files %v, want scripts/spawn.sh among them", r.Files)
	}
}

// The prefix is not cosmetic: a cross-reference left on the unprefixed name
// calls a command that does not exist on a teammate's machine.
func TestCrossReferencesCarryThePrefix(t *testing.T) {
	roles, err := Builtins()
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range roles {
		for _, rel := range r.Files {
			b, err := File(r.ID, rel)
			if err != nil {
				t.Fatal(err)
			}
			for _, other := range wantRoles {
				bare := "/" + strings.TrimPrefix(other, "jaira-")
				for _, line := range strings.Split(string(b), "\n") {
					idx := strings.Index(line, bare)
					if idx < 0 {
						continue
					}
					if idx >= 6 && strings.HasSuffix(line[:idx], "/jaira") {
						continue // already prefixed
					}
					t.Errorf("%s/%s references %s without the jaira- prefix: %s", r.ID, rel, bare, line)
				}
			}
		}
	}
}

func TestInstallWritesEveryFile(t *testing.T) {
	dir := t.TempDir()
	results, err := Install(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	for _, res := range results {
		if res.Action != Written {
			t.Errorf("%s: got %s, want %s on a first install", res.Path, res.Action, Written)
		}
	}
	for _, id := range wantRoles {
		p := filepath.Join(dir, id, skillFile)
		if _, err := os.Stat(p); err != nil {
			t.Errorf("after install: %v", err)
		}
	}
	// A shipped script must arrive executable, or the prompt that calls it is
	// a broken instruction.
	fi, err := os.Stat(filepath.Join(dir, "jaira-teamlead", "scripts", "spawn.sh"))
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm()&0o111 == 0 {
		t.Errorf("spawn.sh installed as %v, want the execute bit", fi.Mode().Perm())
	}
}

func TestSecondRunChangesNothing(t *testing.T) {
	dir := t.TempDir()
	if _, err := Install(dir, false); err != nil {
		t.Fatal(err)
	}
	results, err := Install(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	for _, res := range results {
		if res.Action != Unchanged {
			t.Errorf("%s: got %s on a second run, want %s", res.Path, res.Action, Unchanged)
		}
	}
	if SkippedAny(results) {
		t.Error("a second run reported a skip; an already-installed file is not an edited one")
	}
}

func TestEditedFileIsLeftAloneAndReported(t *testing.T) {
	dir := t.TempDir()
	if _, err := Install(dir, false); err != nil {
		t.Fatal(err)
	}
	mine := filepath.Join(dir, "jaira-role-lane", skillFile)
	const edited = "---\nname: jaira-role-lane\n---\n\nmy own version\n"
	if err := os.WriteFile(mine, []byte(edited), 0o644); err != nil {
		t.Fatal(err)
	}

	results, err := Install(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	if !SkippedAny(results) {
		t.Fatal("an edited file was not reported as skipped")
	}
	if got := actionFor(results, mine); got != Skipped {
		t.Errorf("got %s for the edited file, want %s", got, Skipped)
	}
	if b, _ := os.ReadFile(mine); string(b) != edited {
		t.Error("the edited file was overwritten without --force")
	}

	results, err = Install(dir, true)
	if err != nil {
		t.Fatal(err)
	}
	if got := actionFor(results, mine); got != Overwritten {
		t.Errorf("with --force got %s, want %s", got, Overwritten)
	}
	want, _ := File("jaira-role-lane", skillFile)
	if b, _ := os.ReadFile(mine); string(b) != string(want) {
		t.Error("--force did not restore the embedded version")
	}
}

func actionFor(results []Result, path string) Action {
	for _, r := range results {
		if r.Path == path {
			return r.Action
		}
	}
	return ""
}

func TestProjectTargetIsTheClaudeSkillsDirectory(t *testing.T) {
	root := t.TempDir()
	// One directory, whether or not other agent directories are around: a
	// harness that reads two of them would see the same role registered twice.
	for _, d := range []string{".codex", ".agents"} {
		if err := os.Mkdir(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	got := ProjectTarget(root)
	if want := filepath.Join(root, ".claude", "skills"); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
