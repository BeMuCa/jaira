# A worker closes its own tab; reap.sh cleans up after the merge
Status: Accepted
Date: 2026-11-20

## Context
After a working day finished worker, review and PR tabs and their worktrees stood around.
A give-back within the cache window is cheaper for the same worker than for a new one.

## Decision
close-tab.sh, called by the worker after its report, closes at once in the terminal lane and parks the tab in a person's lane.
reap.sh closes recorded tabs, runs .jaira/worktree-teardown and removes worktree and branch once the PR is merged.
Both serialise with flock on a file descriptor, which the kernel releases when the process dies.

## Rejected
- A fixed ~150k context limit as trigger (decision 2026-10-12: cold cache while waiting for a person).
- Reaping from spawn.sh (whoever starts nothing after a merge would keep the tabs).
- mkdir/hardlink locks with takeover by mtime (DBJHNE, N7JW2Q: replaced by flock).

## Consequences
Anything unrecorded, uncommitted or unanswerable is kept and named, never removed.

## Sources
YD4BE9 (decisions 2026-10-17), QRHB35, N7JW2Q (decisions 2026-11-18), P122GE, JDDCEZ, DBJHNE
