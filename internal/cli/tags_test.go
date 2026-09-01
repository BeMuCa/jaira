package cli

import (
	"errors"
	"os"
	"strings"
	"testing"

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
	h := ""
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	all, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	for _, tk := range all {
		if tk.Title == "first" {
			h = ticket.Handle(tk.ID)
		}
	}
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
	// A coloured tag shows its number, so the file can be hand-edited from
	// what the listing says.
	if row := tagsRow(t, out, "ui"); !strings.ContainsAny(row, "0123456789") {
		t.Errorf("a coloured tag does not show its colour: %q", row)
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
