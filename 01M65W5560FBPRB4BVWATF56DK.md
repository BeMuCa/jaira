---
id: 01M65W5560FBPRB4BVWATF56DK
title: critique-Regeln aus .jaira/lanes in die eingebaute Lane übernehmen
status: testing
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
updated-at: 2026-10-30T07:46:08Z
assignee: AlSa
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-4867
claimed-at: 2026-10-30T07:42:45Z
outcome-what: "Eingebaute critique-Lane (core/lane/builtin/25-critique.md) übernimmt die zwei Regeln der lokalen Lane: Notes vor dem Diff lesen, jeder Durchgang nach dem ersten liest nur die letzten Findings und ihre Reparatur. Neue Version in shipped.sha256, Release-Note, .jaira/lanes auf den eingebauten Stand."
outcome-why: "Jedes andere Board bekam bisher die schwächere critique, deren Schleife Erledigtes neu fand und nicht auf null kam."
outcome-resolves: update --lanes überspringt critique.md hier nicht mehr; neue Boards bekommen die Regeln.
review-summary: none
review-gaps: none
---

# critique-Regeln aus .jaira/lanes in die eingebaute Lane übernehmen

## Definition of Done

- [x] core/lane/builtin/25-critique.md trägt notes in input-requires und beide Regeln im Prompt
  proof: core/lane/builtin/25-critique.md:10 (notes), Prompt-Abschnitte 'Read the notes before you read the diff' und 'Every pass after the first reads less'
- [x] Nach dem Wechsel ist .jaira/lanes/critique.md dieses Repos byte-gleich mit der eingebauten Datei, 'jaira update --lanes' meldet 0 skipped
  proof: cmp .jaira/lanes/critique.md core/lane/builtin/25-critique.md gleich; neues Binary: 'jaira update --lanes' -> lanes: 2 updated, 11 unchanged, 0 skipped
- [x] core/shipped/shipped.sha256 kennt die neue Version, go test ./... grün
  proof: core/shipped/shipped.sha256:27 (7d3e33a0…); go test ./... ohne FAIL, rc=0
- [x] core/release/notes/<ID>.md sagt, dass critique jetzt Notes liest und jede Runde weniger liest
  proof: core/release/notes/TF56DK.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-30 07:42 · Alexander Sacharov** — Keine offenen Entscheidungen: Inhalt ist die lokale critique.md dieses Repos; 'default-board: false' der eingebauten Datei bleibt (Feld für Mehr-Board-Typen), die lokale Datei bekommt es dazu und wird so byte-gleich.
- **2026-10-30 07:45 · Alexander Sacharov** — shipped.sha256 bekam nur die eine critique-Zeile: scripts/shipped-history.sh hätte ~60 alte Versionen anderer Dateien nachgetragen, das gehört nicht zu diesem Ticket. 'update --lanes' hat dabei optimize.md und testing.md auf die eingebaute Version gebracht (nur 'default-board: false' dazu) — mitcommittet, damit das Board sauber ist.
- **2026-10-30 07:45 · Alexander Sacharov** — critique: Diff gelesen (Lane-Datei, eine Manifestzeile, Release-Note, drei .jaira/lanes). notes als Input wird bereits unterstützt (internal/cli/notesinput_test.go), der for-lane-Payload ist complete. Kein Finding.
- **2026-10-30 07:46 · Alexander Sacharov** — optimize: nichts zu entfernen — der Diff ist Text plus eine Hash-Zeile, kein Code.
