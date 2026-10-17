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
updated-at: 2026-10-17T16:40:05Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-9571
claimed-at: 2026-10-17T16:34:08Z
outcome-what: "reap.sh closes a recorded tab only when 'herdr tab get' still labels it <ticket>/<lane> - a reused pane id in a stranger's tab drops the record and closes nothing, and --dry-run applies the same check; records go only after tabs, worktree and branch all went, so a failed git worktree remove or branch -D is retried next round. Two new tests in core/role/reap_test.go."
outcome-why: "Given back from signoff: records outlive tabs by days and herdr may reuse ids after a restart, so a pane id alone could close another ticket's live worker; and a record removed before the worktree left a failed remove standing forever."
outcome-resolves: "Goal of YD4BE9: after the PR is merged, all tabs started for the ticket close and its worktree and branch go without anyone cleaning up by hand - now without risk of closing a tab that is not the ticket's, and with retries when removal fails."
review-summary: "none"
round: 7
review-gaps: "none - no other script reads a tab label ('tab get' is new here); the per-tab rm the change removed has no leftover; 'gone' differs from 'recs' only by the reaper's own pane, so it is not a duplicate list"
test-verdict: "fail: reap.sh --dry-run prints 'would close' for a tab whose label names another ticket, which the real run then leaves open"
review-verdict: "Approve with one reservation. The diff meets all four DoD items and the safety checks (merged head, done lane, dirty worktree, own/main worktree) are sound and tested. Not fully sure about gap 1: if herdr can reuse pane ids after a restart, reap.sh can close a stranger's tab. A label check before 'tab close' would remove that risk — worth deciding before signoff."
review-round: 5
review-check: "1. cd /home/alex/projects/.worktrees/jaira-yd4be9 && go test ./core/role -run 'Reap|Spawn|Teamlead' -count=1 — expect 'ok' (about 50s).  2. In a Herdr session, from the main checkout, start a worker with core/role/builtin/jaira-dispatcher/scripts/spawn.sh <slug> <ID> worker . — then ls .git/jaira/panes/ shows one file with pane=, ticket=, lane=, worktree= lines.  3. Run .claude/skills/jaira-dispatcher/scripts/reap.sh --dry-run — expect 'waiting <ID>: feat/<slug> has no merged pull request'.  4. For a ticket that is done and whose PR is merged, run reap.sh --dry-run — expect 'would close <ID>/<lane> <pane>', 'would remove <worktree>', 'would delete branch feat/<slug>'; nothing changes.  5. Run reap.sh without --dry-run — the tabs close, git worktree list no longer shows the worktree, git branch no longer shows the branch.  6. Restart herdr, open a new tab, then check whether its pane id equals an old id in .git/jaira/panes/ — if yes, gap 1 is real."
---

# Worker-Tabs und ihre Arbeitsbaeume schliessen sich, sobald der Pull Request ihres Tickets gemergt ist

## Definition of Done

- [x] spawn.sh schreibt beim Start Pane-ID, Ticket-ID, Rolle und Worktree in eine Datei, die das Aufraeum-Skript liest; Test
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:174 (record under git common dir jaira/panes/<pane>); test core/role/reap_test.go TestSpawnRecordsItsTabForReap
- [x] Ein Skript (z.B. scripts/reap.sh neben spawn.sh) schliesst fuer jedes Ticket mit gemergtem PR alle seine Panes und entfernt Worktree und lokalen Branch; mit --dry-run zeigt es nur, was es taete; Test
  proof: core/role/builtin/jaira-dispatcher/scripts/reap.sh:169-172 (tab closed only when 'herdr tab get' labels it <ticket>/<lane>),184,195 (worktree remove, branch -D), :200 records go only after all succeeded; --dry-run prints 'would ...'; tests core/role/reap_test.go TestReapClosesAMergedTicket, TestReapNeverRemovesTheWorktreeItRunsIn, TestReapLeavesATabItDidNotOpen, TestReapComesBackToAWorktreeItCouldNotRemove
