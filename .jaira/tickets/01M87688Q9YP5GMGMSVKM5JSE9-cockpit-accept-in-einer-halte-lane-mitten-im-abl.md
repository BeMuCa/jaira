---
id: 01M87688Q9YP5GMGMSVKM5JSE9
title: "Cockpit: Accept in einer Halte-Lane mitten im Ablauf schiebt das Ticket in die nächste Lane"
status: in-progress
ready: true
creator: Sascha
assignee: Sascha
goal: "Eine Person kann ein Ticket in jeder Personen-Lane mit Accept freigeben: in einer Lane vor dem Ende geht es in next_lane weiter, nicht nur von signoff nach done."
context: |-
  Was falsch ist: Accept im Cockpit-Pane 'Wartet auf dich' funktioniert nur, wenn die nächste Lane terminal ist (done).
  In jeder anderen Personen-Lane kommt nur: '<ID> waits in <lane>, which does not end the work: nothing to accept there'.
  Auslöser 24.11.2026: Board 'video' im Repo demo-video, Ticket in Lane script-approval (Halt Skript), next_lane = build. Sascha wollte das Skript freigeben und konnte es nicht; nur 'jaira move --to build' von Hand ging.
  Betrifft auch die Lane 'human' auf dem default-Board, wenn die Frage beantwortet ist.
  Code: core/mod/builtin/jaira-cockpit/hooks/board.ts:158 (acceptTo nur bei next.terminal), :1182 (step notTerminal), hooks/register.tsx:1502, hooks/texts.ts:80. Test, der das heutige Verhalten festschreibt: tests/board.test.ts:854.
  Offen zu klären im Plan: ob jaira move aus einer Personen-Lane ohne --force geht, wenn eine Person es auslöst (Cockpit-Button = Person).
  Geplant danach: die Entscheidungsseite (Skill decision-page, Markierungen 'accept' zweier Personen) soll denselben Weg nutzen, um Tickets aus Personen-Lanes weiterzuschieben.
definition-of-done: In einer Personen-Lane mit nicht-terminaler next_lane bietet das Cockpit Accept an; Accept verschiebt das Ticket nach next_lane
tags:
  - cockpit
  - lanes
blocked-by: []
related: []
commits: []
created-at: 2026-11-24T16:29:33Z
updated-at: 2026-11-24T17:00:27Z
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-51795
claimed-at: 2026-11-24T16:30:00Z
outcome-what: "Cockpit Accept moves a ticket from any person's lane one lane on (script-approval -> build, human -> review), not only signoff -> done; the Execute line names the target mid-board"
outcome-why: "Sascha could not release a script in the video board's script-approval lane from the cockpit; only a hand-typed jaira move went"
outcome-resolves: M5JSE9
review-summary: "Full review at 67b11995a (first review). Cockpit Accept/Execute now moves a ticket from any person's lane one lane on (snapshot acceptTo = next_lane, new acceptEnds = whether that lane is terminal); a PR/landing and the signoff hold run only where it ends the work; the pane line says 'Execute: accept -> <lane>' mid-board and 'noNext' replaces notTerminal where jaira names no next lane. A question still open in a lane that does not end the work (human) stops Execute (needAnswer) until answered; an Explain option marked [accept] re-checks the ticket's lane (stillIn) and clears the question before moving. Cockpit suite 581 pass."
round: 4
back-from: review
back-by: lane
back-count: 1
review-gaps: |-
  blocker: failed gate - go test ./core/role fails TestEveryInstalledFileIsAShippedVersion: core/shipped/shipped.sha256 lacks the new hashes of hooks/board.ts, hooks/register.tsx, hooks/texts.ts (CI runs go test ./...); fix: scripts/shipped-history.sh --write (or append the 3 lines the test prints) and commit with M5JSE9
  note: edge: an Explain option with no action writes decision: but leaves the question set, so Execute keeps answering needAnswer although a decision was noted; either clear the question there too or name that case in t.needAnswer
  note: limit: a gate that refuses the accept after deciding cleared the question leaves the ticket in human with its decision note and no question (decided in the critique note, no rollback)
  note: elsewhere: jaira move never clears question when a ticket leaves human (internal/cli/flow.go:139) - a ticket moved on by CLI carries a stale question into a later person lane, where Execute then says needAnswer; new ticket
