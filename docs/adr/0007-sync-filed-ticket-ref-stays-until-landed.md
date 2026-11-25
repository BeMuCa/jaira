# A filed ticket keeps its ref until it lands; the snapshot run reaps it
Status: Accepted
Date: 2026-09-11

## Context
Deleting the ref at logbook/archive time made the ticket invisible to everybody until the branch was merged,
a window as long as a review in which someone writes the same problem down again.
A first-time cloner has no refs at all, so an unworked backlog would exist only on the remote.

## Decision
logbook and archive write the final state onto the ref and leave it standing (RecordFiled = Record).
jaira snapshot writes all ref tickets under board/ on jaira/board (plumbing, every 72h), then deletes the refs of tickets
in the logbook/archive of a landing branch on the remote (settings, else remote HEAD); no resolvable landing branch deletes nothing.
A ticket not landed after landing-grace (3 days) is reported as stranded; 'jaira snapshot --drop <id>' is the manual way out.

## Rejected
- Deleting the ref immediately on logbook/archive: opens the invisible window.
- Reaping refs automatically without proof the ticket landed: recreates the window.
- Snapshot files under .jaira/tickets/: an accidental merge would add/add every ticket.
- A snapshot commit per write: refs carry their own history; hundreds of commits a day.

## Consequences
The filing clone hides its own filed tickets from the board even though the ref stands (17TJ1H).
Other clones still see a filed ticket until it lands; see the Proposed ADR 0005 (GWK31K).

## Sources
PTQ3XT, RA7PFE (2026-09-11), 17TJ1H, NWFE2B (stamp per clone), 8566KF review-gaps
