---
id: 01M5MNWB3682HRPZXX5Y1TFFG2
title: "spawn.sh startet Worker-Tabs im Berechtigungsmodus der Person und merkt, wenn der Prompt nicht ankommt"
status: critique
ready: true
creator: BeMuCa
goal: "Ein von spawn.sh gestarteter Worker- oder Review-Tab arbeitet sofort los: er haengt nicht an der ersten Bash-Rueckfrage, und spawn.sh meldet 'up' erst, wenn der Rollen-Prompt wirklich angekommen ist."
context: |-
  Beobachtet 22.-23.10.2026 im Teamlead-Betrieb (Feedback-Liste Punkte 1 und 2).
  spawn.sh startet plain 'claude' (core/role/builtin/jaira-dispatcher/scripts/spawn.sh:197) ohne --permission-mode. Jeder neue Tab blieb an seiner ersten Bash-Rueckfrage stehen - dreimal an einem Tag, einmal 49 Minuten, bis ein Mensch hinueberwechselte und erlaubte.
  Der Auto-Modus der Person steht in den User-Settings; ein defaultMode in Projekt- oder Worktree-Settings wurde laut Feedback nicht beachtet (nicht nachgeprueft).
  Einmal kam der Rollen-Prompt gar nicht an (1NATKP: 'never received a start prompt'), spawn.sh meldete trotzdem Erfolg. spawn.sh:236 erkennt einen Freigabe-Dialog nur waehrend des Starts.
  Verwandt: T22Z2C-Nachbar 'Erlaubnisliste fuer Worker' (eigenes Ticket), 3YRPXJ, 7MG5GB.
definition-of-done: "spawn.sh gibt claude den Berechtigungsmodus der startenden Person mit (--permission-mode, aus den User-Settings oder einem Board-Setting), und sagt im Hilfetext, welcher"
tags:
  - cli
blocked-by: []
related: []
commits:
  - 4b920015d063527e9d48d75cc242d8342bf17d2e
  - 1c4ef1ae3addca1b8809f93613402b0a2c14b73a
created-at: 2026-10-23T15:25:41Z
updated-at: 2026-10-26T19:41:16Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: driver-1TFFG2
claimed-at: 2026-10-26T19:41:13Z
outcome-what: "spawn.sh starts every tab's claude in the person's permission mode, resolved in Go by settings.PermissionModeFor (board .jaira/permission-mode without bypassPermissions, then permissions.defaultMode in the claude settings.json, then auto for an autoMode key) and shown by jaira roles list as permission_mode and permission_mode_source; after the enter it checks the role prompt arrived and exits 1 'prompt not received' when not, 4 with the pane printed when an approval dialog waits."
outcome-why: "A worker tab stopped at its first Bash question, once for 49 minutes, and in 1NATKP spawn.sh reported success though the prompt never arrived. Resolved in Go like chat-model and release so there is one reader and roles list says where the mode came from; exit 4 because run-lane.sh already reads 4 as the approval dialog and 3 as timed out."
outcome-resolves: "DoD 1: settings.PermissionModeFor, permission_mode in spawn.sh and usage(), TestPermissionModeFor, TestSpawnStartsClaudeInThePersonsPermissionMode. DoD 2: the arrival check after send-keys (exit 1 or 4), TestSpawnChecksThePromptArrived. DoD 3: Herdr-stub tests and core/release/notes/1TFFG2.md exist, go test ./core/role -race green; the full go test ./... -race is left for the run before the PR."
review-summary: |-
  core/role/builtin/jaira-dispatcher/scripts/spawn.sh:280-311 reads the board's .jaira/permission-mode with its own bash line parser, its own precedence and its own value check, beside the pattern this repo already has for a one-line board file: Go reads it (core/settings/settings.go:178 boardLine, :209 ChatModelFor, :262 ReleaseFor, each returning value and source), 'jaira roles list --json' carries it (internal/cli/roles.go:70-87), and spawn.sh takes it from there (spawn.sh:252, chat_model). Instead: add PermissionModeFor(root) (mode, source) to core/settings/settings.go (board via boardLine with bypassPermissions refused, then permissions.defaultMode in $CLAUDE_CONFIG_DIR or ~/.claude settings.json, then auto for an autoMode key, unknown values dropped like modelName does), carry it as permission_mode and permission_mode_source in roles list (JSON and the text line), and have spawn.sh read it with one 'jaira json get' from a single roles list call it shares with chat_model. The eight table cases of TestSpawnStartsClaudeInThePersonsPermissionMode then become a Go unit test, and one stub case checks that the flag reaches the command line.
  core/role/builtin/jaira-dispatcher/scripts/spawn.sh:427 exits 3 for 'approval dialog waits', but run-lane.sh:69 reads spawn.sh as pane="$(...)" under set -euo pipefail (run-lane.sh:7), so that 3 becomes run-lane.sh's own exit 3, which run-lane.sh:21 and jaira-dispatcher/SKILL.md:353-354 define as 'timed out' (checked with a stub: exit 3 inside x="$(...)" under set -e ends the caller with 3). Instead: use exit 4, which run-lane.sh:101 and SKILL.md:352 already use for an approval dialog. Change usage() and core/release/notes/1TFFG2.md to match, and add both new outcomes (1 = prompt not received, 4 = dialog) to the spawn.sh description in jaira-dispatcher/SKILL.md:291-295, which still ends with 'It prints the pane id'.
