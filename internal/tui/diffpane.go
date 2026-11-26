package tui

import (
	"bytes"
	"errors"
	"fmt"
	"image/color"
	"io/fs"
	"os/exec"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/BeMuCa/jaira/core/gitrepo"
	"github.com/BeMuCa/jaira/core/settings"
	"github.com/BeMuCa/jaira/core/ticket"
)

// diffSplitMinWidth is the terminal width from which a ticket a person owns
// shows its diff beside it. Below it the screen is the single column it has
// always been, so the ticket keeps its full width on a narrow terminal.
const diffSplitMinWidth = 140

// diffLeftWidth is the width of the ticket column when the diff sits beside it.
// It is the width the sign-off screen was laid out at, so the ticket reads as
// it did before the diff was added, and the diff gets what is left.
const diffLeftWidth = 78

// diffPanStep is how many columns left or right moves the open file's code.
const diffPanStep = 8

// diffRow is one rendered row of the diff pane. Every row is one line on the
// screen, so the scroll position is an index into the rows.
type diffRow struct {
	file int    // index into the file list
	path string // repo-relative new path (old path for a deletion)
	line int    // new-side line this row points at; at least 1
	old  int    // old-side line of a '-' or context row; 0 = none
	kind byte   // 'H' file header, 'B' blank row above a header, '@' hunk header, '+', '-', ' ' context
	text string // raw text, sign included, tabs expanded to 4 spaces
	// Set on 'H' rows only: how the file changed when it is not a plain edit
	// (" new", " deleted", " renamed"), and its changed line counts.
	tag        string
	adds, dels int
	// Set on a paired '+' or '-' row: the rune span of text that differs from
	// its partner. from == to means no span.
	from, to int
}

// openInCode starts the IDE command argv in dir. It returns at once, with the
// error Start gives, and a wait function that blocks until the command exits.
// It has no deadline: the Windows 'code' wrapper under WSL took 46 seconds to
// report a broken VS Code server on a loaded machine, and the wait runs in a
// tea.Cmd, so the screen never waits on it. A non-zero exit is an error that
// carries the exit status and the last line the command wrote to stderr. A
// package variable so tests can replace it with a stub instead of an IDE.
var openInCode = func(dir string, argv []string) (func() error, error) {
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return func() error {
		err := cmd.Wait()
		if err == nil {
			return nil
		}
		last := strings.TrimSpace(stderr.String())
		if i := strings.LastIndex(last, "\n"); i >= 0 {
			last = last[i+1:]
		}
		if last == "" {
			return err
		}
		return fmt.Errorf("%w: %s", err, last)
	}, nil
}

// ideCommand is the command that opens file at line with the IDE template. The
// template is split into fields before {file} and {line} are filled in, so a
// path with spaces stays one argument. A template with no {file} gets the path
// as its last argument.
func ideCommand(template, file string, line int) []string {
	var argv []string
	hasFile := false
	for _, f := range strings.Fields(template) {
		hasFile = hasFile || strings.Contains(f, "{file}")
		f = strings.ReplaceAll(f, "{file}", file)
		argv = append(argv, strings.ReplaceAll(f, "{line}", strconv.Itoa(line)))
	}
	if !hasFile {
		argv = append(argv, file)
	}
	return argv
}

// diffEditorDoneMsg reports how the editor opened for a file the IDE does not
// have ended.
type diffEditorDoneMsg struct{ err error }

// diffOpenedMsg reports how the IDE returned from opening text: err is nil when
// it exited cleanly.
type diffOpenedMsg struct {
	text string
	err  error
}

