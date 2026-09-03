---
id: 01M1BQVV8P62NQJMGCNAMR3CN8
title: TestEveryScreenFitsTheTerminal deckt alle Screens ab und misst vor dem Clamp
status: done
ready: true
creator: BeMuCa
goal: "Die fitwindow-Testtabelle prueft auch signoff, edit, settings, lanes, browse, dropboard, message, home und follow-up, und misst Zeilenbreiten auf dem ungeklammerten Render"
context: "Empfehlung des unabhaengigen Reviews zu HREQJR (31.08.): TestEveryScreenFitsTheTerminal (internal/tui/fitwindow_test.go:15-58) deckt board, help, detail, delete, lane focus, compact, projects, move, filter, create ab - keinen der Screens, die der Wrap-Umbau anfasste. Und weil er NACH clampBlock misst, kann er ueberbreite Inhalte prinzipiell nicht sehen (clampBlock schneidet sie vorher ab). Der billige strukturelle Schutz gegen die naechste Instanz der Klasse 'Zeile laeuft rechts raus': Tabelle um die fehlenden Screens erweitern und je Screen zusaetzlich den ungeklammerten Render pruefen (jede Zeile lipgloss.Width <= Breite), wie es wrap_test.go fuer signoff/projects/message vormacht."
definition-of-done: Alle vom Wrap-Umbau geaenderten Screens stehen in der Tabelle; je Screen wird der ungeklammerte Render auf Zeilenbreite geprueft; go test ./... -race gruen
blocked-by: []
commits: []
created-at: 2026-08-31T11:04:25Z
updated-at: 2026-09-03T10:28:51Z
follows: 01M1BNSDH7ZSKJ1GEHEMHREQJR
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-2681933
claimed-at: 2026-08-31T16:51:48Z
assignee: BeMuCa
outcome-what: "home_test.go's TestHomeDoesNotOverflow misst jetzt lipgloss.Width(line) statt len(the(line)); die runenzaehlende Hilfsfunktion the() ist geloescht, da sie keinen anderen Aufrufer hatte"
outcome-why: "Re-Review fand: derselbe Rune-vs-Display-Zellen-Fehler, der in checkLineWidths (fitwindow_test.go) schon behoben wurde, steckte unbehoben in einem zweiten Test desselben Screens (home)"
outcome-resolves: "lipgloss.Width ist jetzt einheitlich das Mass in allen Breitenpruefungen dieses Screens, konsistent mit dem Budget das wrap/truncate/hardBreak selbst verwenden. go test ./internal/tui -count=1 und go test ./... -race -count=1 sind gruen, gofmt-Baseline unveraendert (nur internal/cli/tickets.go)"
review-summary: "none"
review-gaps: "Runde 4: nichts entfernt ausser the() selbst (letzter Aufrufer weg). Alle Funde aus beiden Review-Runden geschlossen; abgelehnt blieb nur das Footer-Reservieren (Feature ueber das Ticket hinaus, im Board notiert)."
review-check: "go test ./internal/tui -run 'TestHomeDoesNotOverflow|TestHomeRendersBeforeItKnowsItsSize|TestEveryScreenFitsTheTerminal|TestScreensWithNoOuterClamp' -count=1 -v: alles gruen. grep -n 'func the(' internal/tui/home_test.go: kein Treffer mehr."
review-verdict: "accept (Opus, 2 Runden + Delta): Clamp-Analyse bestaetigt, Hoehen-Test beisst (Gegenprobe 27/19 Zeilen), lipgloss.Width durchgaengig, the() ohne Waisen geloescht"
---

# TestEveryScreenFitsTheTerminal deckt alle Screens ab und misst vor dem Clamp

## Definition of Done

- [x] Alle vom Wrap-Umbau geaenderten Screens stehen in der Tabelle; je Screen wird der ungeklammerte Render auf Zeilenbreite geprueft; go test ./... -race gruen
  proof: internal/tui/fitwindow_test.go: TestEveryScreenFitsTheTerminal + neue TestScreensWithNoOuterClampFitTheirOwnWidth; go test ./... -race gruen

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Render-Pfade und Test-Setup fuer die 9 fehlenden Screens verifizieren (Explore)
  proof: Explore-Agent + eigene Verifikation via Read/Grep, siehe Notiz
