package pr

import (
	"fmt"
	"strings"
)

// Label names a pull request in a line: "#179 (GPSAR2)".
func Label(p PR) string {
	if len(p.Tickets) == 0 {
		return fmt.Sprintf("#%d", p.Number)
	}
	return fmt.Sprintf("#%d (%s)", p.Number, strings.Join(p.Tickets, ","))
}

// Summary is one pull request's stand in one line, without its label.
func Summary(p PR) string {
	parts := []string{
		fmt.Sprintf("threads: %d waiting, %d answered, %d resolved", p.Threads.Waiting, p.Threads.Answered, p.Threads.Resolved),
		fmt.Sprintf("comments from others: %d", p.Comments),
		"CI: " + p.CI,
	}
	if len(p.ReviewRequested) > 0 {
		parts = append(parts, "review pending: "+strings.Join(p.ReviewRequested, ", "))
	}
	if len(p.ChangesRequested) > 0 {
		parts = append(parts, "changes requested: "+strings.Join(p.ChangesRequested, ", "))
	}
	if len(p.Approved) > 0 {
		parts = append(parts, "approved: "+strings.Join(p.Approved, ", "))
	}
	if p.Settled {
		parts = append(parts, "settled")
	} else {
		parts = append(parts, "not settled")
	}
	return strings.Join(parts, " · ")
}

// Changes says what happened between two reads of the same pull request,
// one line each, oldest kind first: nothing where nothing a teamlead acts on
// changed. The first read is the baseline and is never compared with
// anything — a watch started on a pull request with ten comments reports
// none of them, and a watch restarted reports nothing it reported before.
func Changes(before, after PR) []string {
	var out []string
	label := Label(after)
	if a := after.LastForeign; a != nil && (before.LastForeign == nil || newer(a.At, before.LastForeign.At)) {
		switch a.Kind {
		case "review":
			out = append(out, fmt.Sprintf("%s: review by %s: %s", label, a.Author, strings.ReplaceAll(a.State, "_", " ")))
		case "thread":
			out = append(out, fmt.Sprintf("%s: new review comment by %s — threads waiting: %d", label, a.Author, after.Threads.Waiting))
		default:
			out = append(out, fmt.Sprintf("%s: new comment by %s", label, a.Author))
		}
	}
	if after.CI != before.CI {
		switch {
		case after.CI == "failure":
			out = append(out, label+": CI red")
		case after.CI == "success" && before.CI == "failure":
			out = append(out, label+": CI green again")
		}
	}
	if after.Settled && !before.Settled {
		out = append(out, label+": settled — no thread open, no review pending")
	}
	if after.State != before.State && after.State != "open" {
		out = append(out, fmt.Sprintf("%s: %s", label, after.State))
	}
	return out
}

// Over reports whether a watch on p has ended: merged or closed.
func Over(p PR) bool { return p.State == "merged" || p.State == "closed" }