test-verdict: "pass: round 2 over 1b49fe6..67b1199 (register.tsx deciding, texts.ts, board.test.ts): claude plugin test 581 pass 0 fail RC=0; DoD 7 verified by its two tests (question= set then move --to review --force; a ticket moved on meanwhile: no move); DoD 1-6 tests still green"
hunt-verdict: "pass: bughunt done at 67b11995a, 0 change, 0 pre-existing (round 2 over the fix)"
review-verdict: "send-back: failed gate - core/shipped/shipped.sha256 misses the three changed cockpit files, go test ./core/role is red"
review-round: 3
review-verdict-gaps: |-
  blocker: failed gate - go test ./core/role fails TestEveryInstalledFileIsAShippedVersion: core/shipped/shipped.sha256 lacks the new hashes of hooks/board.ts, hooks/register.tsx, hooks/texts.ts (CI runs go test ./...); fix: scripts/shipped-history.sh --write (or append the 3 lines the test prints) and commit with M5JSE9
  note: edge: an Explain option with no action writes decision: but leaves the question set, so Execute keeps answering needAnswer although a decision was noted; either clear the question there too or name that case in t.needAnswer
  note: limit: a gate that refuses the accept after deciding cleared the question leaves the ticket in human with its decision note and no question (decided in the critique note, no rollback)
  note: elsewhere: jaira move never clears question when a ticket leaves human (internal/cli/flow.go:139) - a ticket moved on by CLI carries a stale question into a later person lane, where Execute then says needAnswer; new ticket
review-check: "1. cd core/mod/builtin/jaira-cockpit && claude plugin test .  -> 581 pass  2. go test ./core/role/  -> ok (today FAIL: shipped.sha256 lacks 3 files)  3. on a scratch board put a ticket in a person's lane whose next lane is not done (human with question cleared: jaira set <id> question=)  4. open the cockpit band, Acceptance on that ticket  -> last line reads 'Execute: accept → review'  5. press Execute  -> ticket is in review, note 'cockpit: accepted by the person, human -> review', no pull request started  6. same with the question still set  -> Execute moves nothing and says 'answer it first'  7. a signoff ticket still says 'Execute: done' and goes to done"
---

# Cockpit: Accept in einer Halte-Lane mitten im Ablauf schiebt das Ticket in die nächste Lane

## Definition of Done

- [x] In einer Personen-Lane mit nicht-terminaler next_lane bietet das Cockpit Accept an; Accept verschiebt das Ticket nach next_lane
  proof: core/mod/builtin/jaira-cockpit/hooks/board.ts:162; test 'Execute moves it one lane on and says where'
- [x] Der Text des Buttons/der Zusammenfassung nennt das Ziel (z. B. 'Accept → build'), damit klar ist, dass das Ticket nicht nach done geht
  proof: core/mod/builtin/jaira-cockpit/hooks/texts.ts:73 ('Execute: accept → review' asserted in that test)
- [x] signoff → done verhält sich unverändert; Give back (note fix / getippter Grund) unverändert
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:731,760 gate PR/hold on acceptEnds; signoff tests ('Execute: done', move --to done --force) and give-back tests unchanged and green
- [x] Tests: board.test.ts deckt beide Fälle ab (terminal und nicht-terminal); der alte notTerminal-Test ist angepasst
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: 'Execute moves it one lane on', 'names no lane after it', executionOf table rows 'in a lane that does not end' / 'with no next lane'; claude plugin test 525 pass
- [x] Release-Note core/release/notes/<ID>.md
  proof: core/release/notes/M5JSE9.md