- [x] Tabelle in fitwindow_test.go um signoff, edit, settings, lanes, browse, dropboard, message, home, follow-up erweitern
  proof: internal/tui/fitwindow_test.go: 7 neue Eintraege in der Tabelle (sign-off, edit, settings, lanes, drop-board, message, follow-up)
- [x] Pro Screen: ungeklammerten Render mit langem Inhalt bei 34 und 100 Spalten auf Zeilenbreite pruefen
  proof: TestScreensWithNoOuterClampFitTheirOwnWidth prueft edit/settings/lanes/browse/home direkt bei w=34,100; signoff/message bereits in wrap_test.go direkt geprueft
- [x] go test ./internal/tui laufen lassen, ueberbreite Zeilen identifizieren
  proof: go test ./internal/tui: FAIL zuerst bei lanes (Zeile 229 breit statt 34/100), siehe Notiz
- [x] Bei Ueberbreite: trivialen wrap()/styleLines()-Fix anwenden oder Screen mit TODO ausschliessen + Notiz
  proof: internal/tui/lanes.go:788 truncate() um promptOf-Titelzeile ergaenzt (analog 45f78f2/8f93bbc)
- [x] go test ./... -race und gofmt pruefen, dann committen und Ticket weiterschieben
  proof: go test ./... -race gruen; gofmt -l core internal = nur internal/cli/tickets.go

