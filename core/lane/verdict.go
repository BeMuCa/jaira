package lane

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// verdictWords normalises the words a lane file declares. A word with no
// letter or digit in it is refused: it folds to nothing, so no verdict could
// ever be read as it.
func verdictWords(ws []string) ([]string, error) {
	var out []string
	for _, w := range ws {
		if strings.TrimSpace(w) == "" {
			continue
		}
		n := normVerdict(w)
		if n == "" {
			return nil, fmt.Errorf("verdict word %q has no letter or digit, and no verdict can be read as it", w)
		}
		out = append(out, n)
	}
	return out, nil
}

func notWordRune(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) }

// unicodeDashes are the hyphens and dashes a word processor puts where a
// person typed '-'.
var unicodeDashes = strings.NewReplacer("\u2010", "-", "\u2011", "-", "\u2012", "-", "\u2013", "-", "\u2014", "-", "\u2015", "-")

// normVerdict folds case, treats a space or any dash like '-' and drops what
// surrounds the word that is neither letter nor digit — markup, quotes, the
// full stop a sentence ends on — so "Send back", "send-back", "Send‐back",
// "**send-back**" and "Accept." are each one word. Lane words and verdicts
// both fold through it, so the two always compare alike.
func normVerdict(s string) string {
	s = strings.TrimFunc(unicodeDashes.Replace(s), notWordRune)
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
// is neither letter nor digit — so "send-back: DoD 2 unmet", "Send back", "Send back, DoD
// 1" and "send-back — DoD 1" all read — markup in front of it too, as in
// "**send-back**: x" or "> send-back: x" — while "send-backs" and a sentence that
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
// w, followed by its end or by anything that is neither letter nor digit.
func beginsWithWord(n, w string) bool {
	rest, ok := strings.CutPrefix(n, w)
	if !ok {
		return false
	}
	r, _ := utf8.DecodeRuneInString(rest)
	return rest == "" || notWordRune(r)
}
