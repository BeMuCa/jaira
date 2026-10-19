---
id: 01M59MQDJ13FSPAJWKA9P6XPKY
title: Ein Merge-Konflikt im Lane-Feld entscheidet nach der Zahl der Schritte des Tickets statt nach precedence der Lane
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Wenn zwei Seiten dasselbe Ticket in verschiedene Lanes bewegt haben, gewinnt die Seite, die mehr Schritte gemacht hat - auch wenn das eine Schleife zurueck war; das Feld precedence verschwindet aus den Lane-Dateien."
context: |-
  Heute: core/merge/merge.go mergeStatus (ab Zeile 361) loest einen echten Konflikt (beide Seiten haben status seit der Basis geaendert) ueber lanes.Precedence: die Lane mit der hoeheren Zahl gewinnt.
  Problem (berkc, 19.10.2026): ein Ticket im 2. Durchlauf von critique ist weiter als eins im 1. Durchlauf von testing, aber testing hat die hoehere precedence. Ausserdem muss jede Lane-Datei eine Zahl tragen, die niemand versteht (siehe FGRQCG).
  Vorschlag: ein Zaehler im Ticket (z.B. Feld 'moves'), jeder jaira move erhoeht ihn um 1. Das ist ein Lamport-Zaehler, die uebliche Loesung fuer 'wer ist weiter' ohne zentrale Uhr (Lamport 1978; CRDT-Register mit Lamport-Zeitstempel). Beim Merge: mehr Schritte seit der Basis gewinnt. Gleichstand = wirklich gleichzeitig -> Tie-Break noetig: Reihenfolge in .jaira/lanes/order, dann updated-at.
  Nicht verifiziert, nur aus Fachwissen: die Literaturangabe; vor dem Bau kurz pruefen.
  claim/assignee verhindert den Konflikt nicht: claim ist laut jaira claim --help nur ein Hinweis (30-min-Lease), kein Lock. Zwei Worktrees oder Ref und Branch koennen weiter beide schreiben.
  Alte Tickets ohne Zaehler: als 0 lesen, Tie-Break greift.
definition-of-done: jaira move erhoeht einen Zaehler im Ticket; Test in core/ticket oder internal/cli
tags:
  - concurrency
blocked-by: []
related: []
commits: []
created-at: 2026-10-19T08:33:53Z
updated-at: 2026-10-19T08:33:53Z
---

# Ein Merge-Konflikt im Lane-Feld entscheidet nach der Zahl der Schritte des Tickets statt nach precedence der Lane

## Definition of Done

- [ ] jaira move erhoeht einen Zaehler im Ticket; Test in core/ticket oder internal/cli
- [ ] mergeStatus waehlt bei Konflikt die Seite mit mehr Schritten seit der Basis, bei Gleichstand nach .jaira/lanes/order; Tests fuer Schleife-zurueck-gewinnt und Gleichstand in core/merge
- [ ] precedence wird aus ausgelieferten Lanes und Vorlage entfernt; eine Lane-Datei mit precedence wird weiter gelesen (ignoriert), kein Fehler; Test
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

