package link

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// store builds an empty board in a temporary directory.
func store(t *testing.T) *ticket.Store {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("JAIRA_HOME", filepath.Join(dir, "home"))
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	return s
}

func lanes(t *testing.T) *lane.Set {
	t.Helper()
	set, err := lane.Load("")
	if err != nil {
		t.Fatal(err)
	}
	return set
}

// make writes one ticket with the given extra frontmatter and returns its id.
func make1(t *testing.T, s *ticket.Store, title, status string, extra map[string]string, lists map[string][]string) string {
	t.Helper()
	fields := map[string]string{
		ticket.FieldID:     ticket.NewID(time.Now()),
		ticket.FieldTitle:  title,
		ticket.FieldStatus: status,
	}
	for k, v := range extra {
		fields[k] = v
	}
	tk, err := s.Create(fields, lists, "")
	if err != nil {
		t.Fatal(err)
	}
	return tk.ID
}

func list(t *testing.T, s *ticket.Store) []*ticket.Ticket {
	t.Helper()
	all, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	return all
}

func kinds(entries []Entry, k Kind) []string {
	var out []string
	for _, e := range entries {
		if e.Kind == k {
			out = append(out, e.Ref.ID)
		}
	}
	return out
}

// A blocker that finished and was filed into the logbook is a dependency that
// cleared. Before the index existed it read as "does not exist in this store"
// — finishing the blocker was the thing that broke the ticket waiting on it.
func TestBlockerInTheLogbookIsSatisfied(t *testing.T) {
	s := store(t)
	blocker := make1(t, s, "the blocker", "done", nil, nil)
	if _, err := s.Logbook(blocker, "as-20260911"); err != nil {
		t.Fatal(err)
	}

	ix := Build(s, lanes(t), list(t, s))
	if !ix.Known(blocker) {
		t.Errorf("a logbooked ticket must still be known, %s was not", ticket.Handle(blocker))
	}
	if !ix.Satisfied(blocker) {
		t.Errorf("a logbooked ticket is finished work; %s must count as satisfied", ticket.Handle(blocker))
	}
	ref, _ := ix.lookup(blocker)
	if ref.Place != PlaceLogbook {
		t.Errorf("place = %q, want %q", ref.Place, PlaceLogbook)
	}
	if ref.Title != "the blocker" {
		t.Errorf("title = %q, want it read back out of the filed file", ref.Title)
	}
}

// The archive is a drawer, not a verdict: a ticket put there unfinished has
// not satisfied anything.
func TestArchivedUnfinishedTicketIsNotSatisfied(t *testing.T) {
	s := store(t)
	id := make1(t, s, "parked", "in-progress", nil, nil)
	if _, err := s.Archive(id); err != nil {
		t.Fatal(err)
	}

	ix := Build(s, lanes(t), list(t, s))
	if !ix.Known(id) {
		t.Error("an archived ticket still exists and must be known")
	}
	if ix.Satisfied(id) {
		t.Error("an archived ticket in a working lane is not finished work")
	}
}

// An id nothing here has ever seen is the one case that is genuinely
// dangling.
func TestUnknownIDIsNeitherKnownNorSatisfied(t *testing.T) {
	s := store(t)
	ix := Build(s, lanes(t), list(t, s))
	ghost := ticket.NewID(time.Now())
	if ix.Known(ghost) || ix.Satisfied(ghost) {
		t.Error("an id that exists nowhere must be neither known nor satisfied")
	}
}

// Children are derived from the parent field on the child, and the whole
// depth comes back in one read — that is what makes an epic an epic here.
func TestChildrenComeBackAsATree(t *testing.T) {
	s := store(t)
	root := make1(t, s, "epic", "todo", nil, nil)
	child := make1(t, s, "child", "todo", map[string]string{ticket.FieldParent: root}, nil)
	grand := make1(t, s, "grandchild", "todo", map[string]string{ticket.FieldParent: child}, nil)

	ix := Build(s, lanes(t), list(t, s))
	got := kinds(ix.Relations(root), KindChild)
	if len(got) != 2 || got[0] != child || got[1] != grand {
		t.Fatalf("children of the root = %v, want %v then %v", got, child, grand)
	}
	var depth int
	for _, e := range ix.Relations(root) {
		if e.Ref.ID == grand {
			depth = e.Ref.Depth
		}
	}
	if depth != 1 {
		t.Errorf("grandchild depth = %d, want 1", depth)
	}
	// And the child knows what it is part of.
	if got := kinds(ix.Relations(child), KindParent); len(got) != 1 || got[0] != root {
		t.Errorf("parent of the child = %v, want %v", got, root)
	}
}

// A child that has been finished and filed is still part of its parent. This
// is the case the whole ticket exists for: an epic whose children have
// shipped must not look empty.
func TestAFiledChildIsStillAChild(t *testing.T) {
	s := store(t)
	root := make1(t, s, "epic", "todo", nil, nil)
	child := make1(t, s, "shipped", "done", map[string]string{ticket.FieldParent: root}, nil)
	if _, err := s.Logbook(child, "as-20260911"); err != nil {
		t.Fatal(err)
	}

	ix := Build(s, lanes(t), list(t, s))
	got := kinds(ix.Relations(root), KindChild)
	if len(got) != 1 || got[0] != child {
		t.Fatalf("children = %v, want the logbooked %v", got, ticket.Handle(child))
	}
}

