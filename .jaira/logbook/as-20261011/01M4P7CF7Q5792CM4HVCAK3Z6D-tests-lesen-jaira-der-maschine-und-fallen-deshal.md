---
id: 01M4P7CF7Q5792CM4HVCAK3Z6D
title: Tests lesen ~/.jaira der Maschine und fallen deshalb auf master durch
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "go test ./... -race ist auf jeder Maschine gruen, egal was in ~/.jaira liegt"
context: |-
  Was falsch ist: auf Alex' Maschine fallen auf master ~16 Tests durch (core/lane, core/move, internal/cli, internal/tui), z.B. 'unrecognized lane "backlog"' und TestArchiveFromTheBoard (nil-Panic).
  Ursache, am 11.10.2026 von mehreren Agenten unabhaengig gefunden: die Tests isolieren HOME nicht und lesen ~/.jaira/default-board.md der Maschine; dort fehlt die Lane backlog.
  Beweis: mit HOME, USERPROFILE und JAIRA_HOME auf ein ABSOLUTES leeres Temp-Verzeichnis ist die ganze Suite gruen (Exit 0, 28 Pakete). Relativer Pfad bricht go mit einem GOMODCACHE-Fehler ab - Exit-Code pruefen, nicht nach FAIL greppen.
  Vorschlag aus der Critique von K1W2JR: HOME/USERPROFILE/JAIRA_HOME einmal pro Paket in TestMain setzen (internal/cli/main_test.go hat schon eins), nicht in einzelnen Fixtures wie claimFixture. Pakete ohne TestMain bekommen eins.
  Folge heute: jeder Agent muss die Basis-Fehler erst von echten unterscheiden; ein DoD 'Suite gruen' ist nicht ehrlich abhakbar.
  CI merkt es nicht, weil der Runner kein ~/.jaira hat.
definition-of-done: "Jedes Testpaket, das ~/.jaira lesen kann, setzt HOME, USERPROFILE und JAIRA_HOME in TestMain auf ein eigenes Temp-Verzeichnis"
tags:
  - ci
blocked-by: []
related: []
commits:
  - 151643cd078ae8bafaf1bff785547054b78e4b1e
created-at: 2026-10-11T19:35:08Z
updated-at: 2026-10-11T20:07:11Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-790024
claimed-at: 2026-10-11T19:45:28Z
outcome-what: "Review: Diff erfuellt beide DoD-Punkte, Suite mit -race auf praepariertem Home gruen, keine Defekte"
outcome-why: "Selbst nachgeprueft: Paketliste per go list -test -deps vollstaendig, Lauf exit 0, Tests schreiben nicht mehr ins Home"
outcome-resolves: "Tests haengen nicht mehr vom ~/.jaira der Maschine ab"
review-summary: "Neues Paket internal/testhome mit Run(m): legt ein leeres Temp-Verzeichnis an, setzt HOME, USERPROFILE und JAIRA_HOME darauf, loescht JAIRA_DEFAULT_BOARD und JAIRA_LANES_DIR, pinnt vorher GOPATH/GOMODCACHE/GOCACHE/GOENV per 'go env' auf die echten Werte (sonst baut internal/cli mit leerem Modul-Cache), fuehrt die Tests aus und raeumt auf. 18 core-Pakete bekommen ein einzeiliges TestMain, internal/cli und internal/tui ersetzen ihr eigenes durch denselben Aufruf. Damit liest kein Test mehr ~/.jaira/default-board.md der Maschine — die Ursache der ~16 Fehlschlaege."
review-gaps: "Nichts Blockierendes. Geprueft: (1) Alle 7 Produktionsstellen mit os.UserHomeDir/JAIRA_HOME/JAIRA_LANES_DIR/JAIRA_DEFAULT_BOARD liegen in lane, project, settings, ticket, identity, role, selfupdate; die 8 Testpakete ohne TestMain (bgrun, board, gitref, hook, notify, release, snapshot, wintrap) haben laut go list -test -deps keines davon als Abhaengigkeit — Paketliste vollstaendig. (2) go test -race -count=1 ./... auf dieser Maschine, deren ~/.jaira/default-board.md keine backlog-Lane hat: exit 0, alle Pakete ok. (3) Mit HOME auf ein Fake-Home mit derselben Datei: exit 0, und nach dem Lauf liegt unter Fake-Home/.jaira nichts ausser default-board.md — die Tests schreiben nicht mehr ins Home. Kleinigkeiten: nichts erzwingt, dass ein kuenftiges Paket mit Home-Zugriff ein TestMain bekommt (bewusst so entschieden, Kommentar in testhome nennt die Regel); JAIRA_USER wird nicht geleert, ist aber kein ~/.jaira-Pfad (in den Notizen begruendet); Plan-Punkt 1 ist nicht abgehakt, obwohl testhome.go existiert — nur Buchhaltung. Kein NOTES.md-Eintrag noetig: nur Tests, nichts fuer Nutzer sichtbar."
test-verdict: "pass: go build ./... RC=0; go test -race -count=1 ./... RC=0, 28 ok, mit HOME auf Temp mit default-board.md ohne backlog UND JAIRA_DEFAULT_BOARD/JAIRA_LANES_DIR/JAIRA_HOME feindlich gesetzt; vorher (ohne Fix) lane, move, cli, tui FAIL; DoD 1 per go list -test -deps-Schleife: kein betroffenes Paket ohne testhome.Run"
review-verdict: "Der Diff erfuellt beide DoD-Punkte; keine Defekte gefunden. Selbst nachgeprueft, nicht nur der Darstellung geglaubt."
review-check: "1. cd /home/alex/projects/.worktrees/jaira-ak3z6d  2. head -3 ~/.jaira/default-board.md — die lanes:-Zeile enthaelt kein backlog  3. go test -race -count=1 ./... ; echo exit $?  4. erwartet: jede Zeile 'ok' oder '[no test files]', kein FAIL, letzte Zeile 'exit 0' (ca. 3 Minuten, internal/tui am laengsten)  5. Gegenprobe ohne den Fix: cd /home/alex/projects/jaira (master) und go test -count=1 ./core/lane  6. erwartet dort: FAIL, u.a. '--- FAIL: TestLoadEmptyProjectDirIsSetUp'"
---

