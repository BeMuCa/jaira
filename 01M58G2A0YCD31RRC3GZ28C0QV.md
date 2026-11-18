---
id: 01M58G2A0YCD31RRC3GZ28C0QV
title: "Eine neu angelegte Lane landet nur in einem Board, nicht im persoenlichen Katalog ~/.jaira/lanes"
status: critique
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
updated-at: 2026-11-18T19:52:27Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-62765
claimed-at: 2026-11-18T19:47:57Z
outcome-what: "Eine mit n angelegte Lane wird nach dem Editor auch nach ~/.jaira/lanes kopiert (catalogueNewLane), ausser sie traegt noch die Platzhalter-id des Templates oder eine Built-in-id; E schreibt in den Katalog durch, solange die Kopien byte-gleich waren; bricht der Editor mit Fehlercode ab, sagt der Lanes-Screen, was das fuer den Katalog heisst."
outcome-why: "Lanes aus einem Repo fehlten in den Settings aller anderen Boards; Signoff wollte zusaetzlich keine my-lane-Platzhalter im Katalog und keinen stillen Abbruch."
outcome-resolves: "Neue und adoptierte Lanes stehen in jedem Board unter '+'; Built-in-Overrides und unbenannte Skelette bleiben auf ihrem Board."
review-summary: "none"
review-gaps: |-
  note: catalogueNewLane erkennt ein unveraendertes Skelett nur bytegenau - speichert der Editor es mit kleiner Aenderung (z.B. Whitespace) und id my-lane, landet my-lane.md im Katalog und wird auf jedem Board unter '+' angeboten
  note: beendet sich der Editor mit Fehlercode (z.B. vim :cq), steigt model.go vor catalogueNewLane/writeThrough aus - die Lane bleibt dann still nur auf dem Board bzw. der Katalog behaelt die alte Kopie; nur die Fehlermeldung des Editors erscheint
  note: ohne Home-Verzeichnis ist lane.UserLanesDir() leer und catalogueTwin/writeThrough rechnen mit dem relativen Pfad <id>.md im aktuellen Verzeichnis (gleiches Altverhalten wie Export fuer Built-ins); praktisch irrelevant
  note: writeThrough prueft den Katalog und schreibt danach ohne Lock - zwei TUI-Sitzungen koennen sich in diesem kurzen Fenster ueberschreiben; fuer eine Ein-Personen-Datei vertretbar
test-verdict: "pass: go build ./... and go test -race -count=1 ./... green (RC=0); DoD 1-5 verified in the tree (all lane-catalogue tests pass incl. TestEditOfBuiltinIDDoesNotWriteThrough; share.go Drift comment read); E/n paths exercised through catalogueTwin/writeThrough/model.Update in tests, not in a live TUI (needs a tty and $EDITOR); failure paths: changed-meanwhile and id-change leave the catalogue alone (tests); this round adds no new write"
hunt-verdict: "pass: bughunt done at 16a8effb, 0 change, 0 pre-existing"
round: 6
review-verdict: "Diff erfuellt alle 5 DoD-Punkte; go vet und go test ./internal/tui ./core/lane gruen. Die beiden Bughunt-Funde (Drift nach E, Built-in-id im Katalog) sind sauber geschlossen und getestet. Kein Defekt gefunden, nur vier kleine Notes zu Randfaellen."
review-round: 4
review-check: |-
  1. Im Worktree: go build -o /tmp/jaira ./cmd/jaira
  2. export JAIRA_HOME=$(mktemp -d); mkdir -p /tmp/r1 /tmp/r2; cd /tmp/r1 && git init -q && /tmp/jaira init; cd /tmp/r2 && git init -q && /tmp/jaira init
  3. cd /tmp/r1 && /tmp/jaira - Lane-Settings oeffnen, n druecken, im Editor id: deploy-check und name: Deploy check setzen, speichern, schliessen
  4. Erwartet: Meldung 'deploy-check is on this board and in your catalogue (.../lanes/deploy-check.md)'; ls $JAIRA_HOME/lanes zeigt deploy-check.md
  5. cd /tmp/r2 && /tmp/jaira - Lane-Settings: deploy-check steht unter '+' (nicht installiert) zur Auswahl
  6. Zurueck in /tmp/r1: deploy-check waehlen, E, Prompt-Text aendern, speichern - Meldung 'saved on this board and in your catalogue'; cat $JAIRA_HOME/lanes/deploy-check.md zeigt den neuen Text, die Lane ist nicht als drifted markiert
  7. n noch einmal, im Editor id: review setzen - Meldung 'on this board only ... review is a built-in lane\'s id'; $JAIRA_HOME/lanes/review.md existiert nicht
  8. Automatisch: go test ./internal/tui -run 'NewLane|Catalogue|WriteThrough|EditOf|Adopted' -v
