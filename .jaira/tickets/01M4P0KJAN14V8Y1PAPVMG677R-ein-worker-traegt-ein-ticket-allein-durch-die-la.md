---
id: 01M4P0KJAN14V8Y1PAPVMG677R
title: "Ein Worker traegt ein Ticket allein durch die Lanes, nur der Review laeuft frisch"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein einziger Agent fuehrt ein Ticket durch alle agentischen Lanes bis zur Review; eine Lane mit 'context: fresh' bekommt einen neuen Agenten ohne Vorwissen; Boards mit eigenen Lanes laufen ohne Migration weiter"
context: |-
  Heute startet jaira-dispatcher fuer JEDE Lane einen neuen Agenten (core/role/builtin/jaira-dispatcher). Das kostet pro Lane ca. 55-60k Tokens Kaltstart, auch fuer critique, die nur den Diff liest.
  A/B-Versuch am 11.10.2026 auf drei Tickets (GDB5YJ, 8WFNQY, K1W2JR), Report: https://claude.ai/artifact/7hu5TjpnChgLF6PkcuugfS
  - Arm A (ein Agent pro Lane): 17 Agenten, 1,11 Mio Tokens, 69 min.
  - Arm B (ein Agent pro Ticket, Review frisch): 6 Agenten, 450k Tokens, 55 min.
  - Blinder Richter: 2:1 fuer A, jedes Mal knapp; Code oft zeilengleich.
  - Fehler in A sassen an den Uebergaben: veralteter --proof nach optimize, --amend nach move, veraltete Kommentare, die keine Lane fixte.
  Arm B lief auf den UNVERAENDERTEN 13 Lanes: der Worker folgt next_lane und holt sich jede Lane mit 'jaira show --for-lane'. Lanes muessen also nicht migriert werden.
  Offen ist nur: woher weiss der Worker, welche Lane er NICHT selbst machen darf? Heute heisst sie review, auf fremden Boards anders.
  Vorschlag: optionales Frontmatter-Feld 'context: fresh' in der Lane-Datei (.jaira/lanes/*.md). Fehlt es ueberall, gilt die letzte agentische Lane vor der ersten Menschen-Lane als fresh. 'jaira update' bietet das Feld fuer die eingebaute review-Lane an, angepasste Lanes bleiben unberuehrt.
  Was erhalten bleiben muss: offene Entscheidungen vor dem Plan zaehlen und mode=conversational setzen, Neustart ueber 'jaira resume', notes an jeder Pause, Menschen-Lanes.
  Risiko aus dem Versuch: ein Langlauf-Agent hat einmal Tests mit relativem HOME laufen lassen, go brach ab, der FAIL-Grep sah es nicht. Deshalb bleibt der frische Reviewer Pflicht.
definition-of-done: "Lane-Frontmatter kennt 'context: fresh' (Schema, Parser, 'jaira lanes' zeigt es, 'show --for-lane --json' gibt es aus); ohne das Feld auf dem ganzen Board gilt die letzte agentische Lane vor der ersten Menschen-Lane als fresh; Tests decken beides"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-11T17:36:41Z
updated-at: 2026-10-11T17:42:04Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-521795
claimed-at: 2026-10-11T17:38:06Z
---

# Ein Worker traegt ein Ticket allein durch die Lanes, nur der Review laeuft frisch

## Definition of Done

- [ ] Lane-Frontmatter kennt 'context: fresh' (Schema, Parser, 'jaira lanes' zeigt es, 'show --for-lane --json' gibt es aus); ohne das Feld auf dem ganzen Board gilt die letzte agentische Lane vor der ersten Menschen-Lane als fresh; Tests decken beides
- [ ] Rolle jaira-worker in core/role/builtin ersetzt jaira-dispatcher: traegt das Ticket selbst durch alle Lanes bis vor die erste fresh-Lane, startet dafuer einen frischen Agenten mit jaira-role-lane, und behaelt Zaehlen offener Entscheidungen, mode=conversational, Neustart ueber resume und notes
- [ ] jaira-teamlead startet jaira-worker statt jaira-dispatcher; der alte Name zeigt einen Hinweis statt still zu verschwinden
- [ ] Zeile in core/release/NOTES.md unter Unreleased; go test ./... -race gruen (mit isoliertem HOME, siehe bekannte Basis-Fehler)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] lane.Lane.Context: parse (fresh|shared), Fehler fuer anderes, lanesEquivalent, Template
- [x] Set: Default-Regel in Load nach der Reihenfolge (ContextImplied)
- [x] Tests in core/lane fuer Feld, Fehler und Default (inkl. shared schaltet ab, Board ohne Menschen-Lane)
- [ ] CLI: jaira lanes (Tabelle+JSON), lanes show, show --for-lane --json context; Tests
- [ ] announce: fresh-Lane in der Lane-Liste markieren
- [ ] Rolle jaira-worker schreiben, Skripte verschieben, Dispatcher-Stub, role_test anpassen
- [ ] teamlead + role-lane + andere Rollen auf worker umstellen; README/docs/roles.go
- [ ] NOTES.md Zeile; go test ./... -race mit isoliertem HOME

## Progress
- **2026-10-11 17:40 · Alexander Sacharov** — Offene Entscheidungen vor dem Plan — alle selbst entschieden, keine braucht den Menschen:
1. 'erste Menschen-Lane' = erste Lane mit requires-human-exit, Lanes die nur eine Frage nehmen (requires-question, HITL) zaehlen nicht — wie lane.Next sie ueberspringt. Woertlich genommen waere HITL die erste und testing fresh; das widerspricht dem Ziel (review frisch). Auf dem Standard-Board ergibt die Regel review.
2. Werte: fresh oder shared. shared ist noetig, damit ein Board ausdruecklich KEINE frische Lane haben kann (sonst greift immer der Default). Leer = nicht gesetzt. Jeder andere Wert ist ein Parse-Fehler wie ein falscher precedence.
3. Der Default greift nur, wenn KEINE Lane des Boards context: setzt.
4. Eingebaute review.md bleibt unveraendert, 'jaira update' bietet nichts an: Default ergibt dasselbe, und ein Feld in der eingebauten Lane wuerde die review.md jedes bestehenden Boards als 'weicht ab' (Overrides) markieren. Nicht in der DoD.
5. jaira-dispatcher bleibt als Stub-SKILL.md, der auf /jaira-worker verweist; spawn.sh/run-lane.sh ziehen nach jaira-worker/scripts. spawn.sh: Lane-Name 'worker' startet /jaira-worker, 'dispatch' bleibt als Alias.
6. Conversational-Modus: die lesende critique neben der Arbeit bleibt, der Worker startet sie.
7. 'jaira lanes': neue Spalte CONTEXT; abgeleitetes fresh mit * und Fussnote. JSON: context + context_implied. show --for-lane --json: context.
8. Die generierte CLAUDE.md-Lane-Liste markiert eine fresh-Lane, damit jeder Agent es liest.
