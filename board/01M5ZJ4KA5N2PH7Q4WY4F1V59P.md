---
id: 01M5ZJ4KA5N2PH7Q4WY4F1V59P
title: "Ein Review-Verdikt 'zurueckschicken' laesst das Ticket nicht nach signoff"
status: backlog
ready: false
creator: BeMuCa
goal: "Steht im review-verdict der Review-Lane 'send back' (oder das Gegenstueck des Boards), nimmt jaira move aus review nur noch in-progress als Ziel; signoff wird mit Exit 3 abgelehnt und nennt das Verdikt."
context: |-
  Was schief ist: die Review-Lane schrieb als Verdikt 'zurueckschicken', das Ticket landete trotzdem in signoff (Ticket 25DZW3 auf einem anderen Board, Teamlead-Session 25.-27.10.2026).
  Heute prueft kein Gate, ob Verdikt und Ziellane zusammenpassen: review muss review-verdict liefern, aber move liest den Inhalt nicht.
  Der Mensch sah ein Ticket zur Abnahme, das der Reviewer selbst nicht freigegeben hatte.
  Bekannt: die Gates leben in core/gate/gate.go; welche Werte review-verdict annimmt, steht in der Lane-Datei review der Boards (grep review-verdict core/lane/).
  Nicht gewollt: ein Freitext-Parser. Ein fester kleiner Wertevorrat (z. B. accept, send-back) reicht; alles andere ist kein Verdikt.
  Zusammenhang: L312/L327 der Feedback-Liste vom 27.10.2026.
definition-of-done: "Ein Move aus review nach signoff mit review-verdict 'send back' wird mit Exit 3 abgelehnt; die Meldung nennt das Verdikt und die erlaubte Ziellane"
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-10-27T20:51:56Z
updated-at: 2026-10-27T20:51:56Z
---

# Ein Review-Verdikt 'zurueckschicken' laesst das Ticket nicht nach signoff

## Definition of Done

- [ ] Ein Move aus review nach signoff mit review-verdict 'send back' wird mit Exit 3 abgelehnt; die Meldung nennt das Verdikt und die erlaubte Ziellane
- [ ] Ein Move aus review nach in-progress mit 'send back' geht ohne --force
- [ ] Die Werte, die review-verdict annimmt, stehen in der Lane-Datei und in jaira show --for-lane review
- [ ] Release-Notiz in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

