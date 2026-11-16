---
id: 01M7K4V7Q92DJA15NGTKKCRS3W
title: Review- und Testing-Lane finden übernommene Fehler; neue Rolle jaira-role-bughunt
status: testing
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Ticket, das alle Lanes grün durchläuft, hat keine Fehler mehr, die ein externer Code-Review (Copilot) danach im selben Code findet"
context: |-
  Was falsch ist: Ein Ticket, das bestehenden Code aus einem anderen Projekt übernahm, lief durch critique, testing und review; das Review fand keine Mängel. Danach schrieb ein externes Copilot-Review auf dem PR 15 Hinweise, 14 davon echte Fehler (Lost Update bei parallelen Schreibzugriffen, Erfolgsmeldung trotz fehlgeschlagenem Schritt, verlorene Ereignisse, Ausgabe unter falscher Überschrift u. a.).
  Ursache 1: core/lane/builtin/40-review.md Zeile 51 fragt nur nach "defects introduced by this change". Fehler, die mit übernommenem oder verschobenem Code mitkommen, gelten als nicht eingeführt und werden nicht gesucht. Zeile 79 nennt none eine "normal, expected answer".
  Ursache 2: 40-review.md input-requires enthält outcome-what und outcome-resolves. Der Reviewer liest den Bericht des Implementierers, bevor er den Code sieht, und prüft durch diese Brille. Copilot hatte den Rahmen nicht.
  Ursache 3: core/lane/builtin/27-testing.md Pass 3 (Function) prüft nur den normalen Ablauf. Fehlerpfade (Befehl schlägt fehl, Abbruch mittendrin) und parallele Schreibzugriffe kommen nicht vor.
  Neu: eine Rolle jaira-role-bughunt (wie jaira-role-tester unter core/role/builtin/). Ein frischer Agent ohne Ticketkontext liest ganze Module (nicht nur den Diff) und sucht Fehler: Nebenläufigkeit, Fehlerpfade, ungewöhnliche Eingaben, Rechte. Jeder Fund wird gegen den Code verifiziert, bevor er gemeldet wird. Funde gehen als note an das Ticket oder als neue Tickets mit follows. Aufruf: /jaira-role-bughunt [ticket-id | pfad...]. Ausdrücklich kein Go-Befehl in jaira: jaira führt keine Modelle aus (Scope-Grenze), die Rolle ist ein Skill.
  Upstream gehört der Kollege (BeMuCa/jaira, Lane testing creator BeMuCa): vor dem Merge abstimmen.
definition-of-done: "40-review.md: Fehler in allem Code suchen, den der Diff enthält, auch übernommenem; berührte Module ganz lesen; Nebenläufigkeit, Fehlerpfade und ungewöhnliche Eingaben ausdrücklich genannt"
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-11-16T21:40:09Z
updated-at: 2026-11-16T22:05:11Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-18351
claimed-at: 2026-11-16T21:45:52Z
outcome-what: "review lane drops outcome-what/outcome-resolves from its input and hunts defects in every module the diff touches, copied or moved code included, naming concurrency, error paths, unusual input and permissions; testing pass 3 now forces a failure path (exit code checked) and two parallel writers; new role jaira-role-bughunt (fresh agent, whole modules, verified findings as notes or tickets); worker, dispatcher and teamlead run it before a code-changing ticket goes to review; release note core/release/notes/KCRS3W.md"
outcome-why: "a ticket taking over code from an older project passed every lane and review, then an external review found 14 real defects: the review looked only for defects introduced by the change, read the implementer's account first, and testing exercised only the happy path"
outcome-resolves: "DoD 1-5 met with proofs; DoD 6 release note written, the agreement with the colleague is the person's step"
review-summary: none
review-gaps: "none: notTickets now one package const shared by Diff and WorktreeDiff instead of a local in WorktreeDiff; nothing else to remove"
test-verdict: "pass: go test -count=1 ./... RC=0; TestDiffLeavesTicketFilesOut fails without the fix, passes with it; for-lane review payload of this ticket built by the new binary carries no ticket file (0 hits on the ticket id), lane and code diffs intact; ticket-only commit yields empty diff, RC=0; failure path: show --for-lane on an unknown ticket exits 5; Diff writes nothing, so no parallel-writer case; bughunt file recipe run on this branch lists the 13 branch files, no foreign ones"
round: 2
---

