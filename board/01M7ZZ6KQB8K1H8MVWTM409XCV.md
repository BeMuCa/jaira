---
id: 01M7ZZ6KQB8K1H8MVWTM409XCV
title: "Accept und Give back im Cockpit melden der Session eine Zeile, damit der Teamlead sofort weitermacht"
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Nimmt die Person im Cockpit an oder gibt zurück, erfährt die Session davon in derselben Minute und kann logbook, PR-Merge oder den nächsten Worker starten, ohne dass die Person es tippt."
context: |-
  Was heute fehlte (21.11.): Alex hat PDZ444 im Cockpit angenommen; die Teamlead-Session bekam nichts davon mit und bat weiter um die Annahme-Zeile.
  Das Cockpit läuft als Plugin in der Session und legt schon heute Zeilen hinein (die gekürzten Worker-Berichte: [shortened by jaira-cockpit] ...).
  Gewünscht: nach erfolgreichem Accept/Give back (core/mod/builtin/jaira-cockpit/hooks/register.tsx, accept/giveBack/acceptFromBand) eine Zeile in die Session, z. B. "jaira-cockpit: Alex accepted PDZ444 (signoff -> done)" bzw. "gave back PDZ444 -> in-progress: <Grund kurz>". Eine Zeile, kein Bericht — Kosten ~50 Token.
  Nur bei Erfolg des Moves; ein already in (gleichzeitiger Klick, MFYM1A) meldet nichts.
definition-of-done: Accept und Give back im Cockpit legen nach erfolgreichem Move genau eine Zeile in die Session (Test)
tags:
  - tui
blocked-by: []
related: []
commits:
  - e733a7ae439026940eee87a52c20d795bfbb048d
  - 8694f3345ada512b7ad9688852f11be1f3ec85bd
  - 1d1990d5b6f54195fe33f39deb1adddf85b2d8a4
created-at: 2026-11-21T21:11:37Z
updated-at: 2026-11-21T21:44:07Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-3494777
claimed-at: 2026-11-21T21:13:34Z
outcome-what: "Accept und Give back im Cockpit legen nach einem Move, der durchging, genau eine Zeile per $.prompt.submit in die Session: angenommen (mit dem, was Accept danach startete) bzw. zurueckgegeben (an wen, Grund einzeilig gekuerzt)."
outcome-why: "Die Teamlead-Session erfuhr nichts von einem Klick im Cockpit und bat weiter um die Annahme-Zeile (PDZ444, 21.11.)."
outcome-resolves: 409XCV
review-summary: "After a move that went through, Accept and Give back in the cockpit now hand this session one line via $.prompt.submit (origin plugin, the same channel decide() already uses). Accept sends it after afterAccept has finished, and afterAccept now returns what it said (PR or landing started, why there is none, or 'next step failed'), so the line names what is already running and the teamlead does not start a second PR. Give back sends it in handedTo once it is known who carries the work (old worker, new worker, nobody), with the person's reason cut to 120 code points (inLine in texts.ts: whitespace collapsed via board.ts oneLine, cut by Array.from so no emoji is split). A refused dry-run, a refused move and 'already in' return before the submit, so nothing is sent. Origin plugin is neither a person's prompt (no cache guard, no thread) nor an agent's (no shortening). In an agent's session the band does not exist. Release note and shipped.sha256 are included."
review-gaps: |-
  note: the line starts a model turn in any person session that runs the cockpit, not only a teamlead's - that is the goal, but in a plain session nobody leads, the turn costs tokens and does nothing useful; a later option is to send only when the session is a teamlead's
  note: void $.prompt.submit(...) drops a rejection without a word (register.tsx:618, :1201, same as the older decide() at :1041) - if submit fails, the session never hears and the cockpit does not say so either; fix: .catch(() => $.ui.toast(...))
  note: pre-existing: cuts by UTF-16 index are still in digest.ts:184, threads.ts:43, board.ts:782 - they could reuse inLine
  note: checked only in the engine harness (claude plugin test, 296 pass), not with a real press in a live session
