# Only a lane's declared send-back word in review-verdict is a rail
Status: Accepted
Date: 2026-11-21

## Context
A review that wrote "send back" still saw its ticket land in signoff; nothing read the verdict.

## Decision
A lane declares verdict-accept and verdict-send-back words; a review-verdict starting with a send-back word refuses every agent move past that lane,
blocked included, and names where to send the ticket back. Any other verdict, free text included, passes.

## Rejected
- Refusing unknown verdict words: the edges kept growing (markdown, numbering, dashes) - "edges grow, pick the smaller design".
- Reading the verdict only on the move out of review: blocked -> signoff walked past it.

## Consequences
A malformed send-back (e.g. "**send-back**") is a known limit; the lane prompt asks for the bare word.
A move made in the board is a person deciding and is not refused.

## Sources
F1V59P (decision notes 2026-11-21 20:10 and 20:40)