# Review- und Testing-Lane finden übernommene Fehler; neue Rolle jaira-role-bughunt

## Definition of Done

- [x] 40-review.md: Fehler in allem Code suchen, den der Diff enthält, auch übernommenem; berührte Module ganz lesen; Nebenläufigkeit, Fehlerpfade und ungewöhnliche Eingaben ausdrücklich genannt
  proof: core/lane/builtin/40-review.md:24 (whole modules, copied/moved code, four named lines)
- [x] 40-review.md: outcome-what und outcome-resolves aus input-requires entfernt; Lane-Schema-Tests und bestehende Boards (jaira update --lanes) laufen weiter
  proof: core/lane/builtin/40-review.md:9; go test ./... RC=0 (TestLanesShowPrintsFullContract, TestEveryBuiltinLaneIsAShippedVersion); jaira update --lanes on a copy of this board: 2 updated, 0 skipped, lanes show review Input: goal, definition-of-done, diff
- [x] 27-testing.md: Pass 3 löst mindestens einen Fehlerpfad aus (Exit-Code geprüft) und, wo geschrieben wird, zwei parallele Schreiber
  proof: core/lane/builtin/27-testing.md:36 and :41
- [x] core/role/builtin/jaira-role-bughunt/SKILL.md: frischer Agent ohne Ticketkontext, ganze Module, Funde verifiziert, Ausgabe als note oder Folgeticket; wird mit jaira update --roles installiert
  proof: core/role/builtin/jaira-role-bughunt/SKILL.md; core/role/role_test.go wantRoles; HOME=scratch jaira update --roles wrote skills/jaira-role-bughunt/SKILL.md
- [x] jaira-teamlead und jaira-worker rufen die Bughunt-Rolle vor der Übergabe an signoff auf, wenn das Ticket Code ändert
  proof: core/role/builtin/jaira-worker/SKILL.md:200, core/role/builtin/jaira-dispatcher/SKILL.md:268, core/role/builtin/jaira-teamlead/SKILL.md:196
