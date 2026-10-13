package cli

import (
	"strings"
	"testing"
)

// spawn.sh and run-lane.sh read single fields out of Herdr's JSON with this
// instead of python3, which a Windows machine may not have (issue #79). A
// string comes back bare, an object as one line of JSON, and a missing path
// either fails — a pane id nobody got must stop spawn.sh — or prints --default.
func TestJSONGetPrintsOneField(t *testing.T) {
	const in = `{"result":{"root_pane":{"pane_id":"p1"},"pane":{"agent":"claude"}},
		"review":{"gaps":"a <gap> ä"},"items":[1,"two"],"q":"","n":null,"num":3}`
	dir := t.TempDir()
	for path, want := range map[string]string{
		"result.root_pane.pane_id": "p1\n",
		"review":                   `{"gaps":"a <gap> ä"}` + "\n",
		"items.1":                  "two\n",
		"q":                        "\n",
		"n":                        "\n",
		"num":                      "3\n",
	} {
		out, err := runCLIStdin(t, dir, in, "json", "get", path)
		if err != nil || out != want {
			t.Errorf("json get %s = %q, %v; want %q", path, out, err, want)
		}
	}
	for _, path := range []string{"result.pane.agent_status", "items.2", "q.x"} {
		if _, err := runCLIStdin(t, dir, in, "json", "get", path); exitCode(err) != ExitNotFound {
			t.Errorf("json get %s: exit %d, want %d", path, exitCode(err), ExitNotFound)
		}
		out, err := runCLIStdin(t, dir, in, "json", "get", path, "--default", "-")
		if err != nil || out != "-\n" {
			t.Errorf("json get %s --default -: %q, %v", path, out, err)
		}
	}
	if _, err := runCLIStdin(t, dir, "Python\n", "json", "get", "x"); exitCode(err) != ExitNotFound {
		t.Errorf("text that is no JSON: exit %d, want %d", exitCode(err), ExitNotFound)
	}
}

// The brief the cockpit's chats get typed in: the newest four notes whole,
// continuation lines included — one that starts "#" or "- " too, a note ends
// only at the next dated heading — how many there are, and nothing of the
// context.
func TestTicketBriefCarriesTheNewestNotes(t *testing.T) {
	body := "\n# T\n\n## Progress\n" +
		"- **2026-10-01 10:00 · A** — first note\n" +
		"- **2026-10-02 10:00 · A** — second note\n  goes on\n" +
		"- **2026-10-03 10:00 · A** — third note\n#79 is the issue\n- a list line\n" +
		"- **2026-10-04 10:00 · A** — fourth note\n" +
		"- **2026-10-05 10:00 · A** — fifth note\n"
	j := map[string]any{"handle": "AB12CD", "title": "T", "goal": "G", "question": "Q?",
		"review": map[string]string{"verdict": "V"}, "context": "LONG CONTEXT"}
	line, err := compactJSON(ticketBrief(j, body))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`{"id":"AB12CD","title":"T","goal":"G","question":"Q?","review":{"verdict":"V"},`,
		`second note\n  goes on`, `third note\n#79 is the issue\n- a list line`, "fifth note", `"notes_total":5}`} {
		if !strings.Contains(line, want) {
			t.Errorf("brief lacks %q:\n%s", want, line)
		}
	}
	for _, not := range []string{"first note", "LONG CONTEXT", "\n"} {
		if strings.Contains(line, not) {
			t.Errorf("brief carries %q:\n%s", not, line)
		}
	}
	if line, _ := compactJSON(ticketBrief(j, "")); !strings.Contains(line, `"notes_last":[],"notes_total":0`) {
		t.Errorf("a ticket without notes: %s", line)
	}
}
