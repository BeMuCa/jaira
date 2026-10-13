package ticket

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// decodeBefore is Decode as it was before the frontmatter was parsed once:
// Validate plus one Doc.Scalar/Doc.List call per field, each its own parse. It
// stays here as the reference the single-parse Decode must match.
func decodeBefore(d *Doc, path string) (*Ticket, error) {
	if err := d.Validate(); err != nil {
		return nil, err
	}
	t := &Ticket{Path: path, doc: d, Body: d.Body()}

	str := func(key string) string {
		v, _, err := d.Scalar(key)
		if err != nil {
			return ""
		}
		return v
	}
	list := func(key string) []string {
		v, err := d.List(key)
		if err != nil {
			return nil
		}
		return v
	}

	t.ID = str(FieldID)
	t.Title = str(FieldTitle)
	if t.Title == "" {
		t.Title = HeadingTitle(t.Body)
	}
	t.Status = str(FieldStatus)
	t.Ready = str(FieldReady) == "true"
	t.Creator = str(FieldCreator)
	t.Assignee = str(FieldAssignee)
	t.ExecutedBy = str(FieldExecutedBy)
	t.Goal = str(FieldGoal)
	t.Context = str(FieldContext)
	t.DoD = str(FieldDoD)
	t.Follows = str(FieldFollows)
	t.Parent = str(FieldParent)
	t.BlockedReason = str(FieldBlockedReason)
	t.ReviewVerdict = str(FieldReviewVerdict)
	t.ReviewSummary = str(FieldReviewSummary)
	t.ReviewGaps = str(FieldReviewGaps)
	t.ReviewCheck = str(FieldReviewCheck)
	t.Round = count(str(FieldRound))
	t.ReviewRound = count(str(FieldReviewRound))
	t.DoDItems = ParseDoDItems(t.Body)
	t.PlanItems = ParsePlanItems(t.Body)
	t.Options = ParseOptions(t.Body)
	t.ModelTier = str(FieldModelTier)
	t.Question = str(FieldQuestion)
	t.Mode = str(FieldMode)
	t.ClaimedBy = str(FieldClaimedBy)
	t.BlockedBy = list(FieldBlockedBy)
	t.Tags = list(FieldTags)
	t.Related = list(FieldRelated)
	t.Commits = list(FieldCommits)
	t.Outcome = Outcome{
		What:     str(FieldOutcomeWhat),
		Why:      str(FieldOutcomeWhy),
		Resolves: str(FieldOutcomeResolves),
	}
	t.CreatedAt = parseTime(str(FieldCreatedAt))
	t.UpdatedAt = parseTime(str(FieldUpdatedAt))
	t.UpdatedBy = str(FieldUpdatedBy)
	t.ClaimedAt = parseTime(str(FieldClaimedAt))

	if t.ID == "" {
		return nil, fmt.Errorf("ticket %s: missing required field %q", path, FieldID)
	}
	if t.Status == "" {
		return nil, fmt.Errorf("ticket %s: missing required field %q", path, FieldStatus)
	}
	return t, nil
}

// compareDecode checks Decode against decodeBefore on one source, and returns
// the counter differences Decode produced.
func compareDecode(t *testing.T, name string, src []byte) (decodes, parses int64) {
	t.Helper()
	d1, err1 := ParseDoc(src)
	d2, err2 := ParseDoc(src)
	if (err1 == nil) != (err2 == nil) {
		t.Fatalf("%s: ParseDoc disagrees: %v vs %v", name, err1, err2)
	}
	if err1 != nil {
		return 0, 0
	}
	want, wantErr := decodeBefore(d1, name)
	dec0, par0 := Counts()
	got, gotErr := Decode(d2, name)
	dec1, par1 := Counts()
	if fmt.Sprint(wantErr) != fmt.Sprint(gotErr) {
		t.Fatalf("%s: error differs\n before: %v\n after:  %v", name, wantErr, gotErr)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("%s: Decode differs from the pre-single-parse Decode\n before: %+v\n after:  %+v", name, want, got)
	}
	return dec1 - dec0, par1 - par0
}

// Decode parses a ticket's frontmatter exactly once and builds the same ticket
// as before, over every ticket file this repository has.
func TestDecodeParsesOnceAndMatchesBefore(t *testing.T) {
	var files []string
	for _, g := range []string{"../../.jaira/tickets/*.md", "../../.jaira/logbook/*/*.md"} {
		m, err := filepath.Glob(g)
		if err != nil {
			t.Fatal(err)
		}
		files = append(files, m...)
	}
	if len(files) < 1 {
		t.Fatal("no ticket files found to compare against")
	}
	for _, f := range files {
		src, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		dec, par := compareDecode(t, f, src)
		if dec != 1 || par != 1 {
			t.Errorf("%s: decodes=%d headerParses=%d, want 1 and 1", f, dec, par)
		}
	}
}

func TestDecodeMatchesBeforeOnEdgeCases(t *testing.T) {
	cases := map[string]string{
		"valid":           "---\nid: A1\nstatus: todo\ntitle: T\n---\nbody\n",
		"broken yaml":     "---\nid: [unclosed\nstatus: todo\n---\n",
		"anchor":          "---\nid: A1\nstatus: &s todo\n---\n",
		"alias":           "---\nid: A1\nstatus: todo\nx: *s\n---\n",
		"root not map":    "---\n- a\n- b\n---\n",
		"root scalar":     "---\njust text\n---\n",
		"single key":      "---\nid: A1\n---\n",
		"single status":   "---\nstatus: todo\n---\n",
		"block scalar":    "---\nid: A1\nstatus: todo\ncontext: |\n  one\n  two\n---\n",
		"null value":      "---\nid: A1\nstatus: todo\ngoal:\n---\n",
		"list for scalar": "---\nid: A1\nstatus: todo\ngoal:\n  - a\n---\n",
		"scalar for list": "---\nid: A1\nstatus: todo\ntags: x\n---\n",
		"map for scalar":  "---\nid: A1\nstatus: todo\ngoal:\n  a: b\n---\n",
		"no id":           "---\nstatus: todo\ntitle: T\n---\n",
		"no status":       "---\nid: A1\ntitle: T\n---\n",
		"empty":           "---\n---\n",
		"heading title":   "---\nid: A1\nstatus: todo\n---\n# From heading\n",
		"lists":           "---\nid: A1\nstatus: todo\ntags:\n  - a\n  - b\nblocked-by: [X, Y]\n---\n",
		"duplicate key":   "---\nid: A1\nid: B2\nstatus: todo\n---\n",
	}
	for name, src := range cases {
		compareDecode(t, name, []byte(src))
	}
}
