# Project Research Summary

**Project:** jAIra
**Domain:** Single-binary terminal TUI kanban board over git-committed markdown tickets, driven by a CLI that AI agents (primarily Claude Code) call from bash
**Researched:** 2026-08-11
**Confidence:** MEDIUM-HIGH overall — HIGH on language/stack/git-integration/feature-landscape fundamentals; MEDIUM-LOW on the two riskiest bets (YAML round-trip fidelity, TodoWrite/Task-tools sync stability); several individual claims are explicitly flagged LOW and are called out below rather than smoothed over.

## Executive Summary

jAIra is a file-based, git-native issue tracker plus a TUI plus an agent-orchestration layer for Claude Code. The adjacent field has been building this exact shape of tool for over 15 years, and the pattern is stark: every tool that kept "the ticket is a single mutable file you hand-edit in place" (Bugs Everywhere, ticgit, gitissius) died, while every long-lived survivor (Fossil, git-bug, Beads) solved the merge-conflict and visibility problems by abandoning that model — replacing it with append-only artifact logs or CRDT object stores. jAIra's hard constraint ("hand-editable, diff-readable markdown frontmatter IS the API") places it structurally on the failure side of that split. The project's answer is not a storage-model change but a narrower, cheaper mitigation: a self-installing, field-aware git merge driver that makes the highest-contention field (`status`) resolve deterministically instead of leaving raw conflict markers. This is a real, well-reasoned mitigation — not the same fix the survivors used, and it carries residual risk that the roadmap needs to carry forward explicitly rather than consider closed.

The recommended technical approach is well-supported and mostly high-confidence: Go 1.26 + Bubble Tea/Lip Gloss/Bubbles for the TUI, Cobra for the CLI, `goccy/go-yaml` for AST-level frontmatter patching, and shelling out to the system `git` binary — all packaged as one static binary with a `core`/`cli`/`tui` split where `core` is the single schema enforcer both interfaces call into. Almost every piece of this stack is HIGH confidence. The one deliberate exception, called out by the stack researcher as the single highest-uncertainty item in the whole project, is YAML round-trip fidelity: both the ticket-store write path and the merge-driver's field-level 3-way merge logic assume that parsing frontmatter to an AST and patching one field back out leaves everything else byte-identical. No Go or Rust library guarantees this for every legal YAML construct. That assumption needs to be proven with a small spike before either the write path or the merge driver is built on top of it — this is a hard dependency ordering, not a nice-to-have.

The remaining risk surface clusters around concurrency and self-certification. "CLI is the single write path" (a PROJECT.md Key Decision) prevents schema drift between the TUI and agents, but by itself does *not* prevent two concurrent CLI invocations on the same machine from silently clobbering each other's writes — that needs explicit locking/atomic-CAS added in the foundation phase, not assumed as a side effect of the single-write-path decision. Separately, PROJECT.md commits unconditionally to two-way TodoWrite sync, but all three research files independently converge on softening that: the underlying Claude Code API has already changed once (TodoWrite → Task tools, officially documented) and is actively still moving, full-list-replace semantics and hook-loop failure classes are real and verified, and the safe path is to ship a one-way mirror first and treat two-way sync as its own hardened, deferred sub-phase. Finally, no lane in the agent pipeline should be allowed to mark a ticket Done on an LLM's own verdict alone — LLM-as-judge research shows judges catch a small fraction of real defects even when reviewing independently — so the Done transition needs a non-LLM-verifiable evidence reference (a real commit, a real test result) as a schema precondition, which changes what every lane's I/O contract must carry from the start.

## Key Findings

### Recommended Stack

Go 1.26.x is the decisive language choice (HIGH confidence) over Rust, on three concrete grounds verified via GitHub/crates.io APIs, not vibes: Go's cross-compilation to macOS/Windows from the stated Linux/WSL2 dev box is a solved, zero-toolchain problem (Rust's `cross` helper is unmaintained since Feb 2023); the Go YAML ecosystem has an actively maintained round-trip-capable library where Rust's canonical `serde_yaml` was archived by its own maintainer in March 2024 with no successor; and Charm's Bubble Tea/Lip Gloss/Bubbles stack is "batteries-included" for exactly this multi-column, detail-pane TUI shape in a way Rust's `ratatui` (an equally solid rendering primitive) leaves more to the app to build.

