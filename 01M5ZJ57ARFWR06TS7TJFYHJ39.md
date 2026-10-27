---
id: 01M5ZJ57ARFWR06TS7TJFYHJ39
title: "Ein Lane-Schreibbefehl aus einem Worktree sagt einmal, wenn der Ref ihn nicht genommen hat, und was zu tun ist"
status: backlog
ready: false
creator: BeMuCa
goal: "Nimmt der Ref einen Schreibbefehl nicht (jemand anders schrieb zuerst), sagt jaira es einmal, klar und mit dem naechsten Schritt (fetch und wiederholen); stimmt der Ref laengst mit der Datei ueberein, kommt keine Meldung."
context: |-
  Was schief ist: 'was not sent: BeMuCa wrote it first' kam bei fast jedem Lane-Schreibbefehl aus einem Ticket-Worktree, obwohl niemand anders schrieb; der Ref hinkte der Datei hinterher und holte beim naechsten Schreiben auf (O1). Bei 1K3RA8 kam die Meldung, obwohl der Ref schon passte. 4N1B66: jaira pull brach mit 'taken while you were pulling it' ab, jaira fetch plus Wiederholung half (O9).
  Agenten lernten, die Meldung zu ignorieren - und damit auch echte Konflikte.
  Bekannt: alle Worktrees teilen ein .git; parallele Driver schreiben denselben Ref (Notiz vom 15.10. zu TS6YCJ-aehnlichen Faellen).
  Zu finden: warum der Vergleich 'jemand anders schrieb zuerst' bei gleichem Inhalt anschlaegt (core/sync oder core/ref: grep 'wrote it first').
definition-of-done: "Ein Schreibbefehl, dessen Inhalt schon auf dem Ref liegt, meldet nichts"
tags:
  - concurrency
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-27T20:52:17Z
updated-at: 2026-10-27T20:52:17Z
---

# Ein Lane-Schreibbefehl aus einem Worktree sagt einmal, wenn der Ref ihn nicht genommen hat, und was zu tun ist

## Definition of Done

- [ ] Ein Schreibbefehl, dessen Inhalt schon auf dem Ref liegt, meldet nichts
- [ ] Ein echter Konflikt wird einmal gemeldet, nennt den Schreiber und sagt 'jaira fetch, dann wiederholen'
- [ ] Ein Test mit zwei Worktrees auf einem .git belegt beide Faelle
- [ ] Release-Notiz in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

