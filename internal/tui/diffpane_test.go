package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/BeMuCa/jaira/core/ticket"
)

const samplePatch = "diff --git a/internal/a.go b/internal/a.go\n" +
	"index 1111111..2222222 100644\n" +
	"--- a/internal/a.go\n" +
	"+++ b/internal/a.go\n" +
	"@@ -10,3 +10,4 @@ func a() {\n" +
	" \tx := 1\n" +
	"-\ty := 2\n" +
	"+\ty := 3\n" +
	"+\tz := 4\n" +
	" \treturn x\n" +
	"diff --git a/docs/b.md b/docs/b.md\n" +
	"new file mode 100644\n" +
	"--- /dev/null\n" +
	"+++ b/docs/b.md\n" +
	"@@ -0,0 +1,2 @@\n" +
	"+first\n" +
	"+second\n"

// waitingTicket is a ticket in the given lane, with enough filled in that the
// sign-off screen has something to lay out.
func waitingTicket(status string) *ticket.Ticket {
	return &ticket.Ticket{
		ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "Rate limit login", Status: status,
		Goal:    "stop credential stuffing",
		Context: "came up while reading the auth logs",
		Outcome: ticket.Outcome{What: "added a token bucket per IP", Why: "no limit", Resolves: "429 above 100/min"},
	}
}

// diffModel opens the ticket with its diff already loaded, so no git runs.
func diffModel(t *testing.T, w int, status string) *Model {
	t.Helper()
	m := newTestModel(t, w, 40)
	tk := waitingTicket(status)
	m.detail = tk
	m.mode = modeDetail
	m.diffFor, m.diffScrollFor, m.diffRows, m.diffNote = tk.ID, tk.ID, parseDiff(samplePatch), "feat/x vs master"
	return m
}

func TestParseDiffSplitsFilesAndNumbersLines(t *testing.T) {
	rows := parseDiff(samplePatch)
	var headers []diffRow
	var added diffRow
	for _, r := range rows {
		if r.kind == 'H' {
			headers = append(headers, r)
		}
		if r.kind == '+' && strings.Contains(r.text, "y := 3") {
			added = r
		}
	}
	if len(headers) != 2 || headers[0].path != "internal/a.go" || headers[1].path != "docs/b.md" {
		t.Fatalf("headers = %+v, want internal/a.go then docs/b.md", headers)
	}
	if !strings.Contains(headers[0].text, "+2 −1") || !strings.Contains(headers[1].text, "+2 −0") {
		t.Errorf("header counts wrong: %q / %q", headers[0].text, headers[1].text)
	}
	if headers[0].line != 10 {
		t.Errorf("first file header line = %d, want 10 (first hunk's new start)", headers[0].line)
	}
	if added.line != 11 || added.file != 0 {
		t.Errorf("added row = %+v, want line 11 of file 0", added)
	}
}

func TestParseDiffDropsCommitNoise(t *testing.T) {
	fallback := "commit 1234abcd\nsubject one\n\n a.txt | 1 +\n 1 file changed, 1 insertion(+)\n\n" +
		"diff --git a/a.txt b/a.txt\nindex 0..1 100644\n--- a/a.txt\n+++ b/a.txt\n" +
		"@@ -1 +1,2 @@\n keep\n+added\n\n" +
		"commit 5678abcd\nsubject two\n\n b.txt | 1 -\n\n" +
		"diff --git a/b.txt b/b.txt\n--- a/b.txt\n+++ b/b.txt\n@@ -1 +0,0 @@\n-gone\n"
	rows := parseDiff(fallback)
	headers := 0
	for _, r := range rows {
		if r.kind == 'H' {
			headers++
		}
		for _, bad := range []string{"subject", "| 1", "commit", "file changed"} {
			if strings.Contains(r.text, bad) {
				t.Errorf("row %+v carries commit noise %q", r, bad)
			}
		}
	}
	if headers != 2 {
		t.Errorf("headers = %d, want 2", headers)
	}
}

func TestPersonLaneShowsDiffBesideTheTicketWhenWide(t *testing.T) {
	m := diffModel(t, 160, "signoff")
	out := stripANSI(m.render())
	for _, want := range []string{"internal/a.go", "+    y := 3", "o open in VS Code"} {
		if !strings.Contains(out, want) {
			t.Errorf("wide sign-off screen misses %q", want)
		}
	}
	// The ticket's handle line and the diff's title share a row: the two panes
	// sit side by side rather than one above the other.
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "7BVRCT") && !strings.Contains(line, "Diff  feat/x") {
			t.Errorf("the ticket and the diff are not side by side: %q", line)
		}
	}
}

