package milestone

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

// ms builds a placed milestone the way a file would read: created n minutes
// after a fixed start, so tests state order by number.
func ms(name string, minute int, after ...string) *Milestone {
	m := New(name, 45, time.Date(2026, 10, 1, 0, minute, 0, 0, time.UTC))
	m.SetAfter(after)
	return m
}

func at(t *testing.T, g Grid, name string) Bubble {
	t.Helper()
	for _, b := range g.Bubbles {
		if b.Name == name {
			return b
		}
	}
	t.Fatalf("no bubble %q in %+v", name, g.Bubbles)
	return Bubble{}
}

func pos(t *testing.T, g Grid, name string) string {
	b := at(t, g, name)
	return fmt.Sprintf("lane%d col%d row%d", b.Lane, b.Col, b.Row)
}

func TestLayoutEmptyAndUnplaced(t *testing.T) {
	if g := Layout(nil); len(g.Bubbles) != 0 || len(g.Lanes) != 0 || len(g.Arrows) != 0 {
		t.Errorf("empty input: %+v", g)
	}
	loose := New("loose", 45, time.Now()) // no after line: not in the layer
	if g := Layout([]*Milestone{loose}); len(g.Bubbles) != 0 {
		t.Errorf("an unplaced milestone must not be drawn: %+v", g.Bubbles)
	}
}

func TestLayoutSingleLaneIsAColumn(t *testing.T) {
	g := Layout([]*Milestone{ms("a", 0), ms("b", 1, "a"), ms("c", 2, "b")})
	for name, want := range map[string]string{"a": "lane0 col0 row0", "b": "lane0 col0 row1", "c": "lane0 col0 row2"} {
		if got := pos(t, g, name); got != want {
			t.Errorf("%s at %s, want %s", name, got, want)
		}
	}
	if len(g.Lanes) != 1 || g.Lanes[0].Top != "a" || g.Lanes[0].Col != 0 || g.Lanes[0].Width != 1 {
		t.Errorf("lanes: %+v", g.Lanes)
	}
	if at(t, g, "c").Parent != "b" || at(t, g, "a").Parent != "" {
		t.Error("parents wrong")
	}
}

// A split: two children side by side under one parent, the parent over the
// first child, children ordered by created-at then name.
func TestLayoutSplit(t *testing.T) {
	g := Layout([]*Milestone{ms("a", 0), ms("z", 1, "a"), ms("b", 1, "a"), ms("late", 5, "a")})
	want := map[string]string{
		"a": "lane0 col0 row0", "b": "lane0 col0 row1", "z": "lane0 col1 row1", "late": "lane0 col2 row1",
	}
	for n, w := range want {
		if got := pos(t, g, n); got != w {
			t.Errorf("%s at %s, want %s", n, got, w)
		}
	}
	if g.Lanes[0].Width != 3 {
		t.Errorf("width %d, want 3", g.Lanes[0].Width)
	}
}

func TestLayoutSubtreeWidthPushesTheSibling(t *testing.T) {
	// b has two children, so c (b's later sibling) starts at column 2.
	g := Layout([]*Milestone{ms("a", 0), ms("b", 1, "a"), ms("c", 2, "a"), ms("b1", 3, "b"), ms("b2", 4, "b")})
	if got := pos(t, g, "c"); got != "lane0 col2 row1" {
		t.Errorf("c at %s", got)
	}
	if got := pos(t, g, "b"); got != "lane0 col0 row1" {
		t.Errorf("b at %s", got)
	}
	if got := pos(t, g, "b2"); got != "lane0 col1 row2" {
		t.Errorf("b2 at %s", got)
	}
}

// A join: h has two predecessors in one lane. The first is the parent, the
// second an extra arrow flagged as a join.
func TestLayoutJoin(t *testing.T) {
	g := Layout([]*Milestone{ms("a", 0), ms("c", 1, "a"), ms("g", 2, "a"), ms("h", 3, "c", "g")})
	if at(t, g, "h").Parent != "c" {
		t.Errorf("parent = %q", at(t, g, "h").Parent)
	}
	var extra []Arrow
	for _, a := range g.Arrows {
		if a.Extra {
			extra = append(extra, a)
		}
	}
	if len(extra) != 1 || extra[0] != (Arrow{From: "g", To: "h", Extra: true}) {
		t.Errorf("extra arrows: %+v", extra)
	}
	if len(g.Arrows) != 4 { // three parent arrows (c,g,h) + the join
		t.Errorf("arrows: %+v", g.Arrows)
	}
}

func TestLayoutCrossLaneArrowAndLaneOrder(t *testing.T) {
	// Lane order follows the top's created-at, not the name: "zed" is older.
	g := Layout([]*Milestone{ms("alpha", 5), ms("zed", 1), ms("x", 6, "zed", "alpha")})
	if g.Lanes[0].Top != "zed" || g.Lanes[1].Top != "alpha" {
		t.Fatalf("lanes: %+v", g.Lanes)
	}
	if g.Lanes[1].Col != 1 {
		t.Errorf("second lane starts at col %d", g.Lanes[1].Col)
	}
	if at(t, g, "x").Lane != 0 {
		t.Error("x belongs to its first predecessor's lane")
	}
	found := false
	for _, a := range g.Arrows {
		if a.From == "alpha" && a.To == "x" && a.Extra {
			found = true
		}
	}
	if !found {
		t.Errorf("no cross-lane arrow alpha->x: %+v", g.Arrows)
	}
}

func TestLayoutLaneTiesBreakByName(t *testing.T) {
	g := Layout([]*Milestone{ms("b", 0), ms("a", 0)})
	if g.Lanes[0].Top != "a" {
		t.Errorf("lanes: %+v", g.Lanes)
	}
}

