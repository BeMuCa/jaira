---
id: 01M4ZAV5XPAN4B89W2V1TPCKZB
title: "Web-Seite: Filter und Tickets aus Logbook und Archiv anzeigen"
status: critique
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
updated-at: 2026-10-19T06:46:22Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-74288
claimed-at: 2026-10-19T06:35:35Z
outcome-what: "jaira web: Filterleiste (Text, Lane, Tag, Person, Datumsbereich) ueber beiden Spalten, Schalter 'All tickets ever' listet Board, Ref, Logbuch und Archiv mit Ort; Detail zeigt Ort und Ablehnung; Filter im URL-Query"
outcome-why: "Alex will jedes Ticket, das es je gab, auf der Seite sehen und nach Personen, Datum und Lane filtern; bisher zeigte die Seite nur offene Tickets ohne Filter"
outcome-resolves: "GET /api/unassigned?all=1 und neue Kartenfelder (web_source.go cards/rejection), Filtermodell features/filter-tickets mit URL-Abgleich, TestWebListsFiledTicketsAndOpensThem; DoD 1-8 mit Proofs"
review-summary: none
round: 2
---

# Web-Seite: Filter und Tickets aus Logbook und Archiv anzeigen

## Definition of Done

- [x] Auf der Seite gibt es ein Textfeld, das Tickets nach Titel, Ziel und ID filtert
  proof: core/web/ui/src/features/filter-tickets/model/filter.ts:71 (title/goal/handle/id, any case), FilterBar.svelte Search field
- [x] Filter nach Lane, Tag und Person (creator, assignee, zuletzt geaendert von) - mehrere Filter wirken zusammen (UND)
  proof: core/web/ui/src/features/filter-tickets/model/filter.ts:70-84 matches() ANDs lane, tag and person (creator|assignee|updated_by); FilterBar.svelte selects
- [x] Tickets aus .jaira/logbook/ sind auf der Seite sichtbar und öffnen sich in der Detailansicht
  proof: internal/cli/web_source.go:120 cards() reads the logbook; TestWebListsFiledTicketsAndOpensThem; live: ?all=1 on this repo lists 129 logbook tickets, 8NHD8Q opens with place logbook
- [x] Tickets aus .jaira/archive/ sind auf der Seite sichtbar und öffnen sich in der Detailansicht
  proof: internal/cli/web_source.go:120 cards() reads the archive; TestWebListsFiledTicketsAndOpensThem; live: 5 archive tickets listed
- [x] core/release/notes/<ID>.md beschreibt die neuen Filter und Ansichten
  proof: core/release/notes/TPCKZB.md
- [x] Filter nach Datum: Zeitraum von-bis auf erstellt, zuletzt geaendert und abgelegt (Logbuch-Datum)
  proof: core/web/ui/src/features/filter-tickets/model/filter.ts:78 date range on created/updated/filed; filed_at from ticket.LogbookFolderDay (web_source.go cards) or rejected-at
- [x] Ein abgelehntes Ticket im Archiv zeigt rejected-kind und rejected-why
  proof: internal/cli/web_source.go:165 rejection(); TicketDetails.svelte:23; TestWebListsFiledTicketsAndOpensThem checks kind and why
- [x] Die gesetzten Filter stehen in der URL; ein Link zeigt dieselbe gefilterte Liste
  proof: core/web/ui/src/features/filter-tickets/model/filter.ts fromQuery/toQuery; ReviewPage.svelte:26 reads, :40 writes via history.replaceState

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/web: Ticket bekommt place, tags, creator, assignee, updated_by, created_at, updated_at, filed_at, filed; GET /api/unassigned?all=1; TicketDetail bekommt place + rejected
- [x] internal/cli/web_source.go: Karten aus Board+Ref+Logbuch+Archiv fuellen, rejected-* ueber Doc.Fields lesen, Tests
- [x] UI: features/filter-tickets (Text, Lane, Tag, Person, Datum, Schalter), Zustand im URL-Query, wirkt auf beide Spalten; Ort-Badge; abgelegte nicht ziehbar; Detail zeigt Ort und Ablehnung
- [x] UI bauen (dist committed), web --help und core/release/notes/TPCKZB.md

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
- **2026-10-19 06:40 · Alexander Sacharov** — Go-Seite steht: web.Ticket traegt place/tags/creator/assignee/updated_by/created_at/updated_at/filed_at/rejected, GET /api/unassigned?all=1 liefert zusaetzlich fertige Board-Tickets und alle abgelegten ohne Epic; TicketDetail traegt place und rejected. rejected-* wird in internal/cli/web_source.go:rejection per Doc.Fields gelesen — schema.go bleibt fuer A0C081 frei; nach dessen Merge kann rejection() auf t.Rejected umgestellt werden. Test: TestWebListsFiledTicketsAndOpensThem. Als naechstes UI.
- **2026-10-19 06:43 · Alexander Sacharov** — UI fertig. Kein JS-Testrunner im Repo (package.json hat nur vite build); filter.ts wurde mit einem Wegwerf-Skript unter node --experimental-strip-types gegen 15 Faelle geprueft (Text, Lane, Tag, Person, Datum je Feld, URL-Roundtrip, ohne filed_at). svelte-check: 0 Warnungen; der eine Fehler in vite.config.ts (node:url-Typen) besteht schon auf master. Live auf diesem Repo: ?all=1 = 222 Tickets in 0,19 s. Die Seite selbst wurde nicht im Browser angesehen — das ist fuer die Abnahme.
- **2026-10-19 06:45 · Alexander Sacharov** — critique (Runde 1): 1) filedPlace in internal/cli/web_source.go ist eine Kopie von placeOf in core/link/link.go:196 — exportieren (link.PlaceOf) statt kopieren. 2) Ortsbezeichnungen stehen doppelt (TicketCard PLACES, TicketDetails if-Kette) — eine Map in entities/ticket/model/place.ts. Stehen gelassen: cards() liest bei jedem /api/epics alle abgelegten Header (gemessen 0,08 s bei 134 Dateien); rejection() liest die Archivdatei ein zweites Mal — beides klein, und rejection() wird nach A0C081 durch t.Rejected ersetzbar; die Doppelung epic.Unassigned / everything-Zweig bleibt, damit die Standardliste genau 'jaira epic ls --unassigned' entspricht.
- **2026-10-19 06:46 · Alexander Sacharov** — critique Runde 1 umgesetzt: link.PlaceOf exportiert (core/link/link.go:196) und in web_source.go benutzt; Ortsbezeichnungen als PLACE_LABELS in entities/ticket/model/place.ts, TicketCard und TicketDetails lesen sie. Proofs 3/4/7 auf neue Zeilen gesetzt.
- **2026-10-19 06:46 · Alexander Sacharov** — critique Runde 2: beide Befunde aus Runde 1 erledigt (filedPlace weg, link.PlaceOf benutzt; PLACE_LABELS in beiden Komponenten). Nichts weiter — gemaess Lane-Regel nur diese Stellen gelesen.
