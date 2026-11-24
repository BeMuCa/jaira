---
id: 01M86HC8QKQYKNWGHCQ7CSQXKR
title: "Cockpit 'Wartet auf dich': Frage schliesst sich selbst oder per x, zeigt ihre Zeit und uebersetzt nur auf Knopfdruck"
status: critique
ready: true
creator: Sascha
assignee: Sascha
goal: "Im Pane 'Wartet auf dich' stehen nur Fragen, die wirklich noch offen sind; jede zeigt, wann sie gestellt wurde; sie erscheint in der Originalsprache und wird nur auf Knopfdruck uebersetzt."
context: |-
  Was falsch ist (Pane aus 1M3X77, core/mod/builtin/jaira-cockpit/hooks/threads.ts + register.tsx, Texte texts.ts:267-274):
  1) Eine Frage verschwindet nur, wenn die Person im Pane eine Option waehlt oder der Klassifizierer eine Chat-Antwort erkennt. Wer im Cockpit per Accept/Give back entschieden hat, behaelt die Frage fuer immer in der Liste. Es gibt keinen Weg, sie zu schliessen.
  2) Es fehlt die Zeit: die Person liest eine Frage und weiss nicht, zu welcher Nachricht sie gehoert.
  3) Das Pane uebersetzt immer (settings language). Sascha (Alex) 24.11.2026: 'Knopf zum Uebersetzen, nicht immer; ich brauche das nicht immer'. Vorbild: FMWBA8 (Abnahme-Pane, Taste t / Translate, ein Aufruf fuer alles, nochmal druecken = Original).
  Ausloeser: Sascha 24.11.2026 im Teamlead-Chat.
definition-of-done: "Eine Frage zu einem Ticket verschwindet von selbst, wenn das Ticket die Lane verlassen hat, in der sie gestellt wurde (z.B. signoff -> done oder zurueck); Test"
tags:
  - cockpit
blocked-by: []
related: []
follows: 01M8537VMKP5N8TJ2YC91M3X77
commits: []
created-at: 2026-11-24T10:24:44Z
updated-at: 2026-11-24T11:28:47Z
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-67597
claimed-at: 2026-11-24T11:23:53Z
outcome-what: "Give-back edges closed: settled() ignores threads dropped by the 30-thread cap; file() keeps askedAt when the same question is repeated; translateAsks re-runs after a successful call when texts arrived meanwhile"
outcome-why: "Signoff give-back: three edges of this change would show wrong state in the Waiting-for-you pane"
outcome-resolves: "DoD 7-9; checked: concurrency (translate re-run reads state after the call), error paths (failed call stops, no loop), input (30-thread cap, repeated question)"
review-summary: "Delta review from 8a7401450 to f1a48e9d4: change() now reads the clock, then computes the new threads from the current ones in the same step, so pane answers, x, a filed turn and the lane check can no longer overwrite each other; the lane check settles a question for a vanished ticket only on 'jaira show' exit 5 and asks again on any other failure. The rest of the change (x to settle, settling on lane change, asked/answered times, Answered section kept 24h, Translate on t) stands as read in the full review."
round: 4
back-from: signoff
back-by: person
back-count: 1
review-gaps: |-
  note: edge: threads.ts file() slice(-THREADS) drops the oldest of 30 threads with its open question, and settled() labels it Answered; exclude threads absent from after by id in settled()
  note: edge: a question filed while a Translate call is out stays untranslated under an Original button until the next thread change; rerun translateAsks after busy clears if texts are missing
  note: edge: file() resets askedAt to now when the assistant repeats the same question on the same thread; keep askedAt when the question text is unchanged
  note: limit: hotkey x only on the first open question; the others have the click target, as DoD 2 allows
test-verdict: "pass: round 3, only the delta f1a48e9d4: cockpit suite 560/560 green (exit 0), incl. new 'a ticket in no list whose jaira show exits 1|5'; go test ./core/role green after shipped.sha256 regenerated; lost update closed by construction (next computed from threads after the clock await, no await between read and write)"
hunt-verdict: "pass: stop-after-fix, no hunt in round 3"
review-verdict: "accept: both fix-now edges closed (change() takes a function, only exit 5 settles); DoD 1-6 met; claude plugin test 560 pass, go test ./core/shipped ok"
review-round: 3
review-verdict-gaps: |-
  note: edge: threads.ts file() slice(-THREADS) drops the oldest of 30 threads with its open question, and settled() labels it Answered; exclude threads absent from after by id in settled()
  note: edge: a question filed while a Translate call is out stays untranslated under an Original button until the next thread change; rerun translateAsks after busy clears if texts are missing
  note: edge: file() resets askedAt to now when the assistant repeats the same question on the same thread; keep askedAt when the question text is unchanged
  note: limit: hotkey x only on the first open question; the others have the click target, as DoD 2 allows
