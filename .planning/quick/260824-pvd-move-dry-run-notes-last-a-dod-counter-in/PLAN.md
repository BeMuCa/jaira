---
quick_id: 260824-pvd
slug: move-dry-run-notes-last-a-dod-counter-in-list-a-stale-claim-warning
date: 2026-08-24
mode: quick
---

# Quick task — the small package

`.planning/NEXT-STEPS.md` item 4. Four separate frictions, one commit each.

## 1. `move --dry-run`

Run the gates, print the violations, write nothing. An agent created and deleted
a throwaway ticket to find out what a gate would do, and one of the user's own
tickets was moved by accident probing a gate that turned out not to refuse.

**The trap:** the real move stages `--what`/`--why`/`--resolves`/the assignee
*before* the gate runs (`internal/cli/flow.go`, so a gate can require them). A
dry run that reused that path would leave those fields behind on a move it then
refused. It stages into the in-memory doc and re-decodes instead, so the gate
sees exactly the ticket it would have seen and nothing reaches the disk.

Same exit code as the real refusal, so a caller branches identically.

## 2. `show --notes-last N`

One ticket is 263 lines / 56 KB because notes append without bound.
`printDetail` dumps the whole body last. The user's register argues for counting
them and showing the last few.

## 3. A DoD counter in `jaira list`

The listing prints handle, title and assignee per lane. `dod_complete` and the
item counts are already in the JSON. Add `DoD n/m` to the human line.

## 4. A stale claim warning

`ClaimActive` silently treats a claim older than `ClaimTTL` as abandoned.
Nothing says so when it happens. Warn where a stale claim is stepped over.

## Not in scope, deliberately

The lane warning already goes to stderr (`internal/cli/root.go:261-263`). The
other session's feedback was wrong about it.

## Done when

Each of the four works on the running binary, `go test ./... -race` is green with
the cache cleared, and `gofmt -l core internal` still lists exactly
`core/gate/gate.go` and `internal/cli/tickets.go`.
