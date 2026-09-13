---
id: 01M293037469EFFW67DAQQ3EX4
title: "Karten sind Baender ohne Rahmen, die Tag-Farbe ist eine Zelle links"
status: done
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
commits:
  - 784ca787e627eec8e1c3f3fd8ac3006ba5afb748
created-at: 2026-09-11T20:36:57Z
updated-at: 2026-09-11T20:44:31Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1327969
claimed-at: 2026-09-11T20:38:22Z
question: "Das hier kippt zwei Entscheidungen von Berk: 81XRXX (jede Karte hat einen Rahmen, 'ich will den rahmen', 08.09.) und VS5DFW (die Tag-Farbe sitzt in der linken Rahmenkante). Beide sind auf ihren Tickets als ueberholt vermerkt, aber Berk hat davon noch nichts gehoert - er hat den Rahmen mit Screenshots eingefordert und zweimal nachgeschaerft. Wer sagt es ihm, und akzeptiert er den Tausch Rahmen gegen zwei Spalten Titel und zwei Zeilen Hoehe? Wenn nicht, kommen beide Tickets zurueck und dieses hier faellt."
outcome-what: "Siehe outcome-what des Uebergangs nach human; seitdem unveraendert"
outcome-why: "Siehe outcome-why des Uebergangs nach human"
outcome-resolves: "Von Alex am 11.09. am laufenden Board abgenommen. review uebersprungen, siehe Notiz."
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
- **2026-09-11 20:44 · Alexander Sacharov** — Angenommen von Alex am 11.09. im Gespraech, Stueck fuer Stueck im laufenden Board angesehen. Die review-Lane wurde dabei uebersprungen, und das ist eine bewusste Luecke, keine erledigte Stufe: review heisst 'ein zweites Modell hat den Diff beurteilt', und der Autor des Codes war dasselbe Modell, das ihn haette pruefen sollen. Was stattfand, war menschliche Abnahme am laufenden Bild, nicht Modell-Review. Wer spaeter einen Fehler in diesen drei Tickets sucht: hier ist die Stelle, an der niemand mit frischen Augen draufgeschaut hat.
- **2026-09-11 20:44 · Alexander Sacharov** — Die offene Frage dieses Tickets - wer sagt es Berk - hat jetzt einen Besitzer und einen Zeitpunkt: Alex, nach Berks Urlaub. Als eigenes Ticket VHQ0F4 auf der Bahn, damit es nicht an diesem hier haengt. 81XRXX und VS5DFW bleiben solange im Backlog mit ihrer Ueberholt-Notiz - nicht archiviert, weil das Wegraeumen seiner Tickets vor dem Gespraech die Spur loeschen wuerde.
