---
id: 01M1EGJFQ2YJ16R844ZPAXTFG3
title: Das Board zeigt Tag-Farben als Karten-Boxen und t blendet die Legende ein
status: done
ready: true
creator: BeMuCa
goal: "Jede Karte mit Tag bekommt eine Box in der Farbe ihres (ersten) Tags; t schaltet ein kleines Popup ein/aus, das Farbe zu Tag-Name zuordnet"
context: "Berk am 01.09., Teil 2 des Tag-Features (folgt auf das Kern-Ticket mit Feld+Registry): auf dem Board sollen die einzelnen Tickets in Boxen umkreist sein, deren Farbe dem Tag entspricht; mit t eine kleine Agenda/Legende als Popup ein- und ausblenden (Farbe -> Tag). Achtung Layout: eine Box um jede Karte kostet ~2 Zeilen Hoehe pro Karte - renderColumn/fitWindow/bodyHeight (internal/tui/view.go) muessen das budgetieren, und die fitwindow-Tests (Vor-Clamp-Messung!) muessen die neuen Formen abdecken. t ist in der Tastenzeile frei (heute: enter v z n m S / ? q). Legende: kleines Overlay wie renderMessage-Muster, toggelt, zeigt Farbfeld+Name aller aktiven Tags des Boards. Farben kommen aus .jaira/tags (Kern-Ticket)."
definition-of-done: "Karten mit Tag sind farbig umrandet (erste Tag-Farbe), ohne Tag unveraendert; t toggelt die Legende, sie listet Farbe+Name; keine Karte und keine Lane laeuft aus dem Terminal (fitwindow-Tests erweitert, Vor-Clamp); Tastenzeile nennt t; go test ./... -race gruen"
blocked-by: []
commits:
  - 441bda210c0df0ccbc796511ad3cb7f8f8d9beaa
  - 84df2a4cf865104a7983723402b547dbb776b1cd
  - d7dcfeeddf35ee09c9e6c45c9f0bf60d82056d62
  - 6c5d81b28e18049c24bafd21a1d639c263474a32
created-at: 2026-09-01T12:54:42Z
updated-at: 2026-09-04T17:30:53Z
follows: 01M1EGJFM7MN1Q08YT1H79GEPW
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-4091163
claimed-at: 2026-09-01T18:58:32Z
assignee: BeMuCa
outcome-what: "Karten mit farbigem Erst-Tag bekommen eine Box in Tag-Farbe (5 statt 3 Zeilen, je Karte budgetiert), t toggelt die Legende (Swatch+Name, farblos als Punkt-Zeile), Tastenzeile nennt t tags; plus Koordinator-Fix: Colour normalisiert den Lookup"
outcome-why: "Berks Wunsch vom 01.09. Teil 2; die /3-Formel haette Box-Karten still halbiert - Budget summiert jetzt echte Hoehen"
outcome-resolves: "go test ./... -race EXIT=0 (voll, Cache geleert); Invariant-Test gegen die alte Formel verifiziert; Untagged-Rendering byte-identisch gepinnt"
review-summary: none
review-verdict: "accept (Koordinator-Review am vollen Diff, Sonnet-Implementierung: Budget-Umbau mit gegen die alte Formel bewiesenem Invariant-Test, Untagged byte-identisch gepinnt, pre-clamp-Faelle board-tags+legend ergaenzt; ein Fund - unnormalisierter Colour-Lookup - selbst gefixt, getestet, mit gruener Voll-Suite committet. Opus war rate-gedeckelt; dein Signoff ist der zweite Mensch-Blick, und die Optik ist ohnehin deine Abnahme)"
review-gaps: "Nichts entfernt. Gelassen: farblose Tags kosten keine Kartenhoehe (Box nur bei Registry-Farbe - Berks Ausweich-Sorge adressiert sich damit teilweise selbst: nur gefaerbte Karten wachsen); q schliesst die Legende statt zu beenden (Popup-Konvention); activeTags liest alle Board-Tickets, nicht nur sichtbare (Legende stabil unter Scroll)."
review-check: "Neu bauen, jaira starten: 1. jaira tag <handle> ui auf 2-3 Tickets, eines davon per jaira set tags=UI (Grossschreibung): alle drei Karten tragen dieselbe farbige Box. 2. t: Legende zeigt Swatch+ui genau einmal; t oder esc schliesst. 3. Karte ohne Tag: unveraendert 3 Zeilen. 4. Spalte mit gemischten Karten bei ~20 Zeilen Terminal: keine halbe Box, ehrliches '+N more'. 5. Tastenzeile nennt 't tags'."
---

