---
id: 01M7ZYEXTPHSC0PNW5STG8HNHT
title: "Ohne Bughunt auf HEAD kein signoff: der Gate verweigert den Schritt, wenn bughunt.state fail ist"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Ticket, das Code geändert hat, kommt nur nach signoff, wenn der letzte Bughunt den Stand HEAD abdeckt."
context: |-
  Was heute schiefging (21.11.): 1MYFKE kam nach signoff, obwohl nach dem letzten Bughunt (10e5c341c) noch ein Commit (c3728ba70) die gejagten Dateien änderte. jaira show 1MYFKE --brief zeigte bughunt.state = fail, niemand hat es aufgehalten.
  Heute prüft das nur der Teamlead von Hand (grep auf "— bughunt: done at <sha>") und das Cockpit zeigt es im Explain an; der Gate selbst kennt bughunt nicht.
  Die Daten gibt es schon: jaira show <id> --brief liefert bughunt.at und bughunt.state (ok | fail | unknown | skipped, null = kein Hunt), eingeführt mit 0.5.15.
  Gewünscht: jaira move <id> --to signoff (und jeder Weg in einen Personen-Lane hinter testing) wird mit exit 3 verweigert, wenn das Ticket Code ausserhalb von .jaira/ geändert hat und bughunt.state fail oder null ist; skipped (stop-after-fix) und ein Ticket ohne Codeänderung gehen durch. Die Meldung nennt den Hunt-Commit und HEAD und sagt: /jaira-role-bughunt <id> laufen lassen.
  unknown (uncommittete Änderungen) ebenfalls verweigern, mit eigener Meldung "erst committen".
  Person mit --force kommt durch.
definition-of-done: "move nach signoff mit bughunt.state fail oder ohne Hunt bei Codeänderung: exit 3 mit Hunt-Commit, HEAD und Hinweis (Test in core/gate)"
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-11-21T20:58:41Z
updated-at: 2026-11-21T21:06:21Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2700790
claimed-at: 2026-11-21T20:59:16Z
outcome-what: "jaira move into signoff refuses a ticket that changed code and whose newest bug hunt does not cover HEAD (needs_bughunt), or whose hunted files are uncommitted (hunt_uncommitted), on boards whose testing lane owes hunt-verdict"
outcome-why: 1MYFKE reached signoff on 21.11. with a commit after its last hunt and bughunt.state fail; only the teamlead checked that by hand
outcome-resolves: "G8HNHT: the gate reads bughunt.state; ok, skipped and tickets without code changes pass; --force overrides"
---

# Ohne Bughunt auf HEAD kein signoff: der Gate verweigert den Schritt, wenn bughunt.state fail ist

## Definition of Done

- [x] move nach signoff mit bughunt.state fail oder ohne Hunt bei Codeänderung: exit 3 mit Hunt-Commit, HEAD und Hinweis (Test in core/gate)
  proof: core/gate/bughunt_test.go TestSignoffNeedsHuntOnHEAD (stale hunt, no hunt); core/gate/gate.go huntAgainstMove; exit 3 via refusalCode→ExitValidation
- [x] skipped, ok und Tickets ohne Codeänderung ausserhalb .jaira/ gehen durch (Tabellentest)
  proof: core/gate/bughunt_test.go TestSignoffNeedsHuntOnHEAD (covers HEAD, stop-after-fix, no code changed); core/gitrepo/huntcheck_test.go TestHuntCheck (board data only)
- [x] unknown wird mit 'erst committen' verweigert; --force der Person kommt durch (Test)
  proof: core/gate/bughunt_test.go TestSignoffNeedsHuntOnHEAD/uncommitted; core/move/move_test.go TestForceOverridesStaleHunt
- [x] Release-Note core/release/notes/<ID>.md
  proof: core/release/notes/G8HNHT.md
- [x] Ein Board ohne Lane, die hunt-verdict produziert, verhält sich wie vorher (Test)
  proof: core/gate/bughunt_test.go TestHuntRuleScope (board without hunt-verdict, Hunt never called)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-21 20:59 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2274544.sock
- **2026-11-21 21:00 · Alexander Sacharov** — decision (worker, settled): the gate guards only lanes with requires-human-exit (signoff), not the human lane — human takes questions, and testing's own max-rounds escape goes into human; blocking a question on a missing hunt would trap that. 'code changed' = a commit of the ticket (CommitsForTicket) touches a path outside .jaira/ (gitrepo anyCode rule); uncommitted code with no commits yields no hunt requirement but the unknown state still refuses. Check is injected into gate.Env like DeriveCommits, so gate stays filesystem-free.
- **2026-11-21 21:00 · Alexander Sacharov** — question (open, for the person): A fresh 'jaira init' board has signoff but no testing lane and no bughunt. Should the gate apply there too? 1) Only on boards with a lane that produces hunt-verdict (Recommended) — boards without hunts unchanged, no surprise for users who never run bughunt. 2) On every board — every code ticket needs a /jaira-role-bughunt note before signoff, or --force.
- **2026-11-21 21:00 · Alexander Sacharov** — decision: Variante 1 (Teamlead 21.11.): der Gate gilt nur auf Boards, deren Lanes hunt-verdict produzieren; Boards ohne Bughunt ändern sich nicht. Gate nur vor Lanes mit requires-human-exit (signoff), nicht vor human.
- **2026-11-21 21:06 · Alexander Sacharov** — in-progress: hunt state logic moved from internal/cli/tickets.go to core/gitrepo/hunt.go (LastHunt, TouchesCode, HuntCheck) so CLI and TUI share it; gitrepo cannot import gate (core/board tests import gitrepo → import cycle), so gitrepo.HuntReport mirrors gate.HuntReport and callers convert. A TUI move (Interactive) is not refused, like the verdict rail: a person in the board decides. 'unknown' also covers git failing to answer; the message says so.
