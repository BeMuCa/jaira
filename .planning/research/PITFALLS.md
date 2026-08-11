# Pitfalls Research

**Domain:** Git-committed markdown ticket tracker with a TUI kanban board and Claude-Code-driven agent pipeline lanes
**Researched:** 2026-08-11
**Confidence:** MEDIUM-HIGH (file-based-tracker history and TodoWrite mechanics are well-evidenced; LLM-judge-reliability numbers come from a small number of studies and should be treated as directional, not exact)

This file ranks pitfalls **Critical → Moderate → Minor**. Critical pitfalls are the ones that can kill adoption or trust in the tool outright; they should be addressed in the earliest phases that touch the relevant subsystem, not deferred as polish.

---

## Critical Pitfalls

### Pitfall 1: File-based trackers have a 15+ year track record of dying, and the reasons are structural, not incidental

**What goes wrong:**
Every prior generation of "tickets as files/artifacts in the repo" tool — ticgit, Bugs Everywhere, gitissius, Fossil's early file-based experiments, and (per Hacker News commentary) git-bug's mainstream traction — ends up abandoned or niche. ticgit's own maintainers marked it dead years ago. A 2018/2022/2024-era Hacker News pattern around git-bug repeatedly surfaces the same objection: management and non-developers cannot see the tracker, cannot pull a report out of it, and it has no gravity outside the one clone that happens to have the tool installed. A survey of the whole distributed-issue-tracking ecosystem ("dead bodies and destruction everywhere" — Verdun analogy) found that nearly every contender in this space is defunct, abandoned even by its original author. (Source: matej.ceplovi.cz survey; confirmed independently by ticgit's own GitHub wiki status and multiple git-bug HN threads. MEDIUM-HIGH confidence — multiple independent sources converge on the same failure pattern.)

**Why it happens:**
Three compounding failure modes recur across every dead project in this space:
1. **No one outside the repo can see it.** A file-based tracker's visibility is exactly the visibility of `git clone` + a specific tool installed locally. Non-technical stakeholders, a teammate who hasn't installed the binary yet, or anyone glancing at GitHub's web UI sees either nothing or a pile of raw markdown files — not a board.
2. **Tooling drift.** The tool that reads the format is a single point of failure. If the maintainer stops updating it, the format becomes an inert pile of files nobody can browse meaningfully (this is literally jAIra's own "unknown lane → passthrough column" mitigation, but for the *entire tool*, not just lanes).
3. **Merge-conflict pain on the exact fields that change most often** (status, ordering) creates enough friction that teams revert to a hosted tracker after the first painful conflict — this is the #1 complaint in Bugs Everywhere / gitissius postmortems.

Crucially: **Fossil**, the most serious and long-lived attempt to solve "put the tracker where the code is," explicitly did **not** store tickets as plain working-tree files for this exact reason. Fossil's own design docs (fossil-scm.org/.../bugtheory.wiki, read directly) state the current state of a ticket is reconstructed by replaying an *append-only* sequence of immutable "ticket change artifacts" keyed by ticket ID and timestamp — never a mutable file that two people edit and 3-way-merge. Fossil's rationale, verbatim in spirit: check-ins are immutable, ticket state is not, so ticket state cannot live inside the same mutable-file model as the checked-out tree, or every concurrent edit becomes a text merge conflict. jAIra's design (`status` as a mutable YAML frontmatter field in a working-tree file, merged via ordinary git 3-way merge) is precisely the model Fossil rejected. This is not a minor implementation detail — it is the single most battle-tested piece of prior art in this exact problem space choosing the opposite architecture. (HIGH confidence — read directly from Fossil's own documentation.)

**How to avoid — and where jAIra's bet is actually more defensible than its predecessors:**
jAIra differs from the graveyard in ways that matter, and the roadmap should lean into these differences rather than assume the graveyard doesn't apply:
- **It targets a single audience that already lives in the terminal** (the developer + Claude Code), not a cross-functional stakeholder audience. It should *not* try to become the tool a PM checks — that pull is what killed nothing here but caused every predecessor to overreach toward "let's add a web view" and then die under that weight instead. Resist that pull explicitly (see Pitfall 7).
- **Visibility problem is mitigated, not solved:** because jAIra is explicitly TUI-only with no server (see Out of Scope), anyone who wants to see the board must clone + run the binary. This is a deliberate tradeoff — document it as one, don't pretend the visibility problem is gone. A GitHub-hosted repo will render the ticket markdown files as ugly raw YAML+markdown in the web UI; that is the honest cost of this architecture, not a bug to fix later.
- **Tooling-drift risk is direct-proportional to how replaceable the schema is.** Because the ticket format is "the API" (per PROJECT.md), treat schema stability as a first-class deliverable of the earliest phase, with a documented, versioned frontmatter schema and a migration story — the predecessor tools that died fastest never separated "format" from "tool" cleanly, so an abandoned tool made the data unreadable too.
- **Take the merge-conflict lesson from Fossil seriously even without adopting Fossil's architecture**: since jAIra can't move to an append-only event log without abandoning "hand-editable, diff-readable markdown" (a stated hard constraint), the mitigation has to be at the *field* level — see Pitfall 2.

**Warning signs:**
- Early users say "I forgot the board existed" after a few days — a sign the TUI isn't being opened enough to justify the format tradeoff.
- Teammates ask "can I just see this on GitHub" within the first week of dogfooding.
- The tool's own maintainers stop running `jaira` and start using raw `git log` / grep on ticket files — the earliest sign of exactly the "tooling drift → inert files" failure mode.

**Phase to address:**
Foundation phase (ticket store + CLI schema). The schema-stability and format/tool separation decisions must be locked before the TUI or agent pipeline is built on top, because every subsequent phase inherits the ticket-store's failure modes.

---

### Pitfall 2: Merge conflicts on the exact frontmatter field that changes most (status/lane), and the CLI-single-write-path mitigation is necessary but not sufficient

