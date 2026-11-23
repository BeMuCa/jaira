---
id: 01M851YE2ACMD981Z6V1Q030WZ
title: Unter JAIRA_USER zaehlt jaira die per E-Mail zugewiesenen Tickets des Checkout-Besitzers als die der anderen Person
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ist JAIRA_USER gesetzt, gilt nur dessen Name samt seinen Aliassen als 'ich' - nicht zusaetzlich die git user.email des Checkouts."
context: |-
  Was falsch ist: core/identity/identity.go:112 (func Aliases) nimmt git user.email immer dazu, auch wenn JAIRA_USER eine andere Person nennt.
  Folge: unter JAIRA_USER=grace zaehlt IsMe die Tickets, die per E-Mail dem Besitzer des Checkouts gehoeren, als grace' Tickets (list --mine, Claim-Pruefungen).
  Gefunden im Review zu DVPAPE am 23.11.2026; gab es schon vorher auf master.
  Erst nach DVPAPE starten: DVPAPE baut core/identity gerade um (Personenliste auf refs/jaira/people), beide Tickets aendern dieselbe Datei.
definition-of-done: "Unter JAIRA_USER liefert Aliases nur JAIRA_USER und dessen Aliasse aus der Personenliste, nicht die git user.email des Checkouts; Test"
tags:
  - cli
blocked-by: []
related: []
follows: 01M67K6X23DJ7TDVQ6DVDVPAPE
commits:
  - aee7831848ef4f20b491752beec50b31c59ab895
  - 5f2c28b16a8537611dca72098cbcdd6b390a4080
  - 6af1ccd240c4a954297997830ef5e987dc3d759f
  - fb027aa0772b78b4e4079563a576263e68bb1a4e
  - 8863ba7af95e3d756c32293ff35d72fb66cd738b
created-at: 2026-11-23T20:35:47Z
updated-at: 2026-11-23T21:33:14Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-37226
claimed-at: 2026-11-23T21:01:32Z
outcome-what: Release note names the real ownership checks instead of the nonexistent list --mine; whoami help points to jaira people alias under JAIRA_USER and --json gains alias_file_read
outcome-why: "testing found the note named a flag jaira does not have; the bughunt found whoami still presented the alias file as in effect"
outcome-resolves: "under JAIRA_USER the checkout owner's git email and alias file no longer make the owner's tickets the other person's"
review-summary: "Full review from aee783184 (DVPAPE merge) to fb027aa07. identity.Aliases now returns right after the JAIRA_USER name and its people-list names when JAIRA_USER is set, so git user.email and ~/.jaira/identity no longer count as 'me'; the fix sits in the callee, so every caller (claim/move gate via identityAliases, release and pull via IsMe, TUI own-ticket view and signoff) picks it up. A new identity.ByEnv() replaces the inline JAIRA_USER checks. whoami says the alias file is not read under JAIRA_USER (text marker, --json alias_file_read) and its help points to 'jaira people alias'. Test TestAliasesUnderJairaUserLeaveOutTheOwnersNames covers owner email, owner git name, alias-file entry (not me) and the people-list email (me)."
review-gaps: "note: limit: someone who sets JAIRA_USER to their own name without being in the people list loses their git email and alias-file matches too; on purpose, the release note tells them to use 'jaira people alias'"
test-verdict: "pass: round 2 — vet + go test -race -count=1 internal/cli green (RC=0), round 1 full suite -race green; release note's commands all run (claim, move, release, pull, people alias); whoami text and --json alias_file_read agree with and without JAIRA_USER; DoD 1 re-verified: grace refused move exit 3, owner exit 0, two parallel sets both land"
hunt-verdict: "pass: bughunt done at fb027aa07, 0 change, 0 pre-existing (round 1's whoami edge fixed in fb027aa07)"
round: 2
back-from: testing
back-by: lane
back-count: 1
review-verdict: "accept: both DoD items met; Aliases under JAIRA_USER returns only that name and its people-list names, tested; no defect found in core/identity or whoami along concurrency, error paths, input or permissions"
review-round: 2
review-verdict-gaps: "note: limit: someone who sets JAIRA_USER to their own name without being in the people list loses their git email and alias-file matches too; on purpose, the release note tells them to use 'jaira people alias'"
review-check: "1. go build -o /tmp/jaira ./cmd/jaira  2. in a scratch repo: git init; git config user.name Owner; git config user.email owner@x.test  3. /tmp/jaira whoami - shows 'Acting as: Owner' and 'Also me: owner@x.test'  4. JAIRA_USER=grace /tmp/jaira whoami - shows only 'Acting as: grace', no 'Also me' line with owner@x.test, and the alias file line ends with '(not read under JAIRA_USER)'  5. JAIRA_USER=grace /tmp/jaira whoami --json - aliases is [\"grace\"], alias_file_read is false  6. go test ./core/identity -run TestAliasesUnderJairaUser - ok"
---

# Unter JAIRA_USER zaehlt jaira die per E-Mail zugewiesenen Tickets des Checkout-Besitzers als die der anderen Person

## Definition of Done

