---
id: 01M6PDF4CMXG25611W5A2DSM0T
title: "Ein Board ohne Git-Repo schiebt Tickets nach done, ohne Commits zu verlangen"
status: optimize
ready: true
creator: BeMuCa
goal: "Liegt das Board in keinem Git-Repository, verlangt die Lane done keine Commits; in einem Repository bleibt die Regel wie heute."
context: |-
  Was schief ist: in einem Ordner ohne Git laufen init, create, move und dod (ausprobiert 2026-11-02). Der Move nach done scheitert:
  'lane "done" requires the commits that carry this change, and git has none for P8QHNY yet'.
  Die Auswege der Meldung (--commits, jaira set commits=) sind ohne Git sinnlos. Ohne done kein Logbook.
  Gewollt (Berk, 2026-11-05): ein Solo-Nutzer ohne Repo kann jaira ganz benutzen. Mit Repo, auch rein lokal ohne Remote, bleibt die Commit-Pflicht.
  Wo: die Gate-Pruefung fuer Commits in core/gate/ (grep 'requires the commits'). core/gitrepo kennt ErrNoGit.
definition-of-done: "Ohne Git-Repo geht ein Ticket nach done und ins Logbook, ohne Commits und ohne --force"
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-11-05T17:52:51Z
updated-at: 2026-11-23T01:32:43Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: worker-2dsm0t
claimed-at: 2026-11-23T01:24:05Z
outcome-what: "A board that lives in no git repository moves tickets into done and the logbook without commits, provided git is installed; inside a repository, remote or not, done still requires them."
outcome-why: "Solo users without a repository could never finish a ticket: the done gate demanded commits that cannot exist. Round 2: the bug hunt showed a symlinked repo with git missing read as 'no repository'."
outcome-resolves: "2DSM0T: gate.Env.OutsideRepo backed by gitrepo.Repo.OutsideRepo - true only when git is on PATH, no .git above the symlink-resolved path and git sees no work tree; wired in CLI and TUI; logbook drops the hand-record hint there. Round 2 fixed the hunt's blocker (EvalSymlinks; no git -> keep the rule). Checked: concurrency n/a (read-only probe; two parallel moves both landed); error paths - any stat/git/symlink error keeps the rule; input - relative, symlinked, subdir paths; permissions - unreadable parent counts as 'maybe a repo'."
review-summary: none
review-gaps: "none: round 2 added an early return and one EvalSymlinks call; the trailing !Available() it replaced is gone, nothing orphaned, no duplicate symlink-resolve helper in core/gitrepo"
test-verdict: "fail: bughunt change blocker - OutsideRepo reads a symlinked repo path with git missing as 'no repository' and waives commits inside a real repo; own passes were green (vet/build rc=0, scratch boards behaved)"
hunt-verdict: "fail: 1 change blocker (core/gitrepo/git.go OutsideRepo, symlink + git missing)"
round: 2
back-from: testing
back-by: lane
back-count: 1
---

# Ein Board ohne Git-Repo schiebt Tickets nach done, ohne Commits zu verlangen

## Definition of Done

- [x] Ohne Git-Repo geht ein Ticket nach done und ins Logbook, ohne Commits und ohne --force
  proof: core/gate/gate.go OutsideRepo waiver in the requires-commits check; core/gitrepo/git.go OutsideRepo; manual: scratch board in plain dir, move to done exit 0 and logbook exit 0 without --force
- [x] In einem Repo ohne Remote verlangt done weiter Commits
  proof: gitrepo.OutsideRepo false when git missing, a .git exists above the real (symlink-resolved) path, or git sees a work tree; manual: git init without remote refused exit 3, also through a symlink, with and without git on PATH; TestOutsideRepo repo cases
- [x] Ein Test je Fall
  proof: TestCommitsWaivedOnlyOutsideARepository (core/gate/commits_test.go), TestOutsideRepo + TestOutsideRepoWithoutGitKeepsTheRule (core/gitrepo/outside_test.go); compiled by go vet, not run (worker rule: no go test)
- [x] Release-Notiz in core/release/notes/<ID>.md
  proof: core/release/notes/2DSM0T.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] gitrepo: OutsideRepo() - true only when no .git in Dir or any parent AND git does not report a work tree
- [x] gate.Env: lazy OutsideRepo func; requires-commits refusal skipped only when derive found nothing and OutsideRepo() is true
- [x] wire it in internal/cli/root.go loadEnv and internal/tui/model.go gateEnv
- [x] tests: gate (nil/false keeps refusal, true passes), gitrepo (plain dir, repo without remote, subdir of repo)
- [x] release note core/release/notes/2DSM0T.md
- [x] manual: scratch board without git reaches done and logbook; repo without remote still refused

