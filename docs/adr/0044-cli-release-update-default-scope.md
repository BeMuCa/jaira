# What a plain 'jaira update' brings up to date
Status: Proposed
Date: 2026-10-27

## Context
Two tickets give one command two different defaults.
Side A, PWAF22 (done): plain 'jaira update' writes only the jaira section; roles and lanes are opt-in (--roles, --lanes).
Side B, GBFM32 (signoff): plain 'jaira update' brings section, roles, cockpit and skill up to date; 'jaira update diff' previews.

## Decision
Open. Today's code follows side A; NH788Q already refreshes unedited lanes without asking.

## Rejected
- none yet

## Consequences
The agent block and the newer-release line both spell out 'jaira update --roles --lanes'; side B would change them.

## Sources
PWAF22, GBFM32, NH788Q
