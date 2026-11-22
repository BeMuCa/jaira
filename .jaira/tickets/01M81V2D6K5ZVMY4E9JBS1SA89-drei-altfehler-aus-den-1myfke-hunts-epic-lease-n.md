---
id: 01M81V2D6K5ZVMY4E9JBS1SA89
title: "Drei Altfehler aus den 1MYFKE-Hunts: Epic-Lease nach ErrRaceLost, CRLF-AGENTS.md driftet immer, TUI verschluckt Fetch-Fehler"
status: testing
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Die drei Altfehler sind behoben, jeder mit einem Test, der ohne den Fix fehlschlägt."
context: |-
  Drei Altfehler, die die Bughunts an 1MYFKE am 22.11.2026 gefunden haben. Sie stammen nicht von 1MYFKE; Alex will sie behoben haben.
  - core/refsync/epic.go:121-137: rebaseEpic hält nach ErrRaceLost den Lease nicht richtig (Rennen zwischen zwei Schreibern auf den Epic-Ref).
  - core/board/announce.go: eine AGENTS.md mit CRLF-Zeilenenden gilt immer als abgewichen (drifted), auch direkt nach jaira update.
  - internal/tui/refs.go:107: die TUI verschluckt den RecordError aus dem Hintergrund-FetchMilestones; ein fehlgeschlagener Fetch wird nirgends gezeigt.
  - Details stehen in den Notizen von 1MYFKE (review-gaps, Notiz "pre-existing, from the hunts").
definition-of-done: "rebaseEpic hält nach ErrRaceLost den Lease korrekt (core/refsync/epic.go, Test)"
tags:
  - concurrency
blocked-by: []
related: []
commits: []
created-at: 2026-11-22T14:37:54Z
updated-at: 2026-11-22T14:46:27Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-3557981
claimed-at: 2026-11-22T14:38:23Z
outcome-what: "rebaseEpic queues an unsent merge on the lease it was merged onto; the agent-file block follows the file's CRLF line endings; the board warns when a fetched milestone could not be queued"
outcome-why: "three pre-existing defects from the 1MYFKE hunts: a teammate's epic change overwritten on the next flush, CRLF AGENTS.md always drifted, a milestone silently never reaching the remote"
outcome-resolves: S1SA89
review-summary: "none"
round: 2
back-from: critique
back-by: lane
back-count: 1
review-gaps: "moved lineEnding above appendEntry's doc comment (it had split the comment from its function); left core/ticket/frontmatter.go:64 alone — it reads the first line's ending on purpose, a different rule in another package; epic.go repeats two lines of RecordEpic (QueueEpic + dirty) rather than widening the epic.Refs interface; no dead code, no repeated work on a hot path"
---

# Drei Altfehler aus den 1MYFKE-Hunts: Epic-Lease nach ErrRaceLost, CRLF-AGENTS.md driftet immer, TUI verschluckt Fetch-Fehler

## Definition of Done

- [x] rebaseEpic hält nach ErrRaceLost den Lease korrekt (core/refsync/epic.go, Test)
  proof: core/refsync/epic.go rebaseEpic oursLease; TestAnEpicThatKeepsLosingIsRequeuedOnTheStateItWasMergedOnto (core/refsync/epic_test.go)
- [x] AGENTS.md mit CRLF gilt nach jaira update nicht als drifted (core/board/announce.go, Test)
  proof: core/board/announce.go managedBlock + announceInAgentFile CRLF; TestACRLFAgentFileIsNotDrifted (core/board/crlf_test.go)
- [x] Ein fehlgeschlagener Hintergrund-FetchMilestones wird in der TUI angezeigt statt verschluckt (internal/tui/refs.go, Test)
  proof: internal/tui/refs.go refFetchedMsg.warn + internal/tui/model.go refFetchedMsg notify; TestTheBoardSaysWhenAMilestoneCouldNotBeQueued (internal/tui/recorderror_test.go)
- [x] Release-Note core/release/notes/<ID>.md, shipped.sha256 falls nötig
  proof: core/release/notes/S1SA89.md; shipped.sha256 unverändert (keine installierte Datei geändert, scripts/shipped-history.sh --check ok)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 14:38 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 14:40 · Alexander Sacharov** — Epic-Lease: Fix = der Text 'ours' wird mit dem Lease des Standes eingereiht, auf den er gemergt wurde (oursLease), nicht mit EpicSHA nach dem Fetch. Test: pre-receive-Hook, der den Ref bei jedem Push weiterschiebt und ablehnt; Hook muss GIT_OBJECT_DIRECTORY/QUARANTINE unsetten, sonst verschwinden seine Objekte.
- **2026-11-22 14:44 · Alexander Sacharov** — CRLF: Regel = Datei enthält irgendwo CRLF -> Block in CRLF (wie gitignore.go appendEntry). Auch der appended-Pfad schreibt CRLF, sonst driftet ein frisch angehängter Block sofort. TUI: RecordError wird als refFetchedMsg.warn durchgereicht und mit dem Wortlaut von 'jaira fetch' per notify gezeigt; Test baut den echten Fehler (alte backend.md + Outbox read-only).
- **2026-11-22 14:44 · Alexander Sacharov** — critique: (1) CRLF-Erkennung dreimal (announce.go x2, gitignore.go:442) -> ein Helper lineEnding(s) in core/board. (2) Warnsatz in TUI und CLI doppelt -> (*milestone.RecordError).Warning() und beide nutzen. epic.go: direktes Box.QueueEpic statt RecordEpic bleibt (RecordEpic ist Teil des epic.Refs-Interfaces, eine Lease-Variante dort wäre mehr Fläche als zwei Zeilen).
- **2026-11-22 14:45 · Alexander Sacharov** — critique Runde 2: beide Befunde erledigt (lineEnding in gitignore.go, RecordError.Warning in write.go, CLI-Ausgabe unverändert); nichts Neues.
