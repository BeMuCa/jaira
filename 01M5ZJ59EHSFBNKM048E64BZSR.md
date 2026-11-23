---
id: 01M5ZJ59EHSFBNKM048E64BZSR
title: "jaira show nennt die optionalen Lanes, an denen next_lane vorbeigeht"
status: critique
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
updated-at: 2026-11-23T01:13:58Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: worker-64bzsr
claimed-at: 2026-11-23T01:10:10Z
outcome-what: "jaira show --json gains skipped_lanes: every unticked optional lane ahead of the ticket with lane, option and the 'jaira dod <id> --option <opt>' that turns it on; the agent block says so in one sentence"
outcome-why: drivers skipped pre-process on tickets that needed a plan because next_lane walks past optional lanes silently (O11)
outcome-resolves: 64BZSR
review-summary: "internal/cli/tickets.go:406 create --json already writes skipped_lanes ([]string, lanes jumped over by create --lane) into the same ticketJSON payload; the new key collides with a different shape and meaning - rename the new field to optional_lanes_off (or similar) everywhere: tickets.go, json_test.go, announce.go, CLAUDE.md, AGENTS.md, release note"
round: 2
back-from: critique
back-by: lane
back-count: 1
---

# jaira show nennt die optionalen Lanes, an denen next_lane vorbeigeht

## Definition of Done

- [x] jaira show <id> --json hat ein Feld, das die uebersprungenen optionalen Lanes und den Befehl zum Einschalten nennt
  proof: internal/cli/tickets.go:1810 unchosen_lanes via unchosenLanesJSON <- lane.Set.SkippedAhead (core/lane/next.go); TestTicketJSONNamesTheUnchosenOptionalLanes; scratch board: backlog lists brainstorm+pre-process with enable lines, after dod --option planning only brainstorm, and next_lane from todo becomes pre-process
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
