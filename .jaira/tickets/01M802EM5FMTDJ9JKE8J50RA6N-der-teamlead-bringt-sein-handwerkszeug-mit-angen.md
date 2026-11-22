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
updated-at: 2026-11-22T11:22:32Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-31706
claimed-at: 2026-11-22T11:08:04Z
outcome-what: "New core/role/builtin/jaira-teamlead/scripts/land.sh <pr>: merges the PR's base into the ticket branch, takes conflicts in files .jaira/land-regenerate names from the base and regenerates them, commits merge+regeneration as one chore(<ID>) commit, pushes, waits for gh pr checks, then merges via gh only where 'jaira roles pr-gate --merge' says yes, else prints the merge line (exit 3). New pr-gate --merge reads .jaira/land off the base (agent only) and reuses the existing every-ticket-in-done check. Teamlead SKILL.md: give back as note+move+spawn, decisions as dod --add, landing/stacking/filing-away section."
outcome-why: "The 2026-11-21 teamlead built these routines by hand (~/.claude/scripts/jaira-land.sh, 20+ PRs landed); a new teamlead would rebuild them, and merging must stay behind a board permission (Alex chose .jaira/land = agent)."
outcome-resolves: "DoD 1-4: land.sh + TestLandMergesOnlyWhereTheBoardLetsIt (9 cases), pr-gate --merge + two prgate tests, SKILL text, release note + manifest. Four categories checked: concurrency — two land.sh on one worktree not guarded (teamlead lands one after another, written in SKILL); error paths — conflict outside generated files aborts the merge, failed generator resets to the start, red CI stops before the gate, each tested where cheap; unusual input — no regenerate file, comment/blank lines, branch already holding the base tested; permissions — regenerate commands and .jaira/land come from the base, never the PR branch."
review-summary: |-
  internal/cli/roles.go:215 pr-gate --merge says only 'does not say agent' when .jaira/land says agent on this checkout's disk but not on the base; the --merge refusal should name that case like the pr_person branch below it does (board.PRWordIn of the disk file), or a teamlead who just wrote the file reads it as broken
  internal/cli/roles.go:147 Short still says 'whether an agent may open this branch's pull request'; it now also answers landing and merging — say open, land or merge
round: 2
back-from: critique
back-by: lane
back-count: 1
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
  proof: internal/cli/roles.go:214 (pr-gate --merge reads .jaira/land off the base); tests TestPRGateMergesOnlyWhereTheBaseSaysAgentAndEveryTicketIsDone, TestLandMergesOnlyWhereTheBoardLetsIt (no file/person/signoff -> exit 3 with merge line; red CI -> exit 1)

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
