# The cockpit polls no pull request; it hands each one to the teamlead once
Status: Accepted
Date: 2026-11-22

## Context
jaira has no daemon. A PR's comments, CI and review requests were tracked in the teamlead's head.

## Decision
The cockpit only tells the session once per PR to run `jaira pr --watch --number N` under a Monitor.
It learns of a PR from the pr lane it started (up to 20 one-minute looks) or from a report naming an open PR of this repository.

## Rejected
- The cockpit polls the forge for every PR - a second watcher beside `jaira pr --watch`.

## Consequences
"Already handed over" is kept in the mod store per project and PR number, so a restart does not repeat it.

## Sources
CRASZS, JG3REW
