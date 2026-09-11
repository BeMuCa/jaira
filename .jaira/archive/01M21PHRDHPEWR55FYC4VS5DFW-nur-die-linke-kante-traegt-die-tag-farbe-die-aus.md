---
id: 01M21PHRDHPEWR55FYC4VS5DFW
title: "Nur die linke Kante traegt die Tag-Farbe, die ausgewaehlte Karte wird ganz farbig"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: Eine nicht ausgewaehlte Karte hat einen weissen Rahmen mit farbiger linker Kante; die ausgewaehlte Karte hat den kompletten Rahmen in der Tag-Farbe und traegt zusaetzlich einen Strich neben sich
context: |-
  Berk am 08.09. mit zwei Screenshots.
  Heute ist der ganze Kartenrahmen in der Tag-Farbe - jede Karte in der Spalte leuchtet, die Spalte wird unruhig.
  Gewuenscht: nur die linke Kante der Box traegt die Tag-Farbe. Der Rest des Rahmens ist weiss.
  Erst wenn eine Karte ausgewaehlt ist, wird ihr ganzer Rahmen zur Tag-Farbe.
  Zusaetzlich, aus demselben Screenshot-Paar: 'bitte doch noch ein strich neben der ausgewaehlten' - die ausgewaehlte Karte bekommt einen Strich/Balken neben sich.
  Das revidiert NJPQWE: dessen review-check Schritt 4 haelt ausdruecklich fest 'Selektierte Karte: Titel gefaerbt und fett, KEIN Balken'. Der Balken kommt jetzt zurueck.
  Karten ohne Tag: NJPQWE hat entschieden, dass auch sie eine Box bekommen, in neutraler Rahmenfarbe (colFaint). Diese Regel bleibt - die linke Kante ist dann eben auch neutral.
  Ort: internal/tui/view.go, renderCardBlock. NJPQWE hat Stapelung bewusst in renderColumn gelegt und renderCardBlock unberuehrt gelassen - dieses Ticket ist genau umgekehrt, es faellt in renderCardBlock.
  Haengt zusammen mit 81XRXX, das die Stapelung angeht. Getrennt gehalten: Layout dort, Farbe hier.
definition-of-done: Eine nicht ausgewaehlte Karte zeigt weissen Rahmen mit farbiger linker Kante; die ausgewaehlte zeigt den ganzen Rahmen in Tag-Farbe plus Strich daneben; tag-lose Karten bleiben neutral geboxt; Tests decken beide Zustaende; go test ./... -race gruen
tags:
  - tui
blocked-by:
  - 01M21PH7EGYZB6WX4CC481XRXX
follows: 01M1KN2HSJ0B32MQ2532NJPQWE
commits:
  - 6bd21c73399d119c3edd146759e47737c2828ef0
  - 784ca787e627eec8e1c3f3fd8ac3006ba5afb748
  - a860e46458d57d2ce80666aed710cba2a8fee06a
created-at: 2026-09-08T23:44:43Z
updated-at: 2026-09-11T20:48:48Z
updated-by: Alexander Sacharov
---

# Nur die linke Kante traegt die Tag-Farbe, die ausgewaehlte Karte wird ganz farbig

## Definition of Done

- [ ] Eine nicht ausgewaehlte Karte zeigt weissen Rahmen mit farbiger linker Kante; die ausgewaehlte zeigt den ganzen Rahmen in Tag-Farbe plus Strich daneben; tag-lose Karten bleiben neutral geboxt; Tests decken beide Zustaende; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-11 20:37 · Alexander Sacharov** — UEBERHOLT durch QQ3EX4, entschieden von Alex am 11.09.: die Tag-Farbe sitzt jetzt in einer ganzen gefuellten Zelle links, nicht in der linken Kante eines Rahmens - es gibt keinen Rahmen mehr. Berks Wunsch vom 08.09. ist inhaltlich erfuellt (die Farbe steht links, nicht rundum), aber anders umgesetzt als beschrieben. Im Mockup war die Rahmenkanten-Variante Nummer 2 und fiel durch: ein Border-Glyph faerbt etwa eine halbe Zelle und liest sich als 'der Rahmen hat eine andere Farbe', nicht als Markierung. Der zusaetzlich gewuenschte Balken neben der ausgewaehlten Karte entfaellt ebenfalls - die Auswahl ist jetzt ganzflaechig gefuellt und im Ton ihres Tags (ANX7Y5, XK4124), ein Balken daneben waere das dritte Merkmal fuer dieselbe Sache.
