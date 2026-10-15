---
id: 01M4ZQ715HGZN2GHX5NW6TCTFV
title: "Ist der letzte Ticket eines Milestones fertig, schlaegt jaira den PR mit /jaira-role-pr vor und bietet an, den Milestone abzulegen"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Sobald alle Tickets eines Milestones fertig sind, sagt jaira (und jeder Teamlead/Dispatcher) dem Menschen: jetzt PR mit /jaira-role-pr, und danach den Milestone ablegen. Geoeffnet wird der PR weiter nur, wenn ein Mensch die Rolle aufruft"
context: |-
  Was heute falsch ist: Der PR-Skill (jaira-role-pr) wird oft vergessen. Keine Lane loest ihn aus, und kein Bericht schlaegt ihn vor.
  Beispiel aus einer anderen Session (requirementsgenie, Ticket 6ESQWR): Das Ticket kam in signoff an, und niemand schlug '/jaira-role-pr 6ESQWR' vor. Der MR !56 musste von Hand nachgezogen werden.
  Absicht von jaira, und die bleibt: Den PR oeffnet nur ein Mensch. jaira-role-pr SKILL.md:11 'Opening the pull request is the maintainer's call, and typing /jaira-role-pr is how they give it'; :16 von einem Agenten aufgerufen pusht er nur. Dispatcher :409 und Teamlead :167 oeffnen nie einen PR.
  Berk am 15.10.: 'wenn wir die Tickets, die unter einem Milestone stehen, alle fertiggestellt haben, ein PR vorschlagen und anbieten, das mit dem jaira pr skill zu machen.'
  Berk fragte auch, warum fertige Milestones (milestone-roadmap, ram-cpu) noch auf der Roadmap stehen. Ablegen geht schon mit 'jaira logbook <milestone>' (Ref traegt dann status: filed), aber niemand schlaegt es vor.
  Andere Ideen, die verworfen bzw. eingeordnet wurden:
  - nur eine Memory-Regel fuer eine Session: zu schmal
  - CLAUDE.md jaira:local: gilt nur fuer ein Repo
  - PostToolUse-Hook: zerbrechlich, nur Claude Code
  Gewaehlt: in jaira selbst. Die CLI-Ausgabe erreicht jeden Agenten, nicht nur Claude (Projekt-Constraint 'CLI must be usable by any bash-capable agent').
  Verwandt: 3SJMV8 (PR-Abschnitt im jaira-Skill), N71NVG (fertigen Milestone archivieren).
definition-of-done: "Macht ein 'jaira move <id> --to done' oder 'jaira logbook <id>' den letzten offenen Ticket eines Milestones fertig, druckt jaira eine Zeile: PR mit '/jaira-role-pr <ids>' vorschlagen und danach 'jaira logbook <milestone>'; --json traegt es als Feld; Test"
tags: []
blocked-by: []
related:
  - 01M4VGHHEKBXMV9120AD3SJMV8
  - 01M2K275PSXDASRN6Y1VN71NVG
commits: []
created-at: 2026-10-15T12:04:57Z
updated-at: 2026-10-15T12:19:30Z
updated-by: BeMuCa
---

# Ist der letzte Ticket eines Milestones fertig, schlaegt jaira den PR mit /jaira-role-pr vor und bietet an, den Milestone abzulegen

## Definition of Done

- [ ] Macht ein 'jaira move <id> --to done' oder 'jaira logbook <id>' den letzten offenen Ticket eines Milestones fertig, druckt jaira eine Zeile: PR mit '/jaira-role-pr <ids>' vorschlagen und danach 'jaira logbook <milestone>'; --json traegt es als Feld; Test
- [ ] jaira-dispatcher und jaira-teamlead nennen in ihrem Bericht den PR-Vorschlag, sobald ein Milestone komplett ist, und rufen /jaira-role-pr nie selbst auf
- [ ] jaira-role-pr nimmt mehrere Tickets bzw. einen Milestone und legt alle angenommenen Tickets vor dem Push ab; fuer einen bestehenden PR/MR aktualisiert er Titel und Beschreibung aus den Tickets
- [ ] Eine Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

