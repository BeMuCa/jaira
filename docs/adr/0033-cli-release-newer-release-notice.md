# A newer release is announced on four commands, and the user decides
Status: Accepted
Date: 2026-10-16

## Context
Agents never see the TUI's release pill, so a session ran an old binary without anyone knowing.
A notice on every command would be noise for the scripts and agents that parse the CLI.

## Decision
A detached background check refreshes ~/.jaira/update-check.json at most once an hour per machine.
'jaira next', 'resume', 'claim' and 'move' print one stderr line when the cache knows a newer release;
the agent runs the two previews, shows the user what they report and upgrades only on a yes.

## Rejected
- A daily check (a release reached running sessions a day late).
- A notice on every command (noise on stdout-parsing scripts).
- The agent upgrading on its own.

## Consequences
Dev builds and JAIRA_NO_UPDATE_CHECK=1 never check or print.
The line is plain stderr text, so it can drown in a chat (open: 4C6R8X).

## Sources
T8A8KM, EDSQZP, 4C6R8X