test-verdict: "pass: round 2 - claude plugin test 296/296 (RC=0), new test covers the surrogate cut; go test ./... RC=0 at round 1, core/shipped green after regen; DoD 1-3 verified"
hunt-verdict: "pass: bughunt done at 1d1990d5b, 0 change, 0 pre-existing (round 1 change fixed: inLine cuts by code points)"
round: 2
back-from: testing
back-by: lane
back-count: 1
review-verdict: "accept: DoD 1-3 are met by the code and the tests (board.test.ts 'the line a press hands this session' plus w.submitted==[] on every failure path); the four checks found no defect: concurrency is the 'already in' case, error paths return before the submit, input is cut by code points, permissions: no band in an agent session"
review-round: 2
review-verdict-gaps: |-
  note: the line starts a model turn in any person session that runs the cockpit, not only a teamlead's - that is the goal, but in a plain session nobody leads, the turn costs tokens and does nothing useful; a later option is to send only when the session is a teamlead's
  note: void $.prompt.submit(...) drops a rejection without a word (register.tsx:618, :1201, same as the older decide() at :1041) - if submit fails, the session never hears and the cockpit does not say so either; fix: .catch(() => $.ui.toast(...))
  note: pre-existing: cuts by UTF-16 index are still in digest.ts:184, threads.ts:43, board.ts:782 - they could reuse inLine
  note: checked only in the engine harness (claude plugin test, 296 pass), not with a real press in a live session
review-check: |-
  1. cd /home/alex/projects/.worktrees/jaira-409xcv/core/mod/builtin/jaira-cockpit && claude plugin test - you should see 296 pass, 0 fail
  2. Start a session with this branch's cockpit: claude --plugin-dir /home/alex/projects/.worktrees/jaira-409xcv/core/mod/builtin/jaira-cockpit, in a repository that has a ticket in signoff
  3. Press Accept on that ticket in the band - a prompt appears in the session: 'jaira-cockpit: the person accepted <ID> in the cockpit (signoff -> done); ...', the model starts answering, and nobody typed anything
  4. Press Give back on another ticket in signoff, type a reason with a line break and an emoji - one line appears: '... gave <ID> back ... (signoff -> in-progress), <who>: <reason on one line>'
  5. Press Accept a second time on a ticket that is already in done (or another session moved it first) - only the error pane 'Accept <ID> failed' appears, no new line in the session
---

# Accept und Give back im Cockpit melden der Session eine Zeile, damit der Teamlead sofort weitermacht

## Definition of Done

