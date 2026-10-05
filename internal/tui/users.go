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

	m.users, m.userIdx = rows, 0
	m.userPicked = map[string]bool{}
	_, picked := splitUserTerm(m.filter)
	for _, name := range picked {
		m.userPicked[name] = true
	}
	m.returnTo = m.mode
	m.mode = modeUsers
}

// splitUserTerm takes the user: condition out of a filter: the other
// conditions as they were typed, and the names it held, lower-cased.
func splitUserTerm(filter string) (rest, names []string) {
	for _, term := range filterTerms(filter) {
		key, val, ok := strings.Cut(term, ":")
		if !ok || !strings.EqualFold(key, "user") {
			rest = append(rest, term)
			continue
		}
		for _, v := range strings.Split(strings.ReplaceAll(val, `"`, ""), ",") {
			if v = strings.TrimSpace(v); v != "" {
				names = append(names, strings.ToLower(v))
			}
		}
	}
	return rest, names
}

// applyUsers writes the ticked people into the board filter in place of any
// user: condition it had, keeping everything else typed there. Like the
// milestone picker it writes the ordinary filter rather than a second kind
// of narrowing, so / shows it and esc on the board clears it.
func (m *Model) applyUsers() {
	rest, _ := splitUserTerm(m.filter)
	var picked []string
	for _, u := range m.users {
		if !m.userPicked[strings.ToLower(u.name)] {
			continue
		}
		name := u.name
		if strings.ContainsAny(name, " \t") {
			name = `"` + name + `"`
		}
		picked = append(picked, name)
	}
	if len(picked) > 0 {
		rest = append(rest, "user:"+strings.Join(picked, ","))
	}
	m.filter = strings.Join(rest, " ")
	m.input = m.filter
	m.rebuild()
}

// keyUsers drives the user picker. esc leaves the filter as it was: the
// ticks live in the picker until enter writes them.
func (m *Model) keyUsers(s string) {
	n := len(m.users)
	switch s {
	case "esc", "u", "q":
		m.mode = m.returnTo
	case "j", "down":
		if n > 0 {
			m.userIdx = (m.userIdx + 1) % n
		}
	case "k", "up":
		if n > 0 {
			m.userIdx = (m.userIdx - 1 + n) % n
		}
	case "space":
		if m.userIdx >= 0 && m.userIdx < n {
			k := strings.ToLower(m.users[m.userIdx].name)
			m.userPicked[k] = !m.userPicked[k]
		}
	case "enter":
		m.applyUsers()
		m.mode = m.returnTo
	case "x":
		m.userPicked = map[string]bool{}
		m.applyUsers()
		m.mode = m.returnTo
	}
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
		if i == m.userIdx {
			marker, name = stySelected.Render("▌ "), stySelected.Render(name)
		}
		tick := styMeta.Render("○")
		if m.userPicked[strings.ToLower(u.name)] {
			tick = styOK.Render("●")
		}
		b.WriteString(marker + tick + " " + name + styMeta.Render(fmt.Sprintf("  %d", u.count)) + "\n")
	}
	for _, l := range wrapHints([]string{"space tick", "enter show only the ticked", "x show everyone again", "esc close"}, max(1, m.width)) {
		b.WriteString("\n" + styMeta.Render(l))
	}
	return b.String()
}
