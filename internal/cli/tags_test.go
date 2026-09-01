package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/tag"
	"github.com/BeMuCa/jaira/core/ticket"
)

// firstHandle is the handle of the board's only ticket, so a test can create
// one and then act on it the way a person would.
func firstHandle(t *testing.T, dir string) string {
	t.Helper()
	return ticket.Handle(createdTicket(t, dir).ID)
}

// tagsRow returns the 'jaira tags' line for one name. The colours are chosen at
// random from the palette, so a test asserts the row's shape rather than pinning
// a number that is deliberately not deterministic.
func tagsRow(t *testing.T, out, name string) string {
	t.Helper()
	for _, l := range strings.Split(out, "\n") {
		for _, f := range strings.Fields(l) {
			if f == name {
				return l
			}
		}
	}
	t.Fatalf("no row for tag %q in:\n%s", name, out)
	return ""
}

// firstTitledHandle finds the handle of the board's ticket with this title.
func firstTitledHandle(t *testing.T, dir, title string) string {
	t.Helper()
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	all, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	for _, tk := range all {
		if tk.Title == title {
			return ticket.Handle(tk.ID)
		}
	}
	t.Fatalf("no ticket titled %q on the board", title)
	return ""
}

// jsonCLI runs a command with --json and decodes the payload, so a test asserts
// against the machine surface an agent actually reads rather than against prose.
func jsonCLI(t *testing.T, dir string, args ...string) map[string]any {
	t.Helper()
	out, err := runCLI(t, dir, append(args, "--json")...)
	if err != nil {
		t.Fatalf("%v: %v\n%s", args, err, out)
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("%v emitted unparseable JSON: %v\n%s", args, err, out)
	}
	return m
}

// strList reads a JSON string array, treating a missing key as absent rather
// than as an empty list, so a test can tell the two apart.
func strList(v any) []string {
	arr, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(arr))
	for _, e := range arr {
		s, _ := e.(string)
		out = append(out, s)
	}
	return out
}

func tagsFile(t *testing.T, dir string) string {
	t.Helper()
	b, err := os.ReadFile(tag.Path(dir))
	if err != nil {
		t.Fatalf("reading %s: %v", tag.Path(dir), err)
	}
	return string(b)
}

