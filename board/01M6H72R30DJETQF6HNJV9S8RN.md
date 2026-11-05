---
id: 01M6H72R30DJETQF6HNJV9S8RN
title: "jaira web bekommt eine Tab-Leiste; Tab 'Alle Tickets' zeigt jedes offene Ticket in der Spalte seiner Lane"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Ein PM sieht im Browser alle offenen Tickets nach Lane geordnet und wechselt per Tab zwischen den Ansichten, die die Folge-Tickets fuellen."
context: |-
  Was heute fehlt: jaira web hat genau eine Ansicht - Epics rechts, Tickets ohne Epic links (core/web/web.go:300-315). Ein Ticket, das in einem Epic steckt, sieht man nur dort.
  Es gibt keine Stelle fuer weitere Ansichten; die Folge-Tickets (Aktiv, Milestones & Epics, Schreiben, Braucht dich) brauchen eine.
  Teil der Aufteilung von CMN1YT (jaira web als Zentrale fuer PMs = Leute ohne Claude und ohne CLI); dort stehen Brainstorm und Antworten des Maintainers.
  Bekannt: Filterleiste (Suche, Lane, Tag, Person, Datum) und Lane-Farben gibt es schon (Hilfetext internal/cli/web.go). Wiederverwenden, nicht neu bauen.
  Entschieden: ein Tab erscheint erst, wenn sein Ticket gebaut ist - keine leeren Tabs.
  Ausgeschlossen: Tickets bearbeiten (eigenes Ticket 'Schreiben').
definition-of-done: "Tab-Leiste oben; die heutige Epic-Ansicht liegt unveraendert unter einem Tab 'Epics'; Test in core/web"
tags:
  - web
blocked-by: []
related: []
follows: 01M67K6Z3NSBTBXAK4EVCMN1YT
commits: []
created-at: 2026-11-03T17:25:01Z
updated-at: 2026-11-03T17:25:01Z
---

# jaira web bekommt eine Tab-Leiste; Tab 'Alle Tickets' zeigt jedes offene Ticket in der Spalte seiner Lane

## Definition of Done

- [ ] Tab-Leiste oben; die heutige Epic-Ansicht liegt unveraendert unter einem Tab 'Epics'; Test in core/web
- [ ] Tab 'Alle Tickets': jedes offene Ticket des Boards genau einmal, in Spalten in der Lane-Reihenfolge von /api/lanes, mit derselben Filterleiste; Test fuer einen neuen GET-Endpunkt
- [ ] Der gewaehlte Tab steht in der Adresse; ein Link oeffnet denselben Tab
- [ ] jaira web --help und docs/COMMANDS.md beschreiben die Tabs; eine Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

