---
phase: quick-260821-u7h
plan: 01
subsystem: ui
tags: [tui, bubbletea, lipgloss, board, kanban]

requires:
  - phase: quick-260821-t7s
    provides: "the versionLine field on Model and its statusBar append (kept last, extended above)"
provides:
  - "z toggle that hides empty lanes on the multi-column board, session-only"
  - "a stored, margin-respecting scroll window (boardFit/windowStart/scrolloff) that replaces the stateless cursor-pinned window"
  - "a status-bar notice naming hidden-lane count and the on-screen window range"
affects: [tui]

actuals:
  tokens: 9640
  tasks: 2
  commits: 2

tech-stack:
  added: []
  patterns:
    - "boardFit(width) is the single decision point for which lanes render and how wide; empty-lane filtering happens first, the window applies to what remains"
    - "windowStart/scrolloff are pure, Model-free functions with their own table tests, following laneFocusFit's precedent in lanefocus.go"
    - "render-writes-view-state: boardFit mutates m.laneStart during render, the same pattern renderColumn already uses for m.scroll"

key-files:
  created:
    - internal/tui/lanewindow_test.go
  modified:
    - internal/tui/model.go
    - internal/tui/view.go
    - internal/tui/hints_test.go
    - internal/tui/updatecheck_test.go
    - README.md
    - docs/COMMANDS.md

key-decisions:
  - "hideEmpty and laneStart are session-only Model fields, not persisted — the board's view state (filter, cursor, per-lane scroll) has never survived a restart, and a saved boolean would be a second write path bought for nothing"
  - "boardFit's keep-rule (i == m.laneIdx) always renders the focused lane even when empty, so a ticket moved there by another shell is never hidden out from under the cursor"
  - "scrolloff is a margin (min(2, (perScreen-1)/2)), not hard centring, so it degrades gracefully as fewer columns fit instead of demanding room that does not exist"
  - "the window's clamp lives inside windowStart and runs on every render from the live width, so tea.WindowSizeMsg needs no resize handler of its own"

requirements-completed: [QT-260821-u7h]

coverage:
  - id: D1
    description: "z toggles empty lanes on/off; the board states the hidden count and that z restores them; the toggle changes only the picture (m.cols ticket counts and lane-ID set unchanged)"
    requirement: QT-260821-u7h
    verification:
      - kind: unit
        ref: "internal/tui/lanewindow_test.go#TestZHidesEmptyLanesAndNotices"
        status: pass
      - kind: unit
        ref: "internal/tui/lanewindow_test.go#TestSecondZRestoresEverything"
        status: pass
      - kind: unit
        ref: "internal/tui/lanewindow_test.go#TestToggleIsDisplayOnly"
        status: pass
    human_judgment: false
  - id: D2
    description: "the focused lane always renders even when empty; moveLane skips hidden lanes on the board but leaves lane-focus's h/l untouched"
    requirement: QT-260821-u7h
    verification:
      - kind: unit
        ref: "internal/tui/lanewindow_test.go#TestFocusedEmptyLaneStaysVisibleAndIsNotCountedHidden"
        status: pass
      - kind: unit
        ref: "internal/tui/lanewindow_test.go#TestLWithToggleOnSkipsHiddenLanes"
        status: pass
      - kind: unit
        ref: "internal/tui/lanewindow_test.go#TestLaneFocusIgnoresTheToggle"
        status: pass
    human_judgment: false
  - id: D3
    description: "columns never render below minColWidth; overflow scrolls with a stored, margin-holding window instead of shrinking every column"
    requirement: QT-260821-u7h
    verification:
      - kind: unit
        ref: "internal/tui/lanewindow_test.go#TestScrolloffTable"
        status: pass
      - kind: unit
        ref: "internal/tui/lanewindow_test.go#TestWindowStartTable"
        status: pass
      - kind: unit
        ref: "internal/tui/lanewindow_test.go#TestMinColWidthKeepsHandleAndAFlagIntact"
        status: pass
    human_judgment: false
  - id: D4
    description: "the two-lane margin holds the cursor mid-window, advances one lane at a time once exhausted, and clamps without padding at both ends"
    requirement: QT-260821-u7h
    verification:
      - kind: unit
        ref: "internal/tui/lanewindow_test.go#TestWalkingLKeepsTheTwoLaneMargin"
        status: pass
      - kind: unit
        ref: "internal/tui/lanewindow_test.go#TestWalkingToEitherEndClampsWithoutPadding"
        status: pass
      - kind: unit
        ref: "internal/tui/lanewindow_test.go#TestSteadyInteriorAt170Columns"
        status: pass
    human_judgment: false
  - id: D5
    description: "a wide terminal (everything fits) and the 125-column toggle-on case (empty-lane toggle alone solving it) both render with no window indicator"
    requirement: QT-260821-u7h
    verification:
      - kind: unit
        ref: "internal/tui/lanewindow_test.go#TestWideBoardHasNoWindowIndicator"
        status: pass
      - kind: unit
        ref: "internal/tui/lanewindow_test.go#TestNarrowBoardWithToggleOnHasNoWindow"
        status: pass
    human_judgment: false
  - id: D6
    description: "resize re-clamps the stored window on the next frame with no dedicated handler, in both directions"
    requirement: QT-260821-u7h
    verification:
      - kind: unit
        ref: "internal/tui/lanewindow_test.go#TestResizeReclampsInBothDirectionsWithNoHandler"
        status: pass
    human_judgment: false
  - id: D7
    description: "an unrelated ticket moved by another shell leaves the window alone; the selected ticket moved off-screen by another shell scrolls the window to include it"
    requirement: QT-260821-u7h
    verification:
      - kind: unit
        ref: "internal/tui/lanewindow_test.go#TestExternalMoveOfUnrelatedTicketDoesNotScrollTheWindow"
        status: pass
      - kind: unit
        ref: "internal/tui/lanewindow_test.go#TestExternalMoveOfSelectedTicketScrollsToInclude"
        status: pass
    human_judgment: false
  - id: D8
    description: "eyeball verification at 125 and 240 columns matches the plan's described shape exactly (window indicator, sliding margin, wide-terminal parity)"
    requirement: QT-260821-u7h
    verification: []
    human_judgment: true
    rationale: "Rendered TUI layout at a real terminal width is a visual/UX judgment call the executor eyeballed via a throwaway test harness (removed after use, not committed); a human should confirm it reads correctly in an actual terminal before shipping."

