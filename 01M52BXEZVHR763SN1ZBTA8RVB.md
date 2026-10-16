---
id: 01M52BXEZVHR763SN1ZBTA8RVB
title: README zeigt die Features mit Use cases und frischen Screenshots
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Wer die README auf GitHub liest, sieht jede wichtige Funktion mit einem Use case in einfachem technischem Englisch und drei aktuelle Screenshots: Hauptboard, Launcher, Milestone-Roadmap"
context: |-
  Die README auf master zeigt alte Bilder (docs/img/board.png, home.png) und erwaehnt mehrere Funktionen gar nicht.
  Fehlend bis 16.10.: Milestone-Roadmap (Taste L) nur in Commands/Keys, Live-Board ohne Neuladen gar nicht, Overview, Explain-Pane, round:, Tag-Farben, follows auf Ticket-Branches.
  Berk am 16.10.: Roadmap auf jeden Fall rein, Use cases in simple technical English, Live-Board rein, drei NEUE Screenshots, Links-Fenster darf auf einem Bild zu sehen sein.
  Die Arbeit liegt schon im Worktree /home/berk/git/jAIra-readme-refresh, Branch docs/readme-refresh, uncommitted: README +103/-2, docs/img/board-links.png, launcher.png, milestone-roadmap.png.
  Bilder entstehen mit scripts/demoworld.sh + scripts/shotgen + termshot.py auf einer Demo-Welt, nie aus echten Tickets.
  Bekannte Fehler in den Bildern: Done-Spalte heisst 'one' (abgeschnitten), oben steht 'jaira dev', das Links-Fenster verdeckt die Mitte des Boards.
definition-of-done: "README hat einen Abschnitt mit allen 16 Features aus der Liste, jedes mit einer Zeile 'Use case:' in einfachem technischem Englisch"
tags:
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-10-16T12:45:13Z
updated-at: 2026-10-16T12:45:59Z
updated-by: BeMuCa
---

# README zeigt die Features mit Use cases und frischen Screenshots

## Definition of Done

- [ ] README hat einen Abschnitt mit allen 16 Features aus der Liste, jedes mit einer Zeile 'Use case:' in einfachem technischem Englisch
- [ ] Milestone-Roadmap und Live-Board ohne Neuladen stehen in der README, mit Quelle in core/release/notes belegt
- [ ] Drei neue Screenshots (Hauptboard, Launcher, Milestone-Roadmap mit Fortschritt): frisch gerendert, sha256 anders als die alten Bilder, kein abgeschnittener Spaltenname, keine 'dev'-Version
- [ ] Ein Screenshot zeigt das Links-Fenster, ohne dass das Hauptboard-Bild dadurch verdeckt ist
- [ ] git diff --check sauber, Code-Fences in README.md gerade, jedes Bild in der README verlinkt existiert

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

