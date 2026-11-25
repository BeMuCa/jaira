# Further boards are .jaira/boards/<name>/lanes/, chosen by a board: field
Status: Accepted
Date: 2026-10-17
## Context
One repository needs boards with different lanes for different subjects.
Existing boards and tickets must keep working without migration.
## Decision
Today's board stays .jaira/lanes/ and is called default only on the CLI and TUI.
A named board exists exactly when .jaira/boards/<name>/lanes/ holds a lane file.
A ticket names its board in board:; absent means default; moving is one line, not a rename.
## Rejected
- Auto-creating a board on an unknown --board: a typo would leave a board behind (exit 2).
- One directory of tickets per board: moving a ticket would be a rename/modify conflict.
## Consequences
list and next without --board show the default board only.
## Sources
S0EG0M (decisions 2026-10-17), 8NHD8Q
