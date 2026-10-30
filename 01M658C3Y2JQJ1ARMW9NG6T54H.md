---
id: 01M658C3Y2JQJ1ARMW9NG6T54H
title: Agents per ticket steht per Default auf worker statt lanes
status: critique
ready: true
creator: berkc
assignee: berkc
goal: "Wer nichts eingestellt hat, bekommt fuer ein ganzes Ticket jaira-worker (ein Agent durch alle Lanes, nur fresh-context-Lanes neu) statt jaira-dispatcher."
context: |-
  Wunsch von berkc am 30.10.2026: 'Agents per ticket ... das soll default auf worker sein'.
  Heute: core/settings/settings.go:278 DispatchMode gibt 'lanes' zurueck, ausser settings.json sagt ausdruecklich 'worker'. Hilfetext internal/cli/roles.go:40 nennt 'lanes' den Default; Settings-Zeile internal/tui/settings.go:56.
  berkc selbst hat schon dispatch: worker in ~/.jaira/settings.json - betrifft also neue Nutzer.
  Wer 'lanes' ausdruecklich gesetzt hat, behaelt es.
definition-of-done: "DispatchMode liefert worker, wenn dispatch leer oder unbekannt ist; 'lanes' nur ausdruecklich; Test"
tags:
  - cli
  - tui
blocked-by: []
related: []
commits:
  - b9872d8a67d695f61db5829a744aac2869f05fbe
  - b9872d8a
  - 297d5915
created-at: 2026-10-30T01:56:44Z
updated-at: 2026-10-30T21:53:47Z
updated-by: berkc
claimed-by: berk-ca-2129199
claimed-at: 2026-10-30T21:40:58Z
outcome-what: "README table rows say worker is the default and dispatcher is the lanes opt-in, fixing the one critique finding"
outcome-why: "the critique found README.md:822 still called worker the opt-in"
outcome-resolves: "G6T54H DoD 2: help and docs name worker as the default"
review-summary: "none"
round: 2
---

# Agents per ticket steht per Default auf worker statt lanes

## Definition of Done

- [x] DispatchMode liefert worker, wenn dispatch leer oder unbekannt ist; 'lanes' nur ausdruecklich; Test
  proof: core/settings/settings.go DispatchMode; core/settings/settings_test.go TestDispatchDefaultsToWorker
- [x] Hilfetexte (roles, Settings) nennen worker als Default
  proof: internal/cli/roles.go roles Long text; internal/tui/settings.go Agents per ticket row; TestSettingsEnterSwitchesDispatch, TestRolesListReportsTheDispatchSetting
- [x] Release-Note-Zeile
  proof: core/release/notes/G6T54H.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] settings.go: DispatchMode returns lanes only for an explicit lanes; comments
- [x] tests: settings_test, roles_test, tui settings_test follow the new default
- [x] help text: roles.go Long, tui settings row desc, teamlead/worker SKILL.md statements of the default
- [x] release note core/release/notes/G6T54H.md
- [x] go vet ./... and go test ./... green, then commit with ticket id

## Progress
- **2026-10-30 21:41 · berkc** — worker: open decisions counted, none the person must take. Settled: explicit 'lanes' keeps lanes; empty or unknown dispatch -> worker (DoD 1). Worktree is ../jaira-G6T54H on feat/G6T54H-worker-default, as the parent asked (the skill's default path is .worktrees/, not used). jaira update warning from 0.5.7 on this board was left alone.
- **2026-10-30 21:45 · berkc** — Plan: flipped DispatchMode to worker (only explicit lanes is lanes; empty/unknown -> worker). Left the TUI and help strings for later; tests flipped with it. Teamlead and worker SKILL.md and README default wording updated too, since they claimed the opposite.
- **2026-10-30 21:45 · berkc** — Correction to the previous note: the TUI row, roles help and their tests were changed in the same pass, not left for later.
- **2026-10-30 21:52 · berkc** — Commit b9872d8a carries the change. The shipped-version index needed the two new skill hashes appended (core/shipped/shipped.sha256 is append-only; TestEveryInstalledFileIsAShippedVersion failed until then). Full go test ./... exit 0 with 33 ok packages; go vet exit 0; gofmt clean.
- **2026-10-30 21:53 · berkc** — critique: README.md:822 table row says /jaira-worker is an opt-in alternative; fix that row only. Stood on: the note about the code diff (DispatchMode single source, no speculative code). The other README lines already fixed in b9872d8a.
