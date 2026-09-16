---
id: 01M2E76PHSYAVX5T5AW23AASF5
title: "Die Tafel kann nicht sagen, was als naechstes vorgeschlagen werden soll"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Teamlead kann jaira fragen, welches Ticket vorzuschlagen ist, wenn Haende frei werden, und bekommt eine Antwort, die nicht bloss 'das aelteste' lautet"
context: |-
  Ein Prioritaetsfeld gibt es nicht. Was es gibt, zielt auf anderes:
  - blocked-by ist eine harte Abhaengigkeit; das Ticket faellt ganz aus actionable heraus.
  - follows sagt 'dieses nach jenem'.
  - sortByProgress in internal/cli/flow.go:511 ordnet nach drei Schluesseln in dieser Reihenfolge: Lane-Precedence, dann wer die Ausgabe der Lane noch schuldet vor dem, der sie schon geliefert hat, dann das aeltere Ticket zuerst. ULIDs sortieren nach Zeit, der letzte Schluessel ist also FIFO.
  - jaira next gibt das erste aus dieser Ordnung.
  Die Ordnung ist gut und absichtlich: Angefangenes fertigmachen, bevor Neues beginnt.
  Die Luecke liegt dahinter. Unter frischen, gleich weit gekommenen Tickets lautet die Antwort heute 'wer zuerst angelegt wurde'. Auf die Frage eines Teamleads - was schlage ich vor, wenn diese hier durch sind - ist das keine Antwort, sondern das Fehlen einer.
  Offene Frage fuer die Brainstorm-Lane, absichtlich nicht entschieden:
  Eine Skala (hoch/mittel/niedrig) ist der naheliegende Weg und der, der in jedem Tracker verrottet: nach zwei Monaten ist alles hoch und die Sortierung wieder sinnlos. Eine Skala ohne Knappheit ist keine Skala.
  Der Gegenvorschlag, der hier zu pruefen ist: kein Grad, sondern ein knappes Kennzeichen - eine Handvoll Tickets traegt 'vorschlagen, wenn Haende frei werden'. Binaer statt abgestuft. Eine Liste von fuenf kann nicht entwerten, wie 'hoch' entwertet.
  Was die Lane entscheiden muss: ob Knappheit erzwungen wird (eine Obergrenze, die jaira durchsetzt) oder nur sichtbar gemacht (zwanzig Kennzeichen sind sofort als kaputt erkennbar). Erzwingen ist ehrlicher und aergert; Sichtbarmachen ist billiger und wird ignoriert.
  Zweitens: ob das Kennzeichen im Ticket steht oder in einer eigenen Datei. Im Ticket heisst, jede Umpriorisierung ist ein Ticket-Commit; in einer Datei heisst, die Reihenfolge ist eine Sache und nicht zwanzig.
  Nicht Teil dieses Tickets: die bestehende Ordnung in sortByProgress umbauen. Angefangenes zuerst bleibt.
definition-of-done: "Ein Teamlead kann per CLI erfragen, welches Ticket als naechstes vorzuschlagen ist, und die Antwort beruht auf mehr als dem Alter; jaira next und jaira list beruecksichtigen das Ergebnis; die gewaehlte Mechanik macht sichtbar oder verhindert, dass sie durch Ueberbenutzung wertlos wird; Angefangenes wird weiterhin vor Neuem fertiggemacht; eine Zeile in core/release/NOTES.md unter ## Unreleased; go test ./... -race gruen"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-09-13T20:26:40Z
updated-at: 2026-09-13T20:26:40Z
---

# Die Tafel kann nicht sagen, was als naechstes vorgeschlagen werden soll

## Definition of Done

- [ ] Ein Teamlead kann per CLI erfragen, welches Ticket als naechstes vorzuschlagen ist, und die Antwort beruht auf mehr als dem Alter; jaira next und jaira list beruecksichtigen das Ergebnis; die gewaehlte Mechanik macht sichtbar oder verhindert, dass sie durch Ueberbenutzung wertlos wird; Angefangenes wird weiterhin vor Neuem fertiggemacht; eine Zeile in core/release/NOTES.md unter ## Unreleased; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

