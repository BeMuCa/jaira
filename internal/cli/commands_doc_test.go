package cli

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestEveryCommandIsDocumented: the README's '## Commands' list and
// docs/COMMANDS.md name every command 'jaira --help' lists. Nothing else ties
// the two documents to the binary, and on 2026-10-13 the README was missing
// nine commands and COMMANDS.md six — a command the docs do not name is one a
// person or an agent reading them never finds.
//
// The list is the one cobra's usage template prints: the root's available
// commands, which leaves out hidden ones like merge-driver. help and completion
// are not on a fresh root — cobra adds them at Execute and documents them
// itself.
func TestEveryCommandIsDocumented(t *testing.T) {
	read := func(path string) string {
		t.Helper()
		b, err := os.ReadFile(filepath.Join("..", "..", path))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	_, section, ok := strings.Cut(read("README.md"), "\n## Commands\n")
	if !ok {
		t.Fatal("README.md has no '## Commands' section")
	}
	section, _, _ = strings.Cut(section, "\n## ")
	reference := read(filepath.Join("docs", "COMMANDS.md"))

	for _, c := range newRoot("test").Commands() {
		if !c.IsAvailableCommand() {
			continue
		}
		name := regexp.QuoteMeta(c.Name())
		if !regexp.MustCompile(`(?m)^jaira ` + name + `( |$)`).MatchString(section) {
			t.Errorf("README.md '## Commands' has no line for 'jaira %s'", c.Name())
		}
		if !regexp.MustCompile("(?m)^\\| `jaira " + name + "[ `]").MatchString(reference) {
			t.Errorf("docs/COMMANDS.md has no table row for 'jaira %s'", c.Name())
		}
	}
}
