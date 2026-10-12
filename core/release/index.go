package release

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// WriteIndex works out which release each note file under <pkgDir>/notes
// first shipped in and writes the answer to <pkgDir>/notes/versions, where
// the next build embeds it.
//
// A note belongs to the oldest tag v* whose tree already holds it. Only tags
// HEAD contains are asked, so rebuilding an old tag never borrows a version
// from a later one, and a note no tag holds yet is left out — Notes() shows
// it as Unreleased. Nothing is committed: the release build runs this in its
// own checkout, and the repository never records a version for a note, so
// releasing needs no commit of its own.
func WriteIndex(pkgDir string) error {
	prefix, err := git(pkgDir, "rev-parse", "--show-prefix")
	if err != nil {
		return err
	}
	tags, err := git(pkgDir, "tag", "--list", "v*", "--merged", "HEAD", "--sort=version:refname")
	if err != nil {
		return err
	}

	first := map[string]string{}
	var versions []string
	for _, tag := range strings.Fields(tags) {
		// A tag from before notes/ existed has no such tree; that is the
		// expected answer for it, not an error. --full-tree, because from a
		// subdirectory ls-tree otherwise filters even a tree object by the
		// working directory's path and lists nothing.
		names, err := git(pkgDir, "ls-tree", "--full-tree", "--name-only", tag+":"+prefix+notesDir)
		if err != nil {
			continue
		}
		version := strings.TrimPrefix(tag, "v")
		for _, name := range strings.Fields(names) {
			if strings.HasSuffix(name, ".md") && first[name] == "" {
				if len(versions) == 0 || versions[len(versions)-1] != version {
					versions = append(versions, version)
				}
				first[name] = version
			}
		}
	}

	var b strings.Builder
	b.WriteString("# Written by go generate (core/release/gen) in the release build; never commit it.\n")
	for i := len(versions) - 1; i >= 0; i-- {
		var names []string
		for name, v := range first {
			if v == versions[i] {
				names = append(names, name)
			}
		}
		sort.Strings(names)
		for _, name := range names {
			fmt.Fprintf(&b, "%s %s\n", versions[i], name)
		}
	}
	return os.WriteFile(filepath.Join(pkgDir, notesDir, indexFile), []byte(b.String()), 0o644)
}

func git(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(out)), nil
}
