---
id: 01M7QT1H3CVSC7KKYTJ3NH788Q
title: Nach einem Upgrade bringt jaira unveraenderte Lanes jedes Boards selbst auf den neuen Stand
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Nach 'jaira self upgrade' muss niemand mehr in jedem Repo 'jaira update --lanes' tippen: Lanes, die niemand von Hand geaendert hat, zieht jaira beim ersten Aufruf im Repo selbst nach und sagt das in einer Zeile; handgeaenderte Lanes bleiben unangetastet."
context: |-
  Was heute falsch ist: nach jedem Release muss Alex in jedem Repo mit Board 'jaira update --lanes' aufrufen. Bei 0.5.14 (neue testing-Lane mit hunt-verdict) wieder: jedes Repo einzeln.
  Bis dahin druckt jeder Befehl nur 'this board was set up by an older version of jaira — run jaira update' (internal/cli/update.go:55).
  Der Kern existiert schon: lane.RefreshShipped (core/lane/refresh.go:30) ersetzt jede Lane-Datei, die genau einer frueher ausgelieferten Fassung gleicht (Hashes in core/shipped/shipped.sha256), und laesst alle anderen als 'skipped' stehen. Aufgerufen wird er heute nur von 'jaira update --lanes' (internal/cli/update.go:126).
  Fehlt: diesen Aufruf automatisch machen, wenn das Board aelter als das Binary ist - dort, wo heute die Hinweiszeile kommt.
  Bewusst NICHT: Lanes global aus dem Binary oder ~/.jaira lesen. Lanes liegen im Repo, damit das ganze Team denselben Ablauf hat; das bleibt so. Die geaenderten Dateien sind normale Aenderungen im Arbeitsbaum, die der Mensch committet.
  Dasselbe gilt fuer den jaira-Block in CLAUDE.md/AGENTS.md, den 'jaira update' mit den Lanes neu schreibt.
  Ersetzt 0S1A7M (lanes use --all), das 'jaira update --lanes' schon erledigt hat.
  Verwandt: 28C0QV (neue Lane auch im persoenlichen Katalog ~/.jaira/lanes).
definition-of-done: "Ein Board mit einer unveraenderten Lane einer aelteren Fassung: der erste jaira-Befehl im Repo nach dem Upgrade ersetzt sie und meldet es in einer Zeile auf stderr; Test"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-11-18T17:07:33Z
updated-at: 2026-11-18T17:24:40Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-72235
claimed-at: 2026-11-18T17:08:08Z
outcome-what: "Der erste jaira-Befehl nach einem Upgrade macht selbst, was 'jaira update --lanes' an Dateien tut: unveraenderte ausgelieferte Lanes aller Boards ersetzen (atomar, unter Lock), den jaira-Block dort erneuern, wo er steht, und das in einer stderr-Zeile sagen. Eigener Stempel lanes-refreshed; der Versions-Stempel bleibt 'jaira update' fuer die Release-Notes. Handgeaenderte Lanes bleiben und werden genannt; ein Befehl mit --dry-run laesst es aus; ein halber Fehlschlag nennt, was schon geschrieben ist."
outcome-why: "Nach jedem Release musste in jedem Repo einzeln 'jaira update --lanes' getippt werden (zuletzt 0.5.14 mit der neuen testing-Lane)."
outcome-resolves: "nudgeIfStale + refreshStaleBoard (internal/cli/update.go), g.dryRun via PersistentPreRun (root.go), release.LanesRefreshed, board.RenewInAgentFiles, lane.RefreshShipped ueber alle Boards via ticket.WriteAtomic; Tests internal/cli/updateauto_test.go"
review-summary: "none"
round: 3
review-gaps: "none — considered: refreshStaleBoard vs the update command body (internal/cli/update.go RunE) share RefreshShipped+Load+announce but differ in stamp, agent-file choice and output, left separate; refreshedLine vs printLaneResults (updatesetup.go) are different outputs; per-command cost is one extra stamp read"
test-verdict: "pass: go build/vet + go clean -testcache + go test -race ./... RC=0; DoD 1-5 verified by running the named tests and on a scratch board (binaries built with -X main.version=9.0.0 and 9.0.1): unedited testing lane replaced on the first 'list', one stderr line, second command back to the old hint, 'jaira update' still shows the notes; edited lane byte-identical (sha1 same) and named with 'jaira update --lanes'; lock dir blocked -> rc 0, old hint, lane untouched; 4 parallel 'list' -> exactly one refresh line, lane replaced, no .jaira-tmp leftovers"
hunt-verdict: "fail: bughunt done at 2a5f81a, 3 change findings (partial failure silent, --dry-run writes, named boards skipped but stamped)"
---

