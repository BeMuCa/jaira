---
id: 01M263P629RST9EAHAA48WFNQY
title: "create --lane sagt, welche Lanes es ueberspringt und was dadurch leer bleibt"
status: done
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
commits:
  - 6bd21c73399d119c3edd146759e47737c2828ef0
  - 2b9c8cf08cea739872e7e882c8610fa723ad4f0a
  - 45367281f3dad6bc728fa1dd71d2e0ac698b8980
  - 88a21c9bf0182e2268ed3c827ae879e4ce7d9d1a
created-at: 2026-09-10T16:51:18Z
updated-at: 2026-10-11T20:47:20Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-13127
claimed-at: 2026-10-11T20:29:01Z
outcome-what: "Review: Diff gegen DoD geprueft, Tests mit isoliertem HOME gruen, Probe auf Wegwerf-Board bestaetigt die Meldung"
outcome-why: "Zweites Modell beurteilt den Diff vor der menschlichen Abnahme"
outcome-resolves: "DoD 1 bestaetigt; keine Defekte, vier nicht blockierende Anmerkungen in review-gaps"
review-summary: "jaira create --lane <spaetere Lane> gibt jetzt zwei Zeilen mehr aus: 'Skipped on the way to <lane>: …' (Lanes zwischen Default-Lane und Ziel, die das Ticket laut Lane.OnPath besucht haette, plus eine option-gesperrte Lane, wenn sie ein leer gebliebenes Input-Feld schreibt - so erscheint pre-process bei --lane in-progress) und 'Empty, though <lane> takes it as input: …' (neue gate.InputMissing: input-requires der Ziel-Lane, die leer sind und nicht per Flag kommen koennten). --json traegt skipped_lanes/empty_inputs nur, wenn nicht leer. set.Next nutzt dieselbe neue Regel Lane.OnPath (reiner Refactor, Verhalten gleich). Parking-Ziel (blocked) meldet nichts. NOTES-Zeile unter Unreleased."
review-gaps: "Keine Defekte gefunden. Kleinigkeiten, nicht blockierend: (1) Auf einem Board mit backlog wird backlog selbst als 'skipped' genannt (Slice beginnt bei der Default-Lane) - wahr, aber Rauschen. (2) brainstorm wird nie genannt, auch wenn --goal fehlt: goal zaehlt als SuppliedField und erscheint stattdessen bei den fehlenden Promotion-Feldern - so gewollt und im Kommentar begruendet. (3) TestCreateIntoALaterLaneReportsItInJSON prueft skipped_lanes nur auf nicht-leer. (4) Der Payload-diff enthaelt den fremden Commit 6bd21c7 (beruehrt nur die Ticketdatei); beurteilt wurde 15f8057..HEAD."
test-verdict: "pass: Binaer aus dieser Branch auf Wegwerf-Board (HOME isoliert). 'create --lane in-progress --goal g --context c --dod d' -> 'Skipped on the way to in-progress: backlog, todo, pre-process' + 'Empty, though in-progress takes it as input: plan'; ohne --lane keine Zusatzzeile; --lane todo -> nur 'Skipped ...: backlog', keine Empty-Zeile; --json traegt skipped_lanes/empty_inputs. go test ./... -race exit 0."
review-verdict: "Der Diff erfuellt die DoD: Meldung bei --lane, ohne --lane unveraendert (Text und JSON), drei Tests, go test ./... -race gruen mit isoliertem HOME (selbst nachgeprueft). Die Richter-Luecke (pre-process fehlte) ist behoben und per Probe bestaetigt. Freigabe empfohlen."
review-check: "1. go build -o /tmp/jb ./cmd/jaira  2. mkdir /tmp/probe && cd /tmp/probe && git init && /tmp/jb init  3. /tmp/jb create Probe --lane in-progress --assignee t --goal g --context c --dod d  -> unter 'Created …' stehen 'Skipped on the way to in-progress: … todo, pre-process' und 'Empty, though in-progress takes it as input: plan'  4. /tmp/jb create Probe2 --assignee t --goal g --context c --dod d  -> keine der beiden Zeilen  5. /tmp/jb --json create Probe3 --lane in-progress --assignee t --goal g --context c --dod d  -> JSON hat skipped_lanes und empty_inputs: [\"plan\"]"
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
