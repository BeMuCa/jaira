---
id: 01M58G2A0YCD31RRC3GZ28C0QV
title: "Eine neu angelegte Lane landet nur in einem Board, nicht im persoenlichen Katalog ~/.jaira/lanes"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Eine Lane, die ich in einem Repo anlege, kann ich in jedem meinem anderen Boards in den Settings auswaehlen; eine geteilte Lane eines Kollegen kann ich dort ebenfalls uebernehmen."
context: |-
  Heute: newLane() (internal/tui/lanes.go:564) schreibt die neue Lane nach lane.ProjectLanesDir (.jaira/lanes/ des Repos). Andere Boards sehen sie nie.
  Gemeldet von berkc: Lanes aus dem videoeditor-Repo fehlen in den Settings der anderen Boards.
  Es gibt schon einen persoenlichen Katalog: lane.UserLanesDir() = ~/.jaira/lanes (core/lane/lane.go:466). Der '+'-Spalte im Lanes-Screen liest daraus. Bei berkc existiert der Ordner noch nicht.
  Teilen gibt es auch schon: publish kopiert nach .jaira/shared/<du>/ (ueber git), adopt kopiert eine geteilte Lane in ~/.jaira/lanes (lanes.go:484, :640).
  Fehlt also nur: beim Anlegen (oder per Taste) zusaetzlich in den Katalog schreiben. Offen: automatisch oder per Taste? Default-Vorschlag: automatisch, weil berkc es so erwartet hat.
  Test mit Alex (zweiter Rechner) erst nach dem Fix: er holt den Branch, oeffnet S -> Lanes, uebernimmt die geteilte Lane mit a und sieht sie unter '+'.
definition-of-done: "Eine im Lanes-Screen neu angelegte Lane liegt danach auch in ~/.jaira/lanes und erscheint in einem zweiten Repo unter '+'; Test mit JAIRA_HOME"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-18T21:53:12Z
updated-at: 2026-10-18T21:53:12Z
---

# Eine neu angelegte Lane landet nur in einem Board, nicht im persoenlichen Katalog ~/.jaira/lanes

## Definition of Done

- [ ] Eine im Lanes-Screen neu angelegte Lane liegt danach auch in ~/.jaira/lanes und erscheint in einem zweiten Repo unter '+'; Test mit JAIRA_HOME
- [ ] Eine Lane aus .jaira/shared/<kollege>/ landet mit a im Katalog und ist danach in jedem Board unter '+' waehlbar; Test
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

