---
id: 01M67RYCQZ131NH1BATCCZWKNY
title: "Release-Notes beginnen mit einer Hook-Zeile, und wer mehrere Versionen ueberspringt, sieht die Notes aller dazwischen"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Jede Release-Note-Datei hat als erste Zeile einen Einzeiler, der wie ein Instagram-Hook neugierig macht; jaira update und das Release-Notes-Popup zeigen ihn zuoberst. Wer von 0.5.4 auf 0.5.8 springt, sieht die Notes von 0.5.5, 0.5.6, 0.5.7 und 0.5.8, nicht nur die der neuesten."
context: |-
  Was fehlt: grosse Aenderungen (Modell pro Lane, Default-Worker, Teamlead-Default-Modelle) gehen zwischen vielen kleinen Zeilen unter.
  Ausloeser: Maintainer 31.10.2026, Punkt 13: 'die erste Linie eine Oneliner beschreibung wie eine instagram hook ... Und ausserdem wenn man nicht eine version, sondern mehrere ueberspringt sollte man trotzdem sehen welche notes bei den anderen relevant war.'
  Bekannt: Format ist ein Zeilenscan (CLAUDE.md, Abschnitt Release-Notes): '- ', '- ! ', MXMDFH bringt '- + ' fuer Features. Ein neues Praefix fuer die Hook-Zeile muss alte Leser nicht brechen. Ob update heute schon alle uebersprungenen Versionen zeigt: pruefen (GK1QGT blaettert pro Version).
definition-of-done: "Ein Praefix fuer die Hook-Zeile ist festgelegt, dokumentiert in CLAUDE.md und core/release, und alte Leser zeigen sie als normale Zeile"
tags:
  - release
blocked-by: []
related:
  - 01M67K70X843TBF0TYRPMXMDFH
commits: []
created-at: 2026-10-31T01:24:49Z
updated-at: 2026-10-31T01:24:49Z
---

# Release-Notes beginnen mit einer Hook-Zeile, und wer mehrere Versionen ueberspringt, sieht die Notes aller dazwischen

## Definition of Done

- [ ] Ein Praefix fuer die Hook-Zeile ist festgelegt, dokumentiert in CLAUDE.md und core/release, und alte Leser zeigen sie als normale Zeile
- [ ] jaira update und das Popup zeigen die Hook-Zeilen zuerst; Test
- [ ] Ein Board, das 3 Versionen zurueckliegt, sieht die Notes aller 3; Test
- [ ] Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

