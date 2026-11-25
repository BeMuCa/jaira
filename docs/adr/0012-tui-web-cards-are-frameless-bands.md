# A card is a frameless filled band with its tag colours in a left bar
Status: Accepted
Date: 2026-09-11

## Context
Boxed cards cost two rows each, so a lane held few tickets, and stacked boxes ran into each other.
Berk had asked for separate frames (81XRXX) and a coloured left edge (VS5DFW); both were overtaken.

## Decision
A card is three rows of filled background with no frame; the background alternates from card to card.
Tag colours sit in a one-column bar at the left, up to three tags (S1VM40); the selected card gets a fill shade.

## Rejected
- Every card boxed in its first tag's colour, neutral box without tag (NJPQWE, AXTFG3).
- Stacked boxes sharing one border row (NJPQWE outcome) or two separate border rows (81XRXX).
- White frame with only the left edge in tag colour, full colour frame when selected (VS5DFW).
- Selected card filled in its tag colour with c toggling to neutral (XK4124): c now runs the agent band (C8V6DQ).

## Consequences
cardHeight is 3 and the lane budget counts rows per card exactly.
State marks must be glyph plus word; colour alone carries no meaning.

## Sources
QQ3EX4, S1VM40, ANX7Y5, NJPQWE, AXTFG3, XK4124, C8V6DQ, archive 81XRXX, VS5DFW, archived notice ticket for Berk (01M293C00…); note 2026-09-11 "UEBERHOLT durch QQ3EX4".
