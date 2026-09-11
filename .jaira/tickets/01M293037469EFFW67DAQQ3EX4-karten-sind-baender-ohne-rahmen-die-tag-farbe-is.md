---
id: 01M293037469EFFW67DAQQ3EX4
title: "Karten sind Baender ohne Rahmen, die Tag-Farbe ist eine Zelle links"
status: human
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Eine Karte ist ein gefuelltes Band ueber die volle Lane-Breite: kein Rahmen, die Tag-Farbe als ganze Zelle links, und der Hintergrund wechselt von Karte zu Karte, sodass zwei gestapelte Karten sich nicht verlaufen"
context: |-
  Alex am 11.09., nach einem Mockup mit 27 Varianten (Rahmen, Balken links, Balken rechts, Zebra ohne Rahmen).

  Was heute stoerte: der Rahmen kostet zwei Spalten pro Zeile, und der Titel wurde bei 17 Zeichen abgeschnitten - 'Ein Ticket zeigt …', 'Das Logbuch wird …'. Auf dem Board war kein einziger Titel zu Ende lesbar.

  Was das Band statt dessen tut: dieselbe Aussage - wo faengt diese Karte an und wo hoert sie auf - mit Zeilen, die ohnehin gezeichnet werden. Titelbreite geht von w-6 auf w-4, also +2 Spalten bei jeder Lane-Breite. Kartenhoehe geht von 5 auf 3, also zwei Zeilen pro Karte zurueck an die Lane.

  Die Tag-Farbe wandert vom Rahmen in eine volle Zelle. Ein Rahmen-Glyph faerbt etwa eine halbe Zelle und liest sich als 'der Rahmen hat eine andere Farbe', nicht als Markierung - das war im Mockup Variante 2 und fiel durch.

  WICHTIG, das hier kippt zwei Entscheidungen von Berk, und zwar bewusst und von Alex:
  - 81XRXX 'Gestapelte Karten teilen keine Border-Reihe mehr' - es gibt keine Border mehr zu teilen.
  - VS5DFW 'Nur die linke Kante traegt die Tag-Farbe' - die Kante ist jetzt eine Zelle, kein Rahmenteil.
  Berk hatte am 08.09. ausdruecklich 'ich will den rahmen' gesagt. Wer das zurueckdreht, dreht beides zurueck.

  Zwei Dinge, die beim Bauen aufliefen:
  - Die Lane-Innenbreite ist w-2, nicht w: columnStyle.Width zaehlt den eigenen Rahmen mit. Mit w umbrachen alle Titel auf zwei Zeilen und die Karte wurde hoeher als cardHeight verspricht.
  - Der gefuellte Titel war vorher fett in der Tag-Farbe. Auf einem Hintergrund, der aus derselben Farbe gemischt ist, liest sich das matschig - der Titel ist jetzt nur noch fett.

  Die Fuelltoene muessen auseinanderliegen: Lane 234/236, Auswahl 239, Glow gemischt ab Grau 237. selBgDark war vorher 236 und damit identisch mit einem der Zebra-Toene - eine ausgewaehlte Karte ohne farbiges Tag waere auf jeder zweiten Position unsichtbar gewesen.
definition-of-done: Karten tragen keinen Rahmen mehr; die Tag-Farbe fuellt eine ganze Zelle links ueber alle drei Zeilen; der Hintergrund wechselt von Karte zu Karte; kein Fuellton der Auswahl faellt mit einem Lane-Ton zusammen; eine Karte ohne Tag-Farbe behaelt ihre Textausrichtung; cardHeight ist 3 und die Budget-Rechnung stimmt; keine Zeile bricht um; go test ./... -race gruen
tags:
  - tui
blocked-by: []
commits: []
created-at: 2026-09-11T20:36:57Z
updated-at: 2026-09-11T20:38:37Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1327969
claimed-at: 2026-09-11T20:38:22Z
question: "Das hier kippt zwei Entscheidungen von Berk: 81XRXX (jede Karte hat einen Rahmen, 'ich will den rahmen', 08.09.) und VS5DFW (die Tag-Farbe sitzt in der linken Rahmenkante). Beide sind auf ihren Tickets als ueberholt vermerkt, aber Berk hat davon noch nichts gehoert - er hat den Rahmen mit Screenshots eingefordert und zweimal nachgeschaerft. Wer sagt es ihm, und akzeptiert er den Tausch Rahmen gegen zwei Spalten Titel und zwei Zeilen Hoehe? Wenn nicht, kommen beide Tickets zurueck und dieses hier faellt."
outcome-what: "Karten tragen keinen Rahmen mehr, sondern sind gefuellte Baender ueber die volle Lane-Innenbreite; die Tag-Farbe ist eine ganze Zelle links ueber alle drei Zeilen; der Fuellton wechselt von Karte zu Karte; cardHeight faellt von 5 auf 3 und die Fuelltoene sind neu verteilt (Lane 234/236, Auswahl 239) damit nichts zusammenfaellt"
outcome-why: "Der Rahmen kostete zwei Spalten in jeder Zeile und schnitt jeden Titel auf dem Board bei 17 Zeichen ab; das Band sagt dasselbe - wo die Karte anfaengt und aufhoert - mit Zeilen, die ohnehin gezeichnet werden"
outcome-resolves: "Titelbreite w-6 auf w-4, Kartenhoehe 5 auf 3, die Tag-Farbe fuellt eine ganze Zelle statt der halben eines Rahmen-Glyphs, und kein Auswahlton faellt mit einem Lane-Ton zusammen. go test ./... -race gruen, 24 Pakete"
---

# Karten sind Baender ohne Rahmen, die Tag-Farbe ist eine Zelle links

## Definition of Done

- [x] Karten tragen keinen Rahmen mehr; die Tag-Farbe fuellt eine ganze Zelle links ueber alle drei Zeilen; der Hintergrund wechselt von Karte zu Karte; kein Fuellton der Auswahl faellt mit einem Lane-Ton zusammen; eine Karte ohne Tag-Farbe behaelt ihre Textausrichtung; cardHeight ist 3 und die Budget-Rechnung stimmt; keine Zeile bricht um; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-11 20:38 · Alexander Sacharov** — Zwei Tests haben ihre Frage verloren, nicht ihre Antwort, und sind darum ersetzt statt geloescht: TestStackedCardsShareOneBorderRow fragte, ob zwei Border-Reihen aufeinandertreffen - jetzt fragt TestStackedCardsAlternateTheirShade, ob zwei Nachbarn denselben Fuellton bekommen, was dieselbe Gefahr ist ohne Rahmen. TestColumnNeverCutsATaggedCardInHalf zaehlte Rahmen-Glyphen; es zaehlt jetzt die Balken-Zellen, weil die pro Kartenzeile genau einmal gezeichnet werden und damit sagen, wieviele Zeilen wirklich auf den Schirm kamen.
