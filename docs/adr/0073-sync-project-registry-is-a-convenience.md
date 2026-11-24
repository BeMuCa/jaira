# projects.json is a convenience list, one entry per repository
Status: Accepted
Date: 2026-11-23

## Context
Every folder jaira ran in became a project: worktrees of one repo and agents' temp boards filled the launcher.
Two jaira processes starting at once interleaved their writes and left an unparsable file that read as an empty list.

## Decision
A board in a linked worktree is recorded as its main checkout (read from .git and commondir, no git process);
boards in the temp directory are not recorded, old worktree entries fold into their repository; at most 20 entries.
Remember writes under a lock file and replaces the file whole; an unparsable file is copied to projects.json.broken[-<time>] first.

## Rejected
- Listing every worktree as a project.
- Treating the registry as a source of truth: a project is just a directory with .jaira.
- Copying a broken file aside in Load: a read alone loses nothing, so only the write path does it.

## Consequences
Deleting projects.json loses only the shortcut list.
R947VK (no vanished paths) and AYBDJ2 remain open for the rest.

## Sources
ZXJ72E, BPR7C6 (decision 2026-11-23 13:30), R947VK (open, same goal), AYBDJ2
