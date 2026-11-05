---
id: 01M6H739H9X7CQJH16WPNEYWJH
title: jaira web startet auch ohne git und ohne Projekt und fuehrt im Browser durch den ersten Start
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Ein PM, der nur jaira installiert hat, bekommt im Browser gesagt, was fehlt, und kann einen git-Pfad eintragen statt im Terminal zu suchen."
context: |-
  Was heute fehlt: ohne git oder ausserhalb eines Repos bricht jaira ab; ein PM ohne Terminal-Erfahrung bleibt dort stehen.
  Teil der Aufteilung von CMN1YT (jaira web als Zentrale fuer PMs = Leute ohne Claude und ohne CLI); dort stehen Brainstorm und Antworten des Maintainers.
  Entschieden (Maintainer): Onboarding im Browser, nicht im curl-Skript. scripts/install.sh bleibt nicht interaktiv, weil Fragen ueber curl|sh unter Windows, WSL und CI leicht brechen.
  Bekannt: git wird an drei Stellen per exec.LookPath gesucht: core/gitref/gitref.go:200, core/gitrepo/git.go:29, core/ticket/store.go:878. Ein eingetragener Pfad muss an allen drei gelten.
  Entschieden: Pfad als 'git-path' in ~/.jaira/settings.json, Umgebungsvariable JAIRA_GIT als Alternative. Keine Shell-Profile veraendern.
  Der Name im Browser ist Teil von DVPAPE, nicht dieses Tickets. Projektwahl ist das naechste Ticket.
definition-of-done: "'jaira web' ausserhalb eines Repos oder ohne git startet trotzdem und zeigt eine Einstiegsseite statt abzubrechen; Test"
tags:
  - web
blocked-by: []
related: []
follows: 01M67K6Z3NSBTBXAK4EVCMN1YT
commits: []
created-at: 2026-11-03T17:25:19Z
updated-at: 2026-11-03T17:25:19Z
---

# jaira web startet auch ohne git und ohne Projekt und fuehrt im Browser durch den ersten Start

## Definition of Done

- [ ] 'jaira web' ausserhalb eines Repos oder ohne git startet trotzdem und zeigt eine Einstiegsseite statt abzubrechen; Test
- [ ] Fehlt git: die Seite nennt git-scm.com/downloads und hat ein Feld 'Pfad zu git'; der Pfad wird geprueft (git --version) und als git-path gespeichert; Test
- [ ] git-path und JAIRA_GIT gelten an allen drei LookPath-Stellen; Test
- [ ] jaira web --help, docs/COMMANDS.md, eine Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

