package epic

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/link"
	"github.com/BeMuCa/jaira/core/ticket"
)

const (
	idA = "01AAAAAAAAAAAAAAAAAAAAAAAA"
	idB = "01BBBBBBBBBBBBBBBBBBBBBBBB"
	idC = "01CCCCCCCCCCCCCCCCCCCCCCCC"
)

// A hand edit survives jaira's edits: the file is the API, so a comment and an
// order somebody chose must come back out of Add and Remove untouched.
func TestEditsKeepEveryOtherLine(t *testing.T) {
	src := "---\ntitle: Login\ndescription: no password\n---\n\nwhy we do this\n\n- " + idA + " the form\n<!-- keep me -->\n- " + idB + "\n"
	e := FromBytes("login", []byte(src))
	if e.Title != "Login" || e.Description != "no password" {
		t.Fatalf("facts = %q / %q", e.Title, e.Description)
	}
	if got := e.Members(); len(got) != 2 || got[0] != idA || got[1] != idB {
		t.Fatalf("members = %v", got)
	}
	e.Remove(idB)
	e.Add(idC)
	want := "---\ntitle: Login\ndescription: no password\n---\n\nwhy we do this\n\n- " + idA + " the form\n<!-- keep me -->\n- " + idC + "\n"
	if got := string(e.Bytes()); got != want {
		t.Errorf("file after edits:\n%s\nwant:\n%s", got, want)
	}
}

// The title is a person's, the name the ref's: an epic created without a title
// reads its name rather than an empty heading.
func TestNewDefaultsTheTitleToTheName(t *testing.T) {
	e := FromBytes("login", New("login", "", "line one\nline two", time.Unix(0, 0)).Bytes())
	if e.Title != "login" || e.Description != "line one line two" {
		t.Errorf("title %q, description %q", e.Title, e.Description)
	}
}

// Both sides of a lost race keep what they did: theirs added B, ours added C
// and took A out and retitled it.
func TestMergeKeepsBothSidesMemberChanges(t *testing.T) {
	base := New("f", "Old", "", time.Unix(0, 0))
	base.Add(idA)
	theirs := FromBytes("f", base.Bytes())
	theirs.Add(idB)
	ours := FromBytes("f", base.Bytes())
	ours.Add(idC)
	ours.Remove(idA)
	ours.setField("title", "New")

	got := FromBytes("f", Merge("f", base.Bytes(), ours.Bytes(), theirs.Bytes()))
	if m := got.Members(); len(m) != 2 || m[0] != idB || m[1] != idC {
		t.Errorf("members = %v; want [B C]", m)
	}
	if got.Title != "New" {
		t.Errorf("title = %q; our retitle was lost", got.Title)
	}
	if strings.Count(string(got.Bytes()), "title:") != 1 {
		t.Errorf("title written twice:\n%s", got.Bytes())
	}
}

// Follow-ups count through follows: and parent:, any number of steps down,
// and a ticket that follows on from a member is not unassigned.
func TestCoveredWalksFollowsAndParent(t *testing.T) {
	board := []*ticket.Ticket{
		{ID: idA, Title: "a"},
		{ID: idB, Title: "b", Follows: idA},
		{ID: idC, Title: "c", Parent: idB},
		{ID: "01DDDDDDDDDDDDDDDDDDDDDDDD", Title: "d"},
	}
	ix := link.Build(nil, nil, board)
	if got := Covered([]string{idA}, ix); len(got) != 3 || got[0] != idA {
		t.Fatalf("covered = %v; want A, B, C", got)
	}
	e := New("f", "", "", time.Unix(0, 0))
	e.Add(idA)
	p := Count(e, ix)
	if p.Total != 3 || p.Done != 0 || !p.Tickets[0].Direct || p.Tickets[1].Direct {
		t.Errorf("progress = %+v", p)
	}
	// Both follow-ups name the member line, even C, which is two steps down:
	// that line is what would have to go to take them out.
	for _, tk := range p.Tickets {
		if want := map[bool]string{true: "", false: idA}[tk.Direct]; tk.Via != want {
			t.Errorf("%s via = %q; want %q", tk.ID, tk.Via, want)
		}
	}
	un := Unassigned([]*Epic{e}, ix, board)
	if len(un) != 1 || un[0].Title != "d" {
		t.Errorf("unassigned = %v; want only d", un)
	}
}

