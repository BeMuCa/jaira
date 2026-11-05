---
id: 01M6NHNN9AF2W8R6MDQK8Z3VRF
title: "jaira <id> startet ein installiertes Plugin"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Befehl, den jaira nicht kennt, startet das gleichnamige installierte Plugin, mit derselben Umgebung auf Linux, macOS und Windows."
context: |-
  Heute: 'jaira jira ...' endet mit 'unknown command'. Es gibt keinen Weg, jaira zu erweitern, ohne den Kern zu aendern.
  Gewollt (Standard in R6X6FH, docs/PLUGINS.md): unbekannter Unterbefehl <id> -> ~/.jaira/plugins/<id>/jaira-<id>(.exe) starten, Argumente, stdin/stdout/stderr und Exit-Code durchreichen.
  Umgebung: JAIRA_ROOT, JAIRA_BOARD, JAIRA_VERSION.
  Eingebaute Befehle gehen immer vor; ein Plugin kann keinen Kernbefehl ueberdecken.
  Kein PATH-Suchen: nur was 'jaira plugin install' abgelegt hat, laeuft (Whitelist).
definition-of-done: "Unbekannter Unterbefehl startet ~/.jaira/plugins/<id>/jaira-<id> mit Argumenten, Umgebung und Exit-Code; Test"
tags:
  - cli
blocked-by:
  - R6X6FH
related: []
commits: []
created-at: 2026-11-05T09:47:04Z
updated-at: 2026-11-05T21:51:07Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-75959
claimed-at: 2026-11-05T21:39:09Z
outcome-what: "jaira <id> startet ~/.jaira/plugins/<id>/jaira-<id> (Argumente, stdio, JAIRA_ROOT/BOARD/VERSION, Exit-Code); eingebaute Befehle gewinnen. create/set/claim/move feuern Plugin-Events aus jaira-plugin.json (move ueber core/move.Moved, also auch aus der TUI). external in show/list --json. Merge uebernimmt external von der Seite, die es geaendert hat."
outcome-why: "Standard docs/PLUGINS.md (R6X6FH) verlangt Dispatch und Events; Teamlead-Luecken (TUI-Moves, external in JSON, merge, AGENTS.md) gehoeren dazu. Der Merge meldete bisher einen Konflikt, sobald nur ein Branch external setzte."
outcome-resolves: Plugins lassen sich ohne Aenderung am Kern starten und hoeren Schreibvorgaenge; Install und Katalog bleiben ein spaeteres Ticket.
review-summary: |-
  core/plugin/plugin.go:30 var Events has no caller (the closed list lives in the schema and the doc) - delete it
  internal/cli/plugin.go:firePlugins checks plugin.Subscribers and plugin.Fire scans the directory again and stats every Program a second time - let Subscribers return nil inside an event (JAIRA_IN_EVENT) and pass its ids into Fire(ids, event, payload, env), so one scan per write and no JSON built inside an event
---

# jaira <id> startet ein installiertes Plugin

## Definition of Done

- [x] Unbekannter Unterbefehl startet ~/.jaira/plugins/<id>/jaira-<id> mit Argumenten, Umgebung und Exit-Code; Test
  proof: internal/cli/plugin.go:runPlugin, core/plugin/plugin.go:Run; Tests TestUnknownCommandStartsTheInstalledPlugin, TestRunPassesArgumentsEnvironmentAndExitCode, TestProgramIsOnlyAnInstalledPlugin
- [x] Ein eingebauter Befehl gleichen Namens gewinnt; Test
  proof: internal/cli/plugin.go:isBuiltin; Test TestBuiltinCommandWinsOverAPluginOfTheSameName (list, help, completion)
- [x] Release-Note-Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/8Z3VRF.md
- [x] Events: create, move, claim und set aus der CLI und Moves aus der TUI rufen jedes abonnierte Plugin als 'jaira-<id> event <name>' auf (stdin = show --json + actor + from_lane, 5 s Timeout, Fehler aendert nichts, JAIRA_IN_EVENT=1 unterdrueckt Events); Test
  proof: core/move/move.go:Moved (CLI+TUI), internal/cli/plugin.go:firePlugins, core/plugin/plugin.go:Fire; Tests TestWritesFireThePluginsEvents, TestAFailingPluginDoesNotFailTheCommand, TestMovedHearsEveryWrittenMove, TestFireHandsThePluginTheEvent, TestFireInsideAnEventCallsNoPlugin, TestFireStopsAHangingPlugin
- [x] Feld external erscheint in show --json und list --json; Test
  proof: internal/cli/tickets.go:ticketJSON 'external', core/ticket/schema.go:External; Test TestExternalIsInShowAndListJSON
- [x] merge_test deckt ein Ticket ab, dessen Frontmatter external als YAML-Block traegt
  proof: core/merge/merge_test.go:TestMergeExternalChangedOnOneSideWins; Fix core/merge/merge.go:mergeOpaque + core/ticket/frontmatter.go:Entry/SetEntry
