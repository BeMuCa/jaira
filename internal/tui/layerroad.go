package tui

import (
	"strconv"
	"strings"

	"github.com/BeMuCa/jaira/core/milestone"
	"github.com/BeMuCa/jaira/core/road"
)

// A swimlane's name is its road's title, and the roadmap shows it in the header
// over the lane. Reaching the header is ↑ from a top bubble; enter on it opens
// one field. A lane whose road exists is renamed. A lane no road names becomes
// a road holding every milestone drawn in it. Esc writes nothing.

// laneLabels is what each swimlane's header says: the title of the road that
// lists it, else Road-n, numbered 1, 2, … over the lanes no road names, left to
// right.
func (m *Model) laneLabels(lanes []milestone.Lane) []string {
	out := make([]string, len(lanes))
	n := 0
	for i, ln := range lanes {
		if ln.Road >= 0 && ln.Road < len(m.roads) {
			out[i] = m.roads[ln.Road].Title
			continue
		}
		n++
		out[i] = "Road-" + strconv.Itoa(n)
	}
	return out
}

// swimlaneIndex is the lane whose top has this ID, or -1.
func swimlaneIndex(lanes []milestone.Lane, top string) int {
	for i, ln := range lanes {
		if ln.Top == top {
			return i
		}
	}
	return -1
}

// startRoadName opens the one field behind enter on a header. It starts with
// the header's own text, Road-n for a lane nobody has named.
func (m *Model) startRoadName() {
	l := m.layer
	sc := m.layerScene()
	i := swimlaneIndex(sc.grid.Lanes, l.sel.id)
	if i < 0 {
		return
	}
	// editID holds the swimlane's top here: the road's name is not known yet
	// for a lane that has none.
	l.form = layerForm{labels: []string{"road"}, vals: []string{sc.labels[i]}, editID: l.sel.id}
	l.err = ""
	l.stage = stageRoadName
}

// keyRoadName drives the field behind a header.
func (m *Model) keyRoadName(s, text string) {
	l := m.layer
	switch s {
	case "esc":
		l.err = ""
		l.stage = stageLayer
	case "backspace":
		if r := []rune(l.form.vals[0]); len(r) > 0 {
			l.form.vals[0] = string(r[:len(r)-1])
		}
	case "enter":
		m.submitRoadName()
	default:
		m.insertText(text)
	}
}

// submitRoadName saves the name. A typed name that is blank is refused with
// nothing written; a lane with a road is renamed (title only), and a lane
// without one gets a road listing every milestone drawn in it.
func (m *Model) submitRoadName() {
	l := m.layer
	title := strings.TrimSpace(l.form.vals[0])
	if title == "" {
		l.err = "a road needs a name"
		return
	}
	sc := m.layerScene()
	i := swimlaneIndex(sc.grid.Lanes, l.form.editID)
	if i < 0 {
		l.err = "that swimlane is no longer on the roadmap"
		l.stage = stageLayer
		return
	}
	ln := sc.grid.Lanes[i]
	var err error
	if ln.Road >= 0 {
		_, _, err = road.Set(m.store, m.refSync, m.roads[ln.Road].Name, &title, nil)
	} else {
		var ids []string
		for _, b := range sc.grid.Bubbles {
			if b.Lane == i {
				ids = append(ids, b.ID)
			}
		}
		_, err = road.Create(m.store, m.refSync, road.NewName(m.roads, title), title, "", ids)
	}
	if err != nil {
		l.err = err.Error()
		return
	}
	l.err = ""
	l.stage = stageLayer
	if err := m.reload(); err != nil {
		m.notify(err.Error(), true)
	}
}