review-check: |-
  1. cd core/mod/builtin/jaira-cockpit && claude plugin test .  - expect 0 fail
  2. claude --plugin-dir core/mod/builtin/jaira-cockpit, ask the assistant something that ends in a question to you about a ticket
  3. Open Waiting for you from the band: question shows asked today HH:MM, in its original language
  4. Press t: translated; t again: original
  5. Press x: question moves below the rule under Answered, dimmed, answered today HH:MM; nothing sent into the chat
  6. Ask another question about a ticket, run jaira move <id> --to <other lane> in a shell: within 60 s it moves to Answered
stop-after-fix: ""
---

# Cockpit 'Wartet auf dich': Frage schliesst sich selbst oder per x, zeigt ihre Zeit und uebersetzt nur auf Knopfdruck

## Definition of Done

- [x] Eine Frage zu einem Ticket gilt von selbst als erledigt, wenn das Ticket die Lane verlassen hat, in der sie gestellt wurde (z.B. signoff -> done oder zurueck); Test
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts 'settling questions, their times and Translate (CSQXKR) > a question on a ticket goes to Answered once the ticket leaves the lane it was asked in' + 'lane after: *' + 'a ticket in no list whose jaira show exits 1|5' (threads.ts moved, register.tsx settleMoved)
- [x] Jede offene Zeile hat x bzw. ein Klick-Ziel 'erledigt', das sie ohne Antwort in den Chat als erledigt markiert; Test
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts 'settling questions, their times and Translate (CSQXKR) > x settles a question without a word into the session; ...' (register.tsx dismissAsk, key ask-N-done)
- [x] Jede Frage zeigt Datum und Uhrzeit, zu der sie gestellt wurde; eine erledigte zusaetzlich, wann sie erledigt wurde (z.B. 'gefragt 07:12 · beantwortet 07:40'); Test
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts 'settling questions, their times and Translate (CSQXKR) > times: today by the clock, another day with its date' + x test matches 'asked today HH:MM' + lane test 'asked today .. · answered today ..'
- [x] Das Pane zeigt die Fragen im Original; Taste t bzw. ein Knopf Translate uebersetzt alle auf einmal, nochmal druecken zeigt wieder das Original; Test
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts 'settling questions, their times and Translate (CSQXKR) > the pane opens on the original; t translates every question at once, t again shows the original' + 'without a language set there is no Translate'
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/CSQXKR.md
- [x] Erledigte Fragen (Antwort im Pane, im Chat, per x oder automatisch) verschwinden nicht sofort: sie stehen unter einer Trennlinie in einem eigenen Abschnitt 'Beantwortet' unter den offenen, gedimmt, und fallen nach 24 Stunden heraus; offene stehen oben; Test
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts 'settling questions, their times and Translate (CSQXKR) > what a change settled carries both times; a day later it is gone' + 'settled questions are kept across a restart, and one settled a day ago is no longer shown' + 1M3X77 'the band opens the pane...' (Answered after a pick)
- [x] Ein Thread, der beim Kuerzen auf 30 wegfaellt, setzt seine offene Frage nicht auf Answered
  proof: core/mod/builtin/jaira-cockpit/hooks/threads.ts:259; test 'a thread dropped to keep thirty is not answered; its open question settles nothing'
