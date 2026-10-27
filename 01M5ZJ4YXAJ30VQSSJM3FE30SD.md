---
id: 01M5ZJ4YXAJ30VQSSJM3FE30SD
title: "Ein Ticket zeigt je Lane, wie lange sie dauerte und was sie kostete"
status: backlog
ready: false
creator: BeMuCa
goal: "Jeder Lane-Durchlauf traegt Minuten (aus den Zeitstempeln von claim/move) und, wenn der Agent es mitgibt, Tokens; jaira show listet sie je Lane und Runde, jaira list --json summiert sie je Ticket."
context: |-
  Was schief ist: zwei Lane-Laeufe kosteten 59k und 49k Tokens statt der angenommenen 18k; 14 Driver fuer 12 kleine Tickets kosteten rund 2,5M Tokens; critique lief bis zu sechs Runden (VPM51D). Es gibt keinen Ort, wo Kosten je Lane stehen (Feedback L319/L329/O8, 27.10.2026).
  Heute: move schreibt Lane und Zeit; Dauer und Tokens stehen nirgends. Tokens kennt nur der Agent (oder sein Cockpit-Report).
  Bekannt: 087E3M will eine Messzahl je Ticket (Ertrag) - das hier ist die Kostenseite, ein anderes Feld.
  Nicht gewollt: Preise in Euro. Minuten und Tokens reichen, der Rest ist Rechnung des Lesers.
definition-of-done: "jaira move nimmt --tokens <n> an und schreibt es zusammen mit der Dauer in Minuten an den Lane-Durchlauf"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-27T20:52:08Z
updated-at: 2026-10-27T20:52:08Z
---

# Ein Ticket zeigt je Lane, wie lange sie dauerte und was sie kostete

## Definition of Done

- [ ] jaira move nimmt --tokens <n> an und schreibt es zusammen mit der Dauer in Minuten an den Lane-Durchlauf
- [ ] jaira show <id> zeigt je Lane und Runde Minuten und Tokens; --json liefert sie strukturiert
- [ ] jaira list --json traegt je Ticket die Summe
- [ ] Release-Notiz in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

