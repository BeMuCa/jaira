---
id: 01M5ZJ555685ZXARWNEF5QR36R
title: Die Abnahme-Rolle fuehrt das Interview im Chat; die Seite ist nur noch auf Wunsch
status: backlog
ready: false
creator: BeMuCa
goal: "jaira-role-acceptance fragt standardmaessig im Chat ab - je Ticket ID, Einzeiler, wie pruefen, Empfehlung, dann 'accept X / explain X / fix X' - und baut die statische Seite nur mit --page."
context: |-
  Was schief ist: die Abnahme-Seite veraltete mit jedem Rebuild, der Read-back lief nie; das Interview im Chat war fuer 17 Tickets schneller als jede Seite (Feedback L324/L335/O24, 27.10.2026).
  Heute: core/role/builtin/jaira-role-acceptance/SKILL.md baut eine Seite und liest Markierungen zurueck; der Chat-Weg steht nirgends.
  Was gut lief: pro Ticket ID, ein Satz was, wie pruefen, Empfehlung, Frage oder Entscheidung; der Mensch antwortet 'accept 1TFFG2' oder stellt die Frage.
  Zusammenhang: Q63EEV (Cockpit-Pane je Ticket, Frage offen in human) - dieses Ticket aendert nur die Rolle, nicht das Cockpit.
definition-of-done: SKILL.md beschreibt das Chat-Interview als Default mit dem Format je Ticket und den drei Antworten
tags:
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-10-27T20:52:15Z
updated-at: 2026-10-30T23:45:25Z
assignee: berkc
updated-by: berkc
---

# Die Abnahme-Rolle fuehrt das Interview im Chat; die Seite ist nur noch auf Wunsch

## Definition of Done

- [ ] SKILL.md beschreibt das Chat-Interview als Default mit dem Format je Ticket und den drei Antworten
- [ ] --page baut die Seite wie bisher; ohne --page wird keine Datei geschrieben
- [ ] Die statische Aufzeichnung der Abnahme entsteht auch im Chat-Weg (gleicher Pfad wie heute)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

