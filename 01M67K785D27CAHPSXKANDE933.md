---
id: 01M67K785D27CAHPSXKANDE933
title: Mit z (leere Lanes schmal) springen links/rechts zwischen Tickets ueber leere Lanes hinweg
status: critique
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
updated-at: 2026-10-31T02:35:14Z
updated-by: berkc
claimed-by: berk-ca-3155695
claimed-at: 2026-10-31T02:31:33Z
outcome-what: stepLane in model.go skips lanes with no tickets when z is on; release note rewritten as an instruction
outcome-why: "with z on, ←/→ from backlog to done jumped through every empty lane in between"
outcome-resolves: "DoD 1 and 2 ticked with proof (tests in lanewindow_test.go:182-301, release note line)"
review-summary: none
---

# Mit z (leere Lanes schmal) springen links/rechts zwischen Tickets ueber leere Lanes hinweg

## Definition of Done

- [x] Mit z: links/rechts ueberspringt leere Lanes; ohne z unveraendert; Test
  proof: internal/tui/model.go:1774 stepLane; tests TestArrowsWithToggleOnSkipEmptyLanes, TestLeftWithToggleOnWrapsToLastFilledLane, TestArrowsWithToggleOnStayWhenOnlyLaneHasTickets, TestArrowsWithToggleOnAllEmptyStillMove, TestArrowsWithToggleOnLeaveEmptyStartLane, TestArrowsWithToggleOffStepOntoEmptyLanes, TestLaneFocusIgnoresTheToggle (internal/tui/lanewindow_test.go:182-301)
- [x] Release-Note-Zeile
  proof: core/release/notes/NDE933.md:1

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Add stepLane in internal/tui/model.go after moveLane; board arrows call it, lane focus keeps moveLane
- [x] Rewrite the z test and add arrow tests in internal/tui/lanewindow_test.go (skip, off, wrap, only-lane, all-empty, empty start)
- [x] Update toggleEmptyLanes doc comment and the z help line in internal/tui/view.go
- [x] Release note line in core/release/notes/NDE933.md
- [x] Verify: go vet and go test ./internal/tui/

## Progress
- **2026-10-31 02:33 · berkc** — Worktree already held 31871f63 (experiment); this pass: stepLane doc names the only-lane case, one comment on the all-empty fallback, release note rewritten as an instruction. Plan items were recorded from plan.md and ticked. Board reports an older jaira version (jaira update not run).
- **2026-10-31 02:35 · berkc** — critique (first pass, fresh agent): no findings. stepLane separate from moveLane is justified (lane focus must keep stepping onto every lane); skip predicate len(tickets)==0 matches the thin-draw predicate in view.go:209; wraparound matches moveLane; help text in view.go updated. Nothing speculative.
