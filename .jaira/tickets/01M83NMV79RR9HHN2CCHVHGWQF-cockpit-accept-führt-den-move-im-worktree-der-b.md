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
updated-at: 2026-11-23T07:52:26Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-70439
claimed-at: 2026-11-23T07:42:00Z
outcome-what: "Cockpit-Accept fuehrt den Move im Worktree der Branch aus, die die Commits des Tickets traegt (acceptTreeFor + ticketHere carried); mehrdeutig: kein Move, Kandidaten genannt. jaira move --to done nennt Branches mit Commits der ID, die HEAD nicht erreicht (gitrepo.BranchesNamingTicket, gate.Env.CommitsElsewhere)."
outcome-why: "Accept lief immer im Haupt-Checkout und fand Commits auf anders benannten Branches (Nachbesserung auf feat/<anderes-ticket>) nicht; die Meldung sagte nur 'commit it', obwohl committet war."
outcome-resolves: "VHGWQF: Accept-Fehler 'git has none for <ID>' bei Arbeit auf fremder Branch"
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
