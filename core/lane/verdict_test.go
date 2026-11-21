package lane

import (
	"strings"
	"testing"
)

// A lane file that names verdict words wrongly is refused at load rather than
// read as some other rule: a word that both accepts and sends back, or words
// on a lane that writes no review-verdict.
func TestVerdictWordsParse(t *testing.T) {
	head := "---\nid: r\nagentic: true\nmodel-tier: strong\n"
	for _, c := range []struct {
		name, front, err string
		accept, back     []string
	}{
		{"lists", "output-produces: [review-verdict]\nverdict-accept: [Accept, ok]\nverdict-send-back: [send back]\n", "", []string{"accept", "ok"}, []string{"send-back"}},
		{"scalars", "output-produces: [review-verdict]\nverdict-accept: accept\nverdict-send-back: send-back\n", "", []string{"accept"}, []string{"send-back"}},
		{"none", "output-produces: [review-verdict]\n", "", nil, nil},
		{"both ways", "output-produces: [review-verdict]\nverdict-accept: [go]\nverdict-send-back: [Go]\n", "begins with", nil, nil},
		{"accept starts with send-back", "output-produces: [review-verdict]\nverdict-accept: [no blockers]\nverdict-send-back: [no]\n", "begins with", nil, nil},
		{"accept only shares letters", "output-produces: [review-verdict]\nverdict-accept: [nothing]\nverdict-send-back: [no]\n", "", []string{"nothing"}, []string{"no"}},
		{"accept alone", "output-produces: [review-verdict]\nverdict-accept: [accept]\n", "come together", nil, nil},
		{"send-back alone", "output-produces: [review-verdict]\nverdict-send-back: [send-back]\n", "come together", nil, nil},
		{"no verdict field", "output-produces: [review-summary]\nverdict-accept: [accept]\nverdict-send-back: [no]\n", "output-produces", nil, nil},
	} {
		l, err := Parse([]byte(head+c.front+"---\nprompt\n"), c.name)
		if c.err != "" {
			if err == nil || !strings.Contains(err.Error(), c.err) {
				t.Errorf("%s: err %v, want one naming %q", c.name, err, c.err)
			}
			continue
		}
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if strings.Join(l.VerdictAccept, ",") != strings.Join(c.accept, ",") || strings.Join(l.VerdictSendBack, ",") != strings.Join(c.back, ",") {
			t.Errorf("%s: read %v / %v, want %v / %v", c.name, l.VerdictAccept, l.VerdictSendBack, c.accept, c.back)
		}
	}
}
