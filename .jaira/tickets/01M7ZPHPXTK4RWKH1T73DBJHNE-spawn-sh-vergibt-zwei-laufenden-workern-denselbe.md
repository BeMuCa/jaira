---
id: 01M7ZPHPXTK4RWKH1T73DBJHNE
title: spawn.sh vergibt zwei laufenden Workern denselben Port-Offset
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: Zwei gleichzeitig laufende Worker bekommen von spawn.sh nie denselben Port-Offset.
context: |-
  Falsch heute: core/role/builtin/jaira-dispatcher/scripts/spawn.sh (Zeile ~186) rechnet den Offset als cksum(slug) % 40 + 1.
  Der Kommentar darüber verspricht "two workers never share a stack" — das stimmt nicht.
  Beispiel: 8r8vts, 32ttss und twb0r6 bekommen alle Offset 5.
  Bei 40 Plätzen ist schon bei 8 parallelen Workern eine Kollision wahrscheinlicher als nicht (Geburtstagsproblem, ~52 %).
  Folge: .jaira/worktree-setup schreibt zwei Worktrees dieselben Ports; der zweite Stack startet nicht oder redet mit dem ersten.
  Gefunden beim Review von 8R8VTS (dort nicht angefasst).
  Idee, nicht entschieden: Offset beim Anlegen des Worktrees vergeben und belegte Offsets der vorhandenen Worktrees überspringen (z. B. aus einer Datei im Worktree lesen), statt Hash.
definition-of-done: "Zwei Slugs mit gleichem cksum-Rest bekommen verschiedene Offsets, solange beide Worktrees existieren (Test in core/role)"
tags:
  - concurrency
blocked-by: []
related: []
follows: 01M7Z6T8TS8CRXA4FH218R8VTS
commits: []
created-at: 2026-11-21T18:40:24Z
updated-at: 2026-11-21T20:06:29Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-987290
claimed-at: 2026-11-21T20:03:04Z
outcome-what: "spawn.sh hands each new worktree a port offset no other existing worktree of the repository holds: starts at the old slug hash, skips taken ones, under a lock written to jaira-offset in the worktree's git dir; a lock older than 120 s is taken over, a lock that cannot be made stops spawn.sh at once; worktrees whose directory is gone hold nothing, relative gitdir paths resolved against the worktree's git dir"
outcome-why: "cksum%40 alone gave 8r8vts, 32ttss and twb0r6 the same offset; two stacks on the same ports. Given back: a killed spawn must not block every later one, an unmakeable lock must not cost 120 s"
outcome-resolves: Zwei gleichzeitig laufende Worker bekommen nie denselben Port-Offset
review-summary: "spawn.sh hands out a worktree's port offset (1-40) once, when the worktree is made: it starts at the old cksum(slug)%40+1 and takes the next offset no other standing worktree of the repo holds, wrapping. What a worktree holds is read from a new jaira-offset file in its git dir, else the offset= in its jaira-setup, else the '# worker stack' line of its .env (worktrees from before). Worktrees whose directory is gone hold nothing; a relative gitdir (git 2.48 useRelativePaths) is resolved against the git dir. A mkdir lock in the git common dir covers look-up, 'git worktree add' and writing jaira-offset, so two parallel spawns cannot pick the same free slot; it waits up to 120 s. With all 40 taken it exits 1 before anything is made, only where a setup or .env needs ports. USAGE and core/release/notes/DBJHNE.md say so; role_test covers collision, keep-on-respawn, freed offset, deleted worktree and relative gitdir."
review-gaps: |-
  note: no test runs two spawn.sh in parallel - the lock is only proven by the bughunt's scratch run; a test with a post-checkout hook that sleeps 1 s and two spawns at once would pin it
  note: mkdir of the lock failing for another reason than "exists" (read-only common dir, permissions) waits the full 120 s and then blames a dead spawn.sh; fix: stop at once when the lock dir does not exist after a failed mkdir
  note: a SIGKILLed spawn leaves jaira-offset.lock, and every later spawn waits 120 s and fails until a person removes it (known, accepted in critique); a lock older than 120 s by mtime could be taken over instead
  note: two spawns of the same slug at once - the second waits for the lock, then fails in git worktree add because the first made the directory; was a race before too, now it is at least serialised
  note: branch feat/dbjhne carries JDDCEZ's 9 commits under its own 3 (git log master..HEAD) - JDDCEZ must land first or the PR shows both
