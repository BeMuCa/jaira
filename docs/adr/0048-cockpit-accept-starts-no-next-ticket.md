# Accept opens the pull request but never starts the next ticket
Status: Accepted
Date: 2026-11-09

## Context
GPFA7B made Accept start the follow-up step itself: open a PR, or else start the next actionable ticket in a Herdr tab.
On 2026-11-09 Accept started an unrelated ticket and the person did not know what ran or why.

## Decision
After done, Accept starts only the pr lane (or the land lane where `.jaira/pr` says none) when the branch needs it.
Which ticket runs next is the teamlead's call; Accept submits one line to the session so the teamlead goes on.

## Rejected
- Accept starts the next ticket from `jaira next` (GPFA7B, 2026-10-12) - hidden work nobody asked for.
- Two buttons "Accept + PR" / "Accept + next" - GPFA7B question 1 chose one rule in code.

## Consequences
A board without a teamlead session does not move on by itself after Accept.

## Sources
X6JHQX, GPFA7B, 409XCV, XB97NJ
