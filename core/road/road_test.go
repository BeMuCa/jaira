package road

import (
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/ticket"
)

const (
	m1 = "ABC123"
	m2 = "DEF456"
	m3 = "GHK789"
)

// mapRefs keeps roads in memory, the way the outbox and the ref keep them.
type mapRefs map[string][]byte

func (m mapRefs) RoadNames() ([]string, error) {
	var out []string
	for name := range m {
		out = append(out, name)
	}
	return out, nil
}

func (m mapRefs) RoadContent(name string) ([]byte, error) {
	if c, ok := m[name]; ok {
		return c, nil
	}
	return nil, ErrNotFound
}

func (m mapRefs) RecordRoad(name string, content []byte) error {
	m[name] = content
	return nil
}

func store(t *testing.T) *ticket.Store {
	return &ticket.Store{Root: t.TempDir()}
}

// A hand edit survives jaira's edits: a comment and an order somebody chose must
// come back out of Add and Remove untouched.
func TestEditsKeepEveryOtherLine(t *testing.T) {
	src := "---\ntitle: Auth rollout\ndescription: one line\n---\n\n<!-- keep me -->\n\n- " + m1 + " the login\n- " + m2 + "\n"
	r := FromBytes("auth-rollout", []byte(src))
	if r.Title != "Auth rollout" || r.Description != "one line" {
		t.Fatalf("facts = %q / %q", r.Title, r.Description)
	}
	if got := r.Members(); !slices.Equal(got, []string{m1, m2}) {
		t.Fatalf("members = %v", got)
	}
	r.Remove(m2)
	r.Add(m3)
	want := "---\ntitle: Auth rollout\ndescription: one line\n---\n\n<!-- keep me -->\n\n- " + m1 + " the login\n- " + m3 + "\n"
	if got := string(r.Bytes()); got != want {
		t.Errorf("file after edits:\n%s\nwant:\n%s", got, want)
	}
}

// A member line is a milestone ID: a ticket ID's shape is not read, a lower-case
// ID is read as the ID it names, and prose with a dash is not a member.
func TestMemberLinesAreMilestoneIDs(t *testing.T) {
	r := FromBytes("x", []byte("- abc123 the first\n- a note about it\n- 01AAAAAAAAAAAAAAAAAAAAAAAA\n"))
	if got := r.Members(); !slices.Equal(got, []string{m1}) {
		t.Errorf("members = %v; want [%s]", got, m1)
	}
}

// The title is what a person reads; an empty one is the name, as for an epic.
func TestNewDefaultsTheTitleToTheName(t *testing.T) {
	r := FromBytes("auth-rollout", New("auth-rollout", "", "line one\nline two", time.Unix(0, 0)).Bytes())
	if r.Title != "auth-rollout" || r.Description != "line one line two" {
		t.Errorf("title %q, description %q", r.Title, r.Description)
	}
	if got := r.Label(); got != `"auth-rollout"` {
		t.Errorf("label %s; a road titled by its name is named once", got)
	}
}

// Create refuses a milestone another road lists and writes nothing for it.
func TestCreateRefusesAMilestoneInAnotherRoad(t *testing.T) {
	refs := mapRefs{}
	if _, err := Create(store(t), refs, "auth", "Auth", "", []string{m1}); err != nil {
		t.Fatalf("create: %v", err)
	}
	before := len(refs["auth"])
	_, err := Create(store(t), refs, "billing", "Billing", "", []string{m1, m2})
	var taken *TakenError
	if !errors.As(err, &taken) || taken.Milestone != m1 || taken.Road != "auth" {
		t.Fatalf("create = %v; want TakenError for %s in auth", err, m1)
	}
	if _, ok := refs["billing"]; ok {
		t.Error("a refused create wrote the road anyway")
	}
	if len(refs["auth"]) != before {
		t.Error("a refused create changed the road that lists the milestone")
	}
}

