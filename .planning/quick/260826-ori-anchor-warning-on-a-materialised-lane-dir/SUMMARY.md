---
quick_id: 260826-ori
slug: anchor-warning-on-a-materialised-lane-dir
date: 2026-08-26
status: complete
---

# A lane with no anchor is making a statement, not a broken reference

Ticket `4DQPMS`, the first thing the board in this repository found about the
tool that hosts it. Every command in this repository printed:

    jaira: warning: lane .jaira/lanes/backlog.md: anchor "" is not installed;
    placed before the terminal lane

## The ticket's hypothesis was wrong

It guessed the export writes an empty `after:`. It does not — `backlog.md` has
no `after:` line at all, which is correct, because the built-in `backlog` has no
anchor either. The defect was in the *loading*, not the writing.

## What was actually wrong

`core/lane/order()` was written when built-in lanes were always pre-placed: it
appends every `Builtin` lane first, then inserts the custom ones relative to
their anchors. `lanes remove` calls `MaterialiseWorkingSet`, which writes
*every* lane out as a file — and `replaceLane` deliberately does not mark a
replacement `Builtin`, because an override must not inherit the shipped slot of
the lane it displaced.

So after `lanes remove`, nothing is `Builtin`, and the placement list starts
empty. Then, for `backlog`:

- `terminalIndex` of an empty list returns `len(ls)`, i.e. `0`
- so the "no anchor: park before the terminal lane" branch computed `idx = -1`
- `idx < 0` is also how an *unresolvable* anchor is detected, and
  `present[""]` is false, so the deliberate case fell into the broken-reference
  branch and warned

The lane still worked. It was placed at the front, which is where it belongs.
Only the warning was wrong — and a validator that cries wolf is how `--strict`
gets turned off, which is the actual cost.

## What changed

`After == ""` is now its own case, ahead of the anchor lookup: it places the
lane at `terminalIndex(out)` — before the terminal lane, or at the front when
nothing is placed yet — and never warns. The two insertion sites collapsed into
one `insertAt` helper, so the clamp exists once.

`TestMaterialisedWorkingSetLoadsWithoutAnchorWarning` walks the real path:
`Load` → `MaterialiseWorkingSet` → `Load` again, and asserts no anchor warning
and that the unanchored lane stays at the front. It fails on the old code.

## Found while verifying, deliberately not fixed here: `QA3GN1`

The same test first asserted the board does not reorder, and caught this:

    before [backlog brainstorm todo pre-process in-progress human review signoff done blocked]
    after  [backlog todo pre-process in-progress human review signoff done blocked brainstorm]

`brainstorm` and `todo` are both `after: backlog`. `order()` puts a lane
*directly* behind its anchor, so the later-inserted sibling pushes the earlier
one right, and the chain then threads in front of it.

**That rule is deliberate and must not be "fixed".** It is exactly how
`critique` lands between `in-progress` and `human` — ahead of the built-in
successor that shares its anchor. Changing it would silently re-route his real
board.

The real gap is that `Materialise` (the `jaira init` path,
`core/lane/defaultboard.go:154`) writes the lane files and never calls
`SaveOrder`, unlike `Add`, `Remove` and `MoveLane`. `b.Lanes` is already an
ordered list, so the information is thrown away rather than missing. Filed as
`QA3GN1` with the measurement above; it is a different symptom on a different
command and deserves its own line in the release notes.

This board never showed it, because it was created by `lanes remove`, which
does write an order file.

## Verified

- `go test ./... -race` with the cache cleared: all packages ok
- `gofmt -l core internal` lists only `internal/cli/tickets.go`. The handoff
  expects `core/gate/gate.go` there too; it is no longer listed, and nothing in
  this change touched that file
- rebuilt `~/.local/bin/jaira` and ran `jaira list` in this repository: the
  warning is gone
