// Command jaira is a git-native kanban board for work done with coding agents.
package main

import (
	"os"

	"github.com/berk/jaira/internal/cli"
)

// version is overridden at build time via -ldflags.
var version = "dev"

func main() {
	os.Exit(cli.Execute(version))
}
