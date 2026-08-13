---
phase: quick/260813-ovw
plan: 01
type: execute
wave: 1
depends_on: []
files_modified:
  - internal/tui/model.go
  - internal/tui/pipeline.go
  - internal/tui/lanefocus.go
  - internal/tui/lanefocus_test.go
  - internal/tui/view.go
  - README.md
  - docs/COMMANDS.md
autonomous: true
requirements: [QUICK-260813-ovw]

must_haves:
  truths:
    - "Pressing enter on a step in the compact view opens that one lane filling the screen, not the multi-column board"
    - "The single-lane view shows more of each ticket than a board column can, because it has the width for it"
    - "Moving between lanes from inside the single-lane view stays in the single-lane view"
    - "esc returns to the compact view it was opened from, not to the board"
    - "enter on a ticket opens the same detail pane the board opens"
  artifacts:
    - path: "internal/tui/lanefocus.go"
      provides: "the single-lane view: render and keys"
  key_links:
    - from: "internal/tui/pipeline.go"
      to: "modeLaneFocus"
      via: "enter on a step opens the lane instead of switching to the board"
      pattern: "modeLaneFocus"
---

<objective>
In the compact view (`v`), `enter` on a step currently does this
(`internal/tui/pipeline.go:261-265`):

    // Opening a step is the full board with that lane focused, which is the
    // block view of what is in it.
    m.mode = modeBoard
    m.rebuild()

So the compact view's one drill-down throws you back to the multi-column board.
The user wants the opposite: the lane they picked, alone, filling the screen.

Purpose: the compact view exists to see the whole flow at a glance. Picking a
step should go deeper into that step, not sideways into a different view of
everything. Output: a single-lane view, reachable from the compact view, that
uses the full width to show more of each ticket than a 22-34 column can.
</objective>

<context>
@CLAUDE.md
@internal/tui/pipeline.go
@internal/tui/model.go
@internal/tui/view.go

<interfaces>
Modes are an iota block in internal/tui/model.go: modeBoard, modeDetail,
modeFilter, modeHelp, modeMove, modeCreate, modeMessage, modeProjects, modeEdit,
modePipeline, modeLanes. Add one.

Existing render pieces worth reusing rather than reinventing:
  (m *Model) renderColumn(idx, w, h int)  — view.go:224, bordered lane column
  (m *Model) renderCard(t, w, selected)   — view.go:276, one card
  (m *Model) rebuild()                    — model.go:162
Board columns are clamped to minColWidth 22 / maxColWidth 34 (view.go:47-50).
The single-lane view is not bound by that — it has the whole terminal.

The keys the board already uses, to stay consistent: j/k or ↓/↑ move card,
h/l or ←/→ move lane, enter open detail, esc back, g/G first/last, y copy id.

Go is not on the default PATH:
  export PATH=$PATH:$HOME/.local/go/bin
</interfaces>
</context>

<tasks>

<task type="auto">
  <name>Task 1: A lane on its own, filling the screen</name>
  <files>internal/tui/model.go, internal/tui/pipeline.go, internal/tui/lanefocus.go, internal/tui/lanefocus_test.go</files>
  <action>
Add `modeLaneFocus` to the mode iota and a `lanefocus.go` holding its render and
key handling, so `view.go` does not grow another screen.

**Entry.** In `pipelineKey`, `enter` sets `m.mode = modeLaneFocus` for the
currently selected `m.laneIdx` instead of `modeBoard`. Rewrite the comment above
it — the current one describes the behaviour being replaced.

**Render.** The lane's name, its ticket count, and its tickets down the screen,
using the full terminal width. Do not simply call `renderColumn` at a wide width:
a card laid out for a 30-column box wastes a 120-column terminal. Decide what the
extra width buys by reading `renderCard` and the ticket fields available, and say
in a comment why you picked those. At minimum a reader should see more per ticket
here than the board shows — the board is for scanning, this is for reading.

Keep the lane's own state visible: if the lane is agentic, its tier; if it is
unknown, that it is read-only; and the same per-ticket markers the board uses for
"waiting on your answer" and "waiting on sign-off", so a ticket does not look
different depending on which screen you are on.

**Keys.**
- `j`/`down`, `k`/`up` — move between tickets in this lane
- `g`/`G` — first/last
- `h`/`left`, `l`/`right` — previous/next lane, staying in this view. This is the
  point of the screen; leaving it to change lane would defeat it.
- `enter` — open the ticket's detail pane, exactly as the board does
- `esc` — back to the compact view it came from, not to the board
- `v` — back to the compact view as well, since `v` is already "toggle the
  compact view" everywhere else
- `q`/`ctrl+c` — quit, as everywhere

An empty lane must render as an empty lane with its name and a plain "no tickets"
line, not as a blank screen. Moving to the next lane from an empty one must work.

**Tests** in `internal/tui/lanefocus_test.go`, following the style of the
existing TUI tests (`newTestModel`, `stripANSI`, driving `m.key(...)`):
- from the compact view, `enter` lands in `modeLaneFocus`, not `modeBoard`
- the rendered output contains the selected lane's name and its tickets' titles,
  and does NOT contain a ticket that belongs to a different lane
- `l` then `h` moves lane and back, and the mode is still `modeLaneFocus`
- `esc` returns to `modePipeline`
- `enter` on a ticket reaches `modeDetail`
- an empty lane renders without panicking and still allows `l` to move on
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go vet ./... && go test ./internal/tui/... -count=1 2>&1 | tail -5</automated>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go test ./... -count=1 2>&1 | grep -c FAIL</automated>
  </verify>
  <done>enter in the compact view opens one lane filling the screen; h/l move lane without leaving it; esc goes back to the compact view; every listed test passes and the whole suite is green.</done>
</task>

<task type="auto">
  <name>Task 2: Say the key exists</name>
  <files>internal/tui/view.go, README.md, docs/COMMANDS.md</files>
  <action>
The `?` help screen (`renderHelp` in view.go) lists the board's keys. Add the
single-lane view where it fits the existing grouping, and correct the compact
view's entry if it claims `enter` opens the board.

`README.md` and `docs/COMMANDS.md` both carry a "Keys" block with a line
describing the compact view. Both currently say `enter` opens the step; make them
say what it now does, and name the keys available inside the single-lane view.
Keep both blocks in their existing terse shape — they are reference cards.
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go test ./internal/tui/... -count=1 2>&1 | tail -3</automated>
    <automated>grep -c "lane" README.md docs/COMMANDS.md</automated>
  </verify>
  <done>The help screen and both key tables describe the single-lane view and no longer claim enter returns to the board.</done>
</task>

</tasks>

<verification>
export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go vet ./... && go test ./... -count=1

Human check, which cannot be automated here and must be reported as unverified:
open the board, press `v`, pick a step, press enter, and confirm the lane fills
the screen and reads well at both a narrow and a wide terminal.
</verification>

<success_criteria>
- `enter` in the compact view opens the selected lane alone, filling the screen.
- `h`/`l` change lane without leaving the view.
- `esc` returns to the compact view, not to the board.
- A reader sees more per ticket than the board column shows.
- `go build ./... && go vet ./... && go test ./... -count=1` passes.
</success_criteria>

<output>
Create `.planning/quick/260813-ovw-one-lane-full-width-from-the-compact-vie/260813-ovw-SUMMARY.md` when done
</output>
