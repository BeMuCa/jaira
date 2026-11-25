# News that arrives in the background is a status line, never a screen
Status: Accepted
Date: 2026-09-11

## Context
A ref fetch can bring tickets assigned to me, milestone warnings or errors while the person is typing.
A screen of its own would interrupt them and swallow the keystroke they were typing.

## Decision
Arrivals ("<handle> is yours: <title>"), fetch warnings and errors go to the board's status line.
A desktop notification goes out as well where notifications are enabled.

## Rejected
- A message screen per arrival (the behaviour before YYR59F).
- Dropping fetch errors silently (fixed in S1SA89).

## Consequences
The fetch interval is one setting (`fetch-every`) shared with the CLI.

## Sources
YYR59F, S1SA89
