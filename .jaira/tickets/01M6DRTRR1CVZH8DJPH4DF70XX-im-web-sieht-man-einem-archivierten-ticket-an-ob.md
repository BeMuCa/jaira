---
id: 01M6DRTRR1CVZH8DJPH4DF70XX
title: "Im Web sieht man einem archivierten Ticket an, ob es fertig, liegen gelassen oder abgelehnt ist, und das Archiv laesst sich filtern"
status: critique
ready: true
creator: Alexander Sacharov
goal: "Ein Ticket im Archiv zeigt im Web auf einen Blick, wie es das Board verlassen hat: fertig (done), ohne Abschluss abgelegt (andere Lane), oder abgelehnt (jaira reject, mit Grund). Das Archiv laesst sich getrennt vom Logbuch filtern."
context: |-
  Was heute falsch ist (gefunden von Alex am 02.11.2026 in jaira web):
  - Nur Tickets mit done=true werden durchgestrichen: core/web/ui/src/entities/ticket/ui/TicketCard.svelte:63 (.done .title line-through).
  - Ein Ticket im Archiv, das NICHT aus done kam, sieht aus wie ein normales Ticket - man sieht nicht, dass es weggelegt wurde.
  - badgeLane (core/web/ui/src/entities/ticket/model/place.ts) zeigt fuer jedes abgelegte Ticket nur 'archive' und versteckt die Lane, aus der es kam.
  - Ein abgelehntes Ticket (jaira reject) bekommt nur den Text ' · rejected: <art>' in der kleinen grauen where-Zeile (TicketCard.svelte:32). Kein eigenes Aussehen, der Grund (Rejection in types.ts) steht nicht auf der Karte.
  - Filter: FilterBar.svelte:32 hat nur 'All tickets ever' (Logbuch UND Archiv zusammen). Archiv allein, Logbuch allein, oder 'nur abgelehnte' geht nicht.
  Was Alex will:
  - drei unterscheidbare Zustaende im Archiv: fertig / liegen gelassen / abgelehnt
  - abgelehnte 'richtig' markiert, nicht nur als Nebentext
  - das Archiv filtern koennen
  Offen fuer brainstorm: wie genau die drei Zustaende aussehen (durchgestrichen, gedimmt, Badge, Farbe) und welche Filter-Optionen es gibt.
definition-of-done: Ein archiviertes Ticket aus done ist durchgestrichen; Test
tags:
  - web
blocked-by: []
related: []
commits: []
created-at: 2026-11-02T09:18:16Z
updated-at: 2026-11-02T09:27:13Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-58445
claimed-at: 2026-11-02T09:18:50Z
outcome-what: "Karten im Web zeigen archivierte Tickets als fertig (durchgestrichen), liegen gelassen (gedimmt, Badge mit Herkunfts-Lane) oder abgelehnt (roter Rand, Zeile Art – Grund); die Checkbox 'All tickets ever' ist durch einen Where-Chip-Filter (board/logbook/archive/rejected) ersetzt, ?all=1 liest als alle vier."
outcome-why: Im Archiv sah ein liegen gelassenes Ticket wie ein offenes aus und ein abgelehntes nur als grauer Nebentext; Logbuch und Archiv liessen sich nicht getrennt zeigen.
outcome-resolves: DF70XX
---

# Im Web sieht man einem archivierten Ticket an, ob es fertig, liegen gelassen oder abgelehnt ist, und das Archiv laesst sich filtern

## Definition of Done

- [x] Ein archiviertes Ticket aus done ist durchgestrichen; Test
  proof: TicketCard.svelte:21 (class:done via fate); place.test.ts 'an archived ticket is finished, left behind or rejected...' (archived from done -> done)
- [x] Ein archiviertes Ticket aus einer anderen Lane sieht anders aus als ein offenes und als ein fertiges, und nennt die Lane, aus der es kam; Test
  proof: TicketCard.svelte:22,71 (.left gedimmt) + place.ts badgeLane (archive -> Herkunfts-Lane); place.test.ts 'an archived ticket shows the lane it came from' und fate 'left'
