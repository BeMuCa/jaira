---
id: 01M7K4V7Q92DJA15NGTKKCRS3W
title: Review- und Testing-Lane finden übernommene Fehler; neue Rolle jaira-role-bughunt
status: critique
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
updated-at: 2026-11-16T21:55:18Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-18351
claimed-at: 2026-11-16T21:45:52Z
outcome-what: "review lane drops outcome-what/outcome-resolves from its input and hunts defects in every module the diff touches, copied or moved code included, naming concurrency, error paths, unusual input and permissions; testing pass 3 now forces a failure path (exit code checked) and two parallel writers; new role jaira-role-bughunt (fresh agent, whole modules, verified findings as notes or tickets); worker, dispatcher and teamlead run it before a code-changing ticket goes to review; release note core/release/notes/KCRS3W.md"
outcome-why: "a ticket taking over code from an older project passed every lane and review, then an external review found 14 real defects: the review looked only for defects introduced by the change, read the implementer's account first, and testing exercised only the happy path"
outcome-resolves: "DoD 1-5 met with proofs; DoD 6 release note written, the agreement with the colleague is the person's step"
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
