---
id: 01M8AF6BDDF8AGAY5T7EGN8NQN
title: "Ein Worker ohne Herdr schickt seine Fragen per SendMessage an main und arbeitet weiter, statt mit der Frage im Bericht stehen zu bleiben"
status: todo
ready: true
creator: Berk
assignee: Berk
goal: "Laeuft ein jaira-Worker oder Dispatcher als Hintergrund-Subagent (ohne Herdr, ohne --parent), erreicht seine Frage die Main-Session sofort, und er arbeitet an allem weiter, was nicht von der Antwort abhaengt."
context: |-
  Was heute falsch ist: Die Skills jaira-worker und jaira-dispatcher sagen, ein Subagent kann die Person nicht fragen; er schreibt die Optionen oben in seinen Abschlussbericht und haelt an. Die Frage kommt erst an, wenn er fertig ist, und die Arbeit steht so lange.
  Mit Herdr geht es schon: spawn.sh gibt den Socket der Elternsitzung als --parent mit, der Worker schickt per SendMessage.
  Am 23.11.2026 getestet: ein Hintergrund-Subagent (Agent-Tool, run_in_background) kann SendMessage mit to="main" aufrufen; die Nachricht kam in der Main-Session an. Seitdem haben zwei Worker so gefragt (EH6Z0N, KQQD1H).
  Laut SendMessage-Doku gilt "main" nur fuer Hintergrund-Subagenten; ein Vordergrund-Subagent kann es nicht - dann bleibt der alte Weg (Optionen im Bericht).
  Wunsch von berkc am 26.11.2026: Ticket anlegen und direkt umsetzen.
  Die Skills werden im Binary ausgeliefert (core/role, Hash in core/shipped/shipped.sha256).
definition-of-done: "jaira-worker sagt: ohne --parent und als Hintergrund-Subagent geht eine Frage per SendMessage to=\"main\" raus (Optionen, Empfehlung zuerst), wird mit jaira note festgehalten, und der Worker arbeitet an Unabhaengigem weiter"
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-11-25T23:03:30Z
updated-at: 2026-11-25T23:05:36Z
updated-by: Berk
claimed-by: w-gn8nqn
claimed-at: 2026-11-25T23:04:48Z
---

# Ein Worker ohne Herdr schickt seine Fragen per SendMessage an main und arbeitet weiter, statt mit der Frage im Bericht stehen zu bleiben

## Definition of Done

- [ ] jaira-worker sagt: ohne --parent und als Hintergrund-Subagent geht eine Frage per SendMessage to="main" raus (Optionen, Empfehlung zuerst), wird mit jaira note festgehalten, und der Worker arbeitet an Unabhaengigem weiter
- [ ] jaira-dispatcher und die Rollen, die heute 'Optionen in den Bericht' sagen, sagen dasselbe; Vordergrund-Subagent behaelt den alten Weg
- [ ] jaira-teamlead sagt, dass so eine Nachricht von einem Subagenten wie eine Dispatcher-Frage behandelt wird (AskUserQuestion an die Person, Antwort per SendMessage zurueck)
- [ ] shipped.sha256 und core/shipped-Tests stimmen; Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

