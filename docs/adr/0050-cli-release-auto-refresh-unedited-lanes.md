# The first command after an upgrade refreshes unedited shipped files itself
Status: Accepted
Date: 2026-11-18

## Context
After 'jaira self upgrade' every repository needed 'jaira update --lanes' by hand, and nobody typed it.
jaira must not overwrite a lane, role or skill a person edited.

## Decision
The first command under a new version replaces lane files whose hash is a shipped version (core/shipped/shipped.sha256),
renews the jaira block only where a file already carries it, and says so in one stderr line.
It has its own stamp (lanes-refreshed); the version stamp stays with 'jaira update' and its release notes.

## Rejected
- Leaving it to 'jaira update --lanes' in each repository.
- Creating CLAUDE.md/AGENTS.md where none exists (an unasked command must not add files).
- Replacing edited lanes.

## Consequences
shipped.sha256 must list every file version of every v* tag; CI enforces it (0K8Y6S).
--dry-run skips the refresh; a half-failed refresh names what it already wrote.

## Sources
NH788Q (decision: 2026-11-18 17:09, 17:14), PWAF22, 0K8Y6S
