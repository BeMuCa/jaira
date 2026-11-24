# Planning objects are stored two ways: milestones as file plus ref, epics on the ref only
Status: Proposed
Date: 2026-11-23
## Context
Milestones keep a gitignored working file mirrored to refs/jaira/milestones/<ID>; epics have no file, refs/jaira/epics/<name> only.
Roads (KQQD1H) and Missions (E3MBCB) are about to add more planning objects; ES625B asks for all shared state on refs.
## Decision
Open. Side A (1MYFKE, 22:05): milestones stay files, never in git, caught up from the ref before every write.
Side B (1MYFKE 20:33/20:55, F914QH, ES625B): every planning object lives on its ref only, like epics.
## Rejected
- (none decided)
## Consequences
Two write paths and two staleness rules for one kind of thing; a new object has to pick one.
The wider question of moving all shared state onto refs is ADR 0039.
## Sources
1MYFKE (2026-11-21 20:33, 20:55, 22:05), F914QH (2026-10-14 08:49), ES625B, KQQD1H, E3MBCB
