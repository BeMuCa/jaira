# A blocker starts as "fix"; accepting it anyway needs a typed reason
Status: Accepted
Date: 2026-11-21

## Context
Accepting from signoff filed follow-up tickets and accepted blockers without anyone saying why.
A blocker is meant to be fixed; leaving one standing is a decision of its own.

## Decision
Preselection: blocker -> fix, every other class -> ok; a follow-up ticket only for a row the person turns to "ticket".
Accepting with a blocker not marked fix requires a typed reason, written as a `decision:` note.

## Rejected
- Every note preselected as follow-up ticket - a plain Accept would file tickets on the shared board.
- Accept with a blocker silently - nobody could later tell why it stayed.

## Consequences
On surfaces without a text field (mobile) a blocker cannot be accepted from the pane.

## Sources
MFYM1A (decision notes 2026-11-21 19:39 and 19:40), 2V50Y5
