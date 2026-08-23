# Project State

## Project Reference

See: .planning/PROJECT.md (updated 2026-08-11)

**Core value:** You never lose track of what a task was for or where it stands — across session boundaries, agent runs, and teammates.
**Current focus:** Phase 1 — Foundation: Ticket Store, Schema & CLI

## Current Position

Phase: 1 of 8 (Foundation — Ticket Store, Schema & CLI)
Plan: TBD (not yet planned)
Status: Ready to plan
Last activity: 2026-08-23 — quick task 260823-khj: jaira derives a ticket's commit list from git itself (union of the ticket file's history and commits naming its id), so no lane change asks for a sha; `jaira sync` files a finished ticket under `.jaira/sync/<initials>-<date>/` with its commits stamped on the way out; and a `jaira:local` area inside the managed CLAUDE.md block survives regeneration. Before that: next --per-lane maps where work waits, and the review lane's semantics were untangled: the user's board treats review as their own inbox, so it is now a human lane there and signoff is gone. Before that: quick task 260820-t7q: review-check (how to check it, as a flow) is enforced by the review lane, the review fields are finally in --json and 'jaira show', next_lane names the route on every ticket, and the optimise lane is optimize. Before that: quick task 260820-g80: rejects-to declares a loop's back edge, critique and optimise ship as catalogue lanes (installed on the req board), and the ticket rides in the code's commit. Before that: quick task 260820-2eo: the writing register gained clean formatting and a reader who knows none of what the writer knows, in all five shipped copies. Before that: quick task 260820-0dr: an accepted ticket leaves the board once its work is pushed, written down in the done lane, SKILL.md, AGENTS.md and the README (documentation, not a gate). Before that: quick tasks 260819-p3m / p9s / pk3, one session: the TUI stops losing your page (single-lane views hold their lane when a ticket is moved from another shell; dialogs and messages dismiss back to the page they were opened on), a gate refusal can be overridden in place with f then y instead of being sent to the CLI, and n on an open ticket writes its follow-up beside it in a split view. Before that: quick task 260819-hco: the agent contracts now demand a board search before create, and a user question on any contradiction. Before that: quick task 260819-djm: q uniformly goes one level back (lane focus no longer quits), the board hint bar names v, the off-screen notice is removed. Before that: quick task 260818-lxx: board lanes cap to the terminal height (flags-row width fix + clampBlock) and the single-lane view scrolls. Before that: quick task 260818-1o3: lane removal confirmation (default no) and catalogue-file delete via x. Before that: Phase 5 (Custom & Portable Lanes) complete, plus eight quick tasks: TUI board setup, ticket id copy, context consolidation, review traceability, version stamp and `jaira update`, the specified-zone gate with a brainstorm step, ticket ownership, and the progress-note contract

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
| 2026-08-23 | 260823-khj | Three habits that worked *around* jaira became part of it (issue #4). Nobody transcribes shas any more: `gitrepo.CommitsForTicket` derives a ticket's commit list from git itself — the union of the ticket file's own `--follow` history and commits whose message names the id, matched as a bounded reference so `fix(A3K9QP):` counts and `XA3K9QPZ1` does not, deduped oldest-first. The gate reaches it through an injected `Env.DeriveCommits` rather than importing gitrepo, so `core/gate` keeps its filesystem-free promise and its old tests still build an Env without git; an explicit `--commits` still wins, and an empty derivation still refuses with today's message. `jaira sync <id>` gives a finished ticket somewhere to go that is not the archive: stamps the commits, files it under `.jaira/sync/<initials>-<date>/`, off the board — terminal lane only, `restore` now searches both places and refuses an ambiguous name. And `<!-- jaira:local -->` inside the managed block keeps a project's own rules verbatim across regeneration, so the instructions an agent reads are one text instead of two that disagree ([260823-khj](./quick/260823-khj-issue-4-derive-ticket-commits-from-git-j/)) |
| 2026-08-21 | 260821-perlane | "critique is not being worked" was starvation, not a bug: `jaira next` returns the furthest-along ticket, so 28 waiting in review outranked 2 in critique every time. `next --per-lane` now maps the front line (every lane with work, pipeline order, agentic flag, one ticket). The finding underneath was bigger: the board used `review` as the user's inbox while the tool held it to be the model's lane, whose four output fields could never be produced (no commits, so no diff) — nothing had ever reached Human Review. On the user's board `review` is now a human lane and `signoff` is removed. `optimize` also produces `review-check`, being the last agent before a person ([260821-perlane](./quick/260821-per-lane-and-the-human-review-lane/)) |
| 2026-08-20 | 260820-t7q | `review-check` is the review's fourth output: how a person checks the change themselves, as a flow, enforced by the lane's contract and printed last on the sign-off screen (read the account, then go and look). Wiring it up surfaced that the review fields lived only in the TUI — `ticketJSON` and `jaira show` carried none of them, so an agent could not hand a reader a review at all; both now carry the whole block. `next_lane` on every `--json` ticket answers the route once instead of every agent deriving it: forward through the column order, skipping opted-out steps, never a parking or question lane. Verified against the real req board, which caught a ticket in HITL claiming `review` and skipping critique and optimize; parked and waiting tickets now report empty. Plus the optimise → optimize rename ([260820-t7q](./quick/260820-t7q-a-review-says-how-to-check-it-and-a-tick/)) |
| 2026-08-20 | 260820-g80 | Loops are declarable and the loop the user wanted exists. `rejects-to:` names a lane's back edge (validated: must resolve, may not be itself; counted as drift; in `lanes show` and `--json`). Two catalogue lanes ship under `lanes/` — `critique` (judges the approach, fixes nothing itself, sends decisions to HITL) and `optimise` (edits: duplication, dead code, fluff, cost) — deliberately NOT built-ins, since built-ins are injected into every board. `core/lane/shipped_test.go` parses them in CI and caught an unquoted colon that had cost a lane its id. Third strand: the ticket now rides in the same commit as the code, in all three agent contracts, because jaira never commits and nothing had said so ([260820-g80](./quick/260820-g80-declare-a-loop-back-edge-and-ship-critiq/)) |
| 2026-08-20 | 260820-2eo | The writing register asked for short, never for tidy, and assumed the reader knew everything the writer knew. Now one sentence in all five shipped places (`announce.go`, `create` help, `note` help, the pre-process lane prompt, `docs/AGENTS.md`): "as if the reader has mild ADHD **and knows none of what you know** - {what} first, short concrete lines with one point each, names and paths rather than adjectives, no jargon and no preamble". The word "dummy" was deliberately not shipped into public CLI help and rendered as "knows none of what you know" plus "no jargon" ([260820-2eo](./quick/260820-2eo-the-writing-register-gains-formatting-an/)) |
| 2026-08-20 | 260820-0dr | Archiving finally has a place in the flow. Nothing said when a ticket should come off the board: `archive` was documented as a command and a key, `push` appeared in no lane file and nowhere in `core/lane` or `core/gate`. Now written in four places (the `done` lane's description, SKILL.md, docs/AGENTS.md, README) with the reason for its order: an accepted ticket leaves the board once its work is pushed, because archiving ahead of the push hands a teammate a board that has forgotten a ticket whose code has not arrived. Documentation only, no gate: archiving is the step after the last lane, so there is nothing for a gate to fire on. Projects with their own copy of `done` need `jaira lanes use done --force` to see the new description ([260820-0dr](./quick/260820-0dr-an-accepted-ticket-leaves-the-board-once/)) |
| 2026-08-19 | 260819-pk3 | `n` on an open ticket writes its follow-up beside it: the screen splits, the ticket it follows stays on the left, so the reason for the new ticket is visible while it is written. The follow-up is a draft in memory until ctrl+s — esc discards it and the board never saw it — then it lands in the default lane with `follows` set, `ready` derived by the gate, and reads as a ticket with an open ticket's keys. `n` again chains and the older ticket slides off. The editor keeps tab for its fields, so the left pane scrolls with shift+arrows; once saved tab moves between panes. No split below 80 columns or 20 rows. `followUpFields` is now shared with the `f`-at-sign-off path ([260819-pk3](./quick/260819-pk3-create-a-follow-up-beside-the-open-ticke/)) |
| 2026-08-19 | 260819-p9s | A gate refusal is overridable where it is read: the message says "press f to override", `f` asks again, `y` writes, and anything else drops the offer. Same mutation and same wording as the CLI's `--force`, covering every refusal code, recording nothing on the ticket because the CLI records nothing either. Underneath: a new `returnTo` means the move picker and every message dismiss to the page they were opened on instead of always the board, and a successful move from an open ticket leaves it open showing the lane it landed in ([260819-p9s](./quick/260819-p9s-force-a-refused-move-from-the-tui-and-di/)) |
| 2026-08-19 | 260819-p3m | The compact view and lane focus no longer get swapped out from under you when the selected ticket is moved from another shell. `rebuild()` re-selected by ID and `selectByID` sets the lane index too — but in those views the lane index is the page, not a cursor. They now hold their lane by ID and let the cursor fall to whatever slid into the gap; the multi-column board still follows the ticket, where the new lane is already on screen ([260819-p3m](./quick/260819-p3m-cursor-stays-in-the-lane-you-are-looking/)) |
| 2026-08-19 | 260819-hco | Agents must check the board before `jaira create`: search with `list -q`/`list --json`, read close hits with `show --json`, and on a contradiction (an existing ticket decided the same question differently) stop and ask the user — both handles, the contradicting line quoted, and the ways forward (adjust, supersede with `--follows`, drop). Landed in SKILL.md and docs/AGENTS.md; deliberately not in the binary — judging "decided differently" is the LLM layer's job, and the CLI already had every query needed ([260819-hco](./quick/260819-hco-before-creating-a-ticket-agents-must-che/)) |
| 2026-08-19 | 260819-djm | `q` is always one level back: lane focus now returns to the compact view instead of quitting (ctrl+c still quits; footer and help updated), matching every sub-screen where q already went back. Board and compact stay equivalent — `v` toggles, `q` from either exits the board program, which lands on the project launcher via runHome's loop. The board's hint bar finally names `v compact`, and the "N lane(s) off-screen" notice is gone by Berk's call ([260819-djm](./quick/260819-djm-q-goes-one-level-back-everywhere-board-h/)) |
| 2026-08-18 | 260818-lxx | A lane never grows past the terminal again. Root cause: `renderCard` truncated its flags row at `w+24` — a stale ANSI fudge from before `truncate` counted display width — so a flag-heavy review card wrapped to 4-5 lines while the scroll math assumed 3, and a full Review lane pushed the board off-screen. Fixed to `w` (meta line too), plus `clampBlock` as a hard w×h backstop on every column. The single-lane view (`v`), which rendered all tickets uncapped, now scrolls with the cursor, shares the board's scroll state, and says "+N more" above and below ([260818-lxx](./quick/260818-lxx-board-lanes-cap-their-height-to-the-term/)) |
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
