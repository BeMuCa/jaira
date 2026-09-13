---
id: 01M28KACAF28348E0SJK51RNC1
title: "Ein Link-Fenster zeigt alles, was mit einem Ticket verbunden ist"
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Von jedem Ticket aus in einem Fenster jede Verbindung sehen — Blocker, Eltern, Kinder, Verwandte, Folgetickets — auch wenn das andere Ticket längst im Logbuch, im Archiv oder nur auf einem git-ref liegt. Dazu zwei neue Beziehungen: parent (woraus Kinder abgeleitet werden) und related."
context: |-
  Heute reisst jede Verbindung ab, sobald das andere Ticket fertig ist.

  - Store.Paths() (core/ticket/store.go:521) liest nur .jaira/tickets/. Logbuch und Archiv sind nicht im Index, also antwortet 'jaira show <id>' auf ein abgelegtes Ticket mit 'not found', obwohl die Datei unter .jaira/logbook/ liegt.
  - core/gate/gate.go:572: ein blocked-by, das nicht im Store steht, ergibt CodeBlocked 'does not exist in this store'. Der Blocker wurde also fertig, wanderte ins Logbuch — und genau dadurch ist das abhaengige Ticket fuer immer blockiert und faellt aus 'jaira list --actionable'.
  - core/validate/validate.go:174 meldet denselben Fall als SeverityError 'the dependency can never clear'. Falsch: die Abhaengigkeit ist erfuellt, nicht kaputt.
  - Es gibt heute nur blocked-by (harte Abhaengigkeit) und follows (Folgeticket). Eine weiche 'haengt damit zusammen'-Beziehung fehlt, und eine Eltern/Kind-Beziehung fehlt ganz.
  - Bereits vorhanden und wiederzuverwenden: Store.FiledAwayIDs() (core/ticket/store.go:629) liefert id -> Pfad fuer Logbuch und Archiv; core/gitref kennt Tickets, die nur auf einem ref liegen (Ticket.ReadOnly); internal/tui/refs.go ist das Vorbild fuer ein eigenes Fenster (mode-Konstanten in internal/tui/model.go:36).
  - Entschieden und nicht neu zu diskutieren: KEIN Feld 'children' auf dem Elternteil. Ein Kind anzulegen wuerde sonst die fremde Datei aendern — die es hier gar nicht immer gibt (Logbuch, nur-ref) — und zwei parallele Sessions kollidieren auf derselben Zeile. Stattdessen 'parent: <id>' auf dem Kind, ein Schreiber pro Datei; Kinder und die ganze Tiefe werden beim Lesen aus dem Index abgeleitet. 'children' bleibt ein Wort in der Oberflaeche, kein Feld im Format.
  - Ebenfalls entschieden: parent ist KEIN Gate. Ein offenes Kind haelt den Elternteil nicht auf. Wer Epics als Zwang will, nimmt blocked-by.
definition-of-done: "Ein Link-Fenster im TUI zeigt zu jedem Ticket alle Verbindungen in beide Richtungen, quer ueber Board, Logbuch, Archiv und git-refs; parent und related existieren als Felder; und ein Blocker im Logbuch blockiert nicht mehr."
tags:
  - tui
  - cli
  - gates
blocked-by: []
commits:
  - ebe4305637eb7a38d99443d3985e6fa95caaf498
  - 24fbb376a0b179f6fe6e6f00bf83972920653b77
  - 28ff4ae23cdbb47d026b14fb328ca59f89876482
  - 4a31d7b85482da16bae02eee09fd3f0dd30d7af8
  - ba77d3781c03c9bfe88a00df8ffec2f611e3bb69
  - 8292a2934e908695c11df5d3df3954e7101832c5
  - 64e96a3a296630c8334ad31e1e270326c6af27bd
  - cefa764e1fa45f8a70f50b8cde61bbd1d9cabc17
  - 70c2e6240f6c45d6c0d702e53ffa8c785f8ff54e
