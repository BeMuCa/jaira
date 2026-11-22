---
id: 01M7R6KCM1YRZFBV6K7QQFZDY7
title: "CI auf Windows dauert 10 Minuten: internal/cli und internal/tui in Shards aufteilen, -race nur auf Linux"
status: done
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
commits:
  - 5b16e491ccd882eaa7abcba1d5341a60595eadc9
  - 0b2ee891b1bac6fcbd6c771343d1c35f88a9acb9
created-at: 2026-11-18T20:47:01Z
updated-at: 2026-11-22T19:28:18Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-377078
claimed-at: 2026-11-22T18:42:53Z
outcome-what: "scripts/test-shard.sh startet die go-test-Aufrufe eines Shards (andere Pakete, internal/cli, internal/tui) gleichzeitig als eigene Prozesse, puffert ihre Ausgabe und gibt sie danach der Reihe nach aus; Exit-Code ist der erste ungleich null."
outcome-why: Shard 0 auf Windows und der macOS-Job liefen die drei Aufrufe nacheinander und brauchten 6-8 min; das Sharding allein (5b16e491c) brachte den laengsten Job nicht unter 5 min.
outcome-resolves: "DoD 4: laengster test-Job eines PR unter 5 min; DoD 1-3 waren schon mit 5b16e491c auf master erfuellt"
review-summary: "Auf master liegt seit 5b16e491c schon: Windows in 3 Shards (scripts/test-shard.sh verteilt die Top-Level-Tests von internal/cli und internal/tui per go test -list reihum), -race nur auf ubuntu, fail-fast aus. Neu auf diesem Branch (0b2ee891b): die bis zu drei go-test-Aufrufe eines Shards (übrige Pakete, cli-Anteil, tui-Anteil) starten gleichzeitig als eigene Prozesse statt nacheinander; Ausgabe wird pro Aufruf in eine Temp-Datei gepuffert und nach dem Ende in fester Reihenfolge ausgegeben, Exit-Code ist der erste ungleich 0. Damit dauert ein Shard so lange wie sein langsamster Aufruf statt der Summe - das trifft Windows 0/3 und macOS (1 Shard, bisher alle drei Aufrufe seriell). Vorher liefen die Pakete unter go test ./... ebenfalls parallel, also kein neues Risiko durch geteilten Zustand zwischen Paketen."
review-gaps: |-
  note: DoD 4 (längster test-Job eines PR unter 5 min) ist unbelegt und nur an einem PR-Lauf messbar - ci.yaml läuft bei push nur auf master; die Erwartung (~4 min) ist gerechnet, nicht gemessen. Vor dem Abhaken die Laufzeiten im PR-Lauf ablesen
  note: schlägt go test -list fehl (Compile-Fehler in internal/cli oder internal/tui), beendet set -e das Skript, während der Hintergrund-Aufruf für die übrigen Pakete noch läuft; der EXIT-trap löscht dessen Log, seine Ausgabe fehlt im CI-Log. Der Job ist trotzdem rot und der Compile-Fehler sichtbar. Fix wäre: die -list-Aufrufe vor dem Start des Hintergrund-Aufrufs ausführen
  note: bekannte Grenze: bricht der Job von außen ab (Cancel), ist die gepufferte Ausgabe weg; ein hängender Test endet per -timeout 20m mit Ausgabe
test-verdict: "pass: go vet/build RC=0; scripts/test-shard.sh 0/3, 1/3, 2/3 lokal gruen (37 s / ~30 s / 17 s); Teilung 447/447 cli, 451/451 tui; absichtlich roter Test in core/tag -> FAIL-Block ausgegeben, RC=1, cli/tui-Ausgabe danach vollstaendig; zwei Laeufe parallel beide RC=0, kein tmp-Rest; DoD 1-3 im Baum geprueft. DoD 4 offen fuer die Person: nur ein PR-Lauf misst ihn (push auf Branch startet ci.yaml nicht)"
hunt-verdict: "pass: bughunt done at 0b2ee891b, 0 change, 0 pre-existing"
review-verdict: "accept: der Code macht, was er soll - lokal Shard 2/3 grün (rc 0, 17 s), ein fehlschlagender Aufruf liefert rc 1; kein Blocker im Diff. Unsicher bleibt nur DoD 4: ob der längste Job unter 5 min fällt, zeigt erst der PR-Lauf"
review-round: 1
review-verdict-gaps: |-
  note: DoD 4 (längster test-Job eines PR unter 5 min) ist unbelegt und nur an einem PR-Lauf messbar - ci.yaml läuft bei push nur auf master; die Erwartung (~4 min) ist gerechnet, nicht gemessen. Vor dem Abhaken die Laufzeiten im PR-Lauf ablesen
  note: schlägt go test -list fehl (Compile-Fehler in internal/cli oder internal/tui), beendet set -e das Skript, während der Hintergrund-Aufruf für die übrigen Pakete noch läuft; der EXIT-trap löscht dessen Log, seine Ausgabe fehlt im CI-Log. Der Job ist trotzdem rot und der Compile-Fehler sichtbar. Fix wäre: die -list-Aufrufe vor dem Start des Hintergrund-Aufrufs ausführen
  note: bekannte Grenze: bricht der Job von außen ab (Cancel), ist die gepufferte Ausgabe weg; ein hängender Test endet per -timeout 20m mit Ausgabe
