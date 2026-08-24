---
quick_id: 260824-piz
slug: show-whether-a-lane-has-been-run
date: 2026-08-24
mode: quick
---

# Quick task — show whether a lane has been run

`.planning/NEXT-STEPS.md` item 2.

## The defect

Nine tickets sat in `critique` uncritiqued and rendered identically to the two
that had been through it. Today the only way to tell is to read `review-summary`
and check it for emptiness — per ticket, by hand.

## Decision (already made): derive it, no new field

A lane's `output-produces` already says what leaving it requires. A ticket in a
lane whose output contract is unsatisfied has not been worked yet.

Verified in the code: `core/gate/gate.go:389-398` already computes exactly this
when refusing a forward move, with `fieldFilled` (`:433`) private. So this is one
extraction plus two call sites, not new logic.

Verified in `core/lane/builtin/`: no non-agentic built-in declares
`output-produces` (backlog, todo, human, signoff, done, blocked declare none), so
the "do not light up every card in todo" trap cannot fire from the built-ins. The
display guard on `Agentic` is still added, because a hand-written non-agentic
lane may declare output and the flag's claim is specifically "this lane's agent
has not run".

## Tasks

1. `core/gate/gate.go`: add `OutputOwed(l *lane.Lane, t *ticket.Ticket) []string`
   — the fields the lane declares it produces that are still empty; nil for a
   nil lane/ticket and for a lane skipped by `requires-option`. Rewrite the
   enforcement loop at `:389` to iterate it, so there is one implementation.
   Verify: gate tests still pass, a new test asserts the refusal and `OutputOwed`
   agree.
2. Card flag in both renderers — `internal/tui/view.go` `renderCard` and
   `internal/tui/lanefocus.go`, which carry the same flag block by design
   ("Same computation as renderCard's flags"). Guarded on `l.Agentic`. Verify: a
   test at both call sites, flag present when the output is owed and absent once
   it is filled.
3. `internal/cli/flow.go` `sortByProgress`: at equal lane precedence, a ticket
   that owes its lane's output sorts first. One change serves plain `next`,
   `next --lane <id>` and `--per-lane` (which takes each lane's first ticket).
   Rationale: across lanes, "finish in-flight work first" still holds — that is
   precedence; within one lane, a ticket whose output is already produced is
   waiting to be *moved*, not worked, so handing it back means re-doing it.
   Verify: a test with two tickets in one lane, one worked, one not.

## Out of scope, deliberately

- No `output_owed` in `ticketJSON` and no extra column in `--per-lane`: the
  ordering already puts an unworked ticket first, which is what the item asks
  for. Add it when something needs the list rather than the choice.

## Done when

A ticket in an agentic lane with its output unfilled is visibly distinct from one
that has it, on the board and in lane focus; `next --lane <id>` hands back an
unworked one first; tests cover both states; `go test ./... -race` green and
`gofmt -l core internal` unchanged.