created-at: 2026-09-11T16:02:57Z
updated-at: 2026-09-11T20:23:31Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-585885
claimed-at: 2026-09-11T16:09:03Z
outcome-what: "Alle neun Review-Befunde behoben, mit Tests fuer jeden, der eine Regression sein kann."
outcome-why: "Das Review hat den Kernpunkt des DoD in der TUI als nicht erfuellt nachgewiesen."
outcome-resolves: "gateEnv setzt Satisfied; gate fragt Satisfied zuerst; Ueberschriften kosten zwei Zeilen; Enter meldet den Filter; create nutzt resolveRef; LoadHeader statt voller Dateien; Optimize-Arbeit committed; CLI- und core/link-Tests ergaenzt; NOTES-Zeile fuer show."
review-summary: "Neues Paket core/link haelt einen Index ueber alle vier Orte, an denen ein Ticket liegen kann - Board, eigener ref, Logbuch, Archiv - und beantwortet jede Beziehung in beide Richtungen: blocked-by/blockiert, parent/enthaelt (als Baum, zyklusfest), related, follows/gefolgt-von. Zwei neue Felder: parent als Skalar auf dem Kind, related als Liste, im Merge-Driver als Union. Kinder werden abgeleitet, nie gespeichert. Aus dem Index werden zwei Praedikate in die reinen Pakete injiziert - gate.Env.Satisfied und validate.Tickets(known) - womit ein fertiger, abgelegter Blocker nicht mehr blockiert und nicht mehr als tot gemeldet wird. CLI: jaira links <id> mit --json, jaira show findet abgelegte Tickets und zeigt part-of/contains/related, set und create loesen Handles zu vollen ids auf und lehnen tote Referenzen ab. TUI: L oeffnet ueber Board und ueber offenem Ticket ein scrollendes Link-Fenster, das als Kasten ueber der Doska liegt statt sie zu ersetzen; derselbe Modal-Helfer traegt jetzt auch Legende und Meldungen."
review-gaps: "none - die drei Review-Runden haben 22 Befunde ergeben, alle behoben und die regressionsfaehigen davon einzeln mit Tests festgenagelt. Offen bleibt nur etwas, das nicht in diesem Diff steht: P1AE82 liegt auf diesem Board gleichzeitig unter tickets/ und im Logbuch, aelter als dieser Zweig, und nur ein Mensch kann sagen, welche Kopie geht."
test-verdict: "go build, go vet und go test ./... gruen, 25 Pakete mit Tests, keine Fehlschlaege. Die Befunde aus Runde 1 sind einzeln nachgestellt: TestTheHintLineSurvivesSeveralGroups faellt mit der alten Zeilenrechnung bei 20 und 30 Zeilen und besteht mit der neuen; TestStaleRefCopyDoesNotOutvoteTheIndex haelt den Gate-Vorrang fest; TestJumpingToAFilteredOutCardSaysSo den stummen Sprung. Die TUI-Luecke (Satisfied fehlte in gateEnv) ist behoben und laeuft jetzt durch dieselbe gate.Env-Konstruktion wie die CLI."
review-verdict: "Die Sache existiert und tut, was das Ticket verlangt, in beiden Oberflaechen; vom Nutzer auf der echten Doska bestaetigt. go build, go vet und go test ./... gruen."
review-check: |-
  1. cd ~/projects/jaira
  2. B=.claude/worktrees/feat-51RNC1-link-window/jaira-links
  3. $B create 'Probe: Epic' --goal g --context c --dod d --mine, dann $B pull auf die gemeldete Handle
  4. $B create 'Probe: Kind' --parent <epic-handle> --goal g --context c --dod d --mine, ebenfalls pullen
  5. $B create 'Probe: Blocker' --lane done --goal g --context c --dod d --mine, pullen, dann $B set <epic> blocked-by=<blocker> und $B logbook <blocker>
  6. $B links <epic> - erwartet: 'waiting on ... logbook · done' und 'contains ... on the board'
  7. $B list --actionable | grep <epic> - der Epic steht drin, obwohl er einen Blocker hat: der Blocker ist fertig
  8. $B starten, Karte 'Probe: Epic' waehlen, L druecken - Kasten ueber der sichtbaren Doska, j/k bewegen, esc zurueck
  9. Terminal auf 12 Zeilen verkleinern, L erneut - die Zeile 'j/k move enter jump esc back' muss sichtbar bleiben und die Auswahl auf dem Schirm
  10. Aufraeumen: Links loeschen (set ... blocked-by= parent=), dann delete --force, den Blocker vorher mit restore zurueckholen
---

## Definition of Done

