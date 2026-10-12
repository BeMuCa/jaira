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
updated-at: 2026-10-12T14:29:51Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1804
claimed-at: 2026-10-12T14:08:49Z
outcome-what: "Gate refuses leaving a requires-human-exit lane only when the target is not earlier in the lane order (core/gate/gate.go:394); agent block, README, docs/AGENTS.md, jaira skill and lane template say the same"
outcome-why: "Sending work back accepts nothing; making it a --force switched off every other refusal on the move, which on 2026-10-12 hid a stale-file overwrite on K1W2JR"
outcome-resolves: "Agents and '!' commands send a ticket from signoff back to in-progress with a plain move; done and later lanes stay a person's (humanexit_test.go)"
review-summary: "core/gate/gate.go:395: the CodeNeedsHuman refusal on leaving a requires-human-exit lane now fires only when the target sits later in the lane order (to > from, the same lane.Index the other checks use). signoff -> in-progress (or any earlier lane) passes for an agent or a '!' command with no --force; signoff -> done and signoff -> blocked (blocked is last in the order) are still refused unless interactive. The refusal message, the agent block (announce.go, CLAUDE.md, AGENTS.md), the generated lane line for human-exit lanes, README, docs/AGENTS.md, the jaira skill and the lane template comment say the same. Three tests in core/gate/humanexit_test.go: back from signoff, forward to done/blocked, and the rule on a second lane flagged requires-human-exit."
review-gaps: "1) core/lane/lane.go:113 doc comment on RequiresHumanExit still says 'no agent may move a ticket out of this lane' — the field's own definition now contradicts the gate. 2) This repo's committed CLAUDE.md:157 and AGENTS.md:130 lane sections still say signoff is 'never out', while line 118/91 of the same files now say send-back is fine; announce.go generates the new text, but the files were not regenerated (jaira update would). 3) core/release/notes/Q484YS.md is not read by this branch's binary: the notes/ directory arrives with AR97K9 (origin/feat/ar97k9), which is not on upstream/master yet — merge AR97K9 first or the line is not shown by jaira update. 4) Minor: the blocked refusal says 'cannot move a ticket on from it' — parking is not 'on', but the rule (blocked stays refused) is the one the ticket decided. No defects in the gate logic: unknown target lanes return earlier, from is never -1 when cur exists."
test-verdict: "pass: go test ./... exit 0; CLI on a scratch board: signoff->in-progress without --force exit 0 and the ticket lands in in-progress, signoff->done exit 3 with the human-checkpoint refusal; the installed old binary refuses the same send-back with exit 3"
review-verdict: "the diff matches all four DoD items and works end to end; approve with two small doc fixes (lane.go:113 comment, regenerate the lane section in CLAUDE.md/AGENTS.md) and merge AR97K9 before this so the release note is read"
review-check: "1. go build -o /tmp/jaira ./cmd/jaira  2. in a scratch repo: git init && git commit --allow-empty -m i && /tmp/jaira init && /tmp/jaira create t --goal g --context c --dod d  (note the id)  3. /tmp/jaira move <id> --to signoff --force  4. /tmp/jaira move <id> --to in-progress  -> prints '<id> → in-progress', exit 0, no --force  5. /tmp/jaira move <id> --to signoff --force, then /tmp/jaira move <id> --to done  -> exit 3, the list includes 'lane \"signoff\" is a human checkpoint ... only back to an earlier lane'  6. go test ./core/gate -run HumanExit -v  -> three PASS"
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
