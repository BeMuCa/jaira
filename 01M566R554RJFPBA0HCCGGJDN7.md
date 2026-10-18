---
id: 01M566R554RJFPBA0HCCGGJDN7
title: "Lanes lassen sich direkt auf dem Board verschieben, nicht nur im Lanes-Bildschirm der Einstellungen"
status: backlog
ready: true
creator: berkc
assignee: BeMuCa
goal: "Auf dem Board verschiebt eine Taste die fokussierte Lane eine Spalte nach links oder rechts, sodass man den Fluss (a->b->c zu a->c->b) dort aendert, wo man ihn sieht."
context: |-
  Berk am 18.10.: 'Wenn man auf einem Board ist und die Lanes sieht, kann man sie nicht verschieben. Man sollte sie auswaehlen und verschieben koennen, um den Fluss des Boards zu aendern.'
  Was es schon gibt: im Lanes-Bildschirm (S -> Lanes, internal/tui/lanes.go:277-283) verschieben H/L die gewaehlte Lane ueber lane.MoveLane; auf der CLI 'jaira lanes move <id> --left|--right'. Auf dem Board selbst fehlt die Taste.
  Gespeichert wird die Reihenfolge schon zentral: .jaira/lanes/order, eine Lane-ID pro Zeile (core/lane/order.go writeIDList). Das 'precedence: 0/5/.../50' in jeder Lane-Datei ist NICHT die Anzeige-Reihenfolge, sondern der Rang, mit dem ein Merge entscheidet, welcher Stand weiter ist (core/lane/lane.go:48-52, 'deliberately separate from display order'). Ein JSON statt der order-Datei ist deshalb nicht noetig; die Textdatei bleibt von Hand editierbar und diff-lesbar (Projekt-Constraint).
  Offen: welche Taste (H/L wie im Lanes-Bildschirm? auf dem Board ist H/L evtl. belegt - pruefen); ob danach precedence angepasst werden muss, wenn der Fluss sich wirklich umdreht (b vor c heisst auch: b ist weniger weit als c) - im brainstorm klaeren.
definition-of-done: Auf dem Board verschiebt eine Taste die fokussierte Lane nach links/rechts und schreibt .jaira/lanes/order (bzw. die order-Datei des aktiven Boards) ueber lane.MoveLane
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-18T00:31:54Z
updated-at: 2026-10-18T00:31:54Z
---

# Lanes lassen sich direkt auf dem Board verschieben, nicht nur im Lanes-Bildschirm der Einstellungen

## Definition of Done

- [ ] Auf dem Board verschiebt eine Taste die fokussierte Lane nach links/rechts und schreibt .jaira/lanes/order (bzw. die order-Datei des aktiven Boards) ueber lane.MoveLane
- [ ] Die Hilfe (?) und die Statusleiste nennen die Taste
- [ ] Ob precedence beim Umdrehen des Flusses mitwandern muss, ist entschieden und umgesetzt oder begruendet abgelehnt
- [ ] Test; Release-Note

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

