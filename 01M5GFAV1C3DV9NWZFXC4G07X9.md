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
updated-at: 2026-10-30T01:56:16Z
updated-by: berkc
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

