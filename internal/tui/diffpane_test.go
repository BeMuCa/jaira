package tui

import (
	"os/exec"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

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

// pairPatch replaces one line with one line, so the removed and the added row
// pair for a word highlight.
const pairPatch = "diff --git a/internal/a.go b/internal/a.go\n" +
	"@@ -10,3 +10,3 @@ func a() {\n" +
	" \tx := 1\n" +
	"-\ty := 2\n" +
	"+\ty := 3\n" +
	" \treturn x\n"

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
	m.diffFor, m.diffScrollFor, m.diffRows, m.diffNote = tk.ID, tk.ID, parseDiff(samplePatch), "feat/x → master"
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
	if headers[0].adds != 2 || headers[0].dels != 1 || headers[1].adds != 2 || headers[1].dels != 0 {
		t.Errorf("header counts wrong: %+v / %+v", headers[0], headers[1])
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

// A successful open must not raise the message modal: it would take the next
// keys. The confirmation sits in the diff title and clears on the next key.
func TestOKeyConfirmsInTheDiffTitleWithoutAModal(t *testing.T) {
	orig := openInCode
	openInCode = func(string, string, int) error { return nil }
	t.Cleanup(func() { openInCode = orig })

	m := diffModel(t, 160, "signoff")
	m.diffDir = "/repo/worktree"
	m.render()
	m.Update(tea.KeyPressMsg{Code: 'o', Text: "o"})
	if m.mode != modeDetail {
		t.Fatalf("mode = %v after a successful open, want modeDetail", m.mode)
	}
	if out := stripANSI(m.render()); !strings.Contains(out, "opened internal/a.go:10") {
		t.Errorf("diff title misses the open confirmation:\n%s", out)
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if out := stripANSI(m.render()); strings.Contains(out, "opened internal/a.go:10") {
		t.Errorf("the open confirmation survives the next key:\n%s", out)
	}
}

// parseDiff numbers both sides: a removed row carries its old-side number, an
// added row its new-side one, and a context row both.
func TestParseDiffNumbersBothSides(t *testing.T) {
	var removed, added, context diffRow
	for _, r := range parseDiff(samplePatch) {
		switch {
		case r.kind == '-' && strings.Contains(r.text, "y := 2"):
			removed = r
		case r.kind == '+' && strings.Contains(r.text, "y := 3"):
			added = r
		case r.kind == ' ' && strings.Contains(r.text, "return x"):
			context = r
		}
	}
	if removed.old != 11 || added.line != 11 || added.old != 0 {
		t.Errorf("removed %+v, added %+v: want old 11 on the removed row and new 11 on the added one", removed, added)
	}
	if context.old != 12 || context.line != 13 {
		t.Errorf("context row %+v, want old 12 and new 13", context)
	}
}

func TestParseDiffTagsNewDeletedAndRenamedFiles(t *testing.T) {
	patch := "diff --git a/n.txt b/n.txt\nnew file mode 100644\n--- /dev/null\n+++ b/n.txt\n@@ -0,0 +1 @@\n+x\n" +
		"diff --git a/d.txt b/d.txt\ndeleted file mode 100644\n--- a/d.txt\n+++ /dev/null\n@@ -1 +0,0 @@\n-y\n" +
		"diff --git a/old.txt b/new.txt\nsimilarity index 100%\nrename from old.txt\nrename to new.txt\n"
	var headers []diffRow
	for _, r := range parseDiff(patch) {
		if r.kind == 'H' {
			headers = append(headers, r)
		}
	}
	want := []string{" new", " deleted", " renamed"}
	if len(headers) != len(want) {
		t.Fatalf("got %d headers, want %d", len(headers), len(want))
	}
	for i, h := range headers {
		if h.tag != want[i] {
			t.Errorf("header %d tag = %q, want %q", i, h.tag, want[i])
		}
	}
	if headers[1].path != "d.txt" {
		t.Errorf("deleted file path = %q, want d.txt", headers[1].path)
	}
}

// spanText is the part of a row that the word highlight paints.
func spanText(r diffRow) string {
	rs := []rune(r.text)
	return string(rs[r.from:r.to])
}

func TestWordPairsMarkOnlyTheChangedMiddle(t *testing.T) {
	var removed, added diffRow
	for _, r := range parseDiff(pairPatch) {
		if r.kind == '-' {
			removed = r
		}
		if r.kind == '+' && strings.Contains(r.text, "y := 3") {
			added = r
		}
	}
	if spanText(removed) != "2" || spanText(added) != "3" {
		t.Errorf("word spans %q and %q, want \"2\" and \"3\"", spanText(removed), spanText(added))
	}

	// Lines that share less than three runes are rewritten, not edited: no span.
	for _, r := range parseDiff("@@ -1 +1 @@\n-abc\n+xyz\n") {
		if r.kind == '-' || r.kind == '+' {
			if r.from != r.to {
				t.Errorf("row %q has a word span although nothing is shared", r.text)
			}
		}
	}
	// Unequal runs are not paired: "-foo bar" would pair with "+foo qux" alone.
	for _, r := range parseDiff("@@ -1,2 +1 @@\n-foo bar\n-foo baz\n+foo qux\n") {
		if r.kind == '-' || r.kind == '+' {
			if r.from != r.to {
				t.Errorf("row %q has a word span although its run is unequal", r.text)
			}
		}
	}
}

// bgOpen is the escape a background of hex opens with, as lipgloss writes it.
func bgOpen(hex string) string {
	s := lipgloss.NewStyle().Background(lipgloss.Color(hex)).Render("x")
	return s[:strings.Index(s, "x")]
}

// Changed rows carry their tint across the whole pane, and the changed words
// carry the stronger one; a context row carries no tint at all.
func TestChangedRowsAreTintedAcrossThePane(t *testing.T) {
	m := diffModel(t, 160, "signoff")
	m.darkBG = true
	m.diffRows = parseDiff(pairPatch)
	width := m.width - diffLeftWidth - 5
	add, del, wordAdd := bgOpen("#16301f"), bgOpen("#3a1d22"), bgOpen("#2b6a3c")
	for _, line := range strings.Split(m.renderDiffPane(width, 20), "\n") {
		plain := stripANSI(line)
		switch {
		case strings.Contains(plain, "y := 3"):
			if !strings.Contains(line, add+"    y := ") || !strings.Contains(line, wordAdd+"3") {
				t.Errorf("added row lacks its tint or its word highlight:\n%q", line)
			}
		case strings.Contains(plain, "y := 2"):
			if !strings.Contains(line, del) {
				t.Errorf("removed row lacks its tint:\n%q", line)
			}
		case strings.Contains(plain, "return x"):
			if strings.Contains(line, "48;") {
				t.Errorf("context row is tinted:\n%q", line)
			}
		default:
			continue
		}
		if lipgloss.Width(line) != width {
			t.Errorf("row is %d wide, want %d: %q", lipgloss.Width(line), width, plain)
		}
	}
}

// The gutter's old-side number sits in the first columns and the new-side one
// after them, so a removed row shows only the first and an added row only the
// second.
func TestGutterShowsOldAndNewNumbersOnTheirSides(t *testing.T) {
	m := diffModel(t, 160, "signoff")
	width := m.width - diffLeftWidth - 5
	lines := strings.Split(stripANSI(m.renderDiffPane(width, 20)), "\n")
	find := func(text string) []rune {
		for _, l := range lines {
			if strings.Contains(l, text) {
				return []rune(l)
			}
		}
		t.Fatalf("no row shows %q", text)
		return nil
	}
	removed, added, context := find("y := 2"), find("y := 3"), find("return x")
	if string(removed[0:3]) != " 11" || string(removed[4:7]) != "   " || removed[10] != '-' {
		t.Errorf("removed row gutter %q, want old 11 and no new number", string(removed[:10]))
	}
	if string(added[0:3]) != "   " || string(added[4:7]) != " 11" || added[10] != '+' {
		t.Errorf("added row gutter %q, want no old number and new 11", string(added[:10]))
	}
	if string(context[0:3]) != " 12" || string(context[4:7]) != " 13" {
		t.Errorf("context row gutter %q, want old 12 and new 13", string(context[:10]))
	}
}

func TestDiffTitleCountsFilesAndTotals(t *testing.T) {
	m := diffModel(t, 160, "signoff")
	if out := stripANSI(m.render()); !strings.Contains(out, "file 1/2  +4 −1") || !strings.Contains(out, "○ Diff  feat/x → master") {
		t.Errorf("title misses the counter or the unfocused mark:\n%s", out)
	}
	m.diffScroll = 8 // the header of docs/b.md
	m.diffFocus = true
	if out := stripANSI(m.render()); !strings.Contains(out, "file 2/2  +4 −1") || !strings.Contains(out, "● Diff  feat/x → master") {
		t.Errorf("title after scrolling to file 2 misses the counter or the focus mark:\n%s", out)
	}
}

// A notice sits below the title, one blank row above it, centred.
func TestEmptyStateIsNoticeBelowTheTitle(t *testing.T) {
	m := diffModel(t, 160, "signoff")
	m.diffRows = nil
	m.diffNotice = paneNotice{head: "No changes", lines: []string{"feat/x has nothing that master lacks"}}
	lines := strings.Split(stripANSI(m.renderDiffPane(77, 10)), "\n")
	if len(lines) < 4 || lines[1] != strings.Repeat(" ", 77) || !strings.Contains(lines[2], "No changes") || !strings.Contains(lines[3], "nothing that master lacks") {
		t.Errorf("notice rows are wrong:\n%s", strings.Join(lines, "\n"))
	}
	if strings.Contains(strings.Join(lines, "\n"), "nothing to show") {
		t.Errorf("the old empty text is still shown")
	}
}

// runGit runs git in dir for a test that needs a real repository.
func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir, "-c", "user.name=t", "-c", "user.email=t@t"}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

// ticketDiff says why a ticket has no rows: no branch names it and no commit
// does, or the branch that names it is already level with the base.
func TestTicketDiffSaysWhyThereAreNoRows(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init", "-q")
	runGit(t, root, "commit", "-q", "--allow-empty", "-m", "base")
	runGit(t, root, "branch", "-M", "master")
	tk := waitingTicket("signoff")
	handle := ticket.Handle(tk.ID)

	_, _, _, notice := ticketDiff(root, tk)
	if notice.head != "No diff yet" || len(notice.lines) == 0 || !strings.Contains(notice.lines[0], "no branch named after "+handle) {
		t.Errorf("no branch: notice = %+v", notice)
	}
	runGit(t, root, "branch", "feat/"+strings.ToLower(handle)+"-x")
	_, _, _, notice = ticketDiff(root, tk)
	if notice.head != "No changes" || len(notice.lines) == 0 || !strings.Contains(notice.lines[0], "feat/") {
		t.Errorf("level branch: notice = %+v", notice)
	}
}

// A blank row sits above every file header but the first, and jumping to a
// file still lands on its header.
func TestBlankRowSeparatesFilesAndJumpsStillLandOnHeaders(t *testing.T) {
	rows := parseDiff(samplePatch)
	if len(rows) < 9 || rows[7].kind != 'B' || rows[8].kind != 'H' {
		t.Fatalf("want a blank row above the second header, got %+v", rows)
	}
	m := diffModel(t, 160, "signoff")
	m.diffScroll = 0
	m.jumpFile(true)
	if m.diffRows[m.diffScroll].kind != 'H' || m.diffRows[m.diffScroll].path != "docs/b.md" {
		t.Errorf("] from the first file lands on %+v, want the docs/b.md header", m.diffRows[m.diffScroll])
	}
}

// A code row longer than the pane ends in "…" and stays exactly one row wide.
func TestLongCodeRowIsCutToThePane(t *testing.T) {
	m := diffModel(t, 160, "signoff")
	m.darkBG = true
	m.diffRows = parseDiff("diff --git a/x.txt b/x.txt\n@@ -1 +1 @@\n+" + strings.Repeat("x", 200) + "\n")
	width := 40
	for _, line := range strings.Split(m.renderDiffPane(width, 6), "\n") {
		if strings.Contains(line, "xxx") {
			if lipgloss.Width(line) != width || !strings.HasSuffix(stripANSI(line), "…") {
				t.Errorf("long row is %d wide, want %d and ending in …: %q", lipgloss.Width(line), width, stripANSI(line))
			}
			return
		}
	}
	t.Error("the long row is not in the pane")
}
