package tui

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/BeMuCa/jaira/core/milestone"
	"github.com/BeMuCa/jaira/core/ticket"
)

// The layer is drawn on a canvas of its own — one rune and one style per
// cell — rather than composed from lipgloss blocks: bubbles, arrows and the
// [+] targets interleave on shared rows, and a cell grid is the one model in
// which "this arrow passes under that bubble" is a matter of drawing order.

const (
	bubbleW   = 16 // a bubble, border included
	bubbleH   = 4
	cellW     = 21 // a grid column: the bubble and the 5-column gap after it
	cellH     = 7  // a grid row: the bubble and three connector rows — the split, the extra arrows' row, the ▼
	extraHead = 4  // where an extra arrow's ▼ lands, three left of the parent arrow's at 7
	laneGap   = 3  // extra columns between swimlanes; the ┊ stands in them
	layerLeft = 2  // room left of the first lane for an arrow's vertical run
	layerTop  = 1  // room above the first row for an arrow into a swimlane top
	fillCells = bubbleW - 2
	fillMix   = 0.55 // how far toward white a bubble's fill is mixed
	nameMax   = 12   // the longest name a bubble shows untruncated
)

// layerScene is the picture's geometry: where each bubble stands and which
// targets the cursor can reach. It is built from the milestones on every key
// and every frame, so a reload is reflected without any bookkeeping — the
// layout is a pure function and cheap.
type layerScene struct {
	grid    milestone.Grid
	byID    map[string]*milestone.Milestone
	bubble  map[string]milestone.Bubble
	at      map[string]image.Point // a bubble's top-left cell
	targets []sceneTarget
	farX    int
	w, h    int
}

type sceneTarget struct {
	layerTarget
	box image.Rectangle
}

// layerScene lays the milestones out and settles the cursor on something that
// is still there: a bubble taken out elsewhere hands the cursor to the first
// bubble, and a [+] its bubble no longer offers hands it back to the bubble.
func (m *Model) layerScene() *layerScene {
	l := m.layer
	sc := &layerScene{
		grid:   milestone.Layout(m.milestones),
		byID:   map[string]*milestone.Milestone{},
		bubble: map[string]milestone.Bubble{},
		at:     map[string]image.Point{},
	}
	for _, ms := range m.milestones {
		sc.byID[ms.ID] = ms
	}
	cols, rows := 0, 1
	for _, ln := range sc.grid.Lanes {
		cols = max(cols, ln.Col+ln.Width)
	}
	hasChild := map[string]bool{}
	for _, b := range sc.grid.Bubbles {
		sc.bubble[b.ID] = b
		sc.at[b.ID] = image.Pt(layerLeft+b.Col*cellW+b.Lane*laneGap, layerTop+b.Row*cellH)
		rows = max(rows, b.Row+1)
		if b.Parent != "" {
			hasChild[b.Parent] = true
		}
	}
	sc.farX = layerLeft + cols*cellW + len(sc.grid.Lanes)*laneGap
	sc.w = sc.farX + len(emptyLayerHint) + 6
	sc.h = layerTop + rows*cellH

	// The cursor first, because which [+] are offered depends on it.
	b, ok := sc.bubble[l.sel.id]
	switch {
	case l.sel.kind == onFar:
	case !ok:
		l.sel = layerTarget{kind: onFar}
		if len(sc.grid.Bubbles) > 0 {
			l.sel = layerTarget{kind: onBubble, id: sc.grid.Bubbles[0].ID}
		}
	case l.sel.kind == onRight && b.Parent == "", l.sel.kind == onBelow && hasChild[b.ID]:
		l.sel.kind = onBubble
	}

	for _, b := range sc.grid.Bubbles {
		p := sc.at[b.ID]
		sc.targets = append(sc.targets, sceneTarget{layerTarget{onBubble, b.ID}, image.Rect(p.X, p.Y, p.X+bubbleW, p.Y+bubbleH)})
	}
	// The [+] of the bubble the cursor is on, or of the bubble whose [+] it
	// is on: drawn in space the layout reserves, so offering them moves
	// nothing. Right of a swimlane top there is none — a sibling of a top is
	// a new swimlane, which the far [+] makes — and below a bubble with
	// children there is none, because a further successor is a sibling of
	// those children.
	//
	// Their boxes are what the cursor steers by, and each is centred on the
	// bubble's own axis — the one beside it on its middle, the one under it
	// on its centre column — so ↑↓ never stops on the [+] beside a bubble
	// and ←→ never on the one under it. "[+]" is drawn at the box's corner.
	if l.sel.kind != onFar {
		b := sc.bubble[l.sel.id]
		p := sc.at[b.ID]
		if b.Parent != "" {
			sc.targets = append(sc.targets, sceneTarget{layerTarget{onRight, b.ID}, image.Rect(p.X+bubbleW+1, p.Y+1, p.X+bubbleW+4, p.Y+3)})
		}
		if !hasChild[b.ID] {
			sc.targets = append(sc.targets, sceneTarget{layerTarget{onBelow, b.ID}, image.Rect(p.X+6, p.Y+bubbleH, p.X+10, p.Y+bubbleH+1)})
		}
	}
	// The far [+] steers on the top row's middle, like the [+] beside a
	// bubble, so ↑ from a top bubble does not jump to it.
	sc.targets = append(sc.targets, sceneTarget{layerTarget{kind: onFar}, image.Rect(sc.farX, layerTop+1, sc.farX+3, layerTop+3)})
	return sc
}

