---
phase: quick-260819-hco
plan: 01
type: execute
wave: 1
depends_on: []
files_modified:
  - .claude/skills/jaira/SKILL.md
  - docs/AGENTS.md
autonomous: true
requirements: [QUICK-260819-hco]

must_haves:
  truths:
    - "An agent reading SKILL.md learns to search the board for overlapping tickets before running `jaira create`."
    - "An agent that finds an existing ticket whose goal/DoD/context decided something different STOPS and asks the user, naming both handles and quoting the contradiction."
    - "The three ways forward are spelled out: honor the existing decision, supersede deliberately with `--follows`, or drop the new ticket."
    - "A pure duplicate is handled by the same check: point at the existing ticket instead of creating a twin."
    - "The same rule reaches non-Claude agents via docs/AGENTS.md — it is not Claude-only."
    - "Every command quoted in the new text exists in the real CLI with those flags."
  artifacts:
    - path: ".claude/skills/jaira/SKILL.md"
      provides: "The pre-create conflict check, in the skill's second-person voice"
      contains: "jaira list -q"
    - path: "docs/AGENTS.md"
      provides: "The same rule for any shell-capable agent"
      contains: "jaira list -q"
  key_links:
    - from: ".claude/skills/jaira/SKILL.md"
      to: "jaira list -q / jaira show --json"
      via: "command block in the new section"
      pattern: "jaira list -q"
    - from: ".claude/skills/jaira/SKILL.md"
      to: "--follows"
      via: "the deliberate-supersession way forward"
      pattern: "--follows"
    - from: "docs/AGENTS.md"
      to: "jaira show <id> --json"
      via: "command block in the new section"
      pattern: "jaira show"
---

<objective>
Teach agents to look before they create: search the board for tickets whose scope
overlaps the one about to be written, and stop to ask the user when an existing
ticket already decided something different.

Purpose: Berk's scenario — "I let Claude create a ticket, but an existing ticket
has decided something different, then the user must be asked again." A silently
created contradiction leaves two tickets pulling opposite ways and no record of
which one is current. That is the exact failure the board exists to prevent.

Output: One new short section in `.claude/skills/jaira/SKILL.md` and its
counterpart in `docs/AGENTS.md`. One commit. No Go code.
</objective>

<execution_context>
@$HOME/.claude/get-shit-done/workflows/execute-plan.md
@$HOME/.claude/get-shit-done/templates/summary.md
</execution_context>

<context>
@.claude/skills/jaira/SKILL.md
@docs/AGENTS.md

**Locked by the orchestrator with the user this session — do not revisit:**

