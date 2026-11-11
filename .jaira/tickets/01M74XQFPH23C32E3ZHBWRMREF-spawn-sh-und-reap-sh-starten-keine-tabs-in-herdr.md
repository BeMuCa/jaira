---
id: 01M74XQFPH23C32E3ZHBWRMREF
title: spawn.sh und reap.sh starten keine Tabs in Herdr App (hc)
status: critique
ready: true
creator: Alexander Sacharov
goal: "spawn.sh, run-lane.sh und reap.sh öffnen, finden und schließen Worker-Tabs auch in Herdr App (hc / herdr-core), im selben Space wie der Aufrufer."
context: |-
  Was kaputt ist: In Herdr App gibt es kein HERDR_ENV, kein herdr-Binary und keine tab/pane-Befehle. spawn.sh bricht mit 'not inside a Herdr pane' ab; Teamlead und Dispatcher fallen auf Subagents zurück.
  Herdr App setzt stattdessen HC_TERMINAL_ID (z. B. t3) und HC_TERMINAL_SOCKET; das CLI heißt hc (~/.local/bin/hc).
  Ein Terminal mit --cwd in einem Nachbarordner (../.worktrees/jaira-<slug>) landet ohne Label im falschen Platz. Mit --label herdr.space=<Space des Aufrufers> legt die App es als neuen Tab im selben Space an (ausprobiert: 'jaira / 2').
  hc run beendet Kinder mit dem Elternprogramm, wenn --outlive-parent fehlt, und tippt Notizen ins Eltern-Terminal, wenn --no-notify fehlt — ein Dispatcher muss den Teamlead überleben.
  Ausgelöst: Alex ist am 2026-11-11 auf Herdr App umgestiegen.
definition-of-done: "spawn.sh startet in Herdr App (HC_TERMINAL_ID gesetzt, kein HERDR_ENV) claude in einem neuen Tab im Space des Aufrufers und tippt den Lane-Befehl"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-11-11T09:06:24Z
updated-at: 2026-11-11T09:31:00Z
assignee: AlSa
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-15715
claimed-at: 2026-11-11T09:06:40Z
outcome-what: "herdr-hc.sh übersetzt die herdr-Aufrufe von spawn.sh, run-lane.sh und reap.sh in Herdr Apps hc; die drei Skripte und die Rollen-Skills erkennen Herdr App an HC_TERMINAL_ID."
outcome-why: "Herdr App setzt kein HERDR_ENV und hat kein herdr-Binary: Worker-Tabs ließen sich nicht mehr öffnen, und Terminals in Nachbarordnern landeten ohne herdr.space-Label außerhalb des eigenen Space."
outcome-resolves: Teamlead und Dispatcher öffnen in Herdr App wieder Tabs im selben Space.
review-summary: "Keine Befunde: Test für den hc-Pfad ergänzt (core/role/role_test.go TestSpawnOpensItsTabThroughHerdrApp); Shim-Form und send-keys-No-op bleiben laut Note der ersten Runde."
round: 3
review-gaps: "none: herdr-hc.sh trägt nur die Aufrufe, die die drei Skripte machen; die drei Selektionsblöcke sind je 6 Zeilen und folgen dem Muster, dass jedes Skript sein herdr selbst bestimmt."
test-verdict: "pass: go test ./... grün; live in Herdr App: spawn.sh öffnete t7 als 'jaira / 2', claude antwortete; Demo-Tab t8 für Alex geöffnet; herdr-hc.sh tab create/get/close und pane get live gegen hc geprüft."
---

# spawn.sh und reap.sh starten keine Tabs in Herdr App (hc)

## Definition of Done

- [x] spawn.sh startet in Herdr App (HC_TERMINAL_ID gesetzt, kein HERDR_ENV) claude in einem neuen Tab im Space des Aufrufers und tippt den Lane-Befehl
  proof: spawn.sh live in Herdr App: tab t7 'jaira / 2', claude answered the discuss prompt; core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh 'tab create'
- [x] reap.sh und run-lane.sh finden und schließen diese Tabs
  proof: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh 'tab close'/'pane get'/'tab get'; reap.sh and run-lane.sh pick it via HC_TERMINAL_ID
- [x] Teamlead-, Dispatcher- und Worker-Skill erkennen Herdr App als Herdr
  proof: core/role/builtin/jaira-teamlead/SKILL.md:152, jaira-dispatcher/SKILL.md transport test, jaira-worker/SKILL.md
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/WRMREF.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-11 09:14 · Alexander Sacharov** — Weg: ein Shim herdr-hc.sh übersetzt die herdr-Aufrufe der drei Skripte in hc, statt drei Skripte doppelt zu verzweigen. Ein hc-Terminal = ein Tab = eine Pane, die Adresse tN ist beide Ids. hc send drückt Enter selbst, deshalb ist send-keys enter im Shim ein No-op. shipped.sha256 nur ergänzt: scripts/shipped-history.sh im Worktree löschte 32 alte Zeilen.
- **2026-11-11 09:23 · Alexander Sacharov** — critique: einziger Befund — kein Test für den hc-Pfad; Shim statt Verzweigung in drei Skripten bleibt (Note 1), send-keys-enter-No-op bleibt, dokumentiert.
- **2026-11-11 09:30 · Alexander Sacharov** — review-Lauf scheiterte: run-lane.sh setzte HERDR_ENV=1; WSLENV exportiert HERDR_ENV leer, also erbte spawn.sh die 1 und rief herdr. Behoben: kein Skript weist HERDR_ENV mehr zu, spawn.sh prüft den Shim selbst, reap.sh nimmt HC_TERMINAL_ID. Test prüft das.
