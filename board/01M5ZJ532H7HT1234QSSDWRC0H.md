---
id: 01M5ZJ532H7HT1234QSSDWRC0H
title: "Der Agent-Block sagt, wie man Goal, Kontext und DoD nach dem Anlegen aendert"
status: backlog
ready: false
creator: BeMuCa
goal: "Der Agent-Block in CLAUDE.md nennt jaira set fuer goal und context und einen dokumentierten Weg, einen DoD-Punkt nachzutragen; niemand muss das per Notiz behelfen."
context: |-
  Was schief ist: Teamlead und Worker glaubten, 'jaira edit fehlt' und ergaenzten Goal und DoD per jaira note (Feedback L323/L330/O20, 27.10.2026).
  Dabei gibt es jaira set <id> goal=... und set --append (jaira set --help, 27.10.2026) - nur steht es nicht im Agent-Block (CLAUDE.md, Abschnitt 'Task tracking: jaira'), und ob set einen DoD-Punkt sauber anhaengt, ist nicht geprueft.
  Zu klaeren im Ticket: wie ein einzelner DoD-Punkt nachgetragen wird (set definition-of-done --append oder ein dod --add), und dass die Nummern bestehender Punkte dabei bleiben, sonst zeigen alte --proof-Eintraege ins Leere.
  Nicht gewollt: ein neuer Befehl jaira edit, wenn set reicht.
definition-of-done: Der Agent-Block nennt jaira set fuer goal und context mit einem Beispiel
tags:
  - docs
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-27T20:52:13Z
updated-at: 2026-10-27T20:52:13Z
---

# Der Agent-Block sagt, wie man Goal, Kontext und DoD nach dem Anlegen aendert

## Definition of Done

- [ ] Der Agent-Block nennt jaira set fuer goal und context mit einem Beispiel
- [ ] Ein DoD-Punkt laesst sich nach dem Anlegen anhaengen, ohne dass bestehende Nummern wandern; ein Test belegt es
- [ ] jaira create --help verweist auf set fuer spaetere Aenderungen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

