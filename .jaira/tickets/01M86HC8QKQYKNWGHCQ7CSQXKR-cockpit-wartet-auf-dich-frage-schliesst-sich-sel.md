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
updated-at: 2026-11-24T12:12:55Z
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-52645
claimed-at: 2026-11-24T11:37:05Z
outcome-what: "Give-back edges closed: settled() ignores threads dropped by the 30-thread cap; file() keeps askedAt when the same question is repeated; translateAsks re-runs after a successful call when texts arrived meanwhile"
outcome-why: "Signoff give-back: three edges of this change would show wrong state in the Waiting-for-you pane"
outcome-resolves: "DoD 7-9; checked: concurrency (translate re-run reads state after the call), error paths (failed call stops, no loop), input (30-thread cap, repeated question)"
review-summary: "none"
round: 7
back-from: testing
back-by: lane
back-count: 2
review-gaps: "none: the three hunt fixes add one condition each and no helper; nothing orphaned"
test-verdict: "pass: claude plugin test 569/569 (exit 0) at e65c788c5, go test ./core/role ok; DoD 1-10 proofs re-pointed, their tests passing; failure paths: jaira show exit 1 or exit 0 without a status keeps the question, exit 5 settles; stale list re-read via show; t x3 and a reopened pane -> one model call; writes: threads store via change() (earlier round), not re-run"
hunt-verdict: "fail: 1 change finding (translateAsks leaves busy set when rephrased throws)"
review-verdict: "accept: DoD 1-9 met with tests; the three give-back edges are closed; claude plugin test 563 pass, go test ./core/shipped ok; two small edges left for the person"
review-round: 4
review-verdict-gaps: |-
  note: edge: settleMoved compares lanes from the jaira list read at the start of refreshBoard; a question filed and titled (fileTurn) with a newer lane while that refresh runs is settled at once. Re-read the ticket with jaira show before settling when its listed lane differs from aboutLane
  note: edge: pressing t while a Translate call is out sets busy=false, so a third press starts a second model call beside the first; keep busy until the running call returns
  note: limit: hotkey x only on the first open question; the others have the click target, as DoD 2 allows
review-check: |-
  1. cd core/mod/builtin/jaira-cockpit && claude plugin test .  - expect 563 pass, 0 fail
  2. claude --plugin-dir core/mod/builtin/jaira-cockpit, let the assistant ask you a question about a ticket
  3. Open 'Waiting for you' from the band: the question shows 'asked today HH:MM', in its original language
  4. Press t: question and options translated; t again: original
  5. Ask the same question once more: 'asked' time stays the old one
  6. Press x: question moves below the rule under 'Answered', dimmed, with 'answered today HH:MM'; nothing is sent into the chat
  7. Ask another question about a ticket, run jaira move <id> --to <other lane> in a shell: within one board refresh it moves to 'Answered'
stop-after-fix: "true"
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
  proof: core/mod/builtin/jaira-cockpit/hooks/threads.ts:270 (settled skips threads absent from after); test 'a thread dropped to keep thirty is not answered by settled()'
- [x] Eine Frage, die waehrend eines laufenden Translate-Aufrufs ankommt, wird danach uebersetzt
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:2127; test 'a question filed while Translate is out is translated once that call is back'
- [x] Dieselbe Frage erneut auf demselben Thread behaelt askedAt
  proof: core/mod/builtin/jaira-cockpit/hooks/threads.ts:150; test 'the same question asked again keeps when it was first asked; another question is asked now'
