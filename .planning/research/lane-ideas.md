# jAIra lane ideas — derived from popular prompt/agent collections

Research method: web search only (no repo cloning, no star counts pulled via API — GitHub
search snippets don't reliably surface star counts, so popularity claims below are qualitative
and flagged where unverified). Existing lanes (not reproposed): backlog, brainstorm, todo,
pre-process, in-progress, critique, optimize, human, review, signoff, done, blocked.

Every idea below is deliberately something that is NOT already covered by the existing lane
list. Where an idea is close to an existing lane (e.g. refactor-smell-scan vs. `optimize`,
pr-description-writer vs. `review`), the distinction is called out explicitly.

---

## Theme: Quality gates

### 1. `lint-gate` — Static Analysis Gate
One sentence: runs linters/static analyzers (golangci-lint, staticcheck, etc.) on the diff and
blocks progression on new violations, without an LLM judging style.
Pipeline position: after `in-progress`, before `critique` — cheap, deterministic, so it should
fail fast before spending a strong-tier critique pass on code that doesn't even lint.
Input fields: `diff`, `outcome-what`.
Output fields: `lint-report` (pass/fail + violation list), `lint-status`.
Model tier: cheap (or no LLM at all — could just shell out to the linter and only invoke a
cheap model to summarize failures into a human-readable note).
Source: general pattern from AI-code-review-tool roundups (SonarQube/CodeClimate/ESLint/RuboCop
described as the deterministic pass that runs before/alongside AI review) —
[Code Smells in Practice](https://kodus.io/en/code-smells/) (vendor blog, popularity not
independently verified) and [12 Best Code Smell Detection Tools in 2026](https://dev.to/rahulxsingh/12-best-code-smell-detection-tools-in-2026-complete-guide-c76)
(community post, popularity unverified).
Draft prompt:
```
Run the project's configured linter/static analyzer against the changed files only.
Report every new violation introduced by this diff (ignore pre-existing violations in
untouched lines). Do not attempt to fix anything — only report. Output a pass/fail status
and, on fail, a list of file:line + rule + message. Keep the summary under 15 lines.
```

### 2. `refactor-smell-scan` — Structural Smell Review
One sentence: flags structural code smells (long method, god class, high cyclomatic
complexity, feature envy) in the new/changed code, independent of dead code or duplication.
Pipeline position: after `in-progress`, parallel to or just before `critique`.
Distinction from `optimize`: `optimize` removes code that already exists elsewhere or that
nobody calls (redundancy/waste); this lane flags code that *works* but is structurally hard to
maintain (complexity), which `optimize`'s charter does not mention.
Input fields: `diff`, `plan`.
Output fields: `smell-report` (list of smells + locations), `refactor-suggestions`.
Model tier: strong (judging "is this method doing too much" needs real code comprehension, not
just pattern matching).
Source: [VoltAgent/awesome-claude-code-subagents — refactoring-specialist.md](https://github.com/VoltAgent/awesome-claude-code-subagents/blob/main/categories/06-developer-experience/refactoring-specialist.md)
— a community-maintained Claude Code subagent collection; adoption/stars not independently
verified in this search session.
Draft prompt:
```
Review only the diff for structural code smells: long methods, god objects/classes, deep
nesting, high cyclomatic complexity, feature envy, primitive obsession. Do not comment on
naming style or formatting (the linter already does that). For each smell found, name the
smell, cite file:line, and give a one-line reason it will cost more to maintain later. If
nothing rises above triviality, say so plainly — do not manufacture findings.
```

---

## Theme: Security

### 3. `threat-model` — STRIDE-style Threat Model
One sentence: before implementation starts, works out what could go wrong from a security
angle (spoofing, tampering, repudiation, info disclosure, DoS, elevation of privilege) for the
planned change.
Pipeline position: after `pre-process` (which produces the plan), before `in-progress` — only
for tickets whose plan touches auth, data handling, external input, or new endpoints/tools.
Input fields: `goal`, `plan`.
Output fields: `threat-model` (per-STRIDE-category findings), `mitigations`.
Model tier: strong — threat modeling requires reasoning about attacker intent, not lookup.
Source: [ASTRIDE: A Security Threat Modeling Platform for Agentic-AI Applications (arXiv)](https://arxiv.org/pdf/2512.04785)
— academic paper, describes STRIDE extended for agentic/LLM systems; and
[fuzzinglabs.com — AI-Driven Threat Modeling](https://fuzzinglabs.com/ai-threat-modeling-arrows/)
(vendor blog, describes per-STRIDE-category LLM analyzer pattern; popularity/adoption not
independently verified).
Draft prompt:
```
Given the plan for this ticket, run a STRIDE pass: for each of Spoofing, Tampering,
Repudiation, Information Disclosure, Denial of Service, and Elevation of Privilege, ask
whether this change opens a new avenue for that threat. Skip categories that plainly do not
apply and say why in one line. For each real finding, state the attack, its impact, and a
concrete mitigation the implementer should build in — not a vague "be careful."
```

### 4. `security-review` — Dedicated Security Review
One sentence: a focused OWASP-Top-10 / secrets / injection audit of the diff, separate from the
general "is this the right implementation" `critique` lane and the general-purpose `review`
lane.
Pipeline position: after `in-progress`, before `critique` (or parallel to `lint-gate`) for
tickets that touch auth, input parsing, file I/O, or shell-out code — relevant to jAIra's own
`os/exec` git shell-outs.
Input fields: `diff`.
Output fields: `security-findings` (severity-tagged list), `security-status`.
Model tier: strong.
Source: [github/awesome-copilot — se-security-reviewer.agent.md](https://github.com/github/awesome-copilot/blob/main/agents/se-security-reviewer.agent.md)
— official GitHub org repo of Copilot agent/prompt definitions, high-confidence adoption signal
since it's GitHub's own curated collection; and [HeadyZhang/agent-audit](https://github.com/HeadyZhang/agent-audit)
(51 rules mapped to OWASP Agentic Top 10 2026, community project, stars not verified).
Draft prompt:
```
Review this diff as a security specialist. Check for: injection (command, SQL, path
traversal), unsafe deserialization, hardcoded secrets/credentials, missing input validation
on anything crossing a trust boundary, and unsafe use of os/exec or file paths built from
user/agent input. Rate each finding Critical/High/Medium/Low. If the diff has no security
surface at all (e.g. pure formatting), say so and stop — do not pad the report.
```

### 5. `secrets-scan` — Secret/Credential Scan
One sentence: a cheap, deterministic-ish gate that greps the diff for likely secrets
(API keys, tokens, private keys) before anything moves further.
Pipeline position: immediately after `in-progress`, ahead of every other gate — cheapest check,
should run first and block hard on a hit.
Distinction from `security-review`: narrower and mechanical (pattern-matching), meant to run on
every ticket, not just security-sensitive ones.
Input fields: `diff`.
Output fields: `secrets-status` (clean/flagged), `secrets-findings`.
Model tier: cheap (pattern-matching augmented by a small model to reduce false positives on
things that look like keys but aren't, e.g. test fixtures).
Source: [PerryLink/dsh-skill-pack-security](https://github.com/PerryLink/dsh-skill-pack-security)
— lists "secret scan" as one of 8 bundled agent skills in a security skill pack; small/community
project, adoption not independently verified.
Draft prompt:
```
Scan the diff for anything that looks like a secret: API keys, access tokens, private keys,
passwords, connection strings with embedded credentials. Flag matches with file:line. If a
match is clearly a placeholder/example value (e.g. "sk-xxxx", test fixture), say so rather
than blocking on it. Report clean/flagged status in one line at the top.
```

### 6. `dependency-audit` — Supply-Chain / License Gate
One sentence: checks new or changed dependencies for known vulnerabilities and license
compatibility before they land.
Pipeline position: after `in-progress`, only triggered when the diff touches go.mod/go.sum (or
equivalent manifest).
Input fields: `diff` (manifest files specifically).
Output fields: `dependency-report` (new deps + CVEs + licenses), `dependency-status`.
Model tier: cheap (mostly tool output summarization — the actual scanning is a deterministic
tool call, e.g. `govulncheck`, not an LLM judgment call).
Source: [bureado/awesome-software-supply-chain-security](https://github.com/bureado/awesome-software-supply-chain-security)
— long-running curated list in the supply-chain security space, widely cross-referenced in the
space though exact popularity not independently verified here; and GitHub's own
[dependency-audit-with-Copilot walkthrough](https://github.blog/developer-skills/github/video-how-to-run-dependency-audits-with-github-copilot/).
Draft prompt:
```
For each dependency added or version-bumped in this diff, report: known CVEs at the new
version, the license, and whether the license is compatible with this project's license.
If nothing in the manifest changed, say "no dependency changes" and stop.
```

---

## Theme: Testing

### 7. `test-writer` — Test-First Authoring
One sentence: writes the failing tests for a ticket's definition-of-done before (or alongside)
implementation, so `in-progress` is implementing against a red test suite rather than writing
tests as an afterthought.
Pipeline position: after `pre-process` (plan exists), before `in-progress` — or as the first
half of `in-progress` if jAIra doesn't want to add a lane for this.
Input fields: `goal`, `dod`, `plan`.
Output fields: `test-plan`, `tests` (the actual test code / diff).
Model tier: cheap — same tier as `in-progress`, since writing tests from a clear DoD is
comparably mechanical to implementing against one.
Source: [obra/superpowers](https://github.com/obra/superpowers/) — "TDD (Test-Driven
Development)" skill, described as enforcing RED-GREEN-REFACTOR (write failing test, watch it
fail, write minimal code, watch it pass, commit). This is distributed as an official plugin on
Anthropic's Claude plugin marketplace, which is a real (if not numeric) adoption signal.
Draft prompt:
```
From the ticket's goal and definition-of-done, write tests that fail today because the
feature doesn't exist yet. Cover the stated DoD items and the obvious edge cases (empty
input, error path, boundary values) — do not speculate about cases the DoD doesn't ask for.
Run the tests and confirm they fail for the right reason (missing feature, not a typo).
Hand off the test files as-is; do not implement the feature.
```

### 8. `coverage-gate` — Diff Coverage Check
One sentence: measures the coverage delta on just the changed lines and blocks if new code
shipped without tests touching it.
Pipeline position: after `in-progress`, before `critique`.
Input fields: `diff`.
Output fields: `coverage-delta`, `coverage-status`.
Model tier: cheap (tool-driven: run tests with coverage, diff against changed line ranges;
LLM only needed to render the summary).
Source: [GitHub Copilot Testing for .NET — devblogs.microsoft.com](https://devblogs.microsoft.com/dotnet/github-copilot-testing-for-dotnet-available-in-visual-studio/)
— official Microsoft blog, describes the @Test agent measuring "code coverage deltas" and
iterating on failures; and [qodo-ai/qodo-cover](https://github.com/qodo-ai/qodo-cover)
(open-source coverage-enhancement tool, stars not independently verified in this session but
it's a named, maintained project under the Qodo org).
Draft prompt:
```
Run the test suite with coverage. Compare coverage on only the lines this diff touched or
added against the rest of the codebase's baseline. Report which changed lines are exercised
by a test and which are not. Do not require 100% — flag genuinely untested new logic
(branches, error paths), not every unexercised line (e.g. trivial getters).
```

### 9. `perf-regression-check` — Performance/Benchmark Regression
One sentence: runs the project's benchmarks (if any) and compares against the pre-change
baseline, flagging regressions before they ship.
Pipeline position: after `in-progress`, for tickets whose plan touches a hot path (file
watching, git shell-out, TUI render loop — all named as perf-sensitive in jAIra's own stack
notes).
Input fields: `diff`, `plan`.
Output fields: `perf-report` (before/after numbers), `perf-status`.
Model tier: cheap (benchmark execution + delta is mechanical; only summarization needs a model).
Source: general pattern described in benchmark/agent literature —
[ProdCodeBench (arXiv)](https://arxiv.org/html/2604.01527v1) notes production-derived agent
tasks "center on debugging from real breakage where developers ask agents to investigate
regressions." No single popular named tool/prompt found for this specific step; treat this
idea as a synthesized pattern rather than a directly-cited popular prompt — flagged as weaker
evidence than the others in this report.
Draft prompt:
```
Run the project's benchmarks. Compare each result against the last recorded baseline for
this branch. Report only benchmarks that regressed beyond noise (>10%), with before/after
numbers. If there are no benchmarks covering the changed code, say so — do not write new
benchmarks here.
```

---

## Theme: Accessibility (its own small theme — distinct enough from generic quality gates)

### 10. `accessibility-review` — a11y / WCAG Audit
One sentence: reviews UI-touching diffs for WCAG compliance (labels, contrast, keyboard nav,
focus management) — relevant to jAIra's Bubble Tea TUI (keyboard-driven, needs equivalent
care around focus/contrast even in a terminal UI).
Pipeline position: after `in-progress`, for tickets whose diff touches TUI rendering
(`View()`/Lip Gloss styling) or any future web/GUI surface.
Input fields: `diff`.
Output fields: `a11y-findings`, `a11y-status`.
Model tier: strong (judging whether a color-contrast choice or keyboard trap is a real problem
needs contextual reasoning, not just a rule lookup).
Source: [guillempuche/ai-agent-a11y-accessibility-reviewer](https://github.com/guillempuche/ai-agent-a11y-accessibility-reviewer)
and [github/awesome-copilot — accessibility.agent.md](https://github.com/github/awesome-copilot)
(official GitHub org collection); also [Community-Access/accessibility-agents](https://github.com/taylorarndt/a11y-agent-team)
describing "eleven specialists that enforce WCAG 2.2 AA compliance." Popularity/star counts not
independently verified in this session, but three independent projects converging on the same
lane shape is a reasonable adoption signal.
Draft prompt:
```
Review this diff for accessibility. For a TUI: check keyboard reachability of every new
interactive element, that focus order is sane, that color is not the only signal for state
(e.g. selected vs. not), and that contrast holds in both light and dark terminal themes. For
a web/GUI surface: check WCAG 2.2 AA — labels, ARIA roles, contrast ratios, focus traps.
Report only real findings, not a boilerplate checklist.
```

---

## Theme: Documentation

### 11. `docs-writer` — Reference/Docs Generation
One sentence: writes or updates README/docstrings/reference docs for the changed code, sorted
into Diataxis categories (tutorial/how-to/reference/explanation) rather than dumped into one
undifferentiated doc.
Pipeline position: after `critique` (implementation is settled) and before `review`, so docs
describe what actually shipped, not what was planned.
Input fields: `diff`, `outcome-what`, `outcome-why`.
Output fields: `docs-diff`.
Model tier: cheap — documentation from a settled diff is close to summarization.
Source: [Diataxis Meets AI — romainlespinasse.dev](https://www.romainlespinasse.dev/posts/diataxis-documentation-skill/)
and [explainx.ai — writing-documentation-with-diataxis skill](https://explainx.ai/skills/sammcj/agentic-coding/writing-documentation-with-diataxis)
— both describe a documentation skill that classifies docs into the four Diataxis categories
and flags pages blending categories; independent blog + skill-marketplace listing, popularity
not independently verified.
Draft prompt:
```
Update documentation for what this diff actually changed. Classify what you're writing:
tutorial (learning-oriented), how-to (task-oriented), reference (information-oriented), or
explanation (understanding-oriented) — and keep each piece in its own lane; don't blend a
how-to with an explanation. Only touch docs affected by this diff. Match existing doc voice
and structure — do not restyle untouched sections.
```

### 12. `changelog-writer` — Changelog / Release Notes Entry
One sentence: produces a user-facing changelog entry for the ticket, distinct from the
commit message and from `outcome-what` (which is for other agents, not end users).
Pipeline position: after `critique`, before `review` or `signoff` — so the entry reflects the
final shipped behavior.
Input fields: `outcome-what`, `outcome-why`, `diff`.
Output fields: `changelog-entry`.
Model tier: cheap.
Source: [CodeRabbit — A Step-by-Step Agentic SDLC Workflow](https://www.coderabbit.ai/blog/agentic-sdlc-workflow)
— vendor blog describing "generate a changelog summary and commit message based on the
modified files" as a standard step in agentic SDLC pipelines; and
[coderabbitai/skills CHANGELOG.md](https://github.com/coderabbitai/skills/blob/main/CHANGELOG.md)
as a live example of a maintained changelog produced this way. Vendor content — treat
popularity as "this vendor considers it standard," not independently verified adoption.
Draft prompt:
```
Write one changelog entry for this ticket, addressed to a user of the tool, not another
agent. State what changed and why it matters to them in plain language — no internal
implementation detail, no jargon from the ticket's internal notes. One or two lines.
```

### 13. `adr-writer` — Architecture Decision Record
One sentence: for tickets whose plan made a real architectural call (not just an
implementation detail), writes a short ADR capturing the decision, alternatives considered,
and why.
Pipeline position: after `pre-process` (a plan exists) — only for tickets flagged as
architecturally significant, otherwise skipped.
Input fields: `plan`, `goal`.
Output fields: `adr` (decision, alternatives, consequences).
Model tier: strong (judging what counts as "architecturally significant" and articulating
tradeoffs needs real reasoning).
Source: [github/awesome-copilot — adr-generator.agent.md](https://github.com/github/awesome-copilot)
— official GitHub org collection of Copilot agents, includes an ADR-generator agent by name.
Draft prompt:
```
If this plan makes a real architectural decision (choice of storage format, protocol,
library with hard-to-reverse consequences, etc.), write a short ADR: the decision, the
alternatives considered and why they were rejected, and the consequences (what becomes
easier, what becomes harder). If this ticket is a routine change with no real architectural
choice, say so and produce nothing.
```

---

## Theme: Planning / decomposition

### 14. `spec-decompose` — SPIDR-style Story Splitting
One sentence: before `pre-process` plans a ticket, checks whether the ticket is actually
several tickets in a trenchcoat and splits it using SPIDR (Spike/Path/Interface/Data/Rules).
Pipeline position: after `brainstorm` (goal exists), before `todo` — if the ticket is too big,
this lane produces child tickets and the parent is closed/redirected rather than proceeding.
Input fields: `goal`.
Output fields: `child-tickets` (list of smaller goals), `split-rationale`.
Model tier: strong (judging story size and where the natural seams are is a planning task, not
mechanical).
Source: [Mountain Goat Software — SPIDR: Five Simple but Powerful Ways to Split User Stories](https://www.mountaingoatsoftware.com/agile/five-simple-but-powerful-ways-to-split-user-stories)
— Mike Cohn's well-established agile technique (long-standing, widely cited in agile practice,
predates LLMs but is now explicitly being applied to AI-assisted planning); and
[arXiv — Splitting User Stories Into Tasks with AI: A Foe or an Ally?](https://arxiv.org/html/2605.07320)
(academic study finding current AI tools aid but don't replace developers at this task —
worth noting as a caution, not just an endorsement).
Draft prompt:
```
Look at this ticket's goal. If it's small enough to implement and verify in one pass, say so
and stop. If not, split it using SPIDR: is there an unknown worth a Spike first? Is there an
obvious simplest Path through it? Can the Interface be built before the full implementation?
Can it be split by Data variant? Can Rules/edge-cases be deferred to a follow-up? Produce one
child ticket per piece, each independently shippable.
```

### 15. `spike` — Time-boxed Research/Prototype
One sentence: for a ticket with a real unknown (unproven library, unclear feasibility), does a
throwaway investigation and reports findings, without producing production code.
Pipeline position: between `brainstorm` and `pre-process`, only when `spec-decompose` (or
`brainstorm` itself) flags an open unknown.
Input fields: `goal`.
Output fields: `spike-findings`, `feasibility-verdict`.
Model tier: cheap-to-strong depending on the unknown; default strong since spikes exist
specifically because the easy answer isn't known yet.
Source: same as above — SPIDR's "Spike" is the named source
([Mountain Goat Software](https://www.mountaingoatsoftware.com/agile/five-simple-but-powerful-ways-to-split-user-stories)).
Draft prompt:
```
Investigate the specific unknown blocking this ticket's plan — try the smallest possible
throwaway experiment to answer it (a script, a one-off API call, reading the library's
source). Do not build production code or worry about code quality. Report what you learned
and a clear feasibility verdict: proceed as planned, proceed with a changed approach, or
this isn't feasible — with your evidence either way.
```

---

## Theme: Release / ops

### 16. `breaking-change-review` — API/Schema Compatibility Check
One sentence: for tickets touching the CLI's `--json` schema, exit codes, or frontmatter
fields, checks whether the change is backward compatible and flags it if not.
Pipeline position: after `in-progress`, before `critique` — specifically relevant to jAIra
given its own stated constraint that CLI exit codes and JSON schemas are "breaking changes"
that agents branch on.
Input fields: `diff`.
Output fields: `compat-report` (breaking/non-breaking + specifics), `compat-status`.
Model tier: cheap (mostly diffing a schema/contract; strong only if the contract is implicit
in code rather than a formal schema).
Source: [Zuplo — API Versioning & Backward Compatibility Best Practices](https://zuplo.com/learning-center/api-versioning-backward-compatibility-best-practices)
and [Aikido — Avoid breaking API contracts](https://www.aikido.dev/code-quality/rules/how-to-avoid-breaking-public-api-contracts-maintaining-backward-compatibility)
— both vendor/practitioner guides describing "automated schema checks in CI should compare
current spec against previous release to flag removed properties, new required fields, enum
contractions, and type changes"; popularity of these specific articles not verified, but the
CI-diff-based pattern they describe is widely used (e.g. OpenAPI diff tooling).
Draft prompt:
```
Compare this diff's public surface (CLI flags, --json output shape, exit codes, ticket
frontmatter fields) against what existed before. Flag anything that would break a script or
agent written against the old surface: removed fields, changed field types, removed flags,
changed exit code meanings. Additive changes (new optional field, new flag) are not
breaking — say so explicitly rather than flagging everything.
```

### 17. `rollout-plan` — Staged Rollout / Rollback Note
One sentence: for a risky change, writes down how it would be rolled back or staged if it goes
wrong in the field — a concrete plan, not a vague "monitor it."
Pipeline position: after `critique`, before `signoff` — only for tickets flagged high-risk.
Input fields: `diff`, `outcome-why`.
Output fields: `rollback-plan`.
Model tier: cheap.
Source: weaker citation — general practice mentioned in passing in
[API Contract Governance for Mobile Platforms — Medium](https://medium.com/@vaibhav.shakya786/api-contract-governance-for-mobile-platforms-versioning-compatibility-deprecation-and-rollback-da92c46bb739)
(rollback strategy as part of contract governance). This is a synthesized idea backed by a
single, not-highly-verified source — flagged as the weakest-evidence idea in this report; ship
it last if at all.
Draft prompt:
```
If this change ships and turns out to be wrong, what is the fastest safe way back? Name the
concrete rollback step (revert commit, flag flip, config change) and anything that would make
rollback hard (a migration that isn't reversible, a schema change already read by old
clients). If rollback is trivial (a normal revert works, nothing else touches it), say so in
one line.
```

---

## Theme: Communication

### 18. `pr-description-writer` — PR Summary for Human Reviewers
One sentence: writes the PR title/description/test-plan checklist that a human sees on GitHub,
distinct from `review` (a second model judging the diff for correctness).
Pipeline position: after `review`, before `signoff` — reviewers should see a clear summary
right when the ticket lands in their queue.
Input fields: `outcome-what`, `outcome-why`, `diff`, `review-summary`.
Output fields: `pr-description`.
Model tier: cheap.
Source: [CodeRabbit — "Meet the Overview Page for Pull Requests"](https://www.coderabbit.ai/blog/introducing-overview)
— vendor blog, describes the PR "Overview page" as "a single PR-level home that answers what a
PR is and whether it can merge," which is the same job as this lane's output; vendor content,
adoption not independently verified beyond CodeRabbit's own product framing.
Draft prompt:
```
Write a PR description for a human reviewer who has not read this ticket. Lead with what
changed and why, in plain language. Include a short test plan: what was run to verify this,
and what a reviewer should check by hand if anything. Do not repeat the full diff or restate
every commit — summarize.
```

### 19. `postmortem` — Blameless Retro for Bounced/Blocked Tickets
One sentence: when a ticket has bounced between `critique`/`optimize` and `in-progress`
repeatedly, or lands in `blocked` for a long time, writes a short retro on why, so the pattern
doesn't repeat on the next ticket.
Pipeline position: triggered from `blocked`, or from a critique-loop counter exceeding a
threshold — writes back into the ticket's notes, doesn't move it to a new lane by itself.
Input fields: `notes` (the ticket's accumulated dead-end log), `outcome-why`.
Output fields: `retro-summary`, `pattern-flag` (something to watch for on future tickets).
Model tier: strong (spotting a real pattern across a messy note history needs judgment).
Source: [Rootly — Streamlined Incident Post-Mortems + AI prompts](https://rootly.com/blog/streamlined-incident-post-mortems-a-concise-template-ai-prompts-for-artefacts)
and [engify.ai — Incident Post-Mortem Facilitator prompt](https://www.engify.ai/prompts/incident-post-mortem-facilitator)
— both describe AI-assisted postmortem generation (timeline, root cause, contributing factors,
what went well, action items); notably, Rootly's own material states the "lessons learned"
section is deliberately left to a human because judgment there is most consequential — this
report's `postmortem` lane should follow that same caution and treat its output as a draft for
a person, not a final verdict.
Draft prompt:
```
This ticket bounced back multiple times (or sat blocked). Read its notes and history. Write
a short, blameless summary: what actually caused the rework or the block, not who. Flag one
concrete pattern worth watching for on similar future tickets. Do not editorialize about
whether anyone should have caught it sooner — describe what happened and what to change
going forward.
```

---

## Top 5 I would ship first

1. **`secrets-scan`** — cheapest possible gate, catches the single worst failure mode (a
   committed credential), runs on every ticket, near-zero cost to add.
2. **`security-review`** — jAIra shells out to git and touches the filesystem constantly
   (`os/exec`, path handling); a dedicated security pass is directly relevant to its own stack,
   not a generic add-on.
3. **`test-writer`** — the single most popular, most independently-corroborated pattern in this
   research (TDD via obra/superpowers, Qodo-Cover, GitHub's own @Test agent); jAIra's
   `in-progress` lane currently implies tests happen but doesn't name a contract for them.
4. **`changelog-writer`** — tiny, cheap, and directly serves jAIra's own stated Core Value
   ("you never lose track of what a task was for") for end users, not just other agents.
5. **`breaking-change-review`** — jAIra's own CLI docs explicitly call CLI exit codes and
   `--json` schemas contracts that agents branch on and states changing them is "a breaking
   change" — this lane operationalizes a rule the project has already committed to in writing.

## Honesty notes on sourcing

- No star counts were pulled via the GitHub API in this research; "popular"/"widely adopted"
  claims above are inferred from search-result framing (official org repos like
  `github/awesome-copilot`, marketplace-distributed plugins like `obra/superpowers`) rather than
  verified numbers. Where a claim rests on a single vendor blog with no independent
  corroboration (`perf-regression-check`, `rollout-plan`), this is flagged explicitly in that
  idea's Source line.
- All URLs above came directly from web search results; none were invented or guessed.