**Core technologies:**
- **Go 1.26.5** — language/toolchain: native cross-compilation, static binaries by default, mature `os/exec` for shelling to git
- **Bubble Tea v2.0.8 + Lip Gloss v2.0.6 + Bubbles v2.1.1** — TUI runtime, styling, and pre-built list/table/viewport components; de facto Go TUI standard, proven at scale (`gh dash`, `glow`, `soft-serve`)
- **Cobra v1.10.2** — CLI framework; industry standard for this shape of tool (`kubectl`, `gh`, `docker`), self-describing command tree an agent-facing skill can introspect
- **goccy/go-yaml v1.19.2** — frontmatter parsing via AST/path-based node editing, not struct marshal — the only actively maintained Go or Rust library purpose-built for single-field patching without reformatting untouched fields (MEDIUM confidence on perfect fidelity — see spike requirement below)
- **Shell out to system `git` via `os/exec`** — not go-git/libgit2/gix — for byte-identical diff/log output a human reviewer will recognize, with zero added binary weight or CGO risk
- **fsnotify v1.10.1** — cross-platform file watching; HIGH confidence on the library, MEDIUM on WSL2 behavior specifically (inotify is unreliable across the `/mnt/c` DrvFs boundary — document that `.jaira/` should live in the WSL2 filesystem)

Supporting: `alecthomas/chroma` v2 for diff syntax highlighting, `glamour` (optional, evaluate against complexity budget) for markdown rendering, `goreleaser` for the cross-platform release pipeline. Explicitly avoid: `serde_yaml`, the original `gopkg.in/yaml.v3` (both unmaintained), any CGO dependency (defeats the static-binary constraint), and full struct marshal/unmarshal as the *write* strategy for tickets (silently drops unrecognized keys — a direct risk to the reserved `external:` block).

### Expected Features

