---
id: 01M28YTCVSCEMRN3G0HBANX7Y5
title: Die ausgewaehlte Karte ist gefuellt und oben geschlossen
status: human
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Die ausgewaehlte Karte hebt sich durch eine Fuellfarbe einen Schritt neben dem Terminalhintergrund ab, und ihr Rahmen ist an allen vier Seiten geschlossen, egal an welcher Stelle der Spalte sie steht"
context: |-
  Alex am 11.09.: 'я не могу увидеть где я' - auf einer vollen Spalte ist nicht zu finden, welche Karte der Cursor haelt.

  Heute markiert nur der Titel die Auswahl: fett und in der Tag-Farbe (renderCard, internal/tui/view.go). Auf einer Spalte, in der jede Karte schon eine Tag-Farbe im Rahmen traegt, faellt das nicht auf.

  Zweiter Befund aus demselben Gespraech: die ausgewaehlte Karte hat gar keine Oberkante, wenn sie nicht die erste sichtbare der Spalte ist. renderColumn schneidet ab der zweiten Karte die Top-Border weg - die Karten teilen sich eine Border-Reihe. Die geteilte Reihe gehoert damit der Karte darueber, und die Fuellung der Auswahl beginnt mitten im Kasten.

  Zwei Fallen, beide beim Bauen aufgetaucht und schon geloest:
  - Die Innenstile der Karte (styMeta, styWarn, ...) enden je auf einem vollen SGR-Reset, der auch den Hintergrund der Box loescht. Ohne Nachfuellen nach jedem Reset ist nur das erste Wort gefuellt. Ein Reset am Zeilenende bleibt unangetastet, sonst laeuft die Farbe ueber den rechten Kartenrand hinaus bis zum Terminalrand.
  - 'Heller als der Hintergrund' hat eine Richtung. Auf einem hellen Terminal gibt es nichts Helleres; dort muss die Fuellung dunkler sein. Bubble Tea liefert das ueber tea.BackgroundColorMsg, angefordert in Init.

  Beruehrt VS5DFW im Backlog, das dieselbe Frage anders beantwortet (ganzer Rahmen in Tag-Farbe plus Balken daneben). Kein Konflikt im Code - die Fuellung ist ein eigener Kanal - aber jemand muss sagen, ob beide bleiben.
definition-of-done: "Die ausgewaehlte Karte ist ganzflaechig gefuellt, Rahmen eingeschlossen; ihre Oberkante steht auch dann, wenn sie nicht die erste Karte der Spalte ist; die Fuellung folgt der Terminalfarbe hell/dunkel; die Fuellung laeuft nicht ueber den Kartenrand hinaus; die Kartenhoehe und die Budget-Rechnung bleiben unveraendert; go test ./... -race gruen"
tags:
  - tui
blocked-by: []
commits: []
created-at: 2026-09-11T19:23:56Z
updated-at: 2026-09-11T19:30:22Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1114279
claimed-at: 2026-09-11T19:25:39Z
outcome-what: "Die ausgewaehlte Karte wird ganzflaechig gefuellt, Rahmen eingeschlossen, einen 256-Farben-Schritt neben dem Terminalhintergrund (236 dunkel, 253 hell, Richtung aus tea.BackgroundColorMsg); die geteilte Border-Reihe gehoert jetzt der ausgewaehlten Karte statt der Karte darueber"
outcome-why: "Auf einer vollen Spalte war der fette Titel als einziges Auswahlmerkmal nicht zu finden, weil jede Karte ohnehin eine Tag-Farbe im Rahmen traegt; und ohne die geteilte Reihe war die Auswahl oben offen, die Fuellung begann mitten im Kasten"
outcome-resolves: "Der Cursor ist ohne Lesen zu finden, und die Auswahl ist an allen vier Seiten geschlossen, egal an welcher Stelle der Spalte sie steht. Alex hat den Stand am 11.09. im laufenden Board angesehen und angenommen: 'край появляется когда тикет выбран чтобы не было прыжков - подходит'. go test ./... -race gruen, 24 Pakete"
question: "Bleibt VS5DFW im Backlog bestehen? Es beantwortet dieselbe Frage anders - nicht ausgewaehlt weisser Rahmen mit farbiger linker Kante, ausgewaehlt ganzer Rahmen in Tag-Farbe plus Balken daneben (Berk, 08.09.). Die Fuellung ist ein eigener Kanal und kollidiert im Code nicht damit, aber beide zusammen sind drei Merkmale fuer eine Auswahl. Entweder VS5DFW bleibt wie beschrieben, oder sein Balken-Teil faellt weg, weil die Fuellung ihn ersetzt."
---

# Die ausgewaehlte Karte ist gefuellt und oben geschlossen

## Definition of Done

- [x] Die ausgewaehlte Karte ist ganzflaechig gefuellt, Rahmen eingeschlossen; ihre Oberkante steht auch dann, wenn sie nicht die erste Karte der Spalte ist; die Fuellung folgt der Terminalfarbe hell/dunkel; die Fuellung laeuft nicht ueber den Kartenrand hinaus; die Kartenhoehe und die Budget-Rechnung bleiben unveraendert; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-11 19:25 · Alexander Sacharov** — Zwei Fallen, beide im Bauen gefunden. Erstens: die Innenstile der Karte enden je auf einem vollen SGR-Reset, der auch den Box-Hintergrund loescht - refill in view.go fuellt nach jedem Reset nach, ausser dem letzten einer Zeile, sonst laeuft die Farbe bis zum Terminalrand. Zweitens, fuer Tests: lipgloss legt Vorder- und Hintergrund in EINE SGR-Sequenz, wo eine Zeile beides hat. Auf den Rahmenzeilen steht '38;5;240;48;5;236m', die Fuellung also nie allein - ein Test, der auf die ganze Sequenz '\x1b[48;5;236m' prueft, findet die Oberkante nicht und meldet faelschlich einen offenen Kasten.
