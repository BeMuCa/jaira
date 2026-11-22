---
id: 01M804FV0R16X73XFM1FKRBGV9
title: Ein Worker macht nach dem Bughunt von selbst weiter statt idle in testing zu stehen
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ist der Bughunt fertig, geht der Worker ohne Anstoß weiter: hunt-verdict, frische Review, signoff."
context: |-
  Problem: der Worker bleibt nach dem Bughunt stehen und macht nicht weiter.
  - Beobachtet 21.11.2026 an Y2CD9X: der Bughunt notierte 22:24 UTC "bughunt: done at b5ec40701 - 0 findings". Danach stand der Worker (Session jaira-y2cd9x-dd) idle, das Ticket blieb in testing, ohne hunt-verdict und ohne Review. Erst um 22:39 schickte der Teamlead eine SendMessage.
  - Ablauf laut core/role/builtin/jaira-worker/SKILL.md, Abschnitt "The bug hunt": die testing-Lane startet den Hunter als Background-Agent, der Worker schreibt test-verdict und wartet ("waiting for bughunt", Notiz 22:19).
  - Vermutung, noch nicht geprüft: der Worker hat seinen Turn mit "warte" beendet. Wenn der Background-Agent fertig wird, wird die Session nicht wieder aufgeweckt, oder die Meldung geht verloren.
  - Folge: jede Nacht kann ein Ticket in testing hängen bleiben, ohne dass es jemand merkt. Im Cockpit sieht das aus wie Arbeit, die läuft.
  - Verwandt: TT3Z11 (Statuszeile an den Elternteil), RH10T7 (Worker wartet auf eine Person).
  - Zuerst: im Transcript von jaira-y2cd9x-dd nachsehen, ob die Task-Notification des Bughunts angekommen ist.
definition-of-done: "Ursache belegt: Transcript oder Repro zeigt, warum der Worker nach dem Bughunt nicht wieder lief"
tags: []
blocked-by: []
related: []
commits: []
created-at: 2026-11-21T22:44:03Z
updated-at: 2026-11-22T11:12:28Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-29770
claimed-at: 2026-11-22T11:07:53Z
outcome-what: "jaira-cockpit recognises a start line Herdr App wraps as <hc from=...>...</hc> (START in hooks/sessions.ts), so a worker's own subagent reports are no longer shortened to a /cockpit link; the jaira-worker skill gains 'When you wait': a waiting: note and a parent line before ending a turn, and a shortened cockpit line counts as the hand-back. Release note, shipped manifest."
outcome-why: "Y2CD9X sat idle in testing for 15 min on 21.11.: the bughunt's report reached the worker as '[shortened by jaira-cockpit] ...' because the cockpit did not see the session as carrying a ticket, and nothing showed that the worker was waiting. Checked: concurrency n/a (pure regex); error path - no match behaves as before; unusual input - attributes, newline after the tag, an hc line that only mentions /jaira-worker (tested); permissions - startOf only runs for non-agent origins, as before."
outcome-resolves: A worker under Herdr App gets its bughunt report whole and goes on to hunt-verdict and review without a SendMessage; a waiting worker says so on the ticket and to its parent.
---

# Ein Worker macht nach dem Bughunt von selbst weiter statt idle in testing zu stehen

## Definition of Done

- [x] Ursache belegt: Transcript oder Repro zeigt, warum der Worker nach dem Bughunt nicht wieder lief
  proof: transcript a68b2392 (jaira-y2cd9x) 21.11. 21:13:25Z: start line arrives as '<hc from="t3">/jaira-worker Y2CD9X ...</hc>' (origin human); START in core/mod/builtin/jaira-cockpit/hooks/sessions.ts anchors on ^\s*/jaira- -> carries() false -> bughunt report shortened to a /cockpit link
- [x] Ein Worker, dessen Bughunt fertig wird, schreibt hunt-verdict und geht ohne SendMessage weiter (Test oder nachgestellter Lauf als Beleg)
  proof: core/mod/builtin/jaira-cockpit/hooks/sessions.ts:39 (START accepts the <hc> wrapper); tests/board.test.ts 'a session that carries a ticket > recognises its start line relayed by Herdr App, and a subagent's report reaches it whole' (+4 siblings; 3 fail without the fix)
- [x] Ein Worker, der auf etwas wartet, sagt das sichtbar (Notiz oder Statuszeile), damit der Teamlead ein Hängen erkennt
  proof: core/role/builtin/jaira-worker/SKILL.md:46 '## When you wait' - waiting: note + line to parent before ending a turn; a shortened cockpit line is a hand-back

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 11:05 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2784.sock
- **2026-11-22 11:06 · Alexander Sacharov** — cause (transcript ~/.claude/projects/-home-alex-projects--worktrees-jaira-y2cd9x/a68b2392…jsonl, 21.11. 22:24:51-52 UTC): the bughunt subagent's final report reached the worker as a peer message (origin.kind=peer, from=<subagent id>) and the jaira-cockpit hook had shortened it to '[shortened by jaira-cockpit] Y2CD9X — full report: /cockpit report Y2CD9X'. The worker read that as 'only a cockpit link, no new instructions' and ended its turn. The task-notification for the same agent followed 0.1 s later with queueTranscriptOnly=true (report already delivered) - it started no turn. Nothing else woke the session until the teamlead's message at 22:39. So: not a lost notification, a hidden hand-back.
- **2026-11-22 11:07 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock (restart; previous worker pid 10355 and parent 2784 gone)
- **2026-11-22 11:09 · Alexander Sacharov** — root cause one level deeper: the cockpit is meant to leave a ticket-carrying session's deliveries alone (register.tsx: 'if (await carries($)) return next(e)'), but under Herdr App every typed line is wrapped as <hc from="tN">...</hc>, and startOf's START regex (hooks/sessions.ts) requires the line to start with /jaira-. So since hc (2026-11-11) no worker/lane/dispatcher session is recognised as carrying: reports from its own subagents get shortened, and no session record is written either. Fix: accept an optional <hc ...> wrapper in START.
- **2026-11-22 11:09 · Alexander Sacharov** — decision: no open decisions. DoD 3 settled as both: a 'waiting: <what> since <time>' note on the ticket and one line to the parent, written in the jaira-worker skill before the worker ends a turn to wait. Reason: the note survives a dead session, the line is what the teamlead actually sees.
