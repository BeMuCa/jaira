---
id: 01M57ZXP95EN0S4B99GPNC323D
title: "Cockpit: Tickets aller Boards und der HITL-Lane in 'Waiting for you'"
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Das Cockpit-Band 'Waiting for you' zeigt die Tickets, die auf einen Menschen warten, auf jedem Board des Repos und auch in der Lane human (HITL), nicht nur in signoff des Default-Boards."
context: |-
  Was falsch ist: Das Cockpit (core/mod/builtin/jaira-cockpit) sieht nur das Default-Board.
  Ursache 1: register.tsx refreshBoard ruft 'jaira lanes --json' und 'jaira list --json' ohne --board. Tickets auf .jaira/boards/<name> (z. B. video-Board, siehe S0EG0M/8NHD8Q) fehlen im Band.
  Ursache 1b: 'jaira lanes' kennt kein --board (internal/cli/tickets.go newLanesCmd); das Cockpit kann die Lanes eines anderen Boards gar nicht lesen.
  Ursache 2: board.ts isPersons liest nur human_exit. Die Lane human hat requires-question, nicht requires-human-exit, also human_exit=false in 'jaira lanes --json'. Tickets in HITL (z. B. 6TCTFV) erscheinen nie im Band.
  Der Kommentar in tickets.go bei human_exit behauptet '(signoff, human)' – stimmt nicht.
  Ausloeser: Alex, 2026-10-18: 'cockpit sieht Tickets anderer Boards nicht; HITL-Tickets sollen auch sichtbar sein, nicht nur signoff'.
  Ausser Scope: Overview-Pane und Workers-Tab bleiben beim Default-Board.
definition-of-done: "'jaira lanes --board <name> --json' liest die Lanes des Boards <name>; unbekanntes Board ist exit 2"
tags:
  - cli
blocked-by: []
related: []
commits:
  - 21385a13f489c299d1b7027fd798a3ccf676f945
created-at: 2026-10-18T17:11:04Z
updated-at: 2026-10-18T17:32:36Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1943142
claimed-at: 2026-10-18T17:11:16Z
outcome-what: jaira lanes kennt --board und gibt requires_question aus; das Cockpit-Band liest alle Boards aus jaira boards und zaehlt Lanes mit requires_question (human) als wartend
outcome-why: "Das Band las nur das Default-Board und erkannte HITL nicht, weil human nur requires-question traegt, nicht human_exit"
outcome-resolves: "Tickets anderer Boards und HITL-Tickets fehlten in 'Waiting for you'"
review-summary: "jaira lanes nimmt jetzt --board <name> (über boardLanes, unbekanntes Board = exit 2 no_such_board) und gibt pro Lane requires_question aus (human=true, signoff=false). Im Cockpit zählt isPersons eine Lane als Personen-Lane, wenn human_exit ODER requires_question true ist (ohne beide Felder Fallback auf die Namen human/signoff). refreshBoard liest zusätzlich 'jaira boards --json' und für jedes Nicht-Default-Board nacheinander lanes/list mit --board, baut je Board einen eigenen Snapshot mit dessen Lanes und hängt waiting/theirs/counts an; ein Board, dessen Lesen scheitert, fällt weg. Accept erscheint nur, wenn next_lane terminal ist, also bei signoff, nicht bei HITL. Tests: TestLanesTakeABoard, zwei neue Cockpit-Tests; Release-Note und COMMANDS.md."
review-gaps: "Keine Lücke gegen die DoD. Hinweise: (1) isPersons wird auch von overview.ts (Overview-Pane, workersOf) benutzt – human zählt dort jetzt ebenfalls als Personen-Lane (nicht mehr in queue/working, eigene Spalte im Workers-Tab). Das stellt das Verhalten vor human_exit wieder her, ist aber weder getestet noch in der Release-Note erwähnt. (2) Give back schiebt fest nach in-progress; auf einem Board ohne diese Lane scheitert es (schon als Notiz offen). (3) Accept/Discuss für Tickets fremder Boards nur über den Test-Stub geprüft, nicht gegen ein echtes zweites Board (dieses Repo hat im Worktree keins)."
test-verdict: "pass: claude plugin test 160 pass, plugin validate passed, go test ./internal/cli ./core/mod ok; live: gebautes jaira lanes --json auf diesem Board: human requires_question=true, signoff human_exit=true; lanes --board nope exit 2. Nicht geprueft: echtes TUI mit zweitem Board"
review-verdict: "Diff erfüllt alle fünf DoD-Punkte; go test ./internal/cli ok, claude plugin test 160 pass. Keine Defekte gefunden; die Overview-Nebenwirkung von isPersons ist plausibel gewollt, sollte aber beim Signoff bewusst abgenickt werden."
review-round: 1
review-check: "1. cd /home/alex/projects/jaira-nc323d && go run ./cmd/jaira lanes --json | grep -B1 -A1 requires_question  → human hat requires_question true, signoff false.  2. go run ./cmd/jaira lanes --board nope --json; echo $?  → Fehler no_such_board, exit 2.  3. In einem Repo mit .jaira/boards/video (z. B. /home/alex/projects/jaira mit S0EG0M-Board): claude --plugin-dir /home/alex/projects/jaira-nc323d/core/mod/builtin/jaira-cockpit starten.  4. Band 'Waiting for you' aufklappen → Tickets aus signoff/approve des video-Boards und Tickets in human (z. B. 6TCTFV) stehen dort; human-Tickets ohne Accept-Button, signoff mit Accept.  5. Overview-Pane öffnen → human-Tickets stehen nicht mehr unter queue/working (Nebenwirkung, bewusst prüfen)."
---

