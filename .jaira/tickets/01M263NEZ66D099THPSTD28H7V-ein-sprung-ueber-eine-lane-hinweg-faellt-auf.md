---
id: 01M263NEZ66D099THPSTD28H7V
title: Ein Sprung ueber eine Lane hinweg faellt auf
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Wer ein Ticket an einer installierten Lane vorbeibewegt, wird darauf hingewiesen; am Ende der Kette ist nachvollziehbar, welche Lanes es nie gesehen hat"
context: |-
  Gemessen am 09.09. auf dem requirementsgenie-Board: 'jaira move JP1PA1 --to review' direkt aus in-progress wird ERLAUBT, EXIT=0, obwohl critique und optimize dazwischen liegen.
  Ursache: core/gate/gate.go:446-455 prueft nur die Lane, die verlassen wird, und die, die betreten wird. Was dazwischen liegt, wird nie gefragt.
  Folge, in den Daten belegbar: von 11 Tickets in review tragen 8 kein review-verdict, und CANFDQ, CNTYMW, ZTR2CW haben woertlich 'none' als review-summary. Die Felder gehoeren critique und optimize, die uebersprungen wurden.
  Ein Agent hat es selbst bemerkt und in review-summary von 9GEGTB hineingeschrieben: 'CRITIQUE (lane run retroactively on 2026-09-09 -- critique and optimize were skipped on the way to review)' - und dann Critique- UND Optimize-Ergebnis in EIN Feld gequetscht, weil kein anderes uebrig war.
  Auch die input-requires der uebersprungenen Lanes verfallen: optimize verlangt 'diff', das nie eingefordert wurde - deshalb tragen 9 der 11 Review-Tickets null Commits.
  Nicht gewuenscht ist ein hartes Verbot: ein Ticket darf legitim an einer Lane vorbei (Parken in blocked, ein Hotfix). Gewuenscht ist, dass es nicht LAUTLOS passiert.
definition-of-done: "Ein Move, der eine installierte Lane ueberspringt, meldet welche und was dadurch ungefragt bleibt; die Information ist am Ticket ablesbar, nicht nur in der Terminalausgabe; ein Test deckt Sprung und Meldung ab; go test ./... -race gruen"
tags:
  - gates
blocked-by: []
commits: []
created-at: 2026-09-10T16:50:54Z
updated-at: 2026-10-07T06:07:42Z
updated-by: Alexander Sacharov
---

# Ein Sprung ueber eine Lane hinweg faellt auf

## Definition of Done

- [ ] Ein Move, der eine installierte Lane ueberspringt, meldet welche und was dadurch ungefragt bleibt; die Information ist am Ticket ablesbar, nicht nur in der Terminalausgabe; ein Test deckt Sprung und Meldung ab; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-07 06:07 · Alexander Sacharov** — Hinweis aus BG5QJ6 (2026-10-07): eine Lane mit requires-question (human) liegt nicht auf dem Weg — core/lane/next.go überspringt sie. Die Sprung-Erkennung darf testing->review nicht als übersprungene Lane melden.
