# Milestones stay local files, never in git, synced only through their refs
Status: Accepted
Date: 2026-11-21

## Context
A logbook commit (`git add -A .jaira`) brought 12 untracked milestone files to master (PR #167 removed them).
A stale milestone file in a worktree overwrote the newer ref and silently dropped tickets (KWQCYZ).

## Decision
Keep .jaira/milestones/<name>-<ID>.md as the working copy; init, share and update put /.jaira/milestones/ into .gitignore,
and update and validate warn about tracked milestone files.
The ref refs/jaira/milestones/<ID> is the sync path; before any milestone write the file is caught up to its ref (CatchUpMilestones).

## Rejected
- Milestones only on refs, no file (decision 2026-11-21 20:33/20:55, built and reverted): superseded at 22:05 the same night.
- Committing milestone files with the code: they travel back and forth between branches.

## Consequences
A hand edit that equals no earlier ref state is kept; one equal to an older ref state is overwritten as stale.
A hand edit must be sent with a milestone command before the next fetch replaces the file with its ref.
Milestones and epics now store differently (see ADR 0071); ES625B and AMJB8G would reopen this.

## Sources
1MYFKE (decisions 2026-11-21 20:33, 20:55, 22:05; note 2026-11-22 11:10), KWQCYZ
