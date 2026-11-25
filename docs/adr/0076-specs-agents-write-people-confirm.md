# Requirements live in EARS specs and decisions in ADRs; agents write them, people confirm
Status: Accepted
Date: 2026-11-24

## Context
Sascha and Berk work in parallel and talk little; they built the same thing twice (X2M486 = AA2C3G) and decided against each other (Road vs Track swimlanes, ADR 0047).
What jaira does today was spread over about 250 tickets, notes and code; no single place said it.

## Decision
Sascha 2026-11-24: Spec-Driven Development with ADRs. docs/specs holds what jaira does today as EARS lines, at most 100 lines per file;
docs/adr holds one decision per file, at most 40 lines. People do not write specs; agents write them, people confirm or reject.
A contradiction between decisions becomes an ADR with Status Proposed for a person to settle.
A decision that changes a requirement goes back into the spec in the same pull request (MYMCDF decision 2026-11-24).

## Rejected
- People write specs by hand: they will not, and a spec nobody keeps is worse than none.
- Long prose design documents: a person reads at most about 100 lines per subject.

## Consequences
Each landed ticket that changes behaviour owes a spec line; until MYMCDF ships, nothing checks it and an agent updates the base by hand.
The gate MYMCDF plans (jaira validate checks that a decision naming a requirement ID changed it) is not built yet.
Migrating another project to this base is ticket 458TPX.

## Sources
7YQ1JC (context, decision 2026-11-24), MYMCDF (decision 2026-11-24 11:37), 458TPX, E3MBCB, docs/specs/README.md
