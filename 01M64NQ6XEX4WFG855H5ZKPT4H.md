---
id: 01M64NQ6XEX4WFG855H5ZKPT4H
title: "Web zeigt den Stand vom Ref, Erstelldatum und wie lange ein Ticket unberührt ist"
status: todo
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Das Board im Browser zeigt für jedes Ticket die Lane, die auf seinem Ref steht, und dazu, wann es angelegt wurde und wie lange es schon niemand angefasst hat."
context: |-
  Gefunden am 29.10.2026 auf einem Board mit geteilten Refs:
  - jaira web zeigt ein Ticket als in-progress, auf seinem Ref steht es seit Tagen in signoff.
  - Ursache: web liest die Ticketdatei im Arbeitsverzeichnis (master). Die Datei kam mit dem letzten Code-Commit dorthin; spätere Lanes ohne Code (critique, review, signoff) schreiben nur den Ref. Die Datei bleibt stehen.
  - Betroffen waren dort über 20 Tickets: Datei sagt critique/testing/in-progress, Ref sagt signoff.
  - Gewünscht (Alex): web nimmt den Ref und führt ihn mit der Datei zusammen – der neuere Stand gewinnt (updated-at), wie 'jaira pull' es schon beim Vergleich tut ('the ticket file here is older than its ref').
  - Außerdem fehlen auf der Karte: Erstelldatum (created) und 'unberührt seit' (Zeit seit updated-at, z.B. '3 Tage').
  - Regel für diese Arbeit: kein Name, Pfad und keine Ticket-ID eines anderen Projekts in Code, Tests, Commits, Notes oder Release-Notes. Testdaten sind erfunden.
definition-of-done: "Ein Ticket, dessen Datei älter ist als sein Ref, erscheint im Web in der Lane des Refs (Test mit Datei=in-progress, Ref=signoff)"
tags:
  - web
blocked-by: []
related: []
commits: []
created-at: 2026-10-29T20:30:44Z
updated-at: 2026-10-29T20:37:17Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1236953
claimed-at: 2026-10-29T20:37:14Z
---

# Web zeigt den Stand vom Ref, Erstelldatum und wie lange ein Ticket unberührt ist

## Definition of Done

- [ ] Ein Ticket, dessen Datei älter ist als sein Ref, erscheint im Web in der Lane des Refs (Test mit Datei=in-progress, Ref=signoff)
- [ ] Ist die Datei neuer als der Ref, gewinnt die Datei
- [ ] Jede Karte zeigt das Erstelldatum
- [ ] Jede Karte zeigt, wie lange das Ticket unberührt ist (seit updated-at), in Tagen/Stunden

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-29 20:35 · Alexander Sacharov** — decision (Alex, 29.10.): keine Erwähnung anderer Projekte – weder in Code, Tests, Fixtures, Commits, Notes noch Release-Notes. Kontext dafür umgeschrieben.
- **2026-10-29 20:37 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/1226588.sock
