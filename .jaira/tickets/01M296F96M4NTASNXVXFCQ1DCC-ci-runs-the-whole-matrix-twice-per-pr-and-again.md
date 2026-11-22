---
id: 01M296F96M4NTASNXVXFCQ1DCC
title: "CI runs the whole matrix twice per PR, and again for ticket-only commits"
status: critique
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
updated-at: 2026-11-22T19:51:37Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-827260
claimed-at: 2026-11-22T19:50:27Z
outcome-what: "ci.yaml: pull_request runs the full set on 3 platforms; push to master is back, but runs only job 'master' (ubuntu, go test ./... -race via test-shard.sh 0 1); a changes job (scripts/ci-code-changed.sh, base/head from the PR or from event.before...sha) skips every test job when only tickets/prose changed; ticket-commits runs on every PR"
outcome-why: "every merged PR ran the 3-platform matrix a second time on master, and ticket-only PRs ran it at all; a Linux run on master stays because two PRs green against the same master can merge into an untested combination (Alex 22.11.)"
outcome-resolves: "DoD 1: merged PR does not re-run the matrix (push runs one Linux job); ticket-only skips tests in PR and push; skipped jobs complete; mixed ticket+code runs everything. DoD 2: push to master = exactly one Linux -race job, none for a ticket-only commit"
review-summary: "DoD 1 holds: ci.yaml runs on pull_request only (no second matrix on the merge commit); a 'changes' job runs scripts/ci-code-changed.sh over base...head (merge base, -z, --no-renames) and shipped-manifest, web-dist, cockpit and the test matrix are gated on code=='true' via needs/if, so a ticket-only PR skips them (skipped = complete) while ticket-commits still runs; mixed ticket+code returns code=true. The skip list (.jaira/, .planning/, docs/ minus COMMANDS.md, AGENTS.md, CLAUDE.md, lanes/README.md) matches what tests read: grep finds only decode_once_test reading ../../.jaira (documented), wintrap scans .go only, the AGENTS/CLAUDE hits in core/board are temp roots. Script fails closed: set -e+pipefail stop on an unknown base before any verdict; empty diff = false. TestCICodeChanged passes locally (go test ./scripts/)."
round: 3
back-from: done
back-by: person
back-count: 1
review-gaps: "1) ci.yaml:10-12 says the PR run 'has already tested that pull request merged into master' - true only for master as of that run; master requires no up-to-date branch, so two PRs green against the same master can merge into an untested combination (bughunt finding). Trade-off DoD 1 asks for; the person decides: accept, or turn on 'require branches to be up to date' (repo setting, no code). The comment could say 'master as of that run'. 2) A non-.md file under docs/ or .planning/ (e.g. a .go file) would be skipped although go build ./... compiles it; none exists today (git ls-files), theoretical. 3) Header lines 3-6 still speak of the release PR's merge into master running no CI because of paths-ignore - now true for every merge; harmless. 4) Workflow itself unverified on GitHub (no actionlint); the first PR run is the proof."
test-verdict: "pass: go build+vet ./... and go test -race -count=1 ./... green (RC=0, 35 pkgs); DoD 1 verified in the tree (ci.yaml: pull_request only, changes job gates 4 jobs; TestCICodeChanged 15 classes x 2 bases); script exercised on real commits b0c1d6359 (logbook-only) code=false, 55e681912 code=true; unknown base rc=128 with no verdict; parallel GITHUB_OUTPUT appends both land. Workflow on GitHub itself only provable by the first PR run"
hunt-verdict: "pass: bughunt done at d41394ac6, 1 change (not a blocker: sequential merges A+B untested together without the push run - the trade-off DoD 1 asks for; note for the person), 0 pre-existing"
review-verdict: accept
review-round: 2
review-check: "read both commits whole (ci.yaml, scripts/ci-code-changed.sh, scripts/ci_code_changed_test.go); git log upstream/master..HEAD = 2 commits = payload commits; grep of *_test.go and go:embed for reads of skipped paths; go test ./scripts/ -run CICodeChanged -count=1 ok"
review-verdict-gaps: "1) ci.yaml:10-12 says the PR run 'has already tested that pull request merged into master' - true only for master as of that run; master requires no up-to-date branch, so two PRs green against the same master can merge into an untested combination (bughunt finding). Trade-off DoD 1 asks for; the person decides: accept, or turn on 'require branches to be up to date' (repo setting, no code). The comment could say 'master as of that run'. 2) A non-.md file under docs/ or .planning/ (e.g. a .go file) would be skipped although go build ./... compiles it; none exists today (git ls-files), theoretical. 3) Header lines 3-6 still speak of the release PR's merge into master running no CI because of paths-ignore - now true for every merge; harmless. 4) Workflow itself unverified on GitHub (no actionlint); the first PR run is the proof."
---

