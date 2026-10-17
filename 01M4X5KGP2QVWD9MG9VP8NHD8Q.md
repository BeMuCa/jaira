---
id: 01M4X5KGP2QVWD9MG9VP8NHD8Q
title: "Der Teamlead legt Tickets auf das Board, das zum Thema passt, und der Worker faehrt sie durch dessen Lanes"
status: in-progress
ready: true
creator: BeMuCa
assignee: Alexander Sacharov
goal: "Aus einer Bitte wie 'in der Datenbank fuer user eine Spalte hinzufuegen' legt der Teamlead ein Ticket an, waehlt das passende Board (hier backend), und der jaira-worker bekommt Ticket und Board und treibt es durch dessen Lanes"
context: |-
  Was heute fehlt: Der Teamlead (Skill jaira-teamlead) legt Tickets auf dem einen Board des Repos an. Ein Board pro Thema gibt es noch nicht (S0EG0M baut es).
  Berks Idee vom 14.10.: 'Wir koennten dann z.B. unserem Teamlead sagen: hey, ich will das und jenes machen, und er erstellt Tickets und matcht sie auch den Boards, die zu dem Thema passen.' Beispiel: 'in unserer Datenbank sollten wir fuer user noch eine Spalte hinzufuegen' -> backend-Board; der jaira-worker kriegt das Ticket und das Board, durch das er es treiben soll.
  Offen fuer brainstorm: woran der Teamlead das Board erkennt (Beschreibung pro Board in seiner Board-Datei? Tags?), und was passiert, wenn kein Board passt (nachfragen statt raten).
  Haengt an S0EG0M (mehrere Boards pro Repo).
definition-of-done: "Jedes Board traegt eine kurze Beschreibung seines Themas, die der Teamlead liest; Test"
tags: []
blocked-by: []
related:
  - 01M4X5JKRN4AV3KX56SES0EG0M
commits: []
created-at: 2026-10-14T12:18:43Z
updated-at: 2026-10-17T14:05:39Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-65566
claimed-at: 2026-10-17T12:36:37Z
outcome-what: "New command 'jaira boards' lists every board with open tickets and its description, read from a plain-text 'description' file in the board's lanes dir (core/lane/boards.go Description/DescriptionPath). The teamlead skill reads it before every create and names board plus reason, or asks when none fits; dispatcher and worker skills read 'board' off the ticket and pass --board to next/list."
outcome-why: "With several boards per repo (S0EG0M) nothing told a teamlead which board a subject belongs to, and dispatcher's 'next --per-lane' saw only today's board, so a ticket on another board was invisible to it."
outcome-resolves: "DoD 1: jaira boards + TestBoardsListsEveryBoardWithItsDescription; DoD 2: teamlead SKILL.md section 'A new ticket goes on the board its subject belongs to'; DoD 3: dispatcher/worker skill text + TestWorkerPathRunsTheTicketsBoard, TestRolesWorkTheTicketsBoard; DoD 4: core/release/notes/8NHD8Q.md"
review-summary: "New command 'jaira boards' (internal/cli/boards.go) lists today's board as 'default' plus every board under .jaira/boards/, each with its open tickets (not in its terminal lane) and a description read from a hand-written plain-text file 'description' in the board's lanes dir (core/lane/boards.go Description/DescriptionPath; wrapped lines joined into one); a board without one prints where to write it; --json gives name/open/description/description_file. The teamlead skill now runs 'jaira boards' before every create: one board - nothing to choose; one description fits - create with --board and name board plus the quoted description; none, two, or a missing description - ask the person. Teamlead also re-runs next/list with --board per extra board. Dispatcher and worker skills read 'board' from 'jaira show --json' and pass --board to 'next --per-lane'/'list'; worker no longer reads 'fresh' from 'jaira lanes' (today's board only). Docs, README and release note added; tests: TestBoardsListsEveryBoardWithItsDescription, TestWorkerPathRunsTheTicketsBoard, TestRolesWorkTheTicketsBoard."
round: 3
review-gaps: "Nothing blocks the DoD. Minor: (1) In any repo with a second board the default board has no description until someone writes .jaira/lanes/description, so by the skill's own rule the teamlead asks the person on every create until then - intended, but worth knowing for demo-video. (2) DoD 2 is proven only by string-presence tests on SKILL.md; real teamlead behaviour is untested by nature. (3) The generated CLAUDE.md agent block does not mention 'jaira boards' or --board, so an agent outside the roles never learns of other boards (already noted as left out). (4) A ticket whose board directory was removed is counted on no row of 'jaira boards'. (5) Testing found 'jaira validate' reports tickets on named boards as read-only - S0EG0M code, outside this ticket."
test-verdict: "pass: go vet + go test -race ./... green (RC=0); DoD 1-4 verified in the tree (TestBoardsListsEveryBoardWithItsDescription, TestWorkerPathRunsTheTicketsBoard, TestRolesWorkTheTicketsBoard, notes file); built binary on a scratch repo: 'jaira boards' lists default + video with descriptions, a missing one names .jaira/lanes/description, --for-lane on the video ticket gives the video lane's tier"
review-verdict: "Approve. The diff matches all four DoD items; go test ./internal/cli ./core/role ./core/lane passes; by hand in a scratch repo 'jaira boards' (text and --json) lists both boards with open counts and descriptions, and 'jaira set <id> board=video' moves a ticket across, falling back to the first lane when the target board lacks its lane. No defects found in the change itself; gaps above are follow-ups."
review-round: 2
review-check: "1. go build -o /tmp/jaira ./cmd/jaira  2. mkdir /tmp/b && cd /tmp/b && git init && git commit --allow-empty -m init && /tmp/jaira init  3. mkdir -p .jaira/boards/video/lanes && cp .jaira/lanes/*.md .jaira/boards/video/lanes/  4. echo 'Videos made from demos' > .jaira/boards/video/lanes/description  5. /tmp/jaira boards - you see 'default 0 open no description yet — write one into .jaira/lanes/description' and 'video 0 open Videos made from demos'  6. /tmp/jaira create x --board video --goal g --context c --dod d, then /tmp/jaira boards - video now says 1 open  7. /tmp/jaira roles install, then open the installed jaira-teamlead SKILL.md - section 'A new ticket goes on the board its subject belongs to' is there"
---

