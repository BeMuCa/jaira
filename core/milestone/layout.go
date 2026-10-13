package milestone

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

// The layer draws placed milestones as bubbles on a grid. Everything here is a
// pure function of the milestones passed in: the board calls it on every
// redraw, so it must give the same answer for the same input and must never
// fail — a hand-edited file with a typo or a loop in it is drawn as well as it
// can be, not refused.

// Everything here is keyed by milestone ID, never by name: a name can change
// or be carried by two milestones, an ID is the one thing that stays.

// Bubble is one placed milestone on the grid. Col is absolute across all
// lanes, Row is the depth below the lane's top, Lane indexes Grid.Lanes.
type Bubble struct {
	ID     string
	Lane   int
	Col    int
	Row    int
	Parent string // ID of the parent; "" for a lane top
}

// Lane is one swimlane: the tree under one top, as a span of columns.
type Lane struct {
	Top   string // ID of the lane's top
	Col   int
	Width int
}

// Arrow runs from a predecessor to its successor. Extra is false for the
// parent arrow every non-top bubble has and true for every further
// predecessor — a join inside one lane, or an arrow across lanes.
type Arrow struct {
	From, To string // IDs
	Extra    bool
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
	if c := strings.Compare(a.Name, b.Name); c != 0 {
		return c
	}
	return strings.Compare(a.ID, b.ID)
}

// resolver reads after: tokens against the milestones on the board. A token is
// an ID when it is one and a milestone has it, and otherwise a name — the
// fallback that keeps a hand-written or older file readable — when exactly one
// milestone carries it. Anything else resolves to nothing.
type resolver struct {
	byID   map[string]*Milestone
	byName map[string][]*Milestone
}

func newResolver(all []*Milestone) resolver {
	r := resolver{byID: make(map[string]*Milestone, len(all)), byName: map[string][]*Milestone{}}
	for _, m := range all {
		r.byID[m.ID] = m
		r.byName[m.Name] = append(r.byName[m.Name], m)
	}
	return r
}

// resolve returns the ID a token names, or "".
func (r resolver) resolve(tok string) string {
	if IsID(tok) {
		if m, ok := r.byID[tok]; ok {
			return m.ID
		}
	}
	if ms := r.byName[tok]; len(ms) == 1 {
		return ms[0].ID
	}
	return ""
}

// spell is a token as the ID it names, or as written when it names nothing.
func (r resolver) spell(tok string) string {
	if id := r.resolve(tok); id != "" {
		return id
	}
	return tok
}

// SpellAsID is SpellAsIDs for one token.
func SpellAsID(all []*Milestone, token string) string { return newResolver(all).spell(token) }

// Resolve returns the ID an after: token names among these milestones — the
// token as an ID, or as the one milestone that carries it as a name — or "".
func Resolve(all []*Milestone, token string) string { return newResolver(all).resolve(token) }

// ids resolves a list of tokens, dropping what names nothing and repeats.
func (r resolver) ids(tokens []string) []string {
	var out []string
	for _, t := range tokens {
		if id := r.resolve(t); id != "" && !slices.Contains(out, id) {
			out = append(out, id)
		}
	}
	return out
}

// usableEdges returns, for every placed milestone, the predecessors that may
// be drawn: placed, known, and not closing a cycle. Milestones are visited
// oldest first and each arrow is kept only if it does not make a loop with the
// ones kept so far, so when a file contains a loop the arrow written by the
// newer milestone is the one dropped — and the same one on every machine.
func usableEdges(placed []*Milestone) map[string][]string {
	r := newResolver(placed)
	kept := make(map[string][]string, len(placed))
	for _, m := range placed {
		for _, p := range r.ids(m.After) {
			if chain(kept, p, m.ID) != nil {
				continue
			}
			kept[m.ID] = append(kept[m.ID], p)
		}
	}
	return kept
}

// chain is the way from `from` to `target` following predecessor edges,
// `from == target` included: the IDs along it, or nil when there is none.
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
	r := newResolver(all)
	preds := make(map[string][]string, len(all))
	for _, m := range all {
		if m.Placed {
			preds[m.ID] = r.ids(m.After)
		}
	}
	return preds
}

