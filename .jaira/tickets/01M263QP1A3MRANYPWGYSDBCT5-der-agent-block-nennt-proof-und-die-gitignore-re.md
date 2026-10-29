---
id: 01M263QP1A3MRANYPWGYSDBCT5
title: "Der Agent-Block nennt --proof, und die Gitignore-Regel verweist auf die Commit-Regel"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Zwei nachgemessene Luecken im generierten Block sind geschlossen: der Nachweis beim Abhaken, und der Widerspruch auf einem ungeteilten Board"
context: |-
  Zwei Befunde aus einer Multi-Agent-Session am 10.09., beide nachgeprueft und beide enger als ursprünglich gemeldet.
  1. --proof fehlt im Block. 'grep -c proof' in der CLAUDE.md eines anderen Projekts ist 0. NICHT zutreffend ist die Meldung, es sei gar nicht auffindbar: 'jaira dod --help' dokumentiert --proof an sechs Stellen samt Beispiel. Es fehlt allein im Block - also genau dort, wo eine Session es liest. Folge: die ersten Tickets wurden ohne Nachweis abgehakt.
  2. Der Gitignore-Vorbehalt und die Commit-Regel widersprechen sich auf einem ungeteilten Board. Der Block sagt 'the ticket rides in the same commit as the code', aber auf einem Board mit gitignoretem .jaira/ KANN die Ticketdatei nicht mitreisen. NICHT zutreffend ist die Meldung, die Reihenfolge sei falsch: der Vorbehalt steht auf Zeile 216, die Regel auf 220 - erst der Zustand, dann die Regel. Was fehlt, ist der Verweis: der Vorbehalt sagt nicht, dass der naechste Punkt fuer dieses Board nicht gilt.
  Ein anderes Projekt ist so ein Board: .gitignore Zeile 42 enthaelt /.jaira/, null getrackte Dateien darunter.
definition-of-done: "Der generierte Block nennt --proof beim Abhaken; der Gitignore-Vorbehalt sagt ausdruecklich, dass die Commit-Regel auf einem ungeteilten Board anders lautet; 'jaira update' schreibt beides auf ein bestehendes Board; go test ./... -race gruen"
tags:
  - docs
blocked-by: []
commits: []
created-at: 2026-09-10T16:52:07Z
updated-at: 2026-09-10T16:52:07Z
---

# Der Agent-Block nennt --proof, und die Gitignore-Regel verweist auf die Commit-Regel

## Definition of Done

- [ ] Der generierte Block nennt --proof beim Abhaken; der Gitignore-Vorbehalt sagt ausdruecklich, dass die Commit-Regel auf einem ungeteilten Board anders lautet; 'jaira update' schreibt beides auf ein bestehendes Board; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

