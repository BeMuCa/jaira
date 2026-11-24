# Roles ship inside the binary and an install never overwrites an edit
Status: Accepted
Date: 2026-09-14

## Context
Role prompts lived in one person's ~/.claude/skills; a teammate who cloned got a board and nobody to drive it.
An old install that nobody edited must still update without --force.

## Decision
core/role/builtin is the only source; `jaira roles install` writes it out, `jaira update --roles` refreshes it.
A file is compared three ways: identical, a shipped version (core/shipped, replaced), or edited (left alone, exit 3).

## Rejected
- Maintaining the roles in ~/.claude/skills and copying them by hand.
- Writing into every directory a harness reads (one command registered twice; E75R4Q).
- Comparing by existence only (os.Stat), which cannot tell a finished install from an edit.

## Consequences
Every shipped file version must be recorded in core/shipped/shipped.sha256, or an old install looks edited.
A role's directory name is its command name, so a rename is a breaking change.

## Sources
C9V7ZV, PMF635, 7KX89C, E75R4Q (decision 2026-10-13), 0K8Y6S, PWAF22
