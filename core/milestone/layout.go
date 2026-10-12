package milestone

import (
	"slices"
	"sort"
)

// The layer draws placed milestones as bubbles on a grid. Everything here is a
// pure function of the milestones passed in: the board calls it on every
// redraw, so it must give the same answer for the same input and must never
// fail — a hand-edited file with a typo or a loop in it is drawn as well as it
// can be, not refused.

// Bubble is one placed milestone on the grid. Col is absolute across all
// lanes, Row is the depth below the lane's top, Lane indexes Grid.Lanes.
type Bubble struct {
	Name   string
	Lane   int
	Col    int
	Row    int
	Parent string // "" for a lane top
}

// Lane is one swimlane: the tree under one top, as a span of columns.
type Lane struct {
	Top   string
	Col   int
	Width int
}

// Arrow runs from a predecessor to its successor. Extra is false for the
// parent arrow every non-top bubble has and true for every further
// predecessor; Cross is set on an extra arrow that leaves its lane, and an
// extra arrow without it is a join inside one lane.
type Arrow struct {
	From, To     string
	Extra, Cross bool
}

// Grid is the layout: bubbles ordered by lane, row, column; lanes left to
// right; arrows ordered by their target.
type Grid struct {
	Bubbles []Bubble
	Lanes   []Lane
	Arrows  []Arrow
}

// byAge is the one order the layout uses wherever it has to choose: oldest
// first, the name breaking a tie. Never file order, so the same milestones draw
// the same way on every machine.
func byAge(a, b *Milestone) int {
	if c := a.CreatedAt.Compare(b.CreatedAt); c != 0 {
		return c
	}
	switch {
	case a.Name < b.Name:
		return -1
	case a.Name > b.Name:
		return 1
	}
	return 0
}

// usableEdges returns, for every placed milestone, the predecessors that may
// be drawn: placed, known, and not closing a cycle. Milestones are visited
// oldest first and each arrow is kept only if it does not make a loop with the
// ones kept so far, so when a file contains a loop the arrow written by the
// newer milestone is the one dropped — and the same one on every machine.
func usableEdges(placed []*Milestone) map[string][]string {
	known := make(map[string]bool, len(placed))
	for _, m := range placed {
		known[m.Name] = true
	}
	kept := make(map[string][]string, len(placed))
	for _, m := range placed {
		for _, p := range m.After {
			if !known[p] || reaches(kept, p, m.Name) {
				continue
			}
			kept[m.Name] = append(kept[m.Name], p)
		}
	}
	return kept
}

// reaches reports whether following predecessor edges from `from` arrives at
// `target`, `from == target` included.
func reaches(preds map[string][]string, from, target string) bool {
	return chain(preds, from, target) != nil
}

// chain is reaches with the way there: the names from `from` to `target`
// following predecessor edges, or nil when there is none.
func chain(preds map[string][]string, from, target string) []string {
	seen := map[string]bool{}
	var walk func(n string) []string
	walk = func(n string) []string {
		if n == target {
			return []string{n}
		}
		if seen[n] {
			return nil
		}
		seen[n] = true
		for _, p := range preds[n] {
			if path := walk(p); path != nil {
				return append([]string{n}, path...)
			}
		}
		return nil
	}
	return walk(from)
}

func predecessors(all []*Milestone) map[string][]string {
	preds := make(map[string][]string, len(all))
	for _, m := range all {
		if m.Placed {
			preds[m.Name] = m.After
		}
	}
	return preds
}

// WouldCycle reports whether making `pred` a predecessor of `succ` would close
// a loop — the check the layer runs before it writes an arrow. Unknown names
// and unplaced milestones contribute no edges, as they do not when drawing.
func WouldCycle(all []*Milestone, succ, pred string) bool {
	return CyclePath(all, succ, pred) != nil
}

// CyclePath names the loop WouldCycle refuses, in the direction the arrows
// run and closed on its first name: for a→b→c and a new arrow c→a it is
// [a b c a]. Nil when there is no loop.
func CyclePath(all []*Milestone, succ, pred string) []string {
	path := chain(predecessors(all), pred, succ) // pred, ..., succ: against the arrows
	if path == nil {
		return nil
	}
	slices.Reverse(path)
	return append(path, succ)
}

