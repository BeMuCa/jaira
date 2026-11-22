package role

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// land.sh lands an accepted pull request: merge the base, regenerate what
// .jaira/land-regenerate names, push, wait for CI, and merge only where the
// base's .jaira/land says agent and every ticket is in done — else it hands
// the merge line back (50RA6N). git is real, origin a bare repository whose
// path reads like GitHub; gh is a stub.
func TestLandMergesOnlyWhereTheBoardLetsIt(t *testing.T) {
	cases := []struct {
		name      string
		land      string // .jaira/land on the base; "" none
		status    string // the ticket's lane
		conflict  bool   // master and the branch both change other.txt
		red       bool   // the pull request's checks fail
		dirty     bool   // a tracked file is changed and not committed
		noRegen   bool   // the base has no .jaira/land-regenerate
		upToDate  bool   // master does not move: nothing to merge or push
		wantExit  int
		wantOut   string
		wantMerge bool
		wantPush  bool // the branch on origin carries master's tip afterwards
	}{
		{name: "no land file", status: "done", wantExit: 3, wantOut: "gh pr merge 7 --repo o/r --merge", wantPush: true},
		{name: "land person", land: "person", status: "done", wantExit: 3, wantOut: "gh pr merge 7 --repo o/r --merge", wantPush: true},
		{name: "agent, ticket in signoff", land: "agent", status: "signoff", wantExit: 3, wantOut: "(signoff)", wantPush: true},
		{name: "agent, ticket done", land: "agent", status: "done", wantExit: 0, wantOut: "merged", wantMerge: true, wantPush: true},
		{name: "conflict outside the generated file", land: "agent", status: "done", conflict: true, wantExit: 1, wantOut: "other.txt"},
		{name: "red CI", land: "agent", status: "done", red: true, wantExit: 1, wantOut: "not green", wantPush: true},
		{name: "no generated files", land: "agent", status: "done", noRegen: true, wantExit: 0, wantOut: "merged", wantMerge: true, wantPush: true},
		{name: "branch already holds the base", land: "agent", status: "done", upToDate: true, wantExit: 0, wantOut: "merged", wantMerge: true, wantPush: true},
		{name: "uncommitted change", land: "agent", status: "done", dirty: true, wantExit: 2, wantOut: "uncommitted"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := newLandRig(t, c.land, !c.noRegen)
			h := r.ticket(t, c.status)
			wt := r.clone(t, h)
			r.write(t, filepath.Join(wt, "gen.txt"), "branch\n")
			if c.conflict {
				r.write(t, filepath.Join(wt, "other.txt"), "branch\n")
			}
			r.git(t, wt, "commit", "-qam", "feat("+h+"): the work")
			r.git(t, wt, "push", "-q", "origin", "HEAD")
			if !c.upToDate {
				r.advanceMaster(t, c.conflict, !c.noRegen)
			}
			if c.red {
				r.write(t, r.checksFail, "")
			}
			if c.dirty {
				r.write(t, filepath.Join(wt, "other.txt"), "dirty\n")
			}

			out, code := r.land(t, wt)
			if code != c.wantExit || !strings.Contains(out, c.wantOut) {
				t.Fatalf("exit %d, want %d with %q:\n%s", code, c.wantExit, c.wantOut, out)
			}
			_, err := os.Stat(r.merged)
			if merged := err == nil; merged != c.wantMerge {
				t.Errorf("merged = %v, want %v\n%s", merged, c.wantMerge, out)
			}
			branch := "feat/" + strings.ToLower(h)
			pushed := exec.Command("git", "-C", r.origin, "merge-base", "--is-ancestor", "master", branch).Run() == nil
			if pushed != c.wantPush {
				t.Errorf("master in origin's %s = %v, want %v\n%s", branch, pushed, c.wantPush, out)
			}
			if c.wantPush && c.noRegen {
				if got := r.git(t, r.origin, "log", "-1", "--format=%s", branch); got != "chore("+h+"): merge origin/master" {
					t.Errorf("last commit on the branch = %q", got)
				}
			}
			if c.wantPush && c.upToDate {
				// Nothing to merge: the regeneration alone, under its own name.
				if got := r.git(t, r.origin, "log", "-1", "--format=%s", branch); got != "chore("+h+"): regenerate gen.txt" {
					t.Errorf("last commit on an up-to-date branch = %q", got)
				}
			}
			if c.wantPush && !c.noRegen && !c.upToDate {
				// The generated file is the generator's, and the merge and the
				// regeneration ride in one commit naming the ticket.
				if got := r.git(t, r.origin, "show", branch+":gen.txt"); got != "regenerated" {
					t.Errorf("gen.txt on the branch = %q, want regenerated", got)
				}
				if got := r.git(t, r.origin, "log", "-1", "--format=%s", branch); !strings.HasPrefix(got, "chore("+h+"): merge origin/master, regenerate gen.txt") {
					t.Errorf("last commit on the branch = %q", got)
				}
			}
			if c.conflict {
				if st := r.git(t, wt, "status", "--porcelain", "--untracked-files=no"); st != "" {
					t.Errorf("a stopped merge left the worktree dirty:\n%s", st)
				}
			}
		})
	}
}

type landRig struct {
	tmp, origin, skills, merged, checksFail, boardDir string
	env                                               []string
	at                                                time.Time
}

