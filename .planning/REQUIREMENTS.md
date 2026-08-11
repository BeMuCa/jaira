# Requirements: jAIra

**Defined:** 2026-08-11
**Core Value:** You never lose track of what a task was for or where it stands — across session boundaries, agent runs, and teammates.

## v1 Requirements

### Ticket Store

- [ ] **STORE-01**: A ticket is a single markdown file with YAML frontmatter under `.jaira/tickets/`, readable and editable by hand
- [ ] **STORE-02**: Ticket IDs are ULIDs, generated without any shared counter, so two parallel sessions never collide
- [ ] **STORE-03**: User can refer to a ticket by an unambiguous ID prefix, the way git accepts short SHAs
- [ ] **STORE-04**: Writing one field leaves every other field, comment, and blank line in the file byte-identical
- [ ] **STORE-05**: An unrecognized frontmatter block (such as the reserved `external:`) survives any write untouched
- [ ] **STORE-06**: Ticket writes are atomic — an interrupted write never leaves a partial or corrupt file
- [ ] **STORE-07**: Two concurrent CLI invocations cannot clobber each other's writes to the same ticket
- [ ] **STORE-08**: Board state persists across sessions — closing and reopening shows the same tickets in the same lanes
- [ ] **STORE-09**: A teammate who clones the repo and runs `jaira` sees the same board with no configuration step
- [ ] **STORE-10**: Frontmatter containing YAML constructs the tool cannot safely rewrite (anchors, aliases) is detected and reported rather than silently mangled

### Ticket Schema

- [ ] **SCHEMA-01**: Ticket records `title`, `status`, `goal`, `context`, `definition-of-done`, `created-at`, `updated-at`
- [ ] **SCHEMA-02**: Ticket records `creator` and `assignee`, both human identities
- [ ] **SCHEMA-03**: Ticket records `executed-by` per run, naming the model that did the work, separately from `assignee`
- [ ] **SCHEMA-04**: Ticket declares `blocked-by` as a flat list of other ticket IDs
- [ ] **SCHEMA-05**: Ticket records an outcome block with three distinct fields — what was done, why it was needed, and how it satisfies the Definition of Done
- [ ] **SCHEMA-06**: Ticket records the commit SHAs produced for it
- [ ] **SCHEMA-07**: Ticket reserves an `external:` block for a future Jira/YouTrack adapter, unused in v1

### CLI

- [ ] **CLI-01**: User can create a ticket from the command line
- [ ] **CLI-02**: User can list tickets, filtered by lane
- [ ] **CLI-03**: User can show one ticket in full
- [ ] **CLI-04**: User can move a ticket to another lane
- [ ] **CLI-05**: Every read command supports `--json` emitting a stable machine-readable schema, never mixed with human text
- [ ] **CLI-06**: Exit codes are stable and documented — success, generic error, usage error, validation failure, dependency-blocked
- [ ] **CLI-07**: Errors under `--json` are emitted as structured JSON on stderr, parseable without regex
- [ ] **CLI-08**: Re-running a command that is already satisfied succeeds as a no-op rather than erroring
- [ ] **CLI-09**: `jaira init` prepares a repo — creates `.jaira/`, `.gitattributes`, and the local gitignore for ephemeral session state
- [ ] **CLI-10**: The CLI is driven entirely through bash by any agent, requiring no MCP server or runtime

### Gates

- [ ] **GATE-01**: A ticket cannot leave Backlog until `goal`, `definition-of-done`, `context`, and `assignee` are all present
- [ ] **GATE-02**: A blocked gate reports exactly which fields are missing, so an agent can fix and retry
- [ ] **GATE-03**: A ticket whose `blocked-by` tickets are unresolved cannot be started
- [ ] **GATE-04**: Gate rules are enforced identically whether the caller is the CLI or the TUI
- [ ] **GATE-05**: Gates are enforced at the moment of mutation, not only at discovery, so an agent cannot bypass them by skipping the discovery step
- [ ] **GATE-06**: `assignee` defaults to the ticket's creator
- [ ] **GATE-07**: A ticket cannot enter Done on an agent's own say-so — the transition requires a non-model signal (passing command or explicit human sign-off)

### Board / TUI

- [ ] **TUI-01**: Board renders lanes as columns with tickets as cards
- [ ] **TUI-02**: User can navigate cards and lanes with vim keys and arrow keys
- [ ] **TUI-03**: User can expand a card to see full ticket detail inline
- [ ] **TUI-04**: User can filter tickets live by typing a query
- [ ] **TUI-05**: User can create a ticket into the Backlog from the board
- [ ] **TUI-06**: User can move a card between lanes from the board
- [ ] **TUI-07**: User can view the combined diff of a ticket's commits without leaving the board
- [ ] **TUI-08**: Board shows a keybinding reference on demand
- [ ] **TUI-09**: Board refreshes automatically when tickets change on disk — another session, a `git pull`, or a manual edit
- [ ] **TUI-10**: A background refresh preserves the user's current selection and any open detail pane
- [ ] **TUI-11**: If a ticket is changed elsewhere while the user is editing it, the conflict is surfaced rather than silently overwritten
- [ ] **TUI-12**: User can switch between projects from within the board
- [ ] **TUI-13**: Board starts instantly — no daemon, no index build, no network

