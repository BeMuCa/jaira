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
updated-at: 2026-10-17T12:42:18Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-65566
claimed-at: 2026-10-17T12:36:37Z
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
