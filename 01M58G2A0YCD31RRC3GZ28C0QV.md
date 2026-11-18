---
id: 01M58G2A0YCD31RRC3GZ28C0QV
title: "Eine neu angelegte Lane landet nur in einem Board, nicht im persoenlichen Katalog ~/.jaira/lanes"
status: testing
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
updated-at: 2026-11-18T17:22:35Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-74923
claimed-at: 2026-11-18T17:08:27Z
outcome-what: Eine mit n angelegte Lane wird nach dem Schliessen des Editors zusaetzlich nach ~/.jaira/lanes kopiert (catalogueNewLane); Tests fuer neue und adoptierte Lanes unter JAIRA_HOME.
outcome-why: "Lanes aus einem Repo fehlten in den Settings aller anderen Boards, weil n nur ins Board-Verzeichnis schrieb."
outcome-resolves: "Neue und adoptierte Lanes stehen in jedem Board unter '+'."
review-summary: "none"
review-gaps: "none: kein bestehender Built-in-id-Check gefunden (internal/cli/lanes.go:418 sammelt nur ids); isBuiltinID hat einen Aufrufer, bleibt wegen Lesbarkeit; catalogueNewLane parst die Datei zweimal (Parse + Adopt) — einmal pro n, gelassen; die alte, verwaiste Kommentarzeile ueber addAvailable ('newLane writes a fresh lane skeleton into the catalogue') ist vorbestehend und blieb"
test-verdict: "pass: go build + go test -race ./... RC=0; DoD 1-3 im Tree geprueft; TUI im tmux auf zwei Scratch-Boards: n in a -> ~/.jaira/lanes/e2e-lane.md + Meldung, Board b zeigt 'E2E lane' als verfuegbare Spalte; zweites n mit gleicher id -> 'on this board only', Katalog-md5 unveraendert; parallele Schreiber: O_EXCL im Skelett + Adopt-Refuse, nicht separat gestartet (TUI-Pfad, ein Tastendruck)"
hunt-verdict: "fail: bughunt done at 6fcd1288, 2 change — E auf einer n-Lane geht nicht in den Katalog und R wirft die Bearbeitung weg; eine n-Lane mit Built-in-id (review) ersetzt den Built-in im Katalog"
round: 3
---

# Eine neu angelegte Lane landet nur in einem Board, nicht im persoenlichen Katalog ~/.jaira/lanes

## Definition of Done

- [x] Eine im Lanes-Screen neu angelegte Lane liegt danach auch in ~/.jaira/lanes und erscheint in einem zweiten Repo unter '+'; Test mit JAIRA_HOME
  proof: internal/tui/lanes.go catalogueNewLane, internal/tui/model.go case newLaneDoneMsg; TestNewLaneLandsInCatalogueAndOnAnotherBoard (JAIRA_HOME)
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
- **2026-11-18 17:13 · Alexander Sacharov** — critique: nichts — Kopie geht ueber lane.Adopt (dasselbe wie adopt() in internal/tui/lanes.go), Aufruf im bestehenden newLaneDoneMsg-Zweig; Template-Check verhindert einen Platzhalter my-lane im Katalog. Drift-Risiko steht als eigene Notiz fuer die Person.
- **2026-11-18 17:15 · Alexander Sacharov** — bughunt: change: internal/tui/lanes.go:605 (catalogueNewLane) with lanes.go:537 (editLane) and lanes.go:648 (refreshDrift) - a lane made with n now exists twice, board file .jaira/lanes/new-lane.md and catalogue copy; a later E on the board opens the board file (l.Source), so the edit never reaches the catalogue (other boards keep the first version), the lane screen labels it 'drifted', and R then overwrites the edit with the old catalogue copy ('pulled the catalogue copy of deploy-check'). refreshDrift/editLane are older (0a15748b) but only this change puts n-made lanes on that path; ran: scratch test n -> catalogueNewLane -> edit l.Source -> newLaneScreen.drift has the id -> refreshDrift -> edited text gone
- **2026-11-18 17:15 · Alexander Sacharov** — bughunt: change: internal/tui/lanes.go:605 - catalogueNewLane copies whatever id the person typed, including a built-in's id (review, testing, ...): Adopt only refuses an existing catalogue FILE, and a built-in has none, so a board-local variant made with n lands as ~/.jaira/lanes/review.md and replaces the built-in on every board whose lane directory is not active; ran: scratch test n with id: review -> catalogueNewLane ok -> lane.Load(fresh root) returns review builtin=false name='Board-local review' from the catalogue file
- **2026-11-18 17:15 · Alexander Sacharov** — bughunt: done at 6fcd1288 - 2 findings (2 change, 0 pre-existing) in 3 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/release/notes/28C0QV.md internal/tui/lanecatalogue_test.go internal/tui/lanes.go internal/tui/model.go
- **2026-11-18 17:16 · Alexander Sacharov** — testing: zurueck wegen der 2 bughunt-change-Funde. Plan: (1) catalogueNewLane katalogisiert keine id eines Built-ins (lane.Builtins) — Meldung 'on this board only'. (2) editLane auf einer Board-Lane, deren Katalogkopie vor dem Edit byte-gleich war, schreibt nach dem Editor auch in den Katalog (write-through, wie der Drift-Kommentar in core/lane/share.go:130 es annimmt); war sie schon verschieden, bleibt der Katalog unberuehrt.
- **2026-11-18 17:19 · Alexander Sacharov** — decision: Bughunt-Fund 1 per write-through geloest, nicht per 'E oeffnet die Katalogkopie': sonst wirkt die Bearbeitung auf diesem Board erst nach R. Write-through nur, wenn Katalog == Board vor dem Edit (catalogueTwin) und der Katalog sich waehrend des Editors nicht geaendert hat; sonst bleibt er unberuehrt. Schreiben passiert vor m.reload(), damit die Drift-Berechnung des neuen Screens den neuen Katalog sieht.
- **2026-11-18 17:19 · Alexander Sacharov** — critique: writeThrough muss die geparste id mit der des Katalogfiles vergleichen; bei geaenderter id Katalog unberuehrt lassen und das melden.
- **2026-11-18 17:21 · Alexander Sacharov** — critique runde 2: id-Wechsel jetzt abgefangen (writeThrough); sonst nichts — newLaneDoneMsg traegt beide Faelle (n/E), weil beide ueber denselben ExecProcess-Rueckweg kommen; ein zweiter Msg-Typ waere ein zweiter case mit identischem Reload.
- **2026-11-18 17:22 · Alexander Sacharov** — bughunt: done at 55ad762d - 0 findings in 3 modules, follow-up hunt (diff since 6fcd1288); concurrency, error paths, input, permissions checked; files: core/release/notes/28C0QV.md internal/tui/lanecatalogue_test.go internal/tui/lanes.go internal/tui/model.go
