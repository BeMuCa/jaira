---
quick_id: 260824-q2i
slug: delete-a-ticket-with-a-double-confirmation
date: 2026-08-24
status: complete
---

# A ticket can be deleted, and it costs two deliberate acts

`.planning/NEXT-STEPS.md` item 5, closed.

## What changed

**`Store.Delete`** — the only irreversible operation in the store. Symlinks are
resolved the way `Archive` resolves them, and both the link and the file it
points at go: a half-deleted ticket is worse than either outcome.

**`jaira delete <id>`** asks for the handle typed back, not a yes. The prompt
goes to stderr so `--json` on stdout stays parseable, and the answer is read from
stdin — a caller with nothing to say gets EOF, which is not the handle, which
deletes nothing. `--force` skips it for scripts and still says what it removed.

**`X` on an open ticket in the TUI**, then the same typed handle. Not `x`: that
is archive on the board and keeps meaning archive everywhere. The ticket stays on
screen while its handle is typed, so what is about to be destroyed is in front of
the person destroying it. `esc` cancels.

**The trap, found while reading `core/validate`:** a dangling `blocked-by` is an
**error** there, not a warning. Deleting a ticket another one depends on would
leave the board failing `jaira validate` — so the delete is refused while
anything still points at it, naming what does, in both `blocked-by` and
`follows`. `--force` goes through; that is the user's call.

**Docs**: `archive.go`'s own help said "nothing here removes anything", which
this made false — corrected, and it now names the one command that does. Plus
`docs/COMMANDS.md`, the README list, the TUI help screen, and one line in
`SKILL.md` telling agents deleting is the user's, not theirs.

## Verified

- `go test ./... -race`, cache cleared: green. `gofmt -l core internal` lists
  exactly the two documented files.
- Nine new tests. The CLI ones cover five wrong answers (empty, `y`, `yes`, a
  prefix of the handle, EOF) leaving every file in place, a correct one removing
  it with `validate` clean afterwards, `--force`, the referenced-ticket refusal,
  and `list`/`next`/`show` still working on a board a ticket was deleted from.
  The TUI ones cover the typed confirmation, the cancel, and `x` still archiving.
- On the running binary: `yes` refused with exit 3 and the file intact, the
  handle accepted with exit 0, `jaira validate` clean.

## Deliberately not done

No `delete` in the board's own key list beyond the open ticket, and no bulk
delete. One ticket, one handle, one deliberate act.
