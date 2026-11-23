---
id: 01M8543JH7M7VXA9THHZW80HNH
title: "Give back: ein gescheitertes 'herdr tab list' gilt als 'kein Tab offen' und startet einen zweiten Worker neben dem lebenden"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Kann das Cockpit die Tabs nicht lesen, startet Give back keinen neuen Worker, sondern sagt der Person, dass es den Zustand nicht pruefen konnte."
context: |-
  Was falsch ist: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1784 - scheitert 'herdr tab list' (Herdr antwortet nicht, Timeout), ist die Liste leer, Give back haelt den alten Worker fuer weg und startet einen neuen. Laeuft der alte doch, arbeiten zwei Worker am selben Ticket und Worktree.
  Gab es schon vor M31BYV (Commit 82aa5654a3); gefunden im Review zu M31BYV am 23.11.2026.
  Erst nach M31BYV starten: beide aendern denselben Give-back-Teil in register.tsx.
definition-of-done: "Scheitert 'herdr tab list', startet Give back keinen Worker und meldet der Person, dass der Tab-Zustand nicht lesbar war; Test"
tags:
  - cockpit
blocked-by: []
related: []
follows: 01M851Y62KHD9DVD2GDAM31BYV
commits: []
created-at: 2026-11-23T21:13:33Z
updated-at: 2026-11-23T21:13:33Z
---

# Give back: ein gescheitertes 'herdr tab list' gilt als 'kein Tab offen' und startet einen zweiten Worker neben dem lebenden

## Definition of Done

- [ ] Scheitert 'herdr tab list', startet Give back keinen Worker und meldet der Person, dass der Tab-Zustand nicht lesbar war; Test
- [ ] Release-Note in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

