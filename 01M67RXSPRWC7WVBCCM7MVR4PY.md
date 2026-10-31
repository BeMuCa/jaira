---
id: 01M67RXSPRWC7WVBCCM7MVR4PY
title: Die Einzeiler des Cockpits lesen Lane und Urteil aus einer festen Schlusszeile statt aus dem ersten passenden Wort
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Ein Worker-Bericht endet mit einer festen Zeile (z.B. 'jaira: <ID> <von> -> <nach> · verdict: <...>'), und digest.ts liest Lane und Urteil nur aus ihr; fehlt sie, zeigt der Einzeiler nur die ID statt zu raten."
context: |-
  Was falsch ist: core/mod/builtin/jaira-cockpit/hooks/digest.ts raet. laneOf nimmt den ersten Pfeil irgendwo im Bericht (Zeile 67), verdictOf das erste Wort wie 'reject' irgendwo (Zeile 77-84). Ein Bericht, der 'reject? reason' oder 'nichts zurueckgewiesen' erwaehnt, wird zu '-> signoff · rejected'.
  Ausloeser: Maintainer 31.10.2026: 'why is our cockpit wrong this often? is it because haiku?' - Nein: kein Modell, reine Regex. Haiku klassifiziert nur die Threads.
  Bekannt: X6HB6V (signoff) korrigiert die Lane. Das Urteil bleibt falsch. Die Schlusszeile muss in jaira-worker, jaira-dispatcher und jaira-role-lane stehen.
definition-of-done: "Worker-, Dispatcher- und Lane-Rolle schreiben die feste Schlusszeile"
tags:
  - tui
blocked-by: []
related:
  - 01M5MNWVKBC8Q193EK7YX6HB6V
commits: []
created-at: 2026-10-31T01:24:29Z
updated-at: 2026-10-31T01:24:29Z
---

# Die Einzeiler des Cockpits lesen Lane und Urteil aus einer festen Schlusszeile statt aus dem ersten passenden Wort

## Definition of Done

- [ ] Worker-, Dispatcher- und Lane-Rolle schreiben die feste Schlusszeile
- [ ] digest.ts liest Lane/Urteil nur aus dieser Zeile; ein Bericht mit 'reject' im Fliesstext ergibt kein 'rejected'; Test
- [ ] Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

