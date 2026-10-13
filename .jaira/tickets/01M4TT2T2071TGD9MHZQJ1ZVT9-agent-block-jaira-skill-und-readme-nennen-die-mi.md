---
id: 01M4TT2T2071TGD9MHZQJ1ZVT9
title: "Agent-Block, jaira-Skill und README nennen die Milestone-Befehle, so wie sie die Tag-Befehle nennen"
status: in-progress
ready: true
creator: BeMuCa
goal: "Ein Agent und ein Mensch finden die Milestone-Befehle dort, wo sie die Tag-Befehle finden: im Agent-Block, im jaira-Skill und in der README."
context: |-
  Die Milestone-Befehle (jaira milestone ls/create/add/rm/set/rename, jaira logbook <milestone>, die Milestone Roadmap mit --after und --no-roadmap) stehen heute nur in docs/COMMANDS.md und in der --help.
  Der Agent-Block (core/board/announce.go, daraus wird der jaira-Abschnitt in CLAUDE.md und AGENTS.md) nennt sie nicht. Der jaira-Skill (.claude/skills/jaira/SKILL.md) nennt sie nicht. Die Befehlsliste der README ('## Commands') nennt sie nicht.
  Folge: Ein Agent weiss nicht, dass er Milestones anlegen, verbinden und umbenennen kann. Berk fragte am 2026-10-13, ob Claude das ueberhaupt kann.
  Vorbild sind die Tags, Commit a208d43b (79GEPW): ein Abschnitt im Skill ('Tags: read the vocabulary before you add to it'), 8 Zeilen im Block ('Read it before you tag anything'), dazu docs/COMMANDS.md.
  Berk am 2026-10-13: 'Wie haben wir das bei den anderen Befehlen gemacht? hier sollten wir es genauso machen. + Update der Readme natuerlich.'
  Baut auf E70WYN auf (8dce517), weil 'jaira milestone rename' dazugehoert.
definition-of-done: "core/board/announce.go: der Block nennt 'jaira milestone ls' als vor dem Anlegen zu lesen (wie jaira tags) und create/add/rm/set/rename samt --after fuer die Roadmap; ein Test in core/board prueft die Zeilen; CLAUDE.md und AGENTS.md dieses Repos sind neu erzeugt"
tags:
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-10-13T14:18:52Z
updated-at: 2026-10-13T17:05:29Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-2510012
claimed-at: 2026-10-13T16:56:10Z
---

# Agent-Block, jaira-Skill und README nennen die Milestone-Befehle, so wie sie die Tag-Befehle nennen

## Definition of Done

