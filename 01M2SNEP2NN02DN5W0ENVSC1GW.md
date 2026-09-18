---
id: 01M2SNEP2NN02DN5W0ENVSC1GW
title: "Die Pruefschleife gehoert ins Binary, nicht in den Katalog"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer jaira installiert, um Agenten arbeiten zu lassen, bekommt die Pruefschleife sofort - statt die Haelfte des Produkts im Katalog vermuten zu muessen"
context: |-
  Alex am 2026-09-18, nachdem 'jaira lanes' auf diesem Board dreizehn Lanes zeigte und zehn davon built-in waren.

  Was ist: 'critique', 'optimize' und 'testing' sind KEINE built-ins. Built-in sind genau zehn: backlog, brainstorm, todo, pre-process, in-progress, human, review, signoff, done, blocked (core/lane/lane.go:34, go:embed builtin/*.md). Die drei liegen im Katalog unter lanes/ und erreichen ein Board nur ueber 'jaira lanes market adopt <id>' plus 'jaira lanes add <id>'. lanes/README.md nennt das ausdruecklich Absicht: built-ins stehen auf jedem Board, Katalog-Lanes nur, wenn jemand sie absichtlich holt.

  Warum das heute stoert: ohne diese drei ist ein frisches Board ein Tracker, kein Agenten-Konveyer. Die Schleife critique -> in-progress ist das, was auf diesem Board die echten Defekte findet - an einem einzigen Tag unter anderem einen Unterscheider, der einen schreibenden Worker still zum stummen Leser gemacht haette, und einen fest verdrahteten Branch-Namen 'master' in einem ausgelieferten Prompt. Wer jaira wegen der Agenten installiert und die drei nicht kennt, bekommt nichts davon und erfaehrt auch nicht, dass es sie gibt: 'jaira lanes' zeigt nur, was installiert ist, und der Katalog meldet sich von selbst nirgends.

  Das Gegenargument steht in CLAUDE.md und ist ernst zu nehmen: 'Every feature is measured against is this smaller than paca - the project fails by growing'. Dreizehn Lanes als Voreinstellung zwingt jedem einen Konveyer auf, auch dem, der jaira als reinen Tracker aufsetzt.

  Die Entscheidung ist deshalb NICHT 'alles einbauen', sondern: was ist die Voreinstellung fuer wen. Moegliche Formen, in der Reihenfolge wachsenden Eingriffs - (1) nur sichtbar machen: ein frisches Board oder 'jaira lanes' nennt die Katalog-Lanes, die es nicht hat; (2) die drei ins Binary einbetten, aber nicht auf die Voreinstellung setzen, also ohne Netz adoptierbar; (3) ein zweites Default-Board 'agentic' neben dem schlanken, bei 'jaira init' waehlbar; (4) sie als Voreinstellung fuer jedes neue Board. Welche - das ist die eigentliche Frage dieses Tickets und gehoert in die brainstorm-Lane.
definition-of-done: "Das Ticket nennt EINE gewaehlte Form mit Begruendung, warum die anderen drei verworfen wurden - an der Scope-Regel aus CLAUDE.md gemessen, nicht am Bauchgefuehl."
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-09-18T07:07:21Z
updated-at: 2026-09-18T07:07:21Z
---

# Die Pruefschleife gehoert ins Binary, nicht in den Katalog

## Definition of Done

- [ ] Das Ticket nennt EINE gewaehlte Form mit Begruendung, warum die anderen drei verworfen wurden - an der Scope-Regel aus CLAUDE.md gemessen, nicht am Bauchgefuehl.
- [ ] Ein frisch mit 'jaira init' angelegtes Board erfaehrt von der Pruefschleife, ohne dass jemand den Katalog kennt. Nachgestellt an einer leeren Testdoska: der Weg von 'jaira init' zu einer arbeitenden critique-Lane ist ohne Vorwissen gehbar.
- [ ] Wer jaira ohne Netz benutzt, kommt an die drei Lanes heran, oder die Fehlermeldung sagt, dass sie aus dem Netz kommen und wie man sie sonst bekommt. 'jaira lanes market' holt heute von GitHub.
- [ ] Ein Board, das jaira als reinen Tracker benutzt, wird nicht mit einem Konveyer beladen, den niemand faehrt - die Wahl bleibt eine Wahl.
- [ ] Je eine Zeile in core/release/NOTES.md fuer das, was ein Benutzer dadurch anders tut.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

