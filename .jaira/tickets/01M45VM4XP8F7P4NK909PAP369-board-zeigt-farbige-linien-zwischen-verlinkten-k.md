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
updated-at: 2026-10-05T11:50:32Z
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-155624
claimed-at: 2026-10-05T11:02:23Z
executed-by: opus
outcome-what: "Critique-Runde 1 umgesetzt: staerkste Art je Partner, Positionen per Zeiger, eindeutige Kurz-Refs, Legende nur in vorhandenem Platz der Statusleiste (fitLegend)"
outcome-why: "Kritik: blocked-by+related wurde grau, Board sprang bei jeder Cursorbewegung auf eine verlinkte Karte, Logbuch-Zwilling stahl den Linienstart, mehrdeutige Kurz-Refs zeichneten zu einem beliebigen Ticket"
outcome-resolves: "Je Befund ein Test, jeder per Mutation als scharf bewiesen (F1 TwoFields, F2 DoesNotMoveTheBoard bei 80/100/150/200 Spalten, F3 CopyUnderTheCursor, F4 AmbiguousHandle); go test ./... -race RC=0; live in tmux: unterer Board-Rand gleiche Zeile mit/ohne Links bei 100 und 180 Spalten"
review-summary: "Grundform stimmt: Kartenpositionen werden in renderColumn gemessen statt gerechnet, Overlay ueber die lipgloss Canvas wie modalOver (modal.go:63); Positionen stimmen in allen geprueften Layouts. Falsch ist dreierlei: (1) cursorLinks verliert die staerkste Art, wenn ein Ticket denselben Wert in zwei Feldern nennt (blocked-by + related -> grau statt rot). (2) Die Legende sitzt im prefix der Statusleiste; wrapHints macht damit jede Hinweiszeile schmaler, das Board schrumpft bei 80 Spalten um bis zu 9 Zeilen und springt bei jeder Cursorbewegung auf/von einer verlinkten Karte. (3) Positionen nach ID statt Zeiger geschluesselt, obwohl logbook.go:33 Board- und Logbuch-Zwillinge per Zeiger trennt. Optional: mehrdeutige Kurz-Refs matchen anders als core/link bySuffix; die Polsterschleife der Statusleiste ist wirkungslos. Zurueck nach in-progress."
---

# Board zeigt farbige Linien zwischen verlinkten Karten

## Definition of Done

- [x] Cursor auf Karte mit Link zu sichtbarer Karte: eine Linie je verlinkter Karte in der Farbe ihres Paars; Screen-Test prueft Glyphe und Farbe an den erwarteten Zellen
  proof: TestLinkLineRunsStraightToABlockerOnTheSameRow, TestLinkLineFindsABlockerNamedByItsHandle, TestTwoFieldsNamingOneTicketKeepTheStrongerKind, TestLinkLineStartsAtTheCopyUnderTheCursor (internal/tui/linklines_test.go); live tmux: fg 203/78/244
- [x] Route: waagrecht auf der Mittelzeile der Cursor-Karte quer ueber Karten dazwischen, Knick im Randspalt vor der Zielspalte, Ecke in die Mittelzeile der Zielkarte; gleiche Lane: senkrecht am rechten Lane-Rand; Test je Fall
  proof: TestLinkLineBendsDownToALowerCard, TestLinkLineRunsLeftToWhatItFollows, TestLinkLineInOneLaneRunsAlongItsBorder (internal/tui/linklines_test.go); drawLinkLines internal/tui/linklines.go
- [x] Karte ohne sichtbaren Link: Board-Ausgabe byte-gleich wie ohne Feature; Test
  proof: TestUnlinkedCardDrawsNothing (internal/tui/linklines_test.go); renderBoard gibt ohne Links b.String() unveraendert zurueck (view.go, 'if len(links) == 0'); alle bestehenden Render-Tests gruen
- [x] Verlinkte Karte nicht auf dem Schirm (gescrollt, gefiltert, schmale Lane, Logbuch): keine Linie, Statusleiste zeigt 'N links off screen · L'; Test
  proof: TestLinkWithoutACardOnScreenIsCounted, TestAmbiguousHandleDrawsNoLine (internal/tui/linklines_test.go); Legende nur in vorhandenem Platz (fitLegend, TestLinkedCardDoesNotMoveTheBoard 80-200 Spalten); live 100 Spalten: '2 links off screen · L'
- [x] Legende der 4 Farben in der Statusleiste nur solange Linien gezeichnet sind; Test
  proof: TestLinkWithoutACardOnScreenIsCounted, TestUnlinkedCardDrawsNothing, TestFitLegendTakesTheFirstThatFits (internal/tui/linklines_test.go)
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
- **2026-10-05 11:50 · BeMuCa** — Critique-Runde 1 umgesetzt: F1 staerkste Art gewinnt (best-Map statt named-Map), F3 Positionen per Zeiger (renderColumn liefert map[*ticket.Ticket]int; Board-Kopie vor Logbuch-Zwilling), F4 Kurz-Ref passt nur eindeutig, sonst kein Ziel (wie bySuffix), F5 Polsterschleife entfaellt.
F2 anders geloest als vorgeschlagen: Vorschlag 'Legende als erste keys' haette auf schmalen Terminals weiter 1 Zeile Sprung gelassen. Jetzt fitLegend: Legende rechtsbuendig in die letzte Zeile der Statusleiste, nur wenn Platz ist; sonst nur 'N links off screen · L'; sonst nichts. Hoehe der Leiste aendert sich nie -> Board springt nie, Reservierung entfaellt. Preis: auf sehr schmalen Terminals fehlt die Legende.
