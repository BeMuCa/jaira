---
id: 01M7ZYRN6S4CGTGB1YYPTT3Z11
title: Worker und Bughunt melden dem Elternteil eine Statuszeile; der volle Bericht steht nur auf dem Ticket
status: critique
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
commits: []
created-at: 2026-11-21T21:04:00Z
updated-at: 2026-11-22T11:50:55Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1588595
claimed-at: 2026-11-22T11:38:15Z
outcome-what: "Worker, dispatcher and bug hunt hand their parent only the status line plus an optional 'needs you:' line; the full report goes onto the ticket as a 'report:' / 'hunt report:' note. Teamlead skill reads details with jaira show. Cockpit digest keeps the needs-you line, framed or not, and ignores an empty or 'none' one."
outcome-why: "Reports of 1-3k tokens each filled the teamlead's session although everything was already on the ticket."
outcome-resolves: "Teamlead context grows by one line per report, not a page."
review-summary: none
review-gaps: "rewrapped three over-long lines the change left in jaira-dispatcher, jaira-worker and jaira-teamlead SKILL.md (one 'them' → 'it'); replaced the intermediate shipped.sha256 hashes instead of adding a second set. Considered and left: needsYouOf beside questionOf in digest.ts — same shape, different key and no-value rule (questionOf drops 'none'/'-'); merging them would need a parameter for one caller each. No dead code, no cost on a hot path."
test-verdict: "fail: suite green (go test ./... RC=0, cockpit 303 pass); DoD 1-3 verified in the tree; but the hunt found two change defects in digest.ts needsYouOf (framed line dropped; empty value crosses a newline) — back to in-progress"
hunt-verdict: "fail: 2 change findings in hooks/digest.ts:103 (framed needs-you line dropped; value not validated)"
round: 2
back-from: testing
back-by: lane
back-count: 1
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
