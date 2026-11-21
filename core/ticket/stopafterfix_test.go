package ticket

import "testing"

// stop-after-fix is read the way jaira set stores it: a hand-written "True"
// or " true " is on, anything that is not true is off.
func TestStopAfterFixReadsLikeSetStoresIt(t *testing.T) {
	for _, tc := range []struct {
		value string
		want  bool
	}{
		{"true", true}, {"True", true}, {"TRUE", true}, {"' true '", true},
		{"false", false}, {"''", false}, {"yes", false},
	} {
		d, err := ParseDoc([]byte("---\nid: 01M7WWZWSDHK37DGGHJFDF0YD6\ntitle: t\nstatus: todo\nstop-after-fix: " + tc.value + "\n---\n"))
		if err != nil {
			t.Fatalf("%s: %v", tc.value, err)
		}
		tk, err := Decode(d, "x.md")
		if err != nil {
			t.Fatalf("%s: %v", tc.value, err)
		}
		if tk.StopAfterFix != tc.want {
			t.Errorf("stop-after-fix: %s reads as %v, want %v", tc.value, tk.StopAfterFix, tc.want)
		}
	}
}
