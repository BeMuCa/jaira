//go:build !windows

package selfupdate

import "os"

// Replace atomically installs bin at target.
//
// A running process on unix keeps its inode open even after the file at its
// path is replaced, so renaming a new binary over a running one is safe: the
// process currently executing target keeps running from the old inode until
// it exits, and the next invocation of the path picks up the new file. There
// is nothing else to do here.
func Replace(target string, bin []byte) error {
	tmp, err := stage(target, bin)
	if err != nil {
		return err
	}
	if err := os.Rename(tmp, target); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}

// Sweep is a no-op on unix. There is nothing to sweep because Replace never
// needs to rename anything out of the way before it can write — unix simply
// does not need this. It takes the target, not a directory, so the Windows
// build can anchor its pattern to jaira's own name.
func Sweep(target string) int { return 0 }
