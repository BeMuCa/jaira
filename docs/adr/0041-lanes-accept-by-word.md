# How does a person's acceptance given in chat reach the board?
Status: Proposed
Date: 2026-10-23

## Context
Leaving signoff forward needs a key press in the board (Interactive); a CLI move is refused even at a person's word.
Today the agent runs `jaira move <id> --to done --force`, which prints the overridden refusals and records nowhere who accepted.

## Decision
Open. Side A (today): a person's acceptance in a shell is --force; the board key press is the only non-forced path.
Side B (BX1MNM, in signoff): `jaira accept <id> --by <person>` moves out of signoff, records who and when, and reads as success.

## Rejected
- None yet.

## Consequences
Side B gives an agent a command that opens the human exit; its safety rests on the agent quoting a real person.

## Sources
BX1MNM, Q484YS, HVHZ7Q, core/gate/gate.go (CodeNeedsHuman)
