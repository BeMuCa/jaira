---
id: 01M83NMV79RR9HHN2CCHVHGWQF
title: "Cockpit-Accept führt den Move im Worktree der Branch aus, die die Commits des Tickets trägt, auch wenn ihr Name nicht die Ticket-ID ist"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Accept im Cockpit findet die Commits eines Tickets, dessen Arbeit auf einer anders benannten Branch liegt (z.B. Nachbesserung auf feat/<anderes-ticket>), und verschiebt nach done ohne Fehler 'git has none for <ID>'."
context: |-
  Was falsch ist (Alex 23.11.2026, anderes Repo, jaira 0.6.2): Accept im Cockpit meldet 'lane "done" requires the commits that carry this change, and git has none for <ID> yet', obwohl Commits mit '(<ID>)' im Subject existieren.
  Ursache: die Commits liegen auf feat/<anderer-slug> in dessen Worktree (der Ticket wurde dort als Nachbesserung gearbeitet), Accept läuft im Haupt-Checkout auf main - von dort sind die Commits nicht erreichbar. Von Hand im Worktree der Branch funktioniert 'jaira move <id> --to done --force'.
  QTQW15 (in 0.6.2) wählt die Branch: eigener Worktree > <prefix>/<slug> > genau ein Holder; ein Holder wird über Commits mit '(<ID>)' gefunden - prüfen, warum er hier nicht gegriffen hat (Accept-Pfad in core/mod/builtin/jaira-cockpit/hooks/register.tsx, carrierFor/holdersOf in board.ts).
  Daneben: der Commit war nicht gepusht; jaira findet Commits nur lokal - das ist ok, aber die Meldung sollte die Branch nennen, wo die Commits liegen.
definition-of-done: "Accept auf einen Ticket, dessen Commits '(<ID>)' nur auf einer anders benannten Branch mit eigenem Worktree liegen: der Move läuft in diesem Worktree und gelingt (Test)"
tags:
  - cockpit
blocked-by: []
related: []
commits: []
created-at: 2026-11-23T07:41:36Z
updated-at: 2026-11-23T08:04:22Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-70439
claimed-at: 2026-11-23T07:42:00Z
outcome-what: "Cockpit-Accept fuehrt den Move im Worktree der Branch aus, die die Commits des Tickets traegt (acceptTreeFor + ticketHere carried); mehrdeutig: kein Move, Kandidaten genannt. jaira move --to done nennt Branches mit Commits der ID, die HEAD nicht erreicht (gitrepo.BranchesNamingTicket, gate.Env.CommitsElsewhere)."
outcome-why: "Accept lief immer im Haupt-Checkout und fand Commits auf anders benannten Branches (Nachbesserung auf feat/<anderes-ticket>) nicht; die Meldung sagte nur 'commit it', obwohl committet war."
outcome-resolves: "VHGWQF: Accept-Fehler 'git has none for <ID>' bei Arbeit auf fremder Branch"
review-summary: "none"
round: 4
back-from: testing
back-by: lane
back-count: 2
review-gaps: "none: fix reuses commitNaming/ticketBranch/acceptTreeFor; carried boolean replaced by the target lane (no dead parameter left); nothing orphaned"
test-verdict: "pass: round 2 - go build ok, core/gate + core/gitrepo green (rc=0), cockpit 475/475 incl. new 'branches named after it but nothing of its own beyond HEAD' (moves in @/repo) and ambiguous text with the target lane; Go side unchanged since round 1 (-race green there); DoD 1-4 proofs still name existing tests"
hunt-verdict: "fail: 1 change finding at 3e2816243 (carrier picked by slug name, not by the branch holding the commit)"
---

# Cockpit-Accept führt den Move im Worktree der Branch aus, die die Commits des Tickets trägt, auch wenn ihr Name nicht die Ticket-ID ist

## Definition of Done

