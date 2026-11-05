---
id: 01M6NJXCBEV521DR64V3W8WH0M
title: jaira web zeigt die Seiten installierter Plugins und leitet ihre API-Aufrufe an das Plugin weiter
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Plugin bringt eigene Seiten in jaira web mit, z.B. eine eigene Jira-Ansicht, ohne dass jaira web dafuer geaendert wird."
context: |-
  Heute kennt jaira web nur seine eingebauten Seiten. Wunsch von Alex am 05.11.2026: ein Plugin soll im Web etwas Eigenes zeigen koennen, z.B. die Jira-Integration anders dargestellt.
  Vertrag steht in docs/PLUGINS.md, Abschnitt Web (R6X6FH):
  - Manifest web.pages [{path, title}] -> je Seite ein Tab/Link in jaira web.
  - Dateien aus ~/.jaira/plugins/<id>/web/ unter /plugins/<id>/ ausliefern.
  - /plugins/<id>/api/<pfad> CGI-artig an 'jaira-<id> web <METHOD> /<pfad>' weiterleiten: Body auf stdin, JSON auf stdout. Kein dauerhaft laufender Plugin-Prozess, kein eigener Port.
  - Nur Plugins mit capability 'web' bekommen Seiten.
  Beruehrt die Web-Tickets von Berk (CMN1YT und Teile): mit ihm abstimmen, wo die Plugin-Tabs in der Tab-Leiste (V9S8RN) sitzen.
definition-of-done: Ein installiertes Test-Plugin mit web.pages erscheint als Tab in jaira web und seine Seite laedt; Test
tags:
  - web
blocked-by:
  - R6X6FH
  - 8Z3VRF
related: []
commits: []
created-at: 2026-11-05T10:08:46Z
updated-at: 2026-11-05T10:08:46Z
---

# jaira web zeigt die Seiten installierter Plugins und leitet ihre API-Aufrufe an das Plugin weiter

## Definition of Done

- [ ] Ein installiertes Test-Plugin mit web.pages erscheint als Tab in jaira web und seine Seite laedt; Test
- [ ] /plugins/<id>/api/x ruft 'jaira-<id> web GET /x' auf und gibt dessen JSON zurueck; ein Plugin ohne capability web bekommt 404; Test
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

