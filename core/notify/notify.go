// Package notify raises a desktop notification through whatever the machine
// already has.
//
// It shells out — notify-send on Linux and the BSDs, osascript on macOS,
// PowerShell's toast on Windows — and brings no dependency of its own. A
// notification library would mean cgo or a Windows API binding, and this tool
// ships as one static binary with no runtime dependency; a feature that
// convenient is not worth that.
//
// Every failure is silent by design. There is no notifier in a container, over
// ssh, or on a headless CI machine, and a board that printed a warning every
// time it could not pop up a bubble would be worse than one that quietly does
// not.
package notify

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// timeout caps how long a notifier may take. osascript in particular can sit
// waiting on a hung WindowServer, and a ticket write must not be held up by the
// desktop.
const timeout = 3 * time.Second

// Available reports whether a notifier could be found. Callers use it to skip
// the work of assembling a message that nothing can show.
func Available() bool { return command("t", "b") != nil }

// Send shows a notification. It returns whether anything was shown, for
// callers that want to fall back to a printed line.
func Send(title, body string) bool {
	cmd := command(title, body)
	if cmd == nil {
		return false
	}
	done := make(chan error, 1)
	if err := cmd.Start(); err != nil {
		return false
	}
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		return err == nil
	case <-time.After(timeout):
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		return false
	}
}

func command(title, body string) *exec.Cmd {
	switch runtime.GOOS {
	case "darwin":
		if path, err := exec.LookPath("osascript"); err == nil {
			script := "display notification " + quote(body) + " with title " + quote(title)
			return exec.Command(path, "-e", script)
		}
	case "windows":
		if path, err := exec.LookPath("powershell.exe"); err == nil {
			// BurntToast is not installed anywhere by default, so this uses
			// the balloon tip that plain .NET provides.
			script := "[reflection.assembly]::LoadWithPartialName('System.Windows.Forms')>$null;" +
				"$n=New-Object System.Windows.Forms.NotifyIcon;" +
				"$n.Icon=[System.Drawing.SystemIcons]::Information;$n.Visible=$true;" +
				"$n.ShowBalloonTip(5000," + psQuote(title) + "," + psQuote(body) + ",'Info')"
			return exec.Command(path, "-NoProfile", "-Command", script)
		}
	default:
		// WSL counts as linux and has no notifier of its own, but it can reach
		// the Windows one, so that is tried before giving up.
		if path, err := exec.LookPath("notify-send"); err == nil {
			return exec.Command(path, "-a", "jaira", title, body)
		}
		if isWSL() {
			if path, err := exec.LookPath("powershell.exe"); err == nil {
				return exec.Command(path, "-NoProfile", "-Command",
					"[reflection.assembly]::LoadWithPartialName('System.Windows.Forms')>$null;"+
						"$n=New-Object System.Windows.Forms.NotifyIcon;"+
						"$n.Icon=[System.Drawing.SystemIcons]::Information;$n.Visible=$true;"+
						"$n.ShowBalloonTip(5000,"+psQuote(title)+","+psQuote(body)+",'Info')")
			}
		}
	}
	return nil
}

func isWSL() bool {
	if os.Getenv("WSL_DISTRO_NAME") != "" || os.Getenv("WSL_INTEROP") != "" {
		return true
	}
	b, err := os.ReadFile("/proc/version")
	return err == nil && strings.Contains(strings.ToLower(string(b)), "microsoft")
}

// quote renders a string as an AppleScript literal.
func quote(s string) string {
	return `"` + strings.ReplaceAll(strings.ReplaceAll(s, `\`, `\\`), `"`, `\"`) + `"`
}

// psQuote renders a string as a PowerShell single-quoted literal, where the
// only escape is a doubled quote.
func psQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}
