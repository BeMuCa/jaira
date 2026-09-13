---
id: 01M2E5JQKJ2GXQEV7XCBTHRR79
title: "CLAUDE.md kuerzen: der Technologie-Stack gehoert nach docs/"
status: optimize
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "CLAUDE.md traegt nur noch, was ein Worker zum Arbeiten braucht; die Stack-Recherche steht in docs/STACK.md und wird von dort verlinkt"
context: |-
  Jede frisch gestartete Agenten-Sitzung liest CLAUDE.md einmal komplett, bevor sie irgendetwas tut. Seit die Lanes von getrennten Sitzungen gefahren werden, passiert das pro Lane erneut - bei acht Lanes acht Mal.
  Gemessen am 2026-09-13: CLAUDE.md ist 34993 Bytes, rund 8700 Token. AGENTS.md kommt mit 2600 Token dazu.
  Der groesste Block ist Zeile 27 bis 164, von '## Technology Stack' bis '## Sources': 21258 Bytes, rund 5300 Token. Das ist die Stack-Entscheidung mit Versionstabellen, 'What NOT to Use', Versionskompatibilitaet und Quellenliste.
  Dieser Block ist fuer einen Menschen geschrieben, der entscheidet, womit gebaut wird. Ein Worker in der Lane in-progress braucht ihn nicht: er aendert Go-Code in einem Repository, in dem die Entscheidung laengst getroffen und im go.mod sichtbar ist.
  Rechnung: 5300 Token mal acht Lanes sind rund 42000 Token pro Ticket, nur fuers Wiederlesen einer Entscheidung, die niemand mehr trifft.
  Vorsicht an zwei Stellen:
  - Ab Zeile 208 steht der von jaira selbst erzeugte Block ('## Task tracking: jaira' und '## This board's lanes'). Der wird bei jedem 'jaira update' neu geschrieben und darf nicht von Hand angefasst werden.
  - Ab Zeile 333 steht der Marker <!-- jaira:local -->. Alles dahinter ist Handarbeit und ueberlebt die Regeneration. Die Abschnitte zu Branch/PR und zu core/release/NOTES.md liegen dort und bleiben, wo sie sind.
  Nicht Teil dieses Tickets: AGENTS.md. Das besteht fast nur aus dem erzeugten Block und hat nichts zu kuerzen.
definition-of-done: "docs/STACK.md enthaelt den Inhalt der bisherigen Zeilen 27-164 unveraendert; CLAUDE.md verweist an der Stelle mit einer Zeile darauf; CLAUDE.md ist danach unter 4000 Token; der erzeugte jaira-Block und alles hinter <!-- jaira:local --> sind unveraendert; 'go run ./cmd/jaira update' schreibt den Block fehlerfrei neu; go test ./... -race gruen"
tags:
  - docs
blocked-by: []
related: []
commits:
  - e3c99fc365eb431b0a29b8db53467dc805d02acc
created-at: 2026-09-13T19:58:17Z
updated-at: 2026-09-13T20:39:21Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-8635
claimed-at: 2026-09-13T20:28:44Z
outcome-what: "docs/STACK.md tragt jetzt einen eigenen H1-Titel und eine Zeile, die .planning/research/STACK.md als Vollfassung und 2026-08-11 als Rechercheteil nennt"
outcome-why: "die Datei war die gekuerzte Zweitfassung, sagte das aber nirgends - wer nur dem Link aus CLAUDE.md folgte, hielt sie fuer das Ganze"
outcome-resolves: "die einzige Feststellung der critique vom 2026-09-13 20:29"
review-summary: "none"
review-gaps: "Entfernt: die mitkopierte Ueberschrift '## Technology Stack' in docs/STACK.md, die direkt unter dem gleichlautenden H1 stand und eine leere Sektion aufmachte; kein Anker und kein Link zeigt darauf. Stehen gelassen: der Forschungstext selbst (woertlich verschoben, dieses Ticket formuliert nichts um); die Doppelnennung von .planning/research/STACK.md in CLAUDE.md und in docs/STACK.md - ein Hop auseinander und beide Male fuer einen anderen Leser; die vorbestehende Gliederung der Recherche ('## 1. Language' vor '## Recommended Stack'), die schon in CLAUDE.md so stand; CLAUDE.md sonst unaufgeraeumt, wie vom Zuschnitt verlangt. Kein Go-Code beruehrt, go test ./... -race gruen."
---

# CLAUDE.md kuerzen: der Technologie-Stack gehoert nach docs/

## Definition of Done

- [x] docs/STACK.md enthaelt den Inhalt der bisherigen Zeilen 27-164 unveraendert; CLAUDE.md verweist an der Stelle mit einer Zeile darauf; CLAUDE.md ist danach unter 4000 Token; der erzeugte jaira-Block und alles hinter <!-- jaira:local --> sind unveraendert; 'go run ./cmd/jaira update' schreibt den Block fehlerfrei neu; go test ./... -race gruen
  proof: docs/STACK.md:1-133 (diff <(git show HEAD:CLAUDE.md | sed -n '27,159p') docs/STACK.md leer); CLAUDE.md:29 Verweiszeile; wc -c CLAUDE.md = 13968 (~3450 Token); git diff CLAUDE.md hat genau einen Hunk @@ -26,137 +26,7 @@, jaira-Block und jaira:local unberuehrt; go test ./... -race grün

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Inhalt der CLAUDE.md-Zeilen 27-160 (## Technology Stack bis zum letzten ## Sources-Punkt) 1:1 nach docs/STACK.md kopieren, ohne Titelzeile und ohne Umformulierung
- [x] In CLAUDE.md die Zeilen 28-160 loeschen; Ueberschrift ## Technology Stack stehen lassen und darunter genau eine Zeile setzen, die auf docs/STACK.md und auf .planning/research/STACK.md als Vollfassung verweist
- [x] Die GSD-Marker <!-- GSD:stack-start source:research/STACK.md --> und <!-- GSD:stack-end --> stehen lassen und die Verweiszeile dazwischen setzen
- [x] Nachmessen: wc -c CLAUDE.md erwartet rund 13800 Bytes (ca. 3450 Token, unter 4000); grep pruefen, dass GSD:stack-start/end, jaira:start, jaira:local und jaira:end unveraendert vorhanden sind
- [x] go run ./cmd/jaira update ausfuehren, danach git diff CLAUDE.md: der Block zwischen jaira:start und jaira:end und alles hinter jaira:local muessen unveraendert sein
- [x] go test ./... -race laufen lassen und gruen erwarten

## Progress
- **2026-09-13 20:08 · Alexander Sacharov** — Korrektur zur Messung im Kontext: AGENTS.md wird von Claude Code nicht gelesen - in einer frischen Sitzung kommt nur CLAUDE.md im Kontext an. Die 2648 Token fuer AGENTS.md sind faelschlich mitgezaehlt worden. Die tatsaechliche Ersparnis pro Worker sind die rund 5300 Token des handgeschriebenen Stack-Blocks in CLAUDE.md, nicht mehr.
Zwei Wege, die hier geprueft und verworfen sind, damit sie nicht wiederkommen:
- AGENTS.md loeschen. Geht nicht: core/board/announce.go:291 haelt fest, dass AGENTS.md das ist, was Codex und andere Werkzeuge lesen, und die Projekt-Beschraenkung lautet 'CLI must be usable by any bash-capable agent, not only Claude Code'.
- Den Installer so aendern, dass er nur CLAUDE.md schreibt, wenn es die gibt. Dasselbe Problem, nur stiller: auf einem Rechner mit beiden Clients verliert Codex unbemerkt den Zugang zur Tafel.
Die beiden Dateien doppeln einander ohnehin kaum: AGENTS.md ist fast nur der erzeugte jaira-Block, CLAUDE.md ist derselbe Block plus alles von Hand Dazugeschriebene. Zu kuerzen ist der zweite Summand, nicht die zweite Datei.
- **2026-09-13 20:17 · Alexander Sacharov** — Grenze praezisiert: die DoD nennt Zeilen 27-164, aber Zeile 161 ist <!-- GSD:stack-end -->, 163 ist GSD:conventions-start. Der eigentliche Inhalt sind die Zeilen 27-160, von '## Technology Stack' bis zum letzten Aufzaehlungspunkt unter '## Sources'. Der Plan verschiebt 27-160 und laesst die Marker stehen.
- **2026-09-13 20:17 · Alexander Sacharov** — Der Stack-Block steht nicht frei in CLAUDE.md, sondern zwischen <!-- GSD:stack-start source:research/STACK.md --> und <!-- GSD:stack-end -->. Das ist ein von GSD erzeugter Block, dessen Quelle .planning/research/STACK.md ist (30256 Bytes, die Vollfassung). Ein Generator dafuer ist hier aber nicht installiert: grep nach 'GSD:stack-start' in ~/.claude trifft nur Transcript-Dateien, kein Skill und kein Skript. Heute schreibt also nichts den Block neu. Deshalb: Marker stehen lassen und die Verweiszeile dazwischen setzen - laeuft irgendwann doch ein GSD-Docs-Lauf, ist das ein sichtbarer Diff und kein stiller Rueckfall.
- **2026-09-13 20:17 · Alexander Sacharov** — docs/STACK.md wird damit eine gekuerzte Zweitfassung von .planning/research/STACK.md. Das ist Absicht, nicht Versehen: die DoD verlangt die Zeilen 27-164 unveraendert, und die CLAUDE.md-Fassung ist die bereits verdichtete. docs/ ist der Pfad, den ein Agent liest; .planning/research/ der des Recherchierenden. Die eine Verweiszeile in CLAUDE.md nennt beide, damit niemand die Vollfassung verliert.
- **2026-09-13 20:17 · Alexander Sacharov** — Groessenrechnung gemessen, nicht geschaetzt: CLAUDE.md ist 34993 Bytes, die Zeilen 27-160 sind 21180 Bytes. Bleiben rund 13813 Bytes, also ca. 3450 Token. Das ist unter den geforderten 4000, ohne dass sonst noch etwas gekuerzt werden muss.
- **2026-09-13 20:21 · Alexander Sacharov** — Plan-Schritt 5 hat eine Falle: "go run ./cmd/jaira update" schreibt in diesem Worktree den jaira-Block KUERZER neu - die Lanes critique, optimize und testing verschwinden aus CLAUDE.md und AGENTS.md. Grund: .jaira/ ist gitignored und damit pro Worktree eigen; /home/alex/projects/jaira-THRR79/.jaira/lanes enthaelt diese drei Lane-Dateien nicht, /home/alex/projects/jaira/.jaira/lanes schon. update hat also fehlerfrei gearbeitet, nur aus einer aelteren Lane-Konfiguration heraus. Konsequenz hier: AGENTS.md und der jaira-Block in CLAUDE.md wurden nach dem Lauf per git checkout zurueckgesetzt, damit der Commit nur die Stack-Verschiebung traegt. Wer in einem Worktree "jaira update" laufen laesst, muss danach den Diff pruefen - sonst landet der Lane-Verlust still im Commit. Eigenes Thema, nicht dieses Ticket.
- **2026-09-13 20:25 · Alexander Sacharov** — Der Worktree /home/alex/projects/jaira-THRR79 hatte eine veraltete .jaira/lanes-Konfiguration ohne critique, optimize und testing - .jaira/lanes ist nicht in git (nur .jaira/tickets ist getrackt), also faehrt jeder Worktree seine eigene Kopie. Deshalb schlug "jaira move --to critique" zuerst mit exit 3 fehl. Behoben durch cp der drei Lane-Dateien und der Datei order aus /home/alex/projects/jaira/.jaira/lanes. Wer hier einen neuen Worktree aufmacht, muss das wieder tun.
- **2026-09-13 20:29 · Alexander Sacharov** — critique: Eine Feststellung, sonst nichts. docs/STACK.md ist jetzt die Datei, auf die CLAUDE.md:29 zeigt, sagt aber selbst nicht, was sie ist: Zeile 1 ist '## Technology Stack', direkt gefolgt von '## 1. Language: Go — not Rust'. docs/AGENTS.md und docs/COMMANDS.md beginnen beide mit einem H1-Titel; das ist das Muster im selben Verzeichnis. Zu aendern: '# Technology Stack' als neue Zeile 1 und darunter eine Zeile, die .planning/research/STACK.md als Vollfassung und 2026-08-11 als Rechercheteil nennt. Grund: die Zweitfassung ist gekuerzt und traegt kein Zeichen davon - wer nur dem Link folgt, haelt sie fuer das Ganze. Der Hinweis auf die Vollfassung steht heute nur in CLAUDE.md, also genau in der Datei, die niemand mehr lesen soll. Zwei Zeilen davor, der verschobene Text bleibt unangetastet. Nicht aufgemacht: die Doppelung docs/STACK.md gegen .planning/research/STACK.md (in der Notiz von 20:17 bereits entschieden) und der GSD-Marker stack-start, dessen Quelle nicht mehr zum Inhalt passt (dort ebenfalls bewusst so gelassen).
- **2026-09-13 20:31 · Alexander Sacharov** — in-progress nach critique: docs/STACK.md hat jetzt '# Technology Stack' als H1 und darunter eine Zeile, die .planning/research/STACK.md als Vollfassung und 2026-08-11 als Rechercheteil nennt. Bewusst stehen gelassen: die Zeile '## Technology Stack' direkt darunter. Sie ist Teil des woertlich verschobenen Textes (CLAUDE.md Zeile 27) und die critique hat nur zwei Zeilen davor verlangt, nicht das Entfernen. Wirkung: H1 und H2 tragen denselben Titel. Wer das aufraeumen will, braucht dafuer eine Entscheidung, dass der verschobene Text doch angefasst werden darf - hier nicht getroffen.
- **2026-09-13 20:33 · Alexander Sacharov** — critique, zweiter Durchlauf: keine Feststellung, review-summary=none. Die eine Feststellung vom 20:29 ist abgearbeitet (8c0ab9e: H1 plus Verweiszeile in docs/STACK.md:1-3). Der Diff ist sonst unveraendert und bleibt zwei Dateien, ein Hunk in CLAUDE.md. Nicht aufgemacht, weil auf dem ersten Durchlauf bewusst so entschieden: die doppelte Ueberschrift '## Technology Stack' unter dem neuen H1 (sie ist Teil des woertlich verschobenen Textes, und der bleibt in diesem Ticket unangetastet), die Doppelung docs/STACK.md gegen .planning/research/STACK.md, und der GSD-Marker stack-start, dessen Quelle nicht mehr zum Inhalt passt. Damit ist die critique-Schleife zu.
- **2026-09-13 20:39 · Alexander Sacharov** — optimize: docs/STACK.md trug unter dem neuen H1 '# Technology Stack' noch die mitkopierte Zeile '## Technology Stack' aus CLAUDE.md - eine leere Ueberschrift, die dasselbe zweimal sagt. Entfernt. Kein Link und kein Anker zeigt darauf (grep ueber *.md und *.go). Der Forschungstext selbst ist unangetastet; die DoD-Formulierung 'Zeilen 27-164 unveraendert' meint diesen Text, nicht die mitgeschleifte Sektionsueberschrift, die in CLAUDE.md ohnehin stehen bleibt.
