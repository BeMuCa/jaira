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
)

// canonicalOrder is the order in which fields are written into a new ticket.
// Existing files are never reordered; this only shapes files jaira creates.
var canonicalOrder = []string{
	FieldID, FieldTitle, FieldStatus, FieldReady,
	FieldCreator, FieldAssignee, FieldExecutedBy,
	FieldGoal, FieldContext, FieldDoD,
	FieldBlockedBy, FieldCommits, FieldModelTier,
	FieldOutcomeWhat, FieldOutcomeWhy, FieldOutcomeResolves,
	FieldQuestion, FieldClaimedBy, FieldClaimedAt,
	FieldCreatedAt, FieldUpdatedAt,
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

	Outcome Outcome

	// Body is the markdown following the frontmatter.
	Body string
	// Path is where this ticket was loaded from.
	Path string
	// doc retains the original bytes so a mutation can be spliced back in.
	doc *Doc
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
	t.Status = str(FieldStatus)
	t.Ready = str(FieldReady) == "true"
	t.Creator = str(FieldCreator)
	t.Assignee = str(FieldAssignee)
	t.ExecutedBy = str(FieldExecutedBy)
	t.Goal = str(FieldGoal)
	t.Context = str(FieldContext)
	t.DoD = str(FieldDoD)
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
			if rawLiteralFields[k] {
				b.WriteString(k + ": " + v + "\n")
			} else {
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
