---
id: 01M2T38EHGQ8F5TH3K7XBQGM6X
title: "Aus einer Menschen-Lane kann ein Mensch nur vorwaerts oder zur Seite, nicht zurueck"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Mensch, der in human oder signoff hinsieht und Nacharbeit will, schickt das Ticket selbst in eine frueher liegende Lane zurueck - ohne --force, ohne Folgeticket, und ohne dass derselbe Weg nach done aufgeht"
context: |-
  Aufgefallen Alex am 2026-09-18 an 7XXT00: Ticket stand in signoff, die review-Lane hatte drei Befunde geschrieben, einer davon ausdruecklich 'gehoert vor die Annahme'. Es gibt keinen Weg, es zur Nacharbeit zurueckzuschicken.

  Was die Doska heute kann. Auf dem signoff-Schirm stehen genau zwei Aktionen (internal/tui/signoff.go:151): 'a' nimmt an und geht direkt in m.lanes.Terminal(), also done (signoff.go:170), und 'f' legt ein Folgeticket an. Der move-Dialog 'm' hilft nicht: der Gate laesst einen Ausgang aus einer requires-human-exit-Lane nur bei req.Interactive durch (core/gate/gate.go:389), und dieses Flag wird ausschliesslich in signoff.go:179 und :189 gesetzt - das sind accept und followUp. Sonst nirgends. Ein Mensch am Gerät bekommt deshalb dieselbe Absage wie ein Agent: 'An agent cannot move a ticket out of it'.

  Warum es so ist, und warum die Begruendung nicht mehr traegt. Absichtlich, siehe fe5a979 und den Kommentar an signoff.go:239: 'Rejecting work without recording what is left undone is how the reason for a ticket gets lost'. Die Entscheidung unterscheidet 'angenommen' von 'braucht eigene Arbeit'. Sie kennt den haeufigsten Fall nicht: 'nicht fertig, die Arbeit geht auf DIESEM Ticket weiter'. Genau das koennen alle Agenten-Lanes - critique, optimize und testing schicken nach in-progress zurueck und wiederholen bis zur Stille. Die Schleife existiert fuer Agenten und fehlt dem Menschen, dessen Ablehnung die einzige ist, die niemand umgehen darf.

  Ein Folgeticket ist hier der falsche Zug, und es widerspricht der Regel dieses Boards, Befunde in das Ticket zu falten, das in der Hand liegt, statt Satellitentickets zu erzeugen. Ein unfertiges Ergebnis ist kein neues Thema.

  Heute bleibt nur 'jaira move <id> --to in-progress --force' aus einem Terminal ohne Agenten-Erkennung - also der Gate mit Gewalt an der Stelle, an der er am wenigsten umgangen werden sollte.
definition-of-done: "Ein Mensch kann ein Ticket aus human und signoff in jede Lane VOR der aktuellen schicken: der move-Dialog der Doska setzt Interactive, und der Gate an core/gate/gate.go:389 laesst diesen Ausgang durch. Ohne --force."
tags:
  - gates
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-09-18T11:08:36Z
updated-at: 2026-09-18T11:12:11Z
updated-by: Alexander Sacharov
---

# Aus einer Menschen-Lane kann ein Mensch nur vorwaerts oder zur Seite, nicht zurueck

## Definition of Done

- [ ] Der Gate an core/gate/gate.go:389 greift nur noch fuer einen Ausgang VORWAERTS: ein Ziel, das in der Lane-Reihenfolge nach der aktuellen liegt, bleibt abgelehnt. Ein Ziel VOR der aktuellen wird durchgelassen - ohne --force und ohne Interactive, also auch fuer einen Agenten. Die Richtung entscheidet, nicht der Actor.
- [ ] Vorwaerts bleibt allein die accept-Taste: ein move in die terminale Lane oder in irgendeine Lane nach der aktuellen wird weiter abgelehnt, aus der Doska wie aus der CLI. Mit Test fuer beide Richtungen.
- [ ] Ein Agent gewinnt dadurch nichts: Interactive wird weiterhin nur von einem Tastendruck in der Doska gesetzt, und ein Test pinnt, dass derselbe move aus der CLI abgelehnt bleibt.
- [ ] Der Grund fuer die Rueckgabe wird festgehalten, damit die Begruendung von fe5a979 gewahrt bleibt: die Rueckgabe schreibt eine Notiz mit der Ziel-Lane und dem, was der Mensch offen sieht, statt stillschweigend die Lane zu wechseln.
- [ ] Die Fusszeile des signoff-Schirms nennt die neue Aktion, und core/release/NOTES.md hat eine Zeile unter ## Unreleased.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

