---
id: 01M5GFAV1C3DV9NWZFXC4G07X9
title: "Boards wie Lanes behandeln: von Hand anlegen, Lanes waehlen, beschreiben, teilen und global speichern"
status: ideas
ready: true
creator: berkc
assignee: berkc
goal: "Ein Board laesst sich im TUI anlegen (Name, Beschreibung, Lanes aus dem Katalog), mit dem Team teilen, ein geteiltes Board in den eigenen Katalog uebernehmen, und jedes Repo sieht die Boards aus dem persoenlichen Katalog."
context: |-
  Wunsch von berkc am 22.10.2026: 'Regarding the Boards we need to handle them similar to the lanes, making it possible to create a board manually and pick lanes to include, write board description and also to share boards with the team or save shared boards. The boards are also saved globally, like for the tickets.'
  Heute: ein Board ist ein Ordner .jaira/boards/<name>/lanes/ mit description-Datei, nur im Repo (S0EG0M). Anlegen geht nur von Hand.
  Vorbild Lanes: persoenlicher Katalog ~/.jaira/lanes (lane.UserLanesDir, core/lane/lane.go:466), publish nach .jaira/shared/<du>/, adopt zurueck (internal/tui/lanes.go:484, :640). Fuer Boards gaebe es analog ~/.jaira/boards und .jaira/shared/<du>/boards/.
  Verwandt: Y2C9QN (Board aus Board-Typ mit einem Befehl, Beschreibung Pflicht), 57TDND (aktive Boards), 28C0QV (neue Lane auch in den Katalog), Ticket zum Board-Fenster auf b.
  Offen: 'saved globally, like for the tickets' - Tickets liegen auf Refs im Repo, nicht global; gemeint ist vermutlich der persoenliche Katalog wie bei Lanes. Bei berkc nachfragen.
definition-of-done: "Im TUI laesst sich ein Board anlegen: Name, Beschreibung (Pflicht), Lanes aus Katalog und Board ankreuzen; Test"
tags:
  - tui
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-22T00:14:19Z
updated-at: 2026-10-30T01:59:36Z
updated-by: berkc
claimed-by: berk-ca-1397571
claimed-at: 2026-10-30T01:57:35Z
---

# Boards wie Lanes behandeln: von Hand anlegen, Lanes waehlen, beschreiben, teilen und global speichern

## Definition of Done

- [ ] Im TUI laesst sich ein Board anlegen: Name, Beschreibung (Pflicht), Lanes aus Katalog und Board ankreuzen; Test
- [ ] Ein Board landet auch im persoenlichen Katalog und ist in jedem Repo waehlbar; Test mit JAIRA_HOME
- [ ] publish/adopt fuer Boards wie fuer Lanes; Test
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-30 01:56 · berkc** — berkc 30.10.2026 (Frage 3/13): Boards und Default-Lanes sollen an einer Stelle global liegen, damit ein 'jaira update' sie fuer alle Sessions und Repos auf einmal aktualisiert und Repo-eigene Lanes nicht ueberschrieben werden. 'global wie Tickets' meint den persoenlichen Katalog ~/.jaira, nicht die Refs. Baut auf Y2C9QN auf (boards create/types, Branch feat/Y2C9QN-boards-create, noch nicht gemergt).
- **2026-10-30 01:59 · berkc** — Blocked-by Y2C9QN (reasoning): this work is built on branch feat/Y2C9QN-boards-create (jaira boards create/types, core/lane/boardtypes.go), not yet on master. Branch feat/4G07X9-boards-catalogue is based on origin/feat/Y2C9QN-boards-create; its PR must land after Y2C9QN's. Not set as blocked_by field so the ticket can run its lanes now.
- **2026-10-30 01:59 · berkc** — Brainstorm. What the code shows: a board is .jaira/boards/<name>/lanes/ (description, order, <id>.md) in one repo only (core/lane/boards.go). Y2C9QN adds board types: embedded dirs core/lane/boardtypes/<type>/ of the same shape, 'jaira boards create --type'. Lanes have a personal catalogue ~/.jaira/lanes (lane.UserLanesDir, honours JAIRA_HOME), 'jaira lanes publish' -> .jaira/shared/<you>/<id>.md, 'jaira lanes adopt <path>' -> catalogue (core/lane/share.go). Nothing for boards. Options: (A) catalogue board = a board-type-shaped dir ~/.jaira/boards/<name>/ with full lane copies; 'boards types' and 'boards create --type' read shipped types plus the catalogue; create also saves into the catalogue; publish/adopt via .jaira/shared/<you>/boards/<name>/. Reuses typeLanes/CreateBoard; repo boards stay copies, never overwritten. (B) repo boards reference the catalogue live (no copy): one edit reaches every repo, but a clone without the person's ~/.jaira shows no lanes - breaks 'teammate clones and sees same state'. (C) jaira update syncs catalogue boards into every repo board: needs shipped-version tracking (master's core/shipped, not on this base) and overwrites repo-owned lanes - against the 30.10 note. Choose A.