# Tests lesen ~/.jaira der Maschine und fallen deshalb auf master durch

## Definition of Done

- [x] Jedes Testpaket, das ~/.jaira lesen kann, setzt HOME, USERPROFILE und JAIRA_HOME in TestMain auf ein eigenes Temp-Verzeichnis
  proof: internal/testhome/testhome.go:33 Run; TestMain in core/{gate,gitrepo,identity,lane,link,market,merge,milestone,move,outbox,project,refsync,role,selfupdate,settings,tag,ticket,validate}/main_test.go, internal/cli/main_test.go, internal/tui/main_test.go — alle 21 Pakete, deren Tests transitiv os.UserHomeDir/JAIRA_HOME erreichen (go list -test -deps)
- [x] go test ./... -race ist gruen mit einem ~/.jaira/default-board.md ohne backlog-Lane (Nachweis: Lauf mit so praeparierter Datei)
  proof: HOME=<tmp mit .jaira/default-board.md, lanes ohne backlog> go test -race ./... -> exit 0, 28 ok; vorher ohne -race mit derselben Datei: lane, move, cli, tui FAIL

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] internal/testhome: Run(m) setzt HOME, USERPROFILE, JAIRA_HOME auf ein absolutes Temp-Verzeichnis, raeumt auf
- [x] TestMain in allen 21 Paketen, deren Tests transitiv lane/project/settings/ticket/identity/role/selfupdate laden (cli und tui erweitern)
- [x] Nachweis: Suite mit HOME auf Verzeichnis, dessen .jaira/default-board.md keine backlog-Lane hat; Exit-Code pruefen

## Progress
- **2026-10-11 19:46 · Alexander Sacharov** — Paketliste ermittelt mit: go list -test -deps <pkg> | grep core/(lane|project|settings|ticket|identity|role|selfupdate) — das sind die Pakete mit os.UserHomeDir/JAIRA_HOME. 21 Testpakete treffen; bgrun, board, gitref, hook, notify, release, snapshot, wintrap nicht. Gemeinsamer Helfer internal/testhome statt 21 Kopien desselben TestMain-Rumpfs.
- **2026-10-11 19:54 · Alexander Sacharov** — Ueber die drei DoD-Variablen hinaus: testhome.Run loescht JAIRA_DEFAULT_BOARD und JAIRA_LANES_DIR (sie zeigen sonst auf Dateien der Maschine) und pinnt vorher GOPATH/GOMODCACHE/GOCACHE/GOENV per 'go env' — internal/cli baut in mergebranches_test.go mit 'go build', das mit leerem HOME sonst einen leeren Modul-Cache haette (Netz) und go env -w verliert. JAIRA_USER bewusst nicht angefasst: kein ~/.jaira-Pfad.
- **2026-10-11 19:54 · Alexander Sacharov** — parent session: jaira-61 (vom Teamlead nachgereicht, Worker lief ohne --parent)
- **2026-10-11 19:54 · Alexander Sacharov** — critique: nichts. Gelassen: 21 TestMain-Dateien (DoD 1 verlangt jedes Paket; Rumpf ist eine Zeile, Logik in internal/testhome); pinGoVars per go-Subprozess in jedem Paket (~10ms je Paket, gebraucht von internal/cli/mergebranches_test.go go build; zentral statt nur in cli, damit ein kuenftiger go-build-Test nicht still das Netz braucht).
- **2026-10-11 20:03 · Alexander Sacharov** — review: Diff gegen DoD geprueft, go test -race ./... auf dieser Maschine (default-board ohne backlog) exit 0; Fake-Home bleibt nach dem Lauf leer bis auf default-board.md. Gegenprobe auf master: core/lane FAIL (TestLoadEmptyProjectDirIsSetUp). Neue state/r1-*-Verzeichnisse in ~/.jaira waehrend des Laufs stammen von einer parallelen Session, nicht von diesen Tests (Fake-Home-Lauf belegt das).
