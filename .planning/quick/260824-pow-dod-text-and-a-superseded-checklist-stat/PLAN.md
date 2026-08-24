---
quick_id: 260824-pow
slug: dod-text-and-a-superseded-checklist-state
date: 2026-08-24
mode: quick
---

# Quick task — reword an item, and mark one superseded

`.planning/NEXT-STEPS.md` item 3.

## The defect

There is no way to reword a checklist item. When two criteria were redefined, the
only route was to tick the old ones with the proof "obsolete, replaced by item 6"
and append new ones. That leaves ticked items that were never met — and a tick
that can mean "never happened" is a tick that means nothing anywhere on the
board.

## Decisions (already made)

- `jaira dod <id> <n> --text "..."` rewrites one item's words, leaving its state
  and its proof alone.
- A fourth state `[-]`, superseded: it did not happen and will not. Distinct from
  done, and **not** blocking completion — that is the whole point.
- Both work on `--plan` too.

## What the code says

- States and `Marker()`: `core/ticket/schema.go:186-221`. `parseCheckbox` at
  `:325` maps the character back, defaulting unknown markers to todo.
- `Checked()` (`:236`) is called in 8 places. Completion is decided by
  `DoDComplete` (`:361`) and `PlanComplete` (`:393`), and the gate re-walks the
  items itself at `core/gate/gate.go:321` and `:371` to name what is still open.
- The TUI never writes a checklist state — `SetItemState` has no caller under
  `internal/tui` — so `[-]` renders through `Marker()` with no TUI change.
  `checklistProgress` (`internal/tui/view.go:366`) is the one exception.
- `replaceMarker` (`core/ticket/checklist.go:173`) changes exactly one character
  and leaves the rest of the line, trailing carriage return included, alone. The
  new text writer must keep that property.

## Tasks

1. `StateSuperseded`: the constant, `String()` → "superseded", `Marker()` → "-",
   and `parseCheckbox` reading `-`. Verify: round-trip test, and an unknown
   marker still reads as todo.
2. `DoDItem.Settled()` — done or superseded, i.e. not outstanding. `Checked()`
   stays exactly what it is (only `[x]`), because a superseded item was not
   achieved and nothing may report that it was. `DoDComplete`, `PlanComplete`
   and the two gate loops switch to `Settled`. Verify: a DoD of one done and one
   superseded item completes; the gate does not list the superseded one as open.
3. `checklistProgress` counts settled, with the comment saying why: the counter
   exists to show what is outstanding, and it follows the gate. Otherwise a
   ticket that can be completed sits at 4/5 forever.
4. `ticket.SetItemText(body, sec, i, text)`: replace everything after `] `,
   keeping indentation, bullet, marker and trailing CR. Whitespace collapsed, as
   `SetItemProof` does, so a newline cannot re-parse as a new checkbox. Verify:
   state and proof survive, the line is otherwise byte-identical.
5. CLI: `--superseded` joins the mutually exclusive state flags; `--text` behaves
   like `--proof` (valid alone, combinable with a state, refused with
   `--add`/`--option`). Help text says what superseded means and that it does not
   block completion.

## Done when

An item can be reworded without touching its state, a superseded item does not
read as achieved, `dod --done` on all remaining items still completes the ticket,
and the round-trip tests cover `[-]`.
