---
id: 01M7ZSB3G8TY655YE0AAPDZ444
title: "Cockpit-Band 'Waiting for you': eine Zeile pro Ticket statt Verdikt und alle Notes"
status: done
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
commits:
  - 1cc926124a34517fcb512533e93eeef9d5a7ac5e
  - c0d6a73b16a435892460bd33698d1d744c09a436
  - d16585f9c242f901b9578aadd980b5b1a6cb59ce
  - 32f9a3a1ebc5874aa79b53457e4509622551901f
  - 0ea312aea6e826298dbe1ad31c624a27291a438d
created-at: 2026-11-21T19:29:13Z
updated-at: 2026-11-21T21:10:09Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2166591
claimed-at: 2026-11-21T20:34:45Z
outcome-what: "The cockpit band shows one row per waiting ticket: id, lane, DoD/bughunt marks, blocker and note counts, old-review and question marks; a new Verdict pane holds title, goal, question, verdict and every review finding; the Explain pane opens with title and goal"
outcome-why: "Three tickets in signoff filled ~10 band lines with clipped verdicts and every review note, unreadable"
outcome-resolves: PDZ444
review-summary: "Die Zeile eines Tickets im Band 'Waiting for you' zeigt jetzt nur noch ID, Lane und Marken: DoD ✓/✗/? (aus dod_complete), bughunt ✓/✗/?/– (aus jaira show --brief, im Worktree des Tickets gelesen; ein fail ohne den Hunt-Commit im Checkout wird zu '?'), Zahl der Blocker (rot) und Notes, 'review of round N' bei älterem Verdikt, '?' bei offener Frage. Verdikt-Text, Notes und Frage stehen nicht mehr unter der Zeile, sondern in einer neuen Verdict-Pane, die mit Titel und Ziel des Tickets beginnt, dann Lane, Frage, Verdikt und alle Blocker/Notes; unbekannte Findings einer früheren Runde werden als '? blockers / ? notes' gezählt statt als 0. Die Explain-Pane beginnt ebenfalls mit Titel und Ziel. Release-Note core/release/notes/PDZ444.md, shipped.sha256 nachgezogen."
round: 3
back-from: testing
back-by: lane
back-count: 1
review-gaps: |-
  note: ein wartendes Ticket ohne jeden abgeschlossenen Bughunt (brief bughunt: null) zeigt ein rotes bughunt ✗ - gewollt nach jairas eigener Semantik, aber kein Test benennt diesen Fall
  note: das Label '[blocker]' in der Verdict-Pane steht als Literal in register.tsx statt in texts.ts wie alle anderen Texte
  note: refreshBoard ruft jetzt pro sichtbarem Ticket jede Minute zusätzlich jaira show --brief und ggf. git merge-base auf - bei BAND_TICKETS Zeilen unkritisch, bei aufgeklapptem Band wächst es linear
  note: pre-existing: core/mod/builtin/jaira-cockpit/node_modules/ ist nicht gitignored und taucht als ?? im Worktree auf
test-verdict: "pass: round 2 - claude plugin test . 286 pass rc=0, go test ./core/shipped rc=0, shipped-history --check rc=0; finding 1 class checked live: F1V59P and AX8K02 read from the main checkout give state=fail with merge-base --is-ancestor rc=1, so the band shows '?' (test 'a hunt's fail is a red cross only in a checkout that holds the hunt's commit' covers held and not held); finding 2: same-round verdict without gaps key still counts 0/0, earlier round without kept copy shows '? blockers | ? notes' and the pane line (test 'an earlier review whose notes were not kept...'); DoD proofs unchanged and still true"
hunt-verdict: "pass: bughunt done at d16585f9c, 0 change, 0 pre-existing; both round-1 change blockers fixed (hunt at c0d6a73b1)"
review-verdict: "Der Diff erfüllt alle fünf DoD-Punkte; claude plugin test: 286 pass, 0 fail; go test ./core/shipped ok. Kein Blocker, nur kleine Notes."
review-round: 3
review-verdict-gaps: |-
  note: ein wartendes Ticket ohne jeden abgeschlossenen Bughunt (brief bughunt: null) zeigt ein rotes bughunt ✗ - gewollt nach jairas eigener Semantik, aber kein Test benennt diesen Fall
  note: das Label '[blocker]' in der Verdict-Pane steht als Literal in register.tsx statt in texts.ts wie alle anderen Texte
  note: refreshBoard ruft jetzt pro sichtbarem Ticket jede Minute zusätzlich jaira show --brief und ggf. git merge-base auf - bei BAND_TICKETS Zeilen unkritisch, bei aufgeklapptem Band wächst es linear
  note: pre-existing: core/mod/builtin/jaira-cockpit/node_modules/ ist nicht gitignored und taucht als ?? im Worktree auf