- [x] Feld `related` (Liste von ids) und `parent` (eine id) im Schema, in der kanonischen Reihenfolge, im Merge-Driver als Union bzw. Skalar (core/ticket/schema.go, core/merge/merge.go:34)
- [x] Ein Beziehungs-Index, der Board, Logbuch, Archiv (Store.FiledAwayIDs) und git-refs zusammenfasst und jede Beziehung in beide Richtungen aufloest: blocked-by/blockiert, parent/Kinder (rekursiv), related, follows/gefolgt-von
- [x] Ein Blocker, der im Logbuch oder Archiv in einer terminalen Lane liegt, gilt als erfuellt: core/gate/gate.go:572 blockiert nicht mehr, core/validate/validate.go:174 meldet keinen Fehler mehr
- [x] `jaira show <id>` findet auch ein abgelegtes Ticket und sagt dazu, wo es liegt
- [x] CLI: Beziehungen setzbar (`jaira set <id> parent=<id>`, `related=<id>,<id>`) und lesbar (`jaira links <id>` mit --json)
- [x] TUI: eine Taste auf einer Karte oeffnet ein Link-Fenster, gruppiert nach Beziehungsart, mit Lane und Ablageort je Eintrag; Enter springt auf ein Ticket, das hier ein File hat; Esc schliesst
- [x] Ein Ticket kann nicht sein eigener Elternteil sein und keinen Zyklus bilden; validate meldet beides
- [x] Eine Zeile in core/release/NOTES.md unter ## Unreleased je sichtbarer Aenderung
- [x] Tests: Beziehungsindex ueber alle vier Quellen, Blocker-im-Logbuch loest den Block, Zyklus wird gemeldet, Golden-Test fuer das Link-Fenster

## Options

- [x] planning

## Plan

- [x] core/ticket: Felder parent (Skalar) und related (Liste) anlegen, in canonicalOrder einsortieren, in Ticket/Decode lesen
- [x] core/merge: related in listFields aufnehmen; parent als gewoehnlichen Skalar bestaetigen (Test)
- [x] internal/cli set: related in die listFields-Map; parent/related in create-Flags und in ticketJSON
- [x] core/ticket/store.go: LoadAnywhere — resolve auch ueber FiledAwayIDs (Logbuch+Archiv), Ergebnis traegt seinen Ablageort
- [x] core/link: neues Paket — Index aus Board (List, enthaelt schon die refs via Source.Extra) plus FiledAwayIDs; Relations(id) liefert beide Richtungen je Art
- [x] core/link: Kinder rekursiv als Baum, mit Zyklusschutz (besuchte ids)
- [x] core/gate: blockedBy fragt den Index — Blocker abgelegt in terminaler Lane = erfuellt, kein CodeBlocked mehr
- [x] core/validate: DanglingDep nur noch wenn die id nirgends existiert; neue Pruefung ParentCycle und Selbst-Elternschaft
- [x] internal/cli: jaira links <id> (Text + --json), jaira show zeigt parent/related/Kinder und sagt wo ein abgelegtes Ticket liegt
- [x] internal/tui: modeLinks neben modeLegend, Taste auf der Karte, Fenster nach Beziehungsart gruppiert, Enter springt, Esc schliesst
- [x] Tests: Index ueber alle vier Quellen, Blocker-im-Logbuch, parent-Zyklus, Golden-Test fuer das Link-Fenster
- [x] core/release/NOTES.md: je sichtbarer Aenderung eine Zeile unter ## Unreleased

## Progress

- **2026-09-11 16:10 · Alexander Sacharov** — Warum der Plan so aussieht:

