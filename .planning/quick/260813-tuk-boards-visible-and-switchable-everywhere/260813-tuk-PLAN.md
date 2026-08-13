---
phase: quick/260813-tuk
plan: 01
type: execute
wave: 1
depends_on: []
files_modified:
  - internal/tui/model.go
  - internal/tui/view.go
  - internal/tui/pipeline.go
  - internal/tui/projects_test.go
  - internal/tui/lanes.go
  - internal/tui/lanes_test.go
  - core/lane/builtin/45-signoff.md
  - README.md
  - docs/COMMANDS.md
autonomous: true
requirements: [QUICK-260813-tuk]

must_haves:
  truths:
    - "The board and the compact view both show the recorded boards, without pressing p first"
    - "1-9 switches board from the board view as well as from the compact view"
    - "A board someone is currently working in is marked, so it is visible at a glance which ones are live"
    - "The sign-off lane reads Human Review"
    - "Copying a built-in lane into a project unchanged does not warn that it overrides a built-in"
---

<objective>
Four things, all found by the user against the running binary.

**The board list is empty unless you press `p`.** `m.projects` is assigned in
exactly one place — `internal/tui/model.go:665-666`, inside `case "p"`. The
compact view renders its project tabs from that slice
(`internal/tui/pipeline.go:129-145`) and its `1`-`9` handler reads it
(`pipeline.go:270`). So the tabs are invisible and the number keys do nothing
until the user has opened the project picker at least once. The footer advertises
"1-9 switch project" regardless.

**The board view has no `1`-`9` at all** — verified, zero bindings. Switching
board from the main view means `p`, then arrow keys, then enter.

**Nothing shows which boards are being worked in.** The session panel
(`renderSessions`) reports sessions for the board you are looking at. With
several agents running across several repos, the question "which boards are live"
has no answer on screen.

**A lane copied into a project warns that it overrides a built-in, even when it
is the built-in.** Pressing `u` on the review lane writes
`.jaira/lanes/review.md`, byte-identical to the shipped file except the `creator:`
stamp `use` adds — and the loader then warns about an override that changes
nothing. Confirmed on the user's own board.
</objective>

<context>
@CLAUDE.md
@internal/tui/model.go
@internal/tui/pipeline.go
@internal/tui/view.go
@internal/tui/lanes.go
@core/lane/lane.go

<interfaces>
project.Load() returns the recorded boards; project.Remember(root) records one.
core/session holds per-board session state; `renderSessions` (view.go) already
reads it for the current board, and `session.Session` has Focus, TicketID, Model
and a Stale() test.
The compact view's tabs and its 1-9 handler: pipeline.go:129-145 and 262-277.
Lane override detection: core/lane/lane.go around 422 sets `Overrides`; the
settings screen renders it at internal/tui/lanes.go:339.

Go is not on the default PATH:
  export PATH=$PATH:$HOME/.local/go/bin
</interfaces>
</context>

<tasks>

<task type="auto">
  <name>Task 1: The boards are always known, and 1-9 works everywhere</name>
  <files>internal/tui/model.go, internal/tui/view.go, internal/tui/pipeline.go, internal/tui/projects_test.go</files>
  <action>
Load the recorded boards when the model is built, not when `p` is pressed, so
every screen that wants the list has it. Keep the reload inside `p` — the list
can change while the TUI is open — but it must no longer be the only place.

Bind `1`-`9` on the board view to the same switch the compact view performs
(`pipeline.go:262-277`), so the two views agree. Factor the switch into one
function rather than copying the branch; two copies of "which board is number
three" will drift.

Show the boards on the board view as well. It has less spare room than the
compact view, so decide what fits and say why in a comment — a single line of
numbered names above the columns is the obvious candidate, dropped when the
terminal is too narrow rather than wrapped.

Tests in internal/tui/projects_test.go:
- a freshly built model already has the recorded boards, before any key
- `2` on the board view switches to the second board; `2` on the current board is
  a no-op rather than a reload
- a number beyond the list is ignored, not a panic
- the compact view's tabs render without `p` having been pressed
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go vet ./... && go test ./... -race -count=1 2>&1 | grep -c FAIL</automated>
  </verify>
  <done>Boards appear on both views without pressing p, 1-9 switches from either, and one function decides what a number means.</done>
</task>

<task type="auto">
  <name>Task 2: Mark the boards that are live</name>
  <files>internal/tui/model.go, internal/tui/view.go, internal/tui/pipeline.go, internal/tui/projects_test.go</files>
  <action>
