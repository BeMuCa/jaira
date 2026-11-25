# A filed or archived ticket leaves other boards when it lands, or at once
Status: Proposed
Date: 2026-09-11

## Context
'jaira archive' and 'jaira logbook' move the file and write the unchanged content to the ref, so other clones keep showing it
as a ref-only card until the file lands in a landing branch and the snapshot run removes the ref.
GWK31K (signoff) wants an archived ticket gone from every board at the next fetch.

## Decision
Open. Side A (RA7PFE, in code): the ref keeps the final state and stays until the ticket lands; only the filing clone hides it at once.
Side B (GWK31K, in signoff): the ref marks the ticket filed, every board hides it on the next fetch, and restore brings it back everywhere.

## Rejected
- Deleting the ref at archive time: reopens the invisible window RA7PFE closed.

## Consequences
A trades visibility of "finished, waiting to land" for a stale card elsewhere.
B needs a "filed" marker on ticket refs, like the milestone "status: filed".
AMJB8G (tickets fully on refs) would change the question.

## Sources
RA7PFE (2026-09-11), GWK31K (2026-10-15, signoff), 17TJ1H, AMJB8G
