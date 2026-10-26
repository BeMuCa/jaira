---
id: 01M5MNWS6C8YAE88PF43VPM51D
title: "Ein Checkout, dessen Ticket anderswo abgelegt wurde, zeigt es als abgelegt statt als geaendert"
status: backlog
ready: false
creator: BeMuCa
goal: "Nach jaira logbook in einem Worktree verwirrt die alte Ticketdatei in den anderen Worktrees und im Haupt-Checkout niemanden mehr: sie gilt als abgelegt, und jaira bietet an, sie zu entfernen."
context: |-
  Beobachtet 22.-23.10.2026 (Feedback Punkt 9): nach einem --force-Move und jaira logbook in einem Worktree zeigten die anderen Worktrees die Ticketdatei weiter als geaendert; der naechste Worker hielt das Ticket fuer offen.
  K6BGB0 hat den Fall 'Datei hinter Ref' geloest; dieser hier ist 'Datei ohne Ref, Ticket upstream im Logbuch' - K6BGB0s Notiz vom 14.10. nennt ihn als dritte Variante, die offen bleibt.
  Verwandt: DPSFGE (git status im Haupt-Checkout ohne Ticket-Rauschen), AYBDJ2 (Arbeitskopie nach logbook aufraeumen).
definition-of-done: "Board, jaira list und jaira show erkennen eine Ticketdatei, deren Ticket im Logbuch eines anderen Checkouts oder auf dem Remote abgelegt ist, und zeigen sie als 'abgelegt'; ein Test baut den Fall mit zwei Worktrees"
tags:
  - concurrency
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T15:25:56Z
updated-at: 2026-10-26T19:08:07Z
assignee: BeMuCa
updated-by: BeMuCa
---

# Ein Checkout, dessen Ticket anderswo abgelegt wurde, zeigt es als abgelegt statt als geaendert

## Definition of Done

- [ ] Board, jaira list und jaira show erkennen eine Ticketdatei, deren Ticket im Logbuch eines anderen Checkouts oder auf dem Remote abgelegt ist, und zeigen sie als 'abgelegt'; ein Test baut den Fall mit zwei Worktrees
- [ ] jaira bietet an, die veraltete Kopie zu entfernen (Befehl in der Meldung); nichts wird still geloescht
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

