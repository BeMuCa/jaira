---
id: 01M4R72WDX2JTP49N1PHQ484YS
title: "Aus signoff zurückschicken geht ohne --force, nur done braucht es"
status: done
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
commits:
  - 8cea1b42e7b0a499bc0010fa363427e23b9379ba
  - 37208bf84f3d814d1430fc51069fd5175ec3ab3f
  - f77567b4dd6493165165ebf9ea151209652e779f
  - b024f8ab095fd083f7d5e27bc0308f2617d8af2a
created-at: 2026-10-12T14:08:23Z
updated-at: 2026-10-12T14:39:00Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1804
claimed-at: 2026-10-12T14:08:49Z
outcome-what: "Gate refuses leaving a requires-human-exit lane only when the target is not earlier in the lane order (core/gate/gate.go:394); agent block, README, docs/AGENTS.md, jaira skill and lane template say the same"
outcome-why: "Sending work back accepts nothing; making it a --force switched off every other refusal on the move, which on 2026-10-12 hid a stale-file overwrite on K1W2JR"
outcome-resolves: "Agents and '!' commands send a ticket from signoff back to in-progress with a plain move; done and later lanes stay a person's (humanexit_test.go)"
review-summary: "CheckAdvance (core/gate/gate.go:395) now refuses leaving a requires-human-exit lane only when the target lane sits later in the lane order (to > from, the same lane.Index pair the outcome and spec-zone checks use); moving back to any earlier lane passes without --force, while done, blocked and any later lane still need --force or the board. The refusal text, the agent block and lane line in announce.go (regenerated into CLAUDE.md/AGENTS.md), the lane.go/template.go contract comments, README, docs/AGENTS.md and the jaira skill all say the same. Three gate tests cover back, forward (done, blocked) and a second human-exit lane set via the contract; Get returns *Lane, so that mutation is real."
review-gaps: "1. The release note sits in core/release/notes/Q484YS.md, but release.go on this branch embeds only NOTES.md - jaira update will not show it until AR97K9 (origin/feat/ar97k9) is merged first; merge order, not code. 2. Round-2 delta b024f8a (lane.go comment, generated lane lines) matches announce.go word for word; nothing left from the earlier review. 3. Untracked .jaira/milestones/ in the worktree is not this ticket's - do not commit it with this branch."
test-verdict: "pass (round 2): go test ./... exit 0 after b024f8a; CLAUDE.md/AGENTS.md lane section identical to what this branch's 'jaira update' writes; gate unchanged since round 1, where signoff->in-progress exit 0 and ->done exit 3 were checked on a scratch board"
review-verdict: "approve: the diff meets all four DoD items, go test ./... is green, and end to end signoff->in-progress and signoff->review pass without --force (exit 0) while signoff->done is refused with the human-checkpoint line (exit 3). Only condition: merge AR97K9 before this so the release note is read."
review-check: "1. go build -o /tmp/jaira ./cmd/jaira  2. mkdir /tmp/q && cd /tmp/q && git init && git commit --allow-empty -m i && /tmp/jaira init  3. /tmp/jaira create t --goal g --context c --dod d  (note the id)  4. /tmp/jaira move <id> --to signoff --force  5. /tmp/jaira move <id> --to in-progress  -> prints '<id> -> in-progress', exit 0, no --force  6. /tmp/jaira move <id> --to signoff --force, then /tmp/jaira move <id> --to done  -> exit 3, list ends with 'lane \"signoff\" is a human checkpoint ... only back to an earlier lane'  7. go test ./core/gate -run HumanExit -v  -> three PASS"
---

# Aus signoff zurückschicken geht ohne --force, nur done braucht es

## Definition of Done

- [x] Move aus signoff in eine Lane mit kleinerer Position in der Lane-Reihenfolge (z. B. in-progress) ohne --force und nicht interaktiv: kein CodeNeedsHuman, Ticket wird verschoben; Test in core/gate
  proof: core/gate/gate.go:395; TestHumanExitLaneSendsBackWithoutForce
- [x] Move aus signoff nach done ohne --force und nicht interaktiv: weiterhin CodeNeedsHuman; Test in core/gate
  proof: core/gate/gate.go:395; TestHumanExitLaneStillRefusesForward
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
- **2026-10-12 14:14 · Alexander Sacharov** — critique 1: gate.go:394 Bedingung auf 'to > from' vereinfachen (req.To != t.Status ist bei eindeutigen Lane-Indizes dasselbe wie to != from). Sonst nichts: Gate ist der richtige Ort (TUI/CLI haben keine eigene Exit-Prüfung, nur Anzeige in internal/tui), from/to existieren schon, Builtins() parst bei jedem Load neu, also leckt die Mutation in TestHumanExitRuleFollowsTheContract nicht. Rollen-Skills (worker/dispatcher 'never out of it') bleiben: sie beschreiben, dass der Agent an einer Menschen-Lane stoppt, nicht das Gate.
- **2026-10-12 14:16 · Alexander Sacharov** — critique 2: nichts mehr — Befund aus critique 1 ist in 37208bf behoben, der Rest steht auf der Notiz zu critique 1.
- **2026-10-12 14:19 · Alexander Sacharov** — review: Gate-Logik korrekt, e2e geprüft (signoff->in-progress exit 0 ohne --force; ->done und ->blocked exit 3). Offen: lane.go:113 Kommentar veraltet; CLAUDE.md:157/AGENTS.md:130 Lane-Abschnitt nicht regeneriert; Release-Notiz wird erst mit AR97K9 gelesen — AR97K9 zuerst mergen.
- **2026-10-12 14:28 · Alexander Sacharov** — Alex 2026-10-12: send back for the two doc gaps from review (lane.go:113 comment; CLAUDE.md/AGENTS.md lane lines via jaira update).
- **2026-10-12 14:29 · Alexander Sacharov** — in-progress 2: lane.go RequiresHumanExit-Kommentar und die signoff-Zeile in CLAUDE.md/AGENTS.md von Hand auf den Generator-Text gesetzt — nicht per 'jaira update' im Worktree, weil das installierte jaira der alte Build ist und den alten Text schreiben würde; gegen 'jaira update' des Branch-Builds in einer Kopie verglichen: identisch. Punkt 3 (AR97K9 zuerst mergen) ist Reihenfolge beim Merge, kein Code hier.
- **2026-10-12 14:30 · Alexander Sacharov** — critique 3 (Runde nach signoff): Delta b024f8a ist nur ein Kommentar und zwei generierte Zeilen, Wortlaut gleich announce.go:221 — nichts zu sagen.
