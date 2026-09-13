---
id: 01M296F96M4NTASNXVXFCQ1DCC
title: "CI runs the whole matrix twice per PR, and again for ticket-only commits"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "One CI run per change, and none at all when a commit touches nothing the binary is built from"
context: |-
  Two wasted CI runs on every change.

  First: the same three-platform matrix runs twice. .github/workflows/ci.yaml triggers on 'pull_request' and on 'push: branches: [main, master]'. A merged PR is a push to master, so ubuntu, macos and windows all run go test -race a second time on the commit they just passed on.

  Second: a commit that touches only .jaira/ or a markdown file still runs the full matrix. That is most commits on this board - a ticket moves lane, a note is written, NOTES.md gains a line - and none of it can change what go test does.

  Known: 'paths-ignore' on the triggers is the mechanism, and the trap is that a required status check that never runs leaves a PR pending forever. If any of the three test jobs is a required check on BeMuCa/jaira, the skip has to be done so the check still reports - the usual way is a job that runs and exits early rather than a filter that stops it existing.

  Do not skip on path alone without checking that: .jaira/ files ride in the same commit as the code they belong to, by this project's own rule, so a mixed commit must still run the tests.
definition-of-done: "A merged PR does not re-run the matrix that already passed on it; a commit touching only .jaira/** or **.md runs no test job; a PR whose commits are all such files still reports its required checks as complete rather than hanging; a commit that mixes ticket files with source still runs the full matrix"
tags:
  - release
blocked-by: []
commits: []
created-at: 2026-09-11T21:37:41Z
updated-at: 2026-09-11T21:37:41Z
---

# CI runs the whole matrix twice per PR, and again for ticket-only commits

## Definition of Done

- [ ] A merged PR does not re-run the matrix that already passed on it; a commit touching only .jaira/** or **.md runs no test job; a PR whose commits are all such files still reports its required checks as complete rather than hanging; a commit that mixes ticket files with source still runs the full matrix

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

