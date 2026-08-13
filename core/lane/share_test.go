package lane

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestExportCopiesVerbatim asserts Export copies a lane's bytes unchanged —
// the file format is the API, and round-tripping through the parsed struct
// would reorder fields.
func TestExportCopiesVerbatim(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	src := writeLane(t, catalogue, "triage.md", `---
id: triage
name: Triage
after: human
precedence: 41
---
A distinctive prompt body.
`)
	want, err := os.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}

	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	l, ok := set.Get("triage")
	if !ok {
		t.Fatal("triage did not load")
	}

	dstDir := t.TempDir()
	dst, err := Export(l, dstDir, false)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Errorf("Export did not copy bytes verbatim:\ngot:  %q\nwant: %q", got, want)
	}
}

// TestExportDerivesFilenameFromID asserts the destination filename comes
// from the lane's validated ID, never the source filename — a source
// filename is not trustworthy input for a path (T-5-03).
func TestExportDerivesFilenameFromID(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	writeLane(t, catalogue, "some-weird-source-name.md", `---
id: myid
name: My Lane
after: human
precedence: 41
---
`)
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	l, ok := set.Get("myid")
	if !ok {
		t.Fatal("myid did not load")
	}
	dstDir := t.TempDir()
	dst, err := Export(l, dstDir, false)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(dst) != "myid.md" {
		t.Errorf("Export named the file %q, want %q", filepath.Base(dst), "myid.md")
	}
}

// TestExportRefusesOverwrite asserts Export never clobbers an existing file
// unless the caller explicitly asks.
func TestExportRefusesOverwrite(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	writeLane(t, catalogue, "triage.md", `---
id: triage
name: Triage
after: human
precedence: 41
---
`)
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	l, _ := set.Get("triage")
	dstDir := t.TempDir()

	if _, err := Export(l, dstDir, false); err != nil {
		t.Fatalf("first export: %v", err)
	}
	if _, err := Export(l, dstDir, false); err == nil {
		t.Fatal("expected a second export without overwrite to refuse")
	}
	if _, err := Export(l, dstDir, true); err != nil {
		t.Fatalf("export with overwrite=true should succeed, got: %v", err)
	}
}

// TestPublishStampsCreatorWhenAbsent asserts Publish writes creator: as a
// line insert, only when the file does not already declare one.
func TestPublishStampsCreatorWhenAbsent(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	writeLane(t, catalogue, "triage.md", `---
id: triage
name: Triage
after: human
precedence: 41
---
`)
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	l, _ := set.Get("triage")
	dstDir := t.TempDir()

	dst, err := Publish(l, dstDir, "alex", false)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "creator: alex") {
		t.Errorf("published file missing stamped creator, got:\n%s", got)
	}
	// The published file must still parse, and the prompt/other fields must
	// be unaffected by the insert.
	pl, err := parse(got, dst, false)
	if err != nil {
		t.Fatalf("published file did not parse: %v", err)
	}
	if pl.Creator != "alex" {
		t.Errorf("Creator = %q, want %q", pl.Creator, "alex")
	}
	if pl.ID != "triage" || pl.Name != "Triage" {
		t.Errorf("publish must not disturb other fields, got id=%q name=%q", pl.ID, pl.Name)
	}
}

// TestPublishDoesNotOverwriteExistingCreator asserts a lane that already
// declares a creator keeps it — publishing is not a way to reassign credit.
func TestPublishDoesNotOverwriteExistingCreator(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	writeLane(t, catalogue, "triage.md", `---
id: triage
name: Triage
after: human
precedence: 41
creator: original-author
---
`)
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	l, _ := set.Get("triage")
	dstDir := t.TempDir()

	dst, err := Publish(l, dstDir, "someone-else", false)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(got), "someone-else") {
		t.Errorf("publish must not override an existing creator: field, got:\n%s", got)
	}
	if !strings.Contains(string(got), "creator: original-author") {
		t.Errorf("published file lost its original creator, got:\n%s", got)
	}
}

// TestPublishRefusesOverwrite mirrors Export's refusal.
func TestPublishRefusesOverwrite(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	writeLane(t, catalogue, "triage.md", `---
id: triage
name: Triage
after: human
precedence: 41
---
`)
	set, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	l, _ := set.Get("triage")
	dstDir := t.TempDir()

	if _, err := Publish(l, dstDir, "alex", false); err != nil {
		t.Fatalf("first publish: %v", err)
	}
	if _, err := Publish(l, dstDir, "alex", false); err == nil {
		t.Fatal("expected a second publish without overwrite to refuse")
	}
}

// TestAdoptParsesBeforeCopying asserts an unparseable file never lands in
// the catalogue.
func TestAdoptParsesBeforeCopying(t *testing.T) {
	srcDir := t.TempDir()
	src := filepath.Join(srcDir, "broken.md")
	if err := os.WriteFile(src, []byte("not frontmatter at all"), 0o644); err != nil {
		t.Fatal(err)
	}
	dstDir := t.TempDir()

	if _, _, err := Adopt(src, dstDir, false); err == nil {
		t.Fatal("expected Adopt to refuse an unparseable file")
	}
	entries, _ := os.ReadDir(dstDir)
	if len(entries) != 0 {
		t.Errorf("an unparseable file must never land in the destination, found: %v", entries)
	}
}

