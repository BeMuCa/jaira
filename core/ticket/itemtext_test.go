package ticket

import (
	"strings"
	"testing"
)

const rewordBody = `## Definition of Done

- [x] 429 returned above 100 req/min
  proof: internal/http/limit.go:44, covered by TestLimit
- [ ] documented in README
`

// Rewording must not disturb the tick or the evidence beneath it: the reason
// this exists is that the old workaround — tick the stale item, append a new one
// — destroyed exactly that pairing.
func TestSetItemTextKeepsStateAndProof(t *testing.T) {
	out, err := SetItemText(rewordBody, SectionDoD, 0, "429 returned above the configured rate")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "- [x] 429 returned above the configured rate\n") {
		t.Errorf("item was not reworded, or its marker changed:\n%s", out)
	}
	if !strings.Contains(out, "  proof: internal/http/limit.go:44, covered by TestLimit") {
		t.Errorf("the proof line did not survive:\n%s", out)
	}
	items := ParseDoDItems(out)
	if len(items) != 2 {
		t.Fatalf("parsed %d items, want 2 — rewording changed the shape of the list", len(items))
	}
	if items[0].State != StateDone || items[0].Proof == "" {
		t.Errorf("item 1 = %+v, want its state and proof intact", items[0])
	}
	if items[1].Text != "documented in README" {
		t.Errorf("item 2 was touched: %q", items[1].Text)
	}
}

// Everything but the words stays put — indentation, bullet, marker.
func TestSetItemTextLeavesTheRestOfTheLine(t *testing.T) {
	body := "## Plan\n\n  * [~] old wording\n"
	out, err := SetItemText(body, SectionPlan, 0, "new wording")
	if err != nil {
		t.Fatal(err)
	}
	if out != "## Plan\n\n  * [~] new wording\n" {
		t.Errorf("got %q", out)
	}
}

func TestSetItemTextPreservesCRLF(t *testing.T) {
	body := "## Definition of Done\r\n\r\n- [ ] first\r\n- [ ] second\r\n"
	out, err := SetItemText(body, SectionDoD, 0, "reworded")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "- [ ] reworded\r\n") {
		t.Errorf("CRLF line ending was lost: %q", out)
	}
	if strings.Count(out, "\r\n") != strings.Count(body, "\r\n") {
		t.Error("line ending count changed")
	}
}

// A newline in the new text would re-parse as a line of its own, and one
// starting with a bullet would come back as another criterion — the list must
// not be able to grow through a reword.
func TestSetItemTextCollapsesWhitespace(t *testing.T) {
	out, err := SetItemText(rewordBody, SectionDoD, 1, "documented\n- [ ] and smuggled in")
	if err != nil {
		t.Fatal(err)
	}
	if n := len(ParseDoDItems(out)); n != 2 {
		t.Errorf("the checklist grew to %d items through a reword:\n%s", n, out)
	}
}

func TestSetItemTextRefusesEmptyText(t *testing.T) {
	if _, err := SetItemText(rewordBody, SectionDoD, 0, "   "); err == nil {
		t.Fatal("expected an error: an empty item says nothing")
	}
}

func TestSetItemTextOutOfRange(t *testing.T) {
	if _, err := SetItemText(rewordBody, SectionDoD, 9, "x"); err == nil {
		t.Fatal("expected an out-of-range error")
	}
	if _, err := SetItemText("# nothing here\n", SectionDoD, 0, "x"); err == nil {
		t.Fatal("expected an error when the section does not exist")
	}
}
