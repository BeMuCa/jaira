# Ticket writes are queued in a local outbox and sent after the command
Status: Accepted
Date: 2026-09-11

## Context
Claiming must work offline, and every write funnels through Store.Mutate/Create, which the TUI and the git merge driver call too.
A synchronous push there would hang the board and make the merge driver push in the middle of a merge.

## Decision
The write path only files the whole new content, with its lease, in a per-working-tree outbox outside the repository.
A process that queued something flushes once after the command (fetch first, oldest first); the TUI flushes before each background fetch.
Offline keeps the queue; a lost ticket race drops the entry and names the winner; milestone and epic losers are merged and requeued.

## Rejected
- Synchronous push inside the write path: blocks TUI and merge driver on the network.
- "claim needs the network, everything else offline" (proposed in the 8566KF body): replaced by the asynchronous outbox.
- Replaying diffs instead of whole files: could apply onto a state nobody reviewed.

## Consequences
Read commands never touch the network. A write can sit "unsent" until the next command; the card says so.
'jaira pull', 'release' and 'reject' are the only synchronous ref writes: with no route there is no ref to take over.

## Sources
8566KF (DoD plan item "Entschieden: asynchron"), 17TJ1H, MMKENA (decision 2026-11-23), S1SA89, core/outbox/outbox.go package doc
