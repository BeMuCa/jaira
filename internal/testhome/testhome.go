// Package testhome keeps a test binary out of the developer's ~/.jaira.
//
// Production code finds per-user state through os.UserHomeDir (HOME, or
// USERPROFILE on Windows) and JAIRA_HOME: the default board, the lane
// catalogue, settings, identity, session state. A test that does not isolate
// all of them reads whatever the machine running it happens to have — on one
// machine a default-board.md without a backlog lane failed sixteen tests that
// pass in CI, where the runner has no ~/.jaira at all.
//
// Isolating it per test is the wrong default: the test that forgets is silently
// wrong on exactly the machines nobody looks at. So every package whose tests
// can reach that code calls Run from its TestMain, and a test that wants a home
// of its own still sets one with t.Setenv.
package testhome

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// goVars are go's own per-user locations. Each defaults to somewhere under
// HOME, so a test that runs 'go build' under the temporary home would start
// from an empty module cache — downloading every dependency — and a cold
// build cache, and would lose whatever 'go env -w' wrote.
var goVars = []string{"GOPATH", "GOMODCACHE", "GOCACHE", "GOENV"}

// Run points HOME, USERPROFILE and JAIRA_HOME at one empty temporary
// directory, clears the variables that override single files under ~/.jaira,
// runs the package's tests and removes the directory again. It returns the
// exit code for os.Exit.
func Run(m *testing.M) int {
	pinGoVars()
	home, err := os.MkdirTemp("", "jaira-test-home")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(home)
	for _, k := range []string{"HOME", "USERPROFILE", "JAIRA_HOME"} {
		os.Setenv(k, home)
	}
	for _, k := range []string{"JAIRA_DEFAULT_BOARD", "JAIRA_LANES_DIR"} {
		os.Unsetenv(k)
	}
	isolateGit()
	return m.Run()
}

// gitRepoVars make git pick a repository regardless of the directory it is
// run in: exported by a git hook or a wrapper, they turn 'git -C <temp dir>'
// into a read of — or a write to — the developer's own checkout, its people
// ref and its .git/config included.
var gitRepoVars = []string{
	"GIT_DIR", "GIT_WORK_TREE", "GIT_COMMON_DIR", "GIT_INDEX_FILE",
	"GIT_OBJECT_DIRECTORY", "GIT_ALTERNATE_OBJECT_DIRECTORIES", "GIT_NAMESPACE",
}

// isolateGit clears gitRepoVars and stops git's upward search at the temp
// directory, so a test's temporary directory is never taken for part of a
// repository that happens to contain TMPDIR.
func isolateGit() {
	for _, k := range gitRepoVars {
		os.Unsetenv(k)
	}
	os.Setenv("GIT_CEILING_DIRECTORIES", os.TempDir())
}

// pinGoVars sets each of goVars that is not set already to the value go
// resolves for it under the real home. Without a go on PATH there is nothing
// to keep pointing at, and nothing is set.
func pinGoVars() {
	out, err := exec.Command("go", append([]string{"env"}, goVars...)...).Output()
	if err != nil {
		return
	}
	for i, v := range strings.Split(strings.TrimRight(string(out), "\r\n"), "\n") {
		if i < len(goVars) && os.Getenv(goVars[i]) == "" {
			os.Setenv(goVars[i], strings.TrimRight(v, "\r"))
		}
	}
}
