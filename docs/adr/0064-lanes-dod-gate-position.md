# Where does an unticked definition of done stop a ticket?
Status: Proposed
Date: 2026-11-22

## Context
The DoD gate sits on the terminal lane only; tickets reach review and signoff with DoD 0/n and look ready to a person.
77B4GA added a stderr warning on leaving in-progress with nothing ticked, and kept the gate "as strict as it is".

## Decision
Open. Side A (77B4GA, code today): warn when leaving the implementing lane; refuse only at the terminal lane.
Side B (48QE0E): refuse already on the move into a lane only a person can leave.

## Rejected
- None yet.

## Consequences
Side B moves the cost of an unticked DoD to before the person looks; it also needs a rule for superseded items and the plan.

## Sources
77B4GA (decision note 2026-11-22 18:50), 48QE0E, NYW4M7
