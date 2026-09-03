---
id: 01M1KN1JEG6D79HXJ20X76WCCW
title: Der verwaltete Block schickt ein gestartetes Ticket bis zum naechsten Human-Step
status: in-progress
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Sagt der User 'starte/arbeite Ticket X', faehrt die Session es Lane fuer Lane durch - Loops (critique/optimize) inklusive - bis es in einer Human-Lane liegt, und nach der Klaerung weiter; sagt der User 'ein Agent soll es bearbeiten', babysittet ein Subagent es genauso durch"
context: "Berk am 03.09.: er promptet das heute jedes Mal von Hand ('babysitte das Ticket mit einem Subagent durch bis zum naechsten human step'). Die Regel soll in den verwalteten jaira-Block der CLAUDE.md/AGENTS.md, damit sie ohne erneutes Tippen gilt. Das ist Mechanismus c aus 88H1P4 (Block-Satz), der explizit auf Berks Go wartete - diese Anfrage ist das Go, erweitert um die Babysit-Formulierung fuer Subagents. WICHTIG: den Block nie von Hand editieren - er wird von 'jaira update' regeneriert; die Aenderung gehoert in den Generator des Blocks im Go-Code. Bekannter Trap: nach der Aenderung ist der Block auf JEDEM Board stale, bis dort einmal 'jaira update' laeuft. Ergaenzend existiert seit 88H1P4 der Stop-Hook ('jaira hook print'), der das Sitzungsende bei wartenden agentischen Lanes verweigert - Block-Satz (Anweisung) und Hook (Durchsetzung) decken zusammen beide Haelften."
definition-of-done: "Der verwaltete Block enthaelt die Regel (Start -> durchfahren bis Human-Lane inkl. Loops, nach Klaerung weiter, 'Agent soll' -> Subagent babysittet); 'jaira update' schreibt sie; ein Test pinnt den neuen Blocktext"
tags: []
blocked-by: []
commits: []
created-at: 2026-09-03T12:49:02Z
updated-at: 2026-09-03T15:55:46Z
claimed-by: EE-3NX6GL3-2626905
claimed-at: 2026-09-03T15:47:55Z
updated-by: BeMuCa
---

# Der verwaltete Block schickt ein gestartetes Ticket bis zum naechsten Human-Step

## Definition of Done

- [x] Der verwaltete Block enthaelt die Regel (Start -> durchfahren bis Human-Lane inkl. Loops, nach Klaerung weiter, 'Agent soll' -> Subagent babysittet); 'jaira update' schreibt sie; ein Test pinnt den neuen Blocktext
  proof: core/board/announce.go:208-220 (laneSection tail, new sentences at 217-220); core/board/lanesection_test.go TestLaneSectionSaysHowToDriveAStartedTicket; regeneration verified on scratch board /tmp/claude-1000/-home-berk-git-jAIra/390b21fd-53c2-4182-9086-cad5a083b878/scratchpad/76wccw-verify/CLAUDE.md:107-109

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-03 15:55 · BeMuCa** — Placed the new rule in laneSection's tail (core/board/announce.go), not in the static agentNote block at the top. Reason: this rule is about driving lane-by-lane, which is exactly what the preceding 'Work one lane to empty...' sentence already covers, and that sentence lives in laneSection (only rendered when the board has lanes) — putting the new rule in agentNote would separate two sentences that belong together and would also show even on a board with zero lanes, where there is nothing to drive. Wording chosen: 'Told to start or work a ticket, drive it this way yourself — lane by lane, loops included — until it sits in a human lane, then continue once the human has answered. Told an agent should work it, hand it to a subagent that babysits the ticket through the same route.' Deliberately did not name critique/optimize specifically (unlike the German ticket text which says 'Loops (critique/optimize)') because laneSection is generic across boards — a board without those lane names still has loops via RejectsTo, and the existing Loop: sentences above already name the actual lanes for this board. Ruled out: touching agentNote's static text (wrong section, see above); touching internal/tui (out of scope per task constraints, and its test only checks 'This board's lanes'/'Order: ' substrings, not this tail, so unaffected). Verified end-to-end by building a scratch jaira binary into the scratchpad and running 'jaira init' in a fresh temp git repo outside this repo — confirmed the two new sentences render correctly in the generated CLAUDE.md. Did NOT run 'jaira update' in this repo per instructions; this repo's own CLAUDE.md/AGENTS.md remain stale until the user accepts and runs it themselves.
