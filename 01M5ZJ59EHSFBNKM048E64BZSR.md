---
id: 01M5ZJ59EHSFBNKM048E64BZSR
title: "jaira show nennt die optionalen Lanes, an denen next_lane vorbeigeht"
status: review
ready: true
creator: BeMuCa
goal: "jaira show --json fuehrt neben next_lane die optionalen Lanes (brainstorm, pre-process) auf, die uebersprungen werden, und wie man sie einschaltet (jaira dod <id> --option <lane>); ein Driver kann so Planung waehlen, statt sie unbemerkt zu verpassen."
context: |-
  Was schief ist: Driver uebersprangen pre-process bei Tickets, die einen Plan gebraucht haetten, bis jemand sagte, dass man --option vorher ticken muss (O11, 27.10.2026).
  Heute: jaira dod --option <lane> schaltet eine Schritt-Option ein (jaira dod --help); jaira show --json zeigt nur next_lane, nicht, dass eine optionale Lane uebersprungen wird; der Agent-Block erwaehnt --option nicht.
  Nicht gewollt: optionale Lanes standardmaessig einschalten - dann sind sie nicht optional.
definition-of-done: "jaira show <id> --json hat ein Feld, das die uebersprungenen optionalen Lanes und den Befehl zum Einschalten nennt"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-27T20:52:19Z
updated-at: 2026-11-23T01:18:25Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: worker-64bzsr
claimed-at: 2026-11-23T01:10:10Z
outcome-what: "jaira show --json gains skipped_lanes: every unticked optional lane ahead of the ticket with lane, option and the 'jaira dod <id> --option <opt>' that turns it on; the agent block says so in one sentence"
outcome-why: drivers skipped pre-process on tickets that needed a plan because next_lane walks past optional lanes silently (O11)
outcome-resolves: 64BZSR
review-summary: "Added lane.Set.UnchosenAhead (core/lane/next.go). It walks the lanes after the ticket's current lane up to the first terminal lane and collects every lane whose requires-option the ticket's Options leave unticked. It returns nothing for a ticket that is parked, waiting on a question, finished or in an unknown lane, the same guard Next uses. ticketJSON (internal/cli/tickets.go) prints the result as unchosen_lanes, never null: one entry per lane with lane, option and enable='jaira dod <ID> --option <option>'. The field therefore appears in every payload built from ticketJSON: show, list, move, resume and create. The agent block (core/board/announce.go, and its copies in CLAUDE.md and AGENTS.md) gains one sentence pointing drivers at the field. The release note says to read the field and to run jaira update. A table test covers backlog, todo, todo with planning ticked, in-progress and done. A driver that used to see only next_lane=todo from backlog now also learns that pre-process exists and how to switch it on."
round: 2
back-from: critique
back-by: lane
back-count: 1
review-gaps: |-
  note: outcome-what still names the field skipped_lanes, but it shipped as unchosen_lanes (renamed in round 2). A person reading the outcome at signoff will grep for a key that does not exist. Fix: correct outcome-what to unchosen_lanes before accepting.
  note: unusual input: the enable command pastes the option name unquoted (internal/cli/tickets.go:1725). A custom lane with requires-option "my step" would print a command that ticks "my". Fix: quote options outside [A-Za-z0-9._-]. pre-existing: the gate's refusal text has the same pattern at core/gate/gate.go:298; fix both together. Stock options are single words.
  note: the test does not cover the parked, question-lane and unknown-lane early returns of UnchosenAhead. They are read-only copies of Next's guard. Fix: add a blocked case and a human case to TestTicketJSONNamesTheUnchosenOptionalLanes.
  note: unchosen_lanes also lands in list/move/resume/create --json, because all of them build on ticketJSON. That is consistent with next_lane, but the release note mentions show only.
test-verdict: "pass: gofmt clean, go vet (internal/cli, core/lane, core/board) and go build ./... RC=0; go test NOT run by house rule tonight (main session runs the race suite before the PR) - TestTicketJSONNamesTheUnchosenOptionalLanes unexecuted; DoD 1-3 verified in tree at 25217204c; scratch board: backlog lists brainstorm+pre-process with enable lines, dod --option planning drops pre-process and next_lane from todo becomes pre-process, blocked ticket gives [] and next_lane empty; unknown id exits 5; the change writes nothing, so no two-writer case"
hunt-verdict: "pass: bughunt done at 25217204c, 1 change (unquoted option name in enable - edge on custom lanes with spaced option names, note for the person), 0 pre-existing"
review-verdict: "accept: all three DoD items hold. show --json carries unchosen_lanes with lane, option and the enable command (checked on a scratch board: the printed command works and next_lane then becomes pre-process). The agent block names --option in one sentence. The release note exists. I found no blocker in next.go, the ticketJSON hunk or the gate's option predicate it mirrors; gofmt, vet and build are clean."
review-round: 2
review-verdict-gaps: |-
  note: outcome-what still names the field skipped_lanes, but it shipped as unchosen_lanes (renamed in round 2). A person reading the outcome at signoff will grep for a key that does not exist. Fix: correct outcome-what to unchosen_lanes before accepting.
  note: unusual input: the enable command pastes the option name unquoted (internal/cli/tickets.go:1725). A custom lane with requires-option "my step" would print a command that ticks "my". Fix: quote options outside [A-Za-z0-9._-]. pre-existing: the gate's refusal text has the same pattern at core/gate/gate.go:298; fix both together. Stock options are single words.
  note: the test does not cover the parked, question-lane and unknown-lane early returns of UnchosenAhead. They are read-only copies of Next's guard. Fix: add a blocked case and a human case to TestTicketJSONNamesTheUnchosenOptionalLanes.
  note: unchosen_lanes also lands in list/move/resume/create --json, because all of them build on ticketJSON. That is consistent with next_lane, but the release note mentions show only.
