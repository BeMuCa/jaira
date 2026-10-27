---
id: 01M5659DPY9FK72XVMZFY2C9QN
title: Ein zweites Board laesst sich mit einem Befehl aus einem Board-Typ anlegen statt Lane-Dateien von Hand zu kopieren
status: critique
ready: true
creator: berkc
assignee: berkc
goal: "'jaira boards create <name> --type <typ>' legt ein Board mit den Lanes und der Beschreibung eines Board-Typs an (z.B. code, research, content), und jaira liefert ein paar Typen mit."
context: |-
  Berk am 18.10.: 'Im Moment haben wir nur 1 Board. Wir wollen mehrere Board-Typen ermoeglichen.'
  Heute: mehrere Boards pro Repo gibt es seit 0.5.3 (Ticket S0EG0M, done): Lane-Dateien und order-Datei nach .jaira/boards/<name>/lanes/, 'jaira create --board <name>', 'jaira boards' listet sie, 'b' im Board wechselt. Berk hat 0.5.2 installiert und sieht das deshalb noch nicht.
  Was fehlt: kein Befehl legt ein Board an ('jaira boards' hat keine Unterbefehle, nur Auflisten). Man muss die Dateien von Hand anlegen. Es gibt keine Board-Typen: nur die eingebauten Code-Lanes, 'jaira lanes default' und 'jaira lanes template --board'.
  Beispiel, dass es gebraucht wird: Berks Repo videoeditor hat ein Board mit eigenen Lanes (Reel Research, Creative Direction, Shoot, Edit ...), von Hand gebaut.
  Verwandt: PP24H1 (Interview-Skill soll Boards passend anlegen - braucht diesen Befehl), YA7KW7 (eine Bitte ueber zwei Boards), BGZ0SW (Ticket-Schema aus Lanes).
  Offen fuer brainstorm: welche Typen mitgeliefert werden; ob ein Typ nur Lanes + Beschreibung ist oder auch Tags/Milestone-Vorlagen; ob eigene Typen im Katalog (~/.jaira) liegen koennen. Kleiner als paca bleiben.
definition-of-done: "'jaira boards create <name> --type <typ>' schreibt .jaira/boards/<name>/lanes/ mit Lanes, order und description des Typs; ein vorhandenes Board wird nicht ueberschrieben (exit 3)"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-18T00:06:22Z
updated-at: 2026-10-27T14:47:43Z
updated-by: berkc
claimed-by: berk-ca-728718
claimed-at: 2026-10-27T14:32:43Z
outcome-what: "Neue Befehle 'jaira boards create <name> --description <...> [--type <typ>]' und 'jaira boards types' (beide mit --json). Board-Typen liegen eingebettet unter core/lane/boardtypes/<typ>/ (description, order, eigene Lane-Dateien; IDs ohne eigene Datei kommen aus den eingebauten Lanes). Mitgeliefert: 'code' = die eingebauten Default-Lanes."
outcome-why: "Ein zweites Board musste man bisher von Hand aus Lane-Dateien zusammenkopieren; Agenten brauchen einen nicht-interaktiven Befehl, und Berk will spaeter eigene Typen (ideas, builder) ausliefern koennen, ohne Go-Code zu schreiben."
outcome-resolves: "Y2C9QN: Board per Befehl aus einem Typ anlegen; vorhandenes Board wird nie ueberschrieben (exit 3), fehlende Beschreibung exit 2."
review-summary: "core/lane/boardtypes.go typeLanes: returns parallel []*Lane and [][]byte, but only l.ID is used and it always equals the id from order - drop the second Parse and the []*Lane slice, return the order ids and the bytes. core/lane/boardtypes.go mustSub: helper with one caller - embed boardtypes and address files as boardtypes/<type>/..., test MapFS keys likewise. core/lane/boardtypes.go boardType: reads the literal \"description\" beside orderFileName - use descriptionFileName from boards.go."
---

# Ein zweites Board laesst sich mit einem Befehl aus einem Board-Typ anlegen statt Lane-Dateien von Hand zu kopieren

## Definition of Done

- [x] 'jaira boards create <name> --type <typ>' schreibt .jaira/boards/<name>/lanes/ mit Lanes, order und description des Typs; ein vorhandenes Board wird nicht ueberschrieben (exit 3)
  proof: core/lane/boardtypes.go CreateBoard, internal/cli/boards.go:94 newBoardsCreateCmd; TestCreateBoardFromATypeWithLanesOfItsOwn, TestBoardsCreate (exit 3 for video and default)