1. This is a docs/skill change only. Semantic contradiction ("an existing ticket
   decided something different") is a judgement only the LLM layer can make. The
   CLI already provides everything needed. No new commands, no new flags, no Go.
2. Scope is **creation only**. Do not add rules for `set`, `move` or edits.
3. The check must not read as bureaucracy on an empty board: one `list` call is
   cheap, and the ask-the-user step fires only on an actual contradiction.
4. Do not touch Go code, lane files, or built-in lane prompt text — changing
   built-in lane text makes every existing project copy warn.

**CLI surface already verified this session** (`grep -n` on the source and both
`--help` commands actually run):

- `internal/cli/tickets.go:415` — `f.StringVarP(&query, "query", "q", "", ...)`,
  so `jaira list -q <substring>` exists.
- `internal/cli/tickets.go:420` — `func matches` iterates
  `t.ID, t.Title, t.Goal, t.Context, t.DoD, t.Assignee, t.Status`, so `-q`
  searches context, definition-of-done, assignee and status too. The `--help`
  one-liner says only "title, goal and id" and understates this: do not repeat
  the help text's narrower claim in the new prose, and do not "fix" the help
  string — that is Go code and out of scope.
- `internal/cli/tickets.go:970` — `ticketJSON` emits `goal` (988), `context`
  (989), `definition_of_done` (990), plus `dod_items` and `follows`, so
  `jaira list --json` returns enough to judge a small board in one call.
- `--json` is a global flag (confirmed in both `list --help` and `show --help`
  output), so `jaira list -q … --json` and `jaira show <handle> --json` are valid.
- `--follows <handle>` on `create` already exists and already resolves-or-exits-5
  (SKILL.md lines 42-54). The supersession path reuses it; introduce nothing new.

**Note for the executor:** this repository has no `.jaira/` board of its own
(`ls -d .jaira` → no such file). Do not add a `jaira validate` step to any
verification; it has no board to validate here.
</context>

<tasks>

<task type="auto">
  <name>Task 1: Confirm the CLI surface the new text will quote</name>
  <files>(no files modified — verification only)</files>
  <action>Before writing a word of prose, confirm the flags the new section will teach actually exist, so the skill never teaches a flag that is not there. Run `go run ./cmd/jaira list --help` and `go run ./cmd/jaira show --help` and confirm: `list` carries `-q, --query`; `--json` is available on both (it is a global flag). Read internal/cli/tickets.go around `func matches` (line 420) to confirm which fields `-q` actually searches, so the new prose describes the real behavior rather than the help string's narrower summary. If any of this disagrees with the Context section above, STOP and report the discrepancy rather than writing prose around it.</action>
  <verify>
    <automated>go run ./cmd/jaira list --help 2>&1 | grep -q -- '-q, --query' && go run ./cmd/jaira show --help 2>&1 | grep -q -- '--json'</automated>
  </verify>
  <done>Both help outputs seen; `-q, --query` present on `list`; `--json` present on both; the field list `-q` searches is known from the source.</done>
</task>

<task type="auto">
  <name>Task 2: Land the pre-create conflict check in SKILL.md and docs/AGENTS.md</name>
  <files>.claude/skills/jaira/SKILL.md, docs/AGENTS.md</files>
  <action>
Write one new short section into each file. Same rule, each file's own voice. One commit for both, because it is one rule landing in two places.

**SKILL.md** — insert a new top-level `##` section immediately after the `## When to create tickets` section ends (after the paragraph about a terminal lane refusing a ticket whose plan is unfinished) and immediately before `## Writing a good ticket from a request`. Do not restructure the surrounding sections. Voice: second person, terse, example-driven, one real command block — match the file exactly.

Substance, in this order:

- Before `jaira create`, look for what the board already says about this. Search the key terms of the ticket you are about to write with `jaira list -q <term>`, one call per term that matters; on a small board `jaira list --json` returns everything, goal, context and definition of done included, so one call is enough. Note that `-q` matches over title, goal, context, definition of done, assignee and status — not the title alone.
- Read any hit that looks close with `jaira show <handle> --json` before judging it.
- The thing you are looking for is not a related topic — related tickets are normal and fine. It is a **contradiction**: an existing ticket whose goal, definition of done or context already decided the same question the other way.
- When you find one, **stop and ask the user**. Name both handles, quote the specific line that contradicts, and give them the ways forward: adjust the new ticket to honor the existing decision, create it anyway as a deliberate supersession — then pass `--follows <handle>` so the chain of why survives — or drop it.
- Never create over a contradiction silently. Two tickets deciding a question opposite ways, with nothing recording which is current, is the failure this board exists to prevent.
- A pure duplicate is the easy case of the same check: the work is already captured, so point at that ticket instead of creating a twin.
- Say plainly that on an empty or tiny board this costs one command and the question never fires — so it does not read as ceremony.

Include one bash command block (the file's established style) showing the search and the read: a `jaira list -q "…"` line and a `jaira show <handle> --json` line, with the same kind of short trailing comments the rest of the file uses. Do not invent flags. Keep the whole section to roughly the size of a neighbouring section — this is one step, not an essay.

**docs/AGENTS.md** — insert a new `##` section titled for creating a ticket, placed after the `## The loop` section's closing paragraphs (the one about the specified zone and the promotion gate) and before `## Leaving a trail`. Match this file's register: it explains the shape to any shell-capable agent, prose paragraphs plus a small command block, slightly more explanatory and less second-person-imperative than SKILL.md. Carry the same substance in fewer words: search before `create`, read the close hits, and when an existing ticket already decided the question differently, stop and put it to the user with both ids and the contradiction quoted — adjust, supersede with `--follows`, or drop. Say that nothing in the binary enforces this, so it is on the agent — do NOT add it to the `## What an agent cannot do` list, which is exclusively about gate-enforced refusals that exit 3.

Do not touch anything else in either file. No Go code, no lane files, no built-in lane text.
  </action>
  <verify>
    <automated>grep -q 'jaira list -q' .claude/skills/jaira/SKILL.md && grep -q 'jaira show' .claude/skills/jaira/SKILL.md && grep -q -- '--follows' .claude/skills/jaira/SKILL.md && grep -q 'jaira list -q' docs/AGENTS.md && grep -q -- '--follows' docs/AGENTS.md && [ "$(git diff --name-only | grep -cv -e '^\.claude/skills/jaira/SKILL\.md$' -e '^docs/AGENTS\.md$')" = 0 ]</automated>
    <human-check>Read both new sections top to bottom. Does each sound like the file it is in rather than like a policy memo? Is the rule one step, not an essay?</human-check>
  </verify>
  <done>Both files carry the rule; every command quoted exists per Task 1; `--follows` named as the supersession path in both; the working tree diff contains no file other than those two; one commit.</done>
</task>

</tasks>

<threat_model>
## Trust Boundaries

| Boundary | Description |
|----------|-------------|
| agent → ticket store | An agent writes tickets from user prose; this change adds a read-before-write step, no new write path |

## STRIDE Threat Register

| Threat ID | Category | Component | Disposition | Mitigation Plan |
|-----------|----------|-----------|-------------|-----------------|
| T-hco-01 | Information disclosure | new prose quoting CLI flags | mitigate | Task 1 verifies every quoted flag against `--help` before it is written, so the skill cannot teach a non-existent surface |
| T-hco-02 | Tampering | built-in lane text / Go code | mitigate | Task 2 verify asserts the working-tree diff contains no file other than the two markdown files |
| T-hco-03 | Repudiation | silent contradictory ticket | mitigate | This change is the mitigation: `--follows` records a deliberate supersession so the chain of why survives |
| T-hco-SC | Tampering | package installs | accept | No package installs in this change — docs only, no dependency touched |
</threat_model>

<verification>
- `go run ./cmd/jaira list --help` shows `-q, --query`; `go run ./cmd/jaira show --help` shows `--json`
- `git diff --name-only` lists exactly `.claude/skills/jaira/SKILL.md` and `docs/AGENTS.md` — nothing else
- Both new sections name `jaira list -q`, `jaira show <handle> --json`, and `--follows`
- SKILL.md's new section sits between `## When to create tickets` and `## Writing a good ticket from a request`
- docs/AGENTS.md's new section sits between `## The loop` and `## Leaving a trail`, and the `## What an agent cannot do` list is unchanged
</verification>

<success_criteria>
An agent handed a request that contradicts an existing ticket searches the board
first, finds it, and asks the user — naming both handles and quoting the
contradiction — instead of creating a second ticket that decides the question the
other way. The rule reads as one cheap step in both files, and no flag it teaches
is invented.
</success_criteria>

<output>
Create `.planning/quick/260819-hco-before-creating-a-ticket-agents-must-che/260819-hco-SUMMARY.md` when done
</output>
