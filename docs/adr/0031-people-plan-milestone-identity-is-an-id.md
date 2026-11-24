# A milestone is identified by a fixed six-character ID; its name is only a label
Status: Accepted
Date: 2026-10-14
## Context
Milestones were keyed by name: file <name>.md, ref per name, after: lines naming predecessors by name.
A rename, or two people renaming or creating at once, could orphan a copy or lose a milestone on another clone.
## Decision
File `<name>-<ID>.md`, ref refs/jaira/milestones/<ID>, `after:` stores IDs; commands accept name, ID or `<name>-<ID>`.
Rename is an ordinary write on the one ref; a filed milestone no longer holds its name.
## Rejected
- Keep names as keys and migrate refs and after: lines on every rename.
- Leave unknown after: tokens verbatim during migration: clones would write different bytes for one old file (2026-10-13).
## Consequences
Two milestones may share a name: warned on stderr, and a writing command given the bare name exits 3.
One ID in two files after a rename is refused by every write (8TT890).
## Sources
EBTX6X (critique decision 2026-10-13 22:28), E70WYN, 8TT890 (decisions 2026-11-22)