review-check: "1. claude --plugin-dir /home/alex/projects/.worktrees/jaira-pdz444/core/mod/builtin/jaira-cockpit  (in einem Repo mit Tickets in signoff/human)  2. Im Band 'Waiting for you': jedes Ticket ist genau eine Zeile 'ID · lane · DoD ✓ · bughunt ✓ · 0 blockers · 2 notes' - kein Verdikt-Text, keine '·'-Notes darunter  3. Verdict drücken: Pane zeigt oben fett den Titel, darunter 'Goal: …', dann Lane, Verdict und alle Review-Notes ([blocker] rot zuerst)  4. Escape oder Close schliesst sie  5. Bei einem Ticket mit Explain: Explain-Pane beginnt ebenfalls mit Titel und Goal  6. Automatisch: cd core/mod/builtin/jaira-cockpit && claude plugin test .  → 286 pass"
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
- **2026-11-21 20:49 · Alexander Sacharov** — testing: note: cost - an expanded band runs one extra 'jaira show --brief --json' per waiting ticket each refresh (16 on this board today), each doing the git hunt-coverage check; fine at 3 shown, watch it if the expanded band feels slow
- **2026-11-21 20:52 · Alexander Sacharov** — bughunt: change: blocker: core/mod/builtin/jaira-cockpit/hooks/register.tsx:387 - class: the band reads the hunt state (jaira show --brief) in the main checkout whenever the ticket has no worktree whose branch/path ends in its slug (worktree removed, work on another ticket's branch, a worktree named otherwise, or git worktree list failing); jaira's HuntCoverage answers fail when the checkout does not contain the hunt's commit, so a hunted ticket shows a red 'bughunt ✗' that is not true. Same fallback feeds the Explain prompt's 'bughunt on HEAD' (ticketDir, register.tsx:516), which predates this ticket. example: F1V59P and AX8K02 read state=ok in their worktrees and state=fail from /home/alex/projects/jaira (main=fail wt=ok). ran: jaira show <id> --brief --json in both checkouts
- **2026-11-21 20:52 · Alexander Sacharov** — bughunt: change: blocker: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1518-1521 - class: an earlier-round verdict without a verdict_gaps copy (rounds stamped since 13PQB3/2026-10-12, copy kept only since DF0YD6/2026-11-21) gets blockers=[] and notes=[] from reviewGaps() (board.ts:229-232, which means 'unknown'), and the new row prints that as a definite '0 blockers | 0 notes' - also the Verdict pane lists no notes; the review's findings are unknown, not zero. example: 1K3RA8 (human, round 10, review_round 9, verdict 'Zurueck nach in-progress wegen Gap 1...', no verdict_gaps) yields blockers=0 notes=0 gaps=''. ran: node --experimental-strip-types on snapshotOf() with the live jaira lanes/list JSON
- **2026-11-21 20:52 · Alexander Sacharov** — bughunt: done at c0d6a73b1 - 2 findings (2 change, 0 pre-existing) in 5 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/mod/builtin/jaira-cockpit/types/index.d.ts core/release/notes/PDZ444.md core/shipped/shipped.sha256
- **2026-11-21 20:53 · Alexander Sacharov** — fix plan for the two hunt findings: (1) a brief's 'fail' is only believed when the hunt's commit is in the checkout it was read in (git merge-base --is-ancestor <at> HEAD in that dir); otherwise the mark is '?' - fixes every no-own-worktree case at the cause, not by finding the worktree better. (2) Waiting gets findingsKnown = reviewGaps(t) !== undefined; unknown -> row '? blockers | ? notes', Verdict pane says the notes of that round were not kept.
- **2026-11-21 20:59 · Alexander Sacharov** — critique 3: the round-2 fix has the right shape - the merge-base check sits where the hunt state is read (refreshBoard), so every no-own-worktree case is covered at the cause; findingsKnown reuses reviewGaps/sameRound, no new rule.
- **2026-11-21 21:01 · Alexander Sacharov** — bughunt: done at d16585f9c - 0 findings in 5 modules, follow-up hunt (diff since c0d6a73b1: board.ts huntOf/findingsKnown, register.tsx refreshBoard ancestor check and band/Verdict pane, texts.ts, board.test.ts, shipped.sha256 hashes match); concurrency, error paths, input, permissions checked; ran: claude plugin test . (286 pass); files: core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/mod/builtin/jaira-cockpit/types/index.d.ts core/release/notes/PDZ444.md core/shipped/shipped.sha256
