---
id: 01M6C37ZE3B9SBJRTH1GXYE0M3
title: Windows-Zip von jaira wird vom Virenscanner gemeldet
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Die Windows-Ausgabe von jaira lässt sich beim Kunden herunterladen und starten, ohne dass ein Virenscanner sie meldet oder in Quarantäne schiebt."
context: |-
  Was ist falsch:
  - Alex (30.10.2026, beim Kunden BMW): der Virenscanner meldet das Windows-Zip von jaira als Virus.
  - Welcher Scanner und welcher Fundname: noch unbekannt – beim nächsten Mal Screenshot/Name festhalten.
  - Alex fragte auch: warum gibt es zwei Windows-Zips im Release? Klären, ob eins reicht.
  Vermutung, nicht geprüft: unsigniertes Binary (typischer Fehlalarm bei selbst gebauten Go/CLI-Programmen).
  Mögliche Wege: Code-Signatur (Zertifikat), Fehlalarm beim Hersteller melden, Build-Optionen prüfen.
  Nicht: dem Kunden raten, den Scanner zu umgehen.
definition-of-done: "Ursache ist benannt (Scanner, Fundname, warum das Binary auffällt)"
tags:
  - release
blocked-by: []
related: []
commits: []
created-at: 2026-11-01T17:41:46Z
updated-at: 2026-11-01T17:41:46Z
---

# Windows-Zip von jaira wird vom Virenscanner gemeldet

## Definition of Done

- [ ] Ursache ist benannt (Scanner, Fundname, warum das Binary auffällt)
- [ ] Das Windows-Release besteht aus einem Zip, oder es ist dokumentiert, wofür das zweite ist
- [ ] Ein frisch heruntergeladenes Windows-Release wird vom Scanner beim Kunden nicht mehr gemeldet

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

