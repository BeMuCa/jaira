//go:build windows

package selfupdate

import (
	"os"
	"path/filepath"
	"testing"
)

// Renaming the running image aside as .old-<pid> is how every self-updating
// tool does this on Windows, so a shared bin directory can hold another one's
// backup. The pattern is anchored to jaira's own name for exactly that reason.
func TestSweepLeavesOtherProgramsLeftoversAlone(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "jaira.exe")
	mine := target + ".old-1234"
	theirs := filepath.Join(dir, "other.exe.old-5678")
	for _, p := range []string{mine, theirs} {
		if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	if n := Sweep(target); n != 1 {
		t.Errorf("Sweep removed %d file(s), want only jaira's own", n)
	}
	if _, err := os.Stat(mine); err == nil {
		t.Error("jaira's own leftover survived")
	}
	if _, err := os.Stat(theirs); err != nil {
		t.Errorf("another program's leftover was deleted: %v", err)
	}
}
