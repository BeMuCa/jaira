---
id: 01M846VR7GYFYHHBZMVEBPR7C6
title: "projects.json: eine kaputte Datei wird nicht mehr still mit einem Eintrag überschrieben"
status: critique
ready: true
creator: Alexander Sacharov
goal: "Ist ~/.jaira/projects.json nicht lesbar, geht die Projektliste nicht verloren: jaira legt eine Kopie an, sagt es auf stderr und überschreibt das Original nicht still."
context: |-
  Was heute falsch ist: core/project/project.go:56 (Load) liest eine projects.json, die sich nicht parsen lässt, als leere Liste. Der nächste Remember() schreibt dann eine Liste mit nur dem aktuellen Projekt - alle anderen Einträge sind weg, ohne Meldung.
  Wann das passiert: Datei von Hand bearbeitet und dabei ein Komma vergessen, oder eine alte zerrissene Datei aus der Zeit vor ZXJ72E (das parallele Schreiben ist seit ZXJ72E atomar und gesperrt).
  Gefunden: Review von ZXJ72E als 'elsewhere'-Notiz, vom Maintainer am 23.11.2026 zur Reparatur freigegeben.
  Bekannt: Remember() schreibt seit ZXJ72E unter ticket.LockFile + ticket.WriteAtomic; core/settings macht für settings.json schon eine .bak-Kopie (settings.json.bak liegt in ~/.jaira).
definition-of-done: "Eine projects.json, die sich nicht parsen lässt, wird vor dem nächsten Schreiben nach projects.json.broken kopiert (eine bestehende .broken wird nicht still überschrieben, z.B. Zeitstempel im Namen)"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-11-23T12:42:28Z
updated-at: 2026-11-23T13:35:42Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-74262
claimed-at: 2026-11-23T13:30:23Z
outcome-what: "Remember copies an unparsable projects.json aside (.broken / .broken-<time>) with one stderr line before writing a new list; a projects.json that exists but cannot be read, or a copy that cannot be written, leaves the file untouched and returns an error"
outcome-why: "Load reads a broken or unreadable file as an empty list, so the next Remember silently kept one board and dropped the rest"
outcome-resolves: "BPR7C6: a broken or unreadable project list is no longer lost without a word"
review-summary: none
review-gaps: "none: no other .broken/copy-aside helper in the repo; Remember now reads projects.json twice (keepBroken, then Load) on every start - a small file, left rather than threading the bytes through Load's ~10 callers; gofmt also reordered worktree_test.go's pre-existing import order, the file this change touches"
test-verdict: "pass: go test -race -count=1 ./... green (rc=0), DoD 1-5 verified by the two new tests"
hunt-verdict: "fail: projects.json that exists but cannot be read (chmod 000) is still replaced by a one-entry list"
round: 2
back-from: testing
back-by: lane
back-count: 1
---

# projects.json: eine kaputte Datei wird nicht mehr still mit einem Eintrag überschrieben

## Definition of Done

- [x] Eine projects.json, die sich nicht parsen lässt, wird vor dem nächsten Schreiben nach projects.json.broken kopiert (eine bestehende .broken wird nicht still überschrieben, z.B. Zeitstempel im Namen)
  proof: core/project/project.go:keepBroken, TestRememberKeepsBrokenList
- [x] jaira meldet auf stderr einmal, dass die Liste nicht lesbar war und wo die Kopie liegt
  proof: core/project/project.go:keepBroken Fprintf to Stderr; TestRememberKeepsBrokenList (one line, then silent)
- [x] Eine leere oder fehlende projects.json verhält sich wie bisher (keine Kopie, keine Meldung)
  proof: TestRememberMissingOrEmptyListIsQuiet
- [x] Tests: kaputte Datei -> Kopie + Meldung + neue Liste; fehlende Datei -> nichts davon
  proof: core/project/worktree_test.go: TestRememberKeepsBrokenList, TestRememberMissingOrEmptyListIsQuiet
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/BPR7C6.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 13:30 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/6958.sock
- **2026-11-23 13:30 · Alexander Sacharov** — decision: no open decisions for the person. Copy goes to projects.json.broken; if that exists, projects.json.broken-<UTC yyyymmddThhmmssZ> (never overwritten). Copy+message only in Remember (the write path), not in Load: a read alone loses nothing. Empty/whitespace-only file counts as missing. Message goes to a package writer (os.Stderr by default, swappable in tests); a TUI board switch mid-session could print over the screen, but the first Remember runs before the TUI starts (internal/cli/home.go openBoardAt) and rewrites the file valid, so that only hits a file broken by hand during a session.
- **2026-11-23 13:32 · Alexander Sacharov** — critique: none. Checked: simpler shape would be Load returning an error, but Load has ~10 read-only callers that should not copy anything - the write path (Remember) is the one place that loses data, so the check sits there. settings.json.bak (core/settings) is a per-write backup, a different job; not reused. Stderr var has one writer but tests need it. TUI-overprint edge stands on the decision note.
- **2026-11-23 13:34 · Alexander Sacharov** — bughunt: change: blocker: core/project/project.go:158 - class: projects.json exists but os.ReadFile fails for any reason other than not-exist (permission denied, EIO, a directory, a file another process holds without share-read on Windows); keepBroken returns nil on every read error, Load (project.go:55-57) reads the same error as an empty list, and Remember's WriteAtomic renames a one-entry list over the file - every other board is gone, no .broken copy, no stderr line, exit 0. Mark change because the new guard on line 158 decides to go ahead on that path; the same loss happened before the ticket too. example: two boards registered, chmod 000 projects.json, then 'jaira projects add b1' prints Registered and exits 0 and the list holds only b1; ran: built ./cmd/jaira with JAIRA_HOME in a scratch dir (exit 0, b2 lost, no projects.json.broken)
- **2026-11-23 13:34 · Alexander Sacharov** — bughunt: change: core/project/project.go:168 - class: the new notice is written straight to os.Stderr from Remember, and Remember is also called while the full-screen TUI is running (internal/tui/model.go:420 switchBoard, internal/tui/browse.go:128/160/173/279); there the line lands on the alternate screen and the next frame paints over it, so the one message saying where the copy went is never seen. Every caller there also ignores Remember's error, so 'could not be copied aside' is dropped as well. The copy is still made, so nothing is lost - edge only. example: the file breaks while a board is open (for example an older binary writing it), then the user switches board in the TUI; ran: read
- **2026-11-23 13:34 · Alexander Sacharov** — hunt report: 2 findings, both change (1 blocker: unreadable file still silently replaces the list; 1 edge: notice invisible inside the TUI), notes on BPR7C6 | read whole: core/project (project.go, worktree.go, discover.go, worktree_test.go), release note BPR7C6.md; plus Remember's callers in internal/tui and internal/cli and ticket.WriteAtomic/LockFile | not checked: the Windows sharing-violation variant (no Windows here); the TUI overpaint was confirmed by reading, not by running
- **2026-11-23 13:34 · Alexander Sacharov** — bughunt: done at ec7a859fe - 2 findings (2 change, 0 pre-existing) in 3 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/project/project.go core/project/worktree_test.go core/release/notes/BPR7C6.md
- **2026-11-23 13:35 · Alexander Sacharov** — testing: class: any read error other than not-exist lets Remember overwrite the list; example: chmod 000 projects.json (bughunt note); fix: keepBroken returns an error for every ReadFile error except fs.ErrNotExist, so Remember writes nothing; test with a directory at the path (portable, works as root too). TUI overprint edge stays a note (decision note above).
