package session

import (
	"os"
	"os/exec"
	"strconv"
	"testing"
)

// The agent's process is recorded with the claim, and the claim is alive here
// exactly as long as that process runs and the claim is still that session's.
func TestAClaimLivesAsLongAsItsAgentProcess(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())
	t.Setenv("CLAUDE_PID", "")
	t.Setenv("JAIRA_AGENT_PID", strconv.Itoa(os.Getpid()))

	if err := RecordClaimProcess("01TICKET", "sess-1"); err != nil {
		t.Fatal(err)
	}
	if !ClaimProcessAlive("01TICKET", "sess-1") {
		t.Error("the claim of a running process reads as dead")
	}
	if ClaimProcessAlive("01TICKET", "sess-2") {
		t.Error("another session's claim borrowed this one's process")
	}

	if err := ForgetClaimProcess("01TICKET"); err != nil {
		t.Fatal(err)
	}
	if ClaimProcessAlive("01TICKET", "sess-1") {
		t.Error("a released claim still reads as alive")
	}
}

func TestAClaimOfAFinishedProcessIsDead(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())
	t.Setenv("CLAUDE_PID", "")
	cmd := exec.Command("true")
	if err := cmd.Run(); err != nil {
		t.Skip("no 'true' to run:", err)
	}
	t.Setenv("JAIRA_AGENT_PID", strconv.Itoa(cmd.Process.Pid))
	if err := RecordClaimProcess("01TICKET", "sess-1"); err != nil {
		t.Fatal(err)
	}
	if ClaimProcessAlive("01TICKET", "sess-1") {
		t.Error("the claim of a process that has exited reads as alive")
	}
}

// A claim taken where no agent process is named leaves nothing behind — and
// removes what an earlier claim of the same ticket left.
func TestAClaimWithoutAnAgentProcessRemovesTheOldRecord(t *testing.T) {
	t.Setenv("JAIRA_HOME", t.TempDir())
	t.Setenv("CLAUDE_PID", "")
	t.Setenv("JAIRA_AGENT_PID", strconv.Itoa(os.Getpid()))
	if err := RecordClaimProcess("01TICKET", "sess-1"); err != nil {
		t.Fatal(err)
	}
	t.Setenv("JAIRA_AGENT_PID", "")
	if err := RecordClaimProcess("01TICKET", "sess-1"); err != nil {
		t.Fatal(err)
	}
	if ClaimProcessAlive("01TICKET", "sess-1") {
		t.Error("a claim with no agent process inherited the earlier one's")
	}
}
