# A milestone is its own file listing its tickets, not a field on each ticket
Status: Accepted
Date: 2026-09-18
## Context
At the end of a round, unfinished work has to move to the next group.
A per-ticket marker means one pull, edit and release per ticket, each on its own ref; nobody does twenty of those.
## Decision
One file per milestone: frontmatter (colour, created-at, goal, after) and one `- <ticket id>` line per member.
jaira edits single lines and keeps comments, blank lines, order and unknown lines verbatim.
## Rejected
- A milestone field on every ticket: twenty ticket writes to regroup twenty tickets.
## Consequences
Regrouping is one edit of one file; one line per member is the smallest thing git merges.
Removing the last ticket keeps the milestone; only `jaira logbook <milestone>` takes it off the board.
## Sources
0YGWXQ, PNC854
