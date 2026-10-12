---
id: 01M4R72WDX2JTP49N1PHQ484YS
title: "Aus signoff zurückschicken geht ohne --force, nur done braucht es"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Agent (auch ein '!'-Befehl in der Claude-Sitzung) kann ein Ticket aus einer Menschen-Lane wie signoff zurück in eine frühere Lane schicken, ohne --force. Nach vorn (done) bleibt die Sperre: dort nur mit --force oder interaktiv."
context: |-
  Was heute falsch ist: core/gate/gate.go:389 sperrt JEDEN Weg aus einer Lane mit RequiresHumanExit (signoff), egal in welche Richtung.
  Folge: Alex muss für 'zurück in in-progress' jedes Mal 'jaira move <id> --to in-progress --force' tippen. Ein '!'-Befehl in Claude gilt als Agent (nicht Interactive), also greift die Sperre.
  --force ist gefährlich: es schaltet ALLE Gate-Refusals ab. Am 2026-10-12 hat ein --force bei K1W2JR einen veralteten Ticket-Stand über den Ref geschrieben, die Warnungen waren verdeckt.
  Zurückschicken ist keine Abnahme: die Sperre soll verhindern, dass ein Agent seine eigene Arbeit abnimmt. Rückwärts entscheidet nichts.
  Auslöser: Alex am 2026-10-12, nach NZZDRH und AR97K9, die beide mit --force zurück in in-progress gingen.
  Offen nach vorn: done (terminal) und jede Lane nach signoff bleiben gesperrt wie heute.
definition-of-done: "Move aus signoff in eine Lane mit kleinerer Position in der Lane-Reihenfolge (z. B. in-progress) ohne --force und nicht interaktiv: kein CodeNeedsHuman, Ticket wird verschoben; Test in core/gate"
tags:
  - gates
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-12T14:08:23Z
updated-at: 2026-10-12T14:14:01Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1804
claimed-at: 2026-10-12T14:08:49Z
outcome-what: "Gate refuses leaving a requires-human-exit lane only when the target is not earlier in the lane order (core/gate/gate.go:394); agent block, README, docs/AGENTS.md, jaira skill and lane template say the same"
outcome-why: "Sending work back accepts nothing; making it a --force switched off every other refusal on the move, which on 2026-10-12 hid a stale-file overwrite on K1W2JR"
outcome-resolves: "Agents and '!' commands send a ticket from signoff back to in-progress with a plain move; done and later lanes stay a person's (humanexit_test.go)"
---

# Aus signoff zurückschicken geht ohne --force, nur done braucht es

## Definition of Done

- [x] Move aus signoff in eine Lane mit kleinerer Position in der Lane-Reihenfolge (z. B. in-progress) ohne --force und nicht interaktiv: kein CodeNeedsHuman, Ticket wird verschoben; Test in core/gate
  proof: core/gate/gate.go:394; TestHumanExitLaneSendsBackWithoutForce
- [x] Move aus signoff nach done ohne --force und nicht interaktiv: weiterhin CodeNeedsHuman; Test in core/gate
  proof: core/gate/gate.go:394; TestHumanExitLaneStillRefusesForward
- [x] Gleiche Regel für jede Lane mit RequiresHumanExit (nicht nur über den Namen signoff); Test mit einer zweiten solchen Lane
  proof: TestHumanExitRuleFollowsTheContract (core/gate/humanexit_test.go)
- [x] Zeile in core/release/notes/<ID>.md, was der Nutzer jetzt anders tun kann
  proof: core/release/notes/Q484YS.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Gate: CodeNeedsHuman nur noch, wenn das Ziel nicht vor der aktuellen Lane liegt (to >= from per lane.Index) — core/gate/gate.go:389
- [x] Tests in core/gate: signoff->in-progress frei, signoff->done gesperrt, signoff->blocked gesperrt, zweite requires-human-exit-Lane
- [x] Agent-Block-Text in core/board/announce.go für requires-human-exit: 'in, und zurück; nie nach vorn' + lanesection_test
- [x] Doku: README.md (REVIEW human checkpoint), docs/AGENTS.md, core/lane/template.go-Kommentar
- [x] core/release/notes/Q484YS.md
- [x] go test ./... grün

## Progress
- **2026-10-12 14:08 · Alexander Sacharov** — parent session: jaira-4c
- **2026-10-12 14:09 · Alexander Sacharov** — Entschieden ohne Rückfrage: (1) Regel = Ziel-Index < Index der aktuellen Lane in der Lane-Reihenfolge (lane.Index), sonst Sperre wie bisher. Damit bleibt auch signoff->blocked gesperrt (blocked steht am Ende der Reihenfolge) — Parken ist kein Zurückschicken, und das Ticket würde nach dem Entblocken wieder in signoff landen, also gibt es dafür keinen Bedarf. (2) Release-Notiz als core/release/notes/Q484YS.md wie im DoD; das Verzeichnis kommt mit AR97K9 (origin/feat/ar97k9, noch nicht auf master) — bis dahin liest release.go die Datei nicht, beim Merge von AR97K9 wird sie automatisch mitgelesen.
- **2026-10-12 14:10 · Alexander Sacharov** — Plan-Begründung: Die Richtung steht schon in CheckAdvance als from/to (env.Lanes.Index), dieselbe Größe benutzen die OutputOwed- und Spec-Zone-Prüfungen — kein neuer Begriff nötig. Lane 'human' ist RequiresQuestion, nicht RequiresHumanExit; ihr 'never out' im Agent-Block bleibt, nur der HumanExit-Fall bekommt neuen Text (core/board/announce.go:220).
