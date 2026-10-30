---
id: 01M5GFCCFBGD9YCC2KE7P3VS02
title: "Ein verlinktes Ticket oeffnet sich aus dem Links-Fenster auch dann, wenn der Filter es verbirgt"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "enter im Links-Fenster oeffnet das verlinkte Ticket immer direkt, egal was der Board-Filter zeigt; eine eigene Taste zeigt das Ticket auf dem Board und setzt dafuer den Filter zurueck."
context: |-
  Gemeldet von berkc am 22.10.2026: in einem offenen Ticket l druecken, auf ein verlinktes Ticket gehen, enter - ist es durch den aktuellen Filter verborgen, oeffnet es sich nicht. 'trying to open the link should always open the ticket without having to go over the board, that way the filter is irrelevant. when i am on a ticket and click b it should navigate me to the ticket on the board (resetting the filter this time)'.
  Ursache (geprueft): internal/tui/links.go:127-157. enter setzt mode=modeBoard und ruft m.selectByID; findet der die Karte nicht (Filter), kommt nur die Meldung 'the filter ... hides it. Clear the filter with esc'. Das Detail wird ueber die Board-Auswahl geoeffnet (m.openDetail nach selectByID), deshalb haengt es am Filter.
  Ausserdem lehnt enter Tickets ab, die nicht auf dem Board liegen (Logbook, nur auf dem Ref) - links.go:134-138. Ein Ticket im Logbook laesst sich also gar nicht aus den Links lesen.
  Achtung Tastenbelegung: b im offenen Ticket springt heute zum Blocked-by-Ticket (Release-Notes), b auf dem Board wechselt das Board. Siehe KVMG0K (Board-Fenster auf b). Vorschlag: im Detail eine andere Taste fuer 'auf dem Board zeigen' oder b nur dann Blocked-by, wenn es eines gibt - mit berkc klaeren.
definition-of-done: "enter im Links-Fenster oeffnet das Ticket direkt im Detail, auch wenn der Filter es verbirgt; esc fuehrt zurueck zum vorigen Ticket; Test in internal/tui"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-22T00:15:09Z
updated-at: 2026-10-30T01:56:17Z
updated-by: berkc
---

# Ein verlinktes Ticket oeffnet sich aus dem Links-Fenster auch dann, wenn der Filter es verbirgt

## Definition of Done

- [ ] enter im Links-Fenster oeffnet das Ticket direkt im Detail, auch wenn der Filter es verbirgt; esc fuehrt zurueck zum vorigen Ticket; Test in internal/tui
- [ ] Auch Tickets aus Logbook und Archiv oeffnen sich read-only; Test
- [ ] Eine Taste im offenen Ticket setzt den Filter zurueck und waehlt die Karte auf dem Board; Belegung mit KVMG0K abgestimmt; Test
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-23 21:04 · berkc** — Entscheidung berkc 23.10.2026: b bleibt fuer Boards (KVMG0K). Der Sprung vom Ticket aufs Board (Filter zuruecksetzen) bekommt eine andere Taste - im brainstorm/pre-process vorschlagen.
- **2026-10-27 14:32 · berkc** — 27.10.2026: b im offenen Ticket ist schon belegt - es folgt dem Blocker, wenn das Ticket blockiert ist (internal/tui/model.go:1608). Das Board-Popup (KVMG0K) liegt auf b in der Board-Ansicht (model.go:1705), kollidiert also nicht. Vorschlag: b im Ticket springt aufs Board (Filter zurueck), wenn das Ticket keinen Blocker hat; mit Blocker folgt b weiter dem Blocker - oder der Blocker-Sprung zieht auf eine andere Taste. Entscheidung berkc offen.
- **2026-10-30 01:56 · berkc** — Entscheidung berkc 30.10.2026: b im offenen Ticket springt aufs Board (Filter zuruecksetzen). Der Blocker-Sprung auf b (internal/tui/model.go:1608) entfaellt - Blocker erreicht man ueber l (Links), es geht also nichts verloren. Board-Popup b auf dem Board (KVMG0K) bleibt.
