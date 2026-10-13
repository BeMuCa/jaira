---
id: 01M4TRC37SNBER05X090DPSFGE
title: "git status im Haupt-Checkout zeigt keine Ticket-Dateien, die kein Commit mehr mitnimmt"
status: backlog
ready: false
creator: BeMuCa
goal: "Ein Blick auf git status zeigt nur Aenderungen, die jemand committen soll; Ticket-Kopien, die nur ihr Ref traegt, stehen dort nicht als Rauschen."
context: |-
  Am 2026-10-13 zeigte git status in /home/berk/git/jAIra (Branch feat/PAP369-link-lines): 3 geaenderte und 6 ungetrackte Dateien unter .jaira/tickets/ und ungetrackt .jaira/milestones/.
  Keine davon gehoert zur Arbeit dieses Branches. Woher sie kommen:
  - Lanes ohne Code committen nach Regel nichts, die Datei bleibt geaendert liegen (CLAUDE.md, 'a lane that changed no code commits nothing').
  - 'jaira pull <id>' und Schreibbefehle auf ein Ticket, das nur auf seinem Ref liegt, legen die Datei in den Checkout, in dem man gerade ist ('pull it first').
  - Milestone-Dateien landen ungetrackt unter .jaira/milestones/.
  Folge: Man sieht nicht, welche Kopie gilt, und eine veraltete Kopie wird zur falschen Anzeige (siehe das Ticket 'Board zeigt den Stand des Refs', JDXV60).
  Offen fuer brainstorm: Dateien nur dort, wo ein Ticket bearbeitet wird; oder jaira raeumt sie weg; oder jaira sagt pro Datei, wohin sie gehoert.
definition-of-done: "In einem Checkout, der ein Ticket nicht bearbeitet, hinterlassen create, pull, note und move auf dieses Ticket nichts in git status, oder jaira nennt fuer jede solche Datei, wohin sie gehoert; die Wahl ist im Ticket notiert; ein Test baut den Fall"
tags:
  - concurrency
blocked-by: []
related: []
commits: []
created-at: 2026-10-13T13:49:00Z
updated-at: 2026-10-13T13:49:00Z
---

# git status im Haupt-Checkout zeigt keine Ticket-Dateien, die kein Commit mehr mitnimmt

## Definition of Done

- [ ] In einem Checkout, der ein Ticket nicht bearbeitet, hinterlassen create, pull, note und move auf dieses Ticket nichts in git status, oder jaira nennt fuer jede solche Datei, wohin sie gehoert; die Wahl ist im Ticket notiert; ein Test baut den Fall
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

