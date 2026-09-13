---
id: 01M21PH7EGYZB6WX4CC481XRXX
title: Gestapelte Karten teilen keine Border-Reihe mehr
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: Zwischen zwei Karten im Board stehen wieder zwei getrennte Rahmen statt einer geteilten Reihe - keine Karte sitzt optisch auf der naechsten
context: |-
  Berk am 08.09. mit Screenshot der Backlog-Spalte: 'die kaestchen sollen nicht overlappen'.
  Das ist eine bewusste Revision von NJPQWE Runde 5.
  Was NJPQWE gebaut hat: ab der zweiten Karte im Fenster faellt deren Top-Border weg, zwei Karten teilen sich eine Border-Reihe (renderColumn), cardsInBudget rechnet gestapelte Karten mit 4 statt 5 Zeilen.
  Warum das damals kam: Berks 4. Screenshot, 'Luecken zwischen den Boxen halbieren'. Gemessen war keine Leerzeile, sondern zwei aneinanderstossende Border-Reihen.
  Was jetzt gilt: die geteilte Reihe liest sich als Ueberlappung. Karten sollen wieder je einen eigenen Rahmen haben.
  Achtung, die Hoehenrechnung haengt dran: cardsInBudget und die Budget-Tests (5+4+4...) sind auf die geteilte Reihe gerechnet und muessen zurueck.
  Achtung, ein Test verbietet das alte Verhalten aktiv: TestStackedCardsShareOneBorderRow in internal/tui verbietet Bottom-ueber-Top-Border boardweit. Der Test dreht sich um.
  Offen und vor dem Bauen zu klaeren ist nichts - Berk hat die Hoehe schon einmal explizit akzeptiert ('das mit den 2 zeilen ist ok, ich will den rahmen').
definition-of-done: Zwischen zwei gestapelten Karten stehen zwei Border-Zeilen; cardsInBudget und die Budget-Tests rechnen wieder mit einheitlicher Kartenhoehe; TestStackedCardsShareOneBorderRow ist umgedreht oder ersetzt; go test ./... -race gruen
tags:
  - tui
blocked-by: []
follows: 01M1KN2HSJ0B32MQ2532NJPQWE
commits:
  - 6bd21c73399d119c3edd146759e47737c2828ef0
  - 784ca787e627eec8e1c3f3fd8ac3006ba5afb748
  - a860e46458d57d2ce80666aed710cba2a8fee06a
created-at: 2026-09-08T23:44:26Z
updated-at: 2026-09-11T20:48:44Z
updated-by: Alexander Sacharov
---

# Gestapelte Karten teilen keine Border-Reihe mehr

## Definition of Done

- [ ] Zwischen zwei gestapelten Karten stehen zwei Border-Zeilen; cardsInBudget und die Budget-Tests rechnen wieder mit einheitlicher Kartenhoehe; TestStackedCardsShareOneBorderRow ist umgedreht oder ersetzt; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-11 20:37 · Alexander Sacharov** — UEBERHOLT durch QQ3EX4, entschieden von Alex am 11.09.: Karten tragen gar keinen Rahmen mehr, es gibt also keine Border-Reihe mehr zu teilen. Berks Anforderung vom 08.09. ('die kaestchen sollen nicht overlappen', 'das mit den 2 zeilen ist ok, ich will den rahmen') ist damit nicht erfuellt, sondern verworfen - der Rahmen kostete zwei Spalten pro Zeile und schnitt jeden Titel bei 17 Zeichen ab. Wer den Rahmen zurueckholt, holt dieses Ticket mit zurueck. Der Test TestStackedCardsShareOneBorderRow ist durch TestStackedCardsAlternateTheirShade ersetzt: was zwei gestapelte Karten trennt, ist jetzt der Wechsel der Fuelltoene.
