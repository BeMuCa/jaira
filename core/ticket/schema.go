package ticket

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Frontmatter field names. These are the API — they appear in every ticket file
// and in git diffs, so they are stable and referenced by name everywhere rather
// than being spelled inline.
const (
	FieldID         = "id"
	FieldTitle      = "title"
	FieldStatus     = "status"
	FieldReady      = "ready"
	FieldCreator    = "creator"
	FieldAssignee   = "assignee"
	FieldGoal       = "goal"
	FieldContext    = "context"
	FieldDoD        = "definition-of-done"
	FieldBlockedBy  = "blocked-by"
	FieldCommits    = "commits"
	FieldModelTier  = "model-tier"
	FieldExecutedBy = "executed-by"
	FieldCreatedAt  = "created-at"
	FieldUpdatedAt  = "updated-at"
	FieldQuestion   = "question"
	FieldClaimedBy  = "claimed-by"
	FieldClaimedAt  = "claimed-at"

	// The outcome is three flat keys rather than a nested mapping. Nesting would
	// force the writer to rewrite an indented block to change one field, and the
	// merge driver resolves per top-level key — flat keys keep both operations
	// single-line, which is the whole point of the format.
	FieldOutcomeWhat     = "outcome-what"
	FieldOutcomeWhy      = "outcome-why"
	FieldOutcomeResolves = "outcome-resolves"

	// Reserved for a future external tracker adapter. jaira never interprets it.
	FieldExternal = "external"

	// FieldFollows links a follow-up ticket back to the one whose review produced
	// it, so the trail from "this was not finished" to "this is the rest" is not
	// lost the moment the first ticket closes.
	FieldFollows = "follows"

	// FieldBlockedReason is why a ticket is parked in the blocked lane. blocked-by
	// covers the case where the blocker is another ticket; this covers everything
	// else — a vendor, a decision, an environment — which is the majority and
	// which the board otherwise records nowhere. A blocked ticket whose reason is
	// unrecorded is the one a board loses: it looks the same on day one and on
	// day ninety.
	FieldBlockedReason = "blocked-reason"

	// FieldReviewVerdict is the reviewer's conclusion, kept separate from
	// outcome-resolves. The implementer writes what it believes it achieved and
	// the reviewer judges the diff; storing both in one field destroys the pair
	// the sign-off view exists to show.
	FieldReviewVerdict = "review-verdict"

	// FieldReviewSummary is what the change actually does, read off the diff.
	// The verdict alone is a judgement with no account of what was judged; a
	// later reader needs the summary to know what shipped without re-reading
	// the diff themselves.
	FieldReviewSummary = "review-summary"

	// FieldReviewGaps is what the reviewer found missing or unconvincing. It
	// must be filled even when there is nothing to report — "none" is a valid
	// answer, an empty field is not — because an empty field must always mean
	// "nobody looked", never "nothing found".
	FieldReviewGaps = "review-gaps"

	// FieldReviewCheck is how to check the change by hand: run this, then that,
	// then look at this. Every other review field is an account of what
	// happened; this one is the only thing the reader can act on, and without it
	// signing off means either trusting the account or reconstructing the steps
	// from the diff. Written for someone who knows none of what the reviewer
	// knows — see the writing register the create and note helps carry.
	FieldReviewCheck = "review-check"
)

// canonicalOrder is the order in which fields are written into a new ticket.
// Existing files are never reordered; this only shapes files jaira creates.
var canonicalOrder = []string{
	FieldID, FieldTitle, FieldStatus, FieldReady,
	FieldCreator, FieldAssignee, FieldExecutedBy,
	FieldGoal, FieldContext, FieldDoD,
	FieldBlockedBy, FieldBlockedReason, FieldFollows, FieldCommits, FieldModelTier,
	FieldOutcomeWhat, FieldOutcomeWhy, FieldOutcomeResolves,
	FieldQuestion, FieldClaimedBy, FieldClaimedAt,
	FieldCreatedAt, FieldUpdatedAt,
}

