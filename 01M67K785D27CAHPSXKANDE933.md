---
id: 01M67K785D27CAHPSXKANDE933
title: Mit z (leere Lanes schmal) springen links/rechts zwischen Tickets ueber leere Lanes hinweg
status: testing
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
updated-at: 2026-11-04T00:51:40Z
updated-by: berkc
claimed-by: berk-ca-516678
claimed-at: 2026-11-04T00:38:18Z
outcome-what: "z is on by default: New() starts thinEmpty true, so empty lanes draw thin and ←/→ skip them from the first frame; one z turns it off for the session. Release note is a single '- !' line (core/release/notes/NDE933.md:1)"
outcome-why: "maintainer asked for z on by default (04.11.2026); z is session-only and not persisted (model.go:292), so no restart behaviour changes"
outcome-resolves: "critique pass 2: stale release-note line 1 ('Press z to draw empty lanes thin') contradicted the default and is deleted; DoD 2 proof re-checked at NDE933.md:1; tmux check still left to the parent"
review-summary: "none"
review-gaps: "removed TestWideBoardShowsAllTenLanesWidened (lanewindow_test.go): with z on by default it pressed z once and checked all ten lane names, the same setup and a subset of TestZWidensEverything; DoD proof line refs updated. Left: stepLane's wrap loop is not a duplicate of moveLane (steps i lanes, not one); TestZWidensEverything's thin check overlaps TestThinEmptyLanesAreOnByDefault's second half by one assertion, kept. No dead code or fluff in model.go/view.go."
test-verdict: "pass: go build ./... RC=0, go test ./... RC=0 (3 known hangers skipped per agent rules), DoD 1-3 proofs checked in the tree, tmux on scratch board: thin on first start, → backlog->done, z widens and → steps onto empty lane, restart thin again"
review-verdict: "pass - diff meets both DoD items (skip with z, unchanged without z, tests, release note); go vet and go test ./internal/tui/ pass; no defects found"
review-round: 1
review-check: "1. cd /home/berkc/code/X_General/jaira-exp-b && go run ./cmd/jaira  2. on the board, put the cursor on a ticket in a lane with empty lanes to its right  3. press → : cursor moves to the adjacent (empty) lane, as before  4. press ← to go back, then press z : empty lanes are drawn thin  5. press → : cursor jumps straight to the next lane that has tickets, skipping the thin ones  6. press z again and → : it steps one lane at a time again. Automated: go test ./internal/tui/ -run 'Arrows|LeftWithToggle' -count=1 passes"
round: 3
---

# Mit z (leere Lanes schmal) springen links/rechts zwischen Tickets ueber leere Lanes hinweg

## Definition of Done

- [x] Mit z: links/rechts ueberspringt leere Lanes; ohne z unveraendert; Test
  proof: internal/tui/model.go:1776 stepLane; tests TestArrowsWithToggleOnSkipEmptyLanes, TestLeftWithToggleOnWrapsToLastFilledLane, TestArrowsWithToggleOnStayWhenOnlyLaneHasTickets, TestArrowsWithToggleOnAllEmptyStillMove, TestArrowsWithToggleOnLeaveEmptyStartLane, TestArrowsWithToggleOffStepOntoEmptyLanes, TestLaneFocusIgnoresTheToggle (internal/tui/lanewindow_test.go:187-315)
- [x] Release-Note-Zeile
  proof: core/release/notes/NDE933.md:1
- [x] z ist standardmaessig an: Board startet mit schmalen leeren Lanes, Pfeile ueberspringen sie, z schaltet aus; ist der z-Zustand heute gespeichert, bleibt ein ausgeschaltetes z ueber Neustarts aus; Test
  proof: internal/tui/model.go:353 New() sets thinEmpty: true; model.go:292 field comment (session-only, not persisted, so every start is thin); test TestThinEmptyLanesAreOnByDefault (internal/tui/lanewindow_test.go:19)

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
- **2026-11-04 00:46 · berkc** — tmux check (parent, binary built from 9237f70d, scratch board: alpha in backlog, omega in done, all lanes between empty): first start shows empty lanes thin, hint 'z widen empty'; → from backlog opens omega (done), ← back to alpha; z -> hint 'z thin empty', lanes wide, → lands on empty Brainstorm; quit + restart -> thin again ('z widen empty'), as decided (session-only).
- **2026-11-04 00:47 · berkc** — critique pass 3 (fresh agent): read only 355c97a7, the answer to critique pass 2's one finding (00:45 note). core/release/notes/NDE933.md now holds the single '- !' line; the stale 'Press z to draw empty lanes thin' line is gone, so nothing contradicts thin-by-default any more. The repair introduces nothing new. Nothing left; the rest stands on the 02:35 and 00:45 notes. review-summary=none.
- **2026-11-04 00:48 · berkc** — optimize pass 2: only finding was a test the default flip made redundant (TestWideBoardShowsAllTenLanesWidened == subset of TestZWidensEverything). Removed; go vet + go test ./internal/tui -count=1 exit 0. Production code (stepLane, New() default, help line) unchanged - nothing to cut.