### Lanes

- [ ] **LANE-01**: Seven base lanes ship inside the binary and need no repo configuration: Backlog, Todo, In Progress, Human, Review, Done, Blocked
- [ ] **LANE-02**: A custom lane is a single portable markdown file in `~/.jaira/lanes/`, shareable by sending the file
- [ ] **LANE-03**: A lane definition specifies its prompt, its model tier, and its input/output contract
- [ ] **LANE-04**: A lane's model tier is a local alias, not a hardcoded model name, so shared lane files survive model renames
- [ ] **LANE-05**: Custom lanes order themselves by anchoring to an existing lane rather than a numeric position
- [ ] **LANE-06**: A lane whose `id` collides with a base lane is rejected with a clear error, never silently overriding it
- [ ] **LANE-07**: A ticket in a lane the local install does not have renders as a read-only passthrough column, never hidden
- [ ] **LANE-08**: Advancing a ticket into or out of an unrecognized lane is refused, since no contract exists to enforce

### Agent Pipeline

- [ ] **PIPE-01**: Agent can query the next actionable ticket, excluding blocked and already-claimed work
- [ ] **PIPE-02**: Agent can take a time-limited claim on a ticket, and a stale claim is treated as abandoned without manual unlocking
- [ ] **PIPE-03**: Agent receives a bounded input assembled by the tool — the lane's required ticket fields, the diff of the ticket's own commits, and the lane prompt — never assembled by the agent itself
- [ ] **PIPE-04**: Agent advances a ticket by supplying structured output, which the tool validates against the lane's declared output schema and rejects with a readable reason if malformed
- [ ] **PIPE-05**: A ticket cannot advance unless the current lane's declared outputs were actually supplied
- [ ] **PIPE-06**: Agent can move a ticket to the Human lane with a question attached
- [ ] **PIPE-07**: Claude works the board in-session, spawning one subagent per lane step at that lane's model tier
- [ ] **PIPE-08**: A skill teaches Claude the CLI surface and the ticket schema

### Conflict Handling

- [ ] **MERGE-01**: Two people moving the same ticket to different lanes resolves automatically rather than producing a conflict
- [ ] **MERGE-02**: Concurrent additions to `blocked-by` or `commits` are both preserved
- [ ] **MERGE-03**: Competing scalar edits resolve deterministically and never silently revert forward lane progress
- [ ] **MERGE-04**: Genuinely conflicting prose edits surface as a conflict scoped to the affected field, not the whole file
- [ ] **MERGE-05**: The merge driver installs itself on first run and announces that it did so, rather than acting silently
- [ ] **MERGE-06**: A merge never produces invalid YAML or loses a ticket

### Session Context

- [ ] **SESSION-01**: Claude records the current focus and reasoning via the CLI at checkpoints
- [ ] **SESSION-02**: Board displays the current session focus as a live panel
- [ ] **SESSION-03**: Multiple concurrent sessions in the same working tree each appear, with stale ones marked rather than deleted
- [ ] **SESSION-04**: Session state is never committed to git

### Task Tool Sync

- [ ] **SYNC-01**: Claude's structured tasks are mirrored into the Backlog as tickets with `ready: false`
- [ ] **SYNC-02**: Mirrored tickets are matched by a stable ID carried in task metadata, so repeated syncs update rather than duplicate
- [ ] **SYNC-03**: A task disappearing from Claude's list never deletes its mirrored ticket
- [ ] **SYNC-04**: Syncing when nothing changed writes nothing and produces no git diff
- [ ] **SYNC-05**: Starting a ticket surfaces it in the in-session task list, so work begun in chat stays visible in both places
- [ ] **SYNC-06**: Task dependencies and jAIra's `blocked-by` stay consistent across a sync
- [ ] **SYNC-07**: The sync adapter is isolated in one module, so a change to Claude's task API is a contained patch
- [ ] **SYNC-08**: A sync failure never breaks the user's Claude session

### Distribution

- [ ] **DIST-01**: The tool is a single static binary with no runtime dependency
- [ ] **DIST-02**: Binaries are produced for Linux, macOS, and Windows on amd64 and arm64
- [ ] **DIST-03**: A teammate can install and run it without a toolchain

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

Populated during roadmap creation.

| Requirement | Phase | Status |
|-------------|-------|--------|
| — | — | Pending |

**Coverage:**
- v1 requirements: 84 total
- Mapped to phases: 0
- Unmapped: 84 ⚠️

---
*Requirements defined: 2026-08-11*
*Last updated: 2026-08-11 after initial definition*
