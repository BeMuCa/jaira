# Roadmap: jAIra

## Overview

jAIra starts as a durable, git-shareable ticket store you can drive entirely from bash — no TUI, no lanes beyond the seven built-ins, no agents — because that slice alone already delivers the core value proposition (never lose track of what a task was for or where it stands) and is independently useful to any bash-capable agent. From there, a read-only board proves the rendering path before any mutation exists, then the TUI gains write parity under the same gates the CLI already enforces. Conflict-tolerant frontmatter lands next, deliberately before any concurrent writer (agents, teammates) becomes real, because retrofitting conflict handling onto a naive model would be a data-format migration. Custom lanes make the pipeline configurable, the agent integration surface turns the board into something Claude can actually work autonomously, and Task-tool sync plus the session/live-refresh polish close the loop on multi-session visibility.

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [ ] **Phase 1: Foundation — Ticket Store, Schema & CLI** - Durable, git-shareable ticket store, driven entirely by a bash CLI, with proven byte-stable YAML writes and atomic/concurrency-safe file handling.
- [ ] **Phase 2: Read-Only TUI Board** - The board renders lanes, cards, detail, diffs, and filtering with no mutation path yet.
- [ ] **Phase 3: TUI Mutation Parity & Gates** - The TUI becomes a full write surface, and the promotion/dependency/Done gates are enforced identically from both interfaces at the moment of mutation.
- [ ] **Phase 4: Conflict-Tolerant Frontmatter** - Concurrent edits to the same ticket resolve automatically for the fields that change most, before any real concurrent writer exists.
- [ ] **Phase 5: Custom & Portable Lanes** - Lanes become readable, writable, reorderable and shareable: a global catalogue, per-project copies that stay private, opt-in sharing through `.jaira/shared/`, and built-ins that may be overridden but never silently.
- [ ] **Phase 6: Agent Pipeline** - Claude (or any bash-capable agent) can discover, claim, and advance tickets through lanes with tool-enforced, bounded I/O.
- [ ] **Phase 7: Task Tool Sync** - Claude's structured Task-tool list and the Backlog stay consistent without duplicating tickets or breaking the user's session.
- [ ] **Phase 8: Session Context & Live Refresh** - The board shows live session focus and picks up external changes automatically.

## Phase Details

### Phase 1: Foundation — Ticket Store, Schema & CLI
**Mode:** mvp
**Goal**: A developer or agent can create, read, list, and move tickets entirely through the CLI; ticket files are durable, git-shareable, byte-stable on write, and safe under concurrent CLI invocations; the tool ships as an installable static binary.
**Depends on**: Nothing (first phase)
**Requirements**: STORE-01, STORE-02, STORE-03, STORE-04, STORE-05, STORE-06, STORE-07, STORE-08, STORE-09, STORE-10, SCHEMA-01, SCHEMA-02, SCHEMA-03, SCHEMA-04, SCHEMA-05, SCHEMA-06, SCHEMA-07, CLI-01, CLI-02, CLI-03, CLI-04, CLI-05, CLI-06, CLI-07, CLI-08, CLI-09, CLI-10, LANE-01, DIST-01, DIST-02, DIST-03
**Success Criteria** (what must be TRUE):
  1. User can run `jaira init` in a fresh repo and get a working `.jaira/` structure (tickets dir, `.gitattributes`, session gitignore) with zero further configuration.
  2. User can create, list, show, and move a ticket via CLI against the seven built-in lanes, and every read command supports a stable `--json` output never mixed with human text.
  3. Editing one field via CLI leaves every other field, comment, blank line, and the reserved `external:` block byte-identical — proven by an AST-patch spike against real-looking ticket fixtures before the rest of the write path is built on it, and frontmatter with anchors/aliases the tool can't safely rewrite is detected and reported rather than silently mangled.
  4. Two concurrent CLI invocations against the same ticket never lose an update or leave a partial/corrupt file, verified by a concurrency test; re-running an already-satisfied command succeeds as a no-op.
  5. A teammate can clone the repo, run a binary built for their OS/arch with no toolchain, and see the same tickets and lanes with no configuration step.
**Plans**: TBD

### Phase 2: Read-Only TUI Board
**Mode:** mvp
**Goal**: Users can visually browse the board — lanes, cards, ticket detail, commit diffs, live filtering, and multiple projects — with no mutation capability yet, proving the core→TUI rendering path before write complexity is added.
**Depends on**: Phase 1
**Requirements**: TUI-01, TUI-02, TUI-03, TUI-04, TUI-07, TUI-08, TUI-12, TUI-13
**Success Criteria** (what must be TRUE):
  1. User can launch the TUI and see all seven base lanes rendered as columns with cards, matching what the CLI reports for the same repo.
  2. User can navigate cards and lanes with vim keys or arrow keys and expand a card to see full ticket detail inline.
  3. User can filter the visible cards live by typing a query, and view a ticket's combined commit diff without leaving the board.
  4. User can bring up an on-demand keybinding reference and switch between projects from within the TUI.
  5. The board starts instantly — no daemon, no index build, no network call.