// SuppliedFields names the lane input-requires values a ticket already
// carries — or that the tool assembling a lane's bounded input derives
// itself — from the moment of creation, independent of any lane's
// output-produces. "diff" is not a frontmatter scalar (it is assembled from
// git commits at request time), so it is named here rather than with a Field
// constant.
//
// "plan" is deliberately not in this set, even though the Plan heading is
// seeded on every new ticket: an empty Plan checklist does not satisfy a
// lane's input-requires (see flow.go's showForLane), so a lane requiring
// "plan" genuinely depends on an earlier lane — pre-process — having
// produced one. Exempting it here would make it impossible for the
// load-order check below to ever catch a lane misordered relative to its
// own producer, which is the one case this check exists for.
//
// core/lane's load-time contract check compares a lane's input-requires
// against this set plus every earlier lane's output-produces, so this is the
// one place both packages read rather than two lists that could drift apart.
var SuppliedFields = []string{
	FieldTitle, FieldGoal, FieldContext, FieldDoD, FieldAssignee, "diff",
}

// Ticket is the decoded view of a ticket, used for reads, filtering and
// rendering. Writes never go through this struct — they go through Doc, so that
// unknown fields survive. Anything read into here is a projection, not the
// source of truth.
type Ticket struct {
	ID         string
	Title      string
	Status     string
	Ready      bool
	Creator    string
	Assignee   string
	ExecutedBy string
	Goal       string
	Context    string
	DoD        string
	BlockedBy  []string
	Commits    []string
	ModelTier  string
	Question   string
	ClaimedBy  string
	ClaimedAt  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// Follows is the ticket whose review produced this one.
	Follows string
	// BlockedReason is why the ticket is parked, when the blocker is not another
	// ticket.
	BlockedReason string
	// ReviewVerdict is the reviewer's conclusion, distinct from the
	// implementer's own account in Outcome.
	ReviewVerdict string
	// ReviewSummary is what the change does, in the reviewer's own words.
	ReviewSummary string
	// ReviewGaps is what the reviewer found missing; "none" is a valid value.
	ReviewGaps string
	// ReviewCheck is the steps a person follows to check the change themselves.
	ReviewCheck string

	Outcome Outcome

	// DoDItems are the checkboxes under a "Definition of Done" heading in the
	// body. A checklist is how people actually write acceptance criteria, and
	// keeping it in the body means it stays readable and editable prose rather
	// than being flattened into a single frontmatter string.
	DoDItems []DoDItem

	// PlanItems are the checkboxes under a "Plan" heading: the method being
	// followed rather than the criteria for acceptance. They carry the
	// in-progress marker, which is what makes "designed but not yet built"
	// visible instead of collapsing into a single in-progress lane.
	PlanItems []DoDItem

	// Options are the checkboxes under an "Options" heading, each naming a step
	// this ticket does or does not need. A lane can require one, which is how a
	// ticket skips planning or review without those lanes being removed.
	Options []DoDItem

	// Body is the markdown following the frontmatter.
	Body string
	// Path is where this ticket was loaded from.
	Path string
	// doc retains the original bytes so a mutation can be spliced back in.
	doc *Doc
}

// State is how far one checklist item has got.
type State int

const (
	// StateTodo is an item not started, or one whose marker is not recognised.
	// An unknown marker is read as unfinished rather than skipped: dropping it
	// would shorten the list and make the checklist easier to complete than the
	// author wrote it.
	StateTodo State = iota
	// StateDoing is the item currently being worked on, written as "[~]". It is
	// progress, not a completion claim, so it does not satisfy the criterion.
	StateDoing
	// StateDone is a finished item, written as "[x]".
	StateDone
	// StateSuperseded is an item that did not happen and will not, written as
	// "[-]". It is not a quieter [x]: the work was never done, so nothing may
	// report it as achieved. It exists because the only way to retire a
	// criterion used to be ticking it with the proof "obsolete, replaced by item
	// 6" — which leaves ticked items that were never met, and a tick that can
	// mean "never happened" means nothing anywhere on the board.
	StateSuperseded
)

