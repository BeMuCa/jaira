# Whether an agent may ever merge a pull request is undecided
Status: Proposed
Date: 2026-11-22

## Context
Side A: 50RA6N added .jaira/land = agent; `pr-gate --merge` and teamlead/scripts/land.sh then merge via gh.
Side B: teamlead SKILL.md ("You never open, merge or approve a pull request"), jaira-role-pr ("never merge"),
`jaira roles list` help ("Merging and approving stay a person's either way") and the CLAUDE.md block say no agent merges.

## Decision
Open. Either the land permission is the documented exception in every text, or land.sh hands back the merge line always.

## Rejected
(none yet)

## Consequences
Until settled, an agent reading the role text and one running land.sh follow two rules for one case.

## Sources
50RA6N (2026-11-22, "Alex chose .jaira/land = agent"), 13VMA8, A3R6YC, JG3REW
