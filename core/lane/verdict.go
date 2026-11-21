package lane

import "strings"

// The three things a review-verdict can mean on a lane that declares its
// words. VerdictNone is an empty field — the lane's output contract already
// refuses that — and VerdictUnknown is a value that starts with neither kind
// of word: no verdict at all.
const (
	VerdictNone     = ""
	VerdictAccept   = "accept"
	VerdictSendBack = "send-back"
	VerdictUnknown  = "unknown"
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

// ReadVerdict classifies a review-verdict against the lane's words. The word is
// everything on the first line before the first ':' — or that whole line when it has none —
// so "send-back: DoD 2 unmet", "Send back" and "accept" all read; a sentence
// that merely contains a word somewhere does not. Returns the word as written
// in the lane file beside its kind.
func (l *Lane) ReadVerdict(v string) (kind, word string) {
	v = strings.TrimSpace(v)
	if v == "" {
		return VerdictNone, ""
	}
	head, _, _ := strings.Cut(v, "\n")
	head, _, _ = strings.Cut(head, ":")
	n := normVerdict(head)
	for _, w := range l.VerdictSendBack {
		if n == w {
			return VerdictSendBack, w
		}
	}
	for _, w := range l.VerdictAccept {
		if n == w {
			return VerdictAccept, w
		}
	}
	return VerdictUnknown, strings.TrimSpace(head)
}
