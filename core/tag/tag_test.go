package tag

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizeStoresOneFormPerName(t *testing.T) {
	for _, c := range []struct {
		raw, want string
		changed   bool
	}{
		{"ui", "ui", false},
		{"UI", "ui", true},
		{"My UI", "my-ui", true},
		{"ui  ux", "ui-ux", true},
		{"needs_review", "needs-review", true},
		{"  spaced  ", "spaced", true},
		{"--ui--", "ui", true},
		{"go1", "go1", false},
	} {
		got, changed, err := Normalize(c.raw)
		if err != nil {
			t.Errorf("Normalize(%q): %v", c.raw, err)
			continue
		}
		if got != c.want {
			t.Errorf("Normalize(%q) = %q, want %q", c.raw, got, c.want)
		}
		if changed != c.changed {
			t.Errorf("Normalize(%q) changed = %v, want %v", c.raw, changed, c.changed)
		}
	}
}

// A character that cannot be stored is refused rather than dropped: a silently
// shortened name is a second name for one subject, which is what the shared
// vocabulary exists to prevent.
func TestNormalizeRefusesRatherThanTrims(t *testing.T) {
	for _, raw := range []string{"ui!", "front/end", "über", "", "  ", "---", "a:b"} {
		if got, _, err := Normalize(raw); err == nil {
			t.Errorf("Normalize(%q) = %q, want a refusal", raw, got)
		}
	}
}

// A registry entry is a plain YAML-safe scalar by construction, which is what
// makes 'tags' a list field nothing has to quote.
func TestNormalizedNamesSurviveFrontmatter(t *testing.T) {
	for _, raw := range []string{"My UI", "backend_api", "GO1"} {
		name, _, err := Normalize(raw)
		if err != nil {
			t.Fatalf("Normalize(%q): %v", raw, err)
		}
		for _, r := range name {
			ok := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-'
			if !ok {
				t.Errorf("Normalize(%q) = %q, which carries %q", raw, name, string(r))
			}
		}
	}
}

// The file belongs to whoever is editing it. A write must add or change one
// line and leave every comment, blank line and unparseable line exactly where
// it was — the same promise the ticket writer makes about frontmatter.
func TestSetPreservesEverythingItDidNotWrite(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".jaira"), 0o755); err != nil {
		t.Fatal(err)
	}
	original := "# my own header\n" +
		"\n" +
		"# the ones I care about\n" +
		"backend: 73  # picked by hand\n" +
		"docs: 178\n" +
		"this line is not an entry\n" +
		"ui: 33\n"
	if err := os.WriteFile(Path(root), []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	reg, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if c, ok := reg.Colour("backend"); !ok || c != 73 {
		t.Fatalf("backend colour = %d,%v, want 73,true", c, ok)
	}
	reg.Set("cli", 45)      // new: inserted alphabetically
	reg.Set("backend", 100) // known: its own line rewritten
	if err := reg.Save(root); err != nil {
		t.Fatal(err)
	}

	b, err := os.ReadFile(Path(root))
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	for _, want := range []string{
		"# my own header",
		"# the ones I care about",
		"this line is not an entry",
		"docs: 178",
		"ui: 33",
		"backend: 100  # picked by hand", // recoloured, comment kept
		"cli: 45",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q after a write:\n%s", want, got)
		}
	}
	if strings.Contains(got, "backend: 73") {
		t.Errorf("the old colour survived the rewrite:\n%s", got)
	}
	if strings.Count(got, "backend:") != 1 {
		t.Errorf("backend has %d lines, want 1:\n%s", strings.Count(got, "backend:"), got)
	}
	// Alphabetical insertion, not an append: two teammates adding two tags at
	// once must not both write the file's last line.
	lines := strings.Split(strings.TrimSpace(got), "\n")
	at := func(prefix string) int {
		for i, l := range lines {
			if strings.HasPrefix(l, prefix) {
				return i
			}
		}
		t.Fatalf("no %q line in:\n%s", prefix, got)
		return -1
	}
	if at("cli:") < at("backend:") || at("cli:") > at("docs:") {
		t.Errorf("cli was not inserted between backend and docs:\n%s", got)
	}
	// Reloading sees exactly what was written.
	again, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if c, ok := again.Colour("cli"); !ok || c != 45 {
		t.Errorf("cli colour after reload = %d,%v, want 45,true", c, ok)
	}
	if c, _ := again.Colour("backend"); c != 100 {
		t.Errorf("backend colour after reload = %d, want 100", c)
	}
}

