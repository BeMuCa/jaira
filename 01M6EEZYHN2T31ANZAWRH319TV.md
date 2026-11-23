---
id: 01M6EEZYHN2T31ANZAWRH319TV
title: "README: Abschnitt 'Jaira PM View - Web only' fuer Leute ohne Coding-Agent"
status: in-progress
ready: true
creator: BeMuCa
goal: "Die README hat einen eigenen Abschnitt, mit dem jemand ohne Agent und ohne Terminal-Board in fuenf Minuten jaira web auf einem geklonten Repo laufen hat: Install-Zeile ohne Rollen/Cockpit, jaira web, ein Screenshot, was man dort kann und was nicht; die Zeile 'Web UI' in der Tabelle 'Not built' ist weg."
context: |-
  Was fehlt: jaira web (Epics sortieren, Tickets lesen) ist fuer PMs und Leute ohne Coding-Agent gedacht, aber nirgends fuer sie erklaert. Es gibt jaira web --help, eine Zeile in der README-Befehlstabelle (Zeile 850) und einen Absatz in docs/COMMANDS.md Zeile 57. Kein Tutorial, kein Screenshot unter docs/img.
  Widerspruch: README Zeile 1022 fuehrt 'Web UI' in der Tabelle 'Not built' als bewusst nicht gebaut.
  Zielgruppe (Wunsch 2026-10-28/11-02): Leute, die jaira nur ziehen, um die Web-Ansicht zu nutzen - ohne Claude, ohne Terminal-Board.
  Bekannt: die Seite steckt im Binary (core/web/web.go:30, go:embed), es braucht kein Node und kein zweites Programm. Der curl-Installer richtet seit V38STG auch Rollen und Cockpit ein; der Abschnitt muss die Variante ohne beides nennen.
  Bekannt: ein frischer Clone hat die Ticket-Refs nicht (origin holt nur refs/heads/*); jaira holt refs/jaira/tickets/* selbst (core/refsync/refsync.go:933). Der Abschnitt sagt, wann die Tickets nach dem Start erscheinen.
  Soll als eigener PR kommen, nur README und Bild.
definition-of-done: "README hat den Abschnitt 'Jaira PM View - Web only' mit Install-Zeile ohne Rollen/Cockpit, jaira web, Screenshot und je einem Satz was geht und was nicht"
tags:
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-11-02T15:45:35Z
updated-at: 2026-11-23T01:44:48Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: worker-h319tv
claimed-at: 2026-11-23T01:38:34Z
---

# README: Abschnitt 'Jaira PM View - Web only' fuer Leute ohne Coding-Agent

## Definition of Done

- [x] README hat den Abschnitt 'Jaira PM View - Web only' mit Install-Zeile ohne Rollen/Cockpit, jaira web, Screenshot und je einem Satz was geht und was nicht
  proof: README.md:300 '## Jaira PM View - Web only' (install lines with JAIRA_NO_ROLES=1, jaira init/web, docs/img/web.png, one sentence each for can/cannot)
- [x] Die Zeile 'Web UI, VS Code extension' in der Tabelle 'Not built' ist entfernt oder sagt, was heute gilt
  proof: README.md:1247 row now '| VS Code extension | ...', Web UI dropped
- [ ] Auf einem frischen Clone nachgestellt: Install-Zeile, jaira web, Tickets erscheinen; die Wartezeit steht im Abschnitt
- [ ] Release-Notiz in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 01:39 · BeMuCa** — Entscheidungen (kein Punkt offen fuer die Person): 1) Abschnitt steht als eigenes '## Jaira PM View - Web only' direkt nach '## Start', nicht im Install-Abschnitt - feat/q9f9w4 aendert den Install-Text, so bleibt der Merge trivial. 2) Install-Zeile mit JAIRA_NO_ROLES=1 (so geschrieben in scripts/install.sh:110 und install.ps1:159). 3) 'Not built': nur 'Web UI' faellt weg, die Zeile bleibt fuer 'VS Code extension' mit ihrem Grund. 4) Release-Notiz wird geschrieben, weil DoD 4 sie verlangt, obwohl der Kontext 'nur README und Bild' sagt und CLAUDE.md README-Aenderungen nicht als client-facing zaehlt - die Person kann sie streichen. 5) Screenshot mit chrome-headless-shell aus ~/.cache/ms-playwright auf einem Scratch-Board, nach docs/img/web.png. Kein Plan-Lane: Option 'planning' nicht angehakt.
