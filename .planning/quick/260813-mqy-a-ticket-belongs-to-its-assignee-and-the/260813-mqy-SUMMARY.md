# Quick Task 260813-mqy: A ticket belongs to its assignee — Summary

One-liner: the gate now refuses a write to someone else's ticket outside the human checkpoint lanes, naming the owner and the two ways forward, with `--force` and hand-over both working exactly as designed.

## What changed

1. **`core/identity`** (0180214) — moved `identity()`'s body out of
   `internal/cli/root.go` into a new `core/identity` package (`Current(dir)`,
   `Slug(name)`), so both the CLI and the gate can share it. `internal/cli`'s
   `identity()` is now a one-line call-through; behavior is unchanged.
2. **Ownership gate** (b8da1d0) — `core/gate.CheckAdvance` refuses a move when
   the ticket has a non-empty `assignee`, the actor isn't that assignee
   (case-insensitive, trimmed), the ticket isn't currently sitting in a
   `RequiresQuestion` or `RequiresHumanExit` lane, and the request isn't a
   hand-over (`Request.NewAssignee`). New violation code `not_owner`. The
   refusal message: `"<id> belongs to <owner> — ask them, take it over with
   'jaira set <id> assignee=<you>', or override with --force"`. Flows through
   the CLI's existing `--force` mechanism with no new flag; exit code stays 3.
3. **`docs/AGENTS.md`** (df75376) — documented the refusal, its exemptions, and
   that it's a guard rail, not a lock (the merge driver is untouched).

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1/2 - Bug / missing functionality] TUI's plain board move never threaded the actor into the gate**
- **Found during:** Task 2, after adding the ownership check to `CheckAdvance`.
- **Issue:** `internal/tui/model.go`'s `applyMove()` (drag-move in the board)
  called `gate.CheckAdvance` with an empty `Request.Actor`. Once ownership
  enforcement existed, this would refuse *every* board move of an assigned
  ticket — including the owner's own — because an empty actor never matches a
  non-empty assignee. `internal/tui/signoff.go`'s `accept()` already passed
  `Actor: identity(m.store.Root)`; `applyMove()` was the one call site that
  didn't.
- **Fix:** Added `Actor: identity(m.store.Root)` to the `gate.Request` built in
  `applyMove()`, matching the existing pattern used elsewhere in the same file.
- **Files modified:** `internal/tui/model.go` (not in the plan's file list,
  but directly caused by the task 2 gate change).
- **Commit:** b8da1d0
- **Test added:** `internal/tui/move_test.go` —
  `TestApplyMoveAllowsOwnerToMoveOwnTicket` and `TestApplyMoveRefusesNonOwner`,
  proving the board's drag-move both allows the owner and refuses a non-owner.

No other deviations. `internal/cli/flow.go` needed no changes: it already
threaded `Actor: identity()` into `gate.Request` for `jaira move`, and the new
`not_owner` code needs no special exit-code handling (it defaults to
`ExitValidation` = 3, same as every other non-blocked refusal).

## Test coverage added

`core/gate/gate_test.go` — six new tests plus a `findViolation` helper:
- `TestOwnershipRefusesWriteByOtherThanAssignee` — refused, message names owner
- `TestOwnershipAllowsLeavingRequiresHumanExitLane` (signoff)
- `TestOwnershipAllowsLeavingRequiresQuestionLane` (human)
- `TestOwnershipAllowsHandOver` (`NewAssignee` set)
- `TestOwnershipAllowsWhenNoAssignee`
- `TestOwnershipAllowsSameActorDifferentCase`

`core/identity/identity_test.go` — `Current` honors `JAIRA_USER`, falls
through to `USER`/etc. when git config is unreachable (HOME redirected to an
empty dir so the real developer's `~/.gitconfig` can't leak into the test);
`Slug` on "BeMuCa", on "Anna Müller" (no leading/trailing separator, only
`[a-z0-9-]`), and on punctuation-only input (non-empty fallback).

`internal/tui/move_test.go` — see deviation above.

## Manual verification (real binary, throwaway repo)

Built the binary, ran `jaira init` in a fresh git repo, created a ticket
assigned to `anna`:

```
== move to todo as anna (owner) ==
AGKXX3 → todo
exit=0

== move to in-progress as bob (not owner) ==
jaira: cannot move AGKXX3 to in-progress:
  - AGKXX3 belongs to anna — ask them, take it over with 'jaira set AGKXX3 assignee=<you>', or override with --force
exit=3

== move to in-progress as bob with --force ==
AGKXX3 → in-progress
Overrode 1 gate refusal(s):
  - AGKXX3 belongs to anna — ask them, take it over with 'jaira set AGKXX3 assignee=<you>', or override with --force
exit=0

== bob takes over via set (no --force needed) ==
Updated AGKXX3
exit=0

== charlie (not owner) advances out of the human checkpoint lane ==
AGKXX3 → review
exit=0
```

Confirmed: refusal names the owner and both ways forward; `--force` overrides
and is recorded; hand-over via `jaira set assignee=` needs no `--force`; a
non-owner can still advance a ticket out of the `human` checkpoint lane.

## Final verification

```
$ export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go vet ./... && go test ./... -count=1
?   	github.com/BeMuCa/jaira/cmd/jaira	[no test files]
ok  	github.com/BeMuCa/jaira/core/board	0.063s
ok  	github.com/BeMuCa/jaira/core/gate	0.694s
?   	github.com/BeMuCa/jaira/core/gitrepo	[no test files]
ok  	github.com/BeMuCa/jaira/core/identity	0.083s
ok  	github.com/BeMuCa/jaira/core/lane	0.596s
ok  	github.com/BeMuCa/jaira/core/merge	0.455s
ok  	github.com/BeMuCa/jaira/core/project	0.132s
ok  	github.com/BeMuCa/jaira/core/release	0.014s
?   	github.com/BeMuCa/jaira/core/session	[no test files]
ok  	github.com/BeMuCa/jaira/core/ticket	0.132s
ok  	github.com/BeMuCa/jaira/core/validate	0.235s
ok  	github.com/BeMuCa/jaira/internal/cli	0.736s
ok  	github.com/BeMuCa/jaira/internal/tui	3.445s
?   	github.com/BeMuCa/jaira/scripts/iconpreview	[no test files]
?   	github.com/BeMuCa/jaira/scripts/shotgen	[no test files]
```

## Commits

- `0180214` feat: move identity into core so the gate can use it too
- `b8da1d0` feat: a ticket belongs to its assignee
- `df75376` docs: say that a ticket has an owner

## Merge driver / rules

Untouched, as required. Ownership is enforced only inside `core/gate` at the
moment of a `jaira move` (CLI or TUI); it never touches `core/merge` or how
git resolves concurrent writes to the same file.

## Self-Check

- `core/identity/identity.go` — FOUND
- `core/identity/identity_test.go` — FOUND
- `core/gate/gate.go` (ownership check) — FOUND
- `core/gate/gate_test.go` (ownership tests) — FOUND
- `internal/tui/model.go` (Actor threaded into `applyMove`) — FOUND
- `internal/tui/move_test.go` — FOUND
- `docs/AGENTS.md` (ownership paragraph) — FOUND
- Commit `0180214` — FOUND
- Commit `b8da1d0` — FOUND
- Commit `df75376` — FOUND

## Self-Check: PASSED
