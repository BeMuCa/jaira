# No signoff without a bug hunt that covers HEAD, on boards that run hunts
Status: Accepted
Date: 2026-11-21

## Context
A commit after the last hunt changed hunted files, and the ticket still landed in signoff with its hunt reading fail.
The testing lane starts the hunt beside its tests and owes hunt-verdict.

## Decision
An agent move of a ticket that changed code into a requires-human-exit lane needs a newest hunt state of ok or skipped covering HEAD.
Otherwise the gate refuses with needs_bughunt, or hunt_uncommitted for uncommitted hunted files.

## Rejected
- Guarding the question lane too: it takes questions, testing's max-rounds escape among them.
- Requiring hunts on every board: a board without a lane producing hunt-verdict would refuse every move.
- Running the hunt as its own lane: it moves no lane; it runs inside testing (Z66Z8E).

## Consequences
stop-after-fix (state skipped) is the ticket's own choice to run no more hunts and passes.

## Sources
G8HNHT (decision note 2026-11-21 21:00), Z66Z8E (2026-11-17), KCRS3W (2026-11-16), DF0YD6
