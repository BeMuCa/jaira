---
id: 01M0YT46B9ECSRW42D79QPJNQP
title: Leere Lane schmal zeichnen statt sie zu verstecken
status: done
ready: true
creator: BeMuCa
goal: "z laesst leere Lanes stehen, nur schmal, statt sie vom Board zu nehmen"
context: "Berk am 25.08.: verstecken kann verwirren - man sieht nicht, dass die Lane existiert. Heute nimmt boardFit (internal/tui/view.go) leere Lanes aus der Spaltenliste, wenn hideEmpty gesetzt ist, und die Statuszeile sagt wie viele weg sind. Alternative: Lane bleibt stehen, colW auf ein Minimum klemmen, Titel gekuerzt. Betrifft boardFit, renderColumn und die z-Tests des Beitragenden in internal/tui/lanewindow_test.go."
definition-of-done: Eine leere Lane ist mit z sichtbar aber schmal; kein Ticket verschwindet vom Board; die bestehenden z-Tests sind angepasst statt geloescht
blocked-by: []
commits:
  - bd5b2eec3587b869a89c7b2db233efd71ea595db
  - f629024b168fe45d1d760a1d438c3a78e28ef760
  - 0c5423cb76c8576413d7cd4868a4602ea16ecf76
  - 6c5d81b28e18049c24bafd21a1d639c263474a32
created-at: 2026-08-26T10:33:48Z
updated-at: 2026-09-03T15:48:25Z
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-3443569
claimed-at: 2026-08-27T15:35:25Z
assignee: BeMuCa
outcome-what: "laneWindow aus view.go nach centre_test.go verschoben; fitWindow traegt die Zentrier-Begruendung jetzt selbst. Bericht praezisiert: zwei z-Tests angepasst (TogglingOn...RightFirst -> TogglingOnLeavesTheCursorOnAnEmptyLane, das zugleich die fokussierte leere Lane als thin prueft), einer ersatzlos (LeftAsFallback: die Cursor-Suche existiert nicht mehr)."
outcome-why: "Produktionscode ohne Produktionsaufrufer ist ein Shim am falschen Ort; die Tests fragen weiter in Einheitskosten, also gehoert der Uebersetzer zu ihnen."
outcome-resolves: "DoD unveraendert erfuellt; Sweep-Test und alle anderen gruen mit -race, Binary neu gebaut. Der Diff ist um eine Funktion kleiner als der, den die Kritik gelesen hat."
review-summary: "Der Diff aendert eine Sache: was z mit leeren Lanes macht. Vorher wurden sie aus boardFits Spaltenliste genommen und die Statuszeile zaehlte sie; jetzt bleiben alle Lanes in der Liste, leere kosten 6 statt 24 Zellen im Budget und werden von renderThinColumn als senkrechter Name gezeichnet. Die Fensterberechnung ist verallgemeinert (fitWindow mit Kosten pro Lane), fuer gleiche Kosten nachweislich identisch. Vier Mechanismen, die nur wegen der Unsichtbarkeit existierten, sind entfernt. Drei Commits: bd5b2ee Feature, f629024 Test-Shim aus dem Produktionscode, 0c5423c Rahmen-Setup gefaltet."
review-gaps: "Gefaltet: renderThinColumn wiederholte die vier Zeilen Rahmen-Setup aus renderColumn (Border, faint, Akzent bei Fokus) - jetzt einmal in columnStyle(idx, w, h), beide rufen es. Gesucht und nicht gefunden: eine zweite Stelle, die Lanes budgetiert oder ein Fenster berechnet (grep fitWindow/laneWindow/perScreen: nur boardFit und die Tests). Toter Code: boardNotice und der hidden-Zaehler sind mit dem Verstecken weggefallen; laneWindow ist in der Kritik-Runde in die Testdatei gewandert. Stehen gelassen und benannt: die vollen Spalten werden seit jeher 2 Zellen breiter budgetiert als sie rendern (lipgloss Width ist die Gesamtbreite) - vorbestehend, der Kommentar an der einen Stelle sagt es jetzt, das Layout ist nicht angefasst. Ebenfalls stehen gelassen: gofmt -l listet weiterhin internal/cli/tickets.go. Kosten: fitWindow laeuft pro Render ueber hoechstens ein Dutzend Lanes - nichts zu holen."
review-check: |-
  1. Im jaira-Repo: go build -o ~/.local/bin/jaira ./cmd/jaira - keine Ausgabe.
  2. Terminal auf etwa 120 Spalten Breite bringen, dann 'jaira' starten. Unten in der Tastenzeile steht 'z thin empty'. Nicht alle 12 Lanes sind sichtbar.
  3. z druecken. Jede Lane ohne Ticket wird 4 Zellen schmal, ihr Name steht senkrecht, Buchstabe unter Buchstabe, kein Zaehler. Die Lanes mit Tickets sind breiter geworden. Es passen mehr Lanes auf den Schirm als vorher. Die Tastenzeile sagt jetzt 'z widen empty'. Nirgends steht 'hidden'.
  4. Mit l nach rechts laufen: der Cursor bleibt auf einer schmalen Lane stehen (ihr Rahmen faerbt sich), statt sie zu ueberspringen.
  5. Noch einmal z: alle Lanes wieder gleich breit, Namen waagerecht.
  6. Gegenprobe: vor dem Bauen mit einem Binary von 7a861b3 (git worktree add, go build) dieselben Schritte - dort verschwinden leere Lanes nach z und unten steht 'N lane(s) hidden - z to show'.
  7. go clean -testcache && go test ./internal/tui/ -run 'TestZ|TestToggling|TestLWithToggle|TestFocusedLane|TestTheWindow' -v - alle PASS.
