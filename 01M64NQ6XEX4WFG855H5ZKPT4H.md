---
id: 01M64NQ6XEX4WFG855H5ZKPT4H
title: "Web zeigt den Stand vom Ref, Erstelldatum und wie lange ein Ticket unberührt ist"
status: in-progress
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
updated-at: 2026-10-29T20:39:01Z
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

- [~] refsync: exportierte Methode Newer(t) – liefert das Ref-Ticket, wenn isStale(t, ref), sonst nil
- [ ] web_source: board() ersetzt veraltete Dateien durch den Ref-Stand, bevor link.Build den Index baut; Ticket(id) ebenso
- [ ] web.Ticket: created_at/updated_at; webSource füllt sie für Epic- und lose Karten
- [ ] UI: TicketCard zeigt created und untouched; dist neu bauen
- [ ] Tests: web_source mit Datei=in-progress/Ref=signoff und umgekehrt; Go-Test für die Felder; Release-Note

## Progress
- **2026-10-29 20:35 · Alexander Sacharov** — decision (Alex, 29.10.): keine Erwähnung anderer Projekte – weder in Code, Tests, Fixtures, Commits, Notes noch Release-Notes. Kontext dafür umgeschrieben.
- **2026-10-29 20:37 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/1226588.sock
- **2026-10-29 20:38 · Alexander Sacharov** — Entscheidungen (selbst getroffen, keine offen für Alex): (1) Regel 'neuer gewinnt' ganz, nicht der Feld-Merge aus refsync.Reconcile – der hält status nie rückwärts und würde DoD 2 (Datei neuer, frühere Lane) brechen. Vergleich über refsync.isStale, dieselbe Definition wie pull/CheckFresh. (2) Nur jaira web; TUI/list lesen weiter die Datei – gleicher Fehler möglich, nicht Teil dieses Tickets. (3) Karte englisch wie der Rest der Seite: 'created 14.10.26' (locale) und 'untouched 3d' / '5h' / '<1h'; Tage ab 24h, sonst Stunden, berechnet im Browser aus updated_at. (4) Auch die Detailansicht nimmt den neueren Stand. (5) Epic-Karten eines abgelegten Tickets: Daten per LoadAnywhere.
