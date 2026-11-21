---
id: 01M7WWYKK8CMJSM2S45WJDDCEZ
title: "Was worktree-setup beim Spawn aufbaut, baut reap.sh beim Aufraeumen wieder ab"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Hat ein Repository .jaira/worktree-setup, kann es ein .jaira/worktree-teardown haben; reap.sh ruft es vor git worktree remove auf (z.B. Docker-Stack des Tickets stoppen und Volumes loeschen)."
context: |-
  Heute: spawn.sh ruft .jaira/worktree-setup <worktree> <slug> <port-offset> <repo-root> auf (~/.claude/skills/jaira-dispatcher/scripts/spawn.sh:47,170).
  reap.sh entfernt Worktree und Branch, kennt aber nichts, was das Setup gestartet hat: Container, Ports, Volumes laufen weiter.
  In einem Projekt liefen zwei Stacks auf demselben Port, einmal mit Ports und Projektname des Hauptstacks. Freie Ports pruefen ist Sache des Setup-Skripts im Projekt; das Aufraeumen braucht aber einen Haken in jaira.
  Alex 20.11.2026: "ok, aber man muss es auch schliessen".
definition-of-done: "reap.sh ruft ein ausfuehrbares .jaira/worktree-teardown mit denselben Argumenten wie das Setup auf, bevor es den Worktree entfernt (Test)"
tags: []
blocked-by:
  - 01M7R61GC6C8D7547N7CN7JW2Q
related:
  - 01M2MFEEE9VZAABPKKS27MG5GB
  - 01M6M72EN4ZZ2JBZJ1Q6DSDPZ4
commits:
  - 8775720a05dc249662c681c742e07e718cd0e5af
  - b3a70827c6941be031e5b4ddb70a1a592d94c9f2 4dc9c55d2530419aaf9ac749d3d92a17bc9721a2 8775720a05dc249662c681c742e07e718cd0e5af
created-at: 2026-11-20T16:34:35Z
updated-at: 2026-11-21T19:26:31Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-3230212
claimed-at: 2026-11-21T19:17:08Z
outcome-what: "spawn.sh records the setup's four arguments in the worktree's git dir (jaira-setup); reap.sh calls an executable .jaira/worktree-teardown with exactly those after closing a merged ticket's tabs and before git worktree remove; no record -> no call; failing or non-executable teardown keeps the ticket with the reason; --help of both describe the pair; release note; manifest"
outcome-why: "worktree-setup starts container stacks, ports and volumes per ticket; reap.sh removed the worktree and left them running, so stacks piled up and collided on ports. Round 2: deriving the slug from the worktree name was wrong for worktrees spawned from another worktree (bughunt)"
outcome-resolves: "Was worktree-setup beim Spawn aufbaut, baut reap.sh beim Aufraeumen wieder ab — DoD 1-3"
review-summary: "none"
review-gaps: "none: no other timeout(1)/gtimeout lookup in the shipped scripts to reuse; nothing orphaned by 46ea1b53; the per-ticket 'command -v timeout' runs only when a teardown exists, left in place"
test-verdict: "pass: go test ./... green (exit 0), go test -race -count=1 ./core/role green, go vet and shipped-history --check green; DoD 1-5 verified in the tree (TestReapTearsDownWhatTheSetupBuilt incl. /hangs: 1s bound, sleep child gone); function: a TERM-ignoring teardown with a background child is killed by timeout -k after 3s (exit 137, no child left); JAIRA_TEARDOWN_TIMEOUT=abc exits 2 with the reason; writers: the change adds no write - the teardown runs inside reap's existing flock"
hunt-verdict: "fail: bughunt done at 46ea1b53 - 1 change: reap.sh:361 asks for timeout(1) before it knows a teardown will run, so a worktree without a jaira-setup record is kept for good on a machine without coreutils"
round: 4
back-from: testing
back-by: lane
back-count: 1
review-verdict: "All three DoD items are met and proven by the table test (go test ./role -run 'Reap|Spawn' and ./shipped pass). The record replaces the earlier name-derived slug, so the wrong-stack teardown the first bughunt found is gone. No blocker; four small edges as notes, the timeout one the most worth closing."
review-round: 2
review-verdict-gaps: |-
  note: reap.sh:379 a teardown that hangs (docker compose down waiting on a stuck container) is called without a timeout, so reap.sh never ends and holds reap.lock - every later round prints 'skipped: reap.sh runs already' and nothing is reaped for any ticket; fix: wrap the call in timeout(1) where it exists and treat expiry as a failed teardown
  note: reap.sh:340 the teardown is always the main checkout's file, while spawn.sh:178 ran $root/.jaira/worktree-setup - for a worktree spawned from another worktree the pair can differ (setup there, teardown not yet on master -> silently no call); fix: record the setup's path in jaira-setup and call the teardown beside it, or say it in --help next to the existing sentence
  note: reap.sh:347 recorded paths are passed through unchecked (bughunt 19:13 still open): root= can name a worktree already reaped; fix: say in reap.sh --help and the release note that $1/$4 may no longer exist, or hand the live $w as $1
  note: the teardown runs in whatever directory reap.sh was started from; fix: cd to the worktree (or say in --help that it must use absolute paths)
