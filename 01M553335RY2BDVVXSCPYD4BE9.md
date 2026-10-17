---
id: 01M553335RY2BDVVXSCPYD4BE9
title: "Worker-Tabs und ihre Arbeitsbaeume schliessen sich, sobald der Pull Request ihres Tickets gemergt ist"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ist der Pull Request eines Tickets gemergt, sind alle Herdr-Tabs, die fuer dieses Ticket gestartet wurden (Worker, Dispatcher, Lane- und Review-Agenten), geschlossen und sein Worktree samt lokalem Branch entfernt - ohne dass jemand von Hand aufraeumt"
context: |-
  Was heute falsch ist: Worker-Tabs bleiben nach getaner Arbeit offen. Niemand beobachtet sie, niemand schliesst sie.
  Beispiel 17.10.: S0EG0M und 8NHD8Q liefen in Herdr-Tabs w3:pHF und w3:pHH, dazu je ein Tab pro frischem Review-Agenten. Nach signoff stehen sie alle noch da, ebenso ihre Worktrees unter ~/projects/.worktrees/jaira-<slug>.
  Wer startet die Tabs: core/role/builtin/jaira-dispatcher/scripts/spawn.sh (Zeile ~220: herdr pane run). spawn.sh merkt sich nirgends, welcher Pane zu welchem Ticket gehoert - es druckt die Pane-ID nur auf stdout.
  Herdr kann schliessen: herdr pane close <pane> (herdr.exe unter WSL, Pfad in HERDR_BIN_PATH).
  Wann schliessen: erst wenn der PR gemergt ist, nicht schon beim Oeffnen - solange der PR offen ist, antwortet der Worker noch auf Review-Kommentare (jaira-role-pr).
  Es gibt keinen Daemon (CLAUDE.md: 'Nothing moves a ticket for you'). Also braucht es ein Shell-Skript neben spawn.sh, das jemand anstoesst: der Teamlead bei jeder Runde, oder ein Hook.
  Merge-Zustand: gh pr view <branch> --json state (GitHub), glab mr view (GitLab) - beide nutzt jaira-role-pr schon.
  Nicht schliessen: einen Tab mit nicht committeter Arbeit im Worktree, oder dessen Branch nicht gemergt ist.
  Verwandt: AYBDJ2 raeumt nach jaira logbook Worktree und Launcher-Eintrag auf - gleicher Zeitpunkt, andere Haelfte; zusammen denken, nicht doppelt bauen.
definition-of-done: "spawn.sh schreibt beim Start Pane-ID, Ticket-ID, Rolle und Worktree in eine Datei, die das Aufraeum-Skript liest; Test"
tags: []
blocked-by: []
related:
  - 01M4X5GJ1ME88ER0DDB2AYBDJ2
commits: []
created-at: 2026-10-17T14:08:43Z
updated-at: 2026-10-17T14:22:24Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-66355
claimed-at: 2026-10-17T14:11:12Z
outcome-what: "spawn.sh writes a record per tab (pane, ticket, lane, worktree) under <git common dir>/jaira/panes/; new jaira-dispatcher/scripts/reap.sh closes every recorded tab of a ticket that is done with a merged PR, removes its worktree and local branch, --dry-run prints only; keeps and names a ticket not done, without merged PR, with commits past the PR head, or with uncommitted changes. Teamlead skill runs it every round. Tests in core/role/reap_test.go."
outcome-why: "Worker, dispatcher and review tabs and their worktrees stayed open after the work was merged; nobody watched or closed them, and spawn.sh never recorded which tab belonged to which ticket."
outcome-resolves: "Goal of YD4BE9: after the PR is merged, all tabs started for the ticket close and its worktree and branch go without anyone cleaning up by hand - at the next teamlead round."
review-summary: |-
  core/role/builtin/jaira-dispatcher/scripts/reap.sh:105-122 checks the lane before the merge, so every ticket still in flight prints 'kept X: in <lane>, not done' each round and the teamlead skill tells it to report every kept line - noise every round. Check the merge first and print an unmerged branch as 'waiting <t>: ...' (still named, DoD 3), keep 'kept' for the merged-but-not-safe cases; teamlead SKILL.md reports closed and kept only
  core/role/builtin/jaira-dispatcher/scripts/reap.sh:126-128 reports 'has commits the merged pull request does not' also when the PR head is simply not fetched locally (merge-base fails); say 'is not at the merged head <sha>' which is true in both cases
  core/role/builtin/jaira-dispatcher/scripts/reap.sh:98 only the main worktree is protected; run from a ticket worktree the script would remove the directory it runs in. Treat $root like $main in the script and drop the prose warning in jaira-teamlead/SKILL.md
