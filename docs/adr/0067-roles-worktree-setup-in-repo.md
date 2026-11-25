# spawn.sh makes the worktree; the repository sets it up
Status: Accepted
Date: 2026-11-22

## Context
spawn.sh copied .env and appended port lines named after one project, values that belong to no other repo.

## Decision
spawn.sh makes the worktree, opens the tab and starts the agent, nothing more.
An executable .jaira/worktree-setup gets <worktree> <slug> <offset> <root>; .jaira/worktree-teardown gets the same at reap.
spawn.sh keeps only what one repo's script cannot promise: a port offset 1-40 unique across standing worktrees.

## Rejected
- Project-specific .env/port writing inside spawn.sh (kept as worktree-setup.example).
- Deriving teardown arguments again at reap time (recorded as jaira-setup instead; a wrong slug tears down another stack).

## Consequences
A repository without a setup gets a bare worktree and no offset; a broken setup stops the spawn before a tab opens.

## Sources
167HY8, 7MG5GB, JDDCEZ (decision 2026-11-21), DBJHNE (decisions 2026-11-21)
