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
// a slug cut through "ü", leaving the lone byte \xc3 — beside the name the
// fixed slug writes, and checks that validate fails on the broken one only,
// quotes it so bash finds the very file, and touches nothing.
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
	keep := filepath.Join(logbook, "01ABC-ampel-gr.md")
	for _, f := range []string{bad, keep} {
		if err := os.WriteFile(f, []byte(f), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	out, err := runCLI(t, dir, "validate")
	if exitCode(err) != ExitValidation {
		t.Fatalf("validate exit = %d, want %d\n%s", exitCode(err), ExitValidation, out)
	}
	if n := strings.Count(out, "not valid UTF-8"); n != 1 {
		t.Fatalf("want exactly one invalid_name line, got %d:\n%s", n, out)
	}
	const marker = "macOS cannot check it out — give "
	i := strings.Index(out, marker)
	if i < 0 {
		t.Fatalf("validate did not name the file:\n%s", out)
	}
	for _, f := range []string{bad, keep} {
		if b, err := os.ReadFile(f); err != nil || string(b) != f {
			t.Fatalf("validate touched %q: %v", f, err)
		}
	}

	// The quoted name, read by bash from another directory, is the file.
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("no bash to read the quoted name")
	}
	quoted, _, _ := strings.Cut(out[i+len(marker):], " a valid name")
	sh := exec.Command("bash", "-c", `printf %s `+quoted)
	sh.Dir = t.TempDir()
	got, err := sh.Output()
	if err != nil {
		t.Fatalf("bash cannot read %s: %v", quoted, err)
	}
	if string(got) != bad {
		t.Fatalf("bash reads %s as %q, want %q", quoted, got, bad)
	}
}
