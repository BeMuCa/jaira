---
id: 01M5MNWDKYJGC78B3EKMX2TGXG
title: "Ein Worker findet seinen Teamlead nach einem Neustart wieder, und eine Frage an die Person kommt als fertige Zeile"
status: testing
ready: true
creator: BeMuCa
goal: "Nach einem Crash oder Neustart des Teamleads erreichen die Nachrichten der Worker ihn weiter, ohne dass jemand den neuen Namen von Hand nachreicht; und eine Frage, die nur die Person per Kommando loesen kann, kommt mit der Zeile, die sie tippen muss."
context: |-
  Beobachtet 22.-23.10.2026 (Feedback Punkte 3 und 4).
  Nach einem Neustart bekommt jede Session einen neuen Namen. Worker behalten das beim Start uebergebene --parent und schicken ihre Berichte ins Leere, bis der Teamlead jedem den neuen Namen nachsendet.
  Zweiter Fall: ein Worker brauchte Zugangsdaten in einem Container; Agenten duerfen das zu Recht nicht kopieren. Die Frage ging an die Person, aber ohne die Zeile, die sie haette tippen koennen - eine Runde verloren. Claude Code kennt '! <befehl>' im Prompt der Person.
  Wo das Elternfeld gesetzt wird: spawn.sh (--parent), jaira-worker/SKILL.md und jaira-dispatcher/SKILL.md (Berichtsweg).
definition-of-done: Der Worker schlaegt die aktuelle Adresse seines Teamleads beim Senden nach (Session-Record oder Board) statt den Startwert einzufrieren; ein Test deckt den Namenswechsel ab
tags:
  - concurrency
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T15:25:44Z
updated-at: 2026-10-26T19:33:11Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-252449
claimed-at: 2026-10-26T19:10:49Z
outcome-what: "Added core/role/builtin/jaira-dispatcher/scripts/teamlead.sh: a teamlead records its messaging socket with --here every round, and a worker or dispatcher runs it with its --parent before every SendMessage to get the address to send to; the teamlead, worker and dispatcher skills say so, and a question that needs a command from the person carries a ready '! <command>' line that the teamlead passes on verbatim."
outcome-why: "After a restart or resume the teamlead has a new socket and name, the --parent a worker was started with names the dead one, and its reports went nowhere until someone handed each worker the new name; and a question that needed a command from the person came without the line to type, which cost a round."
outcome-resolves: "teamlead.sh prints the parent while its socket exists and the address the teamlead recorded once it is gone, covered by TestTeamleadFollowsARestart (restart from a second worktree, second live teamlead keeps its reports, both gone, --here without a socket); the '! <command>' rule is in the worker, dispatcher and teamlead SKILL.md; release note in core/release/notes/X2TGXG.md."
review-summary: none
round: 2
review-gaps: "reused fakeSocket (core/role/role_test.go) in TestTeamleadFollowsARestart instead of its own sock/mk closures, which duplicated it. Left: the git-common-dir lookup and the -e socket check that spawn.sh and reap.sh also carry - the shipped scripts never source each other (run-lane.sh only runs spawn.sh), and a shared file for two lines would be a new shipped script; the usage header, which matches spawn.sh and reap.sh; the test's git closure, since reapRig.git belongs to a rig this test does not build; the two uds: liveness lines in teamlead.sh (a helper is not shorter); the redundant [ -n \"$record\" ] before [ -f \"$record\" ]; the near-identical '! <command>' and teamlead.sh paragraphs in the worker and dispatcher SKILL.md, separate skills as before; 'wrote it there itself on start' in both SKILL.md, narrower than the every-round record but not false"
---

# Ein Worker findet seinen Teamlead nach einem Neustart wieder, und eine Frage an die Person kommt als fertige Zeile

## Definition of Done

- [x] Der Worker schlaegt die aktuelle Adresse seines Teamleads beim Senden nach (Session-Record oder Board) statt den Startwert einzufrieren; ein Test deckt den Namenswechsel ab
  proof: core/role/role_test.go TestTeamleadFollowsARestart
