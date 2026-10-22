//go:build windows

package session

// processRunning cannot be answered cheaply here, so a claim on Windows lives
// for its lease and no longer.
func processRunning(int) bool { return false }