// parseDiff turns a git patch into one row per line of the diff, numbering each
// row by its new-side line so a row can be opened in an editor. Lines outside a
// file's hunks are dropped: the headers git writes between files, and the
// commit lines and stat block of 'git show', which the fallback diff carries.
func parseDiff(patch string) []diffRow {
	var rows []diffRow
	var (
		file     = -1 // index of the file being read; -1 before the first
		header   = -1 // index into rows of that file's header row
		oldPath  string
		path     string
		oldLine  int
		newLine  int
		inHunk   bool
		seenHunk bool   // the file's first hunk has set its header line
		skipping = true // true until a "diff --git" line starts a file
	)
	for _, raw := range strings.Split(patch, "\n") {
		if strings.HasPrefix(raw, "diff --git ") {
			oldPath, path = splitDiffGit(raw)
			if file >= 0 {
				// The blank row sits above this header, so it belongs to the file before.
				rows = append(rows, diffRow{file: file, kind: 'B'})
			}
			file++
			header = len(rows)
			rows = append(rows, diffRow{file: file, path: path, line: 1, kind: 'H'})
			inHunk, skipping, seenHunk = false, false, false
			continue
		}
		if strings.HasPrefix(raw, "commit ") {
			inHunk, skipping = false, true
			continue
		}
		if skipping || raw == "" {
			continue
		}
		if strings.HasPrefix(raw, "@@") {
			inHunk = true
			oldLine = hunkStart(raw, "-")
			newLine = hunkStart(raw, "+")
			if !seenHunk {
				rows[header].line = max(newLine, 1)
				seenHunk = true
			}
			rows = append(rows, diffRow{file: file, path: path, line: max(newLine, 1), kind: '@', text: raw})
			continue
		}
		if !inHunk {
			switch {
			case strings.HasPrefix(raw, "+++ "):
				// A deletion names /dev/null here, and the file is still the old path.
				path = strings.TrimPrefix(raw[4:], "b/")
				if raw[4:] == "/dev/null" {
					path = oldPath
				}
				rows[header].path = path
			case strings.HasPrefix(raw, "new file mode"):
				rows[header].tag = " new"
			case strings.HasPrefix(raw, "deleted file mode"):
				rows[header].tag = " deleted"
			case strings.HasPrefix(raw, "rename from"):
				rows[header].tag = " renamed"
			}
			continue
		}
		text := strings.ReplaceAll(raw, "\t", "    ")
		switch raw[0] {
		case '+':
			rows = append(rows, diffRow{file: file, path: path, line: newLine, kind: '+', text: text})
			newLine++
		case '-':
			rows = append(rows, diffRow{file: file, path: path, line: max(newLine, 1), old: oldLine, kind: '-', text: text})
			oldLine++
		case ' ':
			rows = append(rows, diffRow{file: file, path: path, line: max(newLine, 1), old: oldLine, kind: ' ', text: text})
			oldLine++
			newLine++
		}
	}
	markWordPairs(rows)
	adds, dels := map[int]int{}, map[int]int{}
	for _, r := range rows {
		switch r.kind {
		case '+':
			adds[r.file]++
		case '-':
			dels[r.file]++
		}
	}
	for i := range rows {
		if rows[i].kind == 'H' {
			rows[i].adds, rows[i].dels = adds[rows[i].file], dels[rows[i].file]
		}
	}
	return rows
}

// markWordPairs sets the word span of each '-' row that a '+' row replaces: a
// run of k '-' rows directly followed by a run of exactly k '+' rows pairs the
// i-th '-' with the i-th '+'. Unequal runs are not paired, since then nothing
// says which removed line became which added one.
func markWordPairs(rows []diffRow) {
	for i := 0; i < len(rows); {
		if rows[i].kind != '-' {
			i++
			continue
		}
		j := i
		for j < len(rows) && rows[j].kind == '-' {
			j++
		}
		k := j
		for k < len(rows) && rows[k].kind == '+' {
			k++
		}
		if k-j == j-i {
			for n := 0; n < j-i; n++ {
				markWordPair(&rows[i+n], &rows[j+n])
			}
		}
		i = k
	}
}

// markWordPair spans the text between the common prefix and the common suffix of
// a '-' and a '+' row. Fewer than three shared runes means the lines differ
// throughout, and a highlight would mark nothing useful. Runes are counted from
// the sign, which is rune 0 of text, so the span starts at 1 + the prefix.
func markWordPair(del, add *diffRow) {
	a, b := []rune(del.text[1:]), []rune(add.text[1:])
	p := 0
	for p < len(a) && p < len(b) && a[p] == b[p] {
		p++
	}
	s := 0
	for s < len(a)-p && s < len(b)-p && a[len(a)-1-s] == b[len(b)-1-s] {
		s++
	}
	if p+s < 3 {
		return
	}
	del.from, del.to = 1+p, 1+len(a)-s
	add.from, add.to = 1+p, 1+len(b)-s
}

// splitDiffGit reads "diff --git a/X b/Y" into X and Y.
func splitDiffGit(line string) (oldPath, newPath string) {
	rest := strings.TrimPrefix(line, "diff --git a/")
	i := strings.Index(rest, " b/")
	if i < 0 {
		return rest, rest
	}
	return rest[:i], rest[i+3:]
}

