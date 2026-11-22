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
updated-at: 2026-11-22T11:43:17Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1588595
claimed-at: 2026-11-22T11:38:15Z
outcome-what: "Worker, dispatcher and bug hunt hand their parent only the status line plus an optional 'needs you:' line; the full report goes onto the ticket as a 'report:' / 'hunt report:' note. Teamlead skill reads details with jaira show. Cockpit digest keeps the needs-you line."
outcome-why: "Reports of 1-3k tokens each filled the teamlead's session although everything was already on the ticket."
outcome-resolves: "Teamlead context grows by one line per report, not a page."
---

# Worker und Bughunt melden dem Elternteil eine Statuszeile; der volle Bericht steht nur auf dem Ticket

## Definition of Done

- [x] Worker, Dispatcher und Bughunt senden dem Elternteil nur Statuszeile + optional eine Zeile 'braucht dich' (SKILL-Text, Test in core/role)
  proof: core/role/builtin/jaira-worker/SKILL.md:435, jaira-dispatcher/SKILL.md:564, jaira-role-bughunt/SKILL.md:182; core/role/role_test.go:1489 TestReportsToTheParentAreOneLine; cockpit digest keeps needs-you (hooks/digest.ts needsYouOf, test 'a needs-you line above the status line is kept')
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
