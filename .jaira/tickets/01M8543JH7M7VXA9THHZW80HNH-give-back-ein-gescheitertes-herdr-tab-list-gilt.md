---
id: 01M8543JH7M7VXA9THHZW80HNH
title: "Give back: ein gescheitertes 'herdr tab list' gilt als 'kein Tab offen' und startet einen zweiten Worker neben dem lebenden"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Kann das Cockpit die Tabs nicht lesen, startet Give back keinen neuen Worker, sondern sagt der Person, dass es den Zustand nicht pruefen konnte."
context: |-
  Was falsch ist: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1784 - scheitert 'herdr tab list' (Herdr antwortet nicht, Timeout), ist die Liste leer, Give back haelt den alten Worker fuer weg und startet einen neuen. Laeuft der alte doch, arbeiten zwei Worker am selben Ticket und Worktree.
  Gab es schon vor M31BYV (Commit 82aa5654a3); gefunden im Review zu M31BYV am 23.11.2026.
  Erst nach M31BYV starten: beide aendern denselben Give-back-Teil in register.tsx.
definition-of-done: "Scheitert 'herdr tab list', startet Give back keinen Worker und meldet der Person, dass der Tab-Zustand nicht lesbar war; Test"
tags:
  - cockpit
blocked-by: []
related: []
follows: 01M851Y62KHD9DVD2GDAM31BYV
commits: []
created-at: 2026-11-23T21:13:33Z
updated-at: 2026-11-24T07:08:56Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-46440
claimed-at: 2026-11-24T06:58:04Z
outcome-what: "Give back stops when 'herdr tab list' fails: no new worker, no keys into the old pane, a note and a toast saying the tab state was unreadable"
outcome-why: "A failed tab list read as 'no tab open', so Give back started a second worker beside a live one on the same ticket and worktree"
outcome-resolves: Two workers on one ticket after a Herdr timeout during Give back
review-summary: "givingBack (register.tsx:1790) now stops when 'herdr tab list' exits non-zero or gives no answer: no tab is closed, no text is typed into the recorded pane, spawn.sh is not run; the ticket stays in in-progress, a note 'given back to nobody' carries the Herdr error (last 160 chars), the teamlead session gets a 'none: tab state unreadable' line and a toast (texts.ts tabsUnreadable) tells the person to check the worker tab. One cockpit test covers the timeout case; release note added. Before, a failed list read as 'no tab' and started a second worker."
review-gaps: |-
  blocker: go test ./core/role fails (shipped_test.go:74) - register.tsx and texts.ts changed but core/shipped/shipped.sha256 lacks their new hashes; on master it passes. Fix: scripts/shipped-history.sh --write and commit the sha256 file
  note: elsewhere: herdr-hc.sh:88 / sessions.ts:97-102 - a tab list that exits 0 but is incomplete or unparsable still reads as no tab and starts a second worker; pre-existing (bughunt note), new ticket
test-verdict: "pass: claude plugin test (cockpit) 543/543 RC=0, go test ./core/release/... RC=0; DoD 1 checked at register.tsx:1790 and test 'a failed herdr tab list starts no worker...' (board.test.ts:2527) - herdr tab list exit 1 gives no spawn.sh, no tab close, no send-text, a 'given back to nobody' note and the unreadable toast; DoD 2 core/release/notes/W80HNH.md one line; the cockpit runs only inside Claude Code, so the failure path is exercised through the engine mock, not a live Herdr; new code writes only through the existing jaira note path, no new writer"
hunt-verdict: "pass: bughunt done at b48f82c95, 0 change, 1 pre-existing (same subject: a partial tab list with exit 0 via herdr-hc.sh:88 '|| continue' and tabsOf's [] on bad JSON - note left for the person)"
review-verdict: "send-back: go test ./core/role fails - shipped.sha256 lacks the new hashes of register.tsx and texts.ts; the code change itself meets both DoD items"
review-round: 1
review-verdict-gaps: |-
  blocker: go test ./core/role fails (shipped_test.go:74) - register.tsx and texts.ts changed but core/shipped/shipped.sha256 lacks their new hashes; on master it passes. Fix: scripts/shipped-history.sh --write and commit the sha256 file
  note: elsewhere: herdr-hc.sh:88 / sessions.ts:97-102 - a tab list that exits 0 but is incomplete or unparsable still reads as no tab and starts a second worker; pre-existing (bughunt note), new ticket