# Nach einem Upgrade bringt jaira unveraenderte Lanes jedes Boards selbst auf den neuen Stand

## Definition of Done

- [x] Ein Board mit einer unveraenderten Lane einer aelteren Fassung: der erste jaira-Befehl im Repo nach dem Upgrade ersetzt sie und meldet es in einer Zeile auf stderr; Test
  proof: internal/cli/updateauto_test.go TestFirstCommandAfterUpgradeReplacesAnUneditedLane; internal/cli/update.go refreshStaleBoard
- [x] Eine handgeaenderte Lane bleibt byte-gleich, die Meldung nennt sie mit 'jaira update --lanes' als Hinweis; Test
  proof: internal/cli/updateauto_test.go TestAutomaticRefreshLeavesAnEditedLaneByteForByte; update.go refreshedLine
- [x] Der jaira-Block in CLAUDE.md/AGENTS.md wird dabei mit erneuert, wie bei 'jaira update --lanes'; Test
  proof: internal/cli/updateauto_test.go TestAutomaticRefreshRenewsTheJairaSection, TestAutomaticRefreshCreatesNoAgentFile; core/board/announce.go RenewInAgentFiles
- [x] Zwei jaira-Prozesse gleichzeitig im selben Repo schreiben keine halbe Lane-Datei; Test
  proof: internal/cli/updateauto_test.go TestConcurrentRefreshesWriteNoHalfLane (-race); core/lane/refresh.go ticket.WriteAtomic; internal/cli/update.go refreshStaleBoard Lock board-refresh
