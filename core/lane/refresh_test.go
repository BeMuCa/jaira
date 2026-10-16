package lane

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/shipped"
)

func olderTesting(t *testing.T, base *Lane) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", "testing-catalogue-2fb8f177.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !shipped.Known(shippedKey(base), b) {
		t.Fatal("the testing fixture is not in core/shipped/shipped.sha256; regenerate it: scripts/shipped-history.sh > core/shipped/shipped.sha256")
	}
	return b
}

func builtinByID(t *testing.T, id string) *Lane {
	t.Helper()
	bs, err := Builtins()
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range bs {
		if b.ID == id {
			return b
		}
	}
	t.Fatalf("no built-in lane %q", id)
	return nil
}

// boardWithTesting is a board that holds testing (the default board leaves it
// off) with the older shipped file on disk. It returns the root and that path.
func boardWithTesting(t *testing.T) (string, string) {
	t.Helper()
	t.Setenv("JAIRA_HOME", t.TempDir())
	root := setUpBoard(t)
	set, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	p, _, _, err := Add(root, set, "testing")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, olderTesting(t, builtinByID(t, "testing")), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, p
}

func refreshAction(t *testing.T, got []Refreshed, id string) string {
	t.Helper()
	for _, r := range got {
		if r.ID == id {
			return r.Action
		}
	}
	t.Fatalf("lane %q not in %+v", id, got)
	return ""
}

// Every lane this binary ships is a version jaira lists, or the next release
// reports its own lanes as edited. The failure prints the lines to append.
func TestEveryBuiltinLaneIsAShippedVersion(t *testing.T) {
	bs, err := Builtins()
	if err != nil {
		t.Fatal(err)
	}
	var missing []string
	for _, b := range bs {
		raw, err := Bytes(b)
		if err != nil {
			t.Fatal(err)
		}
		key := shippedKey(b)
		if _, err := os.Stat(filepath.Join("..", "..", filepath.FromSlash(key))); err != nil {
			t.Errorf("%s: the key does not name a file in the repository: %v", b.ID, err)
		}
		if !shipped.Known(key, raw) {
			missing = append(missing, shipped.Line(key, raw))
		}
	}
	if len(missing) > 0 {
		t.Errorf("core/shipped/shipped.sha256 lacks the version of %d file(s) this build installs; run scripts/shipped-history.sh > core/shipped/shipped.sha256, or append:\n%s",
			len(missing), strings.Join(missing, "\n"))
	}
}

func TestRefreshShippedUpdatesAnOlderShippedLane(t *testing.T) {
	root, p := boardWithTesting(t)
	got, err := RefreshShipped(root, false)
	if err != nil {
		t.Fatal(err)
	}
	if a := refreshAction(t, got, "testing"); a != RefreshUpdated {
		t.Errorf("testing: %q, want %q", a, RefreshUpdated)
	}
	for _, r := range got {
		if r.ID != "testing" && r.Action != RefreshUnchanged {
			t.Errorf("%s: %q, want unchanged", r.ID, r.Action)
		}
	}
	want, err := Bytes(builtinByID(t, "testing"))
	if err != nil {
		t.Fatal(err)
	}
	have, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if string(have) != string(want) {
		t.Error("the lane file is not the shipped one after the refresh")
	}
	again, err := RefreshShipped(root, false)
	if err != nil {
		t.Fatal(err)
	}
	if a := refreshAction(t, again, "testing"); a != RefreshUnchanged {
		t.Errorf("second run: testing %q, want unchanged", a)
	}
}

func TestRefreshShippedLeavesAnEditedLaneAlone(t *testing.T) {
	root, p := boardWithTesting(t)
	edited := append(olderTesting(t, builtinByID(t, "testing")), []byte("\nAlso run the linter.\n")...)
	if err := os.WriteFile(p, edited, 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := RefreshShipped(root, false)
	if err != nil {
		t.Fatal(err)
	}
	if a := refreshAction(t, got, "testing"); a != RefreshSkipped {
		t.Errorf("testing: %q, want %q", a, RefreshSkipped)
	}
	have, _ := os.ReadFile(p)
	if string(have) != string(edited) {
		t.Error("an edited lane file was touched")
	}
}

func TestRefreshShippedIgnoresLanesItDoesNotShip(t *testing.T) {
	root, _ := boardWithTesting(t)
	custom := filepath.Join(ProjectLanesDir(root), "mine.md")
	body := "---\nid: mine\nname: Mine\nafter: todo\n---\n"
	if err := os.WriteFile(custom, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := RefreshShipped(root, false)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range got {
		if r.ID == "mine" {
			t.Errorf("a lane jaira does not ship was visited: %+v", r)
		}
	}
	have, _ := os.ReadFile(custom)
	if string(have) != body {
		t.Error("a lane jaira does not ship was touched")
	}
}

func TestRefreshShippedDryRunWritesNothing(t *testing.T) {
	root, p := boardWithTesting(t)
	older := olderTesting(t, builtinByID(t, "testing"))
	got, err := RefreshShipped(root, true)
	if err != nil {
		t.Fatal(err)
	}
	if a := refreshAction(t, got, "testing"); a != RefreshUpdated {
		t.Errorf("testing: %q, want %q", a, RefreshUpdated)
	}
	have, _ := os.ReadFile(p)
	if string(have) != string(older) {
		t.Error("a dry run wrote the lane file")
	}
}
