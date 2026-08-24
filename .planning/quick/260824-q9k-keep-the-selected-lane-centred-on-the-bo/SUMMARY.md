---
quick_id: 260824-q9k
slug: keep-the-selected-lane-centred-on-the-board
date: 2026-08-24
status: complete
---

# The selected lane sits in the middle of the board

`.planning/NEXT-STEPS.md` item 7, closed.

## What changed

`renderBoard` scrolled the row so the focused lane landed on the **right edge**:

    start := 0
    if m.laneIdx >= perScreen { start = m.laneIdx - perScreen + 1 }

So the lanes *after* the focused one were never on screen — and what comes after
a lane is where its work goes next, which is the question the board is being read
to answer. Now:

    start = max(0, min(m.laneIdx-perScreen/2, len(m.cols)-perScreen))

Clamped at both ends, so a first or last lane still gets a full row rather than
blanks padding one side: centring by padding would trade away the same
information the centring was for.

The window moved into `laneWindow(perScreen)` beside `renderColumn` — one named
place, testable at every width and every lane, rather than three lines inside a
60-line render.

## Verified

- `go test ./... -race`, cache cleared: green, including the existing board
  renders at 40 and 170 columns. `gofmt -l core internal` unchanged.
- `TestTheWindowStaysFullAndHoldsTheFocusedLane` sweeps every lane at four row
  capacities and asserts two things a clamp gets wrong: the row is always full,
  and the focused lane is always inside it — scrolling past the lane you are on
  would be worse than the edge-anchoring it replaced.
- `TestCentringSurvivesARealRender` renders the whole board at 90 columns, where
  the lanes do not fit, and checks the lane after the focused one is drawn.
