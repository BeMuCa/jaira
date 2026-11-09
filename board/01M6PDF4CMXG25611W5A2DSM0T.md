---
id: 01M6PDF4CMXG25611W5A2DSM0T
title: "Ein Board ohne Git-Repo schiebt Tickets nach done, ohne Commits zu verlangen"
status: backlog
ready: false
creator: BeMuCa
goal: "Liegt das Board in keinem Git-Repository, verlangt die Lane done keine Commits; in einem Repository bleibt die Regel wie heute."
context: |-
  Was schief ist: in einem Ordner ohne Git laufen init, create, move und dod (ausprobiert 2026-11-02). Der Move nach done scheitert:
  'lane "done" requires the commits that carry this change, and git has none for P8QHNY yet'.
  Die Auswege der Meldung (--commits, jaira set commits=) sind ohne Git sinnlos. Ohne done kein Logbook.
  Gewollt (Berk, 2026-11-05): ein Solo-Nutzer ohne Repo kann jaira ganz benutzen. Mit Repo, auch rein lokal ohne Remote, bleibt die Commit-Pflicht.
  Wo: die Gate-Pruefung fuer Commits in core/gate/ (grep 'requires the commits'). core/gitrepo kennt ErrNoGit.
definition-of-done: "Ohne Git-Repo geht ein Ticket nach done und ins Logbook, ohne Commits und ohne --force"
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-11-05T17:52:51Z
updated-at: 2026-11-05T17:52:51Z
---

# Ein Board ohne Git-Repo schiebt Tickets nach done, ohne Commits zu verlangen

## Definition of Done

- [ ] Ohne Git-Repo geht ein Ticket nach done und ins Logbook, ohne Commits und ohne --force
- [ ] In einem Repo ohne Remote verlangt done weiter Commits
- [ ] Ein Test je Fall
- [ ] Release-Notiz in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