## Progress
- **2026-08-31 17:09 · BeMuCa** — Model.View() klemmt am Ende IMMER via clampBlock (view.go:67) - dadurch ist ein Test, der ueber m.View() misst, fuer signoff/edit/settings/lanes/message strukturell blind fuer Wrap-Bugs (er kann nur Abstuerze/Hoehe pruefen, nie echte Breite). Deshalb zwei Strukturen: die alte Tabelle bekam alle 7 erreichbaren Screens dazu (sign-off, edit, settings, lanes, drop-board, message, follow-up) fuer Hoehen-/Crash-Abdeckung, aber die echte Breitenpruefung braucht einen NEUEN Test (TestScreensWithNoOuterClampFitTheirOwnWidth), der die Render-Funktionen DIREKT aufruft (wie wrap_test.go es fuer signoff/message schon tut - deshalb dort NICHT dupliziert). browse und home sind ueberhaupt nicht ueber *Model erreichbar (browse haengt nur an Home, Home ist ein eigenes Top-Level-tea.Model ohne jeden Clamp in seinem View()) - beide brauchen zwingend eigene Konstruktion (newBrowser direkt, NewHome(nil) direkt). drop-board und follow-up clampen sich selbst am Ende ihrer eigenen render()-Funktion (dropboard.go:195, followup.go:236) - ein direkter Aufruf waere dort genauso blind wie ueber View(), deshalb nur in der alten Tabelle. Echter Fund: lanes.go:788 (promptOf-Titelzeile 'name (id)') war die einzige ungeschuetzte Zeile auf dem lanes-Screen - Test schlug mit 229 statt 34/100 Zeichen fehl. Fix: truncate() ergaenzt, analog zum Muster aus 45f78f2/8f93bbc.
- **2026-08-31 17:18 · BeMuCa** — review (Opus, batch): accept mit Testfunden - (1) der home-Subcase im neuen Breiten-Test kann nicht fehlschlagen (alles laeuft durch Width(h.width), das selbst wrappt); stattdessen das messen, was Home wirklich nicht schuetzt: HOEHE (kein clampBlock in Home) mit h.stats=true und langer msg. (2) checkLineWidths zaehlt Runen, wrap/truncate budgetieren lipgloss.Width - auf lipgloss.Width umstellen. (3) Die Tabellen-Setups teilen kumulativ Zustand (isErr bleibt gesetzt); t.Run mit frischem Modell je Fall macht Fehler benennbar.
- **2026-08-31 17:41 · BeMuCa** — Review-Korrektur, drei Punkte: (1) home-Subcase in TestScreensWithNoOuterClampFitTheirOwnWidth konnte nie fehlschlagen (Width(h.width)-Style wrappt selbst) - ersetzt durch HOEHEN-Pruefung (h.height=10, h.stats=true, langer msg), die vor dem Fix mit 27/19 statt 10 Zeilen tatsaechlich fehlschlug (verifiziert per git stash: Fix rausgenommen, Test schlaegt fehl, Fix zurueck, Test gruen). (2) checkLineWidths zaehlte Runen statt lipgloss.Width - umgestellt, da wrap/truncate/hardBreak alle in lipgloss.Width budgetieren. (3) Tabellen-Setups in TestEveryScreenFitsTheTerminal teilten ein Model ueber alle 17 Faelle (isErr blieb von 'message' gesetzt, settingsScreen/laneScreen/drop blieben non-nil) - auf t.Run mit frischem newTestModel je Fall umgestellt; der 'follow-up muss zuletzt laufen'-Kommentar ist damit hinfaellig und entfernt, da jeder Fall jetzt isoliert ist. Echter Fund/Fix: home.go's render() endete bisher auf 'return b.String()' ohne jeden Clamp (im Gegensatz zu dropboard/follow-up/lane-focus, die alle selbst clampen) - jetzt 'return clampBlock(b.String(), h.width, h.height)', identisches Muster.
- **2026-08-31 17:43 · BeMuCa** — critique Runde 2: Der Home-Hoehenclamp ist richtig, aber er braucht den Zero-Size-Guard aus Model.View (view.go:66-68) - clampBlock(s, 0, 0) ist leer, und Home rendert vor der ersten WindowSizeMsg. Sonst tauscht der Fix einen Overflow bei kleinen Terminals gegen einen Blank-Flash beim Start.
- **2026-08-31 17:55 · BeMuCa** — Review-Korrektur Runde 2: home.go's neuer clampBlock-Aufruf (aus dem letzten Fix) lief unbedingt, aber clampBlock liefert '' bei Breite ODER Hoehe <= 0. Home bekommt seine Groesse nur ueber WindowSizeMsg (setzt width+height gemeinsam), aber der Clamp selbst kannte diese Symmetrie nicht - abgesichert wie Model.View() es tut (view.go:66-68): if h.width > 0 && h.height > 0 { clampBlock(...) } else b.String(). Wichtiger Nebenfund beim Testen: der neue Test darf NICHT nur (0,0) pruefen, denn render() hat schon VOR meinem Fix einen 'if h.width == 0 { return "loading…" }'-Fallback ganz oben - der faengt (0,0) unabhaengig vom Guard ab. Per Hand verifiziert (Guard rausgenommen, Test lief): bei (0,0) bleibt der Test gruen mit ODER ohne den neuen if-Guard - vacuous fuer genau diesen Fix. Der Fall, der tatsaechlich unterscheidet, ist (80,0) - asymmetrisch, Breite bekannt, Hoehe noch nicht: OHNE Guard kam da '' zurueck, MIT Guard nicht. TestHomeRendersBeforeItKnowsItsSize prueft jetzt beide Faelle, aber (80,0) ist der, der den Fix beweist.
- **2026-08-31 18:02 · BeMuCa** — review Runde 2 (accept mit Funden): (L13) home_test.go the() zaehlt Runen wie das gerade gefixte checkLineWidths - auf lipgloss.Width umstellen, the() loeschen. (L12) abgelehnt vom Koordinator: der Hoehen-Test faellt ohne den Clamp (27/19 Zeilen gegen 10, Gegenprobe belegt) - Regressionsschutz, kein Vakuum; Footer-Zeilen reservieren waere ein neues Feature ueber das Ticket hinaus. (L14) dropboard-Zero-Pfad praktisch unerreichbar - nur notiert.
- **2026-08-31 18:10 · BeMuCa** — Review-Korrektur Runde 3 (klein): home_test.go's TestHomeDoesNotOverflow hatte denselben Fehler wie checkLineWidths vorher - len(the(line)) zaehlt Runen statt Display-Zellen. the() hatte keinen anderen Aufrufer, also geloescht statt nur ungenutzt liegen gelassen. Umgestellt auf lipgloss.Width(line), gleiches Budget wie wrap/truncate/hardBreak es verwenden. Keine Verhaltensaenderung fuer normale ASCII-Inhalte erwartet, aber schuetzt jetzt auch bei Breitzeichen wie die anderen Screens.