# Der Teamlead legt Tickets auf das Board, das zum Thema passt, und der Worker faehrt sie durch dessen Lanes

## Definition of Done

- [x] Jedes Board traegt eine kurze Beschreibung seines Themas, die der Teamlead liest; Test
  proof: core/lane/boards.go:141 Description; internal/cli/boards.go jaira boards; test TestBoardsListsEveryBoardWithItsDescription
- [x] Der Teamlead legt ein Ticket mit Board an und nennt seine Wahl mit Grund; passt kein Board, fragt er den Menschen statt zu raten
  proof: core/role/builtin/jaira-teamlead/SKILL.md:89 'A new ticket goes on the board its subject belongs to'; TestRolesWorkTheTicketsBoard
- [x] jaira-worker und Dispatcher lesen das Board vom Ticket und arbeiten dessen Lanes ab, nicht die des Standard-Boards; Test
  proof: core/role/builtin/jaira-dispatcher/SKILL.md:230, jaira-worker/SKILL.md:127; tests TestWorkerPathRunsTheTicketsBoard, TestRolesWorkTheTicketsBoard
- [x] Eine Zeile in core/release/notes/8NHD8Q.md
  proof: core/release/notes/8NHD8Q.md:1

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-17 12:36 · Alexander Sacharov** — Alexander Sacharov took this ticket over from BeMuCa
- **2026-10-17 12:36 · Alexander Sacharov** — decision 17.10. (Alexander, Teamlead-Session): Alexander uebernimmt das Ticket von Berk. Der Branch feat/8nhd8q zweigt von feat/s0eg0m ab, weil er --board und das Feld board: aus S0EG0M braucht; S0EG0M ist beim Abzweigen in review. Wird S0EG0M nach master gemergt, master hier einmergen. Ziel dahinter: in ~/projects/demo-video sagt eine Designerin ohne Programmierkenntnis dem Teamlead 'mach ein Video aus dieser Demo', und der legt das Ticket auf das Board video statt auf das SWE-Board.
- **2026-10-17 12:36 · Alexander Sacharov** — parent session: jaira-79
- **2026-10-17 12:38 · Alexander Sacharov** — found: show/--for-lane/move/resume already read the board off the ticket (S0EG0M). Blind to other boards: 'next --per-lane' and 'list' without --board (default board only), 'lanes' (no --board). No command lists the boards.
- **2026-10-17 12:38 · Alexander Sacharov** — decision (worker, no open decision left for a person): the description is a plain-text file 'description' in the board's lanes directory, beside 'order' - .jaira/lanes/description for today's board, .jaira/boards/<name>/lanes/description for a named one. Why: same idiom as order/corrections, hand-editable, and the lanes dir is what a board is. New command 'jaira boards' (like 'jaira tags') lists name, open tickets and description, --json too. Teamlead: one board -> no choice to name; several -> name the board and the reason in one line; none fits or descriptions missing -> ask the person as a choice. Dispatcher/worker: read 'board' from 'jaira show <id> --json' and pass --board to 'next --per-lane'.
- **2026-10-17 12:42 · Alexander Sacharov** — left out on purpose: the generated agent block in CLAUDE.md still describes only today's lanes and does not name 'jaira boards'; jaira-role-logbook-summary reads 'jaira list --json' (today's board only). Both are outside this DoD. Description lives in the lanes dir, so 'jaira lanes set --clear' (removes the default board dir) would take today's description with it.
- **2026-10-17 12:42 · Alexander Sacharov** — correction to the last note: 'jaira lanes set --clear' removes the default-board template file, not .jaira/lanes/, so it does not touch a description. Ignore that sentence.
- **2026-10-17 12:43 · Alexander Sacharov** — critique round 1: boards.go - replace map rows + type assertions with a row struct (pattern: tagCount in tags.go) and drop the filepath.Rel error fallback. Let stand: reusing readIDList for the description (it is exactly trim-and-drop-blank-lines), the double load of today's lanes (loadEnv + boardLanes), the skill texts.
- **2026-10-17 12:43 · Alexander Sacharov** — critique round 2: the round-1 finding is addressed in 227ade34 (boardRow struct, no Rel fallback); nothing else read, nothing left.
- **2026-10-17 12:46 · Alexander Sacharov** — testing: found outside this change (S0EG0M's code): 'jaira validate' checks every ticket against today's lanes, so a ticket on a named board is reported 'lane "cut" is not installed, so this ticket is read-only' - seen on a scratch board with a ticket in video/cut. Worker/dispatcher do not run validate, so this ticket's DoD is not touched; not fixed here.
- **2026-10-17 12:48 · Alexander Sacharov** — review (fresh context): approve. Tests pass, jaira boards and set board= checked by hand in a scratch repo. Follow-ups only, see review-gaps: default board needs .jaira/lanes/description or the teamlead asks on every create; CLAUDE.md agent block does not name jaira boards.
- **2026-10-17 12:51 · Alexander Sacharov** — done: the validate finding from testing is fixed on feat/s0eg0m (321ae0ad, c160454c: validate checks each ticket against its own board) and merged here in 0ab71820. Conflicts in core/lane/boards.go and internal/cli/boards_test.go were both appends at the file end; kept both. go vet + go test -race ./... green after the merge (RC=0).
- **2026-10-17 14:03 · Alexander Sacharov** — explain: Der Reviewer hat die Aenderung freigegeben (alle vier DoD-Punkte erfuellt, Tests gruen). Offen bleiben nur kleine Folgepunkte: Das Default-Board hat ohne Datei .jaira/lanes/description keine Beschreibung, und der CLAUDE.md-Agentenblock nennt 'jaira boards' nicht.
Warum: Ohne Beschreibung fragt der Teamlead bei jedem create nach dem Board, bis jemand die Datei schreibt. Ein Agent ausserhalb der Rollen erfaehrt nie von weiteren Boards.
1) Zurueckgeben: Beschreibung fuer das Default-Board und 'jaira boards'/--board im CLAUDE.md-Block ergaenzen [give back]
2) Annehmen wie es ist, die Folgepunkte separat behandeln (die Beschreibung bei Bedarf fuer demo-video schreiben) [accept]
Empfehlung: 2 - Nichts blockiert die DoD, die Luecken sind bewusst ausgelassen und gehoeren nicht in dieses Ticket.
- **2026-10-17 14:05 · Alexander Sacharov** — decision: 1) Zurueckgeben: Beschreibung fuer das Default-Board und 'jaira boards'/--board im CLAUDE.md-Block ergaenzen
- **2026-10-17 14:05 · Alexander Sacharov** — Given back from signoff: 1) Zurueckgeben: Beschreibung fuer das Default-Board und 'jaira boards'/--board im CLAUDE.md-Block ergaenzen
