package ticket

import "testing"

// TestParsePRRef runs every input class a blocked-by entry can carry through
// the one parser every reader uses: what is a pull request, what is a
// malformed one, and what is a ticket id.
func TestParsePRRef(t *testing.T) {
	cases := []struct {
		in     string
		n      int
		isPR   bool
		failed bool
	}{
		{"pr#12", 12, true, false},
		{"PR#12", 12, true, false},
		{"Pr#7", 7, true, false},
		{" pr#12 ", 12, true, false},
		{"pr#2147483647", 2147483647, true, false},
		{"pr#", 0, true, true},
		{"pr#0", 0, true, true},
		{"pr#-1", 0, true, true},
		{"pr#abc", 0, true, true},
		{"pr# 12", 0, true, true},
		{"pr#12a", 0, true, true},
		{"pr#1.5", 0, true, true},
		{"pr#１２", 0, true, true}, // full-width digits are no number
		{"pr#2147483648", 0, true, true},
		{"pr#99999999999999999999", 0, true, true},
		{"01M7WWYG213BWB06A7D1CRASZS", 0, false, false},
		{"CRASZS", 0, false, false},
		{"", 0, false, false},
		{"pr", 0, false, false},
		{"#12", 0, false, false},
	}
	for _, c := range cases {
		n, isPR, err := ParsePRRef(c.in)
		if n != c.n || isPR != c.isPR || (err != nil) != c.failed {
			t.Errorf("ParsePRRef(%q) = %d, %v, %v; want %d, %v, failed=%v", c.in, n, isPR, err, c.n, c.isPR, c.failed)
		}
		if got, ok := PRRef(c.in); ok != (c.isPR && !c.failed) || got != c.n {
			t.Errorf("PRRef(%q) = %d, %v", c.in, got, ok)
		}
	}
	if got := FormatPRRef(12); got != "pr#12" {
		t.Errorf("FormatPRRef(12) = %q", got)
	}
}

// A pull request is its own short reference: Handle's tail of a long one
// would be a bare number.
func TestHandleKeepsPRRef(t *testing.T) {
	for in, want := range map[string]string{
		"pr#12":                      "pr#12",
		"pr#1234567":                 "pr#1234567",
		"01M7WWYG213BWB06A7D1CRASZS": "CRASZS",
	} {
		if got := Handle(in); got != want {
			t.Errorf("Handle(%q) = %q, want %q", in, got, want)
		}
	}
}