**Plans**: TBD
**UI hint**: yes

### Phase 3: TUI Mutation Parity & Gates
**Mode:** mvp
**Goal**: The board becomes a full write surface — tickets can be created and moved from the TUI exactly as from the CLI — and every mutation, from either interface, is checked against the promotion, dependency, and Done gates at the moment of mutation, not only at discovery.
**Depends on**: Phase 2
**Requirements**: TUI-05, TUI-06, GATE-01, GATE-02, GATE-03, GATE-04, GATE-05, GATE-06, GATE-07
**Success Criteria** (what must be TRUE):
  1. User can create a ticket into the Backlog and move a card between lanes directly from the TUI, producing the same file state as the equivalent CLI command.
  2. A ticket missing `goal`, `definition-of-done`, `context`, or `assignee` cannot leave Backlog from either the CLI or the TUI, and the rejection names exactly which fields are missing; `assignee` defaults to the ticket's creator.
  3. A ticket with unresolved `blocked-by` entries cannot be started from either interface.
  4. A ticket cannot reach Done on an agent's own say-so — the transition requires a non-model signal (a passing command or explicit human sign-off).
  5. Gate checks fire at the moment of mutation, so calling the mutation command directly (skipping any discovery step) still gets blocked identically from CLI and TUI.
**Plans**: TBD
**UI hint**: yes

### Phase 4: Conflict-Tolerant Frontmatter
**Mode:** mvp
**Goal**: Concurrent edits to the same ticket — across clones, teammates, or parallel sessions — resolve automatically for the fields that change most often, before any real concurrent writer (agents, teammates) makes this load-bearing.
**Depends on**: Phase 3
**Requirements**: MERGE-01, MERGE-02, MERGE-03, MERGE-04, MERGE-05, MERGE-06, TUI-11
**Success Criteria** (what must be TRUE):
  1. Two sides moving the same ticket to different lanes merge automatically with no conflict, resolved by lane order so forward progress is never silently reverted.
  2. Concurrent additions to `blocked-by` or `commits` from different sides are both preserved after a merge.
  3. Competing scalar edits (e.g., `assignee`) resolve deterministically by the ticket's own recency, and a genuinely conflicting prose edit (two rewrites of `goal`) surfaces as a conflict scoped to that one field, never the whole file, and never produces invalid YAML or drops a ticket.
  4. The merge driver installs itself into local git config on first run and prints a one-line notice, rather than acting silently.
  5. If a ticket the user has open in the TUI is changed elsewhere while they're editing it, the conflict is surfaced rather than silently overwritten.
**Plans**: TBD
**UI hint**: yes

### Phase 5: Custom & Portable Lanes
**Mode:** mvp
**Goal**: A user can read, write, reorder and share lane definitions — the pipeline stops being a fixed thing baked into the binary and becomes something a person or an agent can inspect and change.
**Depends on**: Phase 1, Phase 2
**Requirements**: LANE-02, LANE-03, LANE-04, LANE-05, LANE-06, LANE-07, LANE-08
**Design**: `.planning/lane-system-design.md` (agreed with the user 2026-08-13; supersedes the original criterion 4, which forbade overriding a built-in)
**Success Criteria** (what must be TRUE):
  1. A lane is one markdown file. `~/.jaira/lanes/` is the global catalogue of everything this user built or adopted; a project's `.jaira/lanes/` holds copies of only the lanes that project actually uses.
  2. A project's lanes stay private: `.jaira/lanes/` is gitignored even after `jaira share`, which commits `tickets/` and `shared/` but not `lanes/`.
  3. A lane definition's model-tier is a local alias, not a hardcoded model name, so a shared lane file survives a model rename.
  4. Custom lanes anchor to an existing lane id (rather than a numeric position) and order correctly, including when a lane file's anchor isn't present locally.
  5. A custom lane MAY override a built-in, prompt included, and the override is always reported as a warning at load time — powerful, but never silent.
  6. A ticket sitting in a lane the local install doesn't recognize renders as a passthrough column instead of being hidden, and advancing a ticket into or out of an unrecognized lane is refused.
  7. A lane is legible without a ticket in hand: `jaira lanes show <id>` prints it in full including its prompt, `jaira lanes --json` carries the prompt, and `jaira lanes path` names the directory to write into — so an agent can read a lane, write a new one from a template, and verify the result with `jaira lanes`, without a second write path.
  8. A lane can be exported to `.jaira/shared/<user>/` from the TUI lane settings screen; teammates' shared lanes are visible when picking lanes, and adopting one copies it into the adopter's catalogue.
  9. Every lane records a `creator:` signature, so an adopted lane keeps its provenance.
  10. A default board, set once per user from the home screen, decides which lanes a newly initialised board gets and which ticket Options start ticked; a project that changed nothing still carries no lane files, because an absent `.jaira/lanes/` means the built-ins.
  11. An agent can do all of it without the TUI: create a lane, change one, and edit the default board, using `jaira lanes path` to find where files belong, a template to write against, and `jaira lanes` to check the result — including validation of the default board file, not only of lanes.
  12. The board stops claiming that `precedence` is a column position: order follows the `after:` anchor, and the number is labelled as what it is — the rank the merge driver uses to decide which lane wins when two clones moved the same ticket. Separately, a lane whose `input-requires` names a field that no earlier lane produces is reported at load time, the way a cycle already is; today nothing compares the contracts against the order at all.