- [x] docs/AGENTS.md:8 ('no plugin API') und docs/PLUGINS.md (Status, external-Satz) beschreiben den neuen Stand und verlinken docs/PLUGINS.md
  proof: docs/AGENTS.md:8, docs/PLUGINS.md: Status, Writing (external in JSON), Events (Abo aus jaira-plugin.json, TUI-Moves)

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [-] core/ticketjson: ticketJSON samt checklistJSON/rejectionJSON/nextLaneID aus internal/cli/tickets.go herausziehen, damit die TUI denselben Payload bauen kann; external ergaenzen
- [x] core/plugin: Dir (JAIRA_HOME), Lookup(id), Env(root, board, version), Run (stdio, Exit-Code), Subscribers(event) aus jaira-plugin.json, Fire (stdin JSON, 5 s, JAIRA_IN_EVENT)
- [x] cli.Execute: ist os.Args[1] ein Befehl, den cobra nicht kennt, und ein installiertes Plugin -> plugin.Run
- [x] Events: firePlugins neben fireHook in create, claim, set, move; TUI-Moves (model.go x2, signoff.go) ueber einen Callback aus cli
- [x] Tests: Dispatch (Args, Env, Exit-Code, stdin), eingebauter Befehl gewinnt, Events (vier Namen, Timeout, Fehler, JAIRA_IN_EVENT), external in show/list JSON, merge external-Block
- [x] Doku: docs/AGENTS.md:8, docs/PLUGINS.md Status + external-Satz; Release-Note core/release/notes/8Z3VRF.md

## Progress
- **2026-11-05 21:39 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/47752.sock
- **2026-11-05 21:41 · Alexander Sacharov** — Vom Teamlead (47752) am 05.11. nachgereicht, in dieses Ticket gefaltet: AGENTS.md, Events auch aus TUI/web, external in JSON, merge_test external-Block, Install soll Manifest gegen checksums.txt pruefen. Branch ist auf upstream/master eebf52b0 (PR #124, docs/PLUGINS.md) vorgespult.
- **2026-11-05 21:41 · Alexander Sacharov** — Offene Entscheidungen fuer die Person: keine. Selbst entschieden: (1) Plugin-Verzeichnis = $JAIRA_HOME/plugins, sonst ~/.jaira/plugins - wie settings.Path. (2) Nur das erste Argument ist die Plugin-id; 'jaira -C x jira' startet kein Plugin. (3) Eingebaute Befehle inkl. Aliase und help/completion gewinnen: Dispatch nur, wenn cobra das erste Argument nicht kennt. (4) Kein Plugin installiert -> bisheriger 'unknown command'-Fehler, Exit 2. (5) JAIRA_ROOT ausserhalb eines Repos leer, Plugin laeuft trotzdem (z.B. Login). (6) JAIRA_BOARD = Wert von --board/--board= in den Plugin-Argumenten (die unveraendert weitergehen), sonst 'default'. (7) Event-Abos liest jaira aus ~/.jaira/plugins/<id>/jaira-plugin.json (das Install ablegen wird, Punkt 5 des Teamleads) - kein 'manifest --json'-Prozess pro Schreibvorgang. Fehlt/kaputt -> keine Events. (8) Web: core/web hat heute keine Route, die ein Ticket bewegt/claimt/setzt (nur Epics) - es gibt dort nichts zu feuern; wenn eine kommt, ruft sie dieselbe Funktion. (9) TUI feuert nur Plugin-Events, nicht den Nutzer-Hook - der Hook-Vertrag bleibt wie er ist. (10) Pruefung des Manifests gegen checksums.txt gehoert zu 'jaira plugin install', das es noch nicht gibt - nur notiert, kein DoD hier.
- **2026-11-05 21:45 · Alexander Sacharov** — Stand: core/plugin (Program/Run/Subscribers/Fire) und internal/cli/plugin.go (runPlugin in Execute, firePlugins) gebaut. Plan-Schritt 1 (ticketJSON nach core ziehen) verworfen: stattdessen ruft core/move.Move nach jedem geschriebenen Move den Hook move.Moved, den internal/cli in init() setzt - deckt CLI und TUI (model.go x2, signoff.go) mit einer Stelle ab, ohne 30 Aufrufstellen umzuziehen. Gefunden: core/merge mergeOpaque meldet einen Konflikt, wenn nur die andere Seite external aendert - widerspricht docs/PLUGINS.md ('whichever side changed it wins'); wird in DoD 6 behoben.
- **2026-11-05 21:51 · Alexander Sacharov** — critique (Runde 1): core/plugin/plugin.go:30 var Events has no caller (the closed list lives in the schema and the doc) - delete it
internal/cli/plugin.go:firePlugins checks plugin.Subscribers and plugin.Fire scans the directory again and stats every Program a second time - let Subscribers return nil inside an event (JAIRA_IN_EVENT) and pass its ids into Fire(ids, event, payload, env), so one scan per write and no JSON built inside an event. Stehen gelassen: JAIRA_HOME-Logik in plugin.Dir dupliziert settings.Path - das gleiche Muster steht in 7 Paketen (session, selfupdate, project, identity, lane, ticket/store), eigenes Paket ist Konvention. move.Moved als Paketvariable: gleiches Muster wie release.Current, deckt CLI+TUI an einer Stelle.
