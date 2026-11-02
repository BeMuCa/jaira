---
id: 01M6EP2AKFT0BTQEV7388X9SHE
title: "Unter Windows bricht ein jaira-Befehl mit 'Access is denied' ab, statt auf eine gerade freigegebene Sperre zu warten"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Eine Sperre unter .../locks/ wartet unter Windows auch dann, wenn die Sperrdatei gerade geloescht wird, und der Befehl laeuft danach normal durch."
context: |-
  Was falsch ist:
  - core/ticket/store.go:1531 lock() legt die Sperrdatei mit O_CREATE|O_EXCL an. Nur bei os.IsExist wird gewartet und neu versucht; jeder andere Fehler bricht ab.
  - Unter Windows liefert das Anlegen 'Access is denied' (ERROR_ACCESS_DENIED), solange die alte Sperrdatei noch im Loeschen ist (os.Remove im unlock, Zeile 1544). Das ist kein IsExist -> der Befehl bricht ab, statt zu warten.
  Gefunden am 02.11.2026 in CI von PR #120 (182JT1, Lauf 41492357461, test windows-latest):
  - TestRestoreWaitsForTheMilestoneLock (internal/cli/milestones_test.go:559) -> 'restore after unlock: open ...\locks\milestones.lock: Access is denied.'
  - 182JT1 aendert nichts an Sperren oder Milestones; master lief am 02.11. unter Windows gruen. Also flackernd, abhaengig vom Timing.
  Warum es zaehlt: jaira soll mehrere parallele Sessions aushalten. Unter Windows kann jede Sperre (Tickets, Milestones) genau im Moment der Freigabe so abbrechen.
  Naheliegende Loesung: in lock() ERROR_ACCESS_DENIED (unter Windows) wie IsExist behandeln und bis lockTimeout weiter versuchen.
definition-of-done: lock() wartet unter Windows auch bei ERROR_ACCESS_DENIED bis lockTimeout und wirft den Fehler erst danach; Test
tags:
  - concurrency
blocked-by: []
related: []
commits: []
created-at: 2026-11-02T17:49:13Z
updated-at: 2026-11-02T17:49:13Z
---

# Unter Windows bricht ein jaira-Befehl mit 'Access is denied' ab, statt auf eine gerade freigegebene Sperre zu warten

## Definition of Done

- [ ] lock() wartet unter Windows auch bei ERROR_ACCESS_DENIED bis lockTimeout und wirft den Fehler erst danach; Test
- [ ] TestRestoreWaitsForTheMilestoneLock laeuft auf windows-latest stabil (z. B. -count=20 in einem Lauf gruen)
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

