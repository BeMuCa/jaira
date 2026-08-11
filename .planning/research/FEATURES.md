# Feature Research

**Domain:** Terminal kanban board for AI-coding-agent tasks, stored as git-committed markdown tickets
**Researched:** 2026-08-11
**Confidence:** MEDIUM-HIGH (named tools verified against GitHub/official docs/multiple independent write-ups; a few mechanism-level details — exact conflict behavior of `tissue`, `git-issue`, taskwarrior's per-field merge — are single-sourced and flagged LOW where noted)

## Survey Notes (tools actually verified to exist, 2026)

**1. File-based / git-native issue trackers**
- **git-bug** (git-bug/git-bug, active) — issues/comments stored as git objects (not working-tree files); orders operations via git DAG topology + Lamport clocks + hash tie-break, a CRDT-style model. Conflict-free by construction, but you lose the "hand-editable plain file" property — there is no plain-text ticket file to `cat` or diff normally.
- **Fossil** (fossil-scm.org, active, predates git) — tickets are "artifacts" in a append-only "bag of artifacts"; proven to be a G-Set CRDT (grow-only set). Amending a ticket = appending a new artifact with the same ticket ID; "tickets do not branch." This is the most mature real-world proof that append-only + last-values-win-per-field beats line-based text merging for ticket state.
- **Bugs Everywhere** (dead) — stored issues as files in the working tree, edited in place → suffered exactly the collision problem jAIra is worried about; cited by practitioners as *why* this naive approach failed and was abandoned.
- **ticgit** (dead, ~2008, Scott Chacon) — issues as files in an orphan git branch; abandoned once its author moved on to build GitHub. No modern maintenance.
- **tissue** (smeans/tissue and git.systemreboot.net/tissue — two unrelated same-named projects exist; the systemreboot one is a Guile-based minimalist git+plain-text tracker with Xapian search) — issues are free-form text files committed to git. Conflict handling not documented in sources found (LOW confidence on its merge story specifically) — its "minimalist" positioning suggests conflicts are left to git's default text merge, i.e. unresolved on collision.
- **git-issue** (dspinellis/git-issue, active) — issues as plain text files, one per issue directory, changed/shared through git. Decentralized, no server. Merge behavior not verified beyond "plain git merge of files" (LOW confidence).
- **todo.txt** (still used, primarily single-user) — single flat text file, one task per line, plain-text metadata (`+project`, `@context`, priority). Because every task is a full line, git's line-based merge mostly works — collisions only occur if two people edit the *same task line*, which is rare with one-line-per-task granularity. This is the closest structural analog to jAIra's "one field, one line" YAML approach, just simpler (no field structure).
- **Backlog.md** (MrLesk/Backlog.md, active, 2025-2026, directly comparable) — markdown-native, one file per task, explicitly pitched as "collaboration between humans and AI agents in a git ecosystem," with a terminal kanban view *and* an optional web view, MCP integration for Claude Code/Gemini CLI. This is the single closest existing competitor to jAIra's storage model.
- **Beads / bd** (steveyegge/beads and forks, released Oct 2025, very active) — the most relevant AI-agent-focused entrant. Local SQLite is the working store; a background daemon exports every change as a new line to a committed `issues.jsonl` (append-only log, not in-place edits). Hash-based IDs prevent create-collisions between parallel agents. Git may still flag a textual conflict on the JSONL file, but resolution is described as trivial ("keep both lines") because each line is an independent, self-contained change record rather than a shared mutable field. This is the concrete, working answer to the "two agents touch the same ticket" problem — but it requires giving up "the file IS the ticket" (the file is an event log; current ticket state is a derived/materialized view, not something you can hand-edit and trust is current).
- **git-native-issue** (remenoscodes, new 2026 entrant) — explicitly positions itself as fixing what killed Bugs Everywhere by publishing a standalone format spec; too new for adoption evidence (LOW confidence on real-world track record).

**Verdict for jAIra's core tension:** Every tool that succeeded at avoiding conflicts (Fossil, git-bug, Beads) did so by moving away from "the ticket is a single mutable file you edit in place" toward an append-only event log or CRDT object store. Every tool that kept the naive "edit the file's fields in place" model (Bugs Everywhere) failed or never solved the conflict problem. jAIra's stated constraint — hand-editable, diff-readable YAML frontmatter *is* the API — puts it structurally closer to the failure mode (Bugs Everywhere / todo.txt) than to the tools that solved it. The mitigation todo.txt uses (one line = one full unit) is the only pattern that preserves plain-text-editability *and* reduces (not eliminates) collisions: keeping each frontmatter field on its own line already gets jAIra 90% of what todo.txt gets, because a collision only occurs when two people change the *same field* on the *same ticket* in the *same window* — status is the one field with meaningfully higher collision odds since everyone touches it. A custom git merge driver (declared in `.gitattributes`) that applies a small resolution policy specifically to the `status:` line (e.g., defer to whichever side is later in a documented lane order, or just always flag for `jaira` CLI reconciliation rather than trusting git's raw conflict markers) is a smaller, cheaper move than adopting a CRDT engine, and matches the project's own scope discipline. This belongs in PITFALLS.md for depth; noted here because it changes which "conflict handling" feature is table stakes vs. which is overreach (see Anti-Features).

**2. Terminal kanban / TUI task boards**
- **taskwarrior** + **taskwarrior-tui** (kdheepak, active) / **vit** (active, vim-keybound) — not kanban-native (report/filter-based, not column-based), but they set the *interaction* baseline every terminal task tool inherits: vim-style navigation, `/` to filter live, context-sensitive keybinding help via `?`, a details/annotation pane, and multiple saved "reports" (closest analog to saved views). Sync is via an optional Taskserver (`taskd`); conflict resolution is per-task (keyed by UUID) with last-modified-wins on conflicting fields — no field-level merge, whole-task overwrite by recency. Confirms whole-record last-write-wins is the "boring but works" fallback when line/field merging isn't attempted.
- **taskell** (smallhadroncollider/taskell, active) — genuine CLI Kanban board, vim-keybound by default, stores tasks as Markdown with "clean diffs for version control," supports Trello/GitHub Projects import. Structurally the closest non-AI analog to jAIra's TUI: markdown storage + kanban columns + vim keys.
- **kanban.nvim** / **quick-kanban.nvim** — editor-embedded (Neovim) file-based kanban, markdown-backed, Obsidian-kanban-file compatible. Confirms markdown-as-kanban-source-of-truth is an established pattern, but scoped to a single user's editor session, not a shared team artifact.
- **kanban-python** / **daikanban** (PyPI) — smaller/less-maintained terminal kanban tools; existence confirmed via PyPI listings, feature depth not verified beyond package description (LOW confidence, likely low-adoption).
- **General TUI convention benchmark** (lazygit, gitui — not kanban but the dominant modern terminal-tool UX reference): single-key panel actions, `/` for fuzzy filter, vim h/j/k/l plus arrow-key fallback, a persistent keybinding cheat-sheet (`?`), panel-scoped meaning for the same key (context-sensitive commands), and a strong split between "fast tool people leave open all day" (gitui, Rust, startup-speed-focused) vs "deepest feature set" (lazygit). jAIra's "instant startup" constraint puts it in the gitui camp, not the lazygit camp — startup latency competes with feature depth, and users of this class of tool consistently pick speed when forced to choose.

**3. AI-agent task orchestration and tracking tools**
- **Claude Code TodoWrite → Task tools migration (VERIFIED, HIGH confidence, official docs)** — as of Claude Code v2.1.142 / Agent SDK 0.3.142, `TodoWrite` is **no longer the default**; Claude Code now defaults to structured `TaskCreate` / `TaskUpdate` / `TaskGet` / `TaskList` tools. Old shape: one call rewrites a full `todos` array (`content`, `status`, `activeForm`). New shape: per-item calls keyed by a real `taskId`, with `status` (`pending`/`in_progress`/`completed`/`deleted`), `owner`, `addBlocks`, `addBlockedBy`, and a free-form `metadata` field. **This directly threatens the PROJECT.md requirement "Claude's TodoWrite items are mirrored into the Backlog"** — Anthropic has already built a task system with IDs, status, ownership, and dependency fields that is structurally converging with what jAIra wants to layer on top. jAIra's mirror/sync layer must target the current Task tools, not just legacy TodoWrite, and should expect the target API to keep moving. This is the single most important finding for scoping the "Capture and promotion gate" feature set.
- **claude-task-master / task-master-ai** (eyaltoledano/claude-task-master, active) — parses a PRD into tasks, drops into Cursor/Windsurf/Roo/Claude Code, no server, multiple tool-count tiers (7/15/etc tools). Task decomposition + dependency tracking, but no kanban visualization and no per-lane model-tier binding — it is a task *generator*, not a board.
- **Beads (bd)** (also fits here) — explicitly built to solve agent session amnesia ("50 First Dates" problem): persistent structured memory across sessions, `bd ready` computes unblocked work automatically, four dependency types (blocks/related/parent-child/discovered-from). No kanban UI, no lane-as-pipeline-step concept, no human/model attribution split.
- **BMAD-METHOD** (community fork examples verified via GitHub; original by "Brian Madison") — agents defined as single portable markdown files with embedded YAML, installed per-project, invoked via slash-commands in Claude Code/Cursor/Roo/Cline. Confirms "a workflow step as one portable file" is a proven, adopted pattern (directly analogous to jAIra's `~/.jaira/lanes/*.md`) — but BMAD's files define *agent personas*, not board lanes with an I/O contract; there is no ticket/kanban artifact, no persisted state outside the conversation and generated docs.
- **Roo Code** (active) — multi-mode system (Code/Architect/Ask/Debug) plus "Boomerang Tasks," where an orchestrator mode delegates a subtask to a specialized mode and the result returns ("boomerangs") to the parent task. Assigns different models per mode to cut cost — the closest existing precedent for jAIra's "cheap model implements, strong model reviews" lane-tier idea, but it lives entirely inside a single VS Code session; nothing persists as a shared, git-committed artifact.
- **Cline** (active) — simpler Plan/Act two-phase workflow, no multi-mode tree. Confirms two competing philosophies exist (Cline: simple two-phase; Roo: many specialized modes) — evidence that jAIra's "seven fixed base lanes, no repo config" sits deliberately closer to Cline's simplicity than Roo's configurability, which matches the project's own "smaller than paca" discipline.
- **Agent OS** (buildermethods/agent-os, active, v3 released 2026) — not a task tracker; it injects codebase standards into an agent's context and now explicitly defers task/spec planning to the coding tool's own Plan Mode. Relevant mainly as a negative data point: even a well-adopted "AI dev process" framework decided *not* to build its own task/ticket layer, choosing instead to integrate with what Claude Code/Cursor already provide — a caution against jAIra over-building capture mechanics that duplicate what the host agent already does.
- **GitHub Copilot Workspace / Copilot coding agent** (active, "Mission Control" dashboard + repo-level "Agents" tab + CLI `/fleet` command, all late-2025/early-2026 features per Microsoft/GitHub sources) — the clearest evidence of real user demand for *visualizing multiple concurrent agent tasks*: a dashboard for assigning/steering/tracking concurrent agent runs, and `gh agent-task list` for listing them from a terminal. This validates jAIra's session-context panel bet — teams do want to see what agents are doing, not just what tickets exist — but GitHub's version is server/hosted; jAIra's differentiator is doing the equivalent locally, from git state alone, no server.
- **Human-in-the-loop pattern (general, MEDIUM confidence — no tool-specific canonical doc found, cross-referenced across several agent-framework write-ups)** — the common shape is: agent pauses execution, surfaces a specific question/decision, execution resumes only after a human responds; several frameworks implement this as an explicit "approval gate" state. jAIra's "ticket lands in the Human lane with the question attached" is a direct, board-native implementation of this same pattern, made visible as a kanban column instead of a modal prompt or webhook callback — no surveyed tool renders this as a persistent, glanceable board column.

## Feature Landscape

### Table Stakes (Users Expect These)

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| Plain-text/markdown storage that produces clean git diffs | Every surveyed tool in category 1 and taskell/Backlog.md converge here; a ticket store that doesn't diff cleanly breaks the "git is the sync layer" pitch on day one | LOW | Already the project's chosen format; risk is in *how* fields are laid out, not *whether* markdown is used |
| Kanban columns with keyboard-only navigation (vim-style: j/k/arrows, enter to open) | Universal across taskell, taskwarrior-tui, vit, lazygit/gitui; a TUI tool without vim-consistent keys feels foreign to its own audience | LOW-MEDIUM | Should also support arrow keys as fallback per taskwarrior-tui/lazygit convention |
| Inline/pane card detail view | taskell, Backlog.md, and every git TUI (lazygit/gitui panel model) show full detail without leaving the list | LOW | PROJECT.md already specifies this ("pressing a key expands full ticket detail inline") |
| Live filter/search over tickets | `/`-filter is present in taskwarrior-tui, vit, lazygit, gitui; absence makes any board with >20 cards unusable | LOW-MEDIUM | Not yet in PROJECT.md's Active requirements — flag as a likely missing table-stakes item |
| Instant local startup, no daemon required to view the board | gitui explicitly wins adoption over lazygit on raw startup speed for large repos; PROJECT.md already names this as a hard constraint | MEDIUM | Binary distribution + no runtime already chosen; the risk is scope creep in what loads at startup (e.g., don't compute full commit diffs eagerly) |
| Zero-config first run on a fresh clone | Backlog.md and Beads both market "zero setup" (`bd init`, `backlog init`) as a selling point; git-bug's steep object-model learning curve is cited as an adoption barrier by contrast | MEDIUM | Seven built-in base lanes already satisfies this; must ensure first `jaira` run needs no prompts |
| Task dependency declaration (`blocked-by` equivalent) and enforcement | taskwarrior `depends:`, Beads' four dependency types, git-bug issue links — every serious tracker has *some* dependency notion | MEDIUM | PROJECT.md already scopes this narrowly (single `blocked-by` list) — correct scope, see Anti-Features on over-building this |
| Conflict-tolerant sync that never silently loses a ticket | Table stakes given the storage choice — a tool that corrupts or drops tickets on merge conflict fails its core promise faster than any missing feature would | HIGH | Not fully solved by any surveyed plain-file tool (Bugs Everywhere failed here); Fossil/git-bug/Beads solve it by *not* being plain-editable files. See Survey Notes verdict — this is the highest-complexity table-stakes item and needs its own design pass, not a default "just use git merge" assumption |
| Project/board identity tied to the repo, not a separate login | Every category-1 and category-3 tool surveyed is repo-scoped, no accounts | LOW | Comes free from "git is the sync layer" |

### Differentiators (Competitive Advantage)

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| Lanes bind a prompt + model tier + I/O contract (board as a controlled multi-model pipeline) | No surveyed tool combines a kanban lane with a bound prompt/model/I-O contract — Roo Code's per-mode model assignment is the closest precedent, but it's session-local, not a persisted, git-shared lane definition | HIGH | Depends on: ticket store, CLI single-write-path, Claude's ability to spawn a subagent per lane at that lane's tier |
| Promotion gate requiring DoD/goal/context/assignee before leaving Backlog | Every surveyed capture mechanism (TodoWrite/Task tools, Beads' `bd create`, task-master's PRD parse) lets an agent start work on an under-specified item; none gate on completeness before execution begins | MEDIUM | Depends on: ticket schema fields existing; enhances review ergonomics downstream |
| Human-attributed ownership with `executed-by` recorded separately from `assignee` | No surveyed AI-agent tool separates "who owns this outcome" from "which model ran it" — Task tools have an `owner` field but no distinct executor-model field; Beads/task-master don't track authorship at all | LOW-MEDIUM | Straightforward schema addition; the value is in the *convention* (always human-owned) more than the mechanism |
| Session-context panel showing agent focus/reasoning at checkpoints | GitHub's Mission Control/Agents tab proves demand for visualizing what an agent is doing, but that's server-hosted; nothing surveyed does this from local git state alone | MEDIUM-HIGH | Depends on Claude writing structured checkpoints via the CLI; quality of the panel is bounded by how disciplined those checkpoint writes are |
| Two-way sync between the ephemeral in-session task list and durable tickets | Distinguishes jAIra from Beads/task-master, which require an agent to explicitly call *their* tool instead of the one Anthropic ships by default | HIGH | **Risk, not just complexity**: Anthropic replaced TodoWrite with Task tools in Claude Code v2.1.142 — the sync target has already moved once and will likely move again; build against the currently-documented Task tools (`TaskCreate`/`TaskUpdate`/`TaskList`/`TaskGet`), not the deprecated `TodoWrite` shape, and isolate the adapter so the next API change is a small patch |
| Unknown/custom lanes degrade to read-only passthrough columns instead of hiding tickets | No surveyed tool handles "teammate has a lane config you don't" gracefully — Backlog.md and taskell configs are per-user with no fallback story; this is a genuine, simple UX win | LOW | Enhances: portable custom lane files |
| Portable single-file custom lanes (share by sending a file) | Directly analogous to BMAD's "agent as one markdown file," proven adopted pattern — applying it to *board lanes* rather than agent personas is the novel part | LOW-MEDIUM | Depends on lanes-as-agent-steps existing first |
| Definition-of-done-anchored three-field outcome block (what/why/resolves) | Beads and Backlog.md tickets close with free-text description only; no surveyed tool structurally forces a causal link back to the DoD at close time | LOW | Pure schema/CLI convention, cheap to build, high review-quality payoff |
| Board-native "Human" lane for agent questions (vs. a modal prompt or webhook) | Human-in-the-loop is universally implemented as a pause-and-prompt in agent frameworks; rendering it as a persistent, glanceable kanban column is not seen elsewhere in the survey | LOW-MEDIUM | Depends on ticket store + lane model; low complexity because it reuses the existing lane mechanism, just with a fixed built-in lane |

### Anti-Features (Commonly Requested, Often Problematic)

PROJECT.md already excludes these; each is validated or challenged against survey evidence.

| Feature | Why Requested | Why Problematic | Evidence / Alternative |
|---------|---------------|------------------|-------------|
| Sprints | Feels natural coming from Jira-style tools | No surveyed lightweight/file-based or AI-agent tool (taskwarrior, Beads, Backlog.md, task-master) has sprints; it's a full-PM-tool feature, not a tracker feature | **Validated drop.** Nothing in the adjacent-tool survey supports adding this |
| Custom fields | Teams always want "just one more field" | Every surveyed lightweight tool that allows this (Beads' `metadata`, Task tools' `metadata`) treats it as an escape hatch, not a first-class feature — full custom-field UIs belong to paca-class tools | **Validated drop, already satisfied.** PROJECT.md's reserved `external:` block is exactly this escape hatch; no need for a second one |
| Roles/permissions | Teams assume "who can do what" needs configuring | No file-based tracker in the survey has roles — git repo access *is* the permission system; adding roles duplicates git | **Validated drop.** |
| Saved views | Users of taskwarrior-tui/vit lean on named reports/filters daily | Full "views" (Jira-style: shared, configurable, dashboard-embeddable) are heavy; but *personal filter shortcuts* (taskwarrior's saved reports) are a genuinely-used lighter version of the same idea | **Partially challenge.** Drop full saved/shared views for v1 as scoped; treat "save a filter as a named shortcut" as a candidate for v1.x, not v2 — it's structurally just a stored search string, far short of a "view" with columns/fields/sharing |
| Dashboards | Managers want an aggregate status view | No TUI tool surveyed has this — dashboards are inherently a web/server feature (this is what GitHub's Mission Control is) | **Validated drop.** Consistent with "no server" constraint; if ever needed, it's a v2+ hosted add-on, not core |
| Server / accounts / authentication | Assumed necessary for "real" multi-user tools | git-bug, Fossil, Beads, Backlog.md, taskwarrior (core) all operate serverless; Taskserver/taskd exists only as an *optional* sync convenience, never required | **Strongly validated drop.** The entire adjacent field of git-native trackers proves this works without a server |
| Jira/YouTrack sync now | Teams already have an existing tracker | Auth, field-mapping, and conflict resolution across two systems is its own project — no surveyed lightweight tool attempts bidirectional sync with a hosted tracker at v1 | **Validated defer.** The `external:` schema reservation is a good, low-cost hedge; no tool surveyed does better than "reserve space, build later" |
| Web UI / VS Code extension | Backlog.md ships both a terminal and a web view | Splitting effort across UI surfaces is exactly the "bigger than it needs to be" trap; GitHub's hosted dashboard shows the web-view itch is real, but that itch is served by git + terminal for this project's stated audience (a teammate who clones and runs `jaira`) | **Validated drop for v1.** Revisit only if usage data shows people want to check the board without a terminal open — no evidence of that yet |
| Branch-per-ticket | Feels safer, mirrors PR-per-feature convention | Beads, git-bug, Backlog.md, taskwarrior all operate ticket-state independent of branching — coupling ticket lifecycle to branch lifecycle adds merge overhead exactly where the project is trying to remove it, and breaks down further under multiple parallel agents (which PROJECT.md explicitly anticipates) | **Validated drop.** |
| Board-spawned headless `claude -p` background processes | GitHub's `/fleet` and Beads' background daemon prove the opposite pattern (tool-spawned agents) is viable elsewhere | Managing background process lifecycle (crash recovery, orphaned processes, log capture) is a substantial engineering surface that Roo Code and Cline both avoid by keeping orchestration inside the live assistant session — matches jAIra's "queue plus a view" framing | **Validated drop for v1,** with the caveat that it's the one item on this list where a competitor category (hosted multi-agent dashboards) is actively moving the other way — worth re-examining only if users want unattended overnight runs |
| TodoWrite as the backbone of the whole system | It's already there, seems like "just build on top of it" | Beads' own comparison argument (ephemeral session list vs. structured persistent tracker) independently reaches the same conclusion jAIra did; reinforced by the fact Anthropic itself just replaced TodoWrite with a more structured, ID-based Task system, i.e. even Anthropic agrees a flat ephemeral list isn't enough | **Strongly validated drop**, and the promotion gate is the right answer — but see the Task-tools migration risk noted above; the on-ramp must target the current API |
| Full multi-type dependency graph (blocks/related/parent-child/discovered-from, graph queries like Beads' `bd dep tree`) | Beads makes a compelling case that richer dependency typing helps agents plan | A single `blocked-by` list plus "cannot start until satisfied" covers the actual stated need (don't start work prematurely); building graph-query tooling is solving a problem (multi-agent long-horizon planning across hundreds of issues) this project doesn't have at its stated scale | **New anti-feature to add, not in original list.** Keep `blocked-by` as a flat list; do not build dependency-type taxonomy or graph traversal commands |
| CRDT-based conflict-free ticket engine (git-bug/Fossil-style object store) | The merge-conflict problem is explicitly named as critical in scope | Solves conflicts completely, but requires abandoning "the markdown file with YAML frontmatter IS the ticket, hand-editable" — every tool that did this (git-bug, Fossil) stopped being a plain-text-editable format, which is a harder constraint violation than tolerating occasional conflicts | **New anti-feature to add.** Prefer a narrower fix (line-per-field YAML + a custom `.gitattributes` merge driver or a `jaira reconcile` command for the `status:` line specifically) over adopting a CRDT/object-store engine |

## Feature Dependencies

```
Ticket store (markdown + YAML frontmatter, .jaira/tickets/)
    └──requires──> Conflict-tolerant status field
                       (line-per-field layout + merge policy for `status:`)

Promotion gate (goal/DoD/context/assignee required)
    └──requires──> Ticket store
    └──enhances──> Outcome block review ergonomics (what/why/resolves)

CLI single write path
    └──requires──> Ticket store
    └──enables──> Lanes-as-agent-steps (prompt + model tier + I/O contract)
    └──enables──> TodoWrite/Task-tool two-way sync

Lanes-as-agent-steps
    └──requires──> CLI single write path
    └──enables──> Portable custom lane files (~/.jaira/lanes/*.md)
    └──enables──> Session-context panel (visualizes which lane/agent is active)

Portable custom lane files
    └──enhances──> Unknown-lane passthrough columns
                       (passthrough is the fallback when a shared lane file is missing)

TodoWrite/Task-tool two-way sync
    └──requires──> CLI single write path
    └──requires──> Promotion gate (defines what a "promoted" ticket looks like)
    └──depends on external API──> Claude Code Task tools (TaskCreate/TaskUpdate/TaskList/TaskGet),
                                    NOT the deprecated TodoWrite shape

Human-attributed ownership (assignee vs executed-by)
    └──requires──> Ticket store fields
    └──enhances──> Promotion gate (assignee is itself a gate field)

Dependencies (blocked-by)
    └──requires──> Ticket store
    └──enables──> "Cannot start until satisfied" enforcement
    └──conflicts with──> Full multi-type dependency graph (deliberately not built)

Board-native Human lane
    └──requires──> Lanes-as-agent-steps (reuses the same lane mechanism)
    └──enhances──> Human-in-the-loop visibility (vs. a modal prompt)

Filter/search over tickets
    └──requires──> Ticket store
    └──enhances──> Every other view (detail, project switching, board)
```

### Dependency Notes

- **Conflict-tolerant status field must land before or alongside the ticket store**, not after — every failed tool in the survey (Bugs Everywhere) treated conflict handling as an afterthought once the naive file-edit model was already load-bearing; retrofitting it later means a data-format migration.
- **TodoWrite/Task-tool sync depends on an external, moving API.** Build the adapter as an isolated module so the next Anthropic change (there has already been one, TodoWrite → Task tools in v2.1.142) is a contained patch, not a redesign.
- **Lanes-as-agent-steps is the load-bearing differentiator** — the session-context panel, portable lane files, and passthrough-column fallback all enhance it but none of them make sense before it exists. Sequence lanes-as-agent-steps early relative to the other differentiators.
- **Dependencies (blocked-by) and the full dependency graph anti-feature are mutually exclusive by design** — building the simple version is precisely what avoids needing the complex one; don't let the survey's dependency-graph evidence (Beads) pull scope toward the richer version.

## MVP Definition

### Launch With (v1)

- [ ] Markdown + YAML frontmatter ticket store, one file per ticket — the format is the API, nothing works without it
- [ ] Seven built-in base lanes rendered as a TUI kanban board, vim + arrow-key navigation — table stakes for any TUI board
- [ ] Inline card detail view on keypress — table stakes, low complexity
- [ ] Live filter/search across tickets — table stakes, currently missing from PROJECT.md's Active list, should be added
- [ ] Promotion gate (goal/DoD/context/assignee required to leave Backlog) — the primary differentiator that's cheapest to build first
- [ ] `blocked-by` dependency declaration + start-blocking enforcement — table stakes, keep scope to a flat list
- [ ] CLI as single write path — everything else depends on this existing first
- [ ] One-way TodoWrite/Task-tool → Backlog mirror, targeting current Task tools — de-risk the two-way sync by shipping the simpler direction first
- [ ] Human-attributed ownership fields (`assignee`, `executed-by`) — low complexity, high differentiation value, no reason to defer
- [ ] Project switching within the TUI — table stakes given multi-repo usage is explicitly anticipated

### Add After Validation (v1.x)

- [ ] Two-way TodoWrite/Task-tool sync — add once the one-way mirror is proven and the Task-tools API has settled further
- [ ] Lanes-as-agent-steps with bound prompt/model tier/I-O contract — the headline differentiator, but high complexity; validate the simpler board first, then layer the pipeline
- [ ] Portable custom lane files + passthrough-column fallback — depends on lanes-as-agent-steps existing
- [ ] Session-context panel — depends on Claude reliably writing structured checkpoints; needs the CLI/skill to be stable first
- [ ] Named filter shortcuts (lightweight "saved view" alternative) — trigger: users start retyping the same filter string repeatedly

### Future Consideration (v2+)

- [ ] Jira/YouTrack adapter using the reserved `external:` block — defer until a real integration partner/need exists
- [ ] Custom `.gitattributes` merge driver or `jaira reconcile` command for `status:` line conflicts, if real-world usage shows this is a frequent pain point — build reactively, not speculatively
- [ ] Board-spawned headless agent runs — only if unattended/overnight orchestration becomes a stated need; currently explicitly out of scope

## Feature Prioritization Matrix

| Feature | User Value | Implementation Cost | Priority |
|---------|------------|---------------------|----------|
| Markdown ticket store | HIGH | LOW | P1 |
| Kanban TUI with keyboard nav + detail view | HIGH | MEDIUM | P1 |
| Filter/search | HIGH | LOW | P1 |
| Promotion gate | HIGH | MEDIUM | P1 |
| blocked-by dependencies | MEDIUM | MEDIUM | P1 |
| CLI single write path | HIGH | MEDIUM | P1 |
| One-way TodoWrite/Task-tool mirror | HIGH | MEDIUM | P1 |
| Human/executed-by attribution | MEDIUM | LOW | P1 |
| Project switching | MEDIUM | LOW | P1 |
| Conflict-tolerant status field policy | HIGH | HIGH | P1 (design must land in v1, even if implementation is a simple policy) |
| Two-way TodoWrite/Task-tool sync | HIGH | HIGH | P2 |
| Lanes-as-agent-steps (prompt+tier+I/O contract) | HIGH | HIGH | P2 |
| Portable custom lane files + passthrough fallback | MEDIUM | LOW-MEDIUM | P2 |
| Session-context panel | MEDIUM | MEDIUM-HIGH | P2 |
| Named filter shortcuts | LOW-MEDIUM | LOW | P3 |
| Jira/YouTrack adapter | LOW (now) | HIGH | P3 |
| Merge driver / reconcile command | MEDIUM (contingent) | MEDIUM | P3 |

**Priority key:**
- P1: Must have for launch
- P2: Should have, add when possible
- P3: Nice to have, future consideration

## Competitor Feature Analysis

| Feature | Backlog.md | Beads (bd) | jAIra's Approach |
|---------|--------------|--------------|--------------|
| Storage format | One markdown file per task, human-editable | SQLite (working store) + append-only JSONL export (git-committed source of truth) | One markdown file per ticket with YAML frontmatter, human-editable — closer to Backlog.md's model than Beads' |
| Conflict handling | Not documented as solved beyond plain-file git merge (unverified — LOW confidence) | Append-only log makes conflicts "keep both lines" — trivial resolution, but the file is a log, not the current-state view | Line-per-field YAML reduces (not eliminates) collision surface; status-field-specific merge policy recommended, not yet a solved problem |
| AI agent integration | MCP integration for Claude Code/Gemini CLI, "AI-ready" positioning | Built specifically for agent memory/amnesia; CLI-first, agent-callable | CLI + skill teaching Claude the schema; also targets Claude's own Task tools directly, not just a generic MCP hook |
| Kanban visualization | Terminal kanban + optional web view | None (issue list + dependency tree, no board) | Terminal kanban only (by design, no web UI) |
| Dependency tracking | Basic (unverified depth — LOW confidence) | Four relationship types + `bd ready` unblocked-work query | Single flat `blocked-by` list — deliberately simpler than Beads |
| Model-tier / pipeline binding per stage | None | None | Lanes bind prompt + model tier + I/O contract — no surveyed competitor has this |
| Human vs. model attribution | Not documented | Not documented | Explicit `assignee` (human) vs. `executed-by` (model) split — no surveyed competitor has this |
| Promotion/quality gate before work starts | None found | None (issues can be created and picked up immediately) | DoD/goal/context/assignee required before leaving Backlog — no surveyed competitor has this |

## Sources

- [git-bug (GitHub)](https://github.com/git-bug/git-bug) — CRDT/Lamport-clock ordering, active maintenance (MEDIUM-HIGH confidence, official repo + HN discussion)
- [Fossil: Bug-Tracking In Fossil](https://fossil-scm.org/home/doc/tip/www/bugtheory.wiki) — artifact/G-Set CRDT model, ticket amendment mechanics (HIGH confidence, official docs)
- [ticgit.dev](https://ticgit.dev/) and [HN: Show git-bug](https://news.ycombinator.com/item?id=17782121) — history of dead file-based trackers (MEDIUM confidence, community sourced)
- [Backlog.md (MrLesk/Backlog.md)](https://github.com/MrLesk/Backlog.md/) — markdown-native AI-agent-oriented task manager (HIGH confidence, official repo)
- [Beads issue tracker guide — Better Stack](https://betterstack.com/community/guides/ai/beads-issue-tracker-ai-agents/) and [Dicklesworthstone/beads_rust](https://github.com/Dicklesworthstone/beads_rust) and [gastownhall/beads FAQ](https://github.com/steveyegge/beads/blob/main/docs/FAQ.md) (MEDIUM-HIGH confidence, multiple independent sources agree)
- [tissue (git.systemreboot.net)](https://git.systemreboot.net/tissue/) and [smeans/tissue](https://github.com/smeans/tissue) and [dspinellis/git-issue](https://github.com/dspinellis/git-issue) (LOW-MEDIUM confidence on conflict-handling specifics, not documented in sources found)
- [taskwarrior-tui](https://github.com/kdheepak/taskwarrior-tui) and [task-sync(5) man page](https://taskwarrior.org/docs/man/task-sync.5/) — sync/conflict model, UUID + last-modified-wins (HIGH confidence, official docs)
- [taskell (smallhadroncollider/taskell)](https://github.com/smallhadroncollider/taskell) — markdown-backed CLI kanban (HIGH confidence, official repo)
- [lazygit](https://github.com/jesseduffield/lazygit) / [gitui](https://github.com/gitui-org/gitui) comparisons (MEDIUM confidence, multiple blog write-ups agree)
- [Claude Code Todo Lists docs — code.claude.com](https://code.claude.com/docs/en/agent-sdk/todo-tracking) — **HIGH confidence, official docs, directly fetched**: confirms TodoWrite → Task tools (TaskCreate/TaskUpdate/TaskGet/TaskList) migration as of Claude Code v2.1.142 / Agent SDK 0.3.142
- [eyaltoledano/claude-task-master](https://github.com/eyaltoledano/claude-task-master) (HIGH confidence, official repo)
- [buildermethods/agent-os](https://github.com/buildermethods/agent-os) and [Agent OS v3 discussion](https://github.com/buildermethods/agent-os/discussions/310) (HIGH confidence, official repo + maintainer discussion)
- Roo Code / Cline comparisons via [baeseokjae.github.io Cline vs Roo Code 2026](https://baeseokjae.github.io/posts/cline-vs-roo-code-2026/) and [Xebia: Multi Agent Workflow With Roo Code](https://xebia.com/blog/multi-agent-workflow-with-roo-code/) (MEDIUM confidence, secondary sources, cross-checked across two independent write-ups)
- [GitHub Blog: Run multiple agents at once with /fleet](https://github.blog/ai-and-ml/github-copilot/run-multiple-agents-at-once-with-fleet-in-copilot-cli/) and [Visual Studio Magazine: GitHub Agents Tab](https://visualstudiomagazine.com/articles/2026/01/29/hands-on-new-github-agents-tab-for-repo-level-copilot-coding-agent-workflows.aspx) (MEDIUM-HIGH confidence, official blog + trade press)
- `.planning/PROJECT.md` — source of the Out of Scope list validated/challenged above

---
*Feature research for: terminal kanban board for AI-coding-agent task tracking, git-native markdown storage*
*Researched: 2026-08-11*
