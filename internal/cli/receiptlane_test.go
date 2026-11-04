package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// sendLane is a lane that runs a skill and a script, and must leave with a
// receipt. Its files go on the board, so the test runs the real loader on them.
const sendLane = `---
id: send
name: Send
after: in-progress
agentic: true
model-tier: cheap
skill: ponytail
script: scripts/send.sh
receipt: required
---
Send the summary once.
`

// receiptBoard is an empty board with the built-in lanes set up, the lane files
// given, the skill and script they name on disk, and one ticket in each lane.
// It returns the board directory and the handle of each ticket by lane.
func receiptBoard(t *testing.T, lanes map[string]string, lanesFor ...string) (string, map[string]string) {
	t.Helper()
	dir := emptyStore(t)
	if _, err := lane.Load(dir); err != nil {
		t.Fatal(err)
	}
	for name, body := range lanes {
		if err := os.WriteFile(filepath.Join(lane.ProjectLanesDir(dir), name+".md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, p := range []string{
		filepath.Join(dir, ".claude", "skills", "ponytail", "SKILL.md"),
		filepath.Join(dir, "scripts", "send.sh"),
	} {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("skill or script\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	handles := map[string]string{}
	for i, laneID := range lanesFor {
		tk, err := s.Create(map[string]string{
			ticket.FieldID:    ticket.NewID(time.Date(2026, 9, 16, 8+i, 0, 0, 0, time.UTC)),
			ticket.FieldTitle: "t", ticket.FieldStatus: laneID,
		}, nil, "")
		if err != nil {
			t.Fatal(err)
		}
		handles[laneID] = ticket.Handle(tk.ID)
	}
	return dir, handles
}

// A lane that names a skill and a script hands the agent both, with the key a
// provider can use to refuse a second send, and the run line that carries it.
func TestForLaneNamesTheSkillScriptAndKey(t *testing.T) {
	dir, h := receiptBoard(t, map[string]string{"send": sendLane}, "send")
	got := jsonCLI(t, dir, "show", h["send"], "--for-lane", "send")

	if got["skill"] != "ponytail" {
		t.Errorf("skill = %v, want ponytail", got["skill"])
	}
	if p, _ := got["skill_path"].(string); !strings.HasSuffix(p, filepath.Join("ponytail", "SKILL.md")) {
		t.Errorf("skill_path = %v, want the ponytail SKILL.md", got["skill_path"])
	}
	key, _ := got["idempotency_key"].(string)
	if want := got["ticket_id"].(string) + "-send-r1"; key != want {
		t.Errorf("idempotency_key = %q, want %q", key, want)
	}
	if got["script"] != "scripts/send.sh" || got["script_done"] != false {
		t.Errorf("script = %v script_done = %v, want scripts/send.sh and false", got["script"], got["script_done"])
	}
	if run, _ := got["run"].(string); !strings.Contains(run, "JAIRA_IDEMPOTENCY_KEY="+key) || !strings.Contains(run, "./scripts/send.sh") {
		t.Errorf("run = %q, want the key and ./scripts/send.sh", run)
	}
	if got["receipt_required"] != true || got["receipt"] != "" {
		t.Errorf("receipt_required = %v receipt = %v, want true and empty", got["receipt_required"], got["receipt"])
	}
	if got["complete"] != true {
		t.Errorf("complete = %v, want true when the skill and script are on disk: %v", got["complete"], got["missing"])
	}
}

// A lane that declares none of them adds none of the keys, so a reader never
// takes an absent fact for an empty one.
func TestForLaneOfAPlainLaneHasNoSkillKeys(t *testing.T) {
	dir, h := receiptBoard(t, nil, "in-progress")
	got := jsonCLI(t, dir, "show", h["in-progress"], "--for-lane", "in-progress")
	for _, k := range []string{"skill", "skill_path", "script", "script_done", "run", "idempotency_key", "receipt_required", "receipt"} {
		if _, ok := got[k]; ok {
			t.Errorf("a lane with no skill, script or receipt carries %q", k)
		}
	}
}

// The receipt is the record that the side effect happened. Once it is written
// the run line is empty, a fresh agent is told not to run the script, and a second
// receipt for the same round is refused.
func TestReceiptEndsTheRunForThisRound(t *testing.T) {
	dir, h := receiptBoard(t, map[string]string{"send": sendLane}, "send")
	out, err := runCLI(t, dir, "receipt", h["send"], "message", "42")
	if err != nil {
		t.Fatalf("receipt: %v\n%s", err, out)
	}

	got := jsonCLI(t, dir, "show", h["send"], "--for-lane", "send")
	if got["script_done"] != true || got["run"] != "" {
		t.Errorf("after the receipt script_done = %v run = %q, want true and empty", got["script_done"], got["run"])
	}
	if r, _ := got["receipt"].(string); !strings.HasSuffix(r, " message 42") {
		t.Errorf("receipt = %q, want the key then the result", r)
	}
	text, err := runCLI(t, dir, "show", h["send"], "--for-lane", "send")
	if err != nil {
		t.Fatalf("show --for-lane: %v\n%s", err, text)
	}
	if !strings.Contains(text, "already ran for key") {
		t.Errorf("the text payload does not say the script already ran:\n%s", text)
	}

	again, err := runCLI(t, dir, "receipt", h["send"], "second send")
	if exitCode(err) != ExitValidation {
		t.Errorf("a second receipt for the same round: exit %d, want %d:\n%s", exitCode(err), ExitValidation, again)
	}
}

// A receipt for a lane that does not require one would record a side effect the
// gate never asks about, so the command refuses it.
func TestReceiptRefusedWhereTheLaneDoesNotRequireOne(t *testing.T) {
	dir, h := receiptBoard(t, map[string]string{"send": sendLane}, "in-progress")
	out, err := runCLI(t, dir, "receipt", h["in-progress"], "message 42")
	if exitCode(err) != ExitValidation {
		t.Errorf("a receipt on in-progress: exit %d, want %d:\n%s", exitCode(err), ExitValidation, out)
	}
}

// A named skill that is not installed, or a script that is not in the repository,
// makes the payload incomplete and says what is missing: the lane is not skipped.
func TestForLaneReportsAMissingSkillAndScript(t *testing.T) {
	ghost := "---\nid: ghost\nname: Ghost\nafter: in-progress\nagentic: true\nmodel-tier: cheap\nskill: nobody-installed-this\nscript: scripts/absent.sh\n---\n"
	dir, h := receiptBoard(t, map[string]string{"ghost": ghost}, "ghost")
	got := jsonCLI(t, dir, "show", h["ghost"], "--for-lane", "ghost")
	if got["complete"] != false {
		t.Fatalf("complete = %v, want false", got["complete"])
	}
	missing, _ := got["missing"].([]any)
	joined := ""
	for _, m := range missing {
		joined += m.(string) + "\n"
	}
	if !strings.Contains(joined, "skill nobody-installed-this") || !strings.Contains(joined, "script scripts/absent.sh") {
		t.Errorf("missing = %v, want the skill and the script named", missing)
	}
}

// lanes show names a lane's skill, script and receipt, so a person or an agent
// reading the lane's contract sees what it does without opening the file.
func TestLanesShowNamesTheSkillScriptAndReceipt(t *testing.T) {
	dir, _ := receiptBoard(t, map[string]string{"send": sendLane})
	out, err := runCLI(t, dir, "lanes", "show", "send")
	if err != nil {
		t.Fatalf("lanes show: %v\n%s", err, out)
	}
	for _, want := range []string{"Skill:       ponytail", "Script:      scripts/send.sh", "Receipt:     required"} {
		if !strings.Contains(out, want) {
			t.Errorf("lanes show is missing %q:\n%s", want, out)
		}
	}
	got := jsonCLI(t, dir, "lanes", "show", "send")
	if got["skill"] != "ponytail" || got["script"] != "scripts/send.sh" || got["receipt_required"] != true {
		t.Errorf("lanes show --json = skill %v script %v receipt_required %v", got["skill"], got["script"], got["receipt_required"])
	}
}