In the numbered board list, mark the boards that have a session working in them
right now, so a glance answers "where are my agents". Read the same session state
`renderSessions` already uses, for each recorded board rather than only the
current one.

Rules that keep it honest:
- a stale session is not a live one — `session.Stale()` already draws that line,
  and a crashed agent's leftover state must not make a board look busy
- reading another board's session state must never fail the render: a board on a
  disconnected drive, or one deleted since it was recorded, is shown unmarked
  rather than taking the screen down
- the marker must not rely on colour alone; the project's own rendering rule is
  that state is carried by a glyph or a label as well

Do not poll on a timer. Read when the list is loaded and when the board reloads —
this is a board that gets glanced at, not a monitor.

Tests: a board with a live session is marked, one with only a stale session is
not, and an unreadable board renders unmarked without an error.
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go vet ./... && go test ./... -race -count=1 2>&1 | grep -c FAIL</automated>
  </verify>
  <done>Live boards are marked, stale ones are not, an unreadable board does not break the render, and tests cover all three.</done>
</task>

<task type="auto">
  <name>Task 3: Sign-off reads Human Review</name>
  <files>core/lane/builtin/45-signoff.md, README.md, docs/COMMANDS.md</files>
  <action>
Change the sign-off lane's display `name` to `Human Review`. The lane `id` stays
`signoff` — it appears in every ticket's `status:` field and in gate messages, and
renaming it would rewrite history for no gain.

Check every place the old display name is used in prose — README, docs, the
lane's own description, other lanes' descriptions — and update what refers to the
lane by its display name. Leave references to the *id* alone.

Note in your summary, do not fix: the board now shows `Review` (a model's) next
to `Human Review` (a person's), which reads well, but the HITL lane is still
called `HITL`, which is now the least self-explanatory name on the board. That is
the user's call, not yours.
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go test ./... -race -count=1 2>&1 | grep -c FAIL</automated>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go build -o /tmp/jaira-tuk ./cmd/jaira && /tmp/jaira-tuk lanes | grep -c "Human Review"</automated>
  </verify>
  <done>The lane displays as Human Review, its id is unchanged, and no stale display name is left in the docs.</done>
</task>

<task type="auto">
  <name>Task 4: Do not warn about an override that changes nothing</name>
  <files>core/lane/lane.go, core/lane/lane_test.go, internal/tui/lanes.go, internal/tui/lanes_test.go</files>
  <action>
A lane file that matches the built-in it shadows is not an override in any sense
the reader cares about — it is a copy, which is exactly what `jaira lanes use`
produces. Suppress both the load warning and the settings screen's orange
"overrides" label when the shadowing lane is equivalent to the built-in.

Equivalent means: every field that changes behaviour is the same. `creator:` is
stamped by `use`/`publish` and must not count. Decide the comparison by reading
the Lane struct and say in a comment which fields are compared and why the
excluded ones are excluded — a comparison that silently ignores a behavioural
field would hide a real override, which is worse than the noise being fixed.

Keep the protection warning untouched: an override that drops
`requires-human-exit` or `requires-nonmodel-signal` must still shout, even in the
unlikely case the rest matches.

Tests: a byte-identical copy warns about nothing; a copy differing only in
`creator` warns about nothing; a copy with a changed prompt warns as before; a
copy that drops a protection still produces the protection warning.
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go vet ./... && go test ./... -race -count=1 2>&1 | grep -c FAIL</automated>
  </verify>
  <done>A copied built-in produces no warning and no orange label; a genuinely changed one still does; dropping a protection always does.</done>
</task>

</tasks>

<verification>
export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go vet ./... && go test ./... -race -count=1

Against a real binary: in a repo whose `.jaira/lanes/` holds a copy of a built-in
lane made with `jaira lanes use`, `jaira lanes` must print no warning.

Human check, unverifiable here: open the board, confirm the numbered boards are
listed without pressing p, that 1-9 switches, and that a board with a running
agent is marked.
</verification>

<success_criteria>
- Boards are listed on the board and the compact view without pressing p first.
- 1-9 switches board from either view, through one shared function.
- Boards with a live session are marked; stale ones are not.
- The sign-off lane displays as Human Review; its id is unchanged.
- A copied built-in lane produces no override warning; a changed one still does.
- `go build ./... && go vet ./... && go test ./... -race -count=1` passes.
</success_criteria>

<output>
Create `.planning/quick/260813-tuk-boards-visible-and-switchable-everywhere/260813-tuk-SUMMARY.md` when done
</output>
