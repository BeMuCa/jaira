---
id: 01M4TRC0ZE0GNWDWA1XDK6BGB0
title: "Das Board zeigt den Stand des Refs, wenn die Ticket-Datei dahinter zurueckliegt"
status: backlog
ready: false
creator: BeMuCa
goal: "Ein Mensch sieht auf dem Board die Lane, in der ein Ticket wirklich steht, auch wenn die Datei in seinem Checkout veraltet ist."
context: |-
  Am 2026-10-13 stand D0SAHM auf Berks Board (Haupt-Checkout /home/berk/git/jAIra) in in-progress. Sein Ref refs/jaira/tickets/01M1KFMETY9MKH4V38TVD0SAHM sagte schon signoff.
  Grund: Board und 'jaira list' lesen die Datei im Checkout, und die lag hinter dem Ref. Gearbeitet wurde das Ticket in einem anderen Worktree (/home/berk/git/jAIra-D0SAHM).
  'jaira pull D0SAHM' wusste es: 'D0SAHM was behind its ref and is up to date now - was: in-progress ... now: signoff'. jaira erkennt den Fall also schon.
  Es war das zweite Mal in derselben Session. Berk haette fast ein Ticket fuer haengend gehalten, das auf seine Abnahme wartete.
  JDXV60 beschreibt dasselbe Problem zwischen zwei Worktrees (2026-09-15), vor den Refs. Das hier ist der Rest, der mit Refs bleibt.
  Offen fuer brainstorm: Board liest dann den Ref, oder zeigt nur ein Kennzeichen 'Ref ist weiter' an der Karte.
definition-of-done: "Liegt die Datei eines Tickets hinter seinem Ref, zeigt das Board auf der Karte den Stand des Refs oder ein Kennzeichen dafuer; 'jaira list' und 'jaira show' sagen es ebenso; ein Test baut genau diesen Fall"
tags:
  - concurrency
  - tui
blocked-by: []
related: []
follows: 01M2HW4788Z5G0M05G3PJDXV60
commits: []
created-at: 2026-10-13T13:48:57Z
updated-at: 2026-10-13T14:16:08Z
assignee: BeMuCa
updated-by: BeMuCa
---

# Das Board zeigt den Stand des Refs, wenn die Ticket-Datei dahinter zurueckliegt

## Definition of Done

- [ ] Liegt die Datei eines Tickets hinter seinem Ref, zeigt das Board auf der Karte den Stand des Refs oder ein Kennzeichen dafuer; 'jaira list' und 'jaira show' sagen es ebenso; ein Test baut genau diesen Fall
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

