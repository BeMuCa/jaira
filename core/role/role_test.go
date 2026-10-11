package role

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
)

// The roles this binary ships. Named rather than counted, because the
// point of the list is that every one of them arrives, not that a number of
// something did.
var wantRoles = []string{
	"jaira-dispatcher",
	"jaira-role-acceptance",
	"jaira-role-brainstorm",
	"jaira-role-lane",
	"jaira-role-logbook-summary",
	"jaira-role-pr",
	"jaira-role-research",
	"jaira-role-tester",
	"jaira-teamlead",
	"jaira-worker",
}

func TestBuiltinsAreTheShippedRoles(t *testing.T) {
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

// A role's supporting files travel with it: a prompt that calls a script or
// reads a template by a path relative to its own folder is a broken instruction
// when the file stays behind.
var roleFiles = map[string][]string{
	// run-lane.sh is what the dispatcher starts in the background to wait out
	// a lane.
	"jaira-dispatcher": {"scripts/spawn.sh", "scripts/run-lane.sh"},
	// build.py writes the page from template.html and the record from
	// report.html.
	"jaira-role-acceptance": {"build.py", "template.html", "report.html"},
	// collect.py gathers the facts, build.py checks the text and writes the
	// page from template.html; summary.example.md is what the prompt points a
	// person at to shape the summary.
	"jaira-role-logbook-summary": {"collect.py", "build.py", "template.html", "summary.example.md"},
}

func TestRolesShipTheirFiles(t *testing.T) {
	roles, err := Builtins()
	if err != nil {
		t.Fatal(err)
	}
	byID := map[string]Role{}
	for _, r := range roles {
		byID[r.ID] = r
	}
	for id, want := range roleFiles {
		r, ok := byID[id]
		if !ok {
			t.Errorf("no %s role", id)
			continue
		}
		for _, w := range want {
			if !slices.Contains(r.Files, w) {
				t.Errorf("%s: files %v, want %s among them", id, r.Files, w)
			}
		}
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
					// Every occurrence, not just the first: the first hit is
					// often the tail of an already-prefixed name, and stopping
					// there hides a bare reference later in the same line.
					for at := 0; ; {
						idx := strings.Index(line[at:], bare)
						if idx < 0 {
							break
						}
						idx += at
						at = idx + len(bare)
						if idx >= 6 && strings.HasSuffix(line[:idx], "/jaira") {
							continue // already prefixed
						}
						t.Errorf("%s/%s references %s without the jaira- prefix: %s", r.ID, rel, bare, line)
					}
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
	// a broken instruction. Windows has no POSIX execute bit — os.Chmod there
	// only toggles the read-only flag — so the claim is not one this platform
	// can make either way.
	if runtime.GOOS != "windows" {
		for _, rel := range roleFiles["jaira-dispatcher"] {
			fi, err := os.Stat(filepath.Join(dir, "jaira-dispatcher", filepath.FromSlash(rel)))
			if err != nil {
				t.Fatal(err)
			}
			if fi.Mode().Perm()&0o111 == 0 {
				t.Errorf("%s installed as %v, want the execute bit", rel, fi.Mode().Perm())
			}
		}
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

// A spawn.sh that arrived without its execute bit — copied by hand, unzipped,
// checked out on a filesystem that drops the bit — must come back executable.
// os.WriteFile hands the mode to O_CREATE only, so an existing file keeps the
// permissions it had and the repair has to be stated separately.
func TestInstallRestoresTheExecuteBit(t *testing.T) {
	// Windows carries no POSIX permission bits: os.Chmod there only flips the
	// read-only flag, so neither the damage this repairs nor the repair itself
	// can be staged on that platform.
	if runtime.GOOS == "windows" {
		t.Skip("no POSIX execute bit on windows")
	}
	dir := t.TempDir()
	if _, err := Install(dir, false); err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(dir, "jaira-dispatcher", "scripts", "spawn.sh")

	// Same bytes, no execute bit: a plain re-run repairs it.
	if err := os.Chmod(script, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Install(dir, false); err != nil {
		t.Fatal(err)
	}
	if fi, err := os.Stat(script); err != nil {
		t.Fatal(err)
	} else if fi.Mode().Perm() != 0o755 {
		t.Errorf("after a re-run the script is %v, want 0755", fi.Mode().Perm())
	}

	// Edited bytes and no execute bit: --force restores both.
	if err := os.WriteFile(script, []byte("#!/bin/sh\necho mine\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Install(dir, true); err != nil {
		t.Fatal(err)
	}
	fi, err := os.Stat(script)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm() != 0o755 {
		t.Errorf("after --force the script is %v, want 0755", fi.Mode().Perm())
	}
	want, _ := File("jaira-dispatcher", "scripts/spawn.sh")
	if b, _ := os.ReadFile(script); string(b) != string(want) {
		t.Error("--force did not restore the embedded script")
	}
}

// A shipped script is only ever run by a worker, in a pane nobody is watching
// at the moment it breaks. Parse every one of them here instead, and prove the
// dispatcher's own help path does not need Herdr to answer — a caller asking
// what the flags are has no pane yet.
func TestEmbeddedScriptsParse(t *testing.T) {
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skipf("no bash on this machine: %v", err)
	}
	dir := t.TempDir()
	if _, err := Install(dir, false); err != nil {
		t.Fatal(err)
	}

	roles, err := Builtins()
	if err != nil {
		t.Fatal(err)
	}
	scripts := 0
	for _, r := range roles {
		for _, rel := range r.Files {
			if filepath.Ext(rel) != ".sh" {
				continue
			}
			scripts++
			p := filepath.Join(dir, r.ID, filepath.FromSlash(rel))
			if out, err := exec.Command(bash, "-n", p).CombinedOutput(); err != nil {
				t.Errorf("%s/%s does not parse: %v\n%s", r.ID, rel, err, out)
			}
		}
	}
	// Guards the loop itself: a change that stops reporting supporting files
	// would leave every assertion above unrun and the test still green.
	if scripts == 0 {
		t.Fatal("no embedded .sh file was checked")
	}

	// HERDR_ENV unset on purpose: --help is answered before the pane check, so
	// somebody can read the flags on a machine that has no Herdr at all.
	spawn := filepath.Join(dir, "jaira-dispatcher", "scripts", "spawn.sh")
	cmd := exec.Command(bash, spawn, "--help")
	cmd.Env = append(os.Environ(), "HERDR_ENV=")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("spawn.sh --help: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "--no-worktree") {
		t.Errorf("spawn.sh --help does not mention --no-worktree:\n%s", out)
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
	got := ProjectTarget(root)
	if want := filepath.Join(root, ".claude", "skills"); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

// The worker uses the dispatcher's scripts rather than carrying copies that
// drift: both roles are installed side by side, so the paths it names are
// beside it. A rename of either file has to fail here, not in a pane.
func TestWorkerNamesScriptsThatShip(t *testing.T) {
	b, err := File("jaira-worker", skillFile)
	if err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{"scripts/spawn.sh", "scripts/run-lane.sh"} {
		if !strings.Contains(string(b), "jaira-dispatcher/"+rel) {
			t.Errorf("jaira-worker does not name jaira-dispatcher/%s", rel)
		}
		if _, err := File("jaira-dispatcher", rel); err != nil {
			t.Errorf("jaira-dispatcher/%s does not ship: %v", rel, err)
		}
	}
}

// spawn.sh types the role into the tab it opened; "worker" must start
// /jaira-worker and "dispatch" /jaira-dispatcher, --parent carried along.
// Herdr is replaced by a stub that answers the three calls the script reads
// and records what was typed — no tab is opened anywhere.
func TestSpawnTypesTheRoleForItsLaneName(t *testing.T) {
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skipf("no bash on this machine: %v", err)
	}
	if _, err := exec.LookPath("python3"); err != nil {
		t.Skipf("no python3 on this machine: %v", err)
	}
	dir := t.TempDir()
	if _, err := Install(dir, false); err != nil {
		t.Fatal(err)
	}
	typed := filepath.Join(t.TempDir(), "typed")
	stub := filepath.Join(t.TempDir(), "herdr")
	script := `#!/usr/bin/env bash
case "$1 $2" in
  "tab create") echo '{"result":{"root_pane":{"pane_id":"p1"}}}' ;;
  "pane get")   echo '{"result":{"pane":{"agent":"claude","agent_status":"idle","tab_id":"t1"}}}' ;;
  "pane send-text") printf '%s\n' "$4" >> "` + typed + `" ;;
esac
`
	if err := os.WriteFile(stub, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	spawn := filepath.Join(dir, "jaira-dispatcher", "scripts", "spawn.sh")
	for lane, want := range map[string]string{
		"worker":   "/jaira-worker ab12cd --parent lead",
		"dispatch": "/jaira-dispatcher ab12cd --parent lead",
		"review":   "/jaira-role-lane ab12cd review",
	} {
		_ = os.Remove(typed)
		cmd := exec.Command(bash, spawn, "--no-worktree", "--parent", "lead", "ab12cd", "ab12cd", lane)
		cmd.Dir = t.TempDir()
		cmd.Env = append(os.Environ(), "HERDR_ENV=1", "HERDR_BIN_PATH="+stub, "HERDR_WORKSPACE_ID=")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("spawn.sh %s: %v\n%s", lane, err, out)
		}
		got, err := os.ReadFile(typed)
		if err != nil {
			t.Fatal(err)
		}
		if strings.TrimSpace(string(got)) != want {
			t.Errorf("lane %s typed %q, want %q", lane, strings.TrimSpace(string(got)), want)
		}
	}
}
