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
updated-at: 2026-11-23T21:17:41Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-40547
claimed-at: 2026-11-23T20:59:00Z
outcome-what: "Cockpit: band entry 'threads: N open question(s)' is a button (q) opening a 'Waiting for you' pane; each open question with ticket, title and numbered options; a digit sends '<ticket>: <option>' into the session, 'Answer in the chat' fills the prompt"
outcome-why: The person searched the chat for questions the teamlead asked; AskUserQuestion would block the session
outcome-resolves: 1M3X77
---

# Cockpit: Pane 'Wartet auf dich' listet die offenen Fragen der Sitzung, Antwort per Taste statt Suchen im Chat

## Definition of Done

- [x] Der Klassifizierer haelt zu einer offenen Frage das Ticket und die genannten Optionen fest (threads.ts); Test
  proof: core/mod/builtin/jaira-cockpit/hooks/threads.ts:69 (about/options in classifyPrompt, parseFiling, file); tests/board.test.ts:3577 'ride with the question' + options:/about: cases
- [x] Eine Pane 'Wartet auf dich' listet jede offene Frage: Ticket-ID, Ticket-Titel, Frage, nummerierte Optionen; erreichbar vom Band (Klick auf die threads-Anzeige) und per Taste; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:2649 ASKS pane, :2391 band Button hotkey q; tests/board.test.ts:3600 'the band opens the pane'
- [x] Eine Option waehlen schickt sie als Nachricht der Person in die laufende Sitzung, ohne die Sitzung zu blockieren; eine beantwortete Frage verschwindet aus der Liste; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1999 answerAsk (prompt.submit, question cleared); tests/board.test.ts:3600, :3629
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
