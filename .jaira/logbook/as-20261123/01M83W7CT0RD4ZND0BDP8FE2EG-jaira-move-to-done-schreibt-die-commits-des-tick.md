---
id: 01M83W7CT0RD4ZND0BDP8FE2EG
title: "jaira move --to done schreibt die Commits des Tickets, die auf der Branch des Worktrees liegen, statt einer leeren Liste"
status: done
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
commits:
  - fba739900a462db6a10fbbafc438c3797d80817a
  - a4e8d736467e538c2be62e8515c277c8400f67ca
  - e657d614baee1b817294ebf94318aac2b0ba667c
created-at: 2026-11-23T09:36:35Z
updated-at: 2026-11-23T12:02:16Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-68452
claimed-at: 2026-11-23T09:37:10Z
outcome-what: "core/move.Move writes env.DeriveCommits (unioned with recorded shas) in the move's own write when the target lane requires commits"
outcome-why: "the gate passed on derived commits but only a doorway settle stamped them, so a plain move to done recorded none and a later stamp from another checkout wrote an empty list"
outcome-resolves: "accepting into done in a worktree of a differently named branch leaves commits: empty"
review-summary: "core/move.Move now calls env.DeriveCommits (the same git derivation the requires-commits gate uses, reading HEAD of the checkout it runs in) whenever the target lane requires commits, and writes the union with the already recorded shas (ticket.MergeCommits) in the move's own write. Before, only a doorway settle stamped, so a plain move into done passed the gate on derived commits and recorded none. All four write paths (CLI move, TUI moves, signoff accept) go through core/move, so they all get it. Tests: core/move (derived+recorded union, no stamp outside a requires-commits lane) and an internal/cli test with a worktree on branch feat/another-name."
review-gaps: |-
  note: limit: commits recorded at done on a branch that is later squash- or rebase-merged stay as branch shas unreachable from master; MergeCommits keeps recorded shas by design
  note: elsewhere: the worktree holds uncommitted regenerated .jaira/lanes/*.md, AGENTS.md and CLAUDE.md that are not this change; commit them separately or drop them before the PR
test-verdict: "pass: go build+vet, go test -race -count=1 ./... RC=0; DoD 1-3 verified (both named tests run, note file present); scratch repo, worktree on feat/other-name: move --to done --force writes the fix(<HANDLE>) commit + the ticket-file commit; two parallel moves of two tickets into done both land with their own commits; a refused move (exit 3) writes no commits"
hunt-verdict: "pass: bughunt done at a4e8d7364, 0 change, 0 pre-existing (follow-up after chore commit)"
review-verdict: "accept: DoD 1-3 met by core/move/move.go:147-160 and the two new tests; no defect found along concurrency (derive outside the lock, union inside Mutate), error paths (derive failure yields nil = old behaviour), input (empty/duplicate shas handled by MergeCommits) or boundaries"
review-round: 1
review-verdict-gaps: |-
  note: limit: commits recorded at done on a branch that is later squash- or rebase-merged stay as branch shas unreachable from master; MergeCommits keeps recorded shas by design
  note: elsewhere: the worktree holds uncommitted regenerated .jaira/lanes/*.md, AGENTS.md and CLAUDE.md that are not this change; commit them separately or drop them before the PR
review-check: "1. cd into a worktree whose branch name is not the ticket id  2. make a commit with the ticket handle in the subject, e.g. git commit --allow-empty -m 'fix(ABC123): x'  3. jaira move ABC123 --to done --force  4. jaira show ABC123 --json | jq .commits — shows that commit's sha. Automated: go test ./core/move -run TestMoveIntoDone && go test ./internal/cli -run TestMoveToDoneInAWorktreeWritesTheBranchCommits"
question: ""
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
- **2026-11-23 09:41 · Alexander Sacharov** — critique: none. Considered stamping in flow.go's Stage instead — rejected: Stage persists before the gate (would write on a refused move) and misses the TUI/signoff paths, which all go through core/move. Gate and move each call env.DeriveCommits once (two git-log runs per move into done) — let stand: returning the derivation out of gate.CheckAdvance would widen a pure API for one caller. Doorway settle's Prepare stamps again — idempotent union, left.
- **2026-11-23 09:43 · Alexander Sacharov** — waiting: bughunt (background agent) since 10:43
- **2026-11-23 09:44 · Alexander Sacharov** — hunt report: 0 findings (0 change, 0 pre-existing), nothing written as bughunt: notes | read whole: core/move/move.go with its callers' paths (internal/cli/flow.go move, internal/tui/signoff.go accept), core/gitrepo/derive.go CommitsForTicket, core/ticket/trim.go MergeCommits/StampCommits, gate requires-commits check; tests and the regenerated .jaira/lanes, AGENTS.md, CLAUDE.md read as diffs (generated text, no code) | checked --follow rename-jump onto other ticket files across all board and logbook tickets (0 hits); ran go test core/move core/ticket core/gate and internal/cli -run Done|Commit (ok); not checked: squash/rebase merges leaving branch shas recorded at done unreachable from master - MergeCommits keeps recorded shas by design
- **2026-11-23 09:44 · Alexander Sacharov** — bughunt: done at fba739900 - 0 findings (0 change, 0 pre-existing) in 4 modules, first hunt; concurrency, error paths, input, permissions checked; files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/move/move.go core/move/move_test.go core/release/notes/8FE2EG.md internal/cli/donecommits_test.go
- **2026-11-23 09:45 · Alexander Sacharov** — waiting: review (fresh agent, Herdr tab via run-lane.sh) since 10:45
- **2026-11-23 09:46 · Alexander Sacharov** — review: done at fba739900
- **2026-11-23 09:46 · Alexander Sacharov** — review: accept written, move to signoff refused by the hunt gate: the worktree holds uncommitted regenerated .jaira/lanes/{critique,in-progress,review,testing}.md, AGENTS.md, CLAUDE.md, which the hunt at fba739900 read. Question: what happens to them?
1. Commit them separately on this branch as chore(8FE2EG): board regenerated by jaira update, then move to signoff (recommended; generated text, the hunt read it already, no new hunt needed)
2. Drop them (git checkout -- those files) and move on; the regeneration comes back with the next jaira update
3. Move with --force as the person deciding
- **2026-11-23 09:47 · Alexander Sacharov** — decision (worker): review's question settled as option 1 — the files are the board refresh the installed jaira 0.6.4 wrote on its own (version marker + lane prompts), not this change; dropping them only brings them back on the next jaira call. Committed separately as chore(8FE2EG), same as d43db5b25 did for GY8T52.
- **2026-11-23 09:47 · Alexander Sacharov** — waiting: re-hunt at a4e8d73 (signoff hunt gate after the chore commit), background agent since 10:47
- **2026-11-23 09:47 · Alexander Sacharov** — hunt report: 0 findings, nothing written | none read whole; read as a diff since fba739900: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md (generated prose only, lanes still load: jaira lanes exit 0); core/move/move.go, move_test.go, donecommits_test.go and the release note are unchanged since the last round | nothing left unchecked
- **2026-11-23 09:47 · Alexander Sacharov** — bughunt: done at a4e8d7364 - 0 findings (0 change, 0 pre-existing) in 4 modules, follow-up hunt since fba739900 (only generated lane prompts and agent docs changed, no code); concurrency, error paths, input, permissions checked; files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/move/move.go core/move/move_test.go core/release/notes/8FE2EG.md internal/cli/donecommits_test.go
- **2026-11-23 09:48 · Alexander Sacharov** — report: core/move writes the commits its requires-commits gate derived in the move itself (+2 tests, release note); board refresh committed separately as chore | signoff, context unknown | open: 2 review notes (limit: squash-merge leaves branch shas; elsewhere: board refresh, now committed) — person accepts
- **2026-11-23 11:59 · Alexander Sacharov** — Accepted from signoff: angenommen:
- note: limit: commits recorded at done on a branch that is later squash- or rebase-merged stay as branch shas unreachable from master; MergeCommits keeps recorded shas by design
- note: elsewhere: the worktree holds uncommitted regenerated .jaira/lanes/*.md, AGENTS.md and CLAUDE.md that are not this change; commit them separately or drop them before the PR
- **2026-11-23 11:59 · Alexander Sacharov** — cockpit: accepted by the person, signoff -> done
