# Where a board's lane files live: private per machine, or committed with the repo
Status: Proposed
Date: 2026-09-14
## Context
Two rules exist for one thing. jaira share ignores /.jaira/lanes/ (machine-scoped, BNZERQ).
Named boards under .jaira/boards/<name>/lanes/ are not ignored, so they get committed (S0EG0M).
This repository tracks .jaira/lanes/ anyway; JFW22A: fresh clones miss critique/optimize/testing.
ES625B proposes a third answer: lanes on refs.
## Decision
Open. Side A: default lanes stay private per machine. Side B: lanes travel with the repo
(committed, or on refs per ES625B), so every clone runs the same pipeline.
## Rejected
None yet.
## Consequences
Until decided, default and named boards behave differently on jaira share.
## Sources
BNZERQ (2026-09-03), JFW22A (2026-09-14, open), S0EG0M (2026-10-17), ES625B (open), 28C0QV
