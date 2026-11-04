---
id: 01M67K785D27CAHPSXKANDE933
title: Mit z (leere Lanes schmal) springen links/rechts zwischen Tickets ueber leere Lanes hinweg
status: in-progress
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
updated-at: 2026-11-04T00:46:06Z
updated-by: berkc
claimed-by: berk-ca-516678
claimed-at: 2026-11-04T00:38:18Z
outcome-what: "z is on by default: New() starts thinEmpty true, so empty lanes draw thin and ←/→ skip them from the first frame; one z turns it off for the session; ? help line says on by default; tests adjusted (z pressed where a test wants full-width lanes, dropped where it wanted z on) plus TestThinEmptyLanesAreOnByDefault; release note line 2 is a '- !' instruction"
outcome-why: "maintainer asked for z on by default; z is session-only and not persisted (model.go:292), so no restart behaviour changes"
outcome-resolves: "DoD 3 ticked with proof (model.go:353, TestThinEmptyLanesAreOnByDefault); DoD 1 proof re-pointed to model.go:1776 and lanewindow_test.go:201-329; tmux check left to the parent"
review-summary: "core/release/notes/NDE933.md:1 still says \"Press z to draw empty lanes thin\", which is false since 9237f70d made thin the default (z now widens); jaira update shows both lines side by side and they contradict each other. The file is unreleased (directly in core/release/notes/), so delete line 1 and keep only line 2, or rewrite line 1 so it matches line 2. Everything else in pass 2 (New() thinEmpty: true beside anim, field comment, ? help line, zHint tests, no persistence per the 00:37 plan note) needs no change."
review-gaps: |-
  note: skipping wraps round the board end (backlog ← lands on done), same wrap moveLane already had, but a user might expect a stop at the edge
  note: lane focus mode still steps onto empty lanes with z on - deliberate per the doc comment, not covered by a test
test-verdict: "pass: go build ./... RC=0, go vet ./internal/tui RC=0, go test ./internal/tui -count=1 ok (RC=0); DoD 1-2 verified in tree, 7 named tests PASS -v; release note is one '- ' line; behaviour exercised through Model.key(KeyRight/KeyLeft) in tests, not a live TUI run"
review-verdict: "pass - diff meets both DoD items (skip with z, unchanged without z, tests, release note); go vet and go test ./internal/tui/ pass; no defects found"
review-round: 1
review-check: "1. cd /home/berkc/code/X_General/jaira-exp-b && go run ./cmd/jaira  2. on the board, put the cursor on a ticket in a lane with empty lanes to its right  3. press → : cursor moves to the adjacent (empty) lane, as before  4. press ← to go back, then press z : empty lanes are drawn thin  5. press → : cursor jumps straight to the next lane that has tickets, skipping the thin ones  6. press z again and → : it steps one lane at a time again. Automated: go test ./internal/tui/ -run 'Arrows|LeftWithToggle' -count=1 passes"
round: 3
---

# Mit z (leere Lanes schmal) springen links/rechts zwischen Tickets ueber leere Lanes hinweg

## Definition of Done

- [x] Mit z: links/rechts ueberspringt leere Lanes; ohne z unveraendert; Test
  proof: internal/tui/model.go:1776 stepLane; tests TestArrowsWithToggleOnSkipEmptyLanes, TestLeftWithToggleOnWrapsToLastFilledLane, TestArrowsWithToggleOnStayWhenOnlyLaneHasTickets, TestArrowsWithToggleOnAllEmptyStillMove, TestArrowsWithToggleOnLeaveEmptyStartLane, TestArrowsWithToggleOffStepOntoEmptyLanes, TestLaneFocusIgnoresTheToggle (internal/tui/lanewindow_test.go:201-329)
- [x] Release-Note-Zeile
  proof: core/release/notes/NDE933.md:1
