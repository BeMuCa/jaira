---
id: 01M5631PZGXXE3QB21MPX3A8RZ
title: "Das Board laesst rechts einen leeren Streifen, weil jede Spalte zwei Zellen Abstand eingeplant bekommt, die nie gezeichnet werden"
status: in-progress
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
updated-at: 2026-10-17T23:30:28Z
updated-by: berkc
outcome-what: boardFit budgets columns at rendered width (no +2)
outcome-why: lipgloss v2 Width includes border; the +2 left blank strip on the right
outcome-resolves: board row fills terminal width within <full columns cells
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