// Assign refuses a milestone another road lists, and adding one a road already
// lists is no change.
func TestAssignRefusesATakenMilestone(t *testing.T) {
	refs := mapRefs{}
	if _, err := Create(store(t), refs, "auth", "Auth", "", []string{m1}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := Create(store(t), refs, "billing", "Billing", "", []string{m2}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, _, err := Assign(store(t), refs, "billing", []string{m1}, true); err == nil {
		t.Fatal("assign of a milestone another road lists was accepted")
	}
	if got, _ := Load(refs, "billing"); !slices.Equal(got.Members(), []string{m2}) {
		t.Errorf("billing = %v after a refused assign", got.Members())
	}
	_, changed, err := Assign(store(t), refs, "auth", []string{m1}, true)
	if err != nil || len(changed) != 0 {
		t.Errorf("re-adding a listed milestone: changed %v, err %v; want none", changed, err)
	}
	_, changed, err = Assign(store(t), refs, "auth", []string{m1}, false)
	if err != nil || !slices.Equal(changed, []string{m1}) {
		t.Errorf("removing it: changed %v, err %v", changed, err)
	}
}

// Roads are read oldest first, by created-at and then name: that is the order a
// milestone listed in two roads is drawn in.
func TestLoadAllIsOldestFirst(t *testing.T) {
	refs := mapRefs{
		"b": New("b", "B", "", time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)).Bytes(),
		"a": New("a", "A", "", time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)).Bytes(),
		"c": New("c", "C", "", time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)).Bytes(),
	}
	all, err := LoadAll(refs)
	if err != nil {
		t.Fatalf("LoadAll: %v", err)
	}
	var names []string
	for _, r := range all {
		names = append(names, r.Name)
	}
	if !slices.Equal(names, []string{"c", "a", "b"}) {
		t.Errorf("order = %v; want [c a b]", names)
	}
}

// A milestone two roads list (a merge produced it) is reported with the older
// road first, and Owner gives it to that road.
func TestDuplicatesNameTheOlderRoadFirst(t *testing.T) {
	refs := mapRefs{
		"new": New("new", "New", "", time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)).Bytes(),
		"old": New("old", "Old", "", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)).Bytes(),
	}
	for name, id := range map[string]string{"new": m1, "old": m1} {
		r := FromBytes(name, refs[name])
		r.Add(id)
		refs[name] = r.Bytes()
	}
	all, _ := LoadAll(refs)
	dups := Duplicates(all)
	if len(dups) != 1 || dups[0].Milestone != m1 || !slices.Equal(dups[0].Roads, []string{"old", "new"}) {
		t.Fatalf("Duplicates = %+v; want %s in [old new]", dups, m1)
	}
	if owner := Owner(all); owner[m1] != 0 {
		t.Errorf("Owner(%s) = %d; want 0, the older road", m1, owner[m1])
	}
	if g := Groups(all); !slices.Equal(g[0], []string{m1}) || !slices.Equal(g[1], []string{m1}) {
		t.Errorf("Groups = %v", g)
	}
}

// Both sides of a lost race keep what they did: theirs listed B, ours listed C
// and took A out and retitled the road.
func TestMergeKeepsBothSidesMemberChanges(t *testing.T) {
	base := New("auth", "Old", "", time.Unix(0, 0))
	base.Add(m1)
	theirs := FromBytes("auth", base.Bytes())
	theirs.Add(m2)
	ours := FromBytes("auth", base.Bytes())
	ours.Add(m3)
	ours.Remove(m1)
	ours.Title = "New"
	ours.setField("title", "New")

	got := FromBytes("auth", Merge("auth", base.Bytes(), ours.Bytes(), theirs.Bytes()))
	if !slices.Equal(got.Members(), []string{m2, m3}) {
		t.Errorf("merged members = %v; want [%s %s]", got.Members(), m2, m3)
	}
	if got.Title != "New" {
		t.Errorf("merged title = %q; want the title ours changed", got.Title)
	}
}

// A blank title is the name, and a nil field is left as it is.
func TestSetKeepsNilFieldsAndBlankTitleIsTheName(t *testing.T) {
	refs := mapRefs{}
	if _, err := Create(store(t), refs, "auth", "Auth", "first", nil); err != nil {
		t.Fatalf("create: %v", err)
	}
	blank := "  "
	r, written, err := Set(store(t), refs, "auth", &blank, nil)
	if err != nil || !written || r.Title != "auth" || r.Description != "first" {
		t.Errorf("Set blank title = %q, %q, written %v, err %v; want the name and the description kept", r.Title, r.Description, written, err)
	}
	if _, written, _ := Set(store(t), refs, "auth", nil, nil); written {
		t.Error("Set with no fields wrote the road")
	}
}

// A new road's name is its title normalized, and never one already taken.
func TestNewNameIsUniqueAmongRoads(t *testing.T) {
	roads := []*Road{{Name: "auth-rollout"}}
	if got := NewName(roads, "Auth rollout"); got != "auth-rollout-2" {
		t.Errorf("NewName = %q; want auth-rollout-2", got)
	}
	if got := NewName(nil, "Billing"); got != "billing" {
		t.Errorf("NewName = %q; want billing", got)
	}
}
