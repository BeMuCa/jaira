package cli

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/BeMuCa/jaira/core/milestone"
	"github.com/BeMuCa/jaira/core/ticket"
)

// suggestBoard is a board with the shipped lanes and no git, so a move, a
// logbook and a milestone file are all just files on disk.
func suggestBoard(t *testing.T) (string, *ticket.Store) {
	t.Helper()
	t.Setenv("JAIRA_USER", "berk")
	dir := t.TempDir()
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	return dir, s
}

func suggestTicket(t *testing.T, s *ticket.Store, status string) string {
	t.Helper()
	tk, err := s.Create(map[string]string{
		ticket.FieldID:    ticket.NewID(time.Date(2026, 10, 31, 12, 0, 0, 0, time.UTC)),
		ticket.FieldTitle: "t", ticket.FieldStatus: status,
	}, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	return tk.ID
}

// suggestMilestone writes the milestone file the way 'milestone create' and
// 'milestone add' leave it, without a ref.
func suggestMilestone(t *testing.T, dir, name string, ids ...string) {
	t.Helper()
	ms := milestone.New(name, 0, time.Date(2026, 10, 31, 12, 0, 0, 0, time.UTC))
	for _, id := range ids {
		ms.Add(id)
	}
	if err := ms.Save(dir); err != nil {
		t.Fatal(err)
	}
}

func TestMoveToSignoffSuggestsThePullRequest(t *testing.T) {
	dir, s := suggestBoard(t)
	id := suggestTicket(t, s, "in-progress")
	out, err := runCLI(t, dir, "move", ticket.Handle(id), "--to", "signoff", "--force")
	if err != nil {
		t.Fatalf("move: %v\n%s", err, out)
	}
	want := "Next: the pull request. Only a person opens it: /jaira-role-pr " + ticket.Handle(id)
	if !strings.Contains(out, want) {
		t.Errorf("output lacks %q:\n%s", want, out)
	}
}

func TestMoveJSONCarriesSuggestions(t *testing.T) {
	dir, s := suggestBoard(t)
	id := suggestTicket(t, s, "in-progress")
	out, err := runCLI(t, dir, "--json", "move", ticket.Handle(id), "--to", "signoff", "--force")
	if err != nil {
		t.Fatalf("move: %v\n%s", err, out)
	}
	var got struct {
		Suggestions []struct {
			Kind, By string
		} `json:"suggestions"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}
	if len(got.Suggestions) != 1 || got.Suggestions[0].Kind != "pull_request" || got.Suggestions[0].By != "person" {
		t.Errorf("suggestions = %+v, want one pull_request by person\n%s", got.Suggestions, out)
	}
}

func TestMoveIntoAnAgenticLaneSuggestsNothing(t *testing.T) {
	dir, s := suggestBoard(t)
	id := suggestTicket(t, s, "in-progress")
	out, err := runCLI(t, dir, "move", ticket.Handle(id), "--to", "critique", "--force")
	if err != nil {
		t.Fatalf("move: %v\n%s", err, out)
	}
	if strings.Contains(out, "Only a person opens it") {
		t.Errorf("an agentic lane suggested a pull request:\n%s", out)
	}
}

func TestTicketInAnOpenMilestoneSuggestsNoPullRequest(t *testing.T) {
	dir, s := suggestBoard(t)
	a := suggestTicket(t, s, "in-progress")
	b := suggestTicket(t, s, "todo")
	suggestMilestone(t, dir, "round-one", a, b)
	out, err := runCLI(t, dir, "move", ticket.Handle(a), "--to", "signoff", "--force")
	if err != nil {
		t.Fatalf("move: %v\n%s", err, out)
	}
	if strings.Contains(out, "/jaira-role-pr") {
		t.Errorf("a ticket in an open milestone got a pull request line:\n%s", out)
	}
	if !strings.Contains(out, "Milestone round-one: 0 of 2 done") {
		t.Errorf("no progress line:\n%s", out)
	}
}

func TestMoveCompletingAMilestoneSuggestsOnePullRequestAndFiling(t *testing.T) {
	dir, s := suggestBoard(t)
	a := suggestTicket(t, s, "signoff")
	b := suggestTicket(t, s, "done")
	suggestMilestone(t, dir, "round-one", a, b)
	out, err := runCLI(t, dir, "move", ticket.Handle(a), "--to", "done", "--force")
	if err != nil {
		t.Fatalf("move: %v\n%s", err, out)
	}
	for _, want := range []string{
		"Milestone round-one is complete. Next: one pull request for all its tickets. Only a person opens it: /jaira-role-pr " + ticket.Handle(a) + " " + ticket.Handle(b),
		"Once that pull request is merged, take milestone round-one off the board: jaira logbook round-one",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
}

func TestLogbookOfTheLastTicketSuggestsThePullRequest(t *testing.T) {
	dir, s := suggestBoard(t)
	a := suggestTicket(t, s, "done")
	b := suggestTicket(t, s, "done")
	suggestMilestone(t, dir, "round-one", a, b)
	out, err := runCLI(t, dir, "logbook", ticket.Handle(a))
	if err != nil {
		t.Fatalf("logbook: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Milestone round-one is complete") {
		t.Errorf("logbook of the last ticket suggested nothing:\n%s", out)
	}
}

func TestLogbookJSONCarriesSuggestions(t *testing.T) {
	dir, s := suggestBoard(t)
	a := suggestTicket(t, s, "done")
	suggestMilestone(t, dir, "round-one", a)
	out, err := runCLI(t, dir, "--json", "logbook", ticket.Handle(a))
	if err != nil {
		t.Fatalf("logbook: %v\n%s", err, out)
	}
	if !strings.Contains(out, `"suggestions"`) || !strings.Contains(out, `"file_milestone"`) {
		t.Errorf("logbook --json lacks the suggestions:\n%s", out)
	}
}

// The filing the line names is the one 'jaira logbook <milestone>' accepts.
func TestSuggestedFilingIsAccepted(t *testing.T) {
	dir, s := suggestBoard(t)
	a := suggestTicket(t, s, "done")
	suggestMilestone(t, dir, "round-one", a)
	if out, err := runCLI(t, dir, "logbook", "round-one"); err != nil {
		t.Fatalf("the suggested filing was refused: %v\n%s", err, out)
	}
}

func TestLsCarriesTheSuggestionsOfACompleteMilestone(t *testing.T) {
	dir, s := suggestBoard(t)
	a := suggestTicket(t, s, "done")
	suggestMilestone(t, dir, "round-one", a)
	suggestMilestone(t, dir, "empty")

	out, err := runCLI(t, dir, "milestone", "ls")
	if err != nil {
		t.Fatalf("ls: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Milestone round-one is complete") {
		t.Errorf("ls shows no line for the complete milestone:\n%s", out)
	}
	if strings.Contains(out, "Milestone empty is complete") {
		t.Errorf("an empty milestone was called complete:\n%s", out)
	}

	jout, err := runCLI(t, dir, "--json", "milestone", "ls")
	if err != nil {
		t.Fatalf("ls --json: %v\n%s", err, jout)
	}
	var got struct {
		Milestones []struct {
			Name        string
			Suggestions []map[string]any
		} `json:"milestones"`
	}
	if err := json.Unmarshal([]byte(jout), &got); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, jout)
	}
	for _, m := range got.Milestones {
		switch m.Name {
		case "round-one":
			if len(m.Suggestions) != 2 {
				t.Errorf("round-one suggestions = %v, want 2", m.Suggestions)
			}
		case "empty":
			if len(m.Suggestions) != 0 {
				t.Errorf("empty milestone suggestions = %v, want none", m.Suggestions)
			}
		}
	}
}
