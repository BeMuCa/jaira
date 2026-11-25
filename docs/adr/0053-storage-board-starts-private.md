# A board starts private in info/exclude; sharing is an explicit command
Status: Accepted
Date: 2026-11-20
## Context
Committing tickets publishes them to everyone who can read the repository.
Writing the exclude into the tracked .gitignore left every new worktree with a stray change.
## Decision
jaira init excludes .jaira/ in <git-common-dir>/info/exclude (.gitignore only outside git).
jaira share removes that entry, ignores lanes/ and milestones/, writes .gitattributes; no commit.
jaira update never rewrites the ignore state of a board that is already shared.
## Rejected
- Shared by default: the tool cannot know whether notes may be published.
- Entry in the tracked .gitignore: dirty worktrees, refused guarded pushes.
- update re-writing the ignore line: silently hid new tickets of shared boards (ETR0PX).
## Consequences
info/exclude is per clone: after share --undo a teammate runs jaira update once.
## Sources
ETR0PX (2026-09-03), 8W1R94, NNV39Y (decision 2026-11-20), 689E0T
