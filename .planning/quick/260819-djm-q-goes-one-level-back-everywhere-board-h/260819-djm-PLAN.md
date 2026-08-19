---
quick_id: 260819-djm
description: "q goes one level back everywhere (lane focus stops quitting), the board names `v`, the off-screen notice goes"
status: planned
---

# Quick Task: `q` is always one back, the board finally names `v`, and the off-screen notice goes

## Established facts (each one read in the code this session — do not re-investigate)

- `q` is **already** "back" in: `modeDetail` (`model.go:725`), `modeMove` (`629`),
  `modeProjects` (`703`), `modeHelp`/`modeMessage` (`655`), settings
  (`settings.go:52`), lanes (`lanes.go:233,240`), default board
  (`defaultboard.go:80`), browse (`browse.go:98,241`). All eight verified by
  grep. **No changes there.**
- Board (`model.go:797`) and compact view (`model.go:643`) quit the program on
  `q`. That is correct and stays: `runHome`'s for-loop
  (`internal/cli/home.go:22-44`) reopens the launcher whenever the board program
  exits, so `q` already lands on the project view for bare `jaira`; via
  `jaira board` it returns to the shell, which is also one back. **No change.**
- Typing modes (`modeFilter` 576, `modeCreate` 602, `modeEdit` 573) treat `q` as
  a literal character and cancel on `esc`. **Must stay untouched.**
- `esc` stays everywhere as an alias. "Replaces esc" means `q` also works, not
  that `esc` is removed.
- `m.key` (`model.go:568`) dispatches straight into `switch m.mode` — there is no
  global `ctrl+c` handler above it, so each mode's own branch is the only thing
  keeping `ctrl+c` alive.
- The compact view's own hint bar says `v full board` (`pipeline.go:265`), so
  the board's counterpart hint is `v compact`.
- The lane-focus footer is at `lanefocus.go:125` (not 84):
  `"enter open · esc/v back · q quit"`.
- `render_test.go:262-280` (`TestNarrowTerminalStillRenders`) currently
  **asserts** `"off-screen"` is present — it must be inverted, not left.
- The help screen lists `{"esc v", "back to the compact view"}` at
  `view.go:860`.
- Test helpers exist and are reusable: `newTestModel` (`render_test.go:75`),
  `stripANSI` (`render_test.go:341`), `laneFocusModel` (`lanefocus_test.go:10`).
- Baseline is green: `go build ./...` succeeds, `go test ./internal/tui/ -race`
  is `ok` (29s).

## The one real break

`model.go:781-785` intercepts `q` in `modeLaneFocus` and quits the whole program
before `laneFocusKey` ever sees it. Lane focus is one level deeper than the
compact view it opens from, so `q` there must go back to the compact view.

---

<task type="auto" tdd="true">
  <name>Task 1: `q` in lane focus goes back to the compact view</name>
  <files>internal/tui/model.go, internal/tui/lanefocus.go, internal/tui/view.go, internal/tui/lanefocus_test.go</files>
  <behavior>
    - `m.key(tea.KeyPressMsg{Code: 'q', Text: "q"})` in `modeLaneFocus` leaves
      `m.mode == modePipeline` and returns a nil cmd (it does not quit).
    - `m.key(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})` from `modeLaneFocus`
      still returns a non-nil cmd and does not land on `modePipeline`.
      (`tea.ModCtrl` exists in bubbletea v2.0.8 — `mod.go:12` — and `String()`
      renders it as `ctrl+c`, which is what the switch compares against.)
    - `stripANSI(m.renderLaneFocus())` contains `q/esc/v back` and does **not**
      contain `q quit`.
  </behavior>
  <action>
    In `model.go`, in `case modeLaneFocus:` (781-787), narrow the quit branch to
    `ctrl+c` only — `if s == "ctrl+c" { m.Close(); return m, tea.Quit }` — so `q`
    falls through to `m.laneFocusKey(s)`.

    In `lanefocus.go:34`, join `q` into the back case: `case "esc", "v", "q":`.
    Extend the existing comment there to say why `q` belongs with them: one level
    back from lane focus is the compact view, not the shell.

    In `lanefocus.go:125`, change the footer to
    `"enter open · q/esc/v back"`. Do not add a quit hint: `ctrl+c` is not a hint
    this project advertises in any other footer, and lane focus can no longer
    quit.

    In `view.go:860` (help, "Compact view" section), change the key label from
    `"esc v"` to `"q esc v"` so the help matches. Leave its description text
    alone. Note the label column is `%-9s` (`view.go:866`) and `"q esc v"` is 7
    chars, so the alignment still holds.

    In `lanefocus_test.go`, add one test next to
    `TestLaneFocusEscReturnsToPipeline` (~line 79) covering the behaviour above.
    Route it through `m.key(...)`, **not** `m.laneFocusKey("q")` — the bug lived
    in the dispatch, so a test that calls `laneFocusKey` directly would pass
    against the broken code. Use `laneFocusModel(t, 120, 34)`.
  </action>
  <verify>
    <automated>go build ./... &amp;&amp; go test ./internal/tui/ -race -run 'LaneFocus' 2>&amp;1 | tail -20</automated>
  </verify>
  <done>`q` in lane focus lands on the compact view, `ctrl+c` still quits, the footer and help say `q/esc/v back`, lane-focus tests green.</done>
