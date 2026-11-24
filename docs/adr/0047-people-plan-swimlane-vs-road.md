# A roadmap swimlane is derived from after: lines, or is a named Road object
Status: Proposed
Date: 2026-11-03
## Context
Today a swimlane is the tree under a milestone with an empty after:; the first predecessor decides the lane; it has no name.
CLI help and release notes call it "swimlane"; KQQD1H (in signoff, branch feat/KQQD1H-roads, not merged) calls it "Road".
## Decision
Open. Side A (today's code): swimlanes stay derived from after:, unnamed.
Side B (KQQD1H, maintainer notes 2026-11-03): a Road is its own object on a ref with title, description and milestone list,
default Road-1..Road-n; milestones know nothing of their Road.
## Rejected
- (none decided)
## Consequences
Until decided, "swimlane" (CLI, docs) and "Road" (tickets 35Y6N5, VD4ZDW) name one thing.
## Sources
D0SAHM, KQQD1H, 35Y6N5, VD4ZDW
