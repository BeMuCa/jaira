---
id: 01M4TT2T2071TGD9MHZQJ1ZVT9
title: "Agent-Block, jaira-Skill und README nennen die Milestone-Befehle, so wie sie die Tag-Befehle nennen"
status: critique
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
updated-at: 2026-10-13T17:36:13Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: j1zvt9-babysitter
claimed-at: 2026-10-13T17:36:07Z
outcome-what: "Agent block (core/board/announce.go:61-68) gained an 8-line milestone entry after the jaira tags entry: 'jaira milestone ls' read before creating, create/add/rm/rename, 'jaira logbook <name>', 'jaira milestone create|set <name> --after <a,b>' and 'jaira milestone set <name> --no-roadmap'; CLAUDE.md and AGENTS.md regenerated with 'jaira update'. .claude/skills/jaira/SKILL.md gained a Milestones section modeled on the Tags section, docs/AGENTS.md one beside 'Before tagging anything'. README '## Commands' gained 13 lines and docs/COMMANDS.md 9 rows, so both name every command 'jaira --help' lists, help and completion included. Tests: TestAgentNoteNamesTheMilestoneCommands and TestEveryCommandIsDocumented (CRLF-safe). Release note core/release/notes/J1ZVT9.md."
outcome-why: "The milestone commands were only in --help and COMMANDS.md, so an agent never learned it could create, fill, rename or place milestones; Berk asked whether Claude can do that at all. README and COMMANDS.md had drifted from the binary (9 and 6 commands missing) with nothing tying them to it. Review pass 1 found the doc test would fail on windows-latest (unpinned docs checked out with CRLF) and that a bare 'set' in the block reads as 'jaira set'."
outcome-resolves: "DoD 1: the block names ls as read-before-create plus create/add/rm/set/rename and --after, pinned by TestAgentNoteNamesTheMilestoneCommands; CLAUDE.md/AGENTS.md carry exactly those lines (validate: agent_note_stale []). DoD 2: SKILL.md section covers ls first, the commands, roadmap --after/--no-roadmap, rename, logbook <name>. DoD 3 and 5: TestEveryCommandIsDocumented builds the list exactly as cobra's usage template does (available commands plus help, help and completion added as Execute would) and fails for any without a README line or COMMANDS.md row; it passes on a CRLF copy and fails when a line is removed. DoD 4: go vet and go test ./... -race green; release note written."
review-summary: "Two commits. docs: README '## Commands' gains 9 lines, docs/COMMANDS.md 7 rows (links, fetch, pull, release, snapshot, roles list/install), and internal/cli/commands_doc_test.go fails for any top-level command of a fresh root with no README line or COMMANDS.md row. feat: the agent block (core/board/announce.go) gains an 8-line milestone entry after the tags entry, pinned by TestAgentNoteNamesTheMilestoneCommands; CLAUDE.md/AGENTS.md regenerated; .claude/skills/jaira/SKILL.md gains a Milestones section; release note J1ZVT9.md. Content checked against --help and source: accurate."
review-gaps: "1. DEFECT internal/cli/commands_doc_test.go:30 parses README on '\\n' but README.md and docs/COMMANDS.md have no eol=lf pin; on windows-latest (core.autocrlf) the Cut fails - reproduced with a CRLF copy: 'README.md has no ## Commands section'. 2. DoD 3/5 say every command --help lists; help and completion are listed but undocumented and skipped by the test. 3. announce.go: bare '`set --no-roadmap`' reads as the ticket command 'jaira set'. 4. README hook print line calls it the hook; --help says it prints a snippet. 5. release note says agents 'then read' ls - the block only tells them to. 6. README never names 'jaira logbook <milestone>'. 7. docs/AGENTS.md has 'Before tagging anything' (added by 79GEPW in 9236dcb3) but nothing on milestones. 8. test covers top-level commands only; 'jaira hook example' has no COMMANDS.md row (pre-existing)."
test-verdict: "pass: go build/vet + go test ./... -race -count=1 green at 17eae0d2 (29 ok, exit 0); DoD 1-5 verified in the tree; branch build's --help lists 40 commands and README/COMMANDS.md miss only cobra's help and completion; both new tests fail when a README line, a COMMANDS.md row or the block's milestone ls is removed; fresh no-remote 'jaira init' renders the 8 milestone lines into CLAUDE.md and AGENTS.md; on a no-remote copy of this board validate reports agent_note_stale [] and update leaves both files unchanged and lists the J1ZVT9 note"
review-verdict: "fail: windows-latest CI would go red on the new doc test (gap 1); back to in-progress with gaps 1-7"
review-check: "pending: written on the next review pass"
---

# Agent-Block, jaira-Skill und README nennen die Milestone-Befehle, so wie sie die Tag-Befehle nennen

## Definition of Done

- [x] core/board/announce.go: der Block nennt 'jaira milestone ls' als vor dem Anlegen zu lesen (wie jaira tags) und create/add/rm/set/rename samt --after fuer die Roadmap; ein Test in core/board prueft die Zeilen; CLAUDE.md und AGENTS.md dieses Repos sind neu erzeugt
  proof: core/board/announce.go:61-68; TestAgentNoteNamesTheMilestoneCommands (core/board/announce_test.go:169); CLAUDE.md:57-64 and AGENTS.md:30-37 regenerated by 'jaira update' of this branch on a no-remote copy
