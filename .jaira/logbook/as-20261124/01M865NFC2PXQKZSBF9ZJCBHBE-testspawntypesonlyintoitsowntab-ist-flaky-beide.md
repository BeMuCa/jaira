---
id: 01M865NFC2PXQKZSBF9ZJCBHBE
title: "TestSpawnTypesOnlyIntoItsOwnTab ist flaky: beide Starts melden 'pane gehoert jetzt zu t2'"
status: done
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
commits:
  - a3396d44eff10051410343b8813043cba9fe451e
  - 76e949aa4ba358875239f0d9fc46b22aa0c869aa
  - c36da95f3de1aad926cfababed6be96d8ed51763
  - 05d337c1a92fb0c4f7f4cc09a4424b98caf7445f
created-at: 2026-11-24T07:00:03Z
updated-at: 2026-11-24T07:21:07Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-91083
claimed-at: 2026-11-24T07:00:45Z
outcome-what: "Herdr-Stub in TestSpawnTypesOnlyIntoItsOwnTab vergibt den Tab nach Label statt per mkdir-Zaehler mit 10s-Grenze; Verlierer (critique) steht fest. role_test 'second while the first still runs': Stub-Grenze 10s -> 120s."
outcome-why: "Release-Lauf 0.6.6 rot: beide Starts bekamen vom Stub tab t1, Ergebnis hing am Timing des Stubs, nicht an spawn.sh."
outcome-resolves: "Race-Test deterministisch; 200x gruen, auch mit 12s Verzoegerung je Seite. Geprueft: Nebenlaeufigkeit (beide Starts parallel, Verzoegerung beider Seiten), Fehlerpfade (Timeout beeinflusst Ergebnis nicht mehr), Eingabe/Rechte: n/a (nur Test-Stub)."
review-summary: "Delta review from 76e949aa4 (one commit, c36da95f3, 3 lines). The previous review's fix-now is closed: in role_test.go 'second while the first still runs' a t.Cleanup now writes the gate file right after the first spawn.sh starts, so a t.Fatal before the gate (worktree missing, tab never reached) no longer leaves the first spawn's Herdr stub waiting out its new 120s bound. Earlier rounds, already judged: the race stub in spawnpane_test.go fixes the losing tab by label (critique=t1, optimize=t2, both handed p2) instead of an mkdir counter plus a 10s wait, so exactly the critique is refused and the optimize starts, whatever the timing."
review-gaps: |-
  note: edge: role_test.go:1395 the cleanup opens the gate but does not wait for the first spawn.sh (<-first); on an already failing run t.TempDir removal can race the still-running stub and add a second error - drain first in the same cleanup
  note: elsewhere: role_test.go:1394 and :1677 keep fixed 10s test-side deadlines, the same wall-clock class (bughunt left them unconfirmed; QTQW15/5ba3fa8)
test-verdict: "pass: round 2 delta only - go test -race -count=3 TestSpawnTakesUpTheBranchTheWorkRidesOn+TestSpawnTypesOnlyIntoItsOwnTab RC=0, go vet ./core/role RC=0; fix-now (t.Cleanup oeffnet Gate) in role_test.go vorhanden; schreibt nur Test-Code"
hunt-verdict: "pass: stop-after-fix, no hunt in round 2"
review-verdict: "accept: the fix-now finding is answered; DoD 1-3 hold; no blocker in the delta"
review-round: 2
review-verdict-gaps: |-
  note: edge: role_test.go:1395 the cleanup opens the gate but does not wait for the first spawn.sh (<-first); on an already failing run t.TempDir removal can race the still-running stub and add a second error - drain first in the same cleanup
  note: elsewhere: role_test.go:1394 and :1677 keep fixed 10s test-side deadlines, the same wall-clock class (bughunt left them unconfirmed; QTQW15/5ba3fa8)
review-check: "1. cd into the jcbhbe worktree  2. go test ./core/role -run TestSpawnTypesOnlyIntoItsOwnTab -count=200  3. you see 'ok github.com/BeMuCa/jaira/core/role' (about 15s)  4. go test ./core/role -run TestSpawnTakesUpTheBranchTheWorkRidesOn -count=5  5. you see 'ok' again  6. git diff master -- core/role/builtin prints nothing: spawn.sh is unchanged, so no release note is due"
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
- **2026-11-24 07:15 · Alexander Sacharov** — waiting: review round 2 (new fresh agent via run-lane.sh) since 08:15
- **2026-11-24 07:16 · Alexander Sacharov** — review: done at c36da95f3
- **2026-11-24 07:16 · Alexander Sacharov** — report: Race-Stub vergibt Tab nach Label (critique verliert fest), role_test-Gate 120s + t.Cleanup; spawn.sh unveraendert | signoff, context unknown | 2 notes fuer die Person (edge: Cleanup wartet nicht auf <-first; elsewhere: 10s-Fristen role_test.go:1394/1677), Branch gepusht, PR offen
- **2026-11-24 07:18 · Alexander Sacharov** — Accepted from signoff: Folge-Ticket:
- HTSWK6: note: elsewhere: role_test.go:1394 and :1677 keep fixed 10s test-side deadlines, the same wall-clock class (bughunt left them unconfirmed; QTQW15/5ba3fa8)
angenommen:
- note: edge: role_test.go:1395 the cleanup opens the gate but does not wait for the first spawn.sh (<-first); on an already failing run t.TempDir removal can race the still-running stub and add a second error - drain first in the same cleanup
- **2026-11-24 07:18 · Alexander Sacharov** — cockpit: accepted by the person, signoff -> done
