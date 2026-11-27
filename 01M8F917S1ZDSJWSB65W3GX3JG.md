---
id: 01M8F917S1ZDSJWSB65W3GX3JG
title: Worktree ohne .jaira findet das Board des Hauptcheckouts
status: testing
ready: true
creator: Sascha
goal: "jaira (und damit der cockpit) arbeitet in einem git-worktree auch dann, wenn dort kein .jaira liegt: es nimmt das Board des Hauptcheckouts."
context: |-
  Was falsch ist: Accept im cockpit schlägt in einem frischen Projekt fast immer fehl mit 'ticket: no .jaira directory found; run jaira init' (reason no_store).
  Warum: vor 'jaira share' steht .jaira nicht in git. Jeder worktree, den spawn.sh anlegt, hat darum kein .jaira. Dasselbe gilt für Branches, die vor share entstanden sind.
  Der cockpit ruft jaira im worktree des Tickets auf. ticket.Discover (core/ticket/store.go:268) sucht nur aufwärts nach .jaira/tickets und findet nichts.
  Heute hilft nur: im Hauptcheckout 'jaira move <id> --to done --force' von Hand.
  Schon vorhanden: core/project/worktree.go mainCheckout() findet den Hauptcheckout über die commondir-Datei.
  Idee: Discover fällt, wenn aufwärts nichts gefunden wird und dir in einem linked worktree liegt, auf mainCheckout(dir) zurück. Ein worktree MIT eigenem .jaira bleibt bei seinem eigenen.
  Offen für pre-process: ob Schreiben aus dem worktree ins Hauptcheckout-Board Sperren/Refs richtig trifft (Store.Root zeigt dann auf den Hauptcheckout).
definition-of-done: Im worktree ohne .jaira findet jaira list/show/move das Board des Hauptcheckouts (Test in core/ticket oder core/project)
tags:
  - cli
  - cockpit
blocked-by: []
related: []
commits: []
created-at: 2026-11-27T19:52:04Z
updated-at: 2026-11-27T20:56:15Z
assignee: Sascha
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-68058
claimed-at: 2026-11-27T20:33:41Z
outcome-what: "ticket.Discover falls back to the main checkout's board when a linked worktree has none; the git-file reading moved from core/project into a new leaf package core/gitwt shared by both"
outcome-why: "worktrees cut before jaira share have no .jaira, so every jaira call there - the cockpit's Accept included - failed with no_store"
outcome-resolves: jaira list/show/move and cockpit Accept work in a worktree without .jaira; a worktree with its own board keeps it
review-summary: "none"
round: 5
back-from: testing
back-by: lane
back-count: 3
review-gaps: "round 4: nothing to remove - the bound reuses gitwt.Top; the under() check is 1 line and project.under cannot be imported from ticket. Earlier: unexported Store.tree; removed TestDiscoverWorktreeWritesMainBoard. Left: gitwt has no own tests - callers' tests exercise every function"
test-verdict: "pass: go vet + go test -race ./... rc=0 at e1ae3cc55; DoD 1,2,4 verified by the named tests (each red with its fix disabled); binary on scratch repos: list/move/create from a worktree without .jaira use the main board, create follows the worktree's branch ticket, worktree of a repo without a board exits 1 and writes no .jaira, board above the repo is not borrowed (test); two parallel notes from worktree and main both land. DoD 3 (cockpit Accept, seen by a person) left for the person"
hunt-verdict: "fail: 1 change (edge) - inside-check uses the nearest .git, a submodule top in the main checkout refuses the main board from <wt>/sub (round 4)"
question: "testing would send 3GX3JG back a 4th time (jaira refuses it without the teamlead). Round 4 finding: the inside-check measures against gitwt.Top(main), the nearest .git - in a path that sits in an initialized submodule of the main checkout that is the submodule top, so jaira run in <wt>/sub of a worktree whose sub/ is not initialized exits 1 (no board); at <wt> it works. Fails toward ErrNoStore, as master. Options: 1) (Recommended) fix: bound by gitwt.MainCheckout(gitwt.Top(abs)) - the real main top, 1 line + 1 test case; you allow the send-back. 2) accept as known limit, go on to review with it as a note."
---

# Worktree ohne .jaira findet das Board des Hauptcheckouts

## Definition of Done

