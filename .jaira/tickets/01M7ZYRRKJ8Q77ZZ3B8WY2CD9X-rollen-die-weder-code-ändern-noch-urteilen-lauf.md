---
id: 01M7ZYRRKJ8Q77ZZ3B8WY2CD9X
title: "Rollen, die weder Code ändern noch urteilen, laufen auf haiku"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "PR-, CI- und Test-Läufe kosten haiku-Preise, ohne dass ein Urteil schwächer wird."
context: |-
  Idee Alex 21.11.: was keinen Code ändert, auf haiku.
  Abgrenzung (Teamlead): nicht nur "ändert keinen Code", sondern auch "urteilt nicht". critique, review und bughunt ändern keinen Code, sind aber Urteile — die bleiben strong.
  Kandidaten für haiku: jaira-role-pr (push, PR öffnen, Kommentare beantworten ohne Codeänderung), jaira-role-tester (Suite laufen lassen, geerbte von neuen Fehlern trennen), das Warten auf CI, jaira-role-logbook-summary.
  Heute: Lanes haben tier cheap/strong; Rollen legen kein Modell fest (kein model: in core/role/builtin/*/SKILL.md); Chats nutzen .jaira/chat-model (sonnet). spawn.sh startet claude ohne Modellwahl pro Rolle.
  Offen für pre-process: Mechanismus (model: in der SKILL-Frontmatter, Feld pro Rolle in den Board-Settings oder spawn.sh --model) und ob ein PR-Kommentar, der Codeänderung verlangt, an einen Worker auf dem normalen Modell geht.
definition-of-done: "pr-, tester- und logbook-summary-Rolle starten auf haiku, critique/review/bughunt bleiben strong (Test oder roles list zeigt es)"
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-11-21T21:04:04Z
updated-at: 2026-11-21T21:14:31Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-3504002
claimed-at: 2026-11-21T21:13:37Z
outcome-what: "pr, tester and logbook-summary roles carry model: haiku in their SKILL.md header; the pr role hands code-changing review points to an opus agent"
outcome-why: "pushing, running suites and summarising do not need a strong model; judging and code changes still do"
outcome-resolves: "PR, CI and test runs cost haiku prices without weakening critique, review or bughunt"
review-summary: "Three role skills — jaira-role-pr, jaira-role-tester, jaira-role-logbook-summary — get 'model: haiku' in their SKILL.md header, so Claude Code runs them on haiku however they are started. No other role carries a model line, so lane (critique/review), bughunt, worker, dispatcher and teamlead stay on the session's model; TestRoleModels in core/role/role_test.go pins exactly that set. Because the pr role would otherwise fix review comments on haiku, its step 1 now hands a code-changing review point to Agent(subagent_type claude, model opus) with only the thread, and the pr role pushes that commit. Release note core/release/notes/Y2CD9X.md tells users to run jaira update --roles; shipped.sha256 records the three new SKILL.md hashes."
round: 3
back-from: signoff
back-by: person
back-count: 1
review-gaps: |-
  note: jaira-role-pr step 1 hard-codes model "opus"; an account or agent host without opus fails that call - name it as "the strongest model you have, never haiku" or fall back to sonnet if opus is refused
  note: jaira-role-pr step 1 pushes the subagent's commit without checking it - add one line: confirm the printed sha exists on HEAD and its subject names <id> before pushing
  note: model: haiku is installed verbatim for --agent copilot too, where the header and the Agent(model) step likely mean nothing; harmless, but untested (bughunt known limit)
test-verdict: "pass: go test ./... RC=0, go vet, shipped-history --check, check-ticket-commits vs upstream/master all exit 0; DoD 1-3 verified in the tree; roles install --into scratch puts model: haiku on exactly pr/tester/logbook-summary, none on lane/bughunt; install into an unwritable dir exits 1; change writes nothing at runtime (only SKILL.md content), no concurrent-writer path"
hunt-verdict: "pass: bughunt at 7bb4c47db found nothing; three unverifiable limits noted (harness model switch, copilot install, account without opus)"
review-verdict: "Diff meets all three DoD items: exactly pr/tester/logbook-summary carry model: haiku and TestRoleModels fails if any judging role gains a model line; a code-changing PR point goes to an opus worker; release note and shipped.sha256 are present. Checked the other paths in jaira-role-pr that could change code (none-mode merge, conflicts) - it aborts conflicts and stops, so haiku never resolves code there. tester only reads logs. No blocker; three notes about robustness of the opus hand-off."
review-round: 2
review-verdict-gaps: |-
  note: jaira-role-pr step 1 hard-codes model "opus"; an account or agent host without opus fails that call - name it as "the strongest model you have, never haiku" or fall back to sonnet if opus is refused
  note: jaira-role-pr step 1 pushes the subagent's commit without checking it - add one line: confirm the printed sha exists on HEAD and its subject names <id> before pushing
  note: model: haiku is installed verbatim for --agent copilot too, where the header and the Agent(model) step likely mean nothing; harmless, but untested (bughunt known limit)
review-check: "1. cd into this worktree  2. go test ./core/role -run TestRoleModels -count=1  -> ok  3. head -5 core/role/builtin/jaira-role-pr/SKILL.md -> line 'model: haiku' under description; same for jaira-role-tester and jaira-role-logbook-summary  4. grep -l '^model:' core/role/builtin/*/SKILL.md -> exactly those three files, no jaira-role-lane, no jaira-role-bughunt  5. sed -n 484,492p core/role/builtin/jaira-role-pr/SKILL.md -> step 1 sends code-changing points to Agent with model opus  6. Live check (optional): in a fresh claude session run /jaira-role-tester and look at the model shown in the status line -> haiku for that turn"
---

# Rollen, die weder Code ändern noch urteilen, laufen auf haiku

## Definition of Done

- [x] pr-, tester- und logbook-summary-Rolle starten auf haiku, critique/review/bughunt bleiben strong (Test oder roles list zeigt es)
  proof: core/role/role_test.go TestRoleModels; model: haiku in core/role/builtin/jaira-role-{pr,tester,logbook-summary}/SKILL.md
- [x] Ein PR-Kommentar, der eine Codeänderung verlangt, geht an einen Worker auf dem normalen Modell
  proof: core/role/builtin/jaira-role-pr/SKILL.md:484 (step 1: Agent with no model = session's model)
- [x] Release-Note core/release/notes/<ID>.md, shipped.sha256
  proof: core/release/notes/Y2CD9X.md; core/shipped/shipped.sha256 (+3 lines)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-21 21:04 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2274544.sock
- **2026-11-21 21:05 · Alexander Sacharov** — decision: mechanism = model: haiku in the SKILL.md frontmatter of pr, tester, logbook-summary. Claude Code honours model: in a skill header for the turn the skill runs, so it holds however the role is started (spawn.sh tab, Skill tool, Agent). No board setting, no spawn.sh flag: smallest change, nothing new to configure.
- **2026-11-21 21:06 · Alexander Sacharov** — decision: Alex 21.11.: Grenze bestätigt — haiku nur für Rollen, die weder Code ändern noch urteilen (pr ohne Codeänderung, tester, CI-Warten, logbook-summary); critique, review, bughunt bleiben strong.
- **2026-11-21 21:07 · Alexander Sacharov** — decision: code-changing review point goes to Agent(model: opus), named explicitly because a subagent would inherit haiku from the skill. opus = what this account's workers run on; sonnet is only the cockpit chat-model. Not via spawn.sh worker: the ticket is usually in done/signoff by PR time, a full worker run would walk lanes.
- **2026-11-21 21:07 · Alexander Sacharov** — critique: release note line 1 — replace 'run jaira roles install ... remove the line from the installed SKILL.md' with 'jaira update --roles installs it'. Rest stands: model: header is the smallest mechanism (note above); opus named explicitly in pr step 1 because Agent would inherit haiku. Checked: teamlead/dispatcher start pr/tester as their own session, never inline, so haiku does not leak into a strong session.
- **2026-11-21 21:08 · Alexander Sacharov** — testing round 1: vet, shipped-history --check, check-ticket-commits vs upstream/master all exit 0 (vs origin/master it exits 1: the fork's master is stale, not this branch). roles install --into scratch: model: haiku in exactly pr/tester/logbook-summary, none in lane/bughunt; install into /proc/nope exits 1. Full go test ./... running; bughunt started.
- **2026-11-21 21:10 · Alexander Sacharov** — bughunt: done at 7bb4c47db - 0 findings in 6 modules (role package core/role read with role.go, install.go, target.go; the three SKILL.md files; spawn.sh pr lane; shipped list), first hunt; concurrency, error paths, input, permissions checked; files: core/release/notes/Y2CD9X.md core/role/builtin/jaira-role-logbook-summary/SKILL.md core/role/builtin/jaira-role-pr/SKILL.md core/role/builtin/jaira-role-tester/SKILL.md core/role/role_test.go core/shipped/shipped.sha256
- **2026-11-21 21:10 · Alexander Sacharov** — known limits (bughunt, unverifiable here): 1) that Claude Code switches to haiku on model: in the SKILL.md header is harness behaviour, not testable from Go; 2) --agent copilot installs the same header and the pr role's Agent(model: opus) step, which Copilot likely ignores/cannot run; 3) an account without opus makes the pr role's fix-agent call fail.
- **2026-11-21 21:13 · Alexander Sacharov** — Given back from signoff: der Code-Fix-Agent der pr-Rolle darf kein festes opus nennen — auf einem Konto ohne opus scheitert er. Modell weglassen bzw. inherit (das Modell der Session), damit er auf jedem Konto läuft. Rest der Review-Notes bleibt note.
- **2026-11-21 21:13 · Alexander Sacharov** — decision (Alex give-back 21.11.): pr role step 1 names no model on the fix agent; Agent without model runs on the session's model, so it works on any account. Claim 'inherits haiku from this role' removed from the skill text. Release note line reworded to 'the session's model'.