# CI runs the whole matrix twice per PR, and again for ticket-only commits

## Definition of Done

- [x] A merged PR does not re-run the matrix that already passed on it; a commit touching only .jaira/** or **.md runs no test job; a PR whose commits are all such files still reports its required checks as complete rather than hanging; a commit that mixes ticket files with source still runs the full matrix
  proof: .github/workflows/ci.yaml: push to master runs only job 'master' (Linux), the matrix runs on pull_request only (if: event_name == 'pull_request'); changes job + needs/if gate shipped-manifest, web-dist, cockpit, test, master (skipped = complete; no required checks on master); scripts/ci-code-changed.sh classifies; TestCICodeChanged (scripts/ci_code_changed_test.go) covers ticket-only, mixed ticket+code, embedded md, umlaut names, moved-on master, push form
- [x] Push auf master startet genau einen Linux-Lauf (go test ./... -race), keinen macOS/Windows-Lauf; ein reiner Ticket-Commit startet keinen (Beleg: ci.yaml + Test von ci-code-changed.sh) - Alex 22.11.
  proof: .github/workflows/ci.yaml:19 push: master; job 'master' (ci.yaml:175, ubuntu only, test-shard.sh 0 1 -race) is the only job with event push, all PR jobs gated on event_name == 'pull_request'; ticket-only push skips it via changes (ci.yaml:55 event.before...sha); TestCICodeChanged checks old master...merge commit (scripts/ci_code_changed_test.go:99)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 19:39 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 19:40 · Alexander Sacharov** — decisions (none open for a person): (1) BeMuCa/jaira master is unprotected, no required checks, no rulesets (gh api 2026-11-22) - so no hanging check; still skip per job with 'if:' so a skipped job reports 'skipped' = success, should one be made required later. (2) Not a workflow-level paths-ignore: that would also drop the ticket-commits job, which exists exactly for PRs that carry only ticket files. (3) '**.md' narrowed: core/**/*.md and internal/**/*.md are go:embed inputs, lanes/*.md is read by core/lane/shipped_test.go, .claude/skills/jaira/SKILL.md is compared by core/role/jairaskill_test.go. Skippable = .jaira/**, docs/**, .planning/**, and the root README.md/AGENTS.md/CLAUDE.md, lanes/README.md. (4) Known gap: core/ticket/decode_once_test.go globs ../../.jaira/tickets and logbook; a ticket-only PR no longer runs it, the next code PR does. (5) Push trigger dropped entirely: every master commit arrives through a PR whose pull_request run already tested the merge ref.
- **2026-11-22 19:42 · Alexander Sacharov** — Found out: git diff --name-only quotes non-ASCII paths (logbook names have umlauts) - without -z such a ticket-only PR counted as code; TestCICodeChanged/a_logbook_file_with_an_umlaut fails when -z is removed (checked). Not checked: actionlint (go run of it was refused here) - YAML parses and jobs/needs/if read back correctly; the real proof is the first PR run. Why ci-done: a matrix job skipped by job-level if reports under its unexpanded name, so a required 'test (ubuntu-latest)' would hang; ci-done is the one stable check to require. No release note: CI is not observable from the binary.
- **2026-11-22 19:43 · Alexander Sacharov** — critique: drop ci-done (speculative, one-line comment does the job, matches the header's existing 'if ever made required' pattern). Left alone: shell script over dorny/paths-filter - matches scripts/check-ticket-commits.sh; the Go test as package scripts mirrors core/shipped/history_test.go testing shipped-history.sh.
- **2026-11-22 19:44 · Alexander Sacharov** — testing: function pass on real history - b0c1d6359 (logbook-only, umlaut name) code=false; 55e681912 (shipped.sha256) code=true; this branch code=true; unknown base exits 128, prints no verdict; two parallel GITHUB_OUTPUT appends both land. Not exercisable locally: the workflow itself on GitHub (first PR run is the proof).
- **2026-11-22 19:45 · Alexander Sacharov** — bughunt: change: .github/workflows/ci.yaml:10-14 - class: code that lands on master is never tested as merged when master moved on after the pull request's last CI run; the comment says the PR run 'has already tested that pull request merged into master', but refs/pull/N/merge is master as of that run, and master requires no up-to-date branch (gh api repos/BeMuCa/jaira/rules/branches/master returns [], branch protection 404); example: PR A and PR B are both green against the same master, A merges, then B merges; A+B together never ran anywhere, and a semantic conflict between them stays on master until the next PR or the release button fails for a reason unrelated to its own change; the removed push trigger on master was what caught this; not a blocker: the release button tests the tree it releases (release-notes.yaml); fix options: put back a push run limited to master, or require PRs to be up to date before they merge; blame: 1ea6d1ee1 (CQ1DCC) removed push: branches [main, master]; read
- **2026-11-22 19:45 · Alexander Sacharov** — bughunt: done at d41394ac6 - 1 finding (1 change, 0 pre-existing) in 3 modules, first hunt; concurrency, error paths, input, permissions checked; files: .github/workflows/ci.yaml scripts/ci-code-changed.sh scripts/ci_code_changed_test.go
- **2026-11-22 19:46 · Alexander Sacharov** — testing: note: the bughunt change finding is the trade-off DoD 1 asks for (no second matrix run on master); the ci.yaml comment 'has already tested that pull request merged into master' overstates it - true for master as of the PR's last run. For the person: accept, or require up-to-date branches on master (a repo setting, not code).
- **2026-11-22 19:49 · Alexander Sacharov** — Accepted from signoff: angenommen:
- ci.yaml:10-12 says the PR run 'has already tested that pull request merged into master' - true only for master as of that run; master requires no up-to-date branch, so two PRs green against the same master can merge into an untested combination (bughunt finding). Trade-off DoD 1 asks for; the person decides: accept, or turn on 'require branches to be up to date' (repo setting, no code). The comment could say 'master as of that run'. 2) A non-.md file under docs/ or .planning/ (e.g. a .go file) would be skipped although go build ./... compiles it; none exists today (git ls-files), theoretical. 3) Header lines 3-6 still speak of the release PR's merge into master running no CI because of paths-ignore - now true for every merge; harmless. 4) Workflow itself unverified on GitHub (no actionlint); the first PR run is the proof.
- **2026-11-22 19:49 · Alexander Sacharov** — decision (Alex 22.11., give back): nach dem Merge auf master läuft weiter ein kurzer Lauf nur auf ubuntu-latest (go test ./... mit -race, ~2 min) als Absicherung gegen ungetestete Kombinationen zweier grüner PRs; macOS/Windows nur im PR. ci.yaml-Kommentar 10-12 entsprechend korrigieren.
- **2026-11-22 19:50 · Alexander Sacharov** — round 3 plan (DoD 2): push trigger on master back, paths-ignore notes; changes job reads base/head from pull_request or from push (event.before...sha); new job 'master' runs only on push: ubuntu, scripts/test-shard.sh 0 1 -race; ticket-commits, shipped-manifest, web-dist, cockpit, test matrix gated to pull_request. Separate job rather than a dynamic matrix: a fromJSON matrix switch is harder to read than ~15 duplicated setup lines. Test: TestCICodeChanged also checks the push form (old master...merge commit).
