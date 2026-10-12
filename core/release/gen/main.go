// Command gen writes core/release/notes/versions, the index that tells a
// release build which release each note file belongs to. It runs through
// go generate from the release package's directory; see release.WriteIndex.
package main

import (
	"fmt"
	"os"

	"github.com/BeMuCa/jaira/core/release"
)

func main() {
	if err := release.WriteIndex("."); err != nil {
		fmt.Fprintln(os.Stderr, "release notes index:", err)
		os.Exit(1)
	}
}
