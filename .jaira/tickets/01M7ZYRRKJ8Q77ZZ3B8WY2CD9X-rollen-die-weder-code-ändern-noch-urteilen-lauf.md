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
updated-at: 2026-11-21T21:07:09Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2840519
claimed-at: 2026-11-21T21:04:44Z
outcome-what: "pr, tester and logbook-summary roles carry model: haiku in their SKILL.md header; the pr role hands code-changing review points to an opus agent"
outcome-why: "pushing, running suites and summarising do not need a strong model; judging and code changes still do"
outcome-resolves: "PR, CI and test runs cost haiku prices without weakening critique, review or bughunt"
---

# Rollen, die weder Code ändern noch urteilen, laufen auf haiku

## Definition of Done

- [x] pr-, tester- und logbook-summary-Rolle starten auf haiku, critique/review/bughunt bleiben strong (Test oder roles list zeigt es)
  proof: core/role/role_test.go TestRoleModels; model: haiku in core/role/builtin/jaira-role-{pr,tester,logbook-summary}/SKILL.md
- [x] Ein PR-Kommentar, der eine Codeänderung verlangt, geht an einen Worker auf dem normalen Modell
  proof: core/role/builtin/jaira-role-pr/SKILL.md:484 (step 1: Agent model opus for code-changing points)
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
