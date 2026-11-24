# The outcome belongs to the lane that produces it
Status: Accepted
Date: 2026-10-11

## Context
On 2026-10-11 review agents passing --what to jaira move wiped the implementer's account of the change.

## Decision
outcome-what, -why and -resolves are written by the lane whose output-produces names them (in-progress).
A move out of any other lane may fill an empty outcome field but is refused (exit 3, outcome_not_yours) when it would replace one.

## Rejected
- Letting every lane overwrite: the record of why the change was made is lost.
- Refusing any outcome flag from later lanes: an empty field may still be filled.

## Consequences
Replacing the outcome on purpose is `jaira set`, which says so. Later lanes write their findings into their own fields.
Shared review fields (review-summary from critique and review) are still overwritten per round (AHG254, open).

## Sources
QTHD8H (note 2026-10-11 20:37), AHG254