- [x] Release-Note-Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/NH788Q.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-18 17:08 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/6601.sock
- **2026-11-18 17:09 · Alexander Sacharov** — decision: keine offene Frage an Alex. Ausloeser wie vom Teamlead vorgegeben: nudgeIfStale (internal/cli/update.go) - Stempel != Binary -> statt der Hinweiszeile laeuft dasselbe wie 'jaira update --lanes' fuer das Board: lane.RefreshShipped, jaira-Block via board.AnnounceInAgentFiles, dann release.Stamp. Gestempelt wird, damit der Hinweis verschwindet und nicht jeder Befehl Lanes laedt und hasht (Startzeit). Preis: 'jaira update' zeigt danach 'nothing has changed' statt der Release-Notes; die TUI zeigt sie weiter (notes-seen ist ein eigener Stempel). Rollen bleiben bei nudgeSkillsIfStale. Ausgabe nur eine stderr-Zeile, und nur wenn etwas ersetzt, uebersprungen oder der Block neu geschrieben wurde; sonst still stempeln. Dev-Build: nichts (wie bisher). Nebenlaeufigkeit: store.Lock('board-refresh') + Stempel nach dem Lock erneut pruefen; RefreshShipped schreibt via ticket.WriteAtomic. Fehler beim Auto-Refresh brechen keinen Befehl ab: dann wie bisher die Hinweiszeile.
- **2026-11-18 17:14 · Alexander Sacharov** — decision: der jaira-Block wird nur dort erneuert, wo eine Datei ihn schon traegt (board.RenewInAgentFiles), nicht wie 'jaira update' auch neu angelegt: ein ungefragter Befehl soll in einem Repo ohne CLAUDE.md/AGENTS.md keine solche Datei erzeugen. Gefunden: lane.Load ueberspringt eine unlesbare Lane-Datei still, ein chmod 000 auf die Lane provoziert also keinen Fehler - der Fallback-Test blockiert stattdessen das Lock-Verzeichnis.
- **2026-11-18 17:15 · Alexander Sacharov** — critique: RenewInAgentFiles dupliziert die Marker-Suche aus chosenAgentFiles (core/board/announce.go) - in eine Funktion filesWithBlock ziehen, chosenAgentFiles faellt bei leerem Ergebnis auf agentFiles zurueck. Sonst nichts: Ort (cli, neben dem update-Befehl), Lock+Stempel-Recheck, Fallback auf den alten Hinweis stehen.
- **2026-11-18 17:15 · Alexander Sacharov** — decision (Teamlead, 17:15): der Auto-Refresh darf den Versions-Stempel NICHT setzen - 'jaira update' ist der einzige CLI-Ort fuer die Release-Notes inkl. '- !'-Zeilen. Also eigener Stempel <state>/lanes-refreshed; der Versions-Stempel und der Hinweis auf 'jaira update' bleiben bis dahin. Ersetzt meine Notiz zum Stempeln.
- **2026-11-18 17:17 · Alexander Sacharov** — critique pass 2: filesWithBlock beantwortet den Befund; der eigene Stempel lanes-refreshed (Teamlead-Entscheidung) bricht keinen DoD-Punkt - die Zeile bleibt eine. Nichts weiter.
- **2026-11-18 17:20 · Alexander Sacharov** — bughunt: change: internal/cli/update.go:61 (refreshStaleBoard, update.go:97-111) - a failure halfway is silent and loses the report: RefreshShipped has already replaced the lane files when RenewInAgentFiles fails (e.g. a read-only CLAUDE.md), the error is dropped at :61, and only the old 'run jaira update' hint prints; every later command retries, now finds the lanes 'unchanged', and never names them - lane files changed in the working tree with nobody told; ran: scratch test (newStaleSetup + upgrade, chmod 0444 CLAUDE.md, 'jaira list' twice): lane file replaced, no stderr line names it
- **2026-11-18 17:20 · Alexander Sacharov** — bughunt: change: internal/cli/update.go:60 (via openStore, internal/cli/root.go:270) - commands that promise to write nothing now write: the first 'jaira move <id> --to <lane> --dry-run' (also 'jaira sync --dry-run') after an upgrade replaces lane files and rewrites CLAUDE.md, because the refresh runs in openStore and only 'jaira update' is exempted (updateInProgress); ran: scratch test, 'move nosuchid --to todo --dry-run' after upgrade -> lane file rewritten, stderr 'brought this board up to 1.0.1: updated lanes testing; rewrote the jaira section in CLAUDE.md'
- **2026-11-18 17:20 · Alexander Sacharov** — bughunt: change: internal/cli/update.go:97 -> core/lane/refresh.go:40 - lanes of a named board (.jaira/boards/<name>/lanes) are never refreshed, yet the lanes-refreshed stamp is written, so no later command retries and the release note's 'every lane nobody edited' is false there; the default-board-only scope of RefreshShipped is older (5c087b6, also 'update --lanes'), marked change because the automatic path and its promise are this ticket's; ran: scratch test with an old shipped testing.md under .jaira/boards/video/lanes, 'jaira boards' after upgrade -> default lane updated, video lane left at old bytes, second command prints no refresh
- **2026-11-18 17:20 · Alexander Sacharov** — bughunt: done at 2a5f81af - 3 findings (3 change, 0 pre-existing) in 7 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/board/announce.go core/lane/refresh.go core/release/notes/NH788Q.md core/release/release.go internal/cli/update.go internal/cli/update_test.go internal/cli/updateauto_test.go
- **2026-11-18 17:21 · Alexander Sacharov** — testing: bughunt fand 3 change-Befunde (Notizen oben) -> zurueck in in-progress. Plan: (1) refreshStaleBoard gibt auch bei Fehler die Zeile fuer das schon Geschriebene zurueck, nudge druckt sie plus eine Fehlerzeile mit 'jaira update --lanes'; (2) ein Befehl mit --dry-run=true laesst den Auto-Refresh aus (root PersistentPreRun setzt ein Flag, openStore hat kein cmd); (3) RefreshShipped besucht auch die Lanes benannter Boards (Boards/LoadBoard), PreviewRefresh schluesselt dann nach Pfad statt nach Lane-ID.
- **2026-11-18 17:24 · Alexander Sacharov** — Runde 2 in-progress: alle drei bughunt-Befunde behoben. Gefunden: es gibt keinen PersistentPreRun im Root - der neue in root.go ist der einzige; ein Unterbefehl mit eigenem PersistentPreRun wuerde ihn ueberdecken (cobra ruft nur den naechsten). Mutationstest: ohne die g.dryRun-Abfrage faellt TestDryRunCommandDoesNotRefresh. 'jaira update --lanes' erfasst jetzt auch benannte Boards (gleicher Code), daher zweite Release-Note-Zeile.
