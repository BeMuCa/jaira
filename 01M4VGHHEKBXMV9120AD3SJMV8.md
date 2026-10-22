---
id: 01M4VGHHEKBXMV9120AD3SJMV8
title: "Der jaira-Skill sagt, wie ein PR mit seinen Tickets entsteht, und verweist auf die PR-Rolle"
status: todo
ready: true
creator: BeMuCa
goal: "Ein Agent, der fuer ein jaira-Repo einen PR vorbereitet, legt die angenommenen Tickets auf dem Branch ins Logbuch, bevor er pusht, weil der jaira-Skill es ihm sagt."
context: |-
  Am 2026-10-13 hat Claude PR #74 geoeffnet, waehrend dessen Tickets noch in signoff standen. Alex (sashasoft90) setzte 'Changes requested': 'Hi berk schliss mal tickets ab'.
  Die Regel steht in der PR-Rolle core/role/builtin/jaira-role-pr/SKILL.md ('File the finished tickets away': pull, logbook, ein chore-Commit pro Ticket, vor dem Push).
  Der jaira-Skill (.claude/skills/jaira/SKILL.md) erwaehnt die PR-Rolle nirgends und sagt das Gegenteil: 'An accepted ticket leaves the board once its work is pushed ... the logbook waits on it' (SKILL.md:427 auf feat/milestone-roadmap). Der Agent-Block in CLAUDE.md nennt die Rolle auch nicht.
  Optionale Schritte erklaert der Skill schon ('## Optional steps': ein nicht angehakter Schritt wird beim Hineinschieben abgelehnt) - das muss nicht dazu.
  Verwandt: PWAF22 (Skills nach einem Upgrade nachziehen) - ein Skill, der nicht neu installiert wird, hilft nicht.
definition-of-done: "SKILL.md hat einen Abschnitt zum PR: angenommene Tickets werden auf dem PR-Branch mit jaira pull, jaira logbook und einem chore-Commit pro Ticket abgelegt, bevor gepusht wird; er verweist auf /jaira-role-pr"
tags:
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-10-13T20:51:24Z
updated-at: 2026-10-22T07:59:04Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-4188296
claimed-at: 2026-10-22T07:55:16Z
---

# Der jaira-Skill sagt, wie ein PR mit seinen Tickets entsteht, und verweist auf die PR-Rolle

## Definition of Done

- [x] SKILL.md hat einen Abschnitt zum PR: angenommene Tickets werden auf dem PR-Branch mit jaira pull, jaira logbook und einem chore-Commit pro Ticket abgelegt, bevor gepusht wird; er verweist auf /jaira-role-pr
  proof: .claude/skills/jaira/SKILL.md:459 (section Making a pull request)
- [x] Der Satz 'the logbook waits on [the push]' (SKILL.md:427) widerspricht dem nicht mehr: er unterscheidet direkten Push und PR
  proof: .claude/skills/jaira/SKILL.md:439 (direct push vs PR)
- [ ] Der Agent-Block (core/board/announce.go) nennt /jaira-role-pr fuer PRs in einer Zeile; CLAUDE.md und AGENTS.md neu erzeugt; ein Test prueft die Zeile
- [ ] Ein frischer Subagent mit dem Auftrag 'mach einen PR fuer Ticket X' beschreibt das Ablegen vor dem Push, ohne dass man es ihm sagt; das Ergebnis steht als Notiz am Ticket

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