- [x] Eine offene Frage in einer Lane, die die Arbeit nicht beendet (human), hält Accept an: Execute verschiebt nichts und verlangt die Antwort; getippter Text gibt weiterhin zurück
  proof: core/mod/builtin/jaira-cockpit/hooks/board.ts:1191 waitsForAnswer + step needAnswer; test 'a question still open in a lane that does not end the work: Execute moves nothing and asks for the answer'
- [x] Der Entscheidungsweg (Explain-Option [accept]) bewegt nur ein Ticket, das noch in seiner Lane steht, und löscht vor dem Weiterschieben aus einer nicht-terminalen Lane die beantwortete Frage
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1376 deciding; tests 'an option marked [accept] on a ticket waiting on a question clears the question and moves it one lane on', 'an option marked [accept] on a ticket that moved on meanwhile moves nothing'

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-24 16:30 · Sascha** — parent session: uds:/run/user/1000/cc-socks/3417.sock
- **2026-11-24 16:30 · Sascha** — no open decisions: --force question settled from code - accepting() (register.tsx) dry-runs jaira move, otherRefusal drops only needs_human, then --force; the same path serves a mid-flow person lane. Plan: Waiting gets acceptEnds; acceptTo = next_lane whenever it exists; afterAccept (PR/land) and heldInSignoff only when it ends; summary 'Execute: accept -> <lane>' mid-flow, 'Execute: done' unchanged
- **2026-11-24 16:36 · Sascha** — critique: (1) an open question in a mid-board lane must stop Accept - jaira does not clear the question on leaving human, so review would get it unanswered; new step needAnswer. (2) summary's ends default is speculative, make it required
- **2026-11-24 16:39 · Sascha** — critique pass 2: both findings addressed (needAnswer step via waitsForAnswer, shared with lightOf; summary ends required) - nothing left
- **2026-11-24 16:39 · Sascha** — waiting: bughunt (background subagent) since 17:39
- **2026-11-24 16:42 · Sascha** — bughunt: change: blocker: core/mod/builtin/jaira-cockpit/hooks/board.ts:1191 (waitsForAnswer), used at register.tsx:1523, 2591, 2693 - class: 'answered' is read off the ticket's question field, which jaira only ever sets and never clears (a requires_question lane cannot even be entered without one, core/gate/gate.go:325; move only SetScalar's it, internal/cli/flow.go:139), so every ticket in a human lane waits for an answer forever; example: ticket in human with question 'keep the old key?', the person answers (decision: note via the Answer pane), Execute still refuses with needAnswer 'answer it first, then Execute moves it on to review' - a path that never opens, so the release note's 'an answered human question' case is dead and only an Explain option marked [accept] (deciding, register.tsx:1373) moves it; the decision lives in decidedAtom / explanationNotes().decided, which waitsForAnswer does not consult; ran: scratch board, jaira move <id> --to human --question ... then jaira note <id> 'decision: yes', jaira list --json still has question set
- **2026-11-24 16:42 · Sascha** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/board.ts:162 (acceptTo = next_lane) with register.tsx:699-712 accept()/accepting() - class: Accept moves to the lane the snapshot named next without checking the ticket is still in w.lane; since acceptTo is no longer only the terminal lane it can now lie behind the ticket's real lane, and a move backwards is a send-back that the dry-run's needs_human gate never refuses, then --force; example: board shows the ticket in human (acceptTo review) while another session already moved it on to signoff, the person presses Execute/an [accept] option: the ticket goes signoff -> review and gets the note 'accepted by the person, human -> review'; movedBefore only catches 'already in <to>'; before this change acceptTo was always the terminal lane, which nothing is behind; ran: scratch board, ticket in signoff, jaira move <id> --to review --dry-run --json (exit 0), --force (exit 0, status review)
- **2026-11-24 16:42 · Sascha** — hunt report: 2 findings, both change (1 blocker: human-lane Execute never accepts because question is never cleared; 1 edge: stale snapshot turns Accept into a silent send-back), notes on M5JSE9 | read whole: the accept/Execute/decide paths of hooks/register.tsx (accept, accepting, heldInSignoff, ticketHere, deciding, execute, Acceptance pane), hooks/board.ts (snapshotOf, isPersons, lightOf, otherRefusal, executionOf, waitsForAnswer, explanation parsing), hooks/texts.ts, tests/board.test.ts around executionOf, core/release/notes/M5JSE9.md, plus jaira's gate.go/flow.go question handling; first hunt | could not run the cockpit test suite (bun not installed in this shell); the rest of register.tsx (3005 lines) unrelated to acceptTo was not read line by line
- **2026-11-24 16:42 · Sascha** — bughunt: done at 1b49fe6f3 - 2 findings (2 change, 0 pre-existing) in 5 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/M5JSE9.md
- **2026-11-24 16:43 · Sascha** — testing: class: the decide path (deciding -> accept, register.tsx) skips the guards Execute has; example 1: an Explain option marked [accept] on a human ticket writes 'decision:' and moves it to review with the question still set (Execute's needAnswer is bypassed; Discuss is fine, its prompt clears the question - texts.ts discussInThisChat); example 2: a stale band press sends a ticket that already moved on (signoff) back to review, jaira allows it. fix: in deciding's accept branch run stillIn first (as acceptWith does) and, where waitsForAnswer(w), clear the question (the decision note is its answer) before accept; needAnswer's refusal names Discuss as the way to answer
- **2026-11-24 16:49 · Sascha** — critique+optimize round 2 (one pass, fix 34 lines, no lock/parser): stillIn guard reused from acceptWith, waitsForAnswer reused; let stand: a gate that refuses the accept after the question was cleared leaves the ticket in human with its decision note and no question - the decision is the answer, so Execute then works, no rollback needed
- **2026-11-24 16:50 · Sascha** — waiting: bughunt round 2 (background subagent) since 17:50
- **2026-11-24 16:52 · Sascha** — hunt report: 0 findings, nothing written beyond this report | read as a diff since 1b49fe6f3: core/mod/builtin/jaira-cockpit/hooks/register.tsx (deciding with ticketHere, accept, accepting, stillIn), hooks/texts.ts, tests/board.test.ts; board.ts and core/release/notes/M5JSE9.md unchanged since then | checked in a scratch board: jaira set <id> question= clears the question (exit 0) and human -> review --dry-run is then allowed; claude plugin test . 581 pass 0 fail; the nested pull inside accept (ticketHere with carriedTo) was followed by reading only, not run
- **2026-11-24 16:52 · Sascha** — bughunt: done at 67b11995a - 0 findings (0 change, 0 pre-existing) in 5 modules, follow-up hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/M5JSE9.md
- **2026-11-24 16:52 · Sascha** — waiting: review (fresh agent, Herdr tab via run-lane.sh) since 17:52
- **2026-11-24 16:56 · Sascha** — review: done at 67b11995a
- **2026-11-24 16:56 · Sascha** — review: blocker: go test ./core/role fails TestEveryInstalledFileIsAShippedVersion - core/shipped/shipped.sha256 lacks the new hashes of jaira-cockpit hooks/board.ts, register.tsx, texts.ts; run scripts/shipped-history.sh --write (or append the 3 printed lines), commit with M5JSE9, then back through the lanes
- **2026-11-24 17:00 · Sascha** — review round 1 blocker fixed: shipped.sha256 gets the 3 hashes go test ./core/role asked for (appended by hand like CSQXKR, not --write, which also added 8 intermediate branch versions). Full go test ./... green. My miss: 'Before you push' gate list did not include go test ./core/role for a cockpit-only change