- [x] Unter JAIRA_USER liefert Aliases nur JAIRA_USER und dessen Aliasse aus der Personenliste, nicht die git user.email des Checkouts; Test
  proof: core/identity/identity.go Aliases (returns before email/alias file when ByEnv()); TestAliasesUnderJairaUserLeaveOutTheOwnersNames
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/Q030WZ.md (every command it names runs: claim, move, release, pull, people alias)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 21:00 · Alexander Sacharov** — teamlead 23.11.2026: auf Alex' Wunsch sofort gestartet, gestapelt auf feat/dvpape (DVPAPE noch nicht auf master; beide aendern core/identity). Gemergt wird nach DVPAPE. Aliases soll unter JAIRA_USER die Personenliste aus DVPAPE (core/people) nutzen.
- **2026-11-23 21:01 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/48662.sock
- **2026-11-23 21:02 · Alexander Sacharov** — decision (worker): open decisions none. Under JAIRA_USER, Aliases also skips the machine's alias file (~/.jaira/identity, JAIRA_HOME): it belongs to the checkout's owner just like git user.email, and the DoD says 'nur JAIRA_USER und dessen Aliasse aus der Personenliste'. Without JAIRA_USER nothing changes.
- **2026-11-23 21:10 · Alexander Sacharov** — critique+optimize in one pass (49 lines, no lock/parser). critique: none — fix sits in the callee (Aliases), every caller (gate, list --mine, tui, release) picks it up. optimize: four inline JAIRA_USER reads folded into identity.ByEnv().
- **2026-11-23 21:12 · Alexander Sacharov** — bughunt: change: internal/cli/whoami.go:37-39 and :54-56 - class: whoami still presents ~/.jaira/identity as effective while JAIRA_USER makes Aliases skip it; example 1: 'JAIRA_USER=Bob jaira whoami --help' still advises 'Add the name to the alias file rather than passing --force' (useless under JAIRA_USER, the release note says jaira people alias); example 2: 'JAIRA_USER=Bob jaira whoami --json' prints alias_file with the path and no marker, while the text output appends '(not read under JAIRA_USER)' - text and JSON answer differently; lines are older (e51f9ae78c, d296b37b64) but this change made them wrong; not a blocker (help/JSON wording, gates behave correctly); ran: scratch repo with owner git config + alias file, JAIRA_USER=Bob jaira whoami / --json / --help
- **2026-11-23 21:12 · Alexander Sacharov** — hunt report: 1 finding, 1 change (edge, no blocker), as a note on Q030WZ | read whole: core/identity (identity.go, identity_test.go, main_test.go), internal/cli/whoami.go, core/release/notes/Q030WZ.md, plus callers of Aliases/IsMe/Current and people.List.Find/Names | nothing left unchecked; identity and cli tests pass
- **2026-11-23 21:12 · Alexander Sacharov** — bughunt: done at 6af1ccd24 - 1 finding (1 change, 0 pre-existing) in 3 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/identity/identity.go core/identity/identity_test.go core/release/notes/Q030WZ.md internal/cli/whoami.go
- **2026-11-23 21:17 · Alexander Sacharov** — testing: class: user-facing text names commands that do not exist; example: core/release/notes/Q030WZ.md says 'list --mine' — 'jaira list --mine' exits 'unknown flag: --mine'; fix: name what really changes (the ownership checks of claim, move, release and pull, the TUI's own-ticket view) and check every command named in the note runs. Ride along: the bughunt's whoami edge (help still says 'add the name to the alias file', --json alias_file without the JAIRA_USER marker).
- **2026-11-23 21:18 · Alexander Sacharov** — critique+optimize round 2 (delta 6 lines): testing's note-flag finding and the bughunt's whoami edge both addressed; nothing new. DoD 1 proof updated to ByEnv().
- **2026-11-23 21:19 · Alexander Sacharov** — waiting: bughunt round 2 (background subagent) since 22:19
- **2026-11-23 21:19 · Alexander Sacharov** — hunt report: 0 findings, nothing written beyond this note and the closing line | read as a diff since 6af1ccd24: internal/cli/whoami.go, core/release/notes/Q030WZ.md (core/identity/identity.go and identity_test.go unchanged since); release-note claims traced to callers (claim/move via identityAliases, release and pull via IsMe, TUI via Aliases); ran: go test core/identity internal/cli (ok), whoami --json with JAIRA_USER=bob / unset / whitespace (alias_file_read false/true/true, matches ByEnv and raw) | nothing left unchecked
- **2026-11-23 21:19 · Alexander Sacharov** — bughunt: done at fb027aa07 - 0 findings in 4 modules, follow-up hunt (diff since 6af1ccd24); concurrency, error paths, input, permissions checked; files: core/identity/identity.go core/identity/identity_test.go core/release/notes/Q030WZ.md internal/cli/whoami.go
- **2026-11-23 21:20 · Alexander Sacharov** — waiting: review (fresh agent via run-lane.sh) since 22:20
- **2026-11-23 21:23 · Alexander Sacharov** — review: done at fb027aa07
- **2026-11-23 21:24 · Alexander Sacharov** — report: under JAIRA_USER Aliases = that name + its people-list names only (no git email, no ~/.jaira/identity); identity.ByEnv(); whoami marks the alias file unread (+ --json alias_file_read); note Q030WZ.md | signoff, context unknown | review accept, 1 note:limit; stacked on feat/dvpape, merge after DVPAPE; branch pushed to origin, PR is Alex's
- **2026-11-23 21:29 · Alexander Sacharov** — Accepted from signoff: angenommen:
- note: limit: someone who sets JAIRA_USER to their own name without being in the people list loses their git email and alias-file matches too; on purpose, the release note tells them to use 'jaira people alias'
- **2026-11-23 21:29 · Alexander Sacharov** — cockpit: accepted by the person, signoff -> done