---

# Eine neu angelegte Lane landet nur in einem Board, nicht im persoenlichen Katalog ~/.jaira/lanes

## Definition of Done

- [x] Eine im Lanes-Screen neu angelegte Lane liegt danach auch in ~/.jaira/lanes und erscheint in einem zweiten Repo unter '+'; Test mit JAIRA_HOME
  proof: internal/tui/lanes.go catalogueNewLane, internal/tui/model.go case newLaneDoneMsg; TestNewLaneLandsInCatalogueAndOnAnotherBoard (JAIRA_HOME)
- [x] Eine Lane aus .jaira/shared/<kollege>/ landet mit a im Katalog und ist danach in jedem Board unter '+' waehlbar; Test
  proof: TestAdoptedSharedLaneIsOfferedOnEveryBoard (internal/tui/lanecatalogue_test.go) — adopt() schrieb schon in den Katalog
- [x] Release-Note-Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/28C0QV.md
- [x] E auf einer Board-Lane mit der id eines Built-ins schreibt nicht in den Katalog durch; Test
  proof: internal/tui/lanes.go catalogueTwin (isBuiltinID guard); TestEditOfBuiltinIDDoesNotWriteThrough
- [x] Drift-Kommentar in core/lane/share.go sagt, dass nur ein byte-gleiches Paar durchgeschrieben wird
  proof: core/lane/share.go Drift doc comment
- [x] Ein mit n angelegtes File, das noch die Platzhalter-id des Templates (my-lane) traegt, kommt nicht in den Katalog, auch wenn es sonst veraendert wurde; Meldung nennt den Grund; Test
  proof: internal/tui/lanes.go catalogueNewLane + placeholderID; TestNewLaneEditedSkeletonWithPlaceholderIDStaysOutOfCatalogue