review-check: |-
  1. cd core && go test ./role -run 'TestReapTearsDownWhatTheSetupBuilt|TestSpawnLetsTheRepositorySetUpItsWorktree' -v  -> all subtests PASS (none, succeeds, fails, not_executable, no_record, setup)
  2. bash core/role/builtin/jaira-dispatcher/scripts/reap.sh --help  -> a paragraph about .jaira/worktree-teardown before 'Output:', and 'tore down <worktree>' in the output list
  3. bash core/role/builtin/jaira-dispatcher/scripts/spawn.sh --help | grep -A4 'counterpart'  -> names .jaira/worktree-teardown and the jaira-setup record
  4. optional live: in a repo with an executable .jaira/worktree-setup and a .jaira/worktree-teardown that echoes its args, spawn a ticket, then cat "$(git -C <worktree> rev-parse --absolute-git-dir)/jaira-setup"  -> worktree=, slug=, offset=, root= lines; after its PR is merged and the ticket is done, reap.sh prints 'tore down <worktree>' then 'removed <worktree>'
---

# Was worktree-setup beim Spawn aufbaut, baut reap.sh beim Aufraeumen wieder ab

## Definition of Done

- [x] reap.sh ruft ein ausfuehrbares .jaira/worktree-teardown mit denselben Argumenten wie das Setup auf, bevor es den Worktree entfernt (Test)
  proof: core/role/builtin/jaira-dispatcher/scripts/reap.sh:408 calls the teardown with the args from the jaira-setup record core/role/builtin/jaira-dispatcher/scripts/spawn.sh:217 writes; TestReapTearsDownWhatTheSetupBuilt/succeeds (recorded slug bbb, not the name's ab12cd; worktree present, tabs closed first) and TestSpawnLetsTheRepositorySetUpItsWorktree/setup (record = setup args)
- [x] Scheitert das Teardown, behaelt reap.sh den Worktree und nennt den Grund ('kept <ticket>: teardown failed')
  proof: core/role/builtin/jaira-dispatcher/scripts/reap.sh:412 'kept <t>: teardown failed: ...', failed=1 keeps worktree/branch/records; TestReapTearsDownWhatTheSetupBuilt/fails (incl. retry next round), /hangs, /not_executable
- [x] spawn.sh --help und reap.sh --help beschreiben das Paar setup/teardown
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh usage (worktree-teardown paragraph after the setup one); core/role/builtin/jaira-dispatcher/scripts/reap.sh usage (worktree-teardown paragraph before 'Output:')
- [x] reap.sh ruft das Teardown mit timeout(1) auf; Ablauf zaehlt als fehlgeschlagenes Teardown (kept), Kindprozesse gehen mit (Test)
  proof: core/role/builtin/jaira-dispatcher/scripts/reap.sh:408 timeout -k 10 $JAIRA_TEARDOWN_TIMEOUT (default 300), 124/137 -> 'kept <t>: teardown failed: ... (timed out after Ns)'; no timeout/gtimeout -> kept only where a teardown will run (core/role/builtin/jaira-dispatcher/scripts/reap.sh:371); TestReapTearsDownWhatTheSetupBuilt/hangs (sleep 30 child killed after 1s), /no_timeout, /no_timeout_no_record
- [x] reap.sh --help und Release-Note sagen: Teardown kommt aus dem main checkout, $1/$4 koennen schon fehlen, absolute Pfade oder cd "$1" benutzen
  proof: core/role/builtin/jaira-dispatcher/scripts/reap.sh usage paragraph ending ':61 never rely on $4' + timeout paragraph; core/release/notes/JDDCEZ.md lines 4-5 (line 5 marked '- !')

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-20 16:50 · Alexander Sacharov** — Interim-Regeln (Alex 20.11.2026, gelten ab sofort bis TW5FE2 ausgeliefert ist):
1. pre-process: die DoD nennt die Fehlerklasse statt des Symptoms, die Eingabeklassen (leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg) und bei einer Heuristik eine Annahmeregel - als DoD-Punkte (jaira dod <id> --add).
2. in-progress: vor dem Move selbst die vier Kategorien pruefen (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte). Tests ueber die Klasse (Tabellentest), Geschwister-Stellen per grep suchen und mitziehen, Fix an der Ursache statt am Beispiel. Nicht nur Haken setzen: wer umsetzt, denkt die Raender selbst.
3. bughunt/testing: ein Fund nennt Klasse und Beispiel; der Rueckweg verlangt den Fix der Klasse, nicht des einen Beispiels.
4. Schleifenbremse: Blocker sind nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke. Nach einem Fix-Durchgang laeuft der Testlauf und ein Bughunt nur ueber den Diff seit dem letzten Hunt, dann review. Nach 3 Rueckgaben aus derselben Lane: Frage an den Teamlead statt weiter schleifen.
5. Kleine Raender aus review-gaps in derselben Runde schliessen statt als note: stehen lassen - Alex schickt sie sonst ohnehin zurueck. Vorbestehendes ausserhalb des Umfangs als 'note: pre-existing: ...'.
- **2026-11-21 18:58 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/4205.sock
- **2026-11-21 18:58 · Alexander Sacharov** — decision: no open decisions for the person. Teardown args = setup args, derived in reap.sh without changing the pane record: slug = worktree basename minus '<basename main>-' (spawn.sh names it so; exact, so the cksum offset matches); offset = same cksum formula as spawn.sh; repo-root = main checkout (where spawn.sh read .jaira/worktree-setup in the normal case); the teardown file is read from $main/.jaira/worktree-teardown. Runs after the tabs closed (no claude still on the stack), before 'git worktree remove'. Present but not executable: kept before any tab closes. Failure: records stay, so the next round runs it again -> teardown must be idempotent, said in --help. Output to stderr: stdout is reap's line protocol.
- **2026-11-21 19:01 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2274544.sock (claim taken over: previous holder process 2259282 gone, its parent socket 4205 gone)
- **2026-11-21 19:04 · Alexander Sacharov** — found out: CI gate shipped-manifest — every change to a shipped script needs its new sha256 line in core/shipped/shipped.sha256 (TestEveryInstalledFileIsAShippedVersion names it). Four categories checked: concurrency — teardown runs inside reap's flock, one ticket after another; error paths — non-zero exit keeps everything and is retried, so idempotency is demanded in --help (also for a worktree whose setup never ran); input — all paths quoted, empty slug (worktree named exactly '<main>-') refused; permissions — the teardown file comes from the main checkout, same trust as worktree-setup. Teardown's stdin is /dev/null, stdout to stderr (stdout is reap's line protocol). Not covered by a test: reap run from a worktree other than main (repo-root is still $main; equivalent in the rig).
- **2026-11-21 19:05 · Alexander Sacharov** — critique: none. Considered recording slug/offset in the pane record (spawn.sh) instead of deriving them in reap.sh — rejected per the 18:58 decision note: old records lack them and a --no-worktree record carries a slug unrelated to the worktree; the name-derived slug is exact for every worktree spawn.sh makes, and an unmatched name is kept rather than guessed. The second copy of the cksum offset formula (reap.sh vs spawn.sh) is the price; both carry a comment naming the other.
- **2026-11-21 19:08 · Alexander Sacharov** — bughunt: change: blocker: core/role/builtin/jaira-dispatcher/scripts/reap.sh:345 - class: reap.sh rebuilds the slug and offset it hands .jaira/worktree-teardown from the worktree's directory name minus '<physical basename of the main checkout>-', but spawn.sh:169 names the worktree after basename of whatever root it was given (logical path, possibly another worktree) and hands the setup the real slug; wherever those two spellings differ, the teardown gets other arguments than the setup did. example 1: spawn.sh run from a ticket worktree .worktrees/app-aaa with slug bbb (a teamlead or dispatcher sitting in a worktree) creates .worktrees/.worktrees/app-aaa-bbb, setup gets slug bbb offset 9; reap.sh calls the teardown with slug aaa-bbb offset 31, prints 'tore down', removes the worktree and its records - the real stack, ports and volumes stay up and are never torn down again (and a hyphenated free-text slug like aaa-bbb could name another ticket's stack). example 2: main checkout reached through a symlinked leaf (~/code/app -> /data/app-v2): spawn names .worktrees/app-ccc, reap strips 'app-v2-' and keeps the ticket every round ('cannot tell the slug'), a worktree that reaped fine before this change. Same mismatch in the 4th argument: setup got spawn's root, teardown gets main. ran: scratch repo reproducing spawn.sh:166-169 and reap.sh:118,345,351 derivations (reap-slug=aaa-bbb offset=31 vs spawn bbb offset 9; symlink wt=.worktrees/link-ccc vs main basename app); blame: lines 339-351 from 8775720a/4dc9c55d (JDDCEZ)
- **2026-11-21 19:08 · Alexander Sacharov** — bughunt: done at 4dc9c55d - 1 finding (1 change, 0 pre-existing) in 3 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/release/notes/JDDCEZ.md core/role/builtin/jaira-dispatcher/scripts/reap.sh core/role/builtin/jaira-dispatcher/scripts/spawn.sh core/role/reap_test.go core/shipped/shipped.sha256
- **2026-11-21 19:08 · Alexander Sacharov** — decision (worker, overturns the 18:58 note after the bughunt proof): the class is 'setup arguments reconstructed after the fact'. Fix at the cause: spawn.sh writes the arguments it gave the setup into the worktree's own git dir (git rev-parse --git-dir in the worktree -> <common>/worktrees/<name>/jaira-setup, gone with 'git worktree remove'); reap.sh calls the teardown only with those. No record (setup never ran, or a worktree from before this change) -> no teardown call, no guess. Name derivation removed entirely.
- **2026-11-21 19:12 · Alexander Sacharov** — critique round 2: read only the bughunt finding and its answer (b3a70827). The record replaces the name derivation and with it the second copy of the cksum formula (my round-1 note's accepted trade-off is gone). No finding.
- **2026-11-21 19:13 · Alexander Sacharov** — bughunt: change: core/role/builtin/jaira-dispatcher/scripts/reap.sh:347 - class: paths taken from the jaira-setup record are handed to the teardown without checking they still exist, while reap.sh already holds the live worktree in $w; example: ticket bbb spawned from worktree jaira-aaa records root=.../jaira-aaa, aaa is merged and reaped first, then bbb's teardown gets a root that is gone - a teardown that reads its compose file from $4 (or cd's into $1 after a 'git worktree move') fails every round and keeps bbb for good, or one that ignores the error leaks the stack silently; the reap.sh --help text and the release note promise 'exactly as spawn.sh gave them' without saying they may no longer exist; not a blocker; ran: read (spawn.sh:216-218 writes $wt/$root, reap.sh:345-347 passes them through, reap.sh:264-270 canonicalises only $w)
- **2026-11-21 19:13 · Alexander Sacharov** — bughunt: done at b3a70827 - 1 finding (1 change, 0 pre-existing), follow-up hunt since 4dc9c55d; diffs read in 5 modules, role_test.go read as its new hunk only; concurrency, error paths, input, permissions checked; files: core/release/notes/JDDCEZ.md core/role/builtin/jaira-dispatcher/scripts/reap.sh core/role/builtin/jaira-dispatcher/scripts/spawn.sh core/role/reap_test.go core/role/role_test.go core/shipped/shipped.sha256
- **2026-11-21 19:16 · Alexander Sacharov** — Given back from signoff: (1) Teardown mit timeout(1) aufrufen, Ablauf = fehlgeschlagener Teardown (kept) — seit 0.5.15 ruft das Cockpit reap.sh alle 10 min, ein hängender Teardown hielte reap.lock für immer und nichts würde mehr aufgeräumt. (2)-(4) nur in reap.sh --help und Release-Note sagen: Teardown kommt aus dem main checkout, $1/$4 können schon fehlen, absolute Pfade benutzen (oder cd in den Worktree). Kein neues Verhalten dafür.
- **2026-11-21 19:17 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2274544.sock (claim taken over again: holder process 2280022 gone; given back from signoff 19:16, working the 4 points)
- **2026-11-21 19:20 · Alexander Sacharov** — give-back worked: timeout(1) with -k 10 (kills the teardown's process group, so docker compose children go too); 124 and 137 both read as timed out. JAIRA_TEARDOWN_TIMEOUT seconds, default 300 (half the cockpit's 10-min reap interval), validated at start (exit 2 on non-number). macOS without coreutils: no timeout -> ticket kept with 'install coreutils' rather than an unbounded call - not covered by a test (cannot hide /usr/bin/timeout in the rig). shipped.sha256: replaced the branch's unreleased reap.sh line instead of appending (the intermediate version never shipped; shipped-history.sh --check green). Gates run: go test ./..., go vet, go test -race core/role, shipped-history --check.
- **2026-11-21 19:21 · Alexander Sacharov** — critique (give-back round): read only 46ea1b53. none. JAIRA_TEARDOWN_TIMEOUT follows run-lane.sh's JAIRA_LANE_TIMEOUT pattern and is what the test needs to run in 1s; 'kept' without timeout(1) rather than an unbounded call, because an unbounded call is exactly what the give-back forbids.
- **2026-11-21 19:23 · Alexander Sacharov** — bughunt: change: core/role/builtin/jaira-dispatcher/scripts/reap.sh:361 - class: the no-timeout(1) guard sits before the jaira-setup record loop, so it keeps a ticket even when no worktree of it has a record and no teardown would be called; this contradicts the release note that pre-record worktrees get no teardown call and are reaped as before; example: on macOS without coreutils, a merged ticket whose worktree was spawned before this release (no jaira-setup) is kept every round with 'kept <t>: no timeout(1) to bound ...' and never reaped; --dry-run prints the same 'kept' instead of 'would remove'; fix direction: check for tbin only when targs is non-empty; ran: PATH=<PATH without timeout/gtimeout> go test ./core/role -run TestReapTearsDownWhatTheSetupBuilt/no_record (FAIL: worktree gone false, 'kept AB12CD: no timeout(1) ...'); passes with timeout in PATH; blame 46ea1b53 names JDDCEZ; not a blocker: keeps, loses nothing
- **2026-11-21 19:24 · Alexander Sacharov** — bughunt: done at 46ea1b53 - 1 finding (1 change, 0 pre-existing), follow-up hunt since b3a70827; diffs read in 4 modules (reap.sh, reap_test.go, release note, shipped.sha256), spawn.sh and role_test.go unchanged since; concurrency, error paths, input, permissions checked; files: core/release/notes/JDDCEZ.md core/role/builtin/jaira-dispatcher/scripts/reap.sh core/role/builtin/jaira-dispatcher/scripts/spawn.sh core/role/reap_test.go core/role/role_test.go core/shipped/shipped.sha256
- **2026-11-21 19:26 · Alexander Sacharov** — bughunt 46ea1b53 finding fixed at the cause: the timeout(1) lookup now runs after the record scan and only when targs is non-empty, so a worktree without jaira-setup reaps as before on a machine without coreutils. Tests /no_timeout and /no_timeout_no_record run reap.sh with a PATH of links minus timeout/gtimeout (pathWithout in reap_test.go); /no_timeout_no_record fails against 46ea1b53's reap.sh. Not acted on: busybox timeout without -k (hunter could not verify; an old busybox fails the call with exit 1 -> ticket kept, visible, not silent).
