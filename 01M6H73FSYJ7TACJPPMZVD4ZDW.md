---
id: 01M6H73FSYJ7TACJPPMZVD4ZDW
title: "jaira web: die Zuordnung jaira -> Jira (Epic, Milestone, Road, Ticket, Lane) ist auf einer Seite einstellbar, mit Default"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Ein PM stellt einmal ein, wie Epics, Milestones, Roads, Tickets und Lanes in seinem Jira heissen, bevor irgendetwas uebertragen wird."
context: |-
  Was heute fehlt: jaira kennt Jira nicht. Jede Firma nutzt Jira anders (Fix Version oder Sprint, mit oder ohne Initiativen), daher reicht eine feste Zuordnung nicht.
  Teil der Aufteilung von CMN1YT (jaira web als Zentrale fuer PMs = Leute ohne Claude und ohne CLI); dort stehen Brainstorm und Antworten des Maintainers.
  Entschieden (Maintainer): Zuordnung im Web einstellbar, Default wie im Brainstorm vorgeschlagen.
  Default: Epic -> Jira Epic. Milestone -> Fix Version (Alternativen: Sprint, Epic). Road -> Label (Initiative nur mit Jira Premium). Ticket -> Task (Typ aus Tag, sonst Task). DoD-Punkte -> Checkliste im Beschreibungstext, NICHT Subtasks (Subtasks haben eigenen Status, DoD-Punkte nicht). follows/blocked-by -> Issue Links (relates/blocks). Lane -> Jira-Status ueber eine Tabelle (vor in-progress -> To Do, Agenten-Lanes -> In Progress, Review/Signoff -> In Review, Terminal -> Done).
  Entschieden (Empfehlung im CMN1YT-Brainstorm, aenderbar dort): die Zuordnung ist eine committete Board-Datei ohne Adresse und ohne Token, damit alle PMs eines Teams dieselbe benutzen. Die Jira-Adresse steht nur in der lokalen Token-Ablage.
  ABHAENGIGKEIT: die Zeile 'Road' folgt dem Modell aus KQQD1H.
  Ausgeschlossen: GitHub/GitLab/YouTrack - Jira nur auf Anfrage, andere Anbieter erst bei Bedarf (siehe 1YTBPX).
definition-of-done: Einstellungsseite zeigt die Default-Zuordnung als Tabelle; jede Zeile hat eine Auswahl der erlaubten Jira-Typen; Speichern schreibt die Board-Datei; Test
tags:
  - web
blocked-by:
  - 3Q92JY
  - KQQD1H
related: []
follows: 01M67K6Z3NSBTBXAK4EVCMN1YT
commits: []
created-at: 2026-11-03T17:25:25Z
updated-at: 2026-11-03T22:58:14Z
updated-by: berkc
---

# jaira web: die Zuordnung jaira -> Jira (Epic, Milestone, Road, Ticket, Lane) ist auf einer Seite einstellbar, mit Default

## Definition of Done

- [ ] Einstellungsseite zeigt die Default-Zuordnung als Tabelle; jede Zeile hat eine Auswahl der erlaubten Jira-Typen; Speichern schreibt die Board-Datei; Test
- [ ] Die Board-Datei ist von Hand lesbar und enthaelt weder Adresse noch Token; Test
- [ ] Lane-Zeilen werden aus den Lanes des Boards erzeugt, nicht fest verdrahtet; Test mit zwei Boards
- [ ] docs/COMMANDS.md, eine Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

