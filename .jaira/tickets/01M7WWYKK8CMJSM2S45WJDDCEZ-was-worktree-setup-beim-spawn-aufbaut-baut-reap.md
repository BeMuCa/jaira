---
id: 01M7WWYKK8CMJSM2S45WJDDCEZ
title: "Was worktree-setup beim Spawn aufbaut, baut reap.sh beim Aufraeumen wieder ab"
status: testing
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
created-at: 2026-11-20T16:34:35Z
updated-at: 2026-11-21T19:06:29Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2280022
claimed-at: 2026-11-21T19:00:56Z
outcome-what: "reap.sh calls an executable .jaira/worktree-teardown <worktree> <slug> <offset> <main checkout> after closing a merged ticket's tabs and before git worktree remove; a failing, non-executable teardown or an untellable slug keeps the ticket with the reason; spawn.sh/reap.sh --help describe the pair; release note; manifest lines"
outcome-why: "worktree-setup starts container stacks, ports and volumes per ticket; reap.sh removed the worktree and left them running, so stacks piled up and collided on ports"
outcome-resolves: "Was worktree-setup beim Spawn aufbaut, baut reap.sh beim Aufraeumen wieder ab — DoD 1-3, TestReapTearsDownWhatTheSetupBuilt"
review-summary: none
review-gaps: "reap.sh: the untellable-slug case no longer borrows the 'dirty' variable — it prints its kept line and 'continue 2's; reap_test.go offsetOf: strconv.Atoi instead of a hand-written digit fold; shipped.sha256: the unreleased intermediate reap.sh line replaced, not appended. Left: the cksum offset formula exists in spawn.sh and reap.sh (accepted in critique note); no other teardown/cleanup hook exists in the repo (grep worktree-teardown, teardown)"
---

# Was worktree-setup beim Spawn aufbaut, baut reap.sh beim Aufraeumen wieder ab

## Definition of Done

- [x] reap.sh ruft ein ausfuehrbares .jaira/worktree-teardown mit denselben Argumenten wie das Setup auf, bevor es den Worktree entfernt (Test)
  proof: core/role/builtin/jaira-dispatcher/scripts/reap.sh:386 (call before 'git worktree remove'); core/role/reap_test.go TestReapTearsDownWhatTheSetupBuilt/succeeds checks args, worktree still present, tabs closed first
- [x] Scheitert das Teardown, behaelt reap.sh den Worktree und nennt den Grund ('kept <ticket>: teardown failed')
  proof: core/role/builtin/jaira-dispatcher/scripts/reap.sh:388 'kept <t>: teardown failed: ...', failed=1 keeps worktree/branch/records; TestReapTearsDownWhatTheSetupBuilt/fails (incl. retry next round), /not_executable, /slug_not_in_the_name
- [x] spawn.sh --help und reap.sh --help beschreiben das Paar setup/teardown
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh usage (worktree-teardown paragraph after the setup one); core/role/builtin/jaira-dispatcher/scripts/reap.sh usage (worktree-teardown paragraph before 'Output:')

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
