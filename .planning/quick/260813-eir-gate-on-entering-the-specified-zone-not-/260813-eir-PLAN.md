---
phase: quick/260813-eir
plan: 01
type: execute
wave: 1
depends_on: []
files_modified:
  - core/lane/lane.go
  - core/lane/builtin/10-todo.md
  - core/gate/gate.go
  - core/gate/gate_test.go
  - docs/AGENTS.md
autonomous: true
requirements: [QUICK-260813-eir]

must_haves:
  truths:
    - "The promotion gate fires when a ticket enters the specified zone, not when it leaves a lane called backlog"
    - "A lane declares itself the start of that zone with `requires-specified: true`; the built-in todo lane carries it"
    - "A ticket with only a title can move from backlog into a lane that sits before the specified zone, and cannot move past it"
    - "Every existing behaviour of the gate for backlog -> todo, backlog -> in-progress and backlog -> review is unchanged"
    - "A lane set with no lane declaring requires-specified falls back to today's behaviour rather than letting everything through"
  artifacts:
    - path: "core/lane/lane.go"
      provides: "RequiresSpecified on Lane, parsed from frontmatter"
    - path: "core/gate/gate.go"
      provides: "the promotion gate keyed on the specified zone"
  key_links:
    - from: "core/gate/gate.go"
      to: "lane.RequiresSpecified"
      via: "the first lane in precedence order that sets it marks where specification becomes mandatory"
      pattern: "RequiresSpecified"
---

<objective>
`core/gate/gate.go:138` hardcodes the promotion gate to the literal lane id
`"backlog"`:

    leavingBacklog := t.Status == "backlog" && req.To != "backlog"

Any move out of backlog demands goal, definition-of-done, context and assignee.
That makes a thinking lane impossible: a one-line scrap could never move into a
brainstorm lane, because the gate would demand exactly what brainstorming is
supposed to produce.

Purpose: let a board have lanes for work that is not specified yet, without
weakening the guarantee that nothing reaches the working lanes unspecified.
Output: the gate keyed on a lane contract instead of a lane name, exactly like
`requires-question`, `requires-human-exit`, `requires-outcome` and
`requires-option` already are.
</objective>

<context>
@CLAUDE.md
@core/gate/gate.go
@core/lane/lane.go
@core/lane/builtin/10-todo.md

<interfaces>
Lane already carries this family of contract flags, parsed in `parse()` in
core/lane/lane.go and read by the gate:
  RequiresQuestion, RequiresHumanExit, RequiresOutcome, RequiresNonModelSignal,
  RequiresOption

The gate's current promotion block (core/gate/gate.go, around 133-147):
  from := env.Lanes.Index(t.Status)
  to := env.Lanes.Index(req.To)
  leavingBacklog := t.Status == "backlog" && req.To != "backlog"
  if leavingBacklog {
      vs = append(vs, missingPromotionFields(t)...)
      vs = append(vs, blockedBy(env, t)...)
  }
  if to > from && !leavingBacklog {
      vs = append(vs, blockedBy(env, t)...)
  }

`env.Lanes.Index(id)` gives a lane's position in precedence order.
PromotionFields and missingPromotionFields stay exactly as they are.

Go is not on the default PATH:
  export PATH=$PATH:$HOME/.local/go/bin
</interfaces>
</context>

<tasks>

<task type="auto">
  <name>Task 1: A lane can declare where specification becomes mandatory</name>
  <files>core/lane/lane.go, core/lane/builtin/10-todo.md</files>
  <action>
Add `RequiresSpecified bool` to the Lane struct and parse `requires-specified`
from lane frontmatter, following exactly how `requires-question` and
`requires-human-exit` are declared and parsed in the same file. Do not invent a
different spelling or a different parsing path.

Document it in the struct comment as: the first lane in precedence order that
sets this marks the boundary between thinking about a ticket and working on it.
Everything before it is a place a half-formed ticket may sit; nothing at or after
it may hold a ticket that is missing its promotion fields.

Add `requires-specified: true` to `core/lane/builtin/10-todo.md`. Todo is where a
ticket becomes ready to be picked up, so it is the boundary. Do not add the flag
to any other built-in lane — one boundary, not a scattering of them.
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go test ./core/lane/... -count=1 2>&1 | tail -5</automated>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go run ./cmd/jaira lanes --json 2>/dev/null | grep -c "todo"</automated>
  </verify>
  <done>Lane carries RequiresSpecified, parsed from `requires-specified`, and the built-in todo lane sets it. Build and the lane tests pass.</done>
