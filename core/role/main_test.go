package role

import (
	"os"
	"testing"

	"github.com/BeMuCa/jaira/internal/testhome"
)

// TestMain keeps this package's tests out of the machine's ~/.jaira; see
// internal/testhome. It also removes the jaira the script tests built.
func TestMain(m *testing.M) {
	code := testhome.Run(m)
	if built.dir != "" {
		os.RemoveAll(built.dir)
	}
	os.Exit(code)
}