// mapRefs keeps epics in memory, the way the outbox and the ref keep them.
type mapRefs map[string][]byte

func (m mapRefs) EpicNames() ([]string, error) {
	var out []string
	for name := range m {
		out = append(out, name)
	}
	return out, nil
}

func (m mapRefs) EpicContent(name string) ([]byte, error) {
	if c, ok := m[name]; ok {
		return c, nil
	}
	return nil, ErrNotFound
}

func (m mapRefs) RecordEpic(name string, content []byte) error {
	m[name] = content
	return nil
}

func ptr(s string) *string { return &s }

// Rewriting the text leaves the member lines as they were and in their order,
// a hand comment among them included; only the prose around them changes.
func TestSetKeepsTheMembers(t *testing.T) {
	src := "---\ntitle: Login\ndescription: no password\n---\n\n# Login\n\nold why\n\n- " + idB + " the form\n- " + idA + "\n"
	refs := mapRefs{"login": []byte(src)}
	e, changed, err := Set(&ticket.Store{Root: t.TempDir()}, refs, "login", Patch{Body: ptr("# Login\n\n## Goal\n\nnew why\r\nsecond line")})
	if err != nil || !changed {
		t.Fatalf("set: changed %v, %v", changed, err)
	}
	want := "---\ntitle: Login\ndescription: no password\n---\n\n# Login\n\n## Goal\n\nnew why\nsecond line\n\n- " + idB + " the form\n- " + idA + "\n"
	if got := string(refs["login"]); got != want {
		t.Errorf("file:\n%s\nwant:\n%s", got, want)
	}
	if m := e.Members(); len(m) != 2 || m[0] != idB || m[1] != idA {
		t.Errorf("members = %v", m)
	}
	if got := FromBytes("login", refs["login"]).Body(); got != "# Login\n\n## Goal\n\nnew why\nsecond line" {
		t.Errorf("body read back = %q", got)
	}
}

// A new title renames the heading New wrote, and a description is kept on one
// line; asking for what is there already writes nothing.
func TestSetRetitlesAndWritesOnlyWhatChanged(t *testing.T) {
	refs := mapRefs{}
	s := &ticket.Store{Root: t.TempDir()}
	if _, err := Create(s, refs, "f", "Old", "", []string{idA}); err != nil {
		t.Fatal(err)
	}
	e, _, err := Set(s, refs, "f", Patch{Title: ptr("New"), Description: ptr("one\nline")})
	if err != nil {
		t.Fatal(err)
	}
	if e.Title != "New" || e.Description != "one line" || !strings.HasPrefix(e.Body(), "# New\n") || !e.Has(idA) {
		t.Errorf("after set: %q %q %q %v", e.Title, e.Description, e.Body(), e.Members())
	}
	// The web form sends the text back as it read it beside a new title.
	e, _, err = Set(s, refs, "f", Patch{Title: ptr("Newer"), Body: ptr(e.Body() + "\r\n")})
	if err != nil || !strings.HasPrefix(e.Body(), "# Newer\n") {
		t.Fatalf("title beside the unchanged text: %q, %v", e.Body(), err)
	}
	before := string(refs["f"])
	if _, changed, err := Set(s, refs, "f", Patch{Title: ptr("Newer"), Body: ptr(e.Body())}); err != nil || changed {
		t.Errorf("same values: changed %v, %v", changed, err)
	}
	if string(refs["f"]) != before {
		t.Errorf("rewritten without a change")
	}
	if _, _, err := Set(s, refs, "missing", Patch{Title: ptr("x")}); !errors.Is(err, ErrNotFound) {
		t.Errorf("unknown epic: %v", err)
	}
}

