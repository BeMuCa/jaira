---
id: 01M8537VMKP5N8TJ2YC91M3X77
title: "Cockpit: Pane 'Wartet auf dich' listet die offenen Fragen der Sitzung, Antwort per Taste statt Suchen im Chat"
status: critique
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
commits: []
created-at: 2026-11-23T20:58:25Z
updated-at: 2026-11-23T21:42:41Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-40547
claimed-at: 2026-11-23T20:59:00Z
outcome-what: "Cockpit: band entry 'threads: N open question(s)' is a button (q) opening a 'Waiting for you' pane; each open question with ticket, title and numbered options; a digit sends '<ticket>: <option>' into the session, 'Answer in the chat' fills the prompt"
outcome-why: The person searched the chat for questions the teamlead asked; AskUserQuestion would block the session
outcome-resolves: 1M3X77
review-summary: "none"
round: 4
back-from: testing
back-by: lane
back-count: 2
review-gaps: none
hunt-verdict: "fail: round 2, 2 change findings (pick queue misfiles agent turns; timeout stored as no title)"
test-verdict: "fail: suite green; sent back on round-2 hunt findings"
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
