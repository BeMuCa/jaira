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
updated-at: 2026-11-18T17:14:49Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-72235
claimed-at: 2026-11-18T17:08:08Z
outcome-what: "Der erste jaira-Befehl in einem Repo, dessen Stempel aelter ist als das Binary, macht jetzt selbst, was 'jaira update --lanes' tut: unveraenderte ausgelieferte Lanes ersetzen (atomar), den jaira-Block dort erneuern, wo er steht, stempeln, und das in einer stderr-Zeile sagen. Handgeaenderte Lanes bleiben und werden genannt. Scheitert es, kommt wie bisher der Hinweis."
outcome-why: "Nach jedem Release musste in jedem Repo einzeln 'jaira update --lanes' getippt werden (zuletzt 0.5.14 mit der neuen testing-Lane)."
outcome-resolves: "nudgeIfStale -> refreshStaleBoard (internal/cli/update.go), board.RenewInAgentFiles, lane.RefreshShipped via ticket.WriteAtomic; Tests in internal/cli/updateauto_test.go"
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
  proof: internal/cli/updateauto_test.go TestConcurrentRefreshesWriteNoHalfLane (-race); core/lane/refresh.go ticket.WriteAtomic + Lock board-refresh
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
