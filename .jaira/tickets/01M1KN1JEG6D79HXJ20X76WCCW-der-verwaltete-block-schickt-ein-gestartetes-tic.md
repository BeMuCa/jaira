---
id: 01M1KN1JEG6D79HXJ20X76WCCW
title: Der verwaltete Block schickt ein gestartetes Ticket bis zum naechsten Human-Step
status: signoff
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Sagt der User 'starte/arbeite Ticket X', faehrt die Session es Lane fuer Lane durch - Loops (critique/optimize) inklusive - bis es in einer Human-Lane liegt, und nach der Klaerung weiter; sagt der User 'ein Agent soll es bearbeiten', babysittet ein Subagent es genauso durch"
context: "Berk am 03.09.: er promptet das heute jedes Mal von Hand ('babysitte das Ticket mit einem Subagent durch bis zum naechsten human step'). Die Regel soll in den verwalteten jaira-Block der CLAUDE.md/AGENTS.md, damit sie ohne erneutes Tippen gilt. Das ist Mechanismus c aus 88H1P4 (Block-Satz), der explizit auf Berks Go wartete - diese Anfrage ist das Go, erweitert um die Babysit-Formulierung fuer Subagents. WICHTIG: den Block nie von Hand editieren - er wird von 'jaira update' regeneriert; die Aenderung gehoert in den Generator des Blocks im Go-Code. Bekannter Trap: nach der Aenderung ist der Block auf JEDEM Board stale, bis dort einmal 'jaira update' laeuft. Ergaenzend existiert seit 88H1P4 der Stop-Hook ('jaira hook print'), der das Sitzungsende bei wartenden agentischen Lanes verweigert - Block-Satz (Anweisung) und Hook (Durchsetzung) decken zusammen beide Haelften."
definition-of-done: "Der verwaltete Block enthaelt die Regel (Start -> durchfahren bis Human-Lane inkl. Loops, nach Klaerung weiter, 'Agent soll' -> Subagent babysittet); 'jaira update' schreibt sie; ein Test pinnt den neuen Blocktext"
tags: []
blocked-by: []
commits:
  - ddc6c7a0e611d4ba84f5101ad3260cbffbf450d1
created-at: 2026-09-03T12:49:02Z
updated-at: 2026-09-03T15:59:04Z
claimed-by: EE-3NX6GL3-2626905
claimed-at: 2026-09-03T15:47:55Z
updated-by: BeMuCa
outcome-what: "Added two sentences to laneSection's tail in core/board/announce.go (the generator of the managed CLAUDE.md/AGENTS.md block), placed right after the existing 'Work one lane to empty...' sentence: told to start/work a ticket, drive it lane by lane (loops included) until it sits in a human lane, then continue once the human has answered; told an agent should work it, hand it to a subagent that babysits the same route. Extended core/board/lanesection_test.go with TestLaneSectionSaysHowToDriveAStartedTicket pinning the new sentences."
outcome-why: "Berk was typing this instruction by hand every time he started a ticket. This is mechanism c of 88H1P4 (the block-sentence half), which was waiting on his go-ahead — today's request. The stop-hook (jaira hook print) already enforces the agentic half; this closes the instruction half so the block itself tells a session/subagent to keep driving without being re-prompted."
outcome-resolves: "DoD required: (1) the managed block contains the rule — true, core/board/announce.go:217-220 adds it to laneSection's static tail, which every board with lanes renders; (2) 'jaira update' writes it — verified by building a scratch jaira binary into the scratchpad, running 'git init' + 'jaira init' in a fresh temp repo outside this one, and confirming the two new sentences appear verbatim in the generated CLAUDE.md (lines 107-109 there); did not run 'jaira update' in this repo itself, per instructions, so this repo's own CLAUDE.md/AGENTS.md stay stale until the user accepts and regenerates; (3) a test pins the new block text — TestLaneSectionSaysHowToDriveAStartedTicket in core/board/lanesection_test.go, passing. All quality gates green: gofmt -l core internal cmd printed exactly internal/cli/tickets.go (pre-existing, untouched by this change); go build ./... RC=0; go clean -testcache && go test ./... -race RC=0 (all packages ok, internal/tui and internal/cli included, confirming no regression in the areas outside scope)."
review-summary: none
review-gaps: "none — the diff is two appended sentences to laneSection's tail string plus one new test; nothing duplicated, nothing orphaned, no fluff, no hot path touched"
review-check: "1. grep -n \"Told to start or work a ticket\" core/board/announce.go\n   See two sentences right after \"the lane nobody drives is the one that fills up.\" — the babysit rule.\n2. go test ./core/board/... -run TestLaneSectionSaysHowToDriveAStartedTicket -v\n   See PASS.\n3. gofmt -l core internal cmd\n   See exactly one line: internal/cli/tickets.go (pre-existing drift, unrelated to this change).\n4. go build ./... && echo BUILD_OK\n   See BUILD_OK with no errors.\n5. rm -rf /tmp/jaira-check && mkdir /tmp/jaira-check && cd /tmp/jaira-check && git init -q && git config user.email a@b.c && git config user.name a\n   Sets up a throwaway repo outside jAIra so this repo's own CLAUDE.md/AGENTS.md are never touched.\n6. go build -o /tmp/jaira-check-bin -C /home/berk/git/jAIra ./cmd/jaira && /tmp/jaira-check-bin init\n   Creates AGENTS.md and CLAUDE.md in /tmp/jaira-check.\n7. grep -n \"Told to start or work a ticket\" /tmp/jaira-check/CLAUDE.md\n   See the two new sentences present in the freshly generated file — confirms 'jaira update'/'jaira init' actually writes the rule, not just the source string.\n8. rm -rf /tmp/jaira-check /tmp/jaira-check-bin\n   Cleans up the throwaway check."
review-verdict: "accept (Zweitmodell-Review: Implementierung Sonnet-Subagent, Review Fable-Koordinator am selbst gelesenen Diff ddc6c7a). Platzierung richtig: laneSection-Tail, direkt hinter 'work one lane to empty' - gleiche Frage, gleiche Stelle; der statische agentNote-Block waere vom Board-Kontext getrennt gewesen. Wortlaut in Hausstimme, Test pinnt die Saetze verbatim, Scratch-Board-Regeneration vom Implementierer nachgewiesen (CLAUDE.md Zeilen 107-109). Gates im Subagent-Lauf: gofmt nur tickets.go, build ok, ./... -race RC=0. Offen fuer dich: dieses Board bekommt die Saetze erst mit 'jaira update' (Block hier bewusst nicht angefasst)."
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
