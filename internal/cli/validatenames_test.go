package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestValidateNamesFilesWithInvalidUTF8 plants the name an older jaira wrote —
// a slug cut through "ü", leaving the lone byte \xc3 — and checks that
// validate fails on it, offers a rename bash can run as printed, leaves the
// file where it is, and passes once the rename ran.
func TestValidateNamesFilesWithInvalidUTF8(t *testing.T) {
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		t.Skip("this file system cannot hold a name that is not valid UTF-8")
	}
	dir, _ := updateTestStore(t)
	logbook := filepath.Join(dir, ".jaira", "logbook", "someone-2026-11-22")
	if err := os.MkdirAll(logbook, 0o755); err != nil {
		t.Fatal(err)
	}
	bad := filepath.Join(logbook, "01ABC-ampel-gr\xc3.md")
	if err := os.WriteFile(bad, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	out, err := runCLI(t, dir, "validate")
	if exitCode(err) != ExitValidation {
		t.Fatalf("validate exit = %d, want %d\n%s", exitCode(err), ExitValidation, out)
	}
	const marker = "macOS cannot check it out — rename it: git mv "
	i := strings.Index(out, marker)
	if i < 0 {
		t.Fatalf("validate did not name the file:\n%s", out)
	}
	if _, err := os.Stat(bad); err != nil {
		t.Fatalf("validate touched the file: %v", err)
	}

	// The printed rename, run by bash as printed, must hit the real file.
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("no bash to run the printed rename")
	}
	args, _, _ := strings.Cut(out[i+len(marker):], "\n")
	sh := exec.Command("bash", "-c", "mv "+args)
	sh.Dir = dir
	if b, err := sh.CombinedOutput(); err != nil {
		t.Fatalf("the printed rename does not run: %v\n%s\n%s", err, args, b)
	}
	if _, err := os.Stat(filepath.Join(logbook, "01ABC-ampel-gr.md")); err != nil {
		t.Fatalf("rename did not produce the valid name: %v", err)
	}
	if out, err := runCLI(t, dir, "validate"); err != nil || strings.Contains(out, marker) {
		t.Fatalf("validate after the rename: err=%v\n%s", err, out)
	}
}
