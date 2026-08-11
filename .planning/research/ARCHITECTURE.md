# Architecture Research

**Domain:** single-binary TUI kanban over git-committed markdown, agent-driven
**Researched:** 2026-08-11
**Confidence:** MEDIUM-HIGH (core structure and the conflict recommendation are HIGH confidence, general git/software-architecture knowledge; Claude Code hook specifics are MEDIUM — see Sources/Gaps)

## Standard Architecture

### System Overview

```
┌──────────────────────────────────────────────────────────────────────┐
│                         jaira (one static binary)                    │
│                                                                        │
│  ┌───────────────┐   ┌───────────────┐   ┌──────────────────────┐    │
│  │      CLI       │   │      TUI       │   │  Future adapters     │    │
│  │ cmd/*.go(rs)   │   │ internal/tui   │   │ (Jira/YouTrack, v2)   │    │
│  │ parses argv,   │   │ Elm/reducer    │   │ same rule: import     │    │
│  │ formats output │   │ loop, renders  │   │ core, never re-       │    │
│  │ (--json/text)  │   │ lanes+cards    │   │ implement validation  │    │
│  └───────┬────────┘   └───────┬────────┘   └───────────┬───────────┘    │
│          │  (peers, no dependency on each other)        │              │
│          └──────────────────┬────────────────────────────┘             │
│                              ▼                                        │
│  ┌────────────────────────────────────────────────────────────────┐  │
│  │                            core                                  │  │
│  │  schema+validation │ ticket store │ lane resolver │ git access  │  │
│  │  (promotion gate + blocked-by gate live HERE, called by both)   │  │
│  │  merge-driver logic (also invoked standalone as a git driver)   │  │
│  └───────────────────────────┬────────────────────────────────────┘  │
└──────────────────────────────┼────────────────────────────────────────┘
                               ▼
        ┌───────────────────────────────────────────────────┐
        │              .jaira/  (inside the repo)             │
        │  tickets/<ulid>-<slug>.md   (one file per ticket)   │
        │  .sessions/<session-id>.json (gitignored, ephemeral)│
        │  .gitattributes  (declares merge=jaira on tickets)  │
        │  ~/.jaira/lanes/*.md         (custom lanes, portable)│
        └───────────────────────────────────────────────────┘
```

Dependency direction is one-way and absolute: **CLI → core** and **TUI → core**, never the reverse, and **CLI ↛ TUI** / **TUI ↛ CLI**. Both interfaces are thin shells that translate user/agent input into calls on the same core functions and render the same core state. This is the only way "the CLI and the TUI enforce the schema identically" is actually true rather than aspirational — if either interface re-implemented a rule, drift is inevitable the first time one of them is patched and the other isn't.

### Component Responsibilities

