---
id: 01M7NMRQS4W0GMFGFES558DCYM
title: Cockpit-Explain startet ohne Trust-Dialog und schreibt im richtigen Worktree
status: critique
ready: true
creator: Alexander Sacharov
goal: "Ein Druck auf Explain im Cockpit laeuft ohne Eingriff durch: der Tab startet claude ohne Rueckfrage, der Chat schreibt seine explain:-Note, und der Tab schliesst sich danach von selbst."
context: |-
  Was heute kaputt ist (live 17.11., Herdr App, Claude Code 2.1.326, ein Ticket in einem anderen Repo):
  1. spawn.sh startet Explain/Discuss mit 'claude --add-dir /tmp/jaira-chat.XXXX'. Claude Code 2.1.326 fragt dafuer 'Adding a new directory ... trust?' - der Tab bleibt im Dialog stehen, nichts wird getippt.
  2. Nach Bestaetigung von Hand lief der Chat im Herdr-Review-Worktree (Branch review/feat/<id>). Dort liegt das Ticket nicht auf der Platte; 'jaira note <id>' scheitert mit 'on its ref and not on your disk; pull it first'. Ohne explain:-Note schliesst das Cockpit den Tab nie.
  Ursache 2: worktreeFor() in core/mod/builtin/jaira-cockpit/hooks/board.ts nimmt den ersten Worktree, dessen Branch auf /<id> endet - review/feat/<id> matcht genauso wie feat/<id>, wo das Ticket liegt.
  Geprueft: 'claude -p --allowedTools=Read(//tmp/jaira-chat.X/**)' liest die Datei ohne Rueckfrage; ohne Regel verweigert. --allowedTools ist variadisch: nur in der Form mit '=' schreiben.
definition-of-done: spawn.sh startet Explain/Discuss ohne --add-dir; das Lesen der Auftragsdatei erlaubt eine Read-Regel auf genau ihren Ordner (auch auf Windows)
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-11-17T20:56:53Z
updated-at: 2026-11-17T21:12:11Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-27759
claimed-at: 2026-11-17T20:58:46Z
outcome-what: "Explain/Discuss starten claude mit einer Read-Regel auf den Auftragsordner statt --add-dir; das Cockpit waehlt den Worktree des Ticket-Branches statt Herdrs review/feat/<id>"
outcome-why: "Claude Code 2.1.326 haelt den Tab am Trust-Dialog fuer --add-dir an; im Review-Worktree scheitert jaira note, die explain:-Note fehlt und der Tab schliesst nie"
outcome-resolves: DoD 1-4
review-summary: none
review-gaps: "none - pass 2 added one //wintrap:ok comment line, nothing to remove"
test-verdict: "pass: go test ./... RC=0 (wintrap clean); claude plugin test . 176/177 on this branch and on master, same single pre-existing fail (borderStyle round); DoD 1-4 checked in the tree: spawn.sh chat_rule + --allowedTools on both start lines, board.ts ticketsTree with review/feat fixture first, notes/58DCYM.md present"
round: 3
---

# Cockpit-Explain startet ohne Trust-Dialog und schreibt im richtigen Worktree

## Definition of Done

- [x] spawn.sh startet Explain/Discuss ohne --add-dir; das Lesen der Auftragsdatei erlaubt eine Read-Regel auf genau ihren Ordner (auch auf Windows)
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:312 chat_rule + :321/:340 --allowedTools; TestSpawnChatsStartOnTheChatModelWithTheirBrief, TestSpawnStartsClaudeWhereItsShellRuns/git bash chat
- [x] worktreeFor/branchesFor waehlen bei mehreren passenden Worktrees den Branch des Tickets (feat/<id>), nicht review/feat/<id>
  proof: core/mod/builtin/jaira-cockpit/hooks/board.ts ticketsTree() + carrierBranch() via branchParts; tests/board.test.ts WORKTREES fixture (review/feat/mg677r first) and 'a branch left behind beside the Herdr review branch of it'
- [x] Tests in core/role und im Cockpit decken beides ab; go test ./core/... und claude plugin test gruen bis auf vorbestehende Fehler
  proof: go test ./... RC=0 (incl. internal/wintrap, role_test.go:736 marked //wintrap:ok); claude plugin test . 176/177, der eine Fehler (borderStyle round) auch auf master
- [x] Release-Notiz core/release/notes/<ID>.md
  proof: core/release/notes/58DCYM.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] spawn.sh: --add-dir durch '--allowedTools=Read(//<ordner>/**)' ersetzen, Windows-Pfad in POSIX-Form //c/...
