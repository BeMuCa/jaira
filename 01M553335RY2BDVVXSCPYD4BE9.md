---
id: 01M553335RY2BDVVXSCPYD4BE9
title: "Worker-Tabs und ihre Arbeitsbaeume schliessen sich, sobald der Pull Request ihres Tickets gemergt ist"
status: testing
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
updated-at: 2026-10-17T16:49:31Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-9571
claimed-at: 2026-10-17T16:34:08Z
outcome-what: "reap.sh closes a recorded tab only when 'herdr tab get' still labels it <ticket>/<lane> (also under --dry-run); keeps every ticket when herdr does not answer or git status fails; records go only after tabs, worktree and branch all went, so a failed removal is retried next round. Tests in core/role/reap_test.go for each."
outcome-why: "Given back from signoff, then from review: a pane id alone could close another ticket's tab after herdr reuses ids; a record removed before the worktree left a failed remove standing forever; and with herdr unreachable an open tab read as closed, so its worktree and record went and the tab was never closed."
outcome-resolves: "Goal of YD4BE9: after the PR is merged, all tabs started for the ticket close and its worktree and branch go without anyone cleaning up by hand - without closing a stranger's tab, and without losing track of a tab it could not see."
review-summary: "none"
round: 8
review-gaps: "none - one probe and one guard added, nothing duplicated (no other script probes herdr reachability; spawn.sh fails on its own tab create)"
test-verdict: "pass: go test ./... -race rc=0; DoD 1-5 verified in the tree; live on a scratch clone with real herdr.exe and gh (feat/rftqkp = merged #87): dry-run and run close only the tab labelled RFTQKP/worker, a stranger-labelled tab recorded as RFTQKP/review stays open, records gone, worktree and branch removed, rc=0; failed-remove retry covered by TestReapComesBackToAWorktreeItCouldNotRemove"
review-verdict: "Send back. The round-6 fixes are in and work (label check, records kept until removal succeeds, dry-run behind the label check), DoD 1-5 are met and 'go test ./core/role -run Reap|Spawn|Teamlead' passes. But gap 1 is a real defect, reproduced: if reap.sh cannot reach herdr it removes the worktree and drops the record of a tab that may still be open, so that tab is never closed. Small fix (one probe before the loop, one test)."
review-round: 7
review-check: "1. cd /home/alex/projects/.worktrees/jaira-yd4be9 && go test ./core/role -run 'Reap|Spawn|Teamlead' -count=1 - expect 'ok' (about 50s).  2. To see gap 1: in a scratch repo with a worktree on feat/ab12cd, a record file .git/jaira/panes/w1_p1 (pane=w1:p1, ticket=AB12CD, lane=worker, worktree=<wt>), a gh stub answering the branch tip as merged headRefOid and a jaira wrapper whose 'show' exits 5, run: HERDR_BIN_PATH=/nonexistent/herdr bash core/role/builtin/jaira-dispatcher/scripts/reap.sh - you see 'removed <wt>' and 'deleted branch feat/ab12cd', and ls .git/jaira/panes is empty. Expected after a fix: 'kept AB12CD: herdr not reachable', worktree and record still there.  3. After a fix, live in Herdr: spawn a worker with spawn.sh, merge its PR, move the ticket to done, run reap.sh --dry-run then reap.sh - the tab closes, git worktree list and git branch no longer show it."
---

# Worker-Tabs und ihre Arbeitsbaeume schliessen sich, sobald der Pull Request ihres Tickets gemergt ist

## Definition of Done

- [x] spawn.sh schreibt beim Start Pane-ID, Ticket-ID, Rolle und Worktree in eine Datei, die das Aufraeum-Skript liest; Test
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:174 (record under git common dir jaira/panes/<pane>); test core/role/reap_test.go TestSpawnRecordsItsTabForReap
- [x] Ein Skript (z.B. scripts/reap.sh neben spawn.sh) schliesst fuer jedes Ticket mit gemergtem PR alle seine Panes und entfernt Worktree und lokalen Branch; mit --dry-run zeigt es nur, was es taete; Test
  proof: core/role/builtin/jaira-dispatcher/scripts/reap.sh:175-178 (tab closed only when 'herdr tab get' labels it <ticket>/<lane>),190,201 (worktree remove, branch -D), :206 records go only after all succeeded; --dry-run prints 'would ...'; tests core/role/reap_test.go TestReapClosesAMergedTicket, TestReapNeverRemovesTheWorktreeItRunsIn, TestReapLeavesATabItDidNotOpen, TestReapComesBackToAWorktreeItCouldNotRemove
- [x] Ein Worktree mit nicht committeten Aenderungen oder ein nicht gemergter Branch wird nicht angefasst, das Skript nennt ihn und den Grund; Test
  proof: core/role/builtin/jaira-dispatcher/scripts/reap.sh:135 'waiting' (no merged PR), :143,153,160,163 'kept <ticket>: <reason>' (back in a lane, not at merged head, uncommitted changes or git status failed, herdr not reachable); test core/role/reap_test.go TestReapKeepsWhatIsNotFinished
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
- **2026-10-17 16:41 · Alexander Sacharov** — testing round 7: pass, same live run as round 6 - --dry-run now lists only the labelled tab. Scratch tabs created for the run are closed again.
- **2026-10-17 16:41 · Alexander Sacharov** — review handed to a fresh agent via run-lane.sh --no-worktree in this worktree (round 2 of review); testing changed no code, ticket file left modified.
- **2026-10-17 16:44 · Alexander Sacharov** — review round 2: send back. reap.sh reads 'herdr not reachable' as 'tab already gone' (reap.sh:168-170): with HERDR_BIN_PATH unset/wrong it removes the worktree, deletes the branch and drops the record, rc=0, no kept line - a still open tab then has no record and is never closed. Reproduced in a scratch repo with HERDR_BIN_PATH=/nonexistent/herdr. Fix: probe herdr once before the loop, keep everything with 'kept <ticket>: herdr not reachable' when it does not answer, plus a test. Everything else checked: round-6 label check and late record removal correct, tests pass.
- **2026-10-17 16:49 · Alexander Sacharov** — round 7 (review round 2 send-back): reap.sh probes 'herdr tab list' once per round (only when records exist) and keeps every merged ticket with 'kept <t>: herdr not reachable' when it fails - a closed tab and an unreachable herdr look alike through 'pane get'. Probe sits after the dirty check so a kept reason for real work still wins. Also review gap 3: git status now runs on its own, a failure keeps the ticket. Dead end: 'grep <<<"$st"' on an empty status yields one empty line and read every clean worktree as dirty - printf '%s' instead.
- **2026-10-17 16:49 · Alexander Sacharov** — critique round 7: none. Let stand: herdr down keeps per ticket rather than aborting the run - 'waiting' lines still print and a forge error still reads per ticket; one probe per round, only when records exist.