func newLandRig(t *testing.T, land string, regen bool) *landRig {
	t.Helper()
	if _, err := exec.LookPath("bash"); err != nil || runtime.GOOS == "windows" {
		t.Skip("land.sh needs bash and POSIX paths")
	}
	path := jairaPath(t)
	tmp, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	r := &landRig{
		tmp:        tmp,
		origin:     filepath.Join(tmp, "github.com", "o", "r.git"),
		skills:     filepath.Join(tmp, "skills"),
		merged:     filepath.Join(tmp, "merged"),
		checksFail: filepath.Join(tmp, "checks-fail"),
		at:         time.Now(),
	}
	if _, err := install(r.skills, false, true, false); err != nil {
		t.Fatal(err)
	}
	bin := filepath.Join(tmp, "bin")
	gh := `#!/usr/bin/env bash
case "$*" in
  "repo view ` + r.origin + ` --json isFork,parent,nameWithOwner") echo '{"isFork":false,"nameWithOwner":"o/r"}' ;;
  "pr view 7 --repo o/r --json state,headRefName,baseRefName")
    echo "{\"state\":\"OPEN\",\"headRefName\":\"$(git symbolic-ref --short HEAD)\",\"baseRefName\":\"master\"}" ;;
  "pr checks 7 --repo o/r"*) [ -e "` + r.checksFail + `" ] && { printf 'test\tfail\n'; exit 1; }; printf 'test\tpass\n' ;;
  "pr merge 7 --repo o/r --merge") echo "$*" >> "` + r.merged + `" ;;
  *) exit 1 ;;
esac
`
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	r.write(t, filepath.Join(bin, "gh"), gh)
	if err := os.Chmod(filepath.Join(bin, "gh"), 0o755); err != nil {
		t.Fatal(err)
	}
	r.env = append(os.Environ(), "JAIRA_LAND_SETTLE=0", "JAIRA_NO_FETCH=1", "JAIRA_NO_SNAPSHOT=1", "JAIRA_NO_UPDATE_CHECK=1", "HERDR_ENV=", "HC_TERMINAL_ID=",
		"JAIRA_USER=t", "JAIRA_HOME="+filepath.Join(tmp, "home"),
		"PATH="+bin+string(os.PathListSeparator)+path)

	// The base: a generated file, another file, and the board's land files.
	seed := filepath.Join(tmp, "seed")
	r.git(t, tmp, "init", "-q", "--bare", "-b", "master", r.origin)
	r.git(t, tmp, "clone", "-q", r.origin, seed)
	r.git(t, seed, "config", "user.email", "t@example.com")
	r.git(t, seed, "config", "user.name", "t")
	r.write(t, filepath.Join(seed, "gen.txt"), "base\n")
	r.write(t, filepath.Join(seed, "other.txt"), "base\n")
	if regen {
		r.write(t, filepath.Join(seed, ".jaira", "land-regenerate"), "# generated files\n\ngen.txt: echo regenerated > gen.txt\n")
	}
	if land != "" {
		r.write(t, filepath.Join(seed, ".jaira", "land"), land+"\n")
	}
	r.git(t, seed, "add", "-Af")
	r.git(t, seed, "commit", "-qm", "chore: base")
	r.git(t, seed, "push", "-q", "origin", "master")
	return r
}

// ticket files a ticket in a lane on the board the worktree reads and
// returns its handle. The board itself is not committed: tickets stay on
// disk, as on a board not shared yet.
func (r *landRig) ticket(t *testing.T, status string) string {
	t.Helper()
	r.boardDir = filepath.Join(r.tmp, "src", "jaira-wt")
	if err := os.MkdirAll(r.boardDir, 0o755); err != nil {
		t.Fatal(err)
	}
	r.git(t, r.tmp, "clone", "-q", r.origin, r.boardDir)
	s, err := ticket.At(r.boardDir)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("JAIRA_USER", "t")
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	r.at = r.at.Add(time.Millisecond)
	tk, err := s.Create(map[string]string{
		ticket.FieldID: ticket.NewID(r.at), ticket.FieldTitle: "t", ticket.FieldStatus: status,
	}, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	return ticket.Handle(tk.ID)
}

// clone puts the board's clone on the ticket's branch.
func (r *landRig) clone(t *testing.T, h string) string {
	t.Helper()
	wt := r.boardDir
	r.git(t, wt, "config", "user.email", "t@example.com")
	r.git(t, wt, "config", "user.name", "t")
	r.git(t, wt, "checkout", "-qb", "feat/"+strings.ToLower(h))
	return wt
}

// advanceMaster lands another commit on origin's master: a new file, with
// gen also the generated file — a conflict regeneration resolves — and with
// conflict also other.txt, which nothing resolves.
func (r *landRig) advanceMaster(t *testing.T, conflict, gen bool) {
	t.Helper()
	seed := filepath.Join(r.tmp, "seed")
	r.git(t, seed, "pull", "-q", "origin", "master")
	if gen {
		r.write(t, filepath.Join(seed, "gen.txt"), "master moved\n")
	}
	if conflict {
		r.write(t, filepath.Join(seed, "other.txt"), "master moved\n")
	}
	r.write(t, filepath.Join(seed, "third.txt"), "new on master\n")
	r.git(t, seed, "add", "-A")
	r.git(t, seed, "commit", "-qm", "chore: master moves on")
	r.git(t, seed, "push", "-q", "origin", "master")
}

func (r *landRig) land(t *testing.T, wt string) (string, int) {
	t.Helper()
	cmd := exec.Command(filepath.Join(r.skills, "jaira-teamlead", "scripts", "land.sh"), "7")
	cmd.Dir = wt
	cmd.Env = r.env
	out, err := cmd.CombinedOutput()
	if ee, ok := err.(*exec.ExitError); ok {
		return string(out), ee.ExitCode()
	}
	if err != nil {
		t.Fatal(err)
	}
	return string(out), 0
}

func (r *landRig) git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL=/dev/null", "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.com",
		"GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.com")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func (r *landRig) write(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}
