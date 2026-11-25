# An epic is a description plus a ticket list on its own ref, not a ticket and not a milestone
Status: Accepted
Date: 2026-10-14
## Context
A feature can span 40 tickets and needs a place to say what it is for and how far it is.
Its member list is written while agents work member tickets, so it must never share a file with what a branch carries.
## Decision
An epic lives only on refs/jaira/epics/<name> as <name>.md; no working file; members are listed in the epic, never on tickets.
Progress counts members plus every follows:/parent: descendant; done means terminal lane or logbook. Finished epics are closed, not deleted.
## Rejected
- Mirror epics as .jaira/epics/<name>.md like milestones (worker decision, overruled by Alex 2026-10-14 08:49).
- Epic as a ticket with lanes and a definition of done.
- Deleting a finished epic instead of `status: closed` (182JT1).
## Consequences
Epics need a repository with a remote; a board without one has none (exit 3).
## Sources
F914QH (decisions 2026-10-14 08:47, 08:49), 182JT1 (2026-11-02), CRDH6B, PG9P9Y
