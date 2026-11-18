---
id: 01M58G2A0YCD31RRC3GZ28C0QV
title: "Eine neu angelegte Lane landet nur in einem Board, nicht im persoenlichen Katalog ~/.jaira/lanes"
status: in-progress
ready: true
creator: berkc
assignee: Alexander Sacharov
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
updated-at: 2026-11-18T17:12:38Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-74923
claimed-at: 2026-11-18T17:08:27Z
---

# Eine neu angelegte Lane landet nur in einem Board, nicht im persoenlichen Katalog ~/.jaira/lanes

## Definition of Done

- [x] Eine im Lanes-Screen neu angelegte Lane liegt danach auch in ~/.jaira/lanes und erscheint in einem zweiten Repo unter '+'; Test mit JAIRA_HOME
  proof: internal/tui/lanes.go catalogueNewLane, internal/tui/model.go newLaneDoneMsg; TestNewLaneLandsInCatalogueAndOnAnotherBoard (JAIRA_HOME)
- [x] Eine Lane aus .jaira/shared/<kollege>/ landet mit a im Katalog und ist danach in jedem Board unter '+' waehlbar; Test
  proof: TestAdoptedSharedLaneIsOfferedOnEveryBoard (internal/tui/lanecatalogue_test.go) — adopt() schrieb schon in den Katalog
- [x] Release-Note-Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/28C0QV.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-18 17:08 · Alexander Sacharov** — Alexander Sacharov took this ticket over from berkc
- **2026-11-18 17:08 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/6601.sock
- **2026-11-18 17:08 · Alexander Sacharov** — Alex übernimmt den Ticket von berkc am 2026-11-18 (Teamlead-Anweisung).
- **2026-11-18 17:08 · Alexander Sacharov** — decision: neue Lane wird automatisch auch nach ~/.jaira/lanes geschrieben (kein Extra-Key) — vom Teamlead/Alex entschieden. Test auf zweitem Rechner entfällt; Test mit JAIRA_HOME reicht. core/lane/refresh.go nicht anfassen (NH788Q läuft parallel).
- **2026-11-18 17:09 · Alexander Sacharov** — decision: Kopie in den Katalog passiert nach dem Schliessen des Editors (newLaneDoneMsg), nicht beim Schreiben des Skeletts — erst dann steht die echte id im Frontmatter. Weg: lane.Adopt(path, UserLanesDir(), overwrite=false). Gibt es die id im Katalog schon, wird nichts ueberschrieben (Meldung im Lanes-Screen). Ein unveraendertes Skelett (Bytes == lane.Template) wird nicht katalogisiert. DoD 2 (adopt -> Katalog) macht adopt() schon heute; es fehlt nur der Test, dass die Lane danach in einem anderen Board unter '+' (lane.Installable) steht.
- **2026-11-18 17:12 · Alexander Sacharov** — found: E auf einer Projekt-Lane bearbeitet die Projektkopie (l.Source), nicht die Katalogkopie. Ist die Lane jetzt auch im Katalog, meldet Drift sie danach als 'drifted', und R zieht die alte Katalogkopie zurueck ins Board — die Bearbeitung waere weg. Galt schon vorher fuer adoptierte+hinzugefuegte Lanes, wird mit diesem Ticket haeufiger. Nicht gefixt (ausserhalb des Scopes, Drift-Comment in core/lane/share.go:130 nimmt die Gegenrichtung an); fuer die Person zur Entscheidung.