- Die git-refs sind schon drin. Store.List() haengt ueber Source.Extra (core/refsync/refsync.go:793) die Tickets an, die nur auf einem ref liegen, und markiert sie ReadOnly. Fuer den Index heisst das: nur Logbuch und Archiv fehlen, dafuer gibt es Store.FiledAwayIDs() (core/ticket/store.go:629). Keine eigene ref-Abfrage bauen.
- Deshalb ein eigenes Paket core/link statt Logik in core/gate: gate, validate, cli und tui brauchen dieselbe Aufloesung, und drei Kopien davon waeren drei Wahrheiten.
- parent statt children: ein Kind anzulegen darf nie die Datei des Elternteils anfassen, weil es die hier oft gar nicht gibt (Logbuch, nur-ref = ReadOnly, jede Mutation wird abgelehnt). Ein Schreiber pro Datei ist zugleich das, was den Merge-Driver konfliktfrei haelt.
- related landet in core/merge listFields (Zeile 34): zwei Sessions, die verschiedene Verwandte eintragen, haben beide recht; eine Seite zu waehlen verliert stillschweigend Daten. parent ist ein Skalar, da gibt es genau eine Antwort.
- Zyklusschutz ist Pflicht, nicht nett: parent wird beim Rendern rekursiv verfolgt, ein Ring haengt die TUI auf.
- Offen gelassen und bewusst nicht geplant: parent als Gate. Ein offenes Kind haelt den Elternteil nicht auf — wer das will, nimmt blocked-by.
- **2026-09-11 16:25 · Alexander Sacharov** — Umgesetzt. Was beim Bauen herauskam und nirgends sonst steht:

- Die Taste ist L, nicht l: l ist im Board die vim-Bewegung nach rechts. L war frei, seit der Lane-Screen nach S umgezogen ist — der Test TestSettingsOpensFromBoardAndLNoLonger hat genau das festgehalten und musste umgeschrieben werden, weil er 'L tut nichts' behauptete. Er prueft jetzt, dass L nicht mehr den Lane-Screen oeffnet, was sein eigentlicher Zweck war.
- core/gate bekommt keinen Store, sondern eine Funktion Satisfied(id) bool — genau wie DeriveCommits. Das Paket verspricht in seinem Doc-Kommentar, dateisystemfrei zu sein, und seine Tests bauen Env ohne Repository.
- validate.Tickets hat einen dritten Parameter known(id) bool bekommen. 21 Testaufrufe mussten mit nil nachgezogen werden; nil heisst 'nur die uebergebenen Tickets existieren', also genau das alte Verhalten.
- Der Index liest das Logbuch faul: FiledAwayIDs ist nur ein Verzeichnis-Listing, und erst eine Frage, die nur die Datei beantworten kann, oeffnet sie. Rueckwaerts-Links (wer nennt mich als Elternteil) brauchen alle Dateien — das passiert nur beim Oeffnen des Fensters, nie beim Rendern des Boards.
- Logbuch heisst fertig, ohne die Datei zu lesen: dort landet nur Arbeit aus einer terminalen Lane. Das Archiv ist eine Schublade und kein Urteil, deshalb wird dort die Lane nachgelesen.
- Whereabouts und Label leben in core/link, nicht in CLI und TUI je einmal: 'logbook · done' ist der Unterschied zwischen erfuellter und offener Abhaengigkeit, und zwei Schreibweisen davon waeren zwei Antworten.
- **2026-09-11 18:50 · Alexander Sacharov** — critique: vier Befunde, alle mit klarer Korrektur, also zurueck in die Implementierung statt in die HITL-Lane.

1. Doppelte Reihenfolge: kindOrder (internal/cli/links.go) und linkKindOrder (internal/tui/links.go) sind dieselbe Liste zweimal. Gehoert nach core/link, wo schon Kind.Title() und Ref.Whereabouts() liegen - genau aus demselben Grund dorthin gezogen.
2. filedAway() in internal/cli/tickets.go listet das Logbuch erneut durch, obwohl LoadAnywhere drei Zeilen darueber den Pfad schon geliefert hat. Zwei Verzeichnis-Listings fuer eine Frage.
3. Der harte Schnitt in clipHeight() und der weiterlaufende Cursor in linkView widersprechen sich: ein Epic mit vielen Kindern kann eine Auswahl haben, die nicht sichtbar ist. Das ist der Fehler, den der Nutzer zuerst sehen wird.
4. childTree() ist quadratisch, weil es je Knoten alle Tickets scannt. Eine Map parent -> Kinder einmal bauen kostet nichts und macht die Tiefe umsonst.

Nicht beanstandet und bewusst so gelassen: parent auf dem Kind (die Begruendung steht in der Notiz davor), der injizierte Satisfied-Funktionszeiger in core/gate (folgt DeriveCommits) und das faule Lesen des Logbuchs.
- **2026-09-11 18:54 · Alexander Sacharov** — in-progress nach critique: alle vier Befunde erledigt.