## Progress
- **2026-11-23 01:25 · BeMuCa** — Reproduced 2026-11-23 on a scratch board in a plain dir (branch build): move to done refused with 'lane "done" requires the commits ...'. Decision (mine, settle-able): 'no repo' = no .git entry in the board root or any parent AND git does not answer true to rev-parse --is-inside-work-tree (or git is not installed). Why both: a git error that is not 'no repo' (dubious ownership, broken config) must never drop the commit rule inside a real repo; a missing git binary alone must not either when a .git is there. Asked lazily, only on the refusal path, so a repo board pays nothing and render paths never call it. Planning option not ticked, so pre-process is skipped; plan written as plan items anyway.
- **2026-11-23 01:27 · BeMuCa** — Sibling found and fixed in the same change: 'jaira logbook' printed 'Record them by hand with jaira set <id> commits=<sha>' on a git-less board - the same pointless advice. Now silent there (internal/cli/logbook.go). Ownership refusal in manual tests was unrelated: identity falls back to $USER, so the scratch tests set JAIRA_USER. Tests written but not run here (house rule); the main session runs go test.
- **2026-11-23 01:27 · BeMuCa** — critique: considered reusing gitref.Repo.Usable / ErrNoGitRepo (internal/cli/refs.go:158) as the 'no repo' answer - rejected: it reads any rev-parse failure as no repository, which would drop the commit rule on a real repo that git refuses (dubious ownership). The stricter gitrepo.OutsideRepo is a second notion of 'no repo' on purpose. The trailing IsRepo after the .git walk covers GIT_DIR setups; it runs only on the refusal path. Nothing else to change.
- **2026-11-23 01:28 · BeMuCa** — testing (worker): gofmt clean, go vet ./... rc=0, go build ./... rc=0; go test NOT run (house rule tonight - main session runs the race suite before the PR). Branch binary on scratch boards: plain dir -> done exit 0 + logbook exit 0, no hand-record hint; git init without remote -> refused exit 3 needs_commits, and with a fix(<id>) commit -> done + logbook 'Stamped 1 commit(s)'; failure path: repo board with git removed from PATH -> still refused exit 3 (the .git walk decides); plain dir without git on PATH -> done exit 0; two parallel moves to done on the plain-dir board both landed. Bug hunt running in background.
- **2026-11-23 01:30 · BeMuCa** — bughunt: change: blocker: core/gitrepo/git.go:57-71 (OutsideRepo) - class: the repository test is lexical (filepath.Abs keeps symlinks, the .git walk follows the path as written) and its only fallback, !Available() || !IsRepo(), answers true whenever git is missing, so a repository the walk cannot see reads as 'no repository' and the requires-commits gate (core/gate/gate.go:372) is waived inside a real repo - the opposite of the function's own doc (a missing binary must not drop the requirement); same pattern reached from internal/cli/root.go:365, internal/tui/model.go:2069 and the logbook hint at internal/cli/logbook.go:324; example: board at ~/link -> repo/sub (git init'ed repo), git not on PATH: 'jaira move <id> --to done' exit 0, ticket in done with commits: [] - the same repo via its real path refuses with needs_commits, exit 3; fix direction: resolve the path with filepath.EvalSymlinks before the walk (and keep refusing when git is missing and the answer is unsure); ran: scratch board, PATH without git, symlinked cwd (exit 0) vs real path (exit 3)
- **2026-11-23 01:30 · BeMuCa** — bughunt: done at 821aa6f78 - 1 finding (1 change, 0 pre-existing) in 4 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/gate/commits_test.go core/gate/gate.go core/gitrepo/git.go core/gitrepo/outside_test.go core/release/notes/2DSM0T.md internal/cli/logbook.go internal/cli/root.go internal/tui/model.go
- **2026-11-23 01:31 · BeMuCa** — Decision for the fix (mine): without git on PATH OutsideRepo answers false - the requirement stays. Reason: without git nothing can tell a repository reached through a symlink, GIT_DIR or GIT_WORK_TREE from a plain directory, and the rule says a missing binary must never drop the requirement. Cost: a solo user with no git installed at all still cannot reach done without --force; jaira's gitrepo package already assumes git on PATH. Also resolve symlinks before the .git walk so a git that fails for another reason (dubious ownership) behind a symlink does not read as 'no repository'.
- **2026-11-23 01:32 · BeMuCa** — critique round 2: the fix is the smallest that closes the class (path as written vs real path, and a missing git) - early return on !Available, EvalSymlinks before the walk. Nothing else; standing on the round-1 critique note for the gitref.Usable question.
