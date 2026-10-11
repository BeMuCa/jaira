---
id: 01M263P629RST9EAHAA48WFNQY
title: "create --lane sagt, welche Lanes es ueberspringt und was dadurch leer bleibt"
status: testing
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
updated-at: 2026-10-11T20:39:46Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-13127
claimed-at: 2026-10-11T20:29:01Z
outcome-what: "skippedByCreate nennt eine Lane mit nicht angekreuzter Option jetzt doch als uebersprungen, wenn sie ein leer gebliebenes Input-Feld der Ziel-Lane schreibt; Test prueft todo+pre-process, nicht brainstorm; NOTES-Zeile ergaenzt"
outcome-why: "Die Meldung meldete 'plan leer' ohne pre-process zu nennen und widersprach sich selbst (Luecke laut blindem Richter)"
outcome-resolves: "DoD 1: create --lane nennt uebersprungene Lanes und leere Inputs, ohne --lane unveraendert, Tests gruen"
review-summary: none
review-gaps: "Entfernt: die Pfad-Regel (Parking-/Frage-Lane, nicht angekreuzte Option) stand doppelt in internal/cli/tickets.go skippedByCreate und core/lane/next.go Set.Next - jetzt eine Methode Lane.OnPath (core/lane/next.go:54), beide rufen sie. Gelassen: gate.InputMissing hat einen Aufrufer, steht aber bewusst neben gate.OutputOwed (gleiche Form, gleiches Paket, nutzt fieldFilled); show --for-lane (internal/cli/flow.go:578) prueft Inputs eigenstaendig, baut dabei aber den Payload - anderer Zweck, nicht zusammengelegt. Kosten: keine Schleife auf heissem Pfad."
---

# create --lane sagt, welche Lanes es ueberspringt und was dadurch leer bleibt

## Definition of Done

- [x] 'jaira create --lane <spaeter>' nennt die uebersprungenen Lanes und die dadurch leeren input-requires-Felder; ohne --lane aendert sich nichts; ein Test deckt die Meldung ab; go test ./... -race gruen
  proof: internal/cli/tickets.go:484 skippedByCreate (+:494 schreibende Lane ausserhalb der Option), core/lane/next.go:54 Lane.OnPath; Tests TestCreateIntoALaterLaneNamesWhatItSkipped, TestCreateIntoALaterLaneReportsItInJSON, TestCreateWithoutLaneReportsNoSkip; go test ./... -race exit 0 (isoliertes HOME)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-11 19:45 · Alexander Sacharov** — Code aus dem A/B-Versuch vom 11.10.2026 uebernommen (Arm B, ein Worker; blinder Richter knapp fuer diese Variante: sauberer, nutzt ticket.SuppliedFields und die Regeln von set.Next). Bekannte Luecke laut Richter, VOR signoff beheben: im Beispiel des Tickets ('--lane in-progress') fehlt pre-process in der Liste der uebersprungenen Lanes, weil Lanes mit nicht angekreuzter Option ausgelassen werden - dabei ist es genau die Lane, die den Plan schreiben wuerde; die Meldung widerspricht sich damit selbst. Alex hat das Ticket von BeMuCa uebernommen.
- **2026-10-11 20:29 · Alexander Sacharov** — parent session: jaira-61
- **2026-10-11 20:29 · Alexander Sacharov** — Offene Entscheidungen: keine. Luecke aus der Richter-Notiz selbst entschieden: eine Lane, deren Option nicht angekreuzt ist, zaehlt trotzdem als uebersprungen, wenn sie (output-produces) ein Feld schreibt, das in der Ziel-Lane als input leer bleibt. Grund: sonst meldet create 'plan leer', ohne pre-process zu nennen - die Meldung widerspricht sich. Lanes ohne solchen Beitrag (z.B. brainstorm, wenn goal per Flag kam) bleiben weiterhin weg, damit keine Lane gemeldet wird, die der Ticket-Pfad ohnehin nie besucht haette.
- **2026-10-11 20:35 · Alexander Sacharov** — go test ./... -race mit dem echten HOME scheitert in ~25 Tests (core/lane, core/move, internal/cli, internal/tui) - unabhaengig von diesem Ticket: ~/.jaira auf diesem Rechner hat ein 'default board' ohne backlog, und die Tests lesen es. Mit HOME=<leeres absolutes Verzeichnis> (gitconfig hineinkopiert, GOCACHE/GOMODCACHE/GOPATH explizit) ist alles gruen, exit 0. Deshalb prueft TestCreateIntoALaterLaneNamesWhatItSkipped auch nur Enthaltensein (todo, pre-process, kein brainstorm) statt der exakten Zeile: backlog erscheint nur auf den eingebauten Lanes.
- **2026-10-11 20:36 · Alexander Sacharov** — critique: nichts. Geprueft: show --for-lane (internal/cli/flow.go:578) hat eigene Missing-Logik, die aber den Payload baut und Text wie 'plan (the ticket has no Plan checklist)' liefert - kein Duplikat von gate.InputMissing, das neben gate.OutputOwed (gleiche Form) passt. Der Payload-diff der critique-Lane zeigt einen fremden Commit (6bd21c7, Ticketdatei-Historie); beurteilt wurde git diff master...HEAD -- internal/cli/tickets.go core/gate/gate.go.