# Das Board zeigt Tag-Farben als Karten-Boxen und t blendet die Legende ein

## Definition of Done

- [x] Karten mit Tag sind farbig umrandet (erste Tag-Farbe), ohne Tag unveraendert; t toggelt die Legende, sie listet Farbe+Name; keine Karte und keine Lane laeuft aus dem Terminal (fitwindow-Tests erweitert, Vor-Clamp); Tastenzeile nennt t; go test ./... -race gruen
  proof: renderCardBlock: erste Tag-Farbe -> Box+2 Zeilen, sonst unveraendert (Pin-Test); cardsInBudget budgetiert die Spalte pro Kartenhoehe statt fest /3; t oeffnet/schliesst modeLegend (Farbfeld+Name); 't tags' in statusBar; go test ./... -race EXIT=0

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Code lesen: renderColumn/boardFit/renderCard/columnStyle (view.go), Registry-API (core/tag), Model/reload (model.go) - erledigt
  proof: renderColumn/boardFit/renderCard/columnStyle in view.go:128-494 gelesen; Registry-API in core/tag/tag.go (Load/Colour/Names) gelesen; Model-Struct+reload in model.go:54-330 gelesen
- [x] Registry einmal pro reload() laden, in m.tags cachen (Model-Struct, reload())
  proof: model.go: Model.tags *tag.Registry Feld, reload() laedt tag.Load(root) einmal (Zeilen ~54-61, ~290-296)
- [x] cardColor/cardHeight/renderCardBlock: erste Tag-Farbe -> Box mit BorderForeground (5 Zeilen), keine Farbe -> unveraendertes renderCard (3 Zeilen)
  proof: view.go: cardColor/cardHeight/renderCardBlock; renderCard selbst unveraendert gelassen fuer den Pin-Test
- [x] renderColumn auf variable Kartenhoehe umstellen: cardsInBudget statt festem (h-4)/3, Scroll-Fenster fuer Cursor anpassen
  proof: view.go renderColumn: cardsInBudget statt (h-4)/3; TestColumnNeverCutsATaggedCardInHalf beweist die alte Formel bricht (siehe Notiz)
- [x] t-Taste + modeLegend: Board-Key-Handler, key()-Switch (esc/t/q schliesst), render()-Dispatch, renderLegend (Farbfeld+Name aktiver Board-Tags), Tastenzeile 't tags'
  proof: model.go: modeLegend, Board-Case 't', key()-Case modeLegend (esc/t/q); view.go: render()-Dispatch, renderLegend, 't tags' in statusBar keys
- [x] Tests: unveraenderte Karte gepinnt, gefaerbte Box+Hoehe, farblose Tag-Karte ohne Box, fitwindow PRE-clamp bei 40/80 Spalten mit getaggten Karten, Legende toggelt+listet, Tastenzeile nennt t, hints_test erweitert
  proof: internal/tui/tagbox_test.go (neu) + fitwindow_test.go/hints_test.go erweitert, siehe go test Lauf
- [x] go test ./... -race -count=1, gofmt -l core internal pruefen
  proof: go test ./... -race -count=1 EXIT=0; gofmt -l core internal -> nur internal/cli/tickets.go (vorbestehend)