</task>

<task type="auto">
  <name>Task 2: The gate fires on entering the specified zone</name>
  <files>core/gate/gate.go, core/gate/gate_test.go</files>
  <action>
Replace the `leavingBacklog` condition with one keyed on the lane contract.

Compute the boundary once: the index, in precedence order, of the first lane with
`RequiresSpecified`. The promotion fields are required when the move lands at or
after that boundary and the ticket is not already at or after it — i.e. the
ticket is crossing into the specified zone. A move that starts and ends inside
the zone must not re-run the promotion check, because that is today's behaviour
and re-running it would start refusing moves that work now.

Fallback: if no installed lane sets `RequiresSpecified` — an old or hand-built
lane set — fall back to the current rule (leaving a lane whose id is `backlog`).
Fail closed, never open: a board with no boundary declared must not let an
unspecified ticket walk into `done`. Say why in a comment.

Keep `blockedBy` behaviour exactly as it is: checked when crossing into the zone,
and checked on any other forward move, never twice for the same move.

Rename `leavingBacklog` to something that says what it now means — the crossing,
not the lane. Update the comment above it: the current one says "Leaving the
backlog is the moment work becomes real", which stops being true.

Tests in core/gate/gate_test.go. Follow the existing test style in that file; do
not restructure it. Cover:
- backlog -> todo with fields missing is refused, naming the missing fields (this
  is today's behaviour and must be unchanged)
- backlog -> todo with all fields present is allowed
- backlog -> in-progress with fields missing is refused (skipping todo must not
  skip the gate)
- a ticket with only a title moves from backlog into a custom lane placed before
  todo and carrying no `requires-specified`, and is allowed
- that same ticket is then refused when it tries to move from that lane into todo
- a move between two lanes both inside the zone (todo -> in-progress) does not
  re-run the promotion check
- with a lane set where no lane declares `requires-specified`, backlog -> todo
  with fields missing is still refused

Build the custom lane for the test the way core/lane already loads one; if the
test needs a lane set assembled by hand, do that rather than writing to the
user's home directory.
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go vet ./... && go test ./... -count=1 2>&1 | grep -Ev "no test files" | tail -12</automated>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && grep -c "t.Status == \"backlog\" && req.To != \"backlog\"" core/gate/gate.go</automated>
  </verify>
  <done>The promotion gate is keyed on the specified zone; every listed test passes; the whole suite is green; the old hardcoded condition survives only inside the documented fallback.</done>
</task>

<task type="auto">
  <name>Task 3: Say that a board can have lanes before the specified zone</name>
  <files>docs/AGENTS.md</files>
  <action>
docs/AGENTS.md explains the gates to whoever drives the board. Add a short
paragraph, in the register of the surrounding text, saying that lanes before the
first `requires-specified` lane are for tickets that are not specified yet — a
scrap, an idea, something being thought through — and that the promotion gate
fires when a ticket crosses into the specified zone, not when it leaves a lane
with a particular name. Keep it to a few sentences; this file is a reference, not
a tutorial.

Do not add a brainstorm lane. Whether one ships as a built-in is a separate
decision that has not been made.
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go build ./... && grep -c "requires-specified" docs/AGENTS.md</automated>
  </verify>
  <done>docs/AGENTS.md describes the specified zone and what sits before it.</done>
</task>

</tasks>

<verification>
export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go vet ./... && go test ./... -count=1

Manual: in a throwaway repo with JAIRA_HOME redirected, create a ticket with only
a title and confirm it is still refused entry to todo, naming the missing fields.
</verification>

<success_criteria>
- A ticket with only a title can sit in, and move between, lanes before the specified zone.
- Nothing reaches todo or beyond without goal, definition-of-done, context and assignee.
- The gate reads a lane contract, not a lane name.
- A lane set that declares no boundary behaves exactly as it does today.
- `go build ./... && go vet ./... && go test ./... -count=1` passes.
</success_criteria>

<output>
Create `.planning/quick/260813-eir-gate-on-entering-the-specified-zone-not-/260813-eir-SUMMARY.md` when done
</output>
