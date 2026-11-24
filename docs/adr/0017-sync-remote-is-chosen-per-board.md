# The ref remote is chosen per board by a fixed ladder
Status: Accepted
Date: 2026-09-14

## Context
The remote name was a machine-wide setting, but a machine holds boards whose repositories name their remotes differently (forks).
On a board without that remote every ref command failed, or a ticket could land in the wrong repository.

## Decision
RemoteFor walks: git config --local jaira.remote; settings.json "remote" if this repo has it; the repository's only remote;
else the configured name unchanged, so the failure is loud and names it. 'jaira whoami' says which step answered.

## Rejected
- One machine-wide remote name: wrong for forks and multi-remote checkouts.
- Guessing among several remotes: the fork case again.
- A global git config value: the same mistake as the machine-wide setting.

## Consequences
Worktrees of one clone share the --local value. A remote that exists but does not answer is still reported like a missing one (DGF5BW, open).

## Sources
9ET6NC, core/settings/settings.go RemoteSourceFor, DGF5BW (open)