review-check: "1. cd core/mod/builtin/jaira-cockpit && claude plugin test .  2. see 543 pass, among them 'give back > a failed herdr tab list starts no worker, types nothing and says the tab state was unreadable'  3. cd ../../../.. && go test ./core/role  4. see ok (today it fails: shipped.sha256 lacks 2 files)  5. live: in a Herdr session make 'herdr tab list' fail (e.g. a herdr wrapper that exits 1), press Give back on a signoff ticket - no new tab opens, the toast says Herdr's tab list could not be read"
round: 2
back-from: review
back-by: lane
back-count: 1
stop-after-fix: "true"
---

# Give back: ein gescheitertes 'herdr tab list' gilt als 'kein Tab offen' und startet einen zweiten Worker neben dem lebenden

## Definition of Done

- [x] Scheitert 'herdr tab list', startet Give back keinen Worker und meldet der Person, dass der Tab-Zustand nicht lesbar war; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1790; test 'give back > a failed herdr tab list starts no worker, types nothing and says the tab state was unreadable' (tests/board.test.ts:2527)
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/W80HNH.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 21:59 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/48662.sock
- **2026-11-24 06:58 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/3417.sock
- **2026-11-24 07:01 · Alexander Sacharov** — decision: a failed tab list also stops the hand-over to the recorded old pane (no send-text), not only the new worker - without the tab list Give back cannot tell a working or dialog-blocked worker from a gone one, and typing blind could answer a dialog. Garbage stdout with exit 0 still reads as no tab (tabsOf unchanged; Herdr does not answer that way).
- **2026-11-24 07:02 · Alexander Sacharov** — critique+optimize one pass (33 lines): the stop follows the atDialog/place.stop shape in givingBack (handedTo none + say), the 160-char tail like spawn.sh's failure; nothing unused, nothing duplicated. Tests 543/543 via claude plugin test.
- **2026-11-24 07:02 · Alexander Sacharov** — waiting: bughunt (background subagent) since 08:02
- **2026-11-24 07:03 · Alexander Sacharov** — bughunt: pre-existing: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:88 - class: a tab list that exits 0 but is incomplete or unreadable is read by givingBack (register.tsx:1795) as 'the ticket has no worker tab', so Give back closes nothing and starts a second worker beside a live one; also sessions.ts:97-102 (tabsOf turns JSON that does not parse, or has no result.tabs, into [] silently); example: under Herdr App, 'hc get <t>' for the worker's terminal times out -> '|| continue' drops that tab, the list still exits 0 without <ID>/worker; blame: a93866dff (2026-11-11, WRMREF) and 2510f1274/86da7c435 (October), before the ticket; ran: HC_BIN=<fake hc whose get fails for t2> bash herdr-hc.sh tab list -> only t1 listed, exit=0
- **2026-11-24 07:03 · Alexander Sacharov** — hunt report: 1 finding, 0 change, 1 pre-existing, as a note on W80HNH | read whole: givingBack/giveBack and their helpers in register.tsx (run, say, herdr, readRecords, paneBlocked, workPlace, spawnWorker, workers), texts.ts, sessions.ts tabsOf/carrierTab/chooseWorker, herdr-hc.sh tab list, the new test in board.test.ts; the other ~2500 lines of register.tsx and the rest of board.test.ts only skimmed, not read line by line | not checked: whether real herdr (not herdr-hc.sh) ever exits 0 with output that is not JSON - no herdr binary here
- **2026-11-24 07:03 · Alexander Sacharov** — bughunt: done at b48f82c95 - 1 finding (0 change, 1 pre-existing) in 4 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/W80HNH.md
- **2026-11-24 07:04 · Alexander Sacharov** — waiting: review (fresh agent via run-lane.sh) since 08:04
- **2026-11-24 07:06 · Alexander Sacharov** — review: done at b48f82c95
- **2026-11-24 07:06 · Alexander Sacharov** — review: blocker: go test ./core/role fails at shipped_test.go:74 - core/shipped/shipped.sha256 lacks the new hashes of core/mod/builtin/jaira-cockpit/hooks/register.tsx and hooks/texts.ts. Run scripts/shipped-history.sh --write and commit core/shipped/shipped.sha256 with the ticket id.
- **2026-11-24 07:08 · Alexander Sacharov** — review blocker closed: scripts/shipped-history.sh --write recorded register.tsx/texts.ts in core/shipped/shipped.sha256; go test ./... RC=0. stop-after-fix=true set by me (sha file only, no lock, no parser). Lesson: the shipped.sha256 gate belongs in the pre-push run for every cockpit change.