1. link.Order liegt jetzt in core/link neben Kind.Title() und Ref.Whereabouts(); CLI und TUI lesen dieselbe Liste. Der TUI-eigene byID() ist ebenfalls weg - der Titel des Subjekts wird beim Oeffnen mitgenommen, was auch fuer ein Ticket funktioniert, fuer das das Board keine Karte hat.
2. filedAway() nimmt den Pfad, den LoadAnywhere schon zurueckgibt, statt das Logbuch ein zweites Mal zu listen.
3. Das Fenster scrollt jetzt mit dem Cursor und schreibt hin, wie viel oben und unten verborgen ist. Dabei kam heraus, dass die Zeilenzahl mit dem abstimmen muss, was modalOver fuer Inhalt uebrig laesst - sonst schneidet der Modal-Rahmen die Liste ab, und dieser Schnitt kann nicht scrollen. Die Konstante steht deshalb mit Begruendung an linkRows().
4. childrenOf() baut die Map parent -> Kinder einmal; childTree() liest sie nur noch. Vorher quadratisch in der Ticketzahl.

Der Test TestLongLinkListScrollsWithTheCursor haelt Befund 3 fest: er baut ein Epic mit 20 Kindern in einem 20 Zeilen hohen Fenster und prueft nach jedem j, dass die Auswahl sichtbar ist. Er ist beim ersten Lauf fehlgeschlagen und hat genau den Modal-Schnitt gefunden.
- **2026-09-11 19:11 · Alexander Sacharov** — review-Runde 1: neun Befunde, alle behoben.

Der schwerste war echt und haette die Abnahme nicht ueberlebt: internal/tui/model.go gateEnv() hat Satisfied nie gesetzt, also blockierte ein abgelegter Blocker auf dem Board weiter, waehrend die CLI ihn durchliess. Genau der DoD-Punkt, um den es hier geht, und er galt nur in einer der beiden Oberflaechen. Der Index haengt jetzt am Model, wird einmal pro reload gebaut und beim reload verworfen - nicht pro Render, sonst listet jede Karte einzeln das Logbuch.

core/gate fragt Satisfied jetzt zuerst und fuer jede Abhaengigkeit, nicht nur fuer eine, die in env.All fehlt. Ein hier abgelegtes Ticket faehrt weiter auf seinem ref, also steht in All eine Kopie mit dem alten Status - und die hat die fertige Arbeit weiter blockiert. Der alte Test TestOnBoardUnfinishedBlockerStillBlocks behauptete das Gegenteil (All schlaegt Index); er beschrieb eine Welt, die nicht vorkommen kann, weil beide aus denselben Tickets stammen, und ist umgeschrieben.

Das Fenster rechnete je Ueberschrift eine Zeile, druckte aber zwei - ab drei Gruppen schnitt der Modal-Rahmen die Hinweiszeile ab. Der neue Test prueft 20, 30 und 40 Zeilen und faellt mit der alten Rechnung bei 20 und 30 um, so wie es das Review beschrieben hat.

selectByID meldet jetzt, ob es etwas gefunden hat; unter einem Filter sagt Enter, dass der Filter die Karte verdeckt, statt nichts zu tun.

create --parent/--related gehen ueber resolveRef wie set, also auch an einen abgelegten Elternteil.

core/link liest abgelegte Tickets ueber das neue ticket.LoadHeader - nur Frontmatter, gedeckelt, so wie das Board sich selbst liest. Vorher las 'jaira show' die kompletten Bodies des ganzen Logbuchs.

Neu getestet, was vorher gar keinen Test hatte: die vierte Quelle (PlaceRef), der Vorrang Logbuch vor ref-Kopie, bySuffix, und auf CLI-Seite links, links --json, set mit Handle, die Ablehnung einer toten Referenz, show auf einem abgelegten Ticket und die Kinderliste in show.

Nicht im Diff und trotzdem wahr: P1AE82 liegt auf diesem Board gleichzeitig unter tickets/ und im Logbuch. Das ist aelter als dieser Zweig und nur der Mensch kann entscheiden, welche Kopie geht.
- **2026-09-11 19:29 · Alexander Sacharov** — review-Runde 2: acht Befunde, alle behoben.

