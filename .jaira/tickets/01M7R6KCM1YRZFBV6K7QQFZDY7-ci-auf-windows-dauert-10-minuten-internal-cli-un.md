---
id: 01M7R6KCM1YRZFBV6K7QQFZDY7
title: "CI auf Windows dauert 10 Minuten: internal/cli und internal/tui in Shards aufteilen, -race nur auf Linux"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein PR bekommt sein Ergebnis auf allen drei Plattformen in rund 4 statt 10 Minuten, ohne dass ein Test wegfaellt."
context: |-
  Was heute falsch ist: der CI-Job test (windows-latest) braucht ~10,5 min (Lauf 44391790958 am 2026-11-18), macOS ~5, Linux ~2,5. Alex wartet bei jedem PR auf Windows.
  Engpass sind zwei Pakete auf Windows: internal/cli 474 s, internal/tui 252 s, alles andere unter 95 s. Pakete laufen schon parallel, Tests innerhalb von internal/cli nicht (globaler Zustand, siehe Kommentar in .github/workflows/ci.yaml).
  -race verlangsamt jeden Lauf um ein Mehrfaches; Datenrennen sind plattformunabhaengig und werden auf Linux gefunden.
  fail-fast war an: faellt macOS, wird Windows abgebrochen und man sieht dessen Ergebnis nicht (zweimal bei PR #139).
  Loesung: Windows in 3 Jobs aufteilen, jeder nimmt jeden dritten Test aus internal/cli und internal/tui (go test -list), Job 0 dazu alle anderen Pakete; -race nur auf ubuntu; fail-fast: false.
definition-of-done: "test (windows-latest) laeuft in 3 Shards, die zusammen jeden Test aus internal/cli und internal/tui genau einmal laufen lassen und alle anderen Pakete einmal; Beleg im CI-Log"
tags:
  - ci
blocked-by: []
related: []
commits: []
created-at: 2026-11-18T20:47:01Z
updated-at: 2026-11-22T18:45:50Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-377078
claimed-at: 2026-11-22T18:42:53Z
---

# CI auf Windows dauert 10 Minuten: internal/cli und internal/tui in Shards aufteilen, -race nur auf Linux

## Definition of Done

- [x] test (windows-latest) laeuft in 3 Shards, die zusammen jeden Test aus internal/cli und internal/tui genau einmal laufen lassen und alle anderen Pakete einmal; Beleg im CI-Log
  proof: scripts/test-shard.sh (go test -list, (NR-1)%n==i; lokal 447/447 cli, 451/451 tui); .github/workflows/ci.yaml:101-103; CI-Log Lauf 45055768953 'shard i/3: N tests of'
- [x] -race laeuft auf ubuntu-latest weiter, auf macOS und Windows nicht
  proof: .github/workflows/ci.yaml:99-103 race: -race nur bei ubuntu-latest
- [x] fail-fast ist aus: ein rotes macOS bricht Windows nicht ab
  proof: .github/workflows/ci.yaml:96 fail-fast: false
- [ ] Der laengste test-Job eines PR braucht unter 5 Minuten; Beleg: Laufzeit des PR-Laufs

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 18:42 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 18:43 · Alexander Sacharov** — Befund 2026-11-22: Sharding, -race nur Linux und fail-fast:false sind schon auf master (5b16e491c, scripts/test-shard.sh). DoD 1-3 erfuellt. DoD 4 nicht: PR-Laeufe 45055768953/45055352215/45054660701 brauchen test (windows-latest 0/3) 363-498 s und test (macos-latest) 376-414 s. Ursache: test-shard.sh ruft go test fuer 'alle anderen Pakete', internal/cli und internal/tui nacheinander auf (Win 0/3: role 154 s + cli 193 s + tui 56 s seriell). Entscheidung: die drei go-test-Aufrufe eines Shards laufen als eigene Prozesse gleichzeitig, Ausgabe gepuffert und danach in fester Reihenfolge. Kein weiterer Runner, kein Umbau der Matrix. Prozesse sind getrennt, also kein geteilter globaler Zustand; cli und tui liefen vorher als Pakete ohnehin parallel.
- **2026-11-22 18:43 · Alexander Sacharov** — decision (Alex 22.11.): zu Ende führen. Stand: Shards/race/fail-fast seit 5b16e491c auf master (DoD 1-3 mit CI-Beleg abhaken). Fehlt DoD 4: letzter grüner master-Lauf macOS 5:29, Windows-Shards ~5:00, Linux 3:38 - macOS ebenfalls in Shards (Windows ggf. 4), bis der längste test-Job unter 5 min liegt; Beleg: PR-Lauf.
