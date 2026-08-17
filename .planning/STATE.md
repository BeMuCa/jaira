# Project State

## Project Reference

See: .planning/PROJECT.md (updated 2026-08-11)

**Core value:** You never lose track of what a task was for or where it stands — across session boundaries, agent runs, and teammates.
**Current focus:** Phase 1 — Foundation: Ticket Store, Schema & CLI

## Current Position

Phase: 1 of 8 (Foundation — Ticket Store, Schema & CLI)
Plan: TBD (not yet planned)
Status: Ready to plan
Last activity: 2026-08-18 — Quick task 260818-1o3: lane removal confirmation (default no) and catalogue-file delete via x. Before that: Phase 5 (Custom & Portable Lanes) complete, plus eight quick tasks: TUI board setup, ticket id copy, context consolidation, review traceability, version stamp and `jaira update`, the specified-zone gate with a brainstorm step, ticket ownership, and the progress-note contract

Progress: [░░░░░░░░░░] 0%

## Performance Metrics

**Velocity:**
- Total plans completed: 0
- Average duration: N/A
- Total execution time: 0 hours

**By Phase:**

| Phase | Plans | Total | Avg/Plan |
|-------|-------|-------|----------|
| - | - | - | - |

**Recent Trend:**
- Last 5 plans: N/A
- Trend: N/A

*Updated after each plan completion*

## Accumulated Context

### Decisions

Decisions are logged in PROJECT.md Key Decisions table.
Recent decisions affecting current work:

- Roadmap: YAML round-trip fidelity spike is load-bearing and placed inside Phase 1 (not a later phase) — the AST-patch approach must be proven against real ticket fixtures before the write path is built on it.
- Roadmap: Conflict-tolerant frontmatter (Phase 4) is sequenced strictly before the agent pipeline (Phase 6), since concurrency only becomes real once agents work the board.
- Roadmap: Atomic writes and concurrent-CLI safety (STORE-06/07) live in Phase 1 with the store itself, not deferred to the merge-driver phase.
- Roadmap: Task Tool Sync (Phase 7) targets Claude Code's current structured Task tools (TaskCreate/TaskUpdate/TaskGet/TaskList) and identity-matches via a jAIra ULID in task metadata, not content hashing — TodoWrite is confirmed absent from the current tool surface.

### Pending Todos

None yet.

### Blockers/Concerns

- Phase 1: YAML round-trip fidelity via `goccy/go-yaml` AST editing is MEDIUM confidence per STACK.md — no Go/Rust library guarantees byte-perfect round-trip for every legal YAML construct. Mitigated by validating/rejecting anchors/aliases on read (STORE-10) rather than relying on the editor to survive them.
- Phase 7: Claude Code's Task tool API is an external, moving surface (already replaced TodoWrite once). Sync adapter must stay isolated in one module per SYNC-07.

## Quick Tasks Completed

