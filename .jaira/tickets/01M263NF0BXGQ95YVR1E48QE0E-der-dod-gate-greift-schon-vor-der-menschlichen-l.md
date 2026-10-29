---
id: 01M263NF0BXGQ95YVR1E48QE0E
title: "Der DoD-Gate greift schon vor der menschlichen Lane, nicht erst bei done"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Ein Ticket mit unabgehakter Definition of Done kommt nicht kommentarlos in eine Lane, aus der nur ein Mensch herausfuehrt"
context: |-
  Auf dem Board eines anderen Projekts liegen NE23SW, 19PN17 und 022SGY mit 'DoD 0/1' in review - der Lane, aus der nur ein Mensch herauskommt. Berk sieht sie dort als abnahmebereit.
  Der Gate EXISTIERT, er greift nur zu spaet. Nachgemessen auf einem Wegwerf-Board:
  - in-progress -> review mit DoD 0/1: EXIT=0, erlaubt.
  - review -> done mit DoD 0/1: refused, woertlich 'the definition of done is not met: criterion 1 ("d") is still open. Satisfy it, mark it with jaira dod ... --done, then try this move again'.
  Der Gate haengt also an done, nicht am menschlichen Checkpoint davor. Wer vor done abnimmt, nimmt ungeprueft ab.
  Zu entscheiden beim Bauen: verweigern oder warnen? Verweigern ist konsequent, aber es gibt legitime Faelle, in denen ein Mensch genau deshalb draufschauen soll, WEIL der DoD nicht aufgeht. Eine Warnung, die am Ticket sichtbar bleibt, ist vermutlich richtiger als ein Verbot.
definition-of-done: "Ein Move in eine Lane mit requires-human-exit meldet einen offenen DoD; die Meldung ist am Ticket sichtbar, nicht nur im Terminal; der bestehende harte Gate an done bleibt unveraendert; Tests decken beide Lanes ab; go test ./... -race gruen"
tags:
  - gates
blocked-by: []
commits: []
created-at: 2026-09-10T16:50:54Z
updated-at: 2026-09-10T16:50:54Z
---

# Der DoD-Gate greift schon vor der menschlichen Lane, nicht erst bei done

## Definition of Done

- [ ] Ein Move in eine Lane mit requires-human-exit meldet einen offenen DoD; die Meldung ist am Ticket sichtbar, nicht nur im Terminal; der bestehende harte Gate an done bleibt unveraendert; Tests decken beide Lanes ab; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

