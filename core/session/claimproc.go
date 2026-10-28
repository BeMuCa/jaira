package session

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// A claim lives on the ticket's ref and says who took it and when; after
// ClaimTTL nobody can tell whether that agent is still at work. On the machine
// the agent runs on, its process can: the claim file records which process
// holds a ticket, so the board there can keep showing the ticket as worked for
// as long as that process lives. It stays on this machine, because a process id
// means nothing anywhere else, and writing it onto the ticket would change the
// ref on every claim.

// agentPIDEnv are the variables naming the agent's own long-lived process,
// first match wins. Claude Code sets CLAUDE_PID in every shell it opens; any
// other agent can set JAIRA_AGENT_PID.
var agentPIDEnv = []string{"JAIRA_AGENT_PID", "CLAUDE_PID"}

func claimsDir() string {
	if v := os.Getenv("JAIRA_HOME"); v != "" {
		return filepath.Join(v, "claims")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".jaira", "claims")
}

func claimPath(ticketID string) string {
	dir := claimsDir()
	if dir == "" {
		return ""
	}
	return filepath.Join(dir, sanitize(ticketID))
}

// RecordClaimProcess notes that sessionID's agent process holds ticketID. With
// no agent process to name it removes any older note instead, so a claim taken
// from a plain shell never inherits the liveness of an earlier one.
func RecordClaimProcess(ticketID, sessionID string) error {
	pid := 0
	for _, k := range agentPIDEnv {
		if n, err := strconv.Atoi(os.Getenv(k)); err == nil && n > 0 {
			pid = n
			break
		}
	}
	if pid == 0 {
		return ForgetClaimProcess(ticketID)
	}
	path := claimPath(ticketID)
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	line := fmt.Sprintf("%s %d %s\n", sessionID, pid, processName(pid))
	return os.WriteFile(path, []byte(line), 0o644)
}

// ForgetClaimProcess removes the note for ticketID, if there is one.
func ForgetClaimProcess(ticketID string) error {
	path := claimPath(ticketID)
	if path == "" {
		return nil
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// ClaimProcessAlive reports whether the process recorded for ticketID belongs
// to sessionID and is still running on this machine. A recorded process name
// that no longer matches means the id was reused by something else.
func ClaimProcessAlive(ticketID, sessionID string) bool {
	path := claimPath(ticketID)
	if path == "" || sessionID == "" {
		return false
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	f := strings.Fields(string(b))
	if len(f) < 2 || f[0] != sessionID {
		return false
	}
	pid, err := strconv.Atoi(f[1])
	if err != nil || pid <= 0 || !processRunning(pid) {
		return false
	}
	if len(f) >= 3 {
		if now := processName(pid); now != "" && now != f[2] {
			return false
		}
	}
	return true
}

// processName is the running program's name where the system says it cheaply
// (/proc on Linux), "" elsewhere — one word, so it fits the note's line.
func processName(pid int) string {
	b, err := os.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "comm"))
	if err != nil {
		return ""
	}
	return strings.ReplaceAll(strings.TrimSpace(string(b)), " ", "_")
}
