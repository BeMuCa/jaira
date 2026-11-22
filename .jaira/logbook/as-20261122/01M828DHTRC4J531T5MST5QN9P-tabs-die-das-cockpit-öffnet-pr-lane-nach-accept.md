---
id: 01M828DHTRC4J531T5MST5QN9P
title: "Tabs, die das Cockpit öffnet (pr-Lane nach Accept, Review, Give back), melden Fragen und Ergebnis an die Teamlead-Session"
status: done
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
commits:
  - 7803b3e4f972cade09aa8b0631be0b6592dbfc8b
  - f10d21b87ff0733597f43b1e7e4c2849da04d6c2
  - 972332856ac3c89d671863d5713392becfc15734
  - ab67c418e3aa3626789333f4d7e5a5a915bea8b5
  - d7a896e9bf40fb8c9626349a722346c845222728
created-at: 2026-11-22T18:31:11Z
updated-at: 2026-11-22T19:31:21Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2629852
claimed-at: 2026-11-22T19:14:24Z
outcome-what: "jaira-role-pr's description now names --parent and that a started-with-parent tab sends its questions and report to that session"
outcome-why: "give-back from signoff: the description, which is what a caller reads, did not mention --parent"
outcome-resolves: "DoD 5"
review-summary: "Das Cockpit gibt den Tabs, die es selbst startet, seine eigene Session als parent: der pr-Lane nach Accept (afterAccept) und der neue Worker bei Give back bekommen --parent uds:$CLAUDE_CODE_MESSAGING_SOCKET über den neuen Helper parentArgs (register.tsx); Discuss nutzt denselben Helper. Fehlt der Socket oder seine Datei, gibt parentArgs --no-parent zurück und der Tab startet wie früher und fragt in sich selbst - spawn.sh würde einen toten uds:-Parent sonst mit exit 2 ablehnen und gar keinen Tab starten. spawn.sh reicht --parent jetzt auch an den pr-Lane weiter (/jaira-role-pr <id> --by-person --parent <s>), und jaira-role-pr hat einen Abschnitt, der Stopps, Fragen und den Schlussbericht per note + SendMessage an diese Session schickt, mit Rückfall auf den eigenen Tab, wenn sie nicht erreichbar ist. Damit kehrt sich 7V1XPBs --no-parent beim Give back um. Release-Note und shipped.sha256 sind dabei."
review-gaps: |-
  note: known limit - alle Tests mocken fs.exists für den Socket-Pfad; ob $.fs.exists der Mod-Engine einen Pfad außerhalb des Projekts (/tmp/cc-socks-*, /run/user) sieht, ist ungeprüft - sieht es ihn nicht, ist die Funktion still aus (Tabs verhalten sich wie vor dem Ticket, nichts geht kaputt)
  note: known limit - ein alter installierter spawn.sh mit dem neuen Cockpit lässt --parent für den pr-Lane fallen (alter Prompt ohne parent); behoben durch jaira roles install, was die Release-Note sagt
test-verdict: "pass: go test -race ./core/role ./core/shipped RC=0; DoD 1-5 proofs present in the tree, DoD 5 verified: description at SKILL.md:3 names [--parent <session>] and where questions/report go, 784 chars (< 1024 skill limit), its sha256 is in shipped.sha256; change is role text only - no failure path or concurrent writer to exercise"
hunt-verdict: "pass: 0 findings at 972332856 (follow-up since f10d21b87)"
round: 3
back-from: signoff
back-by: person
back-count: 1
review-verdict: "accept: alle 5 DoD-Punkte durch den Diff gedeckt; parentArgs fällt bei fehlendem/totem Socket auf --no-parent zurück, also startet kein Pfad keinen Tab mehr; go test ./core/role ./core/shipped grün, claude plugin test 345/345 grün; kein Defekt bei Concurrency, Fehlerpfaden, Eingaben oder Rechten gefunden"
review-round: 3
review-verdict-gaps: |-
  note: known limit - alle Tests mocken fs.exists für den Socket-Pfad; ob $.fs.exists der Mod-Engine einen Pfad außerhalb des Projekts (/tmp/cc-socks-*, /run/user) sieht, ist ungeprüft - sieht es ihn nicht, ist die Funktion still aus (Tabs verhalten sich wie vor dem Ticket, nichts geht kaputt)
  note: known limit - ein alter installierter spawn.sh mit dem neuen Cockpit lässt --parent für den pr-Lane fallen (alter Prompt ohne parent); behoben durch jaira roles install, was die Release-Note sagt
review-check: "1. jaira roles install im Repo, dann claude im Herdr-Tab mit dem jaira-cockpit starten  2. ein Ticket in signoff im Cockpit mit Accept annehmen  3. im neu geöffneten pr-Tab muss die erste Zeile '/jaira-role-pr <ID> --by-person --parent uds:...' lauten  4. stoppt die Rolle (z.B. Konflikt) oder ist sie fertig, kommt in der Cockpit-Session eine Nachricht mit Frage bzw. Bericht an  5. bei einem anderen signoff-Ticket Give back drücken: der neue Worker-Tab startet mit '/jaira-worker <ID> --parent uds:...'  6. ohne Herdr-Tab ist es nicht prüfbar; Tests: cd core/mod/builtin/jaira-cockpit && claude plugin test . (345 pass), go test ./core/role -run 'TestPRReportsToItsParent|TestSpawn'"
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
- **2026-11-22 19:16 · Alexander Sacharov** — testing round 3: suite green, description checked for the 1024-char skill limit (784); waiting for bughunt
- **2026-11-22 19:16 · Alexander Sacharov** — bughunt: done at 972332856 - 0 findings in 2 modules, follow-up hunt since f10d21b87 (diff only: jaira-role-pr SKILL.md description names --parent, matching body lines 8/39/52; shipped.sha256 entry equals sha256 of the new SKILL.md; go test ./role ./shipped ok); concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/T5QN9P.md core/role/builtin/jaira-dispatcher/scripts/spawn.sh core/role/builtin/jaira-role-pr/SKILL.md core/role/pr_round_test.go core/role/role_test.go core/shipped/shipped.sha256