// related is stored on one side and read from both, because a link that had
// to be written twice would be half-written most of the time.
func TestRelatedIsReadFromBothSides(t *testing.T) {
	s := store(t)
	a := make1(t, s, "a", "todo", nil, nil)
	b := make1(t, s, "b", "todo", nil, map[string][]string{ticket.FieldRelated: {a}})

	ix := Build(s, lanes(t), list(t, s))
	if got := kinds(ix.Relations(a), KindRelated); len(got) != 1 || got[0] != b {
		t.Errorf("a's related = %v, want %v — the side that was not edited must see it too", got, b)
	}
	if got := kinds(ix.Relations(b), KindRelated); len(got) != 1 || got[0] != a {
		t.Errorf("b's related = %v, want %v", got, a)
	}
}

// blocked-by is shown from both ends too: what I wait on, and who waits on me.
func TestBlocksIsTheReverseOfBlockedBy(t *testing.T) {
	s := store(t)
	dep := make1(t, s, "dep", "todo", nil, nil)
	waiter := make1(t, s, "waiter", "todo", nil, map[string][]string{ticket.FieldBlockedBy: {dep}})

	ix := Build(s, lanes(t), list(t, s))
	if got := kinds(ix.Relations(waiter), KindBlockedBy); len(got) != 1 || got[0] != dep {
		t.Errorf("waiting on = %v, want %v", got, dep)
	}
	if got := kinds(ix.Relations(dep), KindBlocks); len(got) != 1 || got[0] != waiter {
		t.Errorf("blocking = %v, want %v", got, waiter)
	}
}

// A parent ring is something a hand-edited file can always produce, and the
// tree is walked recursively wherever it is rendered. Terminating is the
// difference between a wrong view and a hung one.
func TestParentCycleTerminates(t *testing.T) {
	s := store(t)
	a := make1(t, s, "a", "todo", nil, nil)
	b := make1(t, s, "b", "todo", map[string]string{ticket.FieldParent: a}, nil)
	if _, err := s.Mutate(a, func(tk *ticket.Ticket) error {
		return tk.Doc().SetScalar(ticket.FieldParent, b)
	}); err != nil {
		t.Fatal(err)
	}

	ix := Build(s, lanes(t), list(t, s))
	done := make(chan int, 1)
	go func() { done <- len(ix.Relations(a)) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Relations did not terminate on a parent ring")
	}
}

// A link to an id nothing knows is still shown. Hiding it is how a typo
// survives for months.
func TestDanglingLinkIsStillListed(t *testing.T) {
	s := store(t)
	ghost := ticket.NewID(time.Now())
	id := make1(t, s, "a", "todo", nil, map[string][]string{ticket.FieldBlockedBy: {ghost}})

	ix := Build(s, lanes(t), list(t, s))
	for _, e := range ix.Relations(id) {
		if e.Ref.ID == ghost {
			if e.Ref.Place != PlaceUnknown {
				t.Errorf("place = %q, want %q", e.Ref.Place, PlaceUnknown)
			}
			return
		}
	}
	t.Fatal("a dangling blocked-by must still be listed")
}

// A ticket travelling on its own ref, with no file here, is a fourth place a
// link can point at. The store already folds those into its listing; the
// index has to name the place so the reader is told to pull it rather than
// told it is missing.
func TestATicketOnlyOnARefIsNamedAsSuch(t *testing.T) {
	s := store(t)
	onRef := &ticket.Ticket{
		ID:       ticket.NewID(time.Now()),
		Title:    "only on a ref",
		Status:   "todo",
		ReadOnly: true,
	}
	ix := Build(s, lanes(t), []*ticket.Ticket{onRef})

	ref, ok := ix.lookup(onRef.ID)
	if !ok || ref.Place != PlaceRef {
		t.Fatalf("place = %q (found %v), want %q", ref.Place, ok, PlaceRef)
	}
	if ix.Satisfied(onRef.ID) {
		t.Error("a ticket in a working lane is not satisfied, ref or not")
	}
}

// Filing a ticket into the logbook does not delete its ref, so the same
// ticket arrives twice. The filed copy is this clone's own record that the
// work is finished; the ref would otherwise report finished work as "pull it
// to work on it".
func TestTheFiledCopyBeatsTheRefCopy(t *testing.T) {
	s := store(t)
	id := make1(t, s, "finished here", "done", nil, nil)
	if _, err := s.Logbook(id, "as-20260911"); err != nil {
		t.Fatal(err)
	}
	stale := &ticket.Ticket{ID: id, Title: "finished here", Status: "in-progress", ReadOnly: true}

	ix := Build(s, lanes(t), []*ticket.Ticket{stale})

	ref, _ := ix.lookup(id)
	if ref.Place != PlaceLogbook {
		t.Errorf("place = %q, want %q — the filed copy is the live one", ref.Place, PlaceLogbook)
	}
	if !ix.Satisfied(id) {
		t.Error("the work is filed as done; the stale ref must not keep it open")
	}
}

// The file is hand-editable, so a link may be written as the handle somebody
// read off the board. An unambiguous one resolves; an ambiguous one does not,
// because a guess would show a link to whichever was indexed first.
func TestAnUnambiguousHandleResolves(t *testing.T) {
	s := store(t)
	id := make1(t, s, "a", "todo", nil, nil)

	ix := Build(s, lanes(t), list(t, s))
	ref, ok := ix.lookup(ticket.Handle(id))
	if !ok || ref.ID != id {
		t.Errorf("handle %s resolved to %#v, want %s", ticket.Handle(id), ref, id)
	}
	if _, ok := ix.lookup("ZZZZZZ"); ok {
		t.Error("a handle nothing carries must resolve to nothing")
	}
}
