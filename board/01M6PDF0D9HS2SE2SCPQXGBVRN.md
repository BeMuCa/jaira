---
id: 01M6PDF0D9HS2SE2SCPQXGBVRN
title: "Der Milestone-Picker zeigt, wer einen Milestone angelegt hat, und die eigenen stehen oben"
status: backlog
ready: false
creator: BeMuCa
goal: Eine Milestone-Datei traegt created-by; der Picker unter M zeigt den Namen neben jedem Milestone und sortiert die des aktuellen Nutzers nach oben.
context: |-
  Was fehlt (Berk, 2026-11-05): unter M sieht man nicht, wessen Milestone das ist. Auf einem geteilten Board mischen sich die Milestones aller.
  Heute: eine Milestone-Datei (.jaira/milestones/<name>-<ID>.md) hat color, created-at, goal, after - kein created-by. Tickets haben creator.
  Woher der Name kommt: wie bei Tickets (jaira whoami).
  Alte Dateien ohne created-by: kein Name, unten einsortiert, nichts geht verloren.
  Zusammenhang: das Ticket 'Milestones teilen sich ueber einen Ref' - erst dann sieht man ueberhaupt fremde Milestones.
definition-of-done: jaira milestone create schreibt created-by
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-11-05T17:52:47Z
updated-at: 2026-11-05T17:52:47Z
---

# Der Milestone-Picker zeigt, wer einen Milestone angelegt hat, und die eigenen stehen oben

## Definition of Done

- [ ] jaira milestone create schreibt created-by
- [ ] Der Picker unter M zeigt den Namen und sortiert die eigenen oben, danach alphabetisch
- [ ] Eine alte Datei ohne created-by wird angezeigt, ohne Namen, unten
- [ ] Release-Notiz in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

