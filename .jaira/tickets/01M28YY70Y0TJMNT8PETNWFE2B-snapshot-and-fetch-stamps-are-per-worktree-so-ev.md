---
id: 01M28YY70Y0TJMNT8PETNWFE2B
title: "Snapshot and fetch stamps are per worktree, so every new worktree snapshots at once"
status: signoff
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Key the snapshot and fetch schedule stamps by repository instead of by working tree, so the 72h/interval clock is shared across all worktrees of one clone"
context: |-
  The jaira/board snapshot branch gets a commit every 20-30 minutes instead of every 72 hours.

  Cause: the stamp file that records when the background job last ran is stored per working tree. core/ticket/store.go:487 builds the state dir as ~/.jaira/state/<basename>-<sha256(Root)[:4]>, keyed on the working tree path. core/snapshot/schedule.go:23 puts snapshot.json there; core/refsync/schedule.go:34 puts fetch.json there.

  A fresh git worktree is a new path, so it gets a new state dir with no stamp. bgrun.Due (core/bgrun/bgrun.go:106) treats a missing stamp as 'never ran' and fires immediately. Every worktree-agent-* worktree therefore takes a full snapshot on its first jaira command.

  Measured: ~/.jaira/state holds 2694 directories; the agent-* ones each carry their own recent ran_at (17:22, 17:23, 17:24, 18:05, 18:06, 18:34, 18:51, 18:52, 19:17) - exactly the commit spacing on jaira/board.

  The tree-identity guard in core/snapshot/snapshot.go:169 does not help, because tickets really do change between runs, so each run writes a real commit and pushes.

  Known: sessions and locks under the same state dir are correctly per working tree - only the two background-job stamps are repo-wide concerns. git rev-parse --git-common-dir is identical for every worktree of one clone and is the obvious key.
definition-of-done: snapshot.json and fetch.json are keyed by repository (shared by all worktrees of one clone); a second worktree of the same repo does not trigger a snapshot or fetch when the first one ran within the interval; sessions and locks stay per working tree; tests cover the shared-clock behaviour across two worktrees
tags:
  - concurrency
blocked-by: []
commits: []
created-at: 2026-09-11T19:26:01Z
updated-at: 2026-09-11T19:46:24Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1116663
claimed-at: 2026-09-11T19:26:19Z
body: |-
  ## Plan

  - [ ] add Repo.CommonDir to core/gitrepo: git rev-parse --path-format=absolute --git-common-dir, identical from every worktree of one clone
  - [ ] add Store.RepoStateDir in core/ticket: same home directory, keyed on the common dir instead of Root, falling back to the per-tree dir when git cannot answer
  - [ ] point maybeSnapshot (internal/cli/snapshot.go) and maybeFetch (internal/cli/refs.go) at RepoStateDir; leave outbox, refs-seen, sessions and locks on the per-tree dir
  - [ ] test: two worktrees of one fixture repo resolve to the same RepoStateDir, and a per-tree dir for two separate clones stays distinct
  - [ ] test: after a stamped run in worktree A, snapshot.Due is false in worktree B
  - [ ] note in core/release/NOTES.md that the snapshot and fetch clocks are now shared per clone
  - [ ] go test ./... -race
outcome-what: "Keyed the snapshot and fetch stamps by repository: Store.RepoStateDir (core/ticket/store.go) resolves git rev-parse --git-common-dir and keys ~/.jaira/repo/<key> by it, while StateDir stays per working tree; internal/cli/snapshot.go and internal/cli/refs.go spawn against it."
outcome-why: "A git worktree is a new path, so it got a new state dir with no stamp, and bgrun.Due read a missing stamp as never-ran and fired at once — every worktree-agent-* took a full snapshot on its first command, turning a 72h backup into a commit every 20 minutes."
outcome-resolves: "jaira/board gets one snapshot per clone per interval; sessions, locks, outbox and refs-seen stay per checkout; four tests in core/ticket/statedir_test.go cover the shared clock, the per-tree split, two clones and the no-git fallback; go test ./... -race green."
review-summary: "RepoStateDir keys the two background stamps by git --git-common-dir, which is one value per clone and the answer git itself gives; StateDir is untouched, so sessions, locks, outbox and refs-seen stay per checkout. core/bgrun and core/snapshot are unchanged - only the directory handed to them moves - so the spawn, the recursion guard and the tree-identity guard did not have to be re-reasoned. The shell-out sits in core/ticket rather than core/gitrepo because gitrepo/derive.go already imports core/ticket."
review-gaps: "Adds one git rev-parse per command, on the path that already decides whether to spawn a background job - once per process, not per ticket. Upgrading fires one snapshot per clone, because the old per-tree stamps are not migrated: correct once, then quiet. The 2694 stale ~/.jaira/state directories are not cleaned up; separate ticket."
review-verdict: Ready. The change matches the diagnosed cause and leaves the state that is rightly per checkout alone.
review-check: "1. Run: go test ./core/ticket -run StateDir -v -- four tests pass. 2. Add a second checkout of this repo at /tmp/wt-check (worktree add). 3. In /tmp/wt-check run: jaira list. 4. Look at the log of branch jaira/board -- no new 'board: N ticket(s)' commit appeared from that checkout. 5. Run: ls ~/.jaira/repo -- one directory for this clone, holding snapshot.json and fetch.json."
---

# Snapshot and fetch stamps are per worktree, so every new worktree snapshots at once

## Definition of Done

- [x] snapshot.json and fetch.json are keyed by repository (shared by all worktrees of one clone); a second worktree of the same repo does not trigger a snapshot or fetch when the first one ran within the interval; sessions and locks stay per working tree; tests cover the shared-clock behaviour across two worktrees

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-11 19:44 · Alexander Sacharov** — core/ticket cannot import core/gitrepo: gitrepo/derive.go already imports core/ticket, so the obvious Repo.CommonDir helper is an import cycle. RepoStateDir therefore shells out to git itself in store.go (commonGitDir), which is ten lines and no new package edge. Rejected alternatives: reading .git by hand (a worktree points at the main checkout through a gitdir: file, a submodule points elsewhere again — a second, divergent resolver) and hanging the helper off core/gitrepo with a Store argument (inverts who owns the state-dir concept).

Only the two background stamps moved. outbox, refs-seen, sessions and locks stay per working tree on purpose: those are about in-flight work in this checkout, not about the board.
