# Refs are fetched in the background on an interval, never on a read's path
Status: Proposed
Date: 2026-09-11

## Context
A CLI-only reader never fetched, so a ticket assigned to them did not exist until they typed 'jaira fetch'.
A read command must never wait for a remote.

## Decision
Accepted part (PP5SCQ): after every command a detached 'jaira fetch --json' is spawned when the last one is older than fetch-every;
stamp per repository (NWFE2B); JAIRA_NO_FETCH disables it; the TUI uses the same interval.
Open part: the default is 10 minutes in code (PP5SCQ); D1S55M (signoff) proposes 2 minutes so teammates' changes arrive sooner.

## Rejected
- Fetching inside 'jaira list' and other reads: makes reads wait on the network.

## Consequences
At 2 minutes, measured cost is about 0.75% of a core per open board (D1S55M context).

## Sources
PP5SCQ, NWFE2B, 17TJ1H, D1S55M (open)
