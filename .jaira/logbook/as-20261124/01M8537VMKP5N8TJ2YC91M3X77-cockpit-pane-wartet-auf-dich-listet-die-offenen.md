---
id: 01M8537VMKP5N8TJ2YC91M3X77
title: "Cockpit: Pane 'Wartet auf dich' listet die offenen Fragen der Sitzung, Antwort per Taste statt Suchen im Chat"
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Die Person findet jede offene Frage, die die Sitzung ihr gestellt hat, an einer Stelle - mit Ticket, Titel, Frage und Optionen - und beantwortet sie dort per Taste, ohne im Chat zu suchen und ohne dass die Sitzung blockiert."
context: |-
  Was falsch ist: der Teamlead stellt Fragen als Text im Chat ('Naechster Schritt: antworte ja zu VTSF24'). Danach kommen Worker-Meldungen, die Frage rutscht nach oben, die Person scrollt und sucht.
  Ausloeser: Alex 23.11.2026, Screenshot: 'ich suche staendig Antworten im Chat'. Abgelehnte Alternative: jede Frage als AskUserQuestion - blockiert die Sitzung, wenn die Person weggeht, und der Workflow steht. Alex will nicht blockieren.
  Schon da: core/mod/builtin/jaira-cockpit/hooks/threads.ts - nach jedem Turn ordnet ein billiges Modell den Turn einem Thread zu und haelt eine offene Frage fest (Thread.question). Das Band zeigt nur die Zahl: texts.ts:262 'threads: N open question(s), M left lying'. Die Liste selbst sieht man nicht bequem (threadsTitle texts.ts:264).
  Was fehlt: (1) eine Pane/Liste 'Wartet auf dich' aus den Threads mit question, je Eintrag Ticket-ID, Ticket-Titel, die Frage in einer Zeile und, wenn die Frage Optionen nannte, diese nummeriert; (2) erreichbar vom Band aus (Klick auf 'threads: N open' oder eine Taste); (3) Antwort: Taste/Klick schickt die Wahl als Nachricht der Person in die Sitzung (z.B. 'VTSF24: ja'), die Frage verschwindet aus der Liste, wenn der naechste Turn sie beantwortet.
  Die Optionen muss der Klassifizierer mit aus dem Turn ziehen (threads.ts classifyPrompt: Feld options neben question).
  Ueberschneidung: M31BYV aendert register.tsx im Give-back-Teil; diese Arbeit gehoert in threads.ts/eine eigene Pane und nur die Verdrahtung in register.tsx.
definition-of-done: Der Klassifizierer haelt zu einer offenen Frage das Ticket und die genannten Optionen fest (threads.ts); Test
tags:
  - cockpit
blocked-by: []
related: []
commits:
  - 65446e9473f76adc403d7afe981918b378da3323
  - b318241d6d24ec4d9ee0fa85deaf0cc4eda74422
  - 628b531185a10724b51f7a12ce2cd78f1e0c44b7
  - a78bdb00bcffdbf38027e71183cf6c56f87d89bd
