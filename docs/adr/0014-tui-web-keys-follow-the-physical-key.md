# Command keys follow the physical key, not the active layout
Status: Accepted
Date: 2026-09-13

## Context
With a Russian layout active, j/k/q and the other command keys did nothing on the board.

## Decision
The TUI asks the terminal for the physical key behind a press and maps commands from it,
while text fields keep the characters of the active layout.

## Rejected
- A per-layout key table (cannot cover layouts the package does not know).

## Consequences
Lock states (Caps/Num) do not count as a held modifier.

## Sources
Z4W2BX