func (s State) String() string {
	switch s {
	case StateDoing:
		return "doing"
	case StateDone:
		return "done"
	case StateSuperseded:
		return "superseded"
	default:
		return "todo"
	}
}

// Marker is the character this state is written with inside the brackets.
func (s State) Marker() string {
	switch s {
	case StateDoing:
		return "~"
	case StateDone:
		return "x"
	case StateSuperseded:
		return "-"
	default:
		return " "
	}
}

// DoDItem is one checkbox of a checklist, in either the Plan or the definition
// of done.
type DoDItem struct {
	Text  string
	State State
	// Proof is the evidence recorded for this item, written by SetItemProof and
	// read back off the plain indented line directly beneath the item's own —
	// never its own checkbox, since a checkbox there would be misread as
	// another criterion to satisfy rather than support for this one.
	Proof string
}

// Checked reports whether this item is finished. Only StateDone counts: an item
// in progress is outstanding work, and a superseded one was never done at all.
func (i DoDItem) Checked() bool { return i.State == StateDone }

// Settled reports whether this item is off the outstanding list — done, or
// retired as superseded. This is the question completion asks; Checked is the
// question "was this achieved", and the two are deliberately different for a
// superseded item.
func (i DoDItem) Settled() bool { return i.State == StateDone || i.State == StateSuperseded }

// HeadingTitle returns the first level-one heading of a markdown body.
func HeadingTitle(body string) string {
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return ""
}

// dodHeadingRE matches a "Definition of Done" heading in either language, since
// the section is what carries the meaning rather than its exact wording.
var dodHeadings = []string{"definition of done", "definition-of-done", "done when", "akzeptanzkriterien"}

// optionHeadings matches the section that turns steps on and off for one
// ticket. Not every ticket needs planning or review, and forcing every ticket
// through every step is how a pipeline becomes ceremony.
var optionHeadings = []string{"options", "steps needed", "optionen"}

// planHeadings matches the section holding the method — the steps taken to get
// there — as opposed to the criteria that decide whether it worked.
var planHeadings = []string{"plan", "steps", "vorgehen"}

// DoDHeadings and PlanHeadings expose the recognised section names so a renderer
// can tell which parts of a body it has already displayed as a checklist.
func DoDHeadings() []string    { return dodHeadings }
func OptionHeadings() []string { return optionHeadings }
func PlanHeadings() []string   { return planHeadings }

// ParseDoDItems extracts the checkboxes under a Definition of Done heading.
//
// Only that section is read: checkboxes elsewhere in a ticket (an open-questions
// list, say) are not acceptance criteria and must not be mistaken for them.
func ParseDoDItems(body string) []DoDItem { return checklistUnder(body, dodHeadings) }

// ParsePlanItems extracts the checkboxes under a Plan heading.
//
// The Plan is how the work is being done — write the spec, design it, implement
// it — and carries the in-progress marker. It deliberately does not gate the
// terminal lane: following a method is not the same as having met the criteria.
func ParsePlanItems(body string) []DoDItem { return checklistUnder(body, planHeadings) }

// ParseOptions extracts the checkboxes under an Options heading.
func ParseOptions(body string) []DoDItem { return checklistUnder(body, optionHeadings) }

func checklistUnder(body string, headings []string) []DoDItem {
	lines := strings.Split(body, "\n")
	inSection := false
	var out []DoDItem
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			heading := strings.ToLower(strings.TrimLeft(trimmed, "# "))
			inSection = false
			for _, h := range headings {
				if strings.Contains(heading, h) {
					inSection = true
					break
				}
			}
			continue
		}
		if !inSection {
			continue
		}
		item, ok := parseCheckbox(trimmed)
		if ok {
			out = append(out, item)
			continue
		}
		// A proof line is evidence for the item directly above it, not a
		// checkbox of its own — that is the whole reason the format was chosen:
		// parseCheckbox already drops any line that is not a bulleted "[ ]"/
		// "[x]"/"[~]" marker, so this can never be mistaken for another
		// criterion. Last one wins if a hand-edited file somehow carries two;
		// the writer guarantees one.
		if len(out) > 0 && isProofLine(trimmed) {
			out[len(out)-1].Proof = proofText(trimmed)
		}
	}
	return out
}

