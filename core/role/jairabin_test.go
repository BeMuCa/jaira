package role

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
)

// built is this tree's jaira, compiled once for every script test.
var built struct {
	once sync.Once
	dir  string
	err  error
}

// jairaPath returns a PATH with this tree's jaira first. spawn.sh and
// run-lane.sh read their JSON with 'jaira json get' — not python3, which a
// Windows machine may not have (issue #79) — so a test of them needs the
// jaira that has it, not whichever one the machine has installed.
func jairaPath(t *testing.T) string {
	t.Helper()
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skipf("no go toolchain to build jaira with: %v", err)
	}
	built.once.Do(func() {
		if built.dir, built.err = os.MkdirTemp("", "jaira-bin-"); built.err != nil {
			return
		}
		cmd := exec.Command(goBin, "build", "-o", filepath.Join(built.dir, "jaira"), "github.com/BeMuCa/jaira/cmd/jaira")
		if out, err := cmd.CombinedOutput(); err != nil {
			built.err = &buildError{err, string(out)}
		}
	})
	if built.err != nil {
		t.Fatalf("building jaira: %v", built.err)
	}
	return built.dir + string(os.PathListSeparator) + os.Getenv("PATH")
}

type buildError struct {
	err error
	out string
}

func (e *buildError) Error() string { return e.err.Error() + "\n" + e.out }