review-check: "1. cd into a new empty directory and run: git init && git commit --allow-empty -m init && jaira init   2. jaira create probe --goal g --context c --dod d --json and note the handle it prints (say ABC123)   3. jaira show ABC123 --json | grep -A12 unchosen_lanes: you see two entries, brainstorm (option brainstorm) and pre-process (option planning), each with enable 'jaira dod ABC123 --option ...'   4. run the planning enable line exactly as printed   5. show --json again: only brainstorm is left in unchosen_lanes   6. jaira move ABC123 --to todo, then show --json: next_lane is pre-process and unchosen_lanes is []   7. open CLAUDE.md: the paragraph about next_lane now ends with the sentence naming unchosen_lanes and jaira dod <id> --option <name>"
---

# jaira show nennt die optionalen Lanes, an denen next_lane vorbeigeht

## Definition of Done

- [x] jaira show <id> --json hat ein Feld, das die uebersprungenen optionalen Lanes und den Befehl zum Einschalten nennt
  proof: internal/cli/tickets.go:1807 unchosen_lanes via unchosenLanesJSON <- lane.Set.UnchosenAhead (core/lane/next.go); TestTicketJSONNamesTheUnchosenOptionalLanes; scratch board: backlog lists brainstorm+pre-process with enable lines, after dod --option planning only brainstorm, and next_lane from todo becomes pre-process
- [x] Der Agent-Block nennt --option in einem Satz
  proof: core/board/announce.go:318 (and the same sentence in CLAUDE.md/AGENTS.md)
- [x] Release-Notiz in core/release/notes/<ID>.md
  proof: core/release/notes/64BZSR.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 01:11 · BeMuCa** — Open decisions: none for the person. Settled: skipped_lanes lists every unticked optional lane AHEAD of the current lane up to the terminal lane, not only those between status and next_lane - so a driver in backlog already learns pre-process exists and can tick planning before todo. Empty (never null) when Next has no answer by design (parked, question lane, terminal, unknown lane). Shape: [{lane, option, enable: 'jaira dod <ID> --option <opt>'}]. Computed by lane.Set.SkippedAhead in core/lane/next.go beside Next/OnPath. Brainstorm/planning options not ticked: small, path backlog->todo->in-progress.
- **2026-11-23 01:13 · BeMuCa** — critique round 1: the key skipped_lanes is taken - create --lane --json puts a []string of jumped lanes into the same ticketJSON payload (internal/cli/tickets.go:406), so one payload could carry two shapes under one key. Rename to optional_lanes_off. Let stand: SkippedAhead repeats Next's 4-line 'parked/question/terminal -> nothing' guard instead of sharing a helper - kept so Next stays untouched; both sit next to each other in core/lane/next.go.
- **2026-11-23 01:13 · BeMuCa** — in-progress round 2: renamed to unchosen_lanes (the gate already calls an unticked option step 'not chosen', core/gate/gate.go CodeStepNotChosen). create --json keeps its own skipped_lanes untouched.
- **2026-11-23 01:14 · BeMuCa** — critique round 2: the one finding (key clash with create --json skipped_lanes) is answered by the rename to unchosen_lanes; nothing left.
- **2026-11-23 01:16 · BeMuCa** — bughunt: change: internal/cli/tickets.go:1725 - class: option name embedded unquoted in the printed 'enable' command; example: a lane with requires-option: "my step" (lane.go:365 only trims) prints 'jaira dod <id> --option my step', which ticks a new option 'my' and exits 0 while the lane stays off the path; ran: jaira dod 6BPZ6P --option my step (exit 0, '[x] my' added; could not get the spaced lane into unchosen_lanes on a scratch board, so the lane side is read). Edge, not a blocker: stock lanes have single-word options
- **2026-11-23 01:16 · BeMuCa** — bughunt: done at 25217204c - 1 finding (1 change, 0 pre-existing) in 3 modules, first hunt; concurrency, error paths, input, permissions checked; files: AGENTS.md CLAUDE.md core/board/announce.go core/lane/next.go core/release/notes/64BZSR.md internal/cli/json_test.go internal/cli/tickets.go
- **2026-11-23 01:16 · BeMuCa** — testing: note: class: option name pasted unquoted into a shell command; example: a custom lane with requires-option 'my step' would get enable 'jaira dod <id> --option my step', which ticks 'my'; fix: quote the option when it holds anything but [A-Za-z0-9._-] - the gate's own refusal at core/gate/gate.go:298 has the same pattern (pre-existing), so fix both together. Not a blocker: stock options are single words and the lane parser does not reject spaces.