// CyclePath names, by ID, the loop that making `pred` a predecessor of `succ`
// would close, in the direction the arrows run and closed on its first ID: for
// a→b→c and a new arrow c→a it is [a b c a]. Nil when there is no loop.
// Unknown tokens and unplaced milestones contribute no edges, as they do not
// when drawing.
func CyclePath(all []*Milestone, succ, pred string) []string {
	path := chain(predecessors(all), pred, succ) // pred, ..., succ: against the arrows
	if path == nil {
		return nil
	}
	slices.Reverse(path)
	return append(path, succ)
}

// PredecessorError is CheckAfter's refusal: Pred names no milestone on the
// board, or — Cycle set — making it a predecessor of Name closes that loop.
// Cycle is spelled in names, for a person to read.
type PredecessorError struct {
	Name, Pred string
	Cycle      []string // nil when Pred is unknown
}

func (e *PredecessorError) Error() string {
	if e.Cycle == nil {
		return fmt.Sprintf("no milestone %q on this board to come after", e.Pred)
	}
	return fmt.Sprintf("milestone %q cannot come after %q: that closes a cycle, %s",
		e.Name, e.Pred, strings.Join(e.Cycle, " -> "))
}

// CheckAfter refuses what the layer must never be handed: a predecessor no
// milestone on the board carries, and an arrow that would close a loop. The
// drawing ignores both rather than failing, so every write that places a
// milestone — the CLI's and the layer's — asks this first. id is the
// milestone being placed, and each token is an ID or a name as after: reads
// them.
func CheckAfter(all []*Milestone, id string, tokens []string) error {
	r := newResolver(all)
	name := id
	if m, ok := r.byID[id]; ok {
		name = m.Name
	}
	for _, tok := range tokens {
		pid := r.resolve(tok)
		if pid == "" {
			return &PredecessorError{Name: name, Pred: tok}
		}
		if path := CyclePath(all, id, pid); path != nil {
			names := make([]string, len(path))
			for i, p := range path {
				names[i] = p
				if m, ok := r.byID[p]; ok {
					names[i] = m.Name
				}
			}
			return &PredecessorError{Name: name, Pred: r.byID[pid].Name, Cycle: names}
		}
	}
	return nil
}

// ReattachSuccessors computes what taking the milestone with this ID out of
// the layer does to the milestones that come after it: in each one's after
// list the token naming it is replaced by the removed milestone's own
// predecessors, order kept, repeats dropped, so the chain closes up instead of
// breaking. Every token that names a milestone comes back as that milestone's
// ID; one that names nothing is kept as written. A successor left with nothing
// becomes a swimlane top (an empty, non-nil list). Only the milestones that
// name it are in the result, keyed by ID, and writing it back is the caller's
// job.
func ReattachSuccessors(all []*Milestone, id string) map[string][]string {
	r := newResolver(all)
	var own []string
	if m, ok := r.byID[id]; ok {
		for _, p := range m.After {
			own = append(own, r.spell(p))
		}
	}
	out := map[string][]string{}
	for _, m := range all {
		if !m.Placed || m.ID == id || !slices.Contains(r.ids(m.After), id) {
			continue
		}
		next := []string{}
		for _, p := range m.After {
			repl := []string{r.spell(p)}
			if r.resolve(p) == id {
				repl = own
			}
			for _, x := range repl {
				if x != m.ID && !slices.Contains(next, x) {
					next = append(next, x)
				}
			}
		}
		out[m.ID] = next
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
	r := newResolver(placed)

	parent := map[string]string{}
	children := map[string][]string{}
	var tops []string
	for _, m := range placed { // oldest first, so children come out ordered
		if first := ""; len(m.After) > 0 {
			first = r.resolve(m.After[0])
			if first != "" && len(kept[m.ID]) > 0 && kept[m.ID][0] == first {
				parent[m.ID] = first
				children[first] = append(children[first], m.ID)
				continue
			}
		}
		tops = append(tops, m.ID)
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
		g.Bubbles = append(g.Bubbles, Bubble{ID: n, Lane: lane, Col: bc, Row: row, Parent: parent[n]})
	}
	next := 0
	for i, top := range tops {
		w := measure(top)
		g.Lanes = append(g.Lanes, Lane{Top: top, Col: next, Width: w})
		place(top, i, next, 0)
		next += w
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
			g.Arrows = append(g.Arrows, Arrow{From: b.Parent, To: b.ID})
		}
		for _, p := range kept[b.ID] {
			if p == b.Parent {
				continue
			}
			g.Arrows = append(g.Arrows, Arrow{From: p, To: b.ID, Extra: true})
		}
	}
	return g
}