func TestLayoutIgnoresUnknownAndUnplacedPredecessors(t *testing.T) {
	loose := New("loose", 45, time.Now())
	g := Layout([]*Milestone{ms("a", 0, "ghost"), ms("b", 1, "loose", "a"), loose})
	if at(t, g, "a").Parent != "" || len(g.Lanes) != 2 {
		t.Errorf("a with an unknown predecessor must be a lane top: %+v", g)
	}
	// b's first predecessor is unusable, so b is a top, and the usable second
	// one is still an arrow.
	if at(t, g, "b").Parent != "" {
		t.Error("b must be a lane top")
	}
	if len(g.Arrows) != 1 || g.Arrows[0].From != "a" || !g.Arrows[0].Extra {
		t.Errorf("arrows: %+v", g.Arrows)
	}
}

func TestLayoutCycleNeverCrashesAndDropsOneArrow(t *testing.T) {
	g := Layout([]*Milestone{ms("a", 0, "b"), ms("b", 1, "a"), ms("self", 2, "self")})
	if len(g.Bubbles) != 3 {
		t.Fatalf("bubbles: %+v", g.Bubbles)
	}
	n := 0
	for _, a := range g.Arrows {
		if (a.From == "a" && a.To == "b") || (a.From == "b" && a.To == "a") {
			n++
		}
	}
	if n != 1 {
		t.Errorf("exactly one of the two cycle arrows is kept, got %d: %+v", n, g.Arrows)
	}
	if at(t, g, "self").Parent != "" {
		t.Error("a milestone after itself is a lane top")
	}
}

func TestLayoutIsDeterministic(t *testing.T) {
	in := func() []*Milestone {
		return []*Milestone{ms("d", 3, "b", "c"), ms("c", 2, "a"), ms("b", 1, "a"), ms("a", 0), ms("e", 0)}
	}
	first := Layout(in())
	rev := in()
	for i, j := 0, len(rev)-1; i < j; i, j = i+1, j-1 {
		rev[i], rev[j] = rev[j], rev[i]
	}
	if !reflect.DeepEqual(first, Layout(rev)) {
		t.Error("input order changed the layout")
	}
}

func TestCyclePathFindsEveryLoop(t *testing.T) {
	all := []*Milestone{ms("a", 0), ms("b", 1, "a"), ms("c", 2, "b"), ms("x", 3)}
	cases := []struct {
		succ, pred string
		want       bool
	}{
		{"a", "c", true}, // a after c, c is downstream of a
		{"a", "b", true},
		{"a", "a", true},  // itself
		{"c", "a", false}, // a second arrow from upstream is fine
		{"x", "c", false},
		{"a", "x", false},
	}
	for _, c := range cases {
		if got := CyclePath(all, c.succ, c.pred) != nil; got != c.want {
			t.Errorf("CyclePath(%s after %s) != nil is %v, want %v", c.succ, c.pred, got, c.want)
		}
	}
}

func TestReattachSuccessors(t *testing.T) {
	all := []*Milestone{
		ms("p", 0), ms("q", 1),
		ms("m", 2, "p", "q"),
		ms("s1", 3, "m"),           // gets m's predecessors
		ms("s2", 4, "x", "m", "q"), // order kept, q not duplicated
		ms("s3", 5, "p"),           // untouched, not in the result
		ms("t", 6),
		ms("s4", 7, "t", "r"),
	}
	got := ReattachSuccessors(all, "m")
	want := map[string][]string{"s1": {"p", "q"}, "s2": {"x", "p", "q"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	// A successor whose only predecessor was the removed top has none left and
	// becomes a lane top: an empty, non-nil list.
	got = ReattachSuccessors([]*Milestone{ms("m", 0), ms("s", 1, "m")}, "m")
	if v, ok := got["s"]; !ok || v == nil || len(v) != 0 {
		t.Errorf("got %#v, want s -> []", got)
	}
	// A hand-edited file naming itself is not its own successor: taking it
	// out must not put it straight back in as a top.
	if got = ReattachSuccessors([]*Milestone{ms("m", 0, "m")}, "m"); len(got) != 0 {
		t.Errorf("got %#v, want nothing", got)
	}
}

func TestCyclePathNamesTheLoopInArrowDirection(t *testing.T) {
	all := []*Milestone{ms("a", 0), ms("b", 1, "a"), ms("c", 2, "b")}
	if got := strings.Join(CyclePath(all, "a", "c"), " "); got != "a b c a" {
		t.Errorf("got %q", got)
	}
	if CyclePath(all, "c", "a") != nil {
		t.Error("no loop expected")
	}
	if got := CyclePath(all, "a", "a"); len(got) != 2 {
		t.Errorf("self: %v", got)
	}
}

func TestCheckAfter(t *testing.T) {
	all := []*Milestone{ms("a", 0), ms("b", 1, "a"), ms("c", 2, "b")}
	if err := CheckAfter(all, "c", []string{"a", "b"}); err != nil {
		t.Errorf("valid predecessors refused: %v", err)
	}
	var pe *PredecessorError
	if err := CheckAfter(all, "c", []string{"a", "ghost"}); !errors.As(err, &pe) || pe.Pred != "ghost" || pe.Cycle != nil {
		t.Errorf("unknown: %v", err)
	} else if !strings.Contains(err.Error(), `no milestone "ghost"`) {
		t.Errorf("unknown message: %v", err)
	}
	if err := CheckAfter(all, "a", []string{"c"}); !errors.As(err, &pe) || strings.Join(pe.Cycle, " ") != "a b c a" {
		t.Errorf("cycle: %v", err)
	} else if !strings.Contains(err.Error(), "a -> b -> c -> a") {
		t.Errorf("cycle message: %v", err)
	}
}
