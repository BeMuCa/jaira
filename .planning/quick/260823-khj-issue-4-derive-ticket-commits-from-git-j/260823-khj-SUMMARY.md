---
phase: quick-260823-khj
plan: 01
subsystem: cli
tags: [git, go, kanban, cli, tui]

requires: []
provides:
  - "core/gitrepo.Repo.CommitsForTicket — derives the commit list for a ticket from git history"
  - "gate.Env.DeriveCommits — injected closure the gate calls to satisfy requires-commits without a recorded sha"
  - "jaira sync <id> — stamps derived commits and moves a terminal-lane ticket into .jaira/sync/<initials>-<date>/"
  - "Store.Sync / Store.Restore extended to cover .jaira/sync/"
  - "identity.Initials — folder-naming helper reused by jaira sync"
  - "core/board jaira:local marker — preserves project-specific rules across agent-note regeneration"
affects: [cli, tui, board, gate, gitrepo, ticket]

actuals:
  tokens: 16536
  tasks: 3
  commits: 3

tech-stack:
  added: []
  patterns:
    - "Injected derivation closures (gate.Env.DeriveCommits) keep a pure package (core/gate) filesystem- and git-free while still reaching git through the caller"
    - "Fixture git repos built at test runtime with t.TempDir() + os/exec, never a committed nested .git"

key-files:
  created:
    - core/gitrepo/derive.go
    - core/gitrepo/derive_test.go
    - internal/cli/syncout.go
    - internal/cli/syncout_test.go
    - core/ticket/sync_test.go
    - core/board/announce_test.go
  modified:
    - core/gate/gate.go
    - core/gate/commits_test.go
    - core/identity/identity.go
    - core/identity/identity_test.go
    - core/ticket/store.go
    - core/board/announce.go
    - internal/cli/root.go
    - internal/cli/archive.go
    - internal/tui/model.go
    - internal/tui/signoff.go
    - internal/tui/view.go
    - internal/tui/lanefocus.go
    - docs/COMMANDS.md

key-decisions:
  - "Commit list is the union of 'git log --follow' on the ticket file and a --grep over commit messages matching either the full ULID or the handle, deduped and reordered oldest-first via 'git rev-list --no-walk=sorted' reversed in Go (D-01)"
  - "The handle (not just the ULID) is grepped, bounded by a non-alphanumeric character or the line edge, case-insensitive, since the handle is the only id form a person or agent ever writes in a commit message"
  - "Explicit always wins: DeriveCommits is invoked only when the ticket records no commits of its own, and never overrides a recorded list or an explicit --commits"
  - "The gate records nothing on the ticket — stamping happens once, at the moment a ticket actually leaves the board via 'jaira sync' or 'jaira archive' (best-effort for archive)"
  - "Sync folders are named <initials>-<yyyymmdd>, via a new identity.Initials helper, reusing identity.Slug for the single-word/empty-name fallback"
  - "core/board's managed block gained a <!-- jaira:local --> marker: everything from it to the end marker survives regeneration verbatim, so project rules no longer have to live outside the block"

patterns-established:
  - "gate.Env carries injected side-effecting closures (DeriveCommits) rather than importing gitrepo directly, keeping core/gate's 'pure functions, filesystem-free' promise intact while both CLI and TUI wire the same derivation"

requirements-completed: [D-01, D-02, D-03, D-04]

duration: ~25min
completed: 2026-08-23
status: complete
---

# Quick Task 260823-khj: jaira works out a ticket's commits itself, syncs it off the board, and stops overwriting your own rules — Summary

**jaira now derives a ticket's commit list from git (file history union'd with id-mentioning commits, oldest-first), stamps it once when the ticket leaves the board via a new `jaira sync` command, and preserves a project's own rules across agent-note regeneration through a `<!-- jaira:local -->` marker.**

## Performance

