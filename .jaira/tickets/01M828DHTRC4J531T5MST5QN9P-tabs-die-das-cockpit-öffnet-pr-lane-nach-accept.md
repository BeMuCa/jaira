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
updated-at: 2026-11-22T19:15:32Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2629852
claimed-at: 2026-11-22T19:14:24Z
outcome-what: "jaira-role-pr's description now names --parent and that a started-with-parent tab sends its questions and report to that session"
outcome-why: "give-back from signoff: the description, which is what a caller reads, did not mention --parent"
outcome-resolves: "DoD 5"
review-summary: "The cockpit now hands its own session to the tabs it starts as their parent. A new helper parentArgs (register.tsx) reads CLAUDE_CODE_MESSAGING_SOCKET and passes '--parent uds:<socket>' only when that socket file exists; otherwise '--no-parent'. It is used for the pr lane after Accept (before: no parent), the give-back worker (before: always --no-parent, which is why CRASZS's review question stuck) and Discuss (before: a socket was passed unchecked). spawn.sh forwards --parent into the pr prompt (/jaira-role-pr <id> --by-person --parent <s>), and jaira-role-pr gets a section: with --parent, every stop and the final report go onto the ticket as a note and to the parent via SendMessage; an unreachable parent falls back to asking in its own tab. A Review tab is not started by the cockpit (the worker starts it via run-lane.sh), so it reaches the cockpit through its worker, whose parent this now sets."
review-gaps: |-
  note: known limit: every test mocks $.fs.exists; whether the mod engine answers it for a path outside the project (/run/user/1000/cc-socks/*.sock or /tmp) is unverified - if it denies, the .catch turns that into --no-parent and the feature is silently off (tabs behave as before, no breakage); check live once, see review-check
  note: jaira-role-pr frontmatter description does not mention --parent; harmless, the Arguments line does
  note: an older installed spawn.sh with the new cockpit drops --parent for the pr lane (prompt without it) - harmless, the release note already says to rerun jaira roles install
test-verdict: "pass (round 2): cockpit claude plugin test 345 pass incl. 'a socket whose file is gone is no parent' for pr and give-back (neighbour of the example: no socket env at all, also --no-parent); go test ./core/role ./core/shipped ok; round 1 full go -race suite green"
hunt-verdict: "pass: bughunt done at f10d21b87, 0 change, 0 pre-existing"
round: 3
back-from: signoff
back-by: person
back-count: 1
review-verdict: "accept: all four DoD items met in the diff (parentArgs for pr/give-back/discuss, role text + TestPRReportsToItsParent, missing/gone socket -> --no-parent with tests, release note + shipped hashes); go test ./core/role ./core/shipped ok, claude plugin test 345 pass; only open point is the live check that $.fs.exists sees the socket path"
review-round: 2
review-verdict-gaps: |-
  note: known limit: every test mocks $.fs.exists; whether the mod engine answers it for a path outside the project (/run/user/1000/cc-socks/*.sock or /tmp) is unverified - if it denies, the .catch turns that into --no-parent and the feature is silently off (tabs behave as before, no breakage); check live once, see review-check
  note: jaira-role-pr frontmatter description does not mention --parent; harmless, the Arguments line does
  note: an older installed spawn.sh with the new cockpit drops --parent for the pr lane (prompt without it) - harmless, the release note already says to rerun jaira roles install
review-check: "1. jaira roles install in the repo, then restart the Claude session that runs the cockpit (in Herdr)  2. In the cockpit, Give back a ticket that sits in signoff  3. Look at the first line of the new worker tab: it must contain '--parent uds:' (the socket path). If it says --no-parent, the mod cannot see the socket file - the feature is off  4. Have that worker ask a question: it must arrive in the cockpit session as a message, not only in its tab  5. Optional: Accept a ticket with a branch ahead; the pr tab's prompt reads '/jaira-role-pr <ID> --by-person --parent uds:...'"
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
- [x] Die description von jaira-role-pr nennt --parent und wohin Fragen und Bericht dann gehen (Give-back 22.11.)
  proof: core/role/builtin/jaira-role-pr/SKILL.md:3 description; shipped.sha256 +1 line; go test ./core/role ./core/shipped ok

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
- **2026-11-22 18:46 · Alexander Sacharov** — bughunt: done at f10d21b87 - 0 findings in 3 modules, follow-up hunt since 7803b3e4f (diff only: parentArgs socket check, its two tests, shipped hash - hash matches sha256 of register.tsx); concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/T5QN9P.md core/role/builtin/jaira-dispatcher/scripts/spawn.sh core/role/builtin/jaira-role-pr/SKILL.md core/role/pr_round_test.go core/role/role_test.go core/shipped/shipped.sha256
- **2026-11-22 18:46 · Alexander Sacharov** — testing: note: known limit; unverifiable here whether the mod engine's $.fs.exists answers for a path under /run/user or /tmp (sandboxed?). If it cannot, parentArgs quietly gives --no-parent and tabs behave as before this ticket - check live: accept a ticket in the cockpit and look for '--parent uds:' on the pr tab's first line. Native Windows socket visibility also unverified.
- **2026-11-22 18:54 · Alexander Sacharov** — explain: Status: DoD 4/4 ok | tests ok | blockers ok | bughunt on HEAD ok
[known limit] Все тесты подменяют проверку файла сокета; вживую не проверено, видит ли движок модов путь вне проекта, и если нет, фича молча выключена (табы работают как раньше).
[small edge] В описании jaira-role-pr нет упоминания --parent, но в строке Arguments он есть, так что вреда нет.
[known limit] Старый установленный spawn.sh с новым cockpit теряет --parent для pr-лейна; лечится повторным jaira roles install.
1) Принять тикет как есть; живую проверку сделать после принятия: Accept в cockpit и взгляд на первую строку pr-таба [accept]
2) Вернуть и сначала проверить живьём, что $.fs.exists видит сокет, потом принимать [give back]
Recommendation: 1 - блокеров нет, тикет можно принять так; единственная открытая точка (живая проверка $.fs.exists) не ломает ничего, а только оставляет поведение как до тикета.
- **2026-11-22 19:13 · Alexander Sacharov** — decision: [small edge] В описании jaira-role-pr нет упоминания --parent, но в строке Arguments он есть, так что вреда нет.
- **2026-11-22 19:14 · Alexander Sacharov** — Given back from signoff: [small edge] В описании jaira-role-pr нет упоминания --parent, но в строке Arguments он есть, так что вреда нет.
- **2026-11-22 19:14 · Alexander Sacharov** — cockpit: given back to a new worker - nothing recorded about the old worker session (cache and context unknown)
