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

// The refresher re-execs the running binary as "<exe> self upgrade --check",
// which only means anything if that binary is jaira. Inside a Go test binary it
// runs the whole suite again, detached — reaching the network the test suite
// promises not to touch, and on Windows holding its own image open so 'go test'
// cannot remove it after every test has passed. That is how this was found:
// master went red on the Windows runner with every package reporting ok.
func TestSpawnRefreshDoesNothingFromATestBinary(t *testing.T) {
	if err := SpawnRefresh(); err != nil {
		t.Fatalf("SpawnRefresh from a test binary: %v", err)
	}
	// The guard is what makes the call above a no-op; assert it directly too,
	// on both platforms' naming.
	for _, exe := range []string{
		"/tmp/go-build123/b001/tui.test",
		`C:\Users\x\AppData\Local\Temp\go-build1\b301\tui.test.exe`,
	} {
		if !isTestBinary(exe) {
			t.Errorf("isTestBinary(%q) = false, want true", exe)
		}
	}
	for _, exe := range []string{"/home/berk/.local/bin/jaira", `C:\bin\jaira.exe`, "/usr/local/bin/jaira-dev"} {
		if isTestBinary(exe) {
			t.Errorf("isTestBinary(%q) = true, want false — a real install must still refresh", exe)
		}
	}
}
