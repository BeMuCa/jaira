package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BeMuCa/jaira/core/ticket"
)

// newTwoBoards is a repository with today's board and a second board,
// "video", whose lanes are script, render and done — none of them on today's
// board but done.
func newTwoBoards(t *testing.T) string {
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
	lanes := filepath.Join(dir, ticket.DirName, "boards", "video", "lanes")
	if err := os.MkdirAll(lanes, 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"script.md": "---\nid: script\nname: Script\nprecedence: 10\nrequires-specified: true\n---\n",
		"render.md": "---\nid: render\nname: Render\nafter: script\nprecedence: 20\n---\n",
		"done.md":   "---\nid: done\nname: Done\nafter: render\nprecedence: 60\nterminal: true\n---\n",
		"order":     "script\nrender\ndone\n",
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(lanes, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func listed(t *testing.T, m map[string]any) []string {
	t.Helper()
	var out []string
	for _, row := range m["tickets"].([]any) {
		out = append(out, row.(map[string]any)["handle"].(string))
	}
	return out
}

// A ticket created on a board runs through that board's lanes and is listed
// there, and only there: list and next without --board are today's board,
// exactly as before there were several.
func TestCreateListNextTakeABoard(t *testing.T) {
	dir := newTwoBoards(t)

	v := jsonCLI(t, dir, "create", "intro", "--board", "video", "--goal", "g", "--context", "c", "--dod", "d", "--assignee", "berk")
	if v["board"] != "video" || v["status"] != "script" {
		t.Fatalf("created on video: board %v, status %v; want video, script", v["board"], v["status"])
	}
	d := jsonCLI(t, dir, "create", "code", "--goal", "g", "--context", "c", "--dod", "d", "--assignee", "berk")
	if d["board"] != "" || d["status"] != "backlog" {
		t.Fatalf("created without --board: board %v, status %v; want today's board, backlog", d["board"], d["status"])
	}
	raw, err := os.ReadFile(jsonCLI(t, dir, "show", d["handle"].(string))["path"].(string))
	if err != nil || strings.Contains(string(raw), "board:") {
		t.Errorf("a ticket on today's board carries a board: line:\n%s", raw)
	}

	vh, dh := v["handle"].(string), d["handle"].(string)
	if got := listed(t, jsonCLI(t, dir, "list")); len(got) != 1 || got[0] != dh {
		t.Errorf("list: %v, want only %s", got, dh)
	}
	if got := listed(t, jsonCLI(t, dir, "list", "--board", "video")); len(got) != 1 || got[0] != vh {
		t.Errorf("list --board video: %v, want only %s", got, vh)
	}
	if got := listed(t, jsonCLI(t, dir, "next", "--board", "video")); len(got) != 1 || got[0] != vh {
		t.Errorf("next --board video: %v, want %s", got, vh)
	}
	if got := listed(t, jsonCLI(t, dir, "next", "--all")); len(got) != 0 && got[0] == vh {
		t.Errorf("next on today's board handed out the video ticket: %v", got)
	}

	// show needs no --board: it reads the board off the ticket, and judges
	// the ticket by that board's lanes.
	sh := jsonCLI(t, dir, "show", vh)
	if sh["board"] != "video" || sh["status_known"] != true || sh["next_lane"] != "render" {
		t.Errorf("show: board %v, status_known %v, next_lane %v; want video, true, render",
			sh["board"], sh["status_known"], sh["next_lane"])
	}
	if out, err := runCLI(t, dir, "move", vh, "--to", "render"); err != nil {
		t.Fatalf("move on the video board: %v\n%s", err, out)
	}
}

// A board name with nothing behind it is a mistake, not a new board.
func TestUnknownBoardIsRefused(t *testing.T) {
	dir := newTwoBoards(t)
	for _, args := range [][]string{
		{"create", "x", "--board", "vidoe"},
		{"list", "--board", "vidoe"},
		{"next", "--board", "vidoe"},
	} {
		_, err := runCLI(t, dir, args...)
		if exitCode(err) != ExitUsage || !strings.Contains(err.Error(), "default, video") {
			t.Errorf("%v: exit %d (%v), want %d naming the boards there", args, exitCode(err), err, ExitUsage)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, ticket.DirName, "boards", "vidoe")); !os.IsNotExist(err) {
		t.Errorf("a mistyped --board left a directory behind: %v", err)
	}
}

// Changing boards is one line; a lane the new board lacks puts the ticket in
// its first lane, and jaira says so.
func TestSetBoardMovesToFirstLaneWhenItsLaneIsMissing(t *testing.T) {
	dir := newTwoBoards(t)
	d := jsonCLI(t, dir, "create", "code", "--goal", "g", "--context", "c", "--dod", "d")
	h := d["handle"].(string)

	out, err := runCLI(t, dir, "set", h, "board=video")
	if err != nil {
		t.Fatalf("set board=video: %v\n%s", err, out)
	}
	if !strings.Contains(out, `no lane "backlog"`) || !strings.Contains(out, "now in script") {
		t.Errorf("set board did not say where the ticket landed:\n%s", out)
	}
	sh := jsonCLI(t, dir, "show", h)
	if sh["board"] != "video" || sh["status"] != "script" {
		t.Errorf("after set board=video: board %v, status %v; want video, script", sh["board"], sh["status"])
	}

	// done is on both boards, so a ticket in done keeps its lane and is not
	// told anything; back on today's board through its name.
	if out, err := runCLI(t, dir, "move", h, "--to", "done", "--force"); err != nil {
		t.Fatalf("move to done: %v\n%s", err, out)
	}
	out, err = runCLI(t, dir, "set", h, "board=default")
	if err != nil || strings.Contains(out, "first lane") {
		t.Fatalf("set board=default: %v\n%s", err, out)
	}
	sh = jsonCLI(t, dir, "show", h)
	if sh["board"] != "" || sh["status"] != "done" {
		t.Errorf("after set board=default: board %v, status %v; want today's board, done", sh["board"], sh["status"])
	}

	if out, err := runCLI(t, dir, "set", h, "board=vidoe"); exitCode(err) != ExitUsage {
		t.Errorf("set board to a board that does not exist: exit %d\n%s", exitCode(err), out)
	}
}

// A board from before there were several has no boards/ directory and no
// ticket with a board: line, and runs exactly as it did.
func TestSingleBoardRunsUnchanged(t *testing.T) {
	t.Setenv("JAIRA_USER", "berk")
	dir := t.TempDir()
	s, err := ticket.At(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Init(); err != nil {
		t.Fatal(err)
	}
	c := jsonCLI(t, dir, "create", "old", "--goal", "g", "--context", "c", "--dod", "d", "--assignee", "berk")
	h := c["handle"].(string)
	if got := listed(t, jsonCLI(t, dir, "list")); len(got) != 1 || got[0] != h {
		t.Errorf("list: %v, want %s", got, h)
	}
	if got := listed(t, jsonCLI(t, dir, "list", "--board", "default")); len(got) != 1 || got[0] != h {
		t.Errorf("list --board default: %v, want %s", got, h)
	}
	if out, err := runCLI(t, dir, "move", h, "--to", "todo"); err != nil {
		t.Fatalf("move: %v\n%s", err, out)
	}
	if _, err := os.Stat(filepath.Join(dir, ticket.DirName, "boards")); !os.IsNotExist(err) {
		t.Errorf("a single board grew a boards/ directory: %v", err)
	}
}

// jaira boards is what a teamlead reads to choose a board: every board, today's
// as "default", with its open tickets and its description — today's board set
// up from the built-ins says it is for code, a description wrapped by hand
// reads back as one line, a missing one says where it goes.
func TestBoardsListsEveryBoardWithItsDescription(t *testing.T) {
	dir := newTwoBoards(t)
	desc := filepath.Join(dir, ticket.DirName, "boards", "video", "lanes", "description")
	if err := os.WriteFile(desc, []byte("Videos made from demos:\n  script, render.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	jsonCLI(t, dir, "create", "intro", "--board", "video", "--goal", "g", "--context", "c", "--dod", "d")

	v := jsonCLI(t, dir, "boards")
	rows := v["boards"].([]any)
	if len(rows) != 2 {
		t.Fatalf("boards: %v, want default and video", rows)
	}
	def, vid := rows[0].(map[string]any), rows[1].(map[string]any)
	if def["name"] != "default" || def["description"] != "Changes to this repository's code: features, fixes, refactors, tests and docs." || def["open"] != float64(0) {
		t.Errorf("default row: %v", def)
	}
	if vid["name"] != "video" || vid["description"] != "Videos made from demos: script, render." || vid["open"] != float64(1) || vid["description_file"] != desc {
		t.Errorf("video row: %v", vid)
	}

	if err := os.Remove(desc); err != nil {
		t.Fatal(err)
	}
	out, err := runCLI(t, dir, "boards")
	if err != nil {
		t.Fatalf("boards: %v\n%s", err, out)
	}
	if !strings.Contains(out, "no description yet — write one into .jaira/boards/video/lanes/description") {
		t.Errorf("a board without a description does not say where it goes:\n%s", out)
	}
}

// What a worker or a dispatcher runs on a ticket of another board works that
// board's lanes: next --per-lane --board finds it in its own lane, and the
// lane's brief and the move come from the ticket's board without --board.
func TestWorkerPathRunsTheTicketsBoard(t *testing.T) {
	dir := newTwoBoards(t)
	render := filepath.Join(dir, ticket.DirName, "boards", "video", "lanes", "render.md")
	body := "---\nid: render\nname: Render\nafter: script\nprecedence: 20\nagentic: true\nmodel-tier: cheap\noutput-produces: [outcome-what]\n---\n\n# Prompt\n\nRender the film.\n"
	if err := os.WriteFile(render, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	h := jsonCLI(t, dir, "create", "intro", "--board", "video", "--goal", "g", "--context", "c", "--dod", "d", "--assignee", "berk")["handle"].(string)
	if out, err := runCLI(t, dir, "move", h, "--to", "render"); err != nil {
		t.Fatalf("move to render: %v\n%s", err, out)
	}

	per := jsonCLI(t, dir, "next", "--per-lane", "--board", "video")
	if per["count"] != float64(1) || per["lanes"].([]any)[0].(map[string]any)["lane"] != "render" {
		t.Errorf("next --per-lane --board video: %v, want the ticket in render", per)
	}
	if per := jsonCLI(t, dir, "next", "--per-lane"); per["count"] != float64(0) {
		t.Errorf("next --per-lane on today's board handed out the video ticket: %v", per)
	}
	brief := jsonCLI(t, dir, "show", h, "--for-lane", "render")
	if !strings.Contains(brief["prompt"].(string), "Render the film.") || brief["model_tier"] != "cheap" {
		t.Errorf("--for-lane render: %v, want the video board's render lane", brief)
	}
}

// validate judges each ticket by its own board's lanes: a ticket in a lane
// only its board has is valid, not read-only.
func TestValidateChecksEachTicketAgainstItsBoard(t *testing.T) {
	dir := newTwoBoards(t)
	jsonCLI(t, dir, "create", "intro", "--board", "video", "--goal", "g", "--context", "c", "--dod", "d", "--assignee", "berk")
	v := jsonCLI(t, dir, "validate")
	for _, p := range v["problems"].([]any) {
		if p.(map[string]any)["code"] == "unknown_lane" {
			t.Errorf("validate calls a video ticket's lane unknown: %v", p)
		}
	}
}

// The merge driver measures progress in the ticket's own board's lanes:
// today's board knows neither script nor render, and would keep ours however
// far the other side got.
func TestMergeDriverUsesTheTicketsBoard(t *testing.T) {
	dir := newTwoBoards(t)
	doc := func(status string) string {
		return "---\nid: 01M4X5JKRN4AV3KX56SES0EG0M\ntitle: t\nboard: video\nstatus: " + status +
			"\ncreated-at: 2026-10-17T12:00:00Z\nupdated-at: 2026-10-17T12:00:00Z\n---\n"
	}
	paths := map[string]string{}
	for side, status := range map[string]string{"base": "done", "ours": "script", "theirs": "render"} {
		p := filepath.Join(t.TempDir(), side+".md")
		if err := os.WriteFile(p, []byte(doc(status)), 0o644); err != nil {
			t.Fatal(err)
		}
		paths[side] = p
	}
	if out, err := runCLI(t, dir, "merge-driver", paths["base"], paths["ours"], paths["theirs"]); err != nil {
		t.Fatalf("merge-driver: %v\n%s", err, out)
	}
	got, _ := os.ReadFile(paths["ours"])
	if !strings.Contains(string(got), "status: render") {
		t.Errorf("the further lane on the video board did not win:\n%s", got)
	}
}

// lanes --board reads another board's lanes, which is how the cockpit tells
// whose lane a ticket there waits in; a board with nothing behind it is a
// usage error, as on list. requires_question names the human lane, which
// human_exit does not.
func TestLanesTakeABoard(t *testing.T) {
	dir := newTwoBoards(t)
	ids := func(v map[string]any) map[string]map[string]any {
		out := map[string]map[string]any{}
		for _, row := range v["lanes"].([]any) {
			l := row.(map[string]any)
			out[l["id"].(string)] = l
		}
		return out
	}
	video := ids(jsonCLI(t, dir, "lanes", "--board", "video"))
	if len(video) != 3 || video["script"] == nil || video["render"] == nil || video["done"] == nil {
		t.Errorf("lanes --board video: %v, want script, render and done", video)
	}
	today := ids(jsonCLI(t, dir, "lanes"))
	if today["script"] != nil || today["human"]["requires_question"] != true || today["signoff"]["requires_question"] != false {
		t.Errorf("lanes: script %v, human %v, signoff %v; want no script, requires_question on human only", today["script"], today["human"], today["signoff"])
	}
	if _, err := runCLI(t, dir, "lanes", "--board", "vidoe"); exitCode(err) != ExitUsage {
		t.Errorf("lanes --board vidoe: exit %d (%v), want %d", exitCode(err), err, ExitUsage)
	}
}
