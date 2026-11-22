---
id: 01M7ZYRN6S4CGTGB1YYPTT3Z11
title: Worker und Bughunt melden dem Elternteil eine Statuszeile; der volle Bericht steht nur auf dem Ticket
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Der Teamlead-Kontext wächst pro Bericht um eine Zeile, nicht um eine Seite."
context: |-
  Was heute teuer war (21.11.): Bughunt- und Worker-Berichte kamen als ganze Seiten in die Teamlead-Session (je 1-3k Token, viele pro Ticket), obwohl alles schon als Notes auf dem Ticket steht.
  55CS2M hat die feste Statuszeile eingeführt (jaira: <ID> <von> -> <nach> · branch · commits · gates · notes · verdict); das Cockpit liest sie (digest.ts).
  Gewünscht: jaira-worker, jaira-dispatcher und jaira-role-bughunt schicken per SendMessage nur die Statuszeile plus eine Zeile "braucht dich: ..." falls etwas offen ist; der ausführliche Teil geht als jaira note aufs Ticket. Der Teamlead-Skill liest Details mit jaira show bei Bedarf.
  Dateien: core/role/builtin/jaira-worker, jaira-dispatcher, jaira-role-bughunt, jaira-teamlead (SKILL.md).
  Start erst nach P122GE (gleicher Worker-Skill).
definition-of-done: "Worker, Dispatcher und Bughunt senden dem Elternteil nur Statuszeile + optional eine Zeile 'braucht dich' (SKILL-Text, Test in core/role)"
tags:
  - gates
blocked-by: []
related: []
commits:
  - bbb5690d6e4ba510776f56528db425df1d519d57
  - eb5df991fb917eb14fd8718540f776215f04334d
  - f7242cbf0611e81634f3660121726df9644baaa4
  - 252f08ed3a68dae9b9b0771fc93687397773abc0