round: 2
---

# spawn.sh startet Worker-Tabs im Berechtigungsmodus der Person und merkt, wenn der Prompt nicht ankommt

## Definition of Done

- [x] spawn.sh gibt claude den Berechtigungsmodus der startenden Person mit (--permission-mode, aus den User-Settings oder einem Board-Setting), und sagt im Hilfetext, welcher
  proof: core/settings/settings.go PermissionModeFor; spawn.sh permission_mode from roles list + usage(); TestPermissionModeFor, TestSpawnStartsClaudeInThePersonsPermissionMode
- [x] spawn.sh prueft nach dem Start, dass der Rollen-Prompt angekommen ist (die erste Ausgabe der Rolle im Pane), und meldet sonst 'Prompt nicht angekommen' statt 'up'; ein wartender Freigabe-Dialog wird als solcher gemeldet
  proof: spawn.sh prompt-arrival check (exit 1/4); TestSpawnChecksThePromptArrived
- [ ] Ein Test mit einem Herdr-Stub deckt beide Faelle ab; go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-26 19:16 · BeMuCa** — driver: planning lane skipped as the ticket's Options say; the in-progress brief carried the design instead. Found out: Herdr 'pane get' has no field with assistant output (PaneInfo: agent, agent_status idle|working|blocked|done|unknown, agent_session, terminal_title, tokens) - so 'prompt arrived' is read as agent_status working after enter, blocked after enter = dialog waiting. claude --help lists --permission-mode choices acceptEdits, auto, bypassPermissions, manual, dontAsk, plan - no 'default'. Berk's settings.json has no permissions.defaultMode, only an autoMode object. Why the board file may not turn on bypassPermissions: .jaira/permission-mode is repo content every teammate pulls; one commit would otherwise run everybody's workers without any permission check. Existing spawn tests in core/role/role_test.go and reap_test.go stub pane get as idle forever and must answer working after send-keys.
- **2026-10-26 19:26 · BeMuCa** — claude accepts --permission-mode default and manual at parse time though --help omits default; invalid value fails: error: option '--permission-mode <mode>' argument 'bogus' is invalid. spawn.sh passes default through. Stub tests needed CLAUDE_CONFIG_DIR and a per-pane marker file (answer working after send-keys). DoD 3 left open: go test ./core/role -race green, full suite not run.
- **2026-10-26 19:32 · BeMuCa** — critique: two findings, both with a clear fix. (1) The permission-mode lookup in spawn.sh:280-311 is a second way of reading a one-line board file. .jaira/chat-model and .jaira/release are read in Go (core/settings/settings.go boardLine, ChatModelFor, ReleaseFor), shown with a source in 'jaira roles list --json', and spawn.sh reads chat_model from there. Do the same for permission-mode: add settings.PermissionModeFor(root) returning mode and source; roles list carries permission_mode and permission_mode_source; spawn.sh makes one 'jaira json get'. The case table becomes a Go unit test. (2) spawn.sh's new exit 3 (dialog waits) comes out of run-lane.sh:69 under set -e as run-lane's exit 3, and the dispatcher reads that as 'timed out' (SKILL.md:353). Use 4, which run-lane.sh already uses for an approval dialog, and name exit 1 and 4 where SKILL.md:291-295 describes spawn.sh. Deliberately left standing: (a) the floor of 5 polls after the enter. Total wait can then reach about 50 s instead of 40 s, but only when claude used the whole 40 s to come up. Zero polls left would report every slow start as 'not received'. The cockpit gives spawn.sh 120 s. (b) Only bypassPermissions is refused from the board; the other modes may come from there. The requirement says the board may override, and note 19:16 gives the reason for drawing the line at bypass. (c) The hard-coded list of modes: the driver checked that claude refuses an unknown value at parse time, so a typo in the board file would otherwise stop every spawn. (d) The stub line copied into six existing tests in role_test.go and reap_test.go: a surgical change, not a refactor.
- **2026-10-26 19:40 · BeMuCa** — round 2: no Go helper for the claude config dir existed (CLAUDE_CONFIG_DIR only appears in tests and spawn.sh), so PermissionModeFor reads it itself. roles list finds the board from its cwd via ticket.Discover, so the spawn stub test needs .jaira/tickets in the worktree or the board file is never read. Exit 4 for the dialog because run-lane.sh already maps 4 to approval dialog and 3 to timed out.
