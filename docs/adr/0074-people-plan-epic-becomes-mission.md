# An epic becomes a Mission with its own lanes and human spec confirmation
Status: Proposed
Date: 2026-11-24
## Context
Spec-driven development without people writing specs: people only confirm them, through lanes.
Today an epic runs through no lane and has no definition of done (F914QH).
## Decision
Open (E3MBCB): rename Epic to Mission; lanes idea -> spec (agent writes the spec delta and ADRs, a person confirms) ->
tickets (each carries `spec: <ID>`) -> knowledge (validate specs and ADRs are current) -> done.
Side A: epic stays a description plus list (F914QH). Side B: Mission with lanes (E3MBCB).
## Rejected
- (none decided)
## Consequences
Would reverse "an epic runs through no lane"; epic commands, refs and web views would need a new name or an alias.
## Sources
E3MBCB (2026-11-24), F914QH