**Must have (table stakes):**
- Plain-text/markdown storage that produces clean git diffs (already the project's chosen format)
- Kanban columns with vim + arrow-key keyboard navigation, inline detail pane
- Live filter/search across tickets — **verified gap: not yet in PROJECT.md's Active requirements, should be added for v1**
- Instant local startup, zero-config first run on a fresh clone
- `blocked-by` dependency declaration with start-blocking enforcement (flat list — do not build a typed dependency graph, see anti-features)
- A conflict-tolerant `status` field policy — the highest-complexity table-stakes item; no surveyed plain-file tool has fully solved this (see Critical Pitfalls)

**Should have (differentiators — no surveyed competitor combines these):**
- Lanes binding a prompt + model tier + I/O contract (board as a controlled multi-model pipeline) — the project's load-bearing differentiator; sequence it early relative to the other differentiators once the base board is proven
- Promotion gate (DoD/goal/context/assignee required before leaving Backlog) — cheap to build, high review-quality payoff
- Human vs. `executed-by` (model) attribution split — no surveyed tool separates these
- Definition-of-done-anchored three-field outcome block (what/why/resolves)
- Board-native Human lane for agent questions, rendered as a persistent column rather than a modal prompt
- Unknown/custom lanes degrading to read-only passthrough columns instead of hiding tickets

**Defer (v1.x / v2+):**
- Two-way TodoWrite/Task-tool sync — ship the one-way mirror first (see cross-cutting resolution #4 below)
- Session-context panel — depends on Claude reliably writing structured checkpoints; needs the CLI/skill stable first
- Named filter shortcuts (a lighter "saved view") — trigger reactively, not speculatively
- Jira/YouTrack adapter via the reserved `external:` block — no real integration partner yet
- **New anti-features surfaced by research, not originally in PROJECT.md's Out of Scope list:** a full multi-type dependency graph (Beads-style blocks/related/parent-child/discovered-from) — the flat `blocked-by` list already covers the stated need; and a CRDT-based conflict-free engine — solves conflicts completely but requires abandoning the hand-editable-file constraint, which is a harder violation than tolerating occasional conflicts.

### Architecture Approach

One static binary with three peer components — `cmd/cli`, `cmd/tui`, and a shared `core` library that both import and neither wraps — where `core` alone owns schema, validation (promotion gate, `blocked-by` gate), the ticket store, lane resolution, git access, and the merge-driver logic. Dependency direction is one-way and enforced by module/compile boundaries, not convention: this is what actually makes "one schema enforcer" true rather than aspirational. IDs are ULIDs (not a sequential counter file, which would itself become the single hottest merge-conflict target in the repo). Ticket store reads are a bounded partial read of the frontmatter block only; no persistent index/cache in v1 — a directory scan on every process start is well within "instant" at the realistic scale (hundreds to low thousands of tickets per repo) and avoids a second source of truth that has to track git operations happening outside the tool's control.

**Major components:**
1. **core: schema + validation** — ticket struct, YAML (de)serialization, promotion gate, `blocked-by` gate, lane-transition legality — exposed as functions both interfaces call, never reimplemented
2. **core: ticket store** — read/list/write under `.jaira/tickets/`, ULID generation, atomic writes (temp file + rename)
3. **core: git access** — shells to `git`, self-registers the local merge driver on first run, detects conflict markers
4. **core: lane resolver** — 7 embedded base lanes + `~/.jaira/lanes/*.md` discovery, anchor-based topological ordering, unknown-lane passthrough
5. **core: merge-driver logic** — field-aware 3-way merge, invoked both by git externally and by the TUI in-process for live-edit staleness detection (one conflict model, two call sites)
6. **CLI / TUI** — thin shells: argv parsing / rendering + keyboard nav respectively, zero schema logic of their own

**Recommended build order** (from Architecture's dependency graph): core schema + ticket store + minimal CLI → read-only TUI → TUI mutation parity → promotion/blocked-by gates → conflict-tolerant frontmatter (merge driver) → custom lane loading → agent integration surface (`next`/`advance`/`claim`) → TodoWrite sync → session panel/live refresh polish.

## Cross-Cutting Resolutions

The four research files surfaced points that interact or tension with each other and with PROJECT.md. These are resolved (or explicitly left open) below rather than left as separate, unreconciled findings.

**1. Frontmatter conflict model — resolved, with residual risk carried forward, not eliminated.**
Resolution: keep `status` (and other scalar fields) as plain scalars in a mutable working-tree file; do **not** adopt an append-only log or a CRDT engine. Mitigate via a self-installing, field-aware git merge driver built into the `jaira` binary (`jaira internal-merge-driver`), registered into local `.git/config` on first run (the same one-time bootstrap pattern Git LFS uses) and declared via a committed `.gitattributes`. Architecture directly evaluated and rejected the append-only-log-inside-one-file alternative that Pitfalls tentatively raised as a "consider" option (`history:` array with `status` derived from the last entry): splitting into two files violates "one file per ticket," and keeping the log as a section inside the same file breaks `merge=union` whole-file semantics, producing duplicate/invalid YAML keys on any concurrent frontmatter touch. Architecture's rejection is more mechanistic than Pitfalls' tentative suggestion, so the field-aware-driver design is the one to build — but note this is a case where the two files partially disagreed, and Pitfalls' suggestion is documented here rather than silently dropped.

Does the merge driver genuinely answer Fossil's objection? **Partially, not fully.** Fossil's design rejects mutable-file ticket state categorically because concurrent edits become text-merge conflicts by construction; the field-aware driver narrows *which* edits become conflicts (down to genuinely overlapping prose rewrites of the same field) rather than eliminating the category. Residual risk, stated plainly:
- Merge drivers only run for merges git itself performs locally (`git merge`/`pull`/`rebase`). Merges resolved through a hosting provider's web UI (GitHub "Resolve conflicts") do **not** invoke `.gitattributes`/driver config — a ticket conflict resolved that way falls back to raw conflict markers. Acceptable given no server/web component is planned, but must be documented, not silently assumed away.
- Self-registration happens on `jaira` invocation, not on `git clone`. A fresh clone on which a plain `git pull`/`git merge` happens *before* `jaira` is ever run for the first time will not have the driver registered yet, and gets default git behavior for that merge.
- Genuine overlapping prose edits to the same free-text field (`goal`, `context`, `outcome.*`) still produce a real, human-visible conflict requiring manual resolution (`jaira resolve <id>`) — this is the deliberately-kept fallback, not a bug.
- The self-modifying local `.git/config` is a piece of "magic" the design itself flags as needing a visible one-line notice on first run rather than being fully silent — it is a new, tool-owned failure surface (if the driver logic itself has a bug, `git merge` calls into buggy jAIra code, not into a well-tested third-party tool).

**2. Write races — a real gap that "CLI as single write path" does not close; must add locking, in the foundation phase.**
Pitfalls is explicit and correct: a single logical write path can still be invoked concurrently by multiple processes (e.g., two subagents spawned in parallel within one Claude Code session, or two clones on the same machine). Architecture's ticket-store design specifies atomic writes (temp file + rename), which prevents corrupt/partial writes but does **not** prevent a lost update — two processes can both read the old value, both compute a new value, and the second atomic rename simply clobbers the first with no conflict marker and no warning. Resolution: `core: ticket store` must add an OS-level lock (flock on the ticket file, or a per-repo lock directory) around every read-modify-write, in addition to the atomic rename Architecture already specifies. This must land in the **foundation phase**, before any concurrency is possible in practice (before the agent pipeline or TodoWrite sync ships) — retrofitting it after concurrent usage exists in the wild is a correctness bug discovered in production, not a design choice.

**3. YAML round-trip fidelity — a hard dependency: prove the spike before building on it.**
Stack flags this as the single highest-uncertainty item in the whole stack (MEDIUM confidence, not HIGH) and explicitly recommends a spike: write a small AST-edit-and-reprint test against real-looking ticket files with comments, and verify the diff output is a single-line change. Architecture's entire conflict-tolerant-frontmatter design (field-by-field 3-way merge on ancestor/ours/theirs) *assumes* this AST-level field-isolation capability works as advertised. The dependency ordering is therefore: **the YAML AST-patch fidelity spike must be proven before the ticket-store write path is built, and before the merge-driver logic is built on top of it** — both consume the same unproven assumption, and a failure discovered after both are built is a schema-format migration, not a library swap. If the spike reveals gaps (anchors, aliases, exotic flow styles), the mitigation Stack already proposes — validate/normalize frontmatter shape on read and reject/flag exotic constructs — should be designed alongside the spike, not after.

**4. TodoWrite two-way sync — do not commit v1 to this; the underlying API is unstable.**
All three of Stack (implicitly, via not being deep here), Features, and Pitfalls converge, and Features independently verified via official Claude Code docs that `TodoWrite` is no longer the default as of Claude Code v2.1.142 — Anthropic has already migrated to structured `TaskCreate`/`TaskUpdate`/`TaskGet`/`TaskList` tools with per-item IDs, `owner`, and dependency fields. Pitfalls independently corroborates via GitHub issues and a VentureBeat report that this area is actively still moving. Architecture's own hook-payload field-name detail (`{id, text, done}` vs. the more likely `{content, status, activeForm}`) is flagged MEDIUM-LOW confidence and explicitly needs re-verification against a real captured payload before implementation — this is not resolved by any research file.

**Recommendation: two-way sync does not belong in v1.** Ship a one-way `TodoWrite`/Task-tools → Backlog mirror first (targeting the current Task tools API, not the deprecated TodoWrite shape), behind the promotion gate (`ready: false`), with identity matching by content-hash/stable-ID rather than list position (TodoWrite performs a full-list replace every call — this is HIGH confidence, verified against the official tool description and multiple GitHub issues, so any sync logic must diff full snapshots itself, never treat a call as a delta). Treat two-way sync as its own hardened sub-phase, built only after the one-way mirror is proven and isolated behind an internal adapter interface so the next Anthropic API change is a contained patch. **This is a direct conflict with PROJECT.md's Key Decision ("TodoWrite sync is two-way") stated as an unconditional v1 commitment** — flagged below as an open decision needing the user's explicit call, not silently overridden here.

**5. Done cannot be self-certified — a non-LLM signal must gate the transition, and it changes the lane I/O contract.**
LLM-as-judge research (two independent papers, MEDIUM confidence, small number of studies but consistent failure pattern) found judges catch a small fraction of real defects even reviewing independently, and can hallucinate evidence to retroactively justify a "compliant" verdict. An implementer agent self-certifying its own Definition of Done is the worst case of this — not even an independent judge. The gating design: (a) the Review→Done transition must require a *different model-tier invocation* than the one that implemented the ticket to have written the review outcome — enforceable as a schema precondition, the same mechanism the promotion gate already uses; (b) the outcome block's evidence must reference something the CLI can mechanically verify exists (a commit hash, a test result, a specific diff line) rather than accept free-text restatement of the DoD; (c) a lane's Review-stage output schema must therefore carry an evidence-reference field from the start, not as a later addition — this constrains the I/O-contract design (`input.requires`/`output.produces`) introduced in the agent-pipeline phase, and should be designed alongside the lane contract itself, not bolted on afterward. Cap review-retry loops (bounded N failures → force to Human/Blocked) to avoid unattended cost blowup from an implement/critique disagreement loop.

**6. Adoption risk — honestly, this project sits on the failure side of a 15+ year graveyard, and its differentiator is a reasoned bet, not proven.**
The distributed-issue-tracking space has a consistent, well-evidenced failure pattern (MEDIUM-HIGH confidence, multiple independent sources): tools die from (a) no visibility outside `git clone` + the specific tool, (b) tooling drift making an abandoned format's files an inert pile, and (c) merge-conflict pain on exactly the fields that change most. Every long-lived survivor (Fossil, git-bug, Beads) solved this by abandoning the hand-editable-mutable-file model jAIra has chosen for a hard, deliberate constraint. What jAIra offers instead: a narrower target audience (developer + Claude Code, terminal-only, not a cross-functional stakeholder audience that predecessors overreached toward and died under the weight of chasing), a first-class commitment to separating format from tool (a versioned, documented schema so an abandoned tool doesn't also make the data unreadable), and the field-aware merge driver from resolution #1 above. **None of these are proven — they are reasoned mitigations informed by why predecessors specifically died, not evidence that this predecessor will survive.** Worth noting for calibration: the two most credible "solved it" adjacent tools cited in Features research (Beads, released Oct 2025; Backlog.md, active 2025-2026) are themselves too new to have a real longevity track record — this is a young, still-unsettled category, not one with an established recipe jAIra is failing to follow.

## Implications for Roadmap

### Phase 1: Foundation — Schema, Store, and the YAML Spike

**Rationale:** Everything else depends on the ticket store existing and being trustworthy; this is also the "smallest technically-usable end-to-end slice" per Architecture's build order — a teammate can clone and use the CLI before a single pixel of TUI exists.
**Delivers:** Core schema + strict type validation (not just "is it YAML"), referential-integrity checks on `blocked-by`, ULID-based ticket store with atomic writes **and explicit file locking** (resolution #2), minimal CLI (`create`/`show`/`list`/`move` against the 7 hardcoded lanes).
**Hard prerequisite:** the YAML AST-patch fidelity spike (resolution #3) must be run and pass *before* the write path is built on the assumption it works.
**Addresses:** Table-stakes markdown storage, `blocked-by` declaration.
**Avoids:** Write-race lost updates (Pitfall 2), hallucinated-ID/malformed-frontmatter drift (Pitfall 5), and frames the format/tool-separation discipline that mitigates the adoption-graveyard risk (Pitfall 1/resolution #6).

### Phase 2: TUI Board (Read-Only, then Mutation Parity)

**Rationale:** Lowest-risk TUI slice first proves core→TUI rendering and the one-way dependency direction before adding write complexity; combined with Phase 1 this is the smallest slice that feels like the actual product.
**Delivers:** Kanban rendering of the 7 base lanes, keyboard nav (vim + arrows), inline detail pane, live filter/search (verified missing gap — add here), resize/SIGWINCH handling, debounced file-watch refresh — then TUI mutation calling the same core functions the CLI calls.
**Uses:** Bubble Tea/Lip Gloss/Bubbles, fsnotify (debounced).
**Avoids:** TUI terminal/Unicode/WSL2 traps (Pitfall 7) — must be right before the board is trusted as primary interface, not deferred as polish.

### Phase 3: Gates and Conflict Tolerance

**Rationale:** Gates need the schema fields from Phase 1 (present but unenforced); conflict tolerance only matters once concurrent writers are imminent, but must land *before* Phase 4's agent integration, since that's exactly when concurrent writes start happening for real.
**Delivers:** Promotion gate (DoD/goal/context/assignee), `blocked-by` start-blocking enforcement, the field-aware self-registering merge driver (resolution #1) with its documented residual risks (hosted-UI merges, pre-first-run clones, genuine prose conflicts routed to `jaira resolve`).
**Implements:** `core: gate`, `core: merge` per Architecture's project structure.

### Phase 4: Agent Integration Surface (Lanes, Human Lane, Done-Gating)

**Rationale:** Depends on lanes (custom + base), gates, and conflict handling all already existing; this is the project's headline differentiator and its highest-complexity feature.
**Delivers:** `next`/`show`/`claim`/`advance`/`block` commands with two enforcement points (discovery-time and write-time), lane I/O contracts carrying an **evidence-reference field** for Review-stage output (resolution #5) so Done cannot be self-certified, portable custom lane files with anchor-based ordering, unknown-lane passthrough, the Human lane with an unmissable visibility mechanism (Pitfall 8 — pull-based HITL is prone to silent starvation with no server/push channel available).
**Implements:** `core: lane`, the agent-integration commands from Architecture.
**Avoids:** LLM-rubber-stamp Done (Pitfall 3), Human-lane starvation (Pitfall 8).

### Phase 5: One-Way TodoWrite/Task-Tools Mirror

**Rationale:** Self-contained (a hook + one CLI subcommand); only needs ticket creation (Phase 1) and the promotion gate (Phase 3) to exist, so it can run in parallel with Phase 4 rather than strictly after it — but ship the one-way direction only.
**Delivers:** `PostToolUse`-hook-driven mirror into Backlog behind `ready: false`, full-snapshot diffing (not delta-treatment) for identity matching, hook kept deliberately fast/dumb and fail-open.
**Explicitly does not deliver:** two-way sync (resolution #4) — that is a separate, later, hardened sub-phase, gated on the one-way mirror proving stable and the Task-tools API settling further.

### Phase 6: Session Panel, Two-Way Sync Hardening, Polish

**Rationale:** Mostly independent polish, only meaningful once Phase 4 produces real concurrent sessions to visualize/refresh against; two-way TodoWrite sync (if pursued at all) belongs here, isolated behind an adapter interface.
**Delivers:** Session-context panel (`.jaira/.sessions/*.json`, gitignored, per-working-tree), live-refresh staleness reconciliation reusing the merge/staleness mechanism from Phase 3, and — only if the user confirms it should still be pursued — a hardened two-way sync built against the current Task tools API.

### Hard Dependency Ordering (cross-phase, explicit)

1. **YAML AST-patch fidelity spike** must be proven before the ticket-store write path (Phase 1) and before the merge-driver logic (Phase 3) are built — both assume it works.
2. **Atomic writes + file locking** in the ticket store (Phase 1) must exist before any TUI mutation (Phase 2), agent pipeline (Phase 4), or TodoWrite sync (Phase 5) — all three introduce real concurrency.
3. **Strict schema validation + referential integrity** (Phase 1) must land before the promotion gate and `blocked-by` enforcement (Phase 3), which trust the schema they build on.
4. **Promotion gate** (Phase 3) must exist before the TodoWrite mirror lands tickets "behind" it (Phase 5) and before lane-advance logic is meaningful (Phase 4).
5. **Field-aware merge driver** (Phase 3) must exist before the agent-pipeline phase (Phase 4) is exercised with real concurrent writers, per Architecture's own build order.
6. **Lanes-as-agent-steps with the evidence-reference I/O contract** (Phase 4) must exist before the Human lane, portable lane files, and session panel (later in Phase 4/6), which all reuse the lane mechanism.
7. **One-way TodoWrite mirror** (Phase 5) must ship and prove stable before any two-way sync work (Phase 6) is attempted — all three research files converge on this ordering independently.

### Research Flags

Needs deeper research during planning (`--research-phase`):
- **Phase 1:** the YAML round-trip spike itself is a research/prototype task, not standard implementation — treat it as a research-phase deliverable.
- **Phase 4:** Human-in-the-loop visibility mechanisms with no server/push channel — general HITL literature was cross-referenced but not domain-specific; the exact nudge design (shell prompt hook vs. in-session Claude nudge) needs validation against real usage.
- **Phase 5/6:** the current Claude Code hook payload schema for TodoWrite/Task tools (exact field names) is explicitly unverified in Architecture's own Gaps section and must be re-checked against live/current official docs immediately before implementation, since this is a fast-moving product surface that has already changed once.

Phases with standard, well-documented patterns (skip deep research):
- **Phase 2:** Bubble Tea/Lip Gloss TUI patterns (multi-column layout, viewport detail pane, file-watch-driven re-render) are well-trodden in the Charm ecosystem, not novel.
- **Phase 3 (merge driver mechanics):** git's merge-driver framework, local `.git/config` registration, and the git-lfs self-registration precedent are all well-documented, verified patterns.

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | HIGH on language/git-integration/CLI/distribution; MEDIUM on TUI framework specifics; MEDIUM-LOW on YAML round-trip fidelity (explicitly flagged, spike required) | Versions verified via GitHub Releases API / go.dev/dl (primary sources) |
| Features | MEDIUM-HIGH; TodoWrite→Task-tools migration is HIGH (official docs, directly fetched) | A few named tools' exact conflict-handling mechanics (tissue, git-issue) are LOW confidence, single-sourced, not independently verified |
| Architecture | MEDIUM-HIGH; core structure and the conflict-driver recommendation are HIGH (mechanistic reasoning + verified git merge-driver framework docs) | Claude Code hook payload field names and the full hook-event list are explicitly flagged MEDIUM-LOW / LOW — single fetched source, not cross-checked |
| Pitfalls | MEDIUM-HIGH; file-tracker failure history and TodoWrite full-replace/hook-loop mechanics are HIGH (primary GitHub issues, official tool descriptions, Fossil's own docs) | LLM-as-judge defect-recall numbers are MEDIUM (two studies, directional not exact); "Rule of Two" attribution to Meta is explicitly LOW-MEDIUM (secondary citation, not Meta's own docs) |

**Overall confidence:** MEDIUM-HIGH. The stack and feature landscape are on solid, verified ground. The two genuinely load-bearing uncertainties — YAML round-trip fidelity and TodoWrite/Task-tools API stability — are both explicitly flagged by their respective researchers as needing validation before or during early implementation, not as settled facts.

### Gaps to Address

- **YAML AST-patch fidelity** — not independently hand-tested by any researcher; requires the Phase 1 spike before the write path/merge driver are built on the assumption (resolution #3).
- **Exact TodoWrite/Task-tools hook payload schema** — Architecture's own Gaps section flags the fetched field-name example as possibly illustrative rather than accurate; must be re-verified against a real captured payload or current official docs before Phase 5 implementation.
- **Whether hooks can synthesize a tool call on Claude's behalf** — Architecture's "Direction B is soft/convention-based, not hook-enforced" conclusion rests on a single fetched source; if wrong, the board→todo-list direction of sync could be simplified.
- **The scalar-vs-conflict-tolerance decision's residual risk** (resolution #1) is resolved for the roadmap's purposes but not eliminated — the hosted-merge-UI gap and pre-first-run-clone gap should be documented in the README, per Architecture's own recommendation, not silently left implicit.
- **Two-way TodoWrite sync vs. PROJECT.md's Key Decision** (resolution #4) is a genuine, unresolved conflict between what PROJECT.md commits to and what all three researchers recommend — flagged as an open decision below, not resolved unilaterally here.

## Open Decisions Needing User Input

These are not resolved by research and should be explicitly decided (or PROJECT.md's Key Decisions table amended) before or during roadmap creation:

1. **Two-way TodoWrite sync in v1.** PROJECT.md's Key Decisions table states "TodoWrite sync is two-way" as a v1 commitment. All three research files (Stack implicitly, Features explicitly, Pitfalls explicitly) recommend shipping one-way first and treating two-way as a deferred, hardened sub-phase because the underlying Claude Code API is actively unstable. Decision needed: amend the Key Decision to phase two-way sync into v1.x, or accept the added risk of building it in v1 anyway.
2. **Acceptance of the merge driver's residual risk as sufficient for v1**, given it does not fully answer Fossil's structural objection (resolution #1) — vs. investing further (e.g., toward a partial append-only log for just the `status` field, which Architecture rejected on YAML-structural grounds but which could be revisited with a different file layout if the residual risk proves unacceptable in practice).
3. **Live filter/search** was found to be a table-stakes feature missing from PROJECT.md's Active requirements — confirm it should be added to Phase 2 scope.
4. **Whether to formally add the two new anti-features** surfaced by research (full multi-type dependency graph; CRDT-based conflict-free engine) to PROJECT.md's Out of Scope list, to close off scope-creep vectors research specifically identified.

## Sources

### Primary (HIGH confidence)
- GitHub Releases API / go.dev/dl / crates.io API — version verification for the entire stack (bubbletea, lipgloss, bubbles, cobra, fsnotify, goccy/go-yaml, chroma, ratatui, gix)
- [Fossil: Bug-Tracking In Fossil](https://fossil-scm.org/home/doc/tip/www/bugtheory.wiki) — official design rationale for append-only ticket artifacts vs. mutable working-tree files
- [Claude Code Todo Lists docs — code.claude.com](https://code.claude.com/docs/en/agent-sdk/todo-tracking) — official docs confirming TodoWrite → Task tools migration
- [anthropics/claude-code issues #10205, #6674, #2250, #1824, #59962, #6159, #69093](https://github.com/anthropics/claude-code) — hook-loop, full-list-replace, stale-state, Unicode-corruption behaviors
- [charmbracelet/lipgloss issue #562](https://github.com/charmbracelet/lipgloss/issues/562) — emoji/Unicode width layout misalignment
- [go-yaml/yaml issue #709](https://github.com/go-yaml/yaml/issues/709) — documented limitation, no Go YAML library fully round-trips comments
- Git merge-driver framework docs (libgit2 PR #3564, git-lfs-install man page) — merge-driver `%O/%A/%B` mechanics, local-only registration, self-registration precedent

### Secondary (MEDIUM confidence)
- [Current State of the Distributed Issue Tracking survey](https://matej.ceplovi.cz/blog/current-state-of-the-distributed-issue-tracking.html) — dead file-based tracker pattern
- [Backlog.md](https://github.com/MrLesk/Backlog.md/), [Beads/bd](https://github.com/steveyegge/beads) — closest direct competitors, both too new for a real longevity track record
- [Catching One in Five (arXiv 2606.10315)](https://arxiv.org/html/2606.10315), [The Stability Trap (arXiv 2601.11783)](https://arxiv.org/pdf/2601.11783) — LLM-as-judge defect-recall studies, small sample, directional
- Roo Code / Cline comparisons, GitHub Copilot Workspace "Mission Control" coverage — adjacent multi-agent orchestration precedent

### Tertiary (LOW confidence — flagged, not upgraded)
- tissue / git-issue exact conflict-handling mechanics — not documented in sources found
- "Rule of Two" attribution to Meta — secondary citation only, not Meta's own documentation
- Architecture's fetched hook-event list beyond the well-corroborated core set (`PostToolBatch`, `TeammateIdle`, etc.) — single source, not independently cross-checked
- General feature-creep/bloat literature — generic product-management sources, applied as domain-general context only

---
*Research completed: 2026-08-11*
*Ready for roadmap: yes, with the open decisions above flagged for explicit resolution*