- [x] Eine Frage, die waehrend eines laufenden Translate-Aufrufs ankommt, wird danach uebersetzt
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:2118; test 'a question filed while Translate is out is translated once that call is back'
- [x] Dieselbe Frage erneut auf demselben Thread behaelt askedAt
  proof: core/mod/builtin/jaira-cockpit/hooks/threads.ts:150; test 'the same question asked again keeps when it was first asked; another question is asked now'

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-24 10:25 · Sascha** — parent session: uds:/run/user/1000/cc-socks/3417.sock
- **2026-11-24 10:26 · Sascha** — decision: Sascha 24.11.2026: beide Wege (automatisch nach Lane-Wechsel UND manuell per x). Nicht sofort entfernen: erledigte Fragen einen Tag lang unter einer Trennlinie als 'Beantwortet' zeigen, damit klar ist, was neu und was schon beantwortet ist. Zeit der Frage und Zeit der Erledigung anzeigen.
- **2026-11-24 10:26 · Sascha** — decisions (worker, no open question): 1) lane at asking time = status from the same 'jaira show' that looks up the title, only for the ticket the question names (about); the thread's own ticket is not looked up (as today). Every board refresh (60s) compares with 'jaira list' of all boards; a ticket in none of them is asked once with 'jaira show'; unknown or other lane -> done. 2) closed = any open question that is gone after a change (pick, chat filing, x, lane left, a new question on the same thread replacing it) - goes to 'Answered' with answeredAt, dropped after 24h. 3) times in English, cockpit is English only (NZZDRH): 'asked today 07:12 · answered today 07:40', other day '23.11. 21:40'. 4) question is now kept in the assistant's own wording (classifier prompt), Translate = one haiku call for every question+option shown, into settings language; no language set = no button (as FMWBA8); a pick sends the original option text.
- **2026-11-24 10:34 · Sascha** — waiting: cockpit test run (background b0ym9m1cy) since 11:34; code for DoD 1-4,6 written, release note written, shipped.sha256 still to regenerate (go test ./core/role fails until then)
- **2026-11-24 10:40 · Sascha** — in-progress: the pane no longer closes after the last pick (it would hide 'Answered'; the band button only appears while a question is open, so a closed pane could not show it). shipped.sha256 must be regenerated (scripts/shipped-history.sh --write) after every change to the cockpit files, or go test ./core/role fails. Known limit: a question whose ticket is only the thread's own (no 'about') does not settle by lane - its lane is never read (as its title is not).
- **2026-11-24 10:41 · Sascha** — critique: settled() drops options (never shown under Answered, only cost Translate tokens); ASKS pane reads e.props.bodyColumns without the cast, as the Overview pane does. Rest stands: separate asksShown atom (questions are no ticket, the rephrased/translating atoms are keyed by ticket), settling by replacement (decision note above).
- **2026-11-24 10:44 · Sascha** — critique pass 2: both findings of pass 1 answered in 204c40977; nothing else read, nothing left.
- **2026-11-24 10:48 · Sascha** — testing: note: edge: lost update; register.tsx change() takes 'next' computed before its own await $.clock.now(), so a second change landing in that await is overwritten (fileTurn passes file(threads,...)); not seen in the two-writer test; fix: change($, ts => ...) computing next after the await
- **2026-11-24 10:48 · Sascha** — waiting: bughunt agent (background subagent) since 11:48
- **2026-11-24 10:51 · Sascha** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/threads.ts:255 (settled) with threads.ts:165 (file: next.slice(-THREADS)) - class: every way a question leaves the thread list without being settled is shown as settled; settled() treats any open question missing after a change as answered and the pane lists it under 'Answered' with 'answered <time>'; example: with 30 threads the first one (array order, not activity) carries an open question, the next turn files a new thread, the slice drops that thread and the unanswered question appears dimmed under Answered. That the question is dropped is older (slice from 65446e9/86da7c4); labelling it answered is this change. Not a blocker: only after 30 threads, and the question was lost before too; ran: scratch test in a copy of the plugin, file() 30 threads + 1, settled() returned ['still open?'] (claude plugin test, 1 fail as predicted)
- **2026-11-24 10:51 · Sascha** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:2069 (settleMoved) - class: any non-zero exit of 'jaira show' is read as 'jaira no longer knows the ticket' and settles the question; only exit 5 means not found (jaira show ZZZZZZ --json -> exit 5); example: a ticket on a board whose 'jaira list --board' failed this round (skipped with continue at register.tsx:394) and a 'jaira show' that fails for another reason (git or lock error, exit 1) moves an open question to Answered for good. A timeout (undefined) is handled. Not a blocker: needs the ticket missing from every list read plus a failing show; read
- **2026-11-24 10:51 · Sascha** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:2052 (showThreads) with :2097 (translateAsks) - class: a question or settled question that arrives while a Translate call is out is never translated until the next thread change; showThreads skips the call while shown.busy, and translateAsks writes only the texts it started with; example: pane on 'Original' (translation shown), a filing lands during the haiku call, the new question stays in its original language under a button that says the pane is translated, until the next turn or tick changes the threads; read
- **2026-11-24 10:51 · Sascha** — hunt report: 3 findings, all 3 change, none a blocker, notes on CSQXKR | read whole: threads.ts, texts.ts (changed part and t table), types/index.d.ts, release note; register.tsx and board.ts read as the ticket diff plus every function it calls or is called from (refreshBoard, fileTurn, keepThreads, change, settleMoved, showThreads, openAsks, translateAsks, rephrased, answerAsk/dismissAsk, the asks pane, session.start, turn.complete, the 60s tick, lanesOf, standOf); tests/board.test.ts read for the CSQXKR cases and the world mock | not confirmed: change() and fileTurn now await $.clock.now() between reading the global threads and writing it (register.tsx:2003, :2036), so a 60s refreshBoard tick, an x/pick or a filing that runs in that window could overwrite the other's result (a lost or revived question); two scratch tests that interleave them in the mock passed, so the engine may serialise them, which only a live session can show
- **2026-11-24 10:51 · Sascha** — bughunt: done at 8a7401450 - 3 findings (3 change, 0 pre-existing) in 6 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/hooks/threads.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/mod/builtin/jaira-cockpit/types/index.d.ts core/release/notes/CSQXKR.md core/shipped/shipped.sha256
- **2026-11-24 10:52 · Sascha** — waiting: fresh review agent (run-lane.sh --no-worktree, Herdr tab) since 11:52
- **2026-11-24 10:54 · Sascha** — review: done at 8a7401450
- **2026-11-24 10:54 · Sascha** — review: fix-now: (1) register.tsx change(): compute next after await $.clock.now() via a function argument, so concurrent settleMoved/dismissAsk/answerAsk/fileTurn do not overwrite each other; (2) settleMoved: only jaira show exit 5 settles a question, any other non-zero exit skips the round
- **2026-11-24 10:54 · Sascha** — fix round: closing the review's two fix-now (change() lost update; only jaira show exit 5 means ticket gone). Set stop-after-fix=true: under 50 lines, no lock, no parser. Other review notes left for the person.
- **2026-11-24 10:58 · Sascha** — review: fix round on f1a48e9d4; waiting: fresh review agent (run-lane.sh, Herdr tab) since 11:58
- **2026-11-24 11:01 · Sascha** — review: done at f1a48e9d4
- **2026-11-24 11:02 · Sascha** — report: Waiting-for-you pane: x and lane-left settle questions into a dimmed 'Answered' section (24h), asked/answered times, original text with t=Translate; 4 commits 906415bdf..f1a48e9d4 | signoff, context unknown | open: person tries it in the cockpit (review-check); 3 edge notes + 1 limit from review for the person; branch not pushed
- **2026-11-24 11:23 · Sascha** — Given back from signoff: beheben:
- note: edge: threads.ts file() slice(-THREADS) drops the oldest of 30 threads with its open question, and settled() labels it Answered; exclude threads absent from after by id in settled()
- note: edge: a question filed while a Translate call is out stays untranslated under an Original button until the next thread change; rerun translateAsks after busy clears if texts are missing
- note: edge: file() resets askedAt to now when the assistant repeats the same question on the same thread; keep askedAt when the question text is unchanged
angenommen:
- note: limit: hotkey x only on the first open question; the others have the click target, as DoD 2 allows
- **2026-11-24 11:23 · Sascha** — cockpit: given back to a new worker - nothing recorded about the old worker session (cache and context unknown); branch feat/csqxkr is the ticket's worktree's, worktree /home/alex/projects/.worktrees/jaira-csqxkr; reports to uds:/run/user/1000/cc-socks/3417.sock, the ticket's parent session
- **2026-11-24 11:24 · Sascha** — parent session: uds:/run/user/1000/cc-socks/3417.sock; new worker on the give-back: fixing the 3 edges (thread drop, translate while busy, askedAt reset); no open decisions, the give-back names each fix
- **2026-11-24 11:27 · Sascha** — waiting: claude plugin test run (background) since 12:27
- **2026-11-24 11:28 · Sascha** — give-back fixes: test mock model.complete now returns translations only for 'Translate each text' prompts, so a filing during a held Translate call gets w.model; translateAsks re-runs itself after a successful call when texts are missing (not on failure: on=false then, no loop)
