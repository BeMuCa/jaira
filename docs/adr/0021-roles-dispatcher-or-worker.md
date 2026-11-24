# Which role carries a whole ticket, and which is the default
Status: Proposed
Date: 2026-10-11

## Context
Settings -> Agents per ticket chooses jaira-dispatcher (lanes) or jaira-worker (worker) for a whole ticket.
MG677R (decision 2026-10-11) made jaira-worker opt-in; core/settings DispatchMode returns lanes for anything but "worker".
G6T54H (signoff) makes worker the default; docs/COMMANDS.md still says lanes. NV33JM (backlog) removes the dispatcher, keeping its scripts.
XKS048 (signoff) adds critique and optimize as fresh lanes, moving the choice from the role to the lane.

## Decision
Open. Side A: keep both roles with lanes as default (today's code).
Side B: flip the default to worker (G6T54H). Side C: keep only jaira-worker with per-lane context (NV33JM, XKS048).

## Rejected
(none yet)

## Consequences
Two roles with near-identical prompts (dispatcher 584, worker 512 lines) drift; every rule change lands twice.
Until decided, the spec line TUI-73 describes lanes as the default.

## Sources
MG677R, G6T54H (decision note on COMMANDS.md:189), NV33JM, XKS048, core/settings/settings.go DispatchMode
