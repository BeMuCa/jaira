---
id: 01M52BXEZVHR763SN1ZBTA8RVB
title: README zeigt die Features mit Use cases und frischen Screenshots
status: critique
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
updated-at: 2026-10-16T12:59:15Z
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-613775
claimed-at: 2026-10-16T12:46:01Z
outcome-what: "README: Features-Abschnitt mit 16 Funktionen und Use cases; vier frische Bilder (board-main, board-links, launcher, milestone-roadmap) mit jaira 0.5.1; shotgen-Tasten/roadmap, termshot-Grundfarben und die Demo-Welt-Extras ins Repo"
outcome-why: "Die README zeigte alte Bilder und verschwieg Roadmap, Live-Board und weitere Funktionen; die Bilder vom Vortag hatten 'jaira dev', einen abgeschnittenen Done-Titel und ein Links-Fenster mitten im Hauptbild"
outcome-resolves: "Hauptboard ohne Overlay mit allen 10 Lane-Titeln ganz, Links-Fenster als eigenes Bild beim Links-Text, keine 'dev'-Version; scripts/demoworld.sh baut Welt und shotgen mit Release-Version, sodass jeder die Bilder nachbauen kann"
review-summary: |-
  scripts/demoworld.sh:68-69 the comment says the corner reads "jaira 0.5.1" — false from the next release on; say "jaira <newest release>" instead
  scripts/shotgen/main.go:83 a key word that is neither named nor one character presses its first letter ("rigth" -> r) and "" panics on [rune]0; die on any word that is not named and not exactly one rune
  README.md: the two added paragraphs after the logbook section (round:/review-round:) and after the Commands block (follows/tag colours) are single 200-300 character lines; the rest of the README wraps prose at about 80 columns — wrap them the same way
round: 2
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
- [x] README: board-main.png auf den Board-Platz, board-links.png zum Links-Satz unter Follow-ups
- [x] Pruefen: diff --check, Fences, Bildpfade, file, sha256, Bilder ansehen; gofmt/vet/build; commit

## Progress
- **2026-10-16 12:47 · BeMuCa** — Entscheidungen vor der Arbeit (ohne Berk getroffen, Gruende): (1) 'die Liste' der 16 Features = features.md des frueheren Agents im Session-Scratchpad; README-Abschnitt hat genau diese 16. (2) add-world.sh (Milestones, fertige Tickets, Links) wird in scripts/demoworld.sh eingebaut, sonst kann niemand Roadmap- und Links-Bild nachbauen. (3) Version im Bild: shotgen mit -ldflags '-X github.com/BeMuCa/jaira/core/release.Current=0.5.1' und JAIRA_NO_UPDATE_CHECK=1 -> Kopfzeile 'jaira 0.5.1', kein Netz, kein Cache. main.version wirkt auf shotgen nicht, es ist ein eigenes main. (4) Hero oben bleibt launcher.png; board-main.png ersetzt den Board-Platz (frueher board.png, Zeile 23); board-links.png wandert mit dem 'Press l'-Satz unter das Feature 'Follow-ups link themselves'. (5) docs/img/board.png und home.png bleiben liegen, Berk entscheidet.
- **2026-10-16 12:57 · BeMuCa** — in-progress: Funde neben dem Ticket, NICHT behoben (Berk entscheidet ob Ticket): (a) TUI zeichnet die Link-Linien der Cursor-Karte quer ueber die senkrechten Titel schmaler Lanes (z) - deshalb steht der Cursor im board-main auf einer Karte ohne Links ('right right' = Todo-Karte). (b) Im z-Modus bleiben bei 170 Spalten ~20 Spalten rechts leer; Blocked erscheint erst ab 168 Spalten, darunter scrollt das Board seitlich - das war Done='one'. (c) --blocked-by mit Handle: 'jaira links <blocker>' und das Links-Fenster listen den wartenden Ticket NICHT unter blocking (getestet: MEN5J3 --blocked-by KSWEQM), nur die volle ID wirkt; linklines.go loest Handles dagegen auf. (d) docs/img/pipeline.png (alt, nicht Teil des Tickets) zeigt 'Implementi…' und 'Human Revi…' abgeschnitten. (e) Bei 170x24 enden die Lane-Kaesten ungleich hoch; 26 Zeilen sind sauber. Sackgasse: python3 findet PIL nicht, wenn HOME in den Scratchpad zeigt - PYTHONUSERBASE=/home/berk/.local setzen.
- **2026-10-16 12:58 · BeMuCa** — critique pass 1: (1) demoworld.sh Kommentar nennt 0.5.1 fest -> 'jaira <newest release>'; (2) shotgen: unbekanntes Mehrzeichen-Wort drueckt still seinen ersten Buchstaben, '' panict -> die, wenn nicht benannt und nicht genau ein Zeichen; (3) README: zwei neue Absaetze (round:, follows/tag colours) sind Einzeilen mit 200-300 Zeichen -> auf ~80 umbrechen wie der Rest. Stehen gelassen: named-Map mit left/up/down/enter/esc und termshot 40-47/100-107/39/49 obwohl die Rezepte sie nicht brauchen - der Auftrag verlangte die Patches so ins Repo; board-links-Bild unter 'Follow-ups link themselves' statt eigenem 17. Feature, weil der Abschnitt genau die 16 der Liste haben soll.
- **2026-10-16 12:59 · BeMuCa** — in-progress Runde 2: alle drei Critique-Funde behoben - demoworld.sh Kommentar ohne feste Version; shotgen stirbt bei unbekanntem Tastenwort (getestet: 'rigth' und '' -> exit 1); zwei README-Absaetze auf 80 umbrochen. board-main/board-links mit neuem shotgen neu gerendert: byte-gleich (cmp), Proofs bleiben gueltig.
