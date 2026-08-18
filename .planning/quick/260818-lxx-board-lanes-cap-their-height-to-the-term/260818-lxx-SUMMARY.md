---
quick_id: 260818-lxx
subsystem: internal/tui
tags: [tui, rendering, scroll, board, lane-focus]
dependency-graph:
  requires: []
  provides: [board-height-cap, lane-focus-scroll]
  affects: [internal/tui/view.go, internal/tui/lanefocus.go]
tech-stack:
  added: []
  patterns: ["clampBlock hard w x h clamp as a defensive backstop", "budget-based variable-height scroll fit (laneFocusFit) mirroring the board's fixed-height scroll math"]
key-files:
  created: []
  modified:
    - internal/tui/view.go
    - internal/tui/lanefocus.go
    - internal/tui/scroll_test.go
    - internal/tui/lanefocus_test.go
decisions:
  - "renderCard's flags line: w+24 -> w (stale ANSI-length fudge from before truncate became display-width based)"
  - "clampBlock added next to truncate as a hard w x h clamp, applied in renderColumn as a last line of defence"
  - "lane-focus scroll reuses m.scroll[lane.ID] rather than a new field, since the board already recomputes/clamps it for the focused lane every render"
metrics:
  duration: "~45 minutes"
  completed: 2026-08-18
---

# Quick Task 260818-lxx: a lane never renders taller than the window Summary

Fixed the stale `w+24` fudge on the board's flags line that let a flag-heavy
review card render 4-5 lines instead of 3 and push the whole board past the
terminal, added a defensive `clampBlock` hard clamp so a column can never
spill the layout again for any reason, and gave the single-lane view a real
height cap with cursor-following scroll (it previously rendered every ticket
in the lane unconditionally).

## What Was Built

**Task 1 — the board fits, and cannot un-fit again** (`internal/tui/view.go`)

- Root cause: `renderCard`'s flags line was truncated at `w+24` cells instead
  of `w`, a leftover from before `truncate` became `lipgloss.Width`-based.
  Changed to `w`. The `meta` line (handle + assignee), previously written
  with no truncation at all, got the same `truncate(meta, w)` treatment for
  the same reason.
- Added `clampBlock(s string, w, h int) string`: splits on `"\n"`, truncates
  each line to `w`, keeps at most `h` lines, guards `w <= 0`/`h <= 0` by
  returning `""`. Applied in `renderColumn` immediately before styling
  (`style.Render(clampBlock(body.String(), w, h))`) so a column can never
  grow the board past the window regardless of what caused the overflow.
- Added `flagHeavyReviewLane(t, n)` in `scroll_test.go` — a fixture that
  maximises every flag the board can show on a card (sign-off, someone
  else's assignee, executed-by, 3 commits, `Plan`/`DoD` progress, a long
  title) — and `TestBoardFitsTheTerminal`, which asserts the rendered board
  never exceeds `h` lines across `h in {40,30,24,20,16,12,10}` and
  `w in {60,80,120,200}`.

**Task 2 — the single-lane view fits, with cursor-following scroll**
(`internal/tui/lanefocus.go`)

- `renderLaneFocus` previously rendered every ticket in the lane in a plain
  loop with no height cap and no scroll — measured at 155 lines for 30
  tickets at `m.height = 30`.
- Now measures each card's actual rendered height (cards are variable
  height: flags and goal lines are conditional), computes a line budget from
  what the header above has already used and what the footer below reserves,
  and fits as many cards as the budget allows via a new helper,
  `laneFocusFit`, which reserves a line for each "+N more" indicator that
  ends up shown (a two-pass fixed point, since whether "more below" is
  needed can only be known after seeing what fits — the same shape
  `clipToWindow` already uses for its own footer).
- Scroll state is `m.scroll[col.lane.ID]` — no new field. The board already
  recomputes and clamps this per-lane offset on every render of the focused
  column, and lane focus only ever scrolls the lane it is showing, so the
  two views agree on where a long lane is scrolled to.
- Cursor visibility: jump back immediately if the cursor moved above the
  visible window (`cardIdx < first -> first = cardIdx`); page forward one
  card at a time otherwise until the cursor falls inside `[first, last)`.
- "+N more" indicators appear above and/or below whatever is cut off,
  matching the board's wording.
- The final render is wrapped in `clampBlock(result, m.width, m.height)` as
  a last line of defence, same as the board's columns.
- Added three tests in `lanefocus_test.go`: `TestLaneFocusFitsTheTerminal`
  (35-ticket flag-heavy review lane, asserts `<= h` for
  `h in {40,30,24,20,16,12}`), `TestLaneFocusKeepsTheCursorVisible` (walks
  the cursor through all 35 tickets at 120x20, asserts the selected
  ticket's handle is always visible), and `TestLaneFocusSaysWhatIsOffScreen`
  (asserts `"more"` appears when the cursor sits mid-list at 120x20).

## Deviations from Plan

None — plan executed exactly as written. The only implementation detail not
spelled out verbatim in the plan is the two-pass fixed point inside
`laneFocusFit` for deciding how many lines to reserve for off-screen
indicators; this directly implements the plan's "reserve one more line for
each off-screen indicator actually shown" instruction using the same
already-established pattern (`clipToWindow`'s footer computation) rather than
inventing a new one.

## Verification

- `go build ./...` — clean
- `go vet ./...` — clean
- `go test ./... -race` — all packages green, including the new
  `TestBoardFitsTheTerminal`, `TestLaneFocusFitsTheTerminal`,
  `TestLaneFocusKeepsTheCursorVisible`, `TestLaneFocusSaysWhatIsOffScreen`,
  and every pre-existing test (`TestSelectionStaysVisibleWhenScrollingDown`,
  `TestLaneFocusShowsOnlyItsOwnLane`, `TestLaneFocusEmptyLaneDoesNotPanic`,
  etc.) unchanged.
- `gofmt -l` clean on all four changed files.
- Rebuilt `~/.local/bin/jaira` via
  `go build -o ~/.local/bin/jaira ./cmd/jaira` for Berk to verify against a
  real 30+ ticket review lane.

## Commit

- `88b611c` — fix: a lane stops at the bottom of the window — the flags row
  fits the card again, and the single-lane view scrolls
  (`internal/tui/view.go`, `internal/tui/lanefocus.go`,
  `internal/tui/scroll_test.go`, `internal/tui/lanefocus_test.go`)

## Self-Check: PASSED
