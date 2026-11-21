package role

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// build.py reads data.json written by an agent, and a set kept across rounds
// may predate the "person" key. A missing key reads as no items; a malformed
// one is refused by check() with a message, never by a traceback after the
// person has already marked the page. The report lists the open items, or the
// person never reads that an accepted ticket still waits for them.
func TestAcceptanceBuildPersonItems(t *testing.T) {
	py, err := exec.LookPath("python3")
	if err != nil {
		t.Skipf("no python3 on this machine: %v", err)
	}
	dir := t.TempDir()
	if _, err := install(dir, false, true, false); err != nil {
		t.Fatal(err)
	}
	build := filepath.Join(dir, "jaira-role-acceptance", "build.py")

	data := func(person string) string {
		p := ""
		if person != "" {
			p = `"person": ` + person + `, `
		}
		return `{"title": "t", "lang": "en", "tickets": [{"id": "AAAAAA", "title": "x", ` + p +
			`"next_lane": "done", "return_lane": "in-progress"}],
			"blocks": [{"n": 1, "title": "b", "tickets": ["AAAAAA"], "steps": [{"id": "s1", "do": "d", "covers": ["AAAAAA"]}]}]}`
	}
	run := func(args ...string) (string, error) {
		out, err := exec.Command(py, append([]string{build}, args...)...).CombinedOutput()
		return string(out), err
	}
	write := func(name, body string) string {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		return p
	}
	db := write("db.json", `{"results": {"s1": {"status": "ok"}}}`)
	dec := write("dec.json", `[]`)

	t.Run("missing reads as none", func(t *testing.T) {
		d := write("missing.json", data(""))
		out, err := run("plan", d, db)
		if err != nil {
			t.Fatalf("plan: %v\n%s", err, out)
		}
		if !strings.Contains(out, "jaira move AAAAAA --to done --force") || strings.Contains(out, "jaira dod") {
			t.Errorf("plan without person should print the bare move:\n%s", out)
		}
		if out, err := run("page", d, filepath.Join(dir, "page.html")); err != nil {
			t.Errorf("page: %v\n%s", err, out)
		}
	})

	for name, person := range map[string]string{
		"string item": `["seen by a person: it renders"]`,
		"no text":     `[{"n": 1}]`,
		"n as string": `[{"n": "1", "text": "x"}]`,
		"not a list":  `{"n": 1, "text": "x"}`,
	} {
		t.Run("malformed "+name, func(t *testing.T) {
			d := write("bad.json", data(person))
			for _, cmd := range [][]string{{"page", d, filepath.Join(dir, "p.html")}, {"plan", d, db},
				{"report", d, db, dec, filepath.Join(dir, "r.html"), filepath.Join(dir, "r.md")}} {
				out, err := run(cmd...)
				if err == nil {
					t.Errorf("%s accepted person %s:\n%s", cmd[0], person, out)
				}
				if strings.Contains(out, "Traceback") || !strings.Contains(out, "AAAAAA: person must be a list") {
					t.Errorf("%s: want a check() message, got:\n%s", cmd[0], out)
				}
			}
		})
	}

	t.Run("report lists open person items", func(t *testing.T) {
		d := write("ok.json", data(`[{"n": 3, "text": "seen by a person: the column lines up in Windows Terminal"}]`))
		md := filepath.Join(dir, "r.md")
		if out, err := run("report", d, db, dec, filepath.Join(dir, "r.html"), md); err != nil {
			t.Fatalf("report: %v\n%s", err, out)
		}
		b, err := os.ReadFile(md)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), "**For a person to see** (1)") ||
			!strings.Contains(string(b), "- `AAAAAA` 3. seen by a person: the column lines up in Windows Terminal") {
			t.Errorf("report does not list the person item:\n%s", b)
		}
	})
}
