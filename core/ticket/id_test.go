package ticket

import (
	"strings"
	"testing"
	"unicode/utf8"
)

// TestSlugCutsAtRuneBoundary puts a multi-byte letter across byte 48, where a
// byte cut used to leave half a character in the filename. Each case pads
// with ASCII so the wide letter starts at byte 47 (or 46 for four bytes),
// i.e. it straddles the limit.
func TestSlugCutsAtRuneBoundary(t *testing.T) {
	pad := func(n int) string { return strings.Repeat("a", n) }
	cases := []struct {
		name, title, want string
	}{
		{"umlaut", pad(47) + "ü" + "x", pad(47)},
		{"eszett", pad(47) + "ß" + "x", pad(47)},
		{"cyrillic", pad(47) + "ж" + "x", pad(47)},
		{"three-byte letter", pad(46) + "中" + "x", pad(46)},
		{"four-byte letter", pad(45) + "𠀀" + "x", pad(45)},
		// An emoji is no letter: it becomes a dash, so the cut lands after it
		// and the trailing dash is trimmed.
		{"emoji at the limit", pad(47) + "🙂" + "x", pad(47)},
		{"umlaut ends exactly at the limit", pad(46) + "ü" + "x", pad(46) + "ü"},
		{"dash before the cut is trimmed", pad(46) + " " + "ü" + "x", pad(46)},
		{"short title untouched", "Ampel grün", "ampel-grün"},
		// A whole German sentence, cut at the limit.
		{"german sentence", "Die Ampel wird grün, wenn alle Gates in der Pipeline grün sind",
			"die-ampel-wird-grün-wenn-alle-gates-in-der-pipe"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Slug(c.title)
			if !utf8.ValidString(got) {
				t.Fatalf("Slug(%q) = %q, not valid UTF-8", c.title, got)
			}
			if len(got) > slugMax {
				t.Fatalf("Slug(%q) = %q, %d bytes > %d", c.title, got, len(got), slugMax)
			}
			if got != c.want {
				t.Fatalf("Slug(%q) = %q, want %q", c.title, got, c.want)
			}
		})
	}
}

// TestSlugAlwaysValidUTF8 sweeps every cut position for each wide letter,
// so no padding length can land the limit inside one.
func TestSlugAlwaysValidUTF8(t *testing.T) {
	for _, wide := range []string{"ü", "ß", "ж", "中", "𠀀", "🙂", "\xff"} {
		for n := 40; n <= 50; n++ {
			title := strings.Repeat("a", n) + strings.Repeat(wide, 4)
			if got := Slug(title); !utf8.ValidString(got) || len(got) > slugMax {
				t.Fatalf("Slug(%q) = %q (%d bytes)", title, got, len(got))
			}
		}
	}
}