- [x] Ein Worktree mit nicht committeten Aenderungen oder ein nicht gemergter Branch wird nicht angefasst, das Skript nennt ihn und den Grund; Test
  proof: core/role/builtin/jaira-dispatcher/scripts/reap.sh:131 'waiting' (no merged PR), :139,149,157 'kept <ticket>: <reason>' (back in a lane, not at merged head, uncommitted changes); test core/role/reap_test.go TestReapKeepsWhatIsNotFinished
- [x] Der Teamlead-Skill ruft das Skript bei jeder Runde auf und meldet, was geschlossen wurde
  proof: core/role/builtin/jaira-teamlead/SKILL.md:26-38 (reap.sh every round, report closed/kept); test core/role/reap_test.go TestTeamleadReapsEveryRound
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
- **2026-10-17 14:22 · Alexander Sacharov** — critique round 1: (1) reap.sh checks the lane before the merge -> every in-flight ticket prints a kept line each round; check merge first, unmerged prints 'waiting', teamlead reports closed+kept only. (2) 'has commits the merged PR does not' is wrong when the PR head is just not fetched -> 'is not at the merged head <sha>'. (3) protect $root like $main, drop the prose warning. Let stand: tab-close code duplicated from run-lane.sh (two lines, same pattern); shell script instead of Go in jaira logbook - the ticket asks for a script beside spawn.sh.
- **2026-10-17 14:28 · Alexander Sacharov** — critique round 2: reap.sh root was $PWD, not the worktree top -> started in a subdirectory it removed its own worktree. Fixed with rev-parse --show-toplevel, test extended (runs from <wt>/sub). Round 1 findings verified fixed in f67638bc.
- **2026-10-17 14:29 · Alexander Sacharov** — critique round 3: none. Let stand: tab-close lines duplicated from run-lane.sh (same two calls, separate scripts); stale records of tickets never merged stay as 'waiting' lines - harmless, teamlead does not report them.
- **2026-10-17 14:31 · Alexander Sacharov** — testing: live run on a scratch clone (real gh, real herdr.exe, feat/rftqkp at 1e2ace72 = merged #87): with only origin=sashasoft90/jaira every ticket prints 'waiting ... no merged pull request' - gh pr list picks the fork. With an upstream remote added it works: dry-run listed, real run removed worktree (untracked .jaira/x.md let through) and deleted the branch, pane already gone -> record dropped, rc=0. Fix: resolve the target repo once like jaira-role-pr (gh repo view <origin url> --json isFork,parent,nameWithOwner; fork -> parent) and pass --repo to gh pr list. Suite: go test ./... -race rc=0.
- **2026-10-17 14:33 · Alexander Sacharov** — testing fix: reap.sh resolves the PR repo once (gh repo view <origin url>: fork -> parent, else nameWithOwner) and passes --repo; live run on the scratch fork clone without upstream remote now removes feat/rftqkp's worktree and branch. GitLab forks still ask glab without a target project - not exercised, no GitLab remote here.
- **2026-10-17 14:33 · Alexander Sacharov** — critique round 4 (after testing fix): none; the repo lookup mirrors jaira-role-pr's ladder rungs 1-2, rung 3 (a board remote naming a third repo) is not needed to find a merged PR.
- **2026-10-17 14:36 · Alexander Sacharov** — review handed to a fresh agent via run-lane.sh --no-worktree in this worktree; testing lane changed no code, ticket file left modified for the next code commit or the logbook commit.
- **2026-10-17 14:38 · Alexander Sacharov** — review: approve with reservation. reap.sh closes a tab by pane id alone; records outlive their tabs (run-lane.sh closes lane tabs early), so if herdr reuses pane ids after a restart a stranger's tab gets closed. Suggested fix: match the tab label '<ticket>/<lane>' from 'herdr pane get' before 'tab close'. Not verified against real herdr. Minor: records are rm'd before 'git worktree remove' — a failed remove is never retried.
- **2026-10-17 16:33 · Alexander Sacharov** — explain: reap.sh schliesst einen Tab nur anhand der Pane-Id. Die Records ueberleben ihre Tabs oft um Tage, und nach einem herdr-Neustart koennten dieselben Ids neu vergeben werden.
Warum: Dann schliesst reap.sh einen fremden Tab, womoeglich den lebenden Worker eines anderen Tickets. Ausserdem wird der Record vor 'git worktree remove' geloescht, schlaegt das Entfernen fehl, bleiben Worktree und Branch fuer immer stehen.
1) Zurueckgeben: vor 'tab close' das Tab-Label aus 'herdr pane get' mit '<ticket>/<lane>' vergleichen und den Record erst nach erfolgreichem 'git worktree remove' loeschen [give back]
2) So annehmen: alle vier DoD-Punkte sind erfuellt, die Id-Wiederverwendung ist nicht belegt, die Luecken bleiben offen [accept]
Empfehlung: 1 - die Label-Pruefung ist billig und nimmt das einzige Risiko weg, bei dem ein fremder Tab verloren gehen kann.
- **2026-10-17 16:33 · Alexander Sacharov** — decision: 1) Zurueckgeben: vor 'tab close' das Tab-Label aus 'herdr pane get' mit '<ticket>/<lane>' vergleichen und den Record erst nach erfolgreichem 'git worktree remove' loeschen
- **2026-10-17 16:33 · Alexander Sacharov** — Given back from signoff: 1) Zurueckgeben: vor 'tab close' das Tab-Label aus 'herdr pane get' mit '<ticket>/<lane>' vergleichen und den Record erst nach erfolgreichem 'git worktree remove' loeschen
- **2026-10-17 16:34 · Alexander Sacharov** — cockpit: given back to a new worker - old worker's prompt cache cold for 55 min
- **2026-10-17 16:34 · Alexander Sacharov** — round 6 found out: the tab label lives on 'herdr tab get <tab>' -> result.tab.label (checked on real herdr.exe: w3:tGH -> 'YD4BE9/worker'); 'pane get' carries no label. spawn.sh labels '<ticket>/<lane>' exactly as the record's ticket and lane fields, so reap compares against those, case kept. A label that does not match means our tab is gone and the id belongs to someone else: drop the record, close nothing.
- **2026-10-17 16:36 · Alexander Sacharov** — round 6 (given back from signoff): records are no longer removed per tab but all at the end, and only when tabs, worktree and branch all went - a failed remove or branch -D prints 'kept' and the next round retries. The reaper's own pane record is still left alone. Tests verified to fail on the previous reap.sh (both new ones) and pass now.
- **2026-10-17 16:37 · Alexander Sacharov** — critique round 5 (after the signoff give-back): none. Let stand: a branch -D that fails after its worktree went is retried under the default prefix/<ticket> name, since the worktree that named the branch is gone - spawn.sh names the branch that way, so only a hand-renamed branch would be missed. Label read per record (two herdr calls) instead of one 'tab list': the per-pane path is the one run-lane.sh already uses.
- **2026-10-17 16:38 · Alexander Sacharov** — testing round 6: live on a scratch clone (real herdr.exe, gh against BeMuCa/jaira, feat/rftqkp = merged #87): tab labelled RFTQKP/worker closed, tab labelled ZZ99ZZ/worker but recorded as RFTQKP/review left open, records both gone, worktree removed, branch deleted, rc=0. Finding: --dry-run printed 'would close RFTQKP/review w3:pJ1' for the stranger's tab - the dry branch sits before the label check. Fix: look the tab and its label up first, then print 'would close' or close.
- **2026-10-17 16:40 · Alexander Sacharov** — critique round 6: none - the dry branch now sits behind the label check, one line moved; dry-run therefore asks herdr too, which is read-only.