**What goes wrong:**
Git's line-based 3-way merge is actually *fine* for concurrent edits to *different* lines — a common misconception is that any two people touching the same ticket file conflicts. In practice: two teammates (or two parallel agent sessions in separate worktrees/clones) editing *different* fields of the same ticket (one sets `assignee`, another appends a commit hash) auto-merge cleanly. The real failure mode is narrower and more concrete: **two actors change the same line — almost always `status:`** — between the same git base. Example: Session A moves ticket-042 from `Todo` → `In Progress`; Session B (a parallel agent working a different lane pass, or a teammate on another clone) moves the *same* ticket from `Todo` → `Blocked` before pulling A's change. Whoever pushes second gets a conflict on one line, `<<<<<<< status: in-progress ||| status: blocked >>>>>>>`, and now either a human resolves it manually or an unattended agent has to.

A second, less obvious failure mode is **not a git conflict at all**: if the CLI does a naive read-modify-write (read file → parse YAML → mutate field → serialize → write) without file locking, two CLI invocations racing on the same file *on the same machine* (e.g., two subagents spawned in parallel within one Claude Code session, both told to advance different-but-related tickets, one of which touches a shared file) can silently produce a **lost update** before git is ever involved — the second write clobbers the first with no conflict marker, no warning, no commit-time visibility. "CLI is the single write path" (per PROJECT.md's Key Decisions) prevents *schema drift* — it does **not** by itself prevent this race, because a single logical write path can still be invoked concurrently by multiple processes. This is a real gap that must be explicitly assessed, not assumed away: the mitigation for schema drift (one validator) and the mitigation for write races (locking/atomicity) are different problems and both are needed.

**Why it happens:**
- Ticket status is the single field every lane transition touches, so it is the highest-contention line in the highest-contention file, by construction of the domain.
- Multiple parallel Claude sessions are an explicit design goal (per PROJECT.md Constraints: "Multiple parallel Claude sessions and multiple teammates write the same store"), which multiplies the odds of same-tick contention versus a single-writer tool.
- Developers assume "the CLI is the only writer" implies "writes are serialized," which is false unless the CLI enforces it (e.g., via an OS-level file lock, an atomic rename-based write, or a compare-and-swap on a revision field).

**How to avoid:**
- **At the CLI level (same-machine race):** use atomic writes (write to temp file, `rename()` over the target) *and* a lock (flock on the ticket file or a per-repo lock directory) around read-modify-write. This closes the silent-lost-update gap that "single write path" alone does not.
- **At the git level (cross-clone conflict):** keep the frontmatter schema field-per-line (never a nested one-line JSON blob) so git's line-based merge can auto-resolve edits that touch different fields — this is nearly free and should be a hard formatting rule enforced by the CLI's serializer, not a convention.
- **For the `status` line specifically**, since it's the guaranteed hot spot: consider making status transitions **idempotent and order-independent** where possible (e.g., a `history:` append-log of transitions inside the ticket file, with `status` *derived* from the last entry, rather than a single overwritten scalar) — this converts "two conflicting writes to one line" into "two appends to a list," which git handles as a clean auto-merge in the common case (new lines added by both sides) rather than a conflict, and it's a design choice available *without* abandoning plain markdown files.
- **Provide `jaira sync`/`jaira status --resolve` UX** for the conflict that *does* still occur: when a real content conflict happens (e.g., two different target statuses), the CLI should detect the git conflict markers in a ticket file and refuse to treat the file as valid until resolved, rather than silently picking one side.
- Document, explicitly, that git conflicts on ticket files are an expected and normal occurrence, not a defect — set user expectations rather than pretend they're solved.

**Warning signs:**
- Any bug report of the form "a ticket reverted to an earlier status after a pull."
- `git log -p` on `.jaira/tickets/*.md` showing frequent conflict-resolution commits.
- Two agent sessions running concurrently against the same worktree without any lock file present.

**Phase to address:**
Foundation phase for the atomic-write/locking mechanics (must exist before any concurrency is possible, i.e., before the agent pipeline or two-way TodoWrite sync ships). The append-log-vs-scalar status design decision should be made in the same phase, since retrofitting it after the schema is "the API" is a breaking migration.

---

### Pitfall 3: LLM-as-judge review lanes rubber-stamp, and a Definition of Done an agent can self-assess is not a control

**What goes wrong:**
The project's "secondary goal" is explicitly a cheap-implement / expensive-critique pipeline. Recent evaluation research on LLM-as-judge in production, multi-turn agent settings found judges recall well under a quarter of human-confirmed defects — in one batch, a judge flagged **zero** of 23 defects a human review confirmed. A related 2026 finding ("the Stability Trap") shows LLM judges can produce high *verdict agreement* with each other while independently having low *defect-detection accuracy* — meaning two review passes agreeing with each other is not evidence they're catching real problems, it's often evidence they share the same blind spot. Judges have also been observed to **hallucinate evidence to retroactively justify a "compliant" verdict** when asked to review against criteria. (Source: arXiv 2606.10315, arXiv 2601.11783 — MEDIUM confidence, small number of studies, but the failure pattern — judges systematically miss cross-turn/structural issues while catching surface issues — is consistent across independent papers.)

Applied to jAIra's design: an implementer agent marking its own Definition of Done "met" is the single worst-case version of this — it is not even an independent judge, it's the same context self-certifying. A DoD field that is free-text and agent-writable is trivially satisfiable by an agent that just restates the DoD as the outcome. A separate critique lane run by a "strong model" is better than self-certification, but the research above says it is *not* reliable on its own — critique models tend to catch local, surface-level problems (does the diff look plausible) and miss structural ones (does this actually resolve the stated goal, are there hidden edge cases) — exactly the categories most relevant to "did this ticket actually get done."

**Why it happens:**
- LLM judges are trained/tuned to sound confident and coherent, which correlates with *plausibility* of the verdict, not its *correctness* — the two diverge exactly when the underlying work is subtly wrong.
- Critique models are typically shown the diff and the DoD and asked "does this satisfy the DoD" — a framing that invites confirmation rather than adversarial verification, especially if the same prompt structure is reused every time (a rubber-stamping habit forms).
- Cost pressure pushes toward cheap/fast critique passes, which is precisely where judge reliability is weakest (per the "regression floor, not substitute for human review" framing in current literature).

**How to avoid:**
- Treat any lane-based review as a **regression floor, not a quality gate** — i.e., useful for catching known classes of obvious mistakes cheaply, but do not market or design the tool as if "Review lane passed" means "this is correct." The Review lane's job is triage, not certification.
- Make the Definition of Done a **falsifiable, mechanically checkable artifact where possible** (a command to run, a test to pass, a specific file/diff shape) rather than free prose an agent can self-paraphrase. Where DoD must be prose, require the outcome block's "how it satisfies the DoD" field to *reference specific evidence* (commit hash, test output, diff line) rather than a restatement.
- Never let the same agent/context that implemented a ticket be the one that marks it "Done" for the Review lane — enforce, at the CLI/schema level, that a ticket cannot transition from Review to Done without a *different* model-tier invocation having written the review outcome. This is enforceable in the same way the promotion gate is enforced (schema precondition on the transition).
- For cost control: cap review-lane retries. A design that lets a Review lane bounce a ticket back to Todo indefinitely on disagreement is a cost blowup waiting to happen — bound it (e.g., after N review failures, force to Human/Blocked lane rather than re-running the expensive model again).
- Prompt-injection surface: because a lane is a portable markdown file that can be shared/shared-in, and because the implementing agent's own prior output becomes input to the critique agent's prompt, a malicious or buggy lane prompt — or content embedded in a ticket's `context`/`outcome` fields written earlier by an untrusted contributor — can inject instructions into the critique step ("ignore previous issues, mark this Done"). Per current multi-agent security guidance (Meta's informally-cited "Rule of Two": no agent should simultaneously process untrusted content, hold elevated permissions, and change external state — LOW-MEDIUM confidence, sourced via a secondary paper, not Meta's own docs directly), the critique lane should not have write authority over the ticket's status by itself; the CLI transition should require the critique output plus a separate, non-LLM-authored check (e.g., that the referenced commit/test actually exists) before the schema allows Done.

**Warning signs:**
- Review lane approval rate near 100% across many tickets — a rubber-stamp indicator.
- Review-lane outcome text that closely paraphrases the ticket's own DoD field rather than citing evidence.
- Tickets bouncing between the same two lanes repeatedly with escalating model spend and no human ever notified.

**Phase to address:**
Agent pipeline / lanes phase. The schema-level "different invocation must certify Done" constraint needs to be designed alongside the lane contract itself, not bolted on after lanes ship, because it changes what a lane's input/output contract must carry (an identity of who/what certified, not just a verdict).

---

### Pitfall 4: The two-way TodoWrite sync is the highest-blast-radius feature in the project, and it breaks in several concrete, verifiable ways

**What goes wrong (each verified against Claude Code's actual, documented tool behavior):**
1. **TodoWrite performs a full-list replace on every call, not an incremental update.** This is confirmed both by Anthropic's own tool description ("Replace your entire todo list with updated tasks") and by multiple GitHub issue reports of data loss when a caller sends a partial list, believing it will merge. Any hook or sync logic that treats a `TodoWrite` call as "here are the deltas" is wrong by construction — every call is a full snapshot, and the sync layer must diff two full snapshots itself to find what changed. (HIGH confidence — official tool description + issue reports agree.)
2. **Duplicate ticket creation is the default outcome of naive sync**, precisely because of #1: if the sync hook treats every item in a TodoWrite call as "new," and Claude re-emits the same todo text on a later call (which it does routinely — TodoWrite is called repeatedly through a task, each time with the full current list, including unchanged items), a naive mirror-to-Backlog implementation creates a new ticket per call per item. The correct approach requires stable identity — matching todo items to existing tickets by content hash or an ID injected into the todo text, not by list position or literal-text equality (todo text often gets lightly reworded between calls).
3. **PostToolUse-hook-triggered loops are a documented, real failure class in Claude Code**, not a hypothetical: GitHub issue reports describe a hook firing after every Bash call, which itself triggers the agent to call Bash again, re-firing the hook — an amplifying loop, not a one-time bug. The same shape applies directly to a TodoWrite hook that itself calls a tool (e.g., a hook that shells out to `jaira` and that write, if visible to the agent as a side effect or if it re-triggers TodoWrite indirectly) — any hook whose side effect can be observed by the agent and cause it to re-emit a todo is a loop risk. (HIGH confidence — anthropics/claude-code issue #10205, #6674, dev.to writeup of a hook looping 25 times.)
4. **Hook latency is additive and synchronous by default.** Documented guidance states hooks over ~500ms per tool call make a session feel sluggish, and hook execution time is cumulative across all matched hooks for a given event — a TodoWrite-triggered hook that shells out to the `jaira` CLI, parses/writes a markdown file, and possibly does a git operation, on *every single TodoWrite call* (which happens many times per task), is exactly the kind of hook likely to cross that threshold if it isn't kept intentionally minimal (no git operations synchronously, no network, ideally no more than one file write).
5. **Hook failures can break the session, not just the sync.** A hook that exits non-zero or times out can block the tool call it's attached to (depending on hook type/config) — meaning a bug in the jAIra sync hook has the power to make ordinary TodoWrite calls (i.e., ordinary Claude Code task tracking, unrelated to jAIra) fail for the user. This raises the severity of any sync-hook bug from "board is stale" to "Claude Code itself misbehaves."
6. **Board/todo divergence is inevitable, not a corner case**, given #1–#3: any dropped hook invocation (session crash mid-call, hook timeout, git write failure) leaves the todo list and the board disagreeing, with no automatic reconciliation unless one is explicitly built.
7. **API-surface risk: TodoWrite itself may be superseded.** There is credible evidence (GitHub issue titled "TodoWrite tool broken in Claude Code — replaced by TaskCreate/TaskUpdate/TaskList", and VentureBeat coverage of a "Tasks" feature in Claude Code v2.1.16 that persists across context compaction) that Anthropic is actively moving away from the single-call-full-replace TodoWrite model toward a more granular Tasks API. Building deep, hard-coded sync logic against TodoWrite's specific full-replace quirk risks the integration breaking or becoming redundant on a Claude Code upgrade. (MEDIUM confidence — this is an active, evolving area as of the current Claude Code version; treat as a live risk, not a settled one.)

**Why it happens:**
TodoWrite was designed as an ephemeral, in-session scratchpad, not a durable, externally-synced data store — its full-replace semantics make sense for "Claude keeps its own todo list fresh" and become a liability the moment an external system tries to treat each call as an event stream.

**How to avoid:**
- Never treat a TodoWrite call as a delta. Always diff the new full list against the last-known full list (stored by the sync layer, not derived from the hook call alone) to compute adds/removes/updates.
- Assign stable identity at creation time: when a ticket is first mirrored from a todo, embed a short, grep-able ticket-id token in the mirrored todo's text (or maintain an out-of-band map file) so future TodoWrite calls can be matched to existing tickets by ID substring, not fuzzy text similarity.
- Keep the hook itself intentionally dumb and fast: write an event/queue entry synchronously (a single, small, atomic file append), and do the actual git-aware ticket read-modify-write **asynchronously or on next TUI/CLI invocation**, not inline in the hook. This bounds hook latency and means a git failure never blocks the user's Claude Code turn.
- Make the hook fail open, never closed: any exception inside the sync hook should log and no-op rather than return a nonzero/blocking exit that could stall or error the user's actual TodoWrite call. The sync being best-effort is an acceptable tradeoff; the sync breaking Claude Code is not.
- Build an explicit reconciliation path (`jaira sync --check` or similar) that a user can run to detect and fix board/todo drift, since silent eventual divergence is the expected steady state, not an edge case.
- Isolate the TodoWrite-specific integration behind an internal interface so that if/when Claude Code's Tasks API supersedes TodoWrite, the adapter is swapped rather than the whole sync design rebuilt.

**Warning signs:**
- Board shows duplicate tickets with near-identical titles.
- Users report Claude Code feeling slower specifically during tasks with many todo updates.
- A support report of "Claude Code hung" or "todo list stopped updating" that traces back to the jAIra hook.

**Phase to address:**
Capture/promotion-gate phase (per PROJECT.md's grouping) — but treat it as its own sub-phase with dedicated hardening time, separate from and after the one-way "TodoWrite → Backlog" capture, since two-way sync is strictly higher risk than one-way mirroring and should not ship in the same increment.

---

### Pitfall 5: Agent-maintained structured state drifts silently, and "one CLI schema enforcer" doesn't catch every failure mode

**What goes wrong:**
Several distinct sub-failures, each independently documented in adjacent domains (agentic coding, YAML tooling):
- **Silently malformed frontmatter.** YAML is a famously permissive format — a bad indent can reparse a list as a string, a missing quote can turn `no`/`yes`/`on`/`off` into booleans instead of the intended string values, and most YAML parsers do not error on these, they just parse to a different, wrong value. If the CLI's schema validator only checks "does this parse as YAML" rather than "does this match the exact expected shape and types," malformed-but-valid-YAML frontmatter passes silently.
- **Hallucinated ticket IDs.** An agent referencing `blocked-by: TICKET-047` when no such ticket exists (because it half-remembered an ID from earlier context, or invented a plausible-looking one) is a known class of LLM error wherever LLMs reference structured identifiers from memory rather than looking them up. The CLI is the single write path for *creating and mutating* tickets, but nothing stops an agent from writing a plausible-but-wrong reference unless the CLI actively validates that every `blocked-by` (and any other cross-reference) resolves to an existing ticket at write time.
- **Tickets abandoned mid-lane.** Documented, real Claude Code behavior: issue reports describe subagent work completing while the visible task state remains stuck `in_progress`, and separately, Claude prematurely declaring a multi-step plan complete after finishing only part of it. Applied to jAIra: a ticket can get stranded in "In Progress" indefinitely if the agent session ends (context exhaustion, user closes terminal, crash) between "started the work" and "called the CLI to advance the lane."
- **Agents stop updating the board as context gets long.** This is a restatement of the general "context loss is the biggest productivity killer in long sessions" pattern — the more turns since a ticket was picked up, the higher the odds the agent forgets to call the CLI transition at completion, especially if the transition call isn't the natural last action in its own plan.
- **Self-marking Done without meeting DoD** — covered in depth under Pitfall 3, but it is also a *state-integrity* problem, not just a review-quality problem: a wrongly-Done ticket pollutes the record other tickets' `blocked-by` depend on.

**Why it happens:**
A CLI as single write path guarantees *one code path validates every write*, which is real and valuable — it prevents the TUI, Claude, and any future integration from each inventing slightly different serialization and drifting apart. But it does not, by itself, guarantee: (a) semantic validity of cross-references, (b) that a transition actually happens (an agent can simply never call the CLI), or (c) that a claimed state (Done, DoD-met) is true. Schema enforcement at the write boundary is necessary but addresses a different layer than truthfulness and completeness of writes.

**How to avoid:**
- **Strict schema validation, not just "is it YAML."** Validate field types explicitly (status must be one of an enum, `blocked-by` must be a list of strings matching a known ID pattern, etc.) and reject-and-report rather than coerce on ambiguous values.
- **Referential integrity check on every write that touches `blocked-by` or any ID reference** — the CLI should look up the referenced ticket file and fail the write if it doesn't exist, rather than silently accepting an arbitrary string.
- **Timeout/staleness detection for "in progress" tickets** — a ticket sitting In Progress past some heuristic threshold (last modified time, or a session-liveness signal) is a strong candidate to surface on the board distinctly (e.g., a visual "stale" indicator) rather than trusting lane position alone as ground truth.
- **Make "advance the ticket" a natural, hard-to-skip last step of the agent's own plan**, not an optional courtesy call — e.g., the lane's I/O contract should structurally require the CLI transition call as the final output the agent must produce to be considered to have executed the lane at all (this is a prompt/skill design concern, addressed alongside the skill that teaches the CLI).
- **Don't rely on the agent to self-report DoD-met** — see Pitfall 3.

**Warning signs:**
- Tickets that sit In Progress with no corresponding recent commit.
- `blocked-by` references that don't resolve to any file in `.jaira/tickets/`.
- Frontmatter fields with unexpected types when inspected (e.g., a `commits` field that's a string instead of a list because an agent wrote it free-form).

**Phase to address:**
Foundation phase for schema strictness and referential integrity (these are CLI/validator concerns, cheapest to get right before other phases build on the schema). Staleness detection can land with the TUI phase, since it's primarily a board-rendering concern.

---

## Moderate Pitfalls

### Pitfall 6: Scope creep back toward paca is a mechanism, not a risk of carelessness

**What goes wrong:**
The project's own framing names the mechanism precisely: "one good idea, one small feature request at a time." Concretely, for this project, the plausible creep vectors are visible in its own Out of Scope list — every one of those items (custom fields, saved views, sprints, external sync) is the kind of thing a single real user request can make feel reasonable in isolation ("just let me add one custom field," "just let me filter by assignee," "just this one Jira webhook").

**Why it happens:**
Feature bloat happens because 80% of users need only 20% of a tool's features, but every user segment's 20% is different — so satisfying successive individual requests, each locally reasonable, sums to the exact bloated surface being escaped. This is a documented, general pattern in the feature-creep literature, not specific to this domain, but the project is unusually exposed to it because its explicit success metric ("is this smaller than paca?") has no automatic enforcement — it's a review-time judgment call, not a build-time constraint.

**How to avoid:**
- Convert "is this smaller than paca?" from a slogan into a checkable gate: maintain a running comparison (e.g., paca's feature list vs. jAIra's) and require any new feature proposal to explicitly name what it is *not* doing that paca does, at review/PR time.
- Treat every item currently in "Out of Scope" as requiring a deliberate, documented decision to move — not something that can happen implicitly via an accumulation of small unrelated PRs (e.g., "just add a `priority` field" is a custom field in a different name).
- Prefer schema *reservation* over schema *addition* for anything speculative — the project already does this correctly for `external:` (reserved but unbuilt); apply the same pattern discipline to any other future ask before writing code.
- Lane files being user-authored, portable, single-file markdown is itself a good structural guardrail: it pushes configurability out of the core tool and into a plugin-like artifact instead of core feature surface — protect this boundary; don't let core absorb logic that could instead be a shareable lane file.

**Warning signs:**
- A PR description that starts with "just a small addition" touching a field/table not in the current Active requirements.
- Any feature justified by "paca has this too, and it's not that big a deal" — that reasoning is exactly the trap.

**Phase to address:**
Every phase — this is a standing review discipline, not a one-time build task. Make it an explicit exit criterion checked at every phase boundary (the project's own `/gsd-transition` review already includes an Out-of-Scope audit — use it for this specifically).

---

### Pitfall 7: TUI terminal-compatibility traps, with WSL2-specific failure modes

**What goes wrong:**
- **Unicode/emoji column-width miscalculation is a real, currently-open class of bug** in TUI libraries, not a theoretical concern: Charm's lipgloss (a widely-used Go TUI styling library) had an open issue for emoji/Unicode width miscalculation causing layout misalignment, and Claude Code itself hit a real terminal-corruption bug from a wcwidth table not yet classifying a newer Unicode codepoint (Ghostty-specific, but illustrates that "the terminal and the library disagree about width" is an ongoing, unresolved-in-general problem, not a solved one). Any lane/ticket title or status glyph that includes an emoji (a plausible UX choice for a kanban board) risks column misalignment specifically because terminal emulators disagree with each other and with libraries about emoji/grapheme cell width.
- **WSL2-specific rendering issues are documented and recurring**, not a one-off: Windows Terminal + WSL2 has multiple open/historical issues around background-color bleed, truecolor detection ambiguity, and generally slower/sluggish rendering compared to native Linux terminals. Since WSL2 is stated as the primary dev environment for this project, any color-capability detection logic must not assume a clean truecolor signal — test against the actual WSL2 + Windows Terminal combination the developer uses, not just against a native Linux tty.
- **Resize handling** is a classic TUI trap in general (not WSL2-specific): board layouts computed once at startup and not recalculated on `SIGWINCH`/resize events break visibly the first time a user resizes their terminal — for a lane-based kanban board with a fixed number of columns, resize handling directly determines whether the tool degrades gracefully (fewer visible cards, scrollable lanes) or renders garbled/overlapping output.
- **File-watch storms causing render thrash**: chokidar/fs.watch-style watchers are known to fire multiple events for a single logical change (editors often write via temp-file + rename, which is 2+ raw fs events for what is semantically one ticket update), and watching a directory of many small ticket files during a git operation (checkout, merge, rebase) can produce a burst of events across many files nearly simultaneously — a naive "re-render on every fs event" TUI will visibly stutter or flicker during any git operation that touches `.jaira/tickets/`.

**Why it happens:**
Terminal emulators, TUI libraries, and Unicode standards each evolve independently and don't agree with each other on edge cases (emoji width, especially) — this is described directly by TUI authors as terminal emulators having "agreed to disagree" on Unicode interpretation. WSL2 specifically sits behind an extra virtualization/rendering layer (Windows Terminal → WSL2 → Linux tty) which is where the additional quirks originate.

**How to avoid:**
- Avoid emoji/wide-glyph status indicators in card titles or anywhere column width must stay exact; prefer plain ASCII/simple box-drawing characters for structural UI, and if color/glyph status indicators are used, size them defensively (reserve fixed-width cells, don't compute width from the glyph).
- Detect color capability via `$COLORTERM`/`$TERM`/`NO_COLOR` conventions and degrade gracefully to a reduced palette rather than assuming truecolor; test explicitly in the actual WSL2 + Windows Terminal setup used for development, not just a CI Linux runner (CI terminals and dev WSL2 terminals can report different capabilities).
- Implement resize handling as a first-class requirement from the first TUI phase, not an afterthought — recompute layout on resize events, and define an explicit minimum-terminal-width degrade path (e.g., collapse to fewer visible lanes with horizontal scroll) rather than letting layout break unbounded.
- Debounce the file watcher: coalesce bursts of fs events (e.g., 100–300ms debounce window) into a single re-render, and ignore watcher events entirely during a CLI-initiated write (the CLI already knows it just wrote a file; it doesn't need the watcher to tell it) to avoid the watcher and the CLI's own state update racing to re-render twice.

**Warning signs:**
- Column borders visibly shift by one character when a card title contains certain glyphs.
- Visible flicker/stutter specifically during `git pull`/`git checkout` inside a repo with an open jAIra session.
- Bug reports that reproduce only on Windows/WSL2 and not on native Linux.

**Phase to address:**
TUI/board phase for width and resize handling (must be right before the board is trusted as the primary interface); file-watch debounce belongs in the same phase since it's a rendering-pipeline concern, not a data-layer one.

---

### Pitfall 8: Human-in-the-loop lane failure modes — parked forever, missed, or wrongly bypassed

**What goes wrong:**
General human-in-the-loop literature converges on the same handful of failure modes, all directly applicable here: approval/question requests get buried and nobody notices; the person who should answer is unavailable and there's no escalation; the agent (or a naive implementation) simulates "waiting" without a real human ever seeing the question, then proceeds anyway; or, conversely, an agent blocks and waits for input on something it could safely have inferred or proceeded past, wasting a lane slot indefinitely. Applied to jAIra specifically: a ticket in the Human lane with a question attached is only as good as the odds someone opens the TUI *and* notices that specific card *and* answers it — there is no push notification path by design (no server, terminal-only), so the entire mechanism relies on the user voluntarily glancing at the board.

**Why it happens:**
Human-in-the-loop patterns generally assume some notification channel (email, chat ping, a durable workflow engine with SLA timers) — jAIra's constraints (no server, terminal-only, git as the only sync layer) rule out essentially all of the standard mitigations (push notifications, escalation timers backed by a service). The Human lane is pull-based by construction, and pull-based human review is the exact pattern the HITL literature identifies as prone to silent starvation.

**How to avoid:**
- Make the Human lane maximally *visible*, not just present — since push notification is out of scope, the fallback is to make the TUI impossible to open without seeing that something is waiting (e.g., a persistent, unmissable summary line/badge at the very top of the board — "2 tickets waiting on you" — rather than requiring the user to scroll to the Human lane column to notice).
- Consider a lightweight local-only nudge that doesn't require a server: e.g., a shell prompt hook, a `jaira status` one-liner a user can add to their shell prompt/tmux status bar, or a note Claude itself surfaces in-session ("2 tickets are waiting on you in the Human lane") the next time a session starts in that repo — this uses the existing in-session channel rather than inventing infrastructure.
- Define, explicitly, the criteria for when an agent should land a ticket in Human vs. proceed with a reasonable default — an under-specified DoD should not automatically mean "ask a human," but a genuinely ambiguous or destructive decision should. This needs to be a documented judgment rule in the lane prompt/skill, since without one, agents will drift toward either always asking (annoying, defeats automation) or never asking (dangerous).
- Give Human-lane tickets an explicit "asked at" timestamp so staleness is visible on the card itself (e.g., "waiting 6 days") rather than only inferable from git history.

**Warning signs:**
- Tickets sitting in Human lane with no interaction for many sessions.
- Users reporting they "didn't know" a question was pending until much later.
- Agents proceeding past decisions that, on later review, should have gone to Human (a sign the ask/don't-ask heuristic in the lane prompt needs tightening).

**Phase to address:**
Agent pipeline / lanes phase, specifically when the Human lane's contract is defined — the visibility/nudge mechanism should ship in the same phase as the Human lane itself, not as a later polish pass, since an invisible Human lane is arguably worse than no Human lane (work silently stalls instead of visibly failing).

---

## Technical Debt Patterns

| Shortcut | Immediate Benefit | Long-term Cost | When Acceptable |
|----------|-------------------|-----------------|------------------|
| Naive read-modify-write in CLI without file locking | Faster to ship v1 | Silent lost updates under any concurrency (parallel agents, parallel CLI invocations) | Never for the ticket-store CLI, even in MVP — this is core-integrity, not polish |
| TodoWrite sync hook does synchronous git writes inline | Simpler code path | Hook latency risk + hook-failure-blocks-session risk (documented Claude Code behavior) | Never — always defer heavy work out of the synchronous hook path |
| Frontmatter status as a single scalar overwritten in place | Simplest schema | Every concurrent status change is a guaranteed same-line git conflict | Acceptable only if append-log-of-transitions design is deliberately deferred with a documented migration plan, not silently skipped |
| Loose YAML validation ("if it parses, accept it") | Faster CLI implementation | Silently malformed/mistyped fields accepted (YAML's implicit type coercion pitfalls) | Never past MVP — cheap to fix early, expensive once real repos have malformed data in history |
| Skipping resize/SIGWINCH handling in early TUI builds | Faster initial TUI demo | Visibly broken layout on first real-world resize; erodes trust in "lightweight and fast" positioning | Acceptable only for a throwaway prototype, not for anything a teammate will actually run |

## Integration Gotchas

| Integration | Common Mistake | Correct Approach |
|-------------|------------------|-------------------|
| Claude Code TodoWrite | Treating each call as a delta/event | Diff full-list snapshots yourself; TodoWrite always sends the complete list (confirmed by official tool description) |
| Claude Code hooks (PostToolUse/Stop) | Doing heavy/blocking work synchronously inside the hook | Keep hook itself under ~200-500ms; queue heavy work for async processing; fail-open on error |
| git (as sync layer) | Assuming single-CLI-writer prevents all races | Add explicit file locking/atomic writes in the CLI regardless — single write path solves schema drift, not concurrency races |
| Custom lane markdown files (user/community-authored, portable) | Treating lane prompt content as fully trusted just because it's local | Sanitize/limit what a lane prompt can cause the implementing/critique agent to do — a shared lane file is an injectable trust boundary |
| Terminal emulator capability detection | Assuming truecolor/Unicode support uniformly | Detect via `$TERM`/`$COLORTERM`/`NO_COLOR`, test explicitly against WSL2 + Windows Terminal, degrade gracefully |

## Performance Traps

| Trap | Symptoms | Prevention | When It Breaks |
|------|----------|------------|-----------------|
| Full-repo file watch with no debounce | Board flickers/stutters during git operations | Debounce fs events (100-300ms coalescing window); ignore events caused by the CLI's own writes | Any git operation (checkout/merge/rebase) touching many ticket files at once |
| Re-render entire board on every fs event | Sluggish TUI as ticket count grows | Diff-based re-render (only redraw changed cards/lanes) | Noticeable once ticket count is in the low hundreds per repo |
| Synchronous git operations inside the TodoWrite sync hook | Claude Code turns feel slow specifically during multi-todo tasks | Defer git writes out of the hot hook path | As soon as a task generates more than a handful of TodoWrite calls, which is routine |
| Unbounded review-lane retry loop | Cost blowup from repeatedly re-running the expensive critique model | Cap retries; force to Human/Blocked lane after N failures | Any ticket where implement/critique disagree persistently |

## Security Mistakes

| Mistake | Risk | Prevention |
|---------|------|------------|
| Critique lane content or a shared lane-file prompt treated as fully trusted input to the next agent step | Prompt injection causing a lane to falsely certify Done, or causing unintended tool calls | Don't let the critique/review step alone hold write authority over ticket status; require a non-LLM-verifiable signal (existing commit/test) alongside the LLM verdict before allowing a Done transition |
| CLI trusts `blocked-by`/cross-reference fields without validating they resolve to real tickets | Hallucinated IDs corrupt dependency graph silently | Validate referential integrity on every write that includes a reference field |
| Ticket `context`/`outcome` fields rendered/executed as trusted text by a later agent step | Injected instructions embedded in earlier free-text fields could influence a later lane's agent | Treat all stored ticket text as data, not instructions, when constructing later lane prompts — clearly delimit/quote it rather than splicing it inline unescaped |

## UX Pitfalls

| Pitfall | User Impact | Better Approach |
|---------|--------------|-------------------|
| Human lane is present but visually indistinguishable from other lanes at a glance | Users miss pending questions for days | Persistent unmissable badge/summary at the top of the board when anything is waiting on the human |
| Board and TodoWrite silently diverge with no way to tell | User loses trust that either list is accurate | Explicit `jaira sync --check`/reconciliation command, and a visible "last synced" indicator |
| Git conflicts on ticket files surfaced only as raw conflict markers in a markdown file | User has to know git internals to fix a ticket | CLI detects and calls out conflicted ticket files explicitly with a guided resolution path |
| Unknown/custom lanes rendering as inert passthrough columns with no explanation | Confusing "why is this column here and doing nothing" | Passthrough columns should visibly label themselves as unrecognized/read-only, not just render silently |

## "Looks Done But Isn't" Checklist

- [ ] **Ticket store CLI:** Often missing file locking/atomic writes — verify concurrent invocations (two CLI calls racing on the same ticket) don't lose an update.
- [ ] **TodoWrite two-way sync:** Often missing stable identity matching — verify that calling TodoWrite twice with slightly reworded but semantically-same items does not create a duplicate ticket.
- [ ] **Review/critique lane:** Often missing an independent, non-LLM-verifiable signal — verify a ticket cannot reach Done purely on the strength of an LLM's self-reported verdict with no evidence reference.
- [ ] **Human lane:** Often missing a visibility mechanism beyond "it's a column on the board" — verify a user who hasn't opened the TUI in days is still surfaced the pending question somehow (in-session nudge, badge, etc.).
- [ ] **TUI resize handling:** Often missing SIGWINCH/resize recompute — verify resizing the terminal mid-session doesn't garble the layout.
- [ ] **Frontmatter schema validation:** Often missing type-strictness — verify a malformed-but-YAML-valid field (wrong type, unexpected value) is rejected, not silently accepted.
- [ ] **Merge-conflict handling:** Often missing detection of conflict markers left in a ticket file — verify the CLI refuses to treat a conflicted file as valid rather than silently reading garbage.

## Recovery Strategies

| Pitfall | Recovery Cost | Recovery Steps |
|---------|-----------------|-----------------|
| Board/TodoWrite divergence discovered late | LOW | Run reconciliation command; rebuild todo mirror from current ticket store as source of truth |
| Ticket stuck In Progress with abandoned agent session | LOW | Staleness heuristic flags it on the board; user manually reassigns/reverts to Todo |
| Git conflict on a ticket's status field | LOW-MEDIUM | CLI-guided conflict resolution: show both candidate statuses, let user or a fresh agent pass pick one |
| Schema needs to change after tickets already exist in the wild (e.g., status→append-log migration) | HIGH | Requires a versioned migration path baked into the CLI from day one; retrofitting is a breaking change across every existing repo's ticket files |
| Rubber-stamped Done ticket discovered to be wrong later | MEDIUM | Reopen via CLI to a prior lane; the outcome block's evidence references (commit hash) make it possible to audit what actually happened |

## Pitfall-to-Phase Mapping

| Pitfall | Prevention Phase | Verification |
|---------|-------------------|----------------|
| File-based tracker adoption failure (visibility, tooling drift) | Foundation (ticket store + CLI) — architectural framing, not a single fix | Schema is versioned and documented independently of the tool; dogfood test for "do teammates actually open the TUI" |
| Merge conflicts on status field / write races | Foundation (ticket store + CLI) | Concurrency test: two parallel CLI invocations on the same ticket; two clones both moving the same ticket before pulling |
| LLM self-review / rubber-stamp Done | Agent pipeline / lanes phase | Schema enforces different-invocation-certifies-Done; audit a sample of Done tickets against actual evidence |
| Two-way TodoWrite sync hazards | Capture/promotion-gate phase (sync as its own hardened sub-phase) | Stress test: rapid repeated TodoWrite calls with reworded items; hook latency measured under 200-500ms; hook failure doesn't block user's turn |
| Structured-state drift (hallucinated IDs, malformed frontmatter, abandoned mid-lane) | Foundation (schema validation) + TUI phase (staleness surfacing) | Fuzz/malformed-YAML test suite against the validator; staleness indicator visible on board |
| Scope creep toward paca | Every phase (standing review discipline) | Phase-transition checklist explicitly re-audits Out of Scope list |
| TUI terminal/Unicode/WSL2 traps | TUI/board phase | Manual test matrix: WSL2+Windows Terminal, native Linux tty, minimal-capability terminal, resize mid-session |
| Human-in-the-loop starvation | Agent pipeline / lanes phase (Human lane's own contract) | Verify a pending Human-lane question is surfaced without requiring the user to have already opened the TUI |

## Sources

- [Current State of the Distributed Issue Tracking](https://matej.ceplovi.cz/blog/current-state-of-the-distributed-issue-tracking.html) — survey of dead/abandoned DIT tools (ticgit, Bugs Everywhere, gitissius); MEDIUM confidence, single-author survey but corroborated by primary project status pages
- [Fossil: Bug-Tracking In Fossil (bugtheory.wiki)](https://fossil-scm.org/home/doc/tip/www/bugtheory.wiki) — read directly; HIGH confidence, official design rationale for append-only ticket artifacts vs. mutable working-tree files
- [git-bug GitHub repository and issues](https://github.com/git-bug/git-bug) — bridge/sync bugs and third-party integration friction, MEDIUM confidence
- Hacker News threads on git-bug (2018, 2022, 2024 era) — management/visibility objection pattern, MEDIUM confidence (community commentary, not primary data)
- ticgit GitHub wiki — project marked dead by original maintainers, referenced via search; MEDIUM confidence
- [anthropics/claude-code issue #10205 — infinite loop with hooks enabled](https://github.com/anthropics/claude-code/issues/10205) — HIGH confidence, primary issue tracker
- [anthropics/claude-code issue #6674 — hook navigation loop](https://github.com/anthropics/claude-code/issues/6674) — HIGH confidence
- ["I Misconfigured One Claude Code Hook and It Ran 25 Times in a Loop" — dev.to](https://dev.to/ji_ai/writing-a-claude-code-book-with-claude-code-when-posttooluse-hooks-loop-25-times-4h46) — MEDIUM confidence, personal account but mechanism matches official issue reports
- [anthropics/claude-code issue #2250 — TodoWrite overwrites entire list](https://github.com/anthropics/claude-code/issues/2250) and [issue #1824](https://github.com/anthropics/claude-code/issues/1824) — HIGH confidence, confirmed against official tool description
- [obra/superpowers issue #1518 — TodoWrite replaced by TaskCreate/TaskUpdate/TaskList](https://github.com/obra/superpowers/issues/1518) and [VentureBeat coverage of Claude Code Tasks](https://venturebeat.com/orchestration/claude-codes-tasks-update-lets-agents-work-longer-and-coordinate-across) — MEDIUM confidence, evolving/live area
- [anthropics/claude-code issue #59962](https://github.com/anthropics/claude-code/issues/59962) and [issue #6159](https://github.com/anthropics/claude-code/issues/6159) — agent stops mid-task, stale in_progress state; HIGH confidence, primary issue tracker
- [Catching One in Five: LLM-as-Judge Blind Spots in Production Multi-Turn Transaction Agents (arXiv 2606.10315)](https://arxiv.org/html/2606.10315) — MEDIUM confidence, single study with specific numbers, but consistent with broader LLM-judge-reliability literature
- [The Stability Trap: Evaluating the Reliability of LLM-Based Instruction Adherence Auditing (arXiv 2601.11783)](https://arxiv.org/pdf/2601.11783) — MEDIUM confidence
- Multi-agent prompt injection defense literature (arXiv 2509.14285 and related 2026 papers on indirect injection in agent pipelines) — MEDIUM confidence; "Rule of Two" attribution to Meta is via secondary citation, not Meta's own documentation, treat as LOW-MEDIUM confidence specifically for that attribution
- [charmbracelet/lipgloss issue #562 — emoji/Unicode width layout misalignment](https://github.com/charmbracelet/lipgloss/issues/562) — HIGH confidence, primary issue tracker
- [anthropics/claude-code issue #69093 — TUI corruption from unclassified Unicode 16.0 codepoint](https://github.com/anthropics/claude-code/issues/69093) — HIGH confidence
- Microsoft/terminal and Microsoft/WSL GitHub issues on color rendering and WSL2 performance (issues #15347, #10159, #76, warpdotdev/warp #11741) — MEDIUM confidence, multiple independent reports converge on the same class of quirks
- General human-in-the-loop workflow pattern literature (Temporal, Cloudflare Agents docs, Orkes) — MEDIUM confidence, standard patterns for approval-request starvation and escalation, generalized to jAIra's no-server constraint
- General feature-creep/bloat literature (Wikipedia, Mambo, Designli, Shopify partner blog) — LOW-MEDIUM confidence, generic product-management sources, applied here as domain-general context rather than jAIra-specific evidence

---
*Pitfalls research for: git-committed markdown kanban board with agent-pipeline lanes driven by Claude Code*
*Researched: 2026-08-11*
