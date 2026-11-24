# Language of the notes the cockpit writes onto tickets
Status: Proposed
Date: 2026-11-21

## Context
texts.ts and F80QZT: every cockpit text is English, the project language.
The decision lines Give back and Accept write use German headings: "beheben:", "Folge-Ticket:", "angenommen:", "dazu:".

## Decision
Open. Side A: notes the cockpit writes follow the English rule (F80QZT decision 2026-10-16, texts.ts header).
Side B: the headings stay in the words the ticket asked for (MFYM1A, board.ts HEADINGS), and decisionsOf reads them back.

## Rejected
(none yet)

## Consequences
Changing the words breaks decisionsOf on notes already written, unless both spellings are read.

## Sources
MFYM1A, F80QZT, 2V50Y5
