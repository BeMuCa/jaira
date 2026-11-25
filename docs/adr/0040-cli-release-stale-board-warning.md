# How often the stale-board warning is printed
Status: Proposed
Date: 2026-10-23

## Context
While a board's version stamp differs from the binary, every command prints "set up by an older version — run 'jaira update'".
Agents filter it away and people stop reading it; callers that drop stderr would consume a once-only hint unseen.

## Decision
Open (RNQRJB in human): once per board and version, or as today on every command until 'jaira update' runs.

## Rejected
- none yet

## Consequences
A once-only hint needs a place where it is surely seen (e.g. the cockpit), or it is lost to --json pollers.

## Sources
RNQRJB, NH788Q
