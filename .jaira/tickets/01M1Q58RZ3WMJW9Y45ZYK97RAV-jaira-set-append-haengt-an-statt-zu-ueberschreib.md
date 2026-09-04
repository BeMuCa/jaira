---
id: 01M1Q58RZ3WMJW9Y45ZYK97RAV
title: jaira set --append haengt an statt zu ueberschreiben
status: critique
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "jaira set <id> feld=wert --append haengt den Wert mit Leerzeile an den bestehenden an, statt ihn zu ersetzen - Review-Runden behalten ihre Historie im Feld"
context: "Berk am 04.09. (Feedback L119): review-Felder werden je Loop-Durchlauf ueberschrieben, die Runden-Historie lebt nur in Notes. Auf einem Board mit critique/optimize/testing-Loops gehen Begruendungen frueherer Runden im Feld verloren. set validiert bewusst nichts (Invariante 9) - --append aendert nur die Schreibweise: alt + Leerzeile + neu; leeres Altfeld verhaelt sich wie ohne Flag."
definition-of-done: set feld=neu --append ergibt alt+Leerzeile+neu; ohne Altwert identisch zu normalem set; ein Test deckt beides; --help erklaert das Flag
tags: []
blocked-by: []
commits: []
created-at: 2026-09-04T21:30:19Z
updated-at: 2026-09-04T21:35:23Z
claimed-by: EE-3NX6GL3-370754
claimed-at: 2026-09-04T21:31:12Z
updated-by: BeMuCa
outcome-what: "jaira set bekommt --append: Skalare haengen den neuen Wert per Leerzeile unter den alten (leeres Feld = normales set), Listen haengen die neuen Eintraege hinten an; Long-Text dokumentiert es"
outcome-why: "Feedback L119: Review-Felder verlieren je Loop-Runde ihre Historie - Runden sollen im Feld stapeln koennen statt nur in Notes zu ueberleben"
outcome-resolves: "Beide DoD-Faelle je mit Test (Scalar-Stapeln inkl. Overwrite-Gegenprobe, Listen-Anhang); set bleibt validierungsfrei (Invariante 9 unangetastet)"
executed-by: fable
---

# jaira set --append haengt an statt zu ueberschreiben

## Definition of Done

- [x] set feld=neu --append ergibt alt+Leerzeile+neu; ohne Altwert identisch zu normalem set; ein Test deckt beides; --help erklaert das Flag
  proof: internal/cli/tickets.go newSetCmd: --append-Flag, Scalar alt+\n+neu (leer=wie set), Listen alt+neu; Long-Text erklaert es; TestSetAppendKeepsTheOldValue + TestSetAppendExtendsAList gruen (-race)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

