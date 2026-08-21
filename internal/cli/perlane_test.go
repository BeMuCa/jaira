package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/BeMuCa/jaira/core/gate"
	"github.com/BeMuCa/jaira/core/lane"
	"github.com/BeMuCa/jaira/core/ticket"
)

// perLaneOf runs emitPerLane over a set of tickets and returns the decoded JSON.
func perLaneOf(t *testing.T, tickets []*ticket.Ticket) map[string]any {
	t.Helper()
	lanes, err := lane.Load("")
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)
	old := g.jsonOut
	g.jsonOut = true
	defer func() { g.jsonOut = old }()

	if err := emitPerLane(cmd, gate.Env{Lanes: lanes}, tickets); err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("output is not JSON: %v\n%s", err, buf.String())
	}
	return out
}

func at(id, status string) *ticket.Ticket {
	return &ticket.Ticket{ID: id, Title: "t " + id, Status: status}
}

// The default "next" hands out the single furthest-along ticket, so a deep queue
// in a late lane hides every earlier lane completely. That is what --per-lane
// exists to answer, and it is the shape the real board showed: 28 waiting in one
// lane, 2 in another, and the second invisible.
func TestPerLaneShowsEveryLaneWithWork(t *testing.T) {
	var ts []*ticket.Ticket
	for i := range 28 {
		ts = append(ts, at("01KZTT3XZ2YQBX93TTSR7BVR"+string(rune('A'+i%26))+"T", "review"))
	}
	ts = append(ts, at("01KZTT3XZ2YQBX93TTSR7BVRC1", "in-progress"))
	ts = append(ts, at("01KZTT3XZ2YQBX93TTSR7BVRC2", "todo"))

	out := perLaneOf(t, ts)

	lanes, ok := out["lanes"].([]any)
	if !ok {
		t.Fatalf("lanes missing: %#v", out)
	}
	seen := map[string]float64{}
	for _, e := range lanes {
		m := e.(map[string]any)
		seen[m["lane"].(string)] = m["waiting"].(float64)
	}
	if len(seen) != 3 {
		t.Errorf("lanes reported = %v, want the three that have work", seen)
	}
	if seen["review"] != 28 {
		t.Errorf("review waiting = %v, want 28", seen["review"])
	}
	if seen["in-progress"] != 1 || seen["todo"] != 1 {
		t.Errorf("the smaller queues were not reported: %v", seen)
	}
}

// Pipeline order, not progress order: this is a map of the front line, so it
// reads left to right like the board does.
func TestPerLaneIsInPipelineOrder(t *testing.T) {
	out := perLaneOf(t, []*ticket.Ticket{
		at("01KZTT3XZ2YQBX93TTSR7BVRC1", "review"),
		at("01KZTT3XZ2YQBX93TTSR7BVRC2", "todo"),
		at("01KZTT3XZ2YQBX93TTSR7BVRC3", "in-progress"),
	})

	var got []string
	for _, e := range out["lanes"].([]any) {
		got = append(got, e.(map[string]any)["lane"].(string))
	}
	want := []string{"todo", "in-progress", "review"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("order = %v, want %v", got, want)
	}
}

// Each entry says whether an agent may work that lane at all, so the caller does
// not have to look the lane up separately.
func TestPerLaneCarriesTheLanesAgenticFlagAndATicket(t *testing.T) {
	out := perLaneOf(t, []*ticket.Ticket{
		at("01KZTT3XZ2YQBX93TTSR7BVRC1", "in-progress"),
		at("01KZTT3XZ2YQBX93TTSR7BVRC2", "todo"),
	})

	for _, e := range out["lanes"].([]any) {
		m := e.(map[string]any)
		tk, ok := m["ticket"].(map[string]any)
		if !ok || tk["handle"] == "" {
			t.Errorf("lane %v carries no ticket: %#v", m["lane"], m["ticket"])
		}
		switch m["lane"] {
		case "in-progress":
			if m["agentic"] != true {
				t.Error("in-progress is not reported as agentic")
			}
		case "todo":
			if m["agentic"] != false {
				t.Error("todo is reported as agentic")
			}
		}
	}
}

func TestPerLaneOnAnEmptyBoardSaysSo(t *testing.T) {
	out := perLaneOf(t, nil)
	if lanes, _ := out["lanes"].([]any); len(lanes) != 0 {
		t.Errorf("lanes = %v, want none", lanes)
	}
	if out["count"] != float64(0) {
		t.Errorf("count = %v, want 0", out["count"])
	}
}
