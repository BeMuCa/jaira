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
updated-at: 2026-10-31T02:31:26Z
assignee: berkc
updated-by: berkc
claimed-by: berk-ca-3129429
claimed-at: 2026-10-31T01:27:42Z
---

# Die Abnahme-Rolle fuehrt das Interview im Chat; die Seite ist nur noch auf Wunsch

## Definition of Done

- [ ] SKILL.md beschreibt das Chat-Interview als Default mit dem Format je Ticket und den drei Antworten
- [ ] --page baut die Seite wie bisher; ohne --page wird keine Datei geschrieben
- [ ] Die statische Aufzeichnung der Abnahme entsteht auch im Chat-Weg (gleicher Pfad wie heute)

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-30 23:45 · berkc** — Frage berkc 31.10.2026: Das Interview in einer neuen Session bzw. als eigener Interviewer-Subagent fuehren? Vorteil: der Kontext der Hauptsession bleibt schlank, und das Interview laeuft parallel zur normalen Arbeit. Tokens: das Interview kostet gleich viel, aber nicht im teuren, langen Hauptkontext - jede Antwort dort wuerde den ganzen Verlauf neu mitschicken. Empfehlung: ja, eigener Agent/Tab, der nur die Tickets in human/signoff liest; die Hauptsession bekommt am Ende eine Zeile pro Ticket (angenommen/zurueck/Folgeticket).
- **2026-10-31 01:25 · berkc** — Maintainer 31.10.2026: umsetzen. Der Skill heisst interReview (Interview-Abnahme im Chat). Danach geht der Maintainer mit der Session die Tickets in review/signoff damit durch.
