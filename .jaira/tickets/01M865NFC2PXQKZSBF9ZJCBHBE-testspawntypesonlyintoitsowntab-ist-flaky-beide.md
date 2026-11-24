---
id: 01M865NFC2PXQKZSBF9ZJCBHBE
title: "TestSpawnTypesOnlyIntoItsOwnTab ist flaky: beide Starts melden 'pane gehoert jetzt zu t2'"
status: testing
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
updated-at: 2026-11-24T07:09:29Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-91083
claimed-at: 2026-11-24T07:00:45Z
outcome-what: "Herdr-Stub in TestSpawnTypesOnlyIntoItsOwnTab vergibt den Tab nach Label statt per mkdir-Zaehler mit 10s-Grenze; Verlierer (critique) steht fest. role_test 'second while the first still runs': Stub-Grenze 10s -> 120s."
outcome-why: "Release-Lauf 0.6.6 rot: beide Starts bekamen vom Stub tab t1, Ergebnis hing am Timing des Stubs, nicht an spawn.sh."
outcome-resolves: "Race-Test deterministisch; 200x gruen, auch mit 12s Verzoegerung je Seite. Geprueft: Nebenlaeufigkeit (beide Starts parallel, Verzoegerung beider Seiten), Fehlerpfade (Timeout beeinflusst Ergebnis nicht mehr), Eingabe/Rechte: n/a (nur Test-Stub)."
review-summary: none
review-gaps: "spawnpane_test.go: Variable won (fester Wert optimize, zweimal genutzt) eingesetzt; sonst nichts entfernt - Stub hat keine Duplikate, role_test-Stub bleibt eigenstaendig (anderer Zweck: Gate statt Tab-Vergabe)"
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
