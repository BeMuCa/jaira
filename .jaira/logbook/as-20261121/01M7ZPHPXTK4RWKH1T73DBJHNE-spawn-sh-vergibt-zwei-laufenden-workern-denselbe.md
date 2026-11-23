---
id: 01M7ZPHPXTK4RWKH1T73DBJHNE
title: spawn.sh vergibt zwei laufenden Workern denselben Port-Offset
status: done
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
commits:
  - d82c472b83f660d6115616f2712d2734e0e9d043
  - 231d2a909df7b74f00519207b101538d6756e8c8
  - d263a03193c5e87829157deb56fc49bec0ce736f
  - 35d86c99b4b538264c62587c8a1c110377ea05d5
  - 5ba3fa8053a0a87b9c6be55cd15d9b082f423908
  - f280009a9c138b6f54f713da6d60f7f4160d4c10
  - ef8adc9bd699f92f6cdeaa5f3ba2e67cf6a2a8b2
created-at: 2026-11-21T18:40:24Z
updated-at: 2026-11-21T20:34:16Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-987290
claimed-at: 2026-11-21T20:03:04Z
outcome-what: "spawn.sh hands each new worktree a port offset no other existing worktree of the repository holds: starts at the old slug hash, skips taken ones, under a flock on .git/jaira-offset.lock (as reap.sh) that the kernel drops when spawn.sh dies, and writes it to jaira-offset in the worktree's git dir; an unopenable lock stops spawn.sh at once; worktrees whose directory is gone hold nothing, relative gitdir paths resolved against the worktree's git dir"
outcome-why: "cksum%40 alone gave 8r8vts, 32ttss and twb0r6 the same offset; two stacks on the same ports. The mkdir lock had no owner (takeover, put-back and rmdir acted on the path); flock removes stale locks by construction"
outcome-resolves: Zwei gleichzeitig laufende Worker bekommen nie denselben Port-Offset
review-summary: "spawn.sh now hands out a worktree's port offset (1-40) only once, when the worktree is made: it starts at the old cksum%40+1 and takes the next offset no other standing worktree of the repo holds (read from the new <gitdir>/jaira-offset, else legacy jaira-setup offset= or the .env '# worker stack' line; worktrees whose directory is gone hold nothing). Scanning, git worktree add and writing jaira-offset happen under flock on <common>/jaira-offset.lock via fd 9 (perl flock on macOS, unlocked with a warning when neither exists), waiting up to 600 s; git worktree add runs with 9>&- so a hook cannot keep the lock. All 40 taken, or an unopenable lock, stops the spawn before anything is made. Repos with neither setup nor .env get no offset and take no lock. USAGE, release note and manifest updated; two tests in core/role cover collision, keep-on-respawn, parallel spawns, killed spawn and unopenable lock."
review-gaps: |-
  note: the perl flock fallback (macOS) has no test - checked by hand that the lock stays held after perl exits while the shell keeps fd 9; a test with flock hidden from PATH would pin it
  note: spawn.sh ~line 206 - one comment line runs past 150 columns ("or two spawns at once ... held through fd 9 as reap.sh holds its"); rewrap
  note: known limit: a worktree made by hand without jaira-offset, jaira-setup or a marked .env holds no offset, so a spawn can give its ports to a new worker
  note: pre-existing: two spawns of the same slug at once both pass [ ! -d "$wt" ]; the second now waits for the lock and then fails at git worktree add (branch exists) - harmless, but the message is git's, not spawn.sh's
test-verdict: "pass: go test -race -count=1 ./core/role ./core/shipped ./core/lane RC=0 (files changed since 35d86c99b); DoD 1-7 verified by TestSpawnGivesEveryWorktreeItsOwnOffset and TestSpawnOffsetLock - the lock test fails with flock disabled (both got 5) and without 9>&- on worktree add (waited 5 s); live scratch repo: spawn SIGKILLed inside a 20 s checkout hook, next spawn done in 0.07 s rc=0; three spawns at once with a 1 s hook got 5, 6, 7; read-only .git without a lock file exits 1 at once with 'cannot open the lock', no worktree made"
hunt-verdict: "pass: bughunt done at 5ba3fa805, 0 change, 0 pre-existing; flock on drvfs and the perl fallback checked live; Git Bash on real Windows and macOS perl not run (only WSL here)"
round: 5
back-from: human
back-by: person
back-count: 1
review-verdict: "All 7 DoD items are met by the code and the tests (go test ./core/role, ./core/shipped pass). The lock is correct: flock lives on the open file description held by fd 9, closed for git worktree add and released before setup runs. No blocker; four small notes, none a wrong result."
review-round: 5
review-verdict-gaps: |-
  note: the perl flock fallback (macOS) has no test - checked by hand that the lock stays held after perl exits while the shell keeps fd 9; a test with flock hidden from PATH would pin it
  note: spawn.sh ~line 206 - one comment line runs past 150 columns ("or two spawns at once ... held through fd 9 as reap.sh holds its"); rewrap
  note: known limit: a worktree made by hand without jaira-offset, jaira-setup or a marked .env holds no offset, so a spawn can give its ports to a new worker
  note: pre-existing: two spawns of the same slug at once both pass [ ! -d "$wt" ]; the second now waits for the lock and then fails at git worktree add (branch exists) - harmless, but the message is git's, not spawn.sh's
