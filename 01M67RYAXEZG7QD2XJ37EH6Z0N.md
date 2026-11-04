---
id: 01M67RYAXEZG7QD2XJ37EH6Z0N
title: "In human review und signoff sieht man den Code-Diff wie bei einem Merge, mit Sprung zur Datei in VS Code"
status: critique
ready: true
creator: berkc
assignee: berkc
goal: "Steht ein Ticket in human oder signoff, zeigt die Ticketansicht seinen Diff gegen master Datei fuer Datei (wie die Merge-Ansicht von GitHub), und je Datei oeffnet ein Knopf die Datei in VS Code - so wird die Abnahme ein echter Code-Review."
context: |-
  Was fehlt: wer abnimmt, liest Outcome und Review-Text, aber nicht den Code. Den Diff muss man selbst in git oder auf GitHub suchen.
  Ausloeser: Maintainer 31.10.2026, Punkt 8, mit Screenshot (Image #7 - im Agenten-Kontext nicht sichtbar; wo genau die Ansicht hin soll, steht dort): 'Wenn wir im human review step sind, dann sollte man auch die Code changes sehen, wie bei git wenn man mergen will ... mit einem button kann man dann auch in vsc zu dem file wechseln.'
  Bekannt: der Lane-Payload kennt den Diff schon (input-requires: diff in .jaira/lanes/review.md). VS Code oeffnet eine Datei per 'code -g <pfad>:<zeile>'. Ort (Cockpit-Pane aus Q63EEV, TUI-Ticketansicht oder jaira web) im Brainstorm mit dem Maintainer klaeren.
definition-of-done: Brainstorm hat mit dem Maintainer den Ort festgelegt (Screenshot aus Punkt 8)
tags:
  - tui
blocked-by: []
related:
  - 01M5M7JTVR8G7DZ1MGHJQ63EEV
commits: []
created-at: 2026-10-31T01:24:47Z
updated-at: 2026-11-04T01:03:40Z
updated-by: berkc
claimed-by: berk-ca-511849
claimed-at: 2026-11-04T00:37:01Z
outcome-what: "A ticket in signoff or human shows its diff against master beside it in the TUI at 140+ columns: the ticket and the diff each sit in a box, the focused one with the accent border. tab switches focus, ]/[ jump between files, o opens the file at its line in VS Code and confirms inline in the diff title ('opened <path>:<line>', cleared on the next key) with no modal; below 140 columns the screen is unchanged. Scroll and focus survive the 2s board reload and reset only when another ticket is opened."
outcome-why: "The person accepting work judged the outcome and review text but could not see the code; the diff had to be found by hand in git or on GitHub. A confirmation modal after o took the next keys, so the sign-off screen was lost behind it."
outcome-resolves: "Sign-off and human screens show each changed file with its +/- lines beside the ticket, and o opens the file at that line without a modal. Critique round 1 closed: reload keeps scroll and focus (TestReloadKeepsDiffScrollAndFocus, TestOpeningAnotherTicketResetsDiffScroll); one branch-segment split shared by BranchFor and branchTicket; the split is built from splitPane like the follow-up split (TestPersonLaneShowsDiffBesideTheTicketWhenWide). Testing finding closed: a successful o no longer raises a modal (TestOKeyConfirmsInTheDiffTitleWithoutAModal). Tests: go test ./core/gitrepo ./internal/tui -count=1 exit 0; go test ./internal/cli -count=1 -run 'Create|Follows|Branch' exit 0 (earlier run, cli untouched by this fix)."
review-summary: none
round: 3
review-gaps: "folded renderDiffPane's own pad/fill-to-height loop into clipPane (view.go:1490, same code) and dropped the diffScroll/top double bookkeeping there; left: Base() beside gitref.RemoteHead (different lookup, critique let-stand), WorktreeOf/openInCode (no existing twin), parseDiff's second counting pass (runs once per diff load), the focused-title restyle next to the accent border (visual, not a cleanup)"
test-verdict: "fail: a successful 'o' opens a modal NOTE that must be dismissed with esc and swallows the next keys"
---

# In human review und signoff sieht man den Code-Diff wie bei einem Merge, mit Sprung zur Datei in VS Code

## Definition of Done

- [x] Brainstorm hat mit dem Maintainer den Ort festgelegt (Screenshot aus Punkt 8)
  proof: note 2026-11-03 17:18 (maintainer, screenshot): right half of TUI ticket detail view
- [x] Ein Ticket in signoff zeigt seinen Diff je Datei; Test
  proof: internal/tui/diffpane_test.go TestPersonLaneShowsDiffBesideTheTicketWhenWide, TestHumanLaneShowsDiffBesideTheTicketWhenWide, TestNarrowTerminalKeepsTheSingleColumn, TestParseDiffSplitsFilesAndNumbersLines; core/gitrepo/branchdiff_test.go TestBranchDiffIsAgainstMergeBaseAndSkipsJaira
- [x] Ein Knopf/Taste oeffnet die Datei an der Zeile in VS Code
  proof: internal/tui/diffpane_test.go TestOKeyOpensFileAtLine and TestOKeyConfirmsInTheDiffTitleWithoutAModal; internal/tui/model.go case "o", "]", "[" in modeDetail; internal/tui/diffpane.go openDiffRow
- [x] Release-Notiz unter core/release/notes/
  proof: core/release/notes/EH6Z0N.md:1

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/gitrepo/branchdiff.go: BranchFor(handle), Base(), BranchDiff(base,branch), WorktreeOf(branch) + tests in core/gitrepo/branchdiff_test.go
- [x] internal/tui/diffpane.go: parseDiff into per-file rows with new-side line numbers, ticketDiff (branch...master, fallback CommitsForTicket+Diff), memoised per open ticket
- [x] render: renderSignOff and renderDetail (human lane) split into ticket left (78 cols) and diff right when width>=140 and height>=20; else unchanged single column
- [x] keys in modeDetail: tab switches pane focus, arrows/pgup/pgdn scroll focused pane, ] / [ next/prev file, o opens top-visible file:line via 'code -g' (in the branch's worktree if one exists)
- [x] tests internal/tui/diffpane_test.go: parse, wide split shows file header and coloured +/- lines, narrow falls back, ]/[ jump, o calls stubbed launcher with path:line
- [ ] release note core/release/notes/EH6Z0N.md; go test ./core/gitrepo ./internal/tui; tmux check on scratch board

## Progress
- **2026-11-03 17:18 · berkc** — Maintainer 03.11.2026 (Screenshot): Ort = rechte, heute leere Haelfte der TUI-Ticketansicht (neben tags/problem/what/why/summary), dazu im Cockpit-Pane aus Q63EEV - beide zeigen denselben Diff. Brainstorm-Punkt 1 damit erledigt.
- **2026-11-04 00:37 · berkc** — Worker 2026-11-04: Cockpit pane part (Q63EEV) left out - that pane does not exist yet; this ticket delivers only the TUI detail view. Diff = ticket branch vs master, reusing the lane payload's diff code. Narrow terminal: single-column view as today.
- **2026-11-04 00:40 · berkc** — Plan decisions (worker 2026-11-04): (1) Diff source = 'git diff <base>...<branch>' (merge-base, what GitHub shows), base master else main; branch = first local then remote branch whose name has the ticket handle as a segment (same rule as branchTicket in internal/cli/tickets.go). Reason: CommitsForTicket (the lane diff source, core/gitrepo/derive.go) walks HEAD only, so the maintainer's TUI on master would see nothing for an unmerged signoff ticket. Fallback when no branch exists (merged and deleted): the lane derivation CommitsForTicket+Repo.Diff, parsed the same way. .jaira/ is excluded like WorktreeDiff excludes tickets. (2) Split only at width>=140 so the 78-col sign-off column keeps its width and the diff gets >=59 cols; below that the screen is exactly today's. (3) Keys: tab focus, ]/[ file, o = VS Code at the top visible diff line; 'o' and ']' '[' were unused in modeDetail. (4) human lane (requires-question) gets the same right pane via renderDetail, per maintainer 'human or signoff'. Cockpit pane (Q63EEV) out of scope: it does not exist yet.
- **2026-11-04 00:42 · berkc** — Step 1 done: core/gitrepo/branchdiff.go (BranchFor, Base, BranchDiff, WorktreeOf) + branchdiff_test.go pass. Base() returns the short name for origin refs (origin/master). BranchDiff pins --src-prefix/--dst-prefix so diff.noprefix cannot break parseDiff. Fixture needs 'git branch -M master' since fixtureRepo does not name the branch.
- **2026-11-04 00:46 · berkc** — Pause after build: diffpane.go written (parseDiff state machine skips 'commit ' blocks so the fallback's stat lines are not read as context), signoff.go/renderDetail/model keys/followup scrollFocused wired. Not yet tested.
- **2026-11-04 00:48 · berkc** — Tests green: go test ./core/gitrepo ./internal/tui -count=1 exit 0. Decision: renderDiffPane clamps the scroll to the last row, not the last full page, so ] can put a short file's header on top even when the whole diff fits. Dead end: a first draft clamped to len-body and made ] a no-op on short diffs. Not done here: plan item 6 (tmux check on a scratch board, release note already written).
- **2026-11-04 00:50 · berkc** — critique (pass 1): internal/tui/model.go:475 reload() resets diffScroll and diffFocus, and reload runs on every 2s tickMsg (model.go:1133): a scrolled diff jumps back to the top and tab-focus falls back to the ticket within 2 seconds — clear only diffFor in reload, and reset diffScroll/diffFocus inside loadDiff (diffpane.go) only when t.ID differs from the ticket the diff was last loaded for (keep that id in its own field, since diffFor is blanked)
core/gitrepo/branchdiff.go:19 BranchFor re-implements the branch-segment split that branchTicket already has inline (internal/cli/tickets.go:157, same alnum closure as notAlnum) while its comment claims "the same rule" — export one helper from gitrepo (e.g. BranchSegments(name) []string) and use it in both places
internal/tui/diffpane.go:316 withDiff hand-rolls a second two-pane layout with its own focus signal (separator colour + title style) beside renderSplit/splitPane (internal/tui/followup.go:205-260), where the live pane is shown by an accent border — build the split with lipgloss.JoinHorizontal and splitPane so both split screens show focus the same way (shrink diffLeftWidth content/rows by the box frame)
Let stand: fallback via CommitsForTicket+Repo.Diff (needed on boards that work on master without branches; note 00:40); Base() hardcoding master/main; re-running git on each reload (same trade-off as derivedFor memo); WorktreeOf for 'o'.
- **2026-11-04 00:56 · berkc** — Round 2 (worker 2026-11-04): the three critique fixes are in. (1) reload() now clears only diffFor and keeps scroll and focus; loadDiff resets them only when the ticket differs from the new field diffScrollFor. Tests TestReloadKeepsDiffScrollAndFocus and TestOpeningAnotherTicketResetsDiffScroll. (2) gitrepo.BranchSegments is the one segment split; BranchFor and branchTicket both use it, branchTicket behaviour unchanged. (3) withDiff is splitPane boxes joined by lipgloss.JoinHorizontal, rows = height - footer - 2; the hand-built separator is gone. Side effect: TestAMultilineCheckKeepsItsLines (internal/tui/reviewcheck_test.go) matched a newline plus indent, which no longer exists once rows carry box borders at 140 columns; it now asserts step 2 and step 3 never share a row. Dead end: the goal line sits on an empty diff row, so the old check passed only on the border glyph; the new check asserts the handle row also holds the diff title. Open: plan item 6 (tmux check on a scratch board) is still not done.
- **2026-11-04 00:57 · berkc** — critique (pass 2): read only the pass-1 findings and commit 50765410. (1) reload keeps scroll/focus: model.go reload clears only diffFor, diffpane.go loadDiff resets via diffScrollFor - addressed. (2) gitrepo.BranchSegments used by BranchFor and internal/cli/tickets.go branchTicket, notAlnum gone - addressed. (3) withDiff = splitPane boxes + lipgloss.JoinHorizontal, widths 78+2 / 1 / width-83+2 add up to m.width - addressed. Nothing new; pass-1 let-stands still hold. review-summary=none.
- **2026-11-04 00:59 · berkc** — optimize: renderDiffPane hand-padded and height-filled its rows - the same loop clipPane (view.go:1490) runs; now ends in clipPane. Checked and left: Base() vs gitref.RemoteHead (different question: origin default only, Base prefers local master/main - critique let-stand); WorktreeOf has no twin in core (no other 'worktree list' parse); openInCode is the only 'code -g' launcher; parseDiff second pass for +/- counts is cheap and runs once per load; diffFocus title restyle beside the accent border kept (removing it changes the screen).
- **2026-11-04 01:00 · berkc** — testing (worker, tmux on scratch board, build of 50765410): diff split, ] [ tab scroll and narrow fallback (120 cols -> single column) all work; o ran 'code -g a.go:2' at the right line. FINDING: on success o calls m.notify(..., false) (internal/tui/model.go, case "o"), and notify always switches to modeMessage - a modal 'NOTE opened new.go:1 in VS Code / esc dismiss' pops over the screen, and the next keys ([, o) went to the modal and were lost; behind the modal the generic renderDetail screen is drawn instead of the sign-off screen (pre-existing modal behaviour, not this ticket's). FIX: on success do not call notify; show it inline like 'y' shows 'copied' (m.copied) - e.g. a field m.diffOpened string rendered in the diff pane title as 'opened new.go:1', cleared on the next key like m.copied is cleared at the top of case modeDetail. Errors keep notify. Add a test: after 'o' with the stub, m.mode is still modeDetail and the render contains 'opened <path>:<line>'.
- **2026-11-04 01:02 · berkc** — testing finding fixed (worker 2026-11-04): a successful 'o' no longer calls notify. openDiffRow sets Model.diffOpened (model.go, next to diffScrollFor); renderDiffPane prints 'Diff  opened <path>:<line>  <note>' so the opened text survives a narrow title; the key clears it at the top of case modeDetail, and loadDiff clears it when another ticket is opened. Errors still use notify. Test TestOKeyConfirmsInTheDiffTitleWithoutAModal: mode stays modeDetail, title shows the open, next key clears it. Dead end: a first struct edit with its own comment realigned the whole Model field block; kept the field inside the existing aligned block instead. Still open: plan item 6 tmux check on a scratch board, not run in this lane.
- **2026-11-04 01:03 · berkc** — critique (pass 3): read only the testing finding (note 01:00) and its repair 46448074. openDiffRow sets diffOpened instead of notify, cleared at top of case modeDetail beside m.copied and in loadDiff on a ticket change; errors still notify; test TestOKeyConfirmsInTheDiffTitleWithoutAModal - addressed, follows the existing copied pattern (model.go case modeDetail). Pass-1 let-stands and optimize keeps still hold. review-summary=none.
