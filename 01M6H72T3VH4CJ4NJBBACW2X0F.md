---
id: 01M6H72T3VH4CJ4NJBBACW2X0F
title: jaira web legt Tickets an und setzt den Assignee
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Ein PM legt im Browser ein vollstaendiges Ticket an und weist es einem Entwickler zu, ohne Terminal."
context: |-
  Was heute fehlt: im Browser kann man Tickets nur lesen ('Nothing is edited there', internal/cli/web.go). Anlegen und Zuweisen geht nur mit 'jaira create' / 'jaira assign' im Terminal.
  Teil der Aufteilung von CMN1YT (jaira web als Zentrale fuer PMs = Leute ohne Claude und ohne CLI); dort stehen Brainstorm und Antworten des Maintainers.
  Bekannt: ein Ticket reist auf seinem Ref refs/jaira/tickets/*; 'jaira fetch' meldet dem Entwickler ein neu zugewiesenes Ticket schon per Desktop-Benachrichtigung (docs/COMMANDS.md). Der Web-Weg muss genau denselben Code wie die CLI benutzen, keinen zweiten.
  Bekannt: ein fehlgeschlagener Push bleibt in der Warteschlange von jaira web und geht mit der naechsten Aenderung raus.
  Bekannt: Assignee-Namen kommen aus DVPAPE (fester jaira-Name pro Person); ohne DVPAPE zeigt die Auswahl doppelte Namen derselben Person.
  Offen in CMN1YT, nicht hier: ob der PM auch Lanes bewegen darf. Dieses Ticket bewegt keine Lane.
definition-of-done: "Tab 'Schreiben': Formular mit Titel, Ziel, Kontext, DoD-Punkten (mehrere), Tags (Auswahl aus 'jaira tags') und Assignee (Auswahl aus bekannten Namen); fehlt Ziel, Kontext, DoD oder Assignee, nennt die Seite das fehlende Feld; Test"
tags:
  - web
blocked-by:
  - V9S8RN
  - DVPAPE
related: []
follows: 01M67K6Z3NSBTBXAK4EVCMN1YT
commits: []
created-at: 2026-11-03T17:25:03Z
updated-at: 2026-11-03T17:25:03Z
---

# jaira web legt Tickets an und setzt den Assignee

## Definition of Done

- [ ] Tab 'Schreiben': Formular mit Titel, Ziel, Kontext, DoD-Punkten (mehrere), Tags (Auswahl aus 'jaira tags') und Assignee (Auswahl aus bekannten Namen); fehlt Ziel, Kontext, DoD oder Assignee, nennt die Seite das fehlende Feld; Test
- [ ] Das angelegte Ticket ist Byte fuer Byte gleich aufgebaut wie eines aus 'jaira create' mit denselben Werten; Test vergleicht beide
- [ ] In der Ticket-Ansicht laesst sich der Assignee aendern; geschrieben ueber denselben Pfad wie die CLI; Test
- [ ] Wird der Push abgelehnt (z.B. nur Leserecht), zeigt die Seite das lesbar an und das Ticket geht nicht verloren; Test mit einem Remote, das Push verweigert
- [ ] jaira web --help, docs/COMMANDS.md, eine Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

