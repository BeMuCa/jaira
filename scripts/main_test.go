package scripts_test

import (
	"os"
	"testing"

	"github.com/BeMuCa/jaira/internal/testhome"
)

// TestMain keeps this package's tests out of the machine's ~/.jaira and out of
// the developer's git repository; see internal/testhome.
func TestMain(m *testing.M) {
	os.Exit(testhome.Run(m))
}
