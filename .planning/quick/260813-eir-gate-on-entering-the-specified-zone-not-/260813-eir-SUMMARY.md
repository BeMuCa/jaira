# Quick 260813-eir: Gate on entering the specified zone, not leaving backlog Summary

The promotion gate (goal, definition-of-done, context, assignee) used to fire
on the hardcoded condition `t.Status == "backlog" && req.To != "backlog"`.
That made a thinking lane before backlog's exit impossible: any lane a
half-formed ticket sat in before real work started would still be forced to
produce the fields brainstorming is supposed to produce. The gate is now keyed
on a lane contract, `requires-specified`, the same way `requires-question`,
`requires-human-exit`, `requires-outcome` and `requires-option` already work.

## What changed

1. **`core/lane/lane.go`** — added `RequiresSpecified bool` to `Lane`, parsed
   from `requires-specified` in frontmatter, following the same `boolOf`
   pattern as `requires-question`. Documented it as the boundary between
   thinking about a ticket and working on it.
2. **`core/lane/builtin/10-todo.md`** — added `requires-specified: true`. Todo
   is the boundary; no other built-in lane carries it.
3. **`core/gate/gate.go`** — replaced `leavingBacklog` with
   `enteringSpecifiedZone`, computed from a new `specifiedBoundary` helper:
   the display-order index of the first lane that declares
   `RequiresSpecified`. A move triggers the promotion check only when it
   starts before that boundary and lands at or after it (`from < boundary &&
   to >= boundary`) — a crossing, not a same-zone move. If no installed lane
   declares the boundary (`boundary == -1`), the gate falls back to the
   original literal check (`t.Status == "backlog" && req.To != "backlog"`),
   so an old or hand-built lane set still fails closed rather than letting
   everything through. `blockedBy` is still checked exactly once per move,
   either as part of the crossing or as part of the general forward-move
   check — never twice.
4. **`core/gate/gate_test.go`** (new) — 7 tests covering: backlog→todo missing
   fields refused (unchanged behaviour), backlog→todo all fields present
   allowed, backlog→in-progress skipping todo still refused, a title-only
   ticket moving from backlog into a custom pre-zone lane (allowed), that same
   ticket refused crossing from the custom lane into todo, a move fully inside
   the zone (todo→in-progress) not re-running the promotion check even when
   fields are missing, and the fallback rule still refusing when no lane
   declares `requires-specified`.
5. **`docs/AGENTS.md`** — added a short paragraph explaining lanes before the
   specified zone and that the gate fires on crossing into that zone, not on
   leaving a lane by name.

## Verification

`export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go vet ./... && go test ./... -count=1`:

```
?   	github.com/BeMuCa/jaira/cmd/jaira	[no test files]
ok  	github.com/BeMuCa/jaira/core/board	0.006s
ok  	github.com/BeMuCa/jaira/core/gate	0.113s
?   	github.com/BeMuCa/jaira/core/gitrepo	[no test files]
?   	github.com/BeMuCa/jaira/core/lane	[no test files]
ok  	github.com/BeMuCa/jaira/core/merge	0.080s
ok  	github.com/BeMuCa/jaira/core/project	0.022s
ok  	github.com/BeMuCa/jaira/core/release	0.006s
?   	github.com/BeMuCa/jaira/core/session	[no test files]
ok  	github.com/BeMuCa/jaira/core/ticket	0.028s
ok  	github.com/BeMuCa/jaira/core/validate	0.076s
ok  	github.com/BeMuCa/jaira/internal/cli	0.099s
ok  	github.com/BeMuCa/jaira/internal/tui	1.519s
?   	github.com/BeMuCa/jaira/scripts/iconpreview	[no test files]
?   	github.com/BeMuCa/jaira/scripts/shotgen	[no test files]
```

`grep -c "t.Status == \"backlog\" && req.To != \"backlog\"" core/gate/gate.go` → `1`
(the old condition survives only inside the documented fallback branch).

`grep -c "requires-specified" docs/AGENTS.md` → `1`.

**Manual check against a real built binary** — throwaway repo, `JAIRA_HOME`
redirected, a ticket created with only a title (goal, context, and
definition-of-done left empty):

```
$ jaira --json move C5C7KD --to todo
{
  "error": {
    "code": 3,
    "message": "cannot move C5C7KD to todo:\n  - goal is required before this ticket can leave the backlog\n  - context is required before this ticket can leave the backlog",
    "reason": "gate_refused",
    "violations": [
      {"code": "missing_field", "field": "goal", "message": "goal is required before this ticket can leave the backlog"},
      {"code": "missing_field", "field": "context", "message": "context is required before this ticket can leave the backlog"}
    ]
  }
}
exit=3
```

`assignee` was already filled by `create` (defaults to the git user), and the
scaffolded definition-of-done checklist item counts as present under the
existing `fieldFilled` rule for `DoD` (a checklist with items counts as
filled regardless of whether items are checked) — this is pre-existing
behaviour, unrelated to this change, and not something this plan touched. The
refusal still correctly names every field the ticket is genuinely missing
(`goal`, `context`) and refuses with exit code 3, matching the plan's
expectation.

## Deviations from Plan

None — plan executed exactly as written. One correction during manual
verification: the plan's create command was invoked as `jaira create
"<title>"` rather than `jaira new`, since `new` is not a subcommand of this
CLI (confirmed via `jaira --help`).

## Commits

- `5dba91d` — feat: let a lane declare where specification becomes mandatory
- `1e6c4b9` — feat: key the promotion gate on entering the specified zone
- `e3f9333` — docs: say a board can have lanes before the specified zone

## Self-Check: PASSED

- FOUND: core/lane/lane.go (RequiresSpecified field and parsing present)
- FOUND: core/lane/builtin/10-todo.md (requires-specified: true present)
- FOUND: core/gate/gate.go (specifiedBoundary, enteringSpecifiedZone present)
- FOUND: core/gate/gate_test.go (7 tests, all passing)
- FOUND: docs/AGENTS.md (specified-zone paragraph present)
- FOUND commit 5dba91d in `git log --oneline --all`
- FOUND commit 1e6c4b9 in `git log --oneline --all`
- FOUND commit e3f9333 in `git log --oneline --all`
