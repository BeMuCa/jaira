// Package suggest works out what jaira tells the person after a ticket or a
// milestone has reached an end: the pull request to ask for, and the filing
// that follows its merge. It only says what is owed. Nothing here opens a
// pull request or moves a ticket, and the lines are printed by the CLI and
// passed on by the roles as they are written.
//
// A lane counts by what it does, not by its name. Acceptance is a lane a
// person owns (requires-human-exit, no question lane). Done is the terminal
// lane, or no status at all once the ticket has been filed off the board.
package suggest

import (
	"fmt"
	"strings"

	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// Kinds of line. A progress line is shown but asks for nothing.
const (
	KindPullRequest   = "pull_request"
	KindFileMilestone = "file_milestone"
	KindProgress      = "progress"
)

// Member is one ticket of a milestone. Status is its lane, or "" when it is
// off the board (filed, or archived).
type Member struct {
	ID     string
	Status string
}

// Group is a milestone and the tickets it holds.
type Group struct {
	Name    string
	Filed   bool
	Members []Member
}

// Line is one thing jaira says. Text is what a person reads; the rest is what
// --json carries under 'suggestions'.
type Line struct {
	Kind    string `json:"kind"`
	Text    string `json:"text"`
	Command string `json:"command,omitempty"`
	By      string `json:"by,omitempty"`
	After   string `json:"after,omitempty"`
}

// ForTicket returns the lines owed once ticket id stands in status. An empty
// status means it is off the board. groups are the milestones the ticket may
// belong to; the entry for id itself is read at status, not from the group,
// because the group may still show the lane the ticket came from.
func ForTicket(lanes *lane.Set, id, status string, groups []Group) []Line {
	ready := done(lanes, status) || acceptance(lanes, status)
	var out []Line
	inOpen := false
	for _, g := range groups {
		if g.Filed || len(g.Members) == 0 || !g.has(id) {
			continue
		}
		inOpen = true
		if complete(lanes, g, id, status) {
			out = append(out, completeLines(g)...)
			continue
		}
		if ready {
			out = append(out, progressLine(lanes, g, id, status))
		}
	}
	if ready && !inOpen {
		// A ticket outside any open milestone is a group of one: it has no
		// other ticket whose pull request would ever carry it.
		h := ticket.Handle(id)
		out = append(out, Line{
			Kind:    KindPullRequest,
			Text:    "Next: the pull request. Only a person opens it: /jaira-role-pr " + h,
			Command: "/jaira-role-pr " + h,
			By:      "person",
		})
	}
	return out
}

// ForMilestone returns the lines for a milestone that is complete, and none
// otherwise. An empty or filed milestone is never complete.
func ForMilestone(lanes *lane.Set, g Group) []Line {
	if !complete(lanes, g, "", "") {
		return nil
	}
	return completeLines(g)
}

// complete: every member is done. A milestone with no members is not.
func complete(lanes *lane.Set, g Group, id, status string) bool {
	if g.Filed || len(g.Members) == 0 {
		return false
	}
	for _, m := range g.Members {
		if !done(lanes, memberStatus(m, id, status)) {
			return false
		}
	}
	return true
}

func completeLines(g Group) []Line {
	handles := make([]string, 0, len(g.Members))
	for _, m := range g.Members {
		handles = append(handles, ticket.Handle(m.ID))
	}
	all := strings.Join(handles, " ")
	return []Line{
		{
			Kind:    KindPullRequest,
			Text:    fmt.Sprintf("Milestone %s is complete. Next: one pull request for all its tickets. Only a person opens it: /jaira-role-pr %s", g.Name, all),
			Command: "/jaira-role-pr " + all,
			By:      "person",
		},
		{
			Kind:    KindFileMilestone,
			Text:    fmt.Sprintf("Once that pull request is merged, take milestone %s off the board: jaira logbook %s", g.Name, g.Name),
			Command: "jaira logbook " + g.Name,
			After:   "merge",
		},
	}
}

func progressLine(lanes *lane.Set, g Group, id, status string) Line {
	n := 0
	for _, m := range g.Members {
		if done(lanes, memberStatus(m, id, status)) {
			n++
		}
	}
	return Line{
		Kind: KindProgress,
		Text: fmt.Sprintf("Milestone %s: %d of %d done", g.Name, n, len(g.Members)),
	}
}

func (g Group) has(id string) bool {
	for _, m := range g.Members {
		if m.ID == id {
			return true
		}
	}
	return false
}

func memberStatus(m Member, id, status string) string {
	if id != "" && m.ID == id {
		return status
	}
	return m.Status
}

// find walks Lanes: the Set's index is private, and a walk over a list this
// short costs nothing.
func find(lanes *lane.Set, id string) *lane.Lane {
	for _, l := range lanes.Lanes {
		if l.ID == id {
			return l
		}
	}
	return nil
}

func done(lanes *lane.Set, status string) bool {
	if status == "" {
		return true
	}
	l := find(lanes, status)
	return l != nil && l.Terminal
}

// acceptance is a lane a person owns: no agent may move a ticket out of it, and
// it is not a question lane, which a person only answers. Same reading as
// impliedFresh in core/lane.
func acceptance(lanes *lane.Set, status string) bool {
	l := find(lanes, status)
	return l != nil && l.RequiresHumanExit && !l.RequiresQuestion && !l.Agentic
}