// parseCheckbox reads one list item. Any bracketed single character is a
// checkbox: " " and "x" have their usual meanings, "~" marks the item being
// worked on, and anything else is treated as unstarted rather than discarded.
func parseCheckbox(line string) (DoDItem, bool) {
	for _, bullet := range []string{"- ", "* ", "+ "} {
		if !strings.HasPrefix(line, bullet) {
			continue
		}
		rest := strings.TrimPrefix(line, bullet)
		if len(rest) < 4 || rest[0] != '[' || rest[2] != ']' || rest[3] != ' ' {
			continue
		}
		var st State
		switch rest[1] {
		case ' ':
			st = StateTodo
		case '~':
			st = StateDoing
		case 'x', 'X':
			st = StateDone
		case '-':
			st = StateSuperseded
		default:
			st = StateTodo
		}
		return DoDItem{Text: strings.TrimSpace(rest[4:]), State: st}, true
	}
	return DoDItem{}, false
}

// HasDoD reports whether a ticket declares a definition of done at all, in either
// supported form.
func (t *Ticket) HasDoD() bool {
	return len(t.DoDItems) > 0 || strings.TrimSpace(t.DoD) != ""
}

// DoDComplete reports whether every checklist item is finished.
//
// An item in progress counts as remaining. It is a statement that work is under
// way, which is the opposite of the criterion being met, so treating it as
// anything else would let a ticket close with its own checklist saying otherwise.
func (t *Ticket) DoDComplete() (complete bool, remaining []string) {
	if len(t.DoDItems) == 0 {
		return false, nil
	}
	for _, it := range t.DoDItems {
		if !it.Settled() {
			remaining = append(remaining, it.Text)
		}
	}
	return len(remaining) == 0, remaining
}

// OptionSet reports whether a named option is ticked on this ticket.
//
// Matching is on the text, case-insensitively and ignoring surrounding space, so
// "- [x] planning" and "- [x] Planning " both select the planning step. An
// option that is absent is not selected: a step is opt-in, so a ticket that says
// nothing about planning does not get planned.
func (t *Ticket) OptionSet(name string) bool {
	name = strings.ToLower(strings.TrimSpace(name))
	for _, o := range t.Options {
		if strings.ToLower(strings.TrimSpace(o.Text)) == name {
			return o.Checked()
		}
	}
	return false
}

// PlanComplete reports whether every step of the method has been carried out.
//
// A ticket with no plan is vacuously complete: the plan is optional, and its
// absence is not evidence of unfinished work.
func (t *Ticket) PlanComplete() (complete bool, remaining []string) {
	for _, it := range t.PlanItems {
		if !it.Settled() {
			remaining = append(remaining, it.Text)
		}
	}
	return len(remaining) == 0, remaining
}

// Outcome is the three-part close-out: what changed, why it was needed, and the
// causal link back to the Definition of Done. The third field exists because
// "what changed" alone does not let a reviewer judge whether the work is done.
type Outcome struct {
	What     string
	Why      string
	Resolves string
}

// Filled reports whether the outcome has enough substance to close a ticket.
func (o Outcome) Filled() bool {
	return strings.TrimSpace(o.What) != "" &&
		strings.TrimSpace(o.Why) != "" &&
		strings.TrimSpace(o.Resolves) != ""
}

// Doc exposes the underlying document for mutation.
func (t *Ticket) Doc() *Doc { return t.doc }

