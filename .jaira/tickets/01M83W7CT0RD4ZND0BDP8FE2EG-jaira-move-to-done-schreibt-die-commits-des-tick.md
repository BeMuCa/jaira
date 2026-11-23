---
id: 01M83W7CT0RD4ZND0BDP8FE2EG
title: "jaira move --to done schreibt die Commits des Tickets, die auf der Branch des Worktrees liegen, statt einer leeren Liste"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wird ein Ticket im Worktree seiner Branch (Cockpit oder 'jaira move --to done --force') angenommen, stehen danach alle Commits mit der ID im Subject in commits:, ohne Nachtragen von Hand."
context: |-
  Was falsch ist (Teamlead eines anderen Repos, 23.11.2026, jaira 0.6.3): bei 8 Tickets schrieb die Annahme nach done eine leere commits-Liste, obwohl Commits mit der ID im Betreff auf der Branch des Worktrees liegen (Beispiel: Commit '<sha> ... (<ID>) ...' auf feat/<anderer-name>). Nachgetragen mit 'jaira set <id> commits=...'.
  Der Gate 'done requires commits' ging durch (also fand er sie), aber die geschriebene Liste war leer - Gate und Schreiber lesen vermutlich verschieden (Gate: Commits mit ID irgendwo; Schreiber: nur Historie der Ticket-Datei oder nur ab HEAD/Base?).
  Verwandt: VHGWQF (0.6.3: Accept im Worktree der tragenden Branch), core/gate/commits, logbook stempelt Commits.
definition-of-done: "Annahme nach done im Worktree einer anders benannten Branch schreibt alle Commits mit '(<ID>)' bzw. der ID im Subject nach commits: (Test mit Branch-Name != ID)"
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-11-23T09:36:35Z
updated-at: 2026-11-23T09:41:09Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-68452
claimed-at: 2026-11-23T09:37:10Z
outcome-what: "core/move.Move writes env.DeriveCommits (unioned with recorded shas) in the move's own write when the target lane requires commits"
outcome-why: "the gate passed on derived commits but only a doorway settle stamped them, so a plain move to done recorded none and a later stamp from another checkout wrote an empty list"
outcome-resolves: "accepting into done in a worktree of a differently named branch leaves commits: empty"
---

# jaira move --to done schreibt die Commits des Tickets, die auf der Branch des Worktrees liegen, statt einer leeren Liste

## Definition of Done

- [x] Annahme nach done im Worktree einer anders benannten Branch schreibt alle Commits mit '(<ID>)' bzw. der ID im Subject nach commits: (Test mit Branch-Name != ID)
  proof: internal/cli/donecommits_test.go TestMoveToDoneInAWorktreeWritesTheBranchCommits
- [x] Gate und Schreiber benutzen dieselbe Funktion für die Commit-Liste; was der Gate findet, steht danach im Ticket (Test)
  proof: core/move/move.go:154 (env.DeriveCommits, same as gate); core/move/move_test.go TestMoveIntoDoneWritesTheCommitsItsGateFound
- [x] Release-Note core/release/notes/<ID>.md
  proof: core/release/notes/8FE2EG.md:1

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 09:37 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/6958.sock
- **2026-11-23 09:38 · Alexander Sacharov** — decision: no open decisions. Cause: gate.CheckAdvance derives commits via env.DeriveCommits (CommitsForTicket, reads HEAD) for a requires-commits lane, but move.Move writes them only through Prepare, which runs only when a doorway lane files the ticket. A plain move to done therefore records nothing; a later logbook/stamp from another checkout (master, not reaching the branch) finds nothing. Fix: core/move stamps env.DeriveCommits (MergeCommits union) in the move's own write when the target lane requires commits — same function the gate used, same HEAD, no new flag.
