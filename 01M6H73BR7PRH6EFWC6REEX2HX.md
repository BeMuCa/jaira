---
id: 01M6H73BR7PRH6EFWC6REEX2HX
title: "jaira web: Projekt aus gefundenen Repos waehlen oder per Clone-Link holen, und zwischen Projekten umschalten"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: Ein PM bindet ein Projekt im Browser an - aus den Repos auf seiner Platte oder per Clone-Link - und wechselt spaeter zwischen seinen Projekten.
context: |-
  Was heute fehlt: jaira web zeigt nur das Repo, in dem es gestartet wurde.
  Teil der Aufteilung von CMN1YT (jaira web als Zentrale fuer PMs = Leute ohne Claude und ohne CLI); dort stehen Brainstorm und Antworten des Maintainers.
  Bekannt: core/project kann schon Repos mit Board finden (Discover, core/project/discover.go:25, 2 Ebenen tief) und merkt sich Projekte (Load/Remember, core/project/project.go:41/65, Datei ~/.jaira/projects.json).
  Entschieden: der Clone laeuft ueber git selbst; die Anmeldung uebernimmt git mit seinem Credential-Helper. jaira speichert dafuer KEIN Passwort und keinen Token.
  Entschieden: in diesem Ticket nur der Anbieter git. Jira kommt mit eigenem Ticket dazu.
  Voraussetzung: Einstiegsseite aus dem Ticket davor.
definition-of-done: Die Einstiegsseite listet die gefundenen Repos mit Board; ein Klick waehlt eines und merkt es in projects.json; Test
tags:
  - web
blocked-by:
  - NEYWJH
related: []
follows: 01M67K6Z3NSBTBXAK4EVCMN1YT
commits: []
created-at: 2026-11-03T17:25:21Z
updated-at: 2026-11-03T17:25:21Z
---

# jaira web: Projekt aus gefundenen Repos waehlen oder per Clone-Link holen, und zwischen Projekten umschalten

## Definition of Done

- [ ] Die Einstiegsseite listet die gefundenen Repos mit Board; ein Klick waehlt eines und merkt es in projects.json; Test
- [ ] Wird keins gefunden, ist 'Clone-Link einfuegen' die Vorauswahl; git clone in einen waehlbaren Ordner, Fehler von git lesbar angezeigt; Test mit lokalem Remote
- [ ] Eine Projektauswahl oben wechselt zwischen gemerkten Projekten ohne Neustart; Test
- [ ] jaira web --help, eine Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

