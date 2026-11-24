# Two places record "names that mean me": the machine alias file and the shared people list
Status: Proposed
Date: 2026-11-23
## Context
~/.jaira/identity (per machine, read by Aliases) and refs/jaira/people (shared) both say which names are one person.
core/people says a per-machine name is "the problem it solves", yet the alias file is still read and whoami still advises it.
## Decision
Open. Side A: keep both; the alias file covers one's own extra names without touching the shared ref.
Side B: the people list is the only source; drop or migrate ~/.jaira/identity, as JAIRA_USER already ignores it.
## Rejected
- (none decided)
## Consequences
Today an alias in the file makes tickets "mine" on one machine only, while the board shows them under another person.
## Sources
DVPAPE (2026-11-23), Q030WZ (2026-11-23 21:02); core/identity/identity.go AliasesPath, internal/cli/whoami.go help
