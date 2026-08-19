---
quick_id: 260819-pk3
status: complete
commit: ae6cee3
---

# Summary: write a follow-up next to the ticket it is for

## What changed

`n` on an open ticket splits the screen: the ticket stays on the left, the
follow-up is written on the right. New file `internal/tui/followup.go` holds the
whole subject.

- `follow{src, srcScroll, draft, focusLeft}` on the model. `m.detail` stays the
  **right** pane, so every existing detail key keeps acting on the ticket being
  worked. No new mode: drafting is `modeEdit`, a saved follow-up is `modeDetail`,
  and `View()` splits when `m.follow != nil`.
- **A draft is not a file.** The field editor writes each field straight to disk
  (`commitEdit`), which is right for a ticket that exists and wrong for one the
  user may still abandon — so `draft{fields, lists, body}` is edited in memory and
  `store.Create` runs once, on ctrl+s. `esc` discards it and the board never saw
  it (asserted against the store's own count).
- `saveFollowUp()` creates, recomputes `ready` through the gate the way the
  on-disk editor does after every field, reloads, and hands the right pane over as
  a real ticket. `closeFollowUp()` serves both esc paths — discarding an unwritten
  draft and closing a written one are the same move: the ticket it was for comes
  back alone.
- `n` again chains: the ticket just written becomes `src`, the older one slides
  off. Two panes only, as asked.
- `followUpFields()` is now shared with the `f`-at-sign-off path so the two cannot
  drift, and `followUpContext()` takes its lead-in sentence — only sign-off may
  claim a review happened; the split writes "Follows on from &lt;handle&gt;."
- Focus decides **scrolling**, not what the action keys write to. The editor keeps
  `tab` for its fields, so the left pane scrolls with shift+arrows in both states;
  once saved, `tab` moves between panes and the arrows follow. The live pane's
  border carries the accent colour.
- `j`/`k` (walk to the neighbouring ticket) close the split with the ticket.
- Below 80 columns or 20 rows there is no split: the follow-up takes the screen
  through the plain renderers.

Supporting refactors:

- `renderDetail` split into `detailBody(t, width)` + `detailHints(t)`, and
  `renderEdit` into `editBody(w, height)` + its footer, so both render into a pane.
- `clipPane(content, width, height, scroll)` hands a box an exact rectangle —
  cut, then padded, per row. This matters: a sized lipgloss style measures its
  frame *including* the border, so `Width(inner)` on content already `inner` wide
  re-wrapped every full-width row and grew the pane past the terminal (the panes
  came out 39 and 40 rows in a 40-row terminal before this). `splitPane` therefore
  sets neither Width nor Height and lets the border hug the rectangle — the same
  trap `clampBlock`'s comment already documents.
- Help gained the split's keys and the force key from 260819-p9s.

## Verification

`internal/tui/followup_test.go`, 13 cases through `m.key()`: opens without
writing, both tickets on screen, typing touches nothing on disk, discard leaves no
trace, save lands in the default lane with `follows` set and the typed goal, a
saved follow-up renders as a ticket, chaining slides left, esc comes back then
out, tab moves focus and the arrows follow it, shift+arrows scroll while writing,
and both the narrow and the short terminal drop the split.

`go build ./...`, `go test ./... -race` green, `gofmt -l internal/tui/` empty,
binary built and the three states inspected as rendered text.

## Observed, not changed

The seeded body carries a placeholder definition of done
(`- [ ] <What must be true that is not true yet>`), and the gate counts it — so a
saved follow-up is marked `ready` even though nobody has said what done means.
This is pre-existing: the `f`-at-sign-off path has created follow-ups this way all
along. Out of scope here; worth a decision.
