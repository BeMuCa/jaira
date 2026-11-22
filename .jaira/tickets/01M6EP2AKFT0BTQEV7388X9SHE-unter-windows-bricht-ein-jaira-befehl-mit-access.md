---
id: 01M6EP2AKFT0BTQEV7388X9SHE
title: "Unter Windows bricht ein jaira-Befehl mit 'Access is denied' ab, statt auf eine gerade freigegebene Sperre zu warten"
status: testing
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Eine Sperre unter .../locks/ wartet unter Windows auch dann, wenn die Sperrdatei gerade geloescht wird, und der Befehl laeuft danach normal durch."
context: |-
  Was falsch ist:
  - core/ticket/store.go:1531 lock() legt die Sperrdatei mit O_CREATE|O_EXCL an. Nur bei os.IsExist wird gewartet und neu versucht; jeder andere Fehler bricht ab.
  - Unter Windows liefert das Anlegen 'Access is denied' (ERROR_ACCESS_DENIED), solange die alte Sperrdatei noch im Loeschen ist (os.Remove im unlock, Zeile 1544). Das ist kein IsExist -> der Befehl bricht ab, statt zu warten.
  Gefunden am 02.11.2026 in CI von PR #120 (182JT1, Lauf 41492357461, test windows-latest):
  - TestRestoreWaitsForTheMilestoneLock (internal/cli/milestones_test.go:559) -> 'restore after unlock: open ...\locks\milestones.lock: Access is denied.'
  - 182JT1 aendert nichts an Sperren oder Milestones; master lief am 02.11. unter Windows gruen. Also flackernd, abhaengig vom Timing.
  Warum es zaehlt: jaira soll mehrere parallele Sessions aushalten. Unter Windows kann jede Sperre (Tickets, Milestones) genau im Moment der Freigabe so abbrechen.
  Naheliegende Loesung: in lock() ERROR_ACCESS_DENIED (unter Windows) wie IsExist behandeln und bis lockTimeout weiter versuchen.
definition-of-done: lock() wartet unter Windows auch bei ERROR_ACCESS_DENIED bis lockTimeout und wirft den Fehler erst danach; Test
tags:
  - concurrency
blocked-by: []
related: []
commits: []
created-at: 2026-11-02T17:49:13Z
updated-at: 2026-11-22T18:53:41Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-903238
claimed-at: 2026-11-22T18:49:04Z
outcome-what: "Windows-only tests in core/ticket/lock_windows_test.go pin lockFile's wait on a delete-pending lock file: it waits, takes the lock once the old file goes, and gives up with Access is denied at lockTimeout. The fix itself had already landed in 32TTSS."
outcome-why: The 32TTSS fix had no test of its own; a refactor of lockFile could drop the Windows branch and only a flaky CI run would notice.
outcome-resolves: "8X9SHE: lock waits on Windows when the lock file is still being deleted - proven by a deterministic test and 20/20 runs of TestRestoreWaitsForTheMilestoneLock on Windows"
review-summary: none
review-gaps: "replaced the hand-rolled once flag in deletePending with sync.OnceFunc; no duplicate helper exists (no other FileDispositionInfo/CreateFile use in the repo); nothing dead, nothing on a hot path"
---

# Unter Windows bricht ein jaira-Befehl mit 'Access is denied' ab, statt auf eine gerade freigegebene Sperre zu warten

## Definition of Done

- [x] lock() wartet unter Windows auch bei ERROR_ACCESS_DENIED bis lockTimeout und wirft den Fehler erst danach; Test
  proof: core/ticket/lock_windows_test.go TestLockFileWaitsForADeletePendingLock + TestLockFileReportsAPersistentDenialAtTheDeadline: green on Windows 11 (-count=3); both fail with the CI's 'Access is denied' when store.go:1645 is reverted
- [x] TestRestoreWaitsForTheMilestoneLock laeuft auf windows-latest stabil (z. B. -count=20 in einem Lauf gruen)
  proof: internal/cli/milestones_test.go:559 TestRestoreWaitsForTheMilestoneLock -test.count=20 on a local Windows 11 (GOOS=windows test binary run from WSL): 20/20 PASS, exit 0; windows-latest CI runs it again on the PR
- [x] Release-Note-Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/0.5.14/32TTSS.md:3 already tells users of the fix (shipped in 0.5.14); this ticket adds only a test, which is not client-facing, so no 8X9SHE.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 18:49 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 18:49 · Alexander Sacharov** — decision: the lock fix itself already landed in bd77914c7 (32TTSS, 16.11, core/ticket/store.go:1645) and its user line shipped in core/release/notes/0.5.14/32TTSS.md:3. What this ticket still owes: a test that pins the Windows wait (none exists) and a stability run on Windows. A test is not client-facing, so no new core/release/notes/8X9SHE.md - DoD 3 is ticked with the 0.5.14 line as proof. Windows test: make the lock file delete-pending deterministically (SetFileInformationByHandle FileDispositionInfo on an open handle, legacy non-POSIX semantics) - os.Remove alone uses POSIX delete on current Windows and frees the name at once, so it would not reproduce the race. Stability run from WSL: GOOS=windows test binary run on the Windows side.
- **2026-11-22 18:52 · Alexander Sacharov** — critique: considered pulling the busy check into a pure lockBusy(err, goos) helper with a table test on every OS - rejected, it would test the if, not that Windows really answers a delete-pending create with Access is denied; the real-handle test proves both. x/sys/windows is already a dependency (core/session/running_windows.go). The 5 s persistent-denial test runs only on Windows and is the only way to show the deadline with lockTimeout a const.
