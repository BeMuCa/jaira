# Whether tickets, archive and logbook leave the repo and live only on refs
Status: Proposed
Date: 2026-10-22
## Context
State is shared two ways: tickets/milestones/epics on refs/jaira/*; lanes, tags, logbook
(200+ files) and archive as committed files. Archiving needs a commit; jaira update dirties trees.
## Decision
Open. Side A (today, V3S0MW): ticket file rides in the code commit, logbook/archive committed.
Side B (AMJB8G, ES625B): everything shared lives on refs; the repo holds only a ticket snapshot
in the PR; a version guard on a ref keeps the team on one jaira version.
## Rejected
None yet.
## Consequences
Side B removes merge noise but makes refs the only record; snapshot branch becomes critical.
## Sources
V3S0MW (2026-10-06), AMJB8G (2026-10-22, idea), ES625B (2026-11-23, idea), J9H38M
