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
updated-at: 2026-10-19T06:35:50Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-74288
claimed-at: 2026-10-19T06:35:35Z
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
- **2026-10-19 06:35 · Alexander Sacharov** — parent session: jaira-a7
- **2026-10-19 06:35 · Alexander Sacharov** — decision (Alex, 19.10.2026 via teamlead jaira-a7): 'we must see ALL tickets that ever existed and filter by people, date or lane.' Abgelegte Tickets (logbook, archive, refs) kommen in DIESELBE Liste mit einem Schalter 'alle Tickets aller Zeiten'; jedes zeigt, wo es liegt — keine eigene Ansicht. Damit sind die brainstorm-Fragen beantwortet. Offen fuer den Plan: genaue Personenfelder, Server-API vs. Filtern im Browser. core/ticket/store.go Archive/Restore nicht anfassen — A0C081 aendert sie; rejected-kind/why/by/at (flache Keys aus A0C081) anzeigen, wenn vorhanden, nicht auf den Merge warten.
