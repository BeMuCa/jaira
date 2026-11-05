package tui

import (
	"fmt"
	"image/color"
	"os/exec"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/BeMuCa/jaira/core/gitrepo"
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

// openInCode opens a file at a line in VS Code. A package variable so tests can
// replace it with a recorder instead of launching an editor.
var openInCode = func(dir, path string, line int) error {
	cmd := exec.Command("code", "-g", fmt.Sprintf("%s:%d", path, line))
	cmd.Dir = dir
	if err := cmd.Start(); err != nil {
		return err
	}
	go cmd.Wait()
	return nil
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
		m.diffOpened = ""
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

// diffTop is the row at the top of the diff pane, kept inside the rows.
func (m *Model) diffTop() int {
	return min(max(m.diffScroll, 0), max(len(m.diffRows)-1, 0))
}

// diffTarget is the row the file counter, the marker and 'o' act on: the top
// row, or the row below it when the top row is a blank separator, which has no
// path and is never opened.
func (m *Model) diffTarget() int {
	t := m.diffTop()
	if t+1 < len(m.diffRows) && m.diffRows[t].kind == 'B' {
		t++
	}
	return t
}

// jumpFile moves the diff to the next file's header, or to the header of the
// file the top row is in when going back, or the one before it if the top row
// is that header.
func (m *Model) jumpFile(next bool) {
	if len(m.diffRows) == 0 {
		return
	}
	top := m.diffTop()
	if next {
		for i := top + 1; i < len(m.diffRows); i++ {
			if m.diffRows[i].kind == 'H' {
				m.diffScroll = i
				return
			}
		}
		return
	}
	hdr := top
	for hdr > 0 && m.diffRows[hdr].kind != 'H' {
		hdr--
	}
	if hdr < top {
		m.diffScroll = hdr
		return
	}
	for i := hdr - 1; i >= 0; i-- {
		if m.diffRows[i].kind == 'H' {
			m.diffScroll = i
			return
		}
	}
}

// openDiffRow opens the file of the target row in VS Code at its line.
func (m *Model) openDiffRow() {
	if len(m.diffRows) == 0 {
		return
	}
	r := m.diffRows[m.diffTarget()]
	if err := openInCode(m.diffDir, r.path, r.line); err != nil {
		m.notify(err.Error(), true)
		return
	}
	m.diffOpened = "opened " + r.path + ":" + strconv.Itoa(r.line)
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
		bg := lipgloss.NewStyle().Background(fill)
		mark := colFaint
		if r.file == sel {
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
			{strings.Repeat(" ", 2*nw+4), lipgloss.NewStyle()},
			{head, lipgloss.NewStyle().Foreground(colAgentic)},
			{heading, styMeta},
		}
		return padRow(segs, width, lipgloss.NewStyle())
	}
	// A '-' row has no new-side number and a '+' row no old-side one, so each
	// number is shown only on the side the row belongs to.
	num := func(n int) string {
		if n == 0 {
			return strings.Repeat(" ", nw)
		}
		return fmt.Sprintf("%*d", nw, n)
	}
	newNum := num(r.line)
	if r.kind == '-' {
		newNum = num(0)
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
		{num(r.old) + " " + newNum + " │ ", bg.Foreground(colFaint)},
		{r.text[:1], sign},
	}
	rs := []rune(r.text)
	if r.to > r.from {
		segs = append(segs,
			diffSeg{string(rs[1:r.from]), bg},
			diffSeg{string(rs[r.from:r.to]), wordBg},
			diffSeg{string(rs[r.to:]), bg})
	} else {
		segs = append(segs, diffSeg{string(rs[1:]), bg})
	}
	return padRow(segs, width, bg)
}

// renderDiffPane renders the diff into exactly width x height: a title row, then
// the rows from diffScroll, or the notice when there are no rows. The file of
// the top row is marked in its header. Every row is built at width exactly.
func (m *Model) renderDiffPane(width, height int) string {
	// Line numbers are as wide as the largest in the diff, so the gutter does not
	// change width as the diff is scrolled.
	nw, files, adds, dels := 3, 0, 0, 0
	for _, r := range m.diffRows {
		nw = max(nw, len(strconv.Itoa(max(r.line, r.old))))
		if r.kind == 'H' {
			files++
			adds += r.adds
			dels += r.dels
		}
	}
	// The scroll may go to the last row, not the last full page: "]" must be
	// able to put a short file's header at the top even when the whole diff
	// would fit.
	top := m.diffTop()
	m.diffScroll = top
	selFile := -1
	if len(m.diffRows) > 0 {
		selFile = m.diffRows[m.diffTarget()].file
	}

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
		right = []diffSeg{{m.diffOpened, styOK}, {" ", lipgloss.NewStyle()}}
	case files > 0:
		right = []diffSeg{
			{fmt.Sprintf("file %d/%d", selFile+1, files), styMeta},
			{fmt.Sprintf("  +%d", adds), styOK},
			{fmt.Sprintf(" −%d", dels), styErr},
			{" ", lipgloss.NewStyle()},
		}
	}
	lines := []string{spreadRow(left, right, width, lipgloss.NewStyle())}
	if len(m.diffRows) == 0 {
		lines = append(lines, noticeRows(m.diffNotice, width)...)
	}
	body := max(1, height-1)
	end := min(len(m.diffRows), top+body)
	for _, r := range m.diffRows[top:end] {
		lines = append(lines, m.diffLine(r, width, nw, selFile))
	}
	pane, _ := clipPane(strings.Join(lines, "\n"), width, height, 0)
	return pane
}

// withDiff puts the ticket body, built at diffLeftWidth, on the left and the
// diff on the right, with the footer below both. hint is the ticket's own key
// list; the diff's keys are appended to it.
func (m *Model) withDiff(left, hint string) string {
	footer := wrapHints(strings.Split(hint+" · tab diff · ]/[ file · o open in VS Code", " · "), max(1, m.width))
	// A pane costs a row top and bottom for its box, and a column each side and
	// the space between them.
	rows := max(1, m.height-len(footer)-2)
	leftPane, scroll := clipPane(left, diffLeftWidth, rows, m.detailScroll)
	m.detailScroll = scroll
	right := m.renderDiffPane(m.width-diffLeftWidth-5, rows)
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