- [x] Notiz schreiben, committen
  proof: Notiz geschrieben; commit 441bda2 feat(AXTFG3): tag colour boxes on cards, t toggles a legend

## Progress
- **2026-09-01 19:23 · BeMuCa** — lipgloss: Style.Width(n) mit Border() ergibt eine Gesamtbreite von n+2 (Border kommt oben drauf). Deshalb: Kartenbudget w-4 aus renderColumn bleibt gleich, Box bekommt Width(w-4-2), damit boxed und unboxed Karte exakt dieselbe Gesamtbreite w-4 haben.

Farben rendern als ANSI-Codes auch unter go test (kein TTY-Downgrade zu NoColor) - kurz mit einer Probe verifiziert (\x1b[38;5;Nm). Tests koennen also direkt auf '38;5;<farbe>' im ungestrippten Output pruefen statt nur auf Struktur.

Kernproblem beim Hoehenbudget: renderColumn haengt clampBlock(body, w, h) hinten dran, das nur nach Zeilenzahl abschneidet - nicht kartenweise. Ein zu niedrig geschaetztes Budget (z.B. die alte 'visible = (h-4)/3'-Formel unveraendert weiterverwendet fuer 5-zeilige Box-Karten) fuehrt NICHT zu einem Crash oder einer falschen Breite, sondern schneidet eine Karte mitten in der Box ab (oeffnende Ecke ohne schliessende) - das ist genau das 'Karte faellt still raus'-Risiko aus dem Ticket. Deshalb cardsInBudget: summiert echte Kartenhoehen statt Konstante/3. Verifiziert, dass der neue Test TestColumnNeverCutsATaggedCardInHalf bei der alten Formel tatsaechlich rot wird (manuell zurueckgepatcht und wieder hergestellt) - reine Breiten/Hoehen-Bounds-Tests wie TestEveryScreenFitsTheTerminal koennen das NICHT faenger, weil View()'s eigener clampBlock() sie strukturell immer erfuellt (Bug-Risiko aus dem Ticket: 'assertion die nicht fehlschlagen kann').

Entscheidung (wie im Ticket als Lean vorgeschlagen): Tag ohne Registry-Farbe -> keine Box, keine Farbe, rendert wie ungetaggt (renderCardBlock == renderCard). Der Tag zaehlt aber in der Legende, als farblose Zeile ('· name (no colour)'). Erwog Alternative 'faint default box fuer jeden Tag' - verworfen, weil das Kartenhoehe fuer JEDE getaggte Karte auf 5 zwingen wuerde, auch ohne echte Faerbung im Board - widerspricht dem Ticket-Lean direkt.

Legende zeigt Tags aus m.tickets (ganzes Board), nicht nur der sichtbaren Spalte/Scroll-Position - stabiler beim Blaettern, entspricht 'aktive Tags des Boards' aus dem Kontext woertlich.
- **2026-09-01 19:26 · BeMuCa** — critique (Koordinator, adversarial - Opus gedeckelt): ein Fund, selbst gefixt und committet: Registry.Colour normalisierte den Lookup nicht (handgesetztes UI bekam keine Box, Legende zeigte ui farbig + UI farblos) - Colour normalisiert jetzt, activeTags dedupliziert normalisiert, Test gepinnt. Sonst sauber: Karten-Budget summiert echte Hoehen statt /3 (Invariant-Test faellt gegen die alte Formel), Registry einmal je reload geladen, farblose Tags kosten keine Hoehe, t kollidiert nicht.
- **2026-09-01 19:31 · BeMuCa** — Koordinator-Fehler, ehrlich: 84df2a4 wurde mit nicht kompilierendem Testfile gepusht - die Guard-Kette (set -e + test) stoppte nicht; Ursache im Testcode: reg.Set gibt nichts zurueck, als error behandelt. Fix-forward committet, Guards ab jetzt explizit mit || exit. Voller -race-Lauf danach gruen.