**Plans**: 1 plan, 10 tasks
- [ ] `.planning/phase-5-custom-and-portable-lanes/PLAN.md` — decision checkpoint on three open questions, then root-aware lane loading, a legible `jaira lanes`, private project lanes through `share`, ordering/tier/unknown-lane proof, the TUI publish + adopt screens, a per-user default board, the agent-facing lane and default-board surface, and `precedence` as the single ordering mechanism with `after` as a checked constraint

### Phase 6: Agent Pipeline
**Mode:** mvp
**Goal**: Claude (or any bash-capable agent) can work the board end-to-end — discover actionable work, claim it, receive a bounded lane input assembled by the tool, and advance it only by supplying tool-validated structured output.
**Depends on**: Phase 3, Phase 4, Phase 5
**Requirements**: PIPE-01, PIPE-02, PIPE-03, PIPE-04, PIPE-05, PIPE-06, PIPE-07, PIPE-08
**Success Criteria** (what must be TRUE):
  1. Agent can query the next actionable ticket, and blocked or already-claimed tickets never appear in that result.
  2. Agent can take a time-limited claim on a ticket, and an abandoned claim becomes available again automatically with no manual unlocking.
  3. Agent receives a bounded lane input (the lane's required fields, the diff of the ticket's own commits, and the lane prompt) assembled by the tool, never by the agent itself.
  4. A ticket cannot advance unless the current lane's declared outputs were actually supplied, and malformed structured output is rejected with a readable, retryable reason.
  5. An agent can move a ticket to the Human lane with a question attached, and Claude Code can drive this whole loop in-session — spawning one subagent per lane step at that lane's model tier — guided by a skill that teaches the CLI surface and ticket schema.
**Plans**: TBD

### Phase 7: Task Tool Sync
**Mode:** mvp
**Goal**: Claude's structured Task-tool list and the Backlog stay consistent in both directions, matched by a jAIra ULID carried in task metadata, without ever duplicating tickets or breaking the user's session.
**Depends on**: Phase 1, Phase 3
**Requirements**: SYNC-01, SYNC-02, SYNC-03, SYNC-04, SYNC-05, SYNC-06, SYNC-07, SYNC-08
**Success Criteria** (what must be TRUE):
  1. A structured task Claude creates is mirrored into the Backlog as a ticket with `ready: false`, matched by a jAIra ULID stored in the task's `metadata` — not by content hashing.
  2. Running the sync again when nothing changed writes nothing and produces no git diff; repeated syncs update the same mapped ticket rather than creating a duplicate.
  3. A task disappearing from Claude's list never deletes its mirrored ticket, and starting a ticket surfaces it back in the in-session task list.
  4. Task dependencies and jAIra's `blocked-by` stay consistent across a sync, and a sync failure never breaks or blocks the user's Claude session.
  5. The sync adapter lives in one isolated module, so a future change to Claude's task tool API is a contained patch.
**Plans**: TBD

### Phase 8: Session Context & Live Refresh
**Mode:** mvp
**Goal**: The board reflects reality without a manual refresh — showing which sessions are active and what they're focused on, and picking up external changes (another session, a `git pull`, a manual edit) automatically without losing the user's place.
**Depends on**: Phase 6
**Requirements**: SESSION-01, SESSION-02, SESSION-03, SESSION-04, TUI-09, TUI-10
**Success Criteria** (what must be TRUE):
  1. Claude can record its current focus and reasoning via the CLI at checkpoints, and the board displays that focus as a live panel.
  2. Multiple concurrent sessions in the same working tree each appear on the board, with stale ones marked rather than deleted.
  3. Session state never appears in a git diff or commit.
  4. The board refreshes automatically when tickets change on disk — another session, a `git pull`, or a manual edit — without the user asking it to.
  5. A background refresh preserves the user's current selection and any open detail pane rather than resetting the view.
**Plans**: TBD

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 2 → 3 → 4 → 5 → 6 → 7 → 8

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Foundation — Ticket Store, Schema & CLI | 0/TBD | Not started | - |
| 2. Read-Only TUI Board | 0/TBD | Not started | - |
| 3. TUI Mutation Parity & Gates | 0/TBD | Not started | - |
| 4. Conflict-Tolerant Frontmatter | 0/TBD | Not started | - |
| 5. Custom & Portable Lanes | 0/TBD | Not started | - |
| 6. Agent Pipeline | 0/TBD | Not started | - |
| 7. Task Tool Sync | 0/TBD | Not started | - |
| 8. Session Context & Live Refresh | 0/TBD | Not started | - |
