# The version line is always shown, and a dev build says "jaira dev"
Status: Accepted
Date: 2026-09-11

## Context
DNAEPN (2026-09-03) made the version line silent on dev builds, because "dev" answers nothing about updates.
Berk then could not tell which of several local binaries was running (P1AE82).

## Decision
The board and launcher show the running version top left on every build: "jaira dev" on a dev build,
"jaira X.Y.Z · up to date" when the cache knows no newer release, and a pill "^ <version>" when it does (T8A8KM).

## Rejected
- Silence on dev builds (DNAEPN).
- Version only in the footer.

## Consequences
The release cache is polled hourly and re-read by an open board every minute (T8A8KM).

## Sources
P1AE82, DNAEPN, T8A8KM
