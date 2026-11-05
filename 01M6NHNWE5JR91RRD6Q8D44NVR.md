---
id: 01M6NHNWE5JR91RRD6Q8D44NVR
title: "panda-bobr/jaira-releases: plugins.json, CI-Pruefung und CONTRIBUTING"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Autor bringt sein Plugin per PR in den Katalog, und die CI sagt ihm vor jedem Menschen, was fehlt."
context: |-
  Das Repo https://github.com/panda-bobr/jaira-releases ist seit 05.11.2026 leer.
  Inhalt: plugins.json (leeres Array), README, CONTRIBUTING.md (Fork, Eintrag anhaengen, PR), GitHub-Action, die einen PR prueft.
  CI-Pruefung: id eindeutig und ohne 'jaira'; repo existiert und ist oeffentlich; jaira-plugin.json am Default-Branch passt zum Schema (R6X6FH); Release v<version> hat Assets fuer linux/darwin/windows amd64+arm64 und checksums.txt; PR-Autor ist Besitzer des repo (wie Obsidians GitHub-Verknuepfung).
  Arbeit liegt in einem anderen Repo als jaira: Branch und PR dort, das Ticket faehrt mit einem Notiz-Commit hier.
definition-of-done: "plugins.json, README und CONTRIBUTING.md liegen auf main von panda-bobr/jaira-releases"
tags:
  - ci
blocked-by:
  - R6X6FH
related: []
commits: []
created-at: 2026-11-05T09:47:12Z
updated-at: 2026-11-05T09:47:12Z
---

# panda-bobr/jaira-releases: plugins.json, CI-Pruefung und CONTRIBUTING

## Definition of Done

- [ ] plugins.json, README und CONTRIBUTING.md liegen auf main von panda-bobr/jaira-releases
- [ ] Die Action lehnt einen PR mit doppelter id, fehlendem Release oder fehlender Pruefsumme ab; in einem Test-PR gezeigt

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

