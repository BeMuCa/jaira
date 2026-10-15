---
id: 01M4ZAV5XPAN4B89W2V1TPCKZB
title: "Web-Seite: Filter und Tickets aus Logbook und Archiv anzeigen"
status: backlog
ready: false
creator: Alexander Sacharov
goal: "Auf der Seite von 'jaira web' lassen sich Tickets filtern, und abgelegte Tickets aus .jaira/logbook/ und .jaira/archive/ sind dort sichtbar."
context: |-
  Was heute fehlt:
  - 'jaira web' zeigt nur offene Tickets: die Liste 'Tickets in no epic' und die Epics (core/web/web.go, Routen /api/epics und /api/unassigned).
  - Es gibt keinen Filter auf der Seite. Wer ein Ticket sucht, scrollt.
  - Tickets, die 'jaira logbook' abgelegt hat (.jaira/logbook/<wer>-<datum>/), und Tickets in .jaira/archive/ sind auf der Seite gar nicht zu sehen.
  Auslöser: Alex am 2026-10-15 — will alle Tickets aus logbook und archive auf der Seite ansehen können.
  Bekannt:
  - Die UI liegt in core/web/ui/src (Svelte, widgets/ und features/); dist ist eingebettet.
  - 'jaira list' kennt --query, --tag, --milestone, --lane, --assignee — dieselben Filter wären für die Seite naheliegend.
  Offen für brainstorm: welche Filter zuerst (Text, Tag, Lane, Assignee?), und ob logbook/archive eine eigene Ansicht oder ein Schalter 'auch abgelegte zeigen' sind.
definition-of-done: "Auf der Seite gibt es ein Textfeld, das Tickets nach Titel, Ziel und ID filtert"
tags:
  - web
blocked-by: []
related: []
commits: []
created-at: 2026-10-15T08:28:46Z
updated-at: 2026-10-15T08:28:46Z
---

# Web-Seite: Filter und Tickets aus Logbook und Archiv anzeigen

## Definition of Done

- [ ] Auf der Seite gibt es ein Textfeld, das Tickets nach Titel, Ziel und ID filtert
- [ ] Auf der Seite lässt sich nach Tag filtern
- [ ] Tickets aus .jaira/logbook/ sind auf der Seite sichtbar und öffnen sich in der Detailansicht
- [ ] Tickets aus .jaira/archive/ sind auf der Seite sichtbar und öffnen sich in der Detailansicht
- [ ] core/release/notes/<ID>.md beschreibt die neuen Filter und Ansichten

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

