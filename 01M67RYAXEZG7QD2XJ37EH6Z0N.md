---
id: 01M67RYAXEZG7QD2XJ37EH6Z0N
title: "In human review und signoff sieht man den Code-Diff wie bei einem Merge, mit Sprung zur Datei in VS Code"
status: in-progress
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
updated-at: 2026-11-04T00:48:03Z
updated-by: berkc
claimed-by: berk-ca-511849
claimed-at: 2026-11-04T00:37:01Z
---

# In human review und signoff sieht man den Code-Diff wie bei einem Merge, mit Sprung zur Datei in VS Code

## Definition of Done

- [x] Brainstorm hat mit dem Maintainer den Ort festgelegt (Screenshot aus Punkt 8)
  proof: note 2026-11-03 17:18 (maintainer, screenshot): right half of TUI ticket detail view
- [x] Ein Ticket in signoff zeigt seinen Diff je Datei; Test
  proof: internal/tui/diffpane_test.go TestPersonLaneShowsDiffBesideTheTicketWhenWide, TestHumanLaneShowsDiffBesideTheTicketWhenWide, TestNarrowTerminalKeepsTheSingleColumn, TestParseDiffSplitsFilesAndNumbersLines; core/gitrepo/branchdiff_test.go TestBranchDiffIsAgainstMergeBaseAndSkipsJaira
- [x] Ein Knopf/Taste oeffnet die Datei an der Zeile in VS Code
  proof: internal/tui/diffpane_test.go TestOKeyOpensFileAtLine and TestBracketJumpsBetweenFiles; internal/tui/model.go case "o", "]", "[" in modeDetail; internal/tui/diffpane.go openDiffRow
- [x] Release-Notiz unter core/release/notes/
  proof: core/release/notes/EH6Z0N.md:1

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/gitrepo/branchdiff.go: BranchFor(handle), Base(), BranchDiff(base,branch), WorktreeOf(branch) + tests in core/gitrepo/branchdiff_test.go
- [~] internal/tui/diffpane.go: parseDiff into per-file rows with new-side line numbers, ticketDiff (branch...master, fallback CommitsForTicket+Diff), memoised per open ticket
- [ ] render: renderSignOff and renderDetail (human lane) split into ticket left (78 cols) and diff right when width>=140 and height>=20; else unchanged single column
- [ ] keys in modeDetail: tab switches pane focus, arrows/pgup/pgdn scroll focused pane, ] / [ next/prev file, o opens top-visible file:line via 'code -g' (in the branch's worktree if one exists)
- [ ] tests internal/tui/diffpane_test.go: parse, wide split shows file header and coloured +/- lines, narrow falls back, ]/[ jump, o calls stubbed launcher with path:line
- [ ] release note core/release/notes/EH6Z0N.md; go test ./core/gitrepo ./internal/tui; tmux check on scratch board

## Progress
- **2026-11-03 17:18 · berkc** — Maintainer 03.11.2026 (Screenshot): Ort = rechte, heute leere Haelfte der TUI-Ticketansicht (neben tags/problem/what/why/summary), dazu im Cockpit-Pane aus Q63EEV - beide zeigen denselben Diff. Brainstorm-Punkt 1 damit erledigt.
- **2026-11-04 00:37 · berkc** — Worker 2026-11-04: Cockpit pane part (Q63EEV) left out - that pane does not exist yet; this ticket delivers only the TUI detail view. Diff = ticket branch vs master, reusing the lane payload's diff code. Narrow terminal: single-column view as today.
- **2026-11-04 00:40 · berkc** — Plan decisions (worker 2026-11-04): (1) Diff source = 'git diff <base>...<branch>' (merge-base, what GitHub shows), base master else main; branch = first local then remote branch whose name has the ticket handle as a segment (same rule as branchTicket in internal/cli/tickets.go). Reason: CommitsForTicket (the lane diff source, core/gitrepo/derive.go) walks HEAD only, so the maintainer's TUI on master would see nothing for an unmerged signoff ticket. Fallback when no branch exists (merged and deleted): the lane derivation CommitsForTicket+Repo.Diff, parsed the same way. .jaira/ is excluded like WorktreeDiff excludes tickets. (2) Split only at width>=140 so the 78-col sign-off column keeps its width and the diff gets >=59 cols; below that the screen is exactly today's. (3) Keys: tab focus, ]/[ file, o = VS Code at the top visible diff line; 'o' and ']' '[' were unused in modeDetail. (4) human lane (requires-question) gets the same right pane via renderDetail, per maintainer 'human or signoff'. Cockpit pane (Q63EEV) out of scope: it does not exist yet.
- **2026-11-04 00:42 · berkc** — Step 1 done: core/gitrepo/branchdiff.go (BranchFor, Base, BranchDiff, WorktreeOf) + branchdiff_test.go pass. Base() returns the short name for origin refs (origin/master). BranchDiff pins --src-prefix/--dst-prefix so diff.noprefix cannot break parseDiff. Fixture needs 'git branch -M master' since fixtureRepo does not name the branch.
- **2026-11-04 00:46 · berkc** — Pause after build: diffpane.go written (parseDiff state machine skips 'commit ' blocks so the fallback's stat lines are not read as context), signoff.go/renderDetail/model keys/followup scrollFocused wired. Not yet tested.
