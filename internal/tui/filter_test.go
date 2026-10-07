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

	if matches(tk, "assignee:berk", nil) {
		t.Error("assignee:berk matched a ticket assigned to sam whose prose mentions berk")
	}
	if !matches(tk, "context:berk", nil) {
		t.Error("context:berk did not match the context that contains berk")
	}
	if !matches(tk, "lane:review", nil) || !matches(tk, "status:review", nil) {
		t.Error("lane:/status: did not match the ticket's lane")
	}
	if !matches(tk, "ticket:7bvrct", nil) {
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

	if !matches(tk, "http://example.test", nil) {
		t.Error("a query containing a colon with an unknown key did not full-text match")
	}
}

// A known key with an empty ticket field matches nothing rather than leaking
// back into full text.
func TestFilterKnownKeyOnEmptyFieldMatchesNothing(t *testing.T) {
	tk := &ticket.Ticket{
		ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "goal is mentioned here", Status: "todo",
	}

	if matches(tk, "goal:mentioned", nil) {
		t.Error("goal:<q> matched via full text although the goal field is empty")
	}
}

// Every word of a filter is a condition of its own and all of them must
// hold: "assignee:sam 7bvrct" is sam's ticket 7BVRCT, not a ticket whose
// assignee is called "sam 7bvrct".
func TestFilterTermsAllHaveToMatch(t *testing.T) {
	tk := &ticket.Ticket{ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "Fix session cookie", Status: "todo", Assignee: "sam"}

	if !matches(tk, "assignee:sam 7bvrct", nil) {
		t.Error("a field condition and an id did not match together")
	}
	if matches(tk, "assignee:sam zzzzzz", nil) {
		t.Error("a ticket matched although one of the two conditions fails")
	}
}

// A comma lists alternatives for one field.
func TestFilterCommaMeansOr(t *testing.T) {
	tk := &ticket.Ticket{ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "t", Status: "review", Assignee: "sam"}

	if !matches(tk, "assignee:berk,sam", nil) || !matches(tk, "lane:todo,review", nil) {
		t.Error("one of the listed values did not match")
	}
	if matches(tk, "assignee:berk,alex", nil) {
		t.Error("matched although none of the listed values is the assignee")
	}
}

// user: is whoever the ticket belongs to or was written by, by exact name.
func TestFilterUserIsAssigneeOrCreator(t *testing.T) {
	tk := &ticket.Ticket{ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "t", Status: "todo", Assignee: "sam", Creator: "Berk"}

	for _, q := range []string{"user:sam", "user:berk", "user:BERK", "user:alex,sam"} {
		if !matches(tk, q, nil) {
			t.Errorf("%q did not match a ticket assigned to sam and written by Berk", q)
		}
	}
	if matches(tk, "user:be", nil) {
		t.Error("user: matched part of a name; it names people, so it is exact")
	}
}

// Splitting a plain phrase into words only widens it: whatever held the
// phrase holds each word, so no ticket the old filter found goes missing. A
// field value with a space is not a plain phrase — it needs quotes now.
func TestFilterPhraseNeverFindsLess(t *testing.T) {
	tk := &ticket.Ticket{ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "Fix session cookie dropped", Status: "todo"}

	for _, q := range []string{"session cookie", "cookie session", "fix cookie dropped"} {
		if !matches(tk, q, nil) {
			t.Errorf("%q no longer finds a ticket whose title holds it", q)
		}
	}
}

// Double quotes keep words together: a name with a space in it is one
// person, and a quoted phrase is searched whole, the way every filter was
// before words became conditions of their own.
func TestFilterQuotesKeepWordsTogether(t *testing.T) {
	tk := &ticket.Ticket{ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "Fix session cookie", Status: "todo", Creator: "Alexander Sacharov"}

	if !matches(tk, `user:"alexander sacharov"`, nil) || !matches(tk, `user:bemuca,"Alexander Sacharov"`, nil) {
		t.Error("a quoted name with a space did not match its creator")
	}
	if !matches(tk, `"session cookie"`, nil) {
		t.Error("a quoted phrase did not match the title holding it")
	}
	if matches(tk, `"cookie session"`, nil) {
		t.Error("a quoted phrase matched its words in another order")
	}
}

// A key left waiting for its value takes the next word: "tag: ui" is one
// condition, as it was before words became conditions of their own.
func TestFilterSpaceAfterAColonStaysOneCondition(t *testing.T) {
	tk := &ticket.Ticket{ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "t", Status: "todo", Assignee: "sam", Tags: []string{"ui"}}
	other := &ticket.Ticket{ID: "01KZTT3XZ2YQBX93TTSR7BVRCX", Title: "mentions sam", Status: "todo", Assignee: "berk"}

	if !matches(tk, "tag: ui", nil) || !matches(tk, `user: "sam"`, nil) {
		t.Error("a value after a space did not stay with its key")
	}
	if matches(other, "assignee: sam", nil) {
		t.Error("assignee: sam fell apart into every ticket and the word sam")
	}
	// A word that only looks like a key is still prose, space and all.
	prose := &ticket.Ticket{ID: "01KZTT3XZ2YQBX93TTSR7BVRCY", Title: "fix: crash on start", Status: "todo"}
	if !matches(prose, "fix: crash", nil) {
		t.Error("fix: crash no longer finds the title that says it")
	}
}

// An empty alternative is not a match: "title:zzz," while the next name is
// being typed must not let every ticket through.
func TestFilterEmptyAlternativeAddsNothing(t *testing.T) {
	tk := &ticket.Ticket{ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "t", Status: "todo", Assignee: "berk"}

	if matches(tk, "title:zzz,", nil) || matches(tk, "user:sam,", nil) || matches(tk, "user:sam,,", nil) || matches(tk, `user:sam,""`, nil) {
		t.Error("an empty alternative matched")
	}
	if !matches(tk, "user:", nil) || !matches(tk, "title:,", nil) {
		t.Error("a key with no value yet must leave the board as it is")
	}
}

// Quotes keep a comma inside a value: "Doe, John" is one person.
func TestFilterQuotesKeepACommaInAValue(t *testing.T) {
	tk := &ticket.Ticket{ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Title: "then what", Status: "todo", Creator: "Doe, John"}

	if !matches(tk, `user:"Doe, John"`, nil) || !matches(tk, `user:sam,"doe, john"`, nil) {
		t.Error("a quoted name with a comma did not match its creator")
	}
	if matches(tk, `title:"zz, then"`, nil) {
		t.Error("a quoted value was split at its comma")
	}
}
