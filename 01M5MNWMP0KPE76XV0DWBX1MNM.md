---
id: 01M5MNWMP0KPE76XV0DWBX1MNM
title: "jaira accept <id> --by <person>: eine Abnahme auf Wort der Person, die festhaelt, wer abgenommen hat"
status: critique
ready: true
creator: BeMuCa
goal: "Sagt eine Person im Chat 'accept', fuehrt der Agent die Abnahme ehrlich aus: ein Befehl nimmt das Ticket aus signoff, schreibt auf das Ticket, wer wann abgenommen hat, und die Ausgabe liest sich als Erfolg - nicht als uebergangene Verweigerung."
context: |-
  Beobachtet 22.-23.10.2026 (Feedback Punkte 7 und 17): der Teamlead kann ein Ticket auch auf ausdrueckliches Wort der Person nicht aus signoff bewegen. --force geht, druckt aber die Verweigerungen trotzdem ('Overrode 2 gate refusal(s): ...'), was schon am 31.08. einen Tag Umweg kostete (BDV0HM). Wer abgenommen hat, steht nirgends - nur dass ein Agent ueberstimmt hat.
  Die Abnahmeseite (/jaira-role-acceptance) wurde diesmal nicht benutzt; die Person nahm im Chat ab. Entweder wird die Rolle leichter, oder ein getipptes 'accept' reicht dem Teamlead.
  Gegenstueck existiert: A0C081 (jaira reject mit Grund, signoff). Verwandt: 64GNZT (Lane-Namen).
definition-of-done: "jaira accept <id> --by <person> bewegt ein Ticket aus signoff in seine naechste Lane ohne --force, schreibt accepted-by und accepted-at auf das Ticket und eine Notiz 'accepted by <person>'; ohne --by wird abgelehnt"
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T15:25:51Z
updated-at: 2026-10-26T19:30:01Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-268570
claimed-at: 2026-10-26T19:13:06Z
outcome-what: "Added `jaira accept <id> --by <person>`: it moves a ticket out of a human-exit lane (signoff) into that lane's next lane without --force and writes accepted-by, accepted-at and a note \"accepted by <person>\". The human-exit gate is satisfied by a per-request AcceptedBy. A --force move now prints one line \"overrode: <codes>\" before the move line. The teamlead role treats a typed \"accept <id>\" as the acceptance."
outcome-why: "The acceptance rides on the request, not on the ticket, so a ticket restored into signoff with an old accepted-by still waits for a person. It is written inside the move (move.Request.Record), so a refusal by any other gate writes nothing. accept stays on the normal move path, so every other gate still applies."
outcome-resolves: "A person's chat \"accept\" needed --force, which printed \"Overrode N gate refusal(s)\" and recorded no acceptor. accept now succeeds with a success line and records who accepted and when."
---

# jaira accept <id> --by <person>: eine Abnahme auf Wort der Person, die festhaelt, wer abgenommen hat

## Definition of Done

- [x] jaira accept <id> --by <person> bewegt ein Ticket aus signoff in seine naechste Lane ohne --force, schreibt accepted-by und accepted-at auf das Ticket und eine Notiz 'accepted by <person>'; ohne --by wird abgelehnt
  proof: internal/cli/accept.go; TestAcceptMovesOutOfSignoffAndRecordsWho, TestAcceptReachesATicketFiledOnEntry
- [x] Die Ausgabe eines --force-Moves nennt die uebergangenen Gates in einer Zeile als 'uebergangen', nicht im Wortlaut der Verweigerung
  proof: internal/cli/flow.go:250 (violationCodes); TestForcedMoveNamesOverriddenCodesOnOneLine
- [x] jaira-teamlead/SKILL.md sagt: ein getipptes 'accept <id>' der Person ist die Abnahme und wird mit jaira accept ausgefuehrt
  proof: core/role/builtin/jaira-teamlead/SKILL.md:211
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-26 19:15 · BeMuCa** — Driver: pre-process passed by (planning option unticked), so in-progress works from the driver's brief instead of a plan. DoD 4's race suite is run once before the PR, not in lanes; the lanes only write its release note.
- **2026-10-26 19:21 · BeMuCa** — Design: the human-exit gate is satisfied by gate.Request.AcceptedBy (set by move.Request, only by 'jaira accept'), never by reading accepted-by off the ticket, so a ticket restored into signoff with an old accepted-by still waits for a person. Added move.Request.Record, applied inside the write only when the gate passes: Stage persists before the gate and a refused move keeps it, but a refused accept must write nothing. Next lane = lanes.Next(t) (route order, 'done' on the shipped board), not the terminal lane. Not-in-a-human-lane is exit 3 (ExitValidation, 'a gate refused'), missing --by is exit 2 like reject's --why. Dead end: the brief said done files straight into the logbook on entry; the shipped correction done-logbook-on-entry removed that, so on the shipped board the accepted ticket stays in .jaira/tickets with status done. TestAcceptReachesATicketFiledOnEntry covers a board that still has logbook-on-entry. Also touched: README.md and docs/COMMANDS.md (TestEveryCommandIsDocumented demands a row). Seen, not fixed: jaira-role-acceptance SKILL.md still moves with --force; the CLAUDE.md/AGENTS.md agent block still says an agent cannot move work out of the human review lane; the TUI force notice still prints the refusal texts (internal/tui, only the CLI one-liner was asked). DoD 4 (go test ./... -race) left for the pre-PR run; release note is core/release/notes/BX1MNM.md.
