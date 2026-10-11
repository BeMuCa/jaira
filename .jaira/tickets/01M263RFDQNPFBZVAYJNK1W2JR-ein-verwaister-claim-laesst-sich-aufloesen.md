---
id: 01M263RFDQNPFBZVAYJNK1W2JR
title: Ein verwaister Claim laesst sich aufloesen
status: backlog
ready: true
creator: BeMuCa
assignee: Alexander Sacharov
goal: "Ein Claim, dessen Session gestorben ist, laesst sich mit einem Befehl loesen und blockiert das Ticket nicht auf Dauer"
context: |-
  Auf dem requirementsgenie-Board: 'jaira next --per-lane' meldet '5M11NC carries an abandoned claim by EE-3NX6GL3-690539, last renewed 11h ago'. Zwei Tage vorher meldete dasselbe Board '2CQ613 ... last renewed 436h ago' - also 18 Tage.
  jaira erkennt den verwaisten Claim also und benennt ihn praezise. Was fehlt, ist der Befehl daneben: 'jaira --help' listet claim, aber nichts, was einen Claim wieder loest, und er laeuft auch nicht von selbst ab.
  NICHT geprueft und vor dem Bauen zu klaeren: ob es intern doch eine Ablaufzeit gibt, die nur nicht greift - ich habe das aus der Befehlsliste geschlossen, nicht im Code nachgelesen.
  Warum es zaehlt: mehrere parallele Sessions schreiben dieses Board, und eine getoetete Session hat nie einen Zug, ihren Claim zurueckzugeben. Genau der Fall, fuer den das Board gebaut ist.
definition-of-done: "Ein verwaister Claim laesst sich per Befehl loesen; die Meldung in 'jaira next' nennt diesen Befehl; ein Test deckt Loesen und die Meldung ab; go test ./... -race gruen"
tags:
  - cli
blocked-by: []
commits: []
created-at: 2026-09-10T16:52:33Z
updated-at: 2026-10-11T19:45:09Z
updated-by: Alexander Sacharov
---

# Ein verwaister Claim laesst sich aufloesen

## Definition of Done

- [ ] Ein verwaister Claim laesst sich per Befehl loesen; die Meldung in 'jaira next' nennt diesen Befehl; ein Test deckt Loesen und die Meldung ab; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-11 19:45 · Alexander Sacharov** — Code aus dem A/B-Versuch vom 11.10.2026 uebernommen (Arm A; blinder Richter knapp fuer diese Variante; Review: approve). Befund: Claim-Ablauf (ClaimTTL 30 min) und 'jaira claim <id> --release' gab es schon, der Kontext irrte; die Aenderung nennt den Befehl in der Warnung von jaira next. Offen laut Review: der DoD-Proof nannte releaseHint, das inzwischen in warnStaleClaims eingefaltet ist; '--release' ist in 'jaira --help' nicht sichtbar (nicht DoD). Die Fixture-Reparatur in claimrelease_test.go ueberschneidet sich mit AK3Z6D (HOME in TestMain). Alex hat das Ticket von BeMuCa uebernommen.
