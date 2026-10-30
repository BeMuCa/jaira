---
id: 01M65W5560FBPRB4BVWATF56DK
title: critique-Regeln aus .jaira/lanes in die eingebaute Lane übernehmen
status: backlog
ready: true
creator: Alexander Sacharov
goal: "Die eingebaute Lane core/lane/builtin/25-critique.md enthält die beiden Regeln, die bisher nur in .jaira/lanes/critique.md dieses Repos stehen: Notes vor dem Diff lesen, und jeder Durchgang nach dem ersten liest weniger."
context: |-
  Was falsch ist: Jedes Board außer diesem bekommt die schwächere critique-Lane.
  Die Regeln stehen nur in .jaira/lanes/critique.md (Commits 53affda1/X9GJ51 und e442d856/GTQHNH, 17.-18.09.), nie in core/lane/builtin/25-critique.md.
  Aufgefallen: 'jaira update --lanes' am 30.10. hat critique.md als 'edited here' übersprungen; der Diff zeigt die zwei Regeln.
  Regel 1: input-requires bekommt notes; der Prompt sagt 'Read the notes before you read the diff' — was eine Note als repariert oder akzeptiert festhält, wird nicht neu erhoben.
  Regel 2: nur der erste Durchgang liest den ganzen Diff, jeder weitere nur die Findings des Vorgängers und deren Reparatur; ein neuer Fehler in der Reparatur zählt nur, wenn er die DoD bricht.
  Die eingebaute Datei hat 'default-board: false', die lokale nicht — das Feld bleibt.
  Neue Version braucht eine Zeile in core/shipped/shipped.sha256 (scripts/shipped-history.sh), sonst hält 'update --lanes' sie für editiert.
definition-of-done: core/lane/builtin/25-critique.md trägt notes in input-requires und beide Regeln im Prompt
tags: []
blocked-by: []
related: []
commits: []
created-at: 2026-10-30T07:42:27Z
updated-at: 2026-10-30T07:42:41Z
assignee: AlSa
updated-by: Alexander Sacharov
---

# critique-Regeln aus .jaira/lanes in die eingebaute Lane übernehmen

## Definition of Done

- [ ] core/lane/builtin/25-critique.md trägt notes in input-requires und beide Regeln im Prompt
- [ ] Nach dem Wechsel ist .jaira/lanes/critique.md dieses Repos byte-gleich mit der eingebauten Datei, 'jaira update --lanes' meldet 0 skipped
- [ ] core/shipped/shipped.sha256 kennt die neue Version, go test ./... grün
- [ ] core/release/notes/<ID>.md sagt, dass critique jetzt Notes liest und jede Runde weniger liest

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

