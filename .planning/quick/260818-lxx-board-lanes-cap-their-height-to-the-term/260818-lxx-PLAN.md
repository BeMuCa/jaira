---
quick_id: 260818-lxx
description: board lanes and the single-lane view cap their height to the terminal window
status: planned
---

# Quick Task: a lane never renders taller than the window

Berk's report: the Review lane holds 30+ tickets and the lane renders taller
than the terminal. Lanes must adapt to the window height, with the window as
the maximum.

## Root cause (established — do NOT re-investigate)

### The board (`internal/tui/view.go`)

`renderBoard` computes `bodyHeight` (view.go:135) and `renderColumn` caps
visible cards at `visible = max(1, (h-4)/3)` with a `+N more` line — the scroll
math assumes **every card is exactly 3 lines**.

But `renderCard` truncates the flags line at `w+24`:

```go
out += "  " + truncate(strings.Join(flags, " "), w+24) + "\n"   // view.go:337
```

`truncate` is display-width based (view.go:907, `lipgloss.Width`), so that line
can be **22 cells wider than the column's content width**. Verified against
charm.land/lipgloss/v2 v2.0.6 (the version in go.mod): `Style.Width(w)` **wraps**
overlong lines onto extra rows, and `Style.Height(h)` pads but does **not clip**
taller content. So a flag-heavy card renders 4-5 lines instead of 3, the column
exceeds `bodyHeight`, and the whole board grows past the terminal.

Review-lane cards accumulate the most flags (`◆ sign off`, `Plan n/m`,
`DoD n/m`, `✓ commits`, `executedBy`, `@assignee`), which is why the Review lane
is where it shows.

