# Give back reuses the old worker only while its tab lives, its cache is warm and its context is small
Status: Accepted
Date: 2026-10-12

## Context
A ticket given back can go to the worker that carried it (cheap, has context) or to a new one (clean context).
A fixed token limit (150k) would cut off workers on large tasks.

## Decision
Reuse when the tab is alive, the prompt cache is warm and context is below workerContextShare of the window (default 50%).
A worker that is working, or blocked at a dialog, is never replaced; a new worker is started only after the old tab is closed.

## Rejected
- Fixed 150k-token limit - wrong for large windows and large tasks.
- Always a new worker - throws away a warm cache.

## Consequences
The same rule (chooseWorker) answers Give back and the `workers` tool, so teamlead and cockpit agree.

## Sources
QRHB35 (decision 2026-10-12 15:34), QTQW15, M31BYV, W80HNH, P122GE