| Component | Responsibility | Notes |
|-----------|-----------------|-------|
| **core: schema** | Ticket struct, YAML frontmatter (de)serialization, field-level types (`status`, `blocked-by[]`, `outcome{}`, `commits[]`, reserved `external{}`) | Single source of truth for what a valid ticket looks like; versioned so future schema changes can migrate old files |
| **core: validation** | Promotion gate (`goal`/`definition-of-done`/`context`/`assignee` required to leave Backlog), `blocked-by` gate (ticket can't start while blockers unresolved), lane-transition legality | Exposed as functions like `Store.Promote(id)`, `Store.Advance(id, toLane, outcome)` — CLI and TUI both call these, never duplicate the checks |
| **core: ticket store** | Read/list/write ticket files under `.jaira/tickets/`, ID generation, atomic writes (tmp file + rename), partial-frontmatter reads for fast listing | No caching layer in v1 (see Ticket Store Design) |
| **core: git access** | Shell out to `git` (or a library) to stage files after a CLI/TUI write, self-register the local merge driver on first run, detect conflict markers, read `commits[]` diffs | Never touches history/branches — this tool writes files, git does the syncing |
| **core: lane resolver** | Load 7 embedded base lanes + scan `~/.jaira/lanes/*.md`, topologically order, detect id collisions, resolve unknown lane ids to passthrough columns | Same parser/schema for built-in and custom lane files |
| **core: merge-driver logic** | The field-aware 3-way merge for ticket files (see Conflict-Tolerant Frontmatter) | Reusable: called by git as `jaira internal-merge-driver`, and reusable directly for TUI-side optimistic-concurrency detection |
| **CLI** | argv parsing, human/`--json` output formatting, mapping subcommands 1:1 onto core calls | Owns zero schema/validation logic itself |
| **TUI** | Rendering, keyboard nav, panels, file-watch-driven refresh | Calls the exact same core functions as the CLI for every mutation (create/move/edit) |

## Recommended Project Structure

```
jaira/
├── cmd/
│   ├── jaira/            # main package: argv dispatch → cli/ or tui/
│   ├── cli/              # one file per subcommand family (ticket.go, lane.go, todo.go, session.go)
│   └── tui/              # Elm/reducer-style TUI: model.go, update.go, view.go, watch.go
├── core/
│   ├── ticket/           # schema.go, store.go (read/write/list), id.go (ULID gen)
│   ├── lane/             # lane.go (parse+order+merge), embed/ (7 built-in lane .md files)
│   ├── merge/            # field-aware 3-way merge, invoked by git AND by TUI staleness checks
│   ├── gate/             # promotion gate + blocked-by gate (pure functions over ticket.Store)
│   └── gitrepo/          # thin wrapper: stage, self-register driver, detect conflicts
├── .jaira/                # created at `jaira init`, ships into the target repo, not this repo
│   ├── tickets/
│   ├── .sessions/
│   └── .gitattributes
└── lanes/                # embedded built-in lane definitions (go:embed / include_str!)
```

### Structure Rationale

- **`core/` has no imports from `cmd/`:** enforced by module boundaries (Go: internal import cycles fail to compile if reversed; Rust: crate/workspace split with `core` as a library crate `cli`/`tui` depend on). This is the actual mechanism — not just convention — that guarantees "one schema enforcer."
- **`core/gate/` is separate from `core/ticket/`:** gates are policy, store is mechanism. Keeping them apart means the promotion-gate rule can be unit-tested against fixtures without touching the filesystem, and it's obvious there is exactly one file to change when a gate rule changes.
- **`core/merge/` is shared by two very different callers** (git, invoked externally; the TUI, invoked in-process for live-edit staleness) — this is intentional reuse of one conflict-resolution mechanism rather than building two (see Q8 / live refresh).

## Ticket Store Design

**ID generation.** Sequential IDs (`JAIRA-42`) need a coordinator to hand out the next number; without one, two offline/parallel sessions that both scan the directory and compute "next available" will pick the same number. A shared counter file doesn't fix this — it becomes the single hottest merge-conflict line in the whole repo, touched on every single ticket creation, which is worse than the problem it solves. Recommendation: **ULID** as the canonical `id`. 48-bit millisecond timestamp + 80 bits of randomness gives a per-millisecond collision probability low enough to ignore at any realistic creation rate, needs no coordinator, and — usefully — sorts lexicographically by creation time, so `ls .jaira/tickets/` is naturally chronological. (Confidence: HIGH — this is a direct consequence of the "no coordinator" constraint plus well-documented ULID collision math.)

**File naming.** `<ulid>-<slug-of-title-at-creation>.md`. The ULID prefix keeps files unique and time-sorted; the slug is a static breadcrumb for humans running `ls`/`grep` and is **not** kept in sync with later title edits (avoids gratuitous renames/git-history churn every time a title changes). The frontmatter `title:` field is the only field callers should trust for display.

**Human-friendly display.** A raw ULID (`01J8QK3ZXG...`) is not a nice thing to type in `blocked-by: 01J8QK3ZXG...`. Two options: (a) accept it — CLI tab-completion and `jaira next`/`jaira show` make typing full IDs rare; or (b) let the CLI accept an unambiguous prefix (`jaira show 01J8` resolves if unique) the way `git` accepts short SHAs. Recommend (b) as a small, cheap CLI-only UX layer — it changes nothing about the stored ID or the conflict properties, it's purely an input convenience in the CLI package.

**Index/cache vs. scanning.** No persistent index or cache in v1. Board rendering only needs frontmatter, not ticket bodies, and frontmatter can be read with a **bounded partial read** (read up to the second `---` delimiter, stop) rather than loading the whole file — this keeps a full-directory scan fast even as individual ticket bodies grow long with outcome/context prose. For the realistic ticket counts this tool will ever see (a single repo's task list — hundreds, plausibly low thousands over a project's life, not millions), a from-scratch scan on every process start is well within "instant startup" and avoids a second source of truth that has to be kept consistent with git operations happening entirely outside the tool's control (`git pull`, manual edits, another session's writes). Building a cache now would be optimizing a problem that doesn't exist yet, at the cost of a staleness-invalidation problem that does. If ticket volume ever becomes a real bottleneck, the fix is an optional gitignored on-disk cache behind the existing `ticket.Store` interface (mtime/hash-invalidated), added later without touching callers — flagged as a YAGNI-for-v1 item, not a design gap.

