package milestone

import (
	"strings"
	"testing"
)

func TestIDsAreSixCrockfordCharactersWithALetter(t *testing.T) {
	for range 2000 {
		if id := NewID(); !IsID(id) {
			t.Fatalf("NewID gave %q, which IsID rejects", id)
		}
	}
	for _, s := range []string{"K3F9QZ", "ABC123", "0000A0"} {
		if !IsID(s) {
			t.Errorf("IsID(%q) = false", s)
		}
	}
	for _, s := range []string{"", "202611", "k3f9qz", "K3F9Q", "K3F9QZZ", "K3F9QI", "K3F9QU", "stages"} {
		if IsID(s) {
			t.Errorf("IsID(%q) = true", s)
		}
	}
}

func TestSplitFileKeepsDigitAndWordTailsInTheName(t *testing.T) {
	cases := []struct {
		base, name, id string
		ok             bool
	}{
		{"backend-K3F9QZ", "backend", "K3F9QZ", true},
		{"release-stages-K3F9QZ", "release-stages", "K3F9QZ", true},
		{"sprint-202611", "", "", false},
		{"release-stages", "", "", false},
		{"backend", "", "", false},
		{"-K3F9QZ", "", "", false},
		{"K3F9QZ", "", "", false},
	}
	for _, c := range cases {
		name, id, ok := SplitFile(c.base)
		if name != c.name || id != c.id || ok != c.ok {
			t.Errorf("SplitFile(%q) = %q, %q, %v; want %q, %q, %v", c.base, name, id, ok, c.name, c.id, c.ok)
		}
	}
}

func TestLegacyIDIsStableAndNeverDigitsOnly(t *testing.T) {
	if got, again := LegacyID("backend"), LegacyID("backend"); got != again || got != legacyBackend {
		t.Fatalf("LegacyID(backend) = %q then %q, want %q", got, again, legacyBackend)
	}
	for i := range 3000 {
		name := "m" + strings.Repeat("x", i%5) + string(rune('a'+i%26)) + string(rune('a'+i/26%26))
		if id := LegacyID(name); !IsID(id) {
			t.Fatalf("LegacyID(%q) = %q is not an ID", name, id)
		}
	}
}

const legacyBackend = "23G8MG"
