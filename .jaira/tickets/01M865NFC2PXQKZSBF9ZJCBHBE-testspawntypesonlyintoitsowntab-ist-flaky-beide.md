---
id: 01M865NFC2PXQKZSBF9ZJCBHBE
title: "TestSpawnTypesOnlyIntoItsOwnTab ist flaky: beide Starts melden 'pane gehoert jetzt zu t2'"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Der Race-Test aus 3YRPXJ ist deterministisch: genau der Start, dessen Pane wirklich weggenommen wurde, bricht ab; der andere laeuft durch."
context: |-
  Was falsch ist: core/role/spawnpane_test.go:60 (TestSpawnTypesOnlyIntoItsOwnTab, aus 3YRPXJ) schlug am 24.11.2026 im release-notes-Lauf fuer 0.6.6 fehl (Actions-Run 45344252903, ubuntu): BEIDE spawns meldeten 'pane p2 is not in tab t1 (AB12CD/critique bzw. optimize) any more but in t2', erwartet war genau einer.
  Auf den PR-CI-Laeufen (#242, #244) war er gruen - also timing-abhaengig.
  Verdacht: der Herdr-Stub im Test gibt 'pane get' fuer p2 an beide Starts mit tab t2 zurueck, oder beide bekommen p2 von 'tab create'. Entweder ist der Stub nicht atomar oder spawn.sh vergleicht gegen die falsche Tab-Id.
  Zuerst klaeren, ob es der Stub oder spawn.sh ist; ist es spawn.sh, ist es ein echter Fehler aus 3YRPXJ (ein Start, der laufen duerfte, bricht ab).
definition-of-done: Ursache benannt (Stub oder spawn.sh) als Notiz am Ticket
tags:
  - ci
blocked-by: []
related: []
commits: []
created-at: 2026-11-24T07:00:03Z
updated-at: 2026-11-24T07:14:20Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-91083
claimed-at: 2026-11-24T07:00:45Z
outcome-what: "Herdr-Stub in TestSpawnTypesOnlyIntoItsOwnTab vergibt den Tab nach Label statt per mkdir-Zaehler mit 10s-Grenze; Verlierer (critique) steht fest. role_test 'second while the first still runs': Stub-Grenze 10s -> 120s."
outcome-why: "Release-Lauf 0.6.6 rot: beide Starts bekamen vom Stub tab t1, Ergebnis hing am Timing des Stubs, nicht an spawn.sh."
outcome-resolves: "Race-Test deterministisch; 200x gruen, auch mit 12s Verzoegerung je Seite. Geprueft: Nebenlaeufigkeit (beide Starts parallel, Verzoegerung beider Seiten), Fehlerpfade (Timeout beeinflusst Ergebnis nicht mehr), Eingabe/Rechte: n/a (nur Test-Stub)."
review-summary: "Nur Tests geaendert, spawn.sh unveraendert. spawnpane_test.go: der Herdr-Stub vergibt den Tab nicht mehr per mkdir-Zaehler mit 10s-Warten, sondern nach dem --label: AB12CD/critique bekommt t1, optimize t2, beide bekommen Pane p2 (liegt in t2). Damit verliert immer die critique - unabhaengig davon, wer den Stub zuerst erreicht oder ob das Warten auf t2 ablaeuft -, und der Test prueft genau diese Richtung (critique abgelehnt mit 'pane p2 is not in tab t1 (AB12CD/critique)', t1 geschlossen, nur p2 gestartet und mit lane=optimize vermerkt). role_test.go 'second while the first still runs': die Warte-Grenze im tab-create-Stub von 10s auf 120s, damit ein langsamer Runner den ersten Start nicht vorzeitig durchlaesst. Lokal: beide Tests 30x gruen, go vet sauber."
review-gaps: |-
  note: edge: fix-now: role_test.go:1401/1406 - bricht der Test vor os.WriteFile(gate) mit t.Fatal ab, haengt der erste spawn.sh jetzt 120s statt 10s im Stub; Fix: das Gate per t.Cleanup(func(){ os.WriteFile(gate, nil, 0o644) }) direkt nach dem Start der Goroutine schreiben
  note: limit: spawnpane_test.go:30 - startet optimize nie seinen Tab, wartet die critique bis 120s, bevor der Test mit klarer Meldung rot wird; nur auf einem ohnehin roten Pfad
  note: elsewhere: role_test.go:1394 und :1677 - testseitige 10s-Fristen derselben Klasse, vom Hunt schon notiert, nicht Gegenstand dieses Tickets
test-verdict: "pass: go vet ./... RC=0, go test ./... RC=0, go test -race ./core/role RC=0, check-ticket-commits RC=0; DoD 1-3 im Baum geprueft; Test 200x gruen, auch mit 12s Verzoegerung je Seite; Mutation (pane_is_ours ohne Pruefung) macht den Test rot (RC=1, beide Starts durch); schreibt nur Test-Stubs, kein Produktionscode"
hunt-verdict: "pass: bughunt done at 76e949aa4, 0 change, 0 pre-existing"
review-verdict: "accept: DoD 1-3 erfuellt; der Ausgang des Race-Tests haengt nur noch am Label, nicht an Zeit; ein kleiner fix-now in role_test.go"
review-round: 1
review-verdict-gaps: |-
  note: edge: fix-now: role_test.go:1401/1406 - bricht der Test vor os.WriteFile(gate) mit t.Fatal ab, haengt der erste spawn.sh jetzt 120s statt 10s im Stub; Fix: das Gate per t.Cleanup(func(){ os.WriteFile(gate, nil, 0o644) }) direkt nach dem Start der Goroutine schreiben
  note: limit: spawnpane_test.go:30 - startet optimize nie seinen Tab, wartet die critique bis 120s, bevor der Test mit klarer Meldung rot wird; nur auf einem ohnehin roten Pfad
  note: elsewhere: role_test.go:1394 und :1677 - testseitige 10s-Fristen derselben Klasse, vom Hunt schon notiert, nicht Gegenstand dieses Tickets
review-check: "1. cd core  2. go test ./role -run TestSpawnTypesOnlyIntoItsOwnTab -count=200  3. erwartet: ok nach ca. 15s  4. go test ./role -run TestSpawnTakesUpTheBranchTheWorkRidesOn -count=20  5. erwartet: ok  6. git diff master...feat/jcbhbe -- core/role/builtin  7. erwartet: leer (spawn.sh unveraendert, daher keine Release-Note)"
round: 2
back-from: review
back-by: lane
back-count: 1
stop-after-fix: "true"
---

# TestSpawnTypesOnlyIntoItsOwnTab ist flaky: beide Starts melden 'pane gehoert jetzt zu t2'

## Definition of Done

- [x] Ursache benannt (Stub oder spawn.sh) als Notiz am Ticket
  proof: note 'cause:' am Ticket; spawn.sh:458/469
- [x] Test 200x hintereinander gruen (go test ./core/role -run TestSpawnTypesOnlyIntoItsOwnTab -count=200)
  proof: go test ./core/role -run TestSpawnTypesOnlyIntoItsOwnTab -count=200 ok (13.5s); core/role/spawnpane_test.go:16
- [x] Release-Note in core/release/notes/<ID>.md, falls spawn.sh selbst geaendert wird
  proof: spawn.sh nicht geaendert (git diff master -- core/role/builtin leer): keine Release-Note noetig

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-24 07:00 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/3417.sock
- **2026-11-24 07:08 · Alexander Sacharov** — cause: der Stub, nicht spawn.sh. spawn.sh vergleicht den Tab von 'pane get' mit der Tab-Id aus seinem eigenen 'tab create' (spawn.sh:458/469); im CI-Lauf 45344252903 meldeten BEIDE Starts 'tab t1' - der Stub hat also beiden t1 gegeben, und auf diese Eingabe war das Ablehnen beider richtig. Testdauer 10.75s = die 10s-Grenze der Warteschleife im Stub lief ab. Wie zweimal n=1 aus dem mkdir-Zaehler kam, liess sich nicht nachstellen (480 Laeufe unter Last, kuenstliche 11s-Verzoegerung je Seite: alle gruen); darum ist der Zaehler weg statt geflickt: der Tab haengt jetzt am Label (critique=t1, optimize=t2), der Verlierer steht fest, das Warten beeinflusst das Ergebnis nicht mehr.
- **2026-11-24 07:08 · Alexander Sacharov** — sibling: role_test.go 'second while the first still runs' hatte dieselbe 10s-Grenze im tab-create-Stub - laeuft sie ab, kommt der zweite Start durch. Grenze auf 120s (die Gate-Datei schreibt der Test immer, gewartet wird also nie voll). spawn.sh unveraendert, daher keine Release-Note (DoD 3 entfaellt).
- **2026-11-24 07:11 · Alexander Sacharov** — hunt report: 0 findings, nothing filed | read whole: core/role/spawnpane_test.go, plus the spawn subtests of core/role/role_test.go (TestSpawnTakesUpTheBranchTheWorkRidesOn) and spawn.sh's tab create/pane check path they drive; first hunt, no diff-only reads | not confirmed: the test-side 10s deadlines role_test.go:1394 (same subtest whose stub bound went to 120s; blame fc7a4d4, QTQW15) and role_test.go:1677 (5ba3fa8) are the same fixed-wall-clock class and could still fail on a slow runner, but 80/80 runs under 8 parallel processes passed (ran: role.test -test.run 'TestSpawnTypesOnlyIntoItsOwnTab|TestSpawnTakesUpTheBranchTheWorkRidesOn' -test.count=10 x8), so it stays here as unconfirmed
- **2026-11-24 07:11 · Alexander Sacharov** — bughunt: done at 76e949aa4 - 0 findings in 2 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/role/role_test.go core/role/spawnpane_test.go
- **2026-11-24 07:11 · Alexander Sacharov** — waiting: bughunt (background agent) since 08:11
- **2026-11-24 07:11 · Alexander Sacharov** — waiting: review (fresh agent via run-lane.sh) since 08:11
- **2026-11-24 07:13 · Alexander Sacharov** — review: done at 76e949aa4
- **2026-11-24 07:13 · Alexander Sacharov** — review: fix-now: role_test.go 'second while the first still runs' - Gate per t.Cleanup schreiben, damit ein frueher t.Fatal den ersten spawn.sh nicht 120s im Stub haengen laesst
- **2026-11-24 07:14 · Alexander Sacharov** — fix-now aus review geschlossen: role_test.go Gate per t.Cleanup geoeffnet. stop-after-fix=true gesetzt (Fix 3 Zeilen, kein Lock/Parser).