// ReattachSuccessors computes what taking `name` out of the layer does to the
// milestones that come after it: in each one's after list the name is replaced
// by the removed milestone's own predecessors, order kept, repeats dropped, so
// the chain closes up instead of breaking. A successor left with nothing
// becomes a swimlane top (an empty, non-nil list). Only the milestones that
// name it are in the result, and writing it back is the caller's job.
func ReattachSuccessors(all []*Milestone, name string) map[string][]string {
	var own []string
	for _, m := range all {
		if m.Name == name {
			own = m.After
		}
	}
	out := map[string][]string{}
	for _, m := range all {
		if !m.Placed || !slices.Contains(m.After, name) {
			continue
		}
		next := []string{}
		for _, p := range m.After {
			repl := []string{p}
			if p == name {
				repl = own
			}
			for _, r := range repl {
				if r != m.Name && !slices.Contains(next, r) {
					next = append(next, r)
				}
			}
		}
		out[m.Name] = next
	}
	return out
}

// Layout arranges the placed milestones.
//
//  1. A milestone's parent is its first predecessor, when that is a placed
//     milestone and the arrow closes no loop; otherwise it is a swimlane top.
//  2. A top with everything under it is one swimlane; lanes run left to right
//     by the top's age.
//  3. Inside a lane the tree is tidy: row is depth, children are ordered by
//     age, a subtree is as wide as its children together (at least one column),
//     and a parent stands over its first child.
//  4. Every further usable predecessor is an extra arrow.
func Layout(all []*Milestone) Grid {
	var placed []*Milestone
	for _, m := range all {
		if m.Placed {
			placed = append(placed, m)
		}
	}
	slices.SortFunc(placed, byAge)
	kept := usableEdges(placed)

	parent := map[string]string{}
	children := map[string][]string{}
	var tops []string
	for _, m := range placed { // oldest first, so children come out ordered
		if len(m.After) > 0 && len(kept[m.Name]) > 0 && kept[m.Name][0] == m.After[0] {
			parent[m.Name] = m.After[0]
			children[m.After[0]] = append(children[m.After[0]], m.Name)
		} else {
			tops = append(tops, m.Name)
		}
	}

	width := map[string]int{}
	var measure func(n string) int
	measure = func(n string) int {
		w := 0
		for _, c := range children[n] {
			w += measure(c)
		}
		width[n] = max(w, 1)
		return width[n]
	}

	var g Grid
	col := map[string]int{}
	var place func(n string, lane, start, row int)
	place = func(n string, lane, start, row int) {
		x := start
		for _, c := range children[n] {
			place(c, lane, x, row+1)
			x += width[c]
		}
		bc := start
		if len(children[n]) > 0 {
			bc = col[children[n][0]]
		}
		col[n] = bc
		g.Bubbles = append(g.Bubbles, Bubble{Name: n, Lane: lane, Col: bc, Row: row, Parent: parent[n]})
	}
	laneOf := map[string]int{}
	next := 0
	for i, top := range tops {
		w := measure(top)
		g.Lanes = append(g.Lanes, Lane{Top: top, Col: next, Width: w})
		place(top, i, next, 0)
		next += w
	}
	for _, b := range g.Bubbles {
		laneOf[b.Name] = b.Lane
	}
	sort.SliceStable(g.Bubbles, func(i, j int) bool {
		a, b := g.Bubbles[i], g.Bubbles[j]
		if a.Lane != b.Lane {
			return a.Lane < b.Lane
		}
		if a.Row != b.Row {
			return a.Row < b.Row
		}
		return a.Col < b.Col
	})

	for _, b := range g.Bubbles {
		if b.Parent != "" {
			g.Arrows = append(g.Arrows, Arrow{From: b.Parent, To: b.Name})
		}
		for _, p := range kept[b.Name] {
			if p == b.Parent {
				continue
			}
			g.Arrows = append(g.Arrows, Arrow{From: p, To: b.Name, Extra: true, Cross: laneOf[p] != b.Lane})
		}
	}
	return g
}
