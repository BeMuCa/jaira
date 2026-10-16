package role

import (
	"path/filepath"
	"testing"
)

func setHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	return home
}

func installInto(t *testing.T, dir string) {
	t.Helper()
	if _, err := install(dir, false, false, false); err != nil {
		t.Fatal(err)
	}
}

func TestInstalledListsOnlyDirsHoldingRoles(t *testing.T) {
	home := setHome(t)
	root := t.TempDir()
	installInto(t, filepath.Join(root, ".claude", "skills"))
	installInto(t, filepath.Join(home, ".copilot", "skills"))

	got, err := Installed(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("got %+v, want the project claude dir and the personal copilot dir", got)
	}
	if got[0].Agent.Name != "claude" || got[0].Global || got[0].Dir != filepath.Join(root, ".claude", "skills") {
		t.Errorf("first = %+v", got[0])
	}
	if got[1].Agent.Name != "copilot" || !got[1].Global || got[1].Dir != filepath.Join(home, ".copilot", "skills") {
		t.Errorf("second = %+v", got[1])
	}
}

func TestInstalledWithoutRootSkipsTheProject(t *testing.T) {
	home := setHome(t)
	installInto(t, filepath.Join(home, ".claude", "skills"))
	got, err := Installed("")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || !got[0].Global || got[0].Agent.Name != "claude" {
		t.Errorf("got %+v, want the personal claude dir only", got)
	}
}

func TestInstalledListsADirOnce(t *testing.T) {
	home := setHome(t)
	installInto(t, filepath.Join(home, ".claude", "skills"))
	got, err := Installed(home)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Errorf("got %+v, want the directory once", got)
	}
}

func TestBehindCountsWrittenAndUpdatedOnly(t *testing.T) {
	rs := []Result{
		{Action: Written}, {Action: Updated}, {Action: Unchanged},
		{Action: Skipped}, {Action: Overwritten},
	}
	if got := Behind(rs); got != 2 {
		t.Errorf("Behind = %d, want 2", got)
	}
}
