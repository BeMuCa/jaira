---
id: 01M851YE2ACMD981Z6V1Q030WZ
title: Unter JAIRA_USER zaehlt jaira die per E-Mail zugewiesenen Tickets des Checkout-Besitzers als die der anderen Person
status: critique
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
commits: []
created-at: 2026-11-23T20:35:47Z
updated-at: 2026-11-23T21:18:36Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-37226
claimed-at: 2026-11-23T21:01:32Z
outcome-what: Release note names the real ownership checks instead of the nonexistent list --mine; whoami help points to jaira people alias under JAIRA_USER and --json gains alias_file_read
outcome-why: "testing found the note named a flag jaira does not have; the bughunt found whoami still presented the alias file as in effect"
outcome-resolves: "under JAIRA_USER the checkout owner's git email and alias file no longer make the owner's tickets the other person's"
review-summary: none
review-gaps: "folded the four strings.TrimSpace(os.Getenv(JAIRA_USER)) reads (currentIn, Aliases, whoami) into identity.ByEnv(); raw() keeps its own read because it needs the value; nothing dead, no cost on hot paths"
test-verdict: "fail: release note names 'list --mine', a flag jaira does not have (unknown flag: --mine); suite green RC=0 (-race -count=1), DoD 1 verified (grace refused move exit 3, owner moves exit 0, two parallel sets both land)"
hunt-verdict: "pass: bughunt done at 6af1ccd24, 1 change (whoami help/JSON wording, edge, no blocker), 0 pre-existing"
round: 2
back-from: testing
back-by: lane
back-count: 1
---

# Unter JAIRA_USER zaehlt jaira die per E-Mail zugewiesenen Tickets des Checkout-Besitzers als die der anderen Person

## Definition of Done

- [x] Unter JAIRA_USER liefert Aliases nur JAIRA_USER und dessen Aliasse aus der Personenliste, nicht die git user.email des Checkouts; Test
  proof: core/identity/identity.go Aliases (byEnv return); TestAliasesUnderJairaUserLeaveOutTheOwnersNames
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
