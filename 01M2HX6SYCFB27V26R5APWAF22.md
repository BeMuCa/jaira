---
id: 01M2HX6SYCFB27V26R5APWAF22
title: "jaira update bringt das ganze Setup nach, nicht nur den Textblock"
status: pre-process
ready: true
creator: Alexander Sacharov
goal: "Nach einem Upgrade des Binaries genuegt ein Befehl, um ein Board nachzuziehen: Agent-Block, Rollen und Lanes - und er sagt vorher, was er anfassen wuerde."
definition-of-done: "'jaira update' bringt mit den entsprechenden Schaltern Agent-Block, Rollen und Lanes in einem Aufruf auf den Stand des installierten Binaries. Nachgestellt an einem Board, dessen Block, Rollen und Lanes alle drei veraltet sind."
tags:
  - cli
blocked-by: []
related:
  - 01M21PJTWPK4DZ89PRY6T8A8KM
commits: []
created-at: 2026-09-15T06:48:56Z
updated-at: 2026-10-15T12:34:51Z
assignee: "BeMuCa"
updated-by: BeMuCa
context: |-
  Nach 'jaira self update' haengt jedes Board zurueck, und es haengt an drei Stellen zurueck: der erzeugte Agent-Block in CLAUDE.md und AGENTS.md, die Rollen-Prompts, und die Lanes.

  Heute braucht das drei Befehle - 'jaira update', 'jaira roles install --global --force', und je Lane ein 'jaira lanes use <id> --force'. Der Sammelschalter fuer Lanes ist Ticket 0S1A7M; dieses hier ist die Frage, ob ein Befehl alle drei Dinge tut.

  'jaira update' traegt den passenden Namen und verspricht in seiner eigenen Hilfe genau das: 'brings a board's own setup up to date with whichever binary you already have installed'. Es regeneriert aber nur den Textblock. Rollen und Lanes sind ebenso 'a board's own setup'.

  Die Entscheidung, die dieses Ticket treffen muss und die groesser ist als die Umsetzung: was tut der nackte Aufruf ohne Schalter? Wer 'jaira update' heute gewohnheitsmaessig laeuft, erwartet eine Textblock-Regeneration. Faengt der Befehl an, Rollen und Lanes anzufassen, aendert sich das Verhalten unter den Fuessen aller, die ihn schon benutzen. Das ist eine Frage der Vertraeglichkeit, keine der Bequemlichkeit.

  Die zweite Haelfte des Problems, von Alex am 2026-09-15 benannt: auch mit einem Sammelbefehl muss er ihn in JEDEM Projekt einzeln aufrufen. 'jaira projects' kennt die Boards dieses Rechners bereits - acht Eintraege, darunter jaira, requirementsgenie und die Worktrees.

  Aber dieselbe Liste zeigt, warum ein rechnerweiter Lauf heute gefaehrlich waere: drei der acht Eintraege sind Wegwerf-Boards, die testing-Worker im Laufe des Tages unter /tmp/claude-*/scratchpad/ angelegt haben - sandbox, board74, board. Sie stehen dauerhaft in der Liste, obwohl ihre Verzeichnisse Muell sind. Ein Befehl, der ueber alle Boards laeuft, laeuft zuerst in diese hinein.

  Dazu kommt seit Commit 42d8253: die Lane-Dateien dieses Boards liegen unter git. Ein rechnerweiter Lauf macht damit Arbeitsbaeume in mehreren Repositories schmutzig, nicht nur Dateien im Heimverzeichnis.
claimed-by: EE-3NX6GL3-3760272
claimed-at: 2026-10-15T12:05:54Z
---

# jaira update bringt das ganze Setup nach, nicht nur den Textblock

## Definition of Done

- [ ] 'jaira update' bringt mit den entsprechenden Schaltern Agent-Block, Rollen und Lanes in einem Aufruf auf den Stand des installierten Binaries. Nachgestellt an einem Board, dessen Block, Rollen und Lanes alle drei veraltet sind.
- [ ] Der Aufruf ohne Schalter aendert sein Verhalten nicht still: was er heute tut, tut er weiterhin, und alles darueber hinaus wird ausdruecklich verlangt. Die Entscheidung steht begruendet im Ticket, nicht nur im Code.
- [ ] Der Befehl sagt vor dem Schreiben, was er anfassen wuerde, und was er auslaesst weil ein Mensch es geaendert hat - wie 'jaira roles install' es schon tut.
- [ ] Eine Zeile je Aenderung in core/release/notes/PWAF22.md, als Anweisung an den Leser formuliert und wahr
- [ ] jaira self upgrade fuehrt danach jaira update und die Rollen-Installation an den Orten aus, an denen Rollen installiert sind (global und/oder im Projekt), und sagt, was es angefasst hat; Test
- [ ] Sind installierte Skills (jaira, jaira-role-*) aelter als das laufende Binary, warnt jaira einmal und nennt den Befehl, der sie nachzieht; Test

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-13 20:51 · BeMuCa** — Berk 2026-10-13: nach einem Upgrade sollen die jaira-Skills immer neu installiert werden und jaira update automatisch laufen (DoD 5 und 6). Anlass: am 13.10. war ~/.claude/skills/jaira-role-pr noch vom 06.10. und kannte das Ablegen vor dem PR nicht - Alex musste im PR #74 nachfordern. Weder git pull noch jaira update ziehen die Rollen-Skills heute nach (jaira roles install --help). Gehoert mit T8A8KM zusammen.
- **2026-10-15 12:34 · BeMuCa** — Entscheidung DoD 2 - der nackte Aufruf: 'jaira update' ohne Schalter schreibt weiter NUR den Agent-Block und stempelt die Version, wie heute.
- Neue Schalter: --roles (Skills dort, wo die Rollen schon liegen: Projekt-Ordner dieses Boards + persoenlicher Ordner), --lanes (Lane-Dateien dieses Boards, die eine alte ausgelieferte Version sind), --dry-run.
- Kein --all: 'jaira update --roles --lanes' ist schon ein Aufruf.
- Kein --force an update: veraltete Dateien brauchen es nicht mehr. Eine selbst geaenderte Datei ersetzt man gezielt mit 'jaira roles install --force' oder 'jaira lanes use <id> --force'.
- Ohne Schalter kommt hoechstens eine Hinweiszeile dazu ('N Dateien sind aelter - jaira update --roles'). Sie schreibt nichts. Das JSON des nackten Aufrufs bleibt gleich.
- Grund: wer update gewohnheitsmaessig laeuft, bekommt nichts Neues geschrieben. Lane-Dateien liegen auf manchen Boards unter git und gehoeren dem Team.
- **2026-10-15 12:34 · BeMuCa** — Entscheidung DoD 3 - 'vorher sagen': 'jaira update --dry-run' zeigt, was passieren wuerde, und schreibt nichts.
- Gleiche Zeilen wie der echte Lauf, nur 'would update/write' statt 'updated/written'. 'skipped ... (edited here; ...)' wie bei 'jaira roles install'. Letzte Zeile: 'nothing was written (--dry-run)'.
- Gleicher Exit-Code wie der echte Lauf (3, wenn etwas uebersprungen wuerde). So macht es schon 'jaira move --dry-run' (internal/cli/flow.go:317) und 'sync-tasks --dry-run'.
- Verworfen: immer erst einen Plan drucken und dann schreiben - das 'vorher' waere nur die Reihenfolge der Zeilen. Verworfen: Rueckfrage j/n - haelt Agenten an.
- Bekannte Grenze: mit --lanes --dry-run zeigt die Block-Vorschau den Block fuer die Lanes von jetzt, nicht fuer die erneuerten.
