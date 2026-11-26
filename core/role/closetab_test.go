package role

import (
	"fmt"
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

// A worker whose ticket is done closes its own tab once its report went to
// the parent. Board files under .jaira/ are no reason to stay.
func TestCloseTabClosesAWorkerWhoseTicketIsDone(t *testing.T) {
	t.Parallel()
	r, parent := newCloseTabRig(t)
	r.write(t, r.shown, `{"status":"done","question":""}`)
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
	r.waitClosed(t, "tp1")
}

// waitClosed waits for the detached close, which runs after the script has
// returned.
func (r *reapRig) waitClosed(t *testing.T, want string) {
	t.Helper()
	for i := 0; i < 50 && r.closedTabs() != want; i++ {
		time.Sleep(100 * time.Millisecond)
	}
	if got := r.closedTabs(); got != want {
		t.Errorf("closed %q, want %q", got, want)
	}
}

// The worker's own record, as spawn.sh wrote it for pane p1.
func (r *reapRig) workerRecord(t *testing.T) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(r.repo, ".git", "jaira", "panes", "p1"))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// A worker whose ticket waits for a person in a lane it can come back from —
// signoff, or human with the question gone — keeps its tab while its prompt
// cache is warm, so a give-back reaches it: it parks the tab in its record
// and closes nothing. --dry-run parks nothing.
func TestCloseTabParksAWorkerWhoseTicketIsWithAPerson(t *testing.T) {
	t.Parallel()
	for _, lane := range []string{"signoff", "human"} {
		t.Run(lane, func(t *testing.T) {
			t.Parallel()
			r, parent := newCloseTabRig(t)
			r.write(t, r.shown, `{"status":"`+lane+`","question":""}`)
			if out := r.closeTab(t, "p1", "--dry-run", "--parent", parent, "AB12CD"); !strings.Contains(out, "kept AB12CD: cache warm until ") {
				t.Fatalf("--dry-run: %q", out)
			}
			if rec := r.workerRecord(t); strings.Contains(rec, "parked=") {
				t.Fatalf("--dry-run parked: %q", rec)
			}
			before := time.Now().Unix()
			out := r.closeTab(t, "p1", "--parent", parent, "AB12CD")
			if !strings.HasPrefix(out, "kept AB12CD: cache warm until ") || !strings.Contains(out, "reap.sh closes it after") {
				t.Fatalf("output %q", out)
			}
			rec := r.workerRecord(t)
			if !strings.HasPrefix(rec, "pane=p1\nticket=AB12CD\nlane=worker\nworktree="+r.wt+"\nparked=") {
				t.Fatalf("record %q", rec)
			}
			// Parked twice - given back, carried again, finished again - is one
			// line with the newer time, never two.
			r.closeTab(t, "p1", "--parent", parent, "AB12CD")
			if n := strings.Count(r.workerRecord(t), "parked="); n != 1 {
				t.Errorf("%d parked= lines: %q", n, r.workerRecord(t))
			}
			var at int64
			if _, err := fmt.Sscanf(r.workerRecord(t)[strings.Index(r.workerRecord(t), "parked=")+7:], "%d", &at); err != nil || at < before {
				t.Errorf("parked=%d, want a time from %d on (%v)", at, before, err)
			}
			time.Sleep(200 * time.Millisecond)
			if got := r.closedTabs(); got != "" {
				t.Errorf("closed %q", got)
			}
		})
	}
}

