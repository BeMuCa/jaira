---
quick_id: 260824-piz
slug: show-whether-a-lane-has-been-run
date: 2026-08-24
status: complete
---

# A lane that has not been run says so

`.planning/NEXT-STEPS.md` item 2, closed.

## What changed

**`core/gate.OutputOwed(l, t) []string`** — the fields a lane declares it
produces that the ticket has not filled. It is not new logic: the gate already
computed exactly this to refuse a forward move, with `fieldFilled` private, so
the enforcement loop now iterates `OutputOwed` and there is one implementation.
A lane the ticket opted out of (`requires-option`) owes nothing, for the same
reason the move is not refused. Nil lane and nil ticket return nil — the
renderers call this on whatever the board holds.

**`◇ unworked` on the card**, in `renderCard` and in lane focus, which carry the
same flag block by design. Shown only for `l.Agentic`: the claim is specifically
"this lane's agent has not run", and `todo`/`backlog` produce nothing so they
would never light up anyway — but a hand-written non-agentic lane may declare
output, and that is not the same statement.

**`next` prefers an unworked ticket** inside a lane. `sortByProgress` now sorts
by lane precedence, then unworked-before-worked, then oldest ID. One change
covers plain `next`, `--lane <id>` and `--per-lane` (which takes each lane's
first ticket). Across lanes nothing moves: precedence still finishes in-flight
work first. Within a lane, a ticket that has produced the lane's output is
waiting to be *moved*, not worked — handing it back does the lane's work twice.

## Verified

- `go test ./... -race`, cache cleared: green. `gofmt -l core internal` lists
  exactly the two documented files.
- New tests: `core/gate/outputowed_test.go` asserts `OutputOwed` and the gate
  refusal name the same fields (two implementations would drift into a card that
  looks worked beside a move that is refused), plus the skipped-lane and nil
  cases. `internal/tui/unworked_test.go` covers the flag appearing and
  disappearing on the board and in lane focus, and that `todo`/`backlog` cards
  are never marked. `internal/cli/unworked_test.go` gives the worked ticket the
  older ID, so it would win the old tiebreak — the test fails without the change.
- On his real board, read-only (`next --per-lane`, `next --lane critique --all
  --json`): the binary runs and the critique lane hands back an unworked ticket
  first.
- The board flag is verified through the renderer tests, not by eye: the TUI
  needs a TTY.

## What it showed on his board

All **12** tickets in `critique` are unworked — the lane has produced a
`review-summary` for none of them. The two that had been through it are no longer
in the lane. That is the fact the flag exists to make visible.

## Deliberately not done

No `output_owed` in `ticketJSON`, no extra column in `--per-lane`: the ordering
already hands back an unworked ticket first, which is what item 2 asks for. Add
the list when something needs it rather than the choice.
