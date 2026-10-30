---
id: 01M658C60KGB8TEHQYJ5A40V1Z
title: "Das offene Ticket zeigt oben Kopfdaten und DoD-Stand, die Agenten-Notizen darunter"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Wer ein Ticket oeffnet (besonders in human review), sieht zuerst ID, Lane, Assignee, Goal, Context und den Stand der DoD, und erst darunter die Notizen, die fuer Agenten geschrieben sind."
context: |-
  Wunsch von berkc am 30.10.2026: 'Right now the Tickets dont help me alot when i open them. The ticket in human review should still show the meta data of the ticket on top (id, lane, assignee, goal, context etc.) then the state of the DoD. All the notes you take on the ticket is just there for agents to still have it in mind.'
  Betroffen: die Detailansicht im TUI (internal/tui, Detail-Rendering).
  Offen: Notizen eingeklappt oder nur nach unten? Vorschlag: eingeklappt mit Zaehler, Taste zum Aufklappen.
definition-of-done: "Detailansicht: oben ID, Lane, Assignee, Goal, Context; danach DoD mit Haken und Proofs; Notizen zuletzt"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-30T01:56:46Z
updated-at: 2026-10-30T23:45:33Z
updated-by: berkc
---

# Das offene Ticket zeigt oben Kopfdaten und DoD-Stand, die Agenten-Notizen darunter

## Definition of Done

- [ ] Detailansicht: oben ID, Lane, Assignee, Goal, Context; danach DoD mit Haken und Proofs; Notizen zuletzt
- [ ] Gilt in jeder Lane, gezeigt an einem Ticket in human; Test
- [ ] Release-Note-Zeile

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-30 21:39 · berkc** — berkc 30.10.2026: (1) Gemeint sind alle Metadaten, die ein Ticket schon im Backlog zeigt (ID, Lane, Assignee, Creator, Tags, Goal, Context, follows/related, Milestone) - in jeder Lane oben, auch beim Oeffnen vom Board aus. (2) Fuer Tickets, die ein Mensch pruefen soll: oben eine Anleitung, wie er es prueft - das ist das review-check-Feld, das die review-Lane heute schon schreiben muss (review-summary/-gaps/-verdict/-check); es muss sichtbar oben stehen. (3) Die Bloecke 'Outcome / Review / Lane fields - owed by <lane>' zeigen, welche Lane welches Ausgabefeld noch schuldet; im Backlog ist das nur Rauschen: erst zeigen, wenn die Lane erreicht ist, oder einklappen.