---

# Worker-Tabs und ihre Arbeitsbaeume schliessen sich, sobald der Pull Request ihres Tickets gemergt ist

## Definition of Done

- [x] spawn.sh schreibt beim Start Pane-ID, Ticket-ID, Rolle und Worktree in eine Datei, die das Aufraeum-Skript liest; Test
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:174 (record under git common dir jaira/panes/<pane>); test core/role/reap_test.go TestSpawnRecordsItsTabForReap
- [x] Ein Skript (z.B. scripts/reap.sh neben spawn.sh) schliesst fuer jedes Ticket mit gemergtem PR alle seine Panes und entfernt Worktree und lokalen Branch; mit --dry-run zeigt es nur, was es taete; Test
  proof: core/role/builtin/jaira-dispatcher/scripts/reap.sh:146,159,165 (tab close, worktree remove, branch -D), --dry-run prints 'would ...'; test core/role/reap_test.go TestReapClosesAMergedTicket
- [x] Ein Worktree mit nicht committeten Aenderungen oder ein nicht gemergter Branch wird nicht angefasst, das Skript nennt ihn und den Grund; Test
  proof: core/role/builtin/jaira-dispatcher/scripts/reap.sh 'kept <ticket>: <reason>' for no merged PR, commits past the PR head, uncommitted changes, ticket back in a lane; test core/role/reap_test.go TestReapKeepsWhatIsNotFinished
- [x] Der Teamlead-Skill ruft das Skript bei jeder Runde auf und meldet, was geschlossen wurde
  proof: core/role/builtin/jaira-teamlead/SKILL.md:26-39 (reap.sh every round, report closed/kept); test core/role/reap_test.go TestTeamleadReapsEveryRound
- [x] Eine Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/YD4BE9.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] spawn.sh: write the pane registry file right after tab create
- [x] scripts/reap.sh: registry -> merged check -> close tabs, remove worktree and branch; --dry-run; skip dirty/unmerged/came-back with reason
- [x] Go tests in core/role/role_test.go with herdr, gh and jaira stubs over a real git repo
- [x] teamlead SKILL.md: run reap.sh every round, report; reconcile the 'worktree until done' paragraph; dispatcher/worker mention
- [x] release note core/release/notes/YD4BE9.md

## Progress
- **2026-10-17 14:11 · Alexander Sacharov** — parent session: jaira-79
- **2026-10-17 14:12 · Alexander Sacharov** — decisions settled (none open for a person):
1) Registry: one file per pane under $(git common dir)/jaira/panes/ - machine-local like Herdr pane ids, shared by every worktree of the repo, never in git status, no append race between parallel spawns (one file each). Line format: key=value lines pane, ticket, lane, worktree.
2) Merged = gh pr list --head <branch> --state merged (gh pr view <branch> finds nothing in a fork clone - checked on feat/rftqkp: list finds #87, view does not). GitLab: glab mr list --source-branch <b> --merged. Forge ladder as jaira-role-pr: git config jaira.forge, else origin host; neither -> skip with reason.
3) 'Branch not merged' also covers a local branch whose tip is not the merged PR head (commits after the merge): skip and name it. Squash merges make git branch -d useless, so the branch goes with -D after that check.
4) Teamlead skill says worktree stays until done, a merged PR can still come back. Reconciled: reap also skips a ticket that sits on the board in a lane other than done (filed away or done = ok). jaira-role-pr files tickets away before the push, so merged normally means filed.
5) Dirty = git status --porcelain without untracked files under .jaira/ (AYBDJ2's rule). The main worktree is never removed, only its panes closed. The reaper never closes its own pane (HERDR_PANE_ID).
6) Close by tab (pane get -> tab_id -> tab close), the same call run-lane.sh already uses; a pane already gone just drops its registry file.
7) Launcher entry in ~/.jaira/projects.json is AYBDJ2's half, not built here.
- **2026-10-17 14:19 · Alexander Sacharov** — found out: gh pr list --head finds a fork's merged PR, gh pr view <branch> does not. macOS bash 3.2 rules out declare -A/mapfile in reap.sh. Reap tests skip on Windows: Git Bash /c/... vs git.exe C:\ paths never compare. Old tabs/worktrees from before this change have no record and are not reaped - close those by hand once.