// An editor that read the epic before somebody else saved it cannot overwrite
// them: the field it changes is no longer what it read, and nothing is written.
// A field it leaves alone may have changed freely.
func TestSetRefusesAStaleBase(t *testing.T) {
	refs := mapRefs{}
	s := &ticket.Store{Root: t.TempDir()}
	e, _ := Create(s, refs, "f", "T", "d", nil)
	read := e.Fields()
	if _, _, err := Set(s, refs, "f", Patch{Body: ptr("theirs")}); err != nil {
		t.Fatal(err)
	}
	if _, _, err := Set(s, refs, "f", Patch{Body: ptr("mine"), Base: &read}); !errors.Is(err, ErrChanged) {
		t.Fatalf("stale body: %v", err)
	}
	if got := FromBytes("f", refs["f"]).Body(); got != "theirs" {
		t.Fatalf("stale body was written: %q", got)
	}
	if _, changed, err := Set(s, refs, "f", Patch{Title: ptr("T2"), Base: &read}); err != nil || !changed {
		t.Errorf("title over a changed body: changed %v, %v", changed, err)
	}
}

// Two clones rewrite the same text and the same title at once: the one that
// lost the race keeps its own, and the other's is under the text in a marked
// block — nothing either of them wrote is gone, and the members still merge.
func TestMergeKeepsBothSidesOfOneField(t *testing.T) {
	base := New("f", "Old", "why", time.Unix(0, 0))
	base.Add(idA)
	theirs := FromBytes("f", base.Bytes())
	theirs.Add(idB)
	theirs.SetBody("# Old\n\ntheir paragraph")
	theirs.setField("title", "Theirs")
	ours := FromBytes("f", base.Bytes())
	ours.SetBody("# Old\n\nour paragraph")
	ours.setField("title", "Ours")

	got := FromBytes("f", Merge("f", base.Bytes(), ours.Bytes(), theirs.Bytes()))
	body := got.Body()
	if got.Title != "Ours" || !strings.HasPrefix(body, "# Old\n\nour paragraph\n\n") {
		t.Errorf("ours lost: title %q, body:\n%s", got.Title, body)
	}
	for _, want := range []string{"their paragraph", "Their title was:**\n\nTheirs\n", "Their text was:", "<!-- jaira: changed at the same time -->"} {
		if !strings.Contains(body, want) {
			t.Errorf("body lacks %q:\n%s", want, body)
		}
	}
	if m := got.Members(); len(m) != 2 || got.Description != "why" {
		t.Errorf("members %v, description %q", m, got.Description)
	}
}

// A new title on one clone renames the heading; a new paragraph on the other
// is a different edit. Both arrive, and nothing is marked as a clash.
func TestMergeTellsARetitleFromATextEdit(t *testing.T) {
	base := New("f", "Old", "", time.Unix(0, 0))
	theirs := FromBytes("f", base.Bytes())
	theirs.SetBody(theirs.Body() + "\n\nnew paragraph")
	refs := mapRefs{"f": base.Bytes()}
	ours, _, err := Set(&ticket.Store{Root: t.TempDir()}, refs, "f", Patch{Title: ptr("New")})
	if err != nil {
		t.Fatal(err)
	}
	got := FromBytes("f", Merge("f", base.Bytes(), ours.Bytes(), theirs.Bytes()))
	body := got.Body()
	if got.Title != "New" || !strings.HasPrefix(body, "# New\n") || !strings.Contains(body, "new paragraph") || strings.Contains(body, "jaira: changed") {
		t.Errorf("title %q, body:\n%s", got.Title, body)
	}
}

