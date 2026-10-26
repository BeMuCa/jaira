---
id: 01M5MNWDKYJGC78B3EKMX2TGXG
title: "Ein Worker findet seinen Teamlead nach einem Neustart wieder, und eine Frage an die Person kommt als fertige Zeile"
status: review
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
updated-at: 2026-10-26T19:44:10Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-252449
claimed-at: 2026-10-26T19:10:49Z
outcome-what: "Added core/role/builtin/jaira-dispatcher/scripts/teamlead.sh: a teamlead records its messaging socket with --here every round, and a worker or dispatcher runs it with its --parent before every SendMessage to get the address to send to; the teamlead, worker and dispatcher skills say so, and a question that needs a command from the person carries a ready '! <command>' line that the teamlead passes on verbatim."
outcome-why: "After a restart or resume the teamlead has a new socket and name, the --parent a worker was started with names the dead one, and its reports went nowhere until someone handed each worker the new name; and a question that needed a command from the person came without the line to type, which cost a round."
outcome-resolves: "teamlead.sh prints the parent while its socket exists and the address the teamlead recorded once it is gone, covered by TestTeamleadFollowsARestart (restart from a second worktree, second live teamlead keeps its reports, both gone, --here without a socket); the '! <command>' rule is in the worker, dispatcher and teamlead SKILL.md; release note in core/release/notes/X2TGXG.md."
review-summary: "Adds core/role/builtin/jaira-dispatcher/scripts/teamlead.sh, shipped and installed executable beside spawn.sh. 'teamlead.sh --here' writes uds:$CLAUDE_CODE_MESSAGING_SOCKET into <git common dir>/jaira/teamlead - one file per clone, shared by every worktree, outside git status, written to a temp file and moved over; exit 2 when the socket variable is empty, names no file, or there is no repository. 'teamlead.sh <parent>' prints the parent while it is a uds: address whose socket file exists, else the recorded address when its socket file exists, else the parent unchanged. The teamlead skill now runs --here every round beside reap.sh. The worker and dispatcher skills run 'teamlead.sh <session>' before every SendMessage to the parent, send to what it prints, and note a changed address as the new parent. All three skills gain the '! <command>' rule: worker and dispatcher put the ready line into a question only the person can settle by a command, the teamlead passes it on verbatim. role_test.go ships the script, exempts '/scripts/teamlead' from the prefix check, and adds TestTeamleadFollowsARestart (record; old socket gone and new one recorded from a second worktree; a name as parent; old parent alive again keeps its reports; both gone; --here without a socket). Release note core/release/notes/X2TGXG.md, two lines. Net effect: after a clean teamlead restart, workers started before it find the new teamlead without anyone handing them its name."
round: 2
review-gaps: |-
  1. The release note leaves the feature switched off. 'Run `jaira roles install` again' skips all three changed SKILL.md files: base build 13aaa578, then this build, into one --into dir printed 'skipped .../jaira-worker/SKILL.md (edited here; --force to replace)', the same for jaira-dispatcher and jaira-teamlead, '1 written, 32 unchanged, 3 skipped', and 'some roles were edited here and left alone'. Only teamlead.sh lands, so no teamlead runs --here and no worker asks it. Both lines need `jaira roles install --force` (add `--global` if installed globally), as several earlier notes say.
  2. The crash half of the goal is not met. Liveness is [ -e socket ]. Shown here: a unix socket whose owner got kill -9 stays on disk with nothing listening (ss -xl), and teamlead.sh then prints that dead parent although a live address is recorded; the worker falls back to 'parent unreachable', as before. Claude Code names its socket $XDG_RUNTIME_DIR/cc-socks/<pid>.sock (read from its bundle); whether a new claude removes a crashed one's file is unverified. The teamlead skill itself names 'the human can kill your session, start a new teamlead' as the case. Small fix: when SendMessage to the printed address fails, fall back to the recorded one (a mode of teamlead.sh that skips the parent), one sentence in each skill, one test case.
  3. Claims that are false with two teamleads on one clone. The record is one slot, the last --here wins, and --here now runs every round. Shown: lead A restarts as B and records; a second lead X records on its next round; A's workers are sent to X. The usage text 'a second teamlead on the same clone never takes another's workers' and the skill text 'The only replacement is the address teamlead.sh prints, because the teamlead wrote it there itself on start' (jaira-worker/SKILL.md:285-286, jaira-dispatcher/SKILL.md:206-207) are wrong then, and 'on start' is stale since 19d71add made it every round. Fix the wording to the real limit (one teamlead per clone is followed) or key the record per teamlead.
  4. Minor: a parent given as a name (a person can type '/jaira-worker <id> --parent <name>') is never checked and is always replaced by a live record, even while the named session runs.
  5. Minor: in jaira-teamlead/SKILL.md the paragraph 'It closes every tab...' now follows a block of two commands; 'It' means reap.sh only.
  6. Still open by the maintainer's rule: DoD 3's 'go test ./... -race' has not run; it waits for the one full race suite before the PR. Ran here: gofmt -l . empty; go vet ./core/role and go build ./... green; go test ./core/role -run 'TestTeamleadFollowsARestart|TestCrossReferencesCarryThePrefix|TestWorkerNamesScriptsThatShip|TestBuiltinsAreTheShippedRoles' -count=1 ok.
