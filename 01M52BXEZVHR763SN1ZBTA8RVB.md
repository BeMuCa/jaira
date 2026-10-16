---
id: 01M52BXEZVHR763SN1ZBTA8RVB
title: README zeigt die Features mit Use cases und frischen Screenshots
status: in-progress
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
updated-at: 2026-10-16T12:47:44Z
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-613775
claimed-at: 2026-10-16T12:46:01Z
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
- **2026-10-16 12:47 · BeMuCa** — Entscheidungen vor der Arbeit (ohne Berk getroffen, Gruende): (1) 'die Liste' der 16 Features = features.md des frueheren Agents im Session-Scratchpad; README-Abschnitt hat genau diese 16. (2) add-world.sh (Milestones, fertige Tickets, Links) wird in scripts/demoworld.sh eingebaut, sonst kann niemand Roadmap- und Links-Bild nachbauen. (3) Version im Bild: shotgen mit -ldflags '-X github.com/BeMuCa/jaira/core/release.Current=0.5.1' und JAIRA_NO_UPDATE_CHECK=1 -> Kopfzeile 'jaira 0.5.1', kein Netz, kein Cache. main.version wirkt auf shotgen nicht, es ist ein eigenes main. (4) Hero oben bleibt launcher.png; board-main.png ersetzt den Board-Platz (frueher board.png, Zeile 23); board-links.png wandert mit dem 'Press l'-Satz unter das Feature 'Follow-ups link themselves'. (5) docs/img/board.png und home.png bleiben liegen, Berk entscheidet.
