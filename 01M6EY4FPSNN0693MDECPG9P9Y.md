---
id: 01M6EY4FPSNN0693MDECPG9P9Y
title: "Ein Epic traegt einen externen Schluessel mit Link, und Epics und Tickets lassen sich danach suchen und filtern"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: Ein Epic hat ein eigenes Feld fuer den Schluessel einer Aufgabe in einem anderen System (z. B. JIRA-1234) mit URL. Danach kann man in CLI und web Epics und ihre Tickets suchen und filtern.
context: |-
  Warum (Alex, 02.11.2026):
  - Alex will Stunden ueber echte JIRA-Stories abrechnen: pro JIRA-Story ein Epic, die Tickets darin sind die Arbeit, die darueber bezahlt wird.
  - Die Abrechnung soll NICHT in Code-Commits landen. Deshalb das Feld auf dem Epic (eigener Ref refs/jaira/epics/<name>), nicht auf dem Ticket: die Ticket-Datei wird mit dem Code committet, der Epic-Ref nicht.
  - JIRA-Stories und Epics ueberschneiden sich bei der Bezahlung. Ein Ticket darf also in zwei Epics liegen (geht heute schon mit jaira epic add).
  Was heute fehlt:
  - Epic-Frontmatter hat nur title, description, created-at (core/epic/epic.go). Einen Schluessel kann man nur in description schreiben: nicht filterbar, kein Link.
  Entwurf (brainstorm darf aendern):
  - Feld 'external:' im Epic-Frontmatter, eine Liste, je Zeile '<SCHLUESSEL> <URL>'. Allgemein, nicht 'jira' - es gibt auch YouTrack.
  - jaira epic create|set --external '<KEY> <URL>' (wiederholbar), --external-rm <KEY>.
  - jaira epic ls --external <KEY> und jaira list --external <KEY> (Tickets aller Epics mit diesem Schluessel).
  - web: Schluessel als Link am Epic, Filter nach Schluessel.
  Achtung: der Epic-Ref wird zum Remote gepusht - alle mit Zugriff aufs Repo sehen die Schluessel. Fuer Alex ok (bewusst gewaehlt statt einer lokalen Datei).
  Haengt an 182JT1 (PR #120, epic close), das dieselben Dateien in core/epic und web aendert: erst nach dessen Merge anfangen.
  Spaeter, nicht hier: Stundensumme je Schluessel in /jaira-role-logbook-summary.
definition-of-done: "Epic-Frontmatter hat 'external:' (Liste '<KEY> <URL>'); jaira epic create/set --external fuegt hinzu, --external-rm nimmt weg; Datei bleibt von Hand editierbar; Test"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-11-02T20:10:12Z
updated-at: 2026-11-02T20:10:43Z
updated-by: Alexander Sacharov
---

# Ein Epic traegt einen externen Schluessel mit Link, und Epics und Tickets lassen sich danach suchen und filtern

## Definition of Done

- [ ] Epic-Frontmatter hat 'external:' (Liste '<KEY> <URL>'); jaira epic create/set --external fuegt hinzu, --external-rm nimmt weg; Datei bleibt von Hand editierbar; Test
- [ ] jaira epic show und epic ls (auch --json) zeigen die Schluessel; jaira epic ls --external <KEY> filtert; Test
- [ ] jaira list --external <KEY> zeigt die Tickets aller Epics mit diesem Schluessel, Folge-Tickets eingeschlossen; Test
- [ ] web: Schluessel am Epic als klickbarer Link, setzen/entfernen im Epic-Editor, Filter nach Schluessel; Test
- [ ] Kein Ticket-File und kein Code-Commit traegt den Schluessel - er steht nur im Epic-Ref; Test
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