duration: ~50min
completed: 2026-08-21
status: complete
---

# Phase quick-260821-u7h Plan 01: Narrow-terminal board — hide empty lanes, window the columns Summary

**`z` hides lanes with no tickets on the board (session-only), and the multi-column view now scrolls a stored two-lane-margin window instead of shrinking every column below readability.**

## Performance

- **Duration:** ~50min
- **Completed:** 2026-08-21
- **Tasks:** 2/2 completed
- **Files modified:** 7 (1 created)

## Accomplishments

- `z` toggles whether lanes holding zero tickets are drawn on the multi-column board; the board always states the hidden count and that `z` brings them back, and the toggle changes nothing about tickets, moves, or `m.cols` — display only.
- `boardFit` is the single function that decides which lanes render and how wide: it drops empty lanes first (when `hideEmpty` is set, keeping the focused lane regardless so the cursor is never left pointing at something invisible), then applies the scroll window to what remains.
- The scroll window is now model state (`m.laneStart`), moved only when the active lane runs out of a two-lane margin (`scrolloff`/`windowStart`, both pure and table-tested) — columns hold still while the cursor walks the interior, and slide one at a time once the margin is exhausted. Both ends clamp inside the board with no blank padding.
- The status bar's notice line grows a window half (`‹ h  l ›  lanes 4–8 of 10`) joined to the hidden-lanes half with " · ", counted so the two halves never double-count the same lanes.
- A resize re-clamps `m.laneStart` on the next render purely from the live width recomputation — no `tea.WindowSizeMsg` handler was added.
- A wide terminal (everything fits) and the user's actual 125-column terminal with the toggle on (the five real lanes fit) both render with no window indicator at all — untouched from before this plan.
- `minColWidth`'s value (22) is unchanged; it now carries its derivation as a comment instead of a bare constant.

## Task Commits

1. **Task 1: `z` hides empty lanes, end to end** - `776feef` (feat, tdd, tracer)
2. **Task 2: window the columns instead of shrinking them, and say where you are** - `5274cc5` (feat, tdd)

_Both tasks were TDD: tests (`internal/tui/lanewindow_test.go`) were written and run alongside the implementation in each commit, per the plan's `<behavior>` lists._

## Files Created/Modified

