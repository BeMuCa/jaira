---
id: 01M67RWX6BPVT0QKN5ZAX1M1AF
title: "Neue Tickets haben die Option planning standardmaessig angehakt, damit jedes Ticket durch den Planer geht"
status: in-progress
ready: true
creator: berkc
assignee: berkc
goal: "Ein neu angelegtes Ticket hat unter Options 'planning' angehakt, solange die Person keine eigene Default-Board-Datei hat, die etwas anderes sagt; die pre-process-Lane (Opus-Planer) laeuft dadurch fuer jedes Ticket."
context: |-
  Was falsch ist: Neue Tickets haben alle Options unangehakt (core/lane/defaultboard.go:185 ResolveOptions tickt nur, was die Default-Board-Datei nennt; ohne Datei nichts). Die pre-process-Lane (requires-option: planning) wird deshalb fast immer uebersprungen.
  Ausloeser: Experiment 3CBCN6 auf NDE933 (31.10.2026). Nur-Haiku vergass, Hilfezeile und Kommentar mitzuziehen (Urteil 8/7/7/6); Opus-Plan + Haiku-Bau lag vorne (9/9/8/9). Der Maintainer: 'should be on per default der planning tag. damit jeder durch den opus planer geht.'
  Bekannt: Optionen kommen aus requires-option der installierten Lanes (core/lane/lane.go:233 Options). Die Default-Board-Datei pro Person bleibt der Weg, es abzuschalten.
  Haengt zusammen mit CXH98W: dass pre-process wirklich mit Opus laeuft, erzwingt erst CXH98W (Modell pro Lane fuer Worker).
definition-of-done: "jaira create ohne Default-Board-Datei legt ein Ticket an, dessen Options 'planning' angehakt hat; Test dafuer"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-31T01:24:00Z
updated-at: 2026-10-31T02:35:35Z
updated-by: berkc
claimed-by: berk-ca-3129208
claimed-at: 2026-10-31T01:27:35Z
---

# Neue Tickets haben die Option planning standardmaessig angehakt, damit jedes Ticket durch den Planer geht

- [x] core/lane/defaultboard.go: add exported 'var DefaultOptions = []string{"planning"}' with a comment (the Options a new ticket starts with when no default board names options:) and an unexported field 'optionsDefaulted bool' on DefaultBoard; extend the DefaultBoard type comment: an absent options: key means DefaultOptions, options: [] means none
- [x] LoadDefaultBoard: right after 'b := &DefaultBoard{Path: path}' set b.Options = append([]string(nil), DefaultOptions...) and b.optionsDefaulted = true (covers absent file AND unparseable file). Replace the options read with: if d.Has("options") { if opts, err := d.List("options"); err == nil { b.Options, b.optionsDefaulted = opts, false } }. Keep lanes read unchanged. Do NOT touch ResolveOptions or NewBody logic
- [x] Validate: wrap the unknown-option loop in 'if !board.optionsDefaulted' so a board without the pre-process lane gets no spurious 'option planning' warning when the user has no options: line
- [ ] core/ticket/body.go NewBody comment lines 45-50: rewrite 'Unticked by default...' to say callers pass the ticks; without a default board naming options: planning starts ticked (lane.DefaultOptions) so every ticket goes through the pre-process planner; a per-user default board can untick it. Do not change NewBody code or the TestNewBodyMatchesRegressionBaseline test (it passes explicit options, stays green)
- [ ] internal/cli/tickets.go defaultBoardTemplate: change comment under options: to say absent means [planning], [] means none pre-ticked. internal/cli/lanes.go 'lanes default' Long help: add one sentence that with no options: line a new ticket starts with planning ticked
- [ ] tests core/lane/defaultboard_test.go: rename/rewrite TestLoadDefaultBoardAbsentReturnsZeroBoard -> absent file gives Lanes empty, Options == [planning], no warnings; add TestLoadDefaultBoardEmptyOptionsListTicksNothing (file 'options: []' -> Options empty); add TestLoadDefaultBoardNoOptionsKeyDefaultsToPlanning (file with only lanes:); add TestValidateSilentOnDefaultedOptionWithoutPlanningLane (set without pre-process, absent file -> no warning)
- [ ] tests internal/cli/init_test.go: add TestCreateWithoutDefaultBoardTicksPlanning (initTestHome, JAIRA_DEFAULT_BOARD=nonexistent, init+create, body contains '- [x] planning' and '- [ ] brainstorm') and TestCreateWithDefaultBoardOptionsWithoutPlanningLeavesItUnticked (file 'options: [brainstorm]' -> body has '- [ ] planning'; extend existing TestCreateAfterInitTicksDefaultBoardOptions with that assertion instead of a new test is fine)
- [ ] core/release/notes/X1M1AF.md: one line, e.g. '- ! Expect new tickets to start with planning ticked, so they go through pre-process; to keep it off write options: [] (or your own list) into your default board: jaira lanes default --options brainstorm'
- [ ] run go test ./core/lane ./core/ticket ./internal/tui -count=1 and go test ./internal/cli -count=1 -skip 'TestUnreachableRemoteIsOfflineNotARace|TestAnOfflineWriteArrivesWhenTheNetworkIsBack|TestWebPushReportsAQueuedDrop'; any failure from a ticket now carrying [x] planning (e.g. createlane_test comment 'ticks no option', lanes default output 'none pre-ticked'): fix the test's expectation or set JAIRA_DEFAULT_BOARD with options: [] in that test, note which in jaira note
- [ ] commit feat(X1M1AF): ... with ticket file, tick DoD 1-3 with test names as proof

## Definition of Done

- [ ] jaira create ohne Default-Board-Datei legt ein Ticket an, dessen Options 'planning' angehakt hat; Test dafuer
- [ ] Eine Default-Board-Datei, die options ohne planning nennt, laesst planning unangehakt; Test dafuer
- [ ] Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-31 02:33 · berkc** — Planner decision: default lives in LoadDefaultBoard (absent options: key -> [planning]), not in ResolveOptions. Why: the TUI default-board screen (internal/tui/defaultboard.go:62) and 'jaira lanes default --lanes X' (internal/cli/lanes.go) both seed from board.Options and SaveDefaultBoard always writes options:; a default only in ResolveOptions would show planning unticked there and a lanes-only save would write options: [] and silently turn planning off. Distinguish absent key (d.Has) from 'options: []' (= none, full control). Validate must skip the option check when defaulted, else boards without pre-process warn 'option planning not required'. NewBody untouched (pure; regression baseline test passes explicit options). Tests changing: TestLoadDefaultBoardAbsentReturnsZeroBoard (Options now [planning]). Unaffected: body_test, TestResolveOptionsTicksNamedOptions, TUI TogglesAndSaves (opts[0]=brainstorm). Docs: docs/COMMANDS.md does not mention options; only the template in tickets.go:1552 and 'lanes default' help do. Unparseable file is treated as absent -> planning ticked too.
