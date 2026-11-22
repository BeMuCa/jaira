package cli

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestUpdateNamesFilesWithInvalidUTF8 plants the name an older jaira wrote —
// a slug cut through "ü", leaving the lone byte \xc3 — and checks that update
// names it, offers a rename bash can run as printed, and leaves the file where
// it is.
func TestUpdateNamesFilesWithInvalidUTF8(t *testing.T) {
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		t.Skip("this file system cannot hold a name that is not valid UTF-8")
	}
	setCurrent(t, "1.0.0")
	dir, _ := updateTestStore(t)
	logbook := filepath.Join(dir, ".jaira", "logbook", "someone-2026-11-22")
	if err := os.MkdirAll(logbook, 0o755); err != nil {
		t.Fatal(err)
	}
	bad := filepath.Join(logbook, "01ABC-ampel-gr\xc3.md")
	if err := os.WriteFile(bad, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	out, err := runUpdate(t, dir)
	if err != nil {
		t.Fatalf("jaira update: %v\n%s", err, out)
	}
	var line string
	for _, l := range strings.Split(out, "\n") {
		if strings.Contains(l, "macOS cannot check it out — rename it: git mv") {
			line = l
		}
	}
	if line == "" {
		t.Fatalf("update did not name the file:\n%s", out)
	}
	if _, err := os.Stat(bad); err != nil {
		t.Fatalf("update touched the file: %v", err)
	}

	// The printed rename, run by bash as printed, must hit the real file.
	cmd := line[strings.Index(line, "git mv ")+len("git mv "):]
	if _, err := exec.LookPath("bash"); err == nil {
		sh := exec.Command("bash", "-c", "mv "+cmd)
		sh.Dir = dir
		if b, err := sh.CombinedOutput(); err != nil {
			t.Fatalf("the printed rename does not run: %v\n%s\n%s", err, line, b)
		}
		if _, err := os.Stat(filepath.Join(logbook, "01ABC-ampel-gr.md")); err != nil {
			t.Fatalf("rename did not produce the valid name: %v", err)
		}
	}

	out, err = runUpdate(t, dir, "--json")
	if err != nil {
		t.Fatal(err)
	}
	var payload struct {
		InvalidNames []string `json:"invalid_names"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatal(err)
	}
	if len(payload.InvalidNames) != 0 {
		t.Fatalf("after the rename invalid_names = %q, want none", payload.InvalidNames)
	}
}

// TestUpdateQuietOnValidNames: a board with only valid names prints no such
// line and reports an empty list, not null.
func TestUpdateQuietOnValidNames(t *testing.T) {
	setCurrent(t, "1.0.0")
	dir, _ := updateTestStore(t)
	out, err := runUpdate(t, dir, "--json")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"invalid_names": []`) {
		t.Fatalf("want an empty invalid_names list:\n%s", out)
	}
	if out, _ := runUpdate(t, dir); strings.Contains(out, "macOS cannot check it out — rename it: git mv") {
		t.Fatalf("named a valid file:\n%s", out)
	}
}
