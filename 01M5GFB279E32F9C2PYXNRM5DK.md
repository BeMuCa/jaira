---
id: 01M5GFB279E32F9C2PYXNRM5DK
title: "Das Logbook verdichtet sich nach Quartalen: Ticket -> Logbook-Eintrag -> Quartalsuebersicht"
status: ideas
ready: true
creator: berkc
assignee: berkc
goal: "Man sieht pro Quartal auf einer Seite, woran gearbeitet wurde, und geht von dort in die Logbook-Eintraege und weiter bis zum einzelnen Ticket."
context: |-
  Wunsch von berkc (19. und 22.10.2026): eine Baumstruktur - Tickets als unterste, genaueste Knoten, das Logbook als Zusammenfassung dessen, was dort getan wurde, und Zeitraeume (Q1..Q4) fassen die Logbooks zusammen. 'Then if we want to use it we can see in which Quarter what was touched, and if we want more precise information we go deeper.' Am 22.10.: 'yeah i want the ticket'.
  Stand (geprueft): jaira logbook verschiebt die GANZE Ticket-Datei (nicht nur Commits) nach .jaira/logbook/<kuerzel>-<datum>/ und stempelt vorher alle Commits darauf. Der Ref bleibt bestehen und bekommt den Endstand (RecordFiled -> Record, core/refsync/refsync.go:191); nur jaira delete entfernt einen Ref. jaira logbook ohne Argument listet die letzten 28 Tage (logbook-days).
  Es gibt schon die Rolle jaira-role-logbook-summary (core/role/builtin/) - vor dem Bau lesen, was sie tut.
  Vorschlag aus der Diskussion: die Quartalsuebersicht bei Bedarf erzeugen (jaira logbook --quarter 2026-Q3: eine Zeile pro Ticket aus outcome-what, gruppiert nach Milestone), nicht als Datei speichern - eine gespeicherte Zusammenfassung veraltet, sobald jemand ein Ticket wiederherstellt.
definition-of-done: "jaira logbook --quarter <JJJJ-Qn> zeigt pro Milestone die Tickets des Quartals mit einer Zeile Ergebnis; Test"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-22T00:14:26Z
updated-at: 2026-10-22T00:14:26Z
---

# Das Logbook verdichtet sich nach Quartalen: Ticket -> Logbook-Eintrag -> Quartalsuebersicht

## Definition of Done

- [ ] jaira logbook --quarter <JJJJ-Qn> zeigt pro Milestone die Tickets des Quartals mit einer Zeile Ergebnis; Test
- [ ] Von der Quartalsansicht kommt man zum Logbook-Ordner und zum Ticket (im TUI oder per Befehl); Test
- [ ] Entschieden und notiert, ob Zusammenfassungen gespeichert oder erzeugt werden
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

