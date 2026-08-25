package selfupdate

import (
	"strings"
	"testing"
)

// The overrides move the trust root rather than only the address: checksums.txt
// comes from the same base as the archive it certifies, and the end of that
// path writes over the running binary. Cleartext to an arbitrary host would
// leave nothing behind the one protection there is.
func TestReleaseHostOverrideRefusesCleartextAndOddSchemes(t *testing.T) {
	for _, v := range []string{
		"http://evil.example.com",
		"http://192.168.1.5:8080",
		"ftp://example.com",
		"file:///tmp",
		"example.com",
		"  ",
	} {
		t.Setenv("JAIRA_RELEASE_DOWNLOADS", v)
		if got := downloadBase(); strings.Contains(got, "example.com") || strings.Contains(got, "192.168") || strings.Contains(got, "/tmp") {
			t.Errorf("override %q was honoured: downloadBase() = %q", v, got)
		}
		if len(Overridden()) != 0 {
			t.Errorf("override %q is reported as in effect", v)
		}
	}
}

// https is honoured, and so is loopback http — the test suite serves a whole
// release from an httptest.Server rather than reaching the network.
func TestReleaseHostOverrideHonoursHTTPSAndLoopback(t *testing.T) {
	for _, v := range []string{"https://mirror.example.com", "http://127.0.0.1:54321", "http://localhost:8080/x"} {
		t.Setenv("JAIRA_RELEASE_DOWNLOADS", v)
		if got := downloadBase(); got != v {
			t.Errorf("downloadBase() = %q, want %q", got, v)
		}
		if o := Overridden(); len(o) != 1 || !strings.Contains(o[0], v) {
			t.Errorf("Overridden() = %v, want it to name %q", o, v)
		}
	}
}