// Decode projects a parsed document into a Ticket.
func Decode(d *Doc, path string) (*Ticket, error) {
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
		// A hand-written ticket puts its title in the body's first heading rather
		// than in frontmatter. Reading it from there means such a file shows a
		// title on the board instead of a blank card.
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
	t.BlockedReason = str(FieldBlockedReason)
	t.ReviewVerdict = str(FieldReviewVerdict)
	t.ReviewSummary = str(FieldReviewSummary)
	t.ReviewGaps = str(FieldReviewGaps)
	t.ReviewCheck = str(FieldReviewCheck)
	t.DoDItems = ParseDoDItems(t.Body)
	t.PlanItems = ParsePlanItems(t.Body)
	t.Options = ParseOptions(t.Body)
	t.ModelTier = str(FieldModelTier)
	t.Question = str(FieldQuestion)
	t.ClaimedBy = str(FieldClaimedBy)
	t.BlockedBy = list(FieldBlockedBy)
	t.Commits = list(FieldCommits)
	t.Outcome = Outcome{
		What:     str(FieldOutcomeWhat),
		Why:      str(FieldOutcomeWhy),
		Resolves: str(FieldOutcomeResolves),
	}
	t.CreatedAt = parseTime(str(FieldCreatedAt))
	t.UpdatedAt = parseTime(str(FieldUpdatedAt))
	t.ClaimedAt = parseTime(str(FieldClaimedAt))

	if t.ID == "" {
		return nil, fmt.Errorf("ticket %s: missing required field %q", path, FieldID)
	}
	if t.Status == "" {
		return nil, fmt.Errorf("ticket %s: missing required field %q", path, FieldStatus)
	}
	return t, nil
}

func parseTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05Z", "2006-01-02"} {
		if tv, err := time.Parse(layout, s); err == nil {
			return tv
		}
	}
	return time.Time{}
}

// FormatTime renders a timestamp the way tickets store it.
func FormatTime(t time.Time) string { return t.UTC().Format(time.RFC3339) }

// NewDoc builds a fresh ticket document with fields in canonical order. Only
// non-empty fields are written, so a new ticket file stays readable rather than
// being padded with empty keys.
func NewDoc(fields map[string]string, lists map[string][]string, body string) *Doc {
	var b strings.Builder
	b.WriteString("---\n")

	written := map[string]bool{}
	emit := func(k string) {
		if written[k] {
			return
		}
		if v, ok := lists[k]; ok {
			b.WriteString(renderList(k, v, "") + "\n")
			written[k] = true
			return
		}
		if v, ok := fields[k]; ok && v != "" {
			switch {
			case rawLiteralFields[k]:
				b.WriteString(k + ": " + v + "\n")
			case blockable(v):
				b.WriteString(renderBlock(k, v, "", "") + "\n")
			default:
				b.WriteString(k + ": " + encodeScalar(v) + "\n")
			}
			written[k] = true
		}
	}
	for _, k := range canonicalOrder {
		emit(k)
	}
	// Any caller-supplied field not in the canonical list still gets written,
	// in sorted order for determinism.
	var extra []string
	for k := range fields {
		if !written[k] {
			extra = append(extra, k)
		}
	}
	for k := range lists {
		if !written[k] {
			extra = append(extra, k)
		}
	}
	sort.Strings(extra)
	for _, k := range extra {
		emit(k)
	}

	b.WriteString("---\n")
	if body != "" {
		if !strings.HasPrefix(body, "\n") {
			b.WriteString("\n")
		}
		b.WriteString(body)
		if !strings.HasSuffix(body, "\n") {
			b.WriteString("\n")
		}
	}
	d, err := ParseDoc([]byte(b.String()))
	if err != nil {
		// Constructed above from known-safe pieces; a failure here is a bug.
		panic("ticket: NewDoc produced an unparseable document: " + err.Error())
	}
	return d
}

// Touch stamps updated-at. Every mutation path calls this so recency-based
// conflict resolution has something to work with.
func Touch(d *Doc, now time.Time) error {
	return d.SetRaw(FieldUpdatedAt, FormatTime(now))
}

// SetReady records whether the promotion gate is currently satisfied. It is a
// derived convenience for rendering, never the authority: the gate is always
// recomputed before a move.
func SetReady(d *Doc, ready bool) error {
	return d.SetRaw(FieldReady, fmt.Sprintf("%t", ready))
}

// rawLiteralFields are written unquoted because their values are unambiguous
// YAML by construction.
var rawLiteralFields = map[string]bool{
	FieldReady: true, FieldCreatedAt: true, FieldUpdatedAt: true, FieldClaimedAt: true,
}
