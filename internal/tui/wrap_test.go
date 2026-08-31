package tui

// A shrunken terminal has no horizontal scroll, so a sentence cut at the
// right edge is unreadable. Flowing text — messages, paths, prompts, field
// values — wraps to the available width; only deliberate one-liners (cards,
// column heads, status bars) may truncate.

import (
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/project"
)

// A path is one long word: wrap must hard-break it rather than hand back a
// line wider than the terminal.
func TestWrapBreaksOverlongWords(t *testing.T) {
	path := "/home/someone/git/organisation/team/very-long-repository-name/subdir"
	out := wrap(path, 24, 4)
	lines := strings.Split(out, "\n")
	if len(lines) < 2 {
		t.Fatalf("width 24 must force a break, got %q", out)
	}
	joined := ""
	for i, l := range lines {
		content := l
		if i > 0 {
			content = strings.TrimPrefix(l, "    ")
		}
		if w := len([]rune(content)); w > 24 {
			t.Errorf("line %d wider than 24: %q", i, l)
		}
		joined += content
	}
	if joined != path {
		t.Errorf("hard break lost characters:\n got %q\nwant %q", joined, path)
	}
}

// Sentences still break at spaces; a word that fits is never split.
func TestWrapStillBreaksAtSpaces(t *testing.T) {
	out := wrap("the quick brown fox jumps over the lazy dog", 16, 0)
	for _, l := range strings.Split(out, "\n") {
		if len([]rune(l)) > 16 {
			t.Errorf("line wider than 16: %q", l)
		}
	}
	if strings.Contains(out, "qui\nck") || len(strings.Split(out, "\n")) < 3 {
		t.Errorf("unexpected wrapping:\n%s", out)
	}
}

// wrapLines wraps each line of a multi-line text and keeps the line
// structure — paragraphs and blank lines survive.
func TestWrapLinesKeepsLineStructure(t *testing.T) {
	in := "first paragraph that is clearly longer than the width\n\nsecond"
	out := wrapLines(in, 20)
	for _, l := range strings.Split(out, "\n") {
		if len([]rune(l)) > 20 {
			t.Errorf("line wider than 20: %q", l)
		}
	}
	if !strings.Contains(out, "\n\n") {
		t.Errorf("blank line between paragraphs was lost:\n%s", out)
	}
	if !strings.Contains(out, "second") {
		t.Errorf("second paragraph lost:\n%s", out)
	}
}

// An error message on a narrow terminal wraps; the tail of the sentence must
// not be cut off.
func TestRenderMessageWrapsAtNarrowWidth(t *testing.T) {
	m := newTestModel(t, 40, 24)
	m.message = "cannot advance: lane critique declares it produces review-summary, which is still empty on this ticket"
	m.isErr = true
	out := stripANSI(m.renderMessage())
	if strings.Contains(out, "…") {
		t.Errorf("message truncated instead of wrapped:\n%s", out)
	}
	if !strings.Contains(strings.ReplaceAll(out, "\n", " "), "still empty on this ticket") {
		t.Errorf("tail of the message is missing:\n%s", out)
	}
	for _, l := range strings.Split(out, "\n") {
		if len([]rune(l)) > 40 {
			t.Errorf("line wider than the terminal: %q", l)
		}
	}
}

// The board switcher shows each project's full root path, hard-broken over
// lines rather than cut with an ellipsis.
func TestProjectPathsWrapInsteadOfTruncate(t *testing.T) {
	m := newTestModel(t, 36, 24)
	m.projects = []project.Project{{Name: "deep", Root: "/home/someone/git/organisation/team/very-long-repository-name"}}
	out := stripANSI(m.renderProjects())
	if strings.Contains(out, "…") {
		t.Errorf("path truncated instead of wrapped:\n%s", out)
	}
	squashed := strings.ReplaceAll(strings.ReplaceAll(out, "\n", ""), " ", "")
	if !strings.Contains(squashed, "very-long-repository-name") {
		t.Errorf("tail of the path is missing:\n%s", out)
	}
	for _, l := range strings.Split(out, "\n") {
		if n := len([]rune(l)); n > 36 {
			t.Errorf("line wider than the terminal (%d): %q", n, l)
		}
	}
}

// Sign-off is where a person reads carefully: checklist items and their
// proofs wrap instead of vanishing off the right edge.
func TestSignOffWrapsChecklistText(t *testing.T) {
	m := newTestModel(t, 38, 60)
	tk := longTicket()
	tk.DoDItems[0].Text = "every screen of the TUI wraps long sentences to the terminal width instead of cutting them off"
	tk.DoDItems[0].Proof = "internal/tui/wrap_test.go plus a hand check at 40 columns"
	m.detail = tk
	out := stripANSI(m.renderSignOff())
	if strings.Contains(out, "…") {
		t.Errorf("sign-off truncated instead of wrapped:\n%s", out)
	}
	joined := strings.ReplaceAll(out, "\n", " ")
	if !strings.Contains(joined, "instead") {
		t.Errorf("tail of the checklist item is missing:\n%s", out)
	}
	for _, l := range strings.Split(out, "\n") {
		if n := len([]rune(l)); n > 38 {
			t.Errorf("line wider than the terminal (%d): %q", n, l)
		}
	}
}

// A fitting line passes through wrapLines untouched — pre-aligned text like
// a git --stat keeps its column alignment.
func TestWrapLinesLeavesFittingLinesAlone(t *testing.T) {
	in := " internal/tui/view.go         |  67 +++"
	if out := wrapLines(in, 80); out != in {
		t.Errorf("fitting line was rewritten:\n got %q\nwant %q", out, in)
	}
}

// Below wrap's readable minimum an over-wide line is still broken hard
// rather than handed back wider than the pane.
func TestWrapTinyWidthStillBreaksHard(t *testing.T) {
	for _, l := range strings.Split(wrap("abcdefghijklmnop", 6, 0), "\n") {
		if len([]rune(l)) > 6 {
			t.Errorf("line wider than 6: %q", l)
		}
	}
}