`w+24` dates from the first board commit (9e229f9) — a stale ANSI-length fudge
from before `truncate` became width-based. It is **not** a recorded design
decision (checked against the project's design-invariant notes).

**Measured during planning** (probe test, deleted afterwards; fixture = 30
review tickets with `Plan 1/3` + `DoD 1/2` + `✓ 3` + `claude-opus-5` +
`@someoneelse`, terminal 120x30):

| state | board render |
|---|---|
| today | **42 lines** at `m.height = 30` |
| with `w+24` → `w` | 28 lines at h=30; also fits at h=40/24/20/16/12/10 and w=60/80/120/200 |

The whole existing `internal/tui` suite stayed green with that one-token change.

### The single-lane view (`internal/tui/lanefocus.go`)

`renderLaneFocus` renders **all** tickets in a plain `for i, t := range
col.tickets` loop (lanefocus.go:78) — no height cap, no scroll. Measured:
**155 lines at `m.height = 30`**. Its cards are variable height (the flags line
and the goal line are conditional), unlike the board's fixed 3-line cards, so
it cannot borrow the board's `/3` arithmetic.

## Task 1 — the board fits, and cannot un-fit again

Files: `internal/tui/view.go`, `internal/tui/scroll_test.go`

### (a) Root-cause fix

`renderCard` line 337: `w+24` → `w`. The flags line then fits the card like
every other line, exactly as `renderLaneCard` already does
(`truncate(strings.Join(flags, "  "), w)`, lanefocus.go:146). The board's
3-lines-per-card scroll math stays untouched — it becomes true again instead of
being worked around.

Also give the `meta` line the same treatment for the same reason: it is written
as `"  " + meta` with no truncation at all (view.go:335), so a long handle plus
assignee can also exceed `w`. Wrap it: `"  " + truncate(meta, w)`.

### (b) Defensive clamp — a column can never spill the board again

New helper in `view.go`, next to `truncate`:

```go
// clampBlock cuts a rendered block to a hard w x h budget. lipgloss's Width()
// wraps overlong lines onto extra rows and Height() pads without clipping, so
// a block handed to a sized style can silently grow the layout around it —
// which is exactly how a flag-heavy review lane once pushed the board past the
// bottom of the terminal. Clamping here means over-tall content is cut, never
// spilled.
func clampBlock(s string, w, h int) string
```

Behaviour: split on `"\n"`, `truncate(line, w)` each line, keep at most `h`
lines, re-join. Guard `w <= 0` / `h <= 0` (return `""`).

Apply it in `renderColumn` immediately before styling:

```go
return style.Render(clampBlock(body.String(), w, h))
```

(`w`/`h` are the style's content box, so the clamp is the same budget lipgloss
is being told to honour.)

Do **not** change `bodyHeight` or `visible` unless a test in (c) proves the
arithmetic is off by one — the measurements above say it is not.

### (c) Test — `internal/tui/scroll_test.go`

`reviewLane(t, n)` already lives there but its tickets are flag-light (only
`◆ sign off` + `✓ 1`) and **do not reproduce the bug**. Add a sibling helper
`flagHeavyReviewLane(t, n)` that maximises the flags row, since that row is the
whole failure mode:

- `status: review` → `◆ sign off`
- `assignee: someoneelse` (not `m.me`) → the yellow `@name` branch
- `executed-by: claude-opus-5` → the agentic flag
- 3 commits → `✓ 3`
- body with `## Plan` (3 items, 1 checked) and `## Definition of Done`
  (2 items, 1 checked) → `Plan 1/3` + `DoD 1/2`; pass the body as the third arg
  to `s.Create(f, l, body)`, the store parses the checklists
- long titles (`"Review item long title here"`)

`TestBoardFitsTheTerminal`: build `flagHeavyReviewLane(t, 30)`, `New`,
`reload`, focus the review column, then for each `h` in `{40, 30, 24, 20, 16,
12, 10}` and each `w` in `{60, 80, 120, 200}` assert
`len(strings.Split(stripANSI(m.render()), "\n")) <= h`. Follow
`TestDetailFitsTheTerminal` (scroll_test.go:127) for shape and message style.

`TestSelectionStaysVisibleWhenScrollingDown` (scroll_test.go:45) must keep
passing unchanged — the clamp must not eat the selected card. If it does, the
clamp budget is wrong, not the test.

## Task 2 — the single-lane view fits, with cursor-following scroll

Files: `internal/tui/lanefocus.go`, `internal/tui/lanefocus_test.go`

Cap `renderLaneFocus` to `m.height` with a scrolled window that keeps the
cursor visible, consistent with the board's approach but accounting **lines per
card** rather than assuming 3.

Shape:

1. Render each visible card with `m.renderLaneCard(...)` as today. Card height
   is variable, so measure it: `strings.Count(card, "\n")` (+1 for the blank
   separator line the loop writes after each card).
2. Budget: derive the header cost from the builder the way the existing code
   already does (`strings.Count(b.String(), "\n")`, lanefocus.go:85) rather
   than hardcoding it, and reserve the same 2 lines the current footer
   arithmetic reserves (`m.height - used - 2`). Reserve one more line for each
   off-screen indicator actually shown.
3. Scroll state: **reuse `m.scroll[col.lane.ID]`** (model.go:69) — no new
   field. The board recomputes and clamps `first` for the focused lane on every
   render (view.go:257-266) and lane focus only ever scrolls the lane that is
   focused, so the board self-corrects whatever this view leaves behind, and the
   two views agreeing on where a long lane is scrolled to is the desirable
   behaviour anyway.
4. Keep the cursor visible: if `m.cardIdx < first` → `first = m.cardIdx`; while
   the cards from `first` through `m.cardIdx` do not fit the budget →
   `first++`. Write `first` back to `m.scroll[col.lane.ID]`. Clamp a stale
   `first` past the end (the board does this too).
5. Off-screen indicators, so nothing is hidden silently — the board already
   holds this line: `+N more` above when `first > 0`, `+N more` below when the
   last rendered index is not the last ticket. Match the board's wording
   (`styMeta.Render(fmt.Sprintf(" +%d more", rest))`, view.go:275).
6. Last line of defence: return `clampBlock(result, m.width, m.height)` from
   Task 1(b), so this view can no more spill than a column can.

Keep the existing empty-lane path (`no tickets`) and the gap-padding that keeps
the footer pinned to the bottom.

### Tests — `internal/tui/lanefocus_test.go`

Uses `newTestModel` + `stripANSI`; the flag-heavy fixture from Task 1 lives in
the same package, so reuse it (build the model from
`flagHeavyReviewLane(t, 35)` rather than the small `newTestStore` fixture).

1. `TestLaneFocusFitsTheTerminal` — 35 tickets, `mode = modeLaneFocus`, review
   lane focused; for each `h` in `{40, 30, 24, 20, 16, 12}` assert
   `len(strings.Split(stripANSI(m.renderLaneFocus()), "\n")) <= h`.
2. `TestLaneFocusKeepsTheCursorVisible` — at 120x20, walk `m.cardIdx` from 0 to
   `len(tickets)-1` and assert `ticket.Handle(col.tickets[i].ID)` appears in
   every render (the card renders the handle, lanefocus.go:104). Model this on
   `TestSelectionStaysVisibleWhenScrollingDown`.
3. `TestLaneFocusSaysWhatIsOffScreen` — at 120x20 with the cursor mid-list,
   assert `more` appears in the render (nothing hidden without saying so).
4. Existing lane-focus tests must pass unchanged, especially
   `TestLaneFocusShowsOnlyItsOwnLane` (both backlog tickets fit at 120x34) and
   `TestLaneFocusEmptyLaneDoesNotPanic`.

## Verification (all must pass before done)

- `go build ./...`
- `go vet ./...` clean
- `go test ./... -race 2>&1 | tail -20` green — **always `-race`**, CI runs it
  and a plain `go test` once let a failure through
- Rebuild Berk's real binary: `go build -o ~/.local/bin/jaira ./cmd/jaira`
  (his `jaira` runs from `~/.local/bin`; `go install` does NOT update it).
  He verifies against the running binary with a real 30+ ticket review lane.

## Constraints

- Only `internal/tui/view.go`, `internal/tui/lanefocus.go`,
  `internal/tui/scroll_test.go`, `internal/tui/lanefocus_test.go` change.
- Do **not** touch `core/merge`, lane precedence, or any gate logic.
- No new dependency — `truncate`, `lipgloss.Width`, `min`/`max` are already in
  `view.go`.
- Match the file's comment voice: full sentences explaining *why*, the way the
  existing renderer comments do.
- Do not commit `.planning/` files — the orchestrator does that afterwards.

## Commit

One atomic commit for code + tests, on master. Repo style is lowercase
narrative. Suggested:

`fix: a lane stops at the bottom of the window — the flags row fits the card again, and the single-lane view scrolls`

## Threat model

No trust boundary is crossed and no package is installed: the change is pure
terminal rendering of data already loaded from the local ticket store. The one
integrity concern is the opposite of a security one — a clamp that silently
cuts content — and it is mitigated by the `+N more` indicators plus the
cursor-visibility tests, so nothing disappears without the board saying so.
