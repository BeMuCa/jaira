package web

import (
	"os/exec"
	"runtime"

	"github.com/BeMuCa/jaira/core/notify"
)

// OpenBrowser opens url in the default browser and does not wait for it.
//
// It is best effort: a machine with no browser (a container, an ssh session)
// still gets the address printed by the caller, which is the fallback a person
// can always act on.
func OpenBrowser(url string) error {
	cmd := browserCommand(url)
	if cmd == nil {
		return exec.ErrNotFound
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	go cmd.Wait()
	return nil
}

func browserCommand(url string) *exec.Cmd {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url)
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	}
	// WSL counts as linux, and usually has no Linux browser: the Windows one
	// is reached through wslview when wslu is installed, else through cmd.exe.
	if notify.IsWSL() {
		if p, err := exec.LookPath("wslview"); err == nil {
			return exec.Command(p, url)
		}
		if p, err := exec.LookPath("cmd.exe"); err == nil {
			// start treats its first quoted argument as a window title.
			return exec.Command(p, "/c", "start", `""`, url)
		}
	}
	if p, err := exec.LookPath("xdg-open"); err == nil {
		return exec.Command(p, url)
	}
	return nil
}
