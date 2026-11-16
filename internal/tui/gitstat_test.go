package tui

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestGitStatNeverHandsGitAnOption pins that the ticket detail and sign-off
// panes pass git only commit ids from the hand-editable commits: field:
// "--output=<path>" used to make 'git show' overwrite that file.
func TestGitStatNeverHandsGitAnOption(t *testing.T) {
	root := t.TempDir()
	for _, args := range [][]string{
		{"init", "-q"},
		{"-c", "user.name=t", "-c", "user.email=t@t", "commit", "-q", "--allow-empty", "-m", "base"},
	} {
		if out, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	victim := filepath.Join(t.TempDir(), "victim")
	if err := os.WriteFile(victim, []byte("keep me\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	stat, err := (&gitStat{root: root}).of([]string{"--output=" + victim})
	if err != nil || stat != "" {
		t.Errorf("of = %q, %v; want nothing for a value that is no commit id", stat, err)
	}
	if b, err := os.ReadFile(victim); err != nil || string(b) != "keep me\n" {
		t.Errorf("git wrote over the file: %q, %v", b, err)
	}
}
