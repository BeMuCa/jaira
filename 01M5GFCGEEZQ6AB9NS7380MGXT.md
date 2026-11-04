---
id: 01M5GFCGEEZQ6AB9NS7380MGXT
title: "Ein Ticket in done raeumt seinen Worktree weg, sobald der Branch gepusht und sauber ist"
status: todo
ready: true
creator: berkc
assignee: berkc
goal: "Erreicht ein Ticket done, verschwindet sein Worktree (und Herdr-Tab), ohne auf den Merge des PR zu warten; der Branch bleibt, bis der PR gemergt ist."
context: |-
  Entscheidung von berkc am 22.10.2026: Worktrees werden bei done entfernt, nicht bei review.
  Heute (geprueft): reap.sh (core/role/builtin/jaira-dispatcher/scripts/reap.sh, YD4BE9) entfernt Tab, Worktree und lokalen Branch erst, wenn das Ticket done ist UND sein PR gemergt ist. Laeuft nur, wenn der Teamlead es pro Runde startet. AYBDJ2 (Backlog) betrifft nur die Arbeitskopie nach jaira logbook.
  Neu wegen 6TCTFV: PRs entstehen pro Milestone, nicht pro Ticket. Auf den Merge zu warten hiesse, Worktrees liegen bis zum Ende des Milestones herum.
  Sicherheitsregeln: nur entfernen, wenn git status im Worktree sauber ist (oder nur die eigene Ticket-Datei geaendert) und der Branch auf dem Remote denselben Stand hat; sonst stehen lassen und sagen warum. Branch erst nach Merge loeschen.
  Wie andere Sitzungen wissen, woran gearbeitet wird: claim (claimed-by/claimed-at auf dem Ref, 30-min-Lease), jaira next ueberspringt beanspruchte Tickets; C8V6DQ zeigt es auf dem Board.
definition-of-done: "Erreicht ein Ticket done und ist sein Branch gepusht und der Worktree sauber, wird der Worktree entfernt; der Branch bleibt; Test"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-22T00:15:13Z
updated-at: 2026-11-04T21:27:07Z
updated-by: berkc
claimed-by: berk-ca-1279800
claimed-at: 2026-11-04T21:24:25Z
---

# Ein Ticket in done raeumt seinen Worktree weg, sobald der Branch gepusht und sauber ist

## Definition of Done

- [ ] Steht das Ticket in done, ist jeder seiner Worktrees sauber (oder nur die eigene Ticket-Datei geaendert) und steht origin/<branch> auf demselben Commit wie der lokale Branch, schliesst reap.sh seine Tabs und entfernt den Worktree, ohne auf einen PR zu warten; der Branch bleibt; Test in core/role/reap_test.go
- [ ] Ist das Ticket nicht in done, der Worktree nicht sauber, der Branch nicht gepusht oder origin/<branch> auf anderem Stand, bleibt alles (Tabs, Worktree, Branch) und reap.sh druckt je Ticket eine Zeile mit dem Grund; Test je Fall
- [ ] Der lokale Branch wird erst geloescht, wenn sein eigener PR gemergt ist oder sein Stand im Default-Branch von origin liegt (Milestone-PR); reap.sh --dry-run zeigt 'would remove <worktree>' und 'would delete branch <branch>' und aendert nichts; Test
- [ ] Release-Note-Zeile in core/release/notes/80MGXT.md
- [ ] jaira-teamlead/SKILL.md beschreibt den neuen Zeitpunkt (Worktree bei done, Branch nach Merge); neue Hashes von reap.sh und SKILL.md in core/shipped/shipped.sha256 angehaengt, core/role- und core/shipped-Tests gruen

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-31 01:25 · berkc** — Maintainer 31.10.2026 (Punkt 7): Integrationsbranch ok, aber das Ziel ist: wird ein Ticket angenommen und steht in done, wird sein Branch gemerged und sein Worktree geloescht - 'ich will nicht den speicher mit worktrees voll machen die nicht geloescht wurden.' Merge bleibt Sache der Person (PR); das Aufraeumen des Worktrees nach done ist dieses Ticket.
- **2026-11-03 17:18 · berkc** — Maintainer 03.11.2026: Worktree-Entfernen ist auch deshalb sicher, weil die Commits im gepushten Branch liegen. Ergaenzung: das gemeinsame Review laeuft in einem eigenen Integrations-Worktree/-Branch je Milestone (Ticket-Branches hineingemergt, PR von dort); die Ticket-Worktrees verschwinden bei done.
- **2026-11-04 21:26 · berkc** — Entscheidungen (Worker, vor dem Plan): (1) 'gepusht' = refs/remotes/origin/<branch> zeigt auf denselben Commit wie refs/heads/<branch>; lokal gelesen, kein fetch/ls-remote - offline-sicher, Tests haben nur einen Fake-Remote git@github.com. git push setzt den Tracking-Ref ohnehin. (2) Milestone-PR: Branch gilt als gemergt, wenn sein eigener PR gemergt ist (wie heute) ODER sein Tip Vorfahre von refs/remotes/origin/HEAD ist. Squash-gemergter Milestone-PR -> Branch bleibt (sichere Seite, lokal ist nichts verloren, Branch liegt auf origin). (3) Ticket nicht in done -> 'waiting <T>: the ticket is in <lane>, not done' statt 'kept'. (4) Nur Worktrees mit spawn.sh-Record (git-common-dir/jaira/panes) werden geraeumt; Worktrees ohne Record bleiben ausser Scope; spawn.sh wird nicht angefasst (DSDPZ4). (5) Nach dem Worktree-Abbau ersetzt reap.sh die Tab-Records des Tickets durch einen Record ticket=/branch= ohne pane, damit die naechste Runde den Branch kennt und keine Tabs erneut schliesst.
