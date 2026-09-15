---
id: 01M2HWWZ90JKR3749KXS9ZZSFT
title: Zwei von drei Commits aendern nur eine Ticket-Datei
status: pre-process
ready: true
creator: Alexander Sacharov
goal: "Ein Zweig zeigt die Arbeit, nicht die Buchhaltung: wer den Verlauf liest, sieht Aenderungen am Werkzeug und nicht jede Lane, die einen Vermerk hinterlassen hat."
definition-of-done: "Eine Lane, die keinen Code aendert, erzeugt keinen eigenen Commit mehr. Nachgestellt an einem Ticket, das critique, testing und review durchlaeuft: danach steht im Verlauf kein Commit, der nur .jaira/ anfasst."
tags:
  - cli
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-09-15T06:43:33Z
updated-at: 2026-09-15T13:04:20Z
assignee: "Alexander Sacharov"
updated-by: Alexander Sacharov
context: |-
  Gemessen am 2026-09-15 auf dem Zweig feat/9ET6NC-per-board-remote: 92 Commits, davon 30 mit Code. Zweiundsechzig fassen nur .jaira/ an - eine Lane, die ihren Vermerk hinterlaesst, ein Uebergang, eine Notiz.

  Das widerspricht der Regel, die das Projekt selbst aufgeschrieben hat. In CLAUDE.md steht, das Ticket faehrt im SELBEN Commit wie der Code, damit ein Leser die Aenderung und ihren Grund an einer Stelle sieht. Gemacht wird das Gegenteil: je Lane ein eigener Commit, der nichts als die Ticket-Datei bewegt.

  Warum es so kam: jede Lane ist ein eigener Worker, und critique, testing und review aendern gar keinen Code. Sie haben nichts, woran sie ihren Vermerk haengen koennten, also committen sie ihn allein.

  Der Punkt, der die Sache entscheidet: diese Commits sind nicht das, was den Zustand bewahrt. Record() (core/refsync/refsync.go:108) stellt JEDE Schreibung - note, move, dod - in die Outbox und schickt sie auf den Ref, unabhaengig von git. Die Commits dienen der Lesbarkeit des Pull Requests, nicht der Haltbarkeit. Wer sie weglaesst, verliert nichts Dauerhaftes.

  Was dabei nicht kaputtgehen darf: jaira leitet die Commit-Liste eines Tickets aus der Vereinigung zweier Quellen ab - der Historie der Ticket-Datei UND der Commits, die seine Id nennen. Faellt die erste Quelle weg, haengt alles daran, dass Commits die Id im Betreff tragen. Heute tun sie das ohnehin ('fix(KSGSKK): ...'), aber aus einer Gewohnheit wird damit eine Bedingung.
claimed-by: DESKTOP-RFTCH11-41016
claimed-at: 2026-09-15T13:03:43Z
---

# Zwei von drei Commits aendern nur eine Ticket-Datei

## Definition of Done

- [ ] Eine Lane, die keinen Code aendert, erzeugt keinen eigenen Commit mehr. Nachgestellt an einem Ticket, das critique, testing und review durchlaeuft: danach steht im Verlauf kein Commit, der nur .jaira/ anfasst.
- [ ] Die Commit-Liste eines Tickets bleibt vollstaendig, obwohl die Ticket-Datei seltener committet wird. Nachgestellt an einem Ticket, das die Lanes durchlaeuft und danach 'jaira move' in die Endlane erreicht - die abgeleitete Liste nennt jeden Code-Commit, der zu ihm gehoert.
- [ ] Die Regel steht dort, wo ein Agent sie liest: im erzeugten jaira-Block und in den Rollen-Prompts, nicht nur in einem Ticket.
- [ ] Eine Zeile in core/release/NOTES.md unter ## Unreleased, falls sich etwas an der Ableitung oder am Verhalten der Befehle aendert.

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

