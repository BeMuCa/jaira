package epic

import (
	"strings"
	"testing"
	"time"

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
	un := Unassigned([]*Epic{e}, ix, board)
	if len(un) != 1 || un[0].Title != "d" {
		t.Errorf("unassigned = %v; want only d", un)
	}
}
