---
id: 01M5ZJ5FMXTVR0AS5FT77QDQ4W
title: "Das Cockpit-Band frischt sich im Leerlauf auf, nicht nur bei einem Ereignis"
status: backlog
ready: false
creator: BeMuCa
goal: "Das Band im Cockpit liest den Board-Stand auch dann neu, wenn kein Hook-Ereignis kommt, etwa alle 60 Sekunden; ein Ticket, das eine andere Session bewegt hat, erscheint ohne eigenen Befehl."
context: |-
  Was schief ist: ohne eigenes Ereignis blieb das Band stehen; Tickets, die Driver in anderen Sessions bewegten, erschienen erst nach dem naechsten eigenen Befehl (O18, 27.10.2026).
  Heute: core/mod/builtin/jaira-cockpit/hooks/register.tsx pollt um Zeile 281 alle 60s - was genau es neu liest und was das Band daraus macht, ist zu pruefen; im Leerlauf sah der Mensch keinen Wechsel.
  Zusammenhang: D1S55M (Board holt alle 2 Minuten, signoff) - das Board, nicht das Band. Q63EEV (Pane je Ticket, human) bleibt getrennt.
definition-of-done: "Ein Lane-Wechsel aus einer anderen Session steht binnen zwei Minuten im Band, ohne eigenen Befehl"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-27T20:52:25Z
updated-at: 2026-10-27T20:52:25Z
---

# Das Cockpit-Band frischt sich im Leerlauf auf, nicht nur bei einem Ereignis

## Definition of Done

- [ ] Ein Lane-Wechsel aus einer anderen Session steht binnen zwei Minuten im Band, ohne eigenen Befehl
- [ ] Der Poll laeuft nicht, wenn kein Cockpit sichtbar ist
- [ ] Release-Notiz in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

