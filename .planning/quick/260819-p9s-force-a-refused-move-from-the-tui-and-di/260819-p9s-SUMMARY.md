---
quick_id: 260819-p9s
status: complete
commits: 95a04de, 391f642
---

# Summary: `f` overrides a refused move, on the page the move was made from

## What was wrong

Moving a ticket from review to done was refused because its assignee (an email
address) did not read as the user, and the TUI's answer was "Edit the ticket to
supply what is missing, or use the CLI with --force" — leave the board, retype
the move. Underneath that, every overlay dismissed to the board no matter where
it was opened from, so even a refusal read over an open ticket cost the page.

## What changed

### 95a04de — a dialog goes back to the page it was opened from

- New `returnTo mode` field beside `detailFrom`. `openMove()` records it;
  `modeMove`'s esc/q, `modeMessage`'s dismiss keys and `applyMove`'s success path
  all use it. `modeHelp` still goes to the board (`?` is only bound there), so the
  shared `case modeHelp, modeMessage` was split.
- `notify()` records the page only for the modes that *are* a page — board,
  compact view, lane focus, open ticket. From the move picker it leaves `returnTo`
  alone: the picker is on its way out, and the page behind it is what the refusal
  belongs to. Settings, the lane editor and the default-board picker still dismiss
  to the board, exactly as before.
- `finishMove()` reloads an open ticket after a successful move, so it shows the
  lane it landed in instead of the one it left; a failed reload closes it rather
  than leaving a stale one up.
- `viewMode()` resolves whatever is on screen down to the page underneath, and
  `holdsLane()` (from 260819-p3m) is rewritten in terms of it — so a reload under
  a dialog or a message does not jump lanes either.

### 391f642 — the override itself

- `pendingMove` holds the refused move: ticket, target lane, the actor the gate
  was asked on behalf of, whether the move also claims, the refusals, and whether
  `f` has been pressed once.
- The refusal now ends "press f to override — the same override the CLI spells
  --force". `f` asks again ("Move X to Y anyway, overriding N refusal(s)?"), `y`
  writes. `n`, esc, q and enter all cancel and drop the offer with the message, so
  a stale `f` can never fire a move the user walked away from. `enter` is
  deliberately not a yes.
- `moveMutation()` is now shared by the gated and the forced path, so a forced
  move leaves exactly what a clean one would. Nothing is recorded on the ticket —
  the CLI's `--force` writes no audit field either; what was overridden is said
  out loud instead, in the CLI's own wording ("Overrode N gate refusal(s)").
- Force covers every refusal code, not only ownership, exactly as `--force` does
  (proven on a dependency-blocked ticket).
- `renderMessage()`'s footer names whichever keys are live: `esc dismiss`,
  `f override · esc dismiss`, or `y override · n cancel`.
- `switchBoard()` clears the pending move: it names a ticket in the store being
  swapped out.

`core/gate` and `internal/cli` were not touched.

## Verification

`internal/tui/returnto_test.go` (6 cases) and `internal/tui/force_test.go`
(8 cases), all driven through `m.key()` so the dispatch itself is covered. Notable
ones: nothing moves without the second key; the ready flag is recomputed on a
forced move; dismissing drops the offer; a blocked ticket is forceable; the
footer matches the state.

`go build ./...`, `go test ./... -race` green, `gofmt -l internal/tui/` empty.

## Worth knowing

Only one fixture ticket can move legally (the human-checkpoint one) — everything
else trips the promotion gate or its blocker. That is why the move-succeeds test
uses the human lane.