review-verdict: "Angenommen, mit derselben Einschraenkung wie bei 4DQPMS: dieselbe Sitzung hat gebaut, kritisiert, optimiert und beurteilt, kein zweites Modell. Geprueft, nicht behauptet: go test ./... -race mit leerem Cache gruen nach jedem der drei Commits; der Sweep-Test in centre_test.go beweist, dass fitWindow mit Einheitskosten die alte Klemmformel reproduziert; ein Render-Dump bei 200x24 mit z zeigt alle 10 Lanes mit senkrechten Namen. Verhaltensaenderungen, alle gewollt: leere Lanes bleiben sichtbar (das Ziel), h/l landen auf ihnen, der Cursor wird beim Einschalten nicht mehr verschoben, der 'hidden'-Hinweis ist weg. Nicht gemacht: das seit jeher zu grosse Budget der vollen Spalten - benannt in review-gaps, gehoert nicht zu diesem Ticket."
---

# Leere Lane schmal zeichnen statt sie zu verstecken

## Definition of Done

- [x] Eine leere Lane ist mit z sichtbar aber schmal; kein Ticket verschwindet vom Board; die bestehenden z-Tests sind angepasst statt geloescht
  proof: internal/tui/view.go boardFit/fitWindow/renderThinColumn; TestZDrawsEmptyLanesThinAndKeepsThemAll, TestZMakesRoomForLanesThatWereOffScreen; Render-Dump bei 200x24 mit z: alle 10 Lanes, Brainstorm/Pre-process senkrecht

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-08-26 17:07 · BeMuCa** — Entschieden am 26.08.: leere Lanes werden schmal gezeichnet, nicht versteckt. Nicht die Variante 'schmal per Default, z versteckt weiter' - das war nicht die Wahl. Betrifft boardFit, renderColumn und die Tests des Contributors fuer z.
- **2026-08-27 15:46 · BeMuCa** — Zwei Dinge, die das Repo nicht sagt:

1. lipgloss v2 Width() ist die Gesamtbreite INKLUSIVE Rahmen, gemessen: Width(22) rendert 22 Zellen, Width(4) rendert 4. Der alte Kommentar in boardFit ('content width plus two border cells') stimmte nicht; die vollen Spalten reservieren seit jeher 2 Zellen zu viel pro Spalte. Nicht angefasst, nur der Kommentar an meiner Stelle korrigiert. Fuer die schmale Spalte heisst das: Width(3) laesst ' B' umbrechen und erzeugt eine Leerzeile pro Buchstabe - deshalb thinColWidth = 4.

2. Schmal heisst: Lane-Name senkrecht, ein Buchstabe pro Zeile, kein Zaehler. Verworfen: Name waagerecht gekuerzt - bei 12 Zeichen langen Namen (Implementing, Human Review) spart das gegenueber 22 fast nichts, und gekuerzt ist der Name nicht mehr lesbar. Senkrecht spart 18 Zellen pro leerer Lane; bei 5 leeren auf 200 Spalten passen alle 10 statt 8.

Fenster: laneWindow(idx,n,perScreen) ist jetzt fitWindow(idx, costs, budget) mit Einheitskosten - waechst vom Cursor aus abwechselnd links, rechts. Fuer gleiche Kosten identisch mit der alten Klemmformel (der Sweep-Test in centre_test.go beweist es), fuer gemischte Breiten die einzige Definition von 'wie viele passen'.

Weggefallen, absichtlich: die Keep-Regel fuer die fokussierte leere Lane, das Cursor-Verschieben beim Einschalten und das Ueberspringen leerer Lanes bei h/l - alles existierte nur, weil Lanes unsichtbar wurden. Sichtbar ist ansteuerbar.
- **2026-08-27 15:46 · BeMuCa** — critique: laneWindow nach centre_test.go verschieben und aus view.go entfernen - Produktionscode ohne Produktionsaufrufer. Und den Bericht praezisieren: zwei Tests angepasst, einer ersatzlos, weil das Verhalten (Cursor-Suche) nicht mehr existiert.