created-at: 2026-11-21T21:04:00Z
updated-at: 2026-11-22T14:41:01Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1588595
claimed-at: 2026-11-22T11:38:15Z
outcome-what: "Worker, dispatcher and bug hunt hand their parent only the status line plus an optional 'needs you:' line; the full report goes onto the ticket as a 'report:' / 'hunt report:' note. Teamlead skill reads details with jaira show. Cockpit digest keeps the needs-you line, framed or not, and ignores an empty or 'none' one."
outcome-why: "Reports of 1-3k tokens each filled the teamlead's session although everything was already on the ticket."
outcome-resolves: "Teamlead context grows by one line per report, not a page."
review-summary: "The worker, dispatcher and bughunt SKILL.md files now send the parent only the status line, plus one 'needs you: <what>' line above it when something waits for a person. The full three-line report is written onto the ticket as a 'report:' note ('hunt report:' for a bug hunt, which is not a 'bughunt:' note, so the testing lane's grep for findings skips it). A bug hunt on a ticket hands back its closing line without the files list. The teamlead skill says the details are read with jaira show and not out of the message. The cockpit digest (digest.ts needsYouOf) keeps the needs-you line in the shortened line: frame tags are stripped like the status line's, the last such line wins, an empty value is dropped and nothing else is filtered. Covered by TestReportsToTheParentAreOneLine, 4 cockpit cases, a release note and shipped.sha256."
review-gaps: |-
  note: digest.ts needsYouOf clips the value to 140 chars, so a long ask with an acceptance link at its end loses the link in the shortened line (it stays in /cockpit report). Fix: clip only the text before a URL, or tell the worker to put the link first
  note: the bug hunt's short line has no status-line shape, so the cockpit shortens it to handle + lane and drops "2 findings (1 change...)". The critique already accepted this: the testing lane reads the findings off the ticket
  note: pre-existing: digest.ts questionOf uses a prefix test (NO_QUESTION), so "question: no answer from CI..." is dropped as no question. Already like this on master
  note: the worktree carries uncommitted .jaira/lanes/*, AGENTS.md, CLAUDE.md changes from a board regeneration that are not part of this ticket. Keep them out of the logbook commit (git add only the ticket file)
test-verdict: "pass: go test ./... RC=0, cockpit plugin test 306 pass; DoD 1-3 verified in the tree (proof lines re-checked); digest exercised by test cases: unframed, framed, empty, ask starting with 'nothing'; the change writes nothing at runtime (skill text + a pure parser), so no two-writer case"
hunt-verdict: "pass: round 3 at 252f08ed3, 0 findings; rounds 1-2 findings fixed in f7242cbf0 and 252f08ed3"
round: 3
back-from: testing
back-by: lane
back-count: 2
review-verdict: "accept: all 3 DoD items are met in the diff; go test ./core/role ./core/shipped ./core/release ok; claude plugin test (cockpit) 306 pass; no blocker in digest.ts or the four SKILL.md files"
review-round: 3
review-verdict-gaps: |-
  note: digest.ts needsYouOf clips the value to 140 chars, so a long ask with an acceptance link at its end loses the link in the shortened line (it stays in /cockpit report). Fix: clip only the text before a URL, or tell the worker to put the link first
  note: the bug hunt's short line has no status-line shape, so the cockpit shortens it to handle + lane and drops "2 findings (1 change...)". The critique already accepted this: the testing lane reads the findings off the ticket
  note: pre-existing: digest.ts questionOf uses a prefix test (NO_QUESTION), so "question: no answer from CI..." is dropped as no question. Already like this on master
  note: the worktree carries uncommitted .jaira/lanes/*, AGENTS.md, CLAUDE.md changes from a board regeneration that are not part of this ticket. Keep them out of the logbook commit (git add only the ticket file)
review-check: "1. cd core/mod/builtin/jaira-cockpit && claude plugin test .  -> 306 pass, 0 fail; the cases named 'needs-you' pass  2. go test ./core/role -run TestReportsToTheParentAreOneLine  -> ok  3. open core/role/builtin/jaira-worker/SKILL.md, section '## Report'  -> it says: the full report is a 'report:' note on the ticket; the parent gets only 'needs you: ...' (when something waits) plus the 'jaira: <ID> ...' status line  4. after jaira update --roles, start a worker on a small ticket from the teamlead  -> its message to you is at most 2 lines, and jaira show <id> has a 'report:' note"
question: ""
---

# Worker und Bughunt melden dem Elternteil eine Statuszeile; der volle Bericht steht nur auf dem Ticket

## Definition of Done

- [x] Worker, Dispatcher und Bughunt senden dem Elternteil nur Statuszeile + optional eine Zeile 'braucht dich' (SKILL-Text, Test in core/role)
  proof: core/role/builtin/jaira-worker/SKILL.md:435, jaira-dispatcher/SKILL.md:565, jaira-role-bughunt/SKILL.md:182; core/role/role_test.go:1489 TestReportsToTheParentAreOneLine; cockpit digest keeps needs-you (hooks/digest.ts needsYouOf, test 'a needs-you line above the status line is kept')
- [x] Teamlead-Skill sagt: Details aus jaira show, nicht aus dem Bericht
  proof: core/role/builtin/jaira-teamlead/SKILL.md:247; checked by TestReportsToTheParentAreOneLine
- [x] Release-Note core/release/notes/<ID>.md, shipped.sha256
  proof: core/release/notes/TT3Z11.md; core/shipped/shipped.sha256 (+6 lines), TestEveryInstalledFileIsAShippedVersion green

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 11:38 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 11:39 · Alexander Sacharov** — decision (worker, settled without asking): 1) the three report lines plus the context/cache figure go onto the ticket as one note 'report: ...'; the parent gets only the status line and, when something is open, one line 'needs you: <what>' above it. 2) per-lane lines go onto the ticket as notes, not upward — otherwise the teamlead still grows by a line per lane. 3) the cockpit's digest.ts drops every line but the status line and 'question:' — so it learns a 'needs you:' line too, or that line would never reach the teamlead. 4) questions with numbered options stay as they are: a question is not a report. 5) bughunt: on a ticket its closing 'bughunt: done at' note is the report line, the three lines become a note; with paths only (no ticket) the three lines stay the report.
- **2026-11-22 11:40 · Alexander Sacharov** — decision revised: point 2 dropped — per-lane lines stay upward; each is already one line, which is what the goal asks, and they keep the person able to react mid-run. Only the final report shrinks.
- **2026-11-22 11:43 · Alexander Sacharov** — lane in-progress: done; worktree also carries uncommitted .jaira/lanes/*, AGENTS.md, CLAUDE.md changes that this ticket did not make (a lanes/agent-block regeneration, appeared before any edit) — left out of every commit
- **2026-11-22 11:43 · Alexander Sacharov** — critique (pass 1, whole diff): nothing to change. Let stand: (a) a separate 'needs you:' key in digest.ts instead of reusing 'question:' — reusing would need no TS change but would label an acceptance link a question; (b) the bug hunt's short line has no status-line shape, so the cockpit digest shortens it to the handle and drops '2 findings (1 change...)' — that was already so for its three-line report, the testing lane reads the findings off the ticket anyway, and a hunt status line is its own ticket if wanted; (c) 'report:' note with ' | ' separators is new but one line, readable in jaira show.
- **2026-11-22 11:45 · Alexander Sacharov** — waiting: bug hunt (background subagent) since 12:45
- **2026-11-22 11:46 · Alexander Sacharov** — bughunt: change: blocker: core/mod/builtin/jaira-cockpit/hooks/digest.ts:103 - class: needs-you line that shares its physical line with an engine frame tag (<agent-message from=…>, <result>) is dropped; NEEDS_YOU is anchored at ^ and, unlike statusLineOf (FRAME_OPEN, digest.ts:136), strips no frame, and register.tsx:240 hands digest() the raw framed text; since the report is now two lines, the needs-you line is exactly the first line, the one the frame opens on; example: 'Another Claude session sent a message:\n<agent-message from="w">needs you: accept at https://claude.ai/artifact/x1\njaira: 55CS2M testing -> human · notes: 2 · verdict: sent back</agent-message>' (same shape as cockpit.test.ts:118) digests to '55CS2M · testing → human · sent back · notes 2' with no needs-you; same for '<result>needs you: …'; the new test cockpit.test.ts:106 only covers an unframed line; ran: node --experimental-strip-types on a scratch script calling digest()
- **2026-11-22 11:46 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/digest.ts:103 - class: needs-you value is not validated: the \s* after 'needs you:' crosses a newline, and no empty/'nothing' value is filtered (questionOf has NO_QUESTION for this); example: 'needs you:' with nothing after it captures the next line, so the digest reads '… · needs you: jaira: 55CS2M in-progress -> signoff · branch: … · verdict: approve' (status line shown twice, as an ask); 'needs you: nothing' / 'needs you: -' is shown as an ask too, though jaira-worker forbids that line; ran: node --experimental-strip-types on a scratch script calling digest()
- **2026-11-22 11:46 · Alexander Sacharov** — hunt report: 2 findings, both change (1 blocker: digest drops a framed needs-you line), notes on TT3Z11 | read whole: hooks/digest.ts, the TT3Z11 diff of texts.ts, cockpit.test.ts, role_test.go, four role SKILL.md files, release note, shipped.sha256 (hashes match); frame.ts/register.tsx/gitrepo/hunt.go read as callers; uncommitted .jaira/lanes/*, AGENTS.md, CLAUDE.md are board-regeneration prose | not run: the cockpit test suite (needs the claude-code test engine); go test ./core/role ./core/shipped pass
- **2026-11-22 11:46 · Alexander Sacharov** — bughunt: done at eb5df991f - 2 findings (2 change, 0 pre-existing) in 6 modules, first hunt; concurrency, error paths, input, permissions checked; files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/optimize.md .jaira/lanes/pre-process.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/mod/builtin/jaira-cockpit/hooks/digest.ts core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/cockpit.test.ts core/release/notes/TT3Z11.md core/role/builtin/jaira-dispatcher/SKILL.md core/role/builtin/jaira-role-bughunt/SKILL.md core/role/builtin/jaira-teamlead/SKILL.md core/role/builtin/jaira-worker/SKILL.md core/role/role_test.go core/shipped/shipped.sha256
- **2026-11-22 11:50 · Alexander Sacharov** — fix round 2: needsYouOf now reads line by line like statusLineOf (frame tags stripped, last line wins), [ \t]* instead of \s* so an empty value never takes the next line, and NO_QUESTION filters empty/none — NO_QUESTION gained 'nothing', which also makes 'question: nothing' no question. Tests: cockpit.test.ts 3 new cases (framed, none/nothing, empty above status line).
- **2026-11-22 11:51 · Alexander Sacharov** — critique pass 2 (fix only): both hunt findings answered in f7242cbf0; the line loop in needsYouOf mirrors statusLineOf on purpose (same frame rule) — let stand. optimize round 2: nothing to remove.
- **2026-11-22 11:51 · Alexander Sacharov** — waiting: bug hunt round 2 (background subagent) since 12:51
- **2026-11-22 11:52 · Alexander Sacharov** — bughunt: change: blocker: core/mod/builtin/jaira-cockpit/hooks/digest.ts:142 (needsYouOf) - class: NO_QUESTION (digest.ts:90) is a prefix test, not a whole-value test, so any real ask whose first word is none/nothing/no/keine/nein/nichts is dropped as 'no ask' and the cockpit digest loses the only line the teamlead acts on, link included; same prefix test sits in questionOf (digest.ts:96, older than the ticket) for questions; example: 'needs you: nothing to decide, accept at https://claude.ai/artifact/x1', 'needs you: no answer from CI on #149, rerun it', 'needs you: keine Antwort vom CI, bitte neu starten' all return undefined; ran: npx tsx scratch script calling needsYouOf (all three -> undefined); fix direction: treat as no ask only when the whole value is one of those words, optionally with trailing punctuation
- **2026-11-22 11:52 · Alexander Sacharov** — bughunt: done at f7242cbf0 - 1 finding (1 change, 0 pre-existing) in 6 modules, follow-up hunt since eb5df991f (diff only, no new files); concurrency, error paths, input, permissions checked; shipped.sha256 swap checked (removed hash 658716 is branch-only, in no tag; new hash matches digest.ts); files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/optimize.md .jaira/lanes/pre-process.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/mod/builtin/jaira-cockpit/hooks/digest.ts core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/cockpit.test.ts core/release/notes/TT3Z11.md core/role/builtin/jaira-dispatcher/SKILL.md core/role/builtin/jaira-role-bughunt/SKILL.md core/role/builtin/jaira-teamlead/SKILL.md core/role/builtin/jaira-worker/SKILL.md core/role/role_test.go core/shipped/shipped.sha256
- **2026-11-22 11:53 · Alexander Sacharov** — fix round 3: smaller design instead of a better word list — needsYouOf filters only an empty value (the worker writes no needs-you line when nothing waits); the 'nothing' added to NO_QUESTION in round 2 is reverted, so questionOf is as on master. Its pre-existing prefix test (question: no answer from CI… dropped) is left as it was — not this ticket.
- **2026-11-22 11:53 · Alexander Sacharov** — critique pass 3 (fix only): the round-2 finding is answered by removing the word filter, not refining it — nothing to add. optimize round 3: none, the change only removed code.
- **2026-11-22 11:54 · Alexander Sacharov** — waiting: bug hunt round 3 (background subagent) since 12:54
- **2026-11-22 11:56 · Alexander Sacharov** — bughunt: done at 252f08ed3 - 0 findings in 6 modules, follow-up hunt since f7242cbf0 (diff only, no new files); concurrency, error paths, input, permissions checked; NO_QUESTION is back to master's word list, so questionOf is unchanged against master; shipped.sha256 swap checked (new hash 181543 matches digest.ts, removed f48ad1 is in no tag); ran: claude plugin test . (306 pass), go test ./core/role ./core/shipped (ok); files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/optimize.md .jaira/lanes/pre-process.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/mod/builtin/jaira-cockpit/hooks/digest.ts core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/cockpit.test.ts core/release/notes/TT3Z11.md core/role/builtin/jaira-dispatcher/SKILL.md core/role/builtin/jaira-role-bughunt/SKILL.md core/role/builtin/jaira-teamlead/SKILL.md core/role/builtin/jaira-worker/SKILL.md core/role/role_test.go core/shipped/shipped.sha256
- **2026-11-22 11:56 · Alexander Sacharov** — lane testing: pass after 3 rounds (hunt found 3 change defects in digest.ts needsYouOf, all fixed); waiting: review agent (Herdr tab via run-lane.sh) since 12:56
- **2026-11-22 11:59 · Alexander Sacharov** — review: verdict accept, but the move to signoff is refused: the hunt gate sees uncommitted changes in .jaira/lanes/*, AGENTS.md and CLAUDE.md, files the last hunt (252f08ed3) read. They come from a board regeneration, not from this ticket. Question: what happens to them?
1. (recommended) The person moves it themselves: jaira move TT3Z11 --to signoff --force. The changes are prose, the hunt read them, and they stay out of this branch
2. Commit them as a separate chore commit on this branch and run the bughunt once more. They then ship with this PR
3. Revert them (git restore .jaira/lanes AGENTS.md CLAUDE.md), then move. The regeneration is lost and comes back with the next jaira update
- **2026-11-22 12:00 · Alexander Sacharov** — report: reports to the parent shrink to status line + optional 'needs you:'; full report as 'report:' note; cockpit keeps needs-you (3 hunt findings fixed) | review, verdict accept; context unknown | open: question to parent — the uncommitted board-regeneration files block the move to signoff (force / chore commit / revert); 4 review notes for the person
- **2026-11-22 12:00 · Alexander Sacharov** — waiting: parent's answer on the uncommitted regeneration files (uds:/run/user/1000/cc-socks/25712.sock) since 13:00
- **2026-11-22 14:39 · Alexander Sacharov** — answer (Teamlead 22.11.): Board-Regeneration verworfen (git checkout), sie kam über GPSAR2/#179 nach master; Frage erledigt.
