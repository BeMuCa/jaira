package gate

import (
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// testingTicket decodes a ticket sitting in testing from real frontmatter, so
// hunt-verdict — a field this package does not model — is read the way the
// gate reads it on a board: through the document, not a struct field.
func testingTicket(t *testing.T, extra string) *ticket.Ticket {
	t.Helper()
	src := "---\nid: 01KZTT3XZ2YQBX93TTSR7BVRCT\ntitle: t\nstatus: testing\ngoal: g\ncontext: c\nassignee: a\n" +
		"outcome-what: w\noutcome-why: y\noutcome-resolves: r\n" + extra + "---\n"
	d, err := ticket.ParseDoc([]byte(src))
	if err != nil {
		t.Fatalf("ParseDoc: %v", err)
	}
	tk, err := ticket.Decode(d, "")
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	return tk
}

// The testing lane runs the bug hunt beside its tests and owes its result as
// hunt-verdict. A lane that ran the tests and never waited for the hunter has
// written test-verdict alone, and that ticket must not leave testing.
func TestTestingOwesHuntVerdict(t *testing.T) {
	env := testEnv(t)
	l, ok := env.Lanes.Get("testing")
	if !ok {
		t.Fatal("no testing lane")
	}

	tk := testingTicket(t, "test-verdict: \"pass: suite green\"\n")
	owed := OutputOwed(l, tk)
	if len(owed) != 1 || owed[0] != "hunt-verdict" {
		t.Fatalf("testing with only test-verdict owes %v, want [hunt-verdict]", owed)
	}
	if !hasFieldViolation(CheckAdvance(env, tk, Request{To: "human"}), CodeMissingProduces, "hunt-verdict") {
		t.Error("the gate let a ticket leave testing without hunt-verdict")
	}

	tk = testingTicket(t, "test-verdict: \"pass: suite green\"\nhunt-verdict: \"pass: bughunt done at 3f2a91c, 0 change\"\n")
	if owed := OutputOwed(l, tk); len(owed) != 0 {
		t.Errorf("testing with both verdicts still owes %v", owed)
	}
}