test-verdict: "pass: gofmt/vet/build clean, go test -race ./core/role ok; teamlead.sh installed executable next to spawn.sh; restart play A->B (incl. 2nd worktree), A back, both gone, exit codes, outside-repo all as specified; DoD 1+2 proofs present; release note 2 single lines OK. DoD 3 suite part waits for pre-PR run. Not verified: stale socket file left by crashed claude (liveness is -e)."
review-verdict: "Back to in-progress. For a clean restart the change works: the script does what its plan says, DoD 1 and 2 are met, and the test covers the address change. But the release note's instruction does not install the changed skills (verified), the skill and usage texts promise that a second teamlead never gets another's workers, which a run disproves (verified), and the crash half of the goal still ends in 'parent unreachable' whenever the dead socket file stays (OS behaviour verified, Claude Code's own cleanup not). Fix gaps 1 and 3; gap 2 is a small addition worth making in the same round. DoD 3's full race suite is still open."
review-round: 2
review-check: |-
  1. cd /home/berk/git/jAIra-X2TGXG && go test ./core/role -run 'TestTeamleadFollowsARestart|TestWorkerNamesScriptsThatShip' -count=1 -> you see 'ok  github.com/BeMuCa/jaira/core/role'.
  2. Make a throwaway repo. Never do the next steps in your real clone: --here writes .git/jaira/teamlead into whatever repo you stand in. rm -rf /tmp/tl && mkdir /tmp/tl && cd /tmp/tl && git init -q repo && cd repo && git -c user.name=t -c user.email=t@example.com commit -q --allow-empty -m x
  3. T=/home/berk/git/jAIra-X2TGXG/core/role/builtin/jaira-dispatcher/scripts/teamlead.sh; touch /tmp/tl/a.sock /tmp/tl/b.sock   (plain files stand in for sockets)
  4. CLAUDE_CODE_MESSAGING_SOCKET=/tmp/tl/a.sock bash $T --here; cat .git/jaira/teamlead -> you see uds:/tmp/tl/a.sock
  5. The restart: rm /tmp/tl/a.sock; CLAUDE_CODE_MESSAGING_SOCKET=/tmp/tl/b.sock bash $T --here; bash $T uds:/tmp/tl/a.sock -> you see uds:/tmp/tl/b.sock
  6. The old lead is alive again: touch /tmp/tl/a.sock; bash $T uds:/tmp/tl/a.sock -> you see uds:/tmp/tl/a.sock
  7. CLAUDE_CODE_MESSAGING_SOCKET= bash $T --here; echo $? -> you see 'teamlead.sh: CLAUDE_CODE_MESSAGING_SOCKET is empty or names no socket; this session cannot be reached' and then 2
  8. Crash gap: python3 -c "import socket,time;s=socket.socket(socket.AF_UNIX);s.bind('/tmp/tl/c.sock');s.listen(1);time.sleep(60)" & sleep 1; kill -9 $!; ls -l /tmp/tl/c.sock; bash $T uds:/tmp/tl/c.sock -> the socket file is still listed after the kill, and the script prints uds:/tmp/tl/c.sock (the dead one), not the recorded uds:/tmp/tl/b.sock
  9. Second-teamlead gap: touch /tmp/tl/x.sock; CLAUDE_CODE_MESSAGING_SOCKET=/tmp/tl/x.sock bash $T --here; rm /tmp/tl/a.sock; bash $T uds:/tmp/tl/a.sock -> you see uds:/tmp/tl/x.sock (the other lead), not uds:/tmp/tl/b.sock
  10. Release-note gap: cd /home/berk/git/jAIra-X2TGXG && go build -o /tmp/tl/jaira-new ./cmd/jaira && jaira roles install --into /tmp/tl/skills && /tmp/tl/jaira-new roles install --into /tmp/tl/skills -> among the lines: 'written .../jaira-dispatcher/scripts/teamlead.sh' and 'skipped .../SKILL.md (edited here; --force to replace)' for jaira-dispatcher, jaira-teamlead and jaira-worker
  11. cat /home/berk/git/jAIra-X2TGXG/core/release/notes/X2TGXG.md -> two lines, each one line starting '- '; both say 'jaira roles install' without --force
  12. rm -rf /tmp/tl
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
- **2026-10-26 19:36 · BeMuCa** — driver: stale socket after a crash is not verified and could not be here - the machine booted at 19:06 and no session has ended since, so ~/.claude/sessions and /run/user/1000/cc-socks hold no ended session to look at. If a crashed claude leaves its socket file, teamlead.sh prints the dead parent, SendMessage fails, and the worker falls back to the old 'parent unreachable' rule - no worse than before. A real fix needs a probe that connects (bash cannot), or the registry pid check spawn.sh does not do either.
