---
id: 01M827VG20E4159Z9X6R8TT890
title: Zwei Milestone-Dateien mit derselben ID enden mit einer klaren Fehlermeldung statt still liegen zu bleiben
status: critique
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
updated-at: 2026-11-22T19:34:15Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-3443254
claimed-at: 2026-11-22T19:28:38Z
outcome-what: "Every milestone write — create, add, rm, set, rename, and logbook/restore of a milestone — refuses with exit 3 (milestone_twin_files) naming both files when one milestone ID lies in two files"
outcome-why: "After a rename, a stale old copy beside a hand-edited file under the new name stayed on disk unseen and the next write read whichever file sorted first"
outcome-resolves: 8TT890
review-summary: "internal/cli/archive.go:184 parses res.Moved's display string ('old -> new') to find the restored file; where the ID file already exists Migrate only removes the old file and still lists it, so the rollback would move a board file (possibly a twin) into the logbook and lose the restored content. Simpler: run milestone.CheckTwinFiles in unfileMilestone before migrateHeld, so a pair already on the board refuses before anything is migrated and dst->src is the whole rollback; drop the Moved parsing"
review-gaps: |-
  note: internal/cli/archive.go:183 - restore rollback renames dst back to src, but Migrate inside migrateHeld may already have moved an old-format dst to <name>-<newID>.md, so the rollback misses and the file stays on the board marked filed with no logbook copy (bughunt 19:18); fix: run CheckTwinFiles before RestoreFrom moves the file, or roll back the path in res.Moved
  note: internal/cli/archive.go:181 comment says the restored file is one of the two twins - it need not be (an unrelated pre-existing pair refuses too); rollback is still right, reword the comment
  note: reads (jaira milestone ls, the board) still show only the first file of a twin (milestone.scan dedups by ID) - nothing warns until a write; a ls warning line would surface it earlier
  note: refsync.go:571 background Migrate/mergeRefused does not run CheckTwinFiles - next person-facing write refuses, so no wrong result, only a delayed one
test-verdict: "pass: round 2 — go test ./... RC=0 and go vet ./... at b6a0154bc; -race on core/milestone + internal/cli milestone/logbook/restore tests RC=0. Built binary on a scratch board: logbook q4 beside q5-<ID>.md rc=3 naming both, both files stay unfiled; without the twin logbook rc=0; restore of the filed q4-<ID>.md beside a board a1-<ID>.md rc=3, q4 back in the logbook still filed, board keeps a1 only. DoD 1-3 re-checked; round-1 rename/normal-case results stand."
hunt-verdict: "pass: bughunt done at b6a0154bc, 1 change (not a blocker: restore of an old-format q4.md beside a twin pair misses its give-back because Migrate renamed the file first — its bughunt note stands for review and the person), 0 pre-existing"
round: 4
back-from: critique
back-by: lane
back-count: 1
review-verdict: "accept: all three DoD items hold - every milestone write path (core Create/Edit/Rename via lockAndMigrate, CLI logbook/restore via migrateHeld) refuses one ID in two files with exit 3 naming both; the normal rename catch-up tests stay green; release note present. go test ./core/milestone ./core/refsync ./internal/cli green. Open findings are narrow edges, not blockers."
review-round: 2
review-verdict-gaps: |-
  note: internal/cli/archive.go:183 - restore rollback renames dst back to src, but Migrate inside migrateHeld may already have moved an old-format dst to <name>-<newID>.md, so the rollback misses and the file stays on the board marked filed with no logbook copy (bughunt 19:18); fix: run CheckTwinFiles before RestoreFrom moves the file, or roll back the path in res.Moved
  note: internal/cli/archive.go:181 comment says the restored file is one of the two twins - it need not be (an unrelated pre-existing pair refuses too); rollback is still right, reword the comment
  note: reads (jaira milestone ls, the board) still show only the first file of a twin (milestone.scan dedups by ID) - nothing warns until a write; a ls warning line would surface it earlier
  note: refsync.go:571 background Migrate/mergeRefused does not run CheckTwinFiles - next person-facing write refuses, so no wrong result, only a delayed one