review-check: "1. cd /home/alex/projects/.worktrees/jaira-dbjhne  2. go test ./core/role -run 'TestSpawnGivesEveryWorktreeItsOwnOffset|TestSpawnOffsetLock' -count=1 -v  3. you see both PASS: slugs 8r8vts/32ttss get 6/7 beside twb0r6 on 5, a respawn keeps 6, parallel spawns get different offsets, a killed spawn does not block the next, a directory at the lock path fails in under 5 s  4. optional live check: in any repo with a .env, run spawn.sh for two tickets at once and compare 'offset=' in the two worktrees' .env files - they differ"
question: ""
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
  proof: superseded by decision 20:13 (variant 1): no takeover needed - flock dies with the holder; TestSpawnOffsetLock 'A spawn.sh killed inside the lock' (next spawn in <3 s while the killed one's hook still runs)
- [x] mkdir des Locks scheitert und das Lock-Verzeichnis existiert nicht → sofortiger Abbruch mit klarer Meldung
  proof: spawn.sh:214 exec 9>> fails → 'cannot open the lock', exit 1; TestSpawnOffsetLock (directory at the lock path, <5 s, no worktree)
- [x] Test: zwei spawn.sh parallel bekommen verschiedene Offsets
  proof: TestSpawnOffsetLock: 8r8vts+32ttss in parallel with a 1 s post-checkout hook get different offsets; fails with flock disabled
- [x] Offset-Lock ist flock über fd wie reap.sh; mkdir+mtime-Übernahme entfernt; ein getöteter spawn hinterlässt keinen Lock
  proof: spawn.sh:213-231 flock -w 600 9 / perl flock on fd 9, mkdir+mtime code gone; git worktree add runs with 9>&-; TestSpawnOffsetLock fails without 9>&- (waited 5 s)

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
- **2026-11-21 20:07 · Alexander Sacharov** — critique (give-back round): checked reap.sh's flock-on-fd-9 pattern (core/role/builtin/jaira-dispatcher/scripts/reap.sh:105-118) as a simpler shape - the kernel drops a flock when the holder dies, so no stale lock and no unmakeable-dir wait. Not taken: the give-back prescribes mtime takeover explicitly, reap.sh runs unlocked where neither flock(1) nor perl exists (unproven on Git Bash, and for spawn.sh unlocked means the very race this ticket removes), and spawn.sh would have to close fd 9 before every child (git hooks, herdr, claude) or a child keeps the lock. Left standing per the 19:51 critique note: a worktree dir deleted without prune; per the 20:02 give-back: a live spawn whose checkout takes >120 s loses its lock.
- **2026-11-21 20:12 · Alexander Sacharov** — bughunt: change: blocker: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:232-247 and :296 - class: the offset lock has no owner. Takeover, put-back and release act on the path $lock without checking whose lock it is, so two spawns end up inside the lock at once and can both take the same free offset (the hash start collides, e.g. 8r8vts/32ttss). Three places: (1) :232 takes over any lock older than 120 s, including one whose holder is still inside a slow 'git worktree add' (the case the comment at :208 names), and the holder's plain 'rmdir $lock' at :296 (also the EXIT trap at :257) then deletes the newcomer's lock and lets a third spawn in. (2) :236-244: a spawn whose rename caught a fresh lock leaves the path empty, so a third spawn can mkdir in that gap. Then '[ ! -e $lock ]' fails, the renamed lock is rmdir'd, and its owner and the third spawn both hold. (3) :243 '[ ! -e ] && mv' is check-then-act: a lock made in between gets the stale directory moved INSIDE it, so the holder's rmdir fails and every waiting spawn stalls about 120 s, after which a .stale.<pid> directory is left in .git. Example: a holder still working when its lock passes 120 s -> holder, newcomer and a third spawn overlap. ran: harness of the exact lines 214-257 plus a release by rmdir; lock aged with touch -d -130s while A held it -> B 'took over' and entered while A was inside, and C entered while B was inside. 10 parallel spawns on a 200-s-old lock, 25 rounds -> 2 bad rounds (max 2 and max 3 holders at once; in one, the lock contained jaira-offset.lock.stale.<pid> and 7/10 spawns hung). blame: d82c472b8 (release) and 35d86c99b (takeover), both DBJHNE
- **2026-11-21 20:12 · Alexander Sacharov** — bughunt: done at 35d86c99b - 1 finding (1 change, blocker; 0 pre-existing) in 4 files, follow-up hunt since d263a0319 (diff in spawn.sh, role_test.go, shipped.sha256, DBJHNE.md); concurrency, error paths, input, permissions checked; ran: go test ./role -run 'Offset|Spawn' (pass) plus a lock-race harness; files: core/release/notes/DBJHNE.md core/role/builtin/jaira-dispatcher/scripts/spawn.sh core/role/role_test.go core/shipped/shipped.sha256
- **2026-11-21 20:13 · Alexander Sacharov** — testing: third send-back of testing on the lock rule - stopped for the person. question: Третий возврат testing по одному правилу (lock). Bughunt: у mkdir-lock нет владельца — перехват по mtime, возврат на место и rmdir держателя действуют по пути, 2-3 spawn оказываются внутри одновременно (стенд: 2 из 25 прогонов). Варианты: 1) (Recommended) заменить mkdir+mtime на flock через fd, как reap.sh:105-118 — ядро снимает lock при смерти процесса, перехват и проверка возраста не нужны, нельзя открыть файл lock → сразу ошибка; цена: fd закрыть перед дочерними процессами, без flock(1) и perl — без lock с предупреждением. 2) оставить mkdir, добавить владельца (pid в lock, rmdir только своего, перехват только при мёртвом pid) — ещё одна эвристика, новые края. 3) принять как есть, гонку записать как known limit.
- **2026-11-21 20:13 · Alexander Sacharov** — decision: Variante 1 (Alex' Regel 'Ränder wachsen → kleineres Design / jaira-Code wiederverwenden', 21.11.): mkdir+mtime-Lock raus, flock über fd wie reap.sh:105-118. Kein Übernehmen nach mtime mehr — der Kernel gibt den Lock beim Tod frei. Lock-Datei nicht zu öffnen → sofort exit 1. Ohne flock(1)/perl: Warnung, ohne Lock, wie reap.sh.
- **2026-11-21 20:16 · Alexander Sacharov** — variant 1 done: flock on fd 9 (reap.sh pattern), blocking with 600 s cap (flock -w 600, perl alarm 600) - a live slow checkout is waited for, never taken over. 'git worktree add' runs with 9>&- or a checkout hook that keeps running holds the lock (shown: test fails with it removed). The lock file now stays in .git (empty file, harmless; reap.sh's reap.lock does the same). exec 9>> failing (dir at the path, read-only .git) → exit 1 at once. Neither flock nor perl → warning, unlocked, as reap.sh.
- **2026-11-21 20:16 · Alexander Sacharov** — critique (flock round): same pattern as reap.sh:105-118 (exec 9>>, flock(1) else perl, neither → unlocked); differs on purpose in two places: blocking with -w 600 instead of -n (a spawn must get its offset, not skip), and only 'git worktree add' runs with 9>&- instead of a whole block (everything else between lock and release is cat/sed/rev-parse, which end before the release). Nothing simpler left.
- **2026-11-21 20:18 · Alexander Sacharov** — bughunt: done at 5ba3fa805 - 0 findings in 4 files, follow-up hunt since 35d86c99b (diff in spawn.sh, role_test.go, shipped.sha256, DBJHNE.md; no new file); concurrency, error paths, input, permissions checked; ran: go test ./role -run 'Offset|Spawn' -count=3 (pass), flock exclusion+release on drvfs /mnt/c, perl flock fallback holds past perl's exit (waiter waited 2 s for the holder shell), 9>&- on an unopened fd under set -e (exit 0); files: core/release/notes/DBJHNE.md core/role/builtin/jaira-dispatcher/scripts/spawn.sh core/role/role_test.go core/shipped/shipped.sha256