// reap.sh hands every parked tab back to close-tab.sh --reap, which closes it
// once the cache window is over, the ticket is done or filed away - and only
// then, whatever the forge says. A closed or vanished tab is unparked, so the
// next round leaves it alone.
func TestReapClosesAParkedTabOnceItsCacheIsCold(t *testing.T) {
	t.Parallel()
	for name, c := range map[string]struct {
		setup  func(t *testing.T, r *reapRig)
		closes bool
		says   string
	}{
		"cache warm": {func(t *testing.T, r *reapRig) {}, false, "kept AB12CD: cache warm for 60 more min"},
		"window over": {func(t *testing.T, r *reapRig) {
			r.env = append(r.env, "JAIRA_CACHE_WINDOW=0")
		}, true, "closing AB12CD tp1"},
		"ticket done": {func(t *testing.T, r *reapRig) {
			r.write(t, r.shown, `{"status":"done","question":""}`)
		}, true, "closing AB12CD tp1"},
		"ticket filed away": {func(t *testing.T, r *reapRig) {
			if err := os.Remove(r.shown); err != nil {
				t.Fatal(err)
			}
		}, true, "closing AB12CD tp1"},
		"given back to the same worker": {func(t *testing.T, r *reapRig) {
			r.write(t, r.shown, `{"status":"in-progress","question":""}`)
			r.env = append(r.env, "JAIRA_CACHE_WINDOW=0")
		}, false, "kept AB12CD: the ticket is in in-progress, not with a person"},
		"agent working": {func(t *testing.T, r *reapRig) {
			r.write(t, filepath.Join(r.labels, "p1.status"), "working\n")
			r.env = append(r.env, "JAIRA_CACHE_WINDOW=0")
		}, false, "kept AB12CD: its agent is working"},
		"open question": {func(t *testing.T, r *reapRig) {
			r.write(t, r.shown, `{"status":"human","question":"which one?"}`)
			r.env = append(r.env, "JAIRA_CACHE_WINDOW=0")
		}, false, "kept AB12CD: the ticket has an open question"},
		// A herdr that gave no answer keeps the park: unparked, no later
		// round would close the tab.
		"herdr gives no answer": {func(t *testing.T, r *reapRig) {
			r.write(t, filepath.Join(r.labels, "p1.fails"), "echo 'request timed out' >&2; exit 1\n")
			r.env = append(r.env, "JAIRA_CACHE_WINDOW=0")
		}, false, "kept AB12CD: herdr gave no answer for pane p1: request timed out"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			r, parent := newCloseTabRig(t)
			r.write(t, r.shown, `{"status":"signoff","question":""}`)
			r.closeTab(t, "p1", "--parent", parent, "AB12CD")
			c.setup(t, r)
			if out := r.run(t, r.repo, "reap.sh", "--dry-run"); c.closes && !strings.Contains(out, "would close AB12CD tp1") {
				t.Errorf("--dry-run: %q", out)
			}
			if !strings.Contains(r.workerRecord(t), "parked=") {
				t.Fatalf("--dry-run unparked: %q", r.workerRecord(t))
			}
			out := r.run(t, r.repo, "reap.sh")
			if !strings.Contains(out, c.says) {
				t.Fatalf("output %q, want %q", out, c.says)
			}
			if !c.closes {
				time.Sleep(200 * time.Millisecond)
				if got := r.closedTabs(); got != "" {
					t.Errorf("closed %q", got)
				}
				if !strings.Contains(r.workerRecord(t), "parked=") {
					t.Errorf("unparked a tab it kept: %q", r.workerRecord(t))
				}
				return
			}
			r.waitClosed(t, "tp1")
			// The record stays for the merged-ticket reap, which removes the
			// worktree; only its parked= line goes.
			if rec := r.workerRecord(t); strings.Contains(rec, "parked=") {
				t.Errorf("record still parked after the close: %q", rec)
			}
			if out := r.run(t, r.repo, "reap.sh"); strings.Contains(out, "AB12CD tp1") {
				t.Errorf("second round asked about the closed tab again: %q", out)
			}
		})
	}
}

// A parked tab the person closed by hand is unparked on the next round.
func TestReapUnparksATabThatIsGone(t *testing.T) {
	t.Parallel()
	r, parent := newCloseTabRig(t)
	r.write(t, r.shown, `{"status":"signoff","question":""}`)
	r.closeTab(t, "p1", "--parent", parent, "AB12CD")
	r.write(t, filepath.Join(r.labels, "p1"), "\n")
	if out := r.run(t, r.repo, "reap.sh"); !strings.Contains(out, "spawn.sh did not open it") {
		t.Fatalf("output %q", out)
	}
	if rec := r.workerRecord(t); strings.Contains(rec, "parked=") {
		t.Errorf("still parked: %q", rec)
	}
	if got := r.closedTabs(); got != "" {
		t.Errorf("closed %q", got)
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
		}, "p1", "", "herdr gave no answer for pane p1"},
		// Only herdr's own not_found says the tab is gone; a timeout says
		// nothing about it.
		"pane get times out": {func(t *testing.T, r *reapRig) {
			r.write(t, r.shown, `{"status":"signoff","question":""}`)
			r.write(t, filepath.Join(r.labels, "p1.fails"), "echo 'request timed out' >&2; exit 1\n")
		}, "p1", "", "herdr gave no answer for pane p1: request timed out"},
		"tab get times out": {func(t *testing.T, r *reapRig) {
			r.write(t, r.shown, `{"status":"signoff","question":""}`)
			r.write(t, filepath.Join(r.labels, "tp1.fails"), "exit 1\n")
		}, "p1", "", "herdr gave no answer for tab tp1: exit 1, no message"},
		// A step of its own that fails keeps the tab with a kept line too,
		// not an exit 1 the ERR trap never sees inside a function.
		"no temp directory": {func(t *testing.T, r *reapRig) {
			r.write(t, r.shown, `{"status":"signoff","question":""}`)
			r.env = append(r.env, "TMPDIR="+filepath.Join(r.wt, "no-tmp"))
		}, "p1", "", "mktemp failed: cannot ask herdr about pane p1"},
		"pane unknown to hc": {func(t *testing.T, r *reapRig) {
			r.write(t, r.shown, `{"status":"signoff","question":""}`)
			r.write(t, filepath.Join(r.labels, "p1.fails"), `echo '{"error":{"kind":"not_found"}}' >&2; exit 3`+"\n")
		}, "p1", "", "herdr does not know pane p1"},
		"pane unknown to the old herdr": {func(t *testing.T, r *reapRig) {
			r.write(t, r.shown, `{"status":"signoff","question":""}`)
			r.write(t, filepath.Join(r.labels, "p1.fails"), `echo '{"error":{"code":"pane_not_found"}}'; exit 1`+"\n")
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
