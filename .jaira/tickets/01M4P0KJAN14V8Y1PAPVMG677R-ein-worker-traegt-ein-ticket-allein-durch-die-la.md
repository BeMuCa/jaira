---
id: 01M4P0KJAN14V8Y1PAPVMG677R
title: "Ein Worker traegt ein Ticket allein durch die Lanes, nur der Review laeuft frisch"
status: testing
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
updated-at: 2026-10-11T17:57:39Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-521795
claimed-at: 2026-10-11T17:38:06Z
outcome-what: "Lane-Feld context: fresh|shared mit Board-Default, sichtbar in lanes/lanes show/show --for-lane/Agenten-Block; Opt-in-Rolle jaira-worker, Schalter dispatch (settings.json, TUI, roles list), spawn.sh-Lane worker, Teamlead waehlt nach Schalter; critique-Runde-1-Funde behoben"
outcome-why: "Ein Agent pro Lane kostete im A/B-Versuch das 2,5-fache an Tokens bei gleichem Code; der frische Review bleibt Pflicht"
outcome-resolves: "Worker-Modus ohne Lane-Migration, standardmaessig aus"
review-summary: none
review-gaps: "laneFacts gab es zweimal (internal/cli/tickets.go und internal/tui/browse.go); die TUI-Kopie kannte Fresh nicht, sodass ein im TUI neu geschriebener Agenten-Block die fresh-Markierung verlor — beide durch lane.Set.Facts in core/lane/facts.go ersetzt. Lane.Fresh() entfernt (Einzeiler, ein Aufrufer). Gelassen: laneContext und impliedContextNote (zwei Aufrufer bzw. echte Logik), Settings.DispatchRole (zwei Aufrufer in roles.go, getestet)."
---

# Ein Worker traegt ein Ticket allein durch die Lanes, nur der Review laeuft frisch

## Definition of Done

- [x] Lane-Frontmatter kennt 'context: fresh' (Schema, Parser, 'jaira lanes' zeigt es, 'show --for-lane --json' gibt es aus); ohne das Feld auf dem ganzen Board gilt die letzte agentische Lane vor der ersten Menschen-Lane als fresh; Tests decken beides
  proof: core/lane/lane.go parse+impliedFresh; internal/cli/tickets.go lanes CONTEXT column, flow.go for-lane context; tests core/lane/context_test.go, internal/cli/context_test.go TestLanesAndForLaneReportContext, core/board/lanesection_test.go TestLaneSectionMarksAFreshLane
- [x] Rolle jaira-worker in core/role/builtin als Alternative zu jaira-dispatcher (der unveraendert bleibt): traegt das Ticket selbst durch alle Lanes bis vor die erste fresh-Lane, startet dafuer einen frischen Agenten mit jaira-role-lane, und behaelt Zaehlen offener Entscheidungen, mode=conversational, Neustart ueber resume und notes
  proof: core/role/builtin/jaira-worker/SKILL.md (Before the plan, conversational mode, The fresh lane, restartable); TestBuiltinsAreTheShippedRoles, TestWorkerNamesScriptsThatShip
- [x] jaira-teamlead liest die Einstellung und startet jaira-worker oder jaira-dispatcher; ohne Einstellung bleibt alles wie heute (dispatcher, ein Agent pro Lane)
  proof: core/role/builtin/jaira-teamlead/SKILL.md 'Which one is the person's choice' reads jaira --json roles list; internal/cli/roles.go; TestRolesListReportsTheDispatchSetting
- [x] Zeile in core/release/NOTES.md unter Unreleased; go test ./... -race gruen (mit isoliertem HOME, siehe bekannte Basis-Fehler)
  proof: core/release/NOTES.md Unreleased (2 Zeilen); go test ./... -race -count=1 mit HOME/USERPROFILE/JAIRA_HOME = absolutes leeres Temp-Verzeichnis: exit 0
- [x] Opt-in-Schalter in core/settings: ~/.jaira/settings.json "dispatch": "lanes" | "worker", Default "lanes"; im TUI-Einstellungsbildschirm (internal/tui/settings.go) umschaltbar; ueber ein bestehendes CLI-Kommando lesbar (jaira roles list); NOTES.md nennt die Einstellung und wie man sie einschaltet
  proof: core/settings/settings.go Dispatch/DispatchMode/DispatchRole (TestDispatchDefaultsToLanes); internal/tui/settings.go 'Agents per ticket' (TestSettingsEnterSwitchesDispatch); jaira roles list dispatch/dispatch_role; NOTES.md Unreleased line 2
