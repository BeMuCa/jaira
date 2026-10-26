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
updated-at: 2026-10-26T19:53:38Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-252449
claimed-at: 2026-10-26T19:10:49Z
outcome-what: "Added core/role/builtin/jaira-dispatcher/scripts/teamlead.sh: a teamlead records its messaging socket with --here every round; a worker or dispatcher runs 'teamlead.sh <parent>' before every SendMessage (a uds: parent whose socket file is gone is replaced by the record, anything else comes back unchanged) and 'teamlead.sh --recorded' when a send fails; the teamlead, worker and dispatcher skills say so, and a question that needs a command from the person carries a ready '! <command>' line the teamlead passes on verbatim."
outcome-why: "After a restart, resume or crash the teamlead has a new socket and name, the --parent a worker was started with names the dead one, and its reports went nowhere until someone handed each worker the new name; and a question that needed a command from the person came without the line to type, which cost a round."
outcome-resolves: "A worker started under a teamlead that later restarted reaches the new session through the record (TestTeamleadFollowsARestart: restart recorded from a second worktree, live parent kept, name passed through, --recorded beside a stale socket file, record gone); the record is one slot per clone and the skills say so; the '! <command>' rule is in the worker, dispatcher and teamlead SKILL.md; release note core/release/notes/X2TGXG.md tells the reader to run 'jaira roles install --force'."
review-summary: "none"
round: 3
review-gaps: "61be772b only. Removed: teamlead.sh's closing case printed $arg in three places (uds-live branch, catch-all name branch, recorded || echo); now one guarded line (a uds: parent whose socket is gone and a live record -> exit 0) and one echo, same output and exit code in 13 old-vs-new cases (live/gone parent, name, --recorded, no/live/dead/garbage record, outside a repo). In TestTeamleadFollowsARestart the --recorded success check now uses the existing want helper instead of its own run/compare. Considered and left: recorded() (two callers, not a forwarder); the --recorded block (recorded || exit 1; exit 0 is explicit; recorded; exit would lean on set -e); [ -n \"$record\" ] inside recorded() is redundant with [ -f ] but was moved, not added; the exit-1 --recorded check beside the exit-2 --here loop (two shapes, a helper for two uses is not cheaper); the second (jaira note <id> \"parent session: <new>\") in each SKILL.md paragraph repeats the one eight lines up but keeps the command next to the step an agent acts on; the worker/dispatcher wording is duplicated by design (each skill installs alone); the release-note --force lines are needed (review gap 1)."
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
- **2026-10-26 19:44 · BeMuCa** — review: back to in-progress. (1) core/release/notes/X2TGXG.md says 'Run jaira roles install again', but an upgrade without --force skips jaira-worker, jaira-dispatcher and jaira-teamlead SKILL.md (edited here; --force to replace) - only teamlead.sh lands, so the feature stays off; say 'jaira roles install --force' (add --global if installed globally). (2) Two teamleads on one clone: the record is one slot and --here runs every round, so a restarted lead's workers go to whichever lead recorded last (reproduced with plain files); the usage text 'never takes another's workers' and the skill text 'because the teamlead wrote it there itself on start' (worker SKILL.md:285-286, dispatcher SKILL.md:206-207) are false then - correct them to the real limit or key the record. (3) Crash: a unix socket whose owner got kill -9 stays on disk, so teamlead.sh prints the dead parent over a live record; suggest falling back to the recorded address when SendMessage to the printed one fails, plus a test case. Details in review-gaps.
- **2026-10-26 19:45 · BeMuCa** — driver, after review round 1 - fixes for in-progress round 3: (1) release note says 'jaira roles install --force' (add --global if installed globally): plain install skips every changed SKILL.md, core/role/install.go writeFile. (2) crash: a kill -9'd socket file stays on disk, so add 'teamlead.sh --recorded [repo]' (prints the recorded address when its socket file exists, else nothing, exit 1) and make the worker/dispatcher skill say: when SendMessage to the printed address fails, run --recorded and send there if it differs, noting it as the new parent. (3) one slot per clone, last --here wins: fix the wording in usage and both skills to that real limit (with two teamleads on one clone an orphaned worker goes to whichever recorded last) and replace 'on start' by 'every round'; keying per teamlead is not possible, a fresh teamlead has no identity in common with the old one. (4) a parent that is not uds: comes back unchanged from '<parent>' mode - the --recorded fallback covers a dead name. (5) teamlead SKILL.md: 'It closes every tab' -> name reap.sh.
- **2026-10-26 19:48 · BeMuCa** — critique (round 3): all five review gaps (19:44 review, 19:45 driver note) are answered in 61be772b, by shape. (1) release note says 'jaira roles install --force' with --global (both flags exist, internal/cli/roles.go:202,205). (2) teamlead.sh --recorded is the fallback after a failed SendMessage, in worker SKILL.md:283-292 and dispatcher SKILL.md:204-213, with a test case for record found and record gone. (3) usage and both skills now state the one-slot limit and say 'every round'; no 'never takes another's workers' or 'on start' is left. (4) a name parent comes back unchanged; spawn.sh:107-120 always turns a name into uds:, so only a typed name reaches that path, and --recorded covers it when dead. (5) teamlead SKILL.md says 'reap.sh closes'. Per the later-pass rule nothing else was read. Left standing on the 19:15, 19:17 and 19:24 notes. DoD 3's full race suite stays open by the maintainer's rule.
