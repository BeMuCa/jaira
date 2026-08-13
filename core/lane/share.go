package lane

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BeMuCa/jaira/core/ticket"
)

// Bytes returns a lane's file content verbatim: from the embedded binary for
// a built-in, from disk for anything else. Callers that need to copy a lane
// (Export, Publish, Materialise) go through this rather than re-serialising
// the parsed struct, so unknown fields, comments and field order all survive
// — the file format is the API.
func Bytes(l *Lane) ([]byte, error) {
	if l.Builtin {
		name := strings.TrimPrefix(l.Source, "builtin:")
		// path.Join, not filepath.Join: builtinFS is an embedded filesystem and
		// always uses forward slashes (see Builtins).
		return builtinFS.ReadFile(path.Join("builtin", name))
	}
	return os.ReadFile(l.Source)
}

// copyLane is the shared body of Export and Publish: read the lane's bytes
// verbatim, optionally stamp a creator: line, and write them to
// <dstDir>/<id>.md. The filename always comes from the lane's own,
// already-validID-constrained ID — never from the source path — so a
// maliciously named source file cannot escape dstDir (T-5-03).
func copyLane(l *Lane, dstDir, stampCreator string, overwrite bool) (string, error) {
	b, err := Bytes(l)
	if err != nil {
		return "", err
	}
	if stampCreator != "" && l.Creator == "" {
		b = stampCreatorLine(b, stampCreator)
	}
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return "", err
	}
	dst := filepath.Join(dstDir, l.ID+".md")
	if !overwrite {
		if _, statErr := os.Stat(dst); statErr == nil {
			return "", fmt.Errorf("%s already exists; refusing to overwrite", dst)
		} else if !os.IsNotExist(statErr) {
			return "", statErr
		}
	}
	if err := os.WriteFile(dst, b, 0o644); err != nil {
		return "", err
	}
	return dst, nil
}

// stampCreatorLine inserts "creator: <who>" as the first line after the
// opening "---", a line insert rather than a YAML rewrite: parsing and
// re-serialising the frontmatter would reorder fields and drop anything this
// tool does not recognise (the reserved external: block, for one).
func stampCreatorLine(b []byte, who string) []byte {
	s := string(b)
	if !strings.HasPrefix(s, "---") {
		return b // not well-formed frontmatter; leave it alone rather than guess
	}
	nl := strings.IndexByte(s, '\n')
	if nl < 0 {
		return b
	}
	return []byte(s[:nl+1] + "creator: " + who + "\n" + s[nl+1:])
}

// Export copies a lane's file, verbatim, from wherever it was loaded from
// into dstDir — the catalogue-to-project move. It refuses to overwrite an
// existing file unless overwrite is true, so the caller (the lane settings
// screen) decides rather than silently clobbering a project's own copy.
func Export(l *Lane, dstDir string, overwrite bool) (string, error) {
	return copyLane(l, dstDir, "", overwrite)
}

// Publish copies a lane's file into dstDir (normally
// Store.SharedDir()/identity.Slug(who)), stamping creator: with who when the
// file does not already declare one. It refuses to overwrite an existing
// file unless overwrite is true.
func Publish(l *Lane, dstDir, who string, overwrite bool) (string, error) {
	return copyLane(l, dstDir, who, overwrite)
}

// Adopt parses srcPath first — an unparseable file must never land in the
// catalogue — then copies its bytes verbatim into dstDir under the parsed
// lane's own id, never the source filename. It refuses to overwrite an
// existing file unless overwrite is true.
func Adopt(srcPath, dstDir string, overwrite bool) (*Lane, string, error) {
	b, err := os.ReadFile(srcPath)
	if err != nil {
		return nil, "", err
	}
	l, err := parse(b, srcPath, false)
	if err != nil {
		return nil, "", err
	}
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return nil, "", err
	}
	dst := filepath.Join(dstDir, l.ID+".md")
	if !overwrite {
		if _, statErr := os.Stat(dst); statErr == nil {
			return nil, "", fmt.Errorf("%s already exists; refusing to overwrite", dst)
		} else if !os.IsNotExist(statErr) {
			return nil, "", statErr
		}
	}
	if err := os.WriteFile(dst, b, 0o644); err != nil {
		return nil, "", err
	}
	return l, dst, nil
}

