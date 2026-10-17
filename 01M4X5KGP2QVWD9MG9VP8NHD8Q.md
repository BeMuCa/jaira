---
id: 01M4X5KGP2QVWD9MG9VP8NHD8Q
title: "Der Teamlead legt Tickets auf das Board, das zum Thema passt, und der Worker faehrt sie durch dessen Lanes"
status: testing
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
updated-at: 2026-10-17T12:46:19Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-65566
claimed-at: 2026-10-17T12:36:37Z
outcome-what: "New command 'jaira boards' lists every board with open tickets and its description, read from a plain-text 'description' file in the board's lanes dir (core/lane/boards.go Description/DescriptionPath). The teamlead skill reads it before every create and names board plus reason, or asks when none fits; dispatcher and worker skills read 'board' off the ticket and pass --board to next/list."
outcome-why: "With several boards per repo (S0EG0M) nothing told a teamlead which board a subject belongs to, and dispatcher's 'next --per-lane' saw only today's board, so a ticket on another board was invisible to it."
outcome-resolves: "DoD 1: jaira boards + TestBoardsListsEveryBoardWithItsDescription; DoD 2: teamlead SKILL.md section 'A new ticket goes on the board its subject belongs to'; DoD 3: dispatcher/worker skill text + TestWorkerPathRunsTheTicketsBoard, TestRolesWorkTheTicketsBoard; DoD 4: core/release/notes/8NHD8Q.md"
review-summary: "none"
round: 2
review-gaps: "none removed. Considered: the TUI's board cycle (internal/tui/model.go:2233) lists the same names but for switching, not counting - no shared helper worth it; boards.go loads today's lanes twice (loadEnv, then boardLanes for default) and loadEnv builds a link index it does not use - left, same as 'jaira tags' does and on a command run once before a create; reusing readIDList for the description is the existing trim-and-skip-blank reader, not a copy"
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
