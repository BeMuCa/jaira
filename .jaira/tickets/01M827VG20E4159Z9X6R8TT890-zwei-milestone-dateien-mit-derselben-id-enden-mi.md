---
id: 01M827VG20E4159Z9X6R8TT890
title: Zwei Milestone-Dateien mit derselben ID enden mit einer klaren Fehlermeldung statt still liegen zu bleiben
status: testing
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Liegen nach einem Umbenennen zwei Milestone-Dateien mit derselben ID im Checkout, sagt jaira das und lässt keine zweite Datei still liegen."
context: |-
  Was falsch ist (Review-Notiz 1MYFKE, 22.11.2026): CatchUpMilestones (core/refsync/refsync.go ~496, continue) überspringt eine veraltete Datei unter dem alten Namen, wenn unter dem neuen Namen eine von Hand geänderte Datei derselben ID liegt. Dann liegen zwei Dateien mit einer ID da, und milestone.scan sieht beide.
  Review-Empfehlung: an dieser Stelle mit einem AmbiguousError abbrechen, der beide Dateinamen nennt.
  Hintergrund: seit 1MYFKE sind Milestones nicht mehr in git, sondern reisen über Refs; fetch lässt immer den Ref gewinnen (Option b).
definition-of-done: "Zwei Dateien derselben Milestone-ID nach Catch-up: jaira meldet beide Namen und was zu tun ist (Exit 3), keine wird still übersprungen (Test)"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-11-22T18:21:19Z
updated-at: 2026-11-22T19:09:15Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1868655
claimed-at: 2026-11-22T19:03:16Z
outcome-what: "Every milestone write refuses with exit 3 (milestone_twin_files) naming both files when one milestone ID lies in two files, instead of the catch-up silently leaving the stale one beside the hand-edited one"
outcome-why: "After a rename, a stale old copy beside a hand-edited file under the new name stayed on disk unseen and the next write read whichever file sorted first"
outcome-resolves: 8TT890
review-summary: "no findings: the check sits in milestone.lockAndMigrate, the one door every milestone write passes (core/milestone/write.go:234), rather than in refsync's catch-up, so it covers twins however they arose and needs no syncer; own TwinFilesError rather than AmbiguousError because that one tells the reader to pick a <name>-<ID>, which is wrong for one ID (decision note); classify() split out of report() is the existing errors.Is mapping moved, not a new pattern. Let stand: checkTwinFiles repeats the ReadDir+FileIdentity loop of refsync.milestoneFiles — milestone cannot import refsync, and moving that helper is a refactor outside this ticket."
review-gaps: "removed: slices.Sort in checkTwinFiles (os.ReadDir already returns names sorted). Left: checkTwinFiles' ReadDir+FileIdentity loop beside refsync.milestoneFiles and publishUnsent's carriers count (refsync.go) — same listing, but milestone cannot import refsync and publishUnsent skips twins rather than refusing; folding them is a refactor of its own. Cost: one ReadDir per milestone write, under the lock — negligible."
---

# Zwei Milestone-Dateien mit derselben ID enden mit einer klaren Fehlermeldung statt still liegen zu bleiben

## Definition of Done

- [x] Zwei Dateien derselben Milestone-ID nach Catch-up: jaira meldet beide Namen und was zu tun ist (Exit 3), keine wird still übersprungen (Test)
  proof: core/milestone/write.go checkTwinFiles in lockAndMigrate; TestEveryWriteRefusesOneIDInTwoFiles, TestACatchUpLeavesATwinUnderTheNewName, TestMilestoneWritesRefuseOneIDInTwoFiles (exit 3, milestone_twin_files, both names)
- [x] Der normale Fall (eine Datei, umbenannt) bleibt unverändert (Test)
  proof: TestACatchUpRemovesAnOldCopyBesideTheRenamedFile, TestACatchUpAfterARenameInCaseOnlyKeepsTheMilestone (core/refsync/catchup_test.go) still green
- [x] Release-Note core/release/notes/<ID>.md
  proof: core/release/notes/8TT890.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 19:03 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 19:04 · Alexander Sacharov** — decision: open decisions — none. Own error type milestone.TwinFilesError{ID, Files} instead of AmbiguousError (review suggestion): AmbiguousError's message is about one name on two milestones and tells the reader to pick a <name>-<ID>, which is wrong for one ID in two files. Mapped to exit 3 (reason milestone_twin_files) centrally in cli.report so every milestone write (create/add/rm/rename/set) refuses alike. Catch-up still brings the other files up and returns the twin error after the loop.
- **2026-11-22 19:08 · Alexander Sacharov** — in-progress: class = one ID in >1 ID-named file, however it got there (not only the rename case), so the check sits in milestone.lockAndMigrate after catch-up and Migrate and covers create/edit/rename/add/rm/set. Four categories: concurrency — check runs under the milestone lock; error paths — lock released, nothing written/recorded (tests); unusual input — legacy files without ID skipped (Migrate's), case-only names on case-insensitive FS list once; permissions — n/a. Not done: reads (milestone ls, board) still show only the first file of the twin (scan dedups by ID) — reads write nothing, out of scope.
- **2026-11-22 19:08 · Alexander Sacharov** — critique: nothing to change; see review-summary. Considered and rejected: returning the error from CatchUpMilestones itself (covers only the rename case, and a board without remote never runs the catch-up).
