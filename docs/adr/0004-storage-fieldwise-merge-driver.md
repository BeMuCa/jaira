# Ticket files merge field by field through a registered git merge driver
Status: Accepted
Date: 2026-09-11
## Context
Two people moving one ticket both rewrite the status: line, the most frequent operation.
Union merge of text would produce two status: lines, invalid YAML.
## Decision
.jaira/.gitattributes maps tickets/*.md to merge=jaira and /tags to merge=union.
Lists union, status takes the lane further along, other scalars by recency.
Only prose fields can conflict; the losing value is parked in conflict-theirs-<field>.
An empty base (add/add) is merged as a base with no fields.
## Rejected
- Plain line merge: conflicts on every concurrent move.
- Rejecting an empty base: git kept ours silently and the other side was lost (NQCKF5).
## Consequences
The one non-obvious thing jaira does to a repository: the driver must be registered per clone.
## Sources
core/merge (2026-08-11), NQCKF5 (2026-09-11), 79GEPW (tags union)