- [x] Ein abgelehntes Ticket hat ein eigenes Aussehen und zeigt Art und Grund der Ablehnung auf der Karte; Test
  proof: TicketCard.svelte:23,39,72 (roter Rand, Zeile kind – why); web_source.go:153 rejected_why; TestWebListsFiledTicketsAndOpensThem (internal/cli/web_test.go:268); place.test.ts fate 'rejected'
- [x] Der Filter kann Board, Logbuch, Archiv und abgelehnte getrennt zeigen; Test in filter.test
  proof: filter.test.ts:75 'the place filter keeps the board, the logbook, the archive and the rejected apart', :93, :101 (?all=1)
- [x] Release-Note-Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/DF70XX.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-02 09:18 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/43623.sock
- **2026-11-02 09:20 · Alexander Sacharov** — Code gelesen: Karte core/web/ui/src/entities/ticket/ui/TicketCard.svelte, Payload internal/cli/web_source.go cards() (rejected traegt nur kind, why fehlt in der Liste -> Go-Aenderung noetig), Filter features/filter-tickets/model/filter.ts (f.all -> ?all=1 an /unassigned). Offene Entscheidungen vor dem Plan: (1) Aussehen liegen-gelassen vs abgelehnt, (2) Form des Ort-Filters. Fragen gehen an den parent.
- **2026-11-02 09:20 · Alexander Sacharov** — Frage 1 (Aussehen im Archiv; fertig = durchgestrichen steht schon im DoD): 1) (Empfohlen) Liegen gelassen: gedimmt + Badge zeigt die Lane, aus der es kam; abgelehnt: roter linker Rand + Zeile 'abgelehnt: <Art> - <Grund>' in Rot. 2) Beide gedimmt, nur ein Text-Badge 'liegen gelassen' / 'abgelehnt' unterscheidet sie. 3) Abgelehnt rot durchgestrichen, liegen gelassen kursiv-gedimmt.
Frage 2 (Filter statt 'All tickets ever'): 1) (Empfohlen) Chip-Auswahl 'Wo' mit board / logbook / archive / rejected, mehrere waehlbar, Standard nur board; alte Links mit ?all=1 heissen alle vier. 2) Ein Dropdown, genau ein Ort zur Zeit (board | logbook | archive | rejected | alles). 3) Checkbox bleibt, daneben eine zusaetzliche Auswahl nur fuer gefilterte Orte.
- **2026-11-02 09:23 · Alexander Sacharov** — decision (Alex via teamlead, 02.11.): Frage 1 = 1: fertig durchgestrichen; liegen gelassen gedimmt, Badge zeigt die Herkunfts-Lane; abgelehnt roter linker Rand + rote Zeile 'rejected: <Art> - <Grund>'. Frage 2 = 1: Chip-Auswahl 'Where' board/logbook/archive/rejected, mehrere waehlbar, ohne Chip nur board; ?all=1 = alle vier. Mode bleibt autonom.
- **2026-11-02 09:24 · Alexander Sacharov** — Methode: Go web.Ticket bekommt rejected_why (web_source.go cards()); place.ts bekommt fate(t) = open|done|left|rejected und whereOf(t) = board|logbook|archive|rejected (ref zaehlt als board); badgeLane zeigt fuer archive die Herkunfts-Lane, logbook bleibt 'archive'; TicketCard: class left (gedimmt) / rejected (roter Rand + Zeile); filter.ts: 'where' als vierter Chip-Filter statt all, ?all=1 liest als alle vier, Server wird mit all=1 gefragt sobald ein Chip etwas anderes als board meint. Ohne Where-Chip filtert matches() nicht nach Ort, damit Epics wie heute ihre abgelegten Tickets zeigen.
- **2026-11-02 09:27 · Alexander Sacharov** — whereOf/WHERES liegen in features/filter-tickets/model/filter.ts, nicht in entities/ticket: filter.ts darf zur Laufzeit nichts ueber '@/' importieren, sonst bricht 'node --test' (nur import type wird gestrippt). Ohne Where-Chip filtert matches() nicht nach Ort, damit Epics ihre abgelegten Mitglieder weiter zeigen; ein Where-Chip filtert auch die Epics. dist neu gebaut (eingecheckt, go:embed). tsc-Fehler in vite.config.ts (node:url) ist alt, nicht von hier.
