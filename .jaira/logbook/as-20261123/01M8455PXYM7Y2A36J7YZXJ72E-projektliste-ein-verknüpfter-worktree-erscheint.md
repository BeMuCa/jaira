---
id: 01M8455PXYM7Y2A36J7YZXJ72E
title: "Projektliste: ein verknüpfter Worktree erscheint nicht als eigenes Projekt, sein Haupt-Checkout steht dafür"
status: done
ready: true
creator: Alexander Sacharov
goal: "Die Projektliste des Launchers (Projects:) zeigt jedes Repository einmal, unter seinem Haupt-Checkout, statt jeden Worktree und jede Wegwerf-Kopie, in der jaira einmal lief."
context: |-
  Was heute falsch ist: jeder Start von jaira trägt seinen Ordner in ~/.jaira/projects.json ein (core/project/project.go, Remember()).
  Folge: auf der Maschine des Maintainers stehen dort 10 Einträge für eigentlich 2 Repositories - 6 davon sind verknüpfte git-Worktrees unter ~/projects/.worktrees/<repo>-<slug>, 3 sind Wegwerf-Repos in Scratch-Ordnern von Agenten (/tmp/claude-*/.../scratchpad/...) und Abnahme-Sandboxen (~/.claude/acceptance/.../sandbox).
  Der Launcher (internal/tui/home.go, NewHome/boardRoots) listet alle, die Projektauswahl (internal/tui/view.go renderProjects) ebenso.
  Wunsch des Maintainers (23.11.2026): automatisch erkennen - ist jaira im Haupt-Checkout offen und in einem Worktree desselben Repos auch, soll der Worktree nicht als eigenes Board erscheinen.
  Bekannt: core/board/gitignore.go:332 liest 'git worktree list' schon; 'git rev-parse --git-common-dir' sagt für einen verknüpften Worktree, zu welchem Haupt-Repo er gehört.
  Ausgeschlossen: eine Option zum Abschalten - es soll ohne Einstellung stimmen.
definition-of-done: "Ein Start in einem verknüpften Worktree trägt nicht den Worktree, sondern den Haupt-Checkout des Repos in projects.json ein (oder gar nichts, wenn der schon drinsteht)"
tags:
  - tui
blocked-by: []
related: []
commits:
  - 4f8b111f64eee893b44949f6edb626c96e836587
  - 956cfb764ac0ef196925195e8bde1a960c7496c4
  - 8a7c9334ba69288c8c3b3c2d53d817ac1b4d01ef
  - 7fdcf87393b93ede153171f25e4a6ac79acd7ca8
  - e638f3e3a955c7aa54796a21a14b984407ba33e4
