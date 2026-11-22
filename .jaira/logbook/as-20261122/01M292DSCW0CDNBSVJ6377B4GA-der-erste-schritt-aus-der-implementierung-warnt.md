---
id: 01M292DSCW0CDNBSVJ6377B4GA
title: "Der erste Schritt aus der Implementierung warnt, wenn kein DoD-Punkt abgehakt ist"
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer eine Implementierung verlaesst, ohne einen einzigen Punkt der Definition of Done abgehakt zu haben, erfaehrt das sofort - und nicht erst am Schluss, wenn die terminale Lane den Zug verweigert."
context: |-
  Ein Agent hat ein Ticket ueber acht Lanes und drei Review-Runden gefahren und dabei nur die Plan-Checkliste abgehakt, nie die Definition of Done. Aufgefallen ist es erst beim Sprung nach signoff: neun Verweigerungen auf einmal, ganz am Ende.

  - Es gibt zwei Checklisten im Ticket-Body: '## Definition of Done' (die Abnahmekriterien, die das Gate liest) und '## Plan' (die Methode). 'jaira dod <id> <n> --done' meint die erste, mit --plan die zweite.
  - Die Lane-Anweisung sagt beides bereits: .jaira/lanes/in-progress.md, Zeilen 17-23, samt --proof. Der Agent hat sie gelesen und trotzdem nur --plan abgehakt. Text allein hat also nicht gereicht - die beiden Listen sehen gleich aus, und am Ende der Arbeit fuehlt es sich an, als waere schon alles markiert.
  - Der generierte Block in CLAUDE.md/AGENTS.md ist inzwischen praeziser (core/board/announce.go), aber das ist dieselbe Gattung Gegenmittel: Prosa, die vorher gelesen werden muss.
  - Das Gate selbst ist richtig und soll nicht weicher werden. Falsch ist nur der Zeitpunkt, zu dem man es erfaehrt.
  - Gemeint ist eine Warnung, kein Refus: der Zug geht durch, aber die Zeile steht da. Etwa 'Definition of Done: 0 von 9 abgehakt - die terminale Lane liest diese Liste, nicht den Plan'.
  - Wo genau sie erscheint, ist noch nicht entschieden: beim Verlassen der implementierenden Lane, oder bei jedem Zug eines Tickets, das schon Code hat. Das gehoert in den Entwurf.
definition-of-done: "Ein Zug aus der implementierenden Lane heraus schreibt eine Warnung, wenn kein Punkt der Definition of Done abgehakt ist; der Zug selbst geht durch; in --json steht sie nicht auf stdout zwischen den Daten; ein Test deckt beide Faelle ab (kein Punkt abgehakt: Warnung; mindestens einer: keine)"
tags:
  - gates
  - cli
blocked-by: []
related: []
commits:
  - 2143e9ade54856620ae8cb1da999b9d7a3386f3b
  - 791c8e1e433fb3c108246622b2c74b58e6665310
  - 24867546368d339197dad2d9c8656619d6e9dfa7
created-at: 2026-09-11T20:26:57Z
updated-at: 2026-11-22T19:30:07Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-926893
claimed-at: 2026-11-22T18:49:16Z
outcome-what: "A move out of the implementing lane (the lane producing outcome-what), forward and not into a parking lane, with no DoD item settled, prints a warning on stderr; the move lands. Detection in core/move (Result.DoDUnticked), print in internal/cli/flow.go."
outcome-why: An agent ticked only the plan through eight lanes and met nine refusals at signoff; the news belongs at the step that should tick the items.
outcome-resolves: "DoD 1: warning on 0 ticked, none with one ticked, stderr also under --json, move lands; covered by core/move and internal/cli tests."
review-summary: "core/move.Move now computes Result.DoDUnticked after the write: when the ticket leaves the lane whose outputs name outcome-what (in-progress) for a later lane that is not a parking lane, has DoD items, and not one of them is done or superseded, it carries the item count. The CLI's move prints one line on stderr in text and --json mode ('Definition of Done: 0 of N ticked - the terminal lane reads this list, not the plan' plus the jaira dod command); the move lands unchanged and stdout keeps only the data. The terminal gate is untouched. Covered by a table test in core/move (none/one/superseded ticked, send-back, parked, other lane) and a CLI test of both modes; release note core/release/notes/77B4GA.md."
review-gaps: |-
  note: --dry-run of the same move does not show the warning (reportDryRun in internal/cli/flow.go never calls untickedDoD); a dry run is where an agent could learn it before writing - add the line there if wanted
  note: the TUI move paths (internal/tui/model.go, signoff.go) ignore DoDUnticked - deliberate per the decision note (the warning is for agents on the CLI)
  note: pre-existing: Move reads from=t.Status before the locked Mutate, so two concurrent moves out of in-progress both warn (and the Moved hook gets the same stale from) - duplicate line only, no wrong write