</task>

<task type="auto" tdd="true">
  <name>Task 2: the board's hint bar names `v` and drops the off-screen notice</name>
  <files>internal/tui/view.go, internal/tui/render_test.go, internal/tui/hints_test.go</files>
  <behavior>
    - `stripANSI(m.render())` on a board contains `v compact`.
    - `stripANSI(m.render())` at 40×20 with seven lanes does **not** contain
      `off-screen` (inverted assertion in `TestNarrowTerminalStillRenders`), and
      still renders non-empty with no line wider than 60 runes — keep that half
      of the test as is.
    - The narrow-board hint test still finds every hint including `v compact`.
  </behavior>
  <action>
    In `statusBar` (`view.go:543-589`):
    - Delete the `hidden` variable and the `if end-start < len(m.cols)` block
      (564-567) entirely. `prefix` then starts as `""` — keep the warnings block
      that appends to it. `fmt` stays imported either way; it is still used at
      `view.go:575`.
    - Add `"v compact"` to the `keys` slice (572). Place it after `"enter open"`
      and before `"n new"` so the two view keys read as a pair with the compact
      view's own `v full board`. Do not reorder the rest.

    In `renderBoard`'s comment (`view.go:112-116`): the sentence "A board that
    silently truncates lanes would hide tickets, so the header always reports
    how many are off-screen" now describes removed behaviour. Rewrite those
    lines to state what is true: lanes scroll so the focused lane stays visible,
    and the count is deliberately not reported. Keep the second half of the
    comment (columns stretching to fill the row) verbatim.

    In `render_test.go:262-280`, invert the `off-screen` assertion to a
    must-not-contain and update the test's doc comment ("and must say so" is now
    false). The test name `TestNarrowTerminalStillRenders` is still accurate —
    leave it.

    In `hints_test.go`, add `"v compact"` to both key lists (line 15
    `TestWrapHintsKeepsEveryItem` and line 39
    `TestNarrowBoardShowsAllKeysAndFits`) so the wrap test covers the real bar.
    Re-check `TestNarrowBoardShowsAllKeysAndFits` still fits 24 lines at width 30
    with one more hint — if the extra wrapped line breaks the height assertion,
    that is a real regression to report, not a number to bump.
  </action>
  <verify>
    <automated>go build ./... &amp;&amp; go test ./internal/tui/ -race 2>&amp;1 | tail -20</automated>
  </verify>
  <done>Board bar shows `v compact`, no `off-screen` notice anywhere, the renderBoard comment matches the code, all tui tests green.</done>
</task>

## Verification (all must pass before done)

```
go build ./...
go test ./... -race 2>&1 | tail -20
gofmt -l internal/tui/
```

**Pre-existing gofmt state, measured before any change** (so a stale failure is
not mistaken for a new one):

- `internal/tui/view.go` — one stray blank line at EOF (`gofmt -d` shows only
  `-` on line 977). This file is edited by Task 2, so run `gofmt -w` on it and
  that line goes with it.
- `internal/cli/tickets.go` — a map-literal alignment diff, unrelated and
  untouched by this task. **Leave it.** `gofmt -l internal/` therefore still
  lists it; gate on `gofmt -l internal/tui/` being empty instead.

Then run the real binary, since that is how this gets judged:

```
go build -o /tmp/jaira-djm ./cmd/jaira && /tmp/jaira-djm board
```

- Board bar names `v compact`; no "N lane(s) off-screen" even at ~40 columns.
- `v` → compact view → `enter` → lane focus → `q` lands back on the compact
  view (not the shell). `esc` and `v` still do the same.
- `q` on the compact view and on the board still exits the program (and bare
  `jaira` returns to the launcher).

## Constraints

- Do not restructure the mode dispatch — only narrow the `modeLaneFocus` quit
  condition.
- Do not touch `modeFilter` / `modeCreate` / `modeEdit`.
- Do not change what `v` does anywhere.
- Do not touch `internal/cli/home.go` — its loop is already the "one back" for
  the board.
- No new hints beyond `v compact`; no compensating "ctrl+c quit" footer text.

## Commit

`feat(tui): q is one back from lane focus, the board names v, and the off-screen notice goes`