- [ ] CHANGELOG-Eintrag; mit dem Kollegen abgestimmt, bevor es upstream geht

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-16 21:45 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/4208.sock
- **2026-11-16 21:46 · Alexander Sacharov** — decision: no open decisions for the person. (1) Bughunt runs after testing passes and BEFORE the move into review: review is the lane that moves into signoff, so 'before signoff' can only be held by the agent carrying the ticket before it hands review the work; a verified finding sends the ticket back to in-progress like any other loop. (2) Bughunt output: given a ticket id -> notes on that ticket ('bughunt: ...'), no satellite tickets (fold findings into the ticket in hand); given only paths -> one ticket per verified finding via jaira create (follows set by the branch as usual). (3) 'teamlead calls it': the teamlead never runs lanes itself, so the dispatcher (per-lane mode) and the worker (per-ticket mode) both run it, and the teamlead skill says so and refuses to forward a code-changing ticket to signoff without a bughunt note. (4) CHANGELOG = core/release/notes/KCRS3W.md, the repo's release-notes rule; releases belong to 'worker' here. (5) Removing outcome-* from review input: the four review outputs stay; the prompt line 'is anything claimed in the outcome not supported' goes, the person compares outcome and verdict at signoff. (6) DoD 6 'abgestimmt mit dem Kollegen' is a person's act — it is asked in the human lane at the end, not ticked by an agent.
- **2026-11-16 21:48 · Alexander Sacharov** — decision: bughunt is started as a subagent (Agent tool, subagent_type claude), not through spawn.sh in a Herdr tab. Reason: the hunt moves no lane, so run-lane.sh has nothing to wait for, and a subagent returns when it is done on every transport; its context is fresh either way. spawn.sh stays untouched.
- **2026-11-16 21:50 · Alexander Sacharov** — context: Hinweise auf das Kundenprojekt entfernt (Projektname, konkrete Fehlerbilder) - Alex' Regel: keine fremden Daten auf öffentlichen Refs. Inhaltlich unverändert.
- **2026-11-16 21:55 · Alexander Sacharov** — found out: lane and role files are listed by hash in core/shipped/shipped.sha256 (append-only); every edit to a built-in lane or role needs its new sha256 appended, or TestEveryBuiltinLaneIsAShippedVersion/TestEveryInstalledFileIsAShippedVersion fail and jaira update would treat older boards' copies as hand-edited. Appended by sha256sum at the end of each edit round. TestLanesShowPrintsFullContract pins the first review sentence verbatim, so the new hunting instruction starts as a second sentence. DoD 6: release note is core/release/notes/KCRS3W.md; the agreement with the colleague is the person's step and stays open for the human lane.
- **2026-11-16 21:55 · Alexander Sacharov** — critique: none. Let stand on purpose: the bug-hunt motivation paragraph is repeated in worker, dispatcher and teamlead — each skill is read alone, same as the critique-beside section already repeated in worker and dispatcher. Starting the hunt as an Agent subagent stands on the 21:48 decision note. The teamlead's grep for 'bughunt: done' also matches a line from an earlier round; accepted as a cheap tripwire, the worker/dispatcher rule is what actually runs it each round.
- **2026-11-16 21:59 · Alexander Sacharov** — bughunt: core/lane/builtin/40-review.md:9 + core/gitrepo/git.go:110 - dropping outcome-what/outcome-resolves from input-requires does not keep the implementer's account away from review: the diff payload is 'git show' of every ticket commit, and the ticket file rides in those commits, so outcome-what, outcome-why, outcome-resolves and every Progress note arrive inside 'diff' anyway. WorktreeDiff (git.go:211) excludes .jaira/tickets, Diff does not. ran: jaira show KCRS3W --for-lane review --json -> input has only goal+definition-of-done, but diff contains '+outcome-what:', '+outcome-why:', '+outcome-resolves:' and the decision notes
- **2026-11-16 21:59 · Alexander Sacharov** — bughunt: core/role/builtin/jaira-role-bughunt/SKILL.md:27 - 'git merge-base HEAD origin/HEAD' assumes origin is the repo the branch is cut from; in a fork setup (origin = the fork, upstream = the source) origin/HEAD lags and the hunter is handed hundreds of foreign files as 'the branch's files'. ran: in this worktree origin/HEAD -> origin/master is 207 commits behind master; the command lists ~300 files (logbook, web UI, cli...) while the branch's only commit touches 12
- **2026-11-16 21:59 · Alexander Sacharov** — bughunt: core/role/builtin/jaira-role-bughunt/SKILL.md:23-28 - the file list the hunter is told to read whole includes the ticket's own file under .jaira/tickets/ (it rides in the branch commits and shows in git status), so following the recipe reads goal, context, outcome and notes - the story line 23 forbids. ran: both commands list .jaira/tickets/01M7K4V7Q92DJA15NGTKKCRS3W-*.md
- **2026-11-16 21:59 · Alexander Sacharov** — bughunt: core/role/builtin/jaira-teamlead/SKILL.md:209 - 'jaira show <id> | grep "bughunt: done"' matches a closing line from any earlier pass; the worker and dispatcher send a ticket with findings back to in-progress and require a new hunt before it leaves again, so a skipped second hunt still passes the teamlead's check. read: worker SKILL.md 'The bug hunt' (Findings -> run the bug hunt once more) vs the unanchored grep
- **2026-11-16 21:59 · Alexander Sacharov** — bughunt: done - 4 findings in 6 modules; concurrency, error paths, input, permissions checked
- **2026-11-16 22:02 · Alexander Sacharov** — bughunt round 1 fixed: (1) outcome leaked into review through the ticket file committed beside the code — gitrepo.Diff now excludes .jaira/tickets like WorktreeDiff (shared const notTickets), TestDiffLeavesTicketFilesOut fails without it; (2)+(3) bughunt's file list now comes from commits whose message names the id (git log --grep), not a merge base with origin/HEAD (fork 207 commits behind gave ~300 foreign files), and leaves .jaira/tickets out; it no longer runs jaira show at all; (4) closing line is 'bughunt: done at <short sha>' and the teamlead compares it with the branch head, so a line from an earlier round no longer passes.
- **2026-11-16 22:03 · Alexander Sacharov** — critique round 2: read only the bughunt fixes. Found: git show with a pathspec prints nothing at all for a commit that touched only its ticket file (checked in a scratch repo, RC=0, empty output) — no code to judge, so accepted; documented on gitrepo.Diff (comment only) rather than adding a second git call per commit for a header line.
