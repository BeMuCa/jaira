package gate

import (
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// receiptTicket is a ticket read from its file, so it carries the document a
// receipt is read from; ticketWith builds in memory and has none.
func receiptTicket(t *testing.T, receipt string) *ticket.Ticket {
	t.Helper()
	src := "---\nid: 01KZTT3XZ2YQBX93TTSR7BVRCT\ntitle: t\nstatus: review\nassignee: a\n"
	if receipt != "" {
		src += "receipt-review: " + receipt + "\n"
	}
	src += "---\nbody\n"
	d, err := ticket.ParseDoc([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	tk, err := ticket.Decode(d, "ticket.md")
	if err != nil {
		t.Fatal(err)
	}
	return tk
}

// receiptEnv is the built-in board with the review lane made to require a
// receipt. The env is loaded fresh, so changing the lane leaks nowhere.
func receiptEnv(t *testing.T) Env {
	t.Helper()
	env := testEnv(t)
	l, ok := env.Lanes.Get("review")
	if !ok {
		t.Fatal("no review lane")
	}
	l.ReceiptRequired = true
	return env
}

// A lane that requires a receipt holds the ticket until the current round has
// one. Without the refusal the send the lane exists for is never recorded.
func TestReceiptRequiredRefusesWithoutAReceipt(t *testing.T) {
	env := receiptEnv(t)
	tk := receiptTicket(t, "")
	if !hasFieldViolation(CheckAdvance(env, tk, Request{To: "signoff"}), CodeNeedsReceipt, "receipt-review") {
		t.Error("a forward move without a receipt was not refused")
	}
}

// The receipt for this round's key lets the ticket leave.
func TestReceiptForThisRoundPasses(t *testing.T) {
	env := receiptEnv(t)
	tk := receiptTicket(t, ReceiptKey(&ticket.Ticket{ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Round: 1}, "review")+" sent-1")
	if hasFieldViolation(CheckAdvance(env, tk, Request{To: "signoff"}), CodeNeedsReceipt, "receipt-review") {
		t.Error("a receipt for the current round was refused")
	}
}

// A send-back starts a new round on purpose, so last round's receipt must not
// stop this round's send from being recorded.
func TestEarlierRoundReceiptIsRefused(t *testing.T) {
	env := receiptEnv(t)
	tk := receiptTicket(t, ReceiptKey(&ticket.Ticket{ID: "01KZTT3XZ2YQBX93TTSR7BVRCT", Round: 1}, "review")+" sent-1")
	tk.Round = 2
	if !hasFieldViolation(CheckAdvance(env, tk, Request{To: "signoff"}), CodeNeedsReceipt, "receipt-review") {
		t.Error("a round-1 receipt let a round-2 forward move through")
	}
}

// Sending a ticket back is not a move past the lane, so it owes no receipt.
func TestSendBackNeedsNoReceipt(t *testing.T) {
	env := receiptEnv(t)
	tk := receiptTicket(t, "")
	if hasFieldViolation(CheckAdvance(env, tk, Request{To: "in-progress"}), CodeNeedsReceipt, "receipt-review") {
		t.Error("a send-back was refused for a receipt")
	}
}
