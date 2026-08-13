# Requirements: jAIra

**Defined:** 2026-08-11
**Core Value:** You never lose track of what a task was for or where it stands — across session boundaries, agent runs, and teammates.

## v1 Requirements

### Ticket Store

- [x] **STORE-01**: A ticket is a single markdown file with YAML frontmatter under `.jaira/tickets/`, readable and editable by hand
- [x] **STORE-02**: Ticket IDs are ULIDs, generated without any shared counter, so two parallel sessions never collide
- [x] **STORE-03**: User can refer to a ticket by an unambiguous ID prefix, the way git accepts short SHAs
- [x] **STORE-04**: Writing one field leaves every other field, comment, and blank line in the file byte-identical
- [x] **STORE-05**: An unrecognized frontmatter block (such as the reserved `external:`) survives any write untouched
- [x] **STORE-06**: Ticket writes are atomic — an interrupted write never leaves a partial or corrupt file
- [x] **STORE-07**: Two concurrent CLI invocations cannot clobber each other's writes to the same ticket
- [x] **STORE-08**: Board state persists across sessions — closing and reopening shows the same tickets in the same lanes
- [x] **STORE-09**: A teammate who clones the repo and runs `jaira` sees the same board with no configuration step
- [x] **STORE-10**: Frontmatter containing YAML constructs the tool cannot safely rewrite (anchors, aliases) is detected and reported rather than silently mangled

### Ticket Schema

- [x] **SCHEMA-01**: Ticket records `title`, `status`, `goal`, `context`, `definition-of-done`, `created-at`, `updated-at`
- [x] **SCHEMA-02**: Ticket records `creator` and `assignee`, both human identities
- [x] **SCHEMA-03**: Ticket records `executed-by` per run, naming the model that did the work, separately from `assignee`
- [x] **SCHEMA-04**: Ticket declares `blocked-by` as a flat list of other ticket IDs
- [x] **SCHEMA-05**: Ticket records an outcome block with three distinct fields — what was done, why it was needed, and how it satisfies the Definition of Done
- [x] **SCHEMA-06**: Ticket records the commit SHAs produced for it
- [x] **SCHEMA-07**: Ticket reserves an `external:` block for a future Jira/YouTrack adapter, unused in v1

### CLI

- [x] **CLI-01**: User can create a ticket from the command line
- [x] **CLI-02**: User can list tickets, filtered by lane
- [x] **CLI-03**: User can show one ticket in full
- [x] **CLI-04**: User can move a ticket to another lane
- [x] **CLI-05**: Every read command supports `--json` emitting a stable machine-readable schema, never mixed with human text
- [x] **CLI-06**: Exit codes are stable and documented — success, generic error, usage error, validation failure, dependency-blocked
- [x] **CLI-07**: Errors under `--json` are emitted as structured JSON on stderr, parseable without regex
- [x] **CLI-08**: Re-running a command that is already satisfied succeeds as a no-op rather than erroring
- [x] **CLI-09**: `jaira init` prepares a repo — creates `.jaira/`, `.gitattributes`, and the local gitignore for ephemeral session state
- [x] **CLI-10**: The CLI is driven entirely through bash by any agent, requiring no MCP server or runtime

### Gates

- [x] **GATE-01**: A ticket cannot leave Backlog until `goal`, `definition-of-done`, `context`, and `assignee` are all present
- [x] **GATE-02**: A blocked gate reports exactly which fields are missing, so an agent can fix and retry
- [x] **GATE-03**: A ticket whose `blocked-by` tickets are unresolved cannot be started
- [x] **GATE-04**: Gate rules are enforced identically whether the caller is the CLI or the TUI
- [x] **GATE-05**: Gates are enforced at the moment of mutation, not only at discovery, so an agent cannot bypass them by skipping the discovery step
- [x] **GATE-06**: `assignee` defaults to the ticket's creator
- [x] **GATE-07**: A ticket cannot enter Done on an agent's own say-so — the transition requires a non-model signal (passing command or explicit human sign-off)

