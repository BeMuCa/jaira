---
id: 01M5MNWB3682HRPZXX5Y1TFFG2
title: "spawn.sh startet Worker-Tabs im Berechtigungsmodus der Person und merkt, wenn der Prompt nicht ankommt"
status: in-progress
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
commits: []
created-at: 2026-10-23T15:25:41Z
updated-at: 2026-10-26T19:25:59Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-253833
claimed-at: 2026-10-26T19:10:59Z
---

# spawn.sh startet Worker-Tabs im Berechtigungsmodus der Person und merkt, wenn der Prompt nicht ankommt

## Definition of Done

- [x] spawn.sh gibt claude den Berechtigungsmodus der startenden Person mit (--permission-mode, aus den User-Settings oder einem Board-Setting), und sagt im Hilfetext, welcher
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh permission_mode block + usage(); TestSpawnStartsClaudeInThePersonsPermissionMode
- [x] spawn.sh prueft nach dem Start, dass der Rollen-Prompt angekommen ist (die erste Ausgabe der Rolle im Pane), und meldet sonst 'Prompt nicht angekommen' statt 'up'; ein wartender Freigabe-Dialog wird als solcher gemeldet
  proof: spawn.sh prompt-arrival check (exit 1/3); TestSpawnChecksThePromptArrived
- [ ] Ein Test mit einem Herdr-Stub deckt beide Faelle ab; go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-26 19:16 · BeMuCa** — driver: planning lane skipped as the ticket's Options say; the in-progress brief carried the design instead. Found out: Herdr 'pane get' has no field with assistant output (PaneInfo: agent, agent_status idle|working|blocked|done|unknown, agent_session, terminal_title, tokens) - so 'prompt arrived' is read as agent_status working after enter, blocked after enter = dialog waiting. claude --help lists --permission-mode choices acceptEdits, auto, bypassPermissions, manual, dontAsk, plan - no 'default'. Berk's settings.json has no permissions.defaultMode, only an autoMode object. Why the board file may not turn on bypassPermissions: .jaira/permission-mode is repo content every teammate pulls; one commit would otherwise run everybody's workers without any permission check. Existing spawn tests in core/role/role_test.go and reap_test.go stub pane get as idle forever and must answer working after send-keys.
