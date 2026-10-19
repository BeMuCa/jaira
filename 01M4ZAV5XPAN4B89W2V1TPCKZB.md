---
id: 01M4ZAV5XPAN4B89W2V1TPCKZB
title: "Web-Seite: Filter und Tickets aus Logbook und Archiv anzeigen"
status: in-progress
ready: true
creator: Alexander Sacharov
goal: "Auf der Seite von 'jaira web' sieht man jedes Ticket, das es je gab — offen, nur auf dem Ref, im Logbuch oder im Archiv — und filtert es nach Text, Lane, Tag, Person und Datum, mit den Filtern in der URL."
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
updated-at: 2026-10-19T06:40:24Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-74288
claimed-at: 2026-10-19T06:35:35Z
---

# Web-Seite: Filter und Tickets aus Logbook und Archiv anzeigen

## Definition of Done

- [ ] Auf der Seite gibt es ein Textfeld, das Tickets nach Titel, Ziel und ID filtert
- [ ] Filter nach Lane, Tag und Person (creator, assignee, zuletzt geaendert von) - mehrere Filter wirken zusammen (UND)
- [ ] Tickets aus .jaira/logbook/ sind auf der Seite sichtbar und öffnen sich in der Detailansicht
- [ ] Tickets aus .jaira/archive/ sind auf der Seite sichtbar und öffnen sich in der Detailansicht
- [ ] core/release/notes/<ID>.md beschreibt die neuen Filter und Ansichten
- [ ] Filter nach Datum: Zeitraum von-bis auf erstellt, zuletzt geaendert und abgelegt (Logbuch-Datum)
- [ ] Ein abgelehntes Ticket im Archiv zeigt rejected-kind und rejected-why
- [ ] Die gesetzten Filter stehen in der URL; ein Link zeigt dieselbe gefilterte Liste

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/web: Ticket bekommt place, tags, creator, assignee, updated_by, created_at, updated_at, filed_at, filed; GET /api/unassigned?all=1; TicketDetail bekommt place + rejected
- [x] internal/cli/web_source.go: Karten aus Board+Ref+Logbuch+Archiv fuellen, rejected-* ueber Doc.Fields lesen, Tests
- [~] UI: features/filter-tickets (Text, Lane, Tag, Person, Datum, Schalter), Zustand im URL-Query, wirkt auf beide Spalten; Ort-Badge; abgelegte nicht ziehbar; Detail zeigt Ort und Ablehnung
- [ ] UI bauen (dist committed), web --help und core/release/notes/TPCKZB.md

## Progress
- **2026-10-19 06:35 · Alexander Sacharov** — parent session: jaira-a7
- **2026-10-19 06:35 · Alexander Sacharov** — decision (Alex, 19.10.2026 via teamlead jaira-a7): 'we must see ALL tickets that ever existed and filter by people, date or lane.' Abgelegte Tickets (logbook, archive, refs) kommen in DIESELBE Liste mit einem Schalter 'alle Tickets aller Zeiten'; jedes zeigt, wo es liegt — keine eigene Ansicht. Damit sind die brainstorm-Fragen beantwortet. Offen fuer den Plan: genaue Personenfelder, Server-API vs. Filtern im Browser. core/ticket/store.go Archive/Restore nicht anfassen — A0C081 aendert sie; rejected-kind/why/by/at (flache Keys aus A0C081) anzeigen, wenn vorhanden, nicht auf den Merge warten.
- **2026-10-19 06:37 · Alexander Sacharov** — brainstorm — was der Code zeigt:
- core/web/web.go hat nur /api/epics und /api/unassigned; internal/cli/web_source.go baut sie aus boardIndex + epic.Unassigned — also nur Tickets unter .jaira/tickets/ und ref-only (Store.List haengt Source.Extra an).
- Die Detailansicht kann abgelegte Tickets schon: webSource.Ticket nutzt Store.LoadAnywhere (Board, dann logbook, dann archive). Es fehlt nur ein Weg dorthin.
- Filed-Tickets findet Store.FiledAwayIDs() (archive + alle logbook-Ordner, auch der Legacy-Name); Logbuch-Datum steckt im Ordnernamen, ticket.LogbookFolderDay liest es.
- rejected-* gibt es auf master noch nicht (A0C081 ist uncommitted in ../jaira-a0c081, fuegt Ticket.Rejected in core/ticket/schema.go hinzu). Lesbar ohne schema.go anzufassen: Doc.Fields().Scalar("rejected-kind").
- Groesse heute: 134 abgelegte + ~29 offene Tickets.

Wege:
A) Ein Endpoint GET /api/tickets?all=1 liefert jede Karte mit den Filterfeldern (lane, tags, creator, assignee, updated_by, created/updated/filed, place); gefiltert wird im Browser, Filter im URL-Query. Kostet: ein Endpoint, eine Filterleiste. Gibt: sofortiges Filtern, kein Request pro Tastendruck; bei tausenden Tickets wird die Antwort gross.
B) Filter serverseitig (?q=&tag=&lane=...), wie 'jaira list'. Kostet: Filterlogik in Go doppelt zu 'jaira list' oder deren Refactor; jeder Tastendruck ein Request mit Board-Scan. Gibt: kleine Antworten.
C) Eigene Seite /all. Von Alex ausgeschlossen (gleiche Liste, Schalter).

Ich nehme A: die Datenmenge ist klein, Filtern im Browser ist sofort, die URL traegt den Zustand ohne Server-Logik.
- **2026-10-19 06:37 · Alexander Sacharov** — decision (worker, aus Code/Default): Die Filterleiste steht ueber beiden Spalten und wirkt auf beide — auf 'Tickets in no epic' und auf die Ticketzeilen jeder Epic-Karte. Der Schalter 'alle Tickets aller Zeiten' fuegt der linken Liste die abgelegten Tickets ohne Epic hinzu (Epic-Karten zaehlen abgelegte schon heute). Jede Karte zeigt, wo sie liegt: Board / Ref / Logbuch / Archiv. Abgelegte Karten lassen sich nicht ziehen. Personenfilter: EIN Feld 'Person', trifft wenn creator ODER assignee ODER updated_by passt — 'filter by people' ohne drei Dropdowns. Datumsfilter: Auswahl des Feldes (erstellt / geaendert / abgelegt) plus von/bis. Abgelegt-Datum: Logbuch = Datum aus dem Ordnernamen; Archiv = rejected-at, wenn gesetzt, sonst keins. Die Snapshot-Branch (core/snapshot) ist ein Backup und bleibt aussen vor.
