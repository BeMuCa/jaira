---
id: 01M7QT1H3CVSC7KKYTJ3NH788Q
title: Nach einem Upgrade bringt jaira unveraenderte Lanes jedes Boards selbst auf den neuen Stand
status: testing
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
updated-at: 2026-11-18T17:39:06Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-72235
claimed-at: 2026-11-18T17:08:08Z
outcome-what: "Der erste jaira-Befehl nach einem Upgrade macht selbst, was 'jaira update --lanes' an Dateien tut: unveraenderte ausgelieferte Lanes aller Boards ersetzen (atomar, unter Lock), den jaira-Block dort erneuern, wo er steht, und das in einer stderr-Zeile sagen. Eigener Stempel lanes-refreshed; der Versions-Stempel bleibt 'jaira update' fuer die Release-Notes. Handgeaenderte Lanes bleiben und werden genannt; --dry-run laesst es aus; ein halber Fehlschlag nennt, was schon geschrieben ist. Der Block traegt die Version, die seinen Text zuletzt geaendert hat; ein aelteres jaira laesst einen neueren Block byte-gleich."
outcome-why: "Nach jedem Release musste in jedem Repo einzeln 'jaira update --lanes' getippt werden (zuletzt 0.5.14 mit der neuen testing-Lane); ohne Richtungspruefung haetten zwei Versionen im Team den Block hin und her geschrieben."
outcome-resolves: "nudgeIfStale + refreshStaleBoard (internal/cli/update.go), g.dryRun (root.go), release.LanesRefreshed, board.RenewInAgentFiles/BlockVersion/versioned, lane.RefreshShipped ueber alle Boards; Tests internal/cli/updateauto_test.go"
review-summary: "none"
round: 5
review-gaps: "build versionLineRe from versionLinePrefix (the literal was written twice); announceInAgentFile computes withoutVersion(s) once instead of three times; left: refreshedLine's four optional parts, each with one source"
test-verdict: "pass: go clean -testcache + go test -race ./... RC=0 at 0ae7456; DoD 1-5 via the named tests in internal/cli/updateauto_test.go; new TestSkippedNamedBoardLaneIsNotOfferedLanesUse runs 'update --lanes --dry-run' text and --json on an edited named-board lane: by-hand hint, no 'lanes use', board key only on it; scratch-board runs of rounds 1-2 unaffected by this diff (printLaneResults/laneResultsJSON only)"
hunt-verdict: "pass: bughunt done at 0ae7456, 0 change, 0 pre-existing (round 3, after 3+1 change findings fixed in e9601cb and 0ae7456)"
review-verdict: "Der Diff erfuellt alle fuenf DoD-Punkte, Tests in internal/cli, core/lane, core/board, core/release sind gruen, go vet sauber; kein Blocker. Gewichtigster Hinweis ist die Versionsschere: ein aelteres Binary in einem neuen Arbeitsbaum dreht den jaira-Block ungefragt zurueck - die Person sollte entscheiden, ob das vor dem Merge eine Richtungspruefung braucht."
review-round: 4
review-check: "1. go build -ldflags '-X github.com/BeMuCa/jaira/core/release.Current=9.9.9' -o /tmp/jaira-new ./cmd/jaira  2. In einem Testrepo mit dem installierten aelteren jaira 'jaira init' ausfuehren, dann eine Zeile in .jaira/lanes/review.md aendern  3. /tmp/jaira-new list ausfuehren - auf stderr steht eine Zeile 'jaira: brought this board up to 9.9.9: updated lanes ...; left alone, no version jaira shipped: review ...'  4. git diff zeigt geaenderte Lane-Dateien und den jaira-Block in CLAUDE.md, review.md ist byte-gleich  5. /tmp/jaira-new list erneut - nur noch der Hinweis auf 'jaira update', keine Refresh-Zeile  6. /tmp/jaira-new list --dry-run gibt es nicht; stattdessen automatisch: go test ./internal/cli -run 'Refresh|Upgrade|DryRun|NamedBoard' -v - alle PASS"
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
- [x] Ein aelteres jaira auf einem Board mit einem neueren Block (CLAUDE.md/AGENTS.md) laesst die Datei byte-gleich und nennt sie; Test
  proof: internal/cli/updateauto_test.go TestOlderBinaryLeavesANewerBlockByteForByte, TestNewerBinaryWithTheSameTextLeavesTheVersionLine; core/board/announce.go RenewInAgentFiles + BlockVersion

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
- **2026-11-18 17:24 · Alexander Sacharov** — critique pass 3: nur e9601cbe gelesen (Antwort auf die drei bughunt-Befunde). PersistentPreRun im Root fuer --dry-run ist der eine Ort, den jeder Befehl durchlaeuft, wie openStore; benannte Boards im selben RefreshShipped statt einer zweiten Schleife im CLI. Kein Befund.
- **2026-11-18 17:26 · Alexander Sacharov** — bughunt: change: internal/cli/updatesetup.go:122 (+ editedError :64, and the auto-refresh line 'jaira update --lanes names them' in internal/cli/update.go) - since RefreshShipped visits named boards (core/lane/refresh.go:54), a skipped lane of .jaira/boards/<name>/lanes is printed with the hint 'jaira lanes use <id> --force replaces it'; that command writes only into .jaira/lanes (internal/cli/lanes.go:70, lane.ProjectLanesDir), so following the hint overwrites today's board lane of the same id - discarding a person's edit there - and leaves the named board's lane as it was; JSON lane results (updatesetup.go:132) carry no board either; old print line, reached with named-board input only by e9601cbe; ran: temp repo, edited .jaira/lanes/review.md and .jaira/boards/video/lanes/review.md, 'jaira update --lanes --dry-run' printed the same hint for both, 'jaira lanes use review --force' (exit 0) removed the edit from .jaira/lanes/review.md, the named one kept its edit
- **2026-11-18 17:26 · Alexander Sacharov** — bughunt: done at e9601cbe - 1 finding (1 change, 0 pre-existing) in 8 modules, follow-up hunt since 2a5f81af (diff of 5 files, internal/cli/root.go read whole as new); concurrency, error paths, input, permissions checked; files: core/board/announce.go core/lane/refresh.go core/release/notes/NH788Q.md core/release/release.go internal/cli/root.go internal/cli/update.go internal/cli/update_test.go internal/cli/updateauto_test.go
- **2026-11-18 17:27 · Alexander Sacharov** — testing runde 2: bughunt-Befund -> in-progress. 'jaira lanes use' hat kein --board (internal/cli/lanes.go), schreibt nur nach .jaira/lanes. Fix: printLaneResults nennt fuer eine uebersprungene Lane eines benannten Boards keinen lanes-use-Befehl, sondern 'von Hand ersetzen'; editedError sagt, dass lanes use nur today's board betrifft; laneResultsJSON bekommt 'board' fuer benannte Boards.
- **2026-11-18 17:29 · Alexander Sacharov** — critique pass 4: nur 0ae74560 gelesen (Antwort auf den bughunt-Befund runde 2): Hinweis fuer benannte Boards an der einen Stelle, die ihn druckt (printLaneResults), JSON-Feld nur wenn gesetzt. Kein Befund.
- **2026-11-18 17:30 · Alexander Sacharov** — bughunt: done at 0ae74560 - 0 findings in 9 modules, follow-up hunt since e9601cbe (diff of core/release/notes/NH788Q.md internal/cli/updateauto_test.go; internal/cli/updatesetup.go read whole as new); concurrency, error paths, input, permissions checked; files: core/board/announce.go core/lane/refresh.go core/release/notes/NH788Q.md core/release/release.go internal/cli/root.go internal/cli/update.go internal/cli/update_test.go internal/cli/updateauto_test.go internal/cli/updatesetup.go
- **2026-11-18 17:33 · Alexander Sacharov** — worker: review -> signoff. Achtung review-check Schritt 1: '-X github.com/BeMuCa/jaira/core/release.Current=...' greift nicht (Execute ueberschreibt Current mit main.version) - richtig ist go build -ldflags '-X main.version=9.9.9'. Branch feat/nh788q gepusht; die Ticketdatei mit den review-Feldern ist bewusst nicht committet (review aendert keinen Code), sie faehrt mit dem logbook-Commit.
- **2026-11-18 17:35 · Alexander Sacharov** — decision (Teamlead 17:34): ein aelteres Binary schreibt nie zurueck, was ein neueres geschrieben hat - Block wie Lanes. Umsetzung: der Block traegt direkt nach dem Startmarker '<!-- jaira:version X -->'; die Zeile zaehlt beim Vergleich 'aktuell?' nicht mit und wird nur neu geschrieben, wenn sich der Inhalt sonst aendert (sonst bekaeme jedes Repo pro Release einen Einzeilen-Diff). Der Auto-Refresh (RenewInAgentFiles) laesst eine Datei byte-gleich, deren Block-Version neuer als das Binary ist, und nennt sie. Nur der Auto-Pfad: 'jaira update' und Lane-Befehle sind die Tat einer Person und muessen den Block an die Lanes anpassen koennen. Verworfen: Version im Startmarker selbst - aeltere Binaries suchen den Marker woertlich und haengten einen zweiten Block an. Verworfen: Hashliste der Bloecke wie shipped.sha256 - der Block haengt von den Lanes des Boards ab. review-check Schritt 1 muss '-X main.version=9.9.9' heissen; schreibt der naechste frische review neu.
- **2026-11-18 17:38 · Alexander Sacharov** — critique pass 5: nur 7b6fd017 gelesen (Teamlead-Befund Versionsschere). Richtungspruefung sitzt in RenewInAgentFiles, dem einen Auto-Pfad; 'jaira update' bleibt bewusst ausgenommen (Notiz 17:34). Kein Befund.
