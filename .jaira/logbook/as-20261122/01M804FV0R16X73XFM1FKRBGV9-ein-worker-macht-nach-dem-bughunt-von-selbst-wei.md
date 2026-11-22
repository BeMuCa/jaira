---
id: 01M804FV0R16X73XFM1FKRBGV9
title: Ein Worker macht nach dem Bughunt von selbst weiter statt idle in testing zu stehen
status: done
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
commits:
  - fe8a6e1fc5cc00e79e91d1e2396e84fc7bdf6340
created-at: 2026-11-21T22:44:03Z
updated-at: 2026-11-22T11:28:14Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-29770
claimed-at: 2026-11-22T11:07:53Z
outcome-what: "jaira-cockpit recognises a start line Herdr App wraps as <hc from=...>...</hc> (START in hooks/sessions.ts), so a worker's own subagent reports are no longer shortened to a /cockpit link; the jaira-worker skill gains 'When you wait': a waiting: note and a parent line before ending a turn, and a shortened cockpit line counts as the hand-back. Release note, shipped manifest."
outcome-why: "Y2CD9X sat idle in testing for 15 min on 21.11.: the bughunt's report reached the worker as '[shortened by jaira-cockpit] ...' because the cockpit did not see the session as carrying a ticket, and nothing showed that the worker was waiting. Checked: concurrency n/a (pure regex); error path - no match behaves as before; unusual input - attributes, newline after the tag, an hc line that only mentions /jaira-worker (tested); permissions - startOf only runs for non-agent origins, as before."
outcome-resolves: A worker under Herdr App gets its bughunt report whole and goes on to hunt-verdict and review without a SendMessage; a waiting worker says so on the ticket and to its parent.
review-summary: "The jaira-cockpit's START regex (core/mod/builtin/jaira-cockpit/hooks/sessions.ts:39) now accepts an optional Herdr App wrapper '<hc ...>' before /jaira-worker|dispatcher|role-lane. Since hc wraps every line typed from another tab, no worker was recognised as carrying a ticket, so its own subagents' reports (the bug hunt's) were shortened to a '/cockpit report' link and the worker stopped on it. With the wrapper recognised, carries() is true and reports pass whole. Six tests in tests/board.test.ts cover typed, relayed, multi-attribute/newline, lane and a non-start <hc> line. Second, jaira-worker SKILL.md gets '## When you wait': a 'waiting: <what> since <HH:MM>' note plus a line to the parent before ending a turn, and treats a shortened cockpit line as the hand-back to act on in the same turn (fallback for an older cockpit). Release note and shipped manifest updated."
review-gaps: |-
  note: the cockpit fix is proven by tests only - the live run on this ticket used the old installed cockpit and passed through the SKILL fallback; after jaira update --roles a worker started from a Herdr tab should get the bug hunt report whole (check once by hand)
  note: known limit: a relayed '<hc ...>/jaira-worker <id>' landing in the teamlead's own tab marks that session as carrying, exactly as the bare line did before - not new, only now also via hc
  note: [^>]* in the wrapper stops at a '>' inside a quoted attribute value - Herdr App writes none today
test-verdict: "pass: go build ./..., go vet (core/role, core/mod), go test -count=1 ./core/... exit 0; scripts/shipped-history.sh --check exit 0; claude plugin test . (jaira-cockpit) 302 pass 0 fail, the 3 hc cases fail without the fix. DoD 1 checked against transcript a68b2392 (start line '<hc from=\"t3\">/jaira-worker Y2CD9X ...', origin human). DoD 2: tests 'recognises its start line relayed by Herdr App...' etc. DoD 3: SKILL.md '## When you wait'. Function: START at node over <hc>, <hcx> (no match), nested <hc>, text before the slash (no match), short id (no match), role-lane lane kept. Failure path: an hc line that only mentions /jaira-worker leaves the session a person's (test)."
hunt-verdict: "pass: bughunt done at fe8a6e1fc - 0 findings in 4 modules"
review-verdict: "accept: DoD 1-3 met - root cause traced in the transcript, START accepts the <hc> wrapper (302/302 cockpit tests pass, 3 new ones fail without the fix), worker skill writes a waiting note and line; no defects found in sessions.ts, the prompt.submit/session.receive paths in register.tsx or the skill text"
review-round: 1
review-verdict-gaps: |-
  note: the cockpit fix is proven by tests only - the live run on this ticket used the old installed cockpit and passed through the SKILL fallback; after jaira update --roles a worker started from a Herdr tab should get the bug hunt report whole (check once by hand)
  note: known limit: a relayed '<hc ...>/jaira-worker <id>' landing in the teamlead's own tab marks that session as carrying, exactly as the bare line did before - not new, only now also via hc
  note: [^>]* in the wrapper stops at a '>' inside a quoted attribute value - Herdr App writes none today
