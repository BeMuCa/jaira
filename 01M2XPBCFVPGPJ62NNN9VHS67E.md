---
id: 01M2XPBCFVPGPJ62NNN9VHS67E
title: "Wer eine Lane adoptiert und angepasst hat, verliert sie stillschweigend an die eingebettete Fassung"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer eine Katalog-Lane adoptiert und bearbeitet hat, behaelt seine Fassung, oder jaira sagt ihm, dass die eingebettete gewonnen hat - statt sie still zu ersetzen"
context: |-
  Alex am 2026-09-19, aus der review-Lane von VSC1GW heraus aufgemacht.

  Was ist: Installable() in core/lane/order.go:190-205 legt die eingebetteten Lanes in 'seen', BEVOR es ueber ~/.jaira/lanes globbt. Eine eigene Fassung von 'critique', 'optimize' oder 'testing' unter dieser Id wird dadurch nie gesehen. Nachgestellt am gebauten Binary: mit einem JAIRA_LANES_DIR, dessen critique.md eine andere description traegt, zeigt der Fuss von 'jaira lanes' die EINGEBETTETE Beschreibung, und 'jaira lanes add critique' legt die eingebettete Datei auf die Doska.

  Warum das jetzt stoert: fuer die zehn alten built-ins war das immer so. Neu ist, dass VSC1GW genau die drei Lanes eingebettet hat, die es bis gestern NUR ueber 'jaira lanes market adopt' gab - wer sie also adoptiert und angepasst hat, ist genau der Mensch, den es jetzt trifft, und er hat nichts falsch gemacht.

  Was schon entschieden ist und hier NICHT neu aufgemacht wird: die Vorrangregel selbst bleibt, wie sie ist (Alex am 2026-09-19 auf VSC1GW). Die Zeile unter '## Unreleased' in core/release/NOTES.md sagt dem Betroffenen heute, er solle seine eigene Kopie auf eine eigene Id umbenennen.

  Was offen ist: das Umbenennen ist der einzige gebaute Ausweg, und niemand fuehrt ihn. Der Betroffene erfaehrt nichts, solange er die Release-Notes nicht liest - 'jaira lanes' nennt die verdeckte Datei mit keinem Wort. Moegliche Formen, wachsend: (1) 'jaira lanes' meldet eine verdeckte Datei in ~/.jaira/lanes beim Namen; (2) 'jaira validate' warnt; (3) ein Befehl, der die eigene Kopie auf eine freie Id umschreibt. Welche - gehoert in brainstorm.
definition-of-done: "Ein Mensch mit einer eigenen, angepassten Kopie von critique, optimize oder testing unter ~/.jaira/lanes erfaehrt von jaira selbst, dass sie verdeckt ist - nicht nur aus den Release-Notes. Nachgestellt mit einem JAIRA_LANES_DIR, das eine solche Kopie enthaelt."
tags:
  - cli
blocked-by: []
related:
  - 01M2SNEP2NN02DN5W0ENVSC1GW
commits: []
created-at: 2026-09-19T20:39:59Z
updated-at: 2026-09-19T20:39:59Z
---

# Wer eine Lane adoptiert und angepasst hat, verliert sie stillschweigend an die eingebettete Fassung

## Definition of Done

- [ ] Ein Mensch mit einer eigenen, angepassten Kopie von critique, optimize oder testing unter ~/.jaira/lanes erfaehrt von jaira selbst, dass sie verdeckt ist - nicht nur aus den Release-Notes. Nachgestellt mit einem JAIRA_LANES_DIR, das eine solche Kopie enthaelt.
- [ ] Ein Test haelt die Meldung fest und faellt, wenn sie verschwindet.
- [ ] Eine Zeile in core/release/NOTES.md unter '## Unreleased', wenn sich aendert, was der Benutzer sieht oder tut.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

