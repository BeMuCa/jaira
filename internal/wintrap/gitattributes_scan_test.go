package wintrap

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// attributes is the part of .gitattributes this package needs: which paths are
// pinned to LF. It is a matcher of our own rather than a shell-out to
// "git check-attr" so the testdata fixtures can be plain directories instead of
// repositories — and because the only pattern forms this repository uses are
// dir/*.ext and dir/**/*.ext, which path.Match plus a ** segment covers.
type attributes struct {
	rules []attrRule
}

type attrRule struct {
	pattern string
	eolLF   bool
}

// loadAttributes reads root/.gitattributes. A missing file is not an error: it
// means nothing is pinned, which is exactly what the rule 2 check reports on.
func loadAttributes(root string) (*attributes, error) {
	f, err := os.Open(filepath.Join(root, ".gitattributes"))
	if errors.Is(err, fs.ErrNotExist) {
		return &attributes{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	a := &attributes{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		r := attrRule{pattern: fields[0]}
		for _, at := range fields[1:] {
			switch {
			case at == "eol=lf":
				r.eolLF = true
			case at == "-text" || at == "binary":
				r.eolLF = false
			}
		}
		a.rules = append(a.rules, r)
	}
	return a, sc.Err()
}

// eolLF reports whether rel, a slash-separated path relative to the root the
// .gitattributes sits in, is pinned to LF. Git gives the last matching line the
// final say, so the scan runs backwards.
func (a *attributes) eolLF(rel string) bool {
	for i := len(a.rules) - 1; i >= 0; i-- {
		if matchPattern(a.rules[i].pattern, rel) {
			return a.rules[i].eolLF
		}
	}
	return false
}

// matchPattern applies the gitignore-style pattern syntax .gitattributes
// shares: a pattern without a slash matches the basename anywhere, one with a
// slash is anchored at the root, and ** stands for any run of directories.
func matchPattern(pattern, rel string) bool {
	pattern = strings.TrimSuffix(pattern, "/")
	if pattern == "" {
		return false
	}
	if !strings.Contains(strings.TrimPrefix(pattern, "/"), "/") {
		ok, err := path.Match(strings.TrimPrefix(pattern, "/"), path.Base(rel))
		return err == nil && ok
	}
	return matchSegments(strings.Split(strings.TrimPrefix(pattern, "/"), "/"), strings.Split(rel, "/"))
}

func matchSegments(pat, seg []string) bool {
	if len(pat) == 0 {
		return len(seg) == 0
	}
	if pat[0] == "**" {
		// ** matches any run of directories, the empty run included.
		for i := 0; i <= len(seg); i++ {
			if matchSegments(pat[1:], seg[i:]) {
				return true
			}
		}
		return false
	}
	if len(seg) == 0 {
		return false
	}
	ok, err := path.Match(pat[0], seg[0])
	if err != nil || !ok {
		return false
	}
	return matchSegments(pat[1:], seg[1:])
}
