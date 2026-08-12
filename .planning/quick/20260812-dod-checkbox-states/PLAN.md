---
slug: dod-checkbox-states
created: 2026-08-12
status: in-progress
---

# Three-state checklists, and a Plan section

Fix a proven correctness hole in checkbox parsing, then build the three-state
checklist and the separate `## Plan` section on top of it.

## Why

`parseCheckbox` (`core/ticket/schema.go:149`) recognises only `[ ]`, `[x]` and
`[X]`. Any other marker returns `ok=false`, and `ParseDoDItems` drops the item
from the list entirely rather than treating it as unfinished. `DoDComplete` then
counts only the items that survived, so a checklist with unfinished work reports
complete.

Reproduced end to end before writing this plan: a ticket whose Definition of Done
read

```
- [x] placeholder ticked
- [x] first item done
- [~] second item in progress
- [x] third item done
```

moved to the terminal `done` lane with exit 0 and no signal.

This has to be fixed before `[~]` is introduced, because introducing `[~]` on top
of the current parser turns an unfinished ticket into a passing one.

## What changes

1. `DoDItem` gains `State` (`todo` / `doing` / `done`). `Checked` is kept as a
   derived accessor so `gate.go`, `tickets.go` and the TUI keep working.
2. `parseCheckbox` recognises `[ ]` as todo, `[~]` as doing, `[x]`/`[X]` as done,
   and **any other single-character marker as todo** — never dropped.
3. `DoDComplete` counts `doing` as remaining. An item being worked on is not a
   met criterion.
4. A `## Plan` section is parsed separately from `## Definition of Done`, into
   `Ticket.PlanItems`, using the same checkbox grammar. The Plan records method
   (spec, design, test, implement); the DoD records acceptance criteria and is
   the only one that gates the terminal lane.

## Out of scope

- The `jaira dod` CLI verb (next task)
- TUI rendering of either checklist (next task)
- Removing `--signal` (separate decision)

## Verification

- A test reproducing the bypass above fails before the fix and passes after.
- `go test ./... -race`, `go vet ./...`, `gofmt -l .` all clean.
- The end-to-end CLI repro is re-run against a rebuilt binary and now refuses
  the move with exit 3.
