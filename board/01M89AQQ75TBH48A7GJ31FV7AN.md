---
id: 01M89AQQ75TBH48A7GJ31FV7AN
title: "Windows-CI: TestSpawnOffsetLock wartet 5 Minuten nach einem getöteten spawn.sh"
status: backlog
ready: false
creator: Sascha
goal: "TestSpawnOffsetLock läuft auf windows-latest zuverlässig grün, oder der Test sagt, warum er dort nicht gilt."
context: |-
  TestSpawnOffsetLock (core/role/role_test.go:1801) fällt auf windows-latest mit: spawn.sh waited 5m after a killed spawn.sh: the lock outlived it.
  Gesehen: PR #260 (574VCW) zweimal hintereinander am 25.11.2026 (run 45609911982), vorher feat/mmkena run 45345857992 am 24.11. feat/pgxjg0 war am selben Tag grün - also Flake.
  Nicht von 574VCW: dort ändert spawn.sh nur --back und den discuss-Text.
  Master-CI fährt auf push nur Linux; Windows/macOS laufen nur auf PRs, darum sieht master grün aus.
  Unerklärt: warum genau ~300 s. Weder spawn.sh noch der Test hat diesen Timeout (flock -w 600). Verdacht: auf Windows (Git-Bash) gibt ein getöteter Prozess das flock-FD nicht sofort frei, oder ein Kindprozess hält es.
definition-of-done: Ursache der ~300 s Wartezeit gefunden und als Notiz am Ticket
tags:
  - ci
blocked-by: []
related: []
commits: []
created-at: 2026-11-25T12:26:22Z
updated-at: 2026-11-25T12:26:22Z
---

# Windows-CI: TestSpawnOffsetLock wartet 5 Minuten nach einem getöteten spawn.sh

## Definition of Done

- [ ] Ursache der ~300 s Wartezeit gefunden und als Notiz am Ticket
- [ ] TestSpawnOffsetLock auf windows-latest 3 Läufe hintereinander grün

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

