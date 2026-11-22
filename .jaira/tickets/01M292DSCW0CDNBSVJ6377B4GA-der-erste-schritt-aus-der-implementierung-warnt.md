---
id: 01M292DSCW0CDNBSVJ6377B4GA
title: "Der erste Schritt aus der Implementierung warnt, wenn kein DoD-Punkt abgehakt ist"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer eine Implementierung verlaesst, ohne einen einzigen Punkt der Definition of Done abgehakt zu haben, erfaehrt das sofort - und nicht erst am Schluss, wenn die terminale Lane den Zug verweigert."
context: |-
  Ein Agent hat ein Ticket ueber acht Lanes und drei Review-Runden gefahren und dabei nur die Plan-Checkliste abgehakt, nie die Definition of Done. Aufgefallen ist es erst beim Sprung nach signoff: neun Verweigerungen auf einmal, ganz am Ende.

  - Es gibt zwei Checklisten im Ticket-Body: '## Definition of Done' (die Abnahmekriterien, die das Gate liest) und '## Plan' (die Methode). 'jaira dod <id> <n> --done' meint die erste, mit --plan die zweite.
  - Die Lane-Anweisung sagt beides bereits: .jaira/lanes/in-progress.md, Zeilen 17-23, samt --proof. Der Agent hat sie gelesen und trotzdem nur --plan abgehakt. Text allein hat also nicht gereicht - die beiden Listen sehen gleich aus, und am Ende der Arbeit fuehlt es sich an, als waere schon alles markiert.
  - Der generierte Block in CLAUDE.md/AGENTS.md ist inzwischen praeziser (core/board/announce.go), aber das ist dieselbe Gattung Gegenmittel: Prosa, die vorher gelesen werden muss.
  - Das Gate selbst ist richtig und soll nicht weicher werden. Falsch ist nur der Zeitpunkt, zu dem man es erfaehrt.
  - Gemeint ist eine Warnung, kein Refus: der Zug geht durch, aber die Zeile steht da. Etwa 'Definition of Done: 0 von 9 abgehakt - die terminale Lane liest diese Liste, nicht den Plan'.
  - Wo genau sie erscheint, ist noch nicht entschieden: beim Verlassen der implementierenden Lane, oder bei jedem Zug eines Tickets, das schon Code hat. Das gehoert in den Entwurf.
definition-of-done: "Ein Zug aus der implementierenden Lane heraus schreibt eine Warnung, wenn kein Punkt der Definition of Done abgehakt ist; der Zug selbst geht durch; in --json steht sie nicht auf stdout zwischen den Daten; ein Test deckt beide Faelle ab (kein Punkt abgehakt: Warnung; mindestens einer: keine)"
tags:
  - gates
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-09-11T20:26:57Z
updated-at: 2026-11-22T18:53:04Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-926893
claimed-at: 2026-11-22T18:49:16Z
outcome-what: "A move out of the implementing lane (the lane producing outcome-what), forward and not into a parking lane, with no DoD item settled, prints a warning on stderr; the move lands. Detection in core/move (Result.DoDUnticked), print in internal/cli/flow.go."
outcome-why: An agent ticked only the plan through eight lanes and met nine refusals at signoff; the news belongs at the step that should tick the items.
outcome-resolves: "DoD 1: warning on 0 ticked, none with one ticked, stderr also under --json, move lands; covered by core/move and internal/cli tests."
---

# Der erste Schritt aus der Implementierung warnt, wenn kein DoD-Punkt abgehakt ist

## Definition of Done

- [x] Ein Zug aus der implementierenden Lane heraus schreibt eine Warnung, wenn kein Punkt der Definition of Done abgehakt ist; der Zug selbst geht durch; in --json steht sie nicht auf stdout zwischen den Daten; ein Test deckt beide Faelle ab (kein Punkt abgehakt: Warnung; mindestens einer: keine)
  proof: core/move/move.go untickedDoD; internal/cli/flow.go warning on stderr; tests TestMoveReportsUntickedDoDOutOfImplementing, TestMoveOutOfImplementingWarnsOnUntickedDoD

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 18:49 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 18:50 · Alexander Sacharov** — decision: no open design question. DoD fixes the place: a move OUT of the implementing lane, i.e. the lane whose output-produces names outcome-what (in-progress on the shipped board), forward only (later lane index) and not into a parking lane (requires-blocked-reason, e.g. blocked) and not back. Detected in core/move (one place for every caller, field Result.DoDUnticked), printed by the CLI on stderr in both text and --json mode. 'Ticked' = settled (done or superseded): an all-superseded list passes the gate, so warning there would be false. A ticket with no checklist items gets no warning - nothing to tick, the gate says that itself. TUI prints nothing: the warning is for agents, who use the CLI.