// A hand-written "UI: 39" colours the tag "ui", because that is the only name
// the tickets can carry.
func TestLoadNormalizesHandWrittenNames(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".jaira"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(Path(root), []byte("UI: 39\nbad: 999\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	reg, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if c, ok := reg.Colour("ui"); !ok || c != 39 {
		t.Errorf("ui colour = %d,%v, want 39,true", c, ok)
	}
	// 999 is not an ANSI-256 colour, so that line is not an entry — and being
	// no entry is what keeps it in the file untouched.
	if _, ok := reg.Colour("bad"); ok {
		t.Error("an out-of-range colour was read as an entry")
	}
}

func TestLoadOfAnAbsentFileIsNotAnError(t *testing.T) {
	reg, err := Load(t.TempDir())
	if err != nil {
		t.Fatalf("Load on a board with no tags: %v", err)
	}
	if names := reg.Names(); len(names) != 0 {
		t.Errorf("names = %v, want none", names)
	}
}

// The first tag on a board creates the file, and the file says what it is: the
// format is the API, and the person hand-editing it is the audience.
func TestSaveWritesTheHeaderOnceOnly(t *testing.T) {
	root := t.TempDir()
	reg, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	reg.Set("ui", 33)
	if err := reg.Save(root); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(Path(root))
	if !strings.HasPrefix(string(b), "# jaira tag colours") {
		t.Errorf("no header on a fresh file:\n%s", b)
	}

	reg2, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	reg2.Set("docs", 178)
	if err := reg2.Save(root); err != nil {
		t.Fatal(err)
	}
	b2, _ := os.ReadFile(Path(root))
	if n := strings.Count(string(b2), "# jaira tag colours"); n != 1 {
		t.Errorf("header written %d times:\n%s", n, b2)
	}
}

// Every new tag gets a colour no other tag is using, until the palette runs
// out. Two tags rendering identically is the one thing a colour cannot do.
func TestAssignGivesEveryPaletteColourBeforeRepeating(t *testing.T) {
	root := t.TempDir()
	reg, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	used := map[int]bool{}
	for i := range Palette {
		name := "tag" + string(rune('a'+i))
		c := reg.Assign(name)
		if used[c] {
			t.Fatalf("colour %d handed out twice, at tag %d of %d", c, i+1, len(Palette))
		}
		used[c] = true
		reg.Set(name, c)
	}
	if len(used) != len(Palette) {
		t.Errorf("%d distinct colours for %d palette entries", len(used), len(Palette))
	}

	// Past the palette, a repeat — but the same repeat on every machine and
	// every run, so the file stops changing instead of churning a random
	// colour per invocation.
	first := reg.Assign("one-too-many")
	if first != Fallback("one-too-many") {
		t.Errorf("exhausted palette gave %d, want the documented fallback %d", first, Fallback("one-too-many"))
	}
	for range 5 {
		if got := reg.Assign("one-too-many"); got != first {
			t.Errorf("fallback is not stable: %d then %d", first, got)
		}
	}
}

// A palette whose colours are not distinct would silently break the promise
// Assign makes.
func TestPaletteEntriesAreDistinctAnsi256Colours(t *testing.T) {
	seen := map[int]bool{}
	for _, c := range Palette {
		if !ValidColour(c) {
			t.Errorf("palette colour %d is not an ANSI-256 colour", c)
		}
		if seen[c] {
			t.Errorf("palette colour %d appears twice", c)
		}
		seen[c] = true
	}
	// The colours the board already spends on meaning: a tag swatch that reads
	// as a status is worse than an uncoloured tag.
	for _, reserved := range []int{39, 214, 203, 78, 141} {
		if seen[reserved] {
			t.Errorf("palette reuses %d, which the board already spends on status", reserved)
		}
	}
}

// Load lets the last entry for a name win, so Set must rewrite that same last
// line. Rewriting the first would change a line nothing reads and leave the
// colour exactly as it was — and duplicates are not hypothetical: merge=union on
// this file is what produces one when two sides recolour the same tag.
func TestSetRewritesTheLastDuplicateLine(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".jaira"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(Path(root), []byte("ui: 33\ndocs: 178\nui: 45\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	reg, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if c, _ := reg.Colour("ui"); c != 45 {
		t.Fatalf("Load read %d for a duplicated name, want the last line's 45", c)
	}
	reg.Set("ui", 100)
	if err := reg.Save(root); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(Path(root))
	if got := string(b); got != "ui: 33\ndocs: 178\nui: 100\n" {
		t.Errorf("Set rewrote the wrong line:\n%s", got)
	}
	again, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if c, _ := again.Colour("ui"); c != 100 {
		t.Errorf("colour after reload = %d, want 100", c)
	}
}

// Sorted insertion buys merge-friendliness and is only offered where the file
// already shows it is wanted. Sorting into a hand-grouped file puts a backend
// tag under a "# frontend" heading, which makes the comment above a line stop
// describing it — worse than an unsorted file.
func TestInsertionRespectsAHandGroupedFile(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".jaira"), 0o755); err != nil {
		t.Fatal(err)
	}
	grouped := "# frontend\nui: 33\ncss: 45\n\n# backend\nsql: 73\napi: 100\n"
	if err := os.WriteFile(Path(root), []byte(grouped), 0o644); err != nil {
		t.Fatal(err)
	}
	reg, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	reg.Set("docs", 208)
	if err := reg.Save(root); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(Path(root))
	if got := string(b); got != grouped+"docs: 208\n" {
		t.Errorf("an unsorted file was sorted into instead of appended to:\n%s", got)
	}

	// A sorted file still gets sorted insertion — that is what keeps two
	// teammates' additions off the same last line.
	root2 := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root2, ".jaira"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(Path(root2), []byte("api: 33\nui: 45\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	reg2, err := Load(root2)
	if err != nil {
		t.Fatal(err)
	}
	reg2.Set("docs", 208)
	if err := reg2.Save(root2); err != nil {
		t.Fatal(err)
	}
	b2, _ := os.ReadFile(Path(root2))
	if got := string(b2); got != "api: 33\ndocs: 208\nui: 45\n" {
		t.Errorf("a sorted file did not get sorted insertion:\n%s", got)
	}
}

// A reader must never see a half-written registry, and a crash mid-write must
// leave the previous file intact rather than a truncated one. That is what
// tmp-file-plus-rename buys; this pins that Save uses it, by watching the
// directory for the temporary file and by proving the old bytes survive a
// replacement.
func TestSaveWritesViaTempFileAndRename(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".jaira"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(Path(root), []byte("ui: 33\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Pin the file as it stands before Save by giving it a second name. A
	// rename replaces the directory entry Path(root) points at; the file that
	// was there keeps living under the link, so reading the link afterwards
	// still yields the old bytes, whereas an in-place truncate-and-write would
	// show the new ones under both names. Surviving old bytes are the property
	// the atomic write is bought for — a reader caught mid-write sees the
	// previous file whole — so pinning them pins the thing that matters rather
	// than a proxy for it. The link lives outside .jaira so the check further
	// down, that no temporary file is left behind, still sees an empty
	// directory beside the registry.
	//
	// Not os.SameFile over a pair of os.Stat results, which is the obvious way
	// to write this and is unfalsifiable on Windows: there os.Stat records only
	// the path and defers the file identity, which SameFile then loads by
	// reopening that path (os/types_windows.go:287). Both the before and the
	// after FileInfo resolve to whatever the path names at comparison time —
	// the same, already-renamed file — so the assertion passes whatever Save
	// did. It went green on linux and macos and red on windows-latest for
	// exactly that reason.
	//
	// And not an open handle held across the Save, whether from os.Open or from
	// os.Root: on Windows a handle on the rename target makes MoveFileEx fail
	// with "Access is denied", so Save itself dies and the test reports a
	// product fault that is not there. Delete sharing does not rescue it. A
	// hard link holds no handle at all, and NTFS supports it.
	link := filepath.Join(root, "tags-before")
	if err := os.Link(Path(root), link); err != nil {
		t.Fatal(err)
	}

	reg, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	reg.Set("docs", 178)
	if err := reg.Save(root); err != nil {
		t.Fatal(err)
	}

	old, err := os.ReadFile(link)
	if err != nil {
		t.Fatal(err)
	}
	if string(old) != "ui: 33\n" {
		t.Errorf("Save wrote in place: the file linked before Save reads %q, want the pre-Save bytes", old)
	}
	// And no temporary file is left behind in .jaira.
	ents, err := os.ReadDir(filepath.Join(root, ".jaira"))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range ents {
		if strings.HasPrefix(e.Name(), ".jaira-tmp-") {
			t.Errorf("Save left a temporary file behind: %s", e.Name())
		}
	}
	if b, _ := os.ReadFile(Path(root)); !strings.Contains(string(b), "docs: 178") {
		t.Errorf("the new entry is missing:\n%s", b)
	}
}

// Both filters — 'jaira list --tag' and the board's 'tag:' — go through this, so
// exactness is pinned once. A tag is a name from a closed vocabulary: "cur"
// answering with everything tagged "security" is a wrong answer, not a loose one.
func TestMatchesIsExactAndCaseFolded(t *testing.T) {
	have := []string{"security", "My UI"}
	for _, want := range []string{"security", "SECURITY", "my-ui", "My UI"} {
		if !Matches(have, want) {
			t.Errorf("Matches(%v, %q) = false, want true", have, want)
		}
	}
	for _, want := range []string{"cur", "ui", "sec", "securit", "", "front/end"} {
		if Matches(have, want) {
			t.Errorf("Matches(%v, %q) = true, want false", have, want)
		}
	}
}

// Suggest advises the repair Normalize refuses to make silently.
func TestSuggestRepairsWhatNormalizeRefuses(t *testing.T) {
	for raw, want := range map[string]string{
		"front/end": "front-end",
		"ui!":       "ui",
		"a.b.c":     "a-b-c",
	} {
		if _, _, err := Normalize(raw); err == nil {
			t.Errorf("Normalize(%q) no longer refuses; this test is about the ones it does", raw)
		}
		got, ok := Suggest(raw)
		if !ok || got != want {
			t.Errorf("Suggest(%q) = %q,%v, want %q,true", raw, got, ok, want)
		}
	}
	if got, ok := Suggest("///"); ok {
		t.Errorf("Suggest(%q) = %q,true, want no suggestion", "///", got)
	}
}

// Dashing a separator is a repair; dashing a letter is a deletion. The
// suggestion is printed as a runnable command, so "über" repaired to "ber" is a
// wrong tag one paste away — worse than no suggestion at all.
func TestSuggestWithdrawsRatherThanDeleteALetter(t *testing.T) {
	for _, raw := range []string{"über", "naïve", "日本語", "café"} {
		if _, _, err := Normalize(raw); err == nil {
			t.Errorf("Normalize(%q) no longer refuses; this test is about the ones it does", raw)
		}
		if got, ok := Suggest(raw); ok {
			t.Errorf("Suggest(%q) = %q,true — it deleted a letter", raw, got)
		}
	}
	// A separator is still repaired: the withdrawal is about letters, not about
	// giving up on every refused name.
	if got, ok := Suggest("front/end"); !ok || got != "front-end" {
		t.Errorf("Suggest(%q) = %q,%v, want front-end,true", "front/end", got, ok)
	}
	if got, ok := Suggest("ui!"); !ok || got != "ui" {
		t.Errorf("Suggest(%q) = %q,%v, want ui,true", "ui!", got, ok)
	}
}