review-check: "1. cd core/mod/builtin/jaira-cockpit && claude plugin test .  -> '302 pass, 0 fail', including 'recognises its start line relayed by Herdr App'  2. jaira update --roles (brings the new cockpit and jaira-worker skill)  3. from the teamlead tab start a worker on any ticket with spawn.sh, wait until its testing lane starts the bug hunt  4. in the worker tab the hunt report arrives as full text, not '[shortened by jaira-cockpit] ... /cockpit report'  5. jaira show <id> shows a 'waiting: bughunt ... since HH:MM' note, and the worker moves on without you sending anything"
question: ""
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
- **2026-11-22 11:12 · Alexander Sacharov** — critique: none. Considered stripping <hc> in register.tsx before every check instead - startOf is the only reader of the typed line (isPersonsPrompt keys on origin, and 'human' is not in its PERSON set), so the regex is the one place. The skill's 'shortened line is the hand-back' paragraph stays on purpose: it covers a machine whose cockpit is older than its skills, or a start line in a shape not yet known.
- **2026-11-22 11:13 · Alexander Sacharov** — optimize: none - no other parser of the start line (grep /jaira-worker over core/mod, core/role, internal: only spawn.sh writes it, sessions.ts reads it); nothing orphaned; the 'typed' test case overlaps the existing START tests in setup but is the only one that checks a task-notification passes whole.
- **2026-11-22 11:14 · Alexander Sacharov** — waiting: bughunt (background agent, testing lane) since 12:14
- **2026-11-22 11:15 · Alexander Sacharov** — bughunt: done at fe8a6e1fc - 0 findings in 4 modules, first hunt; concurrency, error paths, input, permissions checked; files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/optimize.md .jaira/lanes/pre-process.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/mod/builtin/jaira-cockpit/hooks/sessions.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/KRBGV9.md core/role/builtin/jaira-worker/SKILL.md core/shipped/shipped.sha256
- **2026-11-22 11:15 · Alexander Sacharov** — live check of DoD 2/3 on this very ticket: this worker runs under hc with the old installed cockpit, so the bughunt's report arrived as '[shortened by jaira-cockpit] KRBGV9 — full report: /cockpit report KRBGV9'. Following the new 'When you wait' rule it read the bughunt note off the ticket and went on in the same turn, no SendMessage needed.
- **2026-11-22 11:16 · Alexander Sacharov** — waiting: fresh review (run-lane.sh, own Herdr tab) since 12:16
- **2026-11-22 11:18 · Alexander Sacharov** — review: accept, but the move to signoff is refused - the bughunt gate sees uncommitted jaira-update residue (.jaira/lanes/*.md, AGENTS.md, CLAUDE.md: jaira:version 0.6.0 marker, 3-loop limit text, lane prompt additions), none of it this ticket's code. How to pass the gate?
1. (recommended) move the residue off this branch (WIP commit on a side branch or restore), then move KRBGV9 to signoff - keeps the branch to the ticket's change
2. commit the residue here as chore(KRBGV9): jaira update - simplest, but drags unrelated board config into the PR
3. person passes --force on the move - the hunt at fe8a6e1fc does cover all code; the dirty files are docs/config
- **2026-11-22 11:19 · Alexander Sacharov** — answer to the review's question: option 1. The jaira-update residue (lanes/*.md, AGENTS.md, CLAUDE.md, present before this ticket's worker started) is kept on branch wip/krbgv9-jaira-update-residue (commit 6b67e8986, a git stash create object) and restored out of this worktree; nothing deleted. Bring back with: git checkout wip/krbgv9-jaira-update-residue -- .jaira/lanes AGENTS.md CLAUDE.md
- **2026-11-22 11:19 · Alexander Sacharov** — waiting: fresh review, second agent (run-lane.sh) since 12:19