- [x] Accept auf einen Ticket, dessen Commits '(<ID>)' nur auf einer anders benannten Branch mit eigenem Worktree liegen: der Move läuft in diesem Worktree und gelingt (Test)
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts test 'a ticket worked on another ticket's branch finds that branch by its own commit, the one branch holding it' and describe 'where Accept moves a ticket'
- [x] Findet sich keine eindeutige Branch: kein Move im Haupt-Checkout, sondern das Band nennt die Kandidaten (wie QTQW15) (Test)
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts test 'several branches hold its own commit: no pull request, the press names them'
- [x] Die Fehlermeldung von 'jaira move --to done' ohne erreichbare Commits nennt Branches, auf denen Commits mit der ID liegen, wenn es welche gibt (Test)
  proof: core/gate/commits_test.go TestNeedsCommitsNamesBranchesElsewhere, core/gitrepo/derive_test.go TestBranchesNamingTicket
- [x] Release-Note core/release/notes/<ID>.md, shipped.sha256
  proof: core/release/notes/VHGWQF.md, core/shipped/shipped.sha256 (+3 lines via scripts/shipped-history.sh --write)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 07:42 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/6958.sock
- **2026-11-23 07:43 · Alexander Sacharov** — Ursache: accept() -> ticketHere() sucht nur den Worktree mit dem Ticket-Slug (ownTreeFor); die QTQW15-Carrier-Wahl (ticketBranch) nutzen nur PR-Schritt und Give back. Und Go CommitsForTicket sucht nur 'git log' ab HEAD. Offene Entscheidungen: keine. decision: Carrier-Branch ohne Worktree -> Move wie bisher im Haupt-Checkout (kein neuer Worktree durch Accept), die Go-Meldung nennt dann die Branch (DoD 3). Mehrdeutig -> kein Move, Meldung nennt Kandidaten. Nur Accept bekommt die Carrier-Wahl, andere Presses (note, give back) bleiben wie sie sind.
- **2026-11-23 07:44 · Alexander Sacharov** — Go-Teil fertig: gitrepo.BranchesNamingTicket (git log --branches --not HEAD --grep, dann ein 'git branch --contains' mit allen shas), gate.Env.CommitsElsewhere haengt '. Commits naming X sit on <branches>: run the move in that branch's worktree, or merge it first' an. Als naechstes Cockpit: ticketHere fuer Accept mit Carrier-Wahl.
- **2026-11-23 07:52 · Alexander Sacharov** — decision: Basis fuer die Carrier-Suche bei Accept ist HEAD des Haupt-Checkouts, nicht origin-default: nur Commits, die der Move dort nicht erreicht, schicken ihn woanders hin; in einem Fork mit altem master haette sonst jede spaetere Branch den Commit eines gemergten Tickets und Accept waere mehrdeutig verweigert. Der Ticket-File landet dann im Worktree der Carrier-Branch und bleibt dort (wie von Hand).
- **2026-11-23 07:52 · Alexander Sacharov** — critique: ticketHere liest git worktree list bei Accept zweimal (ownTree + Carrier); einmal lesen, acceptTreeFor entscheidet own-tree zuerst. Sonst nichts: Go-Meldung nutzt namesTicket wie DeriveCommits (Body inklusive) - bewusst gleich wie die Go-Seite, das Cockpit bleibt bei commitNaming (nur eigener Scope, QTQW15).
- **2026-11-23 07:56 · Alexander Sacharov** — waiting: bug hunt (background agent) since 08:56
- **2026-11-23 07:56 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:584 (via ticketBranch :1769-1775 and carrierFor board.ts:90-92) - class: Accept's carried path settles the carrier by branch name before it asks whether the main checkout's HEAD already reaches the ticket's commits, so a name-only ambiguity refuses an Accept that needs no carrier at all; example: feat/zz99zz and fix/zz99zz both still exist, both merged into master, no worktree of either - Accept now fails with acceptAmbiguous ('nothing moved'), where before this change it moved the ticket in the main checkout; ran: scratch repo, git branch --list '*/zz99zz' lists both while git log --branches --not HEAD lists 0 commits; ticketBranch returns the slug ambiguity before reaching the HEAD log (read). Not a blocker: it refuses safely with a message; ticketBranch is QTQW15's, the new caller with base 'HEAD' is this change
- **2026-11-23 07:56 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/texts.ts:107 - class: acceptAmbiguous tells the person a hard-coded move line while Accept moves to w.acceptTo, the board's terminal lane, and names worktrees that may not exist; example: on a board whose terminal lane is not called 'done' the line says 'jaira move <id> --to done --force', which is refused; and where neither candidate branch has a worktree ('in the right one's worktree') there is no worktree to run it in; read. Not a blocker: wrong advice in a refusal text, nothing moves
- **2026-11-23 07:56 · Alexander Sacharov** — hunt report: 2 findings, both change, no blocker, notes on VHGWQF | read whole: register.tsx ticketHere/accept/accepting/ticketBranch/workPlace path, board.ts worktree/carrier helpers (ownTreeFor, branchesFor, worktreeOfBranch, carrierFor, acceptTreeFor, holdersOf, commitNaming), texts.ts, core/gitrepo/derive.go, core/gate/gate.go CheckAdvance commits branch + onBranches, internal/cli/root.go loadBoardEnv, internal/tui/model.go gateEnv + signoff.go accept; .jaira/lanes/*, AGENTS.md, CLAUDE.md are regenerated text, read as diff | not run: the cockpit engine path end to end (needs a live Claude Code engine); the 64-commit cap in BranchesNamingTicket can leave a branch unnamed that holds only older commits - incomplete hint text only, not reported
- **2026-11-23 07:56 · Alexander Sacharov** — bughunt: done at f82525665 - 2 findings (2 change, 0 pre-existing) in 8 modules, first hunt; concurrency, error paths, input, permissions checked; files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/gate/commits_test.go core/gate/gate.go core/gitrepo/derive.go core/gitrepo/derive_test.go core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/VHGWQF.md core/shipped/shipped.sha256 internal/cli/root.go internal/tui/model.go
- **2026-11-23 07:57 · Alexander Sacharov** — fix plan: carried path asks first whether commits of the ticket's own exist outside HEAD (git log --branches --not HEAD + commitNaming); none -> main checkout, no carrier search at all. acceptAmbiguous gets the target lane and no longer assumes a worktree.
- **2026-11-23 07:59 · Alexander Sacharov** — critique round 2: let stand - the pre-check 'git log --branches --not HEAD' repeats once inside ticketBranch when a commit lies beyond HEAD; folding it in would change ticketBranch's name-first order, which is QTQW15's rule for Give back. One extra git log on Accept only.
- **2026-11-23 08:00 · Alexander Sacharov** — waiting: bug hunt round 2 (background agent) since 09:00
- **2026-11-23 08:00 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:587-590 (ticketHere) with register.tsx:1776-1781 (ticketBranch) and board.ts:584-586 (carrierFor) - class: the gate checks that some own commit lies beyond HEAD, but the carrier is then picked among every branch named after the slug, not among the branches that hold that commit; examples: (a) feat/<slug> merged into master and kept, fix/<slug> carries a new fix(<ID>) commit -> Accept refuses as ambiguous although only fix/<slug> is beyond HEAD; (b) feat/<slug> merged and kept, the new own commit sits on another branch -> the merged feat/<slug> is taken as carrier (its worktree, or the main checkout when it has none), so the move does not run where the commit is; not a blocker: the outcome is the old refusal or the old main-checkout move, not a new wrong move; read
- **2026-11-23 08:00 · Alexander Sacharov** — hunt report: 1 finding, 1 change (edge, not a blocker), note on VHGWQF | read as a diff since f82525665: register.tsx (ticketHere and ticketBranch read whole), texts.ts, board.test.ts, board.ts carrierFor/acceptTreeFor/commitNaming as context, lane prompts and AGENTS.md/CLAUDE.md (prose, regenerated block), shipped.sha256 | not run: the cockpit test suite (no runner in this worktree); finding confirmed by reading the call path
- **2026-11-23 08:01 · Alexander Sacharov** — bughunt: done at 3e2816243 - 1 finding (1 change, 0 pre-existing) in 8 modules, follow-up hunt since f82525665; concurrency, error paths, input, permissions checked; files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/gate/commits_test.go core/gate/gate.go core/gitrepo/derive.go core/gitrepo/derive_test.go core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/VHGWQF.md core/shipped/shipped.sha256 internal/cli/root.go internal/tui/model.go
- **2026-11-23 08:01 · Alexander Sacharov** — decision: smaller design after 2 hunt rounds on the same rule - Accept no longer uses ticketBranch's name-first order; it takes only the branches holding the newest own commit beyond HEAD (git branch --contains), one -> its worktree, several -> named. Slug-named branches play no part (own worktree already caught by ownTreeFor).
