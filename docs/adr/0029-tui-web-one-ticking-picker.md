# People and milestones narrow the board through one ticking picker and one filter
Status: Accepted
Date: 2026-10-13

## Context
The user filter (u) let several people be ticked; the milestone picker (M) allowed only one.

## Decision
u and M run on one picker: space ticks, enter writes `user:`/`milestone:` into the board filter, x clears.
Those keys, like `tag:`, match exactly, so the board agrees with `jaira list`.

## Rejected
- A separate single-choice milestone mode.
- Substring matching for tags, people and milestones.

## Consequences
Everything that narrows the board is visible as text in the one filter.

## Sources
PNC854, N172JS
