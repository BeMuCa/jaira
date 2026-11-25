# Filing finished work into the logbook is a person's act, never a lane's side effect
Status: Accepted
Date: 2026-09-18

## Context
done was made a doorway that filed every arriving ticket at once (WXQ9PT), and capped at ten (SGPDYK).
One automatic filing swept forty-nine other people's tickets into one commit named after one handle.

## Decision
Finished tickets wait in the terminal lane; `jaira logbook <id>` or `--all` files them into
logbook/<initials>-<yyyymmdd>/ with commits stamped. archive = not worked, any lane;
reject = archive plus a reason (person only); delete = the one irreversible act.
The shipped done lane carries neither logbook-on-entry nor holds.

## Rejected
- done as a doorway (logbook-on-entry, WXQ9PT, superseded): filing happens when someone accounts for the work.
- A fixed cap on done (holds: 10, SGPDYK): moves files nobody asked to move.
- A single command for finished and abandoned work: the logbook must stay the record of done.

## Consequences
logbook-on-entry and holds remain lane options any board may still declare (core/lane/settle.go).
The board only hints when many done tickets pile up; the folder is who-finished-what-when.
The commits list is stamped at filing and, since 8FE2EG, also written by the move into done.

## Sources
WXQ9PT (2026-09-11), 74VM40 (2026-09-18), SGPDYK, 2MM32Y, A0C081 (decisions 2026-10-19..21), MJ3PVV, 8FE2EG
