---
quick_id: 260819-djm
subsystem: tui
tags: [keybindings, lane-focus, board-hints]
requires: []
provides: [lane-focus-q-goes-back, board-names-v, off-screen-notice-removed]
affects: [internal/tui/model.go, internal/tui/lanefocus.go, internal/tui/view.go]
tech-stack:
  added: []
  patterns: ["m.key(tea.KeyPressMsg{...}) drives dispatch tests, not the mode-local key handler directly, when the bug lives in the dispatch"]
key-files:
  created: []
  modified:
    - internal/tui/model.go
    - internal/tui/lanefocus.go
    - internal/tui/lanefocus_test.go
    - internal/tui/view.go
    - internal/tui/render_test.go
    - internal/tui/hints_test.go
decisions:
  - "q in modeLaneFocus now falls through to laneFocusKey instead of quitting — lane focus is one level deeper than the compact view it opens from, so q there goes back, not out"
  - "board's off-screen lane count is no longer reported at all — matches lane focus, which never reported it either"
metrics:
  duration: "~35 min"
  completed: 2026-08-19
---

# Quick Task 260819-djm: `q` is always one back, the board names `v`, the off-screen notice goes Summary

`q` in the compact view's lane-focus screen no longer quits the whole program — it falls through to the same back-case as `esc`/`v`, landing on the compact view, because lane focus sits one level deeper than the shell that `q` quits to everywhere else.

## What Was Built

**Task 1 — `q` in lane focus goes back, not out.** `model.go`'s `modeLaneFocus` case used to intercept `q` (alongside `ctrl+c`) and quit before `laneFocusKey` ever saw it. Narrowed that branch to `ctrl+c` only, so `q` reaches `laneFocusKey`, where it now joins `esc`/`v` in the back case. Footer changed from `"enter open · esc/v back · q quit"` to `"enter open · q/esc/v back"` (lane focus can no longer quit, so no quit hint). Help screen's compact-view key label changed from `"esc v"` to `"q esc v"`.

**Task 2 — the board's hint bar names `v`, off-screen notice removed.** `statusBar` dropped the `hidden` variable and its `"N lane(s) off-screen"` warning entirely — lanes now scroll silently, same as lane focus already did. `"v compact"` was added to the hint-bar `keys` slice (paired with the compact view's own `"v full board"`). `renderBoard`'s comment was rewritten to state what's true now (off-screen count deliberately not reported) instead of the removed behavior.

## Deviations from Plan

None — plan executed exactly as written. The `start`/`end` parameters on `statusBar` are now unused inside the function body (only `sbLines` is downstream of the removed `hidden` logic); left them in place since the plan's action only specified deleting the `hidden` variable and its `if` block, not the function signature, and Go permits unused parameters.

## Verification

- `go build ./...` — clean
- `go test ./... -race` — all packages `ok` (internal/tui 35s, internal/cli 6s, others cached/fast)
- `gofmt -l internal/tui/` — empty (clean). `gofmt -w` was run on `view.go` since Task 2 edited it, which also picked up the pre-existing stray EOF blank line noted in the plan's baseline. `internal/cli/tickets.go`'s pre-existing gofmt diff was left untouched, per the plan.
- Manual smoke test of the real binary via a disposable tmux session against a throwaway repo:
  - Board hint bar shows `v compact`; at 40×20 with default lanes, no `"off-screen"` text anywhere.
  - `v` → compact view → `enter` → lane focus shows footer `"enter open · q/esc/v back"`.
  - `q` in lane focus landed back on the compact view (session stayed alive, footer read `"...v full board · q quit"` again).
  - `q` on the compact view quit the program (tmux session ended), confirming that path is unchanged.

## Self-Check: PASSED

- FOUND: internal/tui/model.go (modified, committed 98cc92f)
- FOUND: internal/tui/lanefocus.go (modified, committed 98cc92f)
- FOUND: internal/tui/lanefocus_test.go (modified, committed 98cc92f)
- FOUND: internal/tui/view.go (modified, committed c02e49c)
- FOUND: internal/tui/render_test.go (modified, committed c02e49c)
- FOUND: internal/tui/hints_test.go (modified, committed c02e49c)
- FOUND commit 98cc92f in git log
- FOUND commit c02e49c in git log
