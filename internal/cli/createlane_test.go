package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// Creating straight into a later lane is allowed, but it has to say what it
// skipped and which of the target lane's inputs are empty because of it: a
// ticket created into in-progress has no plan, and in-progress is told to
// carry out the plan.
func TestCreateIntoALaterLaneNamesWhatItSkipped(t *testing.T) {
	dir := emptyStore(t)
	out, err := runCLI(t, dir, "create", "probe", "--lane", "in-progress",
		"--assignee", "t", "--goal", "g", "--context", "c", "--dod", "d")
	if err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	// pre-process is option-gated and this ticket ticks no option, but it is
	// the lane that writes the plan reported empty below, so it is named;
	// brainstorm writes only the goal, which a flag gave, so it is not.
	skipped := ""
	for _, line := range strings.Split(out, "\n") {
		if s, ok := strings.CutPrefix(line, "Skipped on the way to in-progress: "); ok {
			skipped = ", " + s + ","
		}
	}
	if !strings.Contains(skipped, ", todo,") || !strings.Contains(skipped, ", pre-process,") || strings.Contains(skipped, "brainstorm") {
		t.Errorf("create --lane in-progress should name todo and pre-process but not brainstorm as skipped:\n%s", out)
	}
	// Only plan: goal, context and the definition of done were given as
	// flags, and a field a flag could have filled is not the skip's doing.
	if !strings.Contains(out, "Empty, though in-progress takes it as input: plan\n") {
		t.Errorf("create --lane in-progress did not say exactly the plan is empty:\n%s", out)
	}
}

func TestCreateIntoALaterLaneReportsItInJSON(t *testing.T) {
	dir := emptyStore(t)
	out, err := runCLI(t, dir, "--json", "create", "probe", "--lane", "in-progress",
		"--assignee", "t", "--goal", "g", "--context", "c", "--dod", "d")
	if err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	var payload struct {
		Skipped []string `json:"skipped_lanes"`
		Empty   []string `json:"empty_inputs"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("decode: %v\n%s", err, out)
	}
	if len(payload.Skipped) == 0 {
		t.Errorf("skipped_lanes is empty:\n%s", out)
	}
	if len(payload.Empty) != 1 || payload.Empty[0] != "plan" {
		t.Errorf("empty_inputs = %v, want [plan]", payload.Empty)
	}
}

// Without --lane nothing changes: no skip report in the text and no new keys
// in the payload.
func TestCreateWithoutLaneReportsNoSkip(t *testing.T) {
	dir := emptyStore(t)
	out, err := runCLI(t, dir, "create", "probe",
		"--assignee", "t", "--goal", "g", "--context", "c", "--dod", "d")
	if err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	if strings.Contains(out, "Skipped") || strings.Contains(out, "takes it as input") {
		t.Errorf("an ordinary create reported a skip:\n%s", out)
	}
	out, err = runCLI(t, dir, "--json", "create", "probe2",
		"--assignee", "t", "--goal", "g", "--context", "c", "--dod", "d")
	if err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	if strings.Contains(out, "skipped_lanes") || strings.Contains(out, "empty_inputs") {
		t.Errorf("an ordinary create's payload carries skip keys:\n%s", out)
	}
}