test-verdict: "pass: go vet ./... and go test ./core/... ./internal/cli/ green (rc=0), -race on core/move+internal/cli green; DoD 1 verified: scratch board, in-progress->review with 0 of 2 ticked prints the warning on stderr, --json stdout is clean JSON, rc=0, move landed; send-back and parking print nothing; refused move (human without question) rc=3, no warning, nothing written; two parallel moves both land, each warns once (note)"
hunt-verdict: "pass: bughunt done at 791c8e1e4, 0 change, 0 pre-existing; stale 'from' under concurrent moves already noted by testing"
review-verdict: "accept: DoD 1 is met - the warning appears on leaving in-progress with nothing ticked, the move lands, --json keeps it off stdout, both cases are tested; no defect found in core/move or the CLI move path"
review-round: 1
review-verdict-gaps: |-
  note: --dry-run of the same move does not show the warning (reportDryRun in internal/cli/flow.go never calls untickedDoD); a dry run is where an agent could learn it before writing - add the line there if wanted
  note: the TUI move paths (internal/tui/model.go, signoff.go) ignore DoDUnticked - deliberate per the decision note (the warning is for agents on the CLI)
  note: pre-existing: Move reads from=t.Status before the locked Mutate, so two concurrent moves out of in-progress both warn (and the Moved hook gets the same stale from) - duplicate line only, no wrong write
review-check: "1. cd into any repo with a jaira board and run: jaira create \"dod probe\" --goal g --context c --dod \"a\" --dod \"b\"  (note the handle it prints, call it H)  2. jaira claim H; jaira move H --to todo; jaira move H --to in-progress  3. jaira move H --to critique --what w --why y --resolves r  4. You see on stderr: \"warning: Definition of Done: 0 of 2 ticked - the terminal lane reads this list, not the plan.\" and the ticket is in critique  5. Repeat with --json added to the move (send it back to in-progress first): the warning still shows, stdout is only the JSON  6. jaira dod H 1 --done --proof x, move again: no warning  7. jaira reject H --why probe --force to remove the probe. Automated: go test ./core/move -run UntickedDoD and go test ./internal/cli -run UntickedDoD"
---

# Der erste Schritt aus der Implementierung warnt, wenn kein DoD-Punkt abgehakt ist

## Definition of Done

- [x] Ein Zug aus der implementierenden Lane heraus schreibt eine Warnung, wenn kein Punkt der Definition of Done abgehakt ist; der Zug selbst geht durch; in --json steht sie nicht auf stdout zwischen den Daten; ein Test deckt beide Faelle ab (kein Punkt abgehakt: Warnung; mindestens einer: keine)
  proof: core/move/move.go untickedDoD; internal/cli/flow.go warning on stderr; tests TestMoveReportsUntickedDoDOutOfImplementing, TestMoveOutOfImplementingWarnsOnUntickedDoD

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 18:49 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 18:50 · Alexander Sacharov** — decision: no open design question. DoD fixes the place: a move OUT of the implementing lane, i.e. the lane whose output-produces names outcome-what (in-progress on the shipped board), forward only (later lane index) and not into a parking lane (requires-blocked-reason, e.g. blocked) and not back. Detected in core/move (one place for every caller, field Result.DoDUnticked), printed by the CLI on stderr in both text and --json mode. 'Ticked' = settled (done or superseded): an all-superseded list passes the gate, so warning there would be false. A ticket with no checklist items gets no warning - nothing to tick, the gate says that itself. TUI prints nothing: the warning is for agents, who use the CLI.
- **2026-11-22 18:53 · Alexander Sacharov** — critique: none. Weighed computing the warning in internal/cli/flow.go alone (one consumer today); kept it in core/move because its package doc makes it the one place for after-a-move rules and Result already carries caller-specific facts (Holds, Trimmed). Lane order via Set.Index is the display order, where blocked sits after done - hence the explicit RequiresBlockedReason exclusion.
- **2026-11-22 18:54 · Alexander Sacharov** — testing: note: two concurrent 'move --to review' from in-progress both land and both print the warning - Move reads 'from' from the ticket loaded before the locked Mutate, so the second writer still sees in-progress. Pre-existing for 'from' (the Moved plugin hook gets the same stale value); here it only means a duplicate warning line, no wrong write.
- **2026-11-22 18:56 · Alexander Sacharov** — bughunt: done at 791c8e1e4 - 0 findings in 5 modules, first hunt; concurrency, error paths, input, permissions checked; files: AGENTS.md CLAUDE.md core/board/announce.go core/move/move.go core/move/move_test.go core/release/NOTES.md core/release/notes/77B4GA.md internal/cli/flow.go internal/cli/untickeddod_test.go
- **2026-11-22 19:26 · Alexander Sacharov** — Accepted from signoff: angenommen:
- note: --dry-run of the same move does not show the warning (reportDryRun in internal/cli/flow.go never calls untickedDoD); a dry run is where an agent could learn it before writing - add the line there if wanted
- note: the TUI move paths (internal/tui/model.go, signoff.go) ignore DoDUnticked - deliberate per the decision note (the warning is for agents on the CLI)
- note: pre-existing: Move reads from=t.Status before the locked Mutate, so two concurrent moves out of in-progress both warn (and the Moved hook gets the same stale from) - duplicate line only, no wrong write