// hunkStart reads the first line of one side of "@@ -a,b +c,d @@": side "-" gives
// a, side "+" gives c.
func hunkStart(h, side string) int {
	f := strings.Fields(h)
	if len(f) < 3 || !strings.HasPrefix(f[1], "-") || !strings.HasPrefix(f[2], "+") {
		return 0
	}
	field := f[1]
	if side == "+" {
		field = f[2]
	}
	n, _ := strconv.Atoi(strings.SplitN(field[1:], ",", 2)[0])
	return n
}

// paneNotice is what the diff pane says in place of its rows when there are
// none: a bold head, then detail lines, or red detail when bad is set.
type paneNotice struct {
	head  string
	lines []string
	bad   bool
}

// ticketDiff assembles the diff of one ticket. The branch named after the
// ticket is compared with the base branch, which is what a pull request shows;
// when no branch carries the handle (merged and deleted, say), the commits that
// name the handle on this checkout stand in. note says where the diff came
// from, notice says why there are no rows, and dir is the checkout 'o' opens
// files in.
func ticketDiff(root string, t *ticket.Ticket) (rows []diffRow, note, dir string, notice paneNotice) {
	repo := &gitrepo.Repo{Dir: root}
	handle := ticket.Handle(t.ID)
	br, base := repo.BranchFor(handle), repo.Base()
	if br != "" && base != "" {
		dir = repo.WorktreeOf(br)
		if dir == "" {
			dir = root
		}
		patch, err := repo.BranchDiff(base, br)
		if err != nil {
			return nil, "", dir, paneNotice{head: "Diff unavailable", lines: []string{err.Error()}, bad: true}
		}
		note = br + " → " + base
		rows = parseDiff(patch)
		if len(rows) == 0 {
			notice = paneNotice{head: "No changes", lines: []string{br + " has nothing that " + base + " lacks — merged already?"}}
		}
		return rows, note, dir, notice
	}
	dir = root
	shas, err := repo.CommitsForTicket(t.Path, t.ID)
	if err != nil {
		return nil, "", dir, paneNotice{head: "Diff unavailable", lines: []string{err.Error()}, bad: true}
	}
	shas = ticket.MergeCommits(shas, t.Commits)
	if len(shas) == 0 {
		return nil, "", dir, paneNotice{head: "No diff yet", lines: []string{
			"no branch named after " + handle + ", and no commit names it",
			"push a branch feat/" + handle + "-… to see it here",
		}}
	}
	patch, unavailable := repo.Diff(shas)
	if len(unavailable) == len(shas) {
		// Repo.Diff has no error return: a commit git cannot show comes back as
		// a placeholder line the parser drops, so with none left it is not "no changes".
		return nil, "", dir, paneNotice{head: "Diff unavailable", lines: []string{
			"git cannot show the commits naming " + handle + " on this checkout: " + strings.Join(unavailable, " "),
		}, bad: true}
	}
	note = "commits naming " + handle + " on this checkout"
	rows = parseDiff(patch)
	if len(rows) == 0 {
		notice = paneNotice{head: "No changes", lines: []string{"the commits naming " + handle + " on this checkout change nothing"}}
	}
	return rows, note, dir, notice
}

// loadDiff computes the diff of the open ticket once, until a reload clears
// diffFor. The scroll and focus are reset only when the ticket changes, so a
// reload that re-reads the same diff leaves the person where they were.
func (m *Model) loadDiff(t *ticket.Ticket) {
	if m.diffFor == t.ID {
		return
	}
	m.diffRows, m.diffNote, m.diffDir, m.diffNotice = ticketDiff(m.store.Root, t)
	m.diffFor = t.ID
	if m.diffScrollFor != t.ID {
		m.diffScroll = 0
		m.diffFocus = false
		m.diffOpened, m.diffOpenErr = "", false
		m.diffSel, m.diffOpen, m.diffHScroll = 0, -1, 0
		m.diffScrollFor = t.ID
	}
}

// diffPaneShown reports whether the open ticket shows its diff beside it: it
// sits in a lane a person owns, no follow-up split is open, and the terminal is
// wide and tall enough for two columns.
func (m *Model) diffPaneShown() bool {
	if m.detail == nil || m.follow != nil || m.width < diffSplitMinWidth || m.height < splitMinHeight {
		return false
	}
	l, ok := m.lanes.Get(m.detail.Status)
	return ok && (l.RequiresHumanExit || l.RequiresQuestion)
}

