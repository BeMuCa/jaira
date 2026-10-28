//go:build !windows

package session

import (
	"errors"
	"syscall"
)

// processRunning asks the kernel whether pid exists; EPERM means it does but
// belongs to someone else.
func processRunning(pid int) bool {
	err := syscall.Kill(pid, 0)
	return err == nil || errors.Is(err, syscall.EPERM)
}
