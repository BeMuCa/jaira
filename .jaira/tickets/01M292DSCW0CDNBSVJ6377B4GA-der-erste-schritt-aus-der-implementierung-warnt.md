---
id: 01M292DSCW0CDNBSVJ6377B4GA
title: "Der erste Schritt aus der Implementierung warnt, wenn kein DoD-Punkt abgehakt ist"
status: backlog
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
updated-at: 2026-09-11T20:27:08Z
updated-by: Alexander Sacharov
---

# Der erste Schritt aus der Implementierung warnt, wenn kein DoD-Punkt abgehakt ist

## Definition of Done

- [ ] Ein Zug aus der implementierenden Lane heraus schreibt eine Warnung, wenn kein Punkt der Definition of Done abgehakt ist; der Zug selbst geht durch; in --json steht sie nicht auf stdout zwischen den Daten; ein Test deckt beide Faelle ab (kein Punkt abgehakt: Warnung; mindestens einer: keine)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

