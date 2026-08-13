---
status: complete
---

# 260813-tuk — boards visible and switchable everywhere

The executor died to a connection loss after committing all four tasks but before
writing this file. Written by the orchestrator from the repository state, which
was verified independently rather than taken from the crashed agent's report.

## Commits

- `f6d6e5c` boards are known from the first frame and 1-9 switches from either view
- `60816cf` mark boards that have an agent working in them right now
- `62316e7` the sign-off lane displays as Human Review
- `e7ba5fe` do not warn about a lane override that changes nothing

## Verified after the crash

- `go build ./... && go vet ./... && go test ./... -race -count=1` — all 11
  packages ok, 0 FAIL
- `git status` — no uncommitted code
- `jaira lanes` on a real board: `signoff  Human Review`
- Override suppression proved end to end: `jaira lanes use review` into a fresh
  project produces **no** warning

## One thing the check turned up

The user's own board still warned after the fix, correctly: task 3 rewrote
"sign-off" to "human review" inside the built-in review prompt, so their project
copy — made before the rename — genuinely differed. Refreshing the copy cleared
it. Worth knowing: **renaming anything inside a built-in lane makes every existing
project copy of it start warning.** There is no migration for that today; the
user has to re-copy or accept the warning.

## Unverified

The interactive checks: that the numbered boards render, that 1-9 switches, and
that a live board is marked, all need a terminal. Covered by tests at the model
level only.
