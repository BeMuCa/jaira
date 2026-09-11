---
id: 01M2943H9MEYF2723ZSVJNYE3T
title: Tests write their state into the real ~/.jaira
status: signoff
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "A full test run leaves no directory under the developer's own ~/.jaira; every test that opens a store points JAIRA_HOME at its own temp dir"
context: |-
  A test run pollutes the real ~/.jaira/state. There are 3139 directories in it on this machine and 3115 of them are named 001-*, 002-* or 003-* - the names Go gives subdirectories of t.TempDir(). Only about 24 belong to real checkouts.

  Where it comes from: Store.stateDir (core/ticket/store.go) falls back to $HOME/.jaira when JAIRA_HOME is unset, and tests that build a ticket.Store without t.Setenv("JAIRA_HOME", ...) therefore write sessions/ and locks/ into the developer's own home, keyed by a temp path that stops existing when the test ends.

  Harmless in content - the leftovers are empty sessions/ and locks/ directories, 38 MB in total - but they bury the ~24 real entries among three thousand dead ones, which is how a stale-directory count got misread as a worktree problem while diagnosing NWFE2B.

  Two parts: stop the leak, then clear what is already there. Deleting the existing ones is safe for any directory whose checkout no longer exists; sessions and locks are per-checkout and short-lived by design.
definition-of-done: go test ./... with HOME pointed at an empty directory creates nothing under it; the tests that leaked are identified and isolated; ~/.jaira/state holds only directories whose checkout still exists
tags:
  - concurrency
blocked-by: []
commits: []
created-at: 2026-09-11T20:56:18Z
updated-at: 2026-09-11T21:01:21Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1358160
claimed-at: 2026-09-11T20:56:31Z
outcome-what: internal/cli now has a TestMain (internal/cli/main_test.go) that points JAIRA_HOME at a temp directory for the whole package and removes it afterwards. Swept ~/.jaira/state of every directory holding no regular file.
outcome-why: "Tests that opened a store without setting JAIRA_HOME wrote sessions/ and locks/ into the real ~/.jaira, keyed by a t.TempDir() path that stops existing when the test ends: 34 per run of that package, 3115 collected on this machine, burying the two dozen entries that belong to real checkouts."
outcome-resolves: "A per-package scan that counts ~/.jaira/state before and after each package now reports no growth anywhere; go test ./... green; ~/.jaira/state is down from 3173 directories and 38 MB to 23 and 636 KB, and the 23 are existing checkouts."
review-summary: "One TestMain replaces a rule nobody could keep: 10 of the 26 test files in internal/cli set JAIRA_HOME and 16 did not, so forgetting was the default and the damage landed in the developer's home rather than in the test. Nothing outside test code changes, so there is nothing here for a user to notice and no NOTES.md line. The leak was measured per package rather than guessed - internal/cli was the only one."
review-gaps: "Nothing stops a future package from leaking the same way; the guard is per package, not global. A stricter version would make the store refuse the real home when it can see it is running under a test binary, as core/bgrun already does for spawning - not done here, because it hides the mistake instead of showing it. The sweep script lives in the scratchpad and is not in the repo: with the leak fixed there is nothing left to sweep periodically."
review-verdict: "Ready. Test-only change, measured before and after, no production code touched."
review-check: "1. Run: ls -A ~/.jaira/state | wc -l  -- note the number. 2. Run: go test -count=1 ./internal/cli. 3. Run: ls -A ~/.jaira/state | wc -l  -- the same number; before this change it grew by 34. 4. Run: go test ./...  -- green."
---

# Tests write their state into the real ~/.jaira

## Definition of Done

- [x] go test ./... with HOME pointed at an empty directory creates nothing under it; the tests that leaked are identified and isolated; ~/.jaira/state holds only directories whose checkout still exists

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-11 21:00 · Alexander Sacharov** — Only internal/cli leaked - measured by running each package on its own and counting ~/.jaira/state before and after, because HOME cannot be redirected in this session. 34 directories per run of that package; every other package already isolates JAIRA_HOME or never opens a store.

Fixed with one TestMain for the package rather than a t.Setenv in each of the 26 test files: 10 of them already set it and 16 did not, so the leak is the default, and a test that forgets is silently wrong somewhere nobody looks. A per-test t.Setenv still overrides it.

The sweep deleted every directory under ~/.jaira/state holding no regular file at all - no session, no lock, no stamp, nothing to lose. 3150 went, 23 real checkouts stayed, 38 MB down to 636 KB. scratchpad/sweep.sh is the script; it is not in the repo, since the leak it cleans up is now fixed.
