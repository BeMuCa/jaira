---
id: 01M121BD400P7CR71GW3BNZERQ
title: "Das Board ist sein Lane-Verzeichnis: nichts wird eingespritzt"
status: done
ready: true
creator: BeMuCa
goal: "Was in .jaira/lanes/ liegt, ist das Board - vollstaendig, in der Reihenfolge der order-Datei, pro Repo und gitignoriert; die eingebauten Lanes sind nur noch ein Angebot im Katalog"
context: |-
  Berks Entscheidung am 27.08.: jeder Nutzer definiert Board und Lane-Reihenfolge pro Repo, gitignoriert, damit auch auf einem geteilten Board jeder sein eigenes Board hat. Eingebaute Lanes gehoeren in den Katalog, koennen hinzugefuegt werden, sind kein Muss.

  Was heute passiert: core/lane/lane.go:398 Load() beginnt mit Builtins() - die zehn Dateien unter core/lane/builtin/*.md, per go:embed in die Binary kompiliert - und legt sie bei JEDEM Befehl auf jedes Board. Eine Projektdatei mit gleicher id ersetzt eine davon, eine mit neuer id kommt dazu, eine fehlende Datei bewirkt nichts. Entfernen kann nur die Datei .jaira/lanes/removed (order.go:22 sagt selbst: sie existiert, weil Load immer einspritzt). Deshalb: jaira init mit einer Auswahl von 4 Lanes zeigt beim naechsten Befehl 10 (QA3GN1), und lanes remove muss erst alle 10 als Dateien schreiben, um eine loeschen zu koennen (4DQPMS).

  Zielbild:
  1. Load(root): hat .jaira/lanes/ Lane-Dateien, sind genau die das Board. Keine Builtins dazu. Hat es keine, aber .jaira/ existiert: Default-Board (~/.jaira/default-board.md, aufgeloest gegen Builtins + Katalog; ohne die Datei: die Builtins) nach .jaira/lanes/ kopieren und order schreiben - einmal, mit Meldung. Das ist die einzige Stelle, an der ein Lesen schreibt; Berk hat sie so gewollt (Punkt 2 seiner Liste vom 27.08.).
  2. init schreibt dieselbe Auswahl sofort.
  3. removed wird weder gelesen noch geschrieben; lanes remove loescht die Datei und den order-Eintrag. Differs/Materialise/MaterialiseWorkingSet fallen zu einer Funktion zusammen.
  4. order ist die Reihenfolge; after: dient nur noch beim Hinzufuegen als Startposition. QA3GN1 erledigt sich.
  5. Katalog = Builtins (kompiliert) + ~/.jaira/lanes/. lanes add und die +-Spalte im Settings-Screen bieten beides an (Installable existiert).
  6. Migration: ein Verzeichnis mit Lane-Dateien, aber ohne order-Datei, ist ein Board von vorher (Dateien wurden zu den eingespritzten Builtins DAZU gelegt). Einmal die fehlenden Builtins daneben materialisieren, order schreiben, sagen was passiert ist. Berks req-Board hat order und removed (von lanes remove), braucht das nicht.
  7. Teammate auf geteiltem Board: .jaira/lanes/ ist ignoriert, er bekommt beim ersten Befehl SEIN Default-Board. Tickets in Lanes, die er nicht hat, erscheinen als '? lane' (Passthrough, lane.go:762, nur lesbar), bis er sie aus .jaira/shared/ adoptiert. Das ist der Preis pro-Nutzer-Boards und gewollt.
  8. Load("") ohne Projekt (Tests, lanes template) liefert weiter Builtins + Katalog: das Angebot.

  Betroffen: core/lane (Load, order, defaultboard, Installable), internal/cli/tickets.go init, lanes.go, core/lane-Tests (viele setzen die Einspritzung voraus), docs (README Abschnitt Lanes, docs/COMMANDS.md, Block-Text).
definition-of-done: "Ein Board mit vier Lane-Dateien in .jaira/lanes/ zeigt vier Lanes, bei jedem Befehl; jaira init mit einem Default-Board von vier Lanes schreibt vier Dateien plus order und der naechste Befehl zeigt vier; die Datei removed wird nirgends mehr gelesen oder geschrieben (grep leer); ein Verzeichnis mit Lanes aber ohne order wird einmal migriert und meldet es; lanes add bietet Builtins und Katalog an; lanes remove loescht nur die Datei und den order-Eintrag; go test ./... -race gruen"
blocked-by: []
commits:
  - 743737f7a731e097c7ab26396293086fae83d3e6
  - 61789a2bdc42d70a3338b9645a3c09b24b4cdf50
  - 6c5d81b28e18049c24bafd21a1d639c263474a32
created-at: 2026-08-27T16:37:48Z
updated-at: 2026-09-03T15:48:28Z
claimed-by: EE-3NX6GL3-3987890
claimed-at: 2026-08-27T20:41:37Z
updated-by: BeMuCa
assignee: BeMuCa
outcome-what: "Label 'overrides' -> 'changed' im Settings-Screen."
outcome-why: "Kein Builtin mehr auf dem Board, das ueberschrieben wuerde; das Feld markiert eine Abweichung von der gelieferten Lane gleichen Namens."
outcome-resolves: "Kritikpunkt erledigt, TUI-Tests gruen, Binary gebaut."
review-summary: "Zwei Commits. 743737f: Load liest bei einem Board nur die Dateien, setUp schreibt sie beim ersten Oeffnen (Default-Board oder Builtins) samt order, migrateLegacy holt alte Verzeichnisse herueber (removed oder fehlende order als Signal), removed/MaterialiseWorkingSet/Differs/Override-Warnung/u-Taste entfernt, lanes use laedt erst das Board und haengt Neues an order, n schreibt ins Board-Verzeichnis, Anker-Warnungen nur ohne order-Datei; init braucht nur noch Load; README/COMMANDS/Hilfen/Paketkommentar. Zweiter Commit: Label 'changed'."
review-gaps: "Entfernt, im Feature-Commit: removed samt LoadRemoved/SaveRemoved, MaterialiseWorkingSet, Differs, die alte Einspritzung in Load, die u-Taste plus use() im Settings-Screen, drei Tests dazu, die Override-Warnung. Gesucht: zweite Stellen, die Lanes einspritzen oder removed lesen - grep Builtins() zeigt nur Load (Angebot), setUp, migrateLegacy, Installable; grep removed zeigt nur migrateLegacy und die Konstante. Stehen gelassen: das Feld Lane.Overrides unter altem Namen (Umbenennung = Churn durch tui und Tests, kein Leser gewinnt); effectiveOrder (order-Datei kann von Hand fehlen); die Konstante removedFileName fuer die Migration; gofmt-Eintrag internal/cli/tickets.go. Kosten: lanes use laedt das Board zweimal (einmal fuer setUp, einmal in refreshAgentNote) - ein Verzeichnis-Glob, nicht angefasst. Nicht gebaut: Aufraeumen der D-03-Formulierungen in ProjectLanesDir/ProjectLanesActive-Kommentaren - inhaltlich noch richtig."
review-check: |-
  1. Wegwerf-Repo: git init && jaira init - Ausgabe 'Wrote the 10 built-in lanes into .../.jaira/lanes'; ls .jaira/lanes zeigt 10 .md-Dateien und order, keine Datei removed.
  2. jaira lanes remove brainstorm - danach 9 Lanes, brainstorm.md ist weg, es gibt weiterhin keine Datei removed.
  3. rm .jaira/lanes/signoff.md und jaira lanes: signoff ist weg, eine Warnung sagt, dass order eine nicht installierte Lane nennt. Zeile aus order loeschen: keine Warnung.
  4. jaira lanes add brainstorm - wieder da, als letzte Spalte.
  5. Default-Board: echo -e '---\nlanes: [backlog, todo, done]\n---' > ~/.jaira/default-board.md (oder JAIRA_DEFAULT_BOARD setzen), neues Repo, jaira init, jaira lanes - genau drei Lanes, auch beim zweiten Aufruf. Danach die Datei wieder entfernen.
  6. Dein req-Board: der erste jaira-Befehl dort meldet einmal 'migrated .jaira/lanes: ...; wrote the 8 shipped lane(s) ...; deleted the obsolete removed file'; ls .jaira/lanes zeigt danach 11 Dateien plus order; das Board sieht aus wie vorher (signoff fehlt weiterhin).
  7. go clean -testcache && go test ./core/lane/ ./internal/cli/ ./internal/tui/ -race - alle ok.
review-verdict: "Angenommen; dieselbe Sitzung, kein zweites Modell. Geprueft, nicht behauptet: go test ./... -race mit leerem Cache gruen nach beiden Commits; vier Live-Szenarien mit dem gebauten Binary (frisches init, remove, Handloeschung, Legacy-Verzeichnis mit removed) und dieses Board selbst (migriert, removed geloescht, zweiter Lauf still). Verhaltensaenderungen, alle gewollt: ein Board ohne Lane-Dateien bekommt sie beim ersten Befehl geschrieben (ein Lesen schreibt - Berks Punkt 2); ein von Hand geloeschtes Lane-File nimmt die Lane vom Board; eine Datei mit Builtin-Namen warnt nicht mehr, nur ein fallengelassener Schutz; lanes use ohne --force lehnt Lanes ab, die das Board schon hat. Eine Folge zum Wissen: Berks req-Board wird beim naechsten jaira-Befehl migriert (8 Dateien geschrieben) - gewollt, gemeldet, einmal."
---

# Das Board ist sein Lane-Verzeichnis: nichts wird eingespritzt

## Definition of Done

- [x] Ein Board mit vier Lane-Dateien in .jaira/lanes/ zeigt vier Lanes, bei jedem Befehl; jaira init mit einem Default-Board von vier Lanes schreibt vier Dateien plus order und der naechste Befehl zeigt vier; die Datei removed wird nirgends mehr gelesen oder geschrieben (grep leer); ein Verzeichnis mit Lanes aber ohne order wird einmal migriert und meldet es; lanes add bietet Builtins und Katalog an; lanes remove loescht nur die Datei und den order-Eintrag; go test ./... -race gruen
  proof: core/lane/lane.go Load/setUp/migrateLegacy; TestABoardIsItsLaneDirectory, TestFirstLoadUsesTheDefaultBoard, TestLegacyDirectoryWithRemovedFileIsMigrated; live: fresh init 10 Dateien+order, lanes remove loescht nur die Datei, Handloeschung bleibt weg, Legacy-Board migriert, dieses Board: removed geloescht, zweiter Lauf still; grep LoadRemoved/SaveRemoved leer

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Load(root): aktives Lane-Verzeichnis = das Board, keine Builtins dazu; root ohne .jaira/ liefert weiter Builtins+Katalog (das Angebot)
- [x] Erstes Laden eines Boards ohne Lane-Dateien: Default-Board (oder Builtins) nach .jaira/lanes/ schreiben plus order, einmal, mit Meldung
- [x] Migration alter Boards: kein order ODER removed vorhanden -> fehlende Builtins aus order (oder alle) daneben schreiben, removed-ids auslassen, order schreiben, removed loeschen, melden
- [x] removed weder lesen noch schreiben; Remove loescht Datei + order-Eintrag; MaterialiseWorkingSet und Differs entfallen; init materialisiert immer
- [x] Override-Warnung streichen (kein Builtin auf dem Board, das ueberschrieben wuerde); Schutz-Verlust-Warnung gegen die gelieferte Lane gleichen Namens behalten
- [x] Tests anpassen (core/lane, defaultboard, order, cli lanes, tui lanes), Docs (README Lanes, COMMANDS lanes remove, Hilfetexte), Block-Text pruefen

## Progress
- **2026-08-27 21:00 · BeMuCa** — critique: Label 'overrides' -> 'changed' in internal/tui/lanes.go:682. Feld Overrides selbst bleibt (Umbenennung waere Churn ohne Leser).
