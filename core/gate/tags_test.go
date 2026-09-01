package gate

import (
	"testing"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// A lane may declare it produces tags — a triage step that files a ticket under
// a subject. Without a case here, tags is a list the gate can never see
// satisfied: fieldFilled would read it as a scalar, find nothing, and the lane
// would owe a field that is already on the ticket.
func TestTagsIsFilledByHavingOne(t *testing.T) {
	if fieldFilled(&ticket.Ticket{ID: "x"}, ticket.FieldTags) {
		t.Error("a ticket with no tags reports the field filled")
	}
	if !fieldFilled(&ticket.Ticket{ID: "x", Tags: []string{"ui"}}, ticket.FieldTags) {
		t.Error("a ticket tagged ui reports the field empty")
	}
}

// The debt a declaring lane carries is the thing that would have gone wrong:
// a tagged ticket must owe nothing.
func TestOwedByClearsTagsOnceTheTicketHasOne(t *testing.T) {
	set := &lane.Set{Lanes: []*lane.Lane{
		{ID: "triage", OutputProduces: []string{ticket.FieldTags}},
	}}
	if l, ok := OwedBy(set, &ticket.Ticket{ID: "x"})[ticket.FieldTags]; !ok || l != "triage" {
		t.Errorf("an untagged ticket owes tags to %q,%v — want triage,true", l, ok)
	}
	tagged := &ticket.Ticket{ID: "x", Tags: []string{"ui"}}
	if l, ok := OwedBy(set, tagged)[ticket.FieldTags]; ok {
		t.Errorf("a tagged ticket still owes tags to %q", l)
	}
}