func (sc *layerScene) find(t layerTarget) (sceneTarget, bool) {
	for _, st := range sc.targets {
		if st.layerTarget == t {
			return st, true
		}
	}
	return sceneTarget{}, false
}

// step moves the cursor to the nearest target in a direction: of the targets
// whose centre lies that way, the one closest along the direction, with the
// distance across it counting double — so ↓ goes to what is under the cursor
// before it goes to what is merely lower somewhere else. Distances are
// measured in grid steps, a column (21 cells) and a row (6) weighing the
// same, so → from a swimlane top reaches the next top before it reaches a
// child one row down.
func (sc *layerScene) step(from layerTarget, dir string) layerTarget {
	cur, ok := sc.find(from)
	if !ok {
		return from
	}
	centre := func(r image.Rectangle) (int, int) { // in 1/(2·cellW·cellH) of a grid step
		return (r.Min.X*2 + r.Dx()) * cellH, (r.Min.Y*2 + r.Dy()) * cellW
	}
	cx, cy := centre(cur.box)
	best, bestScore := from, 0
	for _, t := range sc.targets {
		if t.layerTarget == from {
			continue
		}
		x, y := centre(t.box)
		dx, dy := x-cx, y-cy
		var along, across int
		switch dir {
		case "right":
			along, across = dx, dy
		case "left":
			along, across = -dx, dy
		case "down":
			along, across = dy, dx
		case "up":
			along, across = -dy, dx
		}
		if along <= 0 {
			continue
		}
		score := along + 2*absInt(across)
		if best == from || score < bestScore {
			best, bestScore = t.layerTarget, score
		}
	}
	return best
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// layerCanvas is a grid of cells, each one rune and an index into styles;
// style 0 is no style at all.
type layerCanvas struct {
	w, h   int
	runes  []rune
	sty    []int
	styles []lipgloss.Style
}

func newLayerCanvas(w, h int) *layerCanvas {
	return &layerCanvas{w: w, h: h, runes: make([]rune, w*h), sty: make([]int, w*h), styles: []lipgloss.Style{{}}}
}

func (cv *layerCanvas) style(s lipgloss.Style) int {
	cv.styles = append(cv.styles, s)
	return len(cv.styles) - 1
}

func (cv *layerCanvas) put(x, y int, r rune, st int) {
	if x < 0 || y < 0 || x >= cv.w || y >= cv.h {
		return
	}
	cv.runes[y*cv.w+x], cv.sty[y*cv.w+x] = r, st
}

func (cv *layerCanvas) text(x, y int, s string, st int) {
	for i, r := range []rune(s) {
		cv.put(x+i, y, r, st)
	}
}

// line renders one row of the canvas from column x0, w cells wide, each run
// of one style in one piece.
func (cv *layerCanvas) line(y, x0, w int) string {
	if y < 0 || y >= cv.h {
		return ""
	}
	var b strings.Builder
	var run []rune
	cur := 0
	flush := func() {
		if len(run) == 0 {
			return
		}
		if cur == 0 {
			b.WriteString(string(run))
		} else {
			b.WriteString(cv.styles[cur].Render(string(run)))
		}
		run = run[:0]
	}
	for x := max(0, x0); x < min(cv.w, x0+w); x++ {
		r, st := cv.runes[y*cv.w+x], cv.sty[y*cv.w+x]
		if r == 0 {
			r = ' '
		}
		if st != cur {
			flush()
			cur = st
		}
		run = append(run, r)
	}
	flush()
	return strings.TrimRight(b.String(), " ")
}

// The sides of a cell a line can leave by, as bits, and the glyph for every
// combination: one arrow's corner and another's straight run through the same
// cell make a junction, not a broken line.
const (
	sideUp = 1 << iota
	sideDown
	sideLeft
	sideRight
)

var sideGlyph = [16]rune{
	sideUp: '│', sideDown: '│', sideUp | sideDown: '│',
	sideLeft: '─', sideRight: '─', sideLeft | sideRight: '─',
	sideUp | sideRight: '└', sideUp | sideLeft: '┘', sideDown | sideRight: '┌', sideDown | sideLeft: '┐',
	sideUp | sideDown | sideRight: '├', sideUp | sideDown | sideLeft: '┤',
	sideDown | sideLeft | sideRight: '┬', sideUp | sideLeft | sideRight: '┴',
	sideUp | sideDown | sideLeft | sideRight: '┼',
}

// arrowLines collects the layer's arrows as line pieces per cell — which
// sides of the cell a line leaves by, and whether a parent arrow (strong) or
// only extra arrows use it — and the cells their ▼ sit in.
type arrowLines struct {
	sides  map[image.Point]int
	strong map[image.Point]bool
	heads  map[image.Point]bool // ▼, strong when a parent arrow ends there
}

func newArrowLines() *arrowLines {
	return &arrowLines{sides: map[image.Point]int{}, strong: map[image.Point]bool{}, heads: map[image.Point]bool{}}
}

func (l *arrowLines) mark(p image.Point, side int, strong bool) {
	l.sides[p] |= side
	l.strong[p] = l.strong[p] || strong
}

// link joins two cells on one axis.
func (l *arrowLines) link(a, b image.Point, strong bool) {
	d := image.Pt(sign(b.X-a.X), sign(b.Y-a.Y))
	fwd, back := sideDown, sideUp
	switch {
	case d.Y < 0:
		fwd, back = sideUp, sideDown
	case d.X > 0:
		fwd, back = sideRight, sideLeft
	case d.X < 0:
		fwd, back = sideLeft, sideRight
	}
	for p := a; p != b; p = p.Add(d) {
		l.mark(p, fwd, strong)
		l.mark(p.Add(d), back, strong)
	}
}

// path joins axis-aligned points in order.
func (l *arrowLines) path(pts []image.Point, strong bool) {
	for i := 0; i+1 < len(pts); i++ {
		l.link(pts[i], pts[i+1], strong)
	}
}

func (l *arrowLines) arrow(p image.Point, strong bool) {
	l.heads[p] = l.heads[p] || strong
}

// draw puts the lines on the canvas: plain where a parent arrow runs, dim
// where only extra arrows do. Arrow heads go last, over the line ending in
// them.
func (l *arrowLines) draw(cv *layerCanvas, dim int) {
	style := func(strong bool) int {
		if strong {
			return 0
		}
		return dim
	}
	for p, s := range l.sides {
		cv.put(p.X, p.Y, sideGlyph[s], style(l.strong[p]))
	}
	for p, strong := range l.heads {
		cv.put(p.X, p.Y, '▼', style(strong))
	}
}

func sign(n int) int {
	switch {
	case n > 0:
		return 1
	case n < 0:
		return -1
	}
	return 0
}

// bubbleFill is a bubble's fill: the milestone's colour mixed fillMix of the
// way toward white, and the ink that reads on it. A 24-bit terminal gets the
// mix itself; any other the nearest of the 256 colours. The ink is decided
// on the mix, not on its stand-in, so a bubble keeps black or white text
// across both kinds of terminal.
func (m *Model) bubbleFill(colour int) (fill, ink color.Color) {
	r, g, b := paletteRGB(colour)
	mix := func(c int) int { return int(math.Round(float64(c) + float64(255-c)*fillMix)) }
	fr, fg, fb := mix(r), mix(g), mix(b)
	fill = lipgloss.Color(strconv.Itoa(nearestPalette(fr, fg, fb)))
	if m.trueColor {
		fill = color.RGBA{R: uint8(fr), G: uint8(fg), B: uint8(fb), A: 0xff}
	}
	ink = inkWhite
	if luminance(fr, fg, fb) > 0.35 {
		ink = inkBlack
	}
	return fill, ink
}

var inkBlack, inkWhite = lipgloss.Color("0"), lipgloss.Color("15")

// luminance is the relative luminance of an sRGB colour, 0 for black and 1
// for white — the measure the contrast rules of WCAG are written in.
func luminance(r, g, b int) float64 {
	lin := func(c int) float64 {
		v := float64(c) / 255
		if v <= 0.03928 {
			return v / 12.92
		}
		return math.Pow((v+0.055)/1.055, 2.4)
	}
	return 0.2126*lin(r) + 0.7152*lin(g) + 0.0722*lin(b)
}

// filledCells is how many of a bubble's inner cells are filled for done of
// total. Halfway work never looks untouched or finished: one ticket done of
// forty is one cell, and one open ticket of forty leaves one cell empty.
func filledCells(done, total int) int {
	if total <= 0 {
		return 0
	}
	n := int(math.Round(float64(fillCells*done) / float64(total)))
	if done > 0 && done < total {
		n = min(max(n, 1), fillCells-1)
	}
	return n
}

// drawBubble draws one milestone: the border in its colour (thick and bold
// when selected), its name, done/total, and the fill behind both rows.
func (m *Model) drawBubble(cv *layerCanvas, ms *milestone.Milestone, at image.Point, done, total int, sel bool) {
	fill, ink := m.bubbleFill(ms.Colour)
	bst := cv.style(bubbleBorder(ms, sel))
	tst := 0
	if sel {
		tst = cv.style(lipgloss.NewStyle().Bold(true))
	}
	fst := cv.style(lipgloss.NewStyle().Background(fill).Foreground(ink).Bold(sel))

	tl, tr, bl, br, h, v := '╭', '╮', '╰', '╯', "─", '│'
	if sel {
		tl, tr, bl, br, h, v = '┏', '┓', '┗', '┛', "━", '┃'
	}
	x, y := at.X, at.Y
	cv.put(x, y, tl, bst)
	cv.text(x+1, y, strings.Repeat(h, fillCells), bst)
	cv.put(x+bubbleW-1, y, tr, bst)

	name := []rune(m.milestoneLabel(ms.ID))
	if len(name) > nameMax {
		name = append(name[:nameMax-1], '…')
	}
	count := fmt.Sprintf("%d/%d ", done, total)
	// A ⚠ where the name's leading space is: two milestones share it.
	lead := " "
	if m.isTwin(ms.ID) {
		lead = "⚠"
	}
	rowsText := [2]string{
		fmt.Sprintf("%-*s", fillCells, lead+string(name)),
		fmt.Sprintf("%*s", fillCells, count),
	}
	cells := filledCells(done, total)
	for k, text := range rowsText {
		row := y + 1 + k
		cv.put(x, row, v, bst)
		for i, r := range []rune(text)[:fillCells] {
			st := tst
			if i < cells {
				st = fst
			}
			cv.put(x+1+i, row, r, st)
		}
		cv.put(x+bubbleW-1, row, v, bst)
	}
	cv.put(x, y+3, bl, bst)
	cv.text(x+1, y+3, strings.Repeat(h, fillCells), bst)
	cv.put(x+bubbleW-1, y+3, br, bst)
}

// bubbleExit is where an arrow leaves a bubble's bottom border, drawn in the
// border's own style so the line visibly belongs to it.
func bubbleExit(cv *layerCanvas, at image.Point, dx int, sel bool, ms *milestone.Milestone) {
	tee := '┬'
	if sel {
		tee = '┯'
	}
	cv.put(at.X+dx, at.Y+bubbleH-1, tee, cv.style(bubbleBorder(ms, sel)))
}

// bubbleBorder is a bubble's border style: the milestone's colour, bold when
// selected.
func bubbleBorder(ms *milestone.Milestone, sel bool) lipgloss.Style {
	st := lipgloss.NewStyle().Bold(sel)
	if ms.HasColour() {
		st = st.Foreground(lipgloss.Color(strconv.Itoa(ms.Colour)))
	}
	return st
}

// roadmapTitle is what the layer is called on screen: Berk named it the
// Milestone Roadmap (2026-10-13); the code keeps calling it the layer.
const roadmapTitle = "Milestone Roadmap"

const emptyLayerHint = "enter: place the first milestone"

// drawLayer paints the whole picture: swimlane separators first, then the
// extra arrows, the parent arrows, the bubbles and last the [+] — each layer
// over the one before, so where they meet the more important thing shows.
func (m *Model) drawLayer(sc *layerScene) *layerCanvas {
	l := m.layer
	cv := newLayerCanvas(sc.w, sc.h)
	dim := cv.style(lipgloss.NewStyle().Foreground(colFaint))

	for i, ln := range sc.grid.Lanes {
		if i == 0 {
			continue
		}
		x := layerLeft + ln.Col*cellW + i*laneGap - laneGap
		for y := 0; y < sc.h; y++ {
			cv.put(x, y, '┊', dim)
		}
	}

	// Arrows are collected as line pieces first and turned into glyphs once,
	// so where two arrows meet the cell shows the junction (├ ┤ ┬ ┴ ┼) rather
	// than whichever arrow happened to be drawn last. A cell any parent arrow
	// passes through is drawn plain; one only extra arrows use, dim.
	type exit struct {
		id string
		dx int
	}
	var exits []exit
	lines := newArrowLines()

	// Parent arrows: down from the centre into the first connector row, along
	// it to every child — a split — and down past the second row to the ▼.
	for _, b := range sc.grid.Bubbles {
		if b.Parent == "" {
			continue
		}
		p, t := sc.at[b.Parent], sc.at[b.ID]
		x, row, tx := p.X+7, p.Y+bubbleH, t.X+7
		lines.path([]image.Point{{x, row - 1}, {x, row}, {tx, row}, {tx, row + 2}}, true)
		lines.arrow(image.Pt(tx, row+2), true)
		exits = append(exits, exit{b.Parent, 7})
	}

	// Extra arrows travel the second connector row, which no split uses, so
	// one can only ever cross another arrow's drop (┼) — never run along a
	// split and seem to end at a sibling's junction. They arrive with a ▼ of
	// their own, left of the target's centre: ending in the parent's drop
	// instead made an arrival look like one more crossing (Berk, 2026-10-13).
	// Each leaves right of centre, down to that row, along it — straight to
	// its ▼ when the target is in the next row, a join; else to the gap left
	// of the target, down (or up) that gap to the row above the target, and
	// in to its ▼.
	for _, a := range sc.grid.Arrows {
		if !a.Extra {
			continue
		}
		p, t := sc.at[a.From], sc.at[a.To]
		sx, row, tx := p.X+11, p.Y+bubbleH+1, t.X+extraHead
		pts := []image.Point{{sx, row - 2}, {sx, row}}
		if t.Y == p.Y+cellH {
			pts = append(pts, image.Pt(tx, row), image.Pt(tx, t.Y-1))
		} else {
			gx := t.X - 2
			pts = append(pts, image.Pt(gx, row), image.Pt(gx, t.Y-1), image.Pt(tx, t.Y-1))
		}
		lines.path(pts, false)
		lines.arrow(image.Pt(tx, t.Y-1), false)
		exits = append(exits, exit{a.From, 11})
	}
	lines.draw(cv, dim)

	selBubble := ""
	if l.sel.kind == onBubble {
		selBubble = l.sel.id
	}
	for _, b := range sc.grid.Bubbles {
		ms := sc.byID[b.ID]
		done, total := m.milestoneDone(ms)
		m.drawBubble(cv, ms, sc.at[b.ID], done, total, b.ID == selBubble)
	}
	for _, e := range exits {
		bubbleExit(cv, sc.at[e.id], e.dx, e.id == selBubble, sc.byID[e.id])
	}

	for _, t := range sc.targets {
		if t.kind == onBubble {
			continue
		}
		st := dim
		if t.layerTarget == l.sel {
			st = cv.style(stySelected)
		}
		cv.text(t.box.Min.X, t.box.Min.Y, "[+]", st)
	}
	if len(sc.grid.Bubbles) == 0 {
		cv.text(sc.farX+4, layerTop+1, emptyLayerHint, dim)
	}
	return cv
}

// renderLayer is the layer as it is on screen: the picture, or the open
// milestone, with a dialog over it while one is open.
func (m *Model) renderLayer() string {
	if m.layer == nil {
		return m.renderBoard()
	}
	switch m.layer.stage {
	case stageChoose, stageForm, stagePick, stageRemove:
		return m.modalOver(m.renderLayerDialog(), modeLayer)
	}
	return m.renderLayerBase()
}

// renderLayerBase is what a dialog over the layer covers.
func (m *Model) renderLayerBase() string {
	if m.layer.stage == stageDetail || (m.layer.stage == stageForm && m.formReturn() == stageDetail) {
		return m.renderMilestoneDetail()
	}
	return m.renderLayerGrid()
}

// renderLayerGrid draws the picture through a window the size of the
// terminal, moved only as far as it must to keep the selection in it.
func (m *Model) renderLayerGrid() string {
	l := m.layer
	sc := m.layerScene()
	cv := m.drawLayer(sc)

	hints := []string{"←↑↓→ move", "enter open", "e edit", "c colour", "x remove", "L board"}
	if l.sel.kind != onBubble {
		hints = []string{"←↑↓→ move", "enter new or connect", "L board"}
	}
	width := max(1, m.width)
	footer := wrapHints(hints, width)
	vw, vh := width, max(1, m.height-1-len(footer))

	cur, _ := sc.find(l.sel)
	box := cur.box
	if l.sel.kind == onBubble {
		// The bubble and the [+] beside and under it, so the targets one
		// arrow away are on screen too.
		box.Max = box.Max.Add(image.Pt(5, 2))
	}
	l.panX = min(l.panX, max(0, sc.w-vw))
	l.panY = min(l.panY, max(0, sc.h-vh))
	if box.Max.X > l.panX+vw {
		l.panX = box.Max.X - vw
	}
	if box.Min.X < l.panX {
		l.panX = box.Min.X
	}
	if box.Max.Y > l.panY+vh {
		l.panY = box.Max.Y - vh
	}
	if box.Min.Y < l.panY {
		l.panY = box.Min.Y
	}
	l.panX, l.panY = max(0, l.panX), max(0, l.panY)

	var b strings.Builder
	b.WriteString(m.layerHeader(sc) + "\n")
	for y := l.panY; y < l.panY+vh; y++ {
		b.WriteString(cv.line(y, l.panX, vw) + "\n")
	}
	for i, f := range footer {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(styMeta.Render(f))
	}
	return b.String()
}

// layerHeader names the screen and what the cursor is on: a milestone's goal,
// or what a [+] would create — the one line of words the picture needs.
func (m *Model) layerHeader(sc *layerScene) string {
	l := m.layer
	head := styLaneTitle.Render(roadmapTitle)
	width := max(1, m.width)
	switch l.sel.kind {
	case onBubble:
		ms := sc.byID[l.sel.id]
		goal := "no goal yet — e writes one"
		if ms.Goal != "" {
			goal = ms.Goal
		}
		// Cut before styling: a cut through an escape sequence would leave
		// the colour running into the picture.
		label := m.milestoneLabel(ms.ID)
		goal = truncate(goal, width-lipgloss.Width(roadmapTitle+" ·  · "+label))
		return head + styMeta.Render(" · ") + label + styMeta.Render(" · "+goal)
	default:
		return head + styMeta.Render(truncate(" · [+] a new milestone, "+m.placementText(l.sel), width-len(roadmapTitle)))
	}
}

// renderLayerDialog is the box over the layer: the choice behind a [+], the
// form, the connect picker, or the question before a bubble leaves.
func (m *Model) renderLayerDialog() string {
	l := m.layer
	var b strings.Builder
	errLine := func() {
		if l.err != "" {
			b.WriteString("\n" + styleLines(styErr, wrap(l.err, max(20, m.width-14), 0)) + "\n")
		}
	}
	hints := func(items ...string) {
		for _, h := range wrapHints(items, max(20, m.width-12)) {
			b.WriteString("\n" + styMeta.Render(h))
		}
	}
	switch l.stage {
	case stageChoose:
		b.WriteString(styLaneTitle.Render("[+] "+m.placementText(l.sel)) + "\n\n")
		b.WriteString(styHelpKey.Render("n") + "  a new milestone\n")
		b.WriteString(styHelpKey.Render("c") + "  connect one that exists\n")
		hints("esc back")

	case stageForm:
		f := l.form
		title := "New milestone, " + m.placementText(l.sel)
		if f.edit != "" {
			title = "Edit " + m.milestoneLabel(f.editID)
		}
		b.WriteString(styLaneTitle.Render(title) + "\n\n")
		for i, label := range f.labels {
			val := f.vals[i]
			if i == f.field {
				val += stySelected.Render("▏")
			}
			b.WriteString(styMeta.Render(fmt.Sprintf("%-6s", label)) + " " + val + "\n")
		}
		if f.edit != "" {
			b.WriteString(styMeta.Render("       comma-separated; empty makes it the top of a swimlane") + "\n")
		}
		errLine()
		hints("tab next field", "enter save", "esc cancel")

	case stagePick:
		b.WriteString(styLaneTitle.Render("Connect, "+m.placementText(l.sel)) + "\n\n")
		if len(l.pick) == 0 {
			b.WriteString(styMeta.Render("Nothing to connect here: every other milestone already comes after it, or would close a loop.") + "\n")
		}
		for i, id := range l.pick {
			ms := m.milestoneByID(id)
			marker, label := "  ", m.milestoneLabel(id)
			if i == l.pickIdx {
				marker, label = stySelected.Render("▌ "), stySelected.Render(label)
			}
			note := "not on the roadmap yet: placed here"
			switch {
			case ms == nil || !ms.Placed:
			case len(ms.After) == 0:
				// A top's first predecessor decides its swimlane, so an
				// arrow into it cannot be an extra one: it moves.
				note = "top of a swimlane: moves here, with what comes after it"
			default:
				note = "stays where it is, gains the arrow"
			}
			b.WriteString(marker + milestoneMark(ms) + " " + label + styMeta.Render("  "+note) + "\n")
		}
		errLine()
		hints("↑↓ pick", "enter connect", "esc back")

	case stageRemove:
		b.WriteString(styLaneTitle.Render("Take "+m.milestoneLabel(l.sel.id)+" off the roadmap?") + "\n\n")
		b.WriteString(styMeta.Render("Its tickets stay in it. What came after it moves up to what came before it.") + "\n")
		hints("y take it out", "n keep it")
	}
	return b.String()
}

// renderMilestoneDetail is one milestone in full: goal, place in the layer,
// and a row per ticket — a list that can be longer than the window, so it
// scrolls with its cursor.
func (m *Model) renderMilestoneDetail() string {
	l := m.layer
	ms := m.milestoneByID(l.detailID)
	if ms == nil {
		return m.renderLayerGrid()
	}
	width := max(20, m.width)
	var head []string
	done, total := m.milestoneDone(ms)
	head = append(head, milestoneMark(ms)+" "+styLaneTitle.Render(m.milestoneLabel(ms.ID))+styMeta.Render(fmt.Sprintf("  %d/%d done", done, total)))
	head = append(head, styBar.Render(strings.Repeat("─", min(width, 78))))
	goal := ms.Goal
	if goal == "" {
		goal = styMeta.Render("no goal yet — e writes one")
	}
	head = append(head, styMeta.Render("goal   ")+truncate(goal, width-7))
	after := "top of its own swimlane"
	switch {
	case !ms.Placed:
		after = "not on the roadmap"
	case len(ms.After) > 0:
		after = strings.Join(ms.After, ", ")
	}
	head = append(head, styMeta.Render("after  ")+after, "")
	if len(l.rows) == 0 {
		head = append(head, styMeta.Render("No tickets in it yet — 'jaira milestone add "+m.milestoneLabel(ms.ID)+" <id>' puts one in."))
	}

	const whereW = 16
	titleW := max(10, width-2-7-2-whereW-1)
	var rows []string
	for i, r := range l.rows {
		marker, handle := "  ", fmt.Sprintf("%-7s", ticket.Handle(r.id))
		title := truncate(r.title, titleW)
		if i == l.rowIdx {
			marker, title = stySelected.Render("▌ "), stySelected.Render(title)
		}
		pad := strings.Repeat(" ", max(1, titleW-lipgloss.Width(title)+2))
		rows = append(rows, marker+styHandle.Render(handle)+title+pad+styMeta.Render(truncate(r.where, whereW)))
	}

	footer := wrapHints([]string{"↑↓ pick", "enter open", "e edit", "b board", "esc back"}, width)
	visible := max(1, m.height-1-len(footer)-len(head))
	if l.rowIdx < l.rowTop {
		l.rowTop = l.rowIdx
	}
	if l.rowIdx >= l.rowTop+visible {
		l.rowTop = l.rowIdx - visible + 1
	}
	l.rowTop = max(0, min(l.rowTop, len(rows)-visible))

	var b strings.Builder
	for _, h := range head {
		b.WriteString(h + "\n")
	}
	end := min(len(rows), l.rowTop+visible)
	for _, r := range rows[l.rowTop:end] {
		b.WriteString(r + "\n")
	}
	for i := end - l.rowTop; i < visible; i++ {
		b.WriteString("\n")
	}
	for i, f := range footer {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(styMeta.Render(f))
	}
	return b.String()
}
