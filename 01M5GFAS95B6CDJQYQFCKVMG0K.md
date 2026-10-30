---
id: 01M5GFAS95B6CDJQYQFCKVMG0K
title: "Boards wechseln mit Shift+Pfeil und ein Board-Fenster auf b mit aktiven Boards, Ersteller und Lanes"
status: ideas
ready: true
creator: berkc
assignee: berkc
goal: "Zwischen Boards springt man mit Shift+hoch/runter, und b zeigt alle Boards mit Name, Ersteller und Lanes; dort schaltet die Leertaste ein Board fuer den Teamlead aktiv oder inaktiv."
context: |-
  Wunsch von berkc am 22.10.2026, woertlich: 'pressing shift + Arrow keys (up/down) Rotates through the boards that are active. We said that the active boards are selected in the Settings, but we should also add the hotkey b, which opens up a pop up with all the available boards as long tabs, with the name of the board as a header, creator of the board and the lanes listed below. With the active boards (the one that the team lead can choose from) having a tick symbol. On that view you can select and unselect the boards with the space key. The order of the boards is: active boards, most recently active, all the other boards.'
  Heute: b auf dem Board springt zum naechsten Board (nextBoard(), internal/tui/model.go); im geoeffneten Ticket springt b zum Ticket, von dem es blockiert wird. Das Ticket zum Link-Oeffnen-Fehler will b im Ticket ausserdem fuer 'zeig mir dieses Ticket auf dem Board' - die Belegung von b muss einmal fuer alle drei entschieden werden.
  Verwandt: 57TDND (Settings: aktive Boards fuer den Teamlead waehlen) - das Fenster hier ist eine zweite Stelle fuer dieselbe Auswahl, muss also dieselbe Einstellung schreiben. Y2C9QN (Boards aus Board-Typen). Ein Board hat heute einen description-Datei unter .jaira/boards/<name>/; einen Ersteller speichert es noch nicht.
definition-of-done: Shift+hoch/runter wechselt reihum durch die aktiven Boards; Test in internal/tui
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-22T00:14:17Z
updated-at: 2026-10-30T01:58:32Z
updated-by: berkc
claimed-by: berk-ca-1397863
claimed-at: 2026-10-30T01:57:38Z
---

# Boards wechseln mit Shift+Pfeil und ein Board-Fenster auf b mit aktiven Boards, Ersteller und Lanes

## Definition of Done

- [ ] Shift+hoch/runter wechselt reihum durch die aktiven Boards; Test in internal/tui
- [ ] b oeffnet ein Fenster: pro Board Name als Ueberschrift, Ersteller, Lanes; aktive mit Haken; Reihenfolge aktiv, zuletzt benutzt, Rest; Test
- [ ] Leertaste schaltet aktiv/inaktiv und schreibt dieselbe Einstellung wie 57TDND; Test
- [ ] Die Belegung von b (Board, Blocked-by-Sprung, Ticket auf dem Board zeigen) ist entschieden und in der Hilfe beschrieben
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [x] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-23 21:04 · berkc** — Entscheidung berkc 23.10.2026: Taste b bleibt fuer Boards (Board-Popup). Die Sprung-Funktion aus P3VS02 braucht eine andere Taste.
