# The bug hunt runs as a background subagent inside the testing lane
Status: Accepted
Date: 2026-11-17

## Context
Tickets that passed every lane still had defects an external review found in the same code.
A separate hunt lane would cost a full extra round and a tab nobody waits on.

## Decision
The testing agent starts jaira-role-bughunt as a background subagent (ticket id only) and runs the tests itself meanwhile.
The hunt returns only defects the change introduced, as a class with an example, and later rounds read only the new diff.

## Rejected
- Starting the hunt through spawn.sh in a Herdr tab (it moves no lane, so run-lane.sh has nothing to wait for).
- A hunt after testing and before review as its own step (Z66Z8E folded it into testing).

## Consequences
A move into the person-released lane is refused when no hunt covers HEAD on a board whose lanes owe hunt-verdict (G8HNHT, ADR 0056).

## Sources
KCRS3W (decisions 2026-11-16), Z66Z8E (decision 2026-11-17), TW5FE2, DF0YD6, G8HNHT, KRBGV9