# Cockpit: Tickets aller Boards und der HITL-Lane in 'Waiting for you'

## Definition of Done

- [x] 'jaira lanes --board <name> --json' liest die Lanes des Boards <name>; unbekanntes Board ist exit 2
  proof: internal/cli/tickets.go newLanesCmd --board; TestLanesTakeABoard
- [x] 'jaira lanes --json' gibt pro Lane requires_question aus
  proof: internal/cli/tickets.go requires_question; TestLanesTakeABoard
- [x] Das Cockpit-Band zeigt wartende Tickets aller Boards aus 'jaira boards --json'
  proof: jaira-cockpit/hooks/register.tsx refreshBoard otherBoardsOf; board.test.ts 'tickets waiting on another board are listed too'
- [x] Das Cockpit-Band zeigt Tickets in einer Lane mit requires_question (human) als wartend, ohne Accept
  proof: jaira-cockpit/hooks/board.ts isPersons requires_question; board.test.ts 'the band lists the tickets in a person's lanes' (human lane now requires_question only)
- [x] Release-Note in core/release/notes/<ID>.md; Go- und Cockpit-Tests gruen
  proof: core/release/notes/NC323D.md; claude plugin test 160 pass; go test ./internal/cli ok

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-18 17:17 · Alexander Sacharov** — Ausser Scope geblieben: Overview-Pane und Workers-Tab lesen weiter nur das Default-Board (pro Board andere Lanes, braeuchte eine Board-Auswahl). Give back schiebt fest nach in-progress – auf einem Board ohne in-progress (video) scheitert das mit der Fehlermeldung von jaira move.
- **2026-10-18 17:25 · Alexander Sacharov** — explain: Der Review fand keine Luecke gegen die Definition of Done; alle fuenf Punkte sind erfuellt, Tests sind gruen. Einzige Nebenwirkung: human zaehlt jetzt auch im Overview-Pane als Personen-Lane und steht nicht mehr unter queue/working.
Warum: Das ist plausibel gewollt, aber weder getestet noch in der Release-Note erwaehnt, und ein Mensch sollte es beim Signoff bewusst abnicken.
1) Zurueckgeben: Overview-Verhalten von human testen und in der Release-Note erwaehnen [give back]
2) Annehmen wie es ist; die Overview-Nebenwirkung ist bekannt und bewusst akzeptiert [accept]
Empfehlung: 2 - die Nebenwirkung stellt nur das Verhalten vor human_exit wieder her, und Overview/Workers-Tab sind bereits als ausser Scope notiert.
- **2026-10-18 17:31 · Alexander Sacharov** — decision: 2) Annehmen wie es ist; die Overview-Nebenwirkung ist bekannt und bewusst akzeptiert