created-at: 2026-11-23T12:12:57Z
updated-at: 2026-11-23T12:44:23Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-27553
claimed-at: 2026-11-23T12:13:22Z
outcome-what: "project.Board(root) maps a board in a linked git worktree to the same place in its main checkout (read off .git file + commondir, no git process); Remember records that root and skips the temp dir; Load folds old worktree entries and temp entries; the launcher's boardRoots and 'jaira projects add' use the same mapping"
outcome-why: "~/.jaira/projects.json collected every worktree and agent scratch repo, so the launcher listed one repository many times"
outcome-resolves: "DoD 1-6; checked: concurrency (projects.json unlocked, pre-existing, not widened), error paths (unreadable .git/commondir -> entry kept as is), unusual input (submodule, bare-repo worktree, board in subdir, main without board, relative gitdir), permissions (only reads files git wrote in the repo)"
review-summary: "Full review of 4f8b111f6..7fdcf8739. The project list (projects.json, launcher, board switcher, 'jaira projects') now holds one line per repository: project.Board maps a board in a linked git worktree to the same place in its main checkout (read off .git file + commondir, no git process) when that checkout has a board; Remember records that root, Load folds old worktree entries and duplicates into it (latest open time wins) and drops boards under os.TempDir unless the registry itself lives in temp. So the folded line loses nothing, the TUI compares against Model.projectRoot (Board of the open store) for the current marker, no-op switch and drop guard, and counts agents/liveness across project.Trees (main + every linked worktree with a board, read off .git/worktrees/*/gitdir). Remember now runs its read-modify-write under ticket.LockFile and writes via ticket.WriteAtomic, returns an error, and 'projects add' exits non-zero on it and reports what was actually recorded (added: [] not null)."
review-gaps: |-
  note: limit: inside a linked worktree the switcher treats the main checkout's line as the open board - picking it is a no-op and the drop dialog refuses it as current; reaching the main checkout itself needs 'jaira board' there (internal/tui/pipeline.go:166, model.go:1617)
  note: limit: LoggedPerDay on the folded line counts only the main checkout's logbook; a worktree's logbook files arrive through the merge (internal/tui/home.go:242)
  note: limit: refreshLiveBoards now reads .git/worktrees and opens a store per linked worktree of every listed repository on each reload - plain file reads, 11 worktrees here, no git process (internal/tui/model.go:528)
  note: elsewhere: core/project/project.go:56 - a hand-broken projects.json is still read as empty and the next Remember overwrites it with one entry; there before this change
  note: elsewhere: the uncommitted AGENTS.md, CLAUDE.md and .jaira/lanes/*.md in this worktree are a board regeneration to jaira 0.6.4, not this ticket - keep them out of its commits
test-verdict: "pass: go build/vet ./... green; go test ./... green (RC=0), -race on core/project, internal/tui, internal/cli RC=0; DoD 1-7 verified in the tree; binary: two parallel 'projects add' 30 runs -> 0 torn or lost (before: 6/30), list shows demo-video + jaira main checkouts for two worktrees; unwritable registry dir -> 'could not record ... permission denied' rc=1; temp board --json -> added [] count 0"
hunt-verdict: "pass: bughunt done at 7fdcf8739, 0 change, 0 pre-existing (follow-up hunt since 956cfb764)"
round: 2
back-from: testing
back-by: lane
back-count: 1
review-verdict: "accept: all 7 DoD items hold in the code and tests (go test ./core/project ./internal/tui ./internal/cli green); the bughunt blocker (folded worktree losing agents/current marker) is closed across every comparison site; no blocker found along concurrency, error paths, input or permissions"
review-round: 2
review-verdict-gaps: |-
  note: limit: inside a linked worktree the switcher treats the main checkout's line as the open board - picking it is a no-op and the drop dialog refuses it as current; reaching the main checkout itself needs 'jaira board' there (internal/tui/pipeline.go:166, model.go:1617)
  note: limit: LoggedPerDay on the folded line counts only the main checkout's logbook; a worktree's logbook files arrive through the merge (internal/tui/home.go:242)
  note: limit: refreshLiveBoards now reads .git/worktrees and opens a store per linked worktree of every listed repository on each reload - plain file reads, 11 worktrees here, no git process (internal/tui/model.go:528)
  note: elsewhere: core/project/project.go:56 - a hand-broken projects.json is still read as empty and the next Remember overwrites it with one entry; there before this change
  note: elsewhere: the uncommitted AGENTS.md, CLAUDE.md and .jaira/lanes/*.md in this worktree are a board regeneration to jaira 0.6.4, not this ticket - keep them out of its commits
review-check: "1. cd /home/alex/projects/.worktrees/jaira-zxj72e && go build -o /tmp/jz ./cmd/jaira  2. export JAIRA_HOME=$(mktemp -d -p ~)  3. /tmp/jz projects add .  - prints 'Registered /home/alex/projects/jaira' (the main checkout, not the worktree)  4. /tmp/jz projects --json  - exactly one entry, root /home/alex/projects/jaira  5. /tmp/jz  (launcher) - one 'jaira' line; with an agent checkpointed in any worktree it shows the agent count  6. inside the board press p - the jaira tab carries the ▸ current marker although the worktree is open  7. for i in $(seq 10); do /tmp/jz projects add . & done; wait; python3 -m json.tool $JAIRA_HOME/projects.json - parses, one entry  8. rm -rf $JAIRA_HOME"
question: ""
---

# Projektliste: ein verknüpfter Worktree erscheint nicht als eigenes Projekt, sein Haupt-Checkout steht dafür

## Definition of Done

- [x] Ein Start in einem verknüpften Worktree trägt nicht den Worktree, sondern den Haupt-Checkout des Repos in projects.json ein (oder gar nichts, wenn der schon drinsteht)
  proof: core/project/project.go:97 Remember -> Board -> worktree.go:25 mainCheckout; TestRememberRecordsMainCheckout
- [x] Schon vorhandene Worktree-Einträge in projects.json zeigt der Launcher nicht mehr als eigene Zeile, wenn ihr Haupt-Checkout bekannt ist; ein Worktree, dessen Ordner nicht mehr existiert, verschwindet ebenfalls
  proof: core/project/project.go:66 Load folds via mainCheckout, drops missing .jaira; internal/tui/home.go:128; TestLoadFoldsOldWorktreeEntries
- [x] Repos unter dem Temp-Verzeichnis (os.TempDir) werden nicht eingetragen
  proof: core/project/worktree.go:93 throwaway, used in Remember and Load; TestRememberSkipsTemp
- [x] Ein Haupt-Checkout bleibt eine Zeile, auch wenn gerade mehrere seiner Worktrees jaira offen haben
  proof: core/project/project.go:66 dedupe by main-checkout root; TestRememberRecordsMainCheckout (two worktrees + main = one line)
- [x] Tests für: Worktree -> Haupt-Checkout, Temp-Repo ignoriert, alter Worktree-Eintrag ausgeblendet
  proof: core/project/worktree_test.go: TestBoardIsMainCheckout, TestRememberSkipsTemp, TestLoadFoldsOldWorktreeEntries
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/ZXJ72E.md
- [x] Parallele Starts zerreißen projects.json nicht mehr: Remember() schreibt atomar (Temp-Datei + rename) unter einer Sperre, mit Test für zwei gleichzeitige Schreiber
  proof: core/project/project.go:109 ticket.LockFile + :134 ticket.WriteAtomic; TestRememberParallel (fails without both)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 12:13 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/6958.sock
- **2026-11-23 12:13 · Alexander Sacharov** — decisions (worker, none open for a person): 1) worktree->main detection reads files only (.git file 'gitdir:' + its 'commondir'), no git exec - Remember runs on every start and startup must stay instant. 2) the main checkout replaces a worktree only when it has its own .jaira dir; a worktree whose main checkout has no board (private board only in the worktree) stays as it is, else the board would vanish from the list. 3) Load() does the mapping/filtering, so every reader (launcher, switcher, 'jaira projects') sees the same list; Remember writes the cleaned list back. 4) temp = paths under os.TempDir() after canonical(); ~/.claude/acceptance sandboxes are not covered (DoD names only os.TempDir).
- **2026-11-23 12:18 · Alexander Sacharov** — found out: the shipped tests register t.TempDir() boards (internal/tui/projects_test.go, qback_test.go), so a plain os.TempDir filter would empty their registries; throwaway() stands aside when the registry (JAIRA_HOME) itself lives in temp - also true for any sandbox that points JAIRA_HOME into temp. Launcher extras discovered below cwd are mapped too (boardRoots), so 'jaira' started inside a worktree lists the main checkout; 'jaira board' still opens the worktree's own board. Not touched: projects.json has no lock, two starts at once can lose an entry - pre-existing.
- **2026-11-23 12:19 · Alexander Sacharov** — critique: none. Left standing: worktree.go reads .git/commondir itself although core/ticket/store.go:commonGitDir argues for asking git - Load runs Board once per registry entry on every start (up to 20 git processes, slow on Windows/WSL); the read is narrow and submodule/bare cases are tested. throwaway()'s registry-in-temp exception is what keeps temp-board tests and sandboxes with JAIRA_HOME in temp working.
- **2026-11-23 12:21 · Alexander Sacharov** — waiting: bughunt (background subagent) since 13:21
- **2026-11-23 12:24 · Alexander Sacharov** — testing: note: pre-existing: two jaira processes writing ~/.jaira/projects.json at once leave invalid JSON (torn os.WriteFile without lock/rename), and Load then returns nothing - the whole list is gone until the next start rewrites it. Measured with two parallel 'jaira projects add': released 0.6.4 8/30 runs, this branch 6/30. Not this change (core/project/project.go write path unchanged); fix would be write-to-temp + rename. Also pre-existing: 'projects add' prints 'Registered' when the write fails (registry dir chmod 000 -> rc=0).
- **2026-11-23 12:24 · Alexander Sacharov** — bughunt: change: blocker: internal/tui/home.go:128 (with core/project/project.go:66,95) - class: state that is kept per working tree is now read off the main checkout once a linked worktree is folded into it, so whatever lives in the worktree vanishes from every list; places: launcher describe() counts Agents via session.Load on the main root (internal/tui/home.go:243) and LoggedPerDay from main's disk logbook (home.go:242); board switcher refreshLiveBoards compares p.Root == m.store.Root and reads boardHasLiveSession(main) (internal/tui/model.go:526,535), and the same equality decides the ▸/(current) marker and the no-op switch (internal/tui/pipeline.go:135,166, internal/tui/view.go:1731, model.go:1603); session state is keyed by working-tree root (core/ticket/store.go:895 stateKey) - example: agent checkpoints in worktree wt, launcher/switcher list only main and show 0 agents / no live marker for the repo, and inside wt no tab is marked current; before this change wt had its own row with its live agent; ran: checkpoint --focus in wt -> 'jaira sessions' in wt shows s1 live, in main 'No sessions', 'jaira projects' lists only main
- **2026-11-23 12:24 · Alexander Sacharov** — bughunt: change: core/project/project.go:93,116 + internal/cli/home.go:98,110 - class: Remember returns "" for three different reasons (no registry path, MkdirAll failure, temp-dir board) and 'projects add' reports all of them as the temp-dir case with exit 0; also 'added' is a nil slice so --json prints "added": null instead of [] whenever nothing was recorded - example: JAIRA_HOME=/proc/nope jaira projects add <board> prints 'Nothing registered: boards in the temporary directory are not recorded.' exit 0, and --json gives {"added": null, "count": 0} exit 0; ran: both commands, echo $? = 0
- **2026-11-23 12:24 · Alexander Sacharov** — bughunt: pre-existing: core/project/project.go:99-120 - class: projects.json is read-modify-written without a lock and written in place (os.WriteFile, not temp+rename), and Load returns nil on any unmarshal error, so parallel writers interleave the file and the next Remember overwrites the whole registry with one entry; write errors are also dropped and reported as 'Registered' (line 119) - example: 10 parallel 'jaira projects add bN' leave 0 or 2 of 10 entries, file fails json.load with 'Extra data'; chmod 444 projects.json -> 'Registered <root>' exit 0, file unchanged; triggers: parallel jaira init / jaira board / projects add; blame: bf665b5419 (2026-08-11, before the ticket); ran: 3 rounds of 10 parallel adds (0, 2, 2 entries left)
- **2026-11-23 12:24 · Alexander Sacharov** — hunt report: 3 findings (2 change - 1 blocker: worktree-folded rows lose per-worktree agents/liveness/current marker; 1 edge: misleading 'temp directory' message and added:null - 1 pre-existing: registry write race and dropped write errors), all as notes on ZXJ72E | read whole: core/project (project.go worktree.go discover.go + tests), internal/cli/home.go, internal/tui/home.go, plus callers in internal/tui/model.go pipeline.go view.go browse.go dropboard.go; read as diff: AGENTS.md CLAUDE.md .jaira/lanes/*.md (uncommitted regeneration to jaira 0.6.4, not this ticket's code - should not ride in its commit) and core/release/notes/ZXJ72E.md | not checked live: TUI rendering of the launcher/switcher (agent count and marker confirmed by reading describe/refreshLiveBoards plus 'jaira sessions' per tree); Windows-written gitdir paths (C:/...) under WSL fall back to a separate row - safe, not a defect
- **2026-11-23 12:24 · Alexander Sacharov** — bughunt: done at 956cfb764 - 3 findings (2 change, 1 pre-existing) in 3 modules, first hunt; concurrency, error paths, input, permissions checked; files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/project/project.go core/project/worktree.go core/project/worktree_test.go core/release/notes/ZXJ72E.md internal/cli/home.go internal/tui/home.go
- **2026-11-23 12:24 · Alexander Sacharov** — decision (teamlead 2026-11-23): fix the torn projects.json writes in this ticket, not as a follow-up
- **2026-11-23 12:29 · Alexander Sacharov** — round 2: bughunt blocker fixed as a class - 'is this line the open board' now compares with Model.projectRoot (project.Board of the store root) at every site (model.go refreshLiveBoards/drop, pipeline.go tabs+switch, view.go current); per-tree state (sessions/agents, liveness) is read across project.Trees(root) = main + every linked worktree with a board, read off .git/worktrees/*/gitdir. Not folded: LoggedPerDay stays the main checkout's logbook - worktree logbook files reach main through the merge. Remember now returns an error (no home, mkdir, lock, write); 'projects add' exits 1 on it, and --json added is [] not null. A hand-broken projects.json is still read as empty and overwritten - pre-existing, not the race.
- **2026-11-23 12:29 · Alexander Sacharov** — critique round 2: none. Left standing: describe() and boardHasLiveSession both open a store per tree to read sessions - two callers with different needs (count vs any, panic guard), not worth a shared helper.
- **2026-11-23 12:29 · Alexander Sacharov** — waiting: bughunt round 2 (background subagent) since 13:29
- **2026-11-23 12:31 · Alexander Sacharov** — hunt report: 0 findings, nothing written | read as a diff since 956cfb764: core/project/project.go, core/project/worktree.go, core/project/worktree_test.go, internal/cli/home.go, internal/tui/home.go; new to this round: internal/tui/model.go, pipeline.go, view.go (2600/300/2000 lines; read the hunks plus their whole functions: refreshLiveBoards, boardHasLiveSession, switchBoard, New, the modeProjects key handler, projectTabs, switchToProject, renderProjects, newDropBoard and its current guard), internal/tui/worktree_test.go whole | not checked: the uncommitted .jaira/lanes/*.md, AGENTS.md and CLAUDE.md, which are prose with no code path; Windows behaviour of the lock and WriteAtomic. Ran: go test -race -count=3 ./core/project (ok); 10 parallel 'jaira projects add' against a scratch JAIRA_HOME (all 10 recorded, file parses); 'projects add' with an unwritable JAIRA_HOME (exit 1)
- **2026-11-23 12:31 · Alexander Sacharov** — bughunt: done at 7fdcf8739 - 0 findings (0 change, 0 pre-existing) in 8 modules, follow-up hunt since 956cfb764; concurrency, error paths, input, permissions checked; files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/project/project.go core/project/worktree.go core/project/worktree_test.go core/release/notes/ZXJ72E.md internal/cli/home.go internal/tui/home.go internal/tui/model.go internal/tui/pipeline.go internal/tui/view.go internal/tui/worktree_test.go
- **2026-11-23 12:32 · Alexander Sacharov** — waiting: review lane (fresh agent via run-lane.sh) since 13:32
- **2026-11-23 12:35 · Alexander Sacharov** — review: done at 7fdcf8739
- **2026-11-23 12:35 · Alexander Sacharov** — review: move to signoff refused - the bughunt gate sees uncommitted changes in AGENTS.md, CLAUDE.md, .jaira/lanes/{critique,in-progress,review,testing}.md (a board regeneration to jaira 0.6.4, not this ticket). Review lane changes no code and will not commit them. Options:
1. (recommended) commit them as their own chore on this branch, e.g. chore(ZXJ72E): board regenerated by jaira update, then move review -> signoff - prose only, no new hunt needed
2. revert them here (git checkout -- AGENTS.md CLAUDE.md .jaira/lanes) and redo jaira update on master - keeps this PR to the ticket, but loses the regeneration until then
3. a person moves with --force - leaves the files dangling in the worktree
- **2026-11-23 12:36 · Alexander Sacharov** — question to teamlead: review accepted, but the move to signoff is refused by the bughunt gate over uncommitted AGENTS.md, CLAUDE.md, .jaira/lanes/{critique,in-progress,review,testing}.md - a board regeneration to jaira 0.6.4 that appeared in this worktree, not written by this ticket. 1) (Recommended) revert them here (git checkout -- those 6 files): a regeneration belongs in its own 'jaira update' on master, it is reproducible; 2) commit them as chore(ZXJ72E) on this branch: the PR carries unrelated board text; 3) move with --force, leaving them uncommitted: gate bypassed, worktree stays dirty.
- **2026-11-23 12:36 · Alexander Sacharov** — decision (teamlead 2026-11-23): option 1 - reverted the uncommitted 0.6.4 board regeneration (.jaira/lanes, AGENTS.md, CLAUDE.md); master already has it (#224)
- **2026-11-23 12:36 · Alexander Sacharov** — report: project list folds linked worktrees into their main checkout (project.Board/Trees, file reads only), skips temp boards, keeps worktree agents/live/current marker, registry written under lock + atomically, 'projects add' fails on write errors | signoff, context 177k of 1000k, cache warm 60 min after last answer | open: 5 review notes (limits/elsewhere) for the person; review-check is a hand test in the launcher; branch not pushed
- **2026-11-23 12:41 · Alexander Sacharov** — Accepted from signoff: angenommen:
- note: limit: inside a linked worktree the switcher treats the main checkout's line as the open board - picking it is a no-op and the drop dialog refuses it as current; reaching the main checkout itself needs 'jaira board' there (internal/tui/pipeline.go:166, model.go:1617)
- note: limit: LoggedPerDay on the folded line counts only the main checkout's logbook; a worktree's logbook files arrive through the merge (internal/tui/home.go:242)
- note: limit: refreshLiveBoards now reads .git/worktrees and opens a store per linked worktree of every listed repository on each reload - plain file reads, 11 worktrees here, no git process (internal/tui/model.go:528)
- note: elsewhere: core/project/project.go:56 - a hand-broken projects.json is still read as empty and the next Remember overwrites it with one entry; there before this change
- note: elsewhere: the uncommitted AGENTS.md, CLAUDE.md and .jaira/lanes/*.md in this worktree are a board regeneration to jaira 0.6.4, not this ticket - keep them out of its commits
- **2026-11-23 12:41 · Alexander Sacharov** — cockpit: accepted by the person, signoff -> done
