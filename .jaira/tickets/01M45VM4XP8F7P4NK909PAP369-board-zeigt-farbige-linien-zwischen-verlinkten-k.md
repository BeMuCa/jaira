---
id: 01M45VM4XP8F7P4NK909PAP369
title: Board zeigt farbige Linien zwischen verlinkten Karten
status: critique
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Steht der Cursor auf einer Karte mit Links, sieht man ohne Taste farbige Linien zu jeder verlinkten Karte auf dem Schirm"
context: |-
  Heute sieht man Links (blocked-by, parent, related, follows) nur im L-Fenster, einzeln pro Ticket.
  Berk will (05.10.2026): Cursor auf Karte mit Links -> sofort farbige Linien zu den verlinkten Karten.
  Entschieden von Berk am 05.10.: Linien laufen gerade quer ueber die Karten dazwischen, nicht durch die Raender.
  Entschieden: 4 Farben, eine pro Paar - rot blocked-by/blocks, blau parent/child, gruen follows/followed-by, grau related.
  Geometrie: Spalten liegen direkt aneinander, der Randspalt ist nur '││' (2 Zellen). Karten sind 3 Zeilen hoch und fuellen die Spalte (cardSlots, internal/tui/model.go:1490).
  Overlay wie modalOver: lipgloss Compositor (internal/tui/modal.go:63).
  link.Index.Relations NICHT benutzen: es liest das Logbuch (core/link/link.go:325), zu langsam fuer jede Cursorbewegung.
  Nur direkte Links; Enkel eines Epics bleiben im L-Fenster. Kein Toggle-Key - erst nachruesten, wenn es zu laut ist.
definition-of-done: "Cursor auf Karte mit Link zu sichtbarer Karte: eine Linie je verlinkter Karte in der Farbe ihres Paars; Screen-Test prueft Glyphe und Farbe an den erwarteten Zellen"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-05T11:01:46Z
updated-at: 2026-10-05T11:24:24Z
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-155624
claimed-at: 2026-10-05T11:02:23Z
executed-by: opus
outcome-what: "Board zeichnet von der Cursor-Karte farbige Linien zu jeder direkt verlinkten Karte auf dem Schirm (internal/tui/linklines.go), renderColumn meldet die Mittelzeile jeder Karte, renderBoard legt die Linien per lipgloss Canvas ueber das fertige Board und schreibt eine Legende in die Statusleiste"
outcome-why: "Links waren nur im L-Fenster sichtbar, Ticket fuer Ticket; Berk will beim Navigieren sofort sehen, was womit zusammenhaengt"
outcome-resolves: "DoD 1-2: drawLinkLines routet waagrecht auf der Mittelzeile quer ueber Karten, Knick im Randspalt vor der Zielspalte, gleiche Lane am rechten Rand; Glyphe+Farbe zellgenau in 5 Tests geprueft, live in tmux bestaetigt. DoD 3: ohne Link kehrt renderBoard vor jedem Overlay mit dem unveraenderten String zurueck. DoD 4-5: Legende nennt nur gezeichnete Familien und zaehlt Links ohne Karte; Hoehe der Statusleiste wird mit der laengsten Legende reserviert, damit Kartenpositionen stimmen. DoD 6: Gate m.mode == modeBoard. DoD 7: NOTES.md, -race gruen"
---

# Board zeigt farbige Linien zwischen verlinkten Karten

## Definition of Done

- [x] Cursor auf Karte mit Link zu sichtbarer Karte: eine Linie je verlinkter Karte in der Farbe ihres Paars; Screen-Test prueft Glyphe und Farbe an den erwarteten Zellen
  proof: TestLinkLineRunsStraightToABlockerOnTheSameRow, TestLinkLineFindsABlockerNamedByItsHandle (internal/tui/linklines_test.go); live tmux: fg 203/78/244 an den Linienzellen
- [x] Route: waagrecht auf der Mittelzeile der Cursor-Karte quer ueber Karten dazwischen, Knick im Randspalt vor der Zielspalte, Ecke in die Mittelzeile der Zielkarte; gleiche Lane: senkrecht am rechten Lane-Rand; Test je Fall
  proof: TestLinkLineBendsDownToALowerCard, TestLinkLineRunsLeftToWhatItFollows, TestLinkLineInOneLaneRunsAlongItsBorder (internal/tui/linklines_test.go); drawLinkLines internal/tui/linklines.go
