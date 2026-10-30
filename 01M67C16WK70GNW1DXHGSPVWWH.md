---
id: 01M67C16WK70GNW1DXHGSPVWWH
title: "Tickets, die in einer Agenten-Lane ohne lebenden Claim stehen bleiben, werden gemeldet"
status: ideas
ready: true
creator: berkc
assignee: berkc
goal: "Ein Ticket, das in einer Lane liegt, die ein Agent abarbeiten muss, und das niemand bearbeitet, faellt auf, statt still liegen zu bleiben."
context: |-
  Frage von berkc am 30.10.2026: 'Sollte man einen Mechanismus einbringen, der periodisch schaut, ob Tickets in Lanes warten, wo sie nicht ohne Aufsicht stehen sollten? zb wie bei uns grad auf testing.'
  Beispiel heute: 6TCTFV steht seit 22.10.2026 in pre-process - zurueckgeschickt, danach hat keine Session die Lane gefahren.
  jaira hat keinen Daemon: eine Lane laeuft nur, wenn eine Session sie faehrt.
  Moegliche Orte: das Cockpit-Band (zaehlt schon wartende Tickets), 'jaira next --per-lane', das Board (Marker an der Karte).
  Verwandt: 28PD7N (Mensch aktiv benachrichtigen), C8V6DQ (Band zeigt lebende Claims).
definition-of-done: "Entscheidung: wo es gemeldet wird und ab wann ein Ticket als liegengeblieben gilt"
tags:
  - concurrency
blocked-by: []
related: []
commits: []
created-at: 2026-10-30T21:39:09Z
updated-at: 2026-10-30T21:39:09Z
---

# Tickets, die in einer Agenten-Lane ohne lebenden Claim stehen bleiben, werden gemeldet

## Definition of Done

- [ ] Entscheidung: wo es gemeldet wird und ab wann ein Ticket als liegengeblieben gilt

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

