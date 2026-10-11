---
id: 01M4P7CF7Q5792CM4HVCAK3Z6D
title: Tests lesen ~/.jaira der Maschine und fallen deshalb auf master durch
status: critique
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
commits: []
created-at: 2026-10-11T19:35:08Z
updated-at: 2026-10-11T19:54:19Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-790024
claimed-at: 2026-10-11T19:45:28Z
outcome-what: "Gemeinsamer Helfer internal/testhome.Run setzt HOME, USERPROFILE und JAIRA_HOME auf ein leeres Temp-Verzeichnis; alle 21 betroffenen Testpakete rufen ihn aus TestMain"
outcome-why: "Ein Helfer statt Isolierung pro Fixture: der Test, der es vergisst, ist sonst still falsch auf genau den Maschinen, auf die niemand schaut"
outcome-resolves: "go test ./... -race ist gruen unabhaengig von ~/.jaira, auch mit default-board.md ohne backlog-Lane"
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

- [ ] internal/testhome: Run(m) setzt HOME, USERPROFILE, JAIRA_HOME auf ein absolutes Temp-Verzeichnis, raeumt auf
- [x] TestMain in allen 21 Paketen, deren Tests transitiv lane/project/settings/ticket/identity/role/selfupdate laden (cli und tui erweitern)
- [x] Nachweis: Suite mit HOME auf Verzeichnis, dessen .jaira/default-board.md keine backlog-Lane hat; Exit-Code pruefen

## Progress
- **2026-10-11 19:46 · Alexander Sacharov** — Paketliste ermittelt mit: go list -test -deps <pkg> | grep core/(lane|project|settings|ticket|identity|role|selfupdate) — das sind die Pakete mit os.UserHomeDir/JAIRA_HOME. 21 Testpakete treffen; bgrun, board, gitref, hook, notify, release, snapshot, wintrap nicht. Gemeinsamer Helfer internal/testhome statt 21 Kopien desselben TestMain-Rumpfs.
- **2026-10-11 19:54 · Alexander Sacharov** — Ueber die drei DoD-Variablen hinaus: testhome.Run loescht JAIRA_DEFAULT_BOARD und JAIRA_LANES_DIR (sie zeigen sonst auf Dateien der Maschine) und pinnt vorher GOPATH/GOMODCACHE/GOCACHE/GOENV per 'go env' — internal/cli baut in mergebranches_test.go mit 'go build', das mit leerem HOME sonst einen leeren Modul-Cache haette (Netz) und go env -w verliert. JAIRA_USER bewusst nicht angefasst: kein ~/.jaira-Pfad.
