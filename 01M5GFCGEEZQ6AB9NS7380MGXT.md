---
id: 01M5GFCGEEZQ6AB9NS7380MGXT
title: "Ein Ticket in done raeumt seinen Worktree weg, sobald der Branch gepusht und sauber ist"
status: backlog
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
updated-at: 2026-10-31T01:25:09Z
updated-by: berkc
---

# Ein Ticket in done raeumt seinen Worktree weg, sobald der Branch gepusht und sauber ist

## Definition of Done

- [ ] Erreicht ein Ticket done und ist sein Branch gepusht und der Worktree sauber, wird der Worktree entfernt; der Branch bleibt; Test
- [ ] Ist der Worktree nicht sauber oder der Branch nicht gepusht, bleibt alles und die Ausgabe nennt den Grund; Test
- [ ] Der Branch wird geloescht, sobald der (Milestone-)PR gemergt ist; reap.sh --dry-run zeigt beides
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

