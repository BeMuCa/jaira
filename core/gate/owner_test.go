package gate

import "testing"

// The ownership rail compares strings, and a person is several strings. Without
// the aliases, a ticket assigned under a work email read as somebody else's
// while git reported a user.name, so every move on your own ticket needed
// --force. A gate that is always overridden protects nothing; it trains the
// override.
func TestOwnershipAcceptsTheActorsOtherNames(t *testing.T) {
	tk := ticketWith("")
	tk.Assignee = "berk.calabakan@partner.example"

	refused := func(req Request) bool {
		for _, v := range CheckAdvance(testEnv(t), tk, req) {
			if v.Code == CodeNotOwner {
				return true
			}
		}
		return false
	}

	if !refused(Request{To: "signoff", Actor: "BeMuCa"}) {
		t.Error("without aliases the rail did not fire at all, so this test proves nothing")
	}
	if refused(Request{
		To: "signoff", Actor: "BeMuCa",
		ActorAliases: []string{"berk.calabakan@partner.example"},
	}) {
		t.Error("an alias of the actor was still treated as somebody else")
	}
	// Case and surrounding space are not identity.
	if refused(Request{
		To: "signoff", Actor: "BeMuCa",
		ActorAliases: []string{"  BERK.CALABAKAN@PARTNER.EXAMPLE  "},
	}) {
		t.Error("the comparison is case- or space-sensitive")
	}
}

func TestOwnershipStillRefusesATeammate(t *testing.T) {
	tk := ticketWith("")
	tk.Assignee = "alexander@example.test"

	var found bool
	for _, v := range CheckAdvance(testEnv(t), tk, Request{
		To: "signoff", Actor: "BeMuCa",
		ActorAliases: []string{"berk@example.test", "berk.calabakan@partner.example"},
	}) {
		if v.Code == CodeNotOwner {
			found = true
		}
	}
	if !found {
		t.Error("aliases opened the rail to a ticket that is genuinely someone else's")
	}
}