- [x] Beendet sich der Editor mit Fehlercode, sagt die Meldung im Lanes-Screen zusaetzlich, dass die Lane nur auf dem Board blieb bzw. der Katalog die alte Kopie behaelt; Test
  proof: internal/tui/lanes.go editorFailedMsg, internal/tui/model.go case newLaneDoneMsg; TestEditorErrorSaysCatalogueWasSkipped

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
- **2026-11-18 18:43 · Alexander Sacharov** — explain: Der Review fand keinen Blocker, nur drei Hinweise. Eine mit n angelegte Lane wird jetzt auch in ~/.jaira/lanes kopiert, und E schreibt Änderungen in den Katalog durch. Dabei ist der Kommentar zu Drift in core/lane/share.go:132 veraltet, und E auf einer Lane mit eingebauter id ändert den Katalog-Override für alle Boards.
Warum: Der Kommentar behauptet etwas, das nur noch im byte-gleichen Fall stimmt, und kann spätere Leser in die Irre führen. Die Wirkung auf andere Boards ist konsistent mit der Regel, aber breiter, als die Person auf ihrem Board erwartet.
1) Zurückgeben: Drift-Kommentar in core/lane/share.go:132 korrigieren und E auf eingebaute ids nicht durchschreiben lassen [give back]
2) Zurückgeben: nur den Drift-Kommentar korrigieren, das Durchschreiben auf eingebaute ids lassen wie es ist [give back]
3) Ticket so annehmen, die drei Hinweise bleiben als Notes stehen [accept]
Empfehlung: 3 - beide DoD-Punkte sind erfüllt und getestet, die Hinweise ändern kein Verhalten für den Alltag, der Kommentar lässt sich in einem späteren Commit mitnehmen.
- **2026-11-18 18:44 · Alexander Sacharov** — decision: 1) Zurückgeben: Drift-Kommentar in core/lane/share.go:132 korrigieren und E auf eingebaute ids nicht durchschreiben lassen
- **2026-11-18 18:44 · Alexander Sacharov** — Given back from signoff: 1) Zurückgeben: Drift-Kommentar in core/lane/share.go:132 korrigieren und E auf eingebaute ids nicht durchschreiben lassen
- **2026-11-18 18:45 · Alexander Sacharov** — cockpit: given back to a new worker - nothing recorded about the old worker session (cache and context unknown)
- **2026-11-18 18:46 · Alexander Sacharov** — decision: Built-in-Schutz sitzt in catalogueTwin, nicht in editLane — dort wird syncTo bestimmt, und catalogueNewLane nutzt dasselbe isBuiltinID. Release-Note Zeile 2 nennt jetzt auch E.
- **2026-11-18 18:46 · Alexander Sacharov** — critique runde 3 (nach signoff-Rueckgabe): nichts — beide Rueckgabepunkte beantwortet; der Guard sitzt in catalogueTwin, wo syncTo entsteht, und nutzt dasselbe isBuiltinID wie catalogueNewLane. Dass eine Board-Lane mit Built-in-id nach E als drifted erscheint, ist die gewaehlte Folge (Option 1), kein Fund.
- **2026-11-18 18:48 · Alexander Sacharov** — bughunt: done at 16a8effb - 0 findings in 4 modules, follow-up hunt (diff since 55ad762d; core/lane/share.go new to the list, read whole); concurrency, error paths, input, permissions checked; ran: go test ./internal/tui/ ./core/lane/ (ok); files: core/lane/share.go core/release/notes/28C0QV.md internal/tui/lanecatalogue_test.go internal/tui/lanes.go internal/tui/model.go
- **2026-11-18 18:52 · Alexander Sacharov** — review (fresh context): kein Blocker; DoD 1-5 erfuellt, Tests gruen. Vier Notes in review-gaps (Skelett-Erkennung nur bytegenau, Editor-Fehlercode ueberspringt Katalog, leeres UserLanesDir, Check-then-write ohne Lock).
- **2026-11-18 19:39 · Alexander Sacharov** — explain: Der Review fand keinen Defekt, nur vier Randfaelle: ein leicht veraendertes Skelett (z.B. Whitespace) mit id my-lane landet im Katalog, ein Editor-Abbruch mit Fehlercode ueberspringt den Katalog-Schritt still, fehlendes Home-Verzeichnis nutzt einen relativen Pfad, und zwei TUI-Sitzungen koennen sich beim Write-Through ueberschreiben.
Warum: Alle fuenf DoD-Punkte sind erfuellt und die Tests gruen, die Randfaelle treffen nur seltene Situationen in einer Ein-Personen-Datei und verlieren keine Daten.
1) Ticket so annehmen, die vier Randfaelle bleiben als Notes stehen [accept]
2) Zurueckgeben: Skelett-Erkennung tolerant machen und Editor-Fehlercode melden [give back]
Empfehlung: 1 - Kein Blocker, die Faelle sind selten und harmlos, ein weiterer Durchlauf kostet mehr als er bringt.
- **2026-11-18 19:47 · Alexander Sacharov** — decision: 2) Zurueckgeben: Skelett-Erkennung tolerant machen und Editor-Fehlercode melden
- **2026-11-18 19:47 · Alexander Sacharov** — Given back from signoff: 2) Zurueckgeben: Skelett-Erkennung tolerant machen und Editor-Fehlercode melden
- **2026-11-18 19:47 · Alexander Sacharov** — cockpit: given back to a new worker - nothing recorded about the old worker session (cache and context unknown)
- **2026-11-18 19:48 · Alexander Sacharov** — decision: Skelett-Erkennung tolerant = geparste id gleich der Platzhalter-id des Templates (aus lane.Parse(lane.Template) gelesen, nicht hart my-lane). Unveraendertes Skelett bleibt still; ein veraendertes mit my-lane bekommt eine Meldung. Editor-Fehlercode: Katalog-Schritt bleibt uebersprungen (vim :cq heisst Abbruch), die Meldung haengt nur an, was das fuer den Katalog bedeutet — nicht trotzdem kopieren.
- **2026-11-18 19:51 · Alexander Sacharov** — critique runde 4: Platzhalter-Meldung und Release-Note-Zeile 3 empfehlen E, aber E katalogisiert eine Lane ohne Katalogkopie nie (catalogueTwin). Text korrigieren, kein neues Verhalten fuer E (sonst wuerde jede bearbeitete Board-Lane katalogisiert). Sonst nichts: editorFailedMsg sitzt neben den anderen Meldungs-Helfern in lanes.go, placeholderID liest die id aus dem Template statt sie hart zu kodieren.
- **2026-11-18 19:52 · Alexander Sacharov** — critique runde 5: der eine Fund aus Runde 4 ist beantwortet - Meldung in catalogueNewLane und Release-Note Zeile 3 verweisen jetzt auf n mit eigener id statt auf E. Nichts weiter.