- [x] Herdr bleibt: im Worker-Modus bekommt der Worker einen eigenen Herdr-Tab (Teamlead via spawn.sh, Lane-Name 'worker'), der frische Agent der fresh-Lane ebenfalls einen eigenen Tab via spawn.sh/run-lane.sh, nach seiner Lane geschlossen; ohne HERDR_ENV=1 gilt der bisherige Fallback
  proof: jaira-dispatcher/scripts/spawn.sh lane 'worker' -> /jaira-worker (TestSpawnTypesTheRoleForItsLaneName, Herdr-Stub); jaira-worker/SKILL.md 'Starting another agent' uses run-lane.sh/spawn.sh, Fallback Agent()

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] lane.Lane.Context: parse (fresh|shared), Fehler fuer anderes, lanesEquivalent, Template
- [x] Set: Default-Regel in Load nach der Reihenfolge (ContextImplied)
- [x] Tests in core/lane fuer Feld, Fehler und Default (inkl. shared schaltet ab, Board ohne Menschen-Lane)
- [x] CLI: jaira lanes (Tabelle+JSON), lanes show, show --for-lane --json context; Tests
- [x] announce: fresh-Lane in der Lane-Liste markieren
- [x] Rolle jaira-worker schreiben (nutzt die Skripte von jaira-dispatcher), spawn.sh-Lane 'worker', role_test anpassen
- [x] teamlead liest Schalter; role-lane-Wortlaut; README/docs/roles.go
- [x] NOTES.md Zeile; go test ./... -race mit isoliertem HOME
- [x] Schalter dispatch in core/settings + TUI-Eintrag + 'jaira roles list' (Text+JSON), Tests

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
- **2026-10-11 17:43 · Alexander Sacharov** — Teil 1 (Lane-Feld) committed in 8f1707b. Naechstes: Rolle jaira-worker.
- **2026-10-11 17:45 · Alexander Sacharov** — Scope-Erweiterung von Alex am 2026-10-11 (ueber den Koordinator): (1) Worker-Modus ist opt-in, Default bleibt der Dispatcher mit einem Agenten pro Lane; Schalter 'dispatch' in ~/.jaira/settings.json, im TUI-Settings-Screen, per CLI lesbar; Teamlead startet das Passende. (2) Herdr bleibt: Worker und frischer Review-Agent bekommen je einen eigenen Tab via spawn.sh. DoD 2 und 3 entsprechend umformuliert (vorher: worker ERSETZT dispatcher, alter Name nur noch Hinweis) — das widersprach der Erweiterung. Entscheidung: eigene Rolle jaira-worker statt Zweig im Dispatcher — der Dispatcher-Prompt bleibt byte-gleich, also aendert sich fuer niemanden ohne Schalter etwas, und zwei Arbeitsweisen in einer Datei wuerden jeden Absatz an eine Bedingung haengen. Der Worker benutzt die Skripte von jaira-dispatcher/scripts (gleiches Skills-Verzeichnis, roles install schreibt beide). CLI-Leseweg: 'jaira roles list' (--json: Schluessel dispatch und dispatch_role) — die Stelle, an der der Teamlead ohnehin die Rollen sieht.
- **2026-10-11 17:52 · Alexander Sacharov** — critique Runde 1: internal/cli/tickets.go:1427 impliedContextNote/laneContext stehen zwischen dem Doc-Kommentar von laneFacts und laneFacts — der Kommentar haengt jetzt an der Konstante; Block ueber den Kommentar ziehen
internal/tui/settings.go switchDispatch gibt settingsActionSaved zurueck und loest m.reload() aus, dessen Kommentar (Logbuch-Fenster) nicht passt und das fuer dispatch nichts tut; settingsActionNone zurueckgeben
core/lane/lane.go impliedFresh-Kommentar: eine Zeile ueber 150 Zeichen; umbrechen. internal/cli/roles.go Long: eine ueberlange Zeile; umbrechen. core/board/announce.go: verirrte Leerzeile nach dem tier-Block entfernen
core/role/builtin/jaira-role-lane/SKILL.md: 'A worker starts you for a lane marked fresh' klingt wie der einzige Fall — der Worker startet auch die lesende critique; 'When a worker starts you for ...' schreiben
- **2026-10-11 17:54 · Alexander Sacharov** — critique Runde 2: alle fuenf Funde aus Runde 1 in 5b01536 behoben (laneFacts-Kommentar, settingsActionNone, zwei Umbrueche, Leerzeile, role-lane-Wortlaut); die Reparatur fuehrt nichts Neues ein. Nichts mehr zu sagen.
- **2026-10-11 17:57 · Alexander Sacharov** — optimize: laneFacts gab es zweimal (internal/cli/tickets.go und internal/tui/browse.go); die TUI-Kopie kannte Fresh nicht, sodass ein im TUI neu geschriebener Agenten-Block die fresh-Markierung verlor — beide durch lane.Set.Facts in core/lane/facts.go ersetzt. Lane.Fresh() entfernt (Einzeiler, ein Aufrufer). Gelassen: laneContext und impliedContextNote (zwei Aufrufer bzw. echte Logik), Settings.DispatchRole (zwei Aufrufer in roles.go, getestet). Suite mit isoliertem HOME gruen.