- [x] Accept und Give back im Cockpit legen nach erfolgreichem Move genau eine Zeile in die Session (Test)
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: describe 'the line a press hands this session' (Accept with PR, Accept whose next step fails, Give back to a new worker); explanation-pane [accept]/[give back] tests
- [x] Ein Move, der fehlschlägt oder 'already in' antwortet, meldet nichts (Test)
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: 'a move another session made first: the second press writes nothing' (back/accept/band-accept), 'a refused Accept starts no next step at all', 'a refused move writes no note', 'an Accept the gates refuse writes no decisions' - each expects w.submitted to be []
- [x] Release-Note core/release/notes/<ID>.md, shipped.sha256
  proof: core/release/notes/409XCV.md; core/shipped/shipped.sha256 regenerated by scripts/shipped-history.sh --write; go test ./core/... green

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-21 21:13 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2274544.sock
- **2026-11-21 21:16 · Alexander Sacharov** — decisions (settled, no person needed): (1) channel is $.prompt.submit({text}) - the cockpit's decide() already submits a line to this session that way, and a submitted prompt starts the teamlead's turn, which is the goal. (2) Accept: the line goes after afterAccept finished, so a teamlead reading it never starts a PR role beside the one Accept starts. (3) Give back: the line goes in handedTo, once it is known who got the work (old/new/no worker) - otherwise the teamlead may spawn a second worker. Reason: first line, cut to ~120 chars. (4) No person name in the line ('the person accepted'): the cockpit only knows git identity, and the session knows who presses. (5) branch feat/409xcv sits on feat/pdz444's merge (not yet in upstream/master).
- **2026-11-21 21:21 · Alexander Sacharov** — implemented: accept() submits t.toSessionAccepted after afterAccept (which now returns what it said: PR/landing started, why none, or 'next step failed'); giveBack's handedTo submits t.toSessionGivenBack. oneLine() in texts.ts collapses whitespace and cuts the reason at 120 chars. claude plugin test: 295 pass.
- **2026-11-21 21:21 · Alexander Sacharov** — critique: none. Channel reuses decide()'s $.prompt.submit (register.tsx decide); afterAccept's return replaces reading noticeAtom back, which would couple to say()/fail() ordering; the who-wording in toSessionGivenBack differs from returnedTo on purpose (the session needs 'new worker started' to not spawn one).
- **2026-11-21 21:24 · Alexander Sacharov** — testing: function pass ran in the engine's own harness (claude plugin test, 2.1.x): Accept with PR started, Accept whose next step fails, Give back to a new worker with a multi-line 200+ char reason; failure paths: refused dry-run, refused move (jaira 0.3.6), 'already in' x3 presses, gates refusing - each submits nothing. Not run live in a session with a real press: that is the person's hand-check (claude --plugin-dir core/mod/builtin/jaira-cockpit). The change writes no file: prompt.submit only; two presses at once are the 'already in' case.
- **2026-11-21 21:26 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/texts.ts:10 - class: a text cut by UTF-16 index (slice(0, max-1)) splits a character outside the BMP (emoji, some CJK) and leaves a lone high surrogate in the line prompt.submit sends to the teamlead session; reached by toSessionGivenBack (the reason the person typed, max 120) and toSessionAccepted (the next-step text, max 200); the same cut already sits in older code at digest.ts:184 (report line rewritten into the conversation), threads.ts:43 (clip for the haiku prompt) and board.ts:782 (follow-up title) - pre-existing there, but this change adds two new paths into the session prompt; example: a Give back reason of 118 chars followed by a thumbs-up emoji gives a line ending in 'x\ud83d…', isWellFormed() false; ran: node --experimental-strip-types on t.toSessionGivenBack('MG677R','signoff','new','x'.repeat(118)+'\u{1F44D} passt nicht') (lone surrogate: true). Not a blocker as verified: whether the engine or the API rejects the malformed string (a known 400 'no low surrogate in string' would stall the session) was not reproduced. Fix: cut by code points (Array.from(flat).slice(...)) or drop a trailing high surrogate
- **2026-11-21 21:26 · Alexander Sacharov** — bughunt: done at 8694f3345 - 1 finding (1 change, 0 pre-existing) in 4 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/409XCV.md core/shipped/shipped.sha256
- **2026-11-21 21:28 · Alexander Sacharov** — fix round 2: inLine (texts.ts) cuts by code points (Array.from) - the class is any cut by UTF-16 index on a text handed to the session; both new paths (reason, next-step text) go through inLine. Test 'a reason cut at its limit never splits a character outside the BMP'. Older private clips (digest.ts:184, threads.ts:43, board.ts:782) are pre-existing and not touched here.
- **2026-11-21 21:29 · Alexander Sacharov** — critique round 2: none - the fix sits in inLine, the one place both new session-line paths cut text.
- **2026-11-21 21:30 · Alexander Sacharov** — bughunt: done at 1d1990d5b - 0 findings in 3 modules, follow-up hunt (diff since 8694f3345: texts.ts inLine, board.test.ts, shipped.sha256); concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/409XCV.md core/shipped/shipped.sha256
- **2026-11-21 21:32 · Alexander Sacharov** — review: run-lane.sh failed to start the Herdr tab ('claude did not come up in t19m'); review handed to a fresh subagent of the worker instead (dies with the worker).
- **2026-11-21 21:35 · Alexander Sacharov** — review: accept at 1d1990d5b - no blocker; 4 notes (a turn in any person session, submit rejection swallowed, pre-existing UTF-16 cuts, no live press). Tests: claude plugin test 296 pass; go test core/shipped core/release ok.
