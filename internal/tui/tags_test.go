package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

func taggedTicket() *ticket.Ticket {
	return &ticket.Ticket{
		ID:        "01KZTT3XZ2YQBX93TTSR7BVRCT",
		Title:     "Rate limit the login endpoint",
		Status:    "review",
		Assignee:  "berk",
		Creator:   "berk",
		Goal:      "stop credential stuffing",
		Context:   "came up while reading the auth logs",
		Tags:      []string{"backend", "security"},
		CreatedAt: time.Now().Add(-48 * time.Hour),
		UpdatedAt: time.Now().Add(-time.Hour),
	}
}

// A tag is read as often as who owns a ticket, so it is a base row: shown
// whenever the ticket carries one, in the same dim-label column as the rest.
func TestDetailShowsTheTagsRow(t *testing.T) {
	m := newTestModel(t, 120, 40)
	out := stripANSI(m.detailBody(taggedTicket(), 120))

	if !strings.Contains(out, "tags         backend security") {
		t.Errorf("no tags row:\n%s", out)
	}
	// Among the base rows, before the prose ones: it says which subject the
	// ticket belongs to, which frames everything below it.
	if rowIndex(t, out, "tags") > rowIndex(t, out, "context") {
		t.Errorf("the tags row is not a base row:\n%s", out)
	}
	if rowIndex(t, out, "tags") < rowIndex(t, out, "lane") {
		t.Errorf("the tags row jumped ahead of the lane row:\n%s", out)
	}
}

// An empty field must not leave an empty row behind: tags are optional, and a
// blank "tags" label says nothing at all.
func TestDetailOmitsTheTagsRowWhenThereAreNone(t *testing.T) {
	m := newTestModel(t, 120, 40)
	tk := taggedTicket()
	tk.Tags = nil
	if out := stripANSI(m.detailBody(tk, 120)); strings.Contains(out, "\ntags ") {
		t.Errorf("an untagged ticket still shows a tags row:\n%s", out)
	}
}

// The sign-off screen is where a person accepts work. Which subject it belongs
// to is part of that judgement, so the row appears there too.
func TestSignOffShowsTheTagsRow(t *testing.T) {
	m := newTestModel(t, 120, 40)
	m.detail = taggedTicket()
	out := stripANSI(m.renderSignOff())

	if !strings.Contains(out, "tags         backend security") {
		t.Errorf("no tags row on the sign-off screen:\n%s", out)
	}
	if rowIndex(t, out, "tags") > rowIndex(t, out, "problem") {
		t.Errorf("tags is not shown before the seven questions:\n%s", out)
	}
}

// The point of tagging: pull every ui ticket out of a backlog without reading
// it. The board filter learns tag: beside assignee:, lane: and the rest.
func TestFilterTagHitsAndMisses(t *testing.T) {
	tk := taggedTicket()

	for _, q := range []string{"tag:backend", "tag:security", "tags:backend", "tag:BACKEND"} {
		if !matches(tk, q) {
			t.Errorf("%q did not match a ticket tagged %v", q, tk.Tags)
		}
	}
	for _, q := range []string{"tag:ui", "tag:frontend"} {
		if matches(tk, q) {
			t.Errorf("%q matched a ticket tagged %v", q, tk.Tags)
		}
	}
	// Exact, unlike every other filter key, and deliberately: a tag is a name
	// from a closed vocabulary. "tag:cur" answering with everything tagged
	// "security" is a wrong answer rather than a loose one, and it would have
	// made the board filter disagree with 'jaira list --tag', which is exact.
	for _, q := range []string{"tag:cur", "tag:sec", "tag:security-review", "tag:back"} {
		if matches(tk, q) {
			t.Errorf("%q matched %v: the tag filter is substring, not exact", q, tk.Tags)
		}
	}
	// A hand-edited name still answers to its stored form, since both sides are
	// normalised before comparing.
	shouty := taggedTicket()
	shouty.Tags = []string{"My UI"}
	if !matches(shouty, "tag:my-ui") {
		t.Error("tag:my-ui did not match a ticket carrying \"My UI\"")
	}
	// A known key on an untagged ticket matches nothing rather than leaking
	// back into full text — the rule every other key here follows.
	untagged := taggedTicket()
	untagged.Tags = nil
	untagged.Context = "the backend was fine"
	if matches(untagged, "tag:backend") {
		t.Error("tag:backend matched via full text although the ticket has no tags")
	}
	// A plain search still finds a tag: it is part of what the ticket says.
	if !matches(tk, "security") {
		t.Error("a plain query did not match a tag")
	}
}
