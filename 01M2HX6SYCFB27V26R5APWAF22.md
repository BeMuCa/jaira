---
id: 01M2HX6SYCFB27V26R5APWAF22
title: "jaira update bringt das ganze Setup nach, nicht nur den Textblock"
status: backlog
ready: true
creator: Alexander Sacharov
goal: "Nach einem Upgrade des Binaries genuegt ein Befehl, um ein Board nachzuziehen: Agent-Block, Rollen und Lanes - und er sagt vorher, was er anfassen wuerde."
definition-of-done: "'jaira update' bringt mit den entsprechenden Schaltern Agent-Block, Rollen und Lanes in einem Aufruf auf den Stand des installierten Binaries. Nachgestellt an einem Board, dessen Block, Rollen und Lanes alle drei veraltet sind."
tags:
  - cli
blocked-by: []
related:
  - 01M2HX55H2EPFZ8M3AFW0S1A7M
commits: []
created-at: 2026-09-15T06:48:56Z
updated-at: 2026-09-15T06:49:41Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
context: |-
  Nach 'jaira self update' haengt jedes Board zurueck, und es haengt an drei Stellen zurueck: der erzeugte Agent-Block in CLAUDE.md und AGENTS.md, die Rollen-Prompts, und die Lanes.

  Heute braucht das drei Befehle - 'jaira update', 'jaira roles install --global --force', und je Lane ein 'jaira lanes use <id> --force'. Der Sammelschalter fuer Lanes ist Ticket 0S1A7M; dieses hier ist die Frage, ob ein Befehl alle drei Dinge tut.

  'jaira update' traegt den passenden Namen und verspricht in seiner eigenen Hilfe genau das: 'brings a board's own setup up to date with whichever binary you already have installed'. Es regeneriert aber nur den Textblock. Rollen und Lanes sind ebenso 'a board's own setup'.

  Die Entscheidung, die dieses Ticket treffen muss und die groesser ist als die Umsetzung: was tut der nackte Aufruf ohne Schalter? Wer 'jaira update' heute gewohnheitsmaessig laeuft, erwartet eine Textblock-Regeneration. Faengt der Befehl an, Rollen und Lanes anzufassen, aendert sich das Verhalten unter den Fuessen aller, die ihn schon benutzen. Das ist eine Frage der Vertraeglichkeit, keine der Bequemlichkeit.

  Die zweite Haelfte des Problems, von Alex am 2026-09-15 benannt: auch mit einem Sammelbefehl muss er ihn in JEDEM Projekt einzeln aufrufen. 'jaira projects' kennt die Boards dieses Rechners bereits - acht Eintraege, darunter jaira, requirementsgenie und die Worktrees.

  Aber dieselbe Liste zeigt, warum ein rechnerweiter Lauf heute gefaehrlich waere: drei der acht Eintraege sind Wegwerf-Boards, die testing-Worker im Laufe des Tages unter /tmp/claude-*/scratchpad/ angelegt haben - sandbox, board74, board. Sie stehen dauerhaft in der Liste, obwohl ihre Verzeichnisse Muell sind. Ein Befehl, der ueber alle Boards laeuft, laeuft zuerst in diese hinein.

  Dazu kommt seit Commit 42d8253: die Lane-Dateien dieses Boards liegen unter git. Ein rechnerweiter Lauf macht damit Arbeitsbaeume in mehreren Repositories schmutzig, nicht nur Dateien im Heimverzeichnis.
---

# jaira update bringt das ganze Setup nach, nicht nur den Textblock

## Definition of Done

- [ ] 'jaira update' bringt mit den entsprechenden Schaltern Agent-Block, Rollen und Lanes in einem Aufruf auf den Stand des installierten Binaries. Nachgestellt an einem Board, dessen Block, Rollen und Lanes alle drei veraltet sind.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

