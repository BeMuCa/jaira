---
id: 01M5MQKJT008A580WNKNSJ16MC
title: "Worker-Rolle: kein Version-Bump/Changelog, nur eigene Dateien formatieren"
status: backlog
ready: false
creator: Alexander Sacharov
goal: "Ein Worker ändert weder Version noch Changelog des Repos und formatiert nur Dateien, die er selbst geändert hat."
context: |-
  Gefunden am 23.10.2026 in einer langen Teamlead-Session für requirementsgenie (Notizen: /home/alex/.claude/acceptance/requirementsgenie/teamlead-staerken.md).
  Was falsch ist: Worker bumpten package.json und den Changelog selbst — beim Zusammenführen mehrerer Tickets auf einem Abnahme-Stand gab es bei jedem Ticket Konflikte; der Teamlead vergibt die Version beim Release. Außerdem formatiert 'task code:format' (ruff format) das ganze Repo (238 Dateien), Worker mussten das jedes Mal zurücknehmen.
  Beides wurde per Hinweis an jeden einzelnen Worker gelöst — gehört in die Rolle.
  Wo: Rollen jaira-worker / jaira-role-lane (core/role/builtin). Formatieren: nur 'git diff --name-only' der eigenen Änderung an den Formatter geben. Repos, die den Bump im Ticket wollen, brauchen eine Einstellung (z. B. settings.json 'release': 'teamlead'|'worker'), Standard so wie heute dokumentiert.
definition-of-done: "Rolle sagt: Version/Changelog nicht anfassen, wenn das Repo es dem Teamlead überlässt (Einstellung)"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T15:55:52Z
updated-at: 2026-10-23T15:55:52Z
---

# Worker-Rolle: kein Version-Bump/Changelog, nur eigene Dateien formatieren

## Definition of Done

- [ ] Rolle sagt: Version/Changelog nicht anfassen, wenn das Repo es dem Teamlead überlässt (Einstellung)
- [ ] Rolle sagt: nur eigene geänderte Dateien formatieren

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

