# A write on a file behind its ref: refuse and ask for pull, or catch up by itself
Status: Proposed
Date: 2026-10-12

## Context
Every ticket exists twice: the ref, and a file in each checkout. Another worktree or teammate moves the ref on, the file falls behind.
Today a stale file is refused on every write (exit 3, stale_file) until 'jaira pull'; readers already show the ref state (K6BGB0).

## Decision
Open. Side A (F3VNVC, K6BGB0, in code today): refuse every write on a stale file, also with --force, and name 'jaira pull'.
Side B (1K3RA8 in human, EZ3V6D in backlog): the write command fetches the ref and catches the file up itself, says so,
and refuses only when overwriting would lose something (an unsent write); the pull/behind markers disappear.

## Rejected
- Writing the stale file over the ref: lost a review, its proofs and nine notes in one 'move --force' (F3VNVC).

## Consequences
B removes a manual step agents hit daily (FYHJ39); A keeps every overwrite explicit and the person deciding which copy is right.

## Sources
F3VNVC (2026-10-12), K6BGB0 (signoff), EZ3V6D (2026-10-22, backlog), 1K3RA8 (human), FYHJ39