### Board / TUI

- [x] **TUI-01**: Board renders lanes as columns with tickets as cards
- [x] **TUI-02**: User can navigate cards and lanes with vim keys and arrow keys
- [x] **TUI-03**: User can expand a card to see full ticket detail inline
- [x] **TUI-04**: User can filter tickets live by typing a query
- [x] **TUI-05**: User can create a ticket into the Backlog from the board
- [x] **TUI-06**: User can move a card between lanes from the board
- [x] **TUI-07**: User can view the combined diff of a ticket's commits without leaving the board
- [x] **TUI-08**: Board shows a keybinding reference on demand
- [x] **TUI-09**: Board refreshes automatically when tickets change on disk — another session, a `git pull`, or a manual edit
- [x] **TUI-10**: A background refresh preserves the user's current selection and any open detail pane
- [ ] **TUI-11**: If a ticket is changed elsewhere while the user is editing it, the conflict is surfaced rather than silently overwritten — **not met**: the board has no in-place field editor (only create and move), so there is no edit buffer to clobber. Deferred with the field editor itself.
- [x] **TUI-12**: User can switch between projects from within the board
- [x] **TUI-13**: Board starts instantly — no daemon, no index build, no network

### Lanes

- [x] **LANE-01**: Seven base lanes ship inside the binary and need no repo configuration: Backlog, Todo, In Progress, Human, Review, Done, Blocked
- [x] **LANE-02**: A custom lane is a single portable markdown file in `~/.jaira/lanes/`, shareable by sending the file
- [x] **LANE-03**: A lane definition specifies its prompt, its model tier, and its input/output contract
- [x] **LANE-04**: A lane's model tier is a local alias, not a hardcoded model name, so shared lane files survive model renames
- [x] **LANE-05**: Custom lanes order themselves by anchoring to an existing lane rather than a numeric position
- [x] **LANE-06**: A lane whose `id` collides with a base lane overrides it, and the override is always reported as a warning, never silent
- [x] **LANE-07**: A ticket in a lane the local install does not have renders as a read-only passthrough column, never hidden
- [x] **LANE-08**: Advancing a ticket into or out of an unrecognized lane is refused, since no contract exists to enforce

### Agent Pipeline

- [x] **PIPE-01**: Agent can query the next actionable ticket, excluding blocked and already-claimed work
- [x] **PIPE-02**: Agent can take a time-limited claim on a ticket, and a stale claim is treated as abandoned without manual unlocking
- [x] **PIPE-03**: Agent receives a bounded input assembled by the tool — the lane's required ticket fields, the diff of the ticket's own commits, and the lane prompt — never assembled by the agent itself
- [x] **PIPE-04**: Agent advances a ticket by supplying structured output, which the tool validates against the lane's declared output schema and rejects with a readable reason if malformed
- [x] **PIPE-05**: A ticket cannot advance unless the current lane's declared outputs were actually supplied
- [x] **PIPE-06**: Agent can move a ticket to the Human lane with a question attached
- [x] **PIPE-07**: Claude works the board in-session, spawning one subagent per lane step at that lane's model tier
- [x] **PIPE-08**: A skill teaches Claude the CLI surface and the ticket schema

### Conflict Handling

- [x] **MERGE-01**: Two people moving the same ticket to different lanes resolves automatically rather than producing a conflict
- [x] **MERGE-02**: Concurrent additions to `blocked-by` or `commits` are both preserved
- [x] **MERGE-03**: Competing scalar edits resolve deterministically and never silently revert forward lane progress
- [x] **MERGE-04**: Genuinely conflicting prose edits surface as a conflict scoped to the affected field, not the whole file
- [x] **MERGE-05**: The merge driver installs itself on first run and announces that it did so, rather than acting silently
- [x] **MERGE-06**: A merge never produces invalid YAML or loses a ticket

### Session Context

