---
id: 01M67K785D27CAHPSXKANDE933
title: Mit z (leere Lanes schmal) springen links/rechts zwischen Tickets ueber leere Lanes hinweg
status: pre-process
ready: true
creator: berkc
assignee: berkc
goal: "Ist z aktiv, landet ein Sprung nach links oder rechts im naechsten Ticket statt in jeder leeren Lane dazwischen."
context: "Wunsch von berkc am 31.10.2026 (Punkt 11): 'If we press z to thin empty lanes, i want to also skip the empty lanes when jumping from ticket to ticket. E.g. i am on a ticket in backlog, and there is only one other at done ... since all lanes inbetween are empty i want to just jump over them.'"
definition-of-done: "Mit z: links/rechts ueberspringt leere Lanes; ohne z unveraendert; Test"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-30T23:44:47Z
updated-at: 2026-10-31T02:32:02Z
updated-by: berkc
claimed-by: berk-ca-3155695
claimed-at: 2026-10-31T02:31:33Z
---

# Mit z (leere Lanes schmal) springen links/rechts zwischen Tickets ueber leere Lanes hinweg

## Definition of Done

- [ ] Mit z: links/rechts ueberspringt leere Lanes; ohne z unveraendert; Test
- [ ] Release-Note-Zeile

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Add stepLane in internal/tui/model.go after moveLane; board arrows call it, lane focus keeps moveLane
- [ ] Rewrite the z test and add arrow tests in internal/tui/lanewindow_test.go (skip, off, wrap, only-lane, all-empty, empty start)
- [ ] Update toggleEmptyLanes doc comment and the z help line in internal/tui/view.go
- [ ] Release note line in core/release/notes/NDE933.md
- [ ] Verify: go vet and go test ./internal/tui/

## Progress

