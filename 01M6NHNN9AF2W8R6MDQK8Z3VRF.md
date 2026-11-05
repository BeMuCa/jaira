---
id: 01M6NHNN9AF2W8R6MDQK8Z3VRF
title: "jaira <id> startet ein installiertes Plugin"
status: review
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
updated-at: 2026-11-05T22:06:05Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-75959
claimed-at: 2026-11-05T21:39:09Z
outcome-what: "jaira <id> startet ~/.jaira/plugins/<id>/jaira-<id> (Argumente, stdio, JAIRA_ROOT/BOARD/VERSION, Exit-Code); eingebaute Befehle gewinnen. create/set/claim/move feuern Plugin-Events aus jaira-plugin.json (move ueber core/move.Moved, also auch aus der TUI). external in show/list --json. Merge uebernimmt external von der Seite, die es geaendert hat."
outcome-why: "Standard docs/PLUGINS.md (R6X6FH) verlangt Dispatch und Events; Teamlead-Luecken (TUI-Moves, external in JSON, merge, AGENTS.md) gehoeren dazu. Der Merge meldete bisher einen Konflikt, sobald nur ein Branch external setzte."
outcome-resolves: Plugins lassen sich ohne Aenderung am Kern starten und hoeren Schreibvorgaenge; Install und Katalog bleiben ein spaeteres Ticket.
review-summary: "jaira <id> startet ~/.jaira/plugins/<id>/jaira-<id> (.exe auf Windows) mit allen Argumenten, stdio, JAIRA_ROOT/BOARD/VERSION und gibt dessen Exit-Code zurueck - nur wenn cobra den Befehl nicht kennt (help/completion/Aliase eingeschlossen), nie aus PATH. create/claim/set rufen nach dem Schreiben jedes Plugin, dessen jaira-plugin.json das Event nennt, als 'jaira-<id> event <name>' mit show --json + actor + from_lane auf stdin, 5 s Timeout, Ausgabe verworfen, JAIRA_IN_EVENT=1 verhindert Kettenreaktionen. move feuert aus core/move.Move ueber den Hook move.Moved, damit auch TUI-Moves ankommen; im Terminal-Board laufen die Events ueber eine Warteschlange mit einer Goroutine, die beim Beenden geleert wird; die CLI bleibt synchron. set board=<x>, das die Lane wechselt, feuert zusaetzlich move. external steht in show/list --json; der Merge uebernimmt external jetzt von der Seite, die es geaendert hat (Doc.Entry/SetEntry statt Konflikt). Doku AGENTS.md/PLUGINS.md und Release-Note angepasst."
round: 3
review-gaps: |-
  note: die TUI-Warteschlange hat 64 Plaetze - haengen Plugins laenger, blockiert der 65. Move im Board doch, bis einer frei wird
  note: ein per Signal beendetes Plugin liefert ExitCode -1, jaira beendet sich dann mit 255 statt mit einem eigenen Code
  note: der TUI-Pfad (Warteschlange im echten Board) ist nur per Unit-Test geprueft, nicht im laufenden Board
test-verdict: "pass (round 2): go test -race -timeout 20m ./... RC=0 after go clean -testcache, and RC=0 again after merging upstream/master 3e3e1aee; DoD 8-10 proof tests run green (TestBoardEventsDoNotWaitForThePlugin, TestSetBoardThatChangesTheLaneFiresAMove, TestProgramOnWindowsIsTheExe); DoD 1-7 checked on round 1 with the built binary on a scratch board and still green in the suite. TUI not driven by hand: queueEvents is checked through firePlugins as runBoard wires it"
review-verdict: "Diff erfuellt alle 10 DoD-Punkte; keine Defekte. go test -race core/plugin, core/move, core/merge, core/ticket gruen; internal/cli und internal/tui komplett gruen; Handtest mit Fake-Plugin: Argumente, JAIRA_BOARD aus --board=, Exit-Code 7, unbekannter Befehl ohne Plugin weiter Exit 2, Events create/claim/set/move mit actor und from_lane (backlog -> todo)."
review-round: 3
review-check: "1. cd in dieses Repo, dann: go build -o /tmp/j ./cmd/jaira  2. mkdir -p /tmp/h/plugins/demo  3. Datei /tmp/h/plugins/demo/jaira-demo anlegen: #!/bin/sh / if [ \"$1\" = event ]; then cat > /tmp/h/ev-$2.json; exit 0; fi / echo \"args=$* board=$JAIRA_BOARD\"; exit 7 - dann chmod +x  4. echo '{\"id\":\"demo\",\"events\":[\"create\",\"move\",\"set\"]}' > /tmp/h/plugins/demo/jaira-plugin.json  5. in einem Repo mit .jaira: JAIRA_HOME=/tmp/h /tmp/j demo x --board=foo; echo $? -> zeigt \"args=x --board=foo board=foo\" und 7  6. JAIRA_HOME=/tmp/h /tmp/j create T --goal g --context c --dod d -> /tmp/h/ev-create.json existiert, enthaelt actor und external  7. JAIRA_HOME=/tmp/h /tmp/j move <id> --to todo -> /tmp/h/ev-move.json hat from_lane \"backlog\", status \"todo\"  8. JAIRA_HOME=/tmp/h /tmp/j board, Ticket mit der Tastatur verschieben -> Board reagiert sofort, nach dem Beenden ist ev-move.json neu  9. go test ./core/merge -run TestMergeExternalChangedOnOneSideWins -> ok"
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
- [x] TUI-Moves blockieren die UI nicht: Events laufen in einer Warteschlange im Hintergrund, beim Beenden des Boards abgearbeitet; CLI bleibt synchron; Test
  proof: internal/cli/plugin.go:queueEvents/runBoard (board.go, home.go:openBoardAt); Test TestBoardEventsDoNotWaitForThePlugin (queued returns <1s, stop drains, CLI path synchronous)
