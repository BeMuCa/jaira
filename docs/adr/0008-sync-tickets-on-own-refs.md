# Each ticket travels on its own git ref; the file rides only with code
Status: Accepted
Date: 2026-09-11

## Context
An assignment did not reach a teammate: the ticket file lived in the writer's branch, which may be unpushed, stale or force-pushed.
Scanning every branch is slow and still wrong. jaira has no server; the git remote is the only shared channel.
Backlog-only commits piled up on master and hid the real work in history (V3S0MW).

## Decision
One ref per ticket, refs/jaira/tickets/<id>, whose tree holds the whole ticket file; readable with no checkout.
Every write queues the full file for its ref (outbox) and pushes with --force-with-lease on the SHA last confirmed by the remote;
a lost race and an unreachable remote are distinct errors.
A ticket file enters a commit only together with code, so a reviewer sees change and reason together; the ref answers "where is it now".
A board without a usable remote keeps plain files.

## Rejected
- Scanning all branches for ticket files: expensive and wrong for unpushed or rewritten branches.
- Separate claims/log refs per ticket: one ref plus its history already carries owner, lane and movement.
- Releasing by a commit with parent that keeps the ref: the ticket would stick to its first owner.
- Pushing on the write path: TUI, merge driver and CLI must never wait for the network.
- Committing each created ticket: bookkeeping commits on master (V3S0MW).
- A central server or daemon: against the "git only" constraint.

## Consequences
Refs live outside refs/heads: the default refspec, branches and the hosting UI never show them; jaira must show owners itself.
"Clone and see the board" becomes "clone, then jaira fetch".
Two copies exist (ref and checkout file); a stale file must be detected (see ADR 0027).

## Sources
8566KF (body: "Warum Refs und nicht Branches", "Compare-and-Swap"), PP5SCQ, PXYN99, VM0A76, V3S0MW (2026-10-06), RFC7GA
