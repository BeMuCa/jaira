---
id: 01M5MQKJT008A580WNKNSJ16MC
title: "Worker-Rolle: kein Version-Bump/Changelog, nur eigene Dateien formatieren"
status: in-progress
ready: true
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
updated-at: 2026-10-23T16:01:00Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-54412
claimed-at: 2026-10-23T15:56:31Z
mode: conversational
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

- [x] core/settings: Board-Datei .jaira/release lesen (BoardRelease), Standard worker, Tippfehler = worker
- [x] jaira roles list: release + release_source in --json und Textausgabe, --help nennt die Datei
- [~] Rollen jaira-worker + jaira-role-lane: Abschnitt Version/Changelog (liest release) und Formatieren nur eigener Dateien
- [ ] Tests (settings, roles list, Rollentext) + core/release/notes/SJ16MC.md

## Progress
- **2026-10-23 15:56 · Alexander Sacharov** — parent session: requirementsgenie-92
- **2026-10-23 15:57 · Alexander Sacharov** — question (open decision): Wo lebt die Einstellung 'Version/Changelog gehört dem Teamlead'? 1) (Recommended) Board-Datei .jaira/release (eine Zeile teamlead|worker, committed, wie .jaira/chat-model) — gilt für alle im Repo, Rolle liest sie über 'jaira roles list --json | jaira json get release'; Standard worker = heutiges Verhalten. 2) ~/.jaira/settings.json 'release' — pro Maschine, wie im Ticket vorgeschlagen; gilt dann aber für jedes Repo auf der Maschine (dasselbe Problem wie 'remote'). 3) Keine Einstellung in jaira: Rolle sagt nur 'wenn CLAUDE.md/AGENTS.md des Repos Version/Changelog dem Release überlässt, nicht anfassen' — kein Go-Code. Formatieren (DoD 2) ist nicht offen: nur 'git diff --name-only' der eigenen Änderung an den Formatter, ohne Einstellung.
- **2026-10-23 15:59 · Alexander Sacharov** — decision (parent requirementsgenie-92, 23.10.): Option 1 — Board-Datei .jaira/release, eine Zeile teamlead|worker, committed; Standard worker (heutiges Verhalten). Rolle liest sie über 'jaira roles list --json'. requirementsgenie setzt danach teamlead.
