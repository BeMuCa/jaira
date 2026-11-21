package lane

import (
	"strings"
	"testing"
)

// A lane file that names verdict words wrongly is refused at load rather than
// read as some other rule: a word that both accepts and sends back, or words
// on a lane that writes no review-verdict.
// A declared word and a verdict fold through the same normalisation, so a word
// that starts with markup or a digit reads as itself.
func TestDeclaredWordsMatchTheirVerdicts(t *testing.T) {
	l := &Lane{VerdictSendBack: []string{"reject", "1", "2nd-look"}}
	for _, v := range []string{"#reject: x", "-1: x", "2nd-look, DoD 2", "**Reject**"} {
		if _, back := l.SendsBackBy(v); !back {
			t.Errorf("%q: not read as a send-back", v)
		}
	}
	for _, v := range []string{"12: x", "2nd-looks"} {
		if w, back := l.SendsBackBy(v); back {
			t.Errorf("%q: read as send-back %q", v, w)
		}
	}
}

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
		{"words with markup", "output-produces: [review-verdict]\nverdict-accept: [\"**ok**\"]\nverdict-send-back: [\"#reject\", \"-1\", 2nd-look]\n", "", []string{"ok"}, []string{"reject", "1", "2nd-look"}},
		{"word without letters", "output-produces: [review-verdict]\nverdict-accept: [accept]\nverdict-send-back: [send-back, \"❌\"]\n", "no letter or digit", nil, nil},
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