// diffLeft is the width of the ticket column beside the diff: the width the
// divider was dragged to this session, or diffLeftWidth until it is.
func (m *Model) diffLeft() int {
	if m.diffLeftW == 0 {
		return diffLeftWidth
	}
	return clampDiffLeft(m.diffLeftW, m.width)
}

// clampDiffLeft keeps the ticket column at least 40 columns wide and the diff at
// least 50.
func clampDiffLeft(x, width int) int {
	return min(max(x, 40), width-50)
}

// onDiffDivider reports whether column x is the divider or the column beside it.
// The divider is the gap after the ticket's box, two columns past the width.
func (m *Model) onDiffDivider(x int) bool {
	d := x - (m.diffLeft() + 2)
	return d >= -1 && d <= 1
}

// diffHeaders is the index into diffRows of each file's header row, in file order.
func (m *Model) diffHeaders() []int {
	var hdrs []int
	for i, r := range m.diffRows {
		if r.kind == 'H' {
			hdrs = append(hdrs, i)
		}
	}
	return hdrs
}

// diffFileRows is the index of each row under the header at hdr, up to the next
// file's header. Blank separators are left out: the pane never shows them.
func (m *Model) diffFileRows(hdr int) []int {
	var rows []int
	for i := hdr + 1; i < len(m.diffRows) && m.diffRows[i].kind != 'H'; i++ {
		if m.diffRows[i].kind != 'B' {
			rows = append(rows, i)
		}
	}
	return rows
}

// diffShown settles the selection and scroll for a body of body rows, and returns
// the rows the pane shows as indices into diffRows. With no file open that is one
// header per file from diffScroll, moved so the selected file stays in view. With
// a file open it is that file's header, its rows from diffScroll, then the headers
// of the files after it.
func (m *Model) diffShown(body int) []int {
	hdrs := m.diffHeaders()
	m.diffSel = min(max(m.diffSel, 0), max(len(hdrs)-1, 0))
	if m.diffOpen >= len(hdrs) {
		m.diffOpen = -1
	}
	if m.diffOpen < 0 {
		top := min(max(m.diffScroll, 0), max(len(hdrs)-body, 0))
		if m.diffSel < top {
			top = m.diffSel
		}
		if m.diffSel >= top+body {
			top = m.diffSel - body + 1
		}
		m.diffScroll = top
		return hdrs[top:]
	}
	hdr := hdrs[m.diffOpen]
	rows := m.diffFileRows(hdr)
	m.diffScroll = min(max(m.diffScroll, 0), max(len(rows)-1, 0))
	m.diffHScroll = min(max(m.diffHScroll, 0), m.diffFileWidth(rows))
	shown := append([]int{hdr}, rows[m.diffScroll:]...)
	return append(shown, hdrs[m.diffOpen+1:]...)
}

// diffTargetRow is the row 'o' opens: the header of the selected file, or of the
// open file. A non-empty diff always has a header, since every row follows one.
func (m *Model) diffTargetRow() int {
	hdrs := m.diffHeaders()
	if m.diffOpen >= 0 && m.diffOpen < len(hdrs) {
		rows := m.diffFileRows(hdrs[m.diffOpen])
		if len(rows) == 0 {
			return hdrs[m.diffOpen]
		}
		return rows[min(m.diffScroll, len(rows)-1)]
	}
	return hdrs[min(max(m.diffSel, 0), len(hdrs)-1)]
}

// diffStep moves the pane by d: the file selection while no file is open, the
// open file's rows otherwise. The rows' upper bound is applied when they are
// shown, since that is where the file's rows are known.
func (m *Model) diffStep(d int) {
	if m.diffOpen < 0 {
		m.diffSel = min(max(m.diffSel+d, 0), max(len(m.diffHeaders())-1, 0))
		return
	}
	m.diffScroll = max(m.diffScroll+d, 0)
}

// diffPan moves the open file's code sideways by d columns. diffShown keeps the
// move inside the file's longest row.
func (m *Model) diffPan(d int) {
	if m.diffOpen >= 0 {
		m.diffHScroll = max(m.diffHScroll+d, 0)
	}
}