Der Abschnitt der Hinweiszeile war in Runde 1 nur halb repariert. Zwei Ursachen: das abschliessende Newline jeder Ansicht zaehlt beim Splitten als Zeile und fiel dem Schnitt zum Opfer (modal.go trimmt es jetzt), und die beiden Zeilen 'N more above/below' kosteten Platz, den niemand eingeplant hatte. visible() teilt den Platz jetzt nach Rang: die ausgewaehlte Zeile muss sichtbar sein, dann die Liste, und die beiden Hinweise auf Verborgenes geben ihren Platz wieder her, sobald sie die Auswahl vom Schirm draengen wuerden. Passt gar nichts, beginnt der Ausschnitt bei der Auswahl selbst. Getestet ueber zehn Fensterhoehen, jeweils nach 0 bis 6 Bewegungen.

openLinks ruft jetzt scrollInto, sonst oeffnete das Fenster unter 19 Zeilen mit der Auswahl ausserhalb des Sichtbereichs.

Satisfied und Known oeffnen keine Datei mehr: place() beantwortet aus dem Verzeichnis-Listing, was das Board und das Logbuch schon sagen, und nur das Archiv - eine Schublade, kein Urteil - wird gelesen. Vorher kostete jeder abgelegte Blocker einen Dateizugriff pro Render.

bySuffix normalisiert wie Store.filedPath (NormalizeIDPrefix), sonst loeste ein klein geschriebener Handle in set auf, galt im Fenster aber als tot.

validate.parentCycle folgt der Kette auch ueber Handles, sonst blieb ein so geschriebener Ring ungemeldet, waehrend known() ihn durchliess.

openLinks nimmt den Index des Models statt einen neuen zu bauen.

Neu festgenagelt: der Kernpunkt des DoD auf dem Board selbst (TestAFiledBlockerDoesNotBlockOnTheBoard - faellt, wenn man Satisfied aus gateEnv entfernt, geprueft), related als Union und parent als Skalar im Merge-Driver.
- **2026-09-11 19:50 · Alexander Sacharov** — review-Runde 3: fuenf Befunde. Die Zeilenrechnung des Fensters war zum dritten Mal falsch, also ist sie nicht noch einmal geflickt, sondern ersetzt.

Vorher lag das Budget in linkLines(), der Ausschnitt in visible() und die Verschiebung in scrollInto() - drei Stellen, die sich gegenseitig nicht glaubten. Der Ausweichzweig in visible() zeichnete sogar ab der Auswahl, ohne v.top zu schreiben: danach stand top auf 0, waehrend Zeile 9 oben stand, und der naechste Tastendruck riss die Liste neun Zeilen zurueck.

Jetzt gibt es frame(): es verschiebt top selbst, schreibt es zurueck und liefert denselben Ausschnitt, ueber den es entschieden hat. Die Rangfolge steht im Kommentar und im Code an einer Stelle: Hinweiszeile immer, dann die Auswahl, dann die Ueberschrift darueber, dann der Rest der Liste, dann Titel und Trennlinie, zuletzt die beiden Zeilen ueber Verborgenes - und die werden nur aus Zeilen gezeichnet, die wirklich uebrig sind.

Die Hoehe des Kastens kommt nicht mehr aus zwei Rechnungen: modalHeight()/modalContent() in modal.go sind die einzige Quelle, das Fenster fragt dort. Genau diese Doppelrechnung war die Ursache in allen drei Runden.

validate: eine Handle-Erwaehnung des Elternteils oder eines related-Tickets im Kontext ist keine verschwiegene Abhaengigkeit mehr. Vorher bekam jedes ordentlich geschriebene Kind den Rat, sein Epic in blocked-by einzutragen - im Widerspruch zur Entscheidung dieses Tickets, dass parent kein Gate ist. Und parentCycle weitet nur noch Referenzen, die keine vollstaendige id sind: sonst haette eine Handle-Kollision einen Ring erfunden, den es nicht gibt.

Der Test laeuft jetzt ueber jede Hoehe von 10 bis 45, einmal die ganze Liste hinunter und wieder hinauf, und prueft bei jedem Schritt: Hinweiszeile da, Auswahl sichtbar, die Zaehler stimmen mit dem Ausschnitt ueberein, und der Ausschnitt folgt dem Cursor, ohne ihn zu ueberholen. Mit der alten Markerrechnung faellt er bei 10 Zeilen.
