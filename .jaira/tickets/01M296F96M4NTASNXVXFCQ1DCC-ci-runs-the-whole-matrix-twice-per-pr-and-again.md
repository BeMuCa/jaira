---
id: 01M296F96M4NTASNXVXFCQ1DCC
title: "CI runs the whole matrix twice per PR, and again for ticket-only commits"
status: in-progress
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
updated-at: 2026-11-22T19:42:55Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-133696
claimed-at: 2026-11-22T19:39:21Z
---

# CI runs the whole matrix twice per PR, and again for ticket-only commits

## Definition of Done

- [x] A merged PR does not re-run the matrix that already passed on it; a commit touching only .jaira/** or **.md runs no test job; a PR whose commits are all such files still reports its required checks as complete rather than hanging; a commit that mixes ticket files with source still runs the full matrix
  proof: .github/workflows/ci.yaml: no push trigger (on: pull_request only); changes job + needs/if on shipped-manifest, web-dist, cockpit, test; ci-done aggregate with if: always(); scripts/ci-code-changed.sh classifies; TestCICodeChanged (scripts/ci_code_changed_test.go) covers ticket-only, mixed ticket+code, embedded md, umlaut names, moved-on master

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 19:39 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 19:40 · Alexander Sacharov** — decisions (none open for a person): (1) BeMuCa/jaira master is unprotected, no required checks, no rulesets (gh api 2026-11-22) - so no hanging check; still skip per job with 'if:' so a skipped job reports 'skipped' = success, should one be made required later. (2) Not a workflow-level paths-ignore: that would also drop the ticket-commits job, which exists exactly for PRs that carry only ticket files. (3) '**.md' narrowed: core/**/*.md and internal/**/*.md are go:embed inputs, lanes/*.md is read by core/lane/shipped_test.go, .claude/skills/jaira/SKILL.md is compared by core/role/jairaskill_test.go. Skippable = .jaira/**, docs/**, .planning/**, and the root README.md/AGENTS.md/CLAUDE.md, lanes/README.md. (4) Known gap: core/ticket/decode_once_test.go globs ../../.jaira/tickets and logbook; a ticket-only PR no longer runs it, the next code PR does. (5) Push trigger dropped entirely: every master commit arrives through a PR whose pull_request run already tested the merge ref.
- **2026-11-22 19:42 · Alexander Sacharov** — Found out: git diff --name-only quotes non-ASCII paths (logbook names have umlauts) - without -z such a ticket-only PR counted as code; TestCICodeChanged/a_logbook_file_with_an_umlaut fails when -z is removed (checked). Not checked: actionlint (go run of it was refused here) - YAML parses and jobs/needs/if read back correctly; the real proof is the first PR run. Why ci-done: a matrix job skipped by job-level if reports under its unexpanded name, so a required 'test (ubuntu-latest)' would hang; ci-done is the one stable check to require. No release note: CI is not observable from the binary.
