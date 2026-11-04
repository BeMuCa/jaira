package tui

import (
	"fmt"
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
	kind byte   // 'H' file header, '@' hunk header, '+', '-', ' ' context
	text string // raw text, sign included, tabs expanded to 4 spaces
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
		newLine  int
		inHunk   bool
		skipping = true // true until a "diff --git" line starts a file
	)
	for _, raw := range strings.Split(patch, "\n") {
		if strings.HasPrefix(raw, "diff --git ") {
			oldPath, path = splitDiffGit(raw)
			file++
			header = len(rows)
			rows = append(rows, diffRow{file: file, path: path, line: 1, kind: 'H', text: path})
			inHunk, skipping = false, false
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
			newLine = hunkNewStart(raw)
			if rows[header].line == 1 {
				rows[header].line = max(newLine, 1)
			}
			rows = append(rows, diffRow{file: file, path: path, line: max(newLine, 1), kind: '@', text: raw})
			continue
		}
		if !inHunk {
			if strings.HasPrefix(raw, "+++ ") {
				// A deletion names /dev/null here, and the file is still the old path.
				path = strings.TrimPrefix(raw[4:], "b/")
				if raw[4:] == "/dev/null" {
					path = oldPath
				}
				rows[header].path = path
			}
			continue
		}
		text := strings.ReplaceAll(raw, "\t", "    ")
		switch raw[0] {
		case '+':
			rows = append(rows, diffRow{file: file, path: path, line: newLine, kind: '+', text: text})
			newLine++
		case '-':
			rows = append(rows, diffRow{file: file, path: path, line: max(newLine, 1), kind: '-', text: text})
		case ' ':
			rows = append(rows, diffRow{file: file, path: path, line: max(newLine, 1), kind: ' ', text: text})
			newLine++
		}
	}
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
			rows[i].text = fmt.Sprintf("%s  +%d −%d", rows[i].path, adds[rows[i].file], dels[rows[i].file])
		}
	}
	return rows
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

// hunkNewStart reads the new-side start from "@@ -a,b +c,d @@": c.
func hunkNewStart(h string) int {
	f := strings.Fields(h)
	if len(f) < 3 || !strings.HasPrefix(f[2], "+") {
		return 0
	}
	n, _ := strconv.Atoi(strings.SplitN(strings.TrimPrefix(f[2], "+"), ",", 2)[0])
	return n
}

// ticketDiff assembles the diff of one ticket. The branch named after the
// ticket is compared with the base branch, which is what a pull request shows;
// when no branch carries the handle (merged and deleted, say), the commits that
// name the handle on this checkout stand in. note says where the diff came
// from, or why there is none, and dir is the checkout 'o' opens files in.
func ticketDiff(root string, t *ticket.Ticket) (rows []diffRow, note, dir string) {
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
			return nil, err.Error(), dir
		}
		note = br + " vs " + base
		rows = parseDiff(patch)
	} else {
		dir = root
		shas, _ := repo.CommitsForTicket(t.Path, t.ID)
		shas = ticket.MergeCommits(shas, t.Commits)
		if len(shas) == 0 {
			return nil, "no branch named after " + handle + ", and no commits name it", dir
		}
		patch, _ := repo.Diff(shas)
		note = "commits naming " + handle + " on this checkout"
		rows = parseDiff(patch)
	}
	if len(rows) == 0 {
		note += " — no changes"
	}
	return rows, note, dir
}

// loadDiff computes the diff of the open ticket once, until a reload clears
// diffFor. The scroll and focus are reset only when the ticket changes, so a
// reload that re-reads the same diff leaves the person where they were.
func (m *Model) loadDiff(t *ticket.Ticket) {
	if m.diffFor == t.ID {
		return
	}
	m.diffRows, m.diffNote, m.diffDir = ticketDiff(m.store.Root, t)
	m.diffFor = t.ID
	if m.diffScrollFor != t.ID {
		m.diffScroll = 0
		m.diffFocus = false
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

// openDiffRow opens the file of the top row in VS Code at its line.
func (m *Model) openDiffRow() {
	if len(m.diffRows) == 0 {
		return
	}
	r := m.diffRows[m.diffTop()]
	if err := openInCode(m.diffDir, r.path, r.line); err != nil {
		m.notify(err.Error(), true)
		return
	}
	m.notify("opened "+r.path+":"+strconv.Itoa(r.line)+" in VS Code", false)
}

// renderDiffPane renders the diff into exactly width x height: a title line,
// then the rows from diffScroll. The file of the top row is highlighted in its
// header. Each row is cut to the width before it is styled, since cutting a
// styled string would break its escape codes.
func (m *Model) renderDiffPane(width, height int) string {
	titleSty := styLaneTitle
	if m.diffFocus {
		titleSty = styReview
	}
	lines := []string{titleSty.Render(truncate("Diff  "+m.diffNote, width))}
	body := max(1, height-1)
	// The scroll may go to the last row, not the last full page: "]" must be
	// able to put a short file's header at the top even when the whole diff
	// would fit.
	m.diffScroll = m.diffTop()
	top := m.diffScroll
	selFile := -1
	if len(m.diffRows) > 0 {
		selFile = m.diffRows[top].file
	}
	end := min(len(m.diffRows), m.diffScroll+body)
	for _, r := range m.diffRows[m.diffScroll:end] {
		text := truncate(r.text, width)
		var sty lipgloss.Style
		switch r.kind {
		case 'H':
			sty = styLaneTitle
			if r.file == selFile {
				sty = stySelected
			}
			text = truncate("▌ "+r.text, width)
		case '@':
			sty = styMeta
		case '+':
			sty = styOK
		case '-':
			sty = styErr
		}
		lines = append(lines, sty.Render(text))
	}
	if len(m.diffRows) == 0 {
		lines = append(lines, styMeta.Render(truncate("nothing to show", width)))
	}
	for i, l := range lines {
		if pad := width - lipgloss.Width(l); pad > 0 {
			lines[i] = l + strings.Repeat(" ", pad)
		}
	}
	for len(lines) < height {
		lines = append(lines, strings.Repeat(" ", width))
	}
	return strings.Join(lines[:height], "\n")
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
