# Replacing the binary and re-applying a board's setup are two commands
Status: Accepted
Date: 2026-10-16

## Context
'jaira update' already meant "bring this repository's board setup in step with the installed binary".
Moving to a new release needs a download, a checksum check and an in-place replace of the executable.

## Decision
'jaira self upgrade' downloads, verifies against checksums.txt and replaces the binary, then runs the NEW binary's
'jaira update --roles'. 'jaira update' downloads nothing; --roles, --lanes and --dry-run widen what it re-applies.
Homebrew and 'go install' builds and dev builds are refused with exit 3, naming the right way.

## Rejected
- One 'jaira update' that also downloads the binary (the name was taken by the board command).
- Writing the setup from the old process (its embedded files are the old version).

## Consequences
After an upgrade, lane files still need 'jaira update --lanes' or the first command's automatic refresh (NH788Q).
A pinned downgrade leaves skills and boards untouched.

## Sources
PWAF22, T8A8KM, 5Q20DC
