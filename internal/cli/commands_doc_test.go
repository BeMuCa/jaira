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
// commands and help, which leaves out hidden ones like merge-driver. help and
// completion are not on a fresh root until cobra adds them at Execute, so they
// are added here the same way.
func TestEveryCommandIsDocumented(t *testing.T) {
	read := func(path string) string {
		t.Helper()
		b, err := os.ReadFile(filepath.Join("..", "..", path))
		if err != nil {
			t.Fatal(err)
		}
		// Neither file is pinned eol=lf in .gitattributes, so a Windows
		// checkout with core.autocrlf hands both over with CRLF line endings.
		return strings.ReplaceAll(string(b), "\r\n", "\n")
	}
	_, section, ok := strings.Cut(read("README.md"), "\n## Commands\n")
	if !ok {
		t.Fatal("README.md has no '## Commands' section")
	}
	section, _, _ = strings.Cut(section, "\n## ")
	reference := read(filepath.Join("docs", "COMMANDS.md"))

	root := newRoot("test")
	root.InitDefaultHelpCmd()
	root.InitDefaultCompletionCmd()
	for _, c := range root.Commands() {
		if !c.IsAvailableCommand() && c.Name() != "help" {
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
