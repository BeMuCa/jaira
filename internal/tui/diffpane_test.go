package tui

import (
	"errors"
	"io/fs"
	"os/exec"
	"slices"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/BeMuCa/jaira/core/settings"
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
	m.diffOpen = -1
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
	m.diffOpen = 0
	out := stripANSI(m.render())
	for _, want := range []string{"internal/a.go", "+    y := 3", "o open in IDE"} {
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
	m.diffOpen = 0
	out := stripANSI(m.render())
	if !strings.Contains(out, "+    y := 3") || !strings.Contains(out, "o open in IDE") {
		t.Errorf("human lane screen misses the diff pane:\n%s", out)
	}
}

func TestNarrowTerminalKeepsTheSingleColumn(t *testing.T) {
	m := diffModel(t, 120, "signoff")
	out := stripANSI(m.render())
	if strings.Contains(out, "open in IDE") || strings.Contains(out, "tab diff") {
		t.Errorf("narrow terminal shows the diff pane:\n%s", out)
	}
	m2 := diffModel(t, 120, "signoff")
	hint := "a accept → done · m send back · f follow-up ticket · e edit · E editor · esc back"
	want := stripANSI(m2.clipToWindow(m2.signOffBody(120), hint))
	if out != want {
		t.Errorf("narrow sign-off screen changed:\n got: %q\nwant: %q", out, want)
	}
}

// The diff opens as a list with one row per file, each path from the repository
// root and no blank rows between them.
func TestDiffOpensAsAListOfFiles(t *testing.T) {
	m := diffModel(t, 160, "signoff")
	lines := strings.Split(stripANSI(m.renderDiffPane(77, 8)), "\n")
	if !strings.Contains(lines[1], "internal/a.go") || !strings.Contains(lines[1], "+2") || !strings.Contains(lines[1], "−1") {
		t.Errorf("first file row %q misses its path or counts", lines[1])
	}
	if !strings.Contains(lines[2], "docs/b.md") || !strings.Contains(lines[2], "+2") {
		t.Errorf("second file row %q misses its path or counts", lines[2])
	}
	if strings.TrimSpace(lines[3]) != "" || strings.Contains(strings.Join(lines, "\n"), "y := 3") {
		t.Errorf("the list shows more than one row per file:\n%s", strings.Join(lines, "\n"))
	}
}

// Tab gives the diff the keys. Down selects the next file, enter opens it, and a
// second enter closes it, after which the list can be moved through again.
func TestEnterOpensAndClosesTheSelectedFile(t *testing.T) {
	m := diffModel(t, 160, "signoff")
	m.render()
	m.Update(key("tab"))
	m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if m.diffSel != 1 || m.diffOpen != -1 {
		t.Fatalf("down: sel=%d open=%d, want the second file selected and none open", m.diffSel, m.diffOpen)
	}
	m.Update(tea.KeyPressMsg{Code: ']', Text: "]"})
	if m.diffSel != 1 {
		t.Errorf("] moved the selection to %d; it is no longer a key", m.diffSel)
	}
	m.Update(key("enter"))
	if m.diffOpen != 1 {
		t.Fatalf("enter: open=%d, want 1", m.diffOpen)
	}
	if out := stripANSI(m.render()); !strings.Contains(out, "+first") || strings.Contains(out, "y := 3") {
		t.Errorf("open docs/b.md should show its own rows and not internal/a.go's:\n%s", out)
	}
	m.Update(key("enter"))
	if m.diffOpen != -1 || m.detail == nil {
		t.Fatalf("second enter: open=%d detail=%v, want the file closed and the ticket still open", m.diffOpen, m.detail != nil)
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	if m.diffSel != 0 {
		t.Errorf("up after closing: sel=%d, want the list selectable again", m.diffSel)
	}
}

// Without the diff's focus enter keeps its meaning: it closes the ticket.
func TestEnterWithoutDiffFocusStillClosesTheTicket(t *testing.T) {
	m := diffModel(t, 160, "signoff")
	m.render()
	m.Update(key("enter"))
	if m.detail != nil {
		t.Error("enter with the list not focused left the ticket open")
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
	var gotDir string
	var gotArgv []string
	orig := openInCode
	openInCode = func(dir string, argv []string) (func() error, error) {
		gotDir, gotArgv = dir, argv
		return func() error { return nil }, nil
	}
	t.Cleanup(func() { openInCode = orig })

	m := diffModel(t, 160, "signoff")
	m.diffDir = "/repo/worktree"
	m.render()
	m.Update(tea.KeyPressMsg{Code: 'o', Text: "o"})
	want := []string{"code", "-g", "internal/a.go:10"}
	if !slices.Equal(gotArgv, want) || gotDir != "/repo/worktree" {
		t.Errorf("openInCode(%q, %q), want (/repo/worktree, %q)", gotDir, gotArgv, want)
	}
}

// A successful open must not raise the message modal: it would take the next
// keys. The title says it is opening, then that it opened, and the next key
// clears it.
func TestOKeyConfirmsInTheDiffTitleWithoutAModal(t *testing.T) {
	orig := openInCode
	openInCode = func(string, []string) (func() error, error) { return func() error { return nil }, nil }
	t.Cleanup(func() { openInCode = orig })

	m := diffModel(t, 160, "signoff")
	m.diffDir = "/repo/worktree"
	m.render()
	_, cmd := m.Update(tea.KeyPressMsg{Code: 'o', Text: "o"})
	if m.mode != modeDetail {
		t.Fatalf("mode = %v after o, want modeDetail", m.mode)
	}
	if out := stripANSI(m.render()); !strings.Contains(out, "opening internal/a.go:10") {
		t.Errorf("title misses the opening state:\n%s", out)
	}
	if cmd == nil {
		t.Fatal("o returned no command to wait for the IDE")
	}
	m.Update(cmd())
	if m.mode != modeDetail {
		t.Fatalf("mode = %v after the IDE returned, want modeDetail", m.mode)
	}
	if out := stripANSI(m.render()); !strings.Contains(out, "opened internal/a.go:10") {
		t.Errorf("diff title misses the open confirmation:\n%s", out)
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if out := stripANSI(m.render()); strings.Contains(out, "opened internal/a.go:10") {
		t.Errorf("the open confirmation survives the next key:\n%s", out)
	}
}

// An IDE that exits with an error shows the exit status and the last line of its
// stderr in the title, in red, rather than saying nothing happened.
func TestOKeyShowsTheIDEsFailureInTheTitle(t *testing.T) {
	orig := openInCode
	openInCode = func(string, []string) (func() error, error) {
		return func() error { return errors.New("exit status 127: remote-cli/code: not found") }, nil
	}
	t.Cleanup(func() { openInCode = orig })

	m := diffModel(t, 160, "signoff")
	m.diffDir = "/repo/worktree"
	m.render()
	_, cmd := m.Update(tea.KeyPressMsg{Code: 'o', Text: "o"})
	m.Update(cmd())
	if m.mode != modeDetail || !m.diffOpenErr {
		t.Fatalf("after a failed open mode=%v err=%v, want modeDetail and the error flag", m.mode, m.diffOpenErr)
	}
	if out := stripANSI(m.render()); !strings.Contains(out, "exit status 127: remote-cli/code: not found") {
		t.Errorf("diff title misses the IDE's failure:\n%s", out)
	}
}

// An IDE reply for another ticket must not write into the open ticket's title.
func TestDiffOpenedForAnotherTicketLeavesDiffOpenedUnchanged(t *testing.T) {
	m := diffModel(t, 160, "signoff")
	m.render()
	m.diffOpened, m.diffOpenErr = "opening internal/a.go:10…", false
	m.Update(diffOpenedMsg{id: "ZZZZZZ", text: "opened internal/a.go:10"})
	if m.diffOpened != "opening internal/a.go:10…" || m.diffOpenErr {
		t.Errorf("diffOpened = %q err=%v after another ticket's reply, want it unchanged", m.diffOpened, m.diffOpenErr)
	}
}

// A command that cannot be started at all, for a reason other than a missing
// program, is reported in the title as well.
func TestOKeyShowsALaunchErrorInTheTitle(t *testing.T) {
	orig := openInCode
	openInCode = func(string, []string) (func() error, error) {
		return nil, errors.New("permission denied")
	}
	t.Cleanup(func() { openInCode = orig })

	m := diffModel(t, 160, "signoff")
	m.diffDir = "/repo/worktree"
	m.render()
	m.Update(tea.KeyPressMsg{Code: 'o', Text: "o"})
	if m.mode != modeDetail {
		t.Fatalf("mode = %v after a launch error, want modeDetail, not a modal", m.mode)
	}
	if out := stripANSI(m.render()); !strings.Contains(out, "permission denied") {
		t.Errorf("diff title misses the launch error:\n%s", out)
	}
}

// The real launcher returns at once, and its wait reports the exit status and the
// last stderr line of a command that failed; a clean exit is no error.
func TestOpenInCodeWaitsAndReportsTheExitStatusAndStderr(t *testing.T) {
	wait, err := openInCode(t.TempDir(), []string{"sh", "-c", "echo first >&2; echo boom >&2; exit 3"})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	err = wait()
	if err == nil || !strings.Contains(err.Error(), "exit status 3") || !strings.HasSuffix(err.Error(), ": boom") {
		t.Errorf("wait = %v, want the exit status and the last stderr line boom", err)
	}
	wait, err = openInCode(t.TempDir(), []string{"sh", "-c", "exit 0"})
	if err != nil {
		t.Fatalf("start of a clean command: %v", err)
	}
	if err := wait(); err != nil {
		t.Errorf("wait after a clean exit = %v, want nil", err)
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

// A file's header row takes the start of its first hunk, even when that start
// is line 1; a later hunk of the same file must not overwrite it, and a new
// file starts over.
func TestHeaderLineIsTheFirstHunkOfEachFile(t *testing.T) {
	patch := "diff --git a/x.go b/x.go\n--- a/x.go\n+++ b/x.go\n" +
		"@@ -1,2 +1,2 @@\n a\n-b\n+c\n" +
		"@@ -50,2 +51,2 @@\n d\n-e\n+f\n" +
		"diff --git a/y.go b/y.go\n--- a/y.go\n+++ b/y.go\n" +
		"@@ -7 +7 @@\n-g\n+h\n"
	var headers []diffRow
	for _, r := range parseDiff(patch) {
		if r.kind == 'H' {
			headers = append(headers, r)
		}
	}
	if len(headers) != 2 {
		t.Fatalf("got %d headers, want 2", len(headers))
	}
	if headers[0].line != 1 {
		t.Errorf("first file header line = %d, want 1 (its first hunk starts at +1, not its second at +51)", headers[0].line)
	}
	if headers[1].line != 7 {
		t.Errorf("second file header line = %d, want 7", headers[1].line)
	}
}

// A commit git cannot show leaves no rows; the pane must say the diff is
// unavailable, not that there are no changes.
func TestUnshowableCommitIsDiffUnavailableNotNoChanges(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init", "-q")
	runGit(t, root, "commit", "-q", "--allow-empty", "-m", "base")
	runGit(t, root, "branch", "-M", "master")
	tk := waitingTicket("signoff")
	tk.Commits = []string{"0123456789abcdef0123456789abcdef01234567"}

	_, _, _, notice := ticketDiff(root, tk)
	if notice.head != "Diff unavailable" || !notice.bad {
		t.Errorf("notice = %+v, want the red Diff unavailable state", notice)
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

// bgOpen is the escape a background colour opens with, as lipgloss writes it.
func bgOpen(color string) string {
	s := lipgloss.NewStyle().Background(lipgloss.Color(color)).Render("x")
	return s[:strings.Index(s, "x")]
}

// Changed rows carry their tint across the whole pane, and the changed words
// carry the stronger one; a context row carries no tint at all. The colours are
// 256-palette indices, so a 256-colour terminal shows the same hue.
func TestChangedRowsAreTintedAcrossThePane(t *testing.T) {
	m := diffModel(t, 160, "signoff")
	m.darkBG = true
	m.diffRows = parseDiff(pairPatch)
	m.diffOpen = 0
	width := m.width - diffLeftWidth - 5
	add, del, wordAdd, wordDel := bgOpen("22"), bgOpen("52"), bgOpen("28"), bgOpen("88")
	if add != "\x1b[48;5;22m" || wordDel != "\x1b[48;5;88m" {
		t.Fatalf("palette opens with %q and %q, want the 256-colour indices", add, wordDel)
	}
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
			if !strings.Contains(line, wordDel+"2") {
				t.Errorf("removed row's changed word lacks the word-del colour:\n%q", line)
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

// One number column: a removed row shows its old-side number, every other row its
// new-side one, and the column is as wide as the largest number of the open file.
func TestGutterShowsOneNumberOnTheRowsSide(t *testing.T) {
	m := diffModel(t, 160, "signoff")
	m.diffOpen = 0
	lines := strings.Split(stripANSI(m.renderDiffPane(77, 20)), "\n")
	find := func(text string) string {
		for _, l := range lines {
			if strings.Contains(l, text) {
				return l
			}
		}
		t.Fatalf("no row shows %q", text)
		return ""
	}
	if removed := find("y := 2"); !strings.HasPrefix(removed, "11 │ -") {
		t.Errorf("removed row = %q, want its old number 11 then │ and the sign", removed)
	}
	if added := find("y := 3"); !strings.HasPrefix(added, "11 │ +") {
		t.Errorf("added row = %q, want its new number 11 then │ and the sign", added)
	}
	if context := find("return x"); !strings.HasPrefix(context, "13 │ ") {
		t.Errorf("context row = %q, want its new number 13", context)
	}
	if hunk := find("@@"); !strings.HasPrefix(hunk, strings.Repeat(" ", 5)+"@@") {
		t.Errorf("hunk header %q is not indented past the two-digit gutter", hunk)
	}
}

// A file whose largest number has one digit gets a one-column gutter.
func TestGutterIsOneDigitForAShortFile(t *testing.T) {
	m := diffModel(t, 160, "signoff")
	m.diffRows = parseDiff("diff --git a/s.go b/s.go\n@@ -1,2 +1,2 @@\n a\n-b\n+c\n")
	m.diffOpen = 0
	lines := strings.Split(stripANSI(m.renderDiffPane(77, 8)), "\n")
	var added, context string
	for _, l := range lines {
		switch {
		case strings.Contains(l, "+c"):
			added = l
		case strings.Contains(l, " a"):
			context = l
		}
	}
	if !strings.HasPrefix(added, "2 │ +c") || !strings.HasPrefix(context, "1 │  a") {
		t.Errorf("rows %q and %q, want a one-digit gutter: \"2 │ +c\" and \"1 │  a\"", added, context)
	}
}

func TestDiffTitleCountsFilesAndTotals(t *testing.T) {
	m := diffModel(t, 160, "signoff")
	if out := stripANSI(m.render()); !strings.Contains(out, "file 1/2  +4 −1") || !strings.Contains(out, "○ Diff  feat/x → master") {
		t.Errorf("title misses the counter or the unfocused mark:\n%s", out)
	}
	m.diffSel = 1 // docs/b.md
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

// A code row longer than the pane ends in "…" and stays exactly one row wide.
func TestLongCodeRowIsCutToThePane(t *testing.T) {
	m := diffModel(t, 160, "signoff")
	m.darkBG = true
	m.diffRows = parseDiff("diff --git a/x.txt b/x.txt\n@@ -1 +1 @@\n+" + strings.Repeat("x", 200) + "\n")
	m.diffOpen = 0
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

// A cut that falls exactly on a segment boundary still ends in "…": the segment
// that fills the last column gives up one column for it. Without text after it
// the segment stays whole.
func TestFitSegsEndsAnExactBoundaryCutInEllipsis(t *testing.T) {
	sty := lipgloss.NewStyle()
	text := func(segs []diffSeg) string {
		var b strings.Builder
		for _, s := range segs {
			b.WriteString(s.text)
		}
		return b.String()
	}
	cut := fitSegs([]diffSeg{{"abcd", sty}, {"ef", sty}}, 4)
	if got := text(cut); got != "abc…" {
		t.Errorf("cut at a boundary is %q, want \"abc…\"", got)
	}
	whole := fitSegs([]diffSeg{{"abcd", sty}, {"", sty}}, 4)
	if got := text(whole); got != "abcd" {
		t.Errorf("exact fit with no text after it is %q, want \"abcd\"", got)
	}
}

func TestIDECommandFillsTheTemplate(t *testing.T) {
	cases := []struct {
		name, template, file string
		line                 int
		want                 []string
	}{
		{"default template", "code -g {file}:{line}", "internal/a.go", 10, []string{"code", "-g", "internal/a.go:10"}},
		{"line as its own argument", "idea --line {line} {file}", "internal/a.go", 7, []string{"idea", "--line", "7", "internal/a.go"}},
		{"path with spaces stays one argument", "code -g {file}:{line}", "docs/my notes.md", 3, []string{"code", "-g", "docs/my notes.md:3"}},
		{"no placeholder appends the path", "nvim", "a.go", 1, []string{"nvim", "a.go"}},
	}
	for _, c := range cases {
		if got := ideCommand(c.template, c.file, c.line); !slices.Equal(got, c.want) {
			t.Errorf("%s: ideCommand(%q) = %q, want %q", c.name, c.template, got, c.want)
		}
	}
}

// An IDE that is not installed must not be started: 'o' opens the file in the
// editor instead, and the diff title says which program was missing. Start
// reports a bare program name that is not on PATH as exec.ErrNotFound, and a
// path to a missing program (./ide, /opt/x/idea.sh) as fs.ErrNotExist.
func TestOKeyFallsBackToTheEditorWhenTheIDEIsMissing(t *testing.T) {
	cases := []struct {
		name, prog string
		err        func(prog string) error
	}{
		{"not on PATH", "definitely-not-installed-xyz", func(p string) error { return &exec.Error{Name: p, Err: exec.ErrNotFound} }},
		{"missing path", "/opt/definitely-not-installed-xyz/ide", func(p string) error {
			return &fs.PathError{Op: "fork/exec", Path: p, Err: fs.ErrNotExist}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			orig := openInCode
			openInCode = func(_ string, argv []string) (func() error, error) {
				if argv[0] != c.prog {
					t.Errorf("started %q, want the missing IDE", argv[0])
				}
				return nil, c.err(argv[0])
			}
			t.Cleanup(func() { openInCode = orig })

			// After diffModel: it gives the model its own JAIRA_HOME, which is where the
			// setting has to go.
			m := diffModel(t, 160, "signoff")
			if _, err := settings.Update(func(cur *settings.Settings) { cur.IDE = c.prog + " {file}" }); err != nil {
				t.Fatal(err)
			}
			m.diffDir = "/repo/worktree"
			m.render()
			_, cmd := m.Update(tea.KeyPressMsg{Code: 'o', Text: "o"})
			if cmd == nil {
				t.Fatal("no editor command was returned for the fallback")
			}
			if out := stripANSI(m.render()); !strings.Contains(out, c.prog+" not found — opened in ") {
				t.Errorf("diff title does not say the IDE is missing:\n%s", out)
			}
		})
	}
}

// The open file scrolls with the arrows: down moves its rows under the pinned
// header and stops at its last row, and the files after it show below it.
func TestDownScrollsTheOpenFileAndStopsAtItsLastRow(t *testing.T) {
	m := diffModel(t, 160, "signoff")
	m.render()
	m.Update(key("tab"))
	m.Update(key("enter"))
	out := stripANSI(m.renderDiffPane(77, 4))
	if !strings.Contains(out, "internal/a.go") || !strings.Contains(out, "x := 1") {
		t.Fatalf("open file misses its header or its first row:\n%s", out)
	}
	// The file's rows start at its hunk header, so two downs put the removed
	// line on top and the context line above it is gone.
	m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	out = stripANSI(m.renderDiffPane(77, 4))
	if strings.Contains(out, "x := 1") || !strings.Contains(out, "y := 3") {
		t.Errorf("two downs did not scroll the rows of the file:\n%s", out)
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if out = stripANSI(m.renderDiffPane(77, 4)); !strings.Contains(out, "z := 4") {
		t.Errorf("three downs do not show the added line z := 4:\n%s", out)
	}
	for range 20 {
		m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	}
	m.render()
	if last := len(m.diffFileRows(m.diffHeaders()[0])) - 1; m.diffScroll != last {
		t.Errorf("scroll = %d after many downs, want the last row %d", m.diffScroll, last)
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	if m.diffScroll != len(m.diffFileRows(m.diffHeaders()[0]))-2 {
		t.Errorf("up from the end moved the scroll to %d", m.diffScroll)
	}
}

// Files after the open one show their headers below its rows.
func TestFilesAfterTheOpenOneShowBelowIt(t *testing.T) {
	m := diffModel(t, 160, "signoff")
	m.diffOpen = 0
	if out := stripANSI(m.renderDiffPane(77, 12)); !strings.Contains(out, "docs/b.md") {
		t.Errorf("the file after the open one is not listed:\n%s", out)
	}
}

// Right and left move the open file's code sideways, the sign and the gutter
// staying put, and stop at the longest row of the file.
func TestLeftAndRightScrollTheCodeSideways(t *testing.T) {
	m := diffModel(t, 160, "signoff")
	m.diffRows = parseDiff("diff --git a/x.txt b/x.txt\n@@ -1 +1 @@\n+abcdefghijklmnopqrstuvwxyz\n")
	m.diffOpen = 0
	m.render()
	m.Update(key("tab"))
	m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	if m.diffHScroll != 8 {
		t.Fatalf("right: hscroll = %d, want 8", m.diffHScroll)
	}
	out := stripANSI(m.renderDiffPane(40, 6))
	if !strings.Contains(out, "+ijklmnop") || strings.Contains(out, "abc") {
		t.Errorf("right did not cut the code from column 8:\n%s", out)
	}
	for range 4 {
		m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	}
	m.render()
	if m.diffHScroll != 26 {
		t.Errorf("hscroll = %d after running off the row, want 26 (its 27 runes less the sign)", m.diffHScroll)
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
	if m.diffHScroll != 18 {
		t.Errorf("left: hscroll = %d, want 18", m.diffHScroll)
	}
}

// 'o' opens the row at the top of the open file, not the file's header.
func TestOKeyOpensTheTopRowOfTheOpenFile(t *testing.T) {
	var gotArgv []string
	orig := openInCode
	openInCode = func(_ string, argv []string) (func() error, error) {
		gotArgv = argv
		return func() error { return nil }, nil
	}
	t.Cleanup(func() { openInCode = orig })

	m := diffModel(t, 160, "signoff")
	m.diffDir = "/repo/worktree"
	m.render()
	m.Update(key("tab"))
	m.Update(key("enter"))
	m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m.render()
	m.Update(key("o"))
	if got := gotArgv[len(gotArgv)-1]; got != "internal/a.go:11" {
		t.Errorf("o opened %q, want the top row internal/a.go:11", got)
	}
}

// The divider between the ticket and the diff drags with the mouse, and the
// width it is dragged to stays for the session. A click a column off the divider
// still grabs it; a motion with no drag in progress moves nothing.
func TestDragTheDividerResizesTheSplit(t *testing.T) {
	m := diffModel(t, 160, "signoff")
	m.render()
	if m.View().MouseMode != tea.MouseModeCellMotion {
		t.Error("the diff pane does not ask for mouse events, so the divider cannot be dragged")
	}
	// The divider is at column 80: the ticket's box is 78+2 wide, then a gap.
	m.Update(tea.MouseClickMsg{X: 81, Y: 5, Button: tea.MouseLeft})
	m.Update(tea.MouseMotionMsg{X: 100, Y: 5, Button: tea.MouseLeft})
	if m.diffLeft() != 98 {
		t.Fatalf("dragged to column 100: ticket column = %d, want 98", m.diffLeft())
	}
	m.Update(tea.MouseReleaseMsg{X: 100, Y: 5, Button: tea.MouseLeft})
	m.Update(tea.MouseMotionMsg{X: 120, Y: 5, Button: tea.MouseLeft})
	if m.diffLeft() != 98 {
		t.Errorf("motion after the release moved the divider to %d", m.diffLeft())
	}

	// The drag keeps the ticket at least 40 columns and the diff at least 50.
	m.Update(tea.MouseClickMsg{X: 100, Y: 5, Button: tea.MouseLeft})
	m.Update(tea.MouseMotionMsg{X: 10, Y: 5, Button: tea.MouseLeft})
	if m.diffLeft() != 40 {
		t.Errorf("dragged left past the minimum: ticket column = %d, want 40", m.diffLeft())
	}
	m.Update(tea.MouseMotionMsg{X: 200, Y: 5, Button: tea.MouseLeft})
	if m.diffLeft() != 110 {
		t.Errorf("dragged right past the maximum: ticket column = %d, want 110", m.diffLeft())
	}
	m.Update(tea.MouseReleaseMsg{X: 200, Y: 5, Button: tea.MouseLeft})

	// On screen the split follows the drag: the ticket's box is diffLeft()+2 wide,
	// so its top-right corner sits at diffLeft()+1, and the diff's box starts two
	// columns after that.
	var top []rune
	for _, l := range strings.Split(stripANSI(m.render()), "\n") {
		if strings.Contains(l, "╭") {
			top = []rune(l)
			break
		}
	}
	w := m.diffLeft()
	if len(top) < w+4 || top[w+1] != '╮' || top[w+3] != '╭' {
		t.Errorf("box corners are not at the dragged split (ticket column %d): %q", w, string(top))
	}
}

// Without the diff pane on screen the mouse is left to the terminal, so text can
// still be selected there.
func TestMouseIsLeftToTheTerminalOffTheDiffPane(t *testing.T) {
	m := newTestModel(t, 160, 40)
	m.mode = modeBoard
	if m.View().MouseMode != tea.MouseModeNone {
		t.Error("the board asks for mouse events, which takes text selection away from the terminal")
	}
}

// Git writes a tab after a name with a space, and C-quotes a name with a quote,
// a backslash or a control byte, as the whole "a/…" or "b/…" token. The row
// path must be the name itself, or 'o' opens a file that does not exist.
func TestParseDiffReadsPathsGitWritesInOtherForms(t *testing.T) {
	patch := "diff --git a/my file.txt b/my file.txt\n" +
		"new file mode 100644\n" +
		"--- /dev/null\n" +
		"+++ b/my file.txt\t\n" +
		"@@ -0,0 +1 @@\n" +
		"+hello\n" +
		"diff --git \"a/q\\\"uote.txt\" \"b/q\\\"uote.txt\"\n" +
		"--- \"a/q\\\"uote.txt\"\n" +
		"+++ \"b/q\\\"uote.txt\"\n" +
		"@@ -1 +1 @@\n" +
		"-x\n" +
		"+y\n"
	var headers []string
	for _, r := range parseDiff(patch) {
		if r.kind == 'H' {
			headers = append(headers, r.path)
		}
	}
	want := []string{"my file.txt", `q"uote.txt`}
	if !slices.Equal(headers, want) {
		t.Errorf("header paths = %q, want %q", headers, want)
	}
}

// A file whose name starts with - or + is given to the IDE and to the editor as
// ./name, so it is not read as an option or, in vim, as a command.
func TestDiffArgKeepsAnOptionLikeNameOutOfTheFlags(t *testing.T) {
	cases := map[string]string{
		"-rf.txt":        "./-rf.txt",
		"+!touch pwned":  "./+!touch pwned",
		"internal/a.go":  "internal/a.go",
		"docs/-notes.md": "docs/-notes.md",
	}
	for path, want := range cases {
		if got := diffArg(path); got != want {
			t.Errorf("diffArg(%q) = %q, want %q", path, got, want)
		}
	}
	if got := ideCommand("code -g {file}:{line}", diffArg("-rf.txt"), 4); !slices.Equal(got, []string{"code", "-g", "./-rf.txt:4"}) {
		t.Errorf("IDE argv for -rf.txt = %q", got)
	}
}

// Opening a file whose name starts with - hands the IDE ./name, not the bare name.
func TestOKeyGivesTheIDEADotSlashForAnOptionLikeName(t *testing.T) {
	var gotArgv []string
	orig := openInCode
	openInCode = func(_ string, argv []string) (func() error, error) {
		gotArgv = argv
		return func() error { return nil }, nil
	}
	t.Cleanup(func() { openInCode = orig })

	m := diffModel(t, 160, "signoff")
	m.diffDir = "/repo/worktree"
	m.diffRows = parseDiff("diff --git a/-rf.txt b/-rf.txt\n--- a/-rf.txt\n+++ b/-rf.txt\n@@ -1 +1 @@\n-x\n+y\n")
	m.diffSel, m.diffOpen = 0, -1
	m.render()
	m.Update(tea.KeyPressMsg{Code: 'o', Text: "o"})
	if want := []string{"code", "-g", "./-rf.txt:1"}; !slices.Equal(gotArgv, want) {
		t.Errorf("openInCode argv = %q, want %q", gotArgv, want)
	}
}

// The blank row between two files is never a stop for up and down: the file
// selection walks the file headers only, and 'o' on a header opens that file.
func TestDownNeverStopsOnTheBlankRowBetweenFiles(t *testing.T) {
	m := diffModel(t, 160, "signoff")
	m.render()
	blanks := 0
	for _, r := range m.diffRows {
		if r.kind == 'B' {
			blanks++
		}
	}
	if blanks == 0 {
		t.Fatal("the sample diff has no blank separator row; the test would prove nothing")
	}
	m.Update(key("tab"))
	for range 3 {
		m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
		if kind := m.diffRows[m.diffTargetRow()].kind; kind != 'H' {
			t.Fatalf("down selected a row of kind %q (sel=%d), want a file header", kind, m.diffSel)
		}
	}
	for range 3 {
		m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
		if kind := m.diffRows[m.diffTargetRow()].kind; kind != 'H' {
			t.Fatalf("up selected a row of kind %q (sel=%d), want a file header", kind, m.diffSel)
		}
	}
}
