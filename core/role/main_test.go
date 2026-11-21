package role

import (
	"os"
	"testing"

	"github.com/BeMuCa/jaira/internal/testhome"
)

// TestMain keeps this package's tests out of the machine's ~/.jaira; see
// internal/testhome. It also removes the jaira the script tests built.
// spawn.sh waits seconds for claude to come up in a real tab; every Herdr stub
// here answers at once, so the script tests run it without those pauses.
func TestMain(m *testing.M) {
	os.Setenv("JAIRA_SPAWN_PAUSE", "0")
	code := testhome.Run(m)
	if built.dir != "" {
		os.RemoveAll(built.dir)
	}
	os.Exit(code)
}
