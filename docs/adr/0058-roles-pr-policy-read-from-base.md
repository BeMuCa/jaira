# Who opens a pull request is a board file read from the base branch
Status: Accepted
Date: 2026-11-21

## Context
Agents pushed and stopped; some boards want an agent to open the PR, a solo repo wants no PR at all.
A rule that lives only in a role's text, or on the branch's own disk, lets a branch open its own gate.

## Decision
.jaira/pr says person (default), agent or none; .jaira/pr-when says done (default) or signoff.
`jaira roles pr-gate <base>` reads both as <base> holds them and requires every ticket named on the branch to be accepted.
A person's own invocation (/jaira-role-pr, cockpit --by-person) needs no gate.

## Rejected
- Enforcing the rule in the role prompt alone (decision 2026-11-20 17:38: gate in the binary).
- Reading .jaira/pr off the checked-out branch (decision 2026-11-20 19:48: read from the base).

## Consequences
A change of policy takes effect only after it is merged into the base.
Under none every commit must name a ticket, because nobody reads the branch before it lands.

## Sources
A3R6YC, JG3REW (decisions 2026-11-20), XB97NJ (decisions 2026-11-20/21), CRASZS, 13VMA8
