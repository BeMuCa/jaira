package selfupdate

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestDetectHomebrewCellarElement(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "Cellar", "jaira", "1.0.0", "bin", "jaira")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	kind, instr := Detect(target)
	if kind != Managed {
		t.Fatalf("Detect(%q) kind = %q, want %q", target, kind, Managed)
	}
	if !strings.Contains(instr, "brew upgrade jaira") {
		t.Errorf("instr = %q, want it to name %q", instr, "brew upgrade jaira")
	}
}

func TestDetectHomebrewPrefixEnvVar(t *testing.T) {
	dir := t.TempDir()
	prefix := filepath.Join(dir, "opt", "homebrew")
	target := filepath.Join(prefix, "bin", "jaira")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOMEBREW_PREFIX", prefix)

	kind, instr := Detect(target)
	if kind != Managed {
		t.Fatalf("Detect(%q) kind = %q, want %q", target, kind, Managed)
	}
	if !strings.Contains(instr, "brew upgrade jaira") {
		t.Errorf("instr = %q, want it to name %q", instr, "brew upgrade jaira")
	}
}

func TestDetectHomebrewCellarEnvVar(t *testing.T) {
	dir := t.TempDir()
	prefix := filepath.Join(dir, "some-cellar")
	target := filepath.Join(prefix, "jaira", "1.0.0", "bin", "jaira")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOMEBREW_CELLAR", prefix)

	kind, _ := Detect(target)
	if kind != Managed {
		t.Fatalf("Detect(%q) kind = %q, want %q", target, kind, Managed)
	}
}

// TestDetectHomebrewPrefixDoesNotMatchLookalike asserts a directory sharing a
// string prefix with $HOMEBREW_PREFIX, but not actually nested inside it, is
// not misclassified — the filepath.Rel check, not strings.HasPrefix.
func TestDetectHomebrewPrefixDoesNotMatchLookalike(t *testing.T) {
	dir := t.TempDir()
	prefix := filepath.Join(dir, "opt", "homebrew")
	lookalike := filepath.Join(dir, "opt", "homebrew-mine")
	target := filepath.Join(lookalike, "bin", "jaira")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOMEBREW_PREFIX", prefix)

	kind, _ := Detect(target)
	if kind == Managed {
		t.Errorf("Detect(%q) kind = %q, want it not to match the lookalike prefix %q", target, kind, prefix)
	}
}

func TestDetectGoInstallGOBIN(t *testing.T) {
	dir := t.TempDir()
	gobin := filepath.Join(dir, "gobin")
	target := filepath.Join(gobin, "jaira")
	if err := os.MkdirAll(gobin, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GOBIN", gobin)

	kind, instr := Detect(target)
	if kind != GoInstall {
		t.Fatalf("Detect(%q) kind = %q, want %q", target, kind, GoInstall)
	}
	if !strings.Contains(instr, "go install github.com/BeMuCa/jaira/cmd/jaira@latest") {
		t.Errorf("instr = %q, want it to name the go install command", instr)
	}
}

func TestDetectGoInstallGOPATHBin(t *testing.T) {
	dir := t.TempDir()
	gopath := filepath.Join(dir, "gopath")
	target := filepath.Join(gopath, "bin", "jaira")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GOBIN", "")
	t.Setenv("GOPATH", gopath)

	kind, _ := Detect(target)
	if kind != GoInstall {
		t.Fatalf("Detect(%q) kind = %q, want %q", target, kind, GoInstall)
	}
}

// TestDetectUnwritableDirectoryRefuses asserts a directory the process
// cannot write to classifies as Unwritable and names the directory.
//
// Skipped on Windows (chmod cannot make a directory unreadable/unwritable
// there the way it can on unix) and under root (permission bits do not
// block root), matching TestNudgeIfStaleSkipsSilentlyWhenStateDirUnreadable
// in internal/cli/update_test.go.
func TestDetectUnwritableDirectoryRefuses(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("chmod cannot make a directory unwritable on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("running as root: permission bits do not block writes")
	}
	parent := t.TempDir()
	dir := filepath.Join(parent, "unwritable")
	if err := os.MkdirAll(dir, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(dir, 0o755) })
	target := filepath.Join(dir, "jaira")

	kind, instr := Detect(target)
	if kind != Unwritable {
		t.Fatalf("Detect(%q) kind = %q, want %q", target, kind, Unwritable)
	}
	if !strings.Contains(instr, dir) {
		t.Errorf("instr = %q, want it to name %q", instr, dir)
	}
}

func TestDetectSelfManagedForPlainWritableDir(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "jaira")

	kind, instr := Detect(target)
	if kind != SelfManaged {
		t.Fatalf("Detect(%q) kind = %q, want %q", target, kind, SelfManaged)
	}
	if instr != "" {
		t.Errorf("instr = %q, want empty for %q", instr, SelfManaged)
	}
}

// TestDetectHomebrewWinsOverGoInstall asserts a path matching both
// classifications resolves to Homebrew: the package manager's claim is
// stronger than go install's.
func TestDetectHomebrewWinsOverGoInstall(t *testing.T) {
	dir := t.TempDir()
	gobin := filepath.Join(dir, "Cellar", "jaira", "1.0.0", "bin")
	target := filepath.Join(gobin, "jaira")
	if err := os.MkdirAll(gobin, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GOBIN", gobin)

	kind, _ := Detect(target)
	if kind != Managed {
		t.Fatalf("Detect(%q) kind = %q, want %q (homebrew wins over go-install)", target, kind, Managed)
	}
}
