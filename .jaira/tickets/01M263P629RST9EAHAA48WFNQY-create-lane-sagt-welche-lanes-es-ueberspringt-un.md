---
id: 01M263P629RST9EAHAA48WFNQY
title: "create --lane sagt, welche Lanes es ueberspringt und was dadurch leer bleibt"
status: backlog
ready: true
creator: BeMuCa
assignee: Alexander Sacharov
goal: "Wer ein Ticket direkt in eine spaetere Lane anlegt, erfaehrt sofort, welche Pflichtfelder dadurch ungefuellt bleiben"
context: |-
  Nachgestellt am 10.09. auf einem Wegwerf-Board: 'jaira create "Probe" --lane in-progress --assignee t --goal g --context c --dod d' gibt genau eine Zeile aus: 'Created GBF31Q Probe'. Keine Warnung.
  Das Ticket landet in Implementing, und die '## Plan'-Sektion bleibt der Platzhalter '<Steps, in order - filled in by the pre-process step, or by you.>'.
  Das ist ein Widerspruch im Werkzeug: in-progress deklariert 'input-requires: [goal, definition-of-done, context, plan]', und der Lane-Prompt sagt 'Carry out the plan on this ticket.' Es gibt keinen Plan, und es kann per Konstruktion keinen geben.
  So gemeldet aus einer Multi-Agent-Session am 10.09.: JP1PA1 sah auf dem Board unangefangen aus, obwohl die Arbeit fertig war.
  Kein Verbot gewuenscht - direkt in eine Lane anzulegen ist ein legitimer Schnellweg. Nur soll dabei stehen, was man sich einhandelt.
definition-of-done: "'jaira create --lane <spaeter>' nennt die uebersprungenen Lanes und die dadurch leeren input-requires-Felder; ohne --lane aendert sich nichts; ein Test deckt die Meldung ab; go test ./... -race gruen"
tags:
  - cli
blocked-by: []
commits: []
created-at: 2026-09-10T16:51:18Z
updated-at: 2026-10-11T19:45:02Z
updated-by: Alexander Sacharov
---

# create --lane sagt, welche Lanes es ueberspringt und was dadurch leer bleibt

## Definition of Done

- [ ] 'jaira create --lane <spaeter>' nennt die uebersprungenen Lanes und die dadurch leeren input-requires-Felder; ohne --lane aendert sich nichts; ein Test deckt die Meldung ab; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-11 19:45 · Alexander Sacharov** — Code aus dem A/B-Versuch vom 11.10.2026 uebernommen (Arm B, ein Worker; blinder Richter knapp fuer diese Variante: sauberer, nutzt ticket.SuppliedFields und die Regeln von set.Next). Bekannte Luecke laut Richter, VOR signoff beheben: im Beispiel des Tickets ('--lane in-progress') fehlt pre-process in der Liste der uebersprungenen Lanes, weil Lanes mit nicht angekreuzter Option ausgelassen werden - dabei ist es genau die Lane, die den Plan schreiben wuerde; die Meldung widerspricht sich damit selbst. Alex hat das Ticket von BeMuCa uebernommen.
