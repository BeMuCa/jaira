package role

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// close-tab.sh and reap.sh ask one question - what in this worktree would be
// lost with it - and must give one answer. Every case runs through both, on a
// ticket that is done and merged, so the uncommitted rule is all that differs.
func TestUncommittedIsOneRuleForCloseTabAndReap(t *testing.T) {
	t.Parallel()
	// commit writes and commits a file on the ticket's branch.
	commit := func(t *testing.T, r *reapRig, rel, body string) {
		r.write(t, filepath.Join(r.wt, rel), body)
		r.git(t, r.wt, "add", rel)
		r.git(t, r.wt, "commit", "-q", "-m", "feat(AB12CD): "+rel)
	}
	for name, c := range map[string]struct {
		setup func(t *testing.T, r *reapRig)
		want  string // the paths that count; "" for none
	}{
		"clean": {func(t *testing.T, r *reapRig) {}, ""},
		// The board's refs hold ticket and milestone files in every state
		// they can still exist in.
		"ticket modified, German title": {func(t *testing.T, r *reapRig) {
			commit(t, r, ".jaira/tickets/ab12cd-größe.md", "one\n")
			r.write(t, filepath.Join(r.wt, ".jaira", "tickets", "ab12cd-größe.md"), "two\n")
		}, ""},
		"ticket untracked": {func(t *testing.T, r *reapRig) {
			r.write(t, filepath.Join(r.wt, ".jaira", "tickets", "zz99zz.md"), "x\n")
		}, ""},
		"ticket staged new": {func(t *testing.T, r *reapRig) {
			r.write(t, filepath.Join(r.wt, ".jaira", "tickets", "zz99zz.md"), "x\n")
			r.git(t, r.wt, "add", ".jaira/tickets/zz99zz.md")
		}, ""},
		"milestone staged, then modified": {func(t *testing.T, r *reapRig) {
			r.write(t, filepath.Join(r.wt, ".jaira", "milestones", "m.md"), "x\n")
			r.git(t, r.wt, "add", ".jaira/milestones/m.md")
			r.write(t, filepath.Join(r.wt, ".jaira", "milestones", "m.md"), "y\n")
		}, ""},
		// Nothing else under .jaira/ is on a ref.
		"ticket deleted": {func(t *testing.T, r *reapRig) {
			commit(t, r, ".jaira/tickets/ab12cd.md", "x\n")
			if err := os.Remove(filepath.Join(r.wt, ".jaira", "tickets", "ab12cd.md")); err != nil {
				t.Fatal(err)
			}
		}, ".jaira/tickets/ab12cd.md"},
		"lane modified": {func(t *testing.T, r *reapRig) {
			commit(t, r, ".jaira/lanes/review.md", "one\n")
			r.write(t, filepath.Join(r.wt, ".jaira", "lanes", "review.md"), "two\n")
		}, ".jaira/lanes/review.md"},
		"tags untracked": {func(t *testing.T, r *reapRig) {
			r.write(t, filepath.Join(r.wt, ".jaira", "tags"), "ui red\n")
		}, ".jaira/tags"},
		"logbook untracked": {func(t *testing.T, r *reapRig) {
			r.write(t, filepath.Join(r.wt, ".jaira", "logbook", "me-2026-11-20", "ab12cd.md"), "x\n")
		}, ".jaira/logbook/me-2026-11-20/ab12cd.md"},
		"source modified": {func(t *testing.T, r *reapRig) {
			r.write(t, filepath.Join(r.wt, "README"), "changed\n")
		}, "README"},
		"source untracked": {func(t *testing.T, r *reapRig) {
			r.write(t, filepath.Join(r.wt, "new.go"), "package x\n")
		}, "new.go"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			r, parent := newCloseTabRig(t)
			c.setup(t, r)
			r.write(t, r.shown, `{"status":"done","question":""}`)
			r.mergedAt(t, r.git(t, r.repo, "rev-parse", "feat/ab12cd"))

			tab := r.closeTab(t, "p1", "--dry-run", "--parent", parent, "AB12CD")
			reap := r.run(t, r.repo, "reap.sh", "--dry-run")
			if c.want == "" {
				if !strings.Contains(tab, "would close AB12CD") {
					t.Errorf("close-tab.sh: %q, want it closed", tab)
				}
				if !strings.Contains(reap, "would remove "+r.wt) {
					t.Errorf("reap.sh: %q, want the worktree removed", reap)
				}
				return
			}
			if !strings.Contains(tab, "kept AB12CD: uncommitted changes: "+c.want+"\n") {
				t.Errorf("close-tab.sh: %q, want kept for %q", tab, c.want)
			}
			if !strings.Contains(reap, "kept AB12CD: "+r.wt+" has uncommitted changes: "+c.want+"\n") {
				t.Errorf("reap.sh: %q, want kept for %q", reap, c.want)
			}
		})
	}
}

// Where git cannot answer, uncommitted.sh says so by its exit code: both
// scripts turn that into keeping, never into 'nothing uncommitted'.
func TestUncommittedFailsWhereGitDoes(t *testing.T) {
	t.Parallel()
	r := newReapRig(t)
	for args, want := range map[string]int{"": 2, filepath.Dir(r.repo): 1} {
		cmd := exec.Command(r.bash, filepath.Join(r.scripts, "uncommitted.sh"), args)
		cmd.Env = r.env
		err := cmd.Run()
		got := 0
		if e, ok := err.(*exec.ExitError); ok {
			got = e.ExitCode()
		}
		if got != want {
			t.Errorf("uncommitted.sh %q: exit %d, want %d", args, got, want)
		}
	}
}
