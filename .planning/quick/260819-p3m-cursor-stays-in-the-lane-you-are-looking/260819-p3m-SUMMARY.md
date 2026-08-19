---
quick_id: 260819-p3m
status: complete
commit: 06a1433
---

# Summary: the lane you are looking at stays the lane you are looking at

## What was wrong

`rebuild()` remembered the selected ticket by ID before regrouping and called
`selectByID()` afterwards, and `selectByID` sets **both** `m.laneIdx` and
`m.cardIdx`. In the compact view and in lane focus `m.laneIdx` is not a cursor —
it is the page. So moving the selected ticket from another shell (`jaira move
<id> done`) swapped the whole screen over to the target lane. Reproduced in a
test first: `laneIdx` went 0 → 8 on the fixture.

## What changed

`internal/tui/model.go` only:

- `holdsLane()` — is the view showing one lane rather than all of them?
  `modePipeline`, `modeLaneFocus`, and `modeDetail` resolved through
  `m.detailFrom`, because an open ticket belongs to the view it was opened from.
- `rebuild()` also remembers the held lane, **by ID** — lanes come and go between
  reloads, so an index would not survive.
- `holdLane(laneID, selectedID)` puts the view back on that lane. Inside it the
  cursor still follows the ticket by ID (cards are sorted by `updated_at`, so any
  unrelated edit reorders them); if the ticket left the lane the cursor keeps its
  index and lands on whatever slid into the gap. A lane that disappeared falls
  back to following the ticket — pointing at a lane that no longer exists is
  worse.

`modeBoard` is untouched: there the new lane is already on screen and tracking
the ticket is the answer to "where did it go".

## Verification

`internal/tui/holdlane_test.go`, five cases: lane focus holds, the compact view
holds, an open ticket over lane focus holds, the board still follows the ticket,
and a held lane still follows its ticket when the lane reorders. The last two
passed before the change (they are the regression fence); the first three failed
before it and pass after.

`go build ./...`, `go test ./... -race` green, `gofmt -l internal/tui/` empty.

## Not done

Overlay modes (move picker, message) were left out of `holdsLane` here — they had
no way to record the page behind them yet. Quick task 260819-p9s added
`returnTo` and folded them in.
