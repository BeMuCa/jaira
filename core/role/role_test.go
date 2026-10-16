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
	results, err := install(dir, false, true, false)
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
	// The cockpit mod goes beside the roles, where Claude Code loads a
	// plugin folder by itself.
	for _, rel := range []string{".claude-plugin/plugin.json", "hooks/register.tsx"} {
		if _, err := os.Stat(filepath.Join(dir, "jaira-cockpit", filepath.FromSlash(rel))); err != nil {
			t.Errorf("mod after install: %v", err)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "jaira-cockpit", "tests")); !os.IsNotExist(err) {
		t.Errorf("the mod's tests were installed: %v", err)
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
	if _, err := install(dir, false, true, false); err != nil {
		t.Fatal(err)
	}
	results, err := install(dir, false, true, false)
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
	if _, err := install(dir, false, true, false); err != nil {
		t.Fatal(err)
	}
	mine := filepath.Join(dir, "jaira-role-lane", skillFile)
	const edited = "---\nname: jaira-role-lane\n---\n\nmy own version\n"
	if err := os.WriteFile(mine, []byte(edited), 0o644); err != nil {
		t.Fatal(err)
	}

	results, err := install(dir, false, true, false)
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

	results, err = install(dir, true, true, false)
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
	if _, err := install(dir, false, true, false); err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(dir, "jaira-dispatcher", "scripts", "spawn.sh")

	// Same bytes, no execute bit: a plain re-run repairs it.
	if err := os.Chmod(script, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := install(dir, false, true, false); err != nil {
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
	if _, err := install(dir, true, true, false); err != nil {
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
	if _, err := install(dir, false, true, false); err != nil {
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
	a, ok := AgentNamed(DefaultAgent)
	if !ok {
		t.Fatalf("no agent named %q", DefaultAgent)
	}
	got, also := a.ProjectTarget(root)
	if want := filepath.Join(root, ".claude", "skills"); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if len(also) != 0 {
		t.Errorf("claude reads only .claude/skills, got also %v", also)
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

// The Discuss chat brings its decision back: a note, a move back that jaira
// takes for a ticket that already has its outcome, and a line to the session
// that opened it. A move line carrying --what/--why/--resolves is refused for
// such a ticket (outcome_not_yours), even with --force.
func TestSpawnDiscussReturnsItsDecision(t *testing.T) {
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skipf("no bash on this machine: %v", err)
	}
	dir := t.TempDir()
	if _, err := install(dir, false, true, false); err != nil {
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
	cmd := exec.Command(bash, spawn, "--no-worktree", "--parent", "uds:/tmp/cc-socks/7.sock", "ab12cd", "AB12CD", "discuss")
	cmd.Dir = t.TempDir()
	cmd.Env = append(os.Environ(), "TMPDIR="+t.TempDir(), "HERDR_ENV=1", "HERDR_BIN_PATH="+stub, "HERDR_WORKSPACE_ID=", "PATH="+jairaPath(t))
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("spawn.sh discuss: %v\n%s", err, out)
	}
	got, err := os.ReadFile(typed)
	if err != nil {
		t.Fatal(err)
	}
	line := chatRequest(t, string(got))
	for _, want := range []string{
		`jaira note AB12CD "decision: `,
		"jaira set AB12CD question=",
		"jaira move AB12CD --to in-progress (",
		"SendMessage an uds:/tmp/cc-socks/7.sock",
	} {
		if !strings.Contains(line, want) {
			t.Errorf("discuss prompt lacks %q:\n%s", want, line)
		}
	}
	at := strings.Index(line, "jaira move ")
	if at < 0 {
		t.Fatalf("discuss prompt has no move back:\n%s", line)
	}
	if move, _, _ := strings.Cut(line[at:], "("); strings.Contains(move, "--what") || strings.Contains(move, "--force") {
		t.Errorf("the move back must be a bare move, jaira refuses outcome flags on a ticket with an outcome: %q", move)
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
	path := jairaPath(t)
	dir := t.TempDir()
	if _, err := install(dir, false, true, false); err != nil {
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
		// A person's chats, which the cockpit mod opens: a plain request,
		// never a role, so the tab answers and never works the ticket.
		"explain": "Erklaere mir in einfachen Worten den Review-Befund zu Ticket ab12cd …",
		"discuss": "Ticket ab12cd wartet auf eine Entscheidung von mir …",
		// The cockpit's Accept, a person's press: the flag is what lets
		// jaira-role-pr open the pull request instead of handing it back.
		"pr": "/jaira-role-pr ab12cd --by-person",
	} {
		_ = os.Remove(typed)
		cmd := exec.Command(bash, spawn, "--no-worktree", "--parent", "lead", "ab12cd", "ab12cd", lane)
		cmd.Dir = t.TempDir()
		cmd.Env = append(os.Environ(), "TMPDIR="+t.TempDir(), "HERDR_ENV=1", "HERDR_BIN_PATH="+stub, "HERDR_WORKSPACE_ID=", "PATH="+path)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("spawn.sh %s: %v\n%s", lane, err, out)
		}
		got, err := os.ReadFile(typed)
		if err != nil {
			t.Fatal(err)
		}
		line := strings.TrimSpace(string(got))
		if prefix, ok := strings.CutSuffix(want, " …"); ok {
			line = chatRequest(t, line)
			if !strings.HasPrefix(line, prefix) || strings.HasPrefix(line, "/") {
				t.Errorf("lane %s typed %q, want a plain request starting %q", lane, line, prefix)
			}
		} else if line != want {
			t.Errorf("lane %s typed %q, want %q", lane, line, want)
		}
	}
}

// A pull request shows each ticket's file as the branch last committed it,
// and the lanes after the last code commit commit nothing: jaira-role-pr files
// the finished tickets away before it pushes, by handle, and describes every
// ticket the branch's commits name (GPFA7B, after PRs #62 and #63 showed
// their tickets in critique).
func TestPRFilesTheTicketsAwayBeforeThePush(t *testing.T) {
	body, err := File("jaira-role-pr", "SKILL.md")
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	pull := strings.Index(text, "jaira pull <ID>\njaira logbook <ID>")
	if pull < 0 {
		t.Error("jaira-role-pr must run jaira pull <ID> right before jaira logbook <ID>: a file older than the ref is refused")
	}
	file := strings.Index(text, "jaira logbook <ID>")
	commit := strings.Index(text, `chore(<ID>): file the ticket away`)
	push := strings.Index(text, "## Push the branch")
	if file < 0 || commit < 0 || push < 0 || file > push || commit > push {
		t.Errorf("jaira-role-pr must run jaira logbook <ID> and commit chore(<ID>): file the ticket away before ## Push the branch (at %d, %d, %d)", file, commit, push)
	}
	if strings.Contains(text, "Push what is committed and leave that file alone") || !strings.Contains(text, "Leave that file alone here: **File the finished tickets away**") {
		t.Error("jaira-role-pr's 'Before you push anything' must point at the logbook step instead of leaving the ticket file alone without exception")
	}
	if strings.Contains(text, "git add .jaira/") || !strings.Contains(text, `git add -A -- ':(glob).jaira/**/*<ID>-*.md'`) {
		t.Error("jaira-role-pr must stage the filed ticket by its own path, never .jaira/ whole")
	}
	if !strings.Contains(text, "git log --format=%s <that remote>/<its default branch>..HEAD") || strings.Index(text, "## Which repository it goes to") > file || !strings.Contains(text, "**Every ticket on the branch**") {
		t.Error("jaira-role-pr must read every ticket the branch's commits name, against the target's default branch settled before it, and describe each of them")
	}
}

// The cockpit's chats start on the chat-model setting and get the part of the
// ticket they need typed in, instead of reading the whole of it; every other
// lane starts plain claude as before. jaira and Herdr are stubs: jaira answers
// the two calls spawn.sh makes and hands 'json get' to the real one, Herdr
// records what it is asked to run and type. What goes into the brief is
// TestTicketBriefCarriesTheNewestNotes in internal/cli.
func TestSpawnChatsStartOnTheChatModelWithTheirBrief(t *testing.T) {
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skipf("no bash on this machine: %v", err)
	}
	path := jairaPath(t)
	dir := t.TempDir()
	if _, err := install(dir, false, true, false); err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	ran := filepath.Join(t.TempDir(), "ran")
	typed := filepath.Join(t.TempDir(), "typed")
	model := filepath.Join(t.TempDir(), "model")
	herdr := `#!/usr/bin/env bash
case "$1 $2" in
  "tab create") echo '{"result":{"root_pane":{"pane_id":"p1"}}}' ;;
  "pane get")   echo '{"result":{"pane":{"agent":"claude","agent_status":"idle","tab_id":"t1"}}}' ;;
  "pane run")   printf '%s\n' "$4" > "` + ran + `" ;;
  "pane send-text") printf '%s\n' "$4" > "` + typed + `" ;;
esac
`
	long := "explain: " + strings.Repeat("a finding that goes on ", 200) + "END"
	brief := `{"id":"AB12CD","title":"T","goal":"G","question":"Q?","review":{"verdict":"V","gaps":"a gap"},` +
		`"notes_last":["- **2026-10-02 10:00 · A** — second note\n  goes on","- **2026-10-05 10:00 · A** — fifth note","` + long + `"],"notes_total":5}`
	real := strings.SplitN(path, string(os.PathListSeparator), 2)[0]
	jaira := `#!/usr/bin/env bash
case "$*" in
  "json "*) exec "` + filepath.Join(real, "jaira") + `" "$@" ;;
  "-C "*" roles list --json") printf '{"chat_model": "%s"}\n' "$(cat "` + model + `")" ;;
  *" show AB12CD --brief") cat <<'JSON'
` + brief + `
JSON
  ;;
  *) exit 1 ;;
esac
`
	for name, script := range map[string]string{"herdr": herdr, "jaira": jaira} {
		if err := os.WriteFile(filepath.Join(bin, name), []byte(script), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	spawn := filepath.Join(dir, "jaira-dispatcher", "scripts", "spawn.sh")
	run := func(lane, chatModel string) (start, prompt string) {
		t.Helper()
		if err := os.WriteFile(model, []byte(chatModel), 0o644); err != nil {
			t.Fatal(err)
		}
		cmd := spawnOn(bash, "Linux", spawn, "--no-worktree", "ab12cd", "AB12CD", lane)
		cmd.Dir = t.TempDir()
		cmd.Env = append(os.Environ(), "TMPDIR="+t.TempDir(), "HERDR_ENV=1", "HERDR_BIN_PATH="+filepath.Join(bin, "herdr"),
			"HERDR_WORKSPACE_ID=", "PATH="+bin+string(os.PathListSeparator)+path)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("spawn.sh %s: %v\n%s", lane, err, out)
		}
		r, _ := os.ReadFile(ran)
		p, _ := os.ReadFile(typed)
		return strings.TrimSpace(string(r)), strings.TrimSpace(string(p))
	}

	for _, lane := range []string{"explain", "discuss"} {
		start, typedLine := run(lane, "sonnet")
		prompt := chatRequest(t, typedLine)
		// As spawn.sh spells it, slashes kept: filepath.Dir would turn them
		// on Windows, where Git Bash prints /tmp/... too.
		file := strings.TrimSuffix(strings.TrimPrefix(typedLine, chatLead), chatTail)
		dir := file[:max(strings.LastIndex(file, "/"), 0)]
		if !strings.HasSuffix(start, "&& claude --model 'sonnet' --add-dir '"+dir+"'") {
			t.Errorf("%s started %q, want claude --model 'sonnet' --add-dir '%s'", lane, start, dir)
		}
		// Typed, a long request loses its middle (42C3RZ): the brief's
		// longest note alone runs past 4000 characters, and all of it
		// reaches the file while the tab gets one short line.
		if len(typedLine) > 200 {
			t.Errorf("%s typed %d characters into the tab, want one short line: %q", lane, len(typedLine), typedLine)
		}
		if !strings.Contains(prompt, long) || len(prompt) < 4000 {
			t.Errorf("%s request in the file lost part of the brief (%d characters)", lane, len(prompt))
		}
		for _, want := range []string{`"verdict":"V"`, `"gaps":"a gap"`, `"question":"Q?"`,
			"second note\n  goes on", "fifth note", `"notes_total":5`} {
			if !strings.Contains(prompt, strings.ReplaceAll(want, "\n", `\n`)) {
				t.Errorf("%s prompt lacks %q:\n%s", lane, want, prompt)
			}
		}
		for _, not := range []string{"(jaira show"} {
			if strings.Contains(prompt, not) {
				t.Errorf("%s prompt carries %q, which it does not need:\n%s", lane, not, prompt)
			}
		}
		if strings.Contains(typedLine, "\n") {
			t.Errorf("%s typed line spans lines; a newline typed into the tab sends it early:\n%s", lane, typedLine)
		}
	}
	if start, _ := run("explain", ""); !strings.Contains(start, "&& claude --add-dir '") {
		t.Errorf("no chat-model: started %q, want plain claude, the account default", start)
	}
	if start, _ := run("explain", "x'; rm -rf ~; '"); !strings.Contains(start, "&& claude --add-dir '") {
		t.Errorf("a chat-model that is no model name reached the command line: %q", start)
	}
	if start, prompt := run("review", "sonnet"); !strings.HasSuffix(start, "&& claude") || prompt != "/jaira-role-lane AB12CD review" {
		t.Errorf("a lane must keep plain claude and its role: started %q, typed %q", start, prompt)
	}
}

// A herdr.exe is the Herdr of WSL and of native Windows alike, so spawn.sh
// tells the two apart by the machine it runs on (issue #79): in WSL the tab
// crosses back with wsl.exe --cd; in Git Bash it starts claude where it is,
// in the worktree as Windows spells it. uname and cygpath are stubbed, so both
// cases run on any machine, Windows included.
func TestSpawnStartsClaudeWhereItsShellRuns(t *testing.T) {
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skipf("no bash on this machine: %v", err)
	}
	path := jairaPath(t)
	dir := t.TempDir()
	if _, err := install(dir, false, true, false); err != nil {
		t.Fatal(err)
	}
	spawn := filepath.Join(dir, "jaira-dispatcher", "scripts", "spawn.sh")
	for _, c := range []struct {
		name, uname, distro, lane, wantCwd, wantStart string
	}{
		{"wsl", "Linux", "Ubuntu", "review", "<wt>", "wsl.exe --cd '<wt>' -- bash -lic \"claude\""},
		// cmd and PowerShell 5.1 know no &&: --cwd puts the tab in place.
		{"git bash", "MINGW64_NT-10.0-22631", "", "review", `C:\work\repo`, `claude`},
		// A chat's request file is outside the worktree: claude gets its
		// folder as Windows spells it, in the double quotes cmd and
		// PowerShell both read.
		{"git bash chat", "MINGW64_NT-10.0-22631", "", "explain", `C:\work\repo`, `claude --add-dir "C:\work\repo"`},
	} {
		bin := t.TempDir()
		ran := filepath.Join(bin, "ran")
		created := filepath.Join(bin, "created")
		herdr := `#!/usr/bin/env bash
case "$1 $2" in
  "tab create") printf '%s\n' "$@" > "` + created + `"; echo '{"result":{"root_pane":{"pane_id":"p1"}}}' ;;
  "pane get")   echo '{"result":{"pane":{"agent":"claude","agent_status":"idle","tab_id":"t1"}}}' ;;
  "pane run")   printf '%s\n' "$4" > "` + ran + `" ;;
esac
`
		if err := os.WriteFile(filepath.Join(bin, "herdr.exe"), []byte(herdr), 0o755); err != nil {
			t.Fatal(err)
		}
		wt := t.TempDir()
		cmd := spawnOn(bash, c.uname, spawn, "--no-worktree", "ab12cd", "AB12CD", c.lane, wt)
		cmd.Dir = wt
		cmd.Env = append(os.Environ(), "TMPDIR="+t.TempDir(), "HERDR_ENV=1", "HERDR_BIN_PATH="+filepath.Join(bin, "herdr.exe"),
			"HERDR_WORKSPACE_ID=", "WSL_DISTRO_NAME="+c.distro, "PATH="+bin+string(os.PathListSeparator)+path)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%s: spawn.sh: %v\n%s", c.name, err, out)
		}
		r, _ := os.ReadFile(ran)
		if got, want := strings.TrimSpace(string(r)), strings.ReplaceAll(c.wantStart, "<wt>", wt); got != want {
			t.Errorf("%s: the tab ran %q, want %q", c.name, got, want)
		}
		args, _ := os.ReadFile(created)
		if want := "--cwd\n" + strings.ReplaceAll(c.wantCwd, "<wt>", wt) + "\n"; !strings.Contains(string(args), want) {
			t.Errorf("%s: tab create got %q, want --cwd %s", c.name, args, c.wantCwd)
		}
	}
}

// spawnOn runs spawn.sh in a shell whose uname answers uname and whose cygpath
// answers C:\work\repo. They are exported functions, not stubs on PATH: on
// Windows, Git's bash.exe puts /usr/bin ahead of any PATH it is handed, so a
// stub there is never reached and the script sees the runner's own Git Bash.
func spawnOn(bash, uname, spawn string, args ...string) *exec.Cmd {
	stub := `uname() { echo '` + uname + `'; }; cygpath() { echo 'C:\work\repo'; }; export -f uname cygpath; exec "$BASH" "$@"`
	return exec.Command(bash, append([]string{"-c", stub, "bash", spawn}, args...)...)
}

// What a chat tab is typed: one line pointing at the file its request is in.
// The path is quoted: a %TEMP% with a space in it is one path, not two words.
const chatLead, chatTail = `Lies die Datei "`, `" ganz und tu, was darin steht.`

// chatRequest reads the request a chat tab was pointed at, and fails the test
// when the line typed is not such a pointer.
func chatRequest(t *testing.T, typed string) string {
	t.Helper()
	typed = strings.TrimSpace(typed)
	file, ok := strings.CutPrefix(typed, chatLead)
	file, ok2 := strings.CutSuffix(file, chatTail)
	if !ok || !ok2 {
		t.Fatalf("a chat tab must be typed %q<file>%q, got %q", chatLead, chatTail, typed)
	}
	// On Windows the script runs under Git Bash, whose mktemp prints an MSYS
	// path (/tmp/...) that Go cannot open; Git Bash's cygpath spells it as
	// Windows does.
	if runtime.GOOS == "windows" && strings.HasPrefix(file, "/") {
		bash, err := exec.LookPath("bash")
		if err != nil {
			t.Fatal(err)
		}
		out, err := exec.Command(bash, "-c", `cygpath -w "$1"`, "bash", file).Output()
		if err != nil {
			t.Fatalf("cygpath -w %s: %v", file, err)
		}
		file = strings.TrimSpace(string(out))
	}
	body, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(string(body))
}

// No role runs python3 by name: on Windows it is missing or the Store's stub
// (issue #79). The scripts read JSON with 'jaira json get', a role that runs
// a Python script asks jaira which command to run it with, and the scripts
// write UTF-8 themselves, because a Windows console would not.
func TestRolesNeverRunPython3ByName(t *testing.T) {
	roles, err := Builtins()
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range roles {
		for _, f := range r.Files {
			body, err := File(r.ID, f)
			if err != nil {
				continue
			}
			text := string(body)
			switch {
			case strings.HasSuffix(f, ".sh"):
				if strings.Contains(text, "python3 -c") || strings.Contains(text, "| python") {
					t.Errorf("%s/%s runs python to read JSON; use jaira json get", r.ID, f)
				}
			case strings.HasSuffix(f, ".md"):
				if strings.Contains(text, "python3 <") {
					t.Errorf("%s/%s runs python3 by name; ask $(jaira roles list --json | jaira json get python)", r.ID, f)
				}
			case strings.HasSuffix(f, ".py"):
				if !strings.Contains(text, `sys.stdout.reconfigure(encoding="utf-8")`) {
					t.Errorf("%s/%s does not set its stdout to UTF-8", r.ID, f)
				}
			}
		}
	}
}

// Roles reinstalled, binary not: the jaira on PATH has no 'json get', and
// spawn.sh would die after 'tab create' with an empty tab standing (#79). It
// asks first and stops before Herdr is called at all.
func TestSpawnStopsBeforeATabWhenJairaIsTooOld(t *testing.T) {
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skipf("no bash on this machine: %v", err)
	}
	dir := t.TempDir()
	if _, err := install(dir, false, true, false); err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	called := filepath.Join(bin, "called")
	for name, script := range map[string]string{
		"herdr": "#!/usr/bin/env bash\necho \"$*\" >> '" + called + "'\n",
		"jaira": "#!/usr/bin/env bash\necho 'jaira takes no arguments, received 3' >&2\nexit 2\n",
	} {
		if err := os.WriteFile(filepath.Join(bin, name), []byte(script), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	spawn := filepath.Join(dir, "jaira-dispatcher", "scripts", "spawn.sh")
	cmd := exec.Command(bash, spawn, "--no-worktree", "ab12cd", "AB12CD", "review")
	cmd.Dir = t.TempDir()
	cmd.Env = append(os.Environ(), "TMPDIR="+t.TempDir(), "HERDR_ENV=1", "HERDR_BIN_PATH="+filepath.Join(bin, "herdr"),
		"HERDR_WORKSPACE_ID=", "PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("spawn.sh went on with a jaira that has no json get:\n%s", out)
	}
	if !strings.Contains(string(out), "install the new jaira first") {
		t.Errorf("spawn.sh does not say what to do:\n%s", out)
	}
	if b, _ := os.ReadFile(called); len(b) > 0 {
		t.Errorf("Herdr was called before the check: %s", b)
	}
}
