---
id: 01M846VR7GYFYHHBZMVEBPR7C6
title: "projects.json: eine kaputte Datei wird nicht mehr still mit einem Eintrag überschrieben"
status: done
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
commits:
  - ec7a859fec1cfd2b7c05851e5704bab382d24f16
  - 28fbab3a5595bd25751ff21b1dcf7cf3057423d5
created-at: 2026-11-23T12:42:28Z
updated-at: 2026-11-23T15:55:56Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-74262
claimed-at: 2026-11-23T13:30:23Z
outcome-what: "Remember copies an unparsable projects.json aside (.broken / .broken-<time>) with one stderr line before writing a new list; a projects.json that exists but cannot be read, or a copy that cannot be written, leaves the file untouched and returns an error"
outcome-why: "Load reads a broken or unreadable file as an empty list, so the next Remember silently kept one board and dropped the rest"
outcome-resolves: "BPR7C6: a broken or unreadable project list is no longer lost without a word"
review-summary: "Remember (core/project/project.go) now runs keepBroken under the registry lock before it reads and rewrites projects.json: a file that does not parse as a project list is copied to projects.json.broken (or projects.json.broken-<UTC time> when that exists) and one stderr line names the copy; a missing or whitespace-only file is left alone; a file that exists but cannot be read makes Remember return an error and write nothing, so the list is no longer replaced by a one-entry list. Load stays read-only and lenient. Four tests cover broken, missing/empty, a directory at the path and a 0000 file; release note in core/release/notes/BPR7C6.md."
review-gaps: |-
  note: limit: the stderr line is written over by the full-screen TUI when Remember runs from a board switch (internal/tui/model.go:420, browse.go) - the copy is still made; accepted in the decision note
  note: limit: Remember reads the file twice (keepBroken, then Load); a read error only in the second read would still yield a one-entry list - needs a transient failure between two reads under the lock, not reproducible here (bughunt note)
  note: elsewhere: every caller except 'jaira projects add' (internal/cli/home.go:53, board.go:24, tickets.go:53, internal/tui/*) drops Remember's error, so on an unreadable projects.json a normal start says nothing and the list just stops updating - there before this change
test-verdict: "pass: go test -race ./core/project green (rc=0), full suite green in round 1; DoD 1-5 verified; built binary on scratch boards: broken file -> .broken + one stderr line, second -> .broken-<time>, first copy kept; chmod 000 -> rc=1, file untouched; two parallel 'projects add' both land"
hunt-verdict: "pass: bughunt done at 28fbab3a5, 0 change, 0 pre-existing; round-1 TUI-overprint edge stays a note for the person"
round: 2
back-from: testing
back-by: lane
back-count: 1
review-verdict: "accept: all five DoD items hold in code and tests; go test ./core/project passes; no blocker found"
review-round: 2
review-verdict-gaps: |-
  note: limit: the stderr line is written over by the full-screen TUI when Remember runs from a board switch (internal/tui/model.go:420, browse.go) - the copy is still made; accepted in the decision note
  note: limit: Remember reads the file twice (keepBroken, then Load); a read error only in the second read would still yield a one-entry list - needs a transient failure between two reads under the lock, not reproducible here (bughunt note)
  note: elsewhere: every caller except 'jaira projects add' (internal/cli/home.go:53, board.go:24, tickets.go:53, internal/tui/*) drops Remember's error, so on an unreadable projects.json a normal start says nothing and the list just stops updating - there before this change
review-check: |-
  1. cd /home/alex/projects/.worktrees/jaira-bpr7c6 && go build -o /tmp/jb ./cmd/jaira
  2. export JAIRA_HOME=$(mktemp -d) && printf '[{"root":"/a",}]' > $JAIRA_HOME/projects.json
  3. /tmp/jb projects add .  -> stderr shows one line '... a copy is at .../projects.json.broken'
  4. cat $JAIRA_HOME/projects.json.broken -> the broken text; cat $JAIRA_HOME/projects.json -> a valid list with this board
  5. chmod 000 $JAIRA_HOME/projects.json && /tmp/jb projects add .  -> exit 1, 'could not be read, so it is left as it is'; chmod 644 and the file is unchanged
  6. go test ./core/project -run Remember -> ok
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
- **2026-11-23 13:35 · Alexander Sacharov** — critique round 2: none. The read-error check sits in keepBroken beside the parse check, the one guard before Remember writes; Load stays read-only and lenient (decision note).
- **2026-11-23 13:37 · Alexander Sacharov** — hunt report: 0 findings, nothing written beyond this line | read as a diff since ec7a859fe: core/project/project.go (keepBroken + Remember/Load around it), core/project/worktree_test.go, core/release/notes/BPR7C6.md; ran: go test ./core/project (ok), built binary: projects.json chmod 000 + 'jaira projects add .' -> exit 1, file unchanged; 4 parallel adds over a truncated file -> one projects.json.broken, valid list | not confirmable: Remember still reads the file twice (keepBroken, then Load), and a read error only in Load's read would still yield a one-entry list - needs a transient read failure (e.g. Windows sharing violation) between two reads under the lock, could not be produced here
- **2026-11-23 13:37 · Alexander Sacharov** — bughunt: done at 28fbab3a5 - 0 findings (0 change, 0 pre-existing) in 3 modules, follow-up hunt since ec7a859fe; concurrency, error paths, input, permissions checked; files: core/project/project.go core/project/worktree_test.go core/release/notes/BPR7C6.md
- **2026-11-23 13:38 · Alexander Sacharov** — waiting: review lane, fresh agent via run-lane.sh, since 14:38
- **2026-11-23 13:40 · Alexander Sacharov** — review: done at 28fbab3a5
- **2026-11-23 13:40 · Alexander Sacharov** — report: Remember copies an unparsable projects.json to .broken(/-<time>) with one stderr line before writing; an unreadable one is left untouched with exit 1 | signoff, context unknown | review accept, 3 notes (2 limits, 1 elsewhere: callers drop Remember's error); branch not pushed
