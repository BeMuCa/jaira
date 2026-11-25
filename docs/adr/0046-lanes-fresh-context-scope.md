# Which lanes run in a fresh agent?
Status: Proposed
Date: 2026-10-31

## Context
MG677R decided one worker carries a ticket through every agentic lane and only `context: fresh` lanes get a new agent;
with no lane declaring context, only review is fresh. R2X4TJ lets one agent do critique and optimize in one pass for small changes.

## Decision
Open. Side A (MG677R, R2X4TJ, code today): only review is fresh; critique+optimize may share one pass.
Side B (XKS048, in signoff): critique and optimize each carry `context: fresh`, two more independent looks.

## Rejected
- None yet.

## Consequences
Side B costs two agent starts per round and conflicts with R2X4TJ's single small-change pass; side A keeps self-review in critique.

## Sources
MG677R (2026-10-11), R2X4TJ (2026-11-22), XKS048, 3CBCN6
