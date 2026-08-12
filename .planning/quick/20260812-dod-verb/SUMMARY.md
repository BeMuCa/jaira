---
slug: dod-verb
created: 2026-08-12
status: complete
---

# Summary

`jaira dod` writes checklist state; `core/ticket/checklist.go` does it surgically.

## What changed

- `Section` (`SectionDoD` / `SectionPlan`) and `SetItemState`, which replaces one
  marker character and nothing else. Indexes are counted within the addressed
  section, so a definition-of-done index can never reach a checkbox under another
  heading.
- Marking an item as in progress clears any other in-progress item in that
  section.
- `jaira dod <id> <n> [--plan] --doing|--done|--todo`, 1-based to match `jaira
  show`. Read-only unknown lanes are refused on this path too, matching `set`
  and `move`.

## Verification

Six tests in `core/ticket/checklist_test.go`, written first and failing to
compile. End-to-end against a real board:

```
=== dod 1 --done ===          → 1. [x] returns 429 above 100/min
=== plan 2 --doing ===          1. [ ] write the spec
                              → 2. [~] implement
=== out of range ===          exit=3  definition of done has 2 item(s); there is no item 9
=== two flags ===             exit=2  pass exactly one of --doing, --done or --todo
```

Diff of the ticket file after both writes, ignoring `updated-at`:

```
27c27
< - [ ] returns 429 above 100/min
> - [x] returns 429 above 100/min
36c36
< - [ ] implement
> - [~] implement
```

Two lines, both intended. The gate then refused `done` naming the outstanding
criterion, and accepted it once that criterion was ticked.

`gofmt -l .` empty, `go vet ./...` clean, `go test ./... -race` green.

## Open question raised by the verification

The ticket reached `done` while its Plan still read `[~] implement`. That follows
the decision that the Plan records method and does not gate — but a ticket that
is accepted while its own plan says it is mid-implementation is worth at least a
warning. Not decided; flagged for the TUI task.
