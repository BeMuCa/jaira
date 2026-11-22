---
id: 01M828DHTRC4J531T5MST5QN9P
title: "Tabs, die das Cockpit öffnet (pr-Lane nach Accept, Review, Give back), melden Fragen und Ergebnis an die Teamlead-Session"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Jede Claude-Session, die das Cockpit in einem Herdr-Tab startet, kennt die Teamlead-Session als parent und schickt ihr Fragen und den Schlussbericht - wie Worker, die spawn.sh startet. Keine Frage bleibt still in einem Tab liegen."
context: |-
  Was falsch ist (Alex 22.11.2026): nach Accept öffnet das Cockpit für jedes Ticket einen Tab mit der pr-Lane. Diese Tabs laufen ohne --parent. Fragen landen nur in ihrem eigenen Tab; der Teamlead erfährt nichts.
  Heute gesehen: 1MYFKE- und 76ERHH-pr-Tab hingen mit 'Merge-Konflikt, was tun?', ein CRASZS-Review-Tab mit 'Gate bughunt verweigert signoff' - der Teamlead fand sie erst, als Alex fragte, durch Lesen der Tabs (hc read).
  spawn.sh gibt den Socket der Teamlead-Session als parent weiter (siehe jaira-teamlead SKILL); der Cockpit-Weg (core/mod/builtin/jaira-cockpit/hooks/register.tsx, stepAfterAccept/Give back) tut das nicht.
  Hinweis CRASZS-Worker: 'die pr-Lane laeuft ohne --parent'.
definition-of-done: "Ein vom Cockpit gestarteter Tab (pr-Lane, Review, Give-back-Worker) bekommt die Session, in der das Cockpit läuft, als parent - gleicher Mechanismus wie spawn.sh (Test)"
tags:
  - cockpit
blocked-by: []
related: []
commits: []
created-at: 2026-11-22T18:31:11Z
updated-at: 2026-11-22T18:36:25Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-4095349
claimed-at: 2026-11-22T18:31:38Z
outcome-what: "Tabs the cockpit opens (pr lane after Accept, give-back worker) get the cockpit session as --parent via one helper parentArgs; spawn.sh passes --parent to the pr lane; jaira-role-pr sends stops and report to its parent"
outcome-why: "Questions of cockpit-started tabs (1MYFKE/76ERHH merge conflicts, CRASZS review gate) sat unseen in their own tabs; the teamlead learned of them only by reading tabs"
outcome-resolves: "T5QN9P: no question stays silent in a cockpit-started tab; without a socket the tab asks itself as before"
---

# Tabs, die das Cockpit öffnet (pr-Lane nach Accept, Review, Give back), melden Fragen und Ergebnis an die Teamlead-Session

## Definition of Done

- [x] Ein vom Cockpit gestarteter Tab (pr-Lane, Review, Give-back-Worker) bekommt die Session, in der das Cockpit läuft, als parent - gleicher Mechanismus wie spawn.sh (Test)
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx parentArgs (pr lane in afterAccept, give-back worker); board.test.ts 'a branch ahead of the target without a pull request', 'otherwise closes the old worker's tab first'; role_test.go TestSpawnTypesTheRoleForItsLaneName (pr --parent)
- [x] Eine Frage und der Schlussbericht dieses Tabs kommen in der Teamlead-Session als Nachricht an; die Rolle (jaira-role-pr, Review) fragt die parent-Session statt nur in ihrem Tab (Test oder Rollen-Text + Test)
  proof: core/role/builtin/jaira-role-pr/SKILL.md '## Where your stops and your report go'; pr_round_test.go TestPRReportsToItsParent; review tab: its question goes via run-lane.sh exit 5 to its worker, whose parent is now the cockpit session (note decision 2)
- [x] Ohne erreichbare parent-Session fragt der Tab wie heute in sich selbst (Test)
  proof: board.test.ts 'Herdr App (HC_TERMINAL_ID, no HERDR_ENV)' (pr --no-parent) and 'a session without a messaging socket starts the new worker with no parent'; jaira-role-pr SKILL.md unreachable-parent rule
- [x] Release-Note core/release/notes/<ID>.md, shipped.sha256
  proof: core/release/notes/T5QN9P.md; core/shipped/shipped.sha256 +3 lines; go test ./core/role ./core/shipped ok

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 18:31 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 18:33 · Alexander Sacharov** — decision: open decisions counted, none for the person. (1) This reverses 7V1XPB's decision of 2026-10-23 ('a worker the cockpit starts stays in its own tab', --no-parent at register.tsx give-back): the ticket's own context from Alex 22.11. asks for the opposite, so the give-back worker now gets the cockpit session as parent. (2) The cockpit only starts three kinds of tab: pr lane (afterAccept), give-back worker, explain/discuss. A Review tab is never the cockpit's: the worker starts it via run-lane.sh, and its question already goes question-field -> run-lane exit 5 -> worker -> worker's parent. CRASZS's review question stuck because its worker was a give-back worker with --no-parent; fixing give-back closes that path. (3) Parent = uds:$CLAUDE_CODE_MESSAGING_SOCKET of the cockpit's session, the same address Discuss already passes; no socket in the env -> --no-parent, tab asks itself as today. (4) spawn.sh pr lane takes --parent explicitly only (no auto-parent: only the cockpit starts pr). (5) jaira-role-pr gets a 'Started with --parent' rule: stops/questions and final report go via SendMessage; parent unreachable -> note + ask/report in own tab.
