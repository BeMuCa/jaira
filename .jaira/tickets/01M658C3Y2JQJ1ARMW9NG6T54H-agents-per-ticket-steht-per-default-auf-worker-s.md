---
id: 01M658C3Y2JQJ1ARMW9NG6T54H
title: Agents per ticket steht per Default auf worker statt lanes
status: signoff
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
updated-at: 2026-10-30T22:03:14Z
updated-by: berkc
claimed-by: berk-ca-2129199
claimed-at: 2026-10-30T21:40:58Z
outcome-what: "README table rows say worker is the default and dispatcher is the lanes opt-in, fixing the one critique finding"
outcome-why: "the critique found README.md:822 still called worker the opt-in"
outcome-resolves: "G6T54H DoD 2: help and docs name worker as the default"
review-summary: "Flips the default of the 'dispatch' setting. core/settings/settings.go DispatchMode now returns 'lanes' only when settings.json says exactly 'lanes' (whitespace trimmed); empty, 'worker' or a typo return 'worker', so DispatchRole hands a teamlead /jaira-worker unless the person asked for the dispatcher. Nothing else decides the mode: 'jaira roles list' (text and --json dispatch/dispatch_role) and the TUI 'Agents per ticket' row/toggle read DispatchMode, so they follow without their own logic. Tests flipped with it (settings_test incl. ' lanes ', roles_test with an explicit lanes file, tui settings_test starting from worker). The roles help, the TUI row text, README (feature list, role table, worker paragraph), the teamlead and worker SKILL.md now call worker the default; the two changed SKILL.md hashes are appended to core/shipped/shipped.sha256; release note core/release/notes/G6T54H.md is one line telling users how to get lanes back."
round: 3
review-gaps: |-
  note: docs/COMMANDS.md:189 (the jaira roles list row) still says 'lanes (the default, jaira-dispatcher)' - untouched by the diff, so the outcome's 'help and docs name worker as the default' and the critique note 'grep finds no other claim' are not quite true; a one-line fix, and arguably inside DoD 2
  note: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:36 says the lane name worker is what a teamlead passes when settings.json says 'dispatch': 'worker' - now also when nothing is set
  note: README.md:820 and :849 still say the teamlead starts one dispatcher per ticket; the teamlead SKILL.md defines dispatcher as either role, the README does not
  note: README.md:77-80 says to run 'jaira roles install --force' after switching to lanes; the teamlead reads the setting at run time and the dispatcher is already installed, so that step is not needed for this direction
test-verdict: "pass: go test ./... -race -count=1 exit 0, 33 packages ok, go vet exit 0, gofmt clean; DoD 1-3 verified in the tree (settings.go DispatchMode, TestDispatchDefaultsToWorker, TestSettingsEnterSwitchesDispatch, TestRolesListReportsTheDispatchSetting, core/release/notes/G6T54H.md); behaviour exercised with built binary: no settings.json -> worker/jaira-worker, dispatch lanes -> lanes/jaira-dispatcher, typo lanse -> worker"
review-verdict: "The diff meets all three DoD items: DispatchMode defaults to worker with lanes only when explicit, tested; roles help and the TUI row name worker the default; the release note is there and follows the one-line format. No defects found - vet clean, core/settings, core/role, core/shipped, core/release, internal/tui and internal/cli roles tests pass, and the built binary reports worker with no settings file and lanes/jaira-dispatcher with one. Not a blocker but worth deciding at signoff: docs/COMMANDS.md:189 still documents lanes as the default."
review-round: 3
review-check: "1. cd /home/berkc/code/X_General/jaira-G6T54H && go build -o /tmp/j ./cmd/jaira  2. mkdir -p /tmp/jh && JAIRA_HOME=/tmp/jh /tmp/j roles list | grep dispatch  -> 'dispatch: worker - a ticket handed to an agent goes to /jaira-worker'  3. echo '{\"dispatch\":\"lanes\"}' > /tmp/jh/settings.json && JAIRA_HOME=/tmp/jh /tmp/j --json roles list  -> dispatch \"lanes\", dispatch_role \"jaira-dispatcher\"  4. rm /tmp/jh/settings.json, then JAIRA_HOME=/tmp/jh /tmp/j in a repo with a board, press S, go down to 'Agents per ticket'  -> the row text starts 'worker (the default)'; enter switches it to lanes and writes /tmp/jh/settings.json  5. /tmp/j roles list --help  -> says '\"worker\" (the default) starts jaira-worker'  6. grep -n default docs/COMMANDS.md | grep lanes  -> line 189 still says lanes is the default (the open note)"
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
- **2026-10-30 21:53 · berkc** — critique round 2: README table fixed in 297d5915; grep over the repo finds no other claim that lanes is the default outside historical logbook and NOTES files. Nothing left.
- **2026-10-30 21:54 · berkc** — optimize: nothing to remove. The diff touches only what the DoD names; no duplicate dispatch decision exists in the repo.
- **2026-10-30 22:03 · berkc** — review (fresh context): no blocker. Ran go vet and the tests of core/settings, core/role, core/shipped, core/release, internal/tui, internal/cli roles; built the binary and checked roles list with no settings (worker) and with dispatch lanes (jaira-dispatcher). The one finding worth a decision: docs/COMMANDS.md:189 still says lanes is the default - the critique round's grep missed it. Left as a note, not a blocker, because DoD 2 names the roles and Settings help texts and both are fixed.
