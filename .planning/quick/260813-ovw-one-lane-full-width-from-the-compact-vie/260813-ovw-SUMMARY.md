---
phase: quick/260813-ovw
plan: 01
subsystem: tui
tags: [tui, pipeline, board, lane-focus]
requires: []
provides: [modeLaneFocus]
affects: [internal/tui/pipeline.go, internal/tui/model.go, internal/tui/view.go]
tech-stack:
  added: []
  patterns: ["mode returns to the mode it was opened from, via Model.detailFrom"]
key-files:
  created:
    - internal/tui/lanefocus.go
    - internal/tui/lanefocus_test.go
  modified:
    - internal/tui/model.go
    - internal/tui/pipeline.go
    - internal/tui/pipeline_test.go
    - internal/tui/view.go
    - README.md
    - docs/COMMANDS.md
decisions:
  - "Detail pane now remembers the mode it was opened from (Model.detailFrom) instead of hardcoding modeBoard on the way back, so it works from both the board and the new single-lane view without a second detail screen."
  - "renderLaneCard duplicates renderCard's flag computation rather than sharing it, because task 1's file list did not include view.go and the two functions live in different files by design (lanefocus.go vs view.go)."
metrics:
  duration: "~35m"
  completed: "2026-08-13"
---

# Phase quick/260813-ovw Plan 01: One lane, full width, from the compact view Summary

Enter on a step in the compact view now opens that lane alone, filling the
screen with more per ticket (goal, timespan, and the same status flags the
board shows) than a 22-34 column board card can carry — instead of throwing
the user sideways into the multi-column board.

## What was built

- **`modeLaneFocus`** (`internal/tui/model.go`): new mode in the iota block.
- **`internal/tui/lanefocus.go`**: `laneFocusKey` (navigation: j/k move ticket,
  g/G first/last, h/l change lane without leaving the view, enter opens
  detail, esc/v back to the compact view) and `renderLaneFocus` /
  `renderLaneCard` (the full-width render: lane name + count + agentic tier
  or read-only marker in the header, then one card per ticket showing title,
  handle, assignee, timespan, the board's flags, and the goal).
- **`pipelineKey`** (`internal/tui/pipeline.go`): `enter` now sets
  `m.mode = modeLaneFocus` instead of `modeBoard`; the stale comment describing
  the old behavior was rewritten.
- **`Model.detailFrom`** (`internal/tui/model.go`): the detail pane records
  which mode opened it (`openDetail` sets `m.detailFrom = m.mode`), and
  esc/q/enter/j/k in `modeDetail` return to `m.detailFrom` rather than a
  hardcoded `modeBoard`. Zero value is `modeBoard`, so every pre-existing
  caller is unaffected — this only changes behavior for the new lane-focus
  entry path.
- **`?` help screen** (`internal/tui/view.go`, `renderHelp`): new "Compact
  view" section documenting `v`, `enter`, `h l ← →`, `esc v` — none of these
  were documented in the help screen before this plan.
- **`README.md` / `docs/COMMANDS.md`**: the compact-view key line now says
  `enter` opens the step "full width" rather than implying the board, and a
  new line lists the keys available once inside a lane.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 2 - missing critical functionality] `view.go`'s render dispatch had no case for `modeLaneFocus`**
- **Found during:** Task 1
- **Issue:** Task 1's `<files>` list did not include `view.go`, but without a
  `case modeLaneFocus: return m.renderLaneFocus()` in `render()`, the new mode
  would silently fall through to `renderBoard()` — the exact bug the plan
  exists to fix.
- **Fix:** Added the one-line case to `view.go`'s `render()` switch, in the
  task 1 commit (view.go is in task 2's file list for the help-screen change,
  so it was already destined to be touched by this plan; this pulls one line
  of that forward).
- **Files modified:** internal/tui/view.go
- **Commit:** b4a7bae

**2. [Rule 1 - bug from the plan's own behavior change] `TestPipelineNavigationAndOpening` asserted the old behavior**
- **Found during:** Task 1
- **Issue:** The existing test in `internal/tui/pipeline_test.go` (not listed
  in task 1's `<files>`) asserted `enter` opened `modeBoard`. Changing
  `pipelineKey` without updating it would leave a red test.
- **Fix:** Updated the assertion to expect `modeLaneFocus`.
- **Files modified:** internal/tui/pipeline_test.go
- **Commit:** b4a7bae

**3. [Plan inaccuracy — adapted, not forced through] "correct the compact view's entry [in the `?` help screen] if it claims enter opens the board"**
- **Found during:** Task 2
- **Issue:** Read `renderHelp` in full before editing: it never mentioned the
  compact view, `v`, or `enter`-opens-a-step at all. There was no existing
  entry to correct.
- **Adaptation:** Added a new "Compact view" section documenting `v`, `enter`,
  `h l ← →`, and `esc v`, which satisfies the plan's underlying intent (the
  help screen should describe the single-lane view and not claim enter opens
  the board) without inventing a correction to text that did not exist.
- **Files modified:** internal/tui/view.go
- **Commit:** 41985bb

## Verification

Automated (run verbatim, both passed):

```
$ export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go vet ./... && go test ./... -count=1
?   	github.com/BeMuCa/jaira/cmd/jaira	[no test files]
ok  	github.com/BeMuCa/jaira/core/board	0.022s
ok  	github.com/BeMuCa/jaira/core/gate	0.217s
?   	github.com/BeMuCa/jaira/core/gitrepo	[no test files]
ok  	github.com/BeMuCa/jaira/core/identity	0.015s
ok  	github.com/BeMuCa/jaira/core/lane	0.434s
ok  	github.com/BeMuCa/jaira/core/merge	0.114s
ok  	github.com/BeMuCa/jaira/core/project	0.016s
ok  	github.com/BeMuCa/jaira/core/release	0.013s
?   	github.com/BeMuCa/jaira/core/session	[no test files]
ok  	github.com/BeMuCa/jaira/core/ticket	0.022s
ok  	github.com/BeMuCa/jaira/core/validate	0.084s
ok  	github.com/BeMuCa/jaira/internal/cli	0.200s
ok  	github.com/BeMuCa/jaira/internal/tui	1.884s
?   	github.com/BeMuCa/jaira/scripts/iconpreview	[no test files]
?   	github.com/BeMuCa/jaira/scripts/shotgen	[no test files]
```

```
$ grep -c "lane" README.md docs/COMMANDS.md
README.md:16
docs/COMMANDS.md:8
```

**Human check — NOT run, reported as unverified (no interactive terminal in
this environment):** open the board, press `v`, pick a step, press enter, and
confirm the lane fills the screen and reads well at both a narrow and a wide
terminal. This must be done manually before treating the feature as visually
confirmed.

## Self-Check

- FOUND: internal/tui/lanefocus.go
- FOUND: internal/tui/lanefocus_test.go
- FOUND commit b4a7bae (feat: opening a step from the compact view fills the screen with that lane)
- FOUND commit 41985bb (docs: say the single-lane view exists and what its keys do)

## Self-Check: PASSED