**Keeping reads fast as count grows.** Two cheap, non-cache techniques do most of the work: (1) partial frontmatter reads as above, and (2) lazy body loading (only read/parse the full file, including outcome/context prose, when a card is expanded in the TUI or `jaira show` is called for one ticket). An archival convention for old Done tickets (e.g., a `tickets/archive/` subdirectory) is a plausible future lever but is not needed to hit "instant" at expected scale — noted as a consideration, not a recommendation for v1.

## Conflict-Tolerant Frontmatter — Recommendation

**The recommendation: a self-installing, field-aware custom git merge driver built into the `jaira` binary itself — not a separate file, not a CRDT library, not "just accept conflicts."**

Mechanics:
1. `jaira init` writes `.jaira/tickets/*.md merge=jaira` into a committed `.gitattributes` (this part *does* travel with `git clone`, unlike merge-driver definitions).
2. On any invocation, `jaira` checks whether the local `.git/config` already has `[merge "jaira"]` registered; if not, it registers it transparently (`git config --local merge.jaira.driver "jaira internal-merge-driver %O %A %B"`) and prints a one-line notice. This is the same one-time-per-clone bootstrap pattern Git LFS uses for its filter driver — well-established precedent, not a novel or fragile mechanism. (Verified: git's merge-driver framework takes an external command and `%O/%A/%B` ancestor/ours/theirs paths; verified: driver registration lives in local `.git/config`, not in anything `git clone` transports — so self-registration on first run is required, not optional, for true zero-config UX.)
3. The driver itself (invoked by git, or reused directly by the TUI, see Live Refresh) parses ancestor/ours/theirs as YAML frontmatter + markdown body and merges **field by field**, not line by line:
   - List-typed fields (`blocked-by[]`, `commits[]`): union + dedupe. Adding a blocker or a commit reference from either side is never lost.
   - `status`: resolved by **lane order, not by timestamp** — whichever side's status is further along the lane sequence wins. Silently reverting forward progress is a worse failure than occasionally keeping a slightly-ahead status.
   - Other contested scalars (`assignee`, `model-tier`, `executed-by`): resolved by the ticket's own `updated-at` timestamp — whichever side is newer wins wholesale for that field. (This reuses a field the ticket needs anyway for TUI recency sort — no extra bookkeeping field required.)
   - Free-text fields (`goal`, `context`, `outcome.*`, the markdown body): standard textual 3-way diff. If the two sides' edits don't overlap, they merge cleanly the same way git already merges any text file. If they genuinely overlap (two people rewrote the same sentence of `goal` differently), this is the one case that becomes a real, human-visible conflict — flagged narrowly (only that field, not the whole file) and surfaced via `jaira resolve <id>`, which shows ours/theirs side by side for just the contested field.
4. Because `status`/`assignee`/`blocked-by`/`commits[]` are by far the most frequent kind of concurrent write in this system (every lane advance touches `status`; every commit touches `commits[]`), this design makes the *common* case a guaranteed non-conflict and leaves only the *rare* case (simultaneous prose rewrites) as manual, well-scoped resolution.

**Why not the alternatives:**
- **Append-only event log per ticket:** genuinely the cleanest git-merge behavior (pure appends at the end of a file merge automatically — verified, this is standard 3-way-merge behavior for non-overlapping insertions), but it requires either splitting into two files per ticket (violates the explicit "one file per ticket" requirement, and the tiny status-log file becomes its own hot conflict target since every mutation touches only it) or keeping the log as a section inside the one ticket file — which then can't use git's whole-file `merge=union` mechanics without also unioning the YAML frontmatter above it, producing duplicate/invalid YAML keys on any concurrent frontmatter touch. The field-aware driver gets the same non-conflicting outcome for the fields that matter without either problem.
- **`.gitattributes merge=union` across the whole ticket file:** rejected for the same reason — it's a real, zero-registration built-in git merge driver (verified: no `.git/config` entry needed for `union`), but it operates on whole lines with no structure awareness, so two sides changing `status:` produces **two** `status:` lines in the merged file — invalid/ambiguous YAML, not a resolved value. Only safe for files that are *purely* append-only end to end, which a mixed frontmatter+prose ticket file is not.
- **Separate status file:** rejected — violates "one file per ticket," and doesn't even remove the hotspot, it just relocates it.
- **CRDT-ish approaches (Automerge, Yjs, OT):** rejected as disproportionate. These exist to support continuous, fine-grained collaborative editing of rich structures; here there are a handful of discrete field transitions merged occasionally (at `git pull` time), not a live multi-cursor editing session. Adopting a CRDT dependency would also fight the "file format IS the API, must stay hand-editable" constraint, since CRDT-backed storage typically isn't plain, directly-editable YAML/Markdown without going through the library. This is the kind of weight this project is explicitly trying not to carry.
- **Simply accept conflicts + good tooling:** rejected as the *primary* mechanism (kept as the fallback for the rare prose-overlap case). The constraints explicitly name "multiple parallel Claude sessions... write the same store" as an expected, everyday pattern, not an edge case — and `status` changes are the single most frequent write the whole system makes (every lane advance is one). Leaving the most common operation to manual resolution by default would mean conflicts happen constantly, which directly undermines "you never lose track... across agent runs."

**Known gap, honestly flagged:** merge drivers only run for merges git itself performs locally (`git merge`, `git pull`, `git rebase`). Merges done through a hosting provider's web UI (e.g., "Resolve conflicts" in a GitHub PR) do **not** invoke local `.gitattributes`/merge-driver configuration (verified via search). Given this project has no web/server component by design and git-native local merges are the expected path, this is an acceptable, minor gap rather than a blocker — but it means a ticket conflict resolved through a hosting provider's web UI will fall back to raw line-level conflict markers. Worth a line in the README, not worth building around.

**Weight flag:** this is the one place the design introduces "magic" (a self-modifying local `.git/config`). It adds no external dependency and no extra binary — the merge driver is a hidden subcommand of the same static binary — but it is still doing something non-obvious on first run. This should be visible to the user (a one-line printed notice) rather than fully silent, so it doesn't feel like the tool is secretly rewriting git configuration.

## Lane Resolution

Seven base lanes are embedded in the binary (`go:embed`/`include_str!`) and always load first, guaranteeing a working board with zero repo config. Each base lane file uses the same schema as a custom lane (see below), just shipped inside the binary instead of on disk.

**Discovery order:** (1) load embedded base lanes; (2) glob `~/.jaira/lanes/*.md`; (3) parse each with the same lane schema; (4) merge.

**Ordering:** custom lanes anchor themselves relative to an existing lane id (`after: review`) rather than a raw numeric position, because numeric ordering breaks down for lane files meant to be portable across projects with different other custom lanes installed — two independently-authored lane files could both claim `order: 35`. Anchored ordering + topological sort avoids this collision class entirely. If a lane file's anchor doesn't exist locally (a shared file referencing a custom lane this machine doesn't have installed), it's inserted at a safe fallback position (just before `Done`) and flagged, rather than failing to load.

**ID collisions:** a custom lane whose `id` collides with a base lane id is rejected with a clear error at load time — never silently override a base lane, since "seven base lanes require no repo config" depends on their behavior being predictable everywhere.

**Unknown-lane passthrough:** a ticket whose `status` doesn't match any locally-loaded lane id (base or custom) is never hidden. It renders in a fixed, deterministic position — a reserved "unrecognized lanes" area, lanes sorted alphabetically among themselves for determinism, no attempt to guess intended position from history — and is read-only in that state: the TUI/CLI block `jaira advance` into or out of a lane id the local install doesn't understand, since it has no prompt/model-tier/contract to enforce there. The fix is simply obtaining the lane file, which is by design a single portable markdown file.

**Lane definition file, concretely:**

```yaml
---
id: security-review
name: "Security Review"
after: review              # anchor for topological ordering
model-tier: strong          # alias, not a raw model string — mapped locally so shared
                              # lane files don't break when a model name changes
terminal: false
input:
  requires: [goal, definition-of-done, context, diff]
output:
  produces: [outcome.what, outcome.why, outcome.resolves]
  schema: outcome
---
# Prompt

You are reviewing a change for security issues before it can merge...
{{ticket.goal}} / {{diff}} interpolated by the CLI when it constructs
the subagent's bounded input — never done by the agent itself.
```

`model-tier` as a local alias (not a hardcoded model identifier) matters specifically because lane files are meant to be portable/shareable — a shared file that hardcodes a model name breaks the moment that model is renamed or a teammate has different access; the tier→model mapping should be a small local config, not part of the shared artifact.

## Agent Integration Surface

Core commands an agent (any bash-capable one, not just Claude Code) uses:

| Command | Purpose | Enforcement |
|---------|---------|--------------|
| `jaira next [--lane <id>] [--json]` | Find the next actionable ticket | Filters out tickets with unresolved `blocked-by` and tickets already claimed by another live session — this is the *discovery*-time check |
| `jaira show <id> --json` | Read a ticket: frontmatter + body + resolved commit diffs, formatted as bounded context | Read-only |
| `jaira claim <id>` | Take a short-lived lease (`assignee-session`, `claimed-at`) | TTL-based (e.g. 30 min); a stale claim is treated as abandoned automatically — no server, no permanent lock, no manual unlock needed |
| `jaira advance <id> --to <lane> --json < outcome.json` | Record outcome and move lanes in one atomic call | This is the **write-time** gate: re-checks `blocked-by`, checks the *target* lane's entry requirements, and — critically — checks that the *current* lane's `output.produces` fields were actually supplied. An agent cannot advance without having produced what the lane demanded. |
| `jaira block <id> --question "..."` | Move to Human lane with a question attached | Special case of `advance` |

**Why two enforcement points, not one:** `next` filtering is a convenience (agents following the happy path never even see blocked work), but `advance` is where correctness actually lives, because nothing stops an agent from ignoring `next` and calling `advance` directly. The gate must live at the mutation boundary, not just the discovery boundary — trusting an agent to only ever call `next` first is exactly the kind of thing this design should not trust the agent to do.

**I/O contract, concretely:** a lane subagent's input is deliberately bounded — the ticket fields named in that lane's `input.requires`, the diff scoped to the ticket's own `commits[]` (not the working tree), and the lane's prompt body, all assembled and interpolated by the CLI (`jaira show --for-lane <id>` style), never by the agent itself. Its output is a JSON object matching the lane's declared `output.schema`; the orchestrating session passes that JSON to `jaira advance`, and `jaira advance` — core code, not agent judgment — is what actually validates it against the schema and rejects malformed output with a message the agent can read and retry against. Structure is enforced by the tool; the agent supplies content, not judgment calls about whether it followed the contract.

## Two-Way TodoWrite Sync

**Direction A — TodoWrite → Backlog (safe, hook-driven):** `PostToolUse` hook matched on `TodoWrite` calls `jaira sync-todos --json <payload>`, which only writes ticket files — it never invokes a tool or the model, so it structurally cannot itself trigger another Claude tool call. No loop risk from this direction alone.

**Identity matching:** TodoWrite rewrites its whole list every call, and (per the fetched hook payload example) items carry only content/status, not a guaranteed stable ID across calls — this needs independent verification against a real captured payload before implementation (flagged below). Recommend jAIra maintain its own small ephemeral map (`.jaira/.sessions/<session-id>/todo-map.json`, gitignored) of `content-hash → ticket-id`, built incrementally: an item whose content hash is new → create a Backlog ticket with `ready: false` and record the mapping; an item whose hash is already known → update the existing mapped ticket's status, never create a duplicate; an item that disappears from the list → leave its mirrored ticket alone (todos are ephemeral working memory, tickets are the durable artifact — a todo list changing shouldn't delete board history).

**Direction B — board → in-session todo list:** hooks cannot synthesize a tool call on Claude's behalf (verified against the hooks reference: hooks can only allow/deny/block or inject `additionalContext`, not issue tool calls). So this direction is necessarily softer: `SessionStart` injects context ("N ready tickets, M in progress") and `jaira start <id>`'s own CLI output is what actually nudges the agent to itself issue a TodoWrite call — enforced by convention/skill instruction, not by a hook. This is a real, honest limitation to design around, not a gap to paper over.

**Loop suppression:** the round trip (B populates a todo mirroring an existing ticket → A's hook fires on that TodoWrite call) is not preventable at the event level — the hook fires on every TodoWrite call regardless of what caused it. What actually prevents infinite growth is that Direction A's matching is idempotent: content that already maps to an existing ticket produces an update, not a new ticket, and if nothing actually changed (hash-compare against the last synced state) `sync-todos` is a no-op — no file write, no git-visible diff. The loop converges to a fixed point after one hop rather than oscillating or duplicating; it should be described as "idempotent by construction," not "prevented."

**Weight flag:** none — this is a hook script plus one CLI subcommand, no daemon, no polling.

## Session Context Panel

Recommend storing session/focus state per working tree, not per repo (committed) and not per machine (global) — `.jaira/.sessions/<session-id>.json`, gitignored via a `.jaira/.gitignore` shipped by `jaira init`. Each file: `{session_id, ticket_id, focus, reasoning, updated_at, host, pid}`, written by `jaira checkpoint <id> --focus ... --reasoning ...` at checkpoints. Committing this would pollute git history with high-frequency, ephemeral content — the exact opposite of "diff-readable" — so it must live outside the commit path even though it lives inside the working tree.

**Multi-session visibility:** the TUI globs `.jaira/.sessions/*.json` and renders one row per live session, showing its current ticket and last checkpoint text, dimming/marking sessions whose file hasn't updated in some window (stale, likely ended or crashed) rather than deleting them. This only shows sessions sharing the *same working tree* — sessions on a teammate's separate clone are not visible to each other, which is the correct, honest consequence of "no server": session focus is explicitly transient scratch state, not something that should sync through git, and there is no other sync channel by design.

## Live Refresh

Watch `.jaira/tickets/` (and `.jaira/.sessions/`) with the platform's native file-watch API (inotify/FSEvents/ReadDirectoryChangesW, typically wrapped by a library like `fsnotify` in Go or `notify` in Rust). Debounce roughly 150–300ms before triggering a re-scan — long enough to collapse a `git pull` touching many files into one reload, short enough to feel live. Re-scan the whole directory rather than trying to apply incremental diffs to in-memory state (consistent with the store design above — a full scan is cheap at expected scale, and incremental patching is a second source of truth to keep correct for no measurable benefit).

**Avoiding clobbering in-flight input:** never replace the TUI's list wholesale and reset cursor/selection by index — always re-match by ticket ID so the user's current selection and any open modal survive a background reload. If the ticket underlying an *open edit* was also changed elsewhere while the user was editing (detected by an `updated-at` mismatch against the base the edit started from — the same timestamp already used by the merge driver for last-write-wins), don't silently overwrite the edit buffer; surface it ("this ticket changed elsewhere") and let the user choose to reload or force-save. This reuses the field-aware merge/staleness mechanism built for git conflicts rather than inventing a second, TUI-only concurrency mechanism — one conflict model, two call sites.

## Anti-Patterns

### Anti-Pattern 1: TUI shelling out to the CLI binary for every mutation
**What people do:** implement the TUI as a wrapper that `exec`s `jaira <subcommand>` on every keypress-driven action, treating the CLI as the "real" interface and the TUI as a shell around it.
**Why it's wrong:** fork/exec overhead on every action fights "instant," and it's a second, redundant way of calling into core when the TUI can link core directly since both live in the same binary.
**Instead:** TUI imports `core` directly, exactly like the CLI does; they are peers over the same library, not a dependency of one on the other.

### Anti-Pattern 2: A sequence-counter file for human-friendly IDs
**What people do:** keep a `.jaira/next-id` file, incremented on every ticket creation, to get `JAIRA-1`, `JAIRA-2`, ...
**Why it's wrong:** that single file is now the hottest possible git-conflict target — every ticket creation anywhere touches the same line, guaranteeing exactly the collision this project is trying to avoid, and worse than doing nothing.
**Instead:** ULID as the canonical ID; sacrifice small sequential numbers for guaranteed collision safety with no coordinator.

### Anti-Pattern 3: A background daemon for file watching or agent orchestration
**What people do:** run a long-lived process to watch files or babysit agent sessions "so the TUI/CLI don't have to."
**Why it's wrong:** explicitly ruled out by the "no server" constraint; also reintroduces exactly the process-management weight the project is designed to avoid (background `claude -p` orchestration is separately out of scope for the same reason).
**Instead:** every CLI invocation and TUI startup is stateless and rebuilds its view from the filesystem; the TUI's own process does its own watching only while it's running, and holds no state that must survive it.

## Build Order

The dependency graph, in the order components must exist:

1. **Core schema + ticket store + minimal CLI** (`create`, `show`, `list`, `move` against the 7 hardcoded lanes; ULID generation; atomic writes). No TUI, no custom lanes, no conflict tooling, no agent surface yet. **This is the smallest technically-usable end-to-end slice**: a teammate can clone the repo, run CLI commands, and get durable, git-committed, shareable board state — the actual sentence in the core value proposition — before a single pixel of TUI exists. It is also independently useful for scripting/agent use, satisfying "CLI usable by any bash-capable agent" immediately.
2. **Read-only TUI.** Renders lanes/cards from core with no mutation path. Lowest-risk TUI slice; proves the rendering and the core→TUI dependency direction before adding write complexity. Combined with (1), this is the **smallest slice that feels like the actual product** to a human.
3. **TUI mutation parity.** Move/create/expand in the TUI calling the same core functions the CLI calls. This is where "one schema enforcer" gets proven for real, since both interfaces are now exercised against the same core in the same codebase.
4. **Promotion gate + `blocked-by` gate.** Needs the ticket schema fields from step 1 (already present, just unenforced); layering gates once the basic read/write loop is solid avoids designing validation against a moving target.
5. **Conflict-tolerant frontmatter** (merge-driver logic + self-registration). Deferred past the single-session-development phase deliberately — conflicts only manifest under concurrent writers — but must land *before* step 6, because agent integration is exactly when concurrent writes start happening for real.
6. **Custom lane loading** (`~/.jaira/lanes/*.md` discovery/ordering/passthrough). Depends on the lane resolver skeleton from step 1 (built there for the 7 base lanes) and on TUI rendering (step 2) to make passthrough columns visibly meaningful.
7. **Agent integration surface** (`next`/`advance`/`claim`, I/O contract enforcement). Depends on lanes (6, for input/output contracts), gates (4, enforced at `advance`), and conflict handling (5, since this is where concurrent writers actually appear).
8. **TodoWrite two-way sync.** Depends only on ticket creation (1) and the promotion gate (4, since mirrored todos land behind it as `ready: false`). Otherwise self-contained (a hook script + one CLI subcommand) — can be built in parallel with 6/7 rather than strictly after them.
9. **Session context panel + live refresh.** Mostly independent polish, but only meaningful to build and test once step 7 produces real concurrent sessions to visualize and refresh against.

## Sources

- Git merge-driver framework (`%O/%A/%B`, local `.git/config` registration, not transported by `git clone`): [libgit2 PR #3564](https://github.com/libgit2/libgit2/pull/3564/files), [Julian Burr — custom git merge drivers](https://www.julianburr.de/til/custom-git-merge-drivers), [Charpentier — custom merge driver](https://charpeni.com/blog/use-custom-merge-driver-to-simplify-git-conflicts)
- `merge=union` as a built-in, zero-registration git attribute; its line-level, structure-unaware semantics: [a3nm — automatic git conflict resolution on logs and sets](https://a3nm.net/blog/git_auto_conflicts.html), GitHub Community discussion #9288 on `merge=union` in web UI (confirms web-UI merges don't honor local `.gitattributes` driver config)
- Git-lfs's one-time-per-clone/machine local filter registration as precedent for self-installing config: [git-lfs-install man page](https://www.mankier.com/1/git-lfs-install)
- 3-way merge behavior for non-overlapping insertions (why append-only logs merge cleanly): general git 3-way-merge documentation/tutorials (Atlassian, AlgoMaster) — MEDIUM confidence, standard behavior, not independently reproduced in a test here
- ULID collision properties (48-bit timestamp + 80-bit randomness, no coordinator needed): [ksuid.net ULID/KSUID comparison](https://ksuid.net/compare/ksuid-vs-ulid), general ULID spec — HIGH confidence, well-documented
- `fsnotify`-style debounce patterns for batch reload on burst file events: general Go ecosystem sources (golinuxcloud, dev.to write-ups on fsnotify debounce) — MEDIUM confidence, common technique, not tied to a specific library API verified here
- Claude Code hooks (event list, JSON payload shape, blocking semantics): fetched from a page presenting itself as "Hooks reference - Claude Code Docs" at `code.claude.com/docs/en/hooks` — see Gaps below, confidence is split by event

## Gaps / Unverified — flagged explicitly

- **TodoWrite item schema.** The fetched hooks reference shows a `{"id", "text", "done"}` example for `TodoWrite` payloads. This may be an illustrative example rather than the actual current field names (training-data recollection suggests something closer to `{content, status, activeForm}` with no guaranteed stable `id`). **This must be verified against a real captured hook payload before implementing Direction A's identity matching** — the identity/content-hashing design above is written to be robust to either shape, but the exact field names are not confirmed.
- **Full Claude Code hook event list.** `PreToolUse`, `PostToolUse`, `Stop`, `SubagentStop`, `SessionStart`/`SessionEnd`, `UserPromptSubmit` are HIGH confidence (consistent with training data and multiple sources). Events like `PostToolBatch`, `TeammateIdle`, `TaskCreated`, `ConfigChange` appeared in the single fetched source only and are **not independently corroborated here** — treat as LOW confidence / possibly non-canonical or version-specific, and re-verify against the current official docs immediately before building the TodoWrite-sync phase, since hooks are a fast-moving product surface.
- **Whether hooks can synthesize a TodoWrite call on Claude's behalf.** Fetched source states hooks can only block/allow/inject context, never issue tool calls directly — this drove the "Direction B is soft/prompted, not hook-enforced" conclusion above. Not independently cross-checked against a second source; if wrong, Direction B could be made hook-enforced instead of convention-based, which would simplify that section.
- **Exact static-binary framework choice** (Go+Bubble Tea, Rust+ratatui, etc.) is intentionally left to STACK.md — this document's recommendations (core/CLI/TUI boundary, merge-driver-as-hidden-subcommand, embedded lane files) hold under either.

---
*Architecture research for: single-binary TUI kanban over git-committed markdown, agent-driven*
*Researched: 2026-08-11*
