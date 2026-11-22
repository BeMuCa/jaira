package board

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// crlf turns LF text into what a Windows checkout with core.autocrlf holds.
func crlf(s string) string { return strings.ReplaceAll(s, "\n", "\r\n") }

// bareLF reports whether s has a line ending that is not CRLF.
func bareLF(s string) bool { return strings.Contains(strings.ReplaceAll(s, "\r\n", ""), "\n") }

// An AGENTS.md checked out with CRLF line endings is the same file as its LF
// original. Its block, current in LF, must read as current in CRLF too, and an
// update must keep the file in CRLF rather than leaving the block in LF.
func TestACRLFAgentFileIsNotDrifted(t *testing.T) {
	for _, c := range []struct {
		name, before string
	}{
		{"current block", "# Project\n\n" + jairaMarkerStart + "\n" + noteFor(t.TempDir(), nil) + "\n" + jairaMarkerEnd + "\n\nAfter.\n"},
		{"stale block", "# Project\n\n" + jairaMarkerStart + "\nold text\n" + jairaMarkerEnd + "\n"},
		{"stale block with a local area", "# Project\n\n" + jairaMarkerStart + "\nold text\n\n" + jairaMarkerLocal + "\nmine\n" + jairaMarkerEnd + "\n"},
		{"no block yet", "# Project\n\nSome notes.\n"},
	} {
		t.Run(c.name, func(t *testing.T) {
			root := t.TempDir()
			path := filepath.Join(root, "AGENTS.md")
			if err := os.WriteFile(path, []byte(crlf(c.before)), 0o644); err != nil {
				t.Fatal(err)
			}
			if c.name == "current block" {
				if ok, stale := NoteIsCurrent(root, nil); !ok {
					t.Fatalf("a current block in CRLF reads as drifted: %v", stale)
				}
			}
			if _, _, err := announceInAgentFile(root, "AGENTS.md", nil, false); err != nil {
				t.Fatal(err)
			}
			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if bareLF(string(got)) {
				t.Errorf("the update left LF lines in a CRLF file:\n%q", got)
			}
			if ok, stale := NoteIsCurrent(root, nil); !ok {
				t.Errorf("drifted straight after the update: %v", stale)
			}
			if _, action, err := announceInAgentFile(root, "AGENTS.md", nil, false); err != nil || action != "unchanged" {
				t.Errorf("second update: action %q, err %v; want unchanged", action, err)
			}
			if strings.Contains(c.before, jairaMarkerLocal) && !strings.Contains(string(got), jairaMarkerLocal+"\r\nmine\r\n"+jairaMarkerEnd) {
				t.Errorf("the local area did not survive:\n%q", got)
			}
		})
	}
}

// One CRLF line pasted into an LF file does not make the file CRLF: the block
// stays in LF, as a build that writes only LF would write it too, so two
// builds on one board do not rewrite it back and forth.
func TestAStrayCRLFLineLeavesTheBlockInLF(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "AGENTS.md")
	before := "# Project\r\n\nSome notes.\n\n" + jairaMarkerStart + "\n" + noteFor(root, nil) + "\n" + jairaMarkerEnd + "\n"
	if err := os.WriteFile(path, []byte(before), 0o644); err != nil {
		t.Fatal(err)
	}
	if ok, stale := NoteIsCurrent(root, nil); !ok {
		t.Fatalf("an LF block reads as drifted because of one CRLF line elsewhere: %v", stale)
	}
	if _, action, err := announceInAgentFile(root, "AGENTS.md", nil, false); err != nil || action != "unchanged" {
		t.Errorf("update: action %q, err %v; want unchanged", action, err)
	}
}

// A CRLF file whose block an older build wrote in LF is mixed. The file's own
// lines say CRLF, so the update puts the block in CRLF too, however long the
// block is.
func TestAnLFBlockInACRLFFileBecomesCRLF(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "AGENTS.md")
	before := "# Project\r\n\r\nSome notes.\r\nMore.\r\n\n" + jairaMarkerStart + "\nold text\n" + strings.Repeat("old line\n", 20) + jairaMarkerEnd + "\n"
	if err := os.WriteFile(path, []byte(before), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := announceInAgentFile(root, "AGENTS.md", nil, false); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	start, end, _ := blockBounds(string(got))
	if bareLF(string(got)[start:end]) {
		t.Errorf("the block stayed in LF in a CRLF file:\n%q", got)
	}
	if ok, stale := NoteIsCurrent(root, nil); !ok {
		t.Errorf("drifted straight after the update: %v", stale)
	}
}
