---
quick_id: 260824-q7x
slug: create-mine-assigns-you-at-creation
date: 2026-08-24
status: complete
---

# `create --mine`, and the invariant it must not reverse

`.planning/NEXT-STEPS.md` item 6, closed.

## What changed

`jaira create <title> --mine` sets the assignee to `identity()`. Three lines in
`newCreateCmd`, placed after the board template's own assignee and before the
fields map, so the precedence reads top to bottom: `--assignee`, then a template
that names someone, then `--mine`, then nobody.

## The invariant it is built around

Capture belongs to nobody. Capturing and claiming are two acts, and the claim is
the pull into work — moving an unassigned ticket assigns the mover. `--mine` is
an opt-in and stays one: it says "I am taking this one", which is a claim, and a
claim is exactly what capture is not.

`TestPlainCreateStillLeavesTheAssigneeEmpty` asserts a plain `create` is
unchanged, because a convenience like this is how an invariant quietly becomes a
default. (`TestCreateLeavesTheTicketUnassigned` in `claiming_test.go` already
owns the invariant itself; this one asserts it survives beside the new flag.)

Nothing was added to the TUI: it already claims on the pull into `todo`, so the
normal flow was never missing anything. A key at creation time would have to not
be the default, and there is nothing yet asking for one.

## Verified

- `go test ./... -race`, cache cleared: green. `gofmt -l core internal`
  unchanged.
- On the running binary: `create --mine --json` → `"assignee": "BeMuCa"`, plain
  `create --json` → `""`.