- [x] **SESSION-01**: Claude records the current focus and reasoning via the CLI at checkpoints
- [x] **SESSION-02**: Board displays the current session focus as a live panel
- [x] **SESSION-03**: Multiple concurrent sessions in the same working tree each appear, with stale ones marked rather than deleted
- [x] **SESSION-04**: Session state is never committed to git

### Task Tool Sync

- [x] **SYNC-01**: Claude's structured tasks are mirrored into the Backlog as tickets with `ready: false`
- [x] **SYNC-02**: Mirrored tickets are matched by a stable ID carried in task metadata, so repeated syncs update rather than duplicate
- [x] **SYNC-03**: A task disappearing from Claude's list never deletes its mirrored ticket
- [x] **SYNC-04**: Syncing when nothing changed writes nothing and produces no git diff
- [x] **SYNC-05**: Starting a ticket surfaces it in the in-session task list, so work begun in chat stays visible in both places
- [x] **SYNC-06**: Task dependencies and jAIra's `blocked-by` stay consistent across a sync
- [x] **SYNC-07**: The sync adapter is isolated in one module, so a change to Claude's task API is a contained patch
- [x] **SYNC-08**: A sync failure never breaks the user's Claude session

### Distribution

- [x] **DIST-01**: The tool is a single static binary with no runtime dependency
- [x] **DIST-02**: Binaries are produced for Linux, macOS, and Windows on amd64 and arm64
- [x] **DIST-03**: A teammate can install and run it without a toolchain

## v2 Requirements

### External Sync

- **EXT-01**: Jira adapter reading and writing through the reserved `external:` block
- **EXT-02**: YouTrack adapter
- **EXT-03**: Field mapping configuration between jAIra and the external tracker

### Convenience

- **CONV-01**: Named filter shortcuts, once users are observed retyping the same query
- **CONV-02**: Archival of old Done tickets into a subdirectory, once ticket volume makes scanning slow
- **CONV-03**: An optional gitignored on-disk cache behind the existing store interface, only if startup measurably degrades
- **CONV-04**: An MCP server wrapping the same core, for agents that prefer structured tools over bash

## Out of Scope