created-at: 2026-11-23T20:58:25Z
updated-at: 2026-11-24T06:08:33Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-20021
claimed-at: 2026-11-23T21:49:51Z
outcome-what: "Cockpit: band entry 'threads: N open question(s)' is a button (q) opening a 'Waiting for you' pane; each open question with ticket, title and numbered options; a digit sends '<ticket>: <option>' into the session, 'Answer in the chat' fills the prompt"
outcome-why: The person searched the chat for questions the teamlead asked; AskUserQuestion would block the session
outcome-resolves: 1M3X77
review-summary: "Full review of 6e311f200..628b53118. The thread classifier now also returns, for a question it finds, the ticket it is about ('about', a 6-char id) and the options the assistant offered (max 9, clipped). Both live on the thread and leave with the question (answered()). After filing, fileTurn looks up each named ticket's title once via 'jaira show --json' (titleOf); a timeout is retried after the next turn, a missing title falls back to the thread's title. A new atom 'asks' (asksOf: open questions, newest first) feeds a new pane 'Waiting for you': ticket id + title, the question, one Button per option (the engine draws digits 1-9 in front of them), plus 'Answer in the chat…'. The band's threads entry is now a Button with hotkey q that opens the pane. A pick sends '<ticket>: <option>' with a non-blocking $.prompt.submit, clears the question locally and closes the pane when nothing is left. 'Answer in the chat' fills the prompt with '<ticket>: '. Release note and shipped.sha256 updated; 6 new tests plus parse cases, suite 539 pass."
round: 4
back-from: testing
back-by: lane
back-count: 2
review-gaps: |-
  note: limit: the turn a pick starts is never classified (the engine runs none of the mod's own hooks on what it submits) - a question the assistant asks back in direct reply to a pick does not reach the pane; stated in answerAsk's comment and the notes, not in the release note
  note: limit: the ticket title is looked up only for a ticket the question names (about); a question that only rides on the thread's ticket shows the thread's title (threads.ts asksOf)
  note: limit: band hotkey q - whether it fires while the person types in the composer was not checked; the test harness does not model key routing
  note: edge: a new question in the same thread drops aboutTitle (file() builds it from answered(base)), so a repeated question about the same ticket triggers a second jaira show; harmless cost, no wrong output
  note: edge: answerAsk clears the question even if \$.prompt.submit rejects (void, no catch) - same pattern as decide(); the answer would be lost silently
hunt-verdict: "pass: round 3 at 628b53118, 0 findings"
test-verdict: "pass: claude plugin test 539 pass RC=0, shipped-history --check 0, go test core/shipped core/release ok; DoD 1-4 verified in the tree (pane, hotkey q, digit pick submits '<ticket>: <option>', question leaves, title looked up once); failure path: unknown ticket keeps thread title, timeout retried; writes: thread store only, pick and filing reassign from the current list after each await"
review-verdict: "accept: all four DoD items are met in the code and tested; no blocker found - the remaining gaps are stated limits and two small edges"
review-round: 4
review-verdict-gaps: |-
  note: limit: the turn a pick starts is never classified (the engine runs none of the mod's own hooks on what it submits) - a question the assistant asks back in direct reply to a pick does not reach the pane; stated in answerAsk's comment and the notes, not in the release note
  note: limit: the ticket title is looked up only for a ticket the question names (about); a question that only rides on the thread's ticket shows the thread's title (threads.ts asksOf)
  note: limit: band hotkey q - whether it fires while the person types in the composer was not checked; the test harness does not model key routing
  note: edge: a new question in the same thread drops aboutTitle (file() builds it from answered(base)), so a repeated question about the same ticket triggers a second jaira show; harmless cost, no wrong output
  note: edge: answerAsk clears the question even if \$.prompt.submit rejects (void, no catch) - same pattern as decide(); the answer would be lost silently
review-check: "1. In this worktree start: claude --plugin-dir /home/alex/projects/.worktrees/jaira-1m3x77/core/mod/builtin/jaira-cockpit  2. Ask the session something that makes it ask you back with options, e.g. 'Ask me whether ticket 1M3X77 should go to done or back - give me two options'  3. After the answer the band shows 'threads: 1 open question(s)' in brackets  4. Press q (or click it) - the pane 'Waiting for you' opens with 1M3X77, its title, the question and the options numbered 1: 2:  5. Press 1 - the chat gets '1M3X77: <option 1>' as your message, the pane closes, the band entry disappears  6. Repeat 2-4 and choose 'Answer in the chat…' - the prompt is filled with '1M3X77: ' and waits for your text  7. Automated: cd core/mod/builtin/jaira-cockpit && claude plugin test . shows 539 pass"
---

# Cockpit: Pane 'Wartet auf dich' listet die offenen Fragen der Sitzung, Antwort per Taste statt Suchen im Chat

## Definition of Done

- [x] Der Klassifizierer haelt zu einer offenen Frage das Ticket und die genannten Optionen fest (threads.ts); Test
  proof: core/mod/builtin/jaira-cockpit/hooks/threads.ts:69 (about/options in classifyPrompt, parseFiling, file); tests/board.test.ts:3577 'ride with the question', :3629 'asked about once', options:/about: cases
- [x] Eine Pane 'Wartet auf dich' listet jede offene Frage: Ticket-ID, Ticket-Titel, Frage, nummerierte Optionen; erreichbar vom Band (Klick auf die threads-Anzeige) und per Taste; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:2657 ASKS pane, :2399 band Button hotkey q; tests/board.test.ts:3600 'the band opens the pane', :3643 'one button holds the focus'
- [x] Eine Option waehlen schickt sie als Nachricht der Person in die laufende Sitzung, ohne die Sitzung zu blockieren; eine beantwortete Frage verschwindet aus der Liste; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:2007 answerAsk (prompt.submit, question cleared); tests/board.test.ts:3600, :3659
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/1M3X77.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 20:59 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/48662.sock
- **2026-11-23 21:05 · Alexander Sacharov** — decisions (worker, settled from code, no person needed): (1) classifier gets two new fields 'about' (ticket the question is about) and 'options' (choices the question named); stored on the thread, cleared with the question. (2) ticket title: looked up once with 'jaira show <id> --json' after filing, kept on the thread; falls back to the thread title. (3) pane opened by the band's threads entry turned into a Button with hotkey 'q' (same pattern as code-modernization's band 'show pane' hotkey 'm'); digits pick options like the Explanation pane. (4) pick sends '<ticket>: <option>' via $.prompt.submit (non-blocking, the pattern decide() uses) and clears the question locally at once; 'answer in chat' fills the prompt with '<ticket>: ' via $.prompt.fill. Baseline: claude plugin test . 522 pass.
- **2026-11-23 21:18 · Alexander Sacharov** — critique: fileTurn writes the store twice on every turn; skip the second write when no question names an untitled ticket. Left standing on purpose: two atoms (threads tally keeps the stale count the band needs, asks carries the list) - deriving one from the other would mean the band reads the whole list each render; plugin-origin prompt.submit as the person's answer is the pattern decide() already uses.
- **2026-11-23 21:18 · Alexander Sacharov** — critique round 2: fileTurn early return in place (register.tsx fileTurn); nothing left.
- **2026-11-23 21:19 · Alexander Sacharov** — optimize: searched for an existing answer-into-session helper (decide() submits t.decisionFor and notes on the ticket - different: a thread question has no ticket note to write) and a title reader (storyOf parses the whole story; titleOf is the one-field read beside diskOf). No dead code, nothing orphaned.
- **2026-11-23 21:21 · Alexander Sacharov** — waiting: bughunt (background subagent) since 22:21; gates: claude plugin test 537 pass RC=0, shipped-history --check 0, go test core/shipped core/release ok
- **2026-11-23 21:24 · Alexander Sacharov** — bughunt: change: blocker: core/mod/builtin/jaira-cockpit/hooks/register.tsx:2001 - class: an answer picked in the 'Waiting for you' pane goes in as a plugin-origin prompt, which isPersonsPrompt (hooks/cache.ts:60, used at register.tsx:2174) does not count as the person's, so the turn it starts is never filed by fileTurn; any follow-up question the assistant asks in reply to a pick never reaches the pane, and the thread gets no line for that turn; example: pick 'ja' on 'Antworte ja zu VTSF24?', the assistant answers 'Gut. Und wann?' - the pane says 'No question is waiting for you.'; ran: scratch test in a copy of the mod (claude plugin test): model calls stay 1 -> 1 after the picked turn completes
- **2026-11-23 21:24 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1964 - class: a ticket id whose title lookup fails (jaira show non-zero, or no title) is never remembered as looked up, so untitled() returns it again and jaira show runs for it after every turn for as long as the question stays open, despite the comment 'read once'; example: the cheap model names about 'ZZZZZZ' (a word or an id from another board) - 4 turns run 'jaira show ZZZZZZ --json' 4 times; ran: scratch test in a copy of the mod (4 show calls over 4 turns)
- **2026-11-23 21:24 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:2685 - class: autoFocus is set by two independent rules (n === 1 on the first option across the whole pane, and the chat button of the first question when it has no options), so two buttons claim focus whenever the newest question has no options and an older one has; example: newest 'b?' (no options) and older 'a?' (x, y) - both ask-2-chat and ask-1-0 carry autoFocus, so Enter may send 'x' to the older question the person was not looking at; ran: scratch test in a copy of the mod (props of both buttons autoFocus=true)
- **2026-11-23 21:24 · Alexander Sacharov** — hunt report: 3 findings, all change (1 blocker: a pick's own turn is never filed), notes on 1M3X77 | read whole: hooks/threads.ts, hooks/texts.ts asks entries, types/index.d.ts, release note; register.tsx read through the threads section (threads/fileTurn/keepThreads/answerAsk), the prompt.submit and turn.complete hooks, the band and the asks pane; board.ts titleOf and its parse helpers; tests/board.test.ts 1M3X77 block and harness; shipped.sha256 not read (a hash list) | not checked: whether the engine's band hotkey 'q' fires while the person types in the composer - no engine docs at hand, and the test harness does not model key routing
- **2026-11-23 21:24 · Alexander Sacharov** — bughunt: done at 65446e947 - 3 findings (3 change, 0 pre-existing) in 7 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/hooks/threads.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/mod/builtin/jaira-cockpit/types/index.d.ts core/release/notes/1M3X77.md core/shipped/shipped.sha256
- **2026-11-23 21:32 · Alexander Sacharov** — found: the engine runs none of the mod's own hooks on what the mod submits itself (probed: prompt.submit hook never saw the pick, store marker test). So a pick cannot be recognised in prompt.submit; turn.complete takes it from a queue when the turn had no person's prompt. Known limit: a pick made while an agent-report turn runs is filed against that turn's answer (mild misfile, no loss).
- **2026-11-23 21:36 · Alexander Sacharov** — critique round 3 (fix delta 65446e947..b318241d6): queue in turn.complete is the smallest shape given the engine skips own hooks; nothing left. optimize: delta adds no duplicate or dead code (picked Set + hook branch removed again).
- **2026-11-23 21:36 · Alexander Sacharov** — waiting: bughunt round 2 (background subagent) since 22:36; suite 539 pass RC=0 on b318241d6
- **2026-11-23 21:37 · Alexander Sacharov** — bughunt: change: blocker: core/mod/builtin/jaira-cockpit/hooks/register.tsx:2204 - class: any turn that ends without personSaid takes the oldest picked pane answer and is filed under it, not only the answer's own turn - agent-origin prompts (task-notification, peer, coordinator) return before personSaid is set (2167), slash commands and the mod's own submits (742, 796, 1000, 1016, 1167, 1362, 1785) leave it undefined, and a queued answer the person removes before it runs stays in picked for the next such turn; example: the person picks an option in 'Waiting for you' while a turn started by a worker's report is still running - that report turn's answer is classified with the picked line as the person's words (fileTurn, which is meant never to run on an agent report), and the answer's own turn finds picked empty and is not filed, so a question asked back in reply never reaches the pane; also an interrupted turn (e.answer === '') consumes the entry and drops it; read (prompt.submit 2160-2184, turn.complete 2194-2208, AGENT_ORIGINS filter.ts:34) - no scratch run, the engine queue is not reachable from bun tests
- **2026-11-23 21:37 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1976 - class: every failed title lookup is stored as 'jaira has no title' (aboutTitle '') and never retried - a timeout of run() (15 s, returns undefined), a non-zero exit for any reason, or unparsable stdout all count the same as an unknown ticket; example: jaira show times out once while the repo is busy, the ticket's question shows the thread's title for the rest of the session though the ticket exists; read (run() 214-220, titled threads.ts:214)
- **2026-11-23 21:37 · Alexander Sacharov** — hunt report: 2 findings, both change (1 blocker: picked answer taken by the wrong turn; 1 edge: failed title lookup never retried), notes on 1M3X77 | follow-up hunt since 65446e947, read as a diff with surrounding handlers: register.tsx, threads.ts, board.test.ts, shipped.sha256; no new files | the picked/turn ordering is confirmed by reading only - the engine's prompt queue cannot be driven from the bun tests
- **2026-11-23 21:38 · Alexander Sacharov** — bughunt: done at b318241d6 - 2 findings (2 change, 0 pre-existing) in 4 modules, follow-up hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/hooks/threads.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/mod/builtin/jaira-cockpit/types/index.d.ts core/release/notes/1M3X77.md core/shipped/shipped.sha256
- **2026-11-23 21:42 · Alexander Sacharov** — decision (worker, edges grow -> smaller design): pick queue removed. Probed twice: the mod's own prompt.submit reaches none of its hooks, and $.prompt.submit({wait:true}) returns only {text, origin}, no answer. Nothing tells the pick's turn apart from an agent's, so that turn is not filed. Known limit for the person: a question asked back in direct reply to a pick does not reach the pane. Title lookup: a timeout (run() undefined) is retried after the next turn; a non-zero exit or no title is final.
- **2026-11-23 21:42 · Alexander Sacharov** — critique/optimize round 4 (delta b318241d6..HEAD): the change only removes the queue and adds one guard; nothing left.
- **2026-11-23 21:43 · Alexander Sacharov** — waiting: bughunt round 3 (background subagent) since 22:43; suite 539 pass on 628b53118
- **2026-11-23 21:45 · Alexander Sacharov** — hunt report: 0 findings, nothing filed | read as a diff since b318241d6: hooks/register.tsx (fileTurn title retry, answerAsk, turn.complete), tests/board.test.ts, core/shipped/shipped.sha256; no file new to the change, none read whole | checked against the engine contract that $.process.run rejects on timeout (so the retry fires) and that prompt.submit fires at queue time; cockpit tests not run here, that is the testing lane's
- **2026-11-23 21:45 · Alexander Sacharov** — bughunt: done at 628b53118 - 0 findings (0 change, 0 pre-existing) in 4 modules, follow-up hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/hooks/threads.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/mod/builtin/jaira-cockpit/types/index.d.ts core/release/notes/1M3X77.md core/shipped/shipped.sha256
- **2026-11-23 21:45 · Alexander Sacharov** — waiting: fresh review (Herdr tab via run-lane.sh) since 22:45
- **2026-11-23 21:50 · Alexander Sacharov** — review: done at 628b53118
- **2026-11-23 21:50 · Alexander Sacharov** — report: cockpit 'Waiting for you' pane (band entry -> button, q), open questions with ticket, title, numbered options; digit sends '<ticket>: <option>'; 3 commits, 539 tests | signoff, context 203k of 1000k, cache warm 60 min after last answer | person: try review-check steps; 3 limits + 2 edges in review-gaps (main: a question asked back in reply to a pick is not listed)
- **2026-11-23 21:58 · Alexander Sacharov** — Accepted from signoff: angenommen:
- note: edge: a new question in the same thread drops aboutTitle (file() builds it from answered(base)), so a repeated question about the same ticket triggers a second jaira show; harmless cost, no wrong output
- note: edge: answerAsk clears the question even if \$.prompt.submit rejects (void, no catch) - same pattern as decide(); the answer would be lost silently
- note: limit: the turn a pick starts is never classified (the engine runs none of the mod's own hooks on what it submits) - a question the assistant asks back in direct reply to a pick does not reach the pane; stated in answerAsk's comment and the notes, not in the release note
- note: limit: the ticket title is looked up only for a ticket the question names (about); a question that only rides on the thread's ticket shows the thread's title (threads.ts asksOf)
- note: limit: band hotkey q - whether it fires while the person types in the composer was not checked; the test harness does not model key routing
- **2026-11-23 21:58 · Alexander Sacharov** — cockpit: accepted by the person, signoff -> done
- **2026-11-23 21:58 · Alexander Sacharov** — pr: push stopped — merging upstream/master conflicts in core/mod/builtin/jaira-cockpit/hooks/register.tsx (one hunk; upstream brought DVPAPE, RGFJ0Q, M31BYV). Merge aborted, nothing pushed. Options: 1) (recommended) a worker resolves the conflict on feat/1m3x77, regenerates shipped.sha256, runs the suite, then /jaira-role-pr again; 2) resolve by hand; 3) push the branch without the merge and let the PR show the conflict.
- **2026-11-23 22:00 · Alexander Sacharov** — pr: conflict resolved in a78bdb00b (import list union: parse from master, titleOf from the branch), shipped.sha256 regenerated; check-ticket-commits and shipped-history --check ok; cockpit plugin tests 542 pass. Stopped on usage limit while go vet + go test ./... -race ran: not pushed, no PR. Rerun /jaira-role-pr 1M3X77 --by-person from step 3.