| Date | Task | Outcome |
|------|------|---------|
| 2026-08-18 | 260818-1o3 | Every `x` on the lane settings screen now asks yes/no first — no listed first and preselected, h/l or arrows switch, enter acts, esc cancels — so a hasty enter never deletes. And `x` finally reaches not-installed catalogue lanes: it deletes their file (the stale `my-lane` skeleton a closed-without-saving editor leaves behind was undeletable from the TUI); a built-in without a file errors instead ([260818-1o3](./quick/260818-1o3-lane-settings-x-deletes-catalogue-lane-f/)) |
| 2026-08-12 | dod-verb | `jaira dod` plus a surgical checklist writer: one marker character changes, indexes are scoped to their section, one in-progress item per checklist |
| 2026-08-12 | dod-checkbox-states | Three checklist states, `## Plan` section parsed separately, and the fix for `[~]` items being dropped by the parser — which let a ticket with outstanding work enter the terminal lane |
| 2026-08-13 | 260813-nzq | The working lane never named `jaira note` at all, so nothing told an agent to leave a trail. All three working prompts, the `note` help and the agent note now say what a note contains — what you were doing, what you found and especially what did not work, what the next step is — and when one is due: before stopping, after something failed, when the plan turns out wrong. Deliberately no gate: a note written because a gate demanded one is worthless ([260813-nzq](./quick/260813-nzq-say-what-a-progress-note-must-contain-an/)) |
| 2026-08-13 | 260813-mqy | A ticket belongs to its assignee: writing to someone else's is refused with exit 3, naming the owner. The human checkpoint lanes are exempt by contract (`requires-question`, `requires-human-exit`), because reviewing someone else's work is their purpose; a hand-over is always allowed so an absent owner cannot freeze a ticket; `--force` overrides and is recorded. Deliberately a guard rail, not a lock — git still merges files, so the merge rules are untouched ([260813-mqy](./quick/260813-mqy-a-ticket-belongs-to-its-assignee-and-the/)) |
| 2026-08-13 | 260813-eir | The promotion gate fires on entering the specified zone rather than on leaving a lane named `backlog`, keyed on a new `requires-specified` lane contract that the built-in todo lane carries. That made a thinking lane possible, so a `brainstorm` step ships with it — optional per ticket, and it cannot be left until it produced a `goal` ([260813-eir](./quick/260813-eir-gate-on-entering-the-specified-zone-not-/)) |
| 2026-08-13 | 260813-1br | A board now records which jaira version last prepared it, in the per-clone state directory so a shared board never conflicts on it. Any command whose binary is newer prints one line on stderr pointing at `jaira update`, which re-applies the setup and prints the change notes embedded in the binary; `--json` carries the same for an agent. A `dev` build never nags ([260813-1br](./quick/260813-1br-version-stamp-and-jaira-update-with-embe/)) |
| 2026-08-13 | 260813-1as | The review lane now has to say what the change did, not only whether it passed: `review-summary` and `review-gaps` are in its output contract, so the gate refuses to let a ticket leave review without them. Every definition-of-done and plan item can carry a `proof:` sub-line written with `jaira dod <id> <n> --done --proof "…"`, which is where the evidence that a criterion was really met now lives ([260813-1as](./quick/260813-1as-traceability-review-summary-gaps-and-pro/)) |
| 2026-08-13 | 260813-19h | One place for the problem: `## Description` is gone from the ticket template, because nothing ever read it — no gate required it and no lane's bounded input could reach the markdown body. Everything now lives in `context`, which is written and read as a YAML block literal. That also fixed a silent data bug: a hand-written `context: |` used to read back as the string `"|"`, losing the whole text ([260813-19h](./quick/260813-19h-one-place-for-the-problem-multiline-cont/)) |
| 2026-08-13 | 260813-0z0 | The detail pane shows the full ticket id and `y` copies it to the clipboard over OSC52 (`tea.SetClipboard`, no new dependency), so an id can be pasted into an agent prompt. The full id is also rendered as selectable text, because a terminal with OSC52 disabled would otherwise leave no way to get it ([260813-0z0](./quick/260813-0z0-show-and-copy-the-full-ticket-id-in-the-/)) |
| 2026-08-12 | 260812-wz3 | A board created from the TUI browse screen ("i") was neither gitignored nor announced to any agent — only `jaira init` did that. Privacy and agent-note logic moved into `core/board`, both init paths call `board.Prepare`, and the note now names the whole working loop instead of four commands ([260812-wz3](./quick/260812-wz3-fix-tui-board-creation-to-write-gitignor/)) |

## Deferred Items

Items acknowledged and carried forward from previous milestone close:

| Category | Item | Status | Deferred At |
|----------|------|--------|-------------|
| *(none)* | | | |

## Session Continuity

Last session: 2026-08-11
Stopped at: ROADMAP.md and STATE.md created; REQUIREMENTS.md traceability table updated
Resume file: None

## Build status — 2026-08-11

All eight roadmap phases implemented in one pass. 83 of 84 v1 requirements met and
verified by running code; TUI-11 deferred (see REQUIREMENTS.md) because the board
has no in-place field editor, so there is no edit buffer for a background refresh
to clobber.

Verified rather than asserted:

- byte-faithful single-field writes (16 unit tests; a real `git diff` showing a
  one-line change)
- field-aware merge driver under two real git branches: lane resolved by progress,
  blockers and commits unioned, nothing lost, clean tree
- competing prose rewrites conflict, ticket stays valid YAML, `jaira resolve`
  settles it
- 20 concurrent writers to one ticket: 20/20 writes landed, no corruption, all
  locks released
- promotion gate, dependency gate and the non-model Done signal all refuse with
  documented exit codes
- task-tool sync is idempotent and creates no duplicates on re-sync
- six cross-compile targets build with CGO disabled
- `go vet` clean, `go test ./... -race` green, `gofmt` clean

Known unproven: adoption. File-based trackers have a poor track record and this
one's bet — a field-aware merge driver plus a deliberately small surface — has not
been tested by a real team yet.