| Feature | Reason |
|---------|--------|
| Sprints | Full-PM-tool feature; no lightweight or agent-focused tracker surveyed has them |
| Custom fields | The reserved `external:` block already serves as the escape hatch; a second one is redundant |
| Roles / permissions | Git repository access already *is* the permission system |
| Saved views, dashboards | Inherently server/web features; conflicts with the no-server constraint |
| Server, accounts, authentication | Git is the sync layer; a server would reintroduce exactly the weight being escaped |
| Web UI, VS Code extension | Splitting across UI surfaces is the bloat trap; sharing is solved by git |
| Branch per ticket | Too heavy for small tasks and breaks down under parallel agents |
| Board-spawned headless agent runs | Background process lifecycle is a large surface; orchestration stays in-session |
| Multi-type dependency graph and graph queries | A flat `blocked-by` list covers the actual need; graph tooling solves a problem at a scale this project does not have |
| CRDT-based ticket engine | Solves conflicts completely but abandons hand-editable plain files — a harder constraint violation than tolerating rare conflicts |
| Building on the deprecated TodoWrite shape | Verified replaced by structured Task tools; building on the old shape would target a dead API |
| Feature parity with paca | Explicit anti-goal — the project fails by growing |

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| STORE-01 | Phase 1 | Pending |
| STORE-02 | Phase 1 | Pending |
| STORE-03 | Phase 1 | Pending |
| STORE-04 | Phase 1 | Pending |
| STORE-05 | Phase 1 | Pending |
| STORE-06 | Phase 1 | Pending |
| STORE-07 | Phase 1 | Pending |
| STORE-08 | Phase 1 | Pending |
| STORE-09 | Phase 1 | Pending |
| STORE-10 | Phase 1 | Pending |
| SCHEMA-01 | Phase 1 | Pending |
| SCHEMA-02 | Phase 1 | Pending |
| SCHEMA-03 | Phase 1 | Pending |
| SCHEMA-04 | Phase 1 | Pending |
| SCHEMA-05 | Phase 1 | Pending |
| SCHEMA-06 | Phase 1 | Pending |
| SCHEMA-07 | Phase 1 | Pending |
| CLI-01 | Phase 1 | Pending |
| CLI-02 | Phase 1 | Pending |
| CLI-03 | Phase 1 | Pending |
| CLI-04 | Phase 1 | Pending |
| CLI-05 | Phase 1 | Pending |
| CLI-06 | Phase 1 | Pending |
| CLI-07 | Phase 1 | Pending |
| CLI-08 | Phase 1 | Pending |
| CLI-09 | Phase 1 | Pending |
| CLI-10 | Phase 1 | Pending |
| LANE-01 | Phase 1 | Pending |
| DIST-01 | Phase 1 | Pending |
| DIST-02 | Phase 1 | Pending |
| DIST-03 | Phase 1 | Pending |
| TUI-01 | Phase 2 | Pending |
| TUI-02 | Phase 2 | Pending |
| TUI-03 | Phase 2 | Pending |
| TUI-04 | Phase 2 | Pending |
| TUI-07 | Phase 2 | Pending |
| TUI-08 | Phase 2 | Pending |
| TUI-12 | Phase 2 | Pending |
| TUI-13 | Phase 2 | Pending |
| TUI-05 | Phase 3 | Pending |
| TUI-06 | Phase 3 | Pending |
| GATE-01 | Phase 3 | Pending |
| GATE-02 | Phase 3 | Pending |
| GATE-03 | Phase 3 | Pending |
| GATE-04 | Phase 3 | Pending |
| GATE-05 | Phase 3 | Pending |
| GATE-06 | Phase 3 | Pending |
| GATE-07 | Phase 3 | Pending |
| MERGE-01 | Phase 4 | Pending |
| MERGE-02 | Phase 4 | Pending |
| MERGE-03 | Phase 4 | Pending |
| MERGE-04 | Phase 4 | Pending |
| MERGE-05 | Phase 4 | Pending |
| MERGE-06 | Phase 4 | Pending |
| TUI-11 | Phase 4 | Pending |
| LANE-02 | Phase 5 | Pending |
| LANE-03 | Phase 5 | Pending |
| LANE-04 | Phase 5 | Pending |
| LANE-05 | Phase 5 | Pending |
| LANE-06 | Phase 5 | Pending |
| LANE-07 | Phase 5 | Pending |
| LANE-08 | Phase 5 | Pending |
| PIPE-01 | Phase 6 | Pending |
| PIPE-02 | Phase 6 | Pending |
| PIPE-03 | Phase 6 | Pending |
| PIPE-04 | Phase 6 | Pending |
| PIPE-05 | Phase 6 | Pending |
| PIPE-06 | Phase 6 | Pending |
| PIPE-07 | Phase 6 | Pending |
| PIPE-08 | Phase 6 | Pending |
| SYNC-01 | Phase 7 | Pending |
| SYNC-02 | Phase 7 | Pending |
| SYNC-03 | Phase 7 | Pending |
| SYNC-04 | Phase 7 | Pending |
| SYNC-05 | Phase 7 | Pending |
| SYNC-06 | Phase 7 | Pending |
| SYNC-07 | Phase 7 | Pending |
| SYNC-08 | Phase 7 | Pending |
| SESSION-01 | Phase 8 | Pending |
| SESSION-02 | Phase 8 | Pending |
| SESSION-03 | Phase 8 | Pending |
| SESSION-04 | Phase 8 | Pending |
| TUI-09 | Phase 8 | Pending |
| TUI-10 | Phase 8 | Pending |

**Status after v1 build:**
- Implemented and verified: 83 of 84
- Not met: TUI-11 (see note above)

**Coverage:**
- v1 requirements: 84 total
- Mapped to phases: 84
- Unmapped: 0 ✓

---
*Requirements defined: 2026-08-11*
*Last updated: 2026-08-11 after roadmap creation*
