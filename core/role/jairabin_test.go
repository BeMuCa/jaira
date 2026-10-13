package role

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
	built.once.Do(func() {
		if built.dir, built.err = os.MkdirTemp("", "jaira-bin-"); built.err != nil {
			return
		}
		bin := filepath.Join(built.dir, "jaira")
		if runtime.GOOS == "windows" {
			bin += ".exe"
		}
		//wintrap:ok the .exe suffix is on bin already, from the line above
		if out, err := exec.Command("go", "build", "-o", bin, "github.com/BeMuCa/jaira/cmd/jaira").CombinedOutput(); err != nil {
			built.err = fmt.Errorf("%w\n%s", err, out)
		}
	})
	if built.err != nil {
		t.Fatalf("building jaira: %v", built.err)
	}
	return built.dir + string(os.PathListSeparator) + os.Getenv("PATH")
}
