package cli

import (
	"os"
	"testing"
)

// TestMain gives the whole package a JAIRA_HOME of its own.
//
// Without it, every test that opens a store and does not set the variable
// itself writes sessions/ and locks/ into the developer's actual ~/.jaira,
// keyed by a t.TempDir() path that stops existing the moment the test ends. One
// run of this package left 34 of them behind; the machine this was found on had
// collected 3115, which buried the two dozen directories that belong to real
// checkouts.
//
// Set here rather than in each test because the leak is the default: a test
// that forgets is silently wrong in a place nobody looks, and there is no
// reason for any test in this package to want the real home.
func TestMain(m *testing.M) {
	home, err := os.MkdirTemp("", "jaira-cli-test-home")
	if err != nil {
		panic(err)
	}
	os.Setenv("JAIRA_HOME", home)
	code := m.Run()
	os.RemoveAll(home)
	os.Exit(code)
}