- [x] set board=<x>, das die Lane wechselt, feuert auch 'move' (from_lane = alte Lane) mit den Lanes des neuen Boards; Test
  proof: internal/cli/tickets.go: set fires move when landed != ""; Test TestSetBoardThatChangesTheLaneFiresAMove (set then move, from_lane backlog, status_known/next_lane of the video board)
- [x] .exe-Suffix in plugin.Program durch einen Test abgedeckt, der auch auf Linux laeuft
  proof: core/plugin/plugin.go:program(id, goos); Test TestProgramOnWindowsIsTheExe

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
- **2026-11-05 21:51 · Alexander Sacharov** — critique (Runde 2): beide Befunde aus Runde 1 erledigt (b154df66) - Events geloescht, Subscribers einmal pro Schreibvorgang und leer innerhalb eines Events. Nichts Neues.
- **2026-11-05 21:58 · Alexander Sacharov** — Teamlead (47752) nach Review, vor Signoff: drei Review-Notizen werden behoben. (1) TUI-Moves feuern Events ohne die UI zu blockieren; CLI bleibt synchron wie PLUGINS.md. (2) set board=<x>, das den Status in die erste Lane des neuen Boards setzt, feuert zusaetzlich 'move' mit den Lanes des neuen Boards. (3) Test fuer das .exe-Suffix, der auch auf Linux laeuft (GOOS als Parameter). Plan: TUI -> eine Warteschlange mit einem Worker-Goroutine in internal/cli, Payload wird synchron gebaut (Zustand zum Zeitpunkt des Moves), gestartet vor tui-Lauf, beim Beenden geleert (Reihenfolge bleibt, kein Event geht beim Quit verloren; Quit wartet hoechstens auf laufende Plugins).
- **2026-11-05 22:00 · Alexander Sacharov** — Runde 2 umgesetzt. Befund zu Review-Notiz 2: der set-Payload war schon mit den Lanes des neuen Boards gebaut (envFor liest t.Board nach dem Mutate) - der Test prueft status_known/next_lane jetzt explizit; neu ist nur das zusaetzliche move-Event. Home-Bildschirm (runHome) bewegt keine Tickets, daher nur openBoardAt und 'jaira board' mit Warteschlange. Haengendes Plugin beim Quit: Quit wartet bis zu 5 s je noch laufendem Event - bewusst, statt Events zu verlieren.
- **2026-11-05 22:00 · Alexander Sacharov** — critique (Runde 3, nur 233ffbff): nichts. core/bgrun geprueft - das ist ein abgekoppelter, zeitlich gedrosselter Unterbefehl ohne Reihenfolge und ohne Warten, passt nicht fuer geordnete Events, die beim Quit nicht verloren gehen duerfen.
- **2026-11-05 22:06 · Alexander Sacharov** — review (Runde 3, frischer Kontext): keine Blocker. Geprueft gegen upstream/master...HEAD (5 Commits), alle Tests inkl. -race gruen, Handtest mit Fake-Plugin fuer dispatch und vier Events. Drei Notizen in review-gaps.