- [x] 'jaira boards types' listet die mitgelieferten Typen mit einer Zeile Beschreibung
  proof: internal/cli/boards.go:149 newBoardsTypesCmd, core/lane/boardtypes/code/description; TestBoardsCreate (boards types)
- [x] Ein Board ohne --type bekommt die eingebauten Code-Lanes
  proof: --type defaults to lane.DefaultBoardType=code (core/lane/boardtypes/code/order); TestCodeTypeIsTheBuiltinDefaultLanes, TestBoardsCreate (create without --type)
- [x] Tests; Release-Note
  proof: core/lane/boardtypes_test.go (3 tests), internal/cli/boards_test.go TestBoardsCreate; core/release/notes/Y2C9QN.md; go vet ./... and go test ./... green
- [x] 'boards create' verlangt eine Beschreibung (--description), die sagt, welche Tickets auf dieses Board gehoeren; ohne sie exit 2 mit Hinweis, wofuer der Teamlead sie braucht
  proof: internal/cli/boards.go:111 missing_description exit 2; TestBoardsCreate

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/lane: Board-Typen eingebettet unter core/lane/boardtypes/<typ>/ (description, order, eigene Lanes); BoardTypes() und CreateBoard(root, name, typ, description)
- [x] cli: 'jaira boards create <name> [--type] --description' und 'jaira boards types', beide mit --json; exit 2 ohne Beschreibung / unbekannter Typ, exit 3 bei vorhandenem Board
- [x] Tests (core/lane und internal/cli), Release-Note core/release/notes/Y2C9QN.md, docs falls Befehlsliste gepflegt

## Progress
- **2026-10-27 14:32 · berkc** — Korrektur 27.10.2026 (berkc): Der Kontext beschreibt 'b' so, wie es heute ist (b schaltet zum naechsten Board, internal/tui/model.go:1705). Gewollt ist (KVMG0K): b oeffnet ein Popup mit allen Boards, dort waehlt man mit Leertaste die freigegebenen (aktiven) Boards, aus denen der Teamlead waehlen darf; Shift+hoch/runter wechselt zwischen den aktiven Boards. Eine Einstellung 'aktive Boards' gibt es heute nicht (nur default-board in internal/tui/settings.go). Fuer dieses Ticket aendert das nichts: es geht nur ums Anlegen von Boards per Befehl.
- **2026-10-27 14:34 · berkc** — Entscheidungen (Worker, nach 'kleiner als paca'): (1) Mitgeliefert wird nur der Typ 'code' = die eingebauten Default-Lanes (die Lanes mit default-board != false), also dasselbe wie ohne --type. 'ideas' und 'builder' liefert Berk selbst spaeter; das Format muss sie nur tragen koennen. (2) Ein Typ ist NUR Lanes + order + description - keine Tags, keine Milestone-Vorlagen (niemand hat sie verlangt). (3) Keine eigenen Typen in ~/.jaira in diesem Ticket: ein zweiter Quellort ist ein eigenes Feature; das Format (ein Verzeichnis wie ein Board-lanes/-Verzeichnis) laesst das spaeter mit einer Zeile zu. (4) Format: core/lane/boardtypes/<typ>/ mit 'description' (eine Zeile, was der Typ ist; nur fuer 'boards types'), 'order' (Lane-IDs) und optional eigenen <id>.md Lane-Dateien; eine ID in order ohne eigene Datei kommt aus den eingebauten Lanes. So ist ideas = eigene Lanes (ideate, pick ...), builder = eingebaute Code-Lanes plus eigene. (5) DoD1 'description des Typs' vs DoD5 '--description Pflicht': die Board-description-Datei bekommt den --description-Text (welche Tickets aufs Board gehoeren); die Typ-Beschreibung steht nur in 'boards types'. (6) 'code' ignoriert ~/.jaira/default-board.md bewusst: DoD3 sagt eingebaute Code-Lanes, und ein Befehl, den ein Agent ausfuehrt, soll auf jeder Maschine dasselbe Board bauen.
- **2026-10-27 14:43 · berkc** — in-progress: Agentenblock (core/board/announce.go) nennt jetzt 'boards create'/'boards types'; CLAUDE.md und AGENTS.md des Repos von Hand gleichgezogen wie bei 8NHD8Q/A0C081 ueblich (jaira validate meldet keinen Block-Mismatch). Format fuer ideas/builder ist mit einem Test belegt (TestCreateBoardFromATypeWithLanesOfItsOwn legt per fstest.MapFS einen ideas-Typ mit eigenen Lanes an) - einen neuen Typ ausliefern heisst: Verzeichnis unter core/lane/boardtypes/ anlegen, kein Go-Code.