// diffGutter is the width of the line numbers: as many digits as the largest
// number in the open file, so the gutter keeps its width as the file scrolls. With
// no file open there are no numbers, and the width is one digit.
func (m *Model) diffGutter() int {
	nw := 1
	hdrs := m.diffHeaders()
	if m.diffOpen < 0 || m.diffOpen >= len(hdrs) {
		return nw
	}
	for _, i := range m.diffFileRows(hdrs[m.diffOpen]) {
		r := m.diffRows[i]
		if r.kind == '@' {
			continue
		}
		n := r.line
		if r.kind == '-' {
			n = r.old
		}
		nw = max(nw, len(strconv.Itoa(n)))
	}
	return nw
}

// diffFileWidth is the length in columns of the file's longest row, sign not
// counted: how far its code can be scrolled sideways.
func (m *Model) diffFileWidth(rows []int) int {
	w := 0
	for _, i := range rows {
		w = max(w, len([]rune(m.diffRows[i].text))-1)
	}
	return w
}

// toggleDiffFile opens the selected file, or closes the open one, which gives the
// list its keys back.
func (m *Model) toggleDiffFile() {
	if m.diffOpen >= 0 {
		m.diffOpen, m.diffScroll = -1, 0
		return
	}
	if len(m.diffHeaders()) > 0 {
		m.diffOpen, m.diffScroll = m.diffSel, 0
	}
}

// openDiffRow opens the target row's file in the IDE from settings at its line.
// The IDE is waited for in the returned command, which reports back as a
// diffOpenedMsg, so the title says what happened without a modal. When the IDE's
// program is not installed, the file opens in the editor instead, without a line:
// +N means something different in every editor. That editor command is returned
// too, and nil is returned when the IDE was started.
func (m *Model) openDiffRow() tea.Cmd {
	if len(m.diffRows) == 0 {
		return nil
	}
	r := m.diffRows[m.diffTargetRow()]
	cur := settings.Load()
	argv := ideCommand(cur.IDETemplate(), r.path, r.line)
	where := r.path + ":" + strconv.Itoa(r.line)
	wait, err := openInCode(m.diffDir, argv)
	// A bare program name fails Start with exec.ErrNotFound; a path to a missing
	// program (./ide, /opt/x/idea.sh) fails with fs.ErrNotExist.
	if errors.Is(err, exec.ErrNotFound) || errors.Is(err, fs.ErrNotExist) {
		editor, _ := editorChoice(cur)
		cmd := exec.Command(editor[0], append(editor[1:], r.path)...)
		cmd.Dir = m.diffDir
		m.diffOpened, m.diffOpenErr = argv[0]+" not found — opened in "+editor[0], false
		return tea.ExecProcess(cmd, func(err error) tea.Msg { return diffEditorDoneMsg{err} })
	}
	if err != nil {
		m.diffOpened, m.diffOpenErr = err.Error(), true
		return nil
	}
	m.diffOpened, m.diffOpenErr = "opening "+where+"…", false
	return func() tea.Msg {
		if err := wait(); err != nil {
			return diffOpenedMsg{text: "opened " + where, err: err}
		}
		return diffOpenedMsg{text: "opened " + where}
	}
}

// diffPalette is the background of the changed rows and of the words that
// changed inside them, for a dark and a light terminal.
type diffPalette struct {
	add, del, addWord, delWord color.Color
}

var (
	// 256-colour indices, as view.go's colours are: the hex tints the first
	// version used fell onto the grey ramp of a 256-colour terminal.
	diffPaletteDark  = diffPalette{lipgloss.Color("22"), lipgloss.Color("52"), lipgloss.Color("28"), lipgloss.Color("88")}
	diffPaletteLight = diffPalette{lipgloss.Color("194"), lipgloss.Color("224"), lipgloss.Color("157"), lipgloss.Color("217")}
)

// diffSeg is a run of plain text and the style it is drawn in. A row is built
// from segments so it can be cut while it is still plain text: cutting a styled
// string would break its escape codes.
type diffSeg struct {
	text string
	sty  lipgloss.Style
}

// fitSegs cuts segments to n columns, a cut ending in "…" as truncate does. A
// segment that fills the last column exactly is cut too when text follows it,
// so the row still ends in "…".
func fitSegs(segs []diffSeg, n int) []diffSeg {
	var out []diffSeg
	used := 0
	for i, s := range segs {
		w := lipgloss.Width(s.text)
		if used+w < n || used+w == n && segsWidth(segs[i+1:]) == 0 {
			out = append(out, s)
			used += w
			continue
		}
		if room := n - used; room > 0 {
			out = append(out, diffSeg{truncate(s.text+"…", room), s.sty})
		}
		break
	}
	return out
}

