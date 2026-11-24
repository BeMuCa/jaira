---
id: 01M8543JH7M7VXA9THHZW80HNH
title: "Give back: ein gescheitertes 'herdr tab list' gilt als 'kein Tab offen' und startet einen zweiten Worker neben dem lebenden"
status: critique
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
updated-at: 2026-11-24T07:01:50Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-46440
claimed-at: 2026-11-24T06:58:04Z
outcome-what: "Give back stops when 'herdr tab list' fails: no new worker, no keys into the old pane, a note and a toast saying the tab state was unreadable"
outcome-why: "A failed tab list read as 'no tab open', so Give back started a second worker beside a live one on the same ticket and worktree"
outcome-resolves: Two workers on one ticket after a Herdr timeout during Give back
---

# Give back: ein gescheitertes 'herdr tab list' gilt als 'kein Tab offen' und startet einen zweiten Worker neben dem lebenden

## Definition of Done

- [x] Scheitert 'herdr tab list', startet Give back keinen Worker und meldet der Person, dass der Tab-Zustand nicht lesbar war; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1790; test 'give back > a failed herdr tab list starts no worker, types nothing and says the tab state was unreadable' (tests/board.test.ts:2527)
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/W80HNH.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 21:59 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/48662.sock
- **2026-11-24 06:58 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/3417.sock
- **2026-11-24 07:01 · Alexander Sacharov** — decision: a failed tab list also stops the hand-over to the recorded old pane (no send-text), not only the new worker - without the tab list Give back cannot tell a working or dialog-blocked worker from a gone one, and typing blind could answer a dialog. Garbage stdout with exit 0 still reads as no tab (tabsOf unchanged; Herdr does not answer that way).