- [x] board.ts: unter passenden Worktrees den Branch mit den wenigsten Teilen waehlen (worktreeFor, branchesFor)
- [x] Tests: role_test.go Startzeilen, Cockpit-Fixture mit review/feat/mg677r vor feat/mg677r

## Progress
- **2026-11-17 20:58 · Alexander Sacharov** — Interaktiver Start im Herdr-Tab mit --allowedTools nicht live geprueft, nur claude -p (liest ohne Rueckfrage, ohne Regel verweigert). Windows-Regelform //c/... nach Claude-Code-Doku zu Permission-Pfaden, nicht auf echtem Windows getestet. Fehler 'borderStyle round' in board.test.ts besteht auf master unabhaengig davon.
- **2026-11-17 20:58 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/51013.sock
- **2026-11-17 20:59 · Alexander Sacharov** — critique: nothing raised. board.ts ticketsTree (fewest branch parts) is general rather than Herdr-specific - a fix/<id> or feat/<id> branch beats review/feat/<id>, a detached tree matched by path loses to any branch; spawn.sh keeps one rule variable for both start lines, same shape as the --add-dir it replaces.
- **2026-11-17 21:02 · Alexander Sacharov** — testing: go test ./... RC=1: internal/wintrap TestRepositoryIsClean: 'core/role/role_test.go:736: windows trap 5: a path is glued together with a literal "/" and handed to strings.HasSuffix'. The string is a Claude Code permission rule in a shell line, not a filesystem path - mark the site //wintrap:ok with that reason. go test ./core/... green; cockpit 176/177, the one fail (borderStyle round) pre-existing per earlier note.
- **2026-11-17 21:04 · Alexander Sacharov** — critique pass 2: only the testing finding to check - role_test.go:736 now carries //wintrap:ok with its reason, same marker form as core/role/jairabin_test.go:34. Nothing else reopened.
- **2026-11-17 21:08 · Alexander Sacharov** — bughunt: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:310-312 - chat_rule is Read(/$chat_dir/**) on the path mktemp gives, never resolved; Claude Code 2.1.326 allows a Read only when the rule matches the symlink path AND the resolved one, so wherever TMPDIR runs through a symlink (macOS: /var/folders/... -> /private/var/folders/...) the chat's Read of its request file is refused and the Explain/Discuss tab stops at a permission prompt again - the case --add-dir used to cover. Fix: chat_dir="$(cd "$(mktemp -d ...)" && pwd -P)". ran: claude -p --model haiku with /tmp/jbh-link -> /tmp/jbh-a; rule on the link alone: DENIED (twice), rule on the target alone: DENIED, both rules: file read
- **2026-11-17 21:08 · Alexander Sacharov** — bughunt: core/mod/builtin/jaira-cockpit/hooks/register.tsx:551 with board.ts:418 (carrierBranch) - with no worktree of the ticket left, the PR branch is the first of git branch --list '*/<slug>' alphabetically, not the one with the fewest parts that ticketsTree now prefers; with JAIRA_BRANCH_PREFIX sorting after 'review' (task, work, story) Herdr's review/task/<slug> comes first and Accept opens the pull request from the review branch. ran: scratch repo with branches task/mg677r and review/task/mg677r -> git branch --list '*/mg677r' prints review/task/mg677r first (with feat/ the order is right)
- **2026-11-17 21:08 · Alexander Sacharov** — bughunt: done at 3c409f36 - 2 findings in 4 modules (spawn.sh, cockpit hooks/board.ts with its caller register.tsx, role_test.go, board.test.ts); concurrency, error paths, input, permissions checked
- **2026-11-17 21:12 · Alexander Sacharov** — in-progress round 3 (bughunt): (1) spawn.sh resolves the chat folder with pwd -P before building the Read rule - verified by bughunt that claude needs the rule to match the real path; (2) carrierBranch in board.ts picks the branch with fewest parts too, sharing branchParts with ticketsTree; test 'a branch left behind beside the Herdr review branch' fails on the old code, passes now. Windows 8.3/UNC TEMP left unhandled: not reproducible here, regex falls back to the MSYS path. Parent asked: rebase on upstream/master before review.