// A title written by hand with a double space still owns the heading New
// wrote for it, for a rename and for a merge alike.
func TestTheTitleHeadingIsFoundOneWay(t *testing.T) {
	e := FromBytes("f", []byte("---\ntitle: A  B\n---\n\n# A B\n\ntext\n"))
	if !strings.HasPrefix(e.plain(), titleMark) {
		t.Errorf("plain = %q", e.plain())
	}
	e.retitle("C")
	if got := e.Body(); got != "# C\n\ntext" {
		t.Errorf("body after retitle = %q", got)
	}
}

// What a person reads beside the title leaves out the heading that repeats it
// and the comments New writes for whoever edits the file; an epic nobody wrote
// a text for reads as having none.
func TestTextIsTheBodyAPersonReads(t *testing.T) {
	e := New("f", "Login", "", time.Unix(0, 0))
	if got := e.Text(); got != "" {
		t.Errorf("text of a new epic = %q; want none", got)
	}
	e.SetBody(e.Body() + "\n\n## Goal\n\nno password")
	if got := e.Text(); got != "## Goal\n\nno password" {
		t.Errorf("text = %q", got)
	}
	refs := mapRefs{"f": e.Bytes()}
	if _, _, err := Set(&ticket.Store{Root: t.TempDir()}, refs, "f", Patch{Title: ptr(" ")}); !errors.Is(err, ErrNoTitle) {
		t.Errorf("empty title: %v", err)
	}
}

// A list item in the text that starts with a full ticket id stays text: it
// does not make that ticket a member, and it does not drop out of the text.
func TestATicketIDInTheTextStaysText(t *testing.T) {
	refs := mapRefs{}
	s := &ticket.Store{Root: t.TempDir()}
	if _, err := Create(s, refs, "f", "F", "", []string{idA}); err != nil {
		t.Fatal(err)
	}
	if _, _, err := Set(s, refs, "f", Patch{Body: ptr("see\n\n- " + idB + " the old form\n  - " + idC)}); err != nil {
		t.Fatal(err)
	}
	e := FromBytes("f", refs["f"])
	if m := e.Members(); len(m) != 1 || m[0] != idA {
		t.Errorf("members = %v; want only A", m)
	}
	if got := e.Body(); got != "see\n\n* "+idB+" the old form\n  * "+idC {
		t.Errorf("body = %q", got)
	}
}

// Closing writes "status: closed" into the ref and reopening takes the line
// out again, so a reopened epic reads like one that was never closed. Neither
// deletes anything, and doing either twice writes nothing the second time.
func TestCloseAndReopenWriteOneLine(t *testing.T) {
	s := &ticket.Store{Root: t.TempDir()}
	e := New("f", "F", "", time.Unix(0, 0))
	e.Add(idA)
	open := string(e.Bytes())
	refs := mapRefs{"f": e.Bytes()}
	ix := index(t, &ticket.Ticket{ID: idA, Title: "a", Status: "done"})

	got, written, err := Close(s, refs, "f", ix, false)
	if err != nil || !written || !got.Closed {
		t.Fatalf("close: %v, written %v, closed %v", err, written, got != nil && got.Closed)
	}
	if !strings.Contains(string(refs["f"]), "\nstatus: closed\n") || !FromBytes("f", refs["f"]).Closed {
		t.Fatalf("ref after close:\n%s", refs["f"])
	}
	if m := FromBytes("f", refs["f"]).Members(); len(m) != 1 || m[0] != idA {
		t.Errorf("members after close = %v", m)
	}
	if _, written, _ := Close(s, refs, "f", ix, false); written {
		t.Error("closing a closed epic wrote again")
	}
	if _, written, err := Reopen(s, refs, "f"); err != nil || !written {
		t.Fatalf("reopen: %v, written %v", err, written)
	}
	if got := string(refs["f"]); got != open {
		t.Errorf("reopened epic:\n%s\nwant it as it was:\n%s", got, open)
	}
	if _, written, _ := Reopen(s, refs, "f"); written {
		t.Error("reopening an open epic wrote")
	}
}

