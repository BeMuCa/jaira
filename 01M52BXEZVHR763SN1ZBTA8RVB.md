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
updated-at: 2026-10-16T12:57:17Z
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-613775
claimed-at: 2026-10-16T12:46:01Z
---

# README zeigt die Features mit Use cases und frischen Screenshots

## Definition of Done

- [x] README hat einen Abschnitt mit allen 16 Features aus der Liste, jedes mit einer Zeile 'Use case:' in einfachem technischem Englisch
  proof: README.md:34 '## Features': 16 Eintraege (README.md:36-126), 16 Zeilen ' - Use case:' (awk-Zaehlung des Abschnitts); Liste = features.md des frueheren Agents, gleiche 16 Themen
- [x] Milestone-Roadmap und Live-Board ohne Neuladen stehen in der README, mit Quelle in core/release/notes belegt
  proof: README.md:36 Milestone Roadmap <- core/release/notes/0.5.0/D0SAHM.md (L, Bubbles, enter, b, --after); README.md:113 Live board <- core/release/notes/0.5.0/NA413Y.md (changed ticket within 2 seconds)
- [x] Drei neue Screenshots (Hauptboard, Launcher, Milestone-Roadmap mit Fortschritt): frisch gerendert, sha256 anders als die alten Bilder, kein abgeschnittener Spaltenname, keine 'dev'-Version
  proof: docs/img/board-main.png (board 170 26 z right right), launcher.png (home 110 24), milestone-roadmap.png (roadmap 100 24) aus scripts/demoworld.sh frisch; sha256 bf3fffc9 != board.png 1dd568e7, 0fdba535 != home.png d2c803ef, 629eb98a != vorheriges 9a96dedf; 'jaira dev' 0 Treffer, Kopf 'jaira 0.5.1'; alle 10 Lane-Titel ganz (Bild angesehen)
- [x] Ein Screenshot zeigt das Links-Fenster, ohne dass das Hauptboard-Bild dadurch verdeckt ist
  proof: README.md:23 board-main.png ohne Overlay; README.md:107 board-links.png als eigenes Bild beim Links-Satz unter 'Follow-ups link themselves'
- [x] git diff --check sauber, Code-Fences in README.md gerade, jedes Bild in der README verlinkt existiert
  proof: git diff --check exit 0; grep -c '^```' README.md = 56 (gerade); 10 von 10 img-src in README.md existieren (Skript-Schleife)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] shotgen: roadmap-View und Tasten nach <rows> aus dem Scratch-Overlay in scripts/shotgen/main.go uebernehmen
- [x] termshot.py: SGR 30-37/90-97/40-47/100-107/39/49 aus dem Overlay in scripts/termshot.py uebernehmen
- [x] demoworld.sh: Milestones, fertige Tickets und Links aus add-world.sh einbauen; Render-Rezept mit Release-Version und allen vier Bildern
- [x] Demo-Welt frisch bauen, vier Bilder rendern (board-main, board-links, launcher, milestone-roadmap) mit jaira 0.5.1 und ohne abgeschnittenen Lane-Titel
- [~] README: board-main.png auf den Board-Platz, board-links.png zum Links-Satz unter Follow-ups
- [ ] Pruefen: diff --check, Fences, Bildpfade, file, sha256, Bilder ansehen; gofmt/vet/build; commit

## Progress
- **2026-10-16 12:47 · BeMuCa** — Entscheidungen vor der Arbeit (ohne Berk getroffen, Gruende): (1) 'die Liste' der 16 Features = features.md des frueheren Agents im Session-Scratchpad; README-Abschnitt hat genau diese 16. (2) add-world.sh (Milestones, fertige Tickets, Links) wird in scripts/demoworld.sh eingebaut, sonst kann niemand Roadmap- und Links-Bild nachbauen. (3) Version im Bild: shotgen mit -ldflags '-X github.com/BeMuCa/jaira/core/release.Current=0.5.1' und JAIRA_NO_UPDATE_CHECK=1 -> Kopfzeile 'jaira 0.5.1', kein Netz, kein Cache. main.version wirkt auf shotgen nicht, es ist ein eigenes main. (4) Hero oben bleibt launcher.png; board-main.png ersetzt den Board-Platz (frueher board.png, Zeile 23); board-links.png wandert mit dem 'Press l'-Satz unter das Feature 'Follow-ups link themselves'. (5) docs/img/board.png und home.png bleiben liegen, Berk entscheidet.
- **2026-10-16 12:57 · BeMuCa** — in-progress: Funde neben dem Ticket, NICHT behoben (Berk entscheidet ob Ticket): (a) TUI zeichnet die Link-Linien der Cursor-Karte quer ueber die senkrechten Titel schmaler Lanes (z) - deshalb steht der Cursor im board-main auf einer Karte ohne Links ('right right' = Todo-Karte). (b) Im z-Modus bleiben bei 170 Spalten ~20 Spalten rechts leer; Blocked erscheint erst ab 168 Spalten, darunter scrollt das Board seitlich - das war Done='one'. (c) --blocked-by mit Handle: 'jaira links <blocker>' und das Links-Fenster listen den wartenden Ticket NICHT unter blocking (getestet: MEN5J3 --blocked-by KSWEQM), nur die volle ID wirkt; linklines.go loest Handles dagegen auf. (d) docs/img/pipeline.png (alt, nicht Teil des Tickets) zeigt 'Implementi…' und 'Human Revi…' abgeschnitten. (e) Bei 170x24 enden die Lane-Kaesten ungleich hoch; 26 Zeilen sind sauber. Sackgasse: python3 findet PIL nicht, wenn HOME in den Scratchpad zeigt - PYTHONUSERBASE=/home/berk/.local setzen.
