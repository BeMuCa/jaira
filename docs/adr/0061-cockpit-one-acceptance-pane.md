# One Acceptance pane decides a waiting ticket: per-note ok / fix / ticket, then Execute
Status: Accepted
Date: 2026-11-22

## Context
Accept, Give back and Verdict showed the same review notes in three windows.
The person read notes in Give back (where each could be toggled) and then had to press another button.

## Decision
Each band row has one Acceptance button. The pane shows light, ticket, history and every note with ok / fix / ticket.
Execute derives the step: any fix or typed words -> give back; all ok -> accept into the next lane; ticket -> follow-up ticket.
The next lane is done from signoff, or one lane on from a person's lane mid-board (ADR 0077).

## Rejected
- Three buttons per row (Accept, Give back, Verdict) - same notes read three times (2V50Y5).
- Verdict text and all notes in the band itself - unreadable band; replaced by one line per ticket (PDZ444).

## Consequences
The band row carries only marks; all reading happens in the pane.
The pane stays open with its marks on any refusal and closes once the ticket moved.

## Sources
2V50Y5, PDZ444, MFYM1A, AJ41SC, MK8HD2, M5JSE9
