---
id: 01M263RFEQMYPPT67D40Y5SW81
title: Die blocked-by-Warnung trifft seltener und dafuer richtig
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "'jaira validate' meldet einen nicht deklarierten Handle nur dort, wo es plausibel eine Abhaengigkeit ist - und die echten Befunde gehen nicht mehr darin unter"
context: |-
  Gezaehlt am 09.09. auf dem Board eines anderen Projekts: 'jaira validate' meldet 28 Probleme. 23 davon sind 'Kontext oder Note nennt Handle X, der nicht in blocked-by steht' - 12 aus Notizen, 11 aus dem Kontext. Nur 5 sind echte Befunde ('cannot leave the backlog').
  Bei dieser Quote liest man die Ausgabe nicht mehr. Genau die 5, die zaehlen, verschwinden zwischen 23, die nichts bedeuten.
  Warum die Quote so hoch ist: Tickets nennen einander voellig normal in Prosa - 'revidiert NJPQWE', 'siehe MQ55FZ', 'Kollision mit DNAEPN'. Ein genannter Handle ist meistens ein Verweis, keine Abhaengigkeit.
  Die Warnung selbst ist richtig und hat schon gewirkt: sie hat die echte Abhaengigkeit VS5DFW -> 81XRXX aufgedeckt, die ich daraufhin deklariert habe.
  Moeglichkeiten, vor dem Bauen abzuwaegen: nur warnen, wenn der genannte Handle noch offen ist; einen Vermerk am Ticket, der die Warnung dauerhaft verstummen laesst ('ist kein Blocker'); oder sie aus der Standardausgabe nehmen und nur unter --strict zeigen.
definition-of-done: "'jaira validate' auf dem Board eines anderen Projekts meldet deutlich weniger als 23 Handle-Warnungen, ohne eine echte Abhaengigkeit zu verlieren; die 5 Backlog-Befunde sind ohne Suchen erkennbar; ein Test deckt die neue Regel ab; go test ./... -race gruen"
tags:
  - cli
blocked-by: []
commits: []
created-at: 2026-09-10T16:52:33Z
updated-at: 2026-09-10T16:52:33Z
---

# Die blocked-by-Warnung trifft seltener und dafuer richtig

## Definition of Done

- [ ] 'jaira validate' auf dem Board eines anderen Projekts meldet deutlich weniger als 23 Handle-Warnungen, ohne eine echte Abhaengigkeit zu verlieren; die 5 Backlog-Befunde sind ohne Suchen erkennbar; ein Test deckt die neue Regel ab; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

