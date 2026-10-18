---
id: 01M5631PZGXXE3QB21MPX3A8RZ
title: "Das Board laesst rechts einen leeren Streifen, weil jede Spalte zwei Zellen Abstand eingeplant bekommt, die nie gezeichnet werden"
status: testing
ready: true
creator: berkc
assignee: berkc
goal: Die Lanes fuellen die ganze Terminalbreite; rechts bleibt hoechstens der Rest einer ganzzahligen Teilung frei.
context: |-
  Was heute ist: Berk am 18.10. mit Screenshot - rechts neben Blocked ein leerer Streifen von etwa 25 Zellen, egal wie breit das Fenster ist.
  Ursache (nachgemessen mit tmux, master 2ea98bfe): internal/tui/view.go boardFit budgetiert jede Spalte mit Breite + 2 (Zeilen 197, 200, 211: minColWidth + 2, thinColWidth + 2) und rechnet colW = (width-thinCells)/full - 2. lipgloss v2 Width ist aber die ganze Box inklusive Rahmen (steht selbst im Kommentar an thinColWidth), und JoinHorizontal setzt keinen Abstand. Jede Spalte ist also 2 Zellen schmaler gezeichnet als eingeplant.
  Messung: 200 Spalten breit, 6 volle + 7 duenne Lanes -> gezeichnet 172, 28 Zellen leer. Bei 130 Spalten 112 gezeichnet.
  Die +2 sind ein Rest aus lipgloss v1, wo der Rahmen zur Breite dazukam.
definition-of-done: "boardFit plant jede Spalte mit ihrer gezeichneten Breite ein; ein Test misst, dass die gezeichnete Board-Zeile bei mehreren Breiten und mit duennen Lanes die Terminalbreite bis auf weniger als die Zahl der vollen Spalten ausfuellt und nie ueberschreitet; go test ./... gruen; Release-Note"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-17T23:27:12Z
updated-at: 2026-10-18T00:00:01Z
updated-by: berkc
outcome-what: boardFit budgets columns at rendered width (no +2)
outcome-why: lipgloss v2 Width includes border; the +2 left blank strip on the right
outcome-resolves: board row fills terminal width within <full columns cells
review-summary: "boardFit (internal/tui/view.go:191-220) now budgets each lane at its rendered width (22 full, 4 thin) instead of width+2, and colW=(width-thinCells)/full with no -2. lipgloss v2 Width already includes the border and JoinHorizontal adds no gap, so the old +2 was a phantom margin that left a blank strip right of the last lane. colW comment updated; tagbox_test bound tightened to colW; new TestBoardRowFillsTerminalWidth checks the row never exceeds and fills to < full-column count, with and without thin lanes. Release note added."
review-gaps: "none in the change. Pre-existing, unrelated (also on master 2ea98bfe): a lane whose cards overflow draws its '+N more' line one row below the other columns' bottom border (seen at 79/60/30 wide). Not verified in tmux: inline card-link lines for a linked card (only the links popup); covered by linklines tests, and spots use measured widths, not colW."
test-verdict: "pass: vet clean; go test ./... ok; new test fails on old code; tmux 200 cols 184->198 (thin 174->199)"
review-verdict: "approve: diff matches the DoD; vet clean, go test -race ./... all ok; real binary in tmux at 200/177/131/101/79/60/41/40/30, plain, z thin, cursor walked to the last lane: no line ever wider than the terminal, column row now 198/200 (was 184 on master at 200; 72 vs 78 at 79), focused lane always on screen, links popup inside the frame."
review-round: 1
review-check: "1. cd /home/berkc/code/X_General/jaira-X3A8RZ && go test ./internal/tui -run TestBoardRowFillsTerminalWidth -v  -> PASS  2. go build -o /tmp/jx ./cmd/jaira  3. in a terminal ~200 columns wide, from the jaira repo: JAIRA_NO_FETCH=1 /tmp/jx, press Enter to open the board  4. the last lane's right border sits at (or 1-2 cells before) the right edge - no wide empty strip  5. press z: thin empty lanes, still flush right, nothing wraps onto a second line  6. press the right arrow until the cursor reaches Blocked: the focused lane stays visible  7. shrink the window to ~80 then ~40 columns: lanes drop out, none of the box lines break or wrap"
---

# Das Board laesst rechts einen leeren Streifen, weil jede Spalte zwei Zellen Abstand eingeplant bekommt, die nie gezeichnet werden

## Definition of Done

- [x] boardFit plant jede Spalte mit ihrer gezeichneten Breite ein; ein Test misst, dass die gezeichnete Board-Zeile bei mehreren Breiten und mit duennen Lanes die Terminalbreite bis auf weniger als die Zahl der vollen Spalten ausfuellt und nie ueberschreitet; go test ./... gruen; Release-Note
  proof: internal/tui/view.go boardFit; TestBoardRowFillsTerminalWidth

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-17 23:28 · berkc** — berkc took this ticket over from BeMuCa
- **2026-10-17 23:36 · berkc** — critique: diff is 5 constant edits in boardFit + test; thin/full cost now rendered width, card-link spots and tagbox bounds already used real widths (lipgloss.Width), tagbox_test bound updated to colW. Nothing simpler.
- **2026-10-18 00:00 · berkc** — review: approved. Pre-existing on master, not this ticket: '+N more' line renders one row below the other columns' bottom border.