- [x] Im worktree ohne .jaira findet jaira list/show/move das Board des Hauptcheckouts (Test in core/ticket oder core/project)
  proof: core/ticket/store.go Discover fallback + WorkTree; internal/cli/worktreeboard_test.go TestWorktreeWithoutBoardUsesTheMainCheckouts (list/show/move), TestWorktreeWithoutBoardAsksGitAboutItself (branch, commits); core/ticket/discover_test.go TestDiscoverInLinkedWorktree
- [x] Ein worktree mit eigenem .jaira benutzt weiter sein eigenes (Test)
  proof: core/ticket/discover_test.go TestDiscoverInLinkedWorktree cases 'the worktree's own board wins' (+ from below)
- [ ] cockpit Accept auf ein Ticket, dessen worktree kein .jaira hat, geht durch (gesehen von einer Person)
- [x] Release-Note core/release/notes/<ID>.md
  proof: core/release/notes/3GX3JG.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-27 20:33 · Sascha** — parent session: uds:/run/user/1000/cc-socks/24078.sock
- **2026-11-27 20:34 · Sascha** — decisions (worker, none open for a person): 1) fallback lives in ticket.Discover, only after the upward walk finds nothing - a worktree with its own .jaira keeps it. 2) the git-file reading of project/worktree.go moves into a new leaf package core/gitwt (Top, MainCheckout, Linked) - project imports ticket, so ticket cannot call project; gitwt keeps one copy. 3) writes from the worktree land on the main checkout's board: locks and state are that board's, refs/jaira are shared by all worktrees, CommitsForTicket reads HEAD plus --branches, so branch commits are still found. Known limit: Branch() (create --follows from the branch name) reads the main checkout's branch when run from such a worktree. 4) DoD 3 is the person's check at signoff.
- **2026-11-27 20:39 · Sascha** — critique: core/project/worktree.go:50 Trees' second comment paragraph repeats gitwt.Linked's doc - drop it. Let stand: the fallback walks up from the main-checkout place past the repo top, same as a command run in the main checkout itself would - the worktree behaves as its main checkout; a bound to the repo top would make the two disagree.
- **2026-11-27 20:42 · Sascha** — bughunt: change: blocker: core/ticket/store.go:283 (Discover falls back to gitwt.MainCheckout) - class: once a linked worktree without .jaira opens the main checkout's board, every git question jaira asks through gitrepo.Repo{Dir: s.Root} is answered by the main checkout's HEAD/branch/index, not the worktree's: internal/cli/tickets.go:171 branchTicket (follows), internal/cli/root.go:359 DeriveCommits/CommitsElsewhere/HuntCheck (terminal gate, bughunt freshness), internal/cli/flow.go:675 show --for-lane diff, internal/cli/roles.go:193 pr-gate, tickets.go:865 LastHunt; example 1: main checkout on feat/dpgx0x, 'jaira create' in a worktree on feat/wzktgw printed 'Follows DPGX0X, from branch feat/dpgx0x' (wrong ticket linked; on master no follows at all); example 2: worktree commit 'fix(WZKTGW): work', 'jaira show WZKTGW --for-lane review --json' in the worktree gave commits: [] and a commits_source=worktree diff of the main checkout's untracked AGENTS.md/CLAUDE.md, so a review lane judges the wrong diff; ran: scratch repo + git worktree add + built binary from 807794849; before this change the worktree got ErrNoStore, so the path is newly reached
- **2026-11-27 20:42 · Sascha** — hunt report: 1 finding, 1 change (blocker), note on 3GX3JG | read whole: core/gitwt/gitwt.go core/project/worktree.go core/ticket/store.go (Discover/boardAbove) core/ticket/discover_test.go internal/cli/worktreeboard_test.go core/release/notes/3GX3JG.md, plus callers of Discover and of gitrepo.Repo{Dir: s.Root} | not checked: 'jaira pull' from such a worktree (needs a remote; it likely writes the file into the main checkout, same class)
- **2026-11-27 20:42 · Sascha** — bughunt: done at 807794849 - 1 finding (1 change, 0 pre-existing) in 6 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/gitwt/gitwt.go core/project/worktree.go core/release/notes/3GX3JG.md core/ticket/discover_test.go core/ticket/store.go internal/cli/worktreeboard_test.go
- **2026-11-27 20:45 · Sascha** — fix for bughunt blocker: ticket.Store gets Tree (set by Discover's fallback) and WorkTree(); every gitrepo.Repo that asks about HEAD/branch/diff/commits now uses WorkTree - root.go loadBoardEnv, flow.go diff, roles.go pr-gate, tickets.go branchTicket + LastHunt, tui gateEnv. Left on Root on purpose: questions about the board itself (share/update/tickets.go:83 is .jaira tracked or ignored), pr.Checker (asks by PR number), whoami remotes and refsync (refs and remotes are common to all worktrees), plugin/hook cmd.Dir (documented as the board root). 'jaira pull' in such a worktree writes into the main checkout's board - that is where this board's files live.
- **2026-11-27 20:48 · Sascha** — waiting: bughunt round 2 (background agent) since 21:48
- **2026-11-27 20:48 · Sascha** — bughunt: change: core/gitrepo/git.go:401 (anyCode) reached via internal/cli/root.go:359 and internal/tui/model.go:2108 - class: board-data paths are matched against 'git rev-parse --show-prefix' of r.Dir, and r.Dir is now Store.WorkTree(), which core/ticket/store.go:291 sets to the directory the command ran in (abs of -C/cwd), not the worktree's top; from a subdirectory of a board-less worktree the prefix is that subdirectory, so .jaira/tickets, .jaira/tags, .jaira/logbook, milestones all count as code - HuntCheck (hunt.go:102) then reports CodeChanged for a ticket whose commits carry only board data, and HuntCoverage (git.go:374,381,385) reads a hunt as stale when a later commit of the ticket touched only tags/milestones; example: worktree of a shared board whose .jaira/tickets/ is empty in git (git keeps no empty dir, so Discover falls back), 'cd core && jaira move <id> --to human' on a ticket with only ticket-file commits is refused for a missing hunt; before this change Dir was Store.Root (board top, prefix ''); fix: set tree to the worktree's top-level (git rev-parse --show-toplevel), not abs; edge, not blocker - fails toward demanding a hunt, needs subdir + board-less worktree + board-data-only commits; ran: scratch test in core/gitrepo, (&Repo{Dir: top}).anyCode('.jaira/tags')=false, (&Repo{Dir: top/sub}).anyCode(...)=true
- **2026-11-27 20:48 · Sascha** — hunt report: 1 finding, 1 change (edge, not blocker), as a note on 3GX3JG | read as a diff since 807794849: core/ticket/store.go internal/cli/worktreeboard_test.go; new to the change and read around their changed call (not whole - 6.7k lines of one-line edits): internal/cli/flow.go internal/cli/roles.go internal/cli/root.go internal/cli/tickets.go internal/tui/model.go, plus every gitrepo method those calls reach (git.go, derive.go, hunt.go, named.go) and every remaining gitrepo.Repo{Dir: s.Root} site | not checked end to end: the full move refusal in a board-less worktree with ticket-file-only commits (needs a shared board with an untracked-empty tickets dir); CommitsForTicket's first source now gets an absolute ticket path outside the worktree and git errors silently, judged unreachable since such a worktree has no ticket files in its history
- **2026-11-27 20:48 · Sascha** — bughunt: done at 17c54243e - 1 finding (1 change, 0 pre-existing) in 8 modules, follow-up hunt since 807794849; concurrency, error paths, input, permissions checked; files: core/gitwt/gitwt.go core/project/worktree.go core/release/notes/3GX3JG.md core/ticket/discover_test.go core/ticket/store.go internal/cli/flow.go internal/cli/roles.go internal/cli/root.go internal/cli/tickets.go internal/cli/worktreeboard_test.go internal/tui/model.go
- **2026-11-27 20:49 · Sascha** — fix for hunt r2: tree is now Root mirrored into the worktree (abs + rel(main, root)), not the directory the command ran in - git's --show-prefix then matches the board's place, as with Root. TestDiscoverInLinkedWorktree checks WorkTree per case; red with the old value (2 cases).
- **2026-11-27 20:51 · Sascha** — bughunt: change: blocker: core/ticket/store.go:294-298 (Discover) - class: board found above the main checkout's top (boardAbove(main) walks past the repo, so Rel(main, root) climbs out of it) makes tree = abs+'../..' point outside the linked worktree, so every WorkTree() caller (cli/tickets.go:171 Branch, :865 LastHunt, cli/root.go:359, cli/flow.go:675, cli/roles.go:193, tui/model.go:2108) asks git in a directory that is not the worktree - no repo at all, or a different enclosing repo; example: board at parent/.jaira (not in git), main checkout parent/repo, worktree elsewhere/wt2 on branch feat/6gwvx6 - 'jaira create' no longer sets follows 6GWVX6, which the 17c54243e build does; fix direction: map root via the worktree top (gitwt.Top(abs)) and fall back to that top when root lies outside the main checkout; ran: built HEAD and 17c54243e, jaira create in a scratch layout (HEAD: no Follows line; old: 'Follows 6GWVX6, from branch feat/6gwvx6')
- **2026-11-27 20:51 · Sascha** — hunt report: 1 finding, 1 change (blocker), note on 3GX3JG | read as a diff since 17c54243e: core/ticket/store.go, core/ticket/discover_test.go (no other file changed, none new) | not checked: whether the enclosing-repo variant (board root inside another git repo) gives a wrong branch rather than none - same path, not run
- **2026-11-27 20:51 · Sascha** — bughunt: done at eb1c70358 - 1 finding (1 change, 0 pre-existing) in 2 modules, follow-up hunt since 17c54243e; concurrency, error paths, input, permissions checked; files: core/gitwt/gitwt.go core/project/worktree.go core/release/notes/3GX3JG.md core/ticket/discover_test.go core/ticket/store.go internal/cli/flow.go internal/cli/roles.go internal/cli/root.go internal/cli/tickets.go internal/cli/worktreeboard_test.go internal/tui/model.go
- **2026-11-27 20:51 · Sascha** — question (3rd testing send-back, to teamlead): testing would send 3GX3JG back a 3rd time. Round 3 finding: a board ABOVE the main checkout top (parent/.jaira, repo in parent/repo) makes tree point outside the worktree, so create loses follows. My view: the loop is not deepening the DoD, it is the fallback reaching too far - edges grow, pick the smaller design. Options: 1) (Recommended) fall back only to a board inside the main checkout (root under its top); a board above the repo stays ErrNoStore in a worktree, as before this ticket - ~3 lines + 1 test case. 2) keep the wider fallback, clamp tree to the worktree top when root lies outside the main checkout - ~5 lines, keeps an odd setup working. 3) accept as known limit, move on to review. Send back for 1?
- **2026-11-27 20:52 · Sascha** — waiting: teamlead decision on the 3rd send-back (uds:/run/user/1000/cc-socks/24078.sock) since 21:52
- **2026-11-27 20:52 · Sascha** — decision (teamlead): Rückfall nur auf ein Board im Hauptcheckout selbst; ein Board oberhalb des Repos bleibt im worktree ErrNoStore wie bisher. Kleinere Lösung statt tree-Klammer – die Randfälle wuchsen.
- **2026-11-27 20:54 · Sascha** — fix per teamlead decision (option 1): Discover falls back only to a board inside the main checkout (Rel to gitwt.Top(main) must not climb); a board above the repo stays ErrNoStore in a worktree. Test TestDiscoverIgnoresBoardAboveMainCheckout - red without the bound.
- **2026-11-27 20:55 · Sascha** — waiting: bughunt round 4 (background agent) since 21:55
- **2026-11-27 20:55 · Sascha** — bughunt: change: core/ticket/store.go:294 - class: the 'board inside the main checkout' check measures against gitwt.Top(main), the nearest .git above the mapped path, not the main checkout's top; any path that in the main checkout sits inside a nested repository (initialized submodule, nested clone) gets that repository's top, so the main checkout's own board counts as 'above' and is refused; example: main checkout with submodule sub/ initialized, worktree where sub/ is not initialized, 'jaira list' run in <wt>/sub exits 1 'no .jaira directory found' while the same command at <wt> lists the board; filepath.Dir(common) inside gitwt.MainCheckout is the real main top; ran: scratch repo + submodule + worktree, HEAD binary exit=1, binary at eb1c70358 exit=0 (lists ticket); blame: e1ae3cc55 (names the ticket); not a blocker: master had no fallback at all, so nothing that worked before the ticket breaks
- **2026-11-27 20:55 · Sascha** — hunt report: 1 finding, 1 change (edge, not blocker), as a note on 3GX3JG | read as a diff since eb1c70358: core/ticket/store.go (Discover), core/ticket/discover_test.go; no new files | could not check Windows path case behaviour of the new filepath.Rel test on a real Windows box
- **2026-11-27 20:55 · Sascha** — bughunt: done at e1ae3cc55 - 1 finding (1 change, 0 pre-existing) in 2 modules, follow-up hunt since eb1c70358; concurrency, error paths, input, permissions checked; files: core/gitwt/gitwt.go core/project/worktree.go core/release/notes/3GX3JG.md core/ticket/discover_test.go core/ticket/store.go internal/cli/flow.go internal/cli/roles.go internal/cli/root.go internal/cli/tickets.go internal/cli/worktreeboard_test.go internal/tui/model.go
