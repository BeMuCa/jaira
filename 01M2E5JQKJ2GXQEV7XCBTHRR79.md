---
id: 01M2E5JQKJ2GXQEV7XCBTHRR79
title: "CLAUDE.md kuerzen: der Technologie-Stack gehoert nach docs/"
status: in-progress
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
commits: []
created-at: 2026-09-13T19:58:17Z
updated-at: 2026-09-13T20:20:36Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-60480
claimed-at: 2026-09-13T20:13:55Z
---

# CLAUDE.md kuerzen: der Technologie-Stack gehoert nach docs/

## Definition of Done

- [ ] docs/STACK.md enthaelt den Inhalt der bisherigen Zeilen 27-164 unveraendert; CLAUDE.md verweist an der Stelle mit einer Zeile darauf; CLAUDE.md ist danach unter 4000 Token; der erzeugte jaira-Block und alles hinter <!-- jaira:local --> sind unveraendert; 'go run ./cmd/jaira update' schreibt den Block fehlerfrei neu; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Inhalt der CLAUDE.md-Zeilen 27-160 (## Technology Stack bis zum letzten ## Sources-Punkt) 1:1 nach docs/STACK.md kopieren, ohne Titelzeile und ohne Umformulierung
- [x] In CLAUDE.md die Zeilen 28-160 loeschen; Ueberschrift ## Technology Stack stehen lassen und darunter genau eine Zeile setzen, die auf docs/STACK.md und auf .planning/research/STACK.md als Vollfassung verweist
- [x] Die GSD-Marker <!-- GSD:stack-start source:research/STACK.md --> und <!-- GSD:stack-end --> stehen lassen und die Verweiszeile dazwischen setzen
- [~] Nachmessen: wc -c CLAUDE.md erwartet rund 13800 Bytes (ca. 3450 Token, unter 4000); grep pruefen, dass GSD:stack-start/end, jaira:start, jaira:local und jaira:end unveraendert vorhanden sind
- [ ] go run ./cmd/jaira update ausfuehren, danach git diff CLAUDE.md: der Block zwischen jaira:start und jaira:end und alles hinter jaira:local muessen unveraendert sein
- [ ] go test ./... -race laufen lassen und gruen erwarten

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