- [x] Der 30-Thread-Deckel wirft zuerst die am laengsten nicht beruehrten Threads (lastTurn) und nie einen mit offener Frage (threads.ts:168); Test
  proof: core/mod/builtin/jaira-cockpit/hooks/threads.ts:175 (capped); tests 'thirty threads at most: the one touched longest ago goes, never one whose question waits' + 'thirty threads all waiting on a question: none is dropped' + 'over thirty, the newest thread is kept, so no id is handed out twice'

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
- **2026-11-24 11:29 · Sascha** — set stop-after-fix=true: give-back fix is ~20 changed code lines, no lock, no parser
- **2026-11-24 11:29 · Sascha** — waiting: claude plugin test after optimize rename (background) since 12:29
- **2026-11-24 11:31 · Sascha** — waiting: fresh review agent (run-lane.sh --no-worktree, Herdr tab) since 12:31
- **2026-11-24 11:33 · Sascha** — bughunt: pre-existing: core/mod/builtin/jaira-cockpit/hooks/threads.ts:168 - class: the 30-thread cap (next.slice(-THREADS)) drops threads by creation order, not by lastTurn, so a thread with a question still waiting is dropped however recently it was touched, and its question vanishes from 'Waiting for you' without reaching 'Answered' (settled() now skips dropped threads by design); example: thread 1 asks 'still waiting?', 29 more threads open, thread 1 gets a line at turn 31, the 31st new thread is filed at turn 32 -> asksOf(after) is [] and settled() is []; blame: 86da7c4354 (NZZDRH) for the cap, 65446e9473 (1M3X77) for the questions, both before this ticket; ran: node --experimental-strip-types scratch cap.ts against threads.ts at ea5a79df2
- **2026-11-24 11:33 · Sascha** — hunt report: 1 finding, 0 change, 1 pre-existing (note on CSQXKR) | read as a diff since 8a7401450: hooks/threads.ts (file, settled), hooks/register.tsx (change, fileTurn, settleMoved, translateAsks, answerAsk, dismissAsk), tests/board.test.ts; no new file | not run: claude plugin test (the cockpit runner needs Claude Code; left to the testing lane); jaira show exit 5 for an unknown id and exit 0 for a filed-away ticket checked by hand
- **2026-11-24 11:33 · Sascha** — bughunt: done at ea5a79df2 - 1 finding (0 change, 1 pre-existing) in 3 modules, follow-up hunt since 8a7401450; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/threads.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/shipped/shipped.sha256
- **2026-11-24 11:34 · Sascha** — review: done at ea5a79df2
- **2026-11-24 11:35 · Sascha** — report: give-back closed: dropped thread settles nothing, repeated question keeps askedAt, question filed during Translate gets translated; 6 commits 906415bdf..ea5a79df2 | signoff, context unknown | open: person tries it in cockpit (review-check); 2 edge notes + 1 limit from review; branch not pushed
- **2026-11-24 11:36 · Sascha** — Given back from signoff: beheben:
- note: edge: settleMoved compares lanes from the jaira list read at the start of refreshBoard; a question filed and titled (fileTurn) with a newer lane while that refresh runs is settled at once. Re-read the ticket with jaira show before settling when its listed lane differs from aboutLane
- note: edge: pressing t while a Translate call is out sets busy=false, so a third press starts a second model call beside the first; keep busy until the running call returns
angenommen:
- note: limit: hotkey x only on the first open question; the others have the click target, as DoD 2 allows
- **2026-11-24 11:36 · Sascha** — cockpit: given back to a new worker - nothing recorded about the old worker session (cache and context unknown); branch feat/csqxkr is the ticket's worktree's, worktree /home/alex/projects/.worktrees/jaira-csqxkr; reports to this session (uds:/run/user/1000/cc-socks/3417.sock; does not listen)
- **2026-11-24 11:37 · Sascha** — parent session: uds:/run/user/1000/cc-socks/3417.sock; new worker on give-back 2: settleMoved re-reads with jaira show before settling; t while Translate is out keeps busy. No open decisions.
- **2026-11-24 11:38 · Sascha** — teamlead 24.11.2026: Bughunt-Befund (pre-existing, threads.ts:168) in diese Runde gefaltet: ein offener Fragen-Thread darf nicht still aus 'Wartet auf dich' fallen - DoD-Punkt ergaenzt.
- **2026-11-24 11:43 · Sascha** — give-back 2 fixes: settleMoved trusts the list only where it agrees with every aboutLane of the ticket, otherwise asks jaira show (test 'a listed lane other than the one asked in is read again...'); t back to Original no longer clears busy, so a third press starts no second call (test 't pressed again while Translate is out...'); cap: capped() drops threads without an open question, least recent lastTurn first, never the one this turn touched, so with 30 waiting questions the list grows past 30 (teamlead DoD 10). claude plugin test 567 pass; go test ./core/role ok after shipped-history.sh --write
- **2026-11-24 11:43 · Sascha** — set stop-after-fix=true: ~30 code lines, no lock, no parser
- **2026-11-24 11:44 · Sascha** — cleared stop-after-fix: delta since ea5a79df2 is 92+/10- with tests and checksums, over 50
- **2026-11-24 11:44 · Sascha** — critique (give-back 2 delta since ea5a79df2): none. settleMoved's agreement check stays in the caller - moved() keeps taking plain lanes; capped() is local to file(), its only caller; keeping busy across t-off wastes a returned translation when t is off by then, accepted (next press asks again)
- **2026-11-24 11:45 · Sascha** — waiting: bughunt agent (background subagent) since 12:45
- **2026-11-24 11:55 · Sascha** — bughunt: change: blocker: core/mod/builtin/jaira-cockpit/hooks/threads.ts:174 (capped) with threads.ts:144 (new id = max+1) - class: thread ids are reused once the least-recently-touched rule drops the thread with the highest id; before, slice(-30) always kept the newest-created thread, so max+1 was never handed out twice; a filing that names the old number (thread: N) then writes its line/question onto an unrelated thread; precondition: 30 or more questions waiting, so the cap goes over 30 and drops on a filing that touches another thread; example: threads 1-30 each with a question, plain thread #31 'Plain' added, a filing on #1 drops #31, the next new thread is #31 'Other', a filing {thread:31, line:'meant Plain'} lands on 'Other'; ran: scratch script with file() from hooks/threads.ts (prints 'agent addresses #31 -> Other'); fix direction: derive the next id from a counter kept with the threads (or max over the threads seen so far), not from the threads still kept
- **2026-11-24 11:55 · Sascha** — bughunt: change: blocker: core/mod/builtin/jaira-cockpit/hooks/register.tsx:2079 (settleMoved) - class: jaira show exiting 0 with output standOf cannot parse (truncated stdout - the engine's run result carries isStdoutTruncated - or any non-JSON text) becomes lane null, and moved() settles the question as 'ticket left its lane', against the comment two lines up that only exit 5 settles; the uncommitted change now sends every listed ticket whose listed lane differs from the asked lane down this path, so it is reached on every refresh where the list is stale; example: question about MG677R asked in signoff, list says in-progress, jaira show prints '{"handle":"MG677R","body":"long notes ...' with exit 0 -> the pane says 'No question is waiting for you.' and the question sits under Answered; ran: scratch test in a copy of the mod (claude plugin test, world() jaira show stubbed to the truncated text); unclear: the engine's stdout cap was not found, so how long a ticket must be to hit it in practice is unchecked; fix direction: only a parsed status settles, an unparsed one says nothing (like a timeout)
- **2026-11-24 11:55 · Sascha** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:2088 (openAsks) - class: busy is reset to false on every pane opening while a Translate call may still be out, so the guard in toggleTranslation/translateAsks ('pressing t again starts no second call beside it') is bypassed by closing and reopening the pane; example: press t (call held), Esc, reopen with q, press t -> two model calls out at once, the later one overwrites the earlier one's texts; ran: scratch test in a copy of the mod (claude plugin test): held=2 calls=2; cost and a flicker only, not a wrong result - edge
- **2026-11-24 11:55 · Sascha** — hunt report: 3 findings, all change (2 blocker: thread id reuse threads.ts:174/144, unparsed jaira show settles register.tsx:2079; 1 edge: busy reset on reopen register.tsx:2088), notes on CSQXKR | read as a diff since ea5a79df2: register.tsx, threads.ts, board.test.ts, release notes, shipped.sha256; board.ts (net zero change by CSQXKR), texts.ts and types/index.d.ts are not on the last line's list and were read through their CSQXKR hunks plus lanesOf/standOf, not end to end | not checked: the engine's stdout truncation limit for process.run (type definitions not found locally), so how often finding 2 is hit in practice is open; shipped.sha256 hashes were not recomputed against the files
- **2026-11-24 11:55 · Sascha** — bughunt: done at ea5a79df2 - 3 findings (3 change, 0 pre-existing) in 7 modules, follow-up hunt since ea5a79df2 (uncommitted diff); concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/hooks/threads.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/mod/builtin/jaira-cockpit/types/index.d.ts core/release/notes/CSQXKR.md core/shipped/shipped.sha256
- **2026-11-24 11:58 · Sascha** — hunt fixes (round 3): capped() also keeps the newest thread, so max+1 never reuses an id; settleMoved settles only on exit 5 or a parsed status (a filed-away ticket still shows, status done - checked with jaira show P1AE82); openAsks keeps busy so a reopened pane starts no second Translate call. claude plugin test 569/569
- **2026-11-24 11:58 · Sascha** — critique pass 2: read only the three hunt fixes in e65c788c5; each closes its finding in its callee/caller as before; nothing left
- **2026-11-24 11:59 · Sascha** — waiting: bughunt round 3 (background subagent) since 12:59
- **2026-11-24 12:07 · Sascha** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:2117 - class: the Translate busy flag is set without try/finally and is now cleared only by the call's own return, since openAsks (register.tsx:2092) and toggleTranslation (2104) no longer reset it - any throw or hang between busy:true and the final update (store.get/storeNewest in rephrased sit outside its try, register.tsx:1596/1606; a model.complete that never resolves) leaves Translate dead for the rest of the session; example: a store.set of 'rephrased:asks' that throws - the pane reopens with the label 'translating…', and t starts no call (0 model calls) and shows only the originals; before e65c788c5 reopening the pane freed it; ran: scratch copy of the mod with store.set throwing for rephrased: keys, claude plugin test . (label 'translating…' after reopen, 0 calls after the second t); not a blocker: needs a store failure or a hung model call
- **2026-11-24 12:08 · Sascha** — hunt report: 1 finding, 1 change (edge, not blocker), as a note on CSQXKR | read as a diff since ea5a79df2: hooks/register.tsx (settleMoved, openAsks, toggleTranslation, translateAsks with their callers), hooks/threads.ts (file/capped and every consumer of thread ids), tests/board.test.ts, release/notes/CSQXKR.md; no file new to the change | checked and dropped: two Translate calls at once (t pressed during a filing, both orders, scratch test: never more than one in flight) and settleMoved applying a lane read before its jaira show awaits to a question filed meanwhile (scratch test: the engine runs the turn hook only after the refresh handler returns, nothing settled wrongly); whether the real engine serializes handlers the way the test engine does was not checked
- **2026-11-24 12:08 · Sascha** — bughunt: done at e65c788c5 - 1 finding (1 change, 0 pre-existing) in 6 modules, follow-up hunt since ea5a79df2; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/hooks/threads.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/mod/builtin/jaira-cockpit/types/index.d.ts core/release/notes/CSQXKR.md core/shipped/shipped.sha256
- **2026-11-24 12:12 · Sascha** — hunt fix (round 4): translateAsks catches a throwing rephrased (store.get/storeNewest sit outside its own try), so busy is always freed; test 'a Translate call that throws frees t' with w.storeRefuses. Known limit, left: a model.complete that never returns keeps busy until the session ends - the engine gives no timeout to hang on, and the old code had the same pane stuck on 'translating'. 570/570
- **2026-11-24 12:12 · Sascha** — set stop-after-fix=true: round-4 fix is a try/catch plus one test, under 50 lines, no lock, no parser
