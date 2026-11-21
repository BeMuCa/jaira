---
id: 01M6M7EA5EG5TWJSTVNQH6CD33
title: A lane can say how many times it may send work back before a person decides
status: critique
ready: true
creator: berkc
goal: "A loop lane declares max-rounds: <n>, and once a ticket has been sent back from it n times the next send-back goes to a person instead of to in-progress."
context: |-
  What is wrong today: the limit of three send-backs per loop lane (critique, optimize, testing) lives only in the prose of the worker role (core/role/builtin/jaira-worker/SKILL.md, 'the same lane sent work back three times'). The binary does not know it.
  So an agent that does not read that skill loops forever, and a board cannot choose a different limit for one lane.
  Triggered by: ticket ENCFVZ (skill/script/receipt lanes), decision of the maintainer on 04.11.2026 listing max-rounds as a lane attribute; kept out of ENCFVZ's scope on purpose.
  Already there: the ticket field 'round' counts send-backs (core/move/move.go, sendsBack); review-round stamps the round a verdict was written in.
  Open: whether the gate refuses the n+1-th send-back (exit 3, name the human lane) or only 'jaira show --for-lane' reports rounds left; and whether round counts per lane or per ticket (today it is per ticket).
definition-of-done: "A lane file may declare max-rounds: <n>; the loader refuses a non-positive or non-numeric value"
tags:
  - gates
blocked-by:
  - 01M7WWZWSDHK37DGGHJFDF0YD6
related:
  - 01M7WWZWSDHK37DGGHJFDF0YD6
follows: 01M67K7GEZF694985EHTENCFVZ
commits: []
created-at: 2026-11-04T21:29:03Z
updated-at: 2026-11-21T19:18:07Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2282179
claimed-at: 2026-11-21T19:01:08Z
outcome-what: "A lane file may declare max-rounds: <n>; the gate refuses that lane's (n+1)-th send-back in a row (exit 3, code max_rounds) naming the human lane; --for-lane and lanes show carry max_rounds/rounds_left; the shipped critique, optimize and testing lanes declare 3."
outcome-why: "The three-returns stop lived only in the worker role's prose, so an agent that did not read it looped forever and a board could not pick another limit."
outcome-resolves: "DoD 1-5: parse refusal (TestParseMaxRounds), gate refusal + rounds_left (internal/cli/maxrounds_test.go, core/gate/maxrounds_test.go), README + core/release/notes/H6CD33.md, edge table (person, other lane, force, parking, no limit), shipped lanes. Checked: concurrency (gate reads, move mutates — same race as every gate), error paths (bad value refused, lane skipped with warning), unusual input (empty, list, quoted, huge, negative), permissions (a faked terminal escapes only what --force already does)."
review-summary: "none"
round: 4
back-from: testing
back-by: lane
back-count: 2
review-gaps: "none — the fix is one condition at the cause (gate.go) and one test; nothing to fold or remove"
test-verdict: "pass: go test -race ./... RC=0, go vet, shipped-history --check, windows build; DoD 1-5 verified; scratch board: rounds_left 3→0, 4th send-back exit 3, human with --question lands, bad value → lane skipped, parallel send-backs consistent"
hunt-verdict: "fail: the max-rounds exemption keys on requires-question alone, so an agentic question lane escapes the limit (bughunt change, gate.go:419)"
---

# A lane can say how many times it may send work back before a person decides

## Definition of Done

- [x] A lane file may declare max-rounds: <n>; the loader refuses a non-positive or non-numeric value
  proof: core/lane/lane.go parse max-rounds; TestParseMaxRounds (core/lane/rounds_test.go)
- [x] Sending work back from that lane for the (n+1)-th time is refused with a message naming the human lane, and --for-lane shows the rounds left
  proof: core/gate/gate.go CodeMaxRounds; TestMaxRoundsRefusesTheFourthSendBack (internal/cli/maxrounds_test.go), TestMaxRoundsRefusalNamesThePersonsLane; rounds_left in internal/cli/flow.go for-lane payload
- [x] README lane format and a release note describe it
  proof: README.md 'A loop lane can stop itself'; core/release/notes/H6CD33.md
- [x] Edge classes are refused or passed on purpose, table-tested: a send-back by a person (board or terminal) is never refused and resets the count; one from another lane resets it; --force overrides; parking (blocked) is no send-back; a lane without max-rounds behaves as today
  proof: TestMaxRoundsRefusesOnlyALanesOwnRunPastTheLimit, TestMaxRoundsLetsTheHandOffToAPersonsLaneThrough (core/gate/maxrounds_test.go), TestPersonLaneSkipsAgenticQuestionLanes, TestMaxRoundsLetsAPersonSendItBack
