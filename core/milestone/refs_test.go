package milestone

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fakeLock struct{ held, taken int }

func (f *fakeLock) Lock(string) (func(), error) {
	f.held++
	f.taken++
	return func() { f.held-- }, nil
}

// fakeRec is a Refs over a directory: one file per milestone, named the way a
// ref carries it, so a test sets one up and reads it back as a file. A write
// replaces every file of the same ID, the way a ref holds one name.
type fakeRec struct {
	dir   string
	have  map[string]bool // IDs that have a ref beyond what the directory holds
	names []string
	last  []byte
	err   error
}

// refsAt is the fake over a board's old milestone directory, which is where
// the helpers below put their files.
func refsAt(root string) *fakeRec { return &fakeRec{dir: Dir(root)} }

func (r *fakeRec) HasMilestoneRef(id string) bool { return r.have[id] }

func (r *fakeRec) MilestoneFiles() ([]File, error) {
	entries, err := os.ReadDir(r.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []File
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(r.dir, e.Name()))
		if err != nil {
			continue
		}
		out = append(out, File{Name: e.Name(), Content: b})
	}
	return out, nil
}

func (r *fakeRec) RecordMilestone(key, file string, content []byte) error {
	if r.err != nil {
		return r.err
	}
	r.names = append(r.names, key)
	r.last = content
	if err := os.MkdirAll(r.dir, 0o755); err != nil {
		return err
	}
	entries, _ := os.ReadDir(r.dir)
	for _, e := range entries {
		if _, id, _ := FileIdentity(strings.TrimSuffix(e.Name(), ".md")); id == key && e.Name() != file {
			_ = os.Remove(filepath.Join(r.dir, e.Name()))
		}
	}
	return os.WriteFile(filepath.Join(r.dir, file), content, 0o644)
}

// save records a milestone through the fake at root, as a write would.
func save(t *testing.T, root string, m *Milestone) {
	t.Helper()
	if err := refsAt(root).RecordMilestone(m.ID, m.FileName(), m.Bytes()); err != nil {
		t.Fatal(err)
	}
}

func dirNames(t *testing.T, root string) string {
	t.Helper()
	entries, err := os.ReadDir(Dir(root))
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return strings.Join(names, ",")
}

// writeOld puts a file where the format before IDs kept it: <name>.md.
func writeOld(t *testing.T, root, name, body string) { writeFile(t, root, name+".md", body) }
