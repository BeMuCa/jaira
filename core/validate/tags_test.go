package validate

import (
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

// A tag that cannot be normalised is written on the ticket and reachable by
// nothing: 'jaira tags' skips it when counting and both filters compare
// normalised names, so it is a label only a reader of the raw file ever sees.
// Same shape as a dangling dependency — the value is there, the thing it points
// at is not.
func TestUnstorableTagIsReportedWithTheRename(t *testing.T) {
	bad := tk(ticket.NewID(time.Now()), "mangled tag", "todo")
	bad.Goal, bad.Context, bad.Assignee, bad.DoD = "g", "c", "berk", "d"
	bad.Tags = []string{"ui", "front/end"}

	ps := Tickets([]*ticket.Ticket{bad}, lanes(t))
	if !has(ps, CodeBadTag) {
		t.Fatalf("an unstorable tag produced %v", codes(ps))
	}
	var p Problem
	for _, cand := range ps {
		if cand.Code == CodeBadTag {
			p = cand
		}
	}
	if p.Severity != SeverityWarning {
		t.Errorf("severity = %q, want a warning: the ticket itself is intact", p.Severity)
	}
	if p.Field != ticket.FieldTags {
		t.Errorf("field = %q, want %q", p.Field, ticket.FieldTags)
	}
	if !strings.Contains(p.Message, "front/end") {
		t.Errorf("the message does not name the offending tag: %q", p.Message)
	}
	// The suggestion must be followable verbatim: 'jaira set' replaces a list
	// outright, so naming only the repaired tag would erase ui.
	if !strings.Contains(p.Message, "tags=ui,front-end") {
		t.Errorf("the message does not suggest the whole repaired list: %q", p.Message)
	}
}

// A case difference is not this problem: "UI" normalises to "ui", so the listing
// counts it and both filters match it. Warning about it would be noise about
// something that works.
func TestCaseOnlyTagIsNotReported(t *testing.T) {
	fine := tk(ticket.NewID(time.Now()), "shouty tag", "todo")
	fine.Goal, fine.Context, fine.Assignee, fine.DoD = "g", "c", "berk", "d"
	fine.Tags = []string{"UI", "My UI"}

	if ps := Tickets([]*ticket.Ticket{fine}, lanes(t)); has(ps, CodeBadTag) {
		t.Errorf("a tag that only needs case folding was reported: %v", codes(ps))
	}
}

// Nothing to say about a ticket with no tags.
func TestNoTagsIsNotAProblem(t *testing.T) {
	fine := tk(ticket.NewID(time.Now()), "untagged", "todo")
	fine.Goal, fine.Context, fine.Assignee, fine.DoD = "g", "c", "berk", "d"
	if ps := Tickets([]*ticket.Ticket{fine}, lanes(t)); has(ps, CodeBadTag) {
		t.Errorf("an untagged ticket was reported: %v", codes(ps))
	}
}