- [x] The shipped critique, optimize and testing lanes declare max-rounds: 3, and a project lane file that differs only in max-rounds counts as an override
  proof: core/lane/builtin/25-critique.md:10, 26-optimize.md:10, 27-testing.md:10; TestShippedLoopLanesStopAfterThree, TestMaxRoundsCountsForOverride

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] lane: MaxRounds field, parse with refusal, lanesEquivalent
- [x] gate: ByPerson in Request, CodeMaxRounds check on send-back, human-lane name
- [x] move/cli/tui: pass ByPerson into the gate request
- [x] for-lane: max_rounds and rounds_left in JSON and text
- [x] builtin lanes max-rounds: 3; README; release note; run gates

## Progress
- **2026-11-20 16:48 · Alexander Sacharov** — Interim-Regeln (Alex 20.11.2026, gelten ab sofort bis TW5FE2 ausgeliefert ist):
1. pre-process: die DoD nennt die Fehlerklasse statt des Symptoms, die Eingabeklassen (leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg) und bei einer Heuristik eine Annahmeregel - als DoD-Punkte (jaira dod <id> --add).
2. in-progress: vor dem Move selbst die vier Kategorien pruefen (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte). Tests ueber die Klasse (Tabellentest), Geschwister-Stellen per grep suchen und mitziehen, Fix an der Ursache statt am Beispiel. Nicht nur Haken setzen: wer umsetzt, denkt die Raender selbst.
3. bughunt/testing: ein Fund nennt Klasse und Beispiel; der Rueckweg verlangt den Fix der Klasse, nicht des einen Beispiels.
4. Schleifenbremse: Blocker sind nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke. Nach einem Fix-Durchgang laeuft der Testlauf und ein Bughunt nur ueber den Diff seit dem letzten Hunt, dann review. Nach 3 Rueckgaben aus derselben Lane: Frage an den Teamlead statt weiter schleifen.
5. Kleine Raender aus review-gaps in derselben Runde schliessen statt als note: stehen lassen - Alex schickt sie sonst ohnehin zurueck. Vorbestehendes ausserhalb des Umfangs als 'note: pre-existing: ...'.
- **2026-11-21 18:58 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/4205.sock
- **2026-11-21 19:01 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2274544.sock
- **2026-11-21 19:02 · Alexander Sacharov** — decision (worker, no open question for the person): (1) the limit counts back-count — returns in a row from that lane, as DF0YD6 records them; a person's send-back or one from another lane starts it again, same rule as the worker's third-return stop. (2) the gate refuses the (n+1)-th send-back (DoD 2), it does not redirect; --force overrides like every refusal; a move made by a person (board or terminal, ByPerson/Interactive) is never refused — the person is the one deciding. (3) the shipped loop lanes critique, optimize, testing get max-rounds: 3, the number the worker role already uses, so the binary knows the rule without a board edit. (4) the human lane named is the first rejects-to lane that takes a question, else the board's first question lane. (5) --for-lane carries max_rounds and rounds_left only for a lane that declares max-rounds.
- **2026-11-21 19:08 · Alexander Sacharov** — found out: sendsBack lived in core/move only; moved it to lane.Set.SendsBack so the gate and the round counter share one definition. ByPerson (terminal) now also reaches the gate — TUI pre-check and CLI --dry-run pass it too, else they would predict a refusal the real move does not make. A lane file with a bad max-rounds is skipped with a warning like any other unparsable lane file (readLaneDir), not a board failure. This repo's own .jaira/lanes copies carry no max-rounds until 'jaira update --lanes' — not run here, that is the maintainer's call.
- **2026-11-21 19:09 · Alexander Sacharov** — critique: drop PersonLane's third fallback to a requires-human-exit lane (core/lane/rounds.go) — signoff is for reviewed work, not for an undecided loop; no question lane → generic wording. Everything else let stand: one gate check, back-count reused, SendsBack shared.
- **2026-11-21 19:11 · Alexander Sacharov** — testing: note: scratch board (built binary, lanes add testing): rounds_left 3→2→1→0, 4th send-back exit 3 with the human-lane message; move to human with --question lands once the lane's own outputs are written (with test-verdict/hunt-verdict empty the move is refused for those — expected, a lane writes its verdict before it hands on); max-rounds: zero → warning, lane skipped, lanes exit 0; two parallel send-backs at count 2 → one counted (count 3, round 4), no conflict markers.
- **2026-11-21 19:14 · Alexander Sacharov** — bughunt: change: core/gate/gate.go:416 (with core/lane/rounds.go:45 PersonLane) - class: the max_rounds check refuses every send-back from the exhausted loop lane, including the move into the person's lane that its own refusal message prescribes; whenever that question lane sits earlier on the board than the loop lane (rejects-to may name any installed lane, and the template's own example is rejects-to: [in-progress, human]), the hand-off is itself a send-back and is refused with the same max_rounds code, so the only way out is --force, which an agent is not to pass. example: board with question lane 'backlog' before critique, critique rejects-to [in-progress, backlog], back-count 3 by critique -> move critique->backlog --question 'is the dod wrong?' is refused max_rounds naming 'backlog' again. Not a blocker: on the shipped board and this repo's board 'human' sits after critique/optimize/testing, so the hand-off there is a forward move and passes. Fix direction: exempt a target that takes a question (or isPersons) from the check. ran: scratch go test in core/gate (env.Lanes with backlog.RequiresQuestion=true, human false) logged the max_rounds violation for the move to backlog; blame: b86e9b78 (names H6CD33)
- **2026-11-21 19:14 · Alexander Sacharov** — bughunt: done at c468f6a4 - 1 finding (1 change, 0 pre-existing) in 9 modules, first hunt; concurrency, error paths, input, permissions checked; files: README.md core/board/announce.go core/gate/gate.go core/gate/maxrounds_test.go core/lane/builtin/20-in-progress.md core/lane/builtin/25-critique.md core/lane/builtin/26-optimize.md core/lane/builtin/27-testing.md core/lane/facts.go core/lane/lane.go core/lane/rounds.go core/lane/rounds_test.go core/lane/template.go core/move/move.go core/release/notes/H6CD33.md core/shipped/shipped.sha256 internal/cli/flow.go internal/cli/maxrounds_test.go internal/cli/tickets.go internal/tui/model.go
- **2026-11-21 19:14 · Alexander Sacharov** — testing: class: the max-rounds check refuses every send-back, including the hand-off its own message prescribes, when the question lane is earlier on the board; example: rejects-to [in-progress, backlog] with backlog requires-question; fix: at the cause in gate.go — a target lane that takes a question is a person deciding, never a round; table row for it.
- **2026-11-21 19:16 · Alexander Sacharov** — bughunt: change: core/gate/gate.go:419 - class: the max-rounds exemption keys on the target lane taking a question (target.RequiresQuestion), not on it being a lane a person works (non-agentic), so any agentic lane with requires-question that a loop lane sends back to escapes the limit for good, and a question already on the ticket (gate.go:282 accepts t.Question) satisfies it without a new one; the new test core/gate/maxrounds_test.go:93 even asserts this for in-progress; example: in-progress with requires-question: true, critique at back-count 3..9 moves to in-progress with no --question and is never refused; not a blocker on the shipped board (its only question lane, human, is agentic: false), an edge on boards that mark a work lane requires-question; ran: scratch test in core/gate (CheckAdvance, back-count 3-9, stale t.Question, no CodeMaxRounds)
- **2026-11-21 19:17 · Alexander Sacharov** — bughunt: done at f4648d8b - 1 finding (1 change, 0 pre-existing) in 2 modules, follow-up hunt (diff since c468f6a4: README.md core/gate/gate.go core/gate/maxrounds_test.go; no new files); concurrency, error paths, input, permissions checked; files: README.md core/board/announce.go core/gate/gate.go core/gate/maxrounds_test.go core/lane/builtin/20-in-progress.md core/lane/builtin/25-critique.md core/lane/builtin/26-optimize.md core/lane/builtin/27-testing.md core/lane/facts.go core/lane/lane.go core/lane/rounds.go core/lane/rounds_test.go core/lane/template.go core/move/move.go core/release/notes/H6CD33.md core/shipped/shipped.sha256 internal/cli/flow.go internal/cli/maxrounds_test.go internal/cli/tickets.go internal/tui/model.go
- **2026-11-21 19:17 · Alexander Sacharov** — testing: class: the exemption names 'takes a question' where it means 'a person answers there'; example: in-progress with requires-question escapes max-rounds for good; fix: one predicate (lane takes a question and no agent works it) for both the gate exemption and PersonLane; the test row for in-progress flips to refused.
