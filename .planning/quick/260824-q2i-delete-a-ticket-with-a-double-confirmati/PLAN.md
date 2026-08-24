---
quick_id: 260824-q2i
slug: delete-a-ticket-with-a-double-confirmation
date: 2026-08-24
mode: quick
---

# Quick task — delete a ticket, with a double check

`.planning/NEXT-STEPS.md` item 5.

## Today

`Store` has `Archive` and `Restore` and no delete at all. Archiving is right for
finished work; a ticket created by mistake or a throwaway probe leaves a file
that only `rm` removes.

## Decisions (already made)

- `jaira delete <id>` removes the file. The confirmation is the ticket's handle
  typed back, not a y/n — deletion is the one irreversible operation in a tool
  whose whole promise is that nothing is lost.
- `--force` skips the prompt for scripts, and says what it deleted.
- In the TUI it belongs on the ticket, not the board, and needs the same second
  step. `x` is archive and must not change meaning.
- The help says archive is almost always what you want.

## What the code says

- `Store.Archive` (`core/ticket/store.go:120`) resolves symlinks first so it
  moves the file rather than the link. `Delete` mirrors that, and removes the
  link too when they differ — otherwise a symlinked ticket is half-deleted.
- `x` is bound in `modeBoard` only (`internal/tui/model.go:947`); `modeDetail`
  leaves it free but the decision is that `x` keeps meaning archive everywhere.
  `X` on the open ticket, then a typed handle.
- `m.input` is already the shared line-editor buffer for filter and create, so
  the typed confirmation reuses it rather than adding a second one.
- **The trap:** `core/validate/validate.go:116` makes a dangling `blocked-by` an
  **error**. Deleting a ticket another one depends on would leave the board
  failing `jaira validate` — so the delete is refused when anything still points
  at it, naming what does. `--force` still goes through; that is the user's call,
  and validate will then say so.

## Tasks

1. `Store.Delete(id) (string, error)` — resolve, remove, return the path.
2. `jaira delete <id>`: the handle typed back on stdin, prompt on stderr so
   `--json` on stdout stays parseable; `--force` for scripts; refusal when other
   tickets reference it via `blocked-by` or `follows`.
3. TUI: `X` on the open ticket opens the typed confirmation, esc cancels.
4. Docs: `docs/COMMANDS.md`, the README command list, the TUI help screen, and
   `archive.go`'s own help, which currently says "nothing here removes anything"
   — my change makes that false.

## Done when

A mistyped confirmation deletes nothing, a correct one removes the file,
`jaira validate` is clean afterwards, and no other command starts treating a
missing ticket as an error.
