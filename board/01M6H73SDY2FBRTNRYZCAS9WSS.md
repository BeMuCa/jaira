---
id: 01M6H73SDY2FBRTNRYZCAS9WSS
title: jaira web exportiert ein Epic oder einen Milestone auf Knopfdruck einmalig nach Jira
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Ein PM uebertraegt ein Feature oder eine Runde mit ihren Tickets nach Jira, sieht vorher, was angelegt wird, und erzeugt beim zweiten Mal keine Duplikate."
context: |-
  Was heute fehlt: der Stand in jaira erreicht Jira nicht; Leute, die nur Jira lesen, sehen nichts.
  Teil der Aufteilung von CMN1YT (jaira web als Zentrale fuer PMs = Leute ohne Claude und ohne CLI); dort stehen Brainstorm und Antworten des Maintainers.
  Entschieden (Maintainer): nur auf Anfrage, einmalig, wegen API-Quoten; kein Dauer-Sync in beide Richtungen.
  Bekannt: die Zuordnung kommt aus der einstellbaren Board-Datei (Mapping-Ticket); der Token aus secrets.json.
  Entschieden: was exportiert wurde, bekommt 'external: jira:<KEY>' - an Ticket, Milestone und Epic. Ein zweiter Export aktualisiert diese statt neue anzulegen.
  Ausgeschlossen: Aenderungen aus Jira zurueckholen.
definition-of-done: "Im Tab 'Milestones & Epics' je Epic/Milestone ein Knopf 'nach Jira'; eine Vorschau listet, was angelegt oder geaendert wird, und die Zahl der API-Anfragen; erst Bestaetigen sendet; Test gegen Fake-Server"
tags:
  - web
blocked-by:
  - 35Y6N5
  - VD4ZDW
  - GF2MAB
related: []
follows: 01M67K6Z3NSBTBXAK4EVCMN1YT
commits: []
created-at: 2026-11-03T17:25:35Z
updated-at: 2026-11-03T17:25:35Z
---

# jaira web exportiert ein Epic oder einen Milestone auf Knopfdruck einmalig nach Jira

## Definition of Done

- [ ] Im Tab 'Milestones & Epics' je Epic/Milestone ein Knopf 'nach Jira'; eine Vorschau listet, was angelegt oder geaendert wird, und die Zahl der API-Anfragen; erst Bestaetigen sendet; Test gegen Fake-Server
- [ ] Angelegtes traegt danach 'external'; ein zweiter Export legt nichts doppelt an; Test
- [ ] Bricht der Export ab (z.B. 429), bleibt festgehalten, was schon drueben ist; ein neuer Versuch macht dort weiter; Test
- [ ] jaira web --help, eine Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

