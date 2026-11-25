# A person's lane stops agents only on the way forward; sending back asks why
Status: Accepted
Date: 2026-10-12

## Context
signoff (requires-human-exit) exists so the agent that did the work cannot certify it.
Agents could move a ticket out of it in no direction at all, so incomplete work stayed "ready for acceptance" with only a note.
A person in human or signoff who wanted rework could only move forward or file a follow-up ticket.

## Decision
The gate refuses only a forward move out of a requires-human-exit lane, unless the move is a key press in the board.
A move to an earlier lane needs no --force. The board's move dialog asks the person why (an empty answer keeps the prompt open),
and the note `sent back to <lane> from <lane>: <why>` lands in the same write as the move.

## Rejected
- --force for every send-back: it switches off every other refusal along with this one.
- Locking the lane both ways: a reviewer finding missing work had no way back but a note (PAAAJH).
- A follow-up ticket per send-back.
- Making the reason mandatory in the CLI as well.

## Consequences
An agent can send work back from signoff; accepting stays a person's key press (or a person's --force in a shell).
Skipping past the lane from an earlier one is still not checked (see ADR 0035).

## Sources
Q484YS, BQGM6X (decision notes 2026-11-22, DoD 4), PAAAJH