- [x] jaira-worker/SKILL.md und jaira-dispatcher/SKILL.md sagen: eine Frage, die nur die Person per Kommando loesen kann, enthaelt die fertige '! <befehl>'-Zeile
  proof: core/role/builtin/jaira-worker/SKILL.md, jaira-dispatcher/SKILL.md: '! <command>' paragraph
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] teamlead.sh in core/role/builtin/jaira-dispatcher/scripts/: '--here [repo]' writes uds:$CLAUDE_CODE_MESSAGING_SOCKET to <git common dir>/jaira/teamlead (exit 2 without a live socket); '<parent> [repo]' prints the parent while it is a uds: address whose socket exists, else the recorded address when its socket exists, else the parent unchanged
- [x] Go test in core/role/role_test.go: record A; A gone and B recorded (also from a second worktree) -> B; A alive and B recorded -> A; both gone -> A; --here without a socket -> exit 2
- [x] jaira-teamlead/SKILL.md: run 'teamlead.sh --here' first thing in every session
- [x] jaira-worker/SKILL.md and jaira-dispatcher/SKILL.md: before every SendMessage to the parent run 'teamlead.sh <session>' and send to what it prints; a changed address is noted as the new parent; the 'no other peer' rule names the recorded teamlead as the only replacement
- [x] jaira-worker/SKILL.md and jaira-dispatcher/SKILL.md: a question only the person can resolve by running a command carries the ready-to-type '! <command>' line
- [x] core/release/notes/X2TGXG.md; gofmt/vet/build; go test ./core/role -run <touched tests>; commit code together with the ticket file

## Progress
- **2026-10-26 19:15 · BeMuCa** — design (driver, pre-process skipped by ticket options): the record is one file, <git common dir>/jaira/teamlead, beside spawn.sh's jaira/panes/ - shared by every worktree, machine-local, never in git status. Parent first, record second: a parent whose uds: socket still exists keeps the reports, so a second live teamlead on the same clone never takes another's workers (the rule in both SKILL.md files). Only a gone socket is replaced, and only by the address the teamlead wrote itself on start - not a guess among peers. Rejected: a Go command over RepoStateDir (binary change and new CLI surface for one file); record-first (would route a live teamlead's reports to a second one); spawn.sh writing the record (a dispatcher calling spawn.sh would overwrite it).
- **2026-10-26 19:17 · BeMuCa** — Decision: TestCrossReferencesCarryThePrefix read 'scripts/teamlead.sh' as a bare /teamlead role reference; added an exemption for a path segment '/scripts' in the test rather than renaming the script.
- **2026-10-26 19:24 · BeMuCa** — critique: two findings, both in core/role/builtin/jaira-teamlead/SKILL.md. (1) Lines 165-169: 'teamlead.sh --here first thing in every session' sits in the Herdr spawn paragraph and misses the resume case - claude --resume keeps the session ID (claude --help: only --fork-session makes a new one) but gets a new socket, so a resumed teamlead counts it as already done. Move the call into the every-round routine under 'Read the board first, always' beside reap.sh (lines 31-36) and say every round. (2) Point 4 (lines 72-79) is the only hop that asks the person and recasts questions as AskUserQuestion options; add one sentence that a question carrying a '! <command>' line reaches the person with that line verbatim in the question text. Left standing, on the 19:15 note: the record in <git common dir>/jaira/teamlead beside spawn.sh's panes, parent first, no Go command; a lookup by session ID in ~/.claude/sessions was considered and is not better, because teamlead SKILL.md:174-177 supports starting a fresh teamlead, which a session ID never finds. On the 19:17 note: the /scripts exemption in TestCrossReferencesCarryThePrefix. For testing, unverified: liveness is [ -e socket ]; if a crashed claude leaves its socket file behind, the dead parent wins over the record.
- **2026-10-26 19:28 · BeMuCa** — critique (round 2): both 19:24 findings are repaired in 19d71add - teamlead.sh --here now sits in the every-round block beside reap.sh (jaira-teamlead/SKILL.md:31-48, resume named) and point 4 keeps a '! <command>' line verbatim in the question text (SKILL.md:84-86); release note line 1 says every round to match. Per the later-pass rule nothing else was read. The 19:24 'for testing, unverified' item (liveness is [ -e socket ], a stale socket file of a crashed claude would win over the record) is still for the testing lane.
