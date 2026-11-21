package lane

import "strings"

// verdictWords normalises the words a lane file declares.
func verdictWords(ws []string) []string {
	var out []string
	for _, w := range ws {
		if n := normVerdict(w); n != "" {
			out = append(out, n)
		}
	}
	return out
}

// normVerdict folds case, treats a space like a dash and drops the full stop a
// sentence ends on, so "Send back", "send-back" and "Accept." are each one word.
func normVerdict(s string) string {
	s = strings.TrimRight(strings.TrimSpace(s), ".!;,")
	return strings.Join(strings.Fields(strings.ToLower(strings.ReplaceAll(s, "-", " "))), "-")
}

// DeclaresVerdicts reports whether the lane names the words its review-verdict
// takes. Without them the verdict is free text and nothing reads it.
func (l *Lane) DeclaresVerdicts() bool {
	return len(l.VerdictAccept)+len(l.VerdictSendBack) > 0
}

// SendsBackBy reports whether review-verdict v sends the work back on this
// lane, and the word it does it with. The word is everything on the first line
// before the first ':' — or that whole line when it has none — so
// "send-back: DoD 2 unmet" and "Send back" both read; a sentence that merely
// contains the word somewhere does not. Every other verdict, free text
// included, sends nothing back.
func (l *Lane) SendsBackBy(v string) (string, bool) {
	head, _, _ := strings.Cut(strings.TrimSpace(v), "\n")
	head, _, _ = strings.Cut(head, ":")
	n := normVerdict(head)
	for _, w := range l.VerdictSendBack {
		if n == w {
			return w, true
		}
	}
	return "", false
}