// DriftEntry names one lane whose project copy no longer matches its
// catalogue copy of the same id.
type DriftEntry struct {
	ID            string
	ProjectPath   string
	CataloguePath string
}

// Drift implements D-02: it compares every lane this project's directory
// supplies against a catalogue file of the same id, and reports the ones
// whose bytes differ. Editing a lane writes through to the catalogue
// immediately, so it is the project's copy that goes stale, not the other
// way round — nothing here fixes that automatically; see RefreshDrift.
//
// Checked only when the caller asks (the lane settings screen, on open),
// per the decided design: checking on every command would mean every
// invocation of the loader pays for a comparison most callers never look at.
func Drift(root string, set *Set) ([]DriftEntry, error) {
	if !ProjectLanesActive(root) {
		return nil, nil
	}
	projDir := ProjectLanesDir(root)
	var out []DriftEntry
	for _, l := range set.Lanes {
		if l.Builtin || !strings.HasPrefix(l.Source, projDir) {
			continue
		}
		catPath := filepath.Join(UserLanesDir(), l.ID+".md")
		catBytes, err := os.ReadFile(catPath)
		if err != nil {
			continue // nothing in the catalogue to have drifted from
		}
		projBytes, err := os.ReadFile(l.Source)
		if err != nil {
			continue
		}
		if !bytes.Equal(catBytes, projBytes) {
			out = append(out, DriftEntry{ID: l.ID, ProjectPath: l.Source, CataloguePath: catPath})
		}
	}
	return out, nil
}

// RefreshDrift pulls the catalogue copy into the project, in the direction
// the user asked for — nothing syncs on its own, per D-02.
func RefreshDrift(d DriftEntry) error {
	b, err := os.ReadFile(d.CataloguePath)
	if err != nil {
		return err
	}
	return os.WriteFile(d.ProjectPath, b, 0o644)
}

// SharedLane is one lane found under a project's .jaira/shared/ tree,
// alongside the folder (normally the publisher's identity.Slug) it came
// from.
type SharedLane struct {
	Lane   *Lane
	Folder string
	Path   string
}

// Shared walks <root>/.jaira/shared/*/*.md and returns every lane it can
// parse, alongside the folder it came from and its path. A file that fails
// to read or parse is skipped and reported in warnings rather than failing
// the walk — one teammate's broken file must not hide everyone else's lanes.
//
// Shared lanes are never loaded by Load: they are visible here, adoptable
// through Adopt, and otherwise inert. That separation is the whole security
// story for T-5-02 — a committed shared lane is untrusted prompt content
// that would otherwise run at whatever model-tier it declares — so nothing
// should ever wire this function's result into Load.
func Shared(root string) (lanes []SharedLane, warnings []string, err error) {
	if root == "" {
		return nil, nil, nil
	}
	sharedDir := filepath.Join(root, ticket.DirName, ticket.SharedSubdir)
	matches, globErr := filepath.Glob(filepath.Join(sharedDir, "*", "*.md"))
	if globErr != nil {
		return nil, nil, globErr
	}
	sort.Strings(matches)
	for _, p := range matches {
		b, readErr := os.ReadFile(p)
		if readErr != nil {
			warnings = append(warnings, fmt.Sprintf("could not read shared lane %s: %v", p, readErr))
			continue
		}
		l, parseErr := parse(b, p, false)
		if parseErr != nil {
			warnings = append(warnings, fmt.Sprintf("shared lane %s did not parse and was skipped: %v", p, parseErr))
			continue
		}
		lanes = append(lanes, SharedLane{Lane: l, Folder: filepath.Base(filepath.Dir(p)), Path: p})
	}
	return lanes, warnings, nil
}
