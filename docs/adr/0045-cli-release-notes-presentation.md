# How release notes are shown and marked
Status: Proposed
Date: 2026-10-31

## Context
Three open tickets in signoff describe the release-notes view and the line format differently.
Side A, GK1QGT: one version per page, newest first, left/right pages between versions.
Side B, MXMDFH: one popup, features first, then the rest under 'Others', which needs a third line marker.
CZWKNY: a hook line on top of every note file, shown first by 'jaira update' and the popup.

## Decision
Open. Today's code knows only '- ' and '- ! '; any new marker changes the documented line-scan format in CLAUDE.md.

## Rejected
- none yet

## Consequences
Whichever wins must keep old note files readable and update the CLAUDE.md format rules in the same change.

## Sources
GK1QGT, MXMDFH, CZWKNY, AR97K9