review-check: "1. go build -o /tmp/j ./cmd/jaira (from the worktree)  2. in a scratch repo: /tmp/j init, then /tmp/j milestone create q4  3. cp .jaira/milestones/q4-*.md to a file named q5-<same ID>.md in the same folder  4. /tmp/j milestone create other ; echo $?  -> message names both files (q4-<ID>.md and q5-<ID>.md), says nothing was written, exit 3  5. /tmp/j logbook q4 ; echo $? -> same refusal, exit 3, both files still in .jaira/milestones  6. delete q5-<ID>.md, rerun step 4 -> milestone created, exit 0  Automated: go test ./internal/cli -run 'TwinFiles|OneIDInTwoFiles'"
---

# Zwei Milestone-Dateien mit derselben ID enden mit einer klaren Fehlermeldung statt still liegen zu bleiben

## Definition of Done

- [x] Zwei Dateien derselben Milestone-ID nach Catch-up: jaira meldet beide Namen und was zu tun ist (Exit 3), keine wird still übersprungen (Test)
  proof: core/milestone/write.go checkTwinFiles in lockAndMigrate; TestEveryWriteRefusesOneIDInTwoFiles, TestACatchUpLeavesATwinUnderTheNewName, TestMilestoneWritesRefuseOneIDInTwoFiles (exit 3, milestone_twin_files, both names)
- [x] Der normale Fall (eine Datei, umbenannt) bleibt unverändert (Test)
  proof: TestACatchUpRemovesAnOldCopyBesideTheRenamedFile, TestACatchUpAfterARenameInCaseOnlyKeepsTheMilestone (core/refsync/catchup_test.go) still green
- [x] Release-Note core/release/notes/<ID>.md
  proof: core/release/notes/8TT890.md
