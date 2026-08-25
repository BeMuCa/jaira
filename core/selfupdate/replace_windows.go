//go:build windows

package selfupdate

import (
	"os"
	"path/filepath"
	"strconv"
)

// Replace installs bin at target.
//
// Windows refuses to overwrite the image of a running executable, but it
// will let you rename it out of the way — which is how every self-updating
// Windows tool does this. Sweep runs first so a previous run's leftover
// .old-<pid> file does not sit there forever; if the second rename below
// fails, the old file is renamed back so a failed upgrade leaves the install
// exactly as it was rather than half-swapped.
func Replace(target string, bin []byte) error {
	Sweep(target)

	tmp, err := stage(target, bin)
	if err != nil {
		return err
	}
	old := target + ".old-" + strconv.Itoa(os.Getpid())
	if err := os.Rename(target, old); err != nil {
		os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, target); err != nil {
		os.Rename(old, target)
		os.Remove(tmp)
		return err
	}
	return nil
}

// Sweep removes the .old-<pid> files a previous Replace of target left behind.
// The previous run's process is no longer running by the time a later run
// starts, so this is where those leftovers actually get to die. Removal
// errors are ignored — a still-locked file simply gets swept next time.
//
// It takes the target rather than its directory so the pattern can be anchored
// to jaira's own name. Renaming the running image aside as .old-<pid> is the
// standard way this is done on Windows, so a shared bin directory can hold
// another self-updating tool's backup — and an unanchored "*.old-*" would
// delete it.
func Sweep(target string) int {
	matches, err := filepath.Glob(target + ".old-*")
	if err != nil {
		return 0
	}
	n := 0
	for _, m := range matches {
		if os.Remove(m) == nil {
			n++
		}
	}
	return n
}
