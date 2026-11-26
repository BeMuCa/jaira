---
id: 01M8AF6BDDF8AGAY5T7EGN8NQN
title: "Ein Worker ohne Herdr schickt seine Fragen per SendMessage an main und arbeitet weiter, statt mit der Frage im Bericht stehen zu bleiben"
status: review
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
updated-at: 2026-11-26T09:02:32Z
updated-by: Berk
claimed-by: w-gn8nqn
claimed-at: 2026-11-26T09:02:29Z
stop-after-fix: "true"
outcome-what: "one ask-main order in worker and dispatcher (note and question field first, then SendMessage to main); dispatcher clears the question field after the answer"
outcome-why: "critique round 1 found two gaps: stale question field on dispatcher, and a send-note-send-again reading of the ask-main step"
outcome-resolves: "critique round 1: dispatcher had no clear-the-field sentence; worker/dispatcher ask-main step read as send, note, send again"
review-summary: "Full review (no earlier review note) of 93e7b219..bf15d840, skill text only. jaira-worker and jaira-dispatcher 'How you ask/Where the question goes': the old 'no AskUserQuestion -> options into report and stop' case becomes 'no AskUserQuestion and no parent': a background subagent notes the question with numbered options, sets the question field, then SendMessage to=main and keeps working on what does not depend on the answer; out of work it writes a waiting: note and ends its turn with the options above the status line. A refused SendMessage or no tool (foreground subagent, headless) keeps the old report-and-stop way. Both say the answer from main arrives as a message: note it, clear the question field, go on. Worker 'When you wait' also sends the waiting line to main. jaira-teamlead 'When to stop': a background subagent's SendMessage question is handled like a dispatcher's (AskUserQuestion with options unchanged, jaira note, SendMessage the answer back by name/id). shipped.sha256 gets 5 appended hashes (2 dispatcher, 2 worker, 1 teamlead; current files match 69bbb910/74f035fb/90844161), release note core/release/notes/GN8NQN.md is one '- ' line."
round: 2
back-from: critique
back-by: lane
back-count: 1
review-gaps: "none removed: worker/dispatcher duplicated ask-main paragraph left (existing pattern, round-1 critique let it stand); teamlead SKILL.md:90-94 slight overlap with the next sentence left because it names the return route; no dead code (skill text only)"
hunt-verdict: "pass: stop-after-fix, no hunt in round 2 (skill text, shipped hashes and release note only; no locking or parser)"
test-verdict: "pass: go test ./core/shipped ./core/role ./core/release -count=1 RC=0 (run 2; run 1 red RC=1 on TestSpawnOffsetLock under load avg ~21, passes alone RC=0, rerun green); DoD 1-4 checked in the tree: worker/dispatcher ask-main text, teamlead SendMessage clause, three shipped hashes match sha256sum of current files, release note one '- ' line; function not exercisable beyond go:embed + shipped test (LLM-read skill text); no code written, no commit"
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
  proof: core/shipped/shipped.sha256:777 and :1273 (round-2 hashes, dispatcher and worker SKILL.md); :807, :1222, :1269 (round-1 hashes, shifted by the insert); core/release/notes/GN8NQN.md:1

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
- **2026-11-25 23:32 · Berk** — critique (round 1): core/role/builtin/jaira-dispatcher/SKILL.md:199-208 now has the dispatcher set the question field for its own question to main, but nothing tells it to clear it when the answer arrives (worker got that at jaira-worker/SKILL.md:381-382; dispatcher :222 only clears its workers fields) — a stale field reads as open in the cockpit and close-tab.sh; add the same sentence after the dispatcher list: answer from main comes as a message, jaira note it, jaira set <id> question=, go on
core/role/builtin/jaira-worker/SKILL.md:364-368 and jaira-dispatcher/SKILL.md:200-203 say "tries SendMessage to=main first: if accepted ..." and then "write onto the ticket ..., send the same text to main" — reads as send, then note, then send again; reword to one order: jaira note + question field, then SendMessage to=main; accepted = it reached main, refused/no tool = old way. Left alone: setting the question field beside the note and ending the turn with the options on top of the status line (plan decisions 4 and 5 in the 23:06 note); role-lane/role-pr unchanged (decisions 2, 3); duplicated paragraph in worker+dispatcher matches the existing pattern there.
- **2026-11-25 23:35 · Berk** — waiting: in-progress round 2 (haiku subagent) since 00:35
- **2026-11-25 23:54 · Berk** — round 2: dispatcher line ~210 got the clear-the-field sentence; worker and dispatcher now one order (note + question field, then SendMessage to main). shipped.sha256: core/role test demanded dispatcher 69bbb910 and worker 74f035fb; inserted each in its sorted place in the path group, no old line removed. TestReapTearsDownWhatTheSetupBuilt failed once in the full core/role run (child outlived 1s teardown under load), passes alone; not related to skill text
- **2026-11-26 00:01 · Berk** — critique (round 2): both round-1 findings addressed in bf15d840 — dispatcher SKILL.md ~:210 now clears the question field when main answers; worker :364-373 and dispatcher :200-208 now read in one order (note + question field, then SendMessage to main). No defect introduced by the repair. Change is 47 lines outside .jaira, no untracked files, no locking/parser -> critique+optimize in one pass. Optimize: nothing removed. Duplicated worker/dispatcher paragraph left (round-1 critique let it stand, matches existing pattern). Teamlead :90-94 partly restates the following 'Each answer goes onto the ticket ... SendMessage' sentence but adds the return route (by name or id) — left. go test ./core/role ./core/shipped -run 'Shipped|Installed' ok.
- **2026-11-26 00:25 · Berk** — worker handed back early (harness forced hand-back) at 01:25: ticket in testing, testing subagent still running its background work; left to do: testing verdict, review lane (fresh opus agent), signoff, git push -u origin feat/GN8NQN-ask-main