func TestHumanLaneShowsDiffBesideTheTicketWhenWide(t *testing.T) {
	m := diffModel(t, 160, "human")
	out := stripANSI(m.render())
	if !strings.Contains(out, "+    y := 3") || !strings.Contains(out, "o open in VS Code") {
		t.Errorf("human lane screen misses the diff pane:\n%s", out)
	}
}

func TestNarrowTerminalKeepsTheSingleColumn(t *testing.T) {
	m := diffModel(t, 120, "signoff")
	out := stripANSI(m.render())
	if strings.Contains(out, "open in VS Code") || strings.Contains(out, "tab diff") {
		t.Errorf("narrow terminal shows the diff pane:\n%s", out)
	}
	m2 := diffModel(t, 120, "signoff")
	hint := "a accept → done · f follow-up ticket · e edit · E editor · esc back"
	want := stripANSI(m2.clipToWindow(m2.signOffBody(120), hint))
	if out != want {
		t.Errorf("narrow sign-off screen changed:\n got: %q\nwant: %q", out, want)
	}
}

func TestBracketJumpsBetweenFiles(t *testing.T) {
	m := diffModel(t, 160, "signoff")
	m.render()
	press := func(s string) {
		m.Update(tea.KeyPressMsg{Code: rune(s[0]), Text: s})
		m.render()
	}
	top := func() diffRow { return m.diffRows[m.diffScroll] }

	press("]")
	if got := top(); got.kind != 'H' || got.path != "docs/b.md" {
		t.Fatalf("after ] the top row is %+v, want the docs/b.md header", got)
	}
	press("[")
	if got := top(); got.kind != 'H' || got.path != "internal/a.go" {
		t.Fatalf("after [ the top row is %+v, want the internal/a.go header", got)
	}
}

// The 2s tick reloads the board. A reload must not move a diff the person has
// scrolled, nor take the focus off it: it only marks the diff to be read again.
func TestReloadKeepsDiffScrollAndFocus(t *testing.T) {
	m := diffModel(t, 160, "signoff")
	m.diffScroll, m.diffFocus = 3, true
	if err := m.reload(); err != nil {
		t.Fatal(err)
	}
	if m.diffScroll != 3 || !m.diffFocus {
		t.Fatalf("after reload scroll=%d focus=%v, want 3 and true", m.diffScroll, m.diffFocus)
	}
	m.loadDiff(m.detail)
	if m.diffScroll != 3 || !m.diffFocus {
		t.Fatalf("after reloading the same ticket's diff scroll=%d focus=%v, want 3 and true", m.diffScroll, m.diffFocus)
	}
}

func TestOpeningAnotherTicketResetsDiffScroll(t *testing.T) {
	m := diffModel(t, 160, "signoff")
	m.diffScroll, m.diffFocus = 3, true
	other := waitingTicket("signoff")
	other.ID = "01KZTT3XZ2YQBX93TTSR7BVRCZ"
	m.loadDiff(other)
	if m.diffScroll != 0 || m.diffFocus || m.diffScrollFor != other.ID {
		t.Fatalf("after opening another ticket scroll=%d focus=%v for=%q, want 0, false, %q",
			m.diffScroll, m.diffFocus, m.diffScrollFor, other.ID)
	}
}

func TestOKeyOpensFileAtLine(t *testing.T) {
	var gotDir, gotPath string
	var gotLine int
	orig := openInCode
	openInCode = func(dir, path string, line int) error {
		gotDir, gotPath, gotLine = dir, path, line
		return nil
	}
	t.Cleanup(func() { openInCode = orig })

	m := diffModel(t, 160, "signoff")
	m.diffDir = "/repo/worktree"
	m.render()
	m.Update(tea.KeyPressMsg{Code: 'o', Text: "o"})
	if gotPath != "internal/a.go" || gotLine != 10 || gotDir != "/repo/worktree" {
		t.Errorf("openInCode(%q, %q, %d), want (/repo/worktree, internal/a.go, 10)", gotDir, gotPath, gotLine)
	}
}
