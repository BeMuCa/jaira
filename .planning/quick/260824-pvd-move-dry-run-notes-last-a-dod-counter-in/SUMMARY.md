---
quick_id: 260824-pvd
slug: move-dry-run-notes-last-a-dod-counter-in-list-a-stale-claim-warning
date: 2026-08-24
status: complete
---

# The small package

`.planning/NEXT-STEPS.md` item 4, closed. Four frictions, four commits.

## 1. `move --dry-run` (`a942ad0`)

Asking a gate a question used to mean trying the move. An agent created and
deleted a throwaway ticket to find out; one of the user's tickets was moved by
accident probing a gate that turned out not to refuse.

**The trap, found in the code:** the real move stages `--what`/`--why`/
`--resolves` and the assignee *before* the gate runs, deliberately, so a gate can
require them. A dry run reusing that path would leave those fields on a ticket
whose move it then refused. It stages into the in-memory document and re-decodes,
so the gate sees exactly the ticket it would have seen and nothing reaches the
disk — asserted by comparing the file byte for byte afterwards.

The refusal is the real one down to the exit code (both paths now share
`refusalCode`), so a caller branches identically.

## 2. `show --notes-last N` (`8c195a5`)

**Decision made here, since the item left it open: the default does not change.**
`show` still prints every note. This board's promise is that nothing is lost, so
a note is hidden only when the reader asks — and the hidden count is always
printed, so a trimmed read never looks like a complete one.

Only the Progress section is trimmed; a continuation line travels with its note.
On his 263-line ticket: 144 lines at `--notes-last 3`, reporting 23 hidden.

## 3. `[DoD n/m]` in `jaira list` (`a49c6fe`)

The numbers were already in the JSON and on the board; only the human listing
was missing them. Counts settled items, so a superseded criterion does not leave
a completable ticket reading 1/2 for good.

## 4. A stale-claim warning (`0043cb8`)

`ClaimActive` treats a claim older than `ClaimTTL` as abandoned and says nothing.
New `StaleClaim` answers the other question — whose was it, how long ago — and
`claim`, `next` and `next --per-lane` report it. Stderr, so `--json` on stdout
stays parseable; `took_over` in the claim's own JSON. A claim with no timestamp
is reported as such rather than as "renewed moments ago".

## Verified

- `go test ./... -race`, cache cleared: green. `gofmt -l core internal` lists
  exactly the two documented files.
- On the running binary: a dry run of a refused move and of an allowed one, with
  the ticket unchanged in both; `--notes-last 3` on the real 263-line ticket;
  `[DoD n/m]` across his backlog.
- `docs/COMMANDS.md` carries all four; `SKILL.md` carries `--dry-run` where it
  tells an agent what to do about a refusal, which is where an agent will look.

## Not done, on purpose

The lane warning already goes to stderr (`internal/cli/root.go`). The other
session's feedback was wrong about it, as recorded in HANDOFF.md.
