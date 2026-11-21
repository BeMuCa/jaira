package board

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Only exactly "agent" or "none" on the first line that is neither blank nor
// a comment changes how work arrives; every other input class — missing,
// empty, comment only, case, typo, a second line — stays person, so a slip
// never sends work to the default branch unreviewed (XB97NJ).
func TestPRModeInTakesOnlyTheExactWord(t *testing.T) {
	for in, want := range map[string]string{
		"":                      PRPerson,
		"\n\n":                  PRPerson,
		"# only a comment\n":    PRPerson,
		"none":                  PRNone,
		"none\n":                PRNone,
		"none\r\n":              PRNone,
		"  none \t\n":           PRNone,
		"# solo\n\nnone\n":      PRNone,
		"agent\n":               PRAgent,
		"None\n":                PRPerson,
		"NONE\n":                PRPerson,
		"nnone\n":               PRPerson,
		"none # solo\n":         PRPerson,
		"person\nnone\n":        PRPerson,
		"\xef\xbb\xbfnone\n":    PRNone,
		"\xef\xbb\xbfagent\r\n": PRAgent,
		"none\n\xef\xbb\xbf":    PRNone,
		"none; git push -f\n":   PRPerson,
	} {
		if got := PRModeIn(in); got != want {
			t.Errorf("PRModeIn(%q) = %q, want %q", in, got, want)
		}
	}
	if got := PRModeFor(""); got != PRPerson {
		t.Errorf("PRModeFor(\"\") = %q, want person", got)
	}
	if got := PRModeFor(t.TempDir()); got != PRPerson {
		t.Errorf("PRModeFor(no .jaira) = %q, want person", got)
	}
}

func writePR(t *testing.T, root, value string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(root, ".jaira"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".jaira", "pr"), []byte(value), 0o644); err != nil {
		t.Fatal(err)
	}
}

// A board that says nothing, or says person, gets byte for byte the block it
// got before boards could say how work arrives: nobody who did nothing sees a
// change, and 'jaira validate' finds no drift (XB97NJ, DoD 4).
func TestNoteWithoutPRFileIsUnchanged(t *testing.T) {
	facts := []LaneFact{{ID: "todo"}, {ID: "done", Terminal: true}}
	for _, value := range []string{"", "person\n", "Agent\n"} {
		root := t.TempDir()
		if value != "" {
			writePR(t, root, value)
		}
		if got, want := noteFor(root, facts), agentNote+laneSection(facts); got != want {
			t.Errorf(".jaira/pr %q: the note changed", value)
		}
	}
}

// none and agent each write their own paragraph; a block written before the
// file changed is then out of date, and rewriting it brings it current.
func TestNoteSaysHowWorkArrivesAndDriftsWithTheFile(t *testing.T) {
	root := t.TempDir()
	writePR(t, root, "")
	if _, err := AnnounceInAgentFiles(root, nil); err != nil {
		t.Fatal(err)
	}
	if current, _ := NoteIsCurrent(root, nil); !current {
		t.Fatal("fresh block reported stale")
	}
	for _, c := range []struct{ value, says, never string }{
		{"none\n", "takes no pull requests", "may open that pull request"},
		{"agent\n", "may open that pull request", "takes no pull requests"},
	} {
		writePR(t, root, c.value)
		if current, stale := NoteIsCurrent(root, nil); current || len(stale) == 0 {
			t.Errorf("%q: block not reported stale after .jaira/pr changed", c.value)
		}
		if _, err := AnnounceInAgentFiles(root, nil); err != nil {
			t.Fatal(err)
		}
		for _, name := range []string{"AGENTS.md", "CLAUDE.md"} {
			b, err := os.ReadFile(filepath.Join(root, name))
			if err != nil {
				continue
			}
			s := string(b)
			if !strings.Contains(s, "## How work arrives on this board") || !strings.Contains(s, c.says) || strings.Contains(s, c.never) {
				t.Errorf("%q: %s does not say how work arrives:\n%s", c.value, name, s)
			}
			if strings.Count(s, "## How work arrives on this board") != 1 {
				t.Errorf("%q: %s carries the paragraph more than once", c.value, name)
			}
		}
		if current, _ := NoteIsCurrent(root, nil); !current {
			t.Errorf("%q: rewritten block still reported stale", c.value)
		}
	}
	// Back to nothing: the paragraph goes again.
	if err := os.Remove(filepath.Join(root, ".jaira", "pr")); err != nil {
		t.Fatal(err)
	}
	if _, err := AnnounceInAgentFiles(root, nil); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(filepath.Join(root, "CLAUDE.md"))
	if strings.Contains(string(b), "How work arrives") {
		t.Errorf("paragraph left behind after .jaira/pr was removed:\n%s", b)
	}
}