- [x] Karte ohne sichtbaren Link: Board-Ausgabe byte-gleich wie ohne Feature; Test
  proof: TestUnlinkedCardDrawsNothing (internal/tui/linklines_test.go); renderBoard gibt ohne Links b.String() unveraendert zurueck (view.go, 'if len(links) == 0'); alle bestehenden Render-Tests gruen
- [x] Verlinkte Karte nicht auf dem Schirm (gescrollt, gefiltert, schmale Lane, Logbuch): keine Linie, Statusleiste zeigt 'N links off screen · L'; Test
  proof: TestLinkWithoutACardOnScreenIsCounted (internal/tui/linklines_test.go); live: '1 link off screen · L' fuer Karte hinter '+1 more'
- [x] Legende der 4 Farben in der Statusleiste nur solange Linien gezeichnet sind; Test
  proof: TestLinkWithoutACardOnScreenIsCounted (nur gezeichnete Familien), TestUnlinkedCardDrawsNothing (internal/tui/linklines_test.go)
- [x] Linien nur in der Board-Ansicht; Detail, Pipeline, Lane-Fokus und Modals zeichnen keine; Test
  proof: TestLinkLinesOnlyOnTheBoard: enter, L, m (internal/tui/linklines_test.go); Gate 'm.mode == modeBoard' in renderBoard
- [x] core/release/NOTES.md hat eine Zeile unter ## Unreleased; go test ./... -race gruen
  proof: core/release/NOTES.md ## Unreleased; go test ./... -race RC=0

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] renderColumn gibt die Mittelzeile jeder gezeigten Karte zurueck, gemessen am geschriebenen body; tagbox_test auf zwei Rueckgabewerte
- [x] cursorLinks: direkte Links der Cursor-Karte aus m.tickets + m.logged, staerkste Art je Partner (link.Order); Ziele ohne Karte zaehlen als off screen
- [x] failing tests: Glyphe+Farbe nach rechts, nach links, gleiche Lane; ohne Link byte-gleich; off-screen-Hinweis; Legende; nur modeBoard
- [x] drawLinkLines: Route per Canvas CellAt/SetCell, Hintergrund der Zelle bleibt, breite Zeichen werden uebersprungen, Reihenfolge schwach->stark
- [x] renderBoard: nur in modeBoard; Statusleiste reserviert ihre Hoehe mit der laengsten Legende, dann Overlay
- [x] NOTES.md-Zeile unter ## Unreleased
- [x] go test ./... -race, gofmt, Binary nach ~/.local/bin bauen, im tmux pruefen

## Progress
- **2026-10-05 11:06 · BeMuCa** — Pre-process: Statusleiste bricht um und laesst nie eine Taste weg (view.go ~996). Eine Legende darf ihre Hoehe also nicht nach dem Rendern aendern - sonst stimmen die Kartenpositionen nicht mehr. Loesung: Hoehe mit der laengsten moeglichen Legende reservieren (alle Familien + 'N links off screen'), echte Legende danach einsetzen und auf die reservierte Zeilenzahl auffuellen. Kein Doppel-Rendern.
Links nur aus dem Speicher (m.tickets + m.logged), exakter ID-Vergleich wie Relations' Rueckrichtung: Felder speichern volle IDs (geprueft in .jaira/tickets und logbook). Relations selbst liest das Logbuch - nicht pro Render.
Overlay ohne direkten ultraviolet-Import: Canvas.CellAt liefert *uv.Cell, Kopie per 'cc := *c' aendert Content/Fg und behaelt Bg (Kartenband).
- **2026-10-05 11:24 · BeMuCa** — In-progress: 'jaira create --blocked-by SEFFWC' speichert den Handle wie getippt, parent/follows/related dagegen volle IDs (live im Sandbox-Board gesehen). Exakter ID-Vergleich verpasste so den Blocker - Linie fehlte, blaue Ecke lag oben. Fix: refersTo() matcht kuerzere Werte als ID-Ende, wie core/link bySuffix. Test: TestLinkLineFindsABlockerNamedByItsHandle.
Gleiche Luecke vermutlich in core/link Relations (link.go ~376: 'dep == id' exakt): das L-Fenster des Blockers listet 'blocking' fuer einen per Handle gespeicherten Blocker wohl nicht. Nur im Code gelesen, nicht ausgefuehrt; nicht Teil dieses Tickets.
Test-Helfer border() muss Linienglyphen als Rand akzeptieren - eine senkrechte Linie uebermalt den Rand auf der Titelzeile.
