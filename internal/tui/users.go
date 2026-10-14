package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/BeMuCa/jaira/core/ticket"
)

// userRow is one person in the user picker: the name as the tickets spell
// it, and how many loaded tickets it is assigned or wrote.
type userRow struct {
	name  string
	count int
}

// openUsers opens the user picker. The names come off every loaded ticket,
// not only the ones the filter lets through: a person filtered out has to be
// in the list to be ticked back in. The most tickets first, so the people
// who carry this board are at the top.
func (m *Model) openUsers() {
	at := map[string]int{}
	var rows []userRow
	add := func(t *ticket.Ticket) {
		seen := map[string]bool{}
		for _, name := range []string{t.Assignee, t.Creator} {
			name = strings.TrimSpace(name)
			k := strings.ToLower(name)
			if name == "" || seen[k] {
				continue
			}
			seen[k] = true
			if i, ok := at[k]; ok {
				rows[i].count++
				continue
			}
			at[k] = len(rows)
			rows = append(rows, userRow{name: name, count: 1})
		}
	}
	for _, t := range m.tickets {
		add(t)
	}
	for _, l := range m.logged {
		add(l.Ticket)
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].count != rows[j].count {
			return rows[i].count > rows[j].count
		}
		return strings.ToLower(rows[i].name) < strings.ToLower(rows[j].name)
	})

	m.users = rows
	m.openTicks(modeUsers, "user")
}

// userNames is the user picker's rows as the names they show.
func (m *Model) userNames() []string {
	out := make([]string, len(m.users))
	for i, u := range m.users {
		out[i] = u.name
	}
	return out
}

// renderUsers draws the user picker.
func (m *Model) renderUsers() string {
	var b strings.Builder
	b.WriteString(styLaneTitle.Render("Users") + "\n")
	b.WriteString(styBar.Render(strings.Repeat("─", min(m.width, 40))) + "\n\n")
	if len(m.users) == 0 {
		b.WriteString(styMeta.Render("No ticket on this board names anyone yet.") + "\n")
	}
	for i, u := range m.users {
		marker, name := "  ", u.name
		if i == m.pick.idx {
			marker, name = stySelected.Render("▌ "), stySelected.Render(name)
		}
		b.WriteString(marker + m.tickMark(u.name) + " " + name + styMeta.Render(fmt.Sprintf("  %d", u.count)) + "\n")
	}
	for _, l := range wrapHints([]string{"space tick", "enter show only the ticked", "x show everyone again", "esc close"}, max(1, m.width)) {
		b.WriteString("\n" + styMeta.Render(l))
	}
	return b.String()
}