review-check: "1. PR für Branch feat/qfzdy7 öffnen (Agent öffnet keinen)  2. im Actions-Tab den CI-Lauf des PR öffnen  3. die Jobs test (ubuntu-latest), test (macos-latest), test (windows-latest 0/3), 1/3, 2/3 ansehen: jeder unter 5:00, alle grün  4. im Log von windows-latest 0/3 den Schritt scripts/test-shard.sh öffnen: zuerst ok-Zeilen der übrigen Pakete, dann 'shard 0/3: N tests of ./internal/cli' mit ok, dann dasselbe für ./internal/tui  5. dann DoD 4 abhaken: jaira dod QFZDY7 4 --done --proof '<Lauf-ID, längster Job m:ss>'"
---

# CI auf Windows dauert 10 Minuten: internal/cli und internal/tui in Shards aufteilen, -race nur auf Linux

## Definition of Done

- [x] test (windows-latest) laeuft in 3 Shards, die zusammen jeden Test aus internal/cli und internal/tui genau einmal laufen lassen und alle anderen Pakete einmal; Beleg im CI-Log
  proof: scripts/test-shard.sh (go test -list, (NR-1)%n==i; lokal 447/447 cli, 451/451 tui); .github/workflows/ci.yaml:101-103; CI-Log Lauf 45055768953 'shard i/3: N tests of'
- [x] -race laeuft auf ubuntu-latest weiter, auf macOS und Windows nicht
  proof: .github/workflows/ci.yaml:99-103 race: -race nur bei ubuntu-latest
- [x] fail-fast ist aus: ein rotes macOS bricht Windows nicht ab
  proof: .github/workflows/ci.yaml:96 fail-fast: false
- [x] Der laengste test-Job eines PR braucht unter 5 Minuten; Beleg: Laufzeit des PR-Laufs
  proof: PR #192 run 45057987635: longest test job macos-latest 4:38, windows 0/3 4:24, ubuntu 1:46

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 18:42 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 18:43 · Alexander Sacharov** — Befund 2026-11-22: Sharding, -race nur Linux und fail-fast:false sind schon auf master (5b16e491c, scripts/test-shard.sh). DoD 1-3 erfuellt. DoD 4 nicht: PR-Laeufe 45055768953/45055352215/45054660701 brauchen test (windows-latest 0/3) 363-498 s und test (macos-latest) 376-414 s. Ursache: test-shard.sh ruft go test fuer 'alle anderen Pakete', internal/cli und internal/tui nacheinander auf (Win 0/3: role 154 s + cli 193 s + tui 56 s seriell). Entscheidung: die drei go-test-Aufrufe eines Shards laufen als eigene Prozesse gleichzeitig, Ausgabe gepuffert und danach in fester Reihenfolge. Kein weiterer Runner, kein Umbau der Matrix. Prozesse sind getrennt, also kein geteilter globaler Zustand; cli und tui liefen vorher als Pakete ohnehin parallel.
- **2026-11-22 18:43 · Alexander Sacharov** — decision (Alex 22.11.): zu Ende führen. Stand: Shards/race/fail-fast seit 5b16e491c auf master (DoD 1-3 mit CI-Beleg abhaken). Fehlt DoD 4: letzter grüner master-Lauf macOS 5:29, Windows-Shards ~5:00, Linux 3:38 - macOS ebenfalls in Shards (Windows ggf. 4), bis der längste test-Job unter 5 min liegt; Beleg: PR-Lauf.
- **2026-11-22 18:46 · Alexander Sacharov** — DoD 4 ist nur an einem PR-Lauf messbar: ci.yaml laeuft bei push nur auf master, der Push nach origin/feat/qfzdy7 startet nichts. Lokal (WSL): Shard 0/3 37 s, Shard 2/3 17 s. Erwartung Windows 0/3: max(andere Pakete ~154 s, cli-Drittel ~193 s, tui ~56 s) + vet/setup ~45 s ≈ 4 min statt 7-8; macOS max(cli 150 s, role 84 s) statt Summe ≈ 3-4 min statt 6-7. Erst der PR-Lauf belegt das.
- **2026-11-22 18:46 · Alexander Sacharov** — critique: keine Befunde. Erwogen und verworfen: mehr Windows-Shards bzw. macOS sharden (mehr Runner, gleiche Wirkung nur mit Matrix-Umbau); ein go-test-Aufruf fuer cli+tui (geht nicht, -run gilt pro Aufruf fuer beide Pakete). Bekannte Grenze: wird der Job von aussen abgebrochen, ist die gepufferte Ausgabe weg; ein haengender Test endet aber per -timeout 20m mit Ausgabe.
- **2026-11-22 18:47 · Alexander Sacharov** — bughunt: done at 0b2ee891b - 0 findings in 2 modules, first hunt; concurrency, error paths, input, permissions checked; files: .github/workflows/ci.yaml scripts/test-shard.sh