test-verdict: "pass: round 3 (since 231d2a909): go test -race ./role ./shipped green (RC=0); TestSpawnGivesEveryWorktreeItsOwnOffset now rewrites twb0r6's gitdir relative (as git 2.48 useRelativePaths does) and k2x2yk, hashing to 5, gets 8. Rounds 1-2: full suite -race green, parallel spawns distinct (also with a 12 s post-checkout hook), all-40-taken and foreign lock exit 1 leaving nothing"
hunt-verdict: "pass: bughunt done at d263a0319, 0 change, 0 pre-existing (rounds 1-2 change findings fixed in 231d2a909 and d263a0319)"
round: 4
back-from: signoff
back-by: person
back-count: 1
review-verdict: "the diff meets all three DoD items; no blocker found. Offset allocation is correct for standing, legacy, deleted and relative-gitdir worktrees, and the lock covers the race the goal is about. Open points are edge cases of the lock (no parallel test, stale-lock and permission paths) - known limits, not wrong results. go test ./core/role ./core/shipped passes."
review-round: 3
review-verdict-gaps: |-
  note: no test runs two spawn.sh in parallel - the lock is only proven by the bughunt's scratch run; a test with a post-checkout hook that sleeps 1 s and two spawns at once would pin it
  note: mkdir of the lock failing for another reason than "exists" (read-only common dir, permissions) waits the full 120 s and then blames a dead spawn.sh; fix: stop at once when the lock dir does not exist after a failed mkdir
  note: a SIGKILLed spawn leaves jaira-offset.lock, and every later spawn waits 120 s and fails until a person removes it (known, accepted in critique); a lock older than 120 s by mtime could be taken over instead
  note: two spawns of the same slug at once - the second waits for the lock, then fails in git worktree add because the first made the directory; was a race before too, now it is at least serialised
  note: branch feat/dbjhne carries JDDCEZ's 9 commits under its own 3 (git log master..HEAD) - JDDCEZ must land first or the PR shows both
review-check: "1. cd /home/alex/projects/.worktrees/jaira-dbjhne  2. go test ./core/role -run TestSpawnGivesEveryWorktreeItsOwnOffset -v  3. expect PASS: 8r8vts and 32ttss (same hash, 5 held by twb0r6) get 6 and 7, 8r8vts respawned keeps 6, a deleted worktree frees its offset, k2x2yk with a relative gitdir gets 8  4. live: from a repo with a .env, run spawn.sh for two slugs that share a hash, then: cat <repo>/.git/worktrees/*/jaira-offset - two different numbers"
---

# spawn.sh vergibt zwei laufenden Workern denselben Port-Offset

## Definition of Done

- [x] Zwei Slugs mit gleichem cksum-Rest bekommen verschiedene Offsets, solange beide Worktrees existieren (Test in core/role)
  proof: core/role/role_test.go TestSpawnGivesEveryWorktreeItsOwnOffset (8r8vts/32ttss → 6/7 beside twb0r6 on 5)
- [x] Ein Worktree, der schon steht, behält seinen Offset bei erneutem spawn
  proof: TestSpawnGivesEveryWorktreeItsOwnOffset: 8r8vts spawned again keeps 6; the offset is handed out only inside spawn.sh's 'if [ ! -d "$wt" ]'
- [x] Offset 0 bleibt frei für das Hauptverzeichnis; Bereich 1-40 bleibt, oder die Änderung steht in der USAGE und in einer Release-Note
  proof: spawn.sh: o=$(( (start + i) % 40 + 1 )) stays 1-40; TestSpawnGivesEveryWorktreeItsOwnOffset; USAGE + core/release/notes/DBJHNE.md say the new rule
- [x] Ein Lock älter als 120 s (mtime) wird übernommen statt zu warten und zu scheitern
  proof: spawn.sh:232 lock_age >= 120 → mv + take over; TestSpawnOffsetLock (stale lock 3 min old, 'took over', <20 s)
