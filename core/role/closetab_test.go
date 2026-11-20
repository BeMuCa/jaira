package role

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// The lanes close-tab.sh tells apart: an agent's, the one that takes a
// question, the one only a person moves out of, and the terminal one.
const closeTabLanes = `{"lanes":[
 {"id":"in-progress","human_exit":false,"requires_question":false,"terminal":false},
 {"id":"human","human_exit":false,"requires_question":true,"terminal":false},
 {"id":"signoff","human_exit":true,"requires_question":false,"terminal":false},
 {"id":"done","human_exit":false,"requires_question":false,"terminal":true}]}`

// closeTab runs close-tab.sh in the ticket's worktree as the worker's tab p1
// (labelled AB12CD/worker by spawn.sh), with a parent that listens unless
// parent is given.
func (r *reapRig) closeTab(t *testing.T, pane string, args ...string) string {
	t.Helper()
	cmd := exec.Command(r.bash, append([]string{filepath.Join(r.scripts, "close-tab.sh")}, args...)...)
	cmd.Dir = r.wt
	cmd.Env = append(append([]string{}, r.env...), "HERDR_PANE_ID="+pane, "HC_TERMINAL_ID=")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("close-tab.sh %v: %v\n%s", args, err, out)
	}
	return string(out)
}

func newCloseTabRig(t *testing.T) (*reapRig, string) {
	t.Helper()
	r := newReapRig(t)
	r.write(t, r.lanes, closeTabLanes)
	sock := filepath.Join(filepath.Dir(r.closed), "parent.sock")
	r.write(t, sock, "")
	return r, "uds:" + sock
}

// A worker whose ticket waits for a person — in signoff, in human with the
// question gone, or done — closes its own tab once its report went to the
// parent. Board files under .jaira/ are no reason to stay.
func TestCloseTabClosesAWorkerWhoseTicketIsWithAPerson(t *testing.T) {
	t.Parallel()
	for _, lane := range []string{"signoff", "human", "done"} {
		t.Run(lane, func(t *testing.T) {
			t.Parallel()
			r, parent := newCloseTabRig(t)
			r.write(t, r.shown, `{"status":"`+lane+`","question":""}`)
			// Ticket files are named after their German titles: git quotes a
			// non-ASCII path unless told not to, and the filter missed it.
			ticket := filepath.Join(r.wt, ".jaira", "tickets", "x-größe.md")
			r.write(t, ticket, "in review\n")
			r.git(t, r.wt, "add", ".jaira/tickets")
			r.git(t, r.wt, "commit", "-q", "-m", "ticket")
			r.write(t, ticket, "moved\n")
			r.write(t, filepath.Join(r.wt, ".jaira", "tickets", "x.md"), "moved\n")
			if out := r.closeTab(t, "p1", "--dry-run", "--parent", parent, "ab12cd"); !strings.Contains(out, "would close AB12CD tp1") {
				t.Fatalf("--dry-run: %q", out)
			}
			if r.closedTabs() != "" {
				t.Fatalf("--dry-run closed %q", r.closedTabs())
			}
			out := r.closeTab(t, "p1", "--parent", parent, "AB12CD")
			if !strings.Contains(out, "closing AB12CD tp1") {
				t.Fatalf("output %q", out)
			}
			// The close runs detached, after the script has returned.
			for i := 0; i < 50 && r.closedTabs() == ""; i++ {
				time.Sleep(100 * time.Millisecond)
			}
			if got := r.closedTabs(); got != "tp1" {
				t.Errorf("closed %q, want tp1", got)
			}
		})
	}
}

