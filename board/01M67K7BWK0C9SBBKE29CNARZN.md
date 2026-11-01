---
id: 01M67K7BWK0C9SBBKE29CNARZN
title: "Vor dem Anlegen einer Lane wird nach vorhandenen Lanes gesucht, damit keine Duplikate entstehen"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Wer eine Lane anlegen will - Mensch, Session oder Agent -, bekommt zuerst die vorhandenen Lanes mit aehnlichem Zweck gezeigt und nutzt sie wieder."
context: |-
  Wunsch von berkc am 31.10.2026 (Punkt 26): 'Wir sollten immer wenn wir lanes erstellen wollen mit der session auf existierende schauen damit wir keine duplikate erstellen - das soll nicht nur fuer unsere jetzige session gelten sondern auch fuer andere Leute.'
  Beispiel: der Board-Typ 'projektidee' (CEYMK1) schlug eigene Lanes 'abgelegt' und 'entscheidung' vor, wo done und HITL (human) schon existieren.
  Orte: jaira lanes add/template, Board-Typen, die Agent-Anleitung im CLAUDE.md-Block (core/board/announce.go), die TUI-Lane-Erstellung.
definition-of-done: jaira lanes template/add und das Anlegen im TUI listen aehnliche vorhandene Lanes (Katalog + Board + mitgeliefert) vor dem Anlegen
tags:
  - cli
blocked-by: []
related:
  - 01M67C1AJERNF8P9985F1Z6MRY
  - 01M67C1CS0YSRBA5PY8KCEYMK1
commits: []
created-at: 2026-10-30T23:44:51Z
updated-at: 2026-10-30T23:44:51Z
---

# Vor dem Anlegen einer Lane wird nach vorhandenen Lanes gesucht, damit keine Duplikate entstehen

## Definition of Done

- [ ] jaira lanes template/add und das Anlegen im TUI listen aehnliche vorhandene Lanes (Katalog + Board + mitgeliefert) vor dem Anlegen
- [ ] Der Agenten-Block in CLAUDE.md/AGENTS.md sagt: vor einer neuen Lane 'jaira lanes' lesen und vorhandene wiederverwenden
- [ ] Release-Note-Zeile

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