- [ ] core/board/announce.go: der Block nennt 'jaira milestone ls' als vor dem Anlegen zu lesen (wie jaira tags) und create/add/rm/set/rename samt --after fuer die Roadmap; ein Test in core/board prueft die Zeilen; CLAUDE.md und AGENTS.md dieses Repos sind neu erzeugt
- [ ] .claude/skills/jaira/SKILL.md hat einen Milestone-Abschnitt nach dem Muster des Tags-Abschnitts: Liste vor dem Anlegen lesen, die Befehle, Roadmap mit --after/--no-roadmap, rename, logbook <milestone>
- [ ] README.md '## Commands' nennt jeden Befehl, den 'jaira --help' listet, in einer Zeile (am 2026-10-13 fehlten: hook, links, milestone, note, restore, resume, tag, tags, whoami)
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/
- [ ] docs/COMMANDS.md beschreibt jeden Befehl, den 'jaira --help' listet (am 2026-10-13 fehlten: links, pull, release, roles, snapshot); ein Test vergleicht die Befehlsliste mit README und COMMANDS.md, damit kein neuer Befehl wieder fehlt

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] internal/cli/commands_doc_test.go: test that every command 'jaira --help' lists (newRoot().Commands(), IsAvailableCommand, minus cobra's help/completion) has a line in README '## Commands' and a row in docs/COMMANDS.md; run it red first
- [x] README.md '## Commands': one line each for hook, links, milestone, note, restore, resume, tag, tags, whoami, worded from the branch build's --help
- [x] docs/COMMANDS.md: rows for fetch, links, pull, release, roles, snapshot from --help; leave the milestone rename row as it is; test green
- [x] core/board/announce_test.go: failing test that the block names 'jaira milestone ls' as read-before-create and create/add/rm/set/rename, --after, --no-roadmap
- [x] core/board/announce.go: milestone entry (<= 8 lines) right after the jaira tags entry, modeled on a208d43b; test green
- [x] regenerate CLAUDE.md and AGENTS.md with the branch build's 'jaira update' on a copy without remote (HOME/JAIRA_HOME in scratch), copy both back, diff shows only the new lines
- [x] .claude/skills/jaira/SKILL.md: Milestones section after the Tags section - ls before create, the commands, roadmap --after/--no-roadmap, rename, logbook <name>
- [x] core/release/notes/J1ZVT9.md: one line telling users to run jaira update for the new block lines
- [~] go vet ./... and go test ./... -race green; commit with ticket file, docs(J1ZVT9)/feat(J1ZVT9)

## Progress
- **2026-10-13 16:54 · BeMuCa** — Berk 2026-10-13: Umfang festgelegt - README und docs/COMMANDS.md nennen JEDEN Befehl; Agent-Block und jaira-Skill bekommen nur die Milestone-Befehle, wie damals die Tags (a208d43b, 79GEPW). Grund fuer den kleinen Block: er wird in jede Agent-Sitzung geladen. Kommt in PR #74 (feat/milestone-roadmap); Arbeit im Worktree /home/berk/git/jAIra-J1ZVT9 auf diesem Branch.
- **2026-10-13 17:01 · BeMuCa** — pre-process: what the plan rests on.
- Regeneration: no go generate and no test writes this repo's CLAUDE.md/AGENTS.md. 'jaira update' rewrites the managed block (internal/cli/update.go:105, board.AnnounceInAgentFiles); 'jaira validate' warns on drift (NoteIsCurrent). A dry 'jaira update' on a no-remote copy of the base changed neither file, so the regenerated diff will be only the new lines.
- Copy, not the worktree itself: the CLI flushes ticket refs to origin after every command, and a worktree's .git is shared with the main repo, so 'git remote remove' on a cp of the worktree would hit the real config. The copy is a fresh 'git init' of the tracked files plus .jaira/.
- Gap found: docs/COMMANDS.md also has no row for 'jaira fetch' (the ticket's list named 5; a scan of the branch build's --help found 6). DoD 5 says 'jeden Befehl', so fetch is in.
- The doc test reads the command list from newRoot().Commands() filtered by IsAvailableCommand - the same filter cobra's usage template uses - and leaves out help and completion: cobra adds those at Execute and documents them itself; the ticket's own missing-lists left them out too. Top-level commands only.
- Block entry also names 'jaira logbook <name>': a milestone never leaves on its own, emptying it with rm leaves it standing.
- Base rebased onto 7cbc41de (coordinator: 363d5a94 failed go vet in internal/tui/returnto_test.go).
- **2026-10-13 17:05 · BeMuCa** — in-progress: what the diff does not say.
- Both new tests went red first with exactly the expected misses: README 9 (hook links milestone note restore resume tag tags whoami), COMMANDS.md 6 (fetch links pull release roles snapshot); announce test 8 missing substrings.
- Regenerated CLAUDE.md/AGENTS.md with the branch build's 'jaira update' on a 'git init' copy of the tracked files + .jaira/ (0 remotes, HOME/JAIRA_HOME in scratch); the diff there was exactly the 8 new block lines per file, then both files copied back.
- Every command in the SKILL.md section was run on that copy: 'Round One' -> round-one; set --after nosuch exit 5; a loop exit 3 with the cycle named; --after with --no-roadmap exit 2; logbook <name> exit 3 while a ticket is short of done; rename moved the file.
- Left alone on purpose (scope): README line 'jaira fetch' still says only tickets while its --help says tickets and milestones; 'jaira hook example' has no COMMANDS.md row (the test checks top-level commands only).
