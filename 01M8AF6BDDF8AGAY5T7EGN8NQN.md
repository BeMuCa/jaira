---
id: 01M8AF6BDDF8AGAY5T7EGN8NQN
title: "Ein Worker ohne Herdr schickt seine Fragen per SendMessage an main und arbeitet weiter, statt mit der Frage im Bericht stehen zu bleiben"
status: critique
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
commits:
  - 93e7b219801136cd21e8edde9fafd9a3f381da6f
created-at: 2026-11-25T23:03:30Z
updated-at: 2026-11-25T23:32:44Z
updated-by: Berk
claimed-by: w-gn8nqn
claimed-at: 2026-11-25T23:32:44Z
stop-after-fix: "true"
outcome-what: "Worker, dispatcher and teamlead skills: a background subagent with no parent sends its question to main via SendMessage to=\"main\", notes it, sets the question field and keeps working; refused or no SendMessage keeps the report path. Teamlead handles such a message like a dispatcher question. Shipped hashes appended (3 lines); release note core/release/notes/GN8NQN.md."
outcome-why: "A background subagent's question sat in its final report until it ended, so the work stood; SendMessage to main reaches the main session at once (tested 23.11.2026)."
outcome-resolves: "Goal: a background worker or dispatcher question reaches main immediately and independent work continues; the old report path stays only for foreground or headless agents."
---

# Ein Worker ohne Herdr schickt seine Fragen per SendMessage an main und arbeitet weiter, statt mit der Frage im Bericht stehen zu bleiben

## Definition of Done

- [x] jaira-worker sagt: ohne --parent und als Hintergrund-Subagent geht eine Frage per SendMessage to="main" raus (Optionen, Empfehlung zuerst), wird mit jaira note festgehalten, und der Worker arbeitet an Unabhaengigem weiter
  proof: core/role/builtin/jaira-worker/SKILL.md:363-381 (new no-parent case, SendMessage to main, note + question field, keep working; answer-from-main paragraph); :54 when-you-wait bullet; go test ./core/role/... passes
- [x] jaira-dispatcher und die Rollen, die heute 'Optionen in den Bericht' sagen, sagen dasselbe; Vordergrund-Subagent behaelt den alten Weg
  proof: core/role/builtin/jaira-dispatcher/SKILL.md:199-208 (same case, old report path only for refused SendMessage or no tool); worker's old-way narrowed at jaira-worker/SKILL.md:370; role-lane left unchanged by plan decision
- [x] jaira-teamlead sagt, dass so eine Nachricht von einem Subagenten wie eine Dispatcher-Frage behandelt wird (AskUserQuestion an die Person, Antwort per SendMessage zurueck)
  proof: core/role/builtin/jaira-teamlead/SKILL.md:90-96 (background subagent message = dispatcher question: AskUserQuestion unchanged, jaira note, SendMessage back to agent)
- [x] shipped.sha256 und core/shipped-Tests stimmen; Release-Note-Zeile in core/release/notes/<ID>.md
  proof: core/shipped/shipped.sha256:806, :1221, :1268 (new hashes appended in path order); core/release/notes/GN8NQN.md:1; go test ./core/shipped ./core/role/... ok

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
- **2026-11-25 23:15 · Berk** — Dead end: go test ./core/shipped passes on its own; the role hashes are checked by TestEveryInstalledFileIsAShippedVersion in core/role (the shipped test only checks one fixed path). Did not run scripts/shipped-history.sh --write (regenerates the whole manifest from git history); appended the 3 new hashes by hand in path order, 0 lines removed. Role-lane left unchanged on purpose (its questions go to the dispatcher via the question field and exit 5).
- **2026-11-25 23:22 · Berk** — waiting: critique (opus subagent, background) since 00:22
