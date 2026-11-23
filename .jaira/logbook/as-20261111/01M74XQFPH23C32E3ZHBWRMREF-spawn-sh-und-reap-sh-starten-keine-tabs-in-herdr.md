---
id: 01M74XQFPH23C32E3ZHBWRMREF
title: spawn.sh und reap.sh starten keine Tabs in Herdr App (hc)
status: done
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
commits:
  - 655874a6b9f2c7e6870dfd9ad6b82b004cf39b25
  - 8b51a187975a06b13b805fddd236ef05cd598554
  - 403ae49555ec491a62afa340622116516ed23189
  - a93866dffab04415950c4436efd0870fa35b5eb8
created-at: 2026-11-11T09:06:24Z
updated-at: 2026-11-11T10:00:39Z
assignee: AlSa
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-15715
claimed-at: 2026-11-11T09:06:40Z
outcome-what: "herdr-hc.sh übersetzt die herdr-Aufrufe von spawn.sh, run-lane.sh und reap.sh in Herdr Apps hc; die drei Skripte und die Rollen-Skills erkennen Herdr App an HC_TERMINAL_ID."
outcome-why: "Herdr App setzt kein HERDR_ENV und hat kein herdr-Binary: Worker-Tabs ließen sich nicht mehr öffnen, und Terminals in Nachbarordnern landeten ohne herdr.space-Label außerhalb des eigenen Space."
outcome-resolves: Teamlead und Dispatcher öffnen in Herdr App wieder Tabs im selben Space.
review-summary: "Neues Skript herdr-hc.sh übersetzt die herdr-Aufrufe, die spawn.sh, run-lane.sh und reap.sh machen (tab create/list/get/close, pane get/run/send-text/send-keys), in Aufrufe von Herdr Apps hc: ein hc-Terminal steht für Tab und Pane zugleich (tN). tab create startet 'hc run --bg --outlive-parent --no-notify' mit Label jaira.tab=<ticket/lane> und kopiert das herdr.space-Label des Aufrufers, damit der Tab im selben Space landet. Die drei Skripte nehmen den Shim, wenn HERDR_ENV nicht 1 ist, HC_TERMINAL_ID gesetzt ist und hc auf dem PATH liegt; kein Skript weist HERDR_ENV mehr zu (WSLENV hätte es an Kinder vererbt). reap.sh überspringt den eigenen Terminal über HC_TERMINAL_ID. Dispatcher-, Teamlead- und Worker-Skill erkennen Herdr App am HC_TERMINAL_ID. Ein Go-Test fährt spawn.sh gegen ein Fake-hc."
round: 3
review-gaps: |-
  note: nur spawn.sh ist getestet; run-lane.sh (pane get -> agent_status) und reap.sh (tab list/pane get/tab get/tab close) laufen durch den Shim ohne Test
  note: labelof liest Labels per grep aus kompaktem JSON; gibt hc je Leerzeichen nach dem Doppelpunkt aus, fehlt herdr.space still und der Tab landet im falschen Space
  note: die Skill-Erkennung prüft nur HC_TERMINAL_ID, die Skripte zusätzlich hc auf dem PATH; ohne hc wählt der Skill Herdr und spawn.sh bricht mit 'not inside a Herdr pane' ab statt auf Subagents zu fallen
  note: hc kennt als status.state u. a. working/idle; ob es je 'blocked' oder 'done' meldet, ist offen. Unbekannte Zustände führen in spawn.sh zum Abbruch (sicher), in run-lane.sh nur zum Weiterwarten
test-verdict: "pass: go test ./... grün nach HERDR_ENV-Fix; run-lane.sh-Lauf für review folgt live in Herdr App."
review-verdict: "Diff erfüllt alle vier DoD-Punkte, kein Blocker. Live belegt: dieser Review-Tab (t9) wurde über herdr-hc.sh geöffnet und trägt jaira.tab=WRMREF/review und das herdr.space-Label des Aufrufers. go test ./core/role grün. Nicht live gesehen: das Schließen durch run-lane.sh/reap.sh."
review-round: 3
review-check: "1. In Herdr App einen Teamlead-Tab öffnen, im Repo: echo $HC_TERMINAL_ID zeigt tN  2. jaira roles install  3. ~/.claude/skills/jaira-dispatcher/scripts/run-lane.sh <ID> <lane> mit einem kleinen Ticket starten  4. ein neuer Tab erscheint im selben Space, Label <ID>/<lane>, claude bekommt /jaira-role-lane getippt  5. wenn der Worker die Lane verlässt, schließt run-lane.sh den Tab ('closed tN')  6. go test ./core/role -run TestSpawnOpensItsTabThroughHerdrApp ist grün"
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
- **2026-11-11 09:33 · Alexander Sacharov** — review: kein Blocker. Live gesehen: Review-Tab t9 über herdr-hc.sh geöffnet, Labels jaira.tab und herdr.space korrekt. Nicht live gesehen: Tab-Schließen via run-lane.sh/reap.sh — das prüft der Mensch bei signoff (review-check Schritt 5).
- **2026-11-11 09:35 · Alexander Sacharov** — explain: Der Review fand keinen Blocker, aber Lücken: nur spawn.sh hat einen Test gegen ein Fake-hc, das Schließen der Tabs durch run-lane.sh/reap.sh lief nie live, und labelof liest herdr.space per grep aus kompaktem JSON.
Warum: Gibt hc das JSON mit Leerzeichen aus oder meldet es einen unbekannten status.state, landen Tabs still im falschen Space oder bleiben offen, ohne dass ein Test es merkt.
1) Zurückgeben: Tests für run-lane.sh und reap.sh über den Shim ergänzen und labelof gegen Leerzeichen im JSON robust machen [give back]
2) So annehmen: bei signoff review-check Schritt 5 (Tab schließt sich, 'closed tN') selbst in Herdr App prüfen, Rest als bekannte Grenze stehen lassen [accept]
3) Annehmen und die Lücken als eigenes Folge-Ticket festhalten [accept]
Empfehlung: 2 - der Hauptweg ist live belegt und go test grün; die offene Stelle ist genau Schritt 5, den ein Mensch in einer Minute selbst sieht, und gibt es dort Ärger, geht das Ticket dann zurück.
- **2026-11-11 09:56 · Alexander Sacharov** — decision: 2) So annehmen: bei signoff review-check Schritt 5 (Tab schließt sich, 'closed tN') selbst in Herdr App prüfen, Rest als bekannte Grenze stehen lassen
- **2026-11-11 10:00 · Alexander Sacharov** — Nach done: Cockpit prüfte selbst HERDR_ENV (register.tsx herdr()) - Accept meldete 'needs a pull request ... HERDR_ENV=1'. herdr() nimmt jetzt in Herdr App herdr-hc.sh neben spawn.sh; der Shim liefert tab list im herdr-Format mit jaira.tab-Labels. Test im Cockpit: 'Herdr App (HC_TERMINAL_ID, no HERDR_ENV)'.