// index builds an index over the shipped lanes, in which a ticket in "done"
// is finished.
func index(t *testing.T, board ...*ticket.Ticket) *link.Index {
	t.Helper()
	set, err := lane.Load("")
	if err != nil {
		t.Fatal(err)
	}
	return link.Build(nil, set, board)
}

// An epic still counting an unfinished ticket — a follow-up included — is
// not closed without force; with it, it is.
func TestCloseRefusesUnfinishedWithoutForce(t *testing.T) {
	s := &ticket.Store{Root: t.TempDir()}
	e := New("f", "", "", time.Unix(0, 0))
	e.Add(idA)
	refs := mapRefs{"f": e.Bytes()}
	ix := index(t, &ticket.Ticket{ID: idA, Title: "a", Status: "done"}, &ticket.Ticket{ID: idB, Title: "b", Status: "todo", Follows: idA})
	before := string(refs["f"])
	if _, _, err := Close(s, refs, "f", ix, false); !errors.Is(err, ErrUnfinished) || !strings.Contains(err.Error(), "1 of 2") {
		t.Fatalf("close with an open follow-up: %v; want ErrUnfinished naming 1 of 2", err)
	}
	if string(refs["f"]) != before {
		t.Error("a refused close wrote the ref")
	}
	if _, written, err := Close(s, refs, "f", ix, true); err != nil || !written {
		t.Fatalf("forced close: %v, written %v", err, written)
	}
}

// A closed epic takes no new ticket without force, and with force takes it
// and stays closed. Re-adding a member and taking one out are not refused.
func TestAssignToAClosedEpic(t *testing.T) {
	s := &ticket.Store{Root: t.TempDir()}
	e := New("f", "", "", time.Unix(0, 0))
	e.Add(idA)
	e.setClosed(true)
	refs := mapRefs{"f": e.Bytes()}
	before := string(refs["f"])
	if _, _, err := Assign(s, refs, "f", []string{idA, idB}, true, false); !errors.Is(err, ErrClosed) {
		t.Fatalf("add to closed: %v; want ErrClosed", err)
	}
	if string(refs["f"]) != before {
		t.Error("a refused add wrote the ref")
	}
	if _, changed, err := Assign(s, refs, "f", []string{idA}, true, false); err != nil || len(changed) != 0 {
		t.Errorf("re-adding a member: %v, changed %v; want nothing to do", err, changed)
	}
	got, changed, err := Assign(s, refs, "f", []string{idB}, true, true)
	if err != nil || len(changed) != 1 || !got.Closed || !FromBytes("f", refs["f"]).Closed {
		t.Fatalf("forced add: %v, changed %v; want B in and the epic still closed", err, changed)
	}
	if _, changed, err := Assign(s, refs, "f", []string{idB}, false, false); err != nil || len(changed) != 1 {
		t.Errorf("rm from closed: %v, changed %v", err, changed)
	}
}

// A close on one clone survives a race with a member added on another, and a
// reopen on ours is not undone by theirs having left the status alone.
func TestMergeKeepsAClosing(t *testing.T) {
	base := New("f", "F", "", time.Unix(0, 0))
	theirs := FromBytes("f", base.Bytes())
	theirs.Add(idB)
	ours := FromBytes("f", base.Bytes())
	ours.setClosed(true)
	got := FromBytes("f", Merge("f", base.Bytes(), ours.Bytes(), theirs.Bytes()))
	if !got.Closed || !got.Has(idB) {
		t.Errorf("merged: closed %v, has B %v; want both", got.Closed, got.Has(idB))
	}
	got = FromBytes("f", Merge("f", ours.Bytes(), base.Bytes(), ours.Bytes()))
	if got.Closed {
		t.Error("our reopen was lost")
	}
	got = FromBytes("f", Merge("f", base.Bytes(), theirs.Bytes(), ours.Bytes()))
	if !got.Closed {
		t.Error("their closing was lost when ours left the status alone")
	}
}
