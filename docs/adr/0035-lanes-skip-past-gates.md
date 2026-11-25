# May a forward move skip installed lanes?
Status: Proposed
Date: 2026-10-17

## Context
The gate judges only the lane being left and the lane entered. in-progress -> review passes with critique and optimize never run;
a ticket can jump past a requires-human-exit lane from an earlier one, and --force takes it out of one.

## Decision
Open. Side A (code today): skipping is allowed; next_lane and prompts steer, the gate does not.
Side B (D28H7V): a skip is reported and the ticket records which lanes it never saw.
Side C (VW1WB1): a requires-human-exit lane cannot be jumped over, and --force does not leave it.

## Rejected
- None yet.

## Consequences
Only the verdict rail (verdict_sends_back) looks at lanes between source and target today.

## Sources
D28H7V, VW1WB1, core/gate/gate.go CheckAdvance
