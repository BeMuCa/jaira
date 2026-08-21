package tui

import (
	"os"
	"testing"
)

// TestMain disables the background release check for the whole package by
// default.
//
// A Model or Home built by any of this package's test cases starts with an
// empty, per-test JAIRA_HOME (see newTestStore in render_test.go and
// home_test.go's own setup) — which selfupdate.PollCache reads as "never
// checked" and answers by touching the cache and spawning a detached
// refresher: one real outbound HTTPS request to api.github.com, from every
// single test that builds either. That is exactly the network dependency
// this project's test suite promises not to have, so it is off here by
// default; the handful of tests that actually exercise the version
// indicator (updatecheck_test.go) re-enable it locally with t.Setenv, which
// correctly restores this default once they finish.
func TestMain(m *testing.M) {
	os.Setenv("JAIRA_NO_UPDATE_CHECK", "1")
	os.Exit(m.Run())
}
