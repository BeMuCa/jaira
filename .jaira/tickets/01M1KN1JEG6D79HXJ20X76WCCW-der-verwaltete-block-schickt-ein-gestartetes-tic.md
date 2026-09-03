---
id: 01M1KN1JEG6D79HXJ20X76WCCW
title: Der verwaltete Block schickt ein gestartetes Ticket bis zum naechsten Human-Step
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Sagt der User 'starte/arbeite Ticket X', faehrt die Session es Lane fuer Lane durch - Loops (critique/optimize) inklusive - bis es in einer Human-Lane liegt, und nach der Klaerung weiter; sagt der User 'ein Agent soll es bearbeiten', babysittet ein Subagent es genauso durch"
context: "Berk am 03.09.: er promptet das heute jedes Mal von Hand ('babysitte das Ticket mit einem Subagent durch bis zum naechsten human step'). Die Regel soll in den verwalteten jaira-Block der CLAUDE.md/AGENTS.md, damit sie ohne erneutes Tippen gilt. Das ist Mechanismus c aus 88H1P4 (Block-Satz), der explizit auf Berks Go wartete - diese Anfrage ist das Go, erweitert um die Babysit-Formulierung fuer Subagents. WICHTIG: den Block nie von Hand editieren - er wird von 'jaira update' regeneriert; die Aenderung gehoert in den Generator des Blocks im Go-Code. Bekannter Trap: nach der Aenderung ist der Block auf JEDEM Board stale, bis dort einmal 'jaira update' laeuft. Ergaenzend existiert seit 88H1P4 der Stop-Hook ('jaira hook print'), der das Sitzungsende bei wartenden agentischen Lanes verweigert - Block-Satz (Anweisung) und Hook (Durchsetzung) decken zusammen beide Haelften."
definition-of-done: "Der verwaltete Block enthaelt die Regel (Start -> durchfahren bis Human-Lane inkl. Loops, nach Klaerung weiter, 'Agent soll' -> Subagent babysittet); 'jaira update' schreibt sie; ein Test pinnt den neuen Blocktext"
tags: []
blocked-by: []
commits: []
created-at: 2026-09-03T12:49:02Z
updated-at: 2026-09-03T12:49:02Z
---

# Der verwaltete Block schickt ein gestartetes Ticket bis zum naechsten Human-Step

## Definition of Done

- [ ] Der verwaltete Block enthaelt die Regel (Start -> durchfahren bis Human-Lane inkl. Loops, nach Klaerung weiter, 'Agent soll' -> Subagent babysittet); 'jaira update' schreibt sie; ein Test pinnt den neuen Blocktext

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

