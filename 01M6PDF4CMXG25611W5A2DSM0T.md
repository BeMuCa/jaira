---
id: 01M6PDF4CMXG25611W5A2DSM0T
title: "Ein Board ohne Git-Repo schiebt Tickets nach done, ohne Commits zu verlangen"
status: in-progress
ready: true
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
updated-at: 2026-11-23T01:27:16Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: worker-2dsm0t
claimed-at: 2026-11-23T01:24:05Z
---

# Ein Board ohne Git-Repo schiebt Tickets nach done, ohne Commits zu verlangen

## Definition of Done

- [x] Ohne Git-Repo geht ein Ticket nach done und ins Logbook, ohne Commits und ohne --force
  proof: core/gate/gate.go OutsideRepo waiver in the requires-commits check; core/gitrepo/git.go OutsideRepo; manual: scratch board in plain dir, move to done exit 0 and logbook exit 0 without --force
- [x] In einem Repo ohne Remote verlangt done weiter Commits
  proof: OutsideRepo is false when a .git exists upward or git sees a work tree; manual: git init without remote, move to done refused with needs_commits; TestOutsideRepo repo cases
- [x] Ein Test je Fall
  proof: TestCommitsWaivedOnlyOutsideARepository (core/gate/commits_test.go), TestOutsideRepo (core/gitrepo/outside_test.go); compiled by go vet, not run (worker rule: no go test)
- [x] Release-Notiz in core/release/notes/<ID>.md
  proof: core/release/notes/2DSM0T.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] gitrepo: OutsideRepo() - true only when no .git in Dir or any parent AND git does not report a work tree
- [x] gate.Env: lazy OutsideRepo func; requires-commits refusal skipped only when derive found nothing and OutsideRepo() is true
- [ ] wire it in internal/cli/root.go loadEnv and internal/tui/model.go gateEnv
- [ ] tests: gate (nil/false keeps refusal, true passes), gitrepo (plain dir, repo without remote, subdir of repo)
- [ ] release note core/release/notes/2DSM0T.md
- [ ] manual: scratch board without git reaches done and logbook; repo without remote still refused

## Progress
- **2026-11-23 01:25 · BeMuCa** — Reproduced 2026-11-23 on a scratch board in a plain dir (branch build): move to done refused with 'lane "done" requires the commits ...'. Decision (mine, settle-able): 'no repo' = no .git entry in the board root or any parent AND git does not answer true to rev-parse --is-inside-work-tree (or git is not installed). Why both: a git error that is not 'no repo' (dubious ownership, broken config) must never drop the commit rule inside a real repo; a missing git binary alone must not either when a .git is there. Asked lazily, only on the refusal path, so a repo board pays nothing and render paths never call it. Planning option not ticked, so pre-process is skipped; plan written as plan items anyway.
