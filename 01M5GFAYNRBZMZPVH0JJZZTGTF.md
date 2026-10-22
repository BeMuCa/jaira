---
id: 01M5GFAYNRBZMZPVH0JJZZTGTF
title: "Das Warnzeichen unten links sagt, was los ist, statt nur eine Zahl zu zeigen"
status: ideas
ready: true
creator: berkc
assignee: berkc
goal: "Wer unten links das Warnzeichen sieht, erfaehrt ohne Suchen, was falsch ist und wie man es behebt."
context: |-
  Gemeldet von berkc am 22.10.2026: 'Sometimes there are exclamation marks on the bottom left but nothing is said, so we dont get the information but the signal.'
  Was es ist (geprueft im Code): internal/tui/view.go:1067-1069 zeigt '⚠ N' vor der Tastenzeile, N = Anzahl m.warnings. Die kommen aus den Lane-Dateien (lanes.Warnings, model.go:426) und aus Parse-Problemen (model.go:454). Der Text steht nur in der Hilfe (?) unter 'Warnings' (view.go:1795-1799) - darauf weist nichts hin.
  Bei berkc ist es sehr wahrscheinlich .jaira/lanes/new-lane-2.md: 'missing required field id' (kaputte Lane, siehe FGRQCG).
  Optionen:
  A) Die erste Warnung gekuerzt direkt neben dem Zeichen: '⚠ 1 new-lane-2.md: missing id'.
  B) Eine Taste (z.B. !) oeffnet ein Fenster mit allen Warnungen, jede mit Datei und dem Befehl, der sie behebt; die Tastenzeile sagt '⚠ 1 (! zeigen)'.
  C) Eine neue Warnung blinkt einmal als Meldung (flash), danach nur das Zeichen.
  Empfehlung: A + B. A sagt sofort das Wichtigste, B ist der Ort fuer mehrere.
definition-of-done: Neben dem Warnzeichen steht die erste Warnung gekuerzt; Test in internal/tui
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-22T00:14:22Z
updated-at: 2026-10-22T00:14:22Z
---

# Das Warnzeichen unten links sagt, was los ist, statt nur eine Zahl zu zeigen

## Definition of Done

- [ ] Neben dem Warnzeichen steht die erste Warnung gekuerzt; Test in internal/tui
- [ ] Eine Taste oeffnet alle Warnungen mit Datei und Abhilfe; die Tastenzeile nennt sie; Test
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