// A worker that still has something to do, or whose report would go with the
// tab, stays open and says why.
func TestCloseTabKeepsAWorkerThatIsNotFinished(t *testing.T) {
	t.Parallel()
	for name, c := range map[string]struct {
		setup  func(t *testing.T, r *reapRig)
		pane   string
		parent string // "-": none; "gone": one that does not listen
		why    string
	}{
		"open question": {func(t *testing.T, r *reapRig) {
			r.write(t, r.shown, `{"status":"human","question":"which one?"}`)
		}, "p1", "", "the ticket has an open question"},
		"agent's lane": {func(t *testing.T, r *reapRig) {
			r.write(t, r.shown, `{"status":"in-progress","question":""}`)
		}, "p1", "", "the ticket is in in-progress, not with a person"},
		"no parent": {func(t *testing.T, r *reapRig) {
			r.write(t, r.shown, `{"status":"signoff","question":""}`)
		}, "p1", "-", "no parent: the report is in this tab"},
		"parent gone": {func(t *testing.T, r *reapRig) {
			r.write(t, r.shown, `{"status":"signoff","question":""}`)
		}, "p1", "gone", "does not listen"},
		"another ticket's tab": {func(t *testing.T, r *reapRig) {
			r.write(t, r.shown, `{"status":"signoff","question":""}`)
			r.write(t, filepath.Join(r.labels, "p9"), "ZZ99ZZ/worker\n")
		}, "p9", "", "spawn.sh did not open it"},
		// A 'jaira json get' that fails - an older jaira without --default -
		// is a reason to stay, not an exit code the --help never promised.
		"jaira json get fails": {func(t *testing.T, r *reapRig) {
			r.write(t, r.shown, `{"status":"signoff","question":""}`)
			r.write(t, r.shown+".json-fails", "")
		}, "p1", "", "failed (exit 7)"},
		// The classes the tab's own state can be in: herdr not there at all,
		// a pane whose tab is gone, a lane the board does not list.
		"herdr not reachable": {func(t *testing.T, r *reapRig) {
			r.write(t, r.shown, `{"status":"signoff","question":""}`)
			r.env = append(r.env, "HERDR_BIN_PATH="+filepath.Join(r.wt, "no-herdr"))
		}, "p1", "", "herdr does not know pane p1"},
		"tab already gone": {func(t *testing.T, r *reapRig) {
			r.write(t, r.shown, `{"status":"signoff","question":""}`)
		}, "p42", "", "spawn.sh did not open it"},
		"lane the board does not list": {func(t *testing.T, r *reapRig) {
			r.write(t, r.shown, `{"status":"parked","question":""}`)
		}, "p1", "", "the ticket is in parked, not with a person"},
		"board whose lanes cannot be read": {func(t *testing.T, r *reapRig) {
			r.write(t, r.shown, `{"status":"signoff","question":"","board":"gone"}`)
			if err := os.Remove(r.lanes); err != nil {
				t.Fatal(err)
			}
		}, "p1", "", "jaira lanes failed"},
		"not in Herdr": {func(t *testing.T, r *reapRig) {
			r.write(t, r.shown, `{"status":"signoff","question":""}`)
		}, "", "", "not in a Herdr tab"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			r, parent := newCloseTabRig(t)
			c.setup(t, r)
			args := []string{"--parent", parent}
			switch c.parent {
			case "-":
				args = nil
			case "gone":
				args = []string{"--parent", "uds:" + filepath.Join(r.wt, "no.sock")}
			}
			out := r.closeTab(t, c.pane, append(args, "AB12CD")...)
			if !strings.HasPrefix(out, "kept AB12CD: ") || !strings.Contains(out, c.why) {
				t.Errorf("output %q, want kept with %q", out, c.why)
			}
			time.Sleep(200 * time.Millisecond)
			if got := r.closedTabs(); got != "" {
				t.Errorf("closed %q", got)
			}
		})
	}
}

// Worker and dispatcher both close their tab with close-tab.sh, after the
// report to the parent; a rename of the script has to fail here.
func TestWorkerAndDispatcherCloseTheirTab(t *testing.T) {
	t.Parallel()
	for _, role := range []string{"jaira-worker", "jaira-dispatcher"} {
		b, err := os.ReadFile(filepath.Join("builtin", role, "SKILL.md"))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), "jaira-dispatcher/scripts/close-tab.sh") {
			t.Errorf("%s does not name jaira-dispatcher/scripts/close-tab.sh", role)
		}
	}
	if _, err := File("jaira-dispatcher", "scripts/close-tab.sh"); err != nil {
		t.Errorf("jaira-dispatcher/scripts/close-tab.sh does not ship: %v", err)
	}
}
