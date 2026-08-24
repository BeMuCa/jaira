---
quick_id: 260824-pow
slug: dod-text-and-a-superseded-checklist-state
date: 2026-08-24
status: complete
---

# An item can be reworded, and one can be retired

`.planning/NEXT-STEPS.md` item 3, closed.

## The defect

There was no way to reword a checklist item. When two criteria were redefined,
the only route was ticking the stale ones with the proof "obsolete, replaced by
item 6" and appending new ones — which leaves `[x]` beside work nobody did, and
then no tick anywhere on the board means anything.

## What changed

**`[-]`, a fourth state.** `StateSuperseded`: it did not happen and will not.
`Checked()` stays exactly what it was — only `[x]` — because a superseded item
was not achieved and nothing may report that it was. The new `Settled()` answers
the other question, "is this still outstanding", and that is what completion
uses: `DoDComplete`, `PlanComplete` and the two gate loops that name what is
still open. So a superseded item does not block a terminal lane and is never
listed as work to do.

**`jaira dod <id> <n> --text "..."`** rewrites one item's words and leaves its
marker and its proof line alone. `SetItemText`/`replaceText` keep the same
surgical contract `replaceMarker` has in the other direction: indentation,
bullet, marker and any trailing carriage return survive. The text is
whitespace-collapsed, as `SetItemProof` already does, so a newline cannot
re-parse as a second checkbox — there is a test that tries exactly that.

**The counter follows the gate.** `checklistProgress` counts settled items, not
ticks. Otherwise a ticket that can be completed sits at 4/5 for good.

**Docs**: the state list in README and `docs/COMMANDS.md`, plus both agent
contracts (`SKILL.md`, `docs/AGENTS.md`) — the workaround was an agent's, so the
replacement has to be where the agent reads.

## One deliberate behaviour change

`- [-] cancelled?` used to parse as **todo** — an existing test asserted it, as an
example of an unrecognised marker. It now parses as superseded, which is the
point of the item. Any ticket already carrying `[-]` by hand changes meaning:
that item stops blocking completion. The test was rewritten to use `[/]`, and
unknown markers still read as todo.

## Verified

- `go test ./... -race`, cache cleared: green. `gofmt -l core internal` lists
  exactly the two documented files.
- New tests: `core/ticket/superseded_test.go` (round-trip both ways, not-checked
  but settled, completion with and without an open item beside it, the plan
  checklist, unknown markers still todo), `core/ticket/itemtext_test.go` (state
  and proof survive, the rest of the line survives, CRLF, the smuggled-checkbox
  attempt, empty and out-of-range), `internal/cli/checklist_test.go` (reword
  through the real command, supersede then complete, two state flags refused).
- On the running binary, in a throwaway repo: ticked item 1 with a proof,
  reworded it — `- [x] first thing, reworded` with its proof line intact —
  superseded item 2, and the ticket reports `dod_complete: true`.

## Not done

`checklistJSON` still emits `"done": false` for a superseded item, with
`"state": "superseded"` beside it. `done` means done, and a consumer that wants
"is anything outstanding" has `dod_complete`. Worth revisiting only if an agent
is seen re-opening a retired item.
