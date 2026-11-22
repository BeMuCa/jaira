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
updated-at: 2026-11-22T18:43:10Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-4095349
claimed-at: 2026-11-22T18:31:38Z
outcome-what: "Tabs the cockpit opens (pr lane after Accept, give-back worker) get the cockpit session as --parent via parentArgs, only when its socket file exists; spawn.sh passes --parent to the pr lane; jaira-role-pr sends stops and report to its parent"
outcome-why: "Questions of cockpit-started tabs sat unseen in their own tabs; round 2: a gone socket made spawn.sh refuse and no tab started (bughunt)"
outcome-resolves: "T5QN9P: no question stays silent in a cockpit-started tab; without a live socket the tab asks itself as before"
review-summary: none
review-gaps: "none: one helper parentArgs replaces the inline socket lookup Discuss had (register.tsx openChat), so the socket is read in one place; no orphaned code; spawn.sh's own auto-parent (dispatch/worker only) left as is, it serves teamlead callers"
test-verdict: "pass: go build + go test -race -count=1 ./... RC=0; claude plugin test (cockpit) 343 pass; DoD 1-4 verified in the tree (proof tests run); spawn.sh pr by hand: live parent types '/jaira-role-pr X1 --by-person --parent uds:...' rc=0, --no-parent types it without, dead uds: exits 2 before a tab; the change writes no file or ticket itself (it starts tabs), so no two-writer case"
hunt-verdict: "fail: bughunt done at 7803b3e4f, 1 change: parentArgs passes a socket path without checking it exists; spawn.sh exits 2, Accept starts no pr tab and Give back has already closed the old worker tab (breaks DoD 3)"
round: 2
back-from: testing
back-by: lane
back-count: 1
---

# Tabs, die das Cockpit öffnet (pr-Lane nach Accept, Review, Give back), melden Fragen und Ergebnis an die Teamlead-Session

## Definition of Done

- [x] Ein vom Cockpit gestarteter Tab (pr-Lane, Review, Give-back-Worker) bekommt die Session, in der das Cockpit läuft, als parent - gleicher Mechanismus wie spawn.sh (Test)
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx parentArgs (pr lane in afterAccept, give-back worker); board.test.ts 'a branch ahead of the target without a pull request', 'otherwise closes the old worker's tab first'; role_test.go TestSpawnTypesTheRoleForItsLaneName (pr --parent)
- [x] Eine Frage und der Schlussbericht dieses Tabs kommen in der Teamlead-Session als Nachricht an; die Rolle (jaira-role-pr, Review) fragt die parent-Session statt nur in ihrem Tab (Test oder Rollen-Text + Test)
  proof: core/role/builtin/jaira-role-pr/SKILL.md '## Where your stops and your report go'; pr_round_test.go TestPRReportsToItsParent; review tab: its question goes via run-lane.sh exit 5 to its worker, whose parent is now the cockpit session (note decision 2)
- [x] Ohne erreichbare parent-Session fragt der Tab wie heute in sich selbst (Test)
  proof: register.tsx parentArgs checks $.fs.exists(socket); board.test.ts 'a socket whose file is gone is no parent' (pr and give-back), 'Herdr App (HC_TERMINAL_ID, no HERDR_ENV)', 'a session without a messaging socket starts the new worker with no parent'; jaira-role-pr SKILL.md unreachable-parent rule
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
- **2026-11-22 18:36 · Alexander Sacharov** — critique: none. Considered letting spawn.sh's own auto-parent (CLAUDE_CODE_MESSAGING_SOCKET in the env) cover give-back by just dropping --no-parent: rejected, it would not cover the pr lane and would rely on the mod's run() env; one explicit helper parentArgs (register.tsx) serves pr, give-back and discuss alike. Discuss without a socket now gets --no-parent instead of nothing: same prompt, spawn.sh only auto-parents dispatch/worker.
- **2026-11-22 18:38 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:876 - class: the cockpit hands spawn.sh its messaging socket as --parent without checking it exists, while spawn.sh:140 refuses a uds: parent that does not exist with exit 2 before any tab is made (its own default at spawn.sh:133 checks -e and falls back quietly); places: afterAccept register.tsx:770 (pr lane), giveBack register.tsx:1298 (worker), openChat register.tsx:897 (discuss, which already did this before); example: CLAUDE_CODE_MESSAGING_SOCKET set but the socket file gone or not visible to the shell spawn.sh runs in -> Accept leaves the ticket in done with no pull request started, and Give back has already closed the old worker's tab at register.tsx:1296 and then starts none (handedTo none); before this change both passed no parent / --no-parent and started regardless; fix direction: parentArgs falls back to --no-parent when the socket does not exist ($.fs.exists), as spawn.sh's own default does; ran: HERDR_ENV=1 bash spawn.sh --no-worktree --parent uds:/nonexistent/x.sock x T5QN9P worker <dir> (exit 2, 'no session listens there'); not a blocker by itself: needs a stale or invisible socket path
- **2026-11-22 18:38 · Alexander Sacharov** — bughunt: done at 7803b3e4f - 1 finding (1 change, 0 pre-existing) in 5 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/T5QN9P.md core/role/builtin/jaira-dispatcher/scripts/spawn.sh core/role/builtin/jaira-role-pr/SKILL.md core/role/pr_round_test.go core/role/role_test.go core/shipped/shipped.sha256
- **2026-11-22 18:38 · Alexander Sacharov** — testing: suite green (go -race RC=0, cockpit 343 pass), spawn.sh pr parent exercised by hand; waiting for the bug hunt
- **2026-11-22 18:41 · Alexander Sacharov** — testing: class: a parent address handed on without checking it is reachable; example: CLAUDE_CODE_MESSAGING_SOCKET set but its file gone -> spawn.sh exit 2 'no session listens there', give-back closed the old tab first (register.tsx giveBack), Accept leaves the pr unstarted; fix: parentArgs returns --no-parent unless $.fs.exists(socket), test with a socket env whose file is missing for pr and give-back