- **Duration:** ~25 min
- **Completed:** 2026-08-23
- **Tasks:** 3/3
- **Files modified:** 19 (6 created, 13 modified)

## Accomplishments

- A ticket that records no commits of its own is no longer refused at a `requires-commits` lane: `core/gitrepo.CommitsForTicket` derives the list from git itself, and `gate.Env.DeriveCommits` (injected, so `core/gate` stays filesystem- and git-free) is wired identically into both the CLI (`loadEnv`) and the TUI (`Model.gateEnv`), replacing all four ad-hoc `gate.Env{...}` construction sites in the TUI.
- `jaira sync <id>` files a terminal-lane ticket into `.jaira/sync/<initials>-<yyyymmdd>/`, stamping the derived-plus-explicit commit union onto the ticket file before the move; a non-terminal ticket is refused by name, naming the terminal lane and the `jaira move` that reaches it, with a stable `not_terminal` `--json` reason. `jaira archive` gained the same best-effort stamping.
- `Store.Restore` now searches the archive and every `.jaira/sync/*/` folder, refusing (naming both/all) rather than guessing on an ambiguous name, and resolves only the base name so a traversal attempt cannot escape the store.
- `core/board`'s regenerated agent note no longer instructs an agent to pass `--commits "$(git rev-parse HEAD)"` on every move; it explains the derivation and documents `jaira sync`/`jaira restore`. A new `<!-- jaira:local -->` marker preserves anything a project writes inside the managed block, verbatim, across every future regeneration.
- Verified end-to-end by hand in a throwaway repo: created a ticket, walked it to `done` (the terminal lane) without ever passing `--commits`, confirmed the lane accepted it, then `jaira sync --json` reported the git-derived sha stamped on the moved file, and `jaira restore` brought it back onto the board.

## Task Commits

Each task was committed atomically:

1. **Task 1: git answers "which commits carry this ticket", end to end, from the gate** - `35d9645` (feat)
2. **Task 2: jaira sync — the ticket leaves the board with every commit written down** - `e8bbe2d` (feat)
3. **Task 3: the managed block keeps your rules, and the note stops asking for a sha** - `6b584bd` (feat)

## Deviations from Plan

None — plan executed exactly as written. Two small implementation choices were made where the plan left the exact wording open, both consistent with the plan's intent:

- `Store.Restore`'s ambiguous-match error names all matching sync folders (and the archive, if it's also present) rather than picking one, per the "refuse rather than guess" behavior the plan specified — the `reason` code on `newRestoreCmd` (`not_archived`) was left unchanged since the plan did not ask for a new reason value there.
- `jaira archive`'s stamping step silently no-ops if `loadEnv` fails (e.g. a malformed lane file elsewhere on the board), so a pre-existing archive failure mode is not changed by this addition — consistent with "archive has never required commits and must not start refusing."

## Known Stubs

None.

## Threat Flags

None. The only input crossing a boundary — the `<file>` argument to `jaira restore` — is neutralised with `filepath.Base` in both the archive lookup and the new sync-folder lookup, covered by `TestRestoreCannotEscapeStore`.

## Deferred (out of scope)

- `internal/cli/tickets.go` has a pre-existing `gofmt` drift (two map-literal alignment blocks), unrelated to any file this task touched (confirmed via `git log -1` predating this work). Left as-is per the scope boundary rule; see `deferred-items.md` in this directory.

## Self-Check: PASSED

- `core/gitrepo/derive.go` — FOUND
- `core/gitrepo/derive_test.go` — FOUND
- `internal/cli/syncout.go` — FOUND
- `internal/cli/syncout_test.go` — FOUND
- `core/ticket/sync_test.go` — FOUND
- `core/board/announce_test.go` — FOUND
- Commit `35d9645` — FOUND in `git log --oneline --all`
- Commit `e8bbe2d` — FOUND in `git log --oneline --all`
- Commit `6b584bd` — FOUND in `git log --oneline --all`