- [x] z ist standardmaessig an: Board startet mit schmalen leeren Lanes, Pfeile ueberspringen sie, z schaltet aus; ist der z-Zustand heute gespeichert, bleibt ein ausgeschaltetes z ueber Neustarts aus; Test
  proof: internal/tui/model.go:353 New() sets thinEmpty: true; model.go:292 field comment (session-only, not persisted, so every start is thin); test TestThinEmptyLanesAreOnByDefault (internal/tui/lanewindow_test.go:33)

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
- [x] New() in internal/tui/model.go sets thinEmpty: true beside anim: true; update the thinEmpty field comment (on by default, z turns it off, session-only)
- [x] Fix tests that assumed z starts off (internal/tui/*_test.go): tests that press z to turn it on now get it by default; tests of the off behaviour press z once; add a test that a fresh New() model has thinEmpty on and draws an empty lane thin, and that one z turns it off
- [x] Check the ? help text and hint bar still read right with the default on (view.go zHint)
- [x] Add a second line '- ! ...' to core/release/notes/NDE933.md: empty lanes now start thin and skipped, press z to show them
- [x] Verify: go vet ./internal/tui/ and go test ./internal/tui/ -count=1; TUI via tmux on a scratch board

## Progress
- **2026-10-31 02:33 · berkc** — Worktree already held 31871f63 (experiment); this pass: stepLane doc names the only-lane case, one comment on the all-empty fallback, release note rewritten as an instruction. Plan items were recorded from plan.md and ticked. Board reports an older jaira version (jaira update not run).
- **2026-10-31 02:35 · berkc** — critique (first pass, fresh agent): no findings. stepLane separate from moveLane is justified (lane focus must keep stepping onto every lane); skip predicate len(tickets)==0 matches the thin-draw predicate in view.go:209; wraparound matches moveLane; help text in view.go updated. Nothing speculative.
- **2026-10-31 02:36 · berkc** — optimize: looked for duplicate non-empty-lane scan (none), orphans (none), fluff (none); go test ./internal/tui green
- **2026-10-31 02:38 · berkc** — testing: whole-suite gate in lane prompt not run; per parent instruction and agent rules only internal/tui was tested (build ./... also run, RC=0). proof line says model.go:1774, stepLane is at 1775 (doc comment 1771) - cosmetic.
- **2026-11-04 00:37 · berkc** — Maintainer request 04.11.2026 (before accepting): z must be ON by default - the board starts with empty lanes collapsed/thin and left/right jumps skip them; pressing z turns it off. If the z state is persisted anywhere today (settings, per board), keep that working: a person who turned it off must not get it back on every start.
- **2026-11-04 00:37 · berkc** — Plan decision: z state is not persisted anywhere today - model.go:292-298 documents thinEmpty as session-only (filter, cursor, scroll all start fresh; only lane.SaveDefaultBoard is saved and that is a workflow definition). So nothing to keep working: default flips to on in New() exactly like anim (model.go:352), z turns it off for the session, next start is on again. Adding persistence would be a new write path the maintainer did not ask for; if they want z-off to survive restarts that is a follow-up ticket.
- **2026-11-04 00:38 · berkc** — in-progress pass 2: starting plan 6 (New() thinEmpty default on). Reading model.go New() and the thinEmpty field comment; plan 7 (tests) and 8 (zHint) follow.
- **2026-11-04 00:43 · berkc** — plan 6-8 done: New() thinEmpty default on; tests adjusted (press z where a test wants the old full-width board, drop it where it wanted z on); new TestThinEmptyLanesAreOnByDefault; hint label checks now 'z widen empty' on the default board; ? help line says on by default. go vet + go test ./internal/tui/ -count=1 green (exit 0).
- **2026-11-04 00:43 · berkc** — tmux check (plan 10 second half) is left to the parent: this lane ran go vet and go test ./internal/tui/ -count=1 only. Release note: appended line 2 of core/release/notes/NDE933.md as '- !'; line 1 still reads 'Press z to draw empty lanes thin', now only true as a press-to-widen-then-thin description - left unchanged per brief, parent may rewrite it.
- **2026-11-04 00:45 · berkc** — critique pass 2 (fresh agent): read only 9237f70d, the change since the first critique, which found nothing (02:35 note). One finding: core/release/notes/NDE933.md:1 still says "Press z to draw empty lanes thin", which is false since 9237f70d made thin the default (z now widens); jaira update shows both lines side by side and they contradict each other. The file is unreleased (directly in core/release/notes/), so delete line 1 and keep only line 2, or rewrite line 1 so it matches line 2. Everything else in pass 2 (New() thinEmpty: true beside anim, field comment, ? help line, zHint tests, no persistence per the 00:37 plan note) needs no change.
- **2026-11-04 00:45 · berkc** — in-progress pass 3 (critique pass 2 finding): core/release/notes/NDE933.md line 1 ('Press z to draw empty lanes thin') contradicted the new default. Deleted it; the single '- !' line now carries the change (thin and skipped from the start, z shows them full width, not saved). The file is unreleased, so no version folder is touched. DoD 2 proof core/release/notes/NDE933.md:1 now points at the '- !' line, unchanged.
