---
id: 01M6H731JRQMNRBG1Y9XZMN3MW
title: "jaira web Tab 'Aktiv': ein Avatar je Person, darunter ihre offenen Tickets in der Farbe ihres Stands"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Ein PM sieht auf einen Blick, wer woran arbeitet und wo jedes Ticket steht."
context: |-
  Was heute fehlt: man sieht nur Tickets, nicht Personen. Wer gerade woran sitzt, muss der PM aus dem Filter 'Person' zusammensuchen.
  Teil der Aufteilung von CMN1YT (jaira web als Zentrale fuer PMs = Leute ohne Claude und ohne CLI); dort stehen Brainstorm und Antworten des Maintainers.
  Wunsch: kleine Avatare je Person, darunter ihre Tickets; Farbe nach Stand: backlog grau, in Arbeit orange, in Review orange-gruen, done gruen, blockiert rot.
  Entschieden: Avatar = Initialen auf einer Farbe, die aus dem Namen berechnet wird. Kein externer Bilddienst, weil das Netzzugriff und Datenschutzfragen bringt.
  Entschieden: die Farbe kommt aus der Lane-Reihenfolge des Boards des Tickets, nicht aus festen Lane-Namen - jedes Board hat eigene Lanes. Regel: Lanes vor der ersten Agenten-Lane grau; Agenten-Lanes vor der ersten Personen-Lane orange; ab der ersten Personen-Lane bis vor die Terminal-Lane orange-gruen; Terminal-Lane gruen; Parkspur oder offenes blocked-by rot (gewinnt).
  Bekannt: die heutigen Lane-Farben (internal/cli/web.go, Hilfetext) sind pro Lane; diese fuenf Stufen sind eine Zusammenfassung darueber, keine Ersetzung.
  Voraussetzung: DVPAPE, sonst erscheint dieselbe Person unter zwei Avataren.
definition-of-done: "Tab 'Aktiv': ein Avatar (Initialen, Farbe aus dem Namen) je Person mit mindestens einem offenen Ticket als Assignee; darunter ihre Tickets; Test"
tags:
  - web
blocked-by:
  - V9S8RN
  - DVPAPE
related: []
follows: 01M67K6Z3NSBTBXAK4EVCMN1YT
commits: []
created-at: 2026-11-03T17:25:11Z
updated-at: 2026-11-03T17:25:11Z
---

# jaira web Tab 'Aktiv': ein Avatar je Person, darunter ihre offenen Tickets in der Farbe ihres Stands

## Definition of Done

- [ ] Tab 'Aktiv': ein Avatar (Initialen, Farbe aus dem Namen) je Person mit mindestens einem offenen Ticket als Assignee; darunter ihre Tickets; Test
- [ ] Jedes Ticket traegt eine der fuenf Stufenfarben nach der Regel im Kontext, berechnet aus der Lane-Reihenfolge; Test mit zwei Boards mit unterschiedlichen Lanes
- [ ] Keine Anfrage an einen fremden Server; Test oder Nachweis im Code
- [ ] jaira web --help, eine Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

