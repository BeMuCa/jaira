---
id: 01M5MNWMP0KPE76XV0DWBX1MNM
title: "jaira accept <id> --by <person>: eine Abnahme auf Wort der Person, die festhaelt, wer abgenommen hat"
status: in-progress
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
updated-at: 2026-10-26T19:20:35Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-268570
claimed-at: 2026-10-26T19:13:06Z
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
