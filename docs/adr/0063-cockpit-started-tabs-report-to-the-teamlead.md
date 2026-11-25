# Tabs the cockpit starts report to the ticket's teamlead session
Status: Accepted
Date: 2026-11-22

## Context
A worker the cockpit started after Give back ran with --no-parent and asked its questions in its own tab, which nobody watched.

## Decision
A new worker reports to the session the ticket's newest `parent session:` note names, else this session's messaging socket, else none.
The pr lane and Discuss chats get this session as parent when its socket is live.

## Rejected
- A worker the cockpit starts stays in its own tab with --no-parent (7V1XPB, 2026-10-23) - reversed by T5QN9P.

## Consequences
A parent that no longer listens is skipped (spawn.sh exit 2) and the next one is tried.

## Sources
T5QN9P (decision 2026-11-22 18:33), 7V1XPB, QTQW15
