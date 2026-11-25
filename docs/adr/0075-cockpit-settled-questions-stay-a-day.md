# A settled question leaves "Waiting for you" by itself or by x, and stays a day under Answered
Status: Accepted
Date: 2026-11-24

## Context
The cockpit's "Waiting for you" pane kept a question until an option was picked there.
A question answered elsewhere (in a chat, by moving the ticket) stayed listed as open (CSQXKR).

## Decision
Sascha 2026-11-24: both ways close a question - by itself once its ticket leaves the lane it was asked in, and by hand with x.
A closed question is not removed at once: it stays 24 hours under an "Answered" line, dimmed, with when it was asked and answered.
Questions show in the words they were asked; t translates them all, as the Acceptance pane does (ADR 0068).

## Rejected
- Only by hand (x): a question settled by a lane move would wait for a press nobody remembers.
- Only by itself: a question answered in another chat, with no lane move, never leaves.
- Removing a settled question at once: the person cannot tell what is new from what is already answered.

## Consequences
The cockpit compares each asked-about ticket's lane with the board on every refresh (60 s) and asks `jaira show` for one it does not list.
Past 30 threads the cockpit drops those touched longest ago, never one whose question still waits.

## Sources
CSQXKR (decision 2026-11-24 10:26), 1M3X77, core/mod/builtin/jaira-cockpit/hooks/threads.ts, register.tsx settleMoved