- [x] mkdir des Locks scheitert und das Lock-Verzeichnis existiert nicht → sofortiger Abbruch mit klarer Meldung
  proof: spawn.sh:221-230 'cannot make the lock' exit 1; TestSpawnOffsetLock (file at lock path, fails in <5 s, no worktree made)
- [x] Test: zwei spawn.sh parallel bekommen verschiedene Offsets
  proof: core/role/role_test.go TestSpawnOffsetLock: 8r8vts+32ttss in parallel with a 1 s post-checkout hook get different offsets; fails with the lock removed

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] failing test: two worktrees whose slugs share a cksum rest get different offsets
- [x] spawn.sh: allocate the offset under a lock before 'git worktree add', starting at the hash, skipping offsets the other worktrees hold
- [x] record the offset in the worktree git dir (jaira-offset); read legacy jaira-setup and .env comments too
- [x] all 40 taken: refuse before anything is made, only where ports are needed
- [x] USAGE line, release note, manifest regen, go test ./core/role

## Progress
- **2026-11-21 19:47 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2274544.sock
- **2026-11-21 19:48 · Alexander Sacharov** — decision: no open questions for the person. Offset is allocated when the worktree is made (already the only time spawn.sh computes it, so a standing worktree keeps its own — DoD 2 holds by construction). Start at the old cksum%40+1, take the next free one, wrapping, so a worktree with no collision gets the same number as before. Taken = offsets recorded by the other worktrees of the repo: new file jaira-offset in each worktree git dir (not jaira-setup: its presence is what makes reap.sh run a teardown), plus legacy jaira-setup offset= and the .env '# worker stack: ... offset=' line for worktrees made before this. A mkdir lock in the git common dir covers allocate+worktree add+record, so two spawns at once cannot pick the same free slot. Range stays 1-40, 0 stays the main directory's. All 40 taken: exit 1 before the worktree is made, but only when a setup or .env needs ports.
- **2026-11-21 19:51 · Alexander Sacharov** — critique: checked whether jaira-offset could be dropped for jaira-setup/.env (no: written after setup, outside the lock) and whether the lock is speculative (no: parallel spawns from a teamlead). Accepted: a SIGKILLed spawn leaves jaira-offset.lock and the next spawn stops after 10 s with a message naming it; a worktree dir deleted without 'git worktree prune' keeps its offset taken until pruned.
- **2026-11-21 19:54 · Alexander Sacharov** — bughunt: change: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:209-217 - class: the offset lock is held across 'git worktree add' (line 251), but a waiting spawn gives up after a fixed 10 s, so any checkout slower than that fails every spawn started beside it; example: a repository whose checkout takes 12 s (large tree, LFS, a post-checkout hook, Git Bash on Windows) - two spawns at once, the second exits 1 before making anything and tells the person to remove a lock that is still rightly held; ran: post-checkout hook 'sleep 12' in a scratch repo, two spawn.sh in parallel -> one rc=0, one rc=1 with 'jaira-offset.lock has stood for 10 seconds'
- **2026-11-21 19:54 · Alexander Sacharov** — bughunt: change: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:224-237 - class: every entry under <common>/worktrees/ counts as holding its offset, including prunable ones whose directory is gone (removed with rm -rf rather than 'git worktree remove'); their offsets stay blocked until 'git worktree prune', and with enough of them spawn.sh refuses with 'every port offset 1-40 is held' while far fewer worktrees exist, and its hint (reap.sh) does not free them; example: worktree repo-bbbbbb deleted by hand, still listed as prunable with jaira-offset 8 - slug eeeeee, which starts at 8, gets 9; ran: scratch repo + spawn.sh with a herdr stub, offsets read back from .git/worktrees/*/jaira-offset
- **2026-11-21 19:54 · Alexander Sacharov** — bughunt: done at d82c472b8 - 2 findings (2 change, 0 pre-existing, no blocker) in 4 files, first hunt; concurrency, error paths, input, permissions checked; files: core/release/notes/DBJHNE.md core/role/builtin/jaira-dispatcher/scripts/spawn.sh core/role/role_test.go core/shipped/shipped.sha256
- **2026-11-21 19:55 · Alexander Sacharov** — testing: class: the lock timeout is shorter than what it guards; example: post-checkout hook sleeping 12 s, second parallel spawn exits 1 and says to remove the lock the first one holds; fix: wait as long as a worktree add can reasonably take (120 s) and say in the message to remove it only when no spawn.sh runs. Rides along: skip worktree entries whose directory is gone (prunable) when counting taken offsets.
- **2026-11-21 19:57 · Alexander Sacharov** — bughunt: change: blocker: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:230-231 - class: the path read from $common/worktrees/<id>/gitdir is tested with [ -d ] against spawn.sh's own working directory, but git 2.48+ with worktree.useRelativePaths=true (or worktree add --relative-paths) stores it relative to that worktrees/<id> directory; every such live worktree is skipped as 'deleted', its jaira-offset is not counted, and the next spawn can be handed the same offset - two workers on one stack, the defect this ticket removes. Same relative path also feeds the .env fallback at :237. Before 231d2a909 a worktree with jaira-offset was counted without looking at gitdir. Fix direction: resolve a non-absolute gitdir from $gd (cd "$gd" && cd "$(dirname ...)"); role_test.go covers only absolute gitdir. example: gitdir '../../../../wt/x/.git' with jaira-offset 7, spawn.sh run from $HOME -> skipped; ran: scratch simulation of lines 229-231 (skipped from cwd, held from $gd); git here is 2.43, so git's own relative write was taken from the git-worktree docs (worktree.useRelativePaths, 2.48), not run
- **2026-11-21 19:57 · Alexander Sacharov** — bughunt: done at 231d2a909 - 1 finding (1 change, blocker; 0 pre-existing) in 4 files, follow-up hunt since d82c472b8 (diff in spawn.sh, role_test.go, shipped.sha256; DBJHNE.md unchanged); concurrency, error paths, input, permissions checked; files: core/release/notes/DBJHNE.md core/role/builtin/jaira-dispatcher/scripts/spawn.sh core/role/role_test.go core/shipped/shipped.sha256
- **2026-11-21 19:59 · Alexander Sacharov** — testing: round 3 is the third testing round; the lane sent work back twice (lock timeout + prunable worktrees; relative gitdir). A third send-back stops for the person.
- **2026-11-21 19:59 · Alexander Sacharov** — bughunt: done at d263a0319 - 0 findings in 4 files, follow-up hunt since 231d2a909 (diff in spawn.sh, role_test.go, shipped.sha256; DBJHNE.md unchanged); concurrency, error paths, input, permissions checked; ran: go test ./role -run 'Offset|Spawn' (pass); files: core/release/notes/DBJHNE.md core/role/builtin/jaira-dispatcher/scripts/spawn.sh core/role/role_test.go core/shipped/shipped.sha256
- **2026-11-21 20:02 · Alexander Sacharov** — Given back from signoff: (1) Lock älter als 120 s (mtime) übernehmen statt warten und scheitern — ein per SIGKILL beendeter spawn darf nicht jeden weiteren spawn blockieren, bis jemand von Hand aufräumt. (2) mkdir des Locks scheitert und das Verzeichnis existiert nicht → sofort mit klarer Meldung abbrechen, nicht 120 s warten. (3) Test: zwei spawn.sh parallel bekommen verschiedene Offsets.
- **2026-11-21 20:06 · Alexander Sacharov** — give-back round: stale = lock dir mtime >= 120 s via stat -c %Y / stat -f %m (find -mmin rounds differently on GNU and BSD: GNU -mmin +2 only matched at 180 s). Takeover renames the lock away first (atomic, only one wins); a renamer whose renamed dir turns out fresh puts it back if the place is empty. A live spawn whose checkout takes over 120 s now loses its lock - accepted, that is what the person asked for. Lock that cannot be made: one retry after a failed mkdir with no dir standing, then exit 1. Test setup moved into newOffsetRepo, shared by both offset tests; the parallel test was checked to fail with the lock disabled.
