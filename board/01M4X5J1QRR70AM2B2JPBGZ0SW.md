---
id: 01M4X5J1QRR70AM2B2JPBGZ0SW
title: "Das Ticket-Schema aus der Spec vom 31.08. wird gebaut: Felder kommen aus TEMPLATE.md und den Lanes"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Welche Felder ein Ticket hat und welche eine Lane liefern muss, steht in einer festen Grammatik in TEMPLATE.md und den Lane-Dateien; jaira prueft sie, fuehrt eine History pro Lane-Durchlauf und erzeugt TUI und Anlege-Formular daraus"
context: |-
  Was heute fehlt: Welche Felder ein Ticket traegt und was eine Lane liefern muss, ist im Code verstreut. Eine Lane kann kein eigenes Feld mit eigener Form erklaeren, und ein Loop entscheidet nicht deterministisch, wohin er zurueckgeht.
  Das Design ist fertig und von Berk am 31.08. freigegeben: docs/superpowers/specs/2026-08-31-schema-design.md (Commit 1e8fc57c). Alle Entscheidungen stehen in .planning/schema-brainstorm.md. Erst lesen, dann bauen.
  Kurzfassung der Spec: Die Markdown-Datei ist die Wahrheit, die TUI wird daraus erzeugt. Feste Grammatik mit sechs Formen (<text> <text?> <one-of: a | b> <list> <date> <yes/no>) plus '## Heading'-Abschnitte. Creator-Felder kommen aus .jaira/TEMPLATE.md, Lane-Felder aus 'produces:' jeder Lane. Das Frontmatter haelt den letzten Wert, jeder Lane-Durchlauf haengt einen Block an '## History' an. 'routes-on:' macht das Zurueckschicken in Loops deterministisch.
  Bauordnung A in sechs Schnitten, jeder durch die Lanes des Boards.
  Warum im Milestone jaira-board-fan-out: Berk will mehrere Boards mit eigenen Lanes pro Repo (UI-Board mit Playwright-Pflicht, Video-Board ...). Damit eine Lane eines solchen Boards eigene Felder erklaeren kann, braucht es dieses Schema. NFJCTK (ein Schema sagt, was eine Lane-Datei setzen darf) gehoert dazu.
definition-of-done: "Schnitt 1: Grammatik und Validierung - jaira prueft TEMPLATE.md und 'produces:' gegen die sechs Formen und warnt bei unbekannten Schluesseln; Test"
tags: []
blocked-by: []
related:
  - 01M11YYR0SMVCXNQPZ68NFJCTK
commits: []
created-at: 2026-10-14T12:17:55Z
updated-at: 2026-10-14T12:17:55Z
---

# Das Ticket-Schema aus der Spec vom 31.08. wird gebaut: Felder kommen aus TEMPLATE.md und den Lanes

## Definition of Done

- [ ] Schnitt 1: Grammatik und Validierung - jaira prueft TEMPLATE.md und 'produces:' gegen die sechs Formen und warnt bei unbekannten Schluesseln; Test
- [ ] Schnitt 2: jeder Lane-Durchlauf haengt einen unveraenderlichen Block an '## History' an, das Frontmatter haelt den letzten Wert; Test
- [ ] Schnitt 3: 'routes-on:' entscheidet, ob ein Loop vorwaerts geht oder zurueck; Test
- [ ] Schnitt 4: die TUI faltet die History und zeigt die Runden; Test
- [ ] Schnitt 5: das Anlege-Formular wird aus dem Schema erzeugt; Test
- [ ] Schnitt 6: Doku (TICKET.md, LANES.md, 'jaira ticket template') und Release-Notes

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

