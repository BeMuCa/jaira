---
id: 01M6H73HP1Q0D6Q047PYGF2MAB
title: jaira web uebernimmt ein Ticket auf Knopfdruck aus einem Jira-Key
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Ein PM legt ein jaira-Ticket aus einem bestehenden Jira-Ticket an, mit genau einer Jira-Anfrage pro Klick."
context: |-
  Was heute fehlt: Arbeit, die in Jira beschrieben ist, muss von Hand abgetippt werden; die Verbindung geht verloren.
  Teil der Aufteilung von CMN1YT (jaira web als Zentrale fuer PMs = Leute ohne Claude und ohne CLI); dort stehen Brainstorm und Antworten des Maintainers.
  Entschieden (Maintainer): Jira nur auf Anfrage und einmalig, wegen API-Quoten. Kein Abfragen im Hintergrund, kein Dauer-Sync.
  Bekannt: das Feld 'external' ist im Ticket-Schema reserviert und wird heute nie gelesen (core/ticket/schema.go:77-78). Wert: 'jira:<KEY>'.
  Bekannt: der Token kommt aus der lokalen Token-Ablage (secrets.json); eine HTTP-Anfrage mit der Go-Standardbibliothek, kein Jira-SDK.
  Entschieden: der Onboarding-Schritt 'Anbieter: git oder Jira' kommt mit diesem Ticket, weil erst jetzt Jira etwas kann.
definition-of-done: "Im Tab 'Schreiben': Feld 'Jira-Key' und Knopf 'uebernehmen'; genau eine GET-Anfrage; Titel, Ziel und Kontext werden vorbefuellt, 'external: jira:<KEY>' wird gesetzt; Test gegen einen lokalen Fake-Server"
tags:
  - web
blocked-by:
  - CW2X0F
  - 3Q92JY
related: []
follows: 01M67K6Z3NSBTBXAK4EVCMN1YT
commits: []
created-at: 2026-11-03T17:25:27Z
updated-at: 2026-11-03T17:28:22Z
updated-by: berkc
---

# jaira web uebernimmt ein Ticket auf Knopfdruck aus einem Jira-Key

## Definition of Done

- [ ] Im Tab 'Schreiben': Feld 'Jira-Key' und Knopf 'uebernehmen'; genau eine GET-Anfrage; Titel, Ziel und Kontext werden vorbefuellt, 'external: jira:<KEY>' wird gesetzt; Test gegen einen lokalen Fake-Server
- [ ] 401, 404 und 429 werden lesbar gemeldet; bei 429 nennt die Seite die Wartezeit aus Retry-After; Test
- [ ] Kein Token in Antwort oder Log (Test aus dem Token-Ticket wiederverwendet)
- [ ] Onboarding und Projekt-Einstellungen bieten 'Anbieter: git / Jira'; Jira fuehrt zur Token-Seite; Test
- [ ] jaira web --help, eine Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-03 17:28 · berkc** — Korrektur aus dem Brainstorm (03.11.2026): Epics tragen seit PG9P9Y schon eine Liste 'external' mit Eintraegen {key, url} (core/epic/external.go:23-29). Wo dieses Ticket 'external: jira:<KEY>' schreibt, gilt stattdessen dieselbe Form {key, url} - fuer Ticket und Milestone ebenso, kein zweites Format. Am Ticket ist 'external' heute nur reserviert (core/ticket/schema.go:77-78).