- `internal/tui/lanewindow_test.go` - new: the toggle's behavior list (Task 1) plus `scrolloff`/`windowStart` table tests and the window's render-level behavior (Task 2)
- `internal/tui/model.go` - `hideEmpty`/`laneStart` Model fields, `z` key binding, `toggleEmptyLanes()`, `moveLane`'s hidden-lane skip
- `internal/tui/view.go` - `boardWindow`/`boardFit`, `scrolloff`/`windowStart`, `statusBar` now takes a `boardWindow` and appends `boardNotice(win)`, `renderHelp` names `z`, `minColWidth`'s derivation comment
- `internal/tui/hints_test.go` - the two hardcoded key-hint lists now include `"z hide empty"`
- `internal/tui/updatecheck_test.go` - one call site updated to the new `statusBar(boardWindow)` signature
- `README.md` - `## Keys` board block lists `z`; one sentence on the narrow-terminal scroll behaviour
- `docs/COMMANDS.md` - `## Keys` board block lists `z`

## Decisions Made

- `hideEmpty`/`laneStart` are session-only, not persisted (see key-decisions above) — matches how every other board view state (filter, cursor, scroll) already behaves.
- The window notice's total is counted over the post-hiding lane count (`len(win.cols)`), not the raw lane count, so the hidden-lanes half and the window half of the notice never contradict each other when both apply.
- `scrolloff` is expressed as a margin (`min(2, (perScreen-1)/2)`) rather than hard centring, so it degrades to 1 or 0 at a narrower fit instead of becoming an impossible constraint.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - blocking] `updatecheck_test.go` called the old `statusBar(int, int)` signature directly**
- **Found during:** Task 1, immediately after changing `statusBar`'s signature to take a `boardWindow`
- **Issue:** `TestBoardStatusBarCarriesTheVersionIndicator` called `m.statusBar(0, 0)`, which no longer compiled once `statusBar` took a single `boardWindow` parameter — a build-breaking regression from the signature change, not something the plan's file list anticipated (it lists `model.go`, `view.go`, `lanewindow_test.go`, `hints_test.go`, README, COMMANDS, but not `updatecheck_test.go`).
- **Fix:** Changed the call to `m.statusBar(m.boardFit(m.width))`, computing a real window the same way `renderBoard` does.
- **Files modified:** `internal/tui/updatecheck_test.go`
- **Verification:** `go build ./...` succeeds; `TestBoardStatusBarCarriesTheVersionIndicator` still asserts the version indicator is present and passes under `-race -count=1`.
- **Committed in:** `776feef` (part of Task 1's commit)

---

**Total deviations:** 1 auto-fixed (1 Rule 3)
**Impact on plan:** Necessary to keep the build green after the intentional `statusBar` signature change; no scope creep, no behavior change to the test's own assertion.

## Issues Encountered

None beyond the deviation above. One test-writing correction worth recording for future readers: the plan's Task 1 `<behavior>` bullet "park the cursor on `blocked`, press `z`, and Blocked is still drawn" reads as if `toggleEmptyLanes`'s own cursor-relocation (which moves the cursor off an empty *focused* lane the moment the toggle turns on) would apply — but doing so there empties the very keep-rule the bullet is testing, since the cursor would no longer be on `blocked` by the time `boardFit` runs. The two mechanisms are for different situations: `toggleEmptyLanes`'s relocation guards the moment of the keypress itself (so pressing `z` while parked on an empty lane doesn't strand the cursor there forever), while `boardFit`'s keep-rule (`i == m.laneIdx`) is the general invariant that a cursor sitting on an empty lane — however it got there, including a direct assignment simulating a ticket move — still renders. The test (`TestFocusedEmptyLaneStaysVisibleAndIsNotCountedHidden`) exercises the keep-rule directly: toggle on first (while parked on a non-empty lane, so relocation never fires), then move the cursor onto `blocked`, then assert it renders with a hidden count of 4. Two separate tests (`TestTogglingOnMovesCursorOffEmptyLaneSearchingRightFirst` / `...LeftAsFallback`) cover the relocation mechanism itself.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

The multi-column board now scales from a handful of lanes to a large, open-ended custom-lane count without becoming unreadable: `z` handles the common case (the user's own board has four lanes with real work), and the scroll window handles the rest. No blockers for follow-on work. `docs/COMMANDS.md`'s lane-settings section and any future custom-lane tooling are unaffected — this plan touched only the multi-column board (`renderBoard`/`boardFit`), never `defaultBoardScreen` or the lane settings screen.

---
*Phase: quick-260821-u7h*
*Completed: 2026-08-21*

## Self-Check: PASSED

All files created/modified confirmed present on disk; both task commits
(`776feef`, `5274cc5`) confirmed present in `git log`.
