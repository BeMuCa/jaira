package tui

import (
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// A "key:value" filter narrows to one field: "assignee:berk" must not also
// match every ticket whose prose mentions berk.
func TestFilterKeyNarrowsToTheField(t *testing.T) {
	tk := &ticket.Ticket{
		ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "t", Status: "review",
		Assignee: "sam",
		Context:  "berk reported this while debugging",
	}

	if matches(tk, "assignee:berk") {
		t.Error("assignee:berk matched a ticket assigned to sam whose prose mentions berk")
	}
	if !matches(tk, "context:berk") {
		t.Error("context:berk did not match the context that contains berk")
	}
	if !matches(tk, "lane:review") || !matches(tk, "status:review") {
		t.Error("lane:/status: did not match the ticket's lane")
	}
	if !matches(tk, "ticket:7bvrct") {
		t.Error("ticket:<id suffix> did not match the id")
	}
}

// An unknown key is a search term, not a field — "http:" in a pasted URL must
// fall through to full text instead of matching nothing.
func TestFilterUnknownKeyFallsBackToFullText(t *testing.T) {
	tk := &ticket.Ticket{
		ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "t", Status: "todo",
		Context: "see http://example.test/page",
	}

	if !matches(tk, "http://example.test") {
		t.Error("a query containing a colon with an unknown key did not full-text match")
	}
}

// A known key with an empty ticket field matches nothing rather than leaking
// back into full text.
func TestFilterKnownKeyOnEmptyFieldMatchesNothing(t *testing.T) {
	tk := &ticket.Ticket{
		ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "goal is mentioned here", Status: "todo",
	}

	if matches(tk, "goal:mentioned") {
		t.Error("goal:<q> matched via full text although the goal field is empty")
	}
}