// TestAdoptCopiesUnderParsedID asserts Adopt names the destination from the
// parsed lane's ID, not the source filename.
func TestAdoptCopiesUnderParsedID(t *testing.T) {
	srcDir := t.TempDir()
	src := writeLane(t, srcDir, "someones-export.md", `---
id: adopted-lane
name: Adopted
after: human
precedence: 41
---
`)
	dstDir := t.TempDir()

	l, dst, err := Adopt(src, dstDir, false)
	if err != nil {
		t.Fatal(err)
	}
	if l.ID != "adopted-lane" {
		t.Errorf("Adopt returned lane id %q, want %q", l.ID, "adopted-lane")
	}
	if filepath.Base(dst) != "adopted-lane.md" {
		t.Errorf("Adopt named the file %q, want %q", filepath.Base(dst), "adopted-lane.md")
	}
}

// TestAdoptRefusesOverwriteUnlessConfirmed mirrors Export/Publish, since an
// adopted id colliding with an existing catalogue lane must be a deliberate
// choice.
func TestAdoptRefusesOverwriteUnlessConfirmed(t *testing.T) {
	srcDir := t.TempDir()
	src := writeLane(t, srcDir, "a.md", `---
id: dup-adopt
name: A
after: human
precedence: 41
---
`)
	dstDir := t.TempDir()
	if _, _, err := Adopt(src, dstDir, false); err != nil {
		t.Fatalf("first adopt: %v", err)
	}
	if _, _, err := Adopt(src, dstDir, false); err == nil {
		t.Fatal("expected a second adopt without overwrite to refuse")
	}
	if _, _, err := Adopt(src, dstDir, true); err != nil {
		t.Fatalf("adopt with overwrite=true should succeed, got: %v", err)
	}
}

// TestBytesBuiltinMatchesEmbedded asserts Bytes reads a built-in's original
// file rather than reconstructing one from the parsed struct.
func TestBytesBuiltinMatchesEmbedded(t *testing.T) {
	lanes, err := Builtins()
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range lanes {
		b, err := Bytes(l)
		if err != nil {
			t.Fatalf("Bytes(%s): %v", l.ID, err)
		}
		reparsed, err := parse(b, l.Source, true)
		if err != nil {
			t.Fatalf("Bytes(%s) did not round-trip through parse: %v", l.ID, err)
		}
		if reparsed.ID != l.ID || reparsed.Prompt != l.Prompt {
			t.Errorf("Bytes(%s) does not match the loaded lane", l.ID)
		}
	}
}

// TestDriftReportsChangedProjectCopy covers D-02: a project lane whose bytes
// differ from the catalogue file of the same id is reported.
func TestDriftReportsChangedProjectCopy(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	writeLane(t, catalogue, "triage.md", `---
id: triage
name: Triage (catalogue version)
after: human
precedence: 41
---
`)

	root := t.TempDir()
	projDir := ProjectLanesDir(root)
	writeLane(t, projDir, "triage.md", `---
id: triage
name: Triage (project version, edited)
after: human
precedence: 41
---
`)

	set, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	drift, err := Drift(root, set)
	if err != nil {
		t.Fatal(err)
	}
	if len(drift) != 1 || drift[0].ID != "triage" {
		t.Fatalf("expected one drift entry for triage, got: %v", drift)
	}
}

// TestDriftSilentWhenBytesMatch asserts identical bytes report no drift.
func TestDriftSilentWhenBytesMatch(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	body := `---
id: triage
name: Triage
after: human
precedence: 41
---
`
	writeLane(t, catalogue, "triage.md", body)

	root := t.TempDir()
	projDir := ProjectLanesDir(root)
	writeLane(t, projDir, "triage.md", body)

	set, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	drift, err := Drift(root, set)
	if err != nil {
		t.Fatal(err)
	}
	if len(drift) != 0 {
		t.Errorf("identical bytes must not be reported as drift, got: %v", drift)
	}
}

// TestDriftIgnoresProjectsWithoutAnActiveDirectory asserts Drift is a no-op
// when D-03's project directory is not authoritative for this root.
func TestDriftIgnoresProjectsWithoutAnActiveDirectory(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	writeLane(t, catalogue, "triage.md", `---
id: triage
name: Triage
after: human
precedence: 41
---
`)
	root := t.TempDir() // no .jaira/lanes directory at all
	set, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	drift, err := Drift(root, set)
	if err != nil {
		t.Fatal(err)
	}
	if len(drift) != 0 {
		t.Errorf("a root with no active project directory must report no drift, got: %v", drift)
	}
}

// TestRefreshDriftPullsCatalogueVersion asserts RefreshDrift overwrites the
// project's copy with the catalogue's, in that direction only.
func TestRefreshDriftPullsCatalogueVersion(t *testing.T) {
	catalogue := t.TempDir()
	t.Setenv("JAIRA_LANES_DIR", catalogue)
	catBody := `---
id: triage
name: Triage (catalogue version)
after: human
precedence: 41
---
`
	writeLane(t, catalogue, "triage.md", catBody)

	root := t.TempDir()
	projDir := ProjectLanesDir(root)
	writeLane(t, projDir, "triage.md", `---
id: triage
name: Triage (project version, edited)
after: human
precedence: 41
---
`)

	set, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	drift, err := Drift(root, set)
	if err != nil {
		t.Fatal(err)
	}
	if len(drift) != 1 {
		t.Fatalf("expected drift before refresh, got: %v", drift)
	}
	if err := RefreshDrift(drift[0]); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(drift[0].ProjectPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != catBody {
		t.Errorf("RefreshDrift did not pull the catalogue bytes:\ngot:  %q\nwant: %q", got, catBody)
	}
}