func segsWidth(segs []diffSeg) int {
	w := 0
	for _, s := range segs {
		w += lipgloss.Width(s.text)
	}
	return w
}

func paintSegs(segs []diffSeg) string {
	var b strings.Builder
	for _, s := range segs {
		b.WriteString(s.sty.Render(s.text))
	}
	return b.String()
}

// padRow cuts segs to width and pads the rest with blanks in sty, so the row is
// exactly width wide and a background in sty covers all of it.
func padRow(segs []diffSeg, width int, sty lipgloss.Style) string {
	segs = fitSegs(segs, width)
	if pad := width - segsWidth(segs); pad > 0 {
		segs = append(segs, diffSeg{strings.Repeat(" ", pad), sty})
	}
	return paintSegs(segs)
}

// spreadRow puts left at the start of a row and right at its end, exactly width
// wide. The left part is cut first, so a count on the right is never cut.
func spreadRow(left, right []diffSeg, width int, fill lipgloss.Style) string {
	rw := segsWidth(right)
	left = fitSegs(left, max(0, width-rw-1))
	gap := max(0, width-rw-segsWidth(left))
	return paintSegs(append(append(left, diffSeg{strings.Repeat(" ", gap), fill}), right...))
}

// noticeRows lays a notice out below the title: a blank row, the bold head, then
// the detail lines, each centred in width.
func noticeRows(n paneNotice, width int) []string {
	centred := func(text string, sty lipgloss.Style) string {
		text = truncate(text, width)
		return strings.Repeat(" ", max(0, (width-lipgloss.Width(text))/2)) + sty.Render(text)
	}
	rows := []string{""}
	if n.head != "" {
		rows = append(rows, centred(n.head, styLaneTitle))
	}
	detail := styMeta
	if n.bad {
		detail = styErr
	}
	for _, l := range n.lines {
		rows = append(rows, centred(l, detail))
	}
	return rows
}

// diffLine renders one row of the diff pane exactly width wide. nw is the width
// of a line number and sel the file of the top row, whose header is marked.
func (m *Model) diffLine(r diffRow, width, nw, sel int) string {
	fill, _ := m.selectionFill()
	pal := diffPaletteLight
	if m.darkBG {
		pal = diffPaletteDark
	}
	switch r.kind {
	case 'B':
		return ""
	case 'H':
		// Only the selected file's header takes the selection fill: the rest of
		// the list is plain, so the selection is the one bar on the screen.
		bg := lipgloss.NewStyle()
		mark := colFaint
		if r.file == sel {
			bg = lipgloss.NewStyle().Background(fill)
			mark = colAccent
		}
		dir, base := "", r.path
		if i := strings.LastIndex(r.path, "/"); i >= 0 {
			dir, base = r.path[:i+1], r.path[i+1:]
		}
		left := []diffSeg{
			{" ", bg},
			{"▌", bg.Foreground(mark)},
			{" ", bg},
			{dir, bg.Foreground(colDim)},
			{base, bg.Bold(true)},
			{r.tag, bg.Foreground(colWarn)},
		}
		right := []diffSeg{
			{fmt.Sprintf("+%d", r.adds), bg.Foreground(colOK)},
			{fmt.Sprintf(" −%d", r.dels), bg.Foreground(colErr)},
			{" ", bg},
		}
		return spreadRow(left, right, width, bg)
	case '@':
		// The hunk's own "@@ … @@" is coloured apart from the heading git puts after it.
		head, heading := r.text, ""
		if i := strings.Index(r.text[2:], "@@"); i >= 0 {
			head, heading = r.text[:i+4], r.text[i+4:]
		}
		segs := []diffSeg{
			{strings.Repeat(" ", nw+3), lipgloss.NewStyle()},
			{head, lipgloss.NewStyle().Foreground(colAgentic)},
			{heading, styMeta},
		}
		return padRow(segs, width, lipgloss.NewStyle())
	}
	// A '-' row shows its old-side number and every other row its new-side one.
	num := func(n int) string {
		if n == 0 {
			return strings.Repeat(" ", nw)
		}
		return fmt.Sprintf("%*d", nw, n)
	}
	n := r.line
	if r.kind == '-' {
		n = r.old
	}
	bg, wordBg, sign := lipgloss.NewStyle(), lipgloss.NewStyle(), lipgloss.NewStyle()
	switch r.kind {
	case '+':
		bg, wordBg = bg.Background(pal.add), wordBg.Background(pal.addWord)
		sign = bg.Foreground(colOK)
	case '-':
		bg, wordBg = bg.Background(pal.del), wordBg.Background(pal.delWord)
		sign = bg.Foreground(colErr)
	}
	segs := []diffSeg{
		{num(n) + " │ ", bg.Foreground(colFaint)},
		{r.text[:1], sign},
	}
	// The code after the sign scrolls sideways; the sign and the gutter stay, and
	// the word span moves with the text it marks.
	rs := []rune(r.text)
	off := min(m.diffHScroll, max(len(rs)-1, 0))
	rs = append(rs[:1:1], rs[1+off:]...)
	from, to := max(r.from-off, 1), max(r.to-off, 1)
	if r.to > r.from && to > from {
		segs = append(segs,
			diffSeg{string(rs[1:from]), bg},
			diffSeg{string(rs[from:to]), wordBg},
			diffSeg{string(rs[to:]), bg})
	} else {
		segs = append(segs, diffSeg{string(rs[1:]), bg})
	}
	return padRow(segs, width, bg)
}

