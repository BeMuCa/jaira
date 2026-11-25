---
id: 01M8AF6BDDF8AGAY5T7EGN8NQN
title: "Ein Worker ohne Herdr schickt seine Fragen per SendMessage an main und arbeitet weiter, statt mit der Frage im Bericht stehen zu bleiben"
status: in-progress
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
updated-at: 2026-11-25T23:12:43Z
updated-by: Berk
claimed-by: w-gn8nqn
claimed-at: 2026-11-25T23:07:35Z
stop-after-fix: "true"
---

# Ein Worker ohne Herdr schickt seine Fragen per SendMessage an main und arbeitet weiter, statt mit der Frage im Bericht stehen zu bleiben

## Definition of Done

- [x] jaira-worker sagt: ohne --parent und als Hintergrund-Subagent geht eine Frage per SendMessage to="main" raus (Optionen, Empfehlung zuerst), wird mit jaira note festgehalten, und der Worker arbeitet an Unabhaengigem weiter
  proof: core/role/builtin/jaira-worker/SKILL.md:363-381 (new no-parent case, SendMessage to main, note + question field, keep working; answer-from-main paragraph); :54 when-you-wait bullet; go test ./core/role/... passes
- [ ] jaira-dispatcher und die Rollen, die heute 'Optionen in den Bericht' sagen, sagen dasselbe; Vordergrund-Subagent behaelt den alten Weg
- [ ] jaira-teamlead sagt, dass so eine Nachricht von einem Subagenten wie eine Dispatcher-Frage behandelt wird (AskUserQuestion an die Person, Antwort per SendMessage zurueck)
- [ ] shipped.sha256 und core/shipped-Tests stimmen; Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] jaira-worker 'How you ask': new case between 'AskUserQuestion' and 'none at all' — no --parent, no AskUserQuestion, SendMessage to="main" accepted (background subagent): note the question with numbered options (recommendation first), set the question field, SendMessage to main, keep working on what does not depend on the answer; nothing left -> waiting note and end turn with the options on top of the status line; the answer arrives as a message, note it, clear the field
- [x] jaira-worker: the old 'options into your report' case narrowed to 'SendMessage to main refused or not there (foreground subagent, headless)'; 'When you wait' bullet gets the main case too
- [x] jaira-dispatcher 'Where the question goes': same new case and narrowed fallback; workers' questions it carries go to main the same way
- [x] jaira-teamlead 'When to stop': a message from a subagent you started in the background is handled like a dispatcher question — AskUserQuestion with its options unchanged, jaira note, answer back with SendMessage to that agent
- [x] append new hashes to core/shipped/shipped.sha256 as the core/shipped tests demand; go test ./core/shipped ./core/role -count=1
- [x] release note core/release/notes/GN8NQN.md, one line

## Progress
- **2026-11-25 23:06 · Berk** — plan decisions: (1) how an agent knows it is a background subagent: it does not need to know — it calls SendMessage to="main"; accepted = background, refused/no tool = foreground/headless -> old way. No detection rule to maintain. (2) jaira-role-lane NOT changed: its questions go to the dispatcher via the ticket question field and run-lane.sh exit 5 / its own report to the caller; sending to main would skip the dispatcher that carries it. (3) jaira-role-pr not changed: grep for 'options go into your report' hits only worker and dispatcher; pr's no-parent case reports in its own tab. (4) worker still sets the question field beside the note so a restarted worker and the cockpit see it open. (5) out of independent work: waiting note + end turn with the options on top of the status line, so nothing is lost if the answer never comes. Never pick an option itself.
- **2026-11-25 23:06 · Berk** — stop-after-fix=true: the change is skill text only, no code; a further hunt after a fix cannot find anything worth a round
- **2026-11-25 23:07 · Berk** — waiting: in-progress (haiku subagent, background) since 00:07
