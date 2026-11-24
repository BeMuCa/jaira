# Does an open plan step block the terminal lane?
Status: Proposed
Date: 2026-10-23

## Context
core/gate refuses a move into a terminal lane while any Plan item is unsettled (plan_incomplete): "the criteria cannot have been met while the work is under way".
CLAUDE.md says the plan is "a second and separate list" and the DoD is what the terminal gate reads; GH0GKD with DoD 9/9 was refused for plan step 30.

## Decision
Open. Side A (code today): an open plan item refuses the terminal lane.
Side B (NJBTVC, in signoff): only the definition of done decides; an open plan item becomes a hint on the successful move.

## Rejected
- None yet.

## Consequences
Whichever side wins, CLAUDE.md and the gate comment must say the same thing.

## Sources
NJBTVC, core/gate/gate.go (CodePlanIncomplete), CLAUDE.md "Working a ticket"