- [x] .claude/skills/jaira/SKILL.md hat einen Milestone-Abschnitt nach dem Muster des Tags-Abschnitts: Liste vor dem Anlegen lesen, die Befehle, Roadmap mit --after/--no-roadmap, rename, logbook <milestone>
  proof: .claude/skills/jaira/SKILL.md:204 '## Milestones: read the list before you start one' - ls first, create/add/rm/rename, list --milestone, roadmap --after/--no-roadmap, logbook <name>, restore
- [x] README.md '## Commands' nennt jeden Befehl, den 'jaira --help' listet, in einer Zeile (am 2026-10-13 fehlten: hook, links, milestone, note, restore, resume, tag, tags, whoami)
  proof: README.md:840-876 - 13 new lines incl. help and completion; TestEveryCommandIsDocumented (internal/cli/commands_doc_test.go:21) checks every command --help lists, help and completion included, CRLF-safe
- [x] go test ./... -race gruen; Release-Notiz unter core/release/notes/
  proof: go vet ./... && go test ./... -race -count=1: 29 packages ok, exit 0 (after review pass 1); core/release/notes/J1ZVT9.md
- [x] docs/COMMANDS.md beschreibt jeden Befehl, den 'jaira --help' listet (am 2026-10-13 fehlten: links, pull, release, roles, snapshot); ein Test vergleicht die Befehlsliste mit README und COMMANDS.md, damit kein neuer Befehl wieder fehlt
  proof: docs/COMMANDS.md:60 links, :69 help, :131 completion, :156-159 fetch/pull/release/snapshot, :180-181 roles list/install (fetch was missing too); TestEveryCommandIsDocumented checks README and COMMANDS.md

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
- [x] go vet ./... and go test ./... -race green; commit with ticket file, docs(J1ZVT9)/feat(J1ZVT9)
- [x] review pass 1: commands_doc_test.go normalises \r\n, mirrors --help exactly (help + completion); README/COMMANDS.md lines for help and completion
- [x] review pass 1: block names 'jaira milestone set'/'create|set' in full; README hook print + logbook <milestone> lines; release note wording; docs/AGENTS.md milestone section; regenerate CLAUDE.md/AGENTS.md; suite green; commit

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
- **2026-10-13 17:11 · BeMuCa** — critique pass 1 (whole diff, both commits): nothing found.
- Simpler shape: the doc test is one loop over newRoot().Commands() with two regexes; a parsed --help would be more fragile and say the same. The announce test copies TestAgentNoteNamesTheConversationalMode (core/board/announce_test.go).
- Fit: reading ../../README.md follows core/lane/shipped_test.go; the block entry sits beside the tags entry it was modeled on (a208d43b) and is the same length, 8 lines; the SKILL section follows 'Tags: read the vocabulary...'. No existing doc-vs-binary test to reuse (only market/shipped tests read repo docs).
- Speculative: none. Right place: block text in core/board, the command-list test in internal/cli where newRoot lives.
- Let stand: the SKILL section (49 lines) is longer than Tags (33) because the roadmap flags and the exit codes an agent branches on (5 unknown predecessor, 3 loop/existing name) have no shorter home.
- **2026-10-13 17:12 · BeMuCa** — optimize: nothing to remove.
- Duplication: no helper in the repo lists the root's commands (no other IsAvailableCommand) or reads repo files for tests beyond core/lane/shipped_test.go:16, which inlines its path the same way.
- Dead code / fluff: the read closure has two callers, QuoteMeta'd name is used by both regexes; the announce test is one table loop.
- Cost: two regex compiles per command, about 45 commands, test-only - not worth hoisting.
- **2026-10-13 17:17 · BeMuCa** — testing: what was run besides the suite.
- Mutation in a scratch copy: deleting README 'jaira whoami' and the COMMANDS.md 'jaira snapshot' row turned TestEveryCommandIsDocumented red with exactly those two; renaming 'jaira milestone ls' in announce.go turned TestAgentNoteNamesTheMilestoneCommands red.
- The 'note names X which is not in blocked-by' and 'cannot leave the backlog' warnings from validate on the board copy are other tickets' and predate this branch.
- **2026-10-13 17:29 · BeMuCa** — review pass 1 -> back to in-progress. Must fix: commands_doc_test.go reads README/COMMANDS.md and cuts on '\n'; neither file is pinned eol=lf (.gitattributes pins only embedded files), so the windows-latest job (autocrlf) fails - reproduced on a CRLF copy. Fix in the test: normalise \r\n before parsing, rather than pinning two docs files' checkout. Also take, cheap and inside the goal: (a) test mirrors --help exactly - InitDefaultHelpCmd/InitDefaultCompletionCmd and the usage template's 'IsAvailableCommand or help' - plus README/COMMANDS.md lines for help and completion, so DoD 3/5 hold literally; (b) block: write 'jaira milestone set' / 'jaira milestone create|set' out, bare 'set' reads as 'jaira set'; (c) README hook print = prints a snippet; (d) release note: 'tells your agents to read'; (e) README line for 'jaira logbook <milestone>'; (f) docs/AGENTS.md section next to 'Before tagging anything'. Leave: subcommand coverage (hook example row) - outside the DoD, report as a gap.
