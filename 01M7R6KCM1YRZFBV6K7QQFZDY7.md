---
id: 01M7R6KCM1YRZFBV6K7QQFZDY7
title: "CI auf Windows dauert 10 Minuten: internal/cli und internal/tui in Shards aufteilen, -race nur auf Linux"
status: todo
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
updated-at: 2026-11-18T20:47:01Z
---

# CI auf Windows dauert 10 Minuten: internal/cli und internal/tui in Shards aufteilen, -race nur auf Linux

## Definition of Done

- [ ] test (windows-latest) laeuft in 3 Shards, die zusammen jeden Test aus internal/cli und internal/tui genau einmal laufen lassen und alle anderen Pakete einmal; Beleg im CI-Log
- [ ] -race laeuft auf ubuntu-latest weiter, auf macOS und Windows nicht
- [ ] fail-fast ist aus: ein rotes macOS bricht Windows nicht ab
- [ ] Der laengste test-Job eines PR braucht unter 5 Minuten; Beleg: Laufzeit des PR-Laufs

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

