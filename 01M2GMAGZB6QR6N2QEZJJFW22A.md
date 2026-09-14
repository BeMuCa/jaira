---
id: 01M2GMAGZB6QR6N2QEZJJFW22A
title: "Die drei Schleifen-Lanes liegen nur auf einem Rechner, obwohl sie die Fehler finden"
status: backlog
ready: false
creator: Alexander Sacharov
goal: "critique, optimize und testing kommen aus dem Binary, sodass ein frischer Klon sie hat, ohne dass jemand einen Katalog von Hand mitbringt."
context: |-
  core/lane/builtin liefert zehn Lanes aus: backlog, brainstorm, todo, pre-process, in-progress, human, review, signoff, done, blocked. Genau das meldet CI beim Aufsetzen eines Boards - '10 lane(s) from the built-in lanes'.

  critique, optimize und testing sind nicht dabei. Sie liegen in ~/.jaira/lanes auf Alex' Rechner und auf diesem Board, in keinem Repository. Wer jaira klont, bekommt ein Board ohne Schleifen, und niemand sagt ihm, dass es sie gibt.

  Das ist dieselbe Lage, in der die Rollen-Prompts vor C9V7ZV waren, und derselbe Beschluss sollte gelten: was die Arbeit treibt, gehoert ins Binary, nicht ins Heimverzeichnis.

  Dass die drei ihr Geld wert sind, ist an einem Tag belegt, nicht behauptet. Am 2026-09-14:
  - critique schickte APABM4 dreimal zurueck; Runde 1 und 2 waren echte Fehler (die Regel feuerte auf jedes '-o' und auf ein 'go install -o', das es nicht gibt).
  - critique fand an S1VM40, dass das Feld definition-of-done im Frontmatter noch die alte Zusage trug - und 'show --for-lane' reicht genau dieses Feld an jede spaetere Lane weiter, testing und review haetten gegen eine veraltete Bedingung geprueft.
  - testing fand an APABM4 einen Zweig, den kein Test abdeckte: eine Mutation, die ihn loeschte, blieb gruen.
  - optimize fand an S1VM40, dass cardHeight() die Zeilenzahl als Literal zurueckgab, waehrend cardSlots dieselbe Zahl als Konstante hielt - zwei Werte, die auseinanderlaufen koennen, und dann liest die Leiste ueber ihre Plaetze hinaus.

  Die Entscheidung, die dieses Ticket treffen muss und die nicht dieselbe ist wie sein Ziel: 'eingebaut' und 'auf einem neuen Board vorhanden' sind zwei verschiedene Dinge. 'jaira lanes add' holt eine eingebaute Lane auf ein Board; 'jaira lanes default' bestimmt, womit ein neues Board ueberhaupt startet. Drei Lanes auszuliefern heisst nicht, jedem neuen Board dreizehn Spalten und drei Schleifen zu geben - das waere eine Aussage darueber, wie jeder arbeiten soll, und ist an der Frage zu messen, die im Projekt ohnehin an jeder Funktion haengt: ist das kleiner als paca.

  Nicht Teil dieses Tickets: den Inhalt der drei Prompts zu aendern. Sie werden uebernommen wie sie sind.
definition-of-done: "critique, optimize und testing liegen in core/lane/builtin und sind im Binary eingebettet. Nachgestellt auf einem Rechner ohne ~/.jaira/lanes: 'jaira lanes add critique' bringt die Lane auf ein Board, ohne dass eine Datei von Hand mitgebracht wird."
tags:
  - cli
  - gates
blocked-by: []
related:
  - 01M2E248SM9X1JRZBNTHC9V7ZV
commits: []
created-at: 2026-09-14T18:54:26Z
updated-at: 2026-09-14T18:54:40Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
---

# Die drei Schleifen-Lanes liegen nur auf einem Rechner, obwohl sie die Fehler finden

## Definition of Done

- [ ] critique, optimize und testing liegen in core/lane/builtin und sind im Binary eingebettet. Nachgestellt auf einem Rechner ohne ~/.jaira/lanes: 'jaira lanes add critique' bringt die Lane auf ein Board, ohne dass eine Datei von Hand mitgebracht wird.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

