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
updated-at: 2026-11-23T01:12:50Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: worker-64bzsr
claimed-at: 2026-11-23T01:10:10Z
outcome-what: "jaira show --json gains skipped_lanes: every unticked optional lane ahead of the ticket with lane, option and the 'jaira dod <id> --option <opt>' that turns it on; the agent block says so in one sentence"
outcome-why: drivers skipped pre-process on tickets that needed a plan because next_lane walks past optional lanes silently (O11)
outcome-resolves: 64BZSR
---

# jaira show nennt die optionalen Lanes, an denen next_lane vorbeigeht

## Definition of Done

- [x] jaira show <id> --json hat ein Feld, das die uebersprungenen optionalen Lanes und den Befehl zum Einschalten nennt
  proof: internal/cli/tickets.go:406 1810 skippedLanesJSON from lane.Set.SkippedAhead (core/lane/next.go); TestTicketJSONNamesTheSkippedOptionalLanes; live: show 64BZSR --json in backlog listed brainstorm+pre-process with enable lines
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
