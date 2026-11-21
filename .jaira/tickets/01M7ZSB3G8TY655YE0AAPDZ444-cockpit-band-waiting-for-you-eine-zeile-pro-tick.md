---
id: 01M7ZSB3G8TY655YE0AAPDZ444
title: "Cockpit-Band 'Waiting for you': eine Zeile pro Ticket statt Verdikt und alle Notes"
status: testing
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Im Band 'Waiting for you' steht jedes Ticket in einer kurzen Zeile, die auf einen Blick sagt, ob es angenommen werden kann; die Notes sind erst auf Klick da."
context: |-
  Was heute falsch ist (Screenshot Alex 21.11., 3 Tickets in signoff):
  - Das Band zeigt pro Ticket den vollen Review-Verdikt-Satz, abgeschnitten mit "…", und darunter jede Note aus review-gaps als eigene Zeile.
  - Bei 3 Tickets sind das schon ~10 Zeilen, der Rest steht unter "↓ 6 more"; lesen kann man es nicht.
  - Seit 0.5.15 listet das Cockpit die Review-Notes unter der Zeile (Release-Note "the review's other notes are listed under the ticket's row") — das ist der Teil, der zu lang wurde.
  Code: core/mod/builtin/jaira-cockpit/hooks/register.tsx (Band), board.ts (review.verdict_gaps, Zeile ~212), texts.ts (waiting).
  Vorschlag:
  - Eine Zeile: ID | Lane | Marken (DoD ✓/✗, Bughunt auf HEAD ✓/✗/?, Blocker-Zahl) | Notes: N | [Accept] [Give back] [Verdict]
  - Kein Verdikt-Text, keine Notes im Band. Die Notes zeigt Verdict/Explain (dort stehen sie schon klassifiziert).
  - Marken aus Daten, die es gibt: dod_complete, bughunt.state aus jaira show --brief, blocker:-Zeilen.
  Nach MFYM1A starten: MFYM1A baut das Give-back-Pane in denselben Dateien um.
definition-of-done: Ein Ticket im Band ist genau eine Zeile; Verdikt-Text und Notes stehen nicht im Band (Test)
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-11-21T19:29:13Z
updated-at: 2026-11-21T20:48:30Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2166591
claimed-at: 2026-11-21T20:34:45Z
outcome-what: "The cockpit band shows one row per waiting ticket: id, lane, DoD/bughunt marks, blocker and note counts, old-review and question marks; a new Verdict pane holds title, goal, question, verdict and every review finding; the Explain pane opens with title and goal"
outcome-why: "Three tickets in signoff filled ~10 band lines with clipped verdicts and every review note, unreadable"
outcome-resolves: PDZ444
review-summary: none
round: 2
back-from: critique
back-by: lane
back-count: 1
review-gaps: "folded the two new 'waiting or theirs by id' lookups (Verdict and Explain panes) into waitingTicket() in register.tsx; left the older inline copy in refreshBoard's chat loop (it reads the fresh view, not the atom); left blockersOf/notesOf each splitting the gaps once - per refresh, a handful of lines"
---

# Cockpit-Band 'Waiting for you': eine Zeile pro Ticket statt Verdikt und alle Notes

## Definition of Done

- [x] Ein Ticket im Band ist genau eine Zeile; Verdikt-Text und Notes stehen nicht im Band (Test)
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: 'a ticket in the band is one row: its marks, no verdict, no notes, no question under it'
- [x] Die Zeile zeigt DoD-, Bughunt- und Blocker-Marke und die Zahl der Notes
  proof: hooks/register.tsx marks() in the band render; same test asserts 'MG677R | signoff | DoD ✓ | bughunt ✓ | 1 blocker | 2 notes | review of round 2 | ?'
- [x] Notes und Verdikt sind über Verdict erreichbar, nichts geht verloren
  proof: hooks/register.tsx VERDICT pane (question, verdict, [blocker] + · note lines); tests 'Explain is offered for a blocker or a reservation...' and 'a verdict of an earlier round lists the review's own notes' read them from the pane
- [x] Release-Note core/release/notes/<ID>.md
  proof: core/release/notes/PDZ444.md
- [x] Das Verdict-/Explain-Pane beginnt mit Titel und Ziel (goal) des Tickets, damit man weiss, worum es geht (Alex 21.11.)
  proof: hooks/register.tsx about() at the top of the VERDICT and EXPLANATION panes; test 'the Explain pane starts with the ticket's title and goal'

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Waiting gets goal, dodComplete, blockers count; notes count = all review findings (board.ts) + tests
- [x] refresh: bughunt state per shown ticket via show --brief --json in ticketDir, into an atom
- [x] band row: one line with marks, no verdict/notes/question/reviewedIn lines (register.tsx, texts.ts)
- [x] Verdict pane: title, goal, lane/round, question, verdict, all findings; Explain pane starts with title+goal
- [x] tests: band is one row per ticket; pane holds verdict+notes
- [x] release note core/release/notes/PDZ444.md; shipped.sha256; plugin tests + go test ./core/shipped

## Progress
- **2026-11-21 20:34 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2274544.sock
- **2026-11-21 20:35 · Alexander Sacharov** — decisions (settled, no person needed): (1) one row = 'ID · lane · DoD ✓/✗ · bughunt ✓/✗/? · blockers N · notes N' + buttons; question line, 'reviewed in round N' line and note lines leave the band too (DoD 1 says exactly one line) - the row keeps a dim '?' mark when a question is open, the Discuss button stays. (2) Verdict button shows whenever there is a verdict, notes or a question; it opens a pane: title, goal, then lane/round line, question, verdict, then every review finding (blockers first) - nothing lost. (3) Bughunt state comes from 'jaira show <id> --brief --json' (bughunt.state; null = ✗, unknown/skipped = dim ?/–), run in the ticket's own dir (ticketDir) because HuntCoverage reads that checkout; only for the shown tickets, beside the existing show --json. (4) DoD mark from list --json dod_complete; goal added to Waiting from list --json. (5) Explain pane gets title + goal on top via the waiting entry.
- **2026-11-21 20:45 · Alexander Sacharov** — in-progress: bughunt mark reads 'jaira show <id> --brief --json' per shown ticket in its worktree (one extra jaira call per shown ticket per refresh, max 3 unless expanded). An earlier-round verdict without verdict_gaps counts 0 blockers/0 notes - same rule as before (reviewGaps), the row then shows 'review of round N' in warning colour. Verdict button now shows also for a question alone (human lane) since the question left the row.
- **2026-11-21 20:46 · Alexander Sacharov** — critique 1: drop the stray Text argument of marks(); showVerdict toast falls back to the question when there is no verdict. Left as is: two jaira show calls per shown ticket (show --json carries the body for explanationNotes, only --brief carries bughunt; --brief keeps 4 notes, so it cannot replace the body read).
- **2026-11-21 20:47 · Alexander Sacharov** — critique 2: nothing left; the two round-1 findings are fixed (register.tsx marks call, showVerdict toast).
