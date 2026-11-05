---
id: 01M6DRTRR1CVZH8DJPH4DF70XX
title: "Im Web sieht man einem archivierten Ticket an, ob es fertig, liegen gelassen oder abgelehnt ist, und das Archiv laesst sich filtern"
status: done
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
commits:
  - 510c47d5fb15e817b0cf3f0d1b1ee5f6579fd699
  - 33c52f8ba15187ff4357a633b4b8955b5f566b27
  - e7fab6805b8fa1d305dea702c8d52b73184d1fd0
created-at: 2026-11-02T09:18:16Z
updated-at: 2026-11-02T10:24:31Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-58445
claimed-at: 2026-11-02T09:18:50Z
outcome-what: "Karten im Web zeigen archivierte Tickets als fertig (durchgestrichen), liegen gelassen (gedimmt, Badge mit Herkunfts-Lane) oder abgelehnt (roter Rand, Zeile Art – Grund); die Checkbox 'All tickets ever' ist durch einen Where-Chip-Filter (board/logbook/archive/rejected) ersetzt, ?all=1 liest als alle vier."
outcome-why: Im Archiv sah ein liegen gelassenes Ticket wie ein offenes aus und ein abgelehntes nur als grauer Nebentext; Logbuch und Archiv liessen sich nicht getrennt zeigen.
outcome-resolves: DF70XX
review-summary: "Die Web-Karte unterscheidet jetzt drei Ausgaenge: place.ts fate() liefert open/done/left/rejected; TicketCard streicht fertige durch, dimmt liegen gelassene (opacity 0.6, Zusatz 'left unfinished', Badge zeigt die Herkunfts-Lane statt 'archive') und gibt abgelehnten einen roten linken Rand plus rote Zeile 'rejected: <Art> – <Grund>'. Den Grund liefert Go neu als rejected_why (web.go, web_source.go cards()). Die Checkbox 'All tickets ever' ist durch einen vierten Chip-Filter 'Where' (board/logbook/archive/rejected, auch 'not X') ersetzt; der Server wird nur mit all=1 gefragt, wenn ein Chip mehr als board meint (wantsFiled), und alte Links ?all=1 lesen als alle vier Orte. dist ist neu gebaut und reproduzierbar."
round: 3
review-gaps: |-
  note: ohne Where-Chip zeigen Epics weiter ihre abgelegten Mitglieder - die Entscheidung sagte 'ohne Chip nur board'; bewusst so (Notiz 09:27), gilt aber nur fuer die Liste ohne Epic
  note: TicketDetails.svelte nutzt badgeLane auch - die Detailansicht eines archivierten Tickets zeigt jetzt ebenfalls die Herkunfts-Lane statt 'archive'; konsistent, aber nicht im DoD genannt
  note: die Optik selbst (durchgestrichen, gedimmt, roter Rand) ist nur ueber fate() getestet, nicht ueber die gerenderte Karte; testing hat sie per Screenshot geprueft
test-verdict: "pass: go vet + go test ./... -race RC=0 (33 Pakete), npm test 16/16, vite build ok; DoD 1-5 im Baum geprueft; Scratch-Board mit offen/liegen gelassen/abgelehnt per jaira web + headless chrome: durchgestrichen/gedimmt mit Lane backlog/roter Rand mit 'rejected: wontfix – <Grund>', ?where=rejected zeigt nur das abgelehnte, ?where=archive&where=rejected beide archivierten"
review-verdict: "Der Diff erfuellt alle fuenf DoD-Punkte; npm test (16/16) und go test ./internal/cli -run TestWeb gruen, vite build reproduziert den eingecheckten dist byte-gleich. Keine Defekte gefunden."
review-round: 3
review-check: "1. cd core/web/ui && npm test  -> 16 pass, darunter 'the place filter keeps the board, the logbook, the archive and the rejected apart'  2. go run ./cmd/jaira web in einem Repo mit je einem Ticket im Logbuch, einem per 'jaira archive' aus backlog abgelegten und einem per 'jaira reject --force' abgelehnten  3. im Filter 'Where' archive und rejected waehlen  4. das abgelegte ist gedimmt, Badge 'backlog', Zeile 'in the archive · left unfinished'  5. das abgelehnte hat roten linken Rand und rote Zeile 'rejected: wontfix – <Grund>'  6. Where = logbook: Titel durchgestrichen, Badge 'archive'  7. URL mit ?all=1 oeffnen -> alle vier Where-Chips sind gesetzt"
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
- **2026-11-02 09:27 · Alexander Sacharov** — critique 1: fate/Fate-Export aus entities/ticket/index.ts ist ungenutzt - streichen. Bewusst stehen gelassen: ein alter Link ?all=1 waehlt alle vier Orte, narrows() ist dann wahr und Epics ohne Ticket verschwinden aus der Liste, 'Clear' erscheint - harmlos, Ort-Filter ist ein echter Filter. Abgelehnte Karte hat zwei Zeilen ('in the archive' grau + rote Grund-Zeile) - gewollt, Ort und Grund sind zwei Aussagen.
- **2026-11-02 09:28 · Alexander Sacharov** — critique 2: nichts mehr; Entscheidungen aus critique 1 bleiben stehen.
- **2026-11-02 09:29 · Alexander Sacharov** — testing 1: Screenshot auf Scratch-Board (headless chrome, ?where=board&where=archive&where=rejected): drei Zustaende sichtbar und richtig, aber Svelte trimmt das Leerzeichen vor {#if} in TicketCard.svelte:36/39 -> 'wontfix– nobody needs it' und 'in the archive· left unfinished'. Fix: Trenner als {' – '} / {' · '} schreiben.
- **2026-11-02 09:30 · Alexander Sacharov** — critique 3 / optimize 2: Trenner-Fix ist zwei String-Literale, nichts zu sagen.
- **2026-11-02 09:30 · Alexander Sacharov** — testing 2: Suite lief auf 33c52f8b; e7fab680 aendert nur TicketCard.svelte + dist, danach npm test + build + Screenshot erneut. Screenshots: scratchpad shot2.png/shot3.png (lokal).
- **2026-11-02 10:06 · Alexander Sacharov** — explain: Der Review fand keinen Defekt, nur drei Hinweise: Epics zeigen ohne Where-Chip weiter abgelegte Mitglieder, die Detailansicht zeigt nun die Herkunfts-Lane statt 'archive', und die Optik ist nur ueber fate() getestet, nicht ueber die gerenderte Karte.
Warum: Alle fuenf DoD-Punkte sind erfuellt und die Tests sind gruen; die Hinweise sind Randfaelle ohne Folgen fuer das Ziel, ein Archiv-Ticket auf einen Blick einzuordnen.
1) Zurueckgeben und die Epic-Liste ohne Chip auf 'nur board' bringen, dazu einen Test fuer die gerenderte Karte [give back]
2) Das Ticket so annehmen; die drei Hinweise bleiben als bekannte Randfaelle stehen [accept]
Empfehlung: 2 - die Epic-Entscheidung ist bewusst getroffen (Notiz 09:27), der Rest ist konsistent und per Screenshot geprueft.
- **2026-11-02 10:23 · Alexander Sacharov** — decision: 2) Das Ticket so annehmen; die drei Hinweise bleiben als bekannte Randfaelle stehen
