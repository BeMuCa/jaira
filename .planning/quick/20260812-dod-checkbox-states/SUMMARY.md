---
slug: dod-checkbox-states
created: 2026-08-12
status: complete
---

# Summary

Three checklist states, and a Plan section separate from the definition of done.

## What changed

- `State` (`todo` / `doing` / `done`) with `Marker()`; `DoDItem.Checked` became
  the derived method `Checked()`, true only for `done`.
- `parseCheckbox` reads any `[c] ` marker. `" "` is todo, `"~"` is doing, `"x"`
  and `"X"` are done, **anything else is todo** rather than being discarded.
- `DoDComplete` counts a doing item as remaining.
- `ParsePlanItems` reads a `## Plan` / `## Steps` / `## Vorgehen` section into
  `Ticket.PlanItems`, using the same grammar. Only the DoD gates the terminal
  lane.
- `internal/cli/tickets.go` renders the marker rather than a two-state box.

## Verification

Five tests in `core/ticket/schema_test.go`, written first and failing to compile
against the old `DoDItem`.

The original end-to-end bypass was re-run against a rebuilt binary:

```
## Definition of Done

- [x] first item done
- [~] second item in progress
- [x] third item done

=== move to done : exit=3 ===
jaira: cannot move 6AP7VC to done:
  - definition of done is not met: second item in progress
```

Before the fix the same input gave `exit=0`.

`gofmt -l .` empty, `go vet ./...` clean, `go test ./... -race` green.

## Note for the next task

The first attempt at that repro appended checkboxes to the end of the file, where
they landed under `## Notes` and were correctly ignored. Checkbox position is
load-bearing: a checklist only counts inside a recognised heading. Anything that
writes items — the coming `jaira dod` verb, the TUI editor — has to insert into
the right section rather than appending.
