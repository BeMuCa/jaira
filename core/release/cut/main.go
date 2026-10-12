// Command cut releases the notes: it moves every file directly in
// core/release/notes into core/release/notes/<version>/ and prints the version.
// The release workflow runs it as `go run ./core/release/cut <version>` from
// the repository root; see release.Cut.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BeMuCa/jaira/core/release"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: cut <version>")
		os.Exit(2)
	}
	version, err := release.Cut(filepath.Join("core", "release"), os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "cut:", err)
		os.Exit(1)
	}
	fmt.Println(version)
}
