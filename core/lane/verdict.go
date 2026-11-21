package lane

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

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
// lane, and the word it does it with. The verdict's first line, before the
// first ':', must begin with the word, followed by its end or by anything that
// is not a letter — so "send-back: DoD 2 unmet", "Send back", "Send back, DoD
// 1" and "send-back — DoD 1" all read, while "send-backs" and a sentence that
// merely contains the word somewhere do not. Every other verdict, free text
// included, sends nothing back.
func (l *Lane) SendsBackBy(v string) (string, bool) {
	head, _, _ := strings.Cut(strings.TrimSpace(v), "\n")
	head, _, _ = strings.Cut(head, ":")
	n := normVerdict(head)
	for _, w := range l.VerdictSendBack {
		if beginsWithWord(n, w) {
			return w, true
		}
	}
	return "", false
}

// beginsWithWord reports whether the normalised verdict n begins with the word
// w, followed by its end or by anything that is not a letter.
func beginsWithWord(n, w string) bool {
	rest, ok := strings.CutPrefix(n, w)
	if !ok {
		return false
	}
	r, _ := utf8.DecodeRuneInString(rest)
	return rest == "" || !unicode.IsLetter(r)
}
