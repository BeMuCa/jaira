---
id: 01M296F96M4NTASNXVXFCQ1DCC
title: "CI runs the whole matrix twice per PR, and again for ticket-only commits"
status: done
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
commits:
  - 1ea6d1ee1371e070bc7d9bcd3a6fff7a5c38809e
  - d41394ac684e8b074cb0cdf765e17a8516ed0d6b
  - ecb88f76f438633cad1927c1587400f71bfcd2cf
  - 5306282770f8a01b0250c099348a1d8a1d330156
created-at: 2026-09-11T21:37:41Z
updated-at: 2026-11-22T19:59:57Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-827260
claimed-at: 2026-11-22T19:50:27Z
outcome-what: "ci.yaml: pull_request runs the full set on 3 platforms; push to master is back, but runs only job 'master' (ubuntu, go test ./... -race via test-shard.sh 0 1); a changes job (scripts/ci-code-changed.sh, base/head from the PR or from event.before...sha) skips every test job when only tickets/prose changed; ticket-commits runs on every PR"
outcome-why: "every merged PR ran the 3-platform matrix a second time on master, and ticket-only PRs ran it at all; a Linux run on master stays because two PRs green against the same master can merge into an untested combination (Alex 22.11.)"
outcome-resolves: "DoD 1: merged PR does not re-run the matrix (push runs one Linux job); ticket-only skips tests in PR and push; skipped jobs complete; mixed ticket+code runs everything. DoD 2: push to master = exactly one Linux -race job, none for a ticket-only commit"
review-summary: "ci.yaml: the matrix (ubuntu -race, macOS, Windows x3) plus shipped-manifest, web-dist and cockpit now run only on pull_request; a push to master runs a single new job 'master' (ubuntu, scripts/test-shard.sh 0 1 -race) to catch two PRs green apart but not together. A new 'changes' job runs scripts/ci-code-changed.sh over the PR range (base...head) or the push range (before...sha); when every changed file is .jaira/**, .planning/**, docs/** (except docs/COMMANDS.md), AGENTS.md, CLAUDE.md or lanes/README.md, all test jobs are skipped. ticket-commits still runs on every PR. The script uses the merge base, NUL-separated names (umlaut logbook files), and fails instead of printing a verdict when git fails; TestCICodeChanged covers 15 file sets over PR and post-merge ranges."
round: 3
back-from: done
back-by: person
back-count: 1
review-gaps: |-
  note: a force-push to master leaves github.event.before outside the checkout: changes fails red and master is skipped (loud, not silent) - fall back to HEAD~1 when git cat-file -e "$before" fails
  note: the master job runs go test only, not the go vet ./... and go build ./cmd/jaira the matrix runs - add both steps if a vet-only break on master matters
  note: actionlint was not run on ci.yaml (not installed here); the first PR run on GitHub is the real proof that needs/if resolve as intended
test-verdict: "pass: round 3 over ecb88f76f only - go test -race ./scripts/ RC=0 (TestCICodeChanged incl. push form old-master...merge); DoD 2 read in ci.yaml: push:[master], only job 'master' (ubuntu, test-shard.sh 0 1 -race) has event push, matrix/web-dist/cockpit/shipped-manifest/ticket-commits gated to pull_request; YAML parses (pyyaml). Function on real history as push ranges: 7085dc47f^1..7085dc47f code=true, logbook-only b0c1d6359 code=false; before=000..0 exits 128 with no verdict (job red, not skipped). Writes only GITHUB_OUTPUT (parallel appends checked round 2). Not exercisable locally: the push run on GitHub itself"
hunt-verdict: "pass: bughunt done at ecb88f76f, 1 change not a blocker (force-push on master leaves before unreachable: changes red, master job skipped - loud, not silent; note for the person), 0 pre-existing"
review-verdict: "accept: both DoD items are met by ci.yaml and ci-code-changed.sh, TestCICodeChanged passes locally (go test ./scripts), no defect found beyond the loud force-push edge"
review-round: 3
review-check: |-
  1. go test ./scripts -run TestCICodeChanged -count=1 -v  -> every subtest PASS, incl. 'a logbook file with an umlaut' (code=false) and 'a ticket with its code' (code=true)
  2. bash scripts/ci-code-changed.sh b0c1d6359~1 b0c1d6359  -> prints code=false (logbook-only commit)
  3. bash scripts/ci-code-changed.sh master~5 master  -> prints code=true and a 'needs tests: <file>' line on stderr
  4. After the PR is opened on GitHub: the Checks tab lists changes, ticket-commits, shipped-manifest, web-dist, cockpit and 5 test jobs; 'master' shows as skipped
  5. After the merge: the Actions run for the push to master shows only 'changes' and 'master' (ubuntu) running - no macOS or Windows job
review-verdict-gaps: |-
  note: a force-push to master leaves github.event.before outside the checkout: changes fails red and master is skipped (loud, not silent) - fall back to HEAD~1 when git cat-file -e "$before" fails
  note: the master job runs go test only, not the go vet ./... and go build ./cmd/jaira the matrix runs - add both steps if a vet-only break on master matters
  note: actionlint was not run on ci.yaml (not installed here); the first PR run on GitHub is the real proof that needs/if resolve as intended
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
- **2026-11-22 19:51 · Alexander Sacharov** — critique round 3: read only ecb88f76f (the give-back). Separate 'master' job kept over reusing the matrix: a job-level if cannot see matrix.os, so the only way to keep macOS/Windows off a push is a fromJSON matrix switch or a second job; step-level ifs would still start those runners. Earlier critique note (ci-done dropped) stands.
- **2026-11-22 19:53 · Alexander Sacharov** — bughunt: change: .github/workflows/ci.yaml:55 (changes job, push range github.event.before...github.sha via scripts/ci-code-changed.sh:37) - class: a push to master whose 'before' commit is not in the checkout (a force-push that rewrote master; 'before' is then unreachable and fetch-depth: 0 does not bring it); example: git diff 'X...HEAD' with X absent fails with 'Invalid symmetric difference expression', the changes job exits 128 and the 'master' race job is skipped through needs - the rewritten master gets no test run, only a red changes job. Loud, not silent; not a blocker. ran: GITHUB_OUTPUT= bash scripts/ci-code-changed.sh 1111111111111111111111111111111111111111 HEAD (exit 128); blame: ecb88f76f (CQ1DCC)
- **2026-11-22 19:53 · Alexander Sacharov** — bughunt: done at ecb88f76f - 1 finding (1 change, 0 pre-existing) in 3 modules, follow-up hunt since d41394ac6 (diff of ci.yaml and ci_code_changed_test.go; ci-code-changed.sh unchanged); concurrency, error paths, input, permissions checked; files: .github/workflows/ci.yaml scripts/ci-code-changed.sh scripts/ci_code_changed_test.go
