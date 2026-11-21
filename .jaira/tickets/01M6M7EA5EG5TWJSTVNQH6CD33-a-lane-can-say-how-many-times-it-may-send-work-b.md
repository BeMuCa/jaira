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
updated-at: 2026-11-21T19:08:59Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2282179
claimed-at: 2026-11-21T19:01:08Z
outcome-what: "A lane file may declare max-rounds: <n>; the gate refuses that lane's (n+1)-th send-back in a row (exit 3, code max_rounds) naming the human lane; --for-lane and lanes show carry max_rounds/rounds_left; the shipped critique, optimize and testing lanes declare 3."
outcome-why: "The three-returns stop lived only in the worker role's prose, so an agent that did not read it looped forever and a board could not pick another limit."
outcome-resolves: "DoD 1-5: parse refusal (TestParseMaxRounds), gate refusal + rounds_left (internal/cli/maxrounds_test.go, core/gate/maxrounds_test.go), README + core/release/notes/H6CD33.md, edge table (person, other lane, force, parking, no limit), shipped lanes. Checked: concurrency (gate reads, move mutates — same race as every gate), error paths (bad value refused, lane skipped with warning), unusual input (empty, list, quoted, huge, negative), permissions (a faked terminal escapes only what --force already does)."
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
  proof: TestMaxRoundsRefusesOnlyALanesOwnRunPastTheLimit (core/gate/maxrounds_test.go), TestMaxRoundsLetsAPersonSendItBack
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
