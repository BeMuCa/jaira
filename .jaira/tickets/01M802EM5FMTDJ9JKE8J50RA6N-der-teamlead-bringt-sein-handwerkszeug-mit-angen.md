---
id: 01M802EM5FMTDJ9JKE8J50RA6N
title: "Der Teamlead bringt sein Handwerkszeug mit: angenommene PRs landen, Branches stapeln, saubere Logbook-Commits"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein neuer Teamlead hat die Routinen, die am 21.11. von Hand gebaut wurden, als Skript und Skill-Text, statt sie neu zu erfinden."
context: |-
  Was heute (21.11.) von Hand gebaut wurde und jeder Teamlead wieder braucht (Alex: "speichere es, damit künftige Agenten es haben"):
  1. Einen angenommenen PR landen: master in den Ticket-Worktree mergen, generierte Dateien neu erzeugen (hier: scripts/shipped-history.sh --write, weil core/shipped/shipped.sha256 fast jeden PR in Konflikt bringt), pushen, auf grünes CI warten, mergen, reap.sh. Heute als ~/.claude/scripts/jaira-land.sh (Vorlage), 20+ PRs damit gelandet.
  2. Folgearbeit auf einen angenommenen, noch nicht gemergten Branch stapeln statt auf CI zu warten.
  3. Logbook-Commit nur mit .jaira/tickets und .jaira/logbook (git add -A .jaira hat untracked Milestones nach master gebracht, PR #167).
  4. Give back nach signoff = Note + move + neuer Worker per spawn.sh (alte Tabs schliessen sich seit 0.5.15).
  5. Antworten der Person als dod --add, nicht nur als Note (1MYFKE kam sonst ohne die Entscheidung nach signoff).
  Grenze: Mergen bleibt Sache des Maintainers — die Rolle darf nur landen, wo das Board es erlaubt (wie .jaira/pr agent / none), sonst gibt sie die Zeile zurück.
  Teile 3-5 sind Text in core/role/builtin/jaira-teamlead/SKILL.md; Teil 1 ein Skript neben der Rolle (projektspezifische Generatoren über eine Board-Datei, z. B. .jaira/land-regenerate).
definition-of-done: "Skript neben jaira-teamlead landet einen PR (merge master, Regenerieren per Board-Datei, Push, CI abwarten, merge, reap) nur, wo das Board Landen erlaubt; sonst gibt es die Zeile zurück (Test)"
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-11-21T22:08:26Z
updated-at: 2026-11-22T11:28:41Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-31706
claimed-at: 2026-11-22T11:08:04Z
outcome-what: "New core/role/builtin/jaira-teamlead/scripts/land.sh <pr>: merges the PR's base into the ticket branch, takes conflicts in files .jaira/land-regenerate names from the base and regenerates them, commits merge+regeneration as one chore(<ID>) commit, pushes, waits for gh pr checks, then merges via gh only where 'jaira roles pr-gate --merge' says yes, else prints the merge line (exit 3). New pr-gate --merge reads .jaira/land off the base (agent only) and reuses the existing every-ticket-in-done check. Teamlead SKILL.md: give back as note+move+spawn, decisions as dod --add, landing/stacking/filing-away section."
outcome-why: "The 2026-11-21 teamlead built these routines by hand (~/.claude/scripts/jaira-land.sh, 20+ PRs landed); a new teamlead would rebuild them, and merging must stay behind a board permission (Alex chose .jaira/land = agent)."
outcome-resolves: "DoD 1-4: land.sh + TestLandMergesOnlyWhereTheBoardLetsIt (9 cases), pr-gate --merge + two prgate tests, SKILL text, release note + manifest. Four categories checked: concurrency — two land.sh on one worktree not guarded (teamlead lands one after another, written in SKILL); error paths — conflict outside generated files aborts the merge, failed generator resets to the start, red CI stops before the gate, each tested where cheap; unusual input — no regenerate file, comment/blank lines, branch already holding the base tested; permissions — regenerate commands and .jaira/land come from the base, never the PR branch."
review-summary: none
round: 3
back-from: testing
back-by: lane
back-count: 1
review-gaps: "removed the unused store field and a never-taken clean-tree branch from core/role/land_test.go; land.sh no longer lists origin twice in its remote loop. Left: land.sh's fork/remote resolution repeats reap.sh's (a running script cannot be sourced; critique let it stand); the four MERGE_HEAD checks in land.sh read clearer than a helper."
test-verdict: "fail: go test ./... RC=1 — internal/wintrap TestRepositoryIsClean: core/role/land_test.go:96 windows trap 5 (git refname compared with HasPrefix; needs //wintrap:ok). DoD 1-4 verified in the tree otherwise; land.sh bad arg exit 2, dirty tree exit 2, pr-gate --merge on upstream/master exit 3 with 'a person merges'. land.sh writes only git state: two at once on one worktree meet git's index.lock."
hunt-verdict: "fail: bughunt done at 7ddcf45f3, 3 change (none a blocker alone, ride along with the red gate): --merge skips the unnamed-commit check; gh pr merge without --match-head-commit; a regenerated file the base deleted leaves a half merge"
---

# Der Teamlead bringt sein Handwerkszeug mit: angenommene PRs landen, Branches stapeln, saubere Logbook-Commits

## Definition of Done

- [x] Skript neben jaira-teamlead landet einen PR (merge master, Regenerieren per Board-Datei, Push, CI abwarten, merge, reap) nur, wo das Board Landen erlaubt; sonst gibt es die Zeile zurück (Test)
  proof: core/role/builtin/jaira-teamlead/scripts/land.sh; test TestLandMergesOnlyWhereTheBoardLetsIt (core/role/land_test.go:20)
- [x] jaira-teamlead SKILL.md enthält Stapeln, Logbook-Staging, Give-back-per-spawn und Entscheidungen-als-DoD
  proof: core/role/builtin/jaira-teamlead/SKILL.md:331 (give back: note, move, spawn), :372 (decisions as dod --add), :383 (landing, stacking, filing away by path)
- [x] Release-Note core/release/notes/<ID>.md, shipped.sha256
  proof: core/release/notes/50RA6N.md; core/shipped/shipped.sha256 regenerated (go test ./core/role ./core/shipped green)
- [x] Merge nur bei .jaira/land = agent (von der Basis gelesen), alle Tickets der Branch in done und CI grün; sonst druckt das Skript die Merge-Zeile für die Person (Entscheidung Alex 22.11.)
  proof: internal/cli/roles.go:215 (pr-gate --merge reads .jaira/land off the base, unnamed commits refused); tests TestPRGateMergesOnlyWhereTheBaseSaysAgentAndEveryTicketIsDone, TestPRGateMergeRefusesWhatItCannotJudge, TestLandMergesOnlyWhereTheBoardLetsIt

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 11:05 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2784.sock
- **2026-11-22 11:08 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock (previous parent 2784.sock no longer exists; claim by dead PID 12556 stolen)
- **2026-11-22 11:08 · Alexander Sacharov** — question (open decision before in-progress): Which board setting lets the land script merge a PR? Today .jaira/pr knows person|agent|none and 'agent' explicitly never merges. 1) (Recommended) new one-line file .jaira/land saying 'agent': merge only when every ticket on the branch is in done and CI is green; anything else -> print the merge line for the person. Orthogonal to .jaira/pr, read from the base like pr-gate. 2) new value in .jaira/pr, e.g. 'land' = agent opens AND merges. Fewer files, but mixes opening and merging in one switch. 3) no board file: script merges only with an explicit --by-person flag (person said 'land it' in this run). Simplest, but the board itself never records the permission. Settled by me: regeneration via .jaira/land-regenerate (one shell command per line, run in order from repo root, changed files staged); CI wait via gh only (glab out of scope).
- **2026-11-22 11:11 · Alexander Sacharov** — decision (Alex 22.11.): Variante 1 — neue einzeilige Datei .jaira/land = 'agent': das land-Skript merged nur, wenn jedes Ticket der Branch in done ist und CI grün; sonst druckt es die Merge-Zeile für die Person. Getrennt von .jaira/pr, von der Basis gelesen wie pr-gate.
- **2026-11-22 11:12 · Alexander Sacharov** — plan (in-progress, no plan lane): (a) pr-gate gets --merge: exit 0 only where .jaira/land in <base> says agent and every ticket named on the branch is in done (pr_when ignored, .jaira/pr not read — a PR already exists); refusal reason not_agent / pending. Reuses the existing ticket check instead of a second copy in shell. (b) core/role/builtin/jaira-teamlead/scripts/land.sh <pr>: merge base into the worktree, conflicts in files .jaira/land-regenerate names ('<path>: <command>' per line) are taken from base and regenerated, any other conflict stops; push; gh pr checks --watch; pr-gate --merge; merge + reap.sh, or print the gh pr merge line. (c) Go test with stubbed gh in core/role, like reap_test. (d) SKILL.md text for parts 2-5. (e) release note + shipped manifest.
- **2026-11-22 11:21 · Alexander Sacharov** — in-progress learned: (1) jaira spawns detached fetch/snapshot children that keep writing into a test's TempDir — land_test sets JAIRA_NO_FETCH/NO_SNAPSHOT/NO_UPDATE_CHECK or cleanup fails. (2) shipped-history.sh only hashes tracked files: git add a new installed script before --write. (3) Known limits of land.sh, chosen not fixed: gh only (no glab); a repo with no CI checks reads as 'not green' (gh pr checks exits 1 on 'no checks reported') so the person merges; JAIRA_LAND_SETTLE (30s) is a guess for checks to register. (4) Regenerate commands are read from the base, not the PR branch, so a branch cannot inject a command. (5) pr-gate --merge does not check unnamed commits (unlike none): a PR exists and its diff was read.
- **2026-11-22 11:22 · Alexander Sacharov** — critique: (1) roles.go --merge refusal: name the 'agent here but not in <base>' case like pr_person does. (2) pr-gate Short: open, land or merge. Let stand: land.sh repeats reap.sh's fork/remote resolution (~15 lines) because a script cannot source another that runs on load; without reap's gh fallback for renamed remotes it stops with 'git remote add' — acceptable. The CLAUDE.md arrival note still says merging stays the maintainer's on every board; only the teamlead role knows the .jaira/land exception — workers never merge, so left as is.
- **2026-11-22 11:22 · Alexander Sacharov** — critique pass 2: both findings of pass 1 answered in c826ccc39 (disk-only hint + test case, Short); nothing else re-read.
- **2026-11-22 11:25 · Alexander Sacharov** — bughunt: change: internal/cli/roles.go:263 (with land.sh:145) - class: pr-gate --merge never runs the UnnamedCommits check that the --land path runs, so commits on the branch that name no ticket are merged by an agent with no person deciding; the comment at roles.go:258 states the rule for unattended landing ('every commit must belong to a ticket a person accepted') and --merge is unattended too; example: feat(<ID>) in done plus a later 'wip: no ticket here' commit -> land.sh exit 0, 'merged'; ran: scratch test in core/role on newLandRig (land agent, ticket done, one unnamed commit) -> exit 0 merged=true; not a blocker by itself: a pull request exists and could have been read, but nobody is required to read the commits added after acceptance
- **2026-11-22 11:25 · Alexander Sacharov** — bughunt: change: core/role/builtin/jaira-teamlead/scripts/land.sh:155 - class: time-of-check/time-of-use between the gate and the merge: pr-gate --merge (line 145) judges the local HEAD, CI is watched on whatever head the PR has, and 'gh pr merge' merges the PR's remote head as it is at merge time, without --match-head-commit; example: a second session pushes to feat/<id> during the 30 s settle or the CI watch -> its commits are merged by the agent unjudged by the gate; also the same push makes gh pr checks report the new head's checks, not the pushed one's; fix class: pass --match-head-commit $(git rev-parse HEAD) to gh pr merge; confirmed by read (gh pr merge --help lists --match-head-commit); not a blocker: needs a concurrent push to the ticket's own branch
- **2026-11-22 11:25 · Alexander Sacharov** — bughunt: change: core/role/builtin/jaira-teamlead/scripts/land.sh:110 - class: a conflict on a generated file that 'git checkout --theirs' cannot resolve (modify/delete: the base removed or renamed a file .jaira/land-regenerate still names) kills the script under set -e with a raw git error and leaves the worktree mid-merge (MERGE_HEAD set, 'UD gen.txt'), against the usage text 'aborts the merge and stops (exit 1)'; the next land.sh run then refuses with 'uncommitted changes'; example: branch edits gen.txt, master git rm gen.txt; ran: scratch test on newLandRig -> exit 1, out 'error: path gen.txt does not have their version', status 'UD gen.txt'; edge, not a blocker: nothing is pushed and git merge --abort recovers
- **2026-11-22 11:26 · Alexander Sacharov** — bughunt: done at 7ddcf45f3 - 3 findings (3 change, 0 blocker, 0 pre-existing) in 6 modules, first hunt; concurrency, error paths, input, permissions checked; files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/optimize.md .jaira/lanes/pre-process.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/release/notes/50RA6N.md core/role/builtin/jaira-teamlead/SKILL.md core/role/builtin/jaira-teamlead/scripts/land.sh core/role/land_test.go core/settings/settings.go core/shipped/shipped.sha256 internal/cli/prgate_test.go internal/cli/roles.go
- **2026-11-22 11:26 · Alexander Sacharov** — testing round 1 -> in-progress: (1) //wintrap:ok on land_test.go:96 (git refname). (2) roles.go: run UnnamedCommits under --merge too. (3) land.sh: gh pr merge --match-head-commit <pushed HEAD>. (4) land.sh: checkout --theirs failing (base deleted the generated file) counts as an unresolvable conflict -> abort, exit 1. Test each.