// A new name is new and says so; the same name again is reused and says that.
// The distinction is the whole guard against a board growing one synonym per
// session — an agent that is told "new" when it expected "reused" has just
// invented one.
func TestTagReportsNewThenReused(t *testing.T) {
	dir := emptyStore(t)
	if out, err := runCLI(t, dir, "create", "something"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	h := firstHandle(t, dir)

	out, err := runCLI(t, dir, "tag", h, "ui")
	if err != nil {
		t.Fatalf("tag: %v\n%s", err, out)
	}
	if !strings.Contains(out, "New: ui") {
		t.Errorf("a first-ever tag was not reported as new:\n%s", out)
	}
	if !strings.Contains(out, "jaira tags") {
		t.Errorf("a new tag did not point at the vocabulary listing:\n%s", out)
	}

	out, err = runCLI(t, dir, "tag", h, "ui", "backend")
	if err != nil {
		t.Fatalf("tag: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Reused: ui") {
		t.Errorf("a known tag was not reported as reused:\n%s", out)
	}
	if !strings.Contains(out, "New: backend") {
		t.Errorf("a second new tag was not reported as new:\n%s", out)
	}

	tk := createdTicket(t, dir)
	if strings.Join(tk.Tags, ",") != "ui,backend" {
		t.Errorf("tags = %v, want [ui backend] with no repeat", tk.Tags)
	}
	// Both names have a colour line, and the file carries its own header.
	got := tagsFile(t, dir)
	for _, want := range []string{"# jaira tag colours", "ui:", "backend:"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in .jaira/tags:\n%s", want, got)
		}
	}
}

// A name is stored in the one form the board uses, and the person is told —
// a tag filed under a name you did not type is exactly how a second name for
// one subject appears.
func TestTagNormalizesAndSaysSo(t *testing.T) {
	dir := emptyStore(t)
	if out, err := runCLI(t, dir, "create", "something"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	h := firstHandle(t, dir)

	out, err := runCLI(t, dir, "tag", h, "My UI")
	if err != nil {
		t.Fatalf("tag: %v\n%s", err, out)
	}
	if !strings.Contains(out, `Filed "My UI" as "my-ui"`) {
		t.Errorf("the rename was not reported:\n%s", out)
	}
	if tk := createdTicket(t, dir); strings.Join(tk.Tags, ",") != "my-ui" {
		t.Errorf("tags = %v, want [my-ui]", tk.Tags)
	}

	// Refused rather than trimmed down: a silently shortened name is a second
	// name for one subject.
	out, err = runCLI(t, dir, "tag", h, "front/end")
	if err == nil {
		t.Errorf("an unstorable name was accepted:\n%s", out)
	}
	var ce *codedError
	if !errors.As(err, &ce) || ce.code != ExitUsage {
		t.Errorf("error = %v, want exit code %d for a bad tag name", err, ExitUsage)
	}
}

func TestTagColorFlag(t *testing.T) {
	dir := emptyStore(t)
	if out, err := runCLI(t, dir, "create", "something"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	h := firstHandle(t, dir)

	if out, err := runCLI(t, dir, "tag", h, "ui", "--color", "45"); err != nil {
		t.Fatalf("tag --color: %v\n%s", err, out)
	}
	if got := tagsFile(t, dir); !strings.Contains(got, "ui: 45") {
		t.Errorf("--color did not set the colour:\n%s", got)
	}
	// An explicit colour on a tag that already has one recolours it: refusing
	// would leave hand-editing the file as the only way.
	if out, err := runCLI(t, dir, "tag", h, "ui", "--color", "100"); err != nil {
		t.Fatalf("tag --color again: %v\n%s", err, out)
	}
	got := tagsFile(t, dir)
	if !strings.Contains(got, "ui: 100") || strings.Contains(got, "ui: 45") {
		t.Errorf("--color did not recolour an existing tag:\n%s", got)
	}

	// One colour cannot be meant for several names, and 300 is not a colour.
	if out, err := runCLI(t, dir, "tag", h, "a", "b", "--color", "45"); err == nil {
		t.Errorf("--color with two names was accepted:\n%s", out)
	}
	if out, err := runCLI(t, dir, "tag", h, "a", "--color", "300"); err == nil {
		t.Errorf("--color 300 was accepted:\n%s", out)
	}
}

// The listing is what an agent reads before tagging, so it has to show every
// name the board knows — including one set by hand, which no palette ever
// coloured.
func TestTagsListsCountsAndUncolouredNames(t *testing.T) {
	dir := emptyStore(t)
	if out, err := runCLI(t, dir, "create", "first", "--tag", "ui"); err != nil {
		t.Fatalf("create --tag: %v\n%s", err, out)
	}
	if out, err := runCLI(t, dir, "create", "second", "--tag", "ui", "--tag", "docs"); err != nil {
		t.Fatalf("create --tag: %v\n%s", err, out)
	}
	// Set by hand through the generic list-field writer, which is what makes
	// tags a plain list field worth having: no colour line is created.
	h := firstTitledHandle(t, dir, "first")
	if out, err := runCLI(t, dir, "set", h, "tags=ui,handset"); err != nil {
		t.Fatalf("set tags=: %v\n%s", err, out)
	}

	out, err := runCLI(t, dir, "tags")
	if err != nil {
		t.Fatalf("tags: %v\n%s", err, out)
	}
	for name, wantOpen := range map[string]string{"ui": "2", "docs": "1", "handset": "1"} {
		if row := tagsRow(t, out, name); !strings.Contains(row, wantOpen+" open") {
			t.Errorf("row for %q does not say %s open: %q", name, wantOpen, row)
		}
	}
	// A name nothing coloured still gets a row, with the colour column empty
	// rather than the row missing: that is the least-magic answer to "who
	// colours a tag somebody set by hand" — nobody, until 'jaira tag' does.
	if row := tagsRow(t, out, "handset"); !strings.Contains(row, "-") || !strings.Contains(row, "no colour yet") {
		t.Errorf("a hand-set tag was not shown as uncoloured: %q", row)
	}
	// A coloured tag shows the number the file actually holds, so the registry
	// can be hand-edited from what the listing says. Read the colour out of the
	// file rather than pinning one: the palette choice is random on purpose.
	reg, err := tag.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	colour, ok := reg.Colour("ui")
	if !ok {
		t.Fatal("ui has no colour in .jaira/tags")
	}
	if row := tagsRow(t, out, "ui"); !strings.Contains(row, fmt.Sprintf("%3d", colour)) {
		t.Errorf("row for ui does not carry its colour %d: %q", colour, row)
	}
	if !strings.Contains(out, "synonym") {
		t.Errorf("the listing does not tell the reader to reuse a name:\n%s", out)
	}
	// A read command writes nothing: a listing that edited the board would
	// race every other session reading it.
	before := tagsFile(t, dir)
	if _, err := runCLI(t, dir, "tags"); err != nil {
		t.Fatal(err)
	}
	if after := tagsFile(t, dir); after != before {
		t.Errorf("'jaira tags' rewrote the registry:\nbefore:\n%s\nafter:\n%s", before, after)
	}
	if strings.Contains(before, "handset") {
		t.Errorf("a hand-set tag was given a colour by a read command:\n%s", before)
	}
}

// "open" has to mean outstanding work, or the number is not worth printing.
func TestTagsCountSkipsTheTerminalLane(t *testing.T) {
	dir := emptyStore(t)
	if out, err := runCLI(t, dir, "create", "shipped", "--lane", "done", "--tag", "ui"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	out, err := runCLI(t, dir, "tags")
	if err != nil {
		t.Fatalf("tags: %v\n%s", err, out)
	}
	if !strings.Contains(out, "0 open") {
		t.Errorf("a ticket in the terminal lane was counted as open:\n%s", out)
	}
	// The name is still listed: a tag whose last ticket closed is still the
	// name to reuse for the next one.
	if !strings.Contains(out, "ui") {
		t.Errorf("the tag disappeared with its last open ticket:\n%s", out)
	}
}

func TestListFiltersByTag(t *testing.T) {
	dir := emptyStore(t)
	if out, err := runCLI(t, dir, "create", "a ui thing", "--tag", "ui"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	if out, err := runCLI(t, dir, "create", "a docs thing", "--tag", "docs"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}

	out, err := runCLI(t, dir, "list", "--tag", "ui")
	if err != nil {
		t.Fatalf("list --tag: %v\n%s", err, out)
	}
	if !strings.Contains(out, "a ui thing") {
		t.Errorf("--tag ui did not list the ui ticket:\n%s", out)
	}
	if strings.Contains(out, "a docs thing") {
		t.Errorf("--tag ui listed a ticket tagged docs:\n%s", out)
	}

	out, err = runCLI(t, dir, "list", "--tag", "nothing-carries-this")
	if err != nil {
		t.Fatalf("list --tag: %v\n%s", err, out)
	}
	if !strings.Contains(out, "No tickets match.") {
		t.Errorf("a tag nothing carries still matched:\n%s", out)
	}
}

// 'jaira show' prints the row too, so the field is visible without the TUI.
func TestShowPrintsTheTagsRow(t *testing.T) {
	dir := emptyStore(t)
	if out, err := runCLI(t, dir, "create", "something", "--tag", "ui", "--tag", "backend"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	out, err := runCLI(t, dir, "show", firstHandle(t, dir))
	if err != nil {
		t.Fatalf("show: %v\n%s", err, out)
	}
	if !strings.Contains(out, "tags       ui backend") {
		t.Errorf("no tags row in 'jaira show':\n%s", out)
	}
}

// The listing is what an agent reads, and SKILL.md tells it to read --json. The
// payload has to carry the same three facts the table does.
func TestTagsJSONCarriesNameColorAndCount(t *testing.T) {
	dir := emptyStore(t)
	if out, err := runCLI(t, dir, "create", "first", "--tag", "ui"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	if out, err := runCLI(t, dir, "create", "second", "--tag", "ui"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	h := firstTitledHandle(t, dir, "first")
	if out, err := runCLI(t, dir, "set", h, "tags=ui,handset"); err != nil {
		t.Fatalf("set: %v\n%s", err, out)
	}

	out := jsonCLI(t, dir, "tags")
	rows, ok := out["tags"].([]any)
	if !ok {
		t.Fatalf("tags is not an array: %#v", out["tags"])
	}
	if n, _ := out["count"].(float64); int(n) != len(rows) {
		t.Errorf("count = %v, want %d", out["count"], len(rows))
	}
	byName := map[string]map[string]any{}
	for _, r := range rows {
		m, _ := r.(map[string]any)
		name, _ := m["name"].(string)
		byName[name] = m
	}
	reg, err := tag.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	want, _ := reg.Colour("ui")
	ui, found := byName["ui"]
	if !found {
		t.Fatalf("no ui row in %#v", byName)
	}
	// "color", not "colour": one spelling on the machine surface, matching the
	// --color flag.
	if got, _ := ui["color"].(float64); int(got) != want {
		t.Errorf("ui color = %v, want %d", ui["color"], want)
	}
	if _, wrongKey := ui["colour"]; wrongKey {
		t.Error("the payload still carries the old 'colour' key")
	}
	if got, _ := ui["open"].(float64); int(got) != 2 {
		t.Errorf("ui open = %v, want 2", ui["open"])
	}
	hand, found := byName["handset"]
	if !found {
		t.Fatalf("a hand-set tag is missing from the payload: %#v", byName)
	}
	if hand["color"] != nil {
		t.Errorf("an uncoloured tag reports a colour: %#v", hand["color"])
	}
}

// The new-versus-reused signal is what makes an invented synonym visible, and
// --json is the surface the skill tells agents to drive. It has to be there.
func TestTagJSONCarriesNewAndReused(t *testing.T) {
	dir := emptyStore(t)
	if out, err := runCLI(t, dir, "create", "something"); err != nil {
		t.Fatalf("create: %v\n%s", err, out)
	}
	h := firstHandle(t, dir)

	out := jsonCLI(t, dir, "tag", h, "ui")
	if got := strList(out["tags_new"]); len(got) != 1 || got[0] != "ui" {
		t.Errorf("tags_new = %#v, want [ui]", out["tags_new"])
	}
	if got := strList(out["tags_reused"]); len(got) != 0 {
		t.Errorf("tags_reused = %#v, want empty", out["tags_reused"])
	}
	// Empty arrays, never null: an agent branching on length must not have to
	// handle both.
	if out["tags_reused"] == nil {
		t.Error("tags_reused is null rather than an empty array")
	}
	if got := strList(out["tags"]); len(got) != 1 || got[0] != "ui" {
		t.Errorf("tags = %#v, want [ui]", out["tags"])
	}

	out = jsonCLI(t, dir, "tag", h, "ui", "backend")
	if got := strList(out["tags_reused"]); len(got) != 1 || got[0] != "ui" {
		t.Errorf("tags_reused = %#v, want [ui]", out["tags_reused"])
	}
	if got := strList(out["tags_new"]); len(got) != 1 || got[0] != "backend" {
		t.Errorf("tags_new = %#v, want [backend]", out["tags_new"])
	}
}

// Two sessions tagging at the same moment both read the file, both add their own
// line to what they read, and the second save writes a file assembled before the
// first one's line existed — one tag lost, silently. The reviewer proved that.
// The lock is what orders them; this is the regression test.
func TestRegisterTagsLosesNoTagUnderConcurrentWriters(t *testing.T) {
	dir := emptyStore(t)
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	const writers = 8
	var wg sync.WaitGroup
	errs := make(chan error, writers)
	for i := range writers {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			name := fmt.Sprintf("tag%d", i)
			if _, _, err := registerTags(s, []string{name}, -1); err != nil {
				errs <- err
			}
		}(i)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatalf("registerTags: %v", err)
	}

	reg, err := tag.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	for i := range writers {
		name := fmt.Sprintf("tag%d", i)
		if _, ok := reg.Colour(name); !ok {
			t.Errorf("%s was lost; the file holds %v", name, reg.Names())
		}
	}
}

// And the ordering really is a lock, not luck: while the store's "tags" lock is
// held, registerTags waits rather than reading a file somebody else is about to
// replace.
func TestRegisterTagsWaitsForTheTagsLock(t *testing.T) {
	dir := emptyStore(t)
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	unlock, err := s.Lock(tagsLockName)
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan error, 1)
	go func() {
		_, _, err := registerTags(s, []string{"ui"}, -1)
		done <- err
	}()

	select {
	case err := <-done:
		unlock()
		t.Fatalf("registerTags did not wait for the tags lock (returned %v)", err)
	case <-time.After(150 * time.Millisecond):
	}

	unlock()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("registerTags after unlock: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("registerTags never completed after the lock was released")
	}
	if reg, err := tag.Load(dir); err != nil {
		t.Fatal(err)
	} else if _, ok := reg.Colour("ui"); !ok {
		t.Error("ui was never written")
	}
}

// The tag registry is one independent line per tag, so git's own union driver is
// right for it — and a board configured before the registry existed already
// names the ticket driver, so a single "is anything configured" test would leave
// it without the union line for good.
func TestGitAttributesGainsTheUnionLineOnAnOlderBoard(t *testing.T) {
	dir := emptyStore(t)
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, ticket.DirName, ".gitattributes")
	old := "tickets/*.md merge=" + mergeDriverName + "\n"
	if err := os.WriteFile(path, []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}

	changed, err := writeGitAttributes(s)
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("an older board's attributes file was left without the union line")
	}
	b, _ := os.ReadFile(path)
	got := string(b)
	if !strings.Contains(got, "tags merge=union") {
		t.Errorf("no union line for the registry:\n%s", got)
	}
	if strings.Count(got, "merge="+mergeDriverName) != 1 {
		t.Errorf("the ticket driver line was duplicated:\n%s", got)
	}

	// Idempotent: a second run has nothing left to add.
	changed, err = writeGitAttributes(s)
	if err != nil {
		t.Fatal(err)
	}
	if changed {
		t.Errorf("writeGitAttributes rewrote a file that already said everything:\n%s", got)
	}
}