- [x] Restore einer alten Milestone-Datei ohne ID (q4.md) neben einem Zwilling: die Datei liegt danach wieder im Logbook, filed, nicht auf dem Board (Test)
  proof: internal/cli/archive.go unfileMilestone runs milestone.CheckTwinFiles before migrateHeld; TestRestoreOfAnOldFormatFileBesideATwinGoesBackToTheLogbook (internal/cli/milestones_test.go)

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
- **2026-11-22 19:11 · Alexander Sacharov** — bughunt: change: blocker: internal/cli/logbook.go:357 and internal/cli/archive.go:133 (unfileMilestone, archive.go:175) - class: milestone writes that take the milestone lock themselves instead of going through milestone.lockAndMigrate never reach checkTwinFiles, so with one ID in two files they still read whichever file milestone.scan keeps first (core/milestone/milestone.go:333 dedups by ID) and save/record it; core/release/notes/8TT890.md:1 and the TwinFilesError doc (core/milestone/write.go:73) promise every milestone command refuses. example 1: q4-ID.md + hand-edited q5-ID.md, 'jaira logbook q4' exits 0, files q4-ID.md, leaves q5-ID.md live on the board (on a board with a remote the ref would say filed with q4's content). example 2: logbook holds q4-ID.md, board holds a1-ID.md, 'jaira restore q4-ID.md' exits 0 'Restored', but milestone.Load(id) picks a1, so q4-ID.md keeps status: filed and stays off the board, and recordMilestone queues a1's content. ran: scratch repo with built binary - milestone create other exit 3 (guard works), logbook q4 exit 0, restore exit 0 with q4 still filed. why change and not pre-existing although blame says b54bb3a5e (older): the gap is in the scope this change claims; if its definition of done only covers Edit/Rename/Create, re-sort it as pre-existing
- **2026-11-22 19:11 · Alexander Sacharov** — bughunt: done at 27866b60c - 1 finding (1 change, 0 pre-existing) in 4 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/milestone/write.go core/milestone/write_test.go core/refsync/catchup_test.go core/refsync/refsync.go core/release/notes/8TT890.md internal/cli/milestones_test.go internal/cli/root.go
- **2026-11-22 19:14 · Alexander Sacharov** — testing: class: milestone writes that take the lock themselves (CLI logbook, restore via migrateHeld) skip CheckTwinFiles; example: bughunt note (logbook q4 exit 0 beside q5-<ID>.md); fix: check in migrateHeld, the CLI's lockAndMigrate, so every caller gets it; restore puts its file back in the logbook on refusal.
- **2026-11-22 19:15 · Alexander Sacharov** — in-progress round 2: CheckTwinFiles exported and run in cli.migrateHeld too (logbook, restore, and the pre-lock migrateMilestones of create/edit/rename, which now refuse one step earlier with the same error). restore on refusal renames the file back to the logbook (still filed), test TestLogbookAndRestoreRefuseOneIDInTwoFiles. Left: refsync.mergeRefused (background merge of a refused write) writes the ref's merge to the chosen file without reading the disk twin — not a person-facing command, the next command refuses anyway; note for review.
- **2026-11-22 19:18 · Alexander Sacharov** — bughunt: change: internal/cli/archive.go:183 - class: the twin-refusal rollback renames dst back to src, but migrateHeld (Migrate) may already have moved dst under an ID name, so the rollback misses and the restored file is stranded on the board still marked filed, with no logbook copy left; same stranding, pre-existing, for any other migrateHeld error at archive.go:188; example: logbook holds an old-format filed q4.md (pre-ID build) and the board has q5-<ID>.md + q6-<ID>.md -> 'jaira restore q4.md' exits 3 'nothing was written' plus 'rename .../milestones/q4.md ...: no such file', leaves .jaira/milestones/q4-<newID>.md filed=true (invisible), and a second 'jaira restore q4.md' answers 'not in the archive'; not a blocker: narrow path, refusal itself is right, file recoverable by hand; fix direction: run CheckTwinFiles before RestoreFrom moves the file, or roll back the path Migrate reports in res.Moved; ran: scratch test in internal/cli (removed) restoring old-format q4.md beside a twin pair
- **2026-11-22 19:18 · Alexander Sacharov** — bughunt: done at b6a0154bc - 1 finding (1 change, 0 pre-existing) in 6 modules, follow-up hunt since 27866b60c; concurrency, error paths, input, permissions checked; files: core/milestone/write.go core/milestone/write_test.go core/refsync/catchup_test.go core/refsync/refsync.go core/release/notes/8TT890.md internal/cli/archive.go internal/cli/milestones.go internal/cli/milestones_test.go internal/cli/root.go
- **2026-11-22 19:28 · Alexander Sacharov** — Given back from signoff: beheben:
dazu: 🟡 редкий случай: restore старого q4.md рядом с двойником теряет запись о возврате
angenommen:
- note: internal/cli/archive.go:183 - restore rollback renames dst back to src, but Migrate inside migrateHeld may already have moved an old-format dst to <name>-<newID>.md, so the rollback misses and the file stays on the board marked filed with no logbook copy (bughunt 19:18); fix: run CheckTwinFiles before RestoreFrom moves the file, or roll back the path in res.Moved
- note: internal/cli/archive.go:181 comment says the restored file is one of the two twins - it need not be (an unrelated pre-existing pair refuses too); rollback is still right, reword the comment
- note: reads (jaira milestone ls, the board) still show only the first file of a twin (milestone.scan dedups by ID) - nothing warns until a write; a ls warning line would surface it earlier
- note: refsync.go:571 background Migrate/mergeRefused does not run CheckTwinFiles - next person-facing write refuses, so no wrong result, only a delayed one
- **2026-11-22 19:28 · Alexander Sacharov** — cockpit: given back to a new worker - nothing recorded about the old worker session (cache and context unknown)
- **2026-11-22 19:32 · Alexander Sacharov** — in-progress round 3 (give-back): rollback in unfileMilestone now renames from the path Migrate moved dst to (res.Moved 'old -> new'), not dst; comment reworded (the restored file need not be one of the twins). Chose this over running CheckTwinFiles before RestoreFrom: that would refuse ticket restores too whenever the board has a twin pair, and the twins are only known after catch-up/Migrate. Other migrateHeld errors (archive.go, pre-existing stranding) left as they were — not asked for.
- **2026-11-22 19:32 · Alexander Sacharov** — critique round 3: internal/cli/archive.go:184 parses res.Moved's display string ('old -> new') to find the restored file; where the ID file already exists Migrate only removes the old file and still lists it, so the rollback would move a board file (possibly a twin) into the logbook and lose the restored content. Simpler: run milestone.CheckTwinFiles in unfileMilestone before migrateHeld, so a pair already on the board refuses before anything is migrated and dst->src is the whole rollback; drop the Moved parsing
- **2026-11-22 19:34 · Alexander Sacharov** — in-progress round 4: critique finding applied - twin check before migrateHeld in unfileMilestone, Moved parsing dropped. Residual: a twin that only catch-up creates (inside migrateHeld) beside an old-format restored file still cannot be rolled back (rename fails loudly, joined into the error) - needs a ref rename racing a hand copy, left.
