---
id: 01M6NNAXQ276DMK3A1GH3XXF5B
title: "jaira web: ein Ticket laesst sich per Knopf aus einem Epic nehmen, auch auf dem Handy"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer in jaira web ein Epic ansieht, nimmt jedes Ticket mit einem Klick oder Tipp heraus, ohne auf die CLI ('jaira epic rm') auszuweichen."
context: |-
  Heute nimmt man ein Ticket nur per Drag-and-Drop aus einem Epic: Karte aus der Epic-Liste auf 'Tickets in no epic' ziehen (core/web/ui/src/widgets/epic-list/ui/EpicList.svelte:74, Ziel UnassignedList.svelte).
  Gemeldet von Alex am 05.11.2026: 'manchmal kann ich nicht entfernen, dann muss ich die CLI nehmen'.
  Bekannte Ursachen:
  - Drag-and-Drop geht auf dem Telefon nicht (Touch hat kein HTML5-dragstart). Alex arbeitet gerade vom Telefon aus.
  - Karten mit direct=false (ueber follows/parent mitgezaehlt) sind nicht ziehbar, TicketCard.svelte:23. Sie stehen nicht in der Ticketliste des Epics, also gibt es dort nichts zu entfernen - die Karte muss das sagen und zeigen, welches Ticket sie hineinzieht.
  - Archivierte Tickets (place != board) sind in der Unassigned-Liste nicht ziehbar.
  Die API gibt es schon: DELETE /api/epics/{name}/tickets/{id} (core/web/web.go:280, handleUnassign) - CLI-Gegenstueck 'jaira epic rm'.
  Pruefen, ob ein archiviertes Ticket aus einem Epic genommen werden kann (CLI vs. Web gleich).
definition-of-done: "Jede direkte Ticket-Karte in einem Epic hat einen Entfernen-Knopf, der DELETE /api/epics/{name}/tickets/{id} ruft; funktioniert per Tipp ohne Drag; Test"
tags:
  - web
blocked-by: []
related: []
commits: []
created-at: 2026-11-05T10:51:07Z
updated-at: 2026-11-05T10:53:17Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-85123
claimed-at: 2026-11-05T10:51:26Z
---

# jaira web: ein Ticket laesst sich per Knopf aus einem Epic nehmen, auch auf dem Handy

## Definition of Done

- [ ] Jede direkte Ticket-Karte in einem Epic hat einen Entfernen-Knopf, der DELETE /api/epics/{name}/tickets/{id} ruft; funktioniert per Tipp ohne Drag; Test
- [ ] Eine indirekte Karte (via follows/parent) nennt das Ticket, ueber das sie mitgezaehlt wird, statt still nicht ziehbar zu sein; Test
- [ ] Ein archiviertes Ticket laesst sich im Web genauso aus dem Epic nehmen wie mit 'jaira epic rm', oder die Seite sagt warum nicht; Test
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/epic: Ticket.Via = the member a follow-up is counted through; Covered/Count carry it; Go test
- [ ] core/web + internal/cli: web.Ticket.Via passed through; webSource.load uses LoadAnywhere so archived/logbook tickets can be removed like 'jaira epic rm'; Go test
- [ ] UI: Remove button on direct cards in EpicList calling DELETE via onunassign; indirect card text names the via ticket; vitest
- [ ] release note core/release/notes/3XXF5B.md; build ui dist; go test ./...

## Progress
- **2026-11-05 10:51 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/1770.sock
- **2026-11-05 10:52 · Alexander Sacharov** — decisions (none open for the person, settled from code): (1) Remove button = small text button 'Remove' on every direct card in an epic, no confirm dialog - drag has none either and 'add' undoes it; aria-label names ticket and epic. (2) An indirect card names the epic member it is counted through (the line you would remove), found by walking follows/parent up in epic.Covered; new field 'via' on epic.Ticket and web.Ticket. (3) Archived: root cause found - webSource.load uses Store.Load (board only), 'jaira epic rm' resolves via LoadAnywhere (board+logbook+archive). Web switches to LoadAnywhere, so web and CLI behave the same; no 'why not' message needed.
