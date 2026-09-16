---
id: 01M2E6J3EM92YD0WH33TFTMYSD
title: "Nach Eltern-Ticket filtern: das Board zeigt nur die Kinder"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Eltern-Ticket laesst sich fokussieren: 'jaira list --parent <id>' und ein Fokus-Modus im TUI zeigen nur noch die Tickets, die unter diesem Ticket haengen"
context: |-
  Das Datenmodell kann Eltern und Kinder schon, die Ansicht nicht.
  - Das Feld 'parent' steht in core/ticket/schema.go:195. Ein Kind nennt seinen Elternteil, Kinder werden nie gespeichert, der Baum wird zurueckgelesen.
  - 'jaira create --parent <id>' setzt es.
  - 'jaira links <id>' zeigt den Baum in beide Richtungen, auch ueber Ref, Logbuch und Archiv hinweg.
  Was fehlt, ist das Filtern. 'jaira list' kennt heute --lane, --tag, --assignee, --actionable und --query, aber kein --parent. Im TUI gibt es ueberhaupt keinen Bezug auf parent - internal/tui nennt das Wort nur in einem Kommentar in links.go.
  Folge: wer ein Epic mit zwanzig Kindern hat, sieht sie ueber das ganze Board verstreut und kann nirgends sagen 'zeig mir nur das hier'.
  Zu klaeren in der Brainstorm-Lane, hier absichtlich nicht entschieden:
  - nur direkte Kinder oder der ganze Unterbaum. Ein Epic kann beliebig tief geschachtelt sein.
  - wie man den Fokus betritt und verlaesst, und woran man im TUI sieht, dass man drin ist. Ein gefiltertes Board, das aussieht wie ein leeres Board, ist ein Fehlerbericht.
  - ob der Fokus die Lane-Spalten behaelt oder nur eine Liste wird.
  Spannung mit D0SAHM, die vor dem Bauen aufgeloest gehoert: jenes Ticket gruppiert Epics ueber den TAG, dieses ueber parent. Zwei Gruppierungen nebeneinander sind zwei Arten, dasselbe zu sagen, und der Mensch muss dann raten, welche gilt. Entweder eine gewinnt, oder es wird begruendet, wofuer je eine da ist.
definition-of-done: "'jaira list --parent <id>' listet die Kinder des Tickets, mit --json in gleicher Form wie sonst; das TUI hat einen Fokus-Modus, der das Board auf ein Eltern-Ticket einschraenkt und sichtbar anzeigt, dass er aktiv ist, und sich wieder verlassen laesst; ein Ticket ohne Kinder meldet das statt leer zu wirken; eine Zeile in core/release/NOTES.md unter ## Unreleased; go test ./... -race gruen"
tags:
  - tui
  - cli
blocked-by: []
related:
  - 01M1KFMETY9MKH4V38TVD0SAHM
commits: []
created-at: 2026-09-13T20:15:25Z
updated-at: 2026-09-13T20:15:25Z
---

# Nach Eltern-Ticket filtern: das Board zeigt nur die Kinder

## Definition of Done

- [ ] 'jaira list --parent <id>' listet die Kinder des Tickets, mit --json in gleicher Form wie sonst; das TUI hat einen Fokus-Modus, der das Board auf ein Eltern-Ticket einschraenkt und sichtbar anzeigt, dass er aktiv ist, und sich wieder verlassen laesst; ein Ticket ohne Kinder meldet das statt leer zu wirken; eine Zeile in core/release/NOTES.md unter ## Unreleased; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