// renderDiffPane renders the diff into exactly width x height: a title row, then
// the rows from diffScroll, or the notice when there are no rows. The file of
// the top row is marked in its header. Every row is built at width exactly.
func (m *Model) renderDiffPane(width, height int) string {
	files, adds, dels := 0, 0, 0
	for _, r := range m.diffRows {
		if r.kind == 'H' {
			files++
			adds += r.adds
			dels += r.dels
		}
	}
	shown := m.diffShown(max(1, height-1))
	nw := m.diffGutter()

	mark, markSty := "○ ", styMeta
	if m.diffFocus {
		mark, markSty = "● ", styReview
	}
	label := "Diff"
	if m.diffNote != "" {
		label += "  " + m.diffNote
	}
	left := []diffSeg{{mark, markSty}, {label, markSty}}
	var right []diffSeg
	switch {
	case m.diffOpened != "":
		// The reply replaces the counter; the title's left part is what gets cut.
		sty := styOK
		if m.diffOpenErr {
			sty = styErr
		}
		right = []diffSeg{{truncate(m.diffOpened, width-1), sty}, {" ", lipgloss.NewStyle()}}
	case files > 0:
		right = []diffSeg{
			{fmt.Sprintf("file %d/%d", m.diffSel+1, files), styMeta},
			{fmt.Sprintf("  +%d", adds), styOK},
			{fmt.Sprintf(" −%d", dels), styErr},
			{" ", lipgloss.NewStyle()},
		}
	}
	lines := []string{spreadRow(left, right, width, lipgloss.NewStyle())}
	if len(m.diffRows) == 0 {
		lines = append(lines, noticeRows(m.diffNotice, width)...)
	}
	for _, i := range shown[:min(len(shown), max(1, height-1))] {
		lines = append(lines, m.diffLine(m.diffRows[i], width, nw, m.diffSel))
	}
	pane, _ := clipPane(strings.Join(lines, "\n"), width, height, 0)
	return pane
}

// withDiff puts the ticket body, built at diffLeftWidth, on the left and the
// diff on the right, with the footer below both. hint is the ticket's own key
// list; the diff's keys are appended to it.
func (m *Model) withDiff(left, hint string) string {
	footer := wrapHints(strings.Split(hint+" · tab diff · ↑↓ file · enter open/close · ←→ scroll · o open in IDE", " · "), max(1, m.width))
	// A pane costs a row top and bottom for its box, and a column each side and
	// the space between them.
	rows := max(1, m.height-len(footer)-2)
	leftPane, scroll := clipPane(left, m.diffLeft(), rows, m.detailScroll)
	m.detailScroll = scroll
	right := m.renderDiffPane(m.width-m.diffLeft()-5, rows)
	// The live pane is the one the keys reach, marked by its accent border as in
	// the follow-up split.
	var b strings.Builder
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
		splitPane(leftPane, !m.diffFocus), " ", splitPane(right, m.diffFocus)))
	for _, l := range footer {
		b.WriteString("\n" + styMeta.Render(l))
	}
	return clampBlock(b.String(), m.width, m.height)
}
